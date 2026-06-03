// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeAWSCloudProviderConfig Configuration for running Data Federation in AWS.
type DataLakeAWSCloudProviderConfig struct {
	// Unique identifier associated with the Identity and Access Management (IAM) role that the data lake assumes when accessing the data stores.
	// Read only field.
	ExternalId *string `json:"externalId,omitempty"`
	// Amazon Resource Name (ARN) of the Identity and Access Management (IAM) role that the data lake assumes when accessing data stores.
	// Read only field.
	IamAssumedRoleARN *string `json:"iamAssumedRoleARN,omitempty"`
	// Amazon Resource Name (ARN) of the user that the data lake assumes when accessing data stores.
	// Read only field.
	IamUserARN *string `json:"iamUserARN,omitempty"`
	// Unique identifier of the role that the data lake can use to access the data stores.Required if specifying cloudProviderConfig.
	RoleId string `json:"roleId"`
	// Name of the S3 data bucket that the provided role ID is authorized to access. Required if specifying `cloudProviderConfig`.
	// Write only field.
	TestS3Bucket string `json:"testS3Bucket"`
}

// NewDataLakeAWSCloudProviderConfig instantiates a new DataLakeAWSCloudProviderConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeAWSCloudProviderConfig(roleId string, testS3Bucket string) *DataLakeAWSCloudProviderConfig {
	this := DataLakeAWSCloudProviderConfig{}
	this.RoleId = roleId
	this.TestS3Bucket = testS3Bucket
	return &this
}

// NewDataLakeAWSCloudProviderConfigWithDefaults instantiates a new DataLakeAWSCloudProviderConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeAWSCloudProviderConfigWithDefaults() *DataLakeAWSCloudProviderConfig {
	this := DataLakeAWSCloudProviderConfig{}
	return &this
}

// GetExternalId returns the ExternalId field value if set, zero value otherwise
func (o *DataLakeAWSCloudProviderConfig) GetExternalId() string {
	if o == nil || IsNil(o.ExternalId) {
		var ret string
		return ret
	}
	return *o.ExternalId
}

// GetExternalIdOk returns a tuple with the ExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeAWSCloudProviderConfig) GetExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalId) {
		return nil, false
	}

	return o.ExternalId, true
}

// HasExternalId returns a boolean if a field has been set.
func (o *DataLakeAWSCloudProviderConfig) HasExternalId() bool {
	if o != nil && !IsNil(o.ExternalId) {
		return true
	}

	return false
}

// SetExternalId gets a reference to the given string and assigns it to the ExternalId field.
func (o *DataLakeAWSCloudProviderConfig) SetExternalId(v string) {
	o.ExternalId = &v
}

// GetIamAssumedRoleARN returns the IamAssumedRoleARN field value if set, zero value otherwise
func (o *DataLakeAWSCloudProviderConfig) GetIamAssumedRoleARN() string {
	if o == nil || IsNil(o.IamAssumedRoleARN) {
		var ret string
		return ret
	}
	return *o.IamAssumedRoleARN
}

// GetIamAssumedRoleARNOk returns a tuple with the IamAssumedRoleARN field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeAWSCloudProviderConfig) GetIamAssumedRoleARNOk() (*string, bool) {
	if o == nil || IsNil(o.IamAssumedRoleARN) {
		return nil, false
	}

	return o.IamAssumedRoleARN, true
}

// HasIamAssumedRoleARN returns a boolean if a field has been set.
func (o *DataLakeAWSCloudProviderConfig) HasIamAssumedRoleARN() bool {
	if o != nil && !IsNil(o.IamAssumedRoleARN) {
		return true
	}

	return false
}

// SetIamAssumedRoleARN gets a reference to the given string and assigns it to the IamAssumedRoleARN field.
func (o *DataLakeAWSCloudProviderConfig) SetIamAssumedRoleARN(v string) {
	o.IamAssumedRoleARN = &v
}

// GetIamUserARN returns the IamUserARN field value if set, zero value otherwise
func (o *DataLakeAWSCloudProviderConfig) GetIamUserARN() string {
	if o == nil || IsNil(o.IamUserARN) {
		var ret string
		return ret
	}
	return *o.IamUserARN
}

// GetIamUserARNOk returns a tuple with the IamUserARN field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeAWSCloudProviderConfig) GetIamUserARNOk() (*string, bool) {
	if o == nil || IsNil(o.IamUserARN) {
		return nil, false
	}

	return o.IamUserARN, true
}

// HasIamUserARN returns a boolean if a field has been set.
func (o *DataLakeAWSCloudProviderConfig) HasIamUserARN() bool {
	if o != nil && !IsNil(o.IamUserARN) {
		return true
	}

	return false
}

// SetIamUserARN gets a reference to the given string and assigns it to the IamUserARN field.
func (o *DataLakeAWSCloudProviderConfig) SetIamUserARN(v string) {
	o.IamUserARN = &v
}

// GetRoleId returns the RoleId field value
func (o *DataLakeAWSCloudProviderConfig) GetRoleId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value
// and a boolean to check if the value has been set.
func (o *DataLakeAWSCloudProviderConfig) GetRoleIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RoleId, true
}

// SetRoleId sets field value
func (o *DataLakeAWSCloudProviderConfig) SetRoleId(v string) {
	o.RoleId = v
}

// GetTestS3Bucket returns the TestS3Bucket field value
func (o *DataLakeAWSCloudProviderConfig) GetTestS3Bucket() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TestS3Bucket
}

// GetTestS3BucketOk returns a tuple with the TestS3Bucket field value
// and a boolean to check if the value has been set.
func (o *DataLakeAWSCloudProviderConfig) GetTestS3BucketOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TestS3Bucket, true
}

// SetTestS3Bucket sets field value
func (o *DataLakeAWSCloudProviderConfig) SetTestS3Bucket(v string) {
	o.TestS3Bucket = v
}
