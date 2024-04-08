package searchindex

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func jsonMustEncode(jsn []map[string]interface{}) *apiextensionsv1.JSON {
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
								Source:   &akov2.Source{Collection: "test-collection"},
							},
						}),
						Mapping: &akov2.Mapping{
							Dynamic: pointer.MakePtr(true),
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
					IndexConfigRef: common.ResourceRefNamespaced{},
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
							Tokenizer: &akov2.Tokenizer{
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
								Source:   &akov2.Source{Collection: "test-collection"},
							},
						}),
						Mapping: &akov2.Mapping{
							Dynamic: pointer.MakePtr(true),
							Fields: jsonMustEncode([]map[string]interface{}{
								{"test": "value"},
							}),
						},
						SearchConfigurationRef: common.ResourceRefNamespaced{},
					},
					VectorSearch: &akov2.VectorSearch{Fields: jsonMustEncode([]map[string]interface{}{
						{"test": "value"},
					})},
					IndexConfigRef: common.ResourceRefNamespaced{},
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
							Tokenizer: &akov2.Tokenizer{
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

			got := NewSearchIndexFromAKO(tt.args.index, tt.args.config)
			if diff := cmp.Diff(got, tt.want, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("NewSearchIndexFromAKO() = %v, want %v.\nDiff: %s", got, tt.want, diff)
			}
		})
	}
}

func Test_NewSearchIndexFromAtlas(t *testing.T) {
	type args struct {
		index admin.ClusterSearchIndex
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
				index: admin.ClusterSearchIndex{
					CollectionName: "collection",
					Database:       "db",
					IndexID:        nil,
					Name:           "name",
					Status:         nil,
					Type:           pointer.MakePtr("search"),
					Analyzer:       pointer.MakePtr("lucene.standard"),
					Analyzers: &([]admin.ApiAtlasFTSAnalyzers{
						{
							CharFilters: jsonMustDecode(jsonMustEncode([]map[string]interface{}{
								{"char": "filter"},
							})),
							Name: "name",
							TokenFilters: jsonMustDecode(jsonMustEncode([]map[string]interface{}{
								{"token": "filter"},
							})),
							Tokenizer: admin.ApiAtlasFTSAnalyzersTokenizer{
								MaxGram:        nil,
								MinGram:        nil,
								Type:           pointer.MakePtr("standard"),
								Group:          nil,
								Pattern:        nil,
								MaxTokenLength: pointer.MakePtr(255),
							},
						},
					}),
					Mappings: &admin.ApiAtlasFTSMappings{
						Dynamic: pointer.MakePtr(true),
						Fields:  map[string]interface{}{"field": "value"},
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
			want: &SearchIndex{
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
								Source:   &akov2.Source{Collection: "collection"},
							},
						}),
						Mapping: &akov2.Mapping{
							Dynamic: pointer.MakePtr(true),
							Fields:  jsonMustEncodeMap(map[string]interface{}{"field": "value"}),
						},
						SearchConfigurationRef: common.ResourceRefNamespaced{},
					},
					VectorSearch:   &akov2.VectorSearch{},
					IndexConfigRef: common.ResourceRefNamespaced{},
				},
				AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
					Analyzer: pointer.MakePtr("lucene.standard"),
					Analyzers: &([]akov2.AtlasSearchIndexAnalyzer{
						{
							Name:         "name",
							TokenFilters: jsonMustEncode([]map[string]interface{}{{"token": "filter"}}),
							CharFilters:  jsonMustEncode([]map[string]interface{}{{"char": "filter"}}),
							Tokenizer: &akov2.Tokenizer{
								Type:           pointer.MakePtr("standard"),
								MaxGram:        nil,
								MinGram:        nil,
								Group:          nil,
								Pattern:        nil,
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
			got, err := NewSearchIndexFromAtlas(tt.args.index)
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

func TestSearchIndex_EqualTo(t *testing.T) {
	type fields struct {
		SearchIndex                v1.SearchIndex
		AtlasSearchIndexConfigSpec v1.AtlasSearchIndexConfigSpec
	}
	type args struct {
		value *SearchIndex
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchIndex{
				SearchIndex:                tt.fields.SearchIndex,
				AtlasSearchIndexConfigSpec: tt.fields.AtlasSearchIndexConfigSpec,
			}
			if got := s.EqualTo(tt.args.value); got != tt.want {
				t.Errorf("EqualTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestSearchIndex_ToAtlas(t *testing.T) {
//	type fields struct {
//		SearchIndex                v1.SearchIndex
//		AtlasSearchIndexConfigSpec v1.AtlasSearchIndexConfigSpec
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   *admin.ClusterSearchIndex
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &SearchIndex{
//				SearchIndex:                tt.fields.SearchIndex,
//				AtlasSearchIndexConfigSpec: tt.fields.AtlasSearchIndexConfigSpec,
//			}
//			if got := s.ToAtlas(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ToAtlas() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
