package v1

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type AtlasSearch struct {
	// List of databases with indexes
	Databases []AtlasSearchDatabase `json:"databases"`
	// List of custom-analyzers
	CustomAnalyzers []CustomAnalyzer `json:"customAnalyzers,omitempty"`
}

type AtlasSearchDatabase struct {
	// The name of the database the indexes are in
	Database string `json:"database"`
	// List of collection with indexes
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
	// +kubebuilder:default:lucene.standard
	// +optional
	Analyzer string `json:"analyzer,omitempty"`
	// The analyzer to use for query the collection data
	// +kubebuilder:default:lucene.standard
	// +optional
	SearchAnalyzer string `json:"searchAnalyzer,omitempty"`
	//Synonyms mapping definition to use in this index
	// +optional
	Synonyms []AtlasSearchSynonym `json:"synonyms,omitempty"`
}

// +k8s:deepcopy-gen=false
//
//nolint:stylecheck
type FieldMapping map[string]interface{}

func (in *FieldMapping) DeepCopyInto(out *FieldMapping) {
	if in != nil {
		*out = make(FieldMapping)

		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

func (in FieldMapping) DeepCopy() FieldMapping {
	if in != nil {
		out := new(FieldMapping)
		in.DeepCopyInto(out)
		return *out
	}

	return nil
}

type IndexMapping struct {
	// Flag indicating whether the index uses dynamic or static mappings
	Dynamic bool `json:"dynamic"`
	// Map containing one or more field specifications.
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Fields *FieldMapping `json:"fields,omitempty"`
}

type CustomAnalyzer struct {
	// Name of the custom-analyzer
	Name string `json:"name"`
	// Analyzer on which the custom-analyzer is based
	BaseAnalyzer string `json:"baseAnalyzer"`
	// Specify whether the index is case-sensitive
	// +optional
	IgnoreCase *bool `json:"ignoreCase,omitempty"`
	// Longest text unit to analyze. Atlas Search excludes anything longer from the index
	// +optional
	MaxTokenLength *int `json:"maxTokenLength,omitempty"`
	// Words to exclude from stemming by the language analyzer
	// +optional
	StemExclusionSet []string `json:"stemExclusionSet,omitempty"`
	// Strings to ignore when creating the index
	// +optional
	Stopwords []string `json:"stopwords,omitempty"`
}

type AtlasSearchSynonym struct {
	// Name of the synonym mapping
	Name string `json:"name"`
	// Name of the analyzer to use with this synonym mapping
	Analyzer string `json:"analyzer"`
	// Source collection for synonyms
	Source SynonymSource `json:"source"`
}

type SynonymSource struct {
	// Name of the MongoDB collection that is in the same database as the Atlas Search index
	Collection string `json:"collection"`
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
			Fields:  (*map[string]interface{})(i.Mappings.Fields),
		},
		Synonyms: nil,
	}

	if len(i.Synonyms) > 0 {
		index.Synonyms = make([]map[string]interface{}, 0, len(i.Synonyms))
		for _, synonym := range i.Synonyms {
			s := map[string]interface{}{
				"name":     synonym.Name,
				"analyzer": synonym.Analyzer,
				"source":   synonym.Source,
			}
			index.Synonyms = append(index.Synonyms, s)
		}
	}

	a, _ := json.MarshalIndent(index, "", "    ")
	fmt.Println(string(a))

	return index
}

func (a *CustomAnalyzer) ToAtlas() *mongodbatlas.SearchAnalyzer {
	return &mongodbatlas.SearchAnalyzer{
		Name:             a.Name,
		BaseAnalyzer:     a.BaseAnalyzer,
		IgnoreCase:       a.IgnoreCase,
		MaxTokenLength:   a.MaxTokenLength,
		StemExclusionSet: a.StemExclusionSet,
		Stopwords:        a.Stopwords,
	}
}
