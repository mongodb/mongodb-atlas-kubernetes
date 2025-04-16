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

// nolint:gosec,stylecheck
const AtlasStreamConnectionBySecretIndex = "atlasstreamconnection.spec.kafkaConfig"

type AtlasStreamConnectionBySecretIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasStreamConnectionBySecretIndexer(logger *zap.Logger) *AtlasStreamConnectionBySecretIndexer {
	return &AtlasStreamConnectionBySecretIndexer{
		logger: logger.Named(AtlasStreamConnectionBySecretIndex).Sugar(),
	}
}

func (*AtlasStreamConnectionBySecretIndexer) Object() client.Object {
	return &akov2.AtlasStreamConnection{}
}

func (*AtlasStreamConnectionBySecretIndexer) Name() string {
	return AtlasStreamConnectionBySecretIndex
}

func (a *AtlasStreamConnectionBySecretIndexer) Keys(object client.Object) []string {
	streamConnection, ok := object.(*akov2.AtlasStreamConnection)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasStreamConnection but got %T", object)
		return nil
	}

	var indexes []string

	key, found := credentialSecretKey(streamConnection)
	if found {
		indexes = append(indexes, key)
	}

	key, found = certificateSecretKey(streamConnection)
	if found {
		indexes = append(indexes, key)
	}

	return indexes
}

func credentialSecretKey(connection *akov2.AtlasStreamConnection) (string, bool) {
	if connection == nil || connection.Spec.KafkaConfig == nil || connection.Spec.KafkaConfig.Authentication.Credentials.Name == "" {
		return "", false
	}

	credentialsKey := connection.Spec.KafkaConfig.Authentication.Credentials.GetObject(connection.GetNamespace())

	return credentialsKey.String(), true
}

func certificateSecretKey(connection *akov2.AtlasStreamConnection) (string, bool) {
	if connection == nil || connection.Spec.KafkaConfig == nil || connection.Spec.KafkaConfig.Security.Certificate.Name == "" {
		return "", false
	}

	certificateKey := connection.Spec.KafkaConfig.Security.Certificate.GetObject(connection.GetNamespace())

	return certificateKey.String(), true
}
