// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// SampleDatasetStatus struct for SampleDatasetStatus
type SampleDatasetStatus struct {
	// Unique 24-hexadecimal character string that identifies this sample dataset.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Human-readable label that identifies the cluster into which you loaded the sample dataset.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Date and time when the sample dataset load job completed. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
	// Read only field.
	CompleteDate *time.Time `json:"completeDate,omitempty"`
	// Date and time when you started the sample dataset load job. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
	// Read only field.
	CreateDate *time.Time `json:"createDate,omitempty"`
	// Details of the error returned when MongoDB Cloud loads the sample dataset. This endpoint returns null if state has a value other than FAILED.
	// Read only field.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Status of the sample dataset load job.
	// Read only field.
	State *string `json:"state,omitempty"`
}

// NewSampleDatasetStatus instantiates a new SampleDatasetStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSampleDatasetStatus() *SampleDatasetStatus {
	this := SampleDatasetStatus{}
	return &this
}

// NewSampleDatasetStatusWithDefaults instantiates a new SampleDatasetStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSampleDatasetStatusWithDefaults() *SampleDatasetStatus {
	this := SampleDatasetStatus{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *SampleDatasetStatus) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SampleDatasetStatus) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *SampleDatasetStatus) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *SampleDatasetStatus) SetId(v string) {
	o.Id = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *SampleDatasetStatus) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SampleDatasetStatus) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *SampleDatasetStatus) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *SampleDatasetStatus) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCompleteDate returns the CompleteDate field value if set, zero value otherwise
func (o *SampleDatasetStatus) GetCompleteDate() time.Time {
	if o == nil || IsNil(o.CompleteDate) {
		var ret time.Time
		return ret
	}
	return *o.CompleteDate
}

// GetCompleteDateOk returns a tuple with the CompleteDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SampleDatasetStatus) GetCompleteDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CompleteDate) {
		return nil, false
	}

	return o.CompleteDate, true
}

// HasCompleteDate returns a boolean if a field has been set.
func (o *SampleDatasetStatus) HasCompleteDate() bool {
	if o != nil && !IsNil(o.CompleteDate) {
		return true
	}

	return false
}

// SetCompleteDate gets a reference to the given time.Time and assigns it to the CompleteDate field.
func (o *SampleDatasetStatus) SetCompleteDate(v time.Time) {
	o.CompleteDate = &v
}

// GetCreateDate returns the CreateDate field value if set, zero value otherwise
func (o *SampleDatasetStatus) GetCreateDate() time.Time {
	if o == nil || IsNil(o.CreateDate) {
		var ret time.Time
		return ret
	}
	return *o.CreateDate
}

// GetCreateDateOk returns a tuple with the CreateDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SampleDatasetStatus) GetCreateDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateDate) {
		return nil, false
	}

	return o.CreateDate, true
}

// HasCreateDate returns a boolean if a field has been set.
func (o *SampleDatasetStatus) HasCreateDate() bool {
	if o != nil && !IsNil(o.CreateDate) {
		return true
	}

	return false
}

// SetCreateDate gets a reference to the given time.Time and assigns it to the CreateDate field.
func (o *SampleDatasetStatus) SetCreateDate(v time.Time) {
	o.CreateDate = &v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise
func (o *SampleDatasetStatus) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SampleDatasetStatus) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}

	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *SampleDatasetStatus) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *SampleDatasetStatus) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *SampleDatasetStatus) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SampleDatasetStatus) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *SampleDatasetStatus) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *SampleDatasetStatus) SetState(v string) {
	o.State = &v
}
