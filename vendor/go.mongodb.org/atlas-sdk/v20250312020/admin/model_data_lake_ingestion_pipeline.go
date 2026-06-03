// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DataLakeIngestionPipeline Details of a Data Lake Pipeline.
type DataLakeIngestionPipeline struct {
	// Unique 24-hexadecimal digit string that identifies the Data Lake Pipeline.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Timestamp that indicates when the Data Lake Pipeline was created. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedDate            *time.Time              `json:"createdDate,omitempty"`
	DatasetRetentionPolicy *DatasetRetentionPolicy `json:"datasetRetentionPolicy,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the group.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Timestamp that indicates the last time that the Data Lake Pipeline was updated. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastUpdatedDate *time.Time `json:"lastUpdatedDate,omitempty"`
	// Name of this Data Lake Pipeline.
	Name   *string          `json:"name,omitempty"`
	Sink   *IngestionSink   `json:"sink,omitempty"`
	Source *IngestionSource `json:"source,omitempty"`
	// State of this Data Lake Pipeline.
	// Read only field.
	State *string `json:"state,omitempty"`
	// Fields to be excluded for this Data Lake Pipeline.
	Transformations *[]FieldTransformation `json:"transformations,omitempty"`
}

// NewDataLakeIngestionPipeline instantiates a new DataLakeIngestionPipeline object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeIngestionPipeline() *DataLakeIngestionPipeline {
	this := DataLakeIngestionPipeline{}
	return &this
}

// NewDataLakeIngestionPipelineWithDefaults instantiates a new DataLakeIngestionPipeline object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeIngestionPipelineWithDefaults() *DataLakeIngestionPipeline {
	this := DataLakeIngestionPipeline{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *DataLakeIngestionPipeline) SetId(v string) {
	o.Id = &v
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return *o.CreatedDate
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}

	return o.CreatedDate, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *DataLakeIngestionPipeline) SetCreatedDate(v time.Time) {
	o.CreatedDate = &v
}

// GetDatasetRetentionPolicy returns the DatasetRetentionPolicy field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetDatasetRetentionPolicy() DatasetRetentionPolicy {
	if o == nil || IsNil(o.DatasetRetentionPolicy) {
		var ret DatasetRetentionPolicy
		return ret
	}
	return *o.DatasetRetentionPolicy
}

// GetDatasetRetentionPolicyOk returns a tuple with the DatasetRetentionPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetDatasetRetentionPolicyOk() (*DatasetRetentionPolicy, bool) {
	if o == nil || IsNil(o.DatasetRetentionPolicy) {
		return nil, false
	}

	return o.DatasetRetentionPolicy, true
}

// HasDatasetRetentionPolicy returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasDatasetRetentionPolicy() bool {
	if o != nil && !IsNil(o.DatasetRetentionPolicy) {
		return true
	}

	return false
}

// SetDatasetRetentionPolicy gets a reference to the given DatasetRetentionPolicy and assigns it to the DatasetRetentionPolicy field.
func (o *DataLakeIngestionPipeline) SetDatasetRetentionPolicy(v DatasetRetentionPolicy) {
	o.DatasetRetentionPolicy = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *DataLakeIngestionPipeline) SetGroupId(v string) {
	o.GroupId = &v
}

// GetLastUpdatedDate returns the LastUpdatedDate field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetLastUpdatedDate() time.Time {
	if o == nil || IsNil(o.LastUpdatedDate) {
		var ret time.Time
		return ret
	}
	return *o.LastUpdatedDate
}

// GetLastUpdatedDateOk returns a tuple with the LastUpdatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetLastUpdatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastUpdatedDate) {
		return nil, false
	}

	return o.LastUpdatedDate, true
}

// HasLastUpdatedDate returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasLastUpdatedDate() bool {
	if o != nil && !IsNil(o.LastUpdatedDate) {
		return true
	}

	return false
}

// SetLastUpdatedDate gets a reference to the given time.Time and assigns it to the LastUpdatedDate field.
func (o *DataLakeIngestionPipeline) SetLastUpdatedDate(v time.Time) {
	o.LastUpdatedDate = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataLakeIngestionPipeline) SetName(v string) {
	o.Name = &v
}

// GetSink returns the Sink field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetSink() IngestionSink {
	if o == nil || IsNil(o.Sink) {
		var ret IngestionSink
		return ret
	}
	return *o.Sink
}

// GetSinkOk returns a tuple with the Sink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetSinkOk() (*IngestionSink, bool) {
	if o == nil || IsNil(o.Sink) {
		return nil, false
	}

	return o.Sink, true
}

// HasSink returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasSink() bool {
	if o != nil && !IsNil(o.Sink) {
		return true
	}

	return false
}

// SetSink gets a reference to the given IngestionSink and assigns it to the Sink field.
func (o *DataLakeIngestionPipeline) SetSink(v IngestionSink) {
	o.Sink = &v
}

// GetSource returns the Source field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetSource() IngestionSource {
	if o == nil || IsNil(o.Source) {
		var ret IngestionSource
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetSourceOk() (*IngestionSource, bool) {
	if o == nil || IsNil(o.Source) {
		return nil, false
	}

	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasSource() bool {
	if o != nil && !IsNil(o.Source) {
		return true
	}

	return false
}

// SetSource gets a reference to the given IngestionSource and assigns it to the Source field.
func (o *DataLakeIngestionPipeline) SetSource(v IngestionSource) {
	o.Source = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *DataLakeIngestionPipeline) SetState(v string) {
	o.State = &v
}

// GetTransformations returns the Transformations field value if set, zero value otherwise
func (o *DataLakeIngestionPipeline) GetTransformations() []FieldTransformation {
	if o == nil || IsNil(o.Transformations) {
		var ret []FieldTransformation
		return ret
	}
	return *o.Transformations
}

// GetTransformationsOk returns a tuple with the Transformations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeIngestionPipeline) GetTransformationsOk() (*[]FieldTransformation, bool) {
	if o == nil || IsNil(o.Transformations) {
		return nil, false
	}

	return o.Transformations, true
}

// HasTransformations returns a boolean if a field has been set.
func (o *DataLakeIngestionPipeline) HasTransformations() bool {
	if o != nil && !IsNil(o.Transformations) {
		return true
	}

	return false
}

// SetTransformations gets a reference to the given []FieldTransformation and assigns it to the Transformations field.
func (o *DataLakeIngestionPipeline) SetTransformations(v []FieldTransformation) {
	o.Transformations = &v
}
