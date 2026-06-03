// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AtlasOrganization Details that describe the organization.
type AtlasOrganization struct {
	// Unique 24-hexadecimal digit string that identifies the organization.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Flag that indicates whether this organization has been deleted.
	// Read only field.
	IsDeleted *bool `json:"isDeleted,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the organization.
	Name string `json:"name"`
	// Disables automatic alert creation. When set to true, no organization level alerts will be created automatically.
	SkipDefaultAlertsSettings *bool `json:"skipDefaultAlertsSettings,omitempty"`
}

// NewAtlasOrganization instantiates a new AtlasOrganization object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAtlasOrganization(name string) *AtlasOrganization {
	this := AtlasOrganization{}
	this.Name = name
	var skipDefaultAlertsSettings bool = false
	this.SkipDefaultAlertsSettings = &skipDefaultAlertsSettings
	return &this
}

// NewAtlasOrganizationWithDefaults instantiates a new AtlasOrganization object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAtlasOrganizationWithDefaults() *AtlasOrganization {
	this := AtlasOrganization{}
	var skipDefaultAlertsSettings bool = false
	this.SkipDefaultAlertsSettings = &skipDefaultAlertsSettings
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *AtlasOrganization) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasOrganization) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *AtlasOrganization) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *AtlasOrganization) SetId(v string) {
	o.Id = &v
}

// GetIsDeleted returns the IsDeleted field value if set, zero value otherwise
func (o *AtlasOrganization) GetIsDeleted() bool {
	if o == nil || IsNil(o.IsDeleted) {
		var ret bool
		return ret
	}
	return *o.IsDeleted
}

// GetIsDeletedOk returns a tuple with the IsDeleted field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasOrganization) GetIsDeletedOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDeleted) {
		return nil, false
	}

	return o.IsDeleted, true
}

// HasIsDeleted returns a boolean if a field has been set.
func (o *AtlasOrganization) HasIsDeleted() bool {
	if o != nil && !IsNil(o.IsDeleted) {
		return true
	}

	return false
}

// SetIsDeleted gets a reference to the given bool and assigns it to the IsDeleted field.
func (o *AtlasOrganization) SetIsDeleted(v bool) {
	o.IsDeleted = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *AtlasOrganization) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasOrganization) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *AtlasOrganization) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *AtlasOrganization) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value
func (o *AtlasOrganization) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *AtlasOrganization) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *AtlasOrganization) SetName(v string) {
	o.Name = v
}

// GetSkipDefaultAlertsSettings returns the SkipDefaultAlertsSettings field value if set, zero value otherwise
func (o *AtlasOrganization) GetSkipDefaultAlertsSettings() bool {
	if o == nil || IsNil(o.SkipDefaultAlertsSettings) {
		var ret bool
		return ret
	}
	return *o.SkipDefaultAlertsSettings
}

// GetSkipDefaultAlertsSettingsOk returns a tuple with the SkipDefaultAlertsSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasOrganization) GetSkipDefaultAlertsSettingsOk() (*bool, bool) {
	if o == nil || IsNil(o.SkipDefaultAlertsSettings) {
		return nil, false
	}

	return o.SkipDefaultAlertsSettings, true
}

// HasSkipDefaultAlertsSettings returns a boolean if a field has been set.
func (o *AtlasOrganization) HasSkipDefaultAlertsSettings() bool {
	if o != nil && !IsNil(o.SkipDefaultAlertsSettings) {
		return true
	}

	return false
}

// SetSkipDefaultAlertsSettings gets a reference to the given bool and assigns it to the SkipDefaultAlertsSettings field.
func (o *AtlasOrganization) SetSkipDefaultAlertsSettings(v bool) {
	o.SkipDefaultAlertsSettings = &v
}
