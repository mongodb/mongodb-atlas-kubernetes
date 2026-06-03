// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupShardedClusterSnapshotMember struct for DiskBackupShardedClusterSnapshotMember
type DiskBackupShardedClusterSnapshotMember struct {
	// Human-readable label that identifies the cloud provider.
	// Read only field.
	CloudProvider string `json:"cloudProvider"`
	// Unique 24-hexadecimal digit string that identifies the snapshot.
	// Read only field.
	Id string `json:"id"`
	// Human-readable label that identifies the shard or config host from which MongoDB Cloud took this snapshot.
	// Read only field.
	ReplicaSetName string `json:"replicaSetName"`
}

// NewDiskBackupShardedClusterSnapshotMember instantiates a new DiskBackupShardedClusterSnapshotMember object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupShardedClusterSnapshotMember(cloudProvider string, id string, replicaSetName string) *DiskBackupShardedClusterSnapshotMember {
	this := DiskBackupShardedClusterSnapshotMember{}
	this.CloudProvider = cloudProvider
	this.Id = id
	this.ReplicaSetName = replicaSetName
	return &this
}

// NewDiskBackupShardedClusterSnapshotMemberWithDefaults instantiates a new DiskBackupShardedClusterSnapshotMember object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupShardedClusterSnapshotMemberWithDefaults() *DiskBackupShardedClusterSnapshotMember {
	this := DiskBackupShardedClusterSnapshotMember{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value
func (o *DiskBackupShardedClusterSnapshotMember) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *DiskBackupShardedClusterSnapshotMember) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *DiskBackupShardedClusterSnapshotMember) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetId returns the Id field value
func (o *DiskBackupShardedClusterSnapshotMember) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *DiskBackupShardedClusterSnapshotMember) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *DiskBackupShardedClusterSnapshotMember) SetId(v string) {
	o.Id = v
}

// GetReplicaSetName returns the ReplicaSetName field value
func (o *DiskBackupShardedClusterSnapshotMember) GetReplicaSetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value
// and a boolean to check if the value has been set.
func (o *DiskBackupShardedClusterSnapshotMember) GetReplicaSetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ReplicaSetName, true
}

// SetReplicaSetName sets field value
func (o *DiskBackupShardedClusterSnapshotMember) SetReplicaSetName(v string) {
	o.ReplicaSetName = v
}
