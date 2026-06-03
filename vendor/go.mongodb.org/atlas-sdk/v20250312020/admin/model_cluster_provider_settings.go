// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterProviderSettings Group of cloud provider settings that configure the provisioned MongoDB hosts.
type ClusterProviderSettings struct {
	ProviderName string                  `json:"providerName"`
	AutoScaling  *ClusterFreeAutoScaling `json:"autoScaling,omitempty"`
	// Maximum Disk Input/Output Operations per Second (IOPS) that the database host can perform.
	DiskIOPS *int `json:"diskIOPS,omitempty"`
	// Flag that indicates whether the Amazon Elastic Block Store (EBS) encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the cluster. Clusters always have this setting enabled.
	// Deprecated
	EncryptEBSVolume *bool `json:"encryptEBSVolume,omitempty"`
	// Cluster tier, with a default storage and memory capacity, that applies to all the data-bearing hosts in your cluster. You must set `providerSettings.providerName` to `FLEX` and specify the cloud service provider in `providerSettings.backingProviderName`.
	InstanceSizeName *string `json:"instanceSizeName,omitempty"`
	// Human-readable label that identifies the geographic location of your MongoDB cluster. The region you choose can affect network latency for clients accessing your databases. For a complete list of region names, see [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/#std-label-amazon-aws), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), and [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
	RegionName *string `json:"regionName,omitempty"`
	// Disk Input/Output Operations per Second (IOPS) setting for Amazon Web Services (AWS) storage that you configure only for AWS. Specify whether Disk Input/Output Operations per Second (IOPS) must not exceed the default Input/Output Operations per Second (IOPS) rate for the selected volume size (`STANDARD`), or must fall within the allowable Input/Output Operations per Second (IOPS) range for the selected volume size (`PROVISIONED`). You must set this value to (`PROVISIONED`) for NVMe clusters.
	VolumeType *string `json:"volumeType,omitempty"`
	// Disk type that corresponds to the host's root volume for Azure instances. If omitted, the default disk type for the selected `providerSettings.instanceSizeName` applies.
	DiskTypeName *string `json:"diskTypeName,omitempty"`
	// Cloud service provider on which MongoDB Cloud provisioned the multi-tenant host. The resource returns this parameter when `providerSettings.providerName` is `FLEX` and `providerSetting.instanceSizeName` is `FLEX`.
	BackingProviderName *string `json:"backingProviderName,omitempty"`
	// The true tenant instance size. This is present to support backwards compatibility for deprecated provider types and/or instance sizes.
	// Read only field.
	EffectiveInstanceSizeName *string `json:"effectiveInstanceSizeName,omitempty"`
}

// NewClusterProviderSettings instantiates a new ClusterProviderSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterProviderSettings(providerName string) *ClusterProviderSettings {
	this := ClusterProviderSettings{}
	this.ProviderName = providerName
	var encryptEBSVolume bool = true
	this.EncryptEBSVolume = &encryptEBSVolume
	return &this
}

// NewClusterProviderSettingsWithDefaults instantiates a new ClusterProviderSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterProviderSettingsWithDefaults() *ClusterProviderSettings {
	this := ClusterProviderSettings{}
	var encryptEBSVolume bool = true
	this.EncryptEBSVolume = &encryptEBSVolume
	return &this
}

// GetProviderName returns the ProviderName field value
func (o *ClusterProviderSettings) GetProviderName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetProviderNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderName, true
}

// SetProviderName sets field value
func (o *ClusterProviderSettings) SetProviderName(v string) {
	o.ProviderName = v
}

// GetAutoScaling returns the AutoScaling field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetAutoScaling() ClusterFreeAutoScaling {
	if o == nil || IsNil(o.AutoScaling) {
		var ret ClusterFreeAutoScaling
		return ret
	}
	return *o.AutoScaling
}

// GetAutoScalingOk returns a tuple with the AutoScaling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetAutoScalingOk() (*ClusterFreeAutoScaling, bool) {
	if o == nil || IsNil(o.AutoScaling) {
		return nil, false
	}

	return o.AutoScaling, true
}

// HasAutoScaling returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasAutoScaling() bool {
	if o != nil && !IsNil(o.AutoScaling) {
		return true
	}

	return false
}

// SetAutoScaling gets a reference to the given ClusterFreeAutoScaling and assigns it to the AutoScaling field.
func (o *ClusterProviderSettings) SetAutoScaling(v ClusterFreeAutoScaling) {
	o.AutoScaling = &v
}

// GetDiskIOPS returns the DiskIOPS field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetDiskIOPS() int {
	if o == nil || IsNil(o.DiskIOPS) {
		var ret int
		return ret
	}
	return *o.DiskIOPS
}

// GetDiskIOPSOk returns a tuple with the DiskIOPS field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetDiskIOPSOk() (*int, bool) {
	if o == nil || IsNil(o.DiskIOPS) {
		return nil, false
	}

	return o.DiskIOPS, true
}

// HasDiskIOPS returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasDiskIOPS() bool {
	if o != nil && !IsNil(o.DiskIOPS) {
		return true
	}

	return false
}

// SetDiskIOPS gets a reference to the given int and assigns it to the DiskIOPS field.
func (o *ClusterProviderSettings) SetDiskIOPS(v int) {
	o.DiskIOPS = &v
}

// GetEncryptEBSVolume returns the EncryptEBSVolume field value if set, zero value otherwise
// Deprecated
func (o *ClusterProviderSettings) GetEncryptEBSVolume() bool {
	if o == nil || IsNil(o.EncryptEBSVolume) {
		var ret bool
		return ret
	}
	return *o.EncryptEBSVolume
}

// GetEncryptEBSVolumeOk returns a tuple with the EncryptEBSVolume field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *ClusterProviderSettings) GetEncryptEBSVolumeOk() (*bool, bool) {
	if o == nil || IsNil(o.EncryptEBSVolume) {
		return nil, false
	}

	return o.EncryptEBSVolume, true
}

// HasEncryptEBSVolume returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasEncryptEBSVolume() bool {
	if o != nil && !IsNil(o.EncryptEBSVolume) {
		return true
	}

	return false
}

// SetEncryptEBSVolume gets a reference to the given bool and assigns it to the EncryptEBSVolume field.
// Deprecated
func (o *ClusterProviderSettings) SetEncryptEBSVolume(v bool) {
	o.EncryptEBSVolume = &v
}

// GetInstanceSizeName returns the InstanceSizeName field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetInstanceSizeName() string {
	if o == nil || IsNil(o.InstanceSizeName) {
		var ret string
		return ret
	}
	return *o.InstanceSizeName
}

// GetInstanceSizeNameOk returns a tuple with the InstanceSizeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetInstanceSizeNameOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceSizeName) {
		return nil, false
	}

	return o.InstanceSizeName, true
}

// HasInstanceSizeName returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasInstanceSizeName() bool {
	if o != nil && !IsNil(o.InstanceSizeName) {
		return true
	}

	return false
}

// SetInstanceSizeName gets a reference to the given string and assigns it to the InstanceSizeName field.
func (o *ClusterProviderSettings) SetInstanceSizeName(v string) {
	o.InstanceSizeName = &v
}

// GetRegionName returns the RegionName field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetRegionName() string {
	if o == nil || IsNil(o.RegionName) {
		var ret string
		return ret
	}
	return *o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionName) {
		return nil, false
	}

	return o.RegionName, true
}

// HasRegionName returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasRegionName() bool {
	if o != nil && !IsNil(o.RegionName) {
		return true
	}

	return false
}

// SetRegionName gets a reference to the given string and assigns it to the RegionName field.
func (o *ClusterProviderSettings) SetRegionName(v string) {
	o.RegionName = &v
}

// GetVolumeType returns the VolumeType field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetVolumeType() string {
	if o == nil || IsNil(o.VolumeType) {
		var ret string
		return ret
	}
	return *o.VolumeType
}

// GetVolumeTypeOk returns a tuple with the VolumeType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetVolumeTypeOk() (*string, bool) {
	if o == nil || IsNil(o.VolumeType) {
		return nil, false
	}

	return o.VolumeType, true
}

// HasVolumeType returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasVolumeType() bool {
	if o != nil && !IsNil(o.VolumeType) {
		return true
	}

	return false
}

// SetVolumeType gets a reference to the given string and assigns it to the VolumeType field.
func (o *ClusterProviderSettings) SetVolumeType(v string) {
	o.VolumeType = &v
}

// GetDiskTypeName returns the DiskTypeName field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetDiskTypeName() string {
	if o == nil || IsNil(o.DiskTypeName) {
		var ret string
		return ret
	}
	return *o.DiskTypeName
}

// GetDiskTypeNameOk returns a tuple with the DiskTypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetDiskTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.DiskTypeName) {
		return nil, false
	}

	return o.DiskTypeName, true
}

// HasDiskTypeName returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasDiskTypeName() bool {
	if o != nil && !IsNil(o.DiskTypeName) {
		return true
	}

	return false
}

// SetDiskTypeName gets a reference to the given string and assigns it to the DiskTypeName field.
func (o *ClusterProviderSettings) SetDiskTypeName(v string) {
	o.DiskTypeName = &v
}

// GetBackingProviderName returns the BackingProviderName field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetBackingProviderName() string {
	if o == nil || IsNil(o.BackingProviderName) {
		var ret string
		return ret
	}
	return *o.BackingProviderName
}

// GetBackingProviderNameOk returns a tuple with the BackingProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetBackingProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.BackingProviderName) {
		return nil, false
	}

	return o.BackingProviderName, true
}

// HasBackingProviderName returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasBackingProviderName() bool {
	if o != nil && !IsNil(o.BackingProviderName) {
		return true
	}

	return false
}

// SetBackingProviderName gets a reference to the given string and assigns it to the BackingProviderName field.
func (o *ClusterProviderSettings) SetBackingProviderName(v string) {
	o.BackingProviderName = &v
}

// GetEffectiveInstanceSizeName returns the EffectiveInstanceSizeName field value if set, zero value otherwise
func (o *ClusterProviderSettings) GetEffectiveInstanceSizeName() string {
	if o == nil || IsNil(o.EffectiveInstanceSizeName) {
		var ret string
		return ret
	}
	return *o.EffectiveInstanceSizeName
}

// GetEffectiveInstanceSizeNameOk returns a tuple with the EffectiveInstanceSizeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterProviderSettings) GetEffectiveInstanceSizeNameOk() (*string, bool) {
	if o == nil || IsNil(o.EffectiveInstanceSizeName) {
		return nil, false
	}

	return o.EffectiveInstanceSizeName, true
}

// HasEffectiveInstanceSizeName returns a boolean if a field has been set.
func (o *ClusterProviderSettings) HasEffectiveInstanceSizeName() bool {
	if o != nil && !IsNil(o.EffectiveInstanceSizeName) {
		return true
	}

	return false
}

// SetEffectiveInstanceSizeName gets a reference to the given string and assigns it to the EffectiveInstanceSizeName field.
func (o *ClusterProviderSettings) SetEffectiveInstanceSizeName(v string) {
	o.EffectiveInstanceSizeName = &v
}
