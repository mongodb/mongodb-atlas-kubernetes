// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// Group struct for Group
type Group struct {
	// Quantity of MongoDB Cloud clusters deployed in this project.
	// Read only field.
	ClusterCount int64 `json:"clusterCount"`
	// Date and time when MongoDB Cloud created this project. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created time.Time `json:"created"`
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud project.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the project included in the MongoDB Cloud organization.
	Name string `json:"name"`
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud organization to which the project belongs.
	OrgId string `json:"orgId"`
	// Applies to Atlas for Government only.  In Commercial Atlas, this field will be rejected in requests and missing in responses.  This field sets restrictions on available regions in the project.  `COMMERCIAL_FEDRAMP_REGIONS_ONLY`: Only allows deployments in FedRAMP Moderate regions.  `GOV_REGIONS_ONLY`: Only allows deployments in GovCloud regions.
	RegionUsageRestrictions *string `json:"regionUsageRestrictions,omitempty"`
	// List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the project.
	Tags *[]ResourceTag `json:"tags,omitempty"`
	// Flag that indicates whether to create the project with default alert settings.
	WithDefaultAlertsSettings *bool `json:"withDefaultAlertsSettings,omitempty"`
}

// NewGroup instantiates a new Group object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroup(clusterCount int64, created time.Time, name string, orgId string) *Group {
	this := Group{}
	this.ClusterCount = clusterCount
	this.Created = created
	this.Name = name
	this.OrgId = orgId
	var regionUsageRestrictions string = "COMMERCIAL_FEDRAMP_REGIONS_ONLY"
	this.RegionUsageRestrictions = &regionUsageRestrictions
	var withDefaultAlertsSettings bool = true
	this.WithDefaultAlertsSettings = &withDefaultAlertsSettings
	return &this
}

// NewGroupWithDefaults instantiates a new Group object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupWithDefaults() *Group {
	this := Group{}
	var regionUsageRestrictions string = "COMMERCIAL_FEDRAMP_REGIONS_ONLY"
	this.RegionUsageRestrictions = &regionUsageRestrictions
	var withDefaultAlertsSettings bool = true
	this.WithDefaultAlertsSettings = &withDefaultAlertsSettings
	return &this
}

// GetClusterCount returns the ClusterCount field value
func (o *Group) GetClusterCount() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.ClusterCount
}

// GetClusterCountOk returns a tuple with the ClusterCount field value
// and a boolean to check if the value has been set.
func (o *Group) GetClusterCountOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterCount, true
}

// SetClusterCount sets field value
func (o *Group) SetClusterCount(v int64) {
	o.ClusterCount = v
}

// GetCreated returns the Created field value
func (o *Group) GetCreated() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Created
}

// GetCreatedOk returns a tuple with the Created field value
// and a boolean to check if the value has been set.
func (o *Group) GetCreatedOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Created, true
}

// SetCreated sets field value
func (o *Group) SetCreated(v time.Time) {
	o.Created = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *Group) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Group) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Group) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Group) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *Group) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Group) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *Group) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *Group) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value
func (o *Group) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Group) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Group) SetName(v string) {
	o.Name = v
}

// GetOrgId returns the OrgId field value
func (o *Group) GetOrgId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value
// and a boolean to check if the value has been set.
func (o *Group) GetOrgIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrgId, true
}

// SetOrgId sets field value
func (o *Group) SetOrgId(v string) {
	o.OrgId = v
}

// GetRegionUsageRestrictions returns the RegionUsageRestrictions field value if set, zero value otherwise
func (o *Group) GetRegionUsageRestrictions() string {
	if o == nil || IsNil(o.RegionUsageRestrictions) {
		var ret string
		return ret
	}
	return *o.RegionUsageRestrictions
}

// GetRegionUsageRestrictionsOk returns a tuple with the RegionUsageRestrictions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Group) GetRegionUsageRestrictionsOk() (*string, bool) {
	if o == nil || IsNil(o.RegionUsageRestrictions) {
		return nil, false
	}

	return o.RegionUsageRestrictions, true
}

// HasRegionUsageRestrictions returns a boolean if a field has been set.
func (o *Group) HasRegionUsageRestrictions() bool {
	if o != nil && !IsNil(o.RegionUsageRestrictions) {
		return true
	}

	return false
}

// SetRegionUsageRestrictions gets a reference to the given string and assigns it to the RegionUsageRestrictions field.
func (o *Group) SetRegionUsageRestrictions(v string) {
	o.RegionUsageRestrictions = &v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *Group) GetTags() []ResourceTag {
	if o == nil || IsNil(o.Tags) {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Group) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *Group) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *Group) SetTags(v []ResourceTag) {
	o.Tags = &v
}

// GetWithDefaultAlertsSettings returns the WithDefaultAlertsSettings field value if set, zero value otherwise
func (o *Group) GetWithDefaultAlertsSettings() bool {
	if o == nil || IsNil(o.WithDefaultAlertsSettings) {
		var ret bool
		return ret
	}
	return *o.WithDefaultAlertsSettings
}

// GetWithDefaultAlertsSettingsOk returns a tuple with the WithDefaultAlertsSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Group) GetWithDefaultAlertsSettingsOk() (*bool, bool) {
	if o == nil || IsNil(o.WithDefaultAlertsSettings) {
		return nil, false
	}

	return o.WithDefaultAlertsSettings, true
}

// HasWithDefaultAlertsSettings returns a boolean if a field has been set.
func (o *Group) HasWithDefaultAlertsSettings() bool {
	if o != nil && !IsNil(o.WithDefaultAlertsSettings) {
		return true
	}

	return false
}

// SetWithDefaultAlertsSettings gets a reference to the given bool and assigns it to the WithDefaultAlertsSettings field.
func (o *Group) SetWithDefaultAlertsSettings(v bool) {
	o.WithDefaultAlertsSettings = &v
}
