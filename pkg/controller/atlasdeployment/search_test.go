package atlasdeployment

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
)

func TestSortIndexes(t *testing.T) {
	t.Run("should sort indexes operations", func(t *testing.T) {
		desired := map[string]*mongodbatlas.SearchIndex{
			"db.col.index2": {Name: "index2", Analyzer: "lucene.simple", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
			"db.col.index3": {Name: "index3", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
			"db.col.index4": {Name: "index4", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
		}
		existing := map[string]*mongodbatlas.SearchIndex{
			"db.col.index1": {IndexID: "id1", Name: "index1", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
			"db.col.index2": {IndexID: "id2", Name: "index2", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
			"db.col.index3": {IndexID: "id3", Name: "index3", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
		}
		statuses := []*status.AtlasIndex{
			{ID: "id2", Name: "index2", Database: "db", CollectionName: "col"},
		}

		toCreate, toUpdate, toDelete, syncStatus := sortIndexes(desired, existing, statuses)

		assert.Equal(
			t,
			[]*mongodbatlas.SearchIndex{
				{Name: "index4", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
			},
			toCreate,
		)
		assert.Equal(
			t,
			map[string]*mongodbatlas.SearchIndex{
				"id2": {Name: "index2", Analyzer: "lucene.simple", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}}},
			toUpdate,
		)
		assert.Equal(t, []string{"id1"}, toDelete)
		assert.Equal(
			t,
			[]*mongodbatlas.SearchIndex{
				{IndexID: "id3", Name: "index3", Mappings: &mongodbatlas.IndexMapping{Dynamic: true}},
			},
			syncStatus,
		)
	})
}

func TestHasIndexChanged(t *testing.T) {
	t.Run("should return true when indexes are equal", func(t *testing.T) {
		assert.True(
			t,
			hasIndexChanged(
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
			),
		)
	})

	t.Run("should return false when indexes have different names", func(t *testing.T) {
		assert.False(
			t,
			hasIndexChanged(
				&mongodbatlas.SearchIndex{
					Name:           "myIndex1",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
				&mongodbatlas.SearchIndex{
					Name:           "myIndex2",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
			),
		)
	})

	t.Run("should return false when indexes have different analyzers", func(t *testing.T) {
		assert.False(
			t,
			hasIndexChanged(
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.simple",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
			),
		)
	})

	t.Run("should return false when indexes have different search analyzers", func(t *testing.T) {
		assert.False(
			t,
			hasIndexChanged(
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.simple",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
			),
		)
	})

	t.Run("should return false when indexes have different mappings", func(t *testing.T) {
		assert.False(
			t,
			hasIndexChanged(
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: true,
					},
				},
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
			),
		)
	})

	t.Run("should return false when indexes have different fields", func(t *testing.T) {
		assert.False(
			t,
			hasIndexChanged(
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.simple",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type":    "document",
								"dynamic": true,
							},
						},
					},
				},
				&mongodbatlas.SearchIndex{
					Name:           "myIndex",
					Analyzer:       "lucene.standard",
					SearchAnalyzer: "lucene.standard",
					Mappings: &mongodbatlas.IndexMapping{
						Dynamic: false,
						Fields: &map[string]interface{}{
							"field1": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
			),
		)
	})
}
