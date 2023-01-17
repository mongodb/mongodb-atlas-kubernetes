package v1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"go.mongodb.org/atlas/mongodbatlas"
)

type AtlasSearch struct {
	// The name of the database the collection
	Database string `json:"database"`
	// List of collection with index of the database
	// +kubebuilder:validation:MinItems:1
	Collections []AtlasSearchCollection `json:"collections"`
}

type AtlasSearchCollection struct {
	// The name of the collection the indexes are on
	CollectionName string `json:"collectionName"`
	// List of indexes for the collection
	// +kubebuilder:validation:MinItems:1
	Indexes []SearchIndex `json:"indexes"`
}

type SearchIndex struct {
	// The name of the index
	Name string `json:"name"`
	// Object containing index specifications for the collection fields
	Mappings IndexMapping `json:"mappings"`
	// The analyzer to use for indexing the collection data
	// +kubebuilder:validation:Enum=lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword
	// +kubebuilder:default:lucene.standard
	// +optional
	Analyzer string `json:"analyzer,omitempty"`
	// The analyzer to use for query the collection data
	// +kubebuilder:validation:Enum=lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword
	// +kubebuilder:default:lucene.standard
	// +optional
	SearchAnalyzer string `json:"searchAnalyzer,omitempty"`
}

type IndexMapping struct {
	// Flag indicating whether the index uses dynamic or static mappings
	Dynamic bool `json:"dynamic"`
	// Map containing one or more field specifications.
	Fields map[string][]unstructured.Unstructured `json:"fields,omitempty"`
}

func (i *SearchIndex) ToAtlas(database string, collection string) *mongodbatlas.SearchIndex {
	index := &mongodbatlas.SearchIndex{
		Name:           i.Name,
		Database:       database,
		CollectionName: collection,
		Analyzer:       i.Analyzer,
		SearchAnalyzer: i.SearchAnalyzer,
		Mappings: &mongodbatlas.IndexMapping{
			Dynamic: i.Mappings.Dynamic,
			Fields:  &map[string]interface{}{},
		},
		Synonyms: nil,
	}

	for key, field := range i.Mappings.Fields {
		i.Mappings.Fields[key] = field
	}

	return index
}

func (i *SearchIndex) IsEqual(index *mongodbatlas.SearchIndex) bool {
	if i.Name != index.Name {
		return false
	}

	if i.Analyzer != index.Analyzer {
		return false
	}

	if i.SearchAnalyzer != index.SearchAnalyzer {
		return false
	}

	if i.Mappings.Dynamic != index.Mappings.Dynamic {
		return false
	}

	if !reflect.DeepEqual(i.Mappings.Fields, index.Mappings.Fields) {
		return false
	}

	return true
}
