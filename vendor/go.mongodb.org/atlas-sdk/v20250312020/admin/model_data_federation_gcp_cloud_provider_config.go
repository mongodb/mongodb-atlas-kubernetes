// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataFederationGCPCloudProviderConfig Configuration for running Data Federation in GCP.
type DataFederationGCPCloudProviderConfig struct {
	// The email address of the Google Cloud Platform (GCP) service account created by Atlas which should be authorized to allow Atlas to access Google Cloud Storage.
	// Read only field.
	GcpServiceAccount *string `json:"gcpServiceAccount,omitempty"`
	// Unique identifier of the role that Data Federation can use to access the data stores. Required if specifying `cloudProviderConfig`.
	RoleId string `json:"roleId"`
}

// NewDataFederationGCPCloudProviderConfig instantiates a new DataFederationGCPCloudProviderConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataFederationGCPCloudProviderConfig(roleId string) *DataFederationGCPCloudProviderConfig {
	this := DataFederationGCPCloudProviderConfig{}
	this.RoleId = roleId
	return &this
}

// NewDataFederationGCPCloudProviderConfigWithDefaults instantiates a new DataFederationGCPCloudProviderConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataFederationGCPCloudProviderConfigWithDefaults() *DataFederationGCPCloudProviderConfig {
	this := DataFederationGCPCloudProviderConfig{}
	return &this
}

// GetGcpServiceAccount returns the GcpServiceAccount field value if set, zero value otherwise
func (o *DataFederationGCPCloudProviderConfig) GetGcpServiceAccount() string {
	if o == nil || IsNil(o.GcpServiceAccount) {
		var ret string
		return ret
	}
	return *o.GcpServiceAccount
}

// GetGcpServiceAccountOk returns a tuple with the GcpServiceAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataFederationGCPCloudProviderConfig) GetGcpServiceAccountOk() (*string, bool) {
	if o == nil || IsNil(o.GcpServiceAccount) {
		return nil, false
	}

	return o.GcpServiceAccount, true
}

// HasGcpServiceAccount returns a boolean if a field has been set.
func (o *DataFederationGCPCloudProviderConfig) HasGcpServiceAccount() bool {
	if o != nil && !IsNil(o.GcpServiceAccount) {
		return true
	}

	return false
}

// SetGcpServiceAccount gets a reference to the given string and assigns it to the GcpServiceAccount field.
func (o *DataFederationGCPCloudProviderConfig) SetGcpServiceAccount(v string) {
	o.GcpServiceAccount = &v
}

// GetRoleId returns the RoleId field value
func (o *DataFederationGCPCloudProviderConfig) GetRoleId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value
// and a boolean to check if the value has been set.
func (o *DataFederationGCPCloudProviderConfig) GetRoleIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RoleId, true
}

// SetRoleId sets field value
func (o *DataFederationGCPCloudProviderConfig) SetRoleId(v string) {
	o.RoleId = v
}
