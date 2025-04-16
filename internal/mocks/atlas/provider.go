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
