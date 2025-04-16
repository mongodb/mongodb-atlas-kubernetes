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
	AtlasDeploymentBySearchIndexIndex = "atlasdeployment.spec.deploymentSpec.searchIndexes"
)

type AtlasDeploymentBySearchIndexIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDeploymentBySearchIndexIndexer(logger *zap.Logger) *AtlasDeploymentBySearchIndexIndexer {
	return &AtlasDeploymentBySearchIndexIndexer{
		logger: logger.Named(AtlasDeploymentBySearchIndexIndex).Sugar(),
	}
}

func (*AtlasDeploymentBySearchIndexIndexer) Object() client.Object {
	return &akov2.AtlasDeployment{}
}

func (*AtlasDeploymentBySearchIndexIndexer) Name() string {
	return AtlasDeploymentBySearchIndexIndex
}

func (a *AtlasDeploymentBySearchIndexIndexer) Keys(object client.Object) []string {
	deployment, ok := object.(*akov2.AtlasDeployment)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDeployment but got %T", object)
		return nil
	}

	if deployment.Spec.DeploymentSpec == nil {
		return nil
	}

	if len(deployment.Spec.DeploymentSpec.SearchIndexes) == 0 {
		return nil
	}

	searchIndexes := deployment.Spec.DeploymentSpec.SearchIndexes

	result := make([]string, 0, len(searchIndexes))
	for i := range searchIndexes {
		idx := &searchIndexes[i]
		if idx.Search == nil {
			continue
		}

		// searchIndexConfigKey -> deploymentName
		searchIndexKey := idx.Search.SearchConfigurationRef.GetObject(deployment.GetNamespace())
		result = append(result, searchIndexKey.String())
	}

	return result
}
