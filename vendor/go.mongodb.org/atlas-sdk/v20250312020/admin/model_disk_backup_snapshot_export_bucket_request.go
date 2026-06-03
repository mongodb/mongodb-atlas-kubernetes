// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupSnapshotExportBucketRequest Disk backup snapshot Export Bucket Request.
type DiskBackupSnapshotExportBucketRequest struct {
	// Human-readable label that identifies the cloud provider.
	CloudProvider string `json:"cloudProvider"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the Google Cloud Storage Bucket that the role is authorized to export to.
	BucketName *string `json:"bucketName,omitempty"`
	// Unique 24-hexadecimal character string that identifies the Unified AWS Access role ID that MongoDB Cloud uses to access the AWS S3 bucket.
	IamRoleId *string `json:"iamRoleId,omitempty"`
	// Indicates whether to do exports over PrivateLink as opposed to public IPs. Defaults to False.
	RequirePrivateNetworking *bool `json:"requirePrivateNetworking,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the GCP Cloud Provider Access Role that MongoDB Cloud uses to access the Google Cloud Storage Bucket.
	RoleId *string `json:"roleId,omitempty"`
	// URL of the Azure Storage Account to export to. For example: `https://examplestorageaccount.blob.core.windows.net/exportcontainer`. Only standard endpoints (with `blob.core.windows.net`) are supported.
	ServiceUrl *string `json:"serviceUrl,omitempty"`
	// UUID that identifies the Azure Active Directory Tenant ID. Deprecated: this field is ignored; the `tenantId` of the Cloud Provider Access role (from `roleId`) is used.
	// Deprecated
	TenantId *string `json:"tenantId,omitempty"`
}

// NewDiskBackupSnapshotExportBucketRequest instantiates a new DiskBackupSnapshotExportBucketRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupSnapshotExportBucketRequest(cloudProvider string) *DiskBackupSnapshotExportBucketRequest {
	this := DiskBackupSnapshotExportBucketRequest{}
	this.CloudProvider = cloudProvider
	return &this
}

// NewDiskBackupSnapshotExportBucketRequestWithDefaults instantiates a new DiskBackupSnapshotExportBucketRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupSnapshotExportBucketRequestWithDefaults() *DiskBackupSnapshotExportBucketRequest {
	this := DiskBackupSnapshotExportBucketRequest{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value
func (o *DiskBackupSnapshotExportBucketRequest) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *DiskBackupSnapshotExportBucketRequest) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketRequest) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupSnapshotExportBucketRequest) SetLinks(v []Link) {
	o.Links = &v
}

// GetBucketName returns the BucketName field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketRequest) GetBucketName() string {
	if o == nil || IsNil(o.BucketName) {
		var ret string
		return ret
	}
	return *o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetBucketNameOk() (*string, bool) {
	if o == nil || IsNil(o.BucketName) {
		return nil, false
	}

	return o.BucketName, true
}

// HasBucketName returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasBucketName() bool {
	if o != nil && !IsNil(o.BucketName) {
		return true
	}

	return false
}

// SetBucketName gets a reference to the given string and assigns it to the BucketName field.
func (o *DiskBackupSnapshotExportBucketRequest) SetBucketName(v string) {
	o.BucketName = &v
}

// GetIamRoleId returns the IamRoleId field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketRequest) GetIamRoleId() string {
	if o == nil || IsNil(o.IamRoleId) {
		var ret string
		return ret
	}
	return *o.IamRoleId
}

// GetIamRoleIdOk returns a tuple with the IamRoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetIamRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.IamRoleId) {
		return nil, false
	}

	return o.IamRoleId, true
}

// HasIamRoleId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasIamRoleId() bool {
	if o != nil && !IsNil(o.IamRoleId) {
		return true
	}

	return false
}

// SetIamRoleId gets a reference to the given string and assigns it to the IamRoleId field.
func (o *DiskBackupSnapshotExportBucketRequest) SetIamRoleId(v string) {
	o.IamRoleId = &v
}

// GetRequirePrivateNetworking returns the RequirePrivateNetworking field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketRequest) GetRequirePrivateNetworking() bool {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		var ret bool
		return ret
	}
	return *o.RequirePrivateNetworking
}

// GetRequirePrivateNetworkingOk returns a tuple with the RequirePrivateNetworking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetRequirePrivateNetworkingOk() (*bool, bool) {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		return nil, false
	}

	return o.RequirePrivateNetworking, true
}

// HasRequirePrivateNetworking returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasRequirePrivateNetworking() bool {
	if o != nil && !IsNil(o.RequirePrivateNetworking) {
		return true
	}

	return false
}

// SetRequirePrivateNetworking gets a reference to the given bool and assigns it to the RequirePrivateNetworking field.
func (o *DiskBackupSnapshotExportBucketRequest) SetRequirePrivateNetworking(v bool) {
	o.RequirePrivateNetworking = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketRequest) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *DiskBackupSnapshotExportBucketRequest) SetRoleId(v string) {
	o.RoleId = &v
}

// GetServiceUrl returns the ServiceUrl field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketRequest) GetServiceUrl() string {
	if o == nil || IsNil(o.ServiceUrl) {
		var ret string
		return ret
	}
	return *o.ServiceUrl
}

// GetServiceUrlOk returns a tuple with the ServiceUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketRequest) GetServiceUrlOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceUrl) {
		return nil, false
	}

	return o.ServiceUrl, true
}

// HasServiceUrl returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasServiceUrl() bool {
	if o != nil && !IsNil(o.ServiceUrl) {
		return true
	}

	return false
}

// SetServiceUrl gets a reference to the given string and assigns it to the ServiceUrl field.
func (o *DiskBackupSnapshotExportBucketRequest) SetServiceUrl(v string) {
	o.ServiceUrl = &v
}

// GetTenantId returns the TenantId field value if set, zero value otherwise
// Deprecated
func (o *DiskBackupSnapshotExportBucketRequest) GetTenantId() string {
	if o == nil || IsNil(o.TenantId) {
		var ret string
		return ret
	}
	return *o.TenantId
}

// GetTenantIdOk returns a tuple with the TenantId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *DiskBackupSnapshotExportBucketRequest) GetTenantIdOk() (*string, bool) {
	if o == nil || IsNil(o.TenantId) {
		return nil, false
	}

	return o.TenantId, true
}

// HasTenantId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketRequest) HasTenantId() bool {
	if o != nil && !IsNil(o.TenantId) {
		return true
	}

	return false
}

// SetTenantId gets a reference to the given string and assigns it to the TenantId field.
// Deprecated
func (o *DiskBackupSnapshotExportBucketRequest) SetTenantId(v string) {
	o.TenantId = &v
}
