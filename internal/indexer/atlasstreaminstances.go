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
	AtlasStreamInstanceByConnectionIndex = "atlasstreaminstance.spec.connectionRegistry"
	AtlasStreamInstanceByProjectIndex    = "atlasstreaminstance.spec.projectRef"
)

type AtlasStreamInstanceByConnectionIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasStreamInstanceByConnectionIndexer(logger *zap.Logger) *AtlasStreamInstanceByConnectionIndexer {
	return &AtlasStreamInstanceByConnectionIndexer{
		logger: logger.Named(AtlasStreamInstanceByConnectionIndex).Sugar(),
	}
}

func (*AtlasStreamInstanceByConnectionIndexer) Object() client.Object {
	return &akov2.AtlasStreamInstance{}
}

func (*AtlasStreamInstanceByConnectionIndexer) Name() string {
	return AtlasStreamInstanceByConnectionIndex
}

func (a *AtlasStreamInstanceByConnectionIndexer) Keys(object client.Object) []string {
	streamInstance, ok := object.(*akov2.AtlasStreamInstance)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasStreamInstance but got %T", object)
		return nil
	}

	if len(streamInstance.Spec.ConnectionRegistry) == 0 {
		return nil
	}

	registry := streamInstance.Spec.ConnectionRegistry
	indices := make([]string, 0, len(registry))
	for i := range registry {
		key := registry[i].GetObject(streamInstance.GetNamespace())
		indices = append(indices, key.String())
	}

	return indices
}

type AtlasStreamInstanceByProjectIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasStreamInstanceByProjectIndexer(logger *zap.Logger) *AtlasStreamInstanceByProjectIndexer {
	return &AtlasStreamInstanceByProjectIndexer{
		logger: logger.Named(AtlasStreamInstanceByProjectIndex).Sugar(),
	}
}

func (*AtlasStreamInstanceByProjectIndexer) Object() client.Object {
	return &akov2.AtlasStreamInstance{}
}

func (*AtlasStreamInstanceByProjectIndexer) Name() string {
	return AtlasStreamInstanceByProjectIndex
}

func (a *AtlasStreamInstanceByProjectIndexer) Keys(object client.Object) []string {
	streamInstance, ok := object.(*akov2.AtlasStreamInstance)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasStreamInstance but got %T", object)
		return nil
	}

	if streamInstance.Spec.Project.Name == "" {
		return nil
	}

	key := streamInstance.Spec.Project.GetObject(streamInstance.GetNamespace())

	return []string{key.String()}
}
