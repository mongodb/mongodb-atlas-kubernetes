// Code based on the AtlasAPI V2 OpenAPI file

package admin

// IngestionSink Ingestion destination of a Data Lake Pipeline.
type IngestionSink struct {
	// Type of ingestion destination of this Data Lake Pipeline.
	// Read only field.
	Type *string `json:"type,omitempty"`
	// Target cloud provider for this Data Lake Pipeline.
	MetadataProvider *string `json:"metadataProvider,omitempty"`
	// Target cloud provider region for this Data Lake Pipeline.
	MetadataRegion *string `json:"metadataRegion,omitempty"`
	// Ordered fields used to physically organize data in the destination.
	PartitionFields *[]DataLakePipelinesPartitionField `json:"partitionFields,omitempty"`
}

// NewIngestionSink instantiates a new IngestionSink object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIngestionSink() *IngestionSink {
	this := IngestionSink{}
	return &this
}

// NewIngestionSinkWithDefaults instantiates a new IngestionSink object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIngestionSinkWithDefaults() *IngestionSink {
	this := IngestionSink{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise
func (o *IngestionSink) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSink) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *IngestionSink) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *IngestionSink) SetType(v string) {
	o.Type = &v
}

// GetMetadataProvider returns the MetadataProvider field value if set, zero value otherwise
func (o *IngestionSink) GetMetadataProvider() string {
	if o == nil || IsNil(o.MetadataProvider) {
		var ret string
		return ret
	}
	return *o.MetadataProvider
}

// GetMetadataProviderOk returns a tuple with the MetadataProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSink) GetMetadataProviderOk() (*string, bool) {
	if o == nil || IsNil(o.MetadataProvider) {
		return nil, false
	}

	return o.MetadataProvider, true
}

// HasMetadataProvider returns a boolean if a field has been set.
func (o *IngestionSink) HasMetadataProvider() bool {
	if o != nil && !IsNil(o.MetadataProvider) {
		return true
	}

	return false
}

// SetMetadataProvider gets a reference to the given string and assigns it to the MetadataProvider field.
func (o *IngestionSink) SetMetadataProvider(v string) {
	o.MetadataProvider = &v
}

// GetMetadataRegion returns the MetadataRegion field value if set, zero value otherwise
func (o *IngestionSink) GetMetadataRegion() string {
	if o == nil || IsNil(o.MetadataRegion) {
		var ret string
		return ret
	}
	return *o.MetadataRegion
}

// GetMetadataRegionOk returns a tuple with the MetadataRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSink) GetMetadataRegionOk() (*string, bool) {
	if o == nil || IsNil(o.MetadataRegion) {
		return nil, false
	}

	return o.MetadataRegion, true
}

// HasMetadataRegion returns a boolean if a field has been set.
func (o *IngestionSink) HasMetadataRegion() bool {
	if o != nil && !IsNil(o.MetadataRegion) {
		return true
	}

	return false
}

// SetMetadataRegion gets a reference to the given string and assigns it to the MetadataRegion field.
func (o *IngestionSink) SetMetadataRegion(v string) {
	o.MetadataRegion = &v
}

// GetPartitionFields returns the PartitionFields field value if set, zero value otherwise
func (o *IngestionSink) GetPartitionFields() []DataLakePipelinesPartitionField {
	if o == nil || IsNil(o.PartitionFields) {
		var ret []DataLakePipelinesPartitionField
		return ret
	}
	return *o.PartitionFields
}

// GetPartitionFieldsOk returns a tuple with the PartitionFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSink) GetPartitionFieldsOk() (*[]DataLakePipelinesPartitionField, bool) {
	if o == nil || IsNil(o.PartitionFields) {
		return nil, false
	}

	return o.PartitionFields, true
}

// HasPartitionFields returns a boolean if a field has been set.
func (o *IngestionSink) HasPartitionFields() bool {
	if o != nil && !IsNil(o.PartitionFields) {
		return true
	}

	return false
}

// SetPartitionFields gets a reference to the given []DataLakePipelinesPartitionField and assigns it to the PartitionFields field.
func (o *IngestionSink) SetPartitionFields(v []DataLakePipelinesPartitionField) {
	o.PartitionFields = &v
}
