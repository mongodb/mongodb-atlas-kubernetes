// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DeleteCopiedBackups20240805 Deleted copy setting whose backup copies need to also be deleted.
type DeleteCopiedBackups20240805 struct {
	// Human-readable label that identifies the cloud provider for the deleted copy setting whose backup copies you want to delete.
	// Write only field.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// Target region for the deleted copy setting whose backup copies you want to delete. Please supply the 'Atlas Region'.
	// Write only field.
	RegionName *string `json:"regionName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the zone in a cluster. For global clusters, there can be multiple zones to choose from. For sharded clusters and replica set clusters, there is only one zone in the cluster. To find the Zone Id, do a GET request to Return One Cluster from One Project and consult the `replicationSpecs` array.
	// Write only field.
	ZoneId *string `json:"zoneId,omitempty"`
}

// NewDeleteCopiedBackups20240805 instantiates a new DeleteCopiedBackups20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDeleteCopiedBackups20240805() *DeleteCopiedBackups20240805 {
	this := DeleteCopiedBackups20240805{}
	return &this
}

// NewDeleteCopiedBackups20240805WithDefaults instantiates a new DeleteCopiedBackups20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDeleteCopiedBackups20240805WithDefaults() *DeleteCopiedBackups20240805 {
	this := DeleteCopiedBackups20240805{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *DeleteCopiedBackups20240805) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DeleteCopiedBackups20240805) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *DeleteCopiedBackups20240805) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *DeleteCopiedBackups20240805) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetRegionName returns the RegionName field value if set, zero value otherwise
func (o *DeleteCopiedBackups20240805) GetRegionName() string {
	if o == nil || IsNil(o.RegionName) {
		var ret string
		return ret
	}
	return *o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DeleteCopiedBackups20240805) GetRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionName) {
		return nil, false
	}

	return o.RegionName, true
}

// HasRegionName returns a boolean if a field has been set.
func (o *DeleteCopiedBackups20240805) HasRegionName() bool {
	if o != nil && !IsNil(o.RegionName) {
		return true
	}

	return false
}

// SetRegionName gets a reference to the given string and assigns it to the RegionName field.
func (o *DeleteCopiedBackups20240805) SetRegionName(v string) {
	o.RegionName = &v
}

// GetZoneId returns the ZoneId field value if set, zero value otherwise
func (o *DeleteCopiedBackups20240805) GetZoneId() string {
	if o == nil || IsNil(o.ZoneId) {
		var ret string
		return ret
	}
	return *o.ZoneId
}

// GetZoneIdOk returns a tuple with the ZoneId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DeleteCopiedBackups20240805) GetZoneIdOk() (*string, bool) {
	if o == nil || IsNil(o.ZoneId) {
		return nil, false
	}

	return o.ZoneId, true
}

// HasZoneId returns a boolean if a field has been set.
func (o *DeleteCopiedBackups20240805) HasZoneId() bool {
	if o != nil && !IsNil(o.ZoneId) {
		return true
	}

	return false
}

// SetZoneId gets a reference to the given string and assigns it to the ZoneId field.
func (o *DeleteCopiedBackups20240805) SetZoneId(v string) {
	o.ZoneId = &v
}
