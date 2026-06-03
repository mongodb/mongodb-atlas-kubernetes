// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DropIndexSuggestionsResponse struct for DropIndexSuggestionsResponse
type DropIndexSuggestionsResponse struct {
	// List that contains the documents with information about the hidden indexes that the Performance Advisor suggests to remove.
	// Read only field.
	HiddenIndexes *[]DropIndexSuggestionsIndex `json:"hiddenIndexes,omitempty"`
	// List that contains the documents with information about the redundant indexes that the Performance Advisor suggests to remove.
	// Read only field.
	RedundantIndexes *[]DropIndexSuggestionsIndex `json:"redundantIndexes,omitempty"`
	// List that contains the documents with information about the unused indexes that the Performance Advisor suggests to remove.
	// Read only field.
	UnusedIndexes *[]DropIndexSuggestionsIndex `json:"unusedIndexes,omitempty"`
}

// NewDropIndexSuggestionsResponse instantiates a new DropIndexSuggestionsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDropIndexSuggestionsResponse() *DropIndexSuggestionsResponse {
	this := DropIndexSuggestionsResponse{}
	return &this
}

// NewDropIndexSuggestionsResponseWithDefaults instantiates a new DropIndexSuggestionsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDropIndexSuggestionsResponseWithDefaults() *DropIndexSuggestionsResponse {
	this := DropIndexSuggestionsResponse{}
	return &this
}

// GetHiddenIndexes returns the HiddenIndexes field value if set, zero value otherwise
func (o *DropIndexSuggestionsResponse) GetHiddenIndexes() []DropIndexSuggestionsIndex {
	if o == nil || IsNil(o.HiddenIndexes) {
		var ret []DropIndexSuggestionsIndex
		return ret
	}
	return *o.HiddenIndexes
}

// GetHiddenIndexesOk returns a tuple with the HiddenIndexes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsResponse) GetHiddenIndexesOk() (*[]DropIndexSuggestionsIndex, bool) {
	if o == nil || IsNil(o.HiddenIndexes) {
		return nil, false
	}

	return o.HiddenIndexes, true
}

// HasHiddenIndexes returns a boolean if a field has been set.
func (o *DropIndexSuggestionsResponse) HasHiddenIndexes() bool {
	if o != nil && !IsNil(o.HiddenIndexes) {
		return true
	}

	return false
}

// SetHiddenIndexes gets a reference to the given []DropIndexSuggestionsIndex and assigns it to the HiddenIndexes field.
func (o *DropIndexSuggestionsResponse) SetHiddenIndexes(v []DropIndexSuggestionsIndex) {
	o.HiddenIndexes = &v
}

// GetRedundantIndexes returns the RedundantIndexes field value if set, zero value otherwise
func (o *DropIndexSuggestionsResponse) GetRedundantIndexes() []DropIndexSuggestionsIndex {
	if o == nil || IsNil(o.RedundantIndexes) {
		var ret []DropIndexSuggestionsIndex
		return ret
	}
	return *o.RedundantIndexes
}

// GetRedundantIndexesOk returns a tuple with the RedundantIndexes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsResponse) GetRedundantIndexesOk() (*[]DropIndexSuggestionsIndex, bool) {
	if o == nil || IsNil(o.RedundantIndexes) {
		return nil, false
	}

	return o.RedundantIndexes, true
}

// HasRedundantIndexes returns a boolean if a field has been set.
func (o *DropIndexSuggestionsResponse) HasRedundantIndexes() bool {
	if o != nil && !IsNil(o.RedundantIndexes) {
		return true
	}

	return false
}

// SetRedundantIndexes gets a reference to the given []DropIndexSuggestionsIndex and assigns it to the RedundantIndexes field.
func (o *DropIndexSuggestionsResponse) SetRedundantIndexes(v []DropIndexSuggestionsIndex) {
	o.RedundantIndexes = &v
}

// GetUnusedIndexes returns the UnusedIndexes field value if set, zero value otherwise
func (o *DropIndexSuggestionsResponse) GetUnusedIndexes() []DropIndexSuggestionsIndex {
	if o == nil || IsNil(o.UnusedIndexes) {
		var ret []DropIndexSuggestionsIndex
		return ret
	}
	return *o.UnusedIndexes
}

// GetUnusedIndexesOk returns a tuple with the UnusedIndexes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsResponse) GetUnusedIndexesOk() (*[]DropIndexSuggestionsIndex, bool) {
	if o == nil || IsNil(o.UnusedIndexes) {
		return nil, false
	}

	return o.UnusedIndexes, true
}

// HasUnusedIndexes returns a boolean if a field has been set.
func (o *DropIndexSuggestionsResponse) HasUnusedIndexes() bool {
	if o != nil && !IsNil(o.UnusedIndexes) {
		return true
	}

	return false
}

// SetUnusedIndexes gets a reference to the given []DropIndexSuggestionsIndex and assigns it to the UnusedIndexes field.
func (o *DropIndexSuggestionsResponse) SetUnusedIndexes(v []DropIndexSuggestionsIndex) {
	o.UnusedIndexes = &v
}
