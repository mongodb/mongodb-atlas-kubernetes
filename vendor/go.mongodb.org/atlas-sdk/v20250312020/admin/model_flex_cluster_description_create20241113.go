// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FlexClusterDescriptionCreate20241113 Settings that you can specify when you create a flex cluster.
type FlexClusterDescriptionCreate20241113 struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the instance.
	// Write only field.
	Name             string                             `json:"name"`
	ProviderSettings FlexProviderSettingsCreate20241113 `json:"providerSettings"`
	// List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the instance.
	Tags *[]ResourceTag `json:"tags,omitempty"`
	// Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.
	TerminationProtectionEnabled *bool `json:"terminationProtectionEnabled,omitempty"`
}

// NewFlexClusterDescriptionCreate20241113 instantiates a new FlexClusterDescriptionCreate20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFlexClusterDescriptionCreate20241113(name string, providerSettings FlexProviderSettingsCreate20241113) *FlexClusterDescriptionCreate20241113 {
	this := FlexClusterDescriptionCreate20241113{}
	this.Name = name
	this.ProviderSettings = providerSettings
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	return &this
}

// NewFlexClusterDescriptionCreate20241113WithDefaults instantiates a new FlexClusterDescriptionCreate20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFlexClusterDescriptionCreate20241113WithDefaults() *FlexClusterDescriptionCreate20241113 {
	this := FlexClusterDescriptionCreate20241113{}
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *FlexClusterDescriptionCreate20241113) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionCreate20241113) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *FlexClusterDescriptionCreate20241113) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *FlexClusterDescriptionCreate20241113) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value
func (o *FlexClusterDescriptionCreate20241113) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionCreate20241113) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *FlexClusterDescriptionCreate20241113) SetName(v string) {
	o.Name = v
}

// GetProviderSettings returns the ProviderSettings field value
func (o *FlexClusterDescriptionCreate20241113) GetProviderSettings() FlexProviderSettingsCreate20241113 {
	if o == nil {
		var ret FlexProviderSettingsCreate20241113
		return ret
	}

	return o.ProviderSettings
}

// GetProviderSettingsOk returns a tuple with the ProviderSettings field value
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionCreate20241113) GetProviderSettingsOk() (*FlexProviderSettingsCreate20241113, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderSettings, true
}

// SetProviderSettings sets field value
func (o *FlexClusterDescriptionCreate20241113) SetProviderSettings(v FlexProviderSettingsCreate20241113) {
	o.ProviderSettings = v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *FlexClusterDescriptionCreate20241113) GetTags() []ResourceTag {
	if o == nil || IsNil(o.Tags) {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionCreate20241113) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *FlexClusterDescriptionCreate20241113) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *FlexClusterDescriptionCreate20241113) SetTags(v []ResourceTag) {
	o.Tags = &v
}

// GetTerminationProtectionEnabled returns the TerminationProtectionEnabled field value if set, zero value otherwise
func (o *FlexClusterDescriptionCreate20241113) GetTerminationProtectionEnabled() bool {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		var ret bool
		return ret
	}
	return *o.TerminationProtectionEnabled
}

// GetTerminationProtectionEnabledOk returns a tuple with the TerminationProtectionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexClusterDescriptionCreate20241113) GetTerminationProtectionEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		return nil, false
	}

	return o.TerminationProtectionEnabled, true
}

// HasTerminationProtectionEnabled returns a boolean if a field has been set.
func (o *FlexClusterDescriptionCreate20241113) HasTerminationProtectionEnabled() bool {
	if o != nil && !IsNil(o.TerminationProtectionEnabled) {
		return true
	}

	return false
}

// SetTerminationProtectionEnabled gets a reference to the given bool and assigns it to the TerminationProtectionEnabled field.
func (o *FlexClusterDescriptionCreate20241113) SetTerminationProtectionEnabled(v bool) {
	o.TerminationProtectionEnabled = &v
}
