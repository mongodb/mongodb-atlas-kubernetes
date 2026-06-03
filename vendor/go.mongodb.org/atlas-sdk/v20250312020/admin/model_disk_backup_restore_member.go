// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupRestoreMember struct for DiskBackupRestoreMember
type DiskBackupRestoreMember struct {
	// One Uniform Resource Locator that points to the compressed snapshot files for manual download. MongoDB Cloud returns this parameter when `\"deliveryType\" : \"download\"`.
	// Read only field.
	DownloadUrl *string `json:"downloadUrl,omitempty"`
	// One or more Uniform Resource Locators (URLs) that point to the compressed snapshot files for manual download and the corresponding private endpoint(s). MongoDB Cloud returns this parameter when `\"deliveryType\" : \"download\"` and the download can be performed privately.
	// Read only field.
	PrivateDownloadDeliveryUrls *[]ApiPrivateDownloadDeliveryUrl `json:"privateDownloadDeliveryUrls,omitempty"`
	// Human-readable label that identifies the replica set on the sharded cluster.
	// Read only field.
	ReplicaSetName *string `json:"replicaSetName,omitempty"`
}

// NewDiskBackupRestoreMember instantiates a new DiskBackupRestoreMember object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupRestoreMember() *DiskBackupRestoreMember {
	this := DiskBackupRestoreMember{}
	return &this
}

// NewDiskBackupRestoreMemberWithDefaults instantiates a new DiskBackupRestoreMember object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupRestoreMemberWithDefaults() *DiskBackupRestoreMember {
	this := DiskBackupRestoreMember{}
	return &this
}

// GetDownloadUrl returns the DownloadUrl field value if set, zero value otherwise
func (o *DiskBackupRestoreMember) GetDownloadUrl() string {
	if o == nil || IsNil(o.DownloadUrl) {
		var ret string
		return ret
	}
	return *o.DownloadUrl
}

// GetDownloadUrlOk returns a tuple with the DownloadUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupRestoreMember) GetDownloadUrlOk() (*string, bool) {
	if o == nil || IsNil(o.DownloadUrl) {
		return nil, false
	}

	return o.DownloadUrl, true
}

// HasDownloadUrl returns a boolean if a field has been set.
func (o *DiskBackupRestoreMember) HasDownloadUrl() bool {
	if o != nil && !IsNil(o.DownloadUrl) {
		return true
	}

	return false
}

// SetDownloadUrl gets a reference to the given string and assigns it to the DownloadUrl field.
func (o *DiskBackupRestoreMember) SetDownloadUrl(v string) {
	o.DownloadUrl = &v
}

// GetPrivateDownloadDeliveryUrls returns the PrivateDownloadDeliveryUrls field value if set, zero value otherwise
func (o *DiskBackupRestoreMember) GetPrivateDownloadDeliveryUrls() []ApiPrivateDownloadDeliveryUrl {
	if o == nil || IsNil(o.PrivateDownloadDeliveryUrls) {
		var ret []ApiPrivateDownloadDeliveryUrl
		return ret
	}
	return *o.PrivateDownloadDeliveryUrls
}

// GetPrivateDownloadDeliveryUrlsOk returns a tuple with the PrivateDownloadDeliveryUrls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupRestoreMember) GetPrivateDownloadDeliveryUrlsOk() (*[]ApiPrivateDownloadDeliveryUrl, bool) {
	if o == nil || IsNil(o.PrivateDownloadDeliveryUrls) {
		return nil, false
	}

	return o.PrivateDownloadDeliveryUrls, true
}

// HasPrivateDownloadDeliveryUrls returns a boolean if a field has been set.
func (o *DiskBackupRestoreMember) HasPrivateDownloadDeliveryUrls() bool {
	if o != nil && !IsNil(o.PrivateDownloadDeliveryUrls) {
		return true
	}

	return false
}

// SetPrivateDownloadDeliveryUrls gets a reference to the given []ApiPrivateDownloadDeliveryUrl and assigns it to the PrivateDownloadDeliveryUrls field.
func (o *DiskBackupRestoreMember) SetPrivateDownloadDeliveryUrls(v []ApiPrivateDownloadDeliveryUrl) {
	o.PrivateDownloadDeliveryUrls = &v
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *DiskBackupRestoreMember) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupRestoreMember) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *DiskBackupRestoreMember) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *DiskBackupRestoreMember) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}
