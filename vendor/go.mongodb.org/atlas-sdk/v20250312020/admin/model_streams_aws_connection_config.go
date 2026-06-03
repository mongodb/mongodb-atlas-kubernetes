// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsAWSConnectionConfig AWS configurations for AWS-based connection types.
type StreamsAWSConnectionConfig struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Amazon Resource Name (ARN) that identifies the Amazon Web Services (AWS) Identity and Access Management (IAM) role that MongoDB Cloud assumes when it accesses resources in your AWS account.
	RoleArn *string `json:"roleArn,omitempty"`
	// The name of an S3 bucket used to check authorization of the passed-in IAM role ARN.
	TestBucket *string `json:"testBucket,omitempty"`
}

// NewStreamsAWSConnectionConfig instantiates a new StreamsAWSConnectionConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsAWSConnectionConfig() *StreamsAWSConnectionConfig {
	this := StreamsAWSConnectionConfig{}
	return &this
}

// NewStreamsAWSConnectionConfigWithDefaults instantiates a new StreamsAWSConnectionConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsAWSConnectionConfigWithDefaults() *StreamsAWSConnectionConfig {
	this := StreamsAWSConnectionConfig{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsAWSConnectionConfig) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsAWSConnectionConfig) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsAWSConnectionConfig) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsAWSConnectionConfig) SetLinks(v []Link) {
	o.Links = &v
}

// GetRoleArn returns the RoleArn field value if set, zero value otherwise
func (o *StreamsAWSConnectionConfig) GetRoleArn() string {
	if o == nil || IsNil(o.RoleArn) {
		var ret string
		return ret
	}
	return *o.RoleArn
}

// GetRoleArnOk returns a tuple with the RoleArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsAWSConnectionConfig) GetRoleArnOk() (*string, bool) {
	if o == nil || IsNil(o.RoleArn) {
		return nil, false
	}

	return o.RoleArn, true
}

// HasRoleArn returns a boolean if a field has been set.
func (o *StreamsAWSConnectionConfig) HasRoleArn() bool {
	if o != nil && !IsNil(o.RoleArn) {
		return true
	}

	return false
}

// SetRoleArn gets a reference to the given string and assigns it to the RoleArn field.
func (o *StreamsAWSConnectionConfig) SetRoleArn(v string) {
	o.RoleArn = &v
}

// GetTestBucket returns the TestBucket field value if set, zero value otherwise
func (o *StreamsAWSConnectionConfig) GetTestBucket() string {
	if o == nil || IsNil(o.TestBucket) {
		var ret string
		return ret
	}
	return *o.TestBucket
}

// GetTestBucketOk returns a tuple with the TestBucket field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsAWSConnectionConfig) GetTestBucketOk() (*string, bool) {
	if o == nil || IsNil(o.TestBucket) {
		return nil, false
	}

	return o.TestBucket, true
}

// HasTestBucket returns a boolean if a field has been set.
func (o *StreamsAWSConnectionConfig) HasTestBucket() bool {
	if o != nil && !IsNil(o.TestBucket) {
		return true
	}

	return false
}

// SetTestBucket gets a reference to the given string and assigns it to the TestBucket field.
func (o *StreamsAWSConnectionConfig) SetTestBucket(v string) {
	o.TestBucket = &v
}
