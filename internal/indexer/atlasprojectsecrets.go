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
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

const (
	AtlasProjectBySecretsIndex = "atlasproject.spec.secrets"
)

type AtlasProjectByConnectionSecretIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasProjectByConnectionSecretIndexer(logger *zap.Logger) *AtlasProjectByConnectionSecretIndexer {
	return &AtlasProjectByConnectionSecretIndexer{
		logger: logger.Named(AtlasProjectBySecretsIndex).Sugar(),
	}
}

func (AtlasProjectByConnectionSecretIndexer) Object() client.Object {
	return &akov2.AtlasProject{}
}

func (*AtlasProjectByConnectionSecretIndexer) Name() string {
	return AtlasProjectBySecretsIndex
}

func (a *AtlasProjectByConnectionSecretIndexer) Keys(object client.Object) []string {
	project, ok := object.(*akov2.AtlasProject)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasProject but got %T", object)
		return nil
	}

	result := sets.New[string]()
	addIfNotEmpty := func(ref *common.ResourceRefNamespaced) {
		if !ref.IsEmpty() {
			result.Insert(ref.GetObject(project.Namespace).String())
		}
	}

	if project.Spec.ConnectionSecret != nil {
		addIfNotEmpty(project.Spec.ConnectionSecret)
	}

	if project.Spec.EncryptionAtRest != nil {
		encryptionAtRest := project.Spec.EncryptionAtRest
		addIfNotEmpty(&encryptionAtRest.AwsKms.SecretRef)
		addIfNotEmpty(&encryptionAtRest.AzureKeyVault.SecretRef)
		addIfNotEmpty(&encryptionAtRest.GoogleCloudKms.SecretRef)
	}

	for i := range project.Spec.AlertConfigurations {
		for j := range project.Spec.AlertConfigurations[i].Notifications {
			notification := &project.Spec.AlertConfigurations[i].Notifications[j]
			addIfNotEmpty(&notification.APITokenRef)
			addIfNotEmpty(&notification.DatadogAPIKeyRef)
			addIfNotEmpty(&notification.FlowdockAPITokenRef)
			addIfNotEmpty(&notification.OpsGenieAPIKeyRef)
			addIfNotEmpty(&notification.ServiceKeyRef)
			addIfNotEmpty(&notification.VictorOpsSecretRef)
		}
	}

	return result.UnsortedList()
}
