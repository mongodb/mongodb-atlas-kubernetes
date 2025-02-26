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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

const (
	govAtlasDomain = "mongodbgov.com"
)

type Provider interface {
	Client(ctx context.Context, creds *Credentials, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkClientSet(ctx context.Context, creds *Credentials, log *zap.SugaredLogger) (*ClientSet, string, error)
	IsCloudGov() bool
	IsResourceSupported(resource api.AtlasCustomResource) bool
}

type ClientSet struct {
	SdkClient20231115008 *adminv20231115008.APIClient
	SdkClient20241113001 *adminv20241113001.APIClient
}

type ProductionProvider struct {
	domain string
	dryRun bool
}

// Credentials is the type that holds credentials to authenticate against the Atlas API.
// Currently, only API keys are support but more credential types could be added,
// see https://www.mongodb.com/docs/atlas/configure-api-access/.
type Credentials struct {
	APIKeys *APIKeys
}

// APIKeys is the type that holds Public/Private API keys to authenticate against the Atlas API.
type APIKeys struct {
	OrgID      string
	PublicKey  string
	PrivateKey string
}

func NewProductionProvider(atlasDomain string, dryRun bool) *ProductionProvider {
	return &ProductionProvider{
		domain: atlasDomain,
		dryRun: dryRun,
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

		return !(atlasResource.IsServerless() || atlasResource.IsFlex() || hasSearchNodes)
	}

	return false
}

func (p *ProductionProvider) Client(ctx context.Context, creds *Credentials, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
	clientCfg := []httputil.ClientOpt{
		httputil.Digest(creds.APIKeys.PublicKey, creds.APIKeys.PrivateKey),
		httputil.LoggingTransport(log),
	}

	transport := p.newDryRunTransport(http.DefaultTransport)
	httpClient, err := httputil.DecorateClient(&http.Client{Transport: transport}, clientCfg...)
	if err != nil {
		return nil, "", err
	}

	c, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL(p.domain), mongodbatlas.SetUserAgent(operatorUserAgent()))

	return c, creds.APIKeys.OrgID, err
}

func (p *ProductionProvider) SdkClientSet(ctx context.Context, creds *Credentials, log *zap.SugaredLogger) (*ClientSet, string, error) {
	var transport http.RoundTripper = digest.NewTransport(creds.APIKeys.PublicKey, creds.APIKeys.PrivateKey)
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
	}, creds.APIKeys.OrgID, nil
}

func (p *ProductionProvider) newDryRunTransport(delegate http.RoundTripper) http.RoundTripper {
	if p.dryRun {
		return dryrun.NewDryRunTransport(delegate)
	}

	return delegate
}

func operatorUserAgent() string {
	return fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", version.Version, runtime.GOOS, runtime.GOARCH)
}
