// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// CloudProviderAccessGCPServiceAccountAllOf struct for CloudProviderAccessGCPServiceAccountAllOf
type CloudProviderAccessGCPServiceAccountAllOf struct {
	// Date and time when this Google Service Account was created. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	// List that contains application features associated with this Google Service Account.
	// Read only field.
	FeatureUsages *[]CloudProviderAccessFeatureUsage `json:"featureUsages,omitempty"`
	// Email address for the Google Service Account created by Atlas.
	GcpServiceAccountForAtlas *string `json:"gcpServiceAccountForAtlas,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the role.
	// Read only field.
	RoleId *string `json:"roleId,omitempty"`
	// Provision status of the service account.
	// Read only field.
	Status *string `json:"status,omitempty"`
}

// NewCloudProviderAccessGCPServiceAccountAllOf instantiates a new CloudProviderAccessGCPServiceAccountAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderAccessGCPServiceAccountAllOf() *CloudProviderAccessGCPServiceAccountAllOf {
	this := CloudProviderAccessGCPServiceAccountAllOf{}
	return &this
}

// NewCloudProviderAccessGCPServiceAccountAllOfWithDefaults instantiates a new CloudProviderAccessGCPServiceAccountAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderAccessGCPServiceAccountAllOfWithDefaults() *CloudProviderAccessGCPServiceAccountAllOf {
	this := CloudProviderAccessGCPServiceAccountAllOf{}
	return &this
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return *o.CreatedDate
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}

	return o.CreatedDate, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *CloudProviderAccessGCPServiceAccountAllOf) SetCreatedDate(v time.Time) {
	o.CreatedDate = &v
}

// GetFeatureUsages returns the FeatureUsages field value if set, zero value otherwise
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetFeatureUsages() []CloudProviderAccessFeatureUsage {
	if o == nil || IsNil(o.FeatureUsages) {
		var ret []CloudProviderAccessFeatureUsage
		return ret
	}
	return *o.FeatureUsages
}

// GetFeatureUsagesOk returns a tuple with the FeatureUsages field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetFeatureUsagesOk() (*[]CloudProviderAccessFeatureUsage, bool) {
	if o == nil || IsNil(o.FeatureUsages) {
		return nil, false
	}

	return o.FeatureUsages, true
}

// HasFeatureUsages returns a boolean if a field has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) HasFeatureUsages() bool {
	if o != nil && !IsNil(o.FeatureUsages) {
		return true
	}

	return false
}

// SetFeatureUsages gets a reference to the given []CloudProviderAccessFeatureUsage and assigns it to the FeatureUsages field.
func (o *CloudProviderAccessGCPServiceAccountAllOf) SetFeatureUsages(v []CloudProviderAccessFeatureUsage) {
	o.FeatureUsages = &v
}

// GetGcpServiceAccountForAtlas returns the GcpServiceAccountForAtlas field value if set, zero value otherwise
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetGcpServiceAccountForAtlas() string {
	if o == nil || IsNil(o.GcpServiceAccountForAtlas) {
		var ret string
		return ret
	}
	return *o.GcpServiceAccountForAtlas
}

// GetGcpServiceAccountForAtlasOk returns a tuple with the GcpServiceAccountForAtlas field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetGcpServiceAccountForAtlasOk() (*string, bool) {
	if o == nil || IsNil(o.GcpServiceAccountForAtlas) {
		return nil, false
	}

	return o.GcpServiceAccountForAtlas, true
}

// HasGcpServiceAccountForAtlas returns a boolean if a field has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) HasGcpServiceAccountForAtlas() bool {
	if o != nil && !IsNil(o.GcpServiceAccountForAtlas) {
		return true
	}

	return false
}

// SetGcpServiceAccountForAtlas gets a reference to the given string and assigns it to the GcpServiceAccountForAtlas field.
func (o *CloudProviderAccessGCPServiceAccountAllOf) SetGcpServiceAccountForAtlas(v string) {
	o.GcpServiceAccountForAtlas = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *CloudProviderAccessGCPServiceAccountAllOf) SetRoleId(v string) {
	o.RoleId = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *CloudProviderAccessGCPServiceAccountAllOf) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *CloudProviderAccessGCPServiceAccountAllOf) SetStatus(v string) {
	o.Status = &v
}
