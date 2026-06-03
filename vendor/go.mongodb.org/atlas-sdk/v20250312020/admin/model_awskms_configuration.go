// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AWSKMSConfiguration Amazon Web Services (AWS) KMS configuration details and encryption at rest configuration set for the specified project.
type AWSKMSConfiguration struct {
	// Unique alphanumeric string that identifies an Identity and Access Management (IAM) access key with permissions required to access your Amazon Web Services (AWS) Customer Master Key (CMK).
	AccessKeyID *string `json:"accessKeyID,omitempty"`
	// Unique alphanumeric string that identifies the Amazon Web Services (AWS) Customer Master Key (CMK) you used to encrypt and decrypt the MongoDB master keys.
	CustomerMasterKeyID *string `json:"customerMasterKeyID,omitempty"`
	// Flag that indicates whether someone enabled encryption at rest for the specified project through Amazon Web Services (AWS) Key Management Service (KMS). To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// Physical location where MongoDB Cloud deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases. When MongoDB Cloud deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Cloud creates them as part of the deployment. MongoDB Cloud assigns the VPC a CIDR block. To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.
	Region *string `json:"region,omitempty"`
	// Enable connection to your Amazon Web Services (AWS) Key Management Service (KMS) over private networking.
	RequirePrivateNetworking *bool `json:"requirePrivateNetworking,omitempty"`
	// Unique 24-hexadecimal digit string that identifies an Amazon Web Services (AWS) Identity and Access Management (IAM) role. This IAM role has the permissions required to manage your AWS customer master key.
	// Write only field.
	RoleId *string `json:"roleId,omitempty"`
	// Human-readable label of the Identity and Access Management (IAM) secret access key with permissions required to access your Amazon Web Services (AWS) customer master key.
	// Write only field.
	SecretAccessKey *string `json:"secretAccessKey,omitempty"`
	// Flag that indicates whether the Amazon Web Services (AWS) Key Management Service (KMS) encryption key can encrypt and decrypt data.
	// Read only field.
	Valid *bool `json:"valid,omitempty"`
}

// NewAWSKMSConfiguration instantiates a new AWSKMSConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAWSKMSConfiguration() *AWSKMSConfiguration {
	this := AWSKMSConfiguration{}
	return &this
}

// NewAWSKMSConfigurationWithDefaults instantiates a new AWSKMSConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAWSKMSConfigurationWithDefaults() *AWSKMSConfiguration {
	this := AWSKMSConfiguration{}
	return &this
}

// GetAccessKeyID returns the AccessKeyID field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetAccessKeyID() string {
	if o == nil || IsNil(o.AccessKeyID) {
		var ret string
		return ret
	}
	return *o.AccessKeyID
}

// GetAccessKeyIDOk returns a tuple with the AccessKeyID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetAccessKeyIDOk() (*string, bool) {
	if o == nil || IsNil(o.AccessKeyID) {
		return nil, false
	}

	return o.AccessKeyID, true
}

// HasAccessKeyID returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasAccessKeyID() bool {
	if o != nil && !IsNil(o.AccessKeyID) {
		return true
	}

	return false
}

// SetAccessKeyID gets a reference to the given string and assigns it to the AccessKeyID field.
func (o *AWSKMSConfiguration) SetAccessKeyID(v string) {
	o.AccessKeyID = &v
}

// GetCustomerMasterKeyID returns the CustomerMasterKeyID field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetCustomerMasterKeyID() string {
	if o == nil || IsNil(o.CustomerMasterKeyID) {
		var ret string
		return ret
	}
	return *o.CustomerMasterKeyID
}

// GetCustomerMasterKeyIDOk returns a tuple with the CustomerMasterKeyID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetCustomerMasterKeyIDOk() (*string, bool) {
	if o == nil || IsNil(o.CustomerMasterKeyID) {
		return nil, false
	}

	return o.CustomerMasterKeyID, true
}

// HasCustomerMasterKeyID returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasCustomerMasterKeyID() bool {
	if o != nil && !IsNil(o.CustomerMasterKeyID) {
		return true
	}

	return false
}

// SetCustomerMasterKeyID gets a reference to the given string and assigns it to the CustomerMasterKeyID field.
func (o *AWSKMSConfiguration) SetCustomerMasterKeyID(v string) {
	o.CustomerMasterKeyID = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *AWSKMSConfiguration) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *AWSKMSConfiguration) SetRegion(v string) {
	o.Region = &v
}

// GetRequirePrivateNetworking returns the RequirePrivateNetworking field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetRequirePrivateNetworking() bool {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		var ret bool
		return ret
	}
	return *o.RequirePrivateNetworking
}

// GetRequirePrivateNetworkingOk returns a tuple with the RequirePrivateNetworking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetRequirePrivateNetworkingOk() (*bool, bool) {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		return nil, false
	}

	return o.RequirePrivateNetworking, true
}

// HasRequirePrivateNetworking returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasRequirePrivateNetworking() bool {
	if o != nil && !IsNil(o.RequirePrivateNetworking) {
		return true
	}

	return false
}

// SetRequirePrivateNetworking gets a reference to the given bool and assigns it to the RequirePrivateNetworking field.
func (o *AWSKMSConfiguration) SetRequirePrivateNetworking(v bool) {
	o.RequirePrivateNetworking = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *AWSKMSConfiguration) SetRoleId(v string) {
	o.RoleId = &v
}

// GetSecretAccessKey returns the SecretAccessKey field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetSecretAccessKey() string {
	if o == nil || IsNil(o.SecretAccessKey) {
		var ret string
		return ret
	}
	return *o.SecretAccessKey
}

// GetSecretAccessKeyOk returns a tuple with the SecretAccessKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetSecretAccessKeyOk() (*string, bool) {
	if o == nil || IsNil(o.SecretAccessKey) {
		return nil, false
	}

	return o.SecretAccessKey, true
}

// HasSecretAccessKey returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasSecretAccessKey() bool {
	if o != nil && !IsNil(o.SecretAccessKey) {
		return true
	}

	return false
}

// SetSecretAccessKey gets a reference to the given string and assigns it to the SecretAccessKey field.
func (o *AWSKMSConfiguration) SetSecretAccessKey(v string) {
	o.SecretAccessKey = &v
}

// GetValid returns the Valid field value if set, zero value otherwise
func (o *AWSKMSConfiguration) GetValid() bool {
	if o == nil || IsNil(o.Valid) {
		var ret bool
		return ret
	}
	return *o.Valid
}

// GetValidOk returns a tuple with the Valid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AWSKMSConfiguration) GetValidOk() (*bool, bool) {
	if o == nil || IsNil(o.Valid) {
		return nil, false
	}

	return o.Valid, true
}

// HasValid returns a boolean if a field has been set.
func (o *AWSKMSConfiguration) HasValid() bool {
	if o != nil && !IsNil(o.Valid) {
		return true
	}

	return false
}

// SetValid gets a reference to the given bool and assigns it to the Valid field.
func (o *AWSKMSConfiguration) SetValid(v bool) {
	o.Valid = &v
}
