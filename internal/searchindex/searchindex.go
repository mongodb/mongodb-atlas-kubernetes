package searchindex

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

type SearchIndex struct {
	akov2.SearchIndex
	akov2.AtlasSearchIndexConfigSpec
}

func NewSearchIndexFromAKO(index *akov2.SearchIndex, config *akov2.AtlasSearchIndexConfigSpec) *SearchIndex {
	if index == nil || config == nil {
		return nil
	}
	return &SearchIndex{
		SearchIndex:                *index.DeepCopy(),
		AtlasSearchIndexConfigSpec: *config.DeepCopy(),
	}
}

func NewSearchIndexFromAtlas(index admin.ClusterSearchIndex) (*SearchIndex, error) {
	convertVectorFields := func(in *[]map[string]interface{}) (*apiextensionsv1.JSON, error) {
		if in == nil {
			return nil, nil
		}
		result := &apiextensionsv1.JSON{}
		err := compat.JSONCopy(result, *in)
		return result, err
	}

	convertSynonyms := func(in *[]admin.SearchSynonymMappingDefinition) *[]akov2.Synonym {
		if in == nil {
			return nil
		}

		result := make([]akov2.Synonym, 0, len(*in))

		for i := range *in {
			result = append(result, akov2.Synonym{
				Name:     (*in)[i].Name,
				Analyzer: (*in)[i].Analyzer,
				Source:   &akov2.Source{Collection: (*in)[i].Source.Collection},
			})
		}
		return &result
	}

	convertMappings := func(in *admin.ApiAtlasFTSMappings) (*akov2.Mapping, error) {
		if in == nil {
			return nil, nil
		}
		result := &akov2.Mapping{
			Dynamic: in.Dynamic,
			Fields:  nil,
		}
		if in.Fields == nil {
			return result, nil
		}

		var fields apiextensionsv1.JSON
		err := compat.JSONCopy(&fields, in.Fields)
		result.Fields = &fields
		return result, err
	}

	mappings, mappingsError := convertMappings(index.Mappings)
	if mappingsError != nil {
		return nil, fmt.Errorf("unable to convert mappings: %w", mappingsError)
	}

	search := &akov2.Search{
		Synonyms:               convertSynonyms(index.Synonyms),
		Mapping:                mappings,
		SearchConfigurationRef: common.ResourceRefNamespaced{},
	}

	vectorFields, vectorError := convertVectorFields(index.Fields)
	vectorSearch := &akov2.VectorSearch{
		Fields: vectorFields,
	}
	if mappingsError != nil {
		return nil, fmt.Errorf("unable to convert vector fields: %w", vectorError)
	}

	convertAnalyzers := func(in *[]admin.ApiAtlasFTSAnalyzers) (*[]akov2.AtlasSearchIndexAnalyzer, error) {
		if in == nil {
			return nil, nil
		}
		result := make([]akov2.AtlasSearchIndexAnalyzer, 0, len(*in))

		convertFilters := func(in *[]interface{}) (*apiextensionsv1.JSON, error) {
			if in == nil {
				return nil, nil
			}
			var res apiextensionsv1.JSON
			err := compat.JSONCopy(&res, *in)
			return &res, err
		}

		convertTokenizer := func(in *admin.ApiAtlasFTSAnalyzersTokenizer) (*akov2.Tokenizer, error) {
			if in == nil {
				return nil, nil
			}

			res := &akov2.Tokenizer{}
			err := compat.JSONCopy(res, *in)
			return res, err
		}

		errs := []error{}
		for i := range *in {
			tokenFilters, err := convertFilters((*in)[i].TokenFilters)
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to convert tokenFilters: %w", err))
				continue
			}

			charFilters, err := convertFilters((*in)[i].CharFilters)
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to convert charFilters: %w", err))
				continue
			}

			tokenizer, err := convertTokenizer(&(*in)[i].Tokenizer)
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to convert tokenizer: %w", err))
				continue
			}

			result = append(result, akov2.AtlasSearchIndexAnalyzer{
				Name:         (*in)[i].Name,
				TokenFilters: tokenFilters,
				CharFilters:  charFilters,
				Tokenizer:    tokenizer,
			})
		}
		e := errors.Join(errs...)
		return &result, e
	}

	convertStoredSource := func(in map[string]interface{}) (*apiextensionsv1.JSON, error) {
		val, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		return &apiextensionsv1.JSON{Raw: val}, nil
	}

	var errs []error
	analyzers, err := convertAnalyzers(index.Analyzers)
	if err != nil {
		errs = append(errs, err)
	}

	storedSource, err := convertStoredSource(index.StoredSource)
	if err != nil {
		errs = append(errs, err)
	}

	return &SearchIndex{
		SearchIndex: akov2.SearchIndex{
			Name:           index.Name,
			DBName:         index.Database,
			CollectionName: index.CollectionName,
			Type:           pointer.GetOrDefault(index.Type, ""),
			Search:         search,
			VectorSearch:   vectorSearch,
			IndexConfigRef: common.ResourceRefNamespaced{},
		},
		AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
			Analyzer:       index.Analyzer,
			Analyzers:      analyzers,
			SearchAnalyzer: index.SearchAnalyzer,
			StoredSource:   storedSource,
		},
	}, errors.Join(errs...)
}

func (s *SearchIndex) EqualTo(value *SearchIndex) bool {
	return value != nil
}

//func (s *SearchIndex) ToAtlas() *admin.ClusterSearchIndex {
//	convertTokenizerType := func(t *akov2.Tokenizer) admin.ApiAtlasFTSAnalyzersTokenizer {
//		return admin.ApiAtlasFTSAnalyzersTokenizer{
//			MaxGram:        t.MaxGram,
//			MinGram:        t.MinGram,
//			Type:           t.Type,
//			Group:          t.Group,
//			Pattern:        t.Pattern,
//			MaxTokenLength: t.MaxTokenLength,
//		}
//	}
//
//	convertCharFilters := func(in []akov2.CharFilter) *[]interface{} {
//		result := make([]interface{}, 0, len(in))
//
//		for i := range in {
//			resultItem, _ := json.Marshal((in)[i])
//			result = append(result, resultItem)
//		}
//		return &result
//	}
//
//	convertTokenFilters := func(in []akov2.TokenFilter) *[]interface{} {
//		result := make([]interface{}, 0, len(in))
//		for i := range in {
//			resultItem, _ := json.Marshal((in)[i])
//			result = append(result, resultItem)
//		}
//		return &result
//	}
//
//	convertAnalyzers := func(in *[]akov2.AtlasSearchIndexAnalyzer) *[]admin.ApiAtlasFTSAnalyzers {
//		if in == nil {
//			return nil
//		}
//		result := make([]admin.ApiAtlasFTSAnalyzers, 0, len(*in))
//
//		for i := range *in {
//			result = append(result, admin.ApiAtlasFTSAnalyzers{
//				CharFilters:  convertCharFilters((*in)[i].CharFilters),
//				Name:         (*in)[i].Name,
//				TokenFilters: convertTokenFilters((*in)[i].TokenFilters),
//				Tokenizer:    convertTokenizerType((*in)[i].Tokenizer),
//			})
//		}
//		return &result
//	}
//
//	convertMappings := func(in *akov2.Mapping) *admin.ApiAtlasFTSMappings {
//		resultMap := map[string]interface{}{}
//		for k, v := range in.Fields {
//			resultMap[k] = v
//		}
//		return &admin.ApiAtlasFTSMappings{
//			Dynamic: in.Dynamic,
//			Fields:  resultMap,
//		}
//	}
//
//	convertSynonyms := func(in *[]akov2.Synonym) *[]admin.SearchSynonymMappingDefinition {
//		if in == nil {
//			return nil
//		}
//		result := make([]admin.SearchSynonymMappingDefinition, 0, len(*in))
//
//		for i := range *in {
//			result = append(result, admin.SearchSynonymMappingDefinition{
//				Analyzer: (*in)[i].Analyzer,
//				Name:     (*in)[i].Name,
//				Source: admin.SynonymSource{
//					Collection: (*in)[i].Source.Collection,
//				},
//			})
//		}
//		return &result
//	}
//
//	convertFields := func(in []akov2.VectorSearchField) *[]map[string]interface{} {
//		result := make([]map[string]interface{}, 0, len(in))
//
//		for i := range in {
//			result = append(result, map[string]interface{}{
//				"type":       in[i].Type,
//				"path":       in[i].Path,
//				"dimensions": in[i].Dimensions,
//				"similarity": in[i].Similarity,
//			})
//		}
//		return &result
//	}
//
//	return &admin.ClusterSearchIndex{
//		CollectionName: s.CollectionName,
//		Database:       s.DBName,
//		// Should be nil
//		IndexID:        nil,
//		Name:           s.Name,
//		Status:         nil,
//		Type:           &s.Type,
//		Analyzer:       s.Analyzer,
//		Analyzers:      convertAnalyzers(s.Analyzers),
//		Mappings:       convertMappings(s.Search.Mapping),
//		SearchAnalyzer: s.SearchAnalyzer,
//		Synonyms:       convertSynonyms(s.Search.Synonyms),
//		Fields:         convertFields(s.VectorSearch.Fields),
//	}
//}
