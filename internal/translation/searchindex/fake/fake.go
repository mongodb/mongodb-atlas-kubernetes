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
	if fas.GetIndexFunc == nil {
		panic("unimplemented")
	}
	return fas.GetIndexFunc(ctx, projectID, clusterName, indexName, indexID)
}

func (fas *FakeAtlasSearch) CreateIndex(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
	if fas.CreateIndexFunc == nil {
		panic("unimplemented")
	}
	return fas.CreateIndexFunc(ctx, projectID, clusterName, index)
}

func (fas *FakeAtlasSearch) DeleteIndex(ctx context.Context, projectID, clusterName, indexID string) error {
	if fas.DeleteIndexFunc == nil {
		panic("unimplemented")
	}
	return fas.DeleteIndexFunc(ctx, projectID, clusterName, indexID)
}

func (fas *FakeAtlasSearch) UpdateIndex(ctx context.Context, projectID, clusterName string, index *searchindex.SearchIndex) (*searchindex.SearchIndex, error) {
	if fas.UpdateIndexFunc == nil {
		panic("unimplemented")
	}
	return fas.UpdateIndexFunc(ctx, projectID, clusterName, index)
}
