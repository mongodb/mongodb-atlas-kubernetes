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
)

const (
	AtlasProjectByTeamIndex = "atlasproject.spec.teams"
)

type AtlasProjectByTeamIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasProjectByTeamIndexer(logger *zap.Logger) *AtlasProjectByTeamIndexer {
	return &AtlasProjectByTeamIndexer{
		logger: logger.Named(AtlasProjectByTeamIndex).Sugar(),
	}
}

func (*AtlasProjectByTeamIndexer) Object() client.Object {
	return &akov2.AtlasProject{}
}

func (*AtlasProjectByTeamIndexer) Name() string {
	return AtlasProjectByTeamIndex
}

func (a *AtlasProjectByTeamIndexer) Keys(object client.Object) []string {
	project, ok := object.(*akov2.AtlasProject)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasProject but got %T", object)
		return nil
	}

	result := sets.New[string]()
	for _, team := range project.Spec.Teams {
		if team.TeamRef.IsEmpty() {
			continue
		}
		result.Insert(team.TeamRef.GetObject(project.Namespace).String())
	}

	return result.UnsortedList()
}
