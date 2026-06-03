// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// CloudProviderAccessAWSIAMRole Details that describe the features linked to the Amazon Web Services (AWS) Identity and Access Management (IAM) role.
type CloudProviderAccessAWSIAMRole struct {
	// Amazon Resource Name that identifies the Amazon Web Services (AWS) user account that MongoDB Cloud uses when it assumes the Identity and Access Management (IAM) role.
	// Read only field.
	AtlasAWSAccountArn *string `json:"atlasAWSAccountArn,omitempty"`
	// Unique external ID that MongoDB Cloud uses when it assumes the IAM role in your Amazon Web Services (AWS) account.
	// Read only field.
	AtlasAssumedRoleExternalId *string `json:"atlasAssumedRoleExternalId,omitempty"`
	// Date and time when someone authorized this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	AuthorizedDate *time.Time `json:"authorizedDate,omitempty"`
	// Date and time when someone created this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	// List that contains application features associated with this Amazon Web Services (AWS) Identity and Access Management (IAM) role.
	// Read only field.
	FeatureUsages *[]CloudProviderAccessFeatureUsage `json:"featureUsages,omitempty"`
	// Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account.
	IamAssumedRoleArn *string `json:"iamAssumedRoleArn,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the role.
	// Read only field.
	RoleId *string `json:"roleId,omitempty"`
	// Human-readable label that identifies the cloud provider of the role.
	ProviderName string `json:"providerName"`
	// Unique 24-hexadecimal digit string that identifies the role.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Azure Active Directory Application ID of Atlas. This field is optional and will be derived from the Azure subscription if not provided.
	AtlasAzureAppId *string `json:"atlasAzureAppId,omitempty"`
	// Date and time when this Azure Service Principal was last updated. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastUpdatedDate *time.Time `json:"lastUpdatedDate,omitempty"`
	// UUID string that identifies the Azure Service Principal.
	ServicePrincipalId *string `json:"servicePrincipalId,omitempty"`
	// UUID String that identifies the Azure Active Directory Tenant ID.
	TenantId *string `json:"tenantId,omitempty"`
	// Email address for the Google Service Account created by Atlas.
	GcpServiceAccountForAtlas *string `json:"gcpServiceAccountForAtlas,omitempty"`
	// Provision status of the service account.
	// Read only field.
	Status *string `json:"status,omitempty"`
}

// NewCloudProviderAccessAWSIAMRole instantiates a new CloudProviderAccessAWSIAMRole object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderAccessAWSIAMRole(providerName string) *CloudProviderAccessAWSIAMRole {
	this := CloudProviderAccessAWSIAMRole{}
	this.ProviderName = providerName
	return &this
}

// NewCloudProviderAccessAWSIAMRoleWithDefaults instantiates a new CloudProviderAccessAWSIAMRole object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderAccessAWSIAMRoleWithDefaults() *CloudProviderAccessAWSIAMRole {
	this := CloudProviderAccessAWSIAMRole{}
	return &this
}

// GetAtlasAWSAccountArn returns the AtlasAWSAccountArn field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetAtlasAWSAccountArn() string {
	if o == nil || IsNil(o.AtlasAWSAccountArn) {
		var ret string
		return ret
	}
	return *o.AtlasAWSAccountArn
}

// GetAtlasAWSAccountArnOk returns a tuple with the AtlasAWSAccountArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetAtlasAWSAccountArnOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasAWSAccountArn) {
		return nil, false
	}

	return o.AtlasAWSAccountArn, true
}

// HasAtlasAWSAccountArn returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasAtlasAWSAccountArn() bool {
	if o != nil && !IsNil(o.AtlasAWSAccountArn) {
		return true
	}

	return false
}

// SetAtlasAWSAccountArn gets a reference to the given string and assigns it to the AtlasAWSAccountArn field.
func (o *CloudProviderAccessAWSIAMRole) SetAtlasAWSAccountArn(v string) {
	o.AtlasAWSAccountArn = &v
}

// GetAtlasAssumedRoleExternalId returns the AtlasAssumedRoleExternalId field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetAtlasAssumedRoleExternalId() string {
	if o == nil || IsNil(o.AtlasAssumedRoleExternalId) {
		var ret string
		return ret
	}
	return *o.AtlasAssumedRoleExternalId
}

// GetAtlasAssumedRoleExternalIdOk returns a tuple with the AtlasAssumedRoleExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetAtlasAssumedRoleExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasAssumedRoleExternalId) {
		return nil, false
	}

	return o.AtlasAssumedRoleExternalId, true
}

// HasAtlasAssumedRoleExternalId returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasAtlasAssumedRoleExternalId() bool {
	if o != nil && !IsNil(o.AtlasAssumedRoleExternalId) {
		return true
	}

	return false
}

// SetAtlasAssumedRoleExternalId gets a reference to the given string and assigns it to the AtlasAssumedRoleExternalId field.
func (o *CloudProviderAccessAWSIAMRole) SetAtlasAssumedRoleExternalId(v string) {
	o.AtlasAssumedRoleExternalId = &v
}

// GetAuthorizedDate returns the AuthorizedDate field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetAuthorizedDate() time.Time {
	if o == nil || IsNil(o.AuthorizedDate) {
		var ret time.Time
		return ret
	}
	return *o.AuthorizedDate
}

// GetAuthorizedDateOk returns a tuple with the AuthorizedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetAuthorizedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.AuthorizedDate) {
		return nil, false
	}

	return o.AuthorizedDate, true
}

// HasAuthorizedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasAuthorizedDate() bool {
	if o != nil && !IsNil(o.AuthorizedDate) {
		return true
	}

	return false
}

// SetAuthorizedDate gets a reference to the given time.Time and assigns it to the AuthorizedDate field.
func (o *CloudProviderAccessAWSIAMRole) SetAuthorizedDate(v time.Time) {
	o.AuthorizedDate = &v
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return *o.CreatedDate
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}

	return o.CreatedDate, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *CloudProviderAccessAWSIAMRole) SetCreatedDate(v time.Time) {
	o.CreatedDate = &v
}

// GetFeatureUsages returns the FeatureUsages field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetFeatureUsages() []CloudProviderAccessFeatureUsage {
	if o == nil || IsNil(o.FeatureUsages) {
		var ret []CloudProviderAccessFeatureUsage
		return ret
	}
	return *o.FeatureUsages
}

// GetFeatureUsagesOk returns a tuple with the FeatureUsages field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetFeatureUsagesOk() (*[]CloudProviderAccessFeatureUsage, bool) {
	if o == nil || IsNil(o.FeatureUsages) {
		return nil, false
	}

	return o.FeatureUsages, true
}

// HasFeatureUsages returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasFeatureUsages() bool {
	if o != nil && !IsNil(o.FeatureUsages) {
		return true
	}

	return false
}

// SetFeatureUsages gets a reference to the given []CloudProviderAccessFeatureUsage and assigns it to the FeatureUsages field.
func (o *CloudProviderAccessAWSIAMRole) SetFeatureUsages(v []CloudProviderAccessFeatureUsage) {
	o.FeatureUsages = &v
}

// GetIamAssumedRoleArn returns the IamAssumedRoleArn field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetIamAssumedRoleArn() string {
	if o == nil || IsNil(o.IamAssumedRoleArn) {
		var ret string
		return ret
	}
	return *o.IamAssumedRoleArn
}

// GetIamAssumedRoleArnOk returns a tuple with the IamAssumedRoleArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetIamAssumedRoleArnOk() (*string, bool) {
	if o == nil || IsNil(o.IamAssumedRoleArn) {
		return nil, false
	}

	return o.IamAssumedRoleArn, true
}

// HasIamAssumedRoleArn returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasIamAssumedRoleArn() bool {
	if o != nil && !IsNil(o.IamAssumedRoleArn) {
		return true
	}

	return false
}

// SetIamAssumedRoleArn gets a reference to the given string and assigns it to the IamAssumedRoleArn field.
func (o *CloudProviderAccessAWSIAMRole) SetIamAssumedRoleArn(v string) {
	o.IamAssumedRoleArn = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *CloudProviderAccessAWSIAMRole) SetRoleId(v string) {
	o.RoleId = &v
}

// GetProviderName returns the ProviderName field value
func (o *CloudProviderAccessAWSIAMRole) GetProviderName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetProviderNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderName, true
}

// SetProviderName sets field value
func (o *CloudProviderAccessAWSIAMRole) SetProviderName(v string) {
	o.ProviderName = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CloudProviderAccessAWSIAMRole) SetId(v string) {
	o.Id = &v
}

// GetAtlasAzureAppId returns the AtlasAzureAppId field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetAtlasAzureAppId() string {
	if o == nil || IsNil(o.AtlasAzureAppId) {
		var ret string
		return ret
	}
	return *o.AtlasAzureAppId
}

// GetAtlasAzureAppIdOk returns a tuple with the AtlasAzureAppId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetAtlasAzureAppIdOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasAzureAppId) {
		return nil, false
	}

	return o.AtlasAzureAppId, true
}

// HasAtlasAzureAppId returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasAtlasAzureAppId() bool {
	if o != nil && !IsNil(o.AtlasAzureAppId) {
		return true
	}

	return false
}

// SetAtlasAzureAppId gets a reference to the given string and assigns it to the AtlasAzureAppId field.
func (o *CloudProviderAccessAWSIAMRole) SetAtlasAzureAppId(v string) {
	o.AtlasAzureAppId = &v
}

// GetLastUpdatedDate returns the LastUpdatedDate field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetLastUpdatedDate() time.Time {
	if o == nil || IsNil(o.LastUpdatedDate) {
		var ret time.Time
		return ret
	}
	return *o.LastUpdatedDate
}

// GetLastUpdatedDateOk returns a tuple with the LastUpdatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetLastUpdatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastUpdatedDate) {
		return nil, false
	}

	return o.LastUpdatedDate, true
}

// HasLastUpdatedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasLastUpdatedDate() bool {
	if o != nil && !IsNil(o.LastUpdatedDate) {
		return true
	}

	return false
}

// SetLastUpdatedDate gets a reference to the given time.Time and assigns it to the LastUpdatedDate field.
func (o *CloudProviderAccessAWSIAMRole) SetLastUpdatedDate(v time.Time) {
	o.LastUpdatedDate = &v
}

// GetServicePrincipalId returns the ServicePrincipalId field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetServicePrincipalId() string {
	if o == nil || IsNil(o.ServicePrincipalId) {
		var ret string
		return ret
	}
	return *o.ServicePrincipalId
}

// GetServicePrincipalIdOk returns a tuple with the ServicePrincipalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetServicePrincipalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ServicePrincipalId) {
		return nil, false
	}

	return o.ServicePrincipalId, true
}

// HasServicePrincipalId returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasServicePrincipalId() bool {
	if o != nil && !IsNil(o.ServicePrincipalId) {
		return true
	}

	return false
}

// SetServicePrincipalId gets a reference to the given string and assigns it to the ServicePrincipalId field.
func (o *CloudProviderAccessAWSIAMRole) SetServicePrincipalId(v string) {
	o.ServicePrincipalId = &v
}

// GetTenantId returns the TenantId field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetTenantId() string {
	if o == nil || IsNil(o.TenantId) {
		var ret string
		return ret
	}
	return *o.TenantId
}

// GetTenantIdOk returns a tuple with the TenantId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetTenantIdOk() (*string, bool) {
	if o == nil || IsNil(o.TenantId) {
		return nil, false
	}

	return o.TenantId, true
}

// HasTenantId returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasTenantId() bool {
	if o != nil && !IsNil(o.TenantId) {
		return true
	}

	return false
}

// SetTenantId gets a reference to the given string and assigns it to the TenantId field.
func (o *CloudProviderAccessAWSIAMRole) SetTenantId(v string) {
	o.TenantId = &v
}

// GetGcpServiceAccountForAtlas returns the GcpServiceAccountForAtlas field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetGcpServiceAccountForAtlas() string {
	if o == nil || IsNil(o.GcpServiceAccountForAtlas) {
		var ret string
		return ret
	}
	return *o.GcpServiceAccountForAtlas
}

// GetGcpServiceAccountForAtlasOk returns a tuple with the GcpServiceAccountForAtlas field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetGcpServiceAccountForAtlasOk() (*string, bool) {
	if o == nil || IsNil(o.GcpServiceAccountForAtlas) {
		return nil, false
	}

	return o.GcpServiceAccountForAtlas, true
}

// HasGcpServiceAccountForAtlas returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasGcpServiceAccountForAtlas() bool {
	if o != nil && !IsNil(o.GcpServiceAccountForAtlas) {
		return true
	}

	return false
}

// SetGcpServiceAccountForAtlas gets a reference to the given string and assigns it to the GcpServiceAccountForAtlas field.
func (o *CloudProviderAccessAWSIAMRole) SetGcpServiceAccountForAtlas(v string) {
	o.GcpServiceAccountForAtlas = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *CloudProviderAccessAWSIAMRole) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessAWSIAMRole) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *CloudProviderAccessAWSIAMRole) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *CloudProviderAccessAWSIAMRole) SetStatus(v string) {
	o.Status = &v
}
