// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CreatePushBasedLogExportProjectRequest struct for CreatePushBasedLogExportProjectRequest
type CreatePushBasedLogExportProjectRequest struct {
	// The name of the bucket to which the agent will send the logs to.
	BucketName string `json:"bucketName"`
	// ID of the AWS IAM role that will be used to write to the S3 bucket.
	IamRoleId string `json:"iamRoleId"`
	// S3 directory in which vector will write to in order to store the logs. An empty string denotes the root directory.
	PrefixPath string `json:"prefixPath"`
}

// NewCreatePushBasedLogExportProjectRequest instantiates a new CreatePushBasedLogExportProjectRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreatePushBasedLogExportProjectRequest(bucketName string, iamRoleId string, prefixPath string) *CreatePushBasedLogExportProjectRequest {
	this := CreatePushBasedLogExportProjectRequest{}
	this.BucketName = bucketName
	this.IamRoleId = iamRoleId
	this.PrefixPath = prefixPath
	return &this
}

// NewCreatePushBasedLogExportProjectRequestWithDefaults instantiates a new CreatePushBasedLogExportProjectRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreatePushBasedLogExportProjectRequestWithDefaults() *CreatePushBasedLogExportProjectRequest {
	this := CreatePushBasedLogExportProjectRequest{}
	return &this
}

// GetBucketName returns the BucketName field value
func (o *CreatePushBasedLogExportProjectRequest) GetBucketName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value
// and a boolean to check if the value has been set.
func (o *CreatePushBasedLogExportProjectRequest) GetBucketNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BucketName, true
}

// SetBucketName sets field value
func (o *CreatePushBasedLogExportProjectRequest) SetBucketName(v string) {
	o.BucketName = v
}

// GetIamRoleId returns the IamRoleId field value
func (o *CreatePushBasedLogExportProjectRequest) GetIamRoleId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IamRoleId
}

// GetIamRoleIdOk returns a tuple with the IamRoleId field value
// and a boolean to check if the value has been set.
func (o *CreatePushBasedLogExportProjectRequest) GetIamRoleIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IamRoleId, true
}

// SetIamRoleId sets field value
func (o *CreatePushBasedLogExportProjectRequest) SetIamRoleId(v string) {
	o.IamRoleId = v
}

// GetPrefixPath returns the PrefixPath field value
func (o *CreatePushBasedLogExportProjectRequest) GetPrefixPath() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PrefixPath
}

// GetPrefixPathOk returns a tuple with the PrefixPath field value
// and a boolean to check if the value has been set.
func (o *CreatePushBasedLogExportProjectRequest) GetPrefixPathOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PrefixPath, true
}

// SetPrefixPath sets field value
func (o *CreatePushBasedLogExportProjectRequest) SetPrefixPath(v string) {
	o.PrefixPath = v
}
