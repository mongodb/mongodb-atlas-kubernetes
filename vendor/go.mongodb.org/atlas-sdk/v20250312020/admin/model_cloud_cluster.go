// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudCluster Settings that describe the clusters in each project that the API key is authorized to view.
type CloudCluster struct {
	// Whole number that indicates the quantity of alerts open on the cluster.
	// Read only field.
	AlertCount *int `json:"alertCount,omitempty"`
	// Flag that indicates whether authentication is required to access the nodes in this cluster.
	// Read only field.
	AuthEnabled *bool `json:"authEnabled,omitempty"`
	// Term that expresses how many nodes of the cluster can be accessed when MongoDB Cloud receives this request. This parameter returns `available` when all nodes are accessible, `warning` only when some nodes in the cluster can be accessed, `unavailable` when the cluster can't be accessed, or `dead` when the cluster has been deactivated.
	// Read only field.
	Availability *string `json:"availability,omitempty"`
	// Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses Cloud Backups for dedicated clusters and Shared Cluster Backups for tenant clusters. If set to `false`, the cluster doesn't use MongoDB Cloud backups.
	// Read only field.
	BackupEnabled *bool `json:"backupEnabled,omitempty"`
	// Unique 24-hexadecimal character string that identifies the cluster. Each `clusterId` is used only once across all MongoDB Cloud deployments.
	// Read only field.
	ClusterId *string `json:"clusterId,omitempty"`
	// Total size of the data stored on each node in the cluster. The resource expresses this value in bytes.
	// Read only field.
	DataSizeBytes *int64 `json:"dataSizeBytes,omitempty"`
	// Human-readable label that identifies the cluster.
	// Read only field.
	Name *string `json:"name,omitempty"`
	// Whole number that indicates the quantity of nodes that comprise the cluster.
	// Read only field.
	NodeCount *int `json:"nodeCount,omitempty"`
	// Flag that indicates whether TLS authentication is required to access the nodes in this cluster.
	// Read only field.
	SslEnabled *bool `json:"sslEnabled,omitempty"`
	// Human-readable label that indicates the cluster type.
	// Read only field.
	Type *string `json:"type,omitempty"`
	// List that contains the versions of MongoDB that each node in the cluster runs.
	// Read only field.
	Versions *[]string `json:"versions,omitempty"`
}

// NewCloudCluster instantiates a new CloudCluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudCluster() *CloudCluster {
	this := CloudCluster{}
	return &this
}

// NewCloudClusterWithDefaults instantiates a new CloudCluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudClusterWithDefaults() *CloudCluster {
	this := CloudCluster{}
	return &this
}

// GetAlertCount returns the AlertCount field value if set, zero value otherwise
func (o *CloudCluster) GetAlertCount() int {
	if o == nil || IsNil(o.AlertCount) {
		var ret int
		return ret
	}
	return *o.AlertCount
}

// GetAlertCountOk returns a tuple with the AlertCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetAlertCountOk() (*int, bool) {
	if o == nil || IsNil(o.AlertCount) {
		return nil, false
	}

	return o.AlertCount, true
}

// HasAlertCount returns a boolean if a field has been set.
func (o *CloudCluster) HasAlertCount() bool {
	if o != nil && !IsNil(o.AlertCount) {
		return true
	}

	return false
}

// SetAlertCount gets a reference to the given int and assigns it to the AlertCount field.
func (o *CloudCluster) SetAlertCount(v int) {
	o.AlertCount = &v
}

// GetAuthEnabled returns the AuthEnabled field value if set, zero value otherwise
func (o *CloudCluster) GetAuthEnabled() bool {
	if o == nil || IsNil(o.AuthEnabled) {
		var ret bool
		return ret
	}
	return *o.AuthEnabled
}

// GetAuthEnabledOk returns a tuple with the AuthEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetAuthEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AuthEnabled) {
		return nil, false
	}

	return o.AuthEnabled, true
}

// HasAuthEnabled returns a boolean if a field has been set.
func (o *CloudCluster) HasAuthEnabled() bool {
	if o != nil && !IsNil(o.AuthEnabled) {
		return true
	}

	return false
}

// SetAuthEnabled gets a reference to the given bool and assigns it to the AuthEnabled field.
func (o *CloudCluster) SetAuthEnabled(v bool) {
	o.AuthEnabled = &v
}

// GetAvailability returns the Availability field value if set, zero value otherwise
func (o *CloudCluster) GetAvailability() string {
	if o == nil || IsNil(o.Availability) {
		var ret string
		return ret
	}
	return *o.Availability
}

// GetAvailabilityOk returns a tuple with the Availability field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetAvailabilityOk() (*string, bool) {
	if o == nil || IsNil(o.Availability) {
		return nil, false
	}

	return o.Availability, true
}

// HasAvailability returns a boolean if a field has been set.
func (o *CloudCluster) HasAvailability() bool {
	if o != nil && !IsNil(o.Availability) {
		return true
	}

	return false
}

// SetAvailability gets a reference to the given string and assigns it to the Availability field.
func (o *CloudCluster) SetAvailability(v string) {
	o.Availability = &v
}

// GetBackupEnabled returns the BackupEnabled field value if set, zero value otherwise
func (o *CloudCluster) GetBackupEnabled() bool {
	if o == nil || IsNil(o.BackupEnabled) {
		var ret bool
		return ret
	}
	return *o.BackupEnabled
}

// GetBackupEnabledOk returns a tuple with the BackupEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetBackupEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.BackupEnabled) {
		return nil, false
	}

	return o.BackupEnabled, true
}

// HasBackupEnabled returns a boolean if a field has been set.
func (o *CloudCluster) HasBackupEnabled() bool {
	if o != nil && !IsNil(o.BackupEnabled) {
		return true
	}

	return false
}

// SetBackupEnabled gets a reference to the given bool and assigns it to the BackupEnabled field.
func (o *CloudCluster) SetBackupEnabled(v bool) {
	o.BackupEnabled = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise
func (o *CloudCluster) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}

	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *CloudCluster) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *CloudCluster) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetDataSizeBytes returns the DataSizeBytes field value if set, zero value otherwise
func (o *CloudCluster) GetDataSizeBytes() int64 {
	if o == nil || IsNil(o.DataSizeBytes) {
		var ret int64
		return ret
	}
	return *o.DataSizeBytes
}

// GetDataSizeBytesOk returns a tuple with the DataSizeBytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetDataSizeBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.DataSizeBytes) {
		return nil, false
	}

	return o.DataSizeBytes, true
}

// HasDataSizeBytes returns a boolean if a field has been set.
func (o *CloudCluster) HasDataSizeBytes() bool {
	if o != nil && !IsNil(o.DataSizeBytes) {
		return true
	}

	return false
}

// SetDataSizeBytes gets a reference to the given int64 and assigns it to the DataSizeBytes field.
func (o *CloudCluster) SetDataSizeBytes(v int64) {
	o.DataSizeBytes = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *CloudCluster) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *CloudCluster) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *CloudCluster) SetName(v string) {
	o.Name = &v
}

// GetNodeCount returns the NodeCount field value if set, zero value otherwise
func (o *CloudCluster) GetNodeCount() int {
	if o == nil || IsNil(o.NodeCount) {
		var ret int
		return ret
	}
	return *o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetNodeCountOk() (*int, bool) {
	if o == nil || IsNil(o.NodeCount) {
		return nil, false
	}

	return o.NodeCount, true
}

// HasNodeCount returns a boolean if a field has been set.
func (o *CloudCluster) HasNodeCount() bool {
	if o != nil && !IsNil(o.NodeCount) {
		return true
	}

	return false
}

// SetNodeCount gets a reference to the given int and assigns it to the NodeCount field.
func (o *CloudCluster) SetNodeCount(v int) {
	o.NodeCount = &v
}

// GetSslEnabled returns the SslEnabled field value if set, zero value otherwise
func (o *CloudCluster) GetSslEnabled() bool {
	if o == nil || IsNil(o.SslEnabled) {
		var ret bool
		return ret
	}
	return *o.SslEnabled
}

// GetSslEnabledOk returns a tuple with the SslEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetSslEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.SslEnabled) {
		return nil, false
	}

	return o.SslEnabled, true
}

// HasSslEnabled returns a boolean if a field has been set.
func (o *CloudCluster) HasSslEnabled() bool {
	if o != nil && !IsNil(o.SslEnabled) {
		return true
	}

	return false
}

// SetSslEnabled gets a reference to the given bool and assigns it to the SslEnabled field.
func (o *CloudCluster) SetSslEnabled(v bool) {
	o.SslEnabled = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *CloudCluster) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *CloudCluster) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *CloudCluster) SetType(v string) {
	o.Type = &v
}

// GetVersions returns the Versions field value if set, zero value otherwise
func (o *CloudCluster) GetVersions() []string {
	if o == nil || IsNil(o.Versions) {
		var ret []string
		return ret
	}
	return *o.Versions
}

// GetVersionsOk returns a tuple with the Versions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudCluster) GetVersionsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Versions) {
		return nil, false
	}

	return o.Versions, true
}

// HasVersions returns a boolean if a field has been set.
func (o *CloudCluster) HasVersions() bool {
	if o != nil && !IsNil(o.Versions) {
		return true
	}

	return false
}

// SetVersions gets a reference to the given []string and assigns it to the Versions field.
func (o *CloudCluster) SetVersions(v []string) {
	o.Versions = &v
}
