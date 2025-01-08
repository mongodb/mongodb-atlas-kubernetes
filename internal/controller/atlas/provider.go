package atlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/mongodb-forks/digest"
	adminv20231115008 "go.mongodb.org/atlas-sdk/v20231115008/admin"
	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

const (
	govAtlasDomain = "mongodbgov.com"
	orgIDKey       = "orgId"
	publicAPIKey   = "publicApiKey"
	privateAPIKey  = "privateApiKey"
)

type Provider interface {
	Client(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkClient(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*adminv20231115008.APIClient, string, error)
	SdkClientSet(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*ClientSet, string, error)
	IsCloudGov() bool
	IsResourceSupported(resource api.AtlasCustomResource) bool
}

type ClientSet struct {
	SdkClient20231115008 *adminv20231115008.APIClient
	SdkClient20241113001 *adminv20241113001.APIClient
}

type ProductionProvider struct {
	k8sClient       client.Client
	domain          string
	globalSecretRef client.ObjectKey
	dryRunRecorder  record.EventRecorder
}

type credentialsSecret struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

func NewProductionProvider(atlasDomain string, globalSecretRef client.ObjectKey, k8sClient client.Client, dryRunRecorder record.EventRecorder) *ProductionProvider {
	return &ProductionProvider{
		k8sClient:       k8sClient,
		domain:          atlasDomain,
		globalSecretRef: globalSecretRef,
		dryRunRecorder:  dryRunRecorder,
	}
}

func (p *ProductionProvider) IsCloudGov() bool {
	domainURL, err := url.Parse(p.domain)
	if err != nil {
		return false
	}

	return strings.HasSuffix(domainURL.Hostname(), govAtlasDomain)
}

func (p *ProductionProvider) IsResourceSupported(resource api.AtlasCustomResource) bool {
	if !p.IsCloudGov() {
		return true
	}

	switch atlasResource := resource.(type) {
	case *akov2.AtlasProject,
		*akov2.AtlasTeam,
		*akov2.AtlasBackupSchedule,
		*akov2.AtlasBackupPolicy,
		*akov2.AtlasDatabaseUser,
		*akov2.AtlasSearchIndexConfig,
		*akov2.AtlasBackupCompliancePolicy,
		*akov2.AtlasFederatedAuth,
		*akov2.AtlasPrivateEndpoint:
		return true
	case *akov2.AtlasDataFederation,
		*akov2.AtlasStreamInstance,
		*akov2.AtlasStreamConnection:
		return false
	case *akov2.AtlasDeployment:
		hasSearchNodes := atlasResource.Spec.DeploymentSpec != nil && len(atlasResource.Spec.DeploymentSpec.SearchNodes) > 0
		isServerless := atlasResource.Spec.ServerlessSpec != nil

		return !(isServerless || hasSearchNodes)
	}

	return false
}

func (p *ProductionProvider) Client(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
	secretData, err := getSecrets(ctx, p.k8sClient, secretRef, &p.globalSecretRef)
	if err != nil {
		return nil, "", err
	}

	clientCfg := []httputil.ClientOpt{
		httputil.Digest(secretData.PublicKey, secretData.PrivateKey),
		httputil.LoggingTransport(log),
	}

	transport := p.newDryRunTransport(http.DefaultTransport)
	httpClient, err := httputil.DecorateClient(&http.Client{Transport: transport}, clientCfg...)
	if err != nil {
		return nil, "", err
	}

	c, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(p.domain), mongodbatlas.SetUserAgent(operatorUserAgent()))

	return c, secretData.OrgID, err
}

func (p *ProductionProvider) SdkClient(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*adminv20231115008.APIClient, string, error) {
	clientSet, orgID, err := p.SdkClientSet(ctx, secretRef, log)
	if err != nil {
		return nil, "", err
	}

	// Special case: SdkClient only returns the v20231115008 client.
	return clientSet.SdkClient20231115008, orgID, nil
}

func (p *ProductionProvider) SdkClientSet(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*ClientSet, string, error) {
	secretData, err := getSecrets(ctx, p.k8sClient, secretRef, &p.globalSecretRef)
	if err != nil {
		return nil, "", err
	}

	var transport http.RoundTripper = digest.NewTransport(secretData.PublicKey, secretData.PrivateKey)
	transport = p.newDryRunTransport(transport)
	transport = httputil.NewLoggingTransport(log, false, transport)

	httpClient := &http.Client{Transport: transport}

	clientv20231115008, err := adminv20231115008.NewClient(
		adminv20231115008.UseBaseURL(p.domain),
		adminv20231115008.UseHTTPClient(httpClient),
		adminv20231115008.UseUserAgent(operatorUserAgent()))
	if err != nil {
		return nil, "", err
	}

	clientv20241113001, err := adminv20241113001.NewClient(
		adminv20241113001.UseBaseURL(p.domain),
		adminv20241113001.UseHTTPClient(httpClient),
		adminv20241113001.UseUserAgent(operatorUserAgent()))
	if err != nil {
		return nil, "", err
	}

	return &ClientSet{
		SdkClient20231115008: clientv20231115008,
		SdkClient20241113001: clientv20241113001,
	}, secretData.OrgID, nil
}

func (p *ProductionProvider) newDryRunTransport(delegate http.RoundTripper) http.RoundTripper {
	if p.dryRunRecorder == nil {
		return delegate
	}

	return dryrun.NewDryRunTransport(p.dryRunRecorder, delegate)
}

func getSecrets(ctx context.Context, k8sClient client.Client, secretRef, fallbackRef *client.ObjectKey) (*credentialsSecret, error) {
	if secretRef == nil {
		secretRef = fallbackRef
	}

	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, *secretRef, secret); err != nil {
		return nil, fmt.Errorf("failed to read Atlas API credentials from the secret %s: %w", secretRef.String(), err)
	}

	secretData := credentialsSecret{
		OrgID:      string(secret.Data[orgIDKey]),
		PublicKey:  string(secret.Data[publicAPIKey]),
		PrivateKey: string(secret.Data[privateAPIKey]),
	}

	if missingFields, valid := validateSecretData(&secretData); !valid {
		return nil, fmt.Errorf("the following fields are missing in the secret %v: %v", secretRef, missingFields)
	}

	return &secretData, nil
}

func validateSecretData(secretData *credentialsSecret) ([]string, bool) {
	missingFields := make([]string, 0, 3)

	if secretData.OrgID == "" {
		missingFields = append(missingFields, orgIDKey)
	}

	if secretData.PublicKey == "" {
		missingFields = append(missingFields, publicAPIKey)
	}

	if secretData.PrivateKey == "" {
		missingFields = append(missingFields, privateAPIKey)
	}

	if len(missingFields) > 0 {
		return missingFields, false
	}

	return nil, true
}

func operatorUserAgent() string {
	return fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", version.Version, runtime.GOOS, runtime.GOARCH)
}
