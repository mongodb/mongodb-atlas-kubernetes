// Code based on the AtlasAPI V2 OpenAPI file

package admin

// HardwareSpec20240805 Hardware specifications for all electable nodes deployed in the region. Electable nodes can become the primary and can enable local reads. If you don't specify this option, MongoDB Cloud deploys no electable nodes to the region.
type HardwareSpec20240805 struct {
	// Storage capacity of instance data volumes expressed in gigabytes. Increase this number to add capacity.   This value must be equal for all shards and node types.   This value is not configurable on M0/M2/M5 clusters.   MongoDB Cloud requires this parameter if you set `replicationSpecs`.   If you specify a disk size below the minimum (10 GB), this parameter defaults to the minimum disk size value.    Storage charge calculations depend on whether you choose the default value or a custom value.   The maximum value for disk storage cannot exceed 50 times the maximum RAM for the selected cluster. If you require more storage space, consider upgrading your cluster to a higher tier.
	DiskSizeGB *float64 `json:"diskSizeGB,omitempty"`
	// Target throughput desired for storage attached to your Azure-provisioned cluster. Change this parameter if you:  - set `replicationSpecs[n].regionConfigs[m].providerName` : `Azure`. - set `replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize` : `M40` or greater not including `Mxx_NVME` tiers.  The maximum input/output operations per second (IOPS) depend on the selected `.instanceSize` and `.diskSizeGB`. This parameter defaults to the cluster tier's standard IOPS value. Changing this value impacts cluster cost.
	DiskIOPS *int `json:"diskIOPS,omitempty"`
	// Target throughput desired for storage attached to this hardware. Only returned for Gen 2 instance sizes with Standard (GP3) volume type.
	// Read only field.
	DiskThroughput *int `json:"diskThroughput,omitempty"`
	// Type of storage you want to attach to your AWS-provisioned cluster.  - `STANDARD` volume types can't exceed the default input/output operations per second (IOPS) rate for the selected volume size.   - `PROVISIONED` volume types must fall within the allowable IOPS range for the selected volume size. You must set this value to (`PROVISIONED`) for NVMe clusters.
	EbsVolumeType *string `json:"ebsVolumeType,omitempty"`
	// Hardware specification for the instances in this M0/M2/M5 tier cluster.
	InstanceSize *string `json:"instanceSize,omitempty"`
	// Number of nodes of the given type for MongoDB Cloud to deploy to the region.
	NodeCount *int `json:"nodeCount,omitempty"`
	// The true tenant instance size. This is present to support backwards compatibility for deprecated provider types and/or instance sizes.
	// Read only field.
	EffectiveInstanceSize *string `json:"effectiveInstanceSize,omitempty"`
}

// NewHardwareSpec20240805 instantiates a new HardwareSpec20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHardwareSpec20240805() *HardwareSpec20240805 {
	this := HardwareSpec20240805{}
	var ebsVolumeType string = "STANDARD"
	this.EbsVolumeType = &ebsVolumeType
	return &this
}

// NewHardwareSpec20240805WithDefaults instantiates a new HardwareSpec20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHardwareSpec20240805WithDefaults() *HardwareSpec20240805 {
	this := HardwareSpec20240805{}
	var ebsVolumeType string = "STANDARD"
	this.EbsVolumeType = &ebsVolumeType
	return &this
}

// GetDiskSizeGB returns the DiskSizeGB field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetDiskSizeGB() float64 {
	if o == nil || IsNil(o.DiskSizeGB) {
		var ret float64
		return ret
	}
	return *o.DiskSizeGB
}

// GetDiskSizeGBOk returns a tuple with the DiskSizeGB field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetDiskSizeGBOk() (*float64, bool) {
	if o == nil || IsNil(o.DiskSizeGB) {
		return nil, false
	}

	return o.DiskSizeGB, true
}

// HasDiskSizeGB returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasDiskSizeGB() bool {
	if o != nil && !IsNil(o.DiskSizeGB) {
		return true
	}

	return false
}

// SetDiskSizeGB gets a reference to the given float64 and assigns it to the DiskSizeGB field.
func (o *HardwareSpec20240805) SetDiskSizeGB(v float64) {
	o.DiskSizeGB = &v
}

// GetDiskIOPS returns the DiskIOPS field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetDiskIOPS() int {
	if o == nil || IsNil(o.DiskIOPS) {
		var ret int
		return ret
	}
	return *o.DiskIOPS
}

// GetDiskIOPSOk returns a tuple with the DiskIOPS field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetDiskIOPSOk() (*int, bool) {
	if o == nil || IsNil(o.DiskIOPS) {
		return nil, false
	}

	return o.DiskIOPS, true
}

// HasDiskIOPS returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasDiskIOPS() bool {
	if o != nil && !IsNil(o.DiskIOPS) {
		return true
	}

	return false
}

// SetDiskIOPS gets a reference to the given int and assigns it to the DiskIOPS field.
func (o *HardwareSpec20240805) SetDiskIOPS(v int) {
	o.DiskIOPS = &v
}

// GetDiskThroughput returns the DiskThroughput field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetDiskThroughput() int {
	if o == nil || IsNil(o.DiskThroughput) {
		var ret int
		return ret
	}
	return *o.DiskThroughput
}

// GetDiskThroughputOk returns a tuple with the DiskThroughput field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetDiskThroughputOk() (*int, bool) {
	if o == nil || IsNil(o.DiskThroughput) {
		return nil, false
	}

	return o.DiskThroughput, true
}

// HasDiskThroughput returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasDiskThroughput() bool {
	if o != nil && !IsNil(o.DiskThroughput) {
		return true
	}

	return false
}

// SetDiskThroughput gets a reference to the given int and assigns it to the DiskThroughput field.
func (o *HardwareSpec20240805) SetDiskThroughput(v int) {
	o.DiskThroughput = &v
}

// GetEbsVolumeType returns the EbsVolumeType field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetEbsVolumeType() string {
	if o == nil || IsNil(o.EbsVolumeType) {
		var ret string
		return ret
	}
	return *o.EbsVolumeType
}

// GetEbsVolumeTypeOk returns a tuple with the EbsVolumeType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetEbsVolumeTypeOk() (*string, bool) {
	if o == nil || IsNil(o.EbsVolumeType) {
		return nil, false
	}

	return o.EbsVolumeType, true
}

// HasEbsVolumeType returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasEbsVolumeType() bool {
	if o != nil && !IsNil(o.EbsVolumeType) {
		return true
	}

	return false
}

// SetEbsVolumeType gets a reference to the given string and assigns it to the EbsVolumeType field.
func (o *HardwareSpec20240805) SetEbsVolumeType(v string) {
	o.EbsVolumeType = &v
}

// GetInstanceSize returns the InstanceSize field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetInstanceSize() string {
	if o == nil || IsNil(o.InstanceSize) {
		var ret string
		return ret
	}
	return *o.InstanceSize
}

// GetInstanceSizeOk returns a tuple with the InstanceSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetInstanceSizeOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceSize) {
		return nil, false
	}

	return o.InstanceSize, true
}

// HasInstanceSize returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasInstanceSize() bool {
	if o != nil && !IsNil(o.InstanceSize) {
		return true
	}

	return false
}

// SetInstanceSize gets a reference to the given string and assigns it to the InstanceSize field.
func (o *HardwareSpec20240805) SetInstanceSize(v string) {
	o.InstanceSize = &v
}

// GetNodeCount returns the NodeCount field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetNodeCount() int {
	if o == nil || IsNil(o.NodeCount) {
		var ret int
		return ret
	}
	return *o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetNodeCountOk() (*int, bool) {
	if o == nil || IsNil(o.NodeCount) {
		return nil, false
	}

	return o.NodeCount, true
}

// HasNodeCount returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasNodeCount() bool {
	if o != nil && !IsNil(o.NodeCount) {
		return true
	}

	return false
}

// SetNodeCount gets a reference to the given int and assigns it to the NodeCount field.
func (o *HardwareSpec20240805) SetNodeCount(v int) {
	o.NodeCount = &v
}

// GetEffectiveInstanceSize returns the EffectiveInstanceSize field value if set, zero value otherwise
func (o *HardwareSpec20240805) GetEffectiveInstanceSize() string {
	if o == nil || IsNil(o.EffectiveInstanceSize) {
		var ret string
		return ret
	}
	return *o.EffectiveInstanceSize
}

// GetEffectiveInstanceSizeOk returns a tuple with the EffectiveInstanceSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HardwareSpec20240805) GetEffectiveInstanceSizeOk() (*string, bool) {
	if o == nil || IsNil(o.EffectiveInstanceSize) {
		return nil, false
	}

	return o.EffectiveInstanceSize, true
}

// HasEffectiveInstanceSize returns a boolean if a field has been set.
func (o *HardwareSpec20240805) HasEffectiveInstanceSize() bool {
	if o != nil && !IsNil(o.EffectiveInstanceSize) {
		return true
	}

	return false
}

// SetEffectiveInstanceSize gets a reference to the given string and assigns it to the EffectiveInstanceSize field.
func (o *HardwareSpec20240805) SetEffectiveInstanceSize(v string) {
	o.EffectiveInstanceSize = &v
}
