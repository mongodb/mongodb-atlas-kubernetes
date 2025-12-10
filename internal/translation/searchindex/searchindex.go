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

// Package searchindex contains internal representation of the Atlas SearchIndex resource
package searchindex

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
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

// NewSearchIndex requires both parts of search index: data-related part from a Deployment, and
// index configuration represented in a separate CRD. Partial construction is disabled as it won't work
// for comparing indexes between each other.
func NewSearchIndex(index *akov2.SearchIndex, config *akov2.AtlasSearchIndexConfigSpec) *SearchIndex {
	if index == nil || config == nil {
		return nil
	}
	return &SearchIndex{
		SearchIndex:                *index.DeepCopy(),
		AtlasSearchIndexConfigSpec: *config.DeepCopy(),
	}
}

// fromAtlas returns internal representation of the SearchIndex converted from Atlas
// internals. It can return an error in case some fields are not valid JSON.
func fromAtlas(index admin.SearchIndexResponse) (*SearchIndex, error) {
	mappings, mappingsError := mappingsFromAtlas(index.LatestDefinition.Mappings)
	if mappingsError != nil {
		return nil, fmt.Errorf("unable to convert mappings: %w", mappingsError)
	}
	synonyms := synonymsFromAtlas(index.LatestDefinition.Synonyms)

	var search *akov2.Search
	if mappings != nil || (synonyms != nil && len(*synonyms) > 0) {
		search = &akov2.Search{
			Synonyms: synonyms,
			Mappings: mappings,
		}
	}

	vectorFields, vectorError := convertVectorFields(index.LatestDefinition.Fields)
	if vectorError != nil {
		return nil, fmt.Errorf("unable to convert vector fields: %w", vectorError)
	}

	var vectorSearch *akov2.VectorSearch
	if vectorFields != nil {
		vectorSearch = &akov2.VectorSearch{
			Fields: vectorFields,
		}
	}

	var errs []error
	analyzers, err := analyzersFromAtlas(index.LatestDefinition.Analyzers)
	if err != nil {
		errs = append(errs, err)
	}

	storedSource, err := convertStoredSource(index.LatestDefinition.StoredSource)
	if err != nil {
		errs = append(errs, err)
	}

	return &SearchIndex{
		SearchIndex: akov2.SearchIndex{
			Name:           pointer.GetOrDefault(index.Name, ""),
			DBName:         pointer.GetOrDefault(index.Database, ""),
			CollectionName: pointer.GetOrDefault(index.CollectionName, ""),
			Type:           pointer.GetOrDefault(index.Type, ""),
			Search:         search,
			VectorSearch:   vectorSearch,
		},
		AtlasSearchIndexConfigSpec: akov2.AtlasSearchIndexConfigSpec{
			Analyzer:       index.LatestDefinition.Analyzer,
			Analyzers:      analyzers,
			SearchAnalyzer: index.LatestDefinition.SearchAnalyzer,
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

// toAtlas converts internal SearchIndex representation to the Atlas structure used for API calls
func (s *SearchIndex) toAtlasCreateView() (*admin.SearchIndexCreateRequest, error) {
	storedSource, err := jsonToMap(s.StoredSource)
	if err != nil {
		return nil, fmt.Errorf("unable to convert storedSource: %w", err)
	}

	analyzers, err := analyzersToAtlas(s.Analyzers)
	if err != nil {
		return nil, err
	}

	var mappings *admin.SearchMappings
	if s.Search != nil {
		mappings, err = mappingsToAtlas(s.Search.Mappings)
		if err != nil {
			return nil, err
		}
	}

	var synonyms *[]admin.SearchSynonymMappingDefinition
	if s.Search != nil {
		synonyms = synonymsToAtlas(s.Search.Synonyms)
	}

	var searchFields *[]any
	if s.VectorSearch != nil {
		searchFields, err = jsonToArrayOfAny(s.VectorSearch.Fields)
		if err != nil {
			return nil, err
		}
	}

	result := &admin.SearchIndexCreateRequest{
		CollectionName: s.CollectionName,
		Database:       s.DBName,
		Name:           s.Name,
		Type:           pointer.MakePtr(s.SearchIndex.Type),
		Definition: &admin.BaseSearchIndexCreateRequestDefinition{
			Analyzer:       s.Analyzer,
			Analyzers:      analyzers,
			Mappings:       mappings,
			SearchAnalyzer: s.SearchAnalyzer,
			Synonyms:       synonyms,
			Fields:         searchFields,
		},
	}

	// This is a workaround because of JSON marshaller marshals nil (type map[string]interface{}) to null in JSON
	// which is not accepted by the API endpoint
	if len(storedSource) > 0 {
		result.Definition.StoredSource = storedSource
	}

	return result, nil
}

func (s *SearchIndex) toAtlasUpdateView() (*admin.SearchIndexUpdateRequest, error) {
	storedSource, err := jsonToMap(s.StoredSource)
	if err != nil {
		return nil, fmt.Errorf("unable to convert storedSource: %w", err)
	}

	analyzers, err := analyzersToAtlas(s.Analyzers)
	if err != nil {
		return nil, err
	}

	var mappings *admin.SearchMappings
	if s.Search != nil {
		mappings, err = mappingsToAtlas(s.Search.Mappings)
		if err != nil {
			return nil, err
		}
	}

	var synonyms *[]admin.SearchSynonymMappingDefinition
	if s.Search != nil {
		synonyms = synonymToAtlass(s.Search.Synonyms)
	}

	var searchFields *[]any
	if s.VectorSearch != nil {
		searchFields, err = jsonToArrayOfAny(s.VectorSearch.Fields)
		if err != nil {
			return nil, err
		}
	}

	return &admin.SearchIndexUpdateRequest{
		Definition: admin.SearchIndexUpdateRequestDefinition{
			Analyzer:       s.Analyzer,
			Analyzers:      analyzers,
			Mappings:       mappings,
			SearchAnalyzer: s.SearchAnalyzer,
			StoredSource:   storedSource,
			Synonyms:       synonyms,
			Fields:         searchFields,
		},
	}, nil
}

func convertVectorFields(in *[]any) (*apiextensionsv1.JSON, error) {
	if in == nil {
		return nil, nil
	}
	result := &apiextensionsv1.JSON{}
	err := compat.JSONCopy(result, *in)
	return result, err
}

func synonymsFromAtlas(in *[]admin.SearchSynonymMappingDefinition) *[]akov2.Synonym {
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

func synonymToAtlass(in *[]akov2.Synonym) *[]admin.SearchSynonymMappingDefinition {
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

func mappingsFromAtlas(in *admin.SearchMappings) (*akov2.Mappings, error) {
	if in == nil {
		return nil, nil
	}
	result := &akov2.Mappings{}
	var dynamic, fields apiextensionsv1.JSON
	if err := compat.JSONCopy(&fields, in.Fields); err != nil {
		return nil, err
	}
	if err := compat.JSONCopy(&dynamic, in.Dynamic); err != nil {
		return nil, err
	}
	if len(fields.Raw) > 0 {
		result.Fields = &fields
	}
	if len(dynamic.Raw) > 0 {
		result.Dynamic = &dynamic
	}
	return result, nil
}

func jsonToAny(in *apiextensionsv1.JSON) (any, error) {
	if in == nil {
		return nil, nil
	}
	var result any
	if err := json.Unmarshal(in.Raw, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func jsonToArrayOfAny(in *apiextensionsv1.JSON) (*[]any, error) {
	if in == nil {
		return nil, nil
	}
	var result []any
	if err := json.Unmarshal(in.Raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func jsonToMap(in *apiextensionsv1.JSON) (map[string]any, error) {
	if in == nil {
		return nil, nil
	}
	result := map[string]any{}
	if err := json.Unmarshal(in.Raw, &result); err != nil {
		return result, err
	}
	return result, nil
}

func jsonToInterface(in *apiextensionsv1.JSON) (*[]interface{}, error) {
	if in == nil {
		return pointer.MakePtr([]interface{}{}), nil
	}
	var result []interface{}
	if err := json.Unmarshal(in.Raw, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func analyzersToAtlas(in *[]akov2.AtlasSearchIndexAnalyzer) (*[]admin.AtlasSearchAnalyzer, error) {
	if in == nil {
		return nil, nil
	}

	result := make([]admin.AtlasSearchAnalyzer, 0, len(*in))
	for i := range *in {
		analyzer := (*in)[i]
		charFilters, err := jsonToInterface(analyzer.CharFilters)
		if err != nil {
			return nil, err
		}

		tokenFilters, err := jsonToInterface(analyzer.TokenFilters)
		if err != nil {
			return nil, err
		}

		result = append(result, admin.AtlasSearchAnalyzer{
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
}

func analyzersFromAtlas(in *[]admin.AtlasSearchAnalyzer) (*[]akov2.AtlasSearchIndexAnalyzer, error) {
	if in == nil {
		return nil, nil
	}
	result := make([]akov2.AtlasSearchIndexAnalyzer, 0, len(*in))

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

		tokenizer, err := convertTokenizer((*in)[i].Tokenizer)
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

func mappingsToAtlas(in *akov2.Mappings) (*admin.SearchMappings, error) {
	if in == nil {
		return nil, nil
	}
	fields, err := jsonToMap(in.Fields)
	if err != nil {
		return nil, err
	}
	dynamic, err := jsonToAny(in.Dynamic)
	if err != nil {
		return nil, err
	}
	return &admin.SearchMappings{
		Dynamic: dynamic,
		Fields:  &fields,
	}, nil
}

func synonymsToAtlas(in *[]akov2.Synonym) *[]admin.SearchSynonymMappingDefinition {
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

func convertStoredSource(in any) (*apiextensionsv1.JSON, error) {
	val, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return &apiextensionsv1.JSON{Raw: val}, nil
}

func convertFilters(in *[]interface{}) (*apiextensionsv1.JSON, error) {
	if in == nil {
		return nil, nil
	}
	var res apiextensionsv1.JSON
	if err := compat.JSONCopy(&res, *in); err != nil {
		return nil, err
	}
	return &res, nil
}

func convertTokenizer(in any) (akov2.Tokenizer, error) {
	res := akov2.Tokenizer{}
	if in == nil {
		return res, nil
	}

	if err := compat.JSONCopy(&res, in); err != nil {
		return res, err
	}
	return res, nil
}
