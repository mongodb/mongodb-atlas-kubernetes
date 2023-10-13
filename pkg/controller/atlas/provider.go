package atlas

import (
	"net/url"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/httputil"
)

const govAtlasDomain = "mongodbgov.com"

type Provider interface {
	CreateConnection(secretRef *client.ObjectKey, log *zap.SugaredLogger) (Connection, error)
	CreateClient(connection *Connection, log *zap.SugaredLogger, opts ...httputil.ClientOpt) (mongodbatlas.Client, error)
	IsCloudGov() bool
	IsResourceSupported(resource mdbv1.AtlasCustomResource) bool
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

func (f *ProductionProvider) CreateConnection(secretRef *client.ObjectKey, log *zap.SugaredLogger) (Connection, error) {
	//TODO move implementation here once all controllers are using the manager
	return ReadConnection(log, f.k8sClient, f.globalSecretRef, secretRef)
}

func (f *ProductionProvider) CreateClient(connection *Connection, log *zap.SugaredLogger, opts ...httputil.ClientOpt) (mongodbatlas.Client, error) {
	//TODO move implementation here once all controllers are using the manager
	return Client(f.domain, *connection, log, opts...)
}

func (f *ProductionProvider) IsCloudGov() bool {
	domainURL, err := url.Parse(f.domain)
	if err != nil {
		return false
	}

	return strings.HasSuffix(domainURL.Hostname(), govAtlasDomain)
}

func (f *ProductionProvider) IsResourceSupported(resource mdbv1.AtlasCustomResource) bool {
	if !f.IsCloudGov() {
		return true
	}

	switch atlasResource := resource.(type) {
	case *mdbv1.AtlasProject,
		*mdbv1.AtlasTeam,
		*mdbv1.AtlasBackupSchedule,
		*mdbv1.AtlasBackupPolicy,
		*mdbv1.AtlasDatabaseUser,
		*mdbv1.AtlasFederatedAuth:
		return true
	case *mdbv1.AtlasDataFederation:
		return false
	case *mdbv1.AtlasDeployment:
		if atlasResource.Spec.ServerlessSpec == nil {
			return true
		}
	}

	return false
}
