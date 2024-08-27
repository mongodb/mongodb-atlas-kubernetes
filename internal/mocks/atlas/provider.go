package atlas

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

type TestProvider struct {
	ClientFunc               func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error)
	SdkClientFunc            func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error)
	IsCloudGovFunc           func() bool
	IsSupportedFunc          func() bool
	GlobalFallbackSecretFunc func() *client.ObjectKey
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

func (f *TestProvider) IsResourceSupported(_ api.AtlasCustomResource) bool {
	return f.IsSupportedFunc()
}

func (f *TestProvider) GlobalFallbackSecret() *client.ObjectKey {
	return f.GlobalFallbackSecretFunc()
}
