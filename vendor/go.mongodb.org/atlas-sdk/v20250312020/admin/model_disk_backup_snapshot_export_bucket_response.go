// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupSnapshotExportBucketResponse Disk backup snapshot Export Bucket.
type DiskBackupSnapshotExportBucketResponse struct {
	// Unique 24-hexadecimal character string that identifies the Export Bucket.
	Id string `json:"_id"`
	// The name of the AWS S3 Bucket, Azure Storage Container, or Google Cloud Storage Bucket that Snapshots are exported to.
	BucketName string `json:"bucketName"`
	// Human-readable label that identifies the cloud provider.
	CloudProvider string `json:"cloudProvider"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal character string that identifies the Unified AWS Access role ID that MongoDB Cloud uses to access the AWS S3 bucket.
	IamRoleId *string `json:"iamRoleId,omitempty"`
	// AWS region for the export bucket. This is set by Atlas and is never user-supplied.
	// Read only field.
	Region *string `json:"region,omitempty"`
	// Indicates whether to use private link. User supplied.
	RequirePrivateNetworking *bool `json:"requirePrivateNetworking,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the GCP Cloud Provider Access Role that MongoDB Cloud uses to access the Google Cloud Storage Bucket.
	RoleId *string `json:"roleId,omitempty"`
	// URL of the Azure Storage Account to export to. Only standard endpoints (with `blob.core.windows.net`) are supported.
	ServiceUrl *string `json:"serviceUrl,omitempty"`
	// UUID that identifies the Azure Active Directory Tenant ID used during exports.
	TenantId *string `json:"tenantId,omitempty"`
}

// NewDiskBackupSnapshotExportBucketResponse instantiates a new DiskBackupSnapshotExportBucketResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupSnapshotExportBucketResponse(id string, bucketName string, cloudProvider string) *DiskBackupSnapshotExportBucketResponse {
	this := DiskBackupSnapshotExportBucketResponse{}
	this.Id = id
	this.BucketName = bucketName
	this.CloudProvider = cloudProvider
	return &this
}

// NewDiskBackupSnapshotExportBucketResponseWithDefaults instantiates a new DiskBackupSnapshotExportBucketResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupSnapshotExportBucketResponseWithDefaults() *DiskBackupSnapshotExportBucketResponse {
	this := DiskBackupSnapshotExportBucketResponse{}
	return &this
}

// GetId returns the Id field value
func (o *DiskBackupSnapshotExportBucketResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *DiskBackupSnapshotExportBucketResponse) SetId(v string) {
	o.Id = v
}

// GetBucketName returns the BucketName field value
func (o *DiskBackupSnapshotExportBucketResponse) GetBucketName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetBucketNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BucketName, true
}

// SetBucketName sets field value
func (o *DiskBackupSnapshotExportBucketResponse) SetBucketName(v string) {
	o.BucketName = v
}

// GetCloudProvider returns the CloudProvider field value
func (o *DiskBackupSnapshotExportBucketResponse) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *DiskBackupSnapshotExportBucketResponse) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupSnapshotExportBucketResponse) SetLinks(v []Link) {
	o.Links = &v
}

// GetIamRoleId returns the IamRoleId field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetIamRoleId() string {
	if o == nil || IsNil(o.IamRoleId) {
		var ret string
		return ret
	}
	return *o.IamRoleId
}

// GetIamRoleIdOk returns a tuple with the IamRoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetIamRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.IamRoleId) {
		return nil, false
	}

	return o.IamRoleId, true
}

// HasIamRoleId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasIamRoleId() bool {
	if o != nil && !IsNil(o.IamRoleId) {
		return true
	}

	return false
}

// SetIamRoleId gets a reference to the given string and assigns it to the IamRoleId field.
func (o *DiskBackupSnapshotExportBucketResponse) SetIamRoleId(v string) {
	o.IamRoleId = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *DiskBackupSnapshotExportBucketResponse) SetRegion(v string) {
	o.Region = &v
}

// GetRequirePrivateNetworking returns the RequirePrivateNetworking field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetRequirePrivateNetworking() bool {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		var ret bool
		return ret
	}
	return *o.RequirePrivateNetworking
}

// GetRequirePrivateNetworkingOk returns a tuple with the RequirePrivateNetworking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetRequirePrivateNetworkingOk() (*bool, bool) {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		return nil, false
	}

	return o.RequirePrivateNetworking, true
}

// HasRequirePrivateNetworking returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasRequirePrivateNetworking() bool {
	if o != nil && !IsNil(o.RequirePrivateNetworking) {
		return true
	}

	return false
}

// SetRequirePrivateNetworking gets a reference to the given bool and assigns it to the RequirePrivateNetworking field.
func (o *DiskBackupSnapshotExportBucketResponse) SetRequirePrivateNetworking(v bool) {
	o.RequirePrivateNetworking = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *DiskBackupSnapshotExportBucketResponse) SetRoleId(v string) {
	o.RoleId = &v
}

// GetServiceUrl returns the ServiceUrl field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetServiceUrl() string {
	if o == nil || IsNil(o.ServiceUrl) {
		var ret string
		return ret
	}
	return *o.ServiceUrl
}

// GetServiceUrlOk returns a tuple with the ServiceUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetServiceUrlOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceUrl) {
		return nil, false
	}

	return o.ServiceUrl, true
}

// HasServiceUrl returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasServiceUrl() bool {
	if o != nil && !IsNil(o.ServiceUrl) {
		return true
	}

	return false
}

// SetServiceUrl gets a reference to the given string and assigns it to the ServiceUrl field.
func (o *DiskBackupSnapshotExportBucketResponse) SetServiceUrl(v string) {
	o.ServiceUrl = &v
}

// GetTenantId returns the TenantId field value if set, zero value otherwise
func (o *DiskBackupSnapshotExportBucketResponse) GetTenantId() string {
	if o == nil || IsNil(o.TenantId) {
		var ret string
		return ret
	}
	return *o.TenantId
}

// GetTenantIdOk returns a tuple with the TenantId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotExportBucketResponse) GetTenantIdOk() (*string, bool) {
	if o == nil || IsNil(o.TenantId) {
		return nil, false
	}

	return o.TenantId, true
}

// HasTenantId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotExportBucketResponse) HasTenantId() bool {
	if o != nil && !IsNil(o.TenantId) {
		return true
	}

	return false
}

// SetTenantId gets a reference to the given string and assigns it to the TenantId field.
func (o *DiskBackupSnapshotExportBucketResponse) SetTenantId(v string) {
	o.TenantId = &v
}
