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
	AtlasDatabaseUserBySecretsIndex = "atlasdatabaseuser.spec.passwordSecret"
)

type AtlasDatabaseUserBySecretsIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDatabaseUserBySecretsIndexer(logger *zap.Logger) *AtlasDatabaseUserBySecretsIndexer {
	return &AtlasDatabaseUserBySecretsIndexer{
		logger: logger.Named(AtlasDatabaseUserBySecretsIndex).Sugar(),
	}
}

func (*AtlasDatabaseUserBySecretsIndexer) Object() client.Object {
	return &akov2.AtlasDatabaseUser{}
}

func (*AtlasDatabaseUserBySecretsIndexer) Name() string {
	return AtlasDatabaseUserBySecretsIndex
}

func (a *AtlasDatabaseUserBySecretsIndexer) Keys(object client.Object) []string {
	user, ok := object.(*akov2.AtlasDatabaseUser)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDatabaseUser but got %T", object)
		return nil
	}

	if user.Spec.PasswordSecret == nil || user.Spec.PasswordSecret.Name == "" {
		return nil
	}
	secretKey := user.PasswordSecretObjectKey()
	if secretKey == nil {
		return nil
	}

	return []string{secretKey.String()}
}
