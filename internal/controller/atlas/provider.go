// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/mongodb-forks/digest"
	v20250312 "go.mongodb.org/atlas-sdk/v20250312018/admin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/deprecation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

const (
	govAtlasDomain = "mongodbgov.com"
)

type Provider interface {
	SdkClientSet(ctx context.Context, creds *Credentials, log *zap.SugaredLogger) (*ClientSet, error)
	IsCloudGov() bool
	IsResourceSupported(resource api.AtlasCustomResource) bool
}

type ClientSet struct {
	SdkClient20250312 *v20250312.APIClient
}

type ProductionProvider struct {
	domain       string
	dryRun       bool
	isLogInDebug bool
}

// ConnectionConfig is the type that contains connection configuration to Atlas, including credentials.
type ConnectionConfig struct {
	OrgID       string
	Credentials *Credentials
}

// Credentials is the type that holds credentials to authenticate against the Atlas API.
// Either APIKeys or ServiceAccount must be set, but not both.
// See https://www.mongodb.com/docs/atlas/configure-api-access/.
type Credentials struct {
	APIKeys        *APIKeys
	ServiceAccount *ServiceAccountToken
}

// APIKeys is the type that holds Public/Private API keys to authenticate against the Atlas API.
type APIKeys struct {
	PublicKey  string
	PrivateKey string
}

// ServiceAccountToken holds a pre-fetched OAuth2 bearer token obtained
// by the service-account controller via the client credentials flow.
type ServiceAccountToken struct {
	BearerToken string
}

func NewProductionProvider(atlasDomain string, dryRun, isLogInDebug bool) *ProductionProvider {
	return &ProductionProvider{
		domain:       atlasDomain,
		dryRun:       dryRun,
		isLogInDebug: isLogInDebug,
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
		*akov2.AtlasPrivateEndpoint,
		*akov2.AtlasNetworkContainer,
		*akov2.AtlasNetworkPeering:
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

func (p *ProductionProvider) SdkClientSet(ctx context.Context, creds *Credentials, log *zap.SugaredLogger) (*ClientSet, error) {
	var baseTransport http.RoundTripper
	switch {
	case creds.ServiceAccount != nil:
		baseTransport = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: creds.ServiceAccount.BearerToken,
			}),
		}
	case creds.APIKeys != nil:
		baseTransport = digest.NewTransport(creds.APIKeys.PublicKey, creds.APIKeys.PrivateKey)
	default:
		return nil, fmt.Errorf("no credentials provided")
	}

	transport := p.newTransport(baseTransport, log)
	transport = httputil.NewLoggingTransport(log, false, transport)
	if p.isLogInDebug {
		log.Debug("JSON payload diff is enabled for Atlas API requests (PATCH & PUT)")
		transport = httputil.NewTransportWithDiff(transport, log.Named("payload_diff"))
	}

	httpClient := &http.Client{Transport: transport}

	return NewSDKClientSet(p.domain, httpClient)
}

func (p *ProductionProvider) newTransport(delegate http.RoundTripper, log *zap.SugaredLogger) http.RoundTripper {
	if os.Getenv("AKO_DEPRECATION_WARNINGS") != "" {
		delegate = deprecation.NewLoggingTransport(delegate, log.Desugar())
	}

	if p.dryRun {
		return dryrun.NewDryRunTransport(delegate)
	}

	return delegate
}

func operatorUserAgent() string {
	return fmt.Sprintf("%s/%s (%s;%s)", "MongoDBAtlasKubernetesOperator", version.Version, runtime.GOOS, runtime.GOARCH)
}

func NewSDKClientSet(domain string, httpClient *http.Client) (*ClientSet, error) {
	clientv20250312, err := v20250312.NewClient(
		v20250312.UseBaseURL(domain),
		v20250312.UseHTTPClient(httpClient),
		v20250312.UseUserAgent(operatorUserAgent()))
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		SdkClient20250312: clientv20250312,
	}, nil
}
