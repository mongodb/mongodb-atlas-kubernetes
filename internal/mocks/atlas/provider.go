package atlas

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

type TestProvider struct {
	ClientFunc       func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*mongodbatlas.Client, error)
	SdkClientSetFunc func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error)
	IsCloudGovFunc   func() bool
	IsSupportedFunc  func() bool
}

func (f *TestProvider) Client(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*mongodbatlas.Client, error) {
	return f.ClientFunc(ctx, creds, log)
}

func (f *TestProvider) IsCloudGov() bool {
	return f.IsCloudGovFunc()
}

func (f *TestProvider) SdkClientSet(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
	return f.SdkClientSetFunc(ctx, creds, log)
}

func (f *TestProvider) IsResourceSupported(_ api.AtlasCustomResource) bool {
	return f.IsSupportedFunc()
}
