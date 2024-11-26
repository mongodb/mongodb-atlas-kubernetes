package atlas

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

type TestProvider struct {
	ClientFunc       func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkClientFunc    func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error)
	SdkSetClientFunc func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error)
	IsCloudGovFunc   func() bool
	IsSupportedFunc  func() bool
}

func (f *TestProvider) Client(_ context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
	return f.ClientFunc(secretRef, log)
}

func (f *TestProvider) SdkClient(_ context.Context, secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
	return f.SdkClientFunc(secretRef, log)
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
