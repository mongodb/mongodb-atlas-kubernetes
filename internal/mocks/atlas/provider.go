package atlas

import (
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/httputil"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestProvider struct {
	CreateConnectionFunc func(secretRef *client.ObjectKey) (atlas.Connection, error)
	CreateClientFunc     func() (mongodbatlas.Client, error)
	IsCloudGovFunc       func() bool
	IsSupportedFunc      func() bool
}

func (f *TestProvider) CreateConnection(secretRef *client.ObjectKey, _ *zap.SugaredLogger) (atlas.Connection, error) {
	return f.CreateConnectionFunc(secretRef)
}

func (f *TestProvider) CreateClient(_ *atlas.Connection, _ *zap.SugaredLogger, _ ...httputil.ClientOpt) (mongodbatlas.Client, error) {
	return f.CreateClientFunc()
}

func (f *TestProvider) IsCloudGov() bool {
	return f.IsCloudGovFunc()
}

func (f *TestProvider) IsResourceSupported(_ mdbv1.AtlasCustomResource) bool {
	return f.IsSupportedFunc()
}
