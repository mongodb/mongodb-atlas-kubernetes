// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ReplicationSpec20240805 Details that explain how MongoDB Cloud replicates data on the specified MongoDB database.
type ReplicationSpec20240805 struct {
	// Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster. If you include existing shard replication configurations in the request, you must specify this parameter. If you add a new shard to an existing Cluster, you may specify this parameter. The request deletes any existing shards  in the Cluster that you exclude from the request. This corresponds to Shard ID displayed in the UI.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Hardware specifications for nodes set for a given region. Each `regionConfigs` object must be unique by region and cloud provider within the `replicationSpec`. Each `regionConfigs` object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region. Each `regionConfigs` object must have either an `analyticsSpecs` object, `electableSpecs` object, or `readOnlySpecs` object. Tenant clusters only require `electableSpecs`. Dedicated clusters can specify any of these specifications, but must have at least one `electableSpecs` object within a `replicationSpec`.  **Example:**  If you set `replicationSpecs[n].regionConfigs[m].analyticsSpecs.instanceSize` : `M30`, set `replicationSpecs[n].regionConfigs[m].electableSpecs.instanceSize` : `M30` if you have electable nodes and `replicationSpecs[n].regionConfigs[m].readOnlySpecs.instanceSize` : `M30` if you have read-only nodes.
	RegionConfigs *[]CloudRegionConfig20240805 `json:"regionConfigs,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the zone in a Global Cluster. This value can be used to configure Global Cluster backup policies.
	// Read only field.
	ZoneId *string `json:"zoneId,omitempty"`
	// Human-readable label that describes the zone this shard belongs to in a Global Cluster. Provide this value only if `clusterType` : `GEOSHARDED` but not `selfManagedSharding` : `true`.
	ZoneName *string `json:"zoneName,omitempty"`
}

// NewReplicationSpec20240805 instantiates a new ReplicationSpec20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReplicationSpec20240805() *ReplicationSpec20240805 {
	this := ReplicationSpec20240805{}
	return &this
}

// NewReplicationSpec20240805WithDefaults instantiates a new ReplicationSpec20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReplicationSpec20240805WithDefaults() *ReplicationSpec20240805 {
	this := ReplicationSpec20240805{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ReplicationSpec20240805) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicationSpec20240805) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ReplicationSpec20240805) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ReplicationSpec20240805) SetId(v string) {
	o.Id = &v
}

// GetRegionConfigs returns the RegionConfigs field value if set, zero value otherwise
func (o *ReplicationSpec20240805) GetRegionConfigs() []CloudRegionConfig20240805 {
	if o == nil || IsNil(o.RegionConfigs) {
		var ret []CloudRegionConfig20240805
		return ret
	}
	return *o.RegionConfigs
}

// GetRegionConfigsOk returns a tuple with the RegionConfigs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicationSpec20240805) GetRegionConfigsOk() (*[]CloudRegionConfig20240805, bool) {
	if o == nil || IsNil(o.RegionConfigs) {
		return nil, false
	}

	return o.RegionConfigs, true
}

// HasRegionConfigs returns a boolean if a field has been set.
func (o *ReplicationSpec20240805) HasRegionConfigs() bool {
	if o != nil && !IsNil(o.RegionConfigs) {
		return true
	}

	return false
}

// SetRegionConfigs gets a reference to the given []CloudRegionConfig20240805 and assigns it to the RegionConfigs field.
func (o *ReplicationSpec20240805) SetRegionConfigs(v []CloudRegionConfig20240805) {
	o.RegionConfigs = &v
}

// GetZoneId returns the ZoneId field value if set, zero value otherwise
func (o *ReplicationSpec20240805) GetZoneId() string {
	if o == nil || IsNil(o.ZoneId) {
		var ret string
		return ret
	}
	return *o.ZoneId
}

// GetZoneIdOk returns a tuple with the ZoneId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicationSpec20240805) GetZoneIdOk() (*string, bool) {
	if o == nil || IsNil(o.ZoneId) {
		return nil, false
	}

	return o.ZoneId, true
}

// HasZoneId returns a boolean if a field has been set.
func (o *ReplicationSpec20240805) HasZoneId() bool {
	if o != nil && !IsNil(o.ZoneId) {
		return true
	}

	return false
}

// SetZoneId gets a reference to the given string and assigns it to the ZoneId field.
func (o *ReplicationSpec20240805) SetZoneId(v string) {
	o.ZoneId = &v
}

// GetZoneName returns the ZoneName field value if set, zero value otherwise
func (o *ReplicationSpec20240805) GetZoneName() string {
	if o == nil || IsNil(o.ZoneName) {
		var ret string
		return ret
	}
	return *o.ZoneName
}

// GetZoneNameOk returns a tuple with the ZoneName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReplicationSpec20240805) GetZoneNameOk() (*string, bool) {
	if o == nil || IsNil(o.ZoneName) {
		return nil, false
	}

	return o.ZoneName, true
}

// HasZoneName returns a boolean if a field has been set.
func (o *ReplicationSpec20240805) HasZoneName() bool {
	if o != nil && !IsNil(o.ZoneName) {
		return true
	}

	return false
}

// SetZoneName gets a reference to the given string and assigns it to the ZoneName field.
func (o *ReplicationSpec20240805) SetZoneName(v string) {
	o.ZoneName = &v
}
