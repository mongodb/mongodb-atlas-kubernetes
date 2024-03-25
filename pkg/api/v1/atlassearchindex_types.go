package v1

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

type SearchIndex struct {
	// Human-readable label that identifies this index. Must be unique for a deployment
	// +required
	Name string `json:"name"`
	// Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes
	// +required
	DBName string `json:"DBName"`
	// Human-readable label that identifies the collection that contains one or more Atlas Search indexes
	// +required
	CollectionName string `json:"collectionName"`
	// Type of the index. Default type is search
	// +kubebuilder:validation:Enum:=search;vectorSearch
	Type string `json:"type,omitempty"`
	// +optional
	// Atlas search index configuration
	Search *Search `json:"search,omitempty"`
	// +optional
	// Atlas vector search index configuration
	VectorSearch *VectorSearch `json:"vectorSearch,omitempty"`
	// +required
	// A name and a namespace of the AtlasSearchIndexConfig resource
	IndexConfigRef common.ResourceRefNamespaced `json:"indexConfigRef,omitempty"`
}

//// SearchIndexDTO is intended for internal usage only. This is a temporary object used to compare
//// AKO configured indices and Atlas indices
//// +k8s:deepcopy-gen=false
//type SearchIndexDTO struct {
//	SearchIndex
//	AtlasSearchIndexConfigSpec
//}
//
////type SearchIndexDTO *admin.ClusterSearchIndex
//
//func NewSearchIndexDTOFromAKO(index *SearchIndex, config *AtlasSearchIndexConfigSpec) *SearchIndexDTO {
//	return &SearchIndexDTO{
//		SearchIndex:                *index,
//		AtlasSearchIndexConfigSpec: *config,
//	}
//}
//
//func NewSearchIndexDTOFromAtlas(index *admin.ClusterSearchIndex) *SearchIndexDTO {
//	convertSynonyms := func(in *[]admin.SearchSynonymMappingDefinition) []Synonym {
//
//	}
//	convertMappings := func(in *admin.ApiAtlasFTSMappings) *Mapping {
//
//	}
//	convertAnalyzers := func(in *[]admin.ApiAtlasFTSAnalyzers) []AtlasSearchIndexAnalyzer {
//
//	}
//	return &SearchIndexDTO{
//		Name:           index.Name,
//		DBName:         index.Database,
//		CollectionName: index.CollectionName,
//		Type:           *index.Type,
//		Synonyms:       index.Synonyms, // convert
//		Mapping:        index.Mappings, // convert
//		Fields:         index.Fields,   // convert
//		AtlasSearchIndexConfigSpec: AtlasSearchIndexConfigSpec{
//			Analyzer:       index.Analyzer,
//			Analyzers:      index.Analyzers,
//			SearchAnalyzer: index.SearchAnalyzer,
//			StoredSource:   index.StoredSource,
//		},
//	}
//}
//
//func (si *SearchIndex) ToAtlas(ic *AtlasSearchIndexConfigSpec) *admin.ClusterSearchIndex {
//	convertTokenizerType := func(t *Tokenizer) admin.ApiAtlasFTSAnalyzersTokenizer {
//		result := admin.ApiAtlasFTSAnalyzersTokenizer{}
//
//		switch {
//		case t.EdgeGram != nil:
//			result.Type = pointer.MakePtr("edgeGram")
//			result.MaxGram = &t.EdgeGram.MaxGram
//			result.MinGram = &t.EdgeGram.MinGram
//		case t.NGram != nil:
//			result.Type = pointer.MakePtr("nGram")
//			result.MaxGram = &t.EdgeGram.MaxGram
//			result.MinGram = &t.EdgeGram.MinGram
//		case t.Keyword != nil:
//			result.Type = pointer.MakePtr("keyword")
//		case t.RegexCaptureGroup != nil:
//			result.Type = pointer.MakePtr("regexCaptureGroup")
//			result.Pattern = pointer.MakePtr(t.RegexCaptureGroup.Pattern)
//			result.Group = pointer.MakePtr(t.RegexCaptureGroup.Group)
//		case t.RegexSplit != nil:
//			result.Type = pointer.MakePtr("regexSplit")
//			result.Pattern = pointer.MakePtr(t.RegexSplit.Pattern)
//		case t.Standard != nil:
//			result.Type = pointer.MakePtr("standard")
//			result.MaxTokenLength = pointer.MakePtr(t.Standard.MaxTokenLength)
//		case t.UaxUrlEmail != nil:
//			result.Type = pointer.MakePtr("uaxUrlEmail")
//			result.MaxTokenLength = pointer.MakePtr(t.Standard.MaxTokenLength)
//		case t.Whitespace != nil:
//			result.Type = pointer.MakePtr("whitespace")
//			result.MaxTokenLength = pointer.MakePtr(t.Standard.MaxTokenLength)
//		}
//		return result
//	}
//
//	convertCharFilters := func(in []CharFilter) *[]interface{} {
//		result := make([]interface{}, 0, len(in))
//		for i := range in {
//			resultItem := map[string]interface{}{}
//
//			switch {
//			case in[i].HtmlNormalize != nil:
//				resultItem["ignoreTags"] = in[i].HtmlNormalize.IgnoreTags
//				resultItem["type"] = "htmlNormalize"
//			case in[i].IcuNormalize != nil:
//				resultItem["icuNormalize"] = in[i].IcuNormalize
//				resultItem["type"] = "icuNormalize"
//			case in[i].Mapping != nil:
//				resultItem["mappings"] = in[i].Mapping.Mappings
//				resultItem["type"] = "mapping"
//			case in[i].Persian != nil:
//				resultItem["persian"] = in[i].Persian
//				resultItem["type"] = "persian"
//			}
//			result = append(result, resultItem)
//		}
//		return &result
//	}
//
//	convertTokenFilters := func(in []TokenFilter) *[]interface{} {
//		result := make([]interface{}, 0, len(in))
//		for i := range in {
//			item := map[string]interface{}{}
//
//			switch {
//			case in[i].AsciiFolding != nil:
//				item["type"] = "asciiFolding"
//				item["originalTokens"] = in[i].AsciiFolding.OriginalTokens
//			case in[i].DaitchMokotoffSoundex != nil:
//				item["type"] = "daitchMokotoffSoundex"
//				item["originalTokens"] = in[i].DaitchMokotoffSoundex.OriginalTokens
//			case in[i].EdgeGram != nil:
//				item["type"] = "edgeGram"
//				item["maxGram"] = in[i].EdgeGram.MaxGram
//				item["minGram"] = in[i].EdgeGram.MinGram
//				item["termNotInBounds"] = in[i].EdgeGram.TermNotInBounds
//			case in[i].IcuFolding != nil:
//				item["type"] = "icuFolding"
//			case in[i].IcuNormalizer != nil:
//				item["type"] = "icuNormalizer"
//				item["normalizationForm"] = in[i].IcuNormalizer.NormalizationForm
//			case in[i].Length != nil:
//				item["type"] = "length"
//				item["max"] = in[i].Length.Max
//				item["max"] = in[i].Length.Min
//			case in[i].Lowercase != nil:
//				item["type"] = "lowercase"
//			case in[i].NGram != nil:
//				item["type"] = "nGram"
//				item["maxGram"] = in[i].NGram.MaxGram
//				item["minGram"] = in[i].NGram.MinGram
//				item["termNotInBounds"] = in[i].NGram.TermNotInBounds
//			case in[i].Regex != nil:
//				item["type"] = "regex"
//				item["matches"] = in[i].Regex.Matches
//				item["pattern"] = in[i].Regex.Pattern
//				item["replacement"] = in[i].Regex.Replacement
//			case in[i].Reverse != nil:
//				item["type"] = "reverse"
//			case in[i].Shingle != nil:
//				item["type"] = "shingle"
//				item["maxShingleSize"] = in[i].Shingle.MaxShingleSize
//				item["minShingleSize"] = in[i].Shingle.MinShingleSize
//			case in[i].SnowballStemming != nil:
//				item["type"] = "snowballStemming"
//				item["stemmerName"] = in[i].SnowballStemming.StemmerName
//			case in[i].Stopword != nil:
//				item["type"] = "stopword"
//				item["ignoreCase"] = in[i].Stopword.IgnoreCase
//				item["tokens"] = in[i].Stopword.Tokens
//			case in[i].Trim != nil:
//				item["type"] = "trim"
//			}
//
//			result = append(result, item)
//		}
//		return &result
//	}
//
//	convertAnalyzers := func(in []AtlasSearchIndexAnalyzer) *[]admin.ApiAtlasFTSAnalyzers {
//		result := make([]admin.ApiAtlasFTSAnalyzers, 0, len(in))
//
//		for i := range in {
//			result = append(result, admin.ApiAtlasFTSAnalyzers{
//				CharFilters:  convertCharFilters(in[i].CharFilters),
//				Name:         in[i].Name,
//				TokenFilters: convertTokenFilters(in[i].TokenFilters),
//				Tokenizer:    convertTokenizerType(in[i].Tokenizer),
//			})
//		}
//		return &result
//	}
//
//	convertMappings := func(in *Mapping) *admin.ApiAtlasFTSMappings {
//		return &admin.ApiAtlasFTSMappings{
//			Dynamic: &in.Dynamic,
//			//Fields:  map[string]interface{}(in.Fields),
//		}
//	}
//
//	convertSynonyms := func(in []Synonym) *[]admin.SearchSynonymMappingDefinition {
//		result := make([]admin.SearchSynonymMappingDefinition, 0, len(in))
//
//		for i := range in {
//			result = append(result, admin.SearchSynonymMappingDefinition{
//				Analyzer: in[i].Analyzer,
//				Name:     in[i].Name,
//				Source: admin.SynonymSource{
//					Collection: in[i].Source.Collection,
//				},
//			})
//		}
//		return &result
//	}
//
//	return &admin.ClusterSearchIndex{
//		CollectionName: si.CollectionName,
//		Database:       si.DBName,
//		// Should be nil
//		IndexID:        nil,
//		Name:           si.Name,
//		Status:         nil,
//		Type:           &si.Type,
//		Analyzer:       ic.Analyzer,
//		Analyzers:      convertAnalyzers(ic.Analyzers),
//		Mappings:       convertMappings(si.Search.Mapping),
//		SearchAnalyzer: ic.SearchAnalyzer,
//		Synonyms:       convertSynonyms(si.Search.Synonyms),
//		// TODO: Not described in the docs
//		Fields: nil,
//	}
//}

type Search struct {
	// Rule sets that map words to their synonyms in this index
	// +optional
	Synonyms *[]Synonym `json:"synonyms,omitempty"`
	// Index specifications for the collection's fields
	// +optional
	Mapping *Mapping `json:"mappings,omitempty"`
	// +required
	// A reference to the AtlasSearchIndexConfig custom resource
	SearchConfigurationRef common.ResourceRefNamespaced `json:"searchConfigurationRef"`
}

type Synonym struct {
	// Human-readable label that identifies the synonym definition. Each name must be unique within the same index definition
	// +required
	Name string `json:"name"`
	// Specific pre-defined method chosen to apply to the synonyms to be searched
	// +kubebuilder:validation:Enum:=lucene.standard;lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword;lucene.arabic;lucene.armenian;lucene.basque;lucene.bengali;lucene.brazilian;lucene.bulgarian;lucene.catalan;lucene.chinese;lucene.cjk;lucene.czech;lucene.danish;lucene.dutch;lucene.english;lucene.finnish;lucene.french;lucene.galician;lucene.german;lucene.greek;lucene.hindi;lucene.hungarian;lucene.indonesian;lucene.irish;lucene.italian;lucene.japanese;lucene.korean;lucene.kuromoji;lucene.latvian;lucene.lithuanian;lucene.morfologik;lucene.nori;lucene.norwegian;lucene.persian;lucene.portuguese;lucene.romanian;lucene.russian;lucene.smartcn;lucene.sorani;lucene.spanish;lucene.swedish;lucene.thai;lucene.turkish;lucene.ukrainian
	// +required
	Analyzer string `json:"analyzer"`
	// Data set that stores the mapping one or more words map to one or more synonyms of those words
	Source *Source `json:"source"`
}

type Source struct {
	// Human-readable label that identifies the MongoDB collection that stores words and their applicable synonyms
	Collection string `json:"collection"`
}

type Mapping struct {
	// Flag that indicates whether the index uses dynamic or static mappings. Required if mapping.fields is omitted.
	Dynamic *bool `json:"dynamic,omitempty"`
	// One or more field specifications for the Atlas Search index. Required if mapping.dynamic is omitted or set to false.
	Fields *apiextensions.JSON `json:"fields,omitempty"`
}

type VectorSearch struct {
	// Array of JSON objects. See examples https://dochub.mongodb.org/core/avs-vector-type
	// +required
	//Fields []VectorSearchField `json:"fields,omitempty"`
	Fields *apiextensions.JSON `json:"fields,omitempty"`
}

type VectorSearchField struct {
	// +kubebuilder:validations:Enum:=vector;filter
	// +kubebuilder:default:=vector
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	Path string `json:"path,omitempty"`
	// +optional
	Dimensions int `json:"dimensions,omitempty"`
	// +optional
	Similarity string `json:"similarity,omitempty"`
}
