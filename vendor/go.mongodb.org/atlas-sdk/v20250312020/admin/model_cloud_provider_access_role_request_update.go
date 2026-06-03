// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// CloudProviderAccessRoleRequestUpdate Cloud provider access role.
type CloudProviderAccessRoleRequestUpdate struct {
	// Human-readable label that identifies the cloud provider of the role.
	ProviderName string `json:"providerName"`
	// Amazon Resource Name that identifies the Amazon Web Services (AWS) user account that MongoDB Cloud uses when it assumes the Identity and Access Management (IAM) role.
	// Read only field.
	AtlasAWSAccountArn *string `json:"atlasAWSAccountArn,omitempty"`
	// Unique external ID that MongoDB Cloud uses when it assumes the IAM role in your Amazon Web Services (AWS) account.
	// Read only field.
	AtlasAssumedRoleExternalId *string `json:"atlasAssumedRoleExternalId,omitempty"`
	// Date and time when someone authorized this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	AuthorizedDate *time.Time `json:"authorizedDate,omitempty"`
	// Date and time when this Azure Service Principal was created. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedDate *time.Time `json:"createdDate,omitempty"`
	// List that contains application features associated with this Azure Service Principal.
	// Read only field.
	FeatureUsages *[]CloudProviderAccessFeatureUsage `json:"featureUsages,omitempty"`
	// Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account.
	IamAssumedRoleArn *string `json:"iamAssumedRoleArn,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the role.
	// Read only field.
	RoleId *string `json:"roleId,omitempty"`
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
}

// NewCloudProviderAccessRoleRequestUpdate instantiates a new CloudProviderAccessRoleRequestUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderAccessRoleRequestUpdate(providerName string) *CloudProviderAccessRoleRequestUpdate {
	this := CloudProviderAccessRoleRequestUpdate{}
	this.ProviderName = providerName
	return &this
}

// NewCloudProviderAccessRoleRequestUpdateWithDefaults instantiates a new CloudProviderAccessRoleRequestUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderAccessRoleRequestUpdateWithDefaults() *CloudProviderAccessRoleRequestUpdate {
	this := CloudProviderAccessRoleRequestUpdate{}
	return &this
}

// GetProviderName returns the ProviderName field value
func (o *CloudProviderAccessRoleRequestUpdate) GetProviderName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetProviderNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderName, true
}

// SetProviderName sets field value
func (o *CloudProviderAccessRoleRequestUpdate) SetProviderName(v string) {
	o.ProviderName = v
}

// GetAtlasAWSAccountArn returns the AtlasAWSAccountArn field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetAtlasAWSAccountArn() string {
	if o == nil || IsNil(o.AtlasAWSAccountArn) {
		var ret string
		return ret
	}
	return *o.AtlasAWSAccountArn
}

// GetAtlasAWSAccountArnOk returns a tuple with the AtlasAWSAccountArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetAtlasAWSAccountArnOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasAWSAccountArn) {
		return nil, false
	}

	return o.AtlasAWSAccountArn, true
}

// HasAtlasAWSAccountArn returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasAtlasAWSAccountArn() bool {
	if o != nil && !IsNil(o.AtlasAWSAccountArn) {
		return true
	}

	return false
}

// SetAtlasAWSAccountArn gets a reference to the given string and assigns it to the AtlasAWSAccountArn field.
func (o *CloudProviderAccessRoleRequestUpdate) SetAtlasAWSAccountArn(v string) {
	o.AtlasAWSAccountArn = &v
}

// GetAtlasAssumedRoleExternalId returns the AtlasAssumedRoleExternalId field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetAtlasAssumedRoleExternalId() string {
	if o == nil || IsNil(o.AtlasAssumedRoleExternalId) {
		var ret string
		return ret
	}
	return *o.AtlasAssumedRoleExternalId
}

// GetAtlasAssumedRoleExternalIdOk returns a tuple with the AtlasAssumedRoleExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetAtlasAssumedRoleExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasAssumedRoleExternalId) {
		return nil, false
	}

	return o.AtlasAssumedRoleExternalId, true
}

// HasAtlasAssumedRoleExternalId returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasAtlasAssumedRoleExternalId() bool {
	if o != nil && !IsNil(o.AtlasAssumedRoleExternalId) {
		return true
	}

	return false
}

// SetAtlasAssumedRoleExternalId gets a reference to the given string and assigns it to the AtlasAssumedRoleExternalId field.
func (o *CloudProviderAccessRoleRequestUpdate) SetAtlasAssumedRoleExternalId(v string) {
	o.AtlasAssumedRoleExternalId = &v
}

// GetAuthorizedDate returns the AuthorizedDate field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetAuthorizedDate() time.Time {
	if o == nil || IsNil(o.AuthorizedDate) {
		var ret time.Time
		return ret
	}
	return *o.AuthorizedDate
}

// GetAuthorizedDateOk returns a tuple with the AuthorizedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetAuthorizedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.AuthorizedDate) {
		return nil, false
	}

	return o.AuthorizedDate, true
}

// HasAuthorizedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasAuthorizedDate() bool {
	if o != nil && !IsNil(o.AuthorizedDate) {
		return true
	}

	return false
}

// SetAuthorizedDate gets a reference to the given time.Time and assigns it to the AuthorizedDate field.
func (o *CloudProviderAccessRoleRequestUpdate) SetAuthorizedDate(v time.Time) {
	o.AuthorizedDate = &v
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return *o.CreatedDate
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}

	return o.CreatedDate, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *CloudProviderAccessRoleRequestUpdate) SetCreatedDate(v time.Time) {
	o.CreatedDate = &v
}

// GetFeatureUsages returns the FeatureUsages field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetFeatureUsages() []CloudProviderAccessFeatureUsage {
	if o == nil || IsNil(o.FeatureUsages) {
		var ret []CloudProviderAccessFeatureUsage
		return ret
	}
	return *o.FeatureUsages
}

// GetFeatureUsagesOk returns a tuple with the FeatureUsages field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetFeatureUsagesOk() (*[]CloudProviderAccessFeatureUsage, bool) {
	if o == nil || IsNil(o.FeatureUsages) {
		return nil, false
	}

	return o.FeatureUsages, true
}

// HasFeatureUsages returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasFeatureUsages() bool {
	if o != nil && !IsNil(o.FeatureUsages) {
		return true
	}

	return false
}

// SetFeatureUsages gets a reference to the given []CloudProviderAccessFeatureUsage and assigns it to the FeatureUsages field.
func (o *CloudProviderAccessRoleRequestUpdate) SetFeatureUsages(v []CloudProviderAccessFeatureUsage) {
	o.FeatureUsages = &v
}

// GetIamAssumedRoleArn returns the IamAssumedRoleArn field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetIamAssumedRoleArn() string {
	if o == nil || IsNil(o.IamAssumedRoleArn) {
		var ret string
		return ret
	}
	return *o.IamAssumedRoleArn
}

// GetIamAssumedRoleArnOk returns a tuple with the IamAssumedRoleArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetIamAssumedRoleArnOk() (*string, bool) {
	if o == nil || IsNil(o.IamAssumedRoleArn) {
		return nil, false
	}

	return o.IamAssumedRoleArn, true
}

// HasIamAssumedRoleArn returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasIamAssumedRoleArn() bool {
	if o != nil && !IsNil(o.IamAssumedRoleArn) {
		return true
	}

	return false
}

// SetIamAssumedRoleArn gets a reference to the given string and assigns it to the IamAssumedRoleArn field.
func (o *CloudProviderAccessRoleRequestUpdate) SetIamAssumedRoleArn(v string) {
	o.IamAssumedRoleArn = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *CloudProviderAccessRoleRequestUpdate) SetRoleId(v string) {
	o.RoleId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CloudProviderAccessRoleRequestUpdate) SetId(v string) {
	o.Id = &v
}

// GetAtlasAzureAppId returns the AtlasAzureAppId field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetAtlasAzureAppId() string {
	if o == nil || IsNil(o.AtlasAzureAppId) {
		var ret string
		return ret
	}
	return *o.AtlasAzureAppId
}

// GetAtlasAzureAppIdOk returns a tuple with the AtlasAzureAppId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetAtlasAzureAppIdOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasAzureAppId) {
		return nil, false
	}

	return o.AtlasAzureAppId, true
}

// HasAtlasAzureAppId returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasAtlasAzureAppId() bool {
	if o != nil && !IsNil(o.AtlasAzureAppId) {
		return true
	}

	return false
}

// SetAtlasAzureAppId gets a reference to the given string and assigns it to the AtlasAzureAppId field.
func (o *CloudProviderAccessRoleRequestUpdate) SetAtlasAzureAppId(v string) {
	o.AtlasAzureAppId = &v
}

// GetLastUpdatedDate returns the LastUpdatedDate field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetLastUpdatedDate() time.Time {
	if o == nil || IsNil(o.LastUpdatedDate) {
		var ret time.Time
		return ret
	}
	return *o.LastUpdatedDate
}

// GetLastUpdatedDateOk returns a tuple with the LastUpdatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetLastUpdatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastUpdatedDate) {
		return nil, false
	}

	return o.LastUpdatedDate, true
}

// HasLastUpdatedDate returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasLastUpdatedDate() bool {
	if o != nil && !IsNil(o.LastUpdatedDate) {
		return true
	}

	return false
}

// SetLastUpdatedDate gets a reference to the given time.Time and assigns it to the LastUpdatedDate field.
func (o *CloudProviderAccessRoleRequestUpdate) SetLastUpdatedDate(v time.Time) {
	o.LastUpdatedDate = &v
}

// GetServicePrincipalId returns the ServicePrincipalId field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetServicePrincipalId() string {
	if o == nil || IsNil(o.ServicePrincipalId) {
		var ret string
		return ret
	}
	return *o.ServicePrincipalId
}

// GetServicePrincipalIdOk returns a tuple with the ServicePrincipalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetServicePrincipalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ServicePrincipalId) {
		return nil, false
	}

	return o.ServicePrincipalId, true
}

// HasServicePrincipalId returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasServicePrincipalId() bool {
	if o != nil && !IsNil(o.ServicePrincipalId) {
		return true
	}

	return false
}

// SetServicePrincipalId gets a reference to the given string and assigns it to the ServicePrincipalId field.
func (o *CloudProviderAccessRoleRequestUpdate) SetServicePrincipalId(v string) {
	o.ServicePrincipalId = &v
}

// GetTenantId returns the TenantId field value if set, zero value otherwise
func (o *CloudProviderAccessRoleRequestUpdate) GetTenantId() string {
	if o == nil || IsNil(o.TenantId) {
		var ret string
		return ret
	}
	return *o.TenantId
}

// GetTenantIdOk returns a tuple with the TenantId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessRoleRequestUpdate) GetTenantIdOk() (*string, bool) {
	if o == nil || IsNil(o.TenantId) {
		return nil, false
	}

	return o.TenantId, true
}

// HasTenantId returns a boolean if a field has been set.
func (o *CloudProviderAccessRoleRequestUpdate) HasTenantId() bool {
	if o != nil && !IsNil(o.TenantId) {
		return true
	}

	return false
}

// SetTenantId gets a reference to the given string and assigns it to the TenantId field.
func (o *CloudProviderAccessRoleRequestUpdate) SetTenantId(v string) {
	o.TenantId = &v
}
