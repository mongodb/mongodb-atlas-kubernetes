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

package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasFederatedAuthBySecretsIndex = "atlasfederatedauth.spec.connectionSecret"
)

type AtlasFederatedAuthBySecretsIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasFederatedAuthBySecretsIndexer(logger *zap.Logger) *AtlasFederatedAuthBySecretsIndexer {
	return &AtlasFederatedAuthBySecretsIndexer{
		logger: logger.Named(AtlasFederatedAuthBySecretsIndex).Sugar(),
	}
}

func (*AtlasFederatedAuthBySecretsIndexer) Object() client.Object {
	return &akov2.AtlasFederatedAuth{}
}

func (*AtlasFederatedAuthBySecretsIndexer) Name() string {
	return AtlasFederatedAuthBySecretsIndex
}

func (a *AtlasFederatedAuthBySecretsIndexer) Keys(object client.Object) []string {
	fedAuth, ok := object.(*akov2.AtlasFederatedAuth)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasFederatedAuth but got %T", object)
		return nil
	}

	if fedAuth.Spec.ConnectionSecretRef.IsEmpty() {
		return nil
	}

	return []string{fedAuth.ConnectionSecretObjectKey().String()}
}
