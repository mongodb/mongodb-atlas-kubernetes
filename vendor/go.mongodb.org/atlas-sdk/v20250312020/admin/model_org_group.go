// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgGroup struct for OrgGroup
type OrgGroup struct {
	// Settings that describe the clusters in each project that the API key is authorized to view.
	// Read only field.
	Clusters *[]CloudCluster `json:"clusters,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Human-readable label that identifies the project.
	GroupName *string `json:"groupName,omitempty"`
	// Unique 24-hexadecimal character string that identifies the organization that contains the project.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// Human-readable label that identifies the organization that contains the project.
	OrgName *string `json:"orgName,omitempty"`
	// Human-readable label that indicates the plan type.
	// Read only field.
	PlanType *string `json:"planType,omitempty"`
	// List of human-readable labels that categorize the specified project. MongoDB Cloud returns an empty array.
	// Read only field.
	Tags *[]string `json:"tags,omitempty"`
}

// NewOrgGroup instantiates a new OrgGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgGroup() *OrgGroup {
	this := OrgGroup{}
	return &this
}

// NewOrgGroupWithDefaults instantiates a new OrgGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgGroupWithDefaults() *OrgGroup {
	this := OrgGroup{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise
func (o *OrgGroup) GetClusters() []CloudCluster {
	if o == nil || IsNil(o.Clusters) {
		var ret []CloudCluster
		return ret
	}
	return *o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetClustersOk() (*[]CloudCluster, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}

	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *OrgGroup) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []CloudCluster and assigns it to the Clusters field.
func (o *OrgGroup) SetClusters(v []CloudCluster) {
	o.Clusters = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *OrgGroup) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *OrgGroup) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *OrgGroup) SetGroupId(v string) {
	o.GroupId = &v
}

// GetGroupName returns the GroupName field value if set, zero value otherwise
func (o *OrgGroup) GetGroupName() string {
	if o == nil || IsNil(o.GroupName) {
		var ret string
		return ret
	}
	return *o.GroupName
}

// GetGroupNameOk returns a tuple with the GroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetGroupNameOk() (*string, bool) {
	if o == nil || IsNil(o.GroupName) {
		return nil, false
	}

	return o.GroupName, true
}

// HasGroupName returns a boolean if a field has been set.
func (o *OrgGroup) HasGroupName() bool {
	if o != nil && !IsNil(o.GroupName) {
		return true
	}

	return false
}

// SetGroupName gets a reference to the given string and assigns it to the GroupName field.
func (o *OrgGroup) SetGroupName(v string) {
	o.GroupName = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *OrgGroup) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *OrgGroup) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *OrgGroup) SetOrgId(v string) {
	o.OrgId = &v
}

// GetOrgName returns the OrgName field value if set, zero value otherwise
func (o *OrgGroup) GetOrgName() string {
	if o == nil || IsNil(o.OrgName) {
		var ret string
		return ret
	}
	return *o.OrgName
}

// GetOrgNameOk returns a tuple with the OrgName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetOrgNameOk() (*string, bool) {
	if o == nil || IsNil(o.OrgName) {
		return nil, false
	}

	return o.OrgName, true
}

// HasOrgName returns a boolean if a field has been set.
func (o *OrgGroup) HasOrgName() bool {
	if o != nil && !IsNil(o.OrgName) {
		return true
	}

	return false
}

// SetOrgName gets a reference to the given string and assigns it to the OrgName field.
func (o *OrgGroup) SetOrgName(v string) {
	o.OrgName = &v
}

// GetPlanType returns the PlanType field value if set, zero value otherwise
func (o *OrgGroup) GetPlanType() string {
	if o == nil || IsNil(o.PlanType) {
		var ret string
		return ret
	}
	return *o.PlanType
}

// GetPlanTypeOk returns a tuple with the PlanType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetPlanTypeOk() (*string, bool) {
	if o == nil || IsNil(o.PlanType) {
		return nil, false
	}

	return o.PlanType, true
}

// HasPlanType returns a boolean if a field has been set.
func (o *OrgGroup) HasPlanType() bool {
	if o != nil && !IsNil(o.PlanType) {
		return true
	}

	return false
}

// SetPlanType gets a reference to the given string and assigns it to the PlanType field.
func (o *OrgGroup) SetPlanType(v string) {
	o.PlanType = &v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *OrgGroup) GetTags() []string {
	if o == nil || IsNil(o.Tags) {
		var ret []string
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgGroup) GetTagsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *OrgGroup) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *OrgGroup) SetTags(v []string) {
	o.Tags = &v
}
