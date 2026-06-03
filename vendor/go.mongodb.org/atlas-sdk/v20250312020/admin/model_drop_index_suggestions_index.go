// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DropIndexSuggestionsIndex struct for DropIndexSuggestionsIndex
type DropIndexSuggestionsIndex struct {
	// Usage count (since last restart) of index.
	AccessCount *int64 `json:"accessCount,omitempty"`
	// List that contains documents that specify a key in the index and its sort order.
	Index *[]any `json:"index,omitempty"`
	// Name of index.
	Name *string `json:"name,omitempty"`
	// Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
	Namespace *string `json:"namespace,omitempty"`
	// List that contains strings that specifies the shards where the index is found.
	Shards *[]string `json:"shards,omitempty"`
	// Date of most recent usage of index. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	Since *time.Time `json:"since,omitempty"`
	// Size of index.
	SizeBytes *int64 `json:"sizeBytes,omitempty"`
}

// NewDropIndexSuggestionsIndex instantiates a new DropIndexSuggestionsIndex object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDropIndexSuggestionsIndex() *DropIndexSuggestionsIndex {
	this := DropIndexSuggestionsIndex{}
	return &this
}

// NewDropIndexSuggestionsIndexWithDefaults instantiates a new DropIndexSuggestionsIndex object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDropIndexSuggestionsIndexWithDefaults() *DropIndexSuggestionsIndex {
	this := DropIndexSuggestionsIndex{}
	return &this
}

// GetAccessCount returns the AccessCount field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetAccessCount() int64 {
	if o == nil || IsNil(o.AccessCount) {
		var ret int64
		return ret
	}
	return *o.AccessCount
}

// GetAccessCountOk returns a tuple with the AccessCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetAccessCountOk() (*int64, bool) {
	if o == nil || IsNil(o.AccessCount) {
		return nil, false
	}

	return o.AccessCount, true
}

// HasAccessCount returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasAccessCount() bool {
	if o != nil && !IsNil(o.AccessCount) {
		return true
	}

	return false
}

// SetAccessCount gets a reference to the given int64 and assigns it to the AccessCount field.
func (o *DropIndexSuggestionsIndex) SetAccessCount(v int64) {
	o.AccessCount = &v
}

// GetIndex returns the Index field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetIndex() []any {
	if o == nil || IsNil(o.Index) {
		var ret []any
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetIndexOk() (*[]any, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}

	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given []any and assigns it to the Index field.
func (o *DropIndexSuggestionsIndex) SetIndex(v []any) {
	o.Index = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DropIndexSuggestionsIndex) SetName(v string) {
	o.Name = &v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *DropIndexSuggestionsIndex) SetNamespace(v string) {
	o.Namespace = &v
}

// GetShards returns the Shards field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetShards() []string {
	if o == nil || IsNil(o.Shards) {
		var ret []string
		return ret
	}
	return *o.Shards
}

// GetShardsOk returns a tuple with the Shards field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetShardsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Shards) {
		return nil, false
	}

	return o.Shards, true
}

// HasShards returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasShards() bool {
	if o != nil && !IsNil(o.Shards) {
		return true
	}

	return false
}

// SetShards gets a reference to the given []string and assigns it to the Shards field.
func (o *DropIndexSuggestionsIndex) SetShards(v []string) {
	o.Shards = &v
}

// GetSince returns the Since field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetSince() time.Time {
	if o == nil || IsNil(o.Since) {
		var ret time.Time
		return ret
	}
	return *o.Since
}

// GetSinceOk returns a tuple with the Since field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetSinceOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Since) {
		return nil, false
	}

	return o.Since, true
}

// HasSince returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasSince() bool {
	if o != nil && !IsNil(o.Since) {
		return true
	}

	return false
}

// SetSince gets a reference to the given time.Time and assigns it to the Since field.
func (o *DropIndexSuggestionsIndex) SetSince(v time.Time) {
	o.Since = &v
}

// GetSizeBytes returns the SizeBytes field value if set, zero value otherwise
func (o *DropIndexSuggestionsIndex) GetSizeBytes() int64 {
	if o == nil || IsNil(o.SizeBytes) {
		var ret int64
		return ret
	}
	return *o.SizeBytes
}

// GetSizeBytesOk returns a tuple with the SizeBytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DropIndexSuggestionsIndex) GetSizeBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.SizeBytes) {
		return nil, false
	}

	return o.SizeBytes, true
}

// HasSizeBytes returns a boolean if a field has been set.
func (o *DropIndexSuggestionsIndex) HasSizeBytes() bool {
	if o != nil && !IsNil(o.SizeBytes) {
		return true
	}

	return false
}

// SetSizeBytes gets a reference to the given int64 and assigns it to the SizeBytes field.
func (o *DropIndexSuggestionsIndex) SetSizeBytes(v int64) {
	o.SizeBytes = &v
}
