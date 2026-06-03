// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupSnapshotAWSExportBucketResponse struct for DiskBackupSnapshotAWSExportBucketResponse
type DiskBackupSnapshotAWSExportBucketResponse struct {
	// Unique 24-hexadecimal character string that identifies the Export Bucket.
	Id string `json:"_id"`
	// The name of the AWS S3 Bucket, Azure Storage Container, or Google Cloud Storage Bucket that Snapshots are exported to.
	BucketName string `json:"bucketName"`
	// Human-readable label that identifies the cloud provider.
	CloudProvider string `json:"cloudProvider"`
	// Unique 24-hexadecimal character string that identifies the Unified AWS Access role ID that MongoDB Cloud uses to access the AWS S3 bucket.
	IamRoleId string `json:"iamRoleId"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// AWS region for the export bucket. This is set by Atlas and is never user-supplied.
	// Read only field.
	Region *string `json:"region,omitempty"`
	// Indicates whether to use private link. User supplied.
	RequirePrivateNetworking *bool `json:"requirePrivateNetworking,omitempty"`
}

// NewDiskBackupSnapshotAWSExportBucketResponse instantiates a new DiskBackupSnapshotAWSExportBucketResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupSnapshotAWSExportBucketResponse(id string, bucketName string, cloudProvider string, iamRoleId string) *DiskBackupSnapshotAWSExportBucketResponse {
	this := DiskBackupSnapshotAWSExportBucketResponse{}
	this.Id = id
	this.BucketName = bucketName
	this.CloudProvider = cloudProvider
	this.IamRoleId = iamRoleId
	return &this
}

// NewDiskBackupSnapshotAWSExportBucketResponseWithDefaults instantiates a new DiskBackupSnapshotAWSExportBucketResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupSnapshotAWSExportBucketResponseWithDefaults() *DiskBackupSnapshotAWSExportBucketResponse {
	this := DiskBackupSnapshotAWSExportBucketResponse{}
	return &this
}

// GetId returns the Id field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetId(v string) {
	o.Id = v
}

// GetBucketName returns the BucketName field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetBucketName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetBucketNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BucketName, true
}

// SetBucketName sets field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetBucketName(v string) {
	o.BucketName = v
}

// GetCloudProvider returns the CloudProvider field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetIamRoleId returns the IamRoleId field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetIamRoleId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IamRoleId
}

// GetIamRoleIdOk returns a tuple with the IamRoleId field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetIamRoleIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IamRoleId, true
}

// SetIamRoleId sets field value
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetIamRoleId(v string) {
	o.IamRoleId = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetLinks(v []Link) {
	o.Links = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetRegion(v string) {
	o.Region = &v
}

// GetRequirePrivateNetworking returns the RequirePrivateNetworking field value if set, zero value otherwise
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetRequirePrivateNetworking() bool {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		var ret bool
		return ret
	}
	return *o.RequirePrivateNetworking
}

// GetRequirePrivateNetworkingOk returns a tuple with the RequirePrivateNetworking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) GetRequirePrivateNetworkingOk() (*bool, bool) {
	if o == nil || IsNil(o.RequirePrivateNetworking) {
		return nil, false
	}

	return o.RequirePrivateNetworking, true
}

// HasRequirePrivateNetworking returns a boolean if a field has been set.
func (o *DiskBackupSnapshotAWSExportBucketResponse) HasRequirePrivateNetworking() bool {
	if o != nil && !IsNil(o.RequirePrivateNetworking) {
		return true
	}

	return false
}

// SetRequirePrivateNetworking gets a reference to the given bool and assigns it to the RequirePrivateNetworking field.
func (o *DiskBackupSnapshotAWSExportBucketResponse) SetRequirePrivateNetworking(v bool) {
	o.RequirePrivateNetworking = &v
}
