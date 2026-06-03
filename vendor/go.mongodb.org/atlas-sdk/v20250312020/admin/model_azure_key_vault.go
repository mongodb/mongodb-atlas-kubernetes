// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AzureKeyVault Details that define the configuration of Encryption at Rest using Azure Key Vault (AKV).
type AzureKeyVault struct {
	// Azure environment in which your account credentials reside.
	AzureEnvironment *string `json:"azureEnvironment,omitempty"`
	// Unique 36-hexadecimal character string that identifies an Azure application associated with your Azure Active Directory tenant.
	ClientID *string `json:"clientID,omitempty"`
	// Flag that indicates whether someone enabled encryption at rest for the specified  project. To disable encryption at rest using customer key management and remove the configuration details, pass only this parameter with a value of `false`.
	Enabled *bool `json:"enabled,omitempty"`
	// Web address with a unique key that identifies for your Azure Key Vault.
	KeyIdentifier *string `json:"keyIdentifier,omitempty"`
	// Unique string that identifies the Azure Key Vault that contains your key. This field cannot be modified when you enable and set up private endpoint connections to your Azure Key Vault.
	KeyVaultName *string `json:"keyVaultName,omitempty"`
	// Enable connection to your Azure Key Vault over private networking.
	RequirePrivateNetworking *bool `json:"requirePrivateNetworking,omitempty"`
	// Name of the Azure resource group that contains your Azure Key Vault. This field cannot be modified when you enable and set up private endpoint connections to your Azure Key Vault.
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the Azure Service Principal that MongoDB Cloud uses to access the Azure Key Vault.
	RoleId *string `json:"roleId,omitempty"`
	// Private data that you need secured and that belongs to the specified Azure Key Vault (AKV) tenant (`azureKeyVault.tenantID`). This data can include any type of sensitive data such as passwords, database connection strings, API keys, and the like. AKV stores this information as encrypted binary data.
	// Write only field.
	Secret *string `json:"secret,omitempty"`
	// Unique 36-hexadecimal character string that identifies your Azure subscription. This field cannot be modified when you enable and set up private endpoint connections to your Azure Key Vault.
	SubscriptionID *string `json:"subscriptionID,omitempty"`
	// Unique 36-hexadecimal character string that identifies the Azure Active Directory tenant within your Azure subscription.
	TenantID *string `json:"tenantID,omitempty"`
	// Flag that indicates whether the Azure encryption key can encrypt and decrypt data.
	// Read only field.
	Valid *bool `json:"valid,omitempty"`
}

// NewAzureKeyVault instantiates a new AzureKeyVault object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAzureKeyVault() *AzureKeyVault {
	this := AzureKeyVault{}
	return &this
}

// NewAzureKeyVaultWithDefaults instantiates a new AzureKeyVault object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAzureKeyVaultWithDefaults() *AzureKeyVault {
	this := AzureKeyVault{}
	return &this
}

// GetAzureEnvironment returns the AzureEnvironment field value if set, zero value otherwise
func (o *AzureKeyVault) GetAzureEnvironment() string {
	if o == nil || IsNil(o.AzureEnvironment) {
		var ret string
		return ret
	}
	return *o.AzureEnvironment
}

// GetAzureEnvironmentOk returns a tuple with the AzureEnvironment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetAzureEnvironmentOk() (*string, bool) {
	if o == nil || IsNil(o.AzureEnvironment) {
		return nil, false
	}

	return o.AzureEnvironment, true
}

// HasAzureEnvironment returns a boolean if a field has been set.
func (o *AzureKeyVault) HasAzureEnvironment() bool {
	if o != nil && !IsNil(o.AzureEnvironment) {
		return true
	}

	return false
}

// SetAzureEnvironment gets a reference to the given string and assigns it to the AzureEnvironment field.
func (o *AzureKeyVault) SetAzureEnvironment(v string) {
	o.AzureEnvironment = &v
}

// GetClientID returns the ClientID field value if set, zero value otherwise
func (o *AzureKeyVault) GetClientID() string {
	if o == nil || IsNil(o.ClientID) {
		var ret string
		return ret
	}
	return *o.ClientID
}

// GetClientIDOk returns a tuple with the ClientID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetClientIDOk() (*string, bool) {
	if o == nil || IsNil(o.ClientID) {
		return nil, false
	}

	return o.ClientID, true
}

// HasClientID returns a boolean if a field has been set.
func (o *AzureKeyVault) HasClientID() bool {
	if o != nil && !IsNil(o.ClientID) {
		return true
	}

	return false
}

// SetClientID gets a reference to the given string and assigns it to the ClientID field.
func (o *AzureKeyVault) SetClientID(v string) {
	o.ClientID = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *AzureKeyVault) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *AzureKeyVault) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *AzureKeyVault) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetKeyIdentifier returns the KeyIdentifier field value if set, zero value otherwise
func (o *AzureKeyVault) GetKeyIdentifier() string {
	if o == nil || IsNil(o.KeyIdentifier) {
		var ret string
		return ret
	}
	return *o.KeyIdentifier
}

// GetKeyIdentifierOk returns a tuple with the KeyIdentifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetKeyIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.KeyIdentifier) {
		return nil, false
	}

	return o.KeyIdentifier, true
}

// HasKeyIdentifier returns a boolean if a field has been set.
func (o *AzureKeyVault) HasKeyIdentifier() bool {
	if o != nil && !IsNil(o.KeyIdentifier) {
		return true
	}

	return false
}

// SetKeyIdentifier gets a reference to the given string and assigns it to the KeyIdentifier field.
func (o *AzureKeyVault) SetKeyIdentifier(v string) {
	o.KeyIdentifier = &v
}

// GetKeyVaultName returns the KeyVaultName field value if set, zero value otherwise
func (o *AzureKeyVault) GetKeyVaultName() string {
	if o == nil || IsNil(o.KeyVaultName) {
		var ret string
		return ret
	}
	return *o.KeyVaultName
}

// GetKeyVaultNameOk returns a tuple with the KeyVaultName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetKeyVaultNameOk() (*string, bool) {
	if o == nil || IsNil(o.KeyVaultName) {
		return nil, false
	}

	return o.KeyVaultName, true
}

// HasKeyVaultName returns a boolean if a field has been set.
func (o *AzureKeyVault) HasKeyVaultName() bool {
	if o != nil && !IsNil(o.KeyVaultName) {
		return true
	}

	return false
}

// SetKeyVaultName gets a reference to the given string and assigns it to the KeyVaultName field.
func (o *AzureKeyVault) SetKeyVaultName(v string) {
	o.KeyVaultName = &v
}

// GetRequirePrivateNetworking returns the RequirePrivateNetworking field value if set, zero value otherwise
func (o *AzureKeyVault) GetRequirePrivateNetworking() bool {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		var ret bool
		return ret
	}
	return *o.RequirePrivateNetworking
}

// GetRequirePrivateNetworkingOk returns a tuple with the RequirePrivateNetworking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetRequirePrivateNetworkingOk() (*bool, bool) {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		return nil, false
	}

	return o.RequirePrivateNetworking, true
}

// HasRequirePrivateNetworking returns a boolean if a field has been set.
func (o *AzureKeyVault) HasRequirePrivateNetworking() bool {
	if o != nil && !IsNil(o.RequirePrivateNetworking) {
		return true
	}

	return false
}

// SetRequirePrivateNetworking gets a reference to the given bool and assigns it to the RequirePrivateNetworking field.
func (o *AzureKeyVault) SetRequirePrivateNetworking(v bool) {
	o.RequirePrivateNetworking = &v
}

// GetResourceGroupName returns the ResourceGroupName field value if set, zero value otherwise
func (o *AzureKeyVault) GetResourceGroupName() string {
	if o == nil || IsNil(o.ResourceGroupName) {
		var ret string
		return ret
	}
	return *o.ResourceGroupName
}

// GetResourceGroupNameOk returns a tuple with the ResourceGroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetResourceGroupNameOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceGroupName) {
		return nil, false
	}

	return o.ResourceGroupName, true
}

// HasResourceGroupName returns a boolean if a field has been set.
func (o *AzureKeyVault) HasResourceGroupName() bool {
	if o != nil && !IsNil(o.ResourceGroupName) {
		return true
	}

	return false
}

// SetResourceGroupName gets a reference to the given string and assigns it to the ResourceGroupName field.
func (o *AzureKeyVault) SetResourceGroupName(v string) {
	o.ResourceGroupName = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *AzureKeyVault) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *AzureKeyVault) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *AzureKeyVault) SetRoleId(v string) {
	o.RoleId = &v
}

// GetSecret returns the Secret field value if set, zero value otherwise
func (o *AzureKeyVault) GetSecret() string {
	if o == nil || IsNil(o.Secret) {
		var ret string
		return ret
	}
	return *o.Secret
}

// GetSecretOk returns a tuple with the Secret field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetSecretOk() (*string, bool) {
	if o == nil || IsNil(o.Secret) {
		return nil, false
	}

	return o.Secret, true
}

// HasSecret returns a boolean if a field has been set.
func (o *AzureKeyVault) HasSecret() bool {
	if o != nil && !IsNil(o.Secret) {
		return true
	}

	return false
}

// SetSecret gets a reference to the given string and assigns it to the Secret field.
func (o *AzureKeyVault) SetSecret(v string) {
	o.Secret = &v
}

// GetSubscriptionID returns the SubscriptionID field value if set, zero value otherwise
func (o *AzureKeyVault) GetSubscriptionID() string {
	if o == nil || IsNil(o.SubscriptionID) {
		var ret string
		return ret
	}
	return *o.SubscriptionID
}

// GetSubscriptionIDOk returns a tuple with the SubscriptionID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetSubscriptionIDOk() (*string, bool) {
	if o == nil || IsNil(o.SubscriptionID) {
		return nil, false
	}

	return o.SubscriptionID, true
}

// HasSubscriptionID returns a boolean if a field has been set.
func (o *AzureKeyVault) HasSubscriptionID() bool {
	if o != nil && !IsNil(o.SubscriptionID) {
		return true
	}

	return false
}

// SetSubscriptionID gets a reference to the given string and assigns it to the SubscriptionID field.
func (o *AzureKeyVault) SetSubscriptionID(v string) {
	o.SubscriptionID = &v
}

// GetTenantID returns the TenantID field value if set, zero value otherwise
func (o *AzureKeyVault) GetTenantID() string {
	if o == nil || IsNil(o.TenantID) {
		var ret string
		return ret
	}
	return *o.TenantID
}

// GetTenantIDOk returns a tuple with the TenantID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetTenantIDOk() (*string, bool) {
	if o == nil || IsNil(o.TenantID) {
		return nil, false
	}

	return o.TenantID, true
}

// HasTenantID returns a boolean if a field has been set.
func (o *AzureKeyVault) HasTenantID() bool {
	if o != nil && !IsNil(o.TenantID) {
		return true
	}

	return false
}

// SetTenantID gets a reference to the given string and assigns it to the TenantID field.
func (o *AzureKeyVault) SetTenantID(v string) {
	o.TenantID = &v
}

// GetValid returns the Valid field value if set, zero value otherwise
func (o *AzureKeyVault) GetValid() bool {
	if o == nil || IsNil(o.Valid) {
		var ret bool
		return ret
	}
	return *o.Valid
}

// GetValidOk returns a tuple with the Valid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureKeyVault) GetValidOk() (*bool, bool) {
	if o == nil || IsNil(o.Valid) {
		return nil, false
	}

	return o.Valid, true
}

// HasValid returns a boolean if a field has been set.
func (o *AzureKeyVault) HasValid() bool {
	if o != nil && !IsNil(o.Valid) {
		return true
	}

	return false
}

// SetValid gets a reference to the given bool and assigns it to the Valid field.
func (o *AzureKeyVault) SetValid(v bool) {
	o.Valid = &v
}
