// Code based on the AtlasAPI V2 OpenAPI file

package admin

// TriggerIngestionPipelineRequest struct for TriggerIngestionPipelineRequest
type TriggerIngestionPipelineRequest struct {
	DatasetRetentionPolicy *DatasetRetentionPolicy `json:"datasetRetentionPolicy,omitempty"`
	// Unique 24-hexadecimal character string that identifies the snapshot.
	// Write only field.
	SnapshotId string `json:"snapshotId"`
}

// NewTriggerIngestionPipelineRequest instantiates a new TriggerIngestionPipelineRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTriggerIngestionPipelineRequest(snapshotId string) *TriggerIngestionPipelineRequest {
	this := TriggerIngestionPipelineRequest{}
	this.SnapshotId = snapshotId
	return &this
}

// NewTriggerIngestionPipelineRequestWithDefaults instantiates a new TriggerIngestionPipelineRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTriggerIngestionPipelineRequestWithDefaults() *TriggerIngestionPipelineRequest {
	this := TriggerIngestionPipelineRequest{}
	return &this
}

// GetDatasetRetentionPolicy returns the DatasetRetentionPolicy field value if set, zero value otherwise
func (o *TriggerIngestionPipelineRequest) GetDatasetRetentionPolicy() DatasetRetentionPolicy {
	if o == nil || IsNil(o.DatasetRetentionPolicy) {
		var ret DatasetRetentionPolicy
		return ret
	}
	return *o.DatasetRetentionPolicy
}

// GetDatasetRetentionPolicyOk returns a tuple with the DatasetRetentionPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TriggerIngestionPipelineRequest) GetDatasetRetentionPolicyOk() (*DatasetRetentionPolicy, bool) {
	if o == nil || IsNil(o.DatasetRetentionPolicy) {
		return nil, false
	}

	return o.DatasetRetentionPolicy, true
}

// HasDatasetRetentionPolicy returns a boolean if a field has been set.
func (o *TriggerIngestionPipelineRequest) HasDatasetRetentionPolicy() bool {
	if o != nil && !IsNil(o.DatasetRetentionPolicy) {
		return true
	}

	return false
}

// SetDatasetRetentionPolicy gets a reference to the given DatasetRetentionPolicy and assigns it to the DatasetRetentionPolicy field.
func (o *TriggerIngestionPipelineRequest) SetDatasetRetentionPolicy(v DatasetRetentionPolicy) {
	o.DatasetRetentionPolicy = &v
}

// GetSnapshotId returns the SnapshotId field value
func (o *TriggerIngestionPipelineRequest) GetSnapshotId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value
// and a boolean to check if the value has been set.
func (o *TriggerIngestionPipelineRequest) GetSnapshotIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SnapshotId, true
}

// SetSnapshotId sets field value
func (o *TriggerIngestionPipelineRequest) SetSnapshotId(v string) {
	o.SnapshotId = v
}
