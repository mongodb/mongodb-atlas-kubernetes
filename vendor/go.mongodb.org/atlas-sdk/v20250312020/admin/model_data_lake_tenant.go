// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeTenant struct for DataLakeTenant
type DataLakeTenant struct {
	CloudProviderConfig *DataLakeCloudProviderConfig `json:"cloudProviderConfig,omitempty"`
	DataProcessRegion   *DataLakeDataProcessRegion   `json:"dataProcessRegion,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List that contains the hostnames assigned to the Federated Database Instance.
	// Read only field.
	Hostnames *[]string `json:"hostnames,omitempty"`
	// Human-readable label that identifies the Federated Database Instance.
	Name *string `json:"name,omitempty"`
	// List that contains the sets of private endpoints and hostnames.
	// Read only field.
	PrivateEndpointHostnames *[]PrivateEndpointHostname `json:"privateEndpointHostnames,omitempty"`
	// Label that indicates the status of the Federated Database Instance.
	// Read only field.
	State   *string          `json:"state,omitempty"`
	Storage *DataLakeStorage `json:"storage,omitempty"`
}

// NewDataLakeTenant instantiates a new DataLakeTenant object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeTenant() *DataLakeTenant {
	this := DataLakeTenant{}
	return &this
}

// NewDataLakeTenantWithDefaults instantiates a new DataLakeTenant object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeTenantWithDefaults() *DataLakeTenant {
	this := DataLakeTenant{}
	return &this
}

// GetCloudProviderConfig returns the CloudProviderConfig field value if set, zero value otherwise
func (o *DataLakeTenant) GetCloudProviderConfig() DataLakeCloudProviderConfig {
	if o == nil || IsNil(o.CloudProviderConfig) {
		var ret DataLakeCloudProviderConfig
		return ret
	}
	return *o.CloudProviderConfig
}

// GetCloudProviderConfigOk returns a tuple with the CloudProviderConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetCloudProviderConfigOk() (*DataLakeCloudProviderConfig, bool) {
	if o == nil || IsNil(o.CloudProviderConfig) {
		return nil, false
	}

	return o.CloudProviderConfig, true
}

// HasCloudProviderConfig returns a boolean if a field has been set.
func (o *DataLakeTenant) HasCloudProviderConfig() bool {
	if o != nil && !IsNil(o.CloudProviderConfig) {
		return true
	}

	return false
}

// SetCloudProviderConfig gets a reference to the given DataLakeCloudProviderConfig and assigns it to the CloudProviderConfig field.
func (o *DataLakeTenant) SetCloudProviderConfig(v DataLakeCloudProviderConfig) {
	o.CloudProviderConfig = &v
}

// GetDataProcessRegion returns the DataProcessRegion field value if set, zero value otherwise
func (o *DataLakeTenant) GetDataProcessRegion() DataLakeDataProcessRegion {
	if o == nil || IsNil(o.DataProcessRegion) {
		var ret DataLakeDataProcessRegion
		return ret
	}
	return *o.DataProcessRegion
}

// GetDataProcessRegionOk returns a tuple with the DataProcessRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetDataProcessRegionOk() (*DataLakeDataProcessRegion, bool) {
	if o == nil || IsNil(o.DataProcessRegion) {
		return nil, false
	}

	return o.DataProcessRegion, true
}

// HasDataProcessRegion returns a boolean if a field has been set.
func (o *DataLakeTenant) HasDataProcessRegion() bool {
	if o != nil && !IsNil(o.DataProcessRegion) {
		return true
	}

	return false
}

// SetDataProcessRegion gets a reference to the given DataLakeDataProcessRegion and assigns it to the DataProcessRegion field.
func (o *DataLakeTenant) SetDataProcessRegion(v DataLakeDataProcessRegion) {
	o.DataProcessRegion = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *DataLakeTenant) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *DataLakeTenant) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *DataLakeTenant) SetGroupId(v string) {
	o.GroupId = &v
}

// GetHostnames returns the Hostnames field value if set, zero value otherwise
func (o *DataLakeTenant) GetHostnames() []string {
	if o == nil || IsNil(o.Hostnames) {
		var ret []string
		return ret
	}
	return *o.Hostnames
}

// GetHostnamesOk returns a tuple with the Hostnames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetHostnamesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Hostnames) {
		return nil, false
	}

	return o.Hostnames, true
}

// HasHostnames returns a boolean if a field has been set.
func (o *DataLakeTenant) HasHostnames() bool {
	if o != nil && !IsNil(o.Hostnames) {
		return true
	}

	return false
}

// SetHostnames gets a reference to the given []string and assigns it to the Hostnames field.
func (o *DataLakeTenant) SetHostnames(v []string) {
	o.Hostnames = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DataLakeTenant) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataLakeTenant) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataLakeTenant) SetName(v string) {
	o.Name = &v
}

// GetPrivateEndpointHostnames returns the PrivateEndpointHostnames field value if set, zero value otherwise
func (o *DataLakeTenant) GetPrivateEndpointHostnames() []PrivateEndpointHostname {
	if o == nil || IsNil(o.PrivateEndpointHostnames) {
		var ret []PrivateEndpointHostname
		return ret
	}
	return *o.PrivateEndpointHostnames
}

// GetPrivateEndpointHostnamesOk returns a tuple with the PrivateEndpointHostnames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetPrivateEndpointHostnamesOk() (*[]PrivateEndpointHostname, bool) {
	if o == nil || IsNil(o.PrivateEndpointHostnames) {
		return nil, false
	}

	return o.PrivateEndpointHostnames, true
}

// HasPrivateEndpointHostnames returns a boolean if a field has been set.
func (o *DataLakeTenant) HasPrivateEndpointHostnames() bool {
	if o != nil && !IsNil(o.PrivateEndpointHostnames) {
		return true
	}

	return false
}

// SetPrivateEndpointHostnames gets a reference to the given []PrivateEndpointHostname and assigns it to the PrivateEndpointHostnames field.
func (o *DataLakeTenant) SetPrivateEndpointHostnames(v []PrivateEndpointHostname) {
	o.PrivateEndpointHostnames = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *DataLakeTenant) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *DataLakeTenant) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *DataLakeTenant) SetState(v string) {
	o.State = &v
}

// GetStorage returns the Storage field value if set, zero value otherwise
func (o *DataLakeTenant) GetStorage() DataLakeStorage {
	if o == nil || IsNil(o.Storage) {
		var ret DataLakeStorage
		return ret
	}
	return *o.Storage
}

// GetStorageOk returns a tuple with the Storage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeTenant) GetStorageOk() (*DataLakeStorage, bool) {
	if o == nil || IsNil(o.Storage) {
		return nil, false
	}

	return o.Storage, true
}

// HasStorage returns a boolean if a field has been set.
func (o *DataLakeTenant) HasStorage() bool {
	if o != nil && !IsNil(o.Storage) {
		return true
	}

	return false
}

// SetStorage gets a reference to the given DataLakeStorage and assigns it to the Storage field.
func (o *DataLakeTenant) SetStorage(v DataLakeStorage) {
	o.Storage = &v
}
