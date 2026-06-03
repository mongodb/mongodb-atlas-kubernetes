// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamConfig Configuration options for an Atlas Stream Processing Workspace.
type StreamConfig struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Max tier size for the Stream Workspace. Configures Memory / VCPU allowances.
	MaxTierSize *string `json:"maxTierSize,omitempty"`
	// Selected tier for the Stream Workspace. Configures Memory / VCPU allowances.
	Tier *string `json:"tier,omitempty"`
}

// NewStreamConfig instantiates a new StreamConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamConfig() *StreamConfig {
	this := StreamConfig{}
	return &this
}

// NewStreamConfigWithDefaults instantiates a new StreamConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamConfigWithDefaults() *StreamConfig {
	this := StreamConfig{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamConfig) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamConfig) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamConfig) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamConfig) SetLinks(v []Link) {
	o.Links = &v
}

// GetMaxTierSize returns the MaxTierSize field value if set, zero value otherwise
func (o *StreamConfig) GetMaxTierSize() string {
	if o == nil || IsNil(o.MaxTierSize) {
		var ret string
		return ret
	}
	return *o.MaxTierSize
}

// GetMaxTierSizeOk returns a tuple with the MaxTierSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamConfig) GetMaxTierSizeOk() (*string, bool) {
	if o == nil || IsNil(o.MaxTierSize) {
		return nil, false
	}

	return o.MaxTierSize, true
}

// HasMaxTierSize returns a boolean if a field has been set.
func (o *StreamConfig) HasMaxTierSize() bool {
	if o != nil && !IsNil(o.MaxTierSize) {
		return true
	}

	return false
}

// SetMaxTierSize gets a reference to the given string and assigns it to the MaxTierSize field.
func (o *StreamConfig) SetMaxTierSize(v string) {
	o.MaxTierSize = &v
}

// GetTier returns the Tier field value if set, zero value otherwise
func (o *StreamConfig) GetTier() string {
	if o == nil || IsNil(o.Tier) {
		var ret string
		return ret
	}
	return *o.Tier
}

// GetTierOk returns a tuple with the Tier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamConfig) GetTierOk() (*string, bool) {
	if o == nil || IsNil(o.Tier) {
		return nil, false
	}

	return o.Tier, true
}

// HasTier returns a boolean if a field has been set.
func (o *StreamConfig) HasTier() bool {
	if o != nil && !IsNil(o.Tier) {
		return true
	}

	return false
}

// SetTier gets a reference to the given string and assigns it to the Tier field.
func (o *StreamConfig) SetTier(v string) {
	o.Tier = &v
}
