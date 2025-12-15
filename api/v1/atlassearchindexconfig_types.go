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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasSearchIndexConfig{}, &AtlasSearchIndexConfigList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=atlas,shortName=asic
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
//
// AtlasSearchIndexConfig is the Schema for the AtlasSearchIndexConfig API
type AtlasSearchIndexConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasSearchIndexConfigSpec          `json:"spec,omitempty"`
	Status status.AtlasSearchIndexConfigStatus `json:"status,omitempty"`
}

// AtlasSearchIndexConfigList contains a list of AtlasSearchIndexConfig
// +kubebuilder:object:root=true
type AtlasSearchIndexConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasSearchIndexConfig `json:"items"`
}

// AtlasSearchIndexConfigSpec defines the target state of AtlasSearchIndexConfig.
type AtlasSearchIndexConfigSpec struct {
	// Specific pre-defined method chosen to convert database field text into searchable words. This conversion reduces the text of fields into the smallest units of text.
	// These units are called a term or token. This process, known as tokenization, involves a variety of changes made to the text in fields:
	// - extracting words
	// - removing punctuation
	// - removing accents
	// - hanging to lowercase
	// - removing common words
	// - reducing words to their root form (stemming)
	// - changing words to their base form (lemmatization) MongoDB Cloud uses the selected process to build the Atlas Search index
	// +kubebuilder:validation:Enum:=lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword;lucene.arabic;lucene.armenian;lucene.basque;lucene.bengali;lucene.brazilian;lucene.bulgarian;lucene.catalan;lucene.chinese;lucene.cjk;lucene.czech;lucene.danish;lucene.dutch;lucene.english;lucene.finnish;lucene.french;lucene.galician;lucene.german;lucene.greek;lucene.hindi;lucene.hungarian;lucene.indonesian;lucene.irish;lucene.italian;lucene.japanese;lucene.korean;lucene.kuromoji;lucene.latvian;lucene.lithuanian;lucene.morfologik;lucene.nori;lucene.norwegian;lucene.persian;lucene.portuguese;lucene.romanian;lucene.russian;lucene.smartcn;lucene.sorani;lucene.spanish;lucene.swedish;lucene.thai;lucene.turkish;lucene.ukrainian
	// +optional
	Analyzer *string `json:"analyzer,omitempty"`
	// List of user-defined methods to convert database field text into searchable words.
	// +optional
	Analyzers *[]AtlasSearchIndexAnalyzer `json:"analyzers,omitempty"`
	// Method applied to identify words when searching this index.
	// +optional
	// +kubebuilder:validation:Enum:=lucene.standard;lucene.simple;lucene.whitespace;lucene.keyword;lucene.arabic;lucene.armenian;lucene.basque;lucene.bengali;lucene.brazilian;lucene.bulgarian;lucene.catalan;lucene.chinese;lucene.cjk;lucene.czech;lucene.danish;lucene.dutch;lucene.english;lucene.finnish;lucene.french;lucene.galician;lucene.german;lucene.greek;lucene.hindi;lucene.hungarian;lucene.indonesian;lucene.irish;lucene.italian;lucene.japanese;lucene.korean;lucene.kuromoji;lucene.latvian;lucene.lithuanian;lucene.morfologik;lucene.nori;lucene.norwegian;lucene.persian;lucene.portuguese;lucene.romanian;lucene.russian;lucene.smartcn;lucene.sorani;lucene.spanish;lucene.swedish;lucene.thai;lucene.turkish;lucene.ukrainian
	SearchAnalyzer *string `json:"searchAnalyzer,omitempty"`
	// Flag that indicates whether to store all fields (true) on Atlas Search. By default, Atlas doesn't store (false) the fields on Atlas Search.
	// Alternatively, you can specify an object that only contains the list of fields to store (include) or not store (exclude) on Atlas Search.
	// To learn more, see documentation: https://www.mongodb.com/docs/atlas/atlas-search/stored-source-definition/
	// +optional
	StoredSource *apiextensions.JSON `json:"storedSource,omitempty"`
	// Synonyms and mappings can be found in the AtlasDeployment resource spec
}

type AtlasSearchIndexAnalyzer struct {
	// Human-readable name that identifies the custom analyzer. Names must be unique within an index, and must not start with any of the following strings:
	// "lucene.", "builtin.", "mongodb."
	// +required
	Name string `json:"name"`
	// Filter that performs operations such as:
	// - Stemming, which reduces related words, such as "talking", "talked", and "talks" to their root word "talk".
	// - Redaction, the removal of sensitive information from public documents
	// +optional
	TokenFilters *apiextensions.JSON `json:"tokenFilters,omitempty"`
	// Filters that examine text one character at a time and perform filtering operations.
	// +optional
	CharFilters *apiextensions.JSON `json:"charFilters,omitempty"`
	// Tokenizer that you want to use to create tokens. Tokens determine how Atlas Search splits up text into discrete chunks for indexing.
	// +required
	Tokenizer Tokenizer `json:"tokenizer"`
}

type Tokenizer struct {
	// Human-readable label that identifies this tokenizer type.
	// +kubebuilder:validation:Enum:=whitespace;uaxUrlEmail;standard;regexSplit;regexCaptureGroup;nGram;keyword;edgeGram;
	// +required
	Type *string `json:"type,omitempty"`
	// Characters to include in the longest token that Atlas Search creates.
	// +optional
	MaxGram *int `json:"maxGram,omitempty"`
	// Characters to include in the shortest token that Atlas Search creates.
	// +optional
	MinGram *int `json:"minGram,omitempty"`
	// Index of the character group within the matching expression to extract into tokens. Use `0` to extract all character groups.
	// +optional
	Group *int `json:"group,omitempty"`
	// Regular expression to match against.
	// +optional
	Pattern *string `json:"pattern,omitempty"`
	// Maximum number of characters in a single token. Tokens greater than this length are split at this length into multiple tokens.
	// +optional
	MaxTokenLength *int `json:"maxTokenLength,omitempty"`
}

type CharFilter struct {
	// Human-readable label that identifies this character filter type.
	// +kubebuilder:validation:Enum:=htmlStrip;icuNormalize;mapping;persian
	// +required
	Type string `json:"type,omitempty"`
	// Used when type is "htmlStrip". The HTML tags that you want to exclude from filtering.
	// +optional
	IgnoreTags []string `json:"ignoreTags,omitempty"`
	// Used when type is "mapping". Comma-separated list of mappings. A mapping indicates that one character or group of characters should be substituted for another, using the following format: <original> : <replacement>
	// +optional
	Mappings map[string]string `json:"mappings"`
}

type TokenFilter struct {
	// +required
	// +kubebuilder:validation:Enum:=asciiFolding;daitchMokotoffSoundex;edgeGram;icuFolding;icuNormalizer;length;lowercase;nGram;regex;reverse;shingle;snowballStemming;stopword;trim;
	Type string `json:"type,omitempty"`
	// Value that indicates whether to include or omit the original tokens in the output of the token filter.
	// Choose include if you want to support queries on both the original tokens and the converted forms.
	// +kubebuilder:validation:Enum:=omit;include;
	// +optional
	OriginalTokens string `json:"originalTokens,omitempty"`
	// Value that specifies the maximum length of generated n-grams. This value must be greater than or equal to minGram.
	// +optional
	MaxGram int `json:"maxGram,omitempty"`
	// Value that specifies the minimum length of generated n-grams. This value must be less than or equal to maxGram.
	// +optional
	MinGram int `json:"minGram,omitempty"`
	// Value that indicates whether to index tokens shorter than minGram or longer than maxGram.
	// +kubebuilder:validation:Enum:=omit;include
	// +optional
	TermNotInBounds string `json:"termNotInBounds,omitempty"`
	// Normalization form to apply
	// +kubebuilder:validation:Enum:=nfc;nfkd;nfkc
	// +optional
	NormalizationForm string `json:"normalizationForm,omitempty"`
	// Number that specifies the maximum length of a token. Value must be greater than or equal to min.
	// +optional
	Max int `json:"max,omitempty"`
	// Number that specifies the minimum length of a token. This value must be less than or equal to max.
	// +optional
	Min int `json:"min,omitempty"`
	// Value that indicates whether to replace only the first matching pattern or all matching patterns.
	// +kubebuilder:validation:Enum:=all;first
	// +optional
	Matches string `json:"matches,omitempty"`
	// Regular expression pattern to apply to each token.
	// +optional
	Pattern string `json:"pattern,omitempty"`
	// Replacement string to substitute wherever a matching pattern occurs.
	// +optional
	Replacement string `json:"replacement,omitempty"`
	// Value that specifies the maximum number of tokens per shingle. This value must be greater than or equal to minShingleSize.
	// +optional
	MaxShingleSize int `json:"maxShingleSize,omitempty"`
	// Value that specifies the minimum number of tokens per shingle. This value must be less than or equal to maxShingleSize.
	// +optional
	MinShingleSize int `json:"minShingleSize,omitempty"`
	// Snowball-generated stemmer to use.
	// +kubebuilder:validation:Enum:=arabic;armenian;basque;catalan;danish;dutch;english;finnish;french;german;german2;hungarian;irish;italian;kp;lithuanian;lovins;norwegian;porter;portuguese;romanian;russian;spanish;swedish;turkish
	// +optional
	StemmerName string `json:"stemmerName,omitempty"`
	// The stop words that correspond to the tokens to remove. Value must be one or more stop words.
	// +optional
	Tokens []string `json:"tokens,omitempty"`
	// Flag that indicates whether to ignore the case of stop words when filtering the tokens to remove.
	// +optional
	IgnoreCase bool `json:"ignoreCase,omitempty"`
}

func (s *AtlasSearchIndexConfig) GetStatus() api.Status {
	return s.Status
}

func (s *AtlasSearchIndexConfig) UpdateStatus(conditions []api.Condition, opts ...api.Option) {
	s.Status.Conditions = conditions
	s.Status.ObservedGeneration = s.ObjectMeta.Generation

	for _, o := range opts {
		v := o.(status.AtlasSearchIndexConfigStatusOption)
		v(&s.Status)
	}
}
