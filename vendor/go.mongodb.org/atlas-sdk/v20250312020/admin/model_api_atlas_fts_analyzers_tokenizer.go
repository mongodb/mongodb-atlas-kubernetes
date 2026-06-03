// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasFTSAnalyzersTokenizer Tokenizer that you want to use to create tokens. Tokens determine how Atlas Search splits up text into discrete chunks for indexing.
type ApiAtlasFTSAnalyzersTokenizer struct {
	// Characters to include in the longest token that Atlas Search creates.
	MaxGram *int `json:"maxGram,omitempty"`
	// Characters to include in the shortest token that Atlas Search creates.
	MinGram *int `json:"minGram,omitempty"`
	// Human-readable label that identifies this tokenizer type.
	Type *string `json:"type,omitempty"`
	// Index of the character group within the matching expression to extract into tokens. Use `0` to extract all character groups.
	Group *int `json:"group,omitempty"`
	// Regular expression to match against.
	Pattern *string `json:"pattern,omitempty"`
	// Maximum number of characters in a single token. Tokens greater than this length are split at this length into multiple tokens.
	MaxTokenLength *int `json:"maxTokenLength,omitempty"`
}

// NewApiAtlasFTSAnalyzersTokenizer instantiates a new ApiAtlasFTSAnalyzersTokenizer object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasFTSAnalyzersTokenizer() *ApiAtlasFTSAnalyzersTokenizer {
	this := ApiAtlasFTSAnalyzersTokenizer{}
	var maxTokenLength int = 255
	this.MaxTokenLength = &maxTokenLength
	return &this
}

// NewApiAtlasFTSAnalyzersTokenizerWithDefaults instantiates a new ApiAtlasFTSAnalyzersTokenizer object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasFTSAnalyzersTokenizerWithDefaults() *ApiAtlasFTSAnalyzersTokenizer {
	this := ApiAtlasFTSAnalyzersTokenizer{}
	var maxTokenLength int = 255
	this.MaxTokenLength = &maxTokenLength
	return &this
}

// GetMaxGram returns the MaxGram field value if set, zero value otherwise
func (o *ApiAtlasFTSAnalyzersTokenizer) GetMaxGram() int {
	if o == nil || IsNil(o.MaxGram) {
		var ret int
		return ret
	}
	return *o.MaxGram
}

// GetMaxGramOk returns a tuple with the MaxGram field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) GetMaxGramOk() (*int, bool) {
	if o == nil || IsNil(o.MaxGram) {
		return nil, false
	}

	return o.MaxGram, true
}

// HasMaxGram returns a boolean if a field has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) HasMaxGram() bool {
	if o != nil && !IsNil(o.MaxGram) {
		return true
	}

	return false
}

// SetMaxGram gets a reference to the given int and assigns it to the MaxGram field.
func (o *ApiAtlasFTSAnalyzersTokenizer) SetMaxGram(v int) {
	o.MaxGram = &v
}

// GetMinGram returns the MinGram field value if set, zero value otherwise
func (o *ApiAtlasFTSAnalyzersTokenizer) GetMinGram() int {
	if o == nil || IsNil(o.MinGram) {
		var ret int
		return ret
	}
	return *o.MinGram
}

// GetMinGramOk returns a tuple with the MinGram field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) GetMinGramOk() (*int, bool) {
	if o == nil || IsNil(o.MinGram) {
		return nil, false
	}

	return o.MinGram, true
}

// HasMinGram returns a boolean if a field has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) HasMinGram() bool {
	if o != nil && !IsNil(o.MinGram) {
		return true
	}

	return false
}

// SetMinGram gets a reference to the given int and assigns it to the MinGram field.
func (o *ApiAtlasFTSAnalyzersTokenizer) SetMinGram(v int) {
	o.MinGram = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *ApiAtlasFTSAnalyzersTokenizer) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ApiAtlasFTSAnalyzersTokenizer) SetType(v string) {
	o.Type = &v
}

// GetGroup returns the Group field value if set, zero value otherwise
func (o *ApiAtlasFTSAnalyzersTokenizer) GetGroup() int {
	if o == nil || IsNil(o.Group) {
		var ret int
		return ret
	}
	return *o.Group
}

// GetGroupOk returns a tuple with the Group field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) GetGroupOk() (*int, bool) {
	if o == nil || IsNil(o.Group) {
		return nil, false
	}

	return o.Group, true
}

// HasGroup returns a boolean if a field has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) HasGroup() bool {
	if o != nil && !IsNil(o.Group) {
		return true
	}

	return false
}

// SetGroup gets a reference to the given int and assigns it to the Group field.
func (o *ApiAtlasFTSAnalyzersTokenizer) SetGroup(v int) {
	o.Group = &v
}

// GetPattern returns the Pattern field value if set, zero value otherwise
func (o *ApiAtlasFTSAnalyzersTokenizer) GetPattern() string {
	if o == nil || IsNil(o.Pattern) {
		var ret string
		return ret
	}
	return *o.Pattern
}

// GetPatternOk returns a tuple with the Pattern field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) GetPatternOk() (*string, bool) {
	if o == nil || IsNil(o.Pattern) {
		return nil, false
	}

	return o.Pattern, true
}

// HasPattern returns a boolean if a field has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) HasPattern() bool {
	if o != nil && !IsNil(o.Pattern) {
		return true
	}

	return false
}

// SetPattern gets a reference to the given string and assigns it to the Pattern field.
func (o *ApiAtlasFTSAnalyzersTokenizer) SetPattern(v string) {
	o.Pattern = &v
}

// GetMaxTokenLength returns the MaxTokenLength field value if set, zero value otherwise
func (o *ApiAtlasFTSAnalyzersTokenizer) GetMaxTokenLength() int {
	if o == nil || IsNil(o.MaxTokenLength) {
		var ret int
		return ret
	}
	return *o.MaxTokenLength
}

// GetMaxTokenLengthOk returns a tuple with the MaxTokenLength field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) GetMaxTokenLengthOk() (*int, bool) {
	if o == nil || IsNil(o.MaxTokenLength) {
		return nil, false
	}

	return o.MaxTokenLength, true
}

// HasMaxTokenLength returns a boolean if a field has been set.
func (o *ApiAtlasFTSAnalyzersTokenizer) HasMaxTokenLength() bool {
	if o != nil && !IsNil(o.MaxTokenLength) {
		return true
	}

	return false
}

// SetMaxTokenLength gets a reference to the given int and assigns it to the MaxTokenLength field.
func (o *ApiAtlasFTSAnalyzersTokenizer) SetMaxTokenLength(v int) {
	o.MaxTokenLength = &v
}
