// Code based on the AtlasAPI V2 OpenAPI file

package admin

// MeasurementDiskPartition struct for MeasurementDiskPartition
type MeasurementDiskPartition struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label of the disk or partition to which the measurements apply.
	// Read only field.
	PartitionName *string `json:"partitionName,omitempty"`
}

// NewMeasurementDiskPartition instantiates a new MeasurementDiskPartition object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMeasurementDiskPartition() *MeasurementDiskPartition {
	this := MeasurementDiskPartition{}
	return &this
}

// NewMeasurementDiskPartitionWithDefaults instantiates a new MeasurementDiskPartition object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMeasurementDiskPartitionWithDefaults() *MeasurementDiskPartition {
	this := MeasurementDiskPartition{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *MeasurementDiskPartition) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementDiskPartition) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *MeasurementDiskPartition) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *MeasurementDiskPartition) SetLinks(v []Link) {
	o.Links = &v
}

// GetPartitionName returns the PartitionName field value if set, zero value otherwise
func (o *MeasurementDiskPartition) GetPartitionName() string {
	if o == nil || IsNil(o.PartitionName) {
		var ret string
		return ret
	}
	return *o.PartitionName
}

// GetPartitionNameOk returns a tuple with the PartitionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementDiskPartition) GetPartitionNameOk() (*string, bool) {
	if o == nil || IsNil(o.PartitionName) {
		return nil, false
	}

	return o.PartitionName, true
}

// HasPartitionName returns a boolean if a field has been set.
func (o *MeasurementDiskPartition) HasPartitionName() bool {
	if o != nil && !IsNil(o.PartitionName) {
		return true
	}

	return false
}

// SetPartitionName gets a reference to the given string and assigns it to the PartitionName field.
func (o *MeasurementDiskPartition) SetPartitionName(v string) {
	o.PartitionName = &v
}
