package fake

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
)

type FakeAtlasSearch struct {
	GetIndexFn    func(ctx context.Context, projectID, clusterName, indexName, indexID string) (*searchindex.SearchIndex, error)
	CreateIndexFn func(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error)
	DeleteIndexFn func(ctx context.Context, projectID, clusterName, indexID string) error
	UpdateIndexFn func(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error)
}

func (fas *FakeAtlasSearch) GetIndex(ctx context.Context, projectID, clusterName, indexName, indexID string) (*searchindex.SearchIndex, error) {
	if fas.GetIndexFn == nil {
		panic("unimplemented")
	}
	return fas.GetIndexFn(ctx, projectID, clusterName, indexName, indexID)
}

func (fas *FakeAtlasSearch) CreateIndex(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
	if fas.CreateIndexFn == nil {
		panic("unimplemented")
	}
	return fas.CreateIndexFn(ctx, projectID, clusterName, index)
}

func (fas *FakeAtlasSearch) DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) error {
	if fas.DeleteIndexFn == nil {
		panic("unimplemented")
	}
	return fas.DeleteIndexFn(ctx, projectID, clusterName, indexID)
}

func (fas *FakeAtlasSearch) UpdateIndex(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
	if fas.UpdateIndexFn == nil {
		panic("unimplemented")
	}
	return fas.UpdateIndexFn(ctx, projectID, clusterName, index)
}
