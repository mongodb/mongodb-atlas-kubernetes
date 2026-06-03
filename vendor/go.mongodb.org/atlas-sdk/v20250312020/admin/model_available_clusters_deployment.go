// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AvailableClustersDeployment Deployments that can be migrated to MongoDB Atlas.
type AvailableClustersDeployment struct {
	// Version of MongoDB Agent that monitors/manages the cluster.
	// Read only field.
	AgentVersion *string `json:"agentVersion,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the cluster.
	// Read only field.
	ClusterId *string `json:"clusterId,omitempty"`
	// Size of this database on disk at the time of the request expressed in bytes.
	// Read only field.
	DbSizeBytes *int64 `json:"dbSizeBytes,omitempty"`
	// Version of MongoDB features that this cluster supports.
	// Read only field.
	FeatureCompatibilityVersion string `json:"featureCompatibilityVersion"`
	// Flag that indicates whether Automation manages this cluster.
	// Read only field.
	Managed bool `json:"managed"`
	// Version of MongoDB that this cluster runs.
	// Read only field.
	MongoDBVersion string `json:"mongoDBVersion"`
	// Human-readable label that identifies this cluster.
	// Read only field.
	Name string `json:"name"`
	// Size of the Oplog on disk at the time of the request expressed in MB.
	// Read only field.
	OplogSizeMB *int `json:"oplogSizeMB,omitempty"`
	// Flag that indicates whether someone configured this cluster as a sharded cluster.  - If `true`, this cluster serves as a sharded cluster. - If `false`, this cluster serves as a replica set.
	// Read only field.
	Sharded bool `json:"sharded"`
	// Number of shards that comprise this cluster.
	// Read only field.
	ShardsSize *int `json:"shardsSize,omitempty"`
	// Flag that indicates whether someone enabled TLS for this cluster.
	// Read only field.
	TlsEnabled bool `json:"tlsEnabled"`
}

// NewAvailableClustersDeployment instantiates a new AvailableClustersDeployment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAvailableClustersDeployment(featureCompatibilityVersion string, managed bool, mongoDBVersion string, name string, sharded bool, tlsEnabled bool) *AvailableClustersDeployment {
	this := AvailableClustersDeployment{}
	this.FeatureCompatibilityVersion = featureCompatibilityVersion
	this.Managed = managed
	this.MongoDBVersion = mongoDBVersion
	this.Name = name
	this.Sharded = sharded
	this.TlsEnabled = tlsEnabled
	return &this
}

// NewAvailableClustersDeploymentWithDefaults instantiates a new AvailableClustersDeployment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAvailableClustersDeploymentWithDefaults() *AvailableClustersDeployment {
	this := AvailableClustersDeployment{}
	return &this
}

// GetAgentVersion returns the AgentVersion field value if set, zero value otherwise
func (o *AvailableClustersDeployment) GetAgentVersion() string {
	if o == nil || IsNil(o.AgentVersion) {
		var ret string
		return ret
	}
	return *o.AgentVersion
}

// GetAgentVersionOk returns a tuple with the AgentVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetAgentVersionOk() (*string, bool) {
	if o == nil || IsNil(o.AgentVersion) {
		return nil, false
	}

	return o.AgentVersion, true
}

// HasAgentVersion returns a boolean if a field has been set.
func (o *AvailableClustersDeployment) HasAgentVersion() bool {
	if o != nil && !IsNil(o.AgentVersion) {
		return true
	}

	return false
}

// SetAgentVersion gets a reference to the given string and assigns it to the AgentVersion field.
func (o *AvailableClustersDeployment) SetAgentVersion(v string) {
	o.AgentVersion = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise
func (o *AvailableClustersDeployment) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}

	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *AvailableClustersDeployment) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *AvailableClustersDeployment) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetDbSizeBytes returns the DbSizeBytes field value if set, zero value otherwise
func (o *AvailableClustersDeployment) GetDbSizeBytes() int64 {
	if o == nil || IsNil(o.DbSizeBytes) {
		var ret int64
		return ret
	}
	return *o.DbSizeBytes
}

// GetDbSizeBytesOk returns a tuple with the DbSizeBytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetDbSizeBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.DbSizeBytes) {
		return nil, false
	}

	return o.DbSizeBytes, true
}

// HasDbSizeBytes returns a boolean if a field has been set.
func (o *AvailableClustersDeployment) HasDbSizeBytes() bool {
	if o != nil && !IsNil(o.DbSizeBytes) {
		return true
	}

	return false
}

// SetDbSizeBytes gets a reference to the given int64 and assigns it to the DbSizeBytes field.
func (o *AvailableClustersDeployment) SetDbSizeBytes(v int64) {
	o.DbSizeBytes = &v
}

// GetFeatureCompatibilityVersion returns the FeatureCompatibilityVersion field value
func (o *AvailableClustersDeployment) GetFeatureCompatibilityVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FeatureCompatibilityVersion
}

// GetFeatureCompatibilityVersionOk returns a tuple with the FeatureCompatibilityVersion field value
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetFeatureCompatibilityVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FeatureCompatibilityVersion, true
}

// SetFeatureCompatibilityVersion sets field value
func (o *AvailableClustersDeployment) SetFeatureCompatibilityVersion(v string) {
	o.FeatureCompatibilityVersion = v
}

// GetManaged returns the Managed field value
func (o *AvailableClustersDeployment) GetManaged() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Managed
}

// GetManagedOk returns a tuple with the Managed field value
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetManagedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Managed, true
}

// SetManaged sets field value
func (o *AvailableClustersDeployment) SetManaged(v bool) {
	o.Managed = v
}

// GetMongoDBVersion returns the MongoDBVersion field value
func (o *AvailableClustersDeployment) GetMongoDBVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MongoDBVersion
}

// GetMongoDBVersionOk returns a tuple with the MongoDBVersion field value
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetMongoDBVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MongoDBVersion, true
}

// SetMongoDBVersion sets field value
func (o *AvailableClustersDeployment) SetMongoDBVersion(v string) {
	o.MongoDBVersion = v
}

// GetName returns the Name field value
func (o *AvailableClustersDeployment) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *AvailableClustersDeployment) SetName(v string) {
	o.Name = v
}

// GetOplogSizeMB returns the OplogSizeMB field value if set, zero value otherwise
func (o *AvailableClustersDeployment) GetOplogSizeMB() int {
	if o == nil || IsNil(o.OplogSizeMB) {
		var ret int
		return ret
	}
	return *o.OplogSizeMB
}

// GetOplogSizeMBOk returns a tuple with the OplogSizeMB field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetOplogSizeMBOk() (*int, bool) {
	if o == nil || IsNil(o.OplogSizeMB) {
		return nil, false
	}

	return o.OplogSizeMB, true
}

// HasOplogSizeMB returns a boolean if a field has been set.
func (o *AvailableClustersDeployment) HasOplogSizeMB() bool {
	if o != nil && !IsNil(o.OplogSizeMB) {
		return true
	}

	return false
}

// SetOplogSizeMB gets a reference to the given int and assigns it to the OplogSizeMB field.
func (o *AvailableClustersDeployment) SetOplogSizeMB(v int) {
	o.OplogSizeMB = &v
}

// GetSharded returns the Sharded field value
func (o *AvailableClustersDeployment) GetSharded() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Sharded
}

// GetShardedOk returns a tuple with the Sharded field value
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetShardedOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Sharded, true
}

// SetSharded sets field value
func (o *AvailableClustersDeployment) SetSharded(v bool) {
	o.Sharded = v
}

// GetShardsSize returns the ShardsSize field value if set, zero value otherwise
func (o *AvailableClustersDeployment) GetShardsSize() int {
	if o == nil || IsNil(o.ShardsSize) {
		var ret int
		return ret
	}
	return *o.ShardsSize
}

// GetShardsSizeOk returns a tuple with the ShardsSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetShardsSizeOk() (*int, bool) {
	if o == nil || IsNil(o.ShardsSize) {
		return nil, false
	}

	return o.ShardsSize, true
}

// HasShardsSize returns a boolean if a field has been set.
func (o *AvailableClustersDeployment) HasShardsSize() bool {
	if o != nil && !IsNil(o.ShardsSize) {
		return true
	}

	return false
}

// SetShardsSize gets a reference to the given int and assigns it to the ShardsSize field.
func (o *AvailableClustersDeployment) SetShardsSize(v int) {
	o.ShardsSize = &v
}

// GetTlsEnabled returns the TlsEnabled field value
func (o *AvailableClustersDeployment) GetTlsEnabled() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.TlsEnabled
}

// GetTlsEnabledOk returns a tuple with the TlsEnabled field value
// and a boolean to check if the value has been set.
func (o *AvailableClustersDeployment) GetTlsEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TlsEnabled, true
}

// SetTlsEnabled sets field value
func (o *AvailableClustersDeployment) SetTlsEnabled(v bool) {
	o.TlsEnabled = v
}
