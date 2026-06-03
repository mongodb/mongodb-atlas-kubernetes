// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DatasetRetentionPolicy Dataset Retention Policy for a Scheduled Data Lake Pipeline.
type DatasetRetentionPolicy struct {
	// Date when retention policy was last modified. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastModifiedDate *time.Time `json:"lastModifiedDate,omitempty"`
	// Quantity of time in which the Data Lake Pipeline measures dataset retention.
	Units string `json:"units"`
	// Number that indicates the amount of days, weeks, or months that the Data Lake Pipeline will retain datasets.
	Value int `json:"value"`
}

// NewDatasetRetentionPolicy instantiates a new DatasetRetentionPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatasetRetentionPolicy(units string, value int) *DatasetRetentionPolicy {
	this := DatasetRetentionPolicy{}
	this.Units = units
	this.Value = value
	return &this
}

// NewDatasetRetentionPolicyWithDefaults instantiates a new DatasetRetentionPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatasetRetentionPolicyWithDefaults() *DatasetRetentionPolicy {
	this := DatasetRetentionPolicy{}
	return &this
}

// GetLastModifiedDate returns the LastModifiedDate field value if set, zero value otherwise
func (o *DatasetRetentionPolicy) GetLastModifiedDate() time.Time {
	if o == nil || IsNil(o.LastModifiedDate) {
		var ret time.Time
		return ret
	}
	return *o.LastModifiedDate
}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatasetRetentionPolicy) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastModifiedDate) {
		return nil, false
	}

	return o.LastModifiedDate, true
}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *DatasetRetentionPolicy) HasLastModifiedDate() bool {
	if o != nil && !IsNil(o.LastModifiedDate) {
		return true
	}

	return false
}

// SetLastModifiedDate gets a reference to the given time.Time and assigns it to the LastModifiedDate field.
func (o *DatasetRetentionPolicy) SetLastModifiedDate(v time.Time) {
	o.LastModifiedDate = &v
}

// GetUnits returns the Units field value
func (o *DatasetRetentionPolicy) GetUnits() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Units
}

// GetUnitsOk returns a tuple with the Units field value
// and a boolean to check if the value has been set.
func (o *DatasetRetentionPolicy) GetUnitsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Units, true
}

// SetUnits sets field value
func (o *DatasetRetentionPolicy) SetUnits(v string) {
	o.Units = v
}

// GetValue returns the Value field value
func (o *DatasetRetentionPolicy) GetValue() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *DatasetRetentionPolicy) GetValueOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *DatasetRetentionPolicy) SetValue(v int) {
	o.Value = v
}
