package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasSearchIndexConfig{})
}

// AtlasSearchIndexConfig is the Schema for the AtlasSearchIndexConfig API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasSearchIndexConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasSearchIndexConfigSpec    `json:"spec,omitempty"`
	Status status.AtlasSearchIndexStatus `json:"status,omitempty"`
}

// AtlasSearchIndexConfigList contains a list of AtlasSearchIndexConfig
// +kubebuilder:object:root=true
type AtlasSearchIndexConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasSearchIndexConfig `json:"items"`
}

type AtlasSearchIndexConfigSpec struct {
	Name string `json:"name"`
	// TODO: Fields below are for the AtlasDeploymentSpec!
	Database       string `json:"database"`
	CollectionName string `json:"collectionName"`
	// TODO: end
	// +kubebuilder:validation:Enum:=search;vectorSearch
	Type string `json:"type,omitempty"`
	// +optional
	Search *SearchIndex `json:"search,omitempty"`
	// +optional
	VectorSearch *VectorSearchIndex `json:"vectorSearch,omitempty"`
}

type VectorSearchIndex struct {
	// Array of JSON objects. See examples https://dochub.mongodb.org/core/avs-vector-type
	Fields string `json:"fields,omitempty"`
}

type SearchIndex struct {
	// +kubebuilder:validation:Enum:=lucene.standard;lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword;lucene.arabic;lucene.armenian;lucene.basque;lucene.bengali;lucene.brazilian;lucene.bulgarian;lucene.catalan;lucene.chinese;lucene.cjk;lucene.czech;lucene.danish;lucene.dutch;lucene.english;lucene.finnish;lucene.french;lucene.galician;lucene.german;lucene.greek;lucene.hindi;lucene.hungarian;lucene.indonesian;lucene.irish;lucene.italian;lucene.japanese;lucene.korean;lucene.kuromoji;lucene.latvian;lucene.lithuanian;lucene.morfologik;lucene.nori;lucene.norwegian;lucene.persian;lucene.portuguese;lucene.romanian;lucene.russian;lucene.smartcn;lucene.sorani;lucene.spanish;lucene.swedish;lucene.thai;lucene.turkish;lucene.ukrainian
	// +optional
	Analyzer string `json:"analyzer,omitempty"`
	// +optional
	Analyzers []AtlasSearchIndexAnalyzer `json:"analyzers,omitempty"`
	// +optional
	SearchAnalyzer string `json:"searchAnalyzer,omitempty"`
	// Flag that indicates whether to store all fields (true) on Atlas Search. By default, Atlas doesn't store (false) the fields on Atlas Search. Alternatively, you can specify an object that only contains the list of fields to store (include) or not store (exclude) on Atlas Search. To learn more, see documentation:
	// https://www.mongodb.com/docs/atlas/atlas-search/stored-source-definition/
	// +optional
	StoredSource string `json:"storedSource,omitempty"`
	// Synonyms and mappings can be found in the AtlasDeployment resource spec

}

type AtlasSearchIndexAnalyzer struct {
	// Human-readable name that identifies the custom analyzer. Names must be unique within an index, and must not start with any of the following strings:
	// "lucene.", "builtin.", "mongodb."
	Name string `json:"name"`
	// +optional
	TokenFilters []TokenFilter `json:"tokenFilters,omitempty"`
	// +optional
	CharFilters []CharFilter `json:"charFilters,omitempty"`
	// +required
	Tokenizer *Tokenizer `json:"tokenizer"`
}

type Tokenizer struct {
	// +kubebuilder:validation:Enum:=whiteSpace;uaxUrlEmail;standard;regexSplit;regexCaptureGroup;nGram;keyword;edgeGram
	Type string `json:"type"`
	// +optional
	// Applied for following types: Whitespace, UaxUrlEmail, Standard
	Whitespace *TokenizerMaxLength `json:"whitespace,omitempty"`
	// +optional
	UaxUrlEmail *TokenizerMaxLength `json:"uaxUrlEmail,omitempty"`
	// +optional
	Standard *TokenizerMaxLength `json:"standard,omitempty"`
	// +optional
	RegexSplit *TokenizerRegexSplit `json:"regexSplit,omitempty"`
	// +optional
	RegexCaptureGroup *TokenizerRegexCaptureGroup `json:"regexCaptureGroup,omitempty"`
	// +optional
	NGram *TokenizerNGram `json:"NGram,omitempty"`
	// +optional
	EdgeGram *TokenizerNGram `json:"edgeGram,omitempty"`
}

type TokenizerMaxLength struct {
	// Maximum number of characters in a single token. Tokens greater than this length are split at this length into multiple tokens.
	// +kubebuilder:default:=255
	MaxTokenLength int `json:"maxTokenLength,omitempty"`
}

type TokenizerRegexSplit struct {
	// Regular expression to match against
	Pattern string `json:"pattern"`
}

type TokenizerRegexCaptureGroup struct {
	// Index of the character group within the matching expression to extract into tokens. Use 0 to extract all character groups
	Group int `json:"group"`
	// Regular expression to match against
	Pattern string `json:"pattern"`
}

type TokenizerNGram struct {
	// Characters to include in the longest token that Atlas Search creates
	MaxGram int `json:"maxGram"`
	// Characters to include in the shortest token that Atlas Search creates
	MinGram int `json:"minGram"`
}

type CharFilter struct {
	// +kubebuilder:valudation:Enum:=htmlStip;icuNormalize;mapping;persian
	Type string `json:"type"`
	// +optional
	HtmlNormalize *CharFilterHTMLNormalize `json:"htmlNormalize,omitempty"`
	// +optional
	Mapping *CharFilterMapping `json:"mapping,omitempty"`
	// +optional
	IcuNormalize *string `json:"icuNormalize,omitempty"`
	// +optional
	Persian *string `json:"persian,omitempty"`
}

type CharFilterHTMLNormalize struct {
	// The HTML tags that you want to exclude from filtering.
	// +optional
	IgnoreTags []string `json:"ignoreTags,omitempty"`
}

type CharFilterMapping struct {
	// Comma-separated list of mappings. A mapping indicates that one character or group of characters should be substituted for another, using the following format: <original> : <replacement>
	Mappings string `json:"mappings"`
}

type TokenFilter struct {
	// Human-readable label that identifies this token filter type
	// +kubebuider:validation:Enum:=asciiFolding;daitchMokotoffSoundex;edgeGram;icuFolding;icuNormalizer;length;lowercase;nGram;regex;reverse;shingle;snowballStemming;stopword;trim
	Type string `json:"type,omitempty"`
	// +optional
	AsciiFolding *FilterAsciiFolding `json:"asciiFolding,omitempty"`
	// +optional
	DaitchMokotoffSoundex *FilterDaitchMokotoffSoundex `json:"daitchMokotoffSoundex,omitempty"`
	// +optional
	EdgeGram *FilterNGram `json:"edgeGram,omitempty"`
	// +optional
	NGram *FilterNGram `json:"nGram,omitempty"`
	// +optional
	IcuNormalizer *FilterIcuNormalizer `json:"icuNormalizer,omitempty"`
	// +optional
	Length *FilterLength `json:"length,omitempty"`
	// +optional
	Regex *FilterRegex `json:"regex,omitempty"`
	// +optional
	Shingle *FilterShingle `json:"shingle,omitempty"`
	// +optional
	SnowballStemming *FilterSnowballStemming `json:"snowballStemming,omitempty"`
	// +optional
	Stopword *FilterStopWord `json:"stopword,omitempty"`
}

type FilterAsciiFolding struct {
	// +kubebuilder:default:=omit
	OriginalTokens string `json:"originalTokens,omitempty"`
}

type FilterDaitchMokotoffSoundex struct {
	// +kubebuilder:default:=include
	OriginalTokens string `json:"originalTokens,omitempty"`
}

type FilterNGram struct {
	MaxGram int `json:"maxGram"`
	MinGram int `json:"minGram"`
	// +kubebuilder:validation:Enum:=omit;include
	TermNotInBounds string `json:"termNotInBounds,omitempty"`
}

type FilterIcuNormalizer struct {
	// +kubebuilder:validation:Enum:=nfc;nfkd;nfkc
	// +kubebuilder:default:=nfc
	NormalizationForm string `json:"normalizationForm,omitempty"`
}

type FilterLength struct {
	// kubebuilder:default:=255
	Max int `json:"max,omitempty"`
	// kubebuilder:default:=0
	Min int `json:"min,omitempty"`
}

type FilterRegex struct {
	// +kubebuilder:validation:Enum:=all;first
	Matches     string `json:"matches"`
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

type FilterShingle struct {
	MaxShingleSize int `json:"maxShingleSize"`
	MinShingleSize int `json:"minShingleSize"`
}

type FilterSnowballStemming struct {
	// +kubebuilder:validation:Enum:=arabic;armenian;basque;catalan;danish;dutch;english;finnish;french;german;german2;hungarian;irish;italian;kp;lithuanian;lovins;norwegian;porter;portuguese;romanian;russian;spanish;swedish;turkish
	StemmerName string `json:"stemmerName"`
}

type FilterStopWord struct {
	Tokens []string `json:"tokens"`
	// +kubebuilder:default:=true
	// +optional
	IgnoreCase bool `json:"ignoreCase"`
}
