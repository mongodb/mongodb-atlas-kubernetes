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

package fake

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
)

type FakeAtlasSearch struct {
	GetIndexFunc    func(ctx context.Context, projectID, clusterName, indexName, indexID string) (*searchindex.SearchIndex, error)
	CreateIndexFunc func(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error)
	DeleteIndexFunc func(ctx context.Context, projectID, clusterName, indexID string) error
	UpdateIndexFunc func(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error)
}

func (fas *FakeAtlasSearch) GetIndex(ctx context.Context, projectID, clusterName, indexName, indexID string) (*searchindex.SearchIndex, error) {
	return fas.GetIndexFunc(ctx, projectID, clusterName, indexName, indexID)
}

func (fas *FakeAtlasSearch) CreateIndex(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
	return fas.CreateIndexFunc(ctx, projectID, clusterName, index)
}

func (fas *FakeAtlasSearch) DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) error {
	return fas.DeleteIndexFunc(ctx, projectID, clusterName, indexID)
}

func (fas *FakeAtlasSearch) UpdateIndex(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
	return fas.UpdateIndexFunc(ctx, projectID, clusterName, index)
}
