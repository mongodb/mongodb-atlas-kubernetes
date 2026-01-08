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

package v1

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

// SearchIndex is the CRD to configure part of the Atlas Search Index.
type SearchIndex struct {
	// Human-readable label that identifies this index. Must be unique for a deployment.
	// +required
	Name string `json:"name"`
	// Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.
	// +required
	DBName string `json:"DBName"`
	// Human-readable label that identifies the collection that contains one or more Atlas Search indexes.
	// +required
	CollectionName string `json:"collectionName"`
	// Type of the index.
	// +kubebuilder:validation:Enum:=search;vectorSearch
	// +required
	Type string `json:"type,omitempty"`
	// Atlas search index configuration.
	// +optional
	Search *Search `json:"search,omitempty"`
	// Atlas vector search index configuration.
	// +optional
	VectorSearch *VectorSearch `json:"vectorSearch,omitempty"`
}

// Search represents "search" type of Atlas Search Index.
type Search struct {
	// Rule sets that map words to their synonyms in this index.
	// +optional
	Synonyms *[]Synonym `json:"synonyms,omitempty"`
	// Index specifications for the collection's fields.
	// +required
	Mappings *Mappings `json:"mappings,omitempty"`
	// A reference to the AtlasSearchIndexConfig custom resource.
	// +required
	SearchConfigurationRef common.ResourceRefNamespaced `json:"searchConfigurationRef"`
}

// Synonym represents "Synonym" type of Atlas Search Index.
type Synonym struct {
	// Human-readable label that identifies the synonym definition. Each name must be unique within the same index definition.
	// +required
	Name string `json:"name"`
	// Specific pre-defined method chosen to apply to the synonyms to be searched.
	// +kubebuilder:validation:Enum:=lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword;lucene.arabic;lucene.armenian;lucene.basque;lucene.bengali;lucene.brazilian;lucene.bulgarian;lucene.catalan;lucene.chinese;lucene.cjk;lucene.czech;lucene.danish;lucene.dutch;lucene.english;lucene.finnish;lucene.french;lucene.galician;lucene.german;lucene.greek;lucene.hindi;lucene.hungarian;lucene.indonesian;lucene.irish;lucene.italian;lucene.japanese;lucene.korean;lucene.kuromoji;lucene.latvian;lucene.lithuanian;lucene.morfologik;lucene.nori;lucene.norwegian;lucene.persian;lucene.portuguese;lucene.romanian;lucene.russian;lucene.smartcn;lucene.sorani;lucene.spanish;lucene.swedish;lucene.thai;lucene.turkish;lucene.ukrainian
	// +required
	Analyzer string `json:"analyzer"`
	// Data set that stores the mapping one or more words map to one or more synonyms of those words.
	// +required
	Source Source `json:"source"`
}

// Source represents "Source" type of Atlas Search Index
type Source struct {
	// Human-readable label that identifies the MongoDB collection that stores words and their applicable synonyms.
	Collection string `json:"collection"`
}

// Mappings represents "mappings" type of Atlas Search Index
type Mappings struct {
	// Indicates whether the index uses static, default dynamic, or configurable dynamic mappings.
	// Set to **true** to enable dynamic mapping with default type set or define object to specify the name of the configured type sets for dynamic mapping.
	// If you specify configurable dynamic mappings, you must define the referred type sets in the **typeSets** field.
	// Set to **false** to use only static mappings through **mappings.fields**.
	// See https://www.mongodb.com/docs/atlas/atlas-search/define-field-mappings/#configure-a-typeset for more details.
	Dynamic *apiextensions.JSON `json:"dynamic,omitempty"`
	// One or more field specifications for the Atlas Search index. Required if mapping.dynamic is omitted or set to false.
	Fields *apiextensions.JSON `json:"fields,omitempty"`
}

// VectorSearch represents "vectorSearch" type of Atlas Search Index
type VectorSearch struct {
	// Array of JSON objects. See examples https://dochub.mongodb.org/core/avs-vector-type
	// +required
	Fields *apiextensions.JSON `json:"fields,omitempty"`
}
