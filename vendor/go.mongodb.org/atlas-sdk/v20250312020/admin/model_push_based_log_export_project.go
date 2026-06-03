// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// PushBasedLogExportProject struct for PushBasedLogExportProject
type PushBasedLogExportProject struct {
	// The name of the bucket to which the agent will send the logs to.
	BucketName *string `json:"bucketName,omitempty"`
	// Date and time that this feature was enabled on. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreateDate *time.Time `json:"createDate,omitempty"`
	// ID of the AWS IAM role that will be used to write to the S3 bucket.
	IamRoleId *string `json:"iamRoleId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// S3 directory in which vector will write to in order to store the logs. An empty string denotes the root directory.
	PrefixPath *string `json:"prefixPath,omitempty"`
	// Describes whether or not the feature is enabled and what status it is in.
	// Read only field.
	State *string `json:"state,omitempty"`
}

// NewPushBasedLogExportProject instantiates a new PushBasedLogExportProject object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPushBasedLogExportProject() *PushBasedLogExportProject {
	this := PushBasedLogExportProject{}
	return &this
}

// NewPushBasedLogExportProjectWithDefaults instantiates a new PushBasedLogExportProject object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPushBasedLogExportProjectWithDefaults() *PushBasedLogExportProject {
	this := PushBasedLogExportProject{}
	return &this
}

// GetBucketName returns the BucketName field value if set, zero value otherwise
func (o *PushBasedLogExportProject) GetBucketName() string {
	if o == nil || IsNil(o.BucketName) {
		var ret string
		return ret
	}
	return *o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PushBasedLogExportProject) GetBucketNameOk() (*string, bool) {
	if o == nil || IsNil(o.BucketName) {
		return nil, false
	}

	return o.BucketName, true
}

// HasBucketName returns a boolean if a field has been set.
func (o *PushBasedLogExportProject) HasBucketName() bool {
	if o != nil && !IsNil(o.BucketName) {
		return true
	}

	return false
}

// SetBucketName gets a reference to the given string and assigns it to the BucketName field.
func (o *PushBasedLogExportProject) SetBucketName(v string) {
	o.BucketName = &v
}

// GetCreateDate returns the CreateDate field value if set, zero value otherwise
func (o *PushBasedLogExportProject) GetCreateDate() time.Time {
	if o == nil || IsNil(o.CreateDate) {
		var ret time.Time
		return ret
	}
	return *o.CreateDate
}

// GetCreateDateOk returns a tuple with the CreateDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PushBasedLogExportProject) GetCreateDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateDate) {
		return nil, false
	}

	return o.CreateDate, true
}

// HasCreateDate returns a boolean if a field has been set.
func (o *PushBasedLogExportProject) HasCreateDate() bool {
	if o != nil && !IsNil(o.CreateDate) {
		return true
	}

	return false
}

// SetCreateDate gets a reference to the given time.Time and assigns it to the CreateDate field.
func (o *PushBasedLogExportProject) SetCreateDate(v time.Time) {
	o.CreateDate = &v
}

// GetIamRoleId returns the IamRoleId field value if set, zero value otherwise
func (o *PushBasedLogExportProject) GetIamRoleId() string {
	if o == nil || IsNil(o.IamRoleId) {
		var ret string
		return ret
	}
	return *o.IamRoleId
}

// GetIamRoleIdOk returns a tuple with the IamRoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PushBasedLogExportProject) GetIamRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.IamRoleId) {
		return nil, false
	}

	return o.IamRoleId, true
}

// HasIamRoleId returns a boolean if a field has been set.
func (o *PushBasedLogExportProject) HasIamRoleId() bool {
	if o != nil && !IsNil(o.IamRoleId) {
		return true
	}

	return false
}

// SetIamRoleId gets a reference to the given string and assigns it to the IamRoleId field.
func (o *PushBasedLogExportProject) SetIamRoleId(v string) {
	o.IamRoleId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *PushBasedLogExportProject) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PushBasedLogExportProject) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PushBasedLogExportProject) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *PushBasedLogExportProject) SetLinks(v []Link) {
	o.Links = &v
}

// GetPrefixPath returns the PrefixPath field value if set, zero value otherwise
func (o *PushBasedLogExportProject) GetPrefixPath() string {
	if o == nil || IsNil(o.PrefixPath) {
		var ret string
		return ret
	}
	return *o.PrefixPath
}

// GetPrefixPathOk returns a tuple with the PrefixPath field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PushBasedLogExportProject) GetPrefixPathOk() (*string, bool) {
	if o == nil || IsNil(o.PrefixPath) {
		return nil, false
	}

	return o.PrefixPath, true
}

// HasPrefixPath returns a boolean if a field has been set.
func (o *PushBasedLogExportProject) HasPrefixPath() bool {
	if o != nil && !IsNil(o.PrefixPath) {
		return true
	}

	return false
}

// SetPrefixPath gets a reference to the given string and assigns it to the PrefixPath field.
func (o *PushBasedLogExportProject) SetPrefixPath(v string) {
	o.PrefixPath = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *PushBasedLogExportProject) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PushBasedLogExportProject) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *PushBasedLogExportProject) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *PushBasedLogExportProject) SetState(v string) {
	o.State = &v
}
