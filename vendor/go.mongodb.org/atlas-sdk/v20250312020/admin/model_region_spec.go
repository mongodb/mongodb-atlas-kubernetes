// Code based on the AtlasAPI V2 OpenAPI file

package admin

// RegionSpec Physical location where MongoDB Cloud provisions cluster nodes.
type RegionSpec struct {
	// Number of analytics nodes in the region. Analytics nodes handle analytic data such as reporting queries from MongoDB Connector for Business Intelligence on MongoDB Cloud. Analytics nodes are read-only, and can never become the primary. Use `replicationSpecs[n].{region}.analyticsNodes` instead.
	AnalyticsNodes *int `json:"analyticsNodes,omitempty"`
	// Number of electable nodes to deploy in the specified region. Electable nodes can become the primary and can facilitate local reads. Use `replicationSpecs[n].{region}.electableNodes` instead.
	ElectableNodes *int `json:"electableNodes,omitempty"`
	// Number that indicates the election priority of the region. To identify the Preferred Region of the cluster, set this parameter to `7`. The primary node runs in the **Preferred Region**. To identify a read-only region, set this parameter to `0`.
	Priority *int `json:"priority,omitempty"`
	// Number of read-only nodes in the region. Read-only nodes can never become the primary member, but can facilitate local reads. Use `replicationSpecs[n].{region}.readOnlyNodes` instead.
	ReadOnlyNodes *int `json:"readOnlyNodes,omitempty"`
}

// NewRegionSpec instantiates a new RegionSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegionSpec() *RegionSpec {
	this := RegionSpec{}
	return &this
}

// NewRegionSpecWithDefaults instantiates a new RegionSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegionSpecWithDefaults() *RegionSpec {
	this := RegionSpec{}
	return &this
}

// GetAnalyticsNodes returns the AnalyticsNodes field value if set, zero value otherwise
func (o *RegionSpec) GetAnalyticsNodes() int {
	if o == nil || IsNil(o.AnalyticsNodes) {
		var ret int
		return ret
	}
	return *o.AnalyticsNodes
}

// GetAnalyticsNodesOk returns a tuple with the AnalyticsNodes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RegionSpec) GetAnalyticsNodesOk() (*int, bool) {
	if o == nil || IsNil(o.AnalyticsNodes) {
		return nil, false
	}

	return o.AnalyticsNodes, true
}

// HasAnalyticsNodes returns a boolean if a field has been set.
func (o *RegionSpec) HasAnalyticsNodes() bool {
	if o != nil && !IsNil(o.AnalyticsNodes) {
		return true
	}

	return false
}

// SetAnalyticsNodes gets a reference to the given int and assigns it to the AnalyticsNodes field.
func (o *RegionSpec) SetAnalyticsNodes(v int) {
	o.AnalyticsNodes = &v
}

// GetElectableNodes returns the ElectableNodes field value if set, zero value otherwise
func (o *RegionSpec) GetElectableNodes() int {
	if o == nil || IsNil(o.ElectableNodes) {
		var ret int
		return ret
	}
	return *o.ElectableNodes
}

// GetElectableNodesOk returns a tuple with the ElectableNodes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RegionSpec) GetElectableNodesOk() (*int, bool) {
	if o == nil || IsNil(o.ElectableNodes) {
		return nil, false
	}

	return o.ElectableNodes, true
}

// HasElectableNodes returns a boolean if a field has been set.
func (o *RegionSpec) HasElectableNodes() bool {
	if o != nil && !IsNil(o.ElectableNodes) {
		return true
	}

	return false
}

// SetElectableNodes gets a reference to the given int and assigns it to the ElectableNodes field.
func (o *RegionSpec) SetElectableNodes(v int) {
	o.ElectableNodes = &v
}

// GetPriority returns the Priority field value if set, zero value otherwise
func (o *RegionSpec) GetPriority() int {
	if o == nil || IsNil(o.Priority) {
		var ret int
		return ret
	}
	return *o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RegionSpec) GetPriorityOk() (*int, bool) {
	if o == nil || IsNil(o.Priority) {
		return nil, false
	}

	return o.Priority, true
}

// HasPriority returns a boolean if a field has been set.
func (o *RegionSpec) HasPriority() bool {
	if o != nil && !IsNil(o.Priority) {
		return true
	}

	return false
}

// SetPriority gets a reference to the given int and assigns it to the Priority field.
func (o *RegionSpec) SetPriority(v int) {
	o.Priority = &v
}

// GetReadOnlyNodes returns the ReadOnlyNodes field value if set, zero value otherwise
func (o *RegionSpec) GetReadOnlyNodes() int {
	if o == nil || IsNil(o.ReadOnlyNodes) {
		var ret int
		return ret
	}
	return *o.ReadOnlyNodes
}

// GetReadOnlyNodesOk returns a tuple with the ReadOnlyNodes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RegionSpec) GetReadOnlyNodesOk() (*int, bool) {
	if o == nil || IsNil(o.ReadOnlyNodes) {
		return nil, false
	}

	return o.ReadOnlyNodes, true
}

// HasReadOnlyNodes returns a boolean if a field has been set.
func (o *RegionSpec) HasReadOnlyNodes() bool {
	if o != nil && !IsNil(o.ReadOnlyNodes) {
		return true
	}

	return false
}

// SetReadOnlyNodes gets a reference to the given int and assigns it to the ReadOnlyNodes field.
func (o *RegionSpec) SetReadOnlyNodes(v int) {
	o.ReadOnlyNodes = &v
}
