package atlas

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

type TestProvider struct {
	ClientFunc       func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkSetClientFunc func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error)
	IsCloudGovFunc   func() bool
	IsSupportedFunc  func() bool
}

func (f *TestProvider) Client(_ context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
	return f.ClientFunc(secretRef, log)
}

func (f *TestProvider) IsCloudGov() bool {
	return f.IsCloudGovFunc()
}

func (f *TestProvider) SdkClientSet(ctx context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
	return f.SdkSetClientFunc(secretRef, log)
}

func (f *TestProvider) IsResourceSupported(_ api.AtlasCustomResource) bool {
	return f.IsSupportedFunc()
}
