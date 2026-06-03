// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FlexClusterDescriptionUpdate20241113 Settings that you can specify when you update a flex cluster.
type FlexClusterDescriptionUpdate20241113 struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the instance.
	Tags *[]ResourceTag `json:"tags,omitempty"`
	// Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.
	TerminationProtectionEnabled *bool `json:"terminationProtectionEnabled,omitempty"`
}

// NewFlexClusterDescriptionUpdate20241113 instantiates a new FlexClusterDescriptionUpdate20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFlexClusterDescriptionUpdate20241113() *FlexClusterDescriptionUpdate20241113 {
	this := FlexClusterDescriptionUpdate20241113{}
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	return &this
}

// NewFlexClusterDescriptionUpdate20241113WithDefaults instantiates a new FlexClusterDescriptionUpdate20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFlexClusterDescriptionUpdate20241113WithDefaults() *FlexClusterDescriptionUpdate20241113 {
	this := FlexClusterDescriptionUpdate20241113{}
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *FlexClusterDescriptionUpdate20241113) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionUpdate20241113) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *FlexClusterDescriptionUpdate20241113) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *FlexClusterDescriptionUpdate20241113) SetLinks(v []Link) {
	o.Links = &v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *FlexClusterDescriptionUpdate20241113) GetTags() []ResourceTag {
	if o == nil || IsNil(o.Tags) {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionUpdate20241113) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *FlexClusterDescriptionUpdate20241113) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *FlexClusterDescriptionUpdate20241113) SetTags(v []ResourceTag) {
	o.Tags = &v
}

// GetTerminationProtectionEnabled returns the TerminationProtectionEnabled field value if set, zero value otherwise
func (o *FlexClusterDescriptionUpdate20241113) GetTerminationProtectionEnabled() bool {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		var ret bool
		return ret
	}
	return *o.TerminationProtectionEnabled
}

// GetTerminationProtectionEnabledOk returns a tuple with the TerminationProtectionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionUpdate20241113) GetTerminationProtectionEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		return nil, false
	}

	return o.TerminationProtectionEnabled, true
}

// HasTerminationProtectionEnabled returns a boolean if a field has been set.
func (o *FlexClusterDescriptionUpdate20241113) HasTerminationProtectionEnabled() bool {
	if o != nil && !IsNil(o.TerminationProtectionEnabled) {
		return true
	}

	return false
}

// SetTerminationProtectionEnabled gets a reference to the given bool and assigns it to the TerminationProtectionEnabled field.
func (o *FlexClusterDescriptionUpdate20241113) SetTerminationProtectionEnabled(v bool) {
	o.TerminationProtectionEnabled = &v
}
