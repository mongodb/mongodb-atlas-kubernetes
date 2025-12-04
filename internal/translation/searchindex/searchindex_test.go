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

package searchindex

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func jsonMustEncode(jsn any) *apiextensionsv1.JSON {
	val, err := json.Marshal(jsn)
	if err != nil {
		panic(err)
	}
	return &apiextensionsv1.JSON{Raw: val}
}

func jsonMustEncodeMap(jsn map[string]interface{}) *apiextensionsv1.JSON {
	val, err := json.Marshal(jsn)
	if err != nil {
		panic(err)
	}
	return &apiextensionsv1.JSON{Raw: val}
}

func jsonMustDecode(jsn *apiextensionsv1.JSON) *[]interface{} {
	var val []interface{}
	err := json.Unmarshal(jsn.Raw, &val)
	if err != nil {
		panic(err)
	}
	return &val
}

func Test_NewSearchIndexFromAKO(t *testing.T) {
	type args struct {
		index  *akov2.SearchIndex
		config *akov2.AtlasSearchIndexConfigSpec
	}
	tests := []struct {
		name string
		args args
		want *SearchIndex
	}{
		{
			name: "Convert Atlas index to AKO internal",
			args: args{
				index: &akov2.SearchIndex{
					Name:           "TestIndex",
					DBName:         "TestDBName",
					CollectionName: "TestCollectionName",
					Type:           "search",
					Search: &akov2.Search{
						Synonyms: &([]akov2.Synonym{
							{
								Name:     "MySynonym",
								Analyzer: "lucene.standard",
								Source:   akov2.Source{Collection: "test-collection"},
							},
						}),
						Mappings: &akov2.Mappings{
							Dynamic: jsonMustEncode(true),
							Fields: jsonMustEncode([]map[string]interface{}{
								{
									"test": "value",
								},
							}),
						},
						SearchConfigurationRef: common.ResourceRefNamespaced{},
					},
					VectorSearch: &akov2.VectorSearch{Fields: jsonMustEncode([]map[string]interface{}{
						{"test": "value"},
					})},
				},
				config: &akov2.AtlasSearchIndexConfigSpec{
					Analyzer: pointer.MakePtr("lucene.standard"),
					Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
						{
							Name: "MyAnalyzer",
							TokenFilters: jsonMustEncode([]map[string]interface{}{
								{"token": "value"},
							}),
							CharFilters: jsonMustEncode([]map[string]interface{}{
								{"filter": "value"},
							}),
							Tokenizer: akov2.Tokenizer{
								Type:           pointer.MakePtr("standard"),
								MaxGram:        nil,
								MinGram:        nil,
								Group:          nil,
								Pattern:        nil,
								MaxTokenLength: pointer.MakePtr(255),
							},
						},
					}),
					SearchAnalyzer: pointer.MakePtr("lucene.standard"),
					StoredSource:   jsonMustEncodeMap(map[string]interface{}{"include": "test"}),
				},
			},
			want: &SearchIndex{
				SearchIndex: akov2.SearchIndex{
					Name:           "TestIndex",
					DBName:         "TestDBName",
					CollectionName: "TestCollectionName",
					Type:           "search",
					Search: &akov2.Search{
						Synonyms: &([]akov2.Synonym{
							{
								Name:     "MySynonym",
								Analyzer: "lucene.standard",
								Source:   akov2.Source{Collection: "test-collection"},
							},
						}),
						Mappings: &akov2.Mappings{
							Dynamic: jsonMustEncode(true),
							Fields: jsonMustEncode([]map[string]interface{}{
								{"test": "value"},
							}),
						},
						SearchConfigurationRef: common.ResourceRefNamespaced{},
					},
					VectorSearch: &akov2.VectorSearch{Fields: jsonMustEncode([]map[string]interface{}{
						{"test": "value"},
					})},
				},
				AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
					Analyzer: pointer.MakePtr("lucene.standard"),
					Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
						{
							Name: "MyAnalyzer",
							TokenFilters: jsonMustEncode([]map[string]interface{}{
								{"token": "value"},
							}),
							CharFilters: jsonMustEncode([]map[string]interface{}{
								{"filter": "value"},
							}),
							Tokenizer: akov2.Tokenizer{
								Type:           pointer.MakePtr("standard"),
								MaxGram:        nil,
								MinGram:        nil,
								Group:          nil,
								Pattern:        nil,
								MaxTokenLength: pointer.MakePtr(255),
							},
						},
					}),
					SearchAnalyzer: pointer.MakePtr("lucene.standard"),
					StoredSource:   jsonMustEncodeMap(map[string]interface{}{"include": "test"}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSearchIndex(tt.args.index, tt.args.config)
			if diff := cmp.Diff(got, tt.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("NewSearchIndexFromAKO() = %v, want %v.\nDiff: %s", got, tt.want, diff)
			}
		})
	}
}

func Test_NewSearchIndexFromAtlas(t *testing.T) {
	type args struct {
		index admin.SearchIndexResponse
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchIndex
		wantErr bool
	}{
		{
			name: "Convert from Atlas",
			args: args{
				index: admin.SearchIndexResponse{
					CollectionName: pointer.MakePtr("collection"),
					Database:       pointer.MakePtr("db"),
					IndexID:        pointer.MakePtr("indexID"),
					Name:           pointer.MakePtr("name"),
					Status:         pointer.MakePtr("ACTIVE"),
					Type:           pointer.MakePtr("search"),
					LatestDefinition: &admin.BaseSearchIndexResponseLatestDefinition{
						Analyzer: pointer.MakePtr("lucene.standard"),
						Analyzers: &([]admin.AtlasSearchAnalyzer{
							{
								CharFilters: jsonMustDecode(jsonMustEncode([]map[string]interface{}{
									{"char": "filter"},
								})),
								Name: "name",
								TokenFilters: jsonMustDecode(jsonMustEncode([]map[string]interface{}{
									{"token": "filter"},
								})),
								Tokenizer: admin.ApiAtlasFTSAnalyzersTokenizer{
									MaxGram:        pointer.MakePtr(20),
									MinGram:        pointer.MakePtr(10),
									Type:           pointer.MakePtr("standard"),
									Group:          pointer.MakePtr(10),
									Pattern:        pointer.MakePtr("testRegex"),
									MaxTokenLength: pointer.MakePtr(255),
								},
							},
						}),
						Mappings: &admin.SearchMappings{
							Dynamic: pointer.MakePtr(true),
							Fields:  &map[string]interface{}{"field": "value"},
						},
						SearchAnalyzer: pointer.MakePtr("search-analyzer"),
						Synonyms: &([]admin.SearchSynonymMappingDefinition{
							{
								Analyzer: "analyzer",
								Name:     "name",
								Source: admin.SynonymSource{
									Collection: "collection",
								},
							},
						}),
						Fields: &([]any{
							map[string]string{
								"testKey": "testValue",
							},
						}),
						StoredSource: map[string]interface{}{"include": "test"},
					},
				},
			},
			want: &SearchIndex{
				ID:     pointer.MakePtr("indexID"),
				Status: pointer.MakePtr("ACTIVE"),
				SearchIndex: akov2.SearchIndex{
					Name:           "name",
					DBName:         "db",
					CollectionName: "collection",
					Type:           "search",
					Search: &akov2.Search{
						Synonyms: &([]akov2.Synonym{
							{
								Name:     "name",
								Analyzer: "analyzer",
								Source:   akov2.Source{Collection: "collection"},
							},
						}),
						Mappings: &akov2.Mappings{
							Dynamic: jsonMustEncode(true),
							Fields:  jsonMustEncodeMap(map[string]interface{}{"field": "value"}),
						},
						SearchConfigurationRef: common.ResourceRefNamespaced{},
					},
					VectorSearch: &akov2.VectorSearch{
						Fields: jsonMustEncode([]map[string]interface{}{
							{
								"testKey": "testValue",
							},
						}),
					},
				},
				AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
					Analyzer: pointer.MakePtr("lucene.standard"),
					Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
						{
							Name:         "name",
							TokenFilters: jsonMustEncode([]map[string]interface{}{{"token": "filter"}}),
							CharFilters:  jsonMustEncode([]map[string]interface{}{{"char": "filter"}}),
							Tokenizer: akov2.Tokenizer{
								MaxGram:        pointer.MakePtr(20),
								MinGram:        pointer.MakePtr(10),
								Type:           pointer.MakePtr("standard"),
								Group:          pointer.MakePtr(10),
								Pattern:        pointer.MakePtr("testRegex"),
								MaxTokenLength: pointer.MakePtr(255),
							},
						},
					}),
					SearchAnalyzer: pointer.MakePtr("search-analyzer"),
					StoredSource:   jsonMustEncodeMap(map[string]interface{}{"include": "test"}),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromAtlas(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSearchIndexFromAtlas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("NewSearchIndexFromAtlas() = %v, want %v.\nDiff: %s", got, tt.want, diff)
			}
		})
	}
}

//nolint:dupl
func TestSearchIndex_EqualTo(t *testing.T) {
	t.Run("Indexes should be equal", func(t *testing.T) {
		idx1 := &SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "",
				DBName:         "",
				CollectionName: "",
				Type:           "",
				Search: &akov2.Search{
					Synonyms: &([]akov2.Synonym{}),
					Mappings: &akov2.Mappings{
						Dynamic: nil,
						Fields:  nil,
					},
					SearchConfigurationRef: common.ResourceRefNamespaced{},
				},
				VectorSearch: nil,
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("analyzer"),
				Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
					{
						Name:         "test",
						TokenFilters: jsonMustEncode([]map[string]interface{}{{"token": "filter"}, {"token1": "filter1"}}),
						CharFilters:  jsonMustEncode([]map[string]interface{}{{"char": "filter"}, {"char1": "filter1"}}),
						Tokenizer: akov2.Tokenizer{
							Type:           pointer.MakePtr("type"),
							MaxGram:        nil,
							MinGram:        nil,
							Group:          nil,
							Pattern:        nil,
							MaxTokenLength: nil,
						},
					},
				}),
				SearchAnalyzer: pointer.MakePtr("searchAnalyzer"),
				StoredSource:   nil,
			},
			ID:     nil,
			Status: nil,
		}
		idx2 := &SearchIndex{
			SearchIndex: akov2.SearchIndex{
				Name:           "",
				DBName:         "",
				CollectionName: "",
				Type:           "",
				Search: &akov2.Search{
					Synonyms: &([]akov2.Synonym{}),
					Mappings: &akov2.Mappings{
						Dynamic: nil,
						Fields:  nil,
					},
					SearchConfigurationRef: common.ResourceRefNamespaced{},
				},
				VectorSearch: nil,
			},
			AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
				Analyzer: pointer.MakePtr("analyzer"),
				Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
					{
						Name:         "test",
						TokenFilters: jsonMustEncode([]map[string]interface{}{{"token": "filter"}, {"token1": "filter1"}}),
						CharFilters:  jsonMustEncode([]map[string]interface{}{{"char": "filter"}, {"char1": "filter1"}}),
						Tokenizer: akov2.Tokenizer{
							Type:           pointer.MakePtr("type"),
							MaxGram:        nil,
							MinGram:        nil,
							Group:          nil,
							Pattern:        nil,
							MaxTokenLength: nil,
						},
					},
				}),
				SearchAnalyzer: pointer.MakePtr("searchAnalyzer"),
				StoredSource:   nil,
			},
			ID:     nil,
			Status: nil,
		}
		isEqual, err := idx1.EqualTo(idx2)
		assert.Nil(t, err)
		assert.True(t, isEqual)
	})
}

//nolint:dupl
func TestSearchIndex_ToAtlas(t *testing.T) {
	type fields struct {
		SearchIndex                akov2.SearchIndex
		AtlasSearchIndexConfigSpec akov2.AtlasSearchIndexConfigSpec
		ID                         *string
		Status                     *string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *admin.SearchIndexCreateRequest
		wantErr bool
	}{
		{
			name: "Should convert to index to Atlas using a valid index",
			fields: fields{
				SearchIndex: akov2.SearchIndex{
					Name:           "name",
					DBName:         "db",
					CollectionName: "collection",
					Type:           "search",
					Search: &akov2.Search{
						Synonyms: &([]akov2.Synonym{
							{
								Name:     "name",
								Analyzer: "analyzer",
								Source:   akov2.Source{Collection: "collection"},
							},
						}),
						Mappings: &akov2.Mappings{
							Dynamic: jsonMustEncode(true),
							Fields:  jsonMustEncodeMap(map[string]interface{}{"field": "value"}),
						},
						SearchConfigurationRef: common.ResourceRefNamespaced{},
					},
					VectorSearch: &akov2.VectorSearch{},
				},
				AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
					Analyzer: pointer.MakePtr("lucene.standard"),
					Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
						{
							Name:         "name",
							TokenFilters: jsonMustEncode([]map[string]interface{}{{"token": "filter"}}),
							CharFilters:  jsonMustEncode([]map[string]interface{}{{"char": "filter"}}),
							Tokenizer: akov2.Tokenizer{
								MaxGram:        pointer.MakePtr(20),
								MinGram:        pointer.MakePtr(10),
								Type:           pointer.MakePtr("standard"),
								Group:          pointer.MakePtr(10),
								Pattern:        pointer.MakePtr("testRegex"),
								MaxTokenLength: pointer.MakePtr(255),
							},
						},
					}),
					SearchAnalyzer: pointer.MakePtr("search-analyzer"),
					StoredSource:   jsonMustEncodeMap(map[string]interface{}{"include": "test"}),
				},
			},
			want: &admin.SearchIndexCreateRequest{
				CollectionName: "collection",
				Database:       "db",
				Name:           "name",
				Type:           pointer.MakePtr("search"),
				Definition: &admin.BaseSearchIndexCreateRequestDefinition{
					Analyzer: pointer.MakePtr("lucene.standard"),
					Analyzers: &([]admin.AtlasSearchAnalyzer{
						{
							CharFilters: jsonMustDecode(jsonMustEncode([]map[string]interface{}{
								{"char": "filter"},
							})),
							Name: "name",
							TokenFilters: jsonMustDecode(jsonMustEncode([]map[string]interface{}{
								{"token": "filter"},
							})),
							Tokenizer: admin.ApiAtlasFTSAnalyzersTokenizer{
								MaxGram:        pointer.MakePtr(20),
								MinGram:        pointer.MakePtr(10),
								Type:           pointer.MakePtr("standard"),
								Group:          pointer.MakePtr(10),
								Pattern:        pointer.MakePtr("testRegex"),
								MaxTokenLength: pointer.MakePtr(255),
							},
						},
					}),
					Mappings: &admin.SearchMappings{
						Dynamic: true,
						Fields:  &map[string]interface{}{"field": "value"},
					},
					SearchAnalyzer: pointer.MakePtr("search-analyzer"),
					Synonyms: &([]admin.SearchSynonymMappingDefinition{
						{
							Analyzer: "analyzer",
							Name:     "name",
							Source: admin.SynonymSource{
								Collection: "collection",
							},
						},
					}),
					Fields:       nil,
					StoredSource: map[string]interface{}{"include": "test"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchIndex{
				SearchIndex:                tt.fields.SearchIndex,
				AtlasSearchIndexConfigSpec: tt.fields.AtlasSearchIndexConfigSpec,
				ID:                         tt.fields.ID,
				Status:                     tt.fields.Status,
			}
			got, err := s.toAtlasCreateView()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToAtlas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ToAtlas() = %v, want %v.\nDiff: %s", got, tt.want, diff)
			}
		})
	}
}
