// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudProviderAccessRoles struct for CloudProviderAccessRoles
type CloudProviderAccessRoles struct {
	// List that contains the Amazon Web Services (AWS) IAM roles registered and authorized with MongoDB Cloud.
	AwsIamRoles *[]CloudProviderAccessAWSIAMRole `json:"awsIamRoles,omitempty"`
	// List that contains the Azure Service Principals registered with MongoDB Cloud.
	AzureServicePrincipals *[]CloudProviderAccessAzureServicePrincipal `json:"azureServicePrincipals,omitempty"`
	// List that contains the Google Service Accounts registered and authorized with MongoDB Cloud.
	GcpServiceAccounts *[]CloudProviderAccessGCPServiceAccount `json:"gcpServiceAccounts,omitempty"`
}

// NewCloudProviderAccessRoles instantiates a new CloudProviderAccessRoles object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderAccessRoles() *CloudProviderAccessRoles {
	this := CloudProviderAccessRoles{}
	return &this
}

// NewCloudProviderAccessRolesWithDefaults instantiates a new CloudProviderAccessRoles object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderAccessRolesWithDefaults() *CloudProviderAccessRoles {
	this := CloudProviderAccessRoles{}
	return &this
}

// GetAwsIamRoles returns the AwsIamRoles field value if set, zero value otherwise
func (o *CloudProviderAccessRoles) GetAwsIamRoles() []CloudProviderAccessAWSIAMRole {
	if o == nil || IsNil(o.AwsIamRoles) {
		var ret []CloudProviderAccessAWSIAMRole
		return ret
	}
	return *o.AwsIamRoles
}

// GetAwsIamRolesOk returns a tuple with the AwsIamRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoles) GetAwsIamRolesOk() (*[]CloudProviderAccessAWSIAMRole, bool) {
	if o == nil || IsNil(o.AwsIamRoles) {
		return nil, false
	}

	return o.AwsIamRoles, true
}

// HasAwsIamRoles returns a boolean if a field has been set.
func (o *CloudProviderAccessRoles) HasAwsIamRoles() bool {
	if o != nil && !IsNil(o.AwsIamRoles) {
		return true
	}

	return false
}

// SetAwsIamRoles gets a reference to the given []CloudProviderAccessAWSIAMRole and assigns it to the AwsIamRoles field.
func (o *CloudProviderAccessRoles) SetAwsIamRoles(v []CloudProviderAccessAWSIAMRole) {
	o.AwsIamRoles = &v
}

// GetAzureServicePrincipals returns the AzureServicePrincipals field value if set, zero value otherwise
func (o *CloudProviderAccessRoles) GetAzureServicePrincipals() []CloudProviderAccessAzureServicePrincipal {
	if o == nil || IsNil(o.AzureServicePrincipals) {
		var ret []CloudProviderAccessAzureServicePrincipal
		return ret
	}
	return *o.AzureServicePrincipals
}

// GetAzureServicePrincipalsOk returns a tuple with the AzureServicePrincipals field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoles) GetAzureServicePrincipalsOk() (*[]CloudProviderAccessAzureServicePrincipal, bool) {
	if o == nil || IsNil(o.AzureServicePrincipals) {
		return nil, false
	}

	return o.AzureServicePrincipals, true
}

// HasAzureServicePrincipals returns a boolean if a field has been set.
func (o *CloudProviderAccessRoles) HasAzureServicePrincipals() bool {
	if o != nil && !IsNil(o.AzureServicePrincipals) {
		return true
	}

	return false
}

// SetAzureServicePrincipals gets a reference to the given []CloudProviderAccessAzureServicePrincipal and assigns it to the AzureServicePrincipals field.
func (o *CloudProviderAccessRoles) SetAzureServicePrincipals(v []CloudProviderAccessAzureServicePrincipal) {
	o.AzureServicePrincipals = &v
}

// GetGcpServiceAccounts returns the GcpServiceAccounts field value if set, zero value otherwise
func (o *CloudProviderAccessRoles) GetGcpServiceAccounts() []CloudProviderAccessGCPServiceAccount {
	if o == nil || IsNil(o.GcpServiceAccounts) {
		var ret []CloudProviderAccessGCPServiceAccount
		return ret
	}
	return *o.GcpServiceAccounts
}

// GetGcpServiceAccountsOk returns a tuple with the GcpServiceAccounts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoles) GetGcpServiceAccountsOk() (*[]CloudProviderAccessGCPServiceAccount, bool) {
	if o == nil || IsNil(o.GcpServiceAccounts) {
		return nil, false
	}

	return o.GcpServiceAccounts, true
}

// HasGcpServiceAccounts returns a boolean if a field has been set.
func (o *CloudProviderAccessRoles) HasGcpServiceAccounts() bool {
	if o != nil && !IsNil(o.GcpServiceAccounts) {
		return true
	}

	return false
}

// SetGcpServiceAccounts gets a reference to the given []CloudProviderAccessGCPServiceAccount and assigns it to the GcpServiceAccounts field.
func (o *CloudProviderAccessRoles) SetGcpServiceAccounts(v []CloudProviderAccessGCPServiceAccount) {
	o.GcpServiceAccounts = &v
}
