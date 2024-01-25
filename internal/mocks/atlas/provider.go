package atlas

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestProvider struct {
	ClientFunc      func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkClientFunc   func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error)
	IsCloudGovFunc  func() bool
	IsSupportedFunc func() bool
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

func (f *TestProvider) IsResourceSupported(_ mdbv1.AtlasCustomResource) bool {
	return f.IsSupportedFunc()
}
