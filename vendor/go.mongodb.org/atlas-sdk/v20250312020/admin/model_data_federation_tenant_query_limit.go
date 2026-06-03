// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DataFederationTenantQueryLimit Details of a tenant-level query limit for Data Federation. Query limit is the limit on the amount of usage during a time period based on cost.
type DataFederationTenantQueryLimit struct {
	// Amount that indicates the current usage of the limit.
	// Read only field.
	CurrentUsage *int64 `json:"currentUsage,omitempty"`
	// Default value of the limit.
	// Read only field.
	DefaultLimit *int64 `json:"defaultLimit,omitempty"`
	// Only used for Data Federation limits. Timestamp that indicates when this usage limit was last modified. This field uses the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastModifiedDate *time.Time `json:"lastModifiedDate,omitempty"`
	// Maximum value of the limit.
	// Read only field.
	MaximumLimit *int64 `json:"maximumLimit,omitempty"`
	// Human-readable label that identifies the user-managed limit to modify.
	// Read only field.
	Name string `json:"name"`
	// Only used for Data Federation limits. Action to take when the usage limit is exceeded. If limit span is set to QUERY, this is ignored because MongoDB Cloud stops the query when it exceeds the usage limit.
	OverrunPolicy *string `json:"overrunPolicy,omitempty"`
	// Human-readable label that identifies the Federated Database Instance. If specified, the usage limit is for the specified federated database instance only. If omitted, the usage limit is for all federated database instances in the project.
	// Read only field.
	TenantName *string `json:"tenantName,omitempty"`
	// Amount to set the limit to.
	Value int64 `json:"value"`
}

// NewDataFederationTenantQueryLimit instantiates a new DataFederationTenantQueryLimit object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataFederationTenantQueryLimit(name string, value int64) *DataFederationTenantQueryLimit {
	this := DataFederationTenantQueryLimit{}
	this.Name = name
	this.Value = value
	return &this
}

// NewDataFederationTenantQueryLimitWithDefaults instantiates a new DataFederationTenantQueryLimit object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataFederationTenantQueryLimitWithDefaults() *DataFederationTenantQueryLimit {
	this := DataFederationTenantQueryLimit{}
	return &this
}

// GetCurrentUsage returns the CurrentUsage field value if set, zero value otherwise
func (o *DataFederationTenantQueryLimit) GetCurrentUsage() int64 {
	if o == nil || IsNil(o.CurrentUsage) {
		var ret int64
		return ret
	}
	return *o.CurrentUsage
}

// GetCurrentUsageOk returns a tuple with the CurrentUsage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetCurrentUsageOk() (*int64, bool) {
	if o == nil || IsNil(o.CurrentUsage) {
		return nil, false
	}

	return o.CurrentUsage, true
}

// HasCurrentUsage returns a boolean if a field has been set.
func (o *DataFederationTenantQueryLimit) HasCurrentUsage() bool {
	if o != nil && !IsNil(o.CurrentUsage) {
		return true
	}

	return false
}

// SetCurrentUsage gets a reference to the given int64 and assigns it to the CurrentUsage field.
func (o *DataFederationTenantQueryLimit) SetCurrentUsage(v int64) {
	o.CurrentUsage = &v
}

// GetDefaultLimit returns the DefaultLimit field value if set, zero value otherwise
func (o *DataFederationTenantQueryLimit) GetDefaultLimit() int64 {
	if o == nil || IsNil(o.DefaultLimit) {
		var ret int64
		return ret
	}
	return *o.DefaultLimit
}

// GetDefaultLimitOk returns a tuple with the DefaultLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetDefaultLimitOk() (*int64, bool) {
	if o == nil || IsNil(o.DefaultLimit) {
		return nil, false
	}

	return o.DefaultLimit, true
}

// HasDefaultLimit returns a boolean if a field has been set.
func (o *DataFederationTenantQueryLimit) HasDefaultLimit() bool {
	if o != nil && !IsNil(o.DefaultLimit) {
		return true
	}

	return false
}

// SetDefaultLimit gets a reference to the given int64 and assigns it to the DefaultLimit field.
func (o *DataFederationTenantQueryLimit) SetDefaultLimit(v int64) {
	o.DefaultLimit = &v
}

// GetLastModifiedDate returns the LastModifiedDate field value if set, zero value otherwise
func (o *DataFederationTenantQueryLimit) GetLastModifiedDate() time.Time {
	if o == nil || IsNil(o.LastModifiedDate) {
		var ret time.Time
		return ret
	}
	return *o.LastModifiedDate
}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastModifiedDate) {
		return nil, false
	}

	return o.LastModifiedDate, true
}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *DataFederationTenantQueryLimit) HasLastModifiedDate() bool {
	if o != nil && !IsNil(o.LastModifiedDate) {
		return true
	}

	return false
}

// SetLastModifiedDate gets a reference to the given time.Time and assigns it to the LastModifiedDate field.
func (o *DataFederationTenantQueryLimit) SetLastModifiedDate(v time.Time) {
	o.LastModifiedDate = &v
}

// GetMaximumLimit returns the MaximumLimit field value if set, zero value otherwise
func (o *DataFederationTenantQueryLimit) GetMaximumLimit() int64 {
	if o == nil || IsNil(o.MaximumLimit) {
		var ret int64
		return ret
	}
	return *o.MaximumLimit
}

// GetMaximumLimitOk returns a tuple with the MaximumLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetMaximumLimitOk() (*int64, bool) {
	if o == nil || IsNil(o.MaximumLimit) {
		return nil, false
	}

	return o.MaximumLimit, true
}

// HasMaximumLimit returns a boolean if a field has been set.
func (o *DataFederationTenantQueryLimit) HasMaximumLimit() bool {
	if o != nil && !IsNil(o.MaximumLimit) {
		return true
	}

	return false
}

// SetMaximumLimit gets a reference to the given int64 and assigns it to the MaximumLimit field.
func (o *DataFederationTenantQueryLimit) SetMaximumLimit(v int64) {
	o.MaximumLimit = &v
}

// GetName returns the Name field value
func (o *DataFederationTenantQueryLimit) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *DataFederationTenantQueryLimit) SetName(v string) {
	o.Name = v
}

// GetOverrunPolicy returns the OverrunPolicy field value if set, zero value otherwise
func (o *DataFederationTenantQueryLimit) GetOverrunPolicy() string {
	if o == nil || IsNil(o.OverrunPolicy) {
		var ret string
		return ret
	}
	return *o.OverrunPolicy
}

// GetOverrunPolicyOk returns a tuple with the OverrunPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetOverrunPolicyOk() (*string, bool) {
	if o == nil || IsNil(o.OverrunPolicy) {
		return nil, false
	}

	return o.OverrunPolicy, true
}

// HasOverrunPolicy returns a boolean if a field has been set.
func (o *DataFederationTenantQueryLimit) HasOverrunPolicy() bool {
	if o != nil && !IsNil(o.OverrunPolicy) {
		return true
	}

	return false
}

// SetOverrunPolicy gets a reference to the given string and assigns it to the OverrunPolicy field.
func (o *DataFederationTenantQueryLimit) SetOverrunPolicy(v string) {
	o.OverrunPolicy = &v
}

// GetTenantName returns the TenantName field value if set, zero value otherwise
func (o *DataFederationTenantQueryLimit) GetTenantName() string {
	if o == nil || IsNil(o.TenantName) {
		var ret string
		return ret
	}
	return *o.TenantName
}

// GetTenantNameOk returns a tuple with the TenantName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetTenantNameOk() (*string, bool) {
	if o == nil || IsNil(o.TenantName) {
		return nil, false
	}

	return o.TenantName, true
}

// HasTenantName returns a boolean if a field has been set.
func (o *DataFederationTenantQueryLimit) HasTenantName() bool {
	if o != nil && !IsNil(o.TenantName) {
		return true
	}

	return false
}

// SetTenantName gets a reference to the given string and assigns it to the TenantName field.
func (o *DataFederationTenantQueryLimit) SetTenantName(v string) {
	o.TenantName = &v
}

// GetValue returns the Value field value
func (o *DataFederationTenantQueryLimit) GetValue() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *DataFederationTenantQueryLimit) GetValueOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *DataFederationTenantQueryLimit) SetValue(v int64) {
	o.Value = v
}
