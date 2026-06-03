// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId Identifying characteristics about the Amazon Web Services (AWS) Simple Storage Service (S3) export bucket linked to this AWS Identity and Access Management (IAM) role.
type CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId struct {
	// Name of the AWS S3 bucket to which your logs will be exported to.
	// Read only field.
	BucketName *string `json:"bucketName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies your project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
}

// NewCloudProviderAccessFeatureUsagePushBasedLogExportFeatureId instantiates a new CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderAccessFeatureUsagePushBasedLogExportFeatureId() *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId {
	this := CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId{}
	return &this
}

// NewCloudProviderAccessFeatureUsagePushBasedLogExportFeatureIdWithDefaults instantiates a new CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderAccessFeatureUsagePushBasedLogExportFeatureIdWithDefaults() *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId {
	this := CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId{}
	return &this
}

// GetBucketName returns the BucketName field value if set, zero value otherwise
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) GetBucketName() string {
	if o == nil || IsNil(o.BucketName) {
		var ret string
		return ret
	}
	return *o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) GetBucketNameOk() (*string, bool) {
	if o == nil || IsNil(o.BucketName) {
		return nil, false
	}

	return o.BucketName, true
}

// HasBucketName returns a boolean if a field has been set.
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) HasBucketName() bool {
	if o != nil && !IsNil(o.BucketName) {
		return true
	}

	return false
}

// SetBucketName gets a reference to the given string and assigns it to the BucketName field.
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) SetBucketName(v string) {
	o.BucketName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) SetGroupId(v string) {
	o.GroupId = &v
}
