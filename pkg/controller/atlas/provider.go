package atlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/httputil"
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
	IsCloudGov() bool
	IsResourceSupported(resource akov2.AtlasCustomResource) bool
}

type ProductionProvider struct {
	k8sClient       client.Client
	domain          string
	globalSecretRef client.ObjectKey
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
	if secretRef == nil {
		secretRef = &p.globalSecretRef
	}

	secret := &corev1.Secret{}
	if err := p.k8sClient.Get(ctx, *secretRef, secret); err != nil {
		return nil, "", fmt.Errorf("failed to read Atlas API credentials from the secret %s: %w", secretRef.String(), err)
	}

	secretData := make(map[string]string)
	for k, v := range secret.Data {
		secretData[k] = string(v)
	}

	if missingFields, valid := validateSecretData(secretData); !valid {
		return nil, "", fmt.Errorf("the following fields are missing in the secret %v: %v", secretRef, missingFields)
	}

	clientCfg := []httputil.ClientOpt{
		httputil.Digest(secretData[publicAPIKey], secretData[privateAPIKey]),
		httputil.LoggingTransport(log),
	}
	httpClient, err := httputil.DecorateClient(&http.Client{Transport: http.DefaultTransport}, clientCfg...)
	if err != nil {
		return nil, "", err
	}

	c, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(p.domain), mongodbatlas.SetUserAgent(operatorUserAgent()))

	return c, secretData[orgIDKey], err
}

func validateSecretData(secretData map[string]string) ([]string, bool) {
	var missingFields []string
	requiredKeys := []string{orgIDKey, publicAPIKey, privateAPIKey}

	for _, key := range requiredKeys {
		if _, ok := secretData[key]; !ok {
			missingFields = append(missingFields, key)
		}
	}

	if len(missingFields) > 0 {
		return missingFields, false
	}

	return nil, true
}

func operatorUserAgent() string {
	return fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", version.Version, runtime.GOOS, runtime.GOARCH)
}
