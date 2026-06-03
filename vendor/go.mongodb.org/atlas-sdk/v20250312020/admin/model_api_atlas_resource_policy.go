// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ApiAtlasResourcePolicy struct for ApiAtlasResourcePolicy
type ApiAtlasResourcePolicy struct {
	CreatedByUser *ApiAtlasUserMetadata `json:"createdByUser,omitempty"`
	// Date and time in UTC when the atlas resource policy was created.
	// Read only field.
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	// Description of the atlas resource policy.
	// Read only field.
	Description *string `json:"description,omitempty"`
	// Unique 24-hexadecimal character string that identifies the atlas resource policy.
	// Read only field.
	Id                *string               `json:"id,omitempty"`
	LastUpdatedByUser *ApiAtlasUserMetadata `json:"lastUpdatedByUser,omitempty"`
	// Date and time in UTC when the atlas resource policy was last updated.
	// Read only field.
	LastUpdatedDate *time.Time `json:"lastUpdatedDate,omitempty"`
	// Human-readable label that describes the atlas resource policy.
	// Read only field.
	Name *string `json:"name,omitempty"`
	// Unique 24-hexadecimal character string that identifies the organization the atlas resource policy belongs to.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// List of policies that make up the atlas resource policy.
	// Read only field.
	Policies *[]ApiAtlasPolicy `json:"policies,omitempty"`
	// A string that identifies the version of the atlas resource policy.
	// Read only field.
	Version *string `json:"version,omitempty"`
}

// NewApiAtlasResourcePolicy instantiates a new ApiAtlasResourcePolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasResourcePolicy() *ApiAtlasResourcePolicy {
	this := ApiAtlasResourcePolicy{}
	return &this
}

// NewApiAtlasResourcePolicyWithDefaults instantiates a new ApiAtlasResourcePolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasResourcePolicyWithDefaults() *ApiAtlasResourcePolicy {
	this := ApiAtlasResourcePolicy{}
	return &this
}

// GetCreatedByUser returns the CreatedByUser field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetCreatedByUser() ApiAtlasUserMetadata {
	if o == nil || IsNil(o.CreatedByUser) {
		var ret ApiAtlasUserMetadata
		return ret
	}
	return *o.CreatedByUser
}

// GetCreatedByUserOk returns a tuple with the CreatedByUser field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetCreatedByUserOk() (*ApiAtlasUserMetadata, bool) {
	if o == nil || IsNil(o.CreatedByUser) {
		return nil, false
	}

	return o.CreatedByUser, true
}

// HasCreatedByUser returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasCreatedByUser() bool {
	if o != nil && !IsNil(o.CreatedByUser) {
		return true
	}

	return false
}

// SetCreatedByUser gets a reference to the given ApiAtlasUserMetadata and assigns it to the CreatedByUser field.
func (o *ApiAtlasResourcePolicy) SetCreatedByUser(v ApiAtlasUserMetadata) {
	o.CreatedByUser = &v
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return *o.CreatedDate
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}

	return o.CreatedDate, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *ApiAtlasResourcePolicy) SetCreatedDate(v time.Time) {
	o.CreatedDate = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *ApiAtlasResourcePolicy) SetDescription(v string) {
	o.Description = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiAtlasResourcePolicy) SetId(v string) {
	o.Id = &v
}

// GetLastUpdatedByUser returns the LastUpdatedByUser field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetLastUpdatedByUser() ApiAtlasUserMetadata {
	if o == nil || IsNil(o.LastUpdatedByUser) {
		var ret ApiAtlasUserMetadata
		return ret
	}
	return *o.LastUpdatedByUser
}

// GetLastUpdatedByUserOk returns a tuple with the LastUpdatedByUser field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetLastUpdatedByUserOk() (*ApiAtlasUserMetadata, bool) {
	if o == nil || IsNil(o.LastUpdatedByUser) {
		return nil, false
	}

	return o.LastUpdatedByUser, true
}

// HasLastUpdatedByUser returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasLastUpdatedByUser() bool {
	if o != nil && !IsNil(o.LastUpdatedByUser) {
		return true
	}

	return false
}

// SetLastUpdatedByUser gets a reference to the given ApiAtlasUserMetadata and assigns it to the LastUpdatedByUser field.
func (o *ApiAtlasResourcePolicy) SetLastUpdatedByUser(v ApiAtlasUserMetadata) {
	o.LastUpdatedByUser = &v
}

// GetLastUpdatedDate returns the LastUpdatedDate field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetLastUpdatedDate() time.Time {
	if o == nil || IsNil(o.LastUpdatedDate) {
		var ret time.Time
		return ret
	}
	return *o.LastUpdatedDate
}

// GetLastUpdatedDateOk returns a tuple with the LastUpdatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetLastUpdatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastUpdatedDate) {
		return nil, false
	}

	return o.LastUpdatedDate, true
}

// HasLastUpdatedDate returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasLastUpdatedDate() bool {
	if o != nil && !IsNil(o.LastUpdatedDate) {
		return true
	}

	return false
}

// SetLastUpdatedDate gets a reference to the given time.Time and assigns it to the LastUpdatedDate field.
func (o *ApiAtlasResourcePolicy) SetLastUpdatedDate(v time.Time) {
	o.LastUpdatedDate = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ApiAtlasResourcePolicy) SetName(v string) {
	o.Name = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *ApiAtlasResourcePolicy) SetOrgId(v string) {
	o.OrgId = &v
}

// GetPolicies returns the Policies field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetPolicies() []ApiAtlasPolicy {
	if o == nil || IsNil(o.Policies) {
		var ret []ApiAtlasPolicy
		return ret
	}
	return *o.Policies
}

// GetPoliciesOk returns a tuple with the Policies field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetPoliciesOk() (*[]ApiAtlasPolicy, bool) {
	if o == nil || IsNil(o.Policies) {
		return nil, false
	}

	return o.Policies, true
}

// HasPolicies returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasPolicies() bool {
	if o != nil && !IsNil(o.Policies) {
		return true
	}

	return false
}

// SetPolicies gets a reference to the given []ApiAtlasPolicy and assigns it to the Policies field.
func (o *ApiAtlasResourcePolicy) SetPolicies(v []ApiAtlasPolicy) {
	o.Policies = &v
}

// GetVersion returns the Version field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicy) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicy) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}

	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicy) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *ApiAtlasResourcePolicy) SetVersion(v string) {
	o.Version = &v
}
