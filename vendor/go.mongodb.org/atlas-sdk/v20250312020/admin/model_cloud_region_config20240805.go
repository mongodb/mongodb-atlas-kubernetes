// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudRegionConfig20240805 Cloud service provider on which MongoDB Cloud provisions the hosts.
type CloudRegionConfig20240805 struct {
	ElectableSpecs *HardwareSpec20240805 `json:"electableSpecs,omitempty"`
	// Precedence is given to this region when a primary election occurs. If your `regionConfigs` has only `readOnlySpecs`, `analyticsSpecs`, or both, set this value to `0`. If you have multiple `regionConfigs` objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order. The highest priority is `7`.  **Example:** If you have three regions, their priorities would be `7`, `6`, and `5` respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be `4` and `3` respectively.
	Priority *int `json:"priority,omitempty"`
	// Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.
	ProviderName *string `json:"providerName,omitempty"`
	// Physical location of your MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. The region name is only returned in the response for single-region clusters. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. It assigns the VPC a Classless Inter-Domain Routing (CIDR) block. To limit a new VPC peering connection to one Classless Inter-Domain Routing (CIDR) block and region, create the connection first. Deploy the cluster after the connection starts. GCP Clusters and Multi-region clusters require one VPC peering connection for each region. MongoDB nodes can use only the peering connection that resides in the same region as the nodes to communicate with the peered VPC.
	RegionName              *string                        `json:"regionName,omitempty"`
	AnalyticsAutoScaling    *AdvancedAutoScalingSettings   `json:"analyticsAutoScaling,omitempty"`
	AnalyticsSpecs          *DedicatedHardwareSpec20240805 `json:"analyticsSpecs,omitempty"`
	AutoScaling             *AdvancedAutoScalingSettings   `json:"autoScaling,omitempty"`
	EffectiveAnalyticsSpecs *DedicatedHardwareSpec20240805 `json:"effectiveAnalyticsSpecs,omitempty"`
	EffectiveElectableSpecs *DedicatedHardwareSpec20240805 `json:"effectiveElectableSpecs,omitempty"`
	EffectiveReadOnlySpecs  *DedicatedHardwareSpec20240805 `json:"effectiveReadOnlySpecs,omitempty"`
	ReadOnlySpecs           *DedicatedHardwareSpec20240805 `json:"readOnlySpecs,omitempty"`
	// Cloud service provider on which MongoDB Cloud provisioned the multi-tenant cluster. The resource returns this parameter when `providerName` is `TENANT` and `electableSpecs.instanceSize` is `M0`, `M2` or `M5`.   Please note that  using an `instanceSize` of `M2` or `M5` will create a Flex cluster instead. Support for the `instanceSize` of `M2` or `M5` will be discontinued in January 2026. We recommend using the Create Flex Cluster API for such configurations moving forward.
	BackingProviderName *string `json:"backingProviderName,omitempty"`
}

// NewCloudRegionConfig20240805 instantiates a new CloudRegionConfig20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudRegionConfig20240805() *CloudRegionConfig20240805 {
	this := CloudRegionConfig20240805{}
	return &this
}

// NewCloudRegionConfig20240805WithDefaults instantiates a new CloudRegionConfig20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudRegionConfig20240805WithDefaults() *CloudRegionConfig20240805 {
	this := CloudRegionConfig20240805{}
	return &this
}

// GetElectableSpecs returns the ElectableSpecs field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetElectableSpecs() HardwareSpec20240805 {
	if o == nil || IsNil(o.ElectableSpecs) {
		var ret HardwareSpec20240805
		return ret
	}
	return *o.ElectableSpecs
}

// GetElectableSpecsOk returns a tuple with the ElectableSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetElectableSpecsOk() (*HardwareSpec20240805, bool) {
	if o == nil || IsNil(o.ElectableSpecs) {
		return nil, false
	}

	return o.ElectableSpecs, true
}

// HasElectableSpecs returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasElectableSpecs() bool {
	if o != nil && !IsNil(o.ElectableSpecs) {
		return true
	}

	return false
}

// SetElectableSpecs gets a reference to the given HardwareSpec20240805 and assigns it to the ElectableSpecs field.
func (o *CloudRegionConfig20240805) SetElectableSpecs(v HardwareSpec20240805) {
	o.ElectableSpecs = &v
}

// GetPriority returns the Priority field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetPriority() int {
	if o == nil || IsNil(o.Priority) {
		var ret int
		return ret
	}
	return *o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetPriorityOk() (*int, bool) {
	if o == nil || IsNil(o.Priority) {
		return nil, false
	}

	return o.Priority, true
}

// HasPriority returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasPriority() bool {
	if o != nil && !IsNil(o.Priority) {
		return true
	}

	return false
}

// SetPriority gets a reference to the given int and assigns it to the Priority field.
func (o *CloudRegionConfig20240805) SetPriority(v int) {
	o.Priority = &v
}

// GetProviderName returns the ProviderName field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetProviderName() string {
	if o == nil || IsNil(o.ProviderName) {
		var ret string
		return ret
	}
	return *o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderName) {
		return nil, false
	}

	return o.ProviderName, true
}

// HasProviderName returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasProviderName() bool {
	if o != nil && !IsNil(o.ProviderName) {
		return true
	}

	return false
}

// SetProviderName gets a reference to the given string and assigns it to the ProviderName field.
func (o *CloudRegionConfig20240805) SetProviderName(v string) {
	o.ProviderName = &v
}

// GetRegionName returns the RegionName field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetRegionName() string {
	if o == nil || IsNil(o.RegionName) {
		var ret string
		return ret
	}
	return *o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionName) {
		return nil, false
	}

	return o.RegionName, true
}

// HasRegionName returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasRegionName() bool {
	if o != nil && !IsNil(o.RegionName) {
		return true
	}

	return false
}

// SetRegionName gets a reference to the given string and assigns it to the RegionName field.
func (o *CloudRegionConfig20240805) SetRegionName(v string) {
	o.RegionName = &v
}

// GetAnalyticsAutoScaling returns the AnalyticsAutoScaling field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetAnalyticsAutoScaling() AdvancedAutoScalingSettings {
	if o == nil || IsNil(o.AnalyticsAutoScaling) {
		var ret AdvancedAutoScalingSettings
		return ret
	}
	return *o.AnalyticsAutoScaling
}

// GetAnalyticsAutoScalingOk returns a tuple with the AnalyticsAutoScaling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetAnalyticsAutoScalingOk() (*AdvancedAutoScalingSettings, bool) {
	if o == nil || IsNil(o.AnalyticsAutoScaling) {
		return nil, false
	}

	return o.AnalyticsAutoScaling, true
}

// HasAnalyticsAutoScaling returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasAnalyticsAutoScaling() bool {
	if o != nil && !IsNil(o.AnalyticsAutoScaling) {
		return true
	}

	return false
}

// SetAnalyticsAutoScaling gets a reference to the given AdvancedAutoScalingSettings and assigns it to the AnalyticsAutoScaling field.
func (o *CloudRegionConfig20240805) SetAnalyticsAutoScaling(v AdvancedAutoScalingSettings) {
	o.AnalyticsAutoScaling = &v
}

// GetAnalyticsSpecs returns the AnalyticsSpecs field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetAnalyticsSpecs() DedicatedHardwareSpec20240805 {
	if o == nil || IsNil(o.AnalyticsSpecs) {
		var ret DedicatedHardwareSpec20240805
		return ret
	}
	return *o.AnalyticsSpecs
}

// GetAnalyticsSpecsOk returns a tuple with the AnalyticsSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetAnalyticsSpecsOk() (*DedicatedHardwareSpec20240805, bool) {
	if o == nil || IsNil(o.AnalyticsSpecs) {
		return nil, false
	}

	return o.AnalyticsSpecs, true
}

// HasAnalyticsSpecs returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasAnalyticsSpecs() bool {
	if o != nil && !IsNil(o.AnalyticsSpecs) {
		return true
	}

	return false
}

// SetAnalyticsSpecs gets a reference to the given DedicatedHardwareSpec20240805 and assigns it to the AnalyticsSpecs field.
func (o *CloudRegionConfig20240805) SetAnalyticsSpecs(v DedicatedHardwareSpec20240805) {
	o.AnalyticsSpecs = &v
}

// GetAutoScaling returns the AutoScaling field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetAutoScaling() AdvancedAutoScalingSettings {
	if o == nil || IsNil(o.AutoScaling) {
		var ret AdvancedAutoScalingSettings
		return ret
	}
	return *o.AutoScaling
}

// GetAutoScalingOk returns a tuple with the AutoScaling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetAutoScalingOk() (*AdvancedAutoScalingSettings, bool) {
	if o == nil || IsNil(o.AutoScaling) {
		return nil, false
	}

	return o.AutoScaling, true
}

// HasAutoScaling returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasAutoScaling() bool {
	if o != nil && !IsNil(o.AutoScaling) {
		return true
	}

	return false
}

// SetAutoScaling gets a reference to the given AdvancedAutoScalingSettings and assigns it to the AutoScaling field.
func (o *CloudRegionConfig20240805) SetAutoScaling(v AdvancedAutoScalingSettings) {
	o.AutoScaling = &v
}

// GetEffectiveAnalyticsSpecs returns the EffectiveAnalyticsSpecs field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetEffectiveAnalyticsSpecs() DedicatedHardwareSpec20240805 {
	if o == nil || IsNil(o.EffectiveAnalyticsSpecs) {
		var ret DedicatedHardwareSpec20240805
		return ret
	}
	return *o.EffectiveAnalyticsSpecs
}

// GetEffectiveAnalyticsSpecsOk returns a tuple with the EffectiveAnalyticsSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetEffectiveAnalyticsSpecsOk() (*DedicatedHardwareSpec20240805, bool) {
	if o == nil || IsNil(o.EffectiveAnalyticsSpecs) {
		return nil, false
	}

	return o.EffectiveAnalyticsSpecs, true
}

// HasEffectiveAnalyticsSpecs returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasEffectiveAnalyticsSpecs() bool {
	if o != nil && !IsNil(o.EffectiveAnalyticsSpecs) {
		return true
	}

	return false
}

// SetEffectiveAnalyticsSpecs gets a reference to the given DedicatedHardwareSpec20240805 and assigns it to the EffectiveAnalyticsSpecs field.
func (o *CloudRegionConfig20240805) SetEffectiveAnalyticsSpecs(v DedicatedHardwareSpec20240805) {
	o.EffectiveAnalyticsSpecs = &v
}

// GetEffectiveElectableSpecs returns the EffectiveElectableSpecs field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetEffectiveElectableSpecs() DedicatedHardwareSpec20240805 {
	if o == nil || IsNil(o.EffectiveElectableSpecs) {
		var ret DedicatedHardwareSpec20240805
		return ret
	}
	return *o.EffectiveElectableSpecs
}

// GetEffectiveElectableSpecsOk returns a tuple with the EffectiveElectableSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetEffectiveElectableSpecsOk() (*DedicatedHardwareSpec20240805, bool) {
	if o == nil || IsNil(o.EffectiveElectableSpecs) {
		return nil, false
	}

	return o.EffectiveElectableSpecs, true
}

// HasEffectiveElectableSpecs returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasEffectiveElectableSpecs() bool {
	if o != nil && !IsNil(o.EffectiveElectableSpecs) {
		return true
	}

	return false
}

// SetEffectiveElectableSpecs gets a reference to the given DedicatedHardwareSpec20240805 and assigns it to the EffectiveElectableSpecs field.
func (o *CloudRegionConfig20240805) SetEffectiveElectableSpecs(v DedicatedHardwareSpec20240805) {
	o.EffectiveElectableSpecs = &v
}

// GetEffectiveReadOnlySpecs returns the EffectiveReadOnlySpecs field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetEffectiveReadOnlySpecs() DedicatedHardwareSpec20240805 {
	if o == nil || IsNil(o.EffectiveReadOnlySpecs) {
		var ret DedicatedHardwareSpec20240805
		return ret
	}
	return *o.EffectiveReadOnlySpecs
}

// GetEffectiveReadOnlySpecsOk returns a tuple with the EffectiveReadOnlySpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetEffectiveReadOnlySpecsOk() (*DedicatedHardwareSpec20240805, bool) {
	if o == nil || IsNil(o.EffectiveReadOnlySpecs) {
		return nil, false
	}

	return o.EffectiveReadOnlySpecs, true
}

// HasEffectiveReadOnlySpecs returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasEffectiveReadOnlySpecs() bool {
	if o != nil && !IsNil(o.EffectiveReadOnlySpecs) {
		return true
	}

	return false
}

// SetEffectiveReadOnlySpecs gets a reference to the given DedicatedHardwareSpec20240805 and assigns it to the EffectiveReadOnlySpecs field.
func (o *CloudRegionConfig20240805) SetEffectiveReadOnlySpecs(v DedicatedHardwareSpec20240805) {
	o.EffectiveReadOnlySpecs = &v
}

// GetReadOnlySpecs returns the ReadOnlySpecs field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetReadOnlySpecs() DedicatedHardwareSpec20240805 {
	if o == nil || IsNil(o.ReadOnlySpecs) {
		var ret DedicatedHardwareSpec20240805
		return ret
	}
	return *o.ReadOnlySpecs
}

// GetReadOnlySpecsOk returns a tuple with the ReadOnlySpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetReadOnlySpecsOk() (*DedicatedHardwareSpec20240805, bool) {
	if o == nil || IsNil(o.ReadOnlySpecs) {
		return nil, false
	}

	return o.ReadOnlySpecs, true
}

// HasReadOnlySpecs returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasReadOnlySpecs() bool {
	if o != nil && !IsNil(o.ReadOnlySpecs) {
		return true
	}

	return false
}

// SetReadOnlySpecs gets a reference to the given DedicatedHardwareSpec20240805 and assigns it to the ReadOnlySpecs field.
func (o *CloudRegionConfig20240805) SetReadOnlySpecs(v DedicatedHardwareSpec20240805) {
	o.ReadOnlySpecs = &v
}

// GetBackingProviderName returns the BackingProviderName field value if set, zero value otherwise
func (o *CloudRegionConfig20240805) GetBackingProviderName() string {
	if o == nil || IsNil(o.BackingProviderName) {
		var ret string
		return ret
	}
	return *o.BackingProviderName
}

// GetBackingProviderNameOk returns a tuple with the BackingProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudRegionConfig20240805) GetBackingProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.BackingProviderName) {
		return nil, false
	}

	return o.BackingProviderName, true
}

// HasBackingProviderName returns a boolean if a field has been set.
func (o *CloudRegionConfig20240805) HasBackingProviderName() bool {
	if o != nil && !IsNil(o.BackingProviderName) {
		return true
	}

	return false
}

// SetBackingProviderName gets a reference to the given string and assigns it to the BackingProviderName field.
func (o *CloudRegionConfig20240805) SetBackingProviderName(v string) {
	o.BackingProviderName = &v
}
