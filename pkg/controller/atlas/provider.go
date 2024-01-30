package atlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

const (
	govAtlasDomain = "mongodbgov.com"
	orgIDKey       = "orgId"
	publicAPIKey   = "publicApiKey"
	privateAPIKey  = "privateApiKey"
)

type Provider interface {
	Client(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkClient(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error)
	IsCloudGov() bool
	IsResourceSupported(resource akov2.AtlasCustomResource) bool
}

type ProductionProvider struct {
	k8sClient       client.Client
	domain          string
	globalSecretRef client.ObjectKey
}

type credentialsSecret struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

func NewProductionProvider(atlasDomain string, globalSecretRef client.ObjectKey, k8sClient client.Client) *ProductionProvider {
	return &ProductionProvider{
		k8sClient:       k8sClient,
		domain:          atlasDomain,
		globalSecretRef: globalSecretRef,
	}
}

func (p *ProductionProvider) IsCloudGov() bool {
	domainURL, err := url.Parse(p.domain)
	if err != nil {
		return false
	}

	return strings.HasSuffix(domainURL.Hostname(), govAtlasDomain)
}

func (p *ProductionProvider) IsResourceSupported(resource akov2.AtlasCustomResource) bool {
	if !p.IsCloudGov() {
		return true
	}

	switch atlasResource := resource.(type) {
	case *akov2.AtlasProject,
		*akov2.AtlasTeam,
		*akov2.AtlasBackupSchedule,
		*akov2.AtlasBackupPolicy,
		*akov2.AtlasDatabaseUser,
		*akov2.AtlasFederatedAuth:
		return true
	case *akov2.AtlasDataFederation:
		return false
	case *akov2.AtlasDeployment:
		if atlasResource.Spec.ServerlessSpec == nil {
			return true
		}
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
	httpClient, err := httputil.DecorateClient(&http.Client{Transport: http.DefaultTransport}, clientCfg...)
	if err != nil {
		return nil, "", err
	}

	c, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(p.domain), mongodbatlas.SetUserAgent(operatorUserAgent()))

	return c, secretData.OrgID, err
}

func (p *ProductionProvider) SdkClient(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
	secretData, err := getSecrets(ctx, p.k8sClient, secretRef, &p.globalSecretRef)
	if err != nil {
		return nil, "", err
	}

	// TODO review we need add a custom logger to http client
	//httpClientWithCustomLogger := http.DefaultClient
	//err = httputil.LoggingTransport(log)(http.DefaultClient)
	//if err != nil {
	//	return nil, "", err
	//}

	c, err := NewClient(p.domain, secretData.PublicKey, secretData.PrivateKey)
	if err != nil {
		return nil, "", err
	}

	return c, secretData.OrgID, nil
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
