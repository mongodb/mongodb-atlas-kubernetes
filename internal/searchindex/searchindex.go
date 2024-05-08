// Package searchindex contains internal representation of the Atlas SearchIndex resource
package searchindex

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

// SearchIndex is the internal representation of the Atlas SearchIndex resource for the AKO usage
// SearchIndexes represented differently in AKO as CRDs and in Atlas as atlas internal structures.
// Having a separate representation allows for simpler testing of the internal logic, not tied to
// AKO and Atlas structures
type SearchIndex struct {
	akov2.SearchIndex
	akov2.AtlasSearchIndexConfigSpec
	ID     *string
	Status *string
}

func (s *SearchIndex) GetID() string {
	return pointer.GetOrDefault(s.ID, "")
}

func (s *SearchIndex) GetStatus() string {
	return pointer.GetOrDefault(s.Status, "")
}

// NewSearchIndexFromAKO requires both parts of search index: data-related part from a Deployment, and
// index configuration represented in a separate CRD. Partial construction is disabled as it won't work
// for comparing indexes between each other.
func NewSearchIndexFromAKO(index *akov2.SearchIndex, config *akov2.AtlasSearchIndexConfigSpec) *SearchIndex {
	if index == nil || config == nil {
		return nil
	}
	return &SearchIndex{
		SearchIndex:                *index.DeepCopy(),
		AtlasSearchIndexConfigSpec: *config.DeepCopy(),
	}
}

// NewSearchIndexFromAtlas returns internal representation of the SearchIndex converted from Atlas
// internals. It can return an error in case some fields are not valid JSON.
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
				Source:   akov2.Source{Collection: (*in)[i].Source.Collection},
			})
		}
		return &result
	}

	convertMappings := func(in *admin.ApiAtlasFTSMappings) (*akov2.Mappings, error) {
		if in == nil {
			return nil, nil
		}
		result := &akov2.Mappings{
			Dynamic: in.Dynamic,
			Fields:  nil,
		}
		if in.Fields == nil {
			return result, nil
		}

		var fields apiextensionsv1.JSON
		if err := compat.JSONCopy(&fields, in.Fields); err != nil {
			return nil, err
		}
		result.Fields = &fields
		return result, nil
	}

	mappings, mappingsError := convertMappings(index.Mappings)
	if mappingsError != nil {
		return nil, fmt.Errorf("unable to convert mappings: %w", mappingsError)
	}
	synonyms := convertSynonyms(index.Synonyms)

	var search *akov2.Search
	if mappings != nil || (synonyms != nil && len(*synonyms) > 0) {
		search = &akov2.Search{
			Synonyms: synonyms,
			Mappings: mappings,
		}
	}

	vectorFields, vectorError := convertVectorFields(index.Fields)
	if vectorError != nil {
		return nil, fmt.Errorf("unable to convert vector fields: %w", vectorError)
	}

	var vectorSearch *akov2.VectorSearch
	if vectorFields != nil {
		vectorSearch = &akov2.VectorSearch{
			Fields: vectorFields,
		}
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
			if err := compat.JSONCopy(&res, *in); err != nil {
				return nil, err
			}
			return &res, nil
		}

		convertTokenizer := func(in *admin.ApiAtlasFTSAnalyzersTokenizer) (akov2.Tokenizer, error) {
			res := akov2.Tokenizer{}
			if in == nil {
				return res, nil
			}

			if err := compat.JSONCopy(&res, *in); err != nil {
				return res, err
			}
			return res, nil
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
		},
		AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
			Analyzer:       index.Analyzer,
			Analyzers:      analyzers,
			SearchAnalyzer: index.SearchAnalyzer,
			StoredSource:   storedSource,
		},
		ID:     index.IndexID,
		Status: index.Status,
	}, errors.Join(errs...)
}

// cleanup normalizes the search index for comparison
func (s *SearchIndex) cleanup(cleaners ...indexCleaner) {
	s.ID = nil
	s.Status = nil

	for _, cleaner := range cleaners {
		cleaner(s)
	}
}

type indexCleaner func(s *SearchIndex)

// There is no need to compare references to configuration indexes as their data
// already represented as SearchIndexConfiguration field
func cleanConfigRef() indexCleaner {
	return func(s *SearchIndex) {
		if s.Search == nil {
			s.Search = &akov2.Search{}
		}
		s.Search.SearchConfigurationRef.Name = ""
		s.Search.SearchConfigurationRef.Namespace = ""
	}
}

// This is because API never returns nil
func cleanVectorSearch() indexCleaner {
	return func(s *SearchIndex) {
		if s.VectorSearch == nil {
			s.VectorSearch = &akov2.VectorSearch{}
		}
	}
}

// This is because API never returns nil
func cleanSynonyms() indexCleaner {
	return func(s *SearchIndex) {
		if s.Search == nil {
			return
		}
		if s.Search.Synonyms == nil {
			s.Search.Synonyms = &([]akov2.Synonym{})
		}
	}
}

// This is because API returns "null" as a valid JSON when no value is set
func cleanStoredSource() indexCleaner {
	return func(s *SearchIndex) {
		if s.StoredSource == nil {
			return
		}
		if string(s.StoredSource.Raw) == "null" {
			s.StoredSource = nil
		}
	}
}

// EqualTo compares two SearchIndexes using SemanticEqual method
func (s *SearchIndex) EqualTo(value *SearchIndex) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("caller is nil")
	}
	if s == nil {
		return false, fmt.Errorf("value is nil")
	}
	copySelf := &SearchIndex{
		SearchIndex:                *s.SearchIndex.DeepCopy(),
		AtlasSearchIndexConfigSpec: *s.AtlasSearchIndexConfigSpec.DeepCopy(),
		ID:                         nil,
		Status:                     nil,
	}
	copyValue := &SearchIndex{
		SearchIndex:                *value.SearchIndex.DeepCopy(),
		AtlasSearchIndexConfigSpec: *value.AtlasSearchIndexConfigSpec.DeepCopy(),
		ID:                         nil,
		Status:                     nil,
	}
	copySelf.cleanup(
		cleanConfigRef(),
		cleanVectorSearch(),
		cleanSynonyms(),
		cleanStoredSource(),
	)
	copyValue.cleanup(
		cleanConfigRef(),
		cleanVectorSearch(),
		cleanSynonyms(),
		cleanStoredSource(),
	)

	return cmp.SemanticEqual(copySelf, copyValue)
}

func (s *SearchIndex) Normalize() (*SearchIndex, error) {
	// TODO: Refactor interface to return error!
	err := cmp.Normalize(s)
	return s, err
}

// ToAtlas converts internal SearchIndex representation to the Atlas structure used for API calls
func (s *SearchIndex) ToAtlas() (*admin.ClusterSearchIndex, error) {
	convertJSONToListOfMaps := func(in *apiextensionsv1.JSON) (*[]map[string]interface{}, error) {
		if in == nil {
			return nil, nil
		}
		var result []map[string]interface{}
		if err := json.Unmarshal(in.Raw, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	convertJSONToMap := func(in *apiextensionsv1.JSON) (map[string]interface{}, error) {
		if in == nil {
			return nil, nil
		}
		result := map[string]interface{}{}
		if err := json.Unmarshal(in.Raw, &result); err != nil {
			return result, err
		}
		return result, nil
	}
	convertJSONToInterface := func(in *apiextensionsv1.JSON) (*[]interface{}, error) {
		if in == nil {
			return pointer.MakePtr([]interface{}{}), nil
		}
		var result []interface{}
		if err := json.Unmarshal(in.Raw, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	storedSource, err := convertJSONToMap(s.StoredSource)
	if err != nil {
		return nil, fmt.Errorf("unable to convert storedSource: %w", err)
	}

	analyzers, err := func(in *[]akov2.AtlasSearchIndexAnalyzer) (*[]admin.ApiAtlasFTSAnalyzers, error) {
		if in == nil {
			return nil, nil
		}

		result := make([]admin.ApiAtlasFTSAnalyzers, 0, len(*in))
		for i := range *in {
			analyzer := (*in)[i]
			charFilters, err := convertJSONToInterface(analyzer.CharFilters)
			if err != nil {
				return nil, err
			}

			tokenFilters, err := convertJSONToInterface(analyzer.TokenFilters)
			if err != nil {
				return nil, err
			}

			result = append(result, admin.ApiAtlasFTSAnalyzers{
				CharFilters:  charFilters,
				Name:         analyzer.Name,
				TokenFilters: tokenFilters,
				Tokenizer: admin.ApiAtlasFTSAnalyzersTokenizer{
					MaxGram:        analyzer.Tokenizer.MaxGram,
					MinGram:        analyzer.Tokenizer.MinGram,
					Type:           analyzer.Tokenizer.Type,
					Group:          analyzer.Tokenizer.Group,
					Pattern:        analyzer.Tokenizer.Pattern,
					MaxTokenLength: analyzer.Tokenizer.MaxTokenLength,
				},
			})
		}
		return &result, nil
	}(s.Analyzers)
	if err != nil {
		return nil, err
	}

	var mappings *admin.ApiAtlasFTSMappings
	convertMappings := func(in *akov2.Mappings) (*admin.ApiAtlasFTSMappings, error) {
		if in == nil {
			return nil, nil
		}
		fields, err := convertJSONToMap(in.Fields)
		if err != nil {
			return nil, err
		}
		return &admin.ApiAtlasFTSMappings{
			Dynamic: in.Dynamic,
			Fields:  fields,
		}, nil
	}
	if s.Search != nil {
		mappings, err = convertMappings(s.Search.Mappings)
		if err != nil {
			return nil, err
		}
	}

	var synonyms *[]admin.SearchSynonymMappingDefinition
	convertSynonyms := func(in *[]akov2.Synonym) *[]admin.SearchSynonymMappingDefinition {
		if in == nil {
			return nil
		}

		result := make([]admin.SearchSynonymMappingDefinition, 0, len(*in))
		for i := range *in {
			syn := &(*in)[i]

			result = append(result, admin.SearchSynonymMappingDefinition{
				Analyzer: syn.Analyzer,
				Name:     syn.Name,
				Source:   admin.SynonymSource{Collection: syn.Source.Collection},
			})
		}

		return &result
	}
	if s.Search != nil {
		synonyms = convertSynonyms(s.Search.Synonyms)
	}

	var searchFields *[]map[string]interface{}
	if s.VectorSearch != nil {
		searchFields, err = convertJSONToListOfMaps(s.VectorSearch.Fields)
		if err != nil {
			return nil, err
		}
	}

	return &admin.ClusterSearchIndex{
		CollectionName: s.CollectionName,
		Database:       s.DBName,
		IndexID:        s.ID,
		Name:           s.Name,
		Status:         s.Status,
		Type:           pointer.MakePtr(s.SearchIndex.Type),
		Analyzer:       s.Analyzer,
		Analyzers:      analyzers,
		Mappings:       mappings,
		SearchAnalyzer: s.SearchAnalyzer,
		StoredSource:   storedSource,
		Synonyms:       synonyms,
		Fields:         searchFields,
	}, nil
}
