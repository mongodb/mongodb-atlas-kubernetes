// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// NetworkPermissionEntry struct for NetworkPermissionEntry
type NetworkPermissionEntry struct {
	// Unique string of the Amazon Web Services (AWS) security group that you want to add to the project's IP access list. Your IP access list entry can be one `awsSecurityGroup`, one `cidrBlock`, or one `ipAddress`. You must configure Virtual Private Connection (VPC) peering for your project before you can add an AWS security group to an IP access list. You cannot set AWS security groups as temporary access list entries. Don't set this parameter if you set `cidrBlock` or `ipAddress`.
	AwsSecurityGroup *string `json:"awsSecurityGroup,omitempty"`
	// Range of IP addresses in Classless Inter-Domain Routing (CIDR) notation that you want to add to the project's IP access list. Your IP access list entry can be one `awsSecurityGroup`, one `cidrBlock`, or one `ipAddress`. Don't set this parameter if you set `awsSecurityGroup` or `ipAddress`.
	CidrBlock *string `json:"cidrBlock,omitempty"`
	// Remark that explains the purpose or scope of this IP access list entry.
	Comment *string `json:"comment,omitempty"`
	// Date and time after which MongoDB Cloud deletes the temporary access list entry. This parameter expresses its value in the ISO 8601 timestamp format in UTC and can include the time zone designation. The date must be later than the current date but no later than one week after you submit this request. The resource returns this parameter if you specified an expiration date when creating this IP access list entry.
	DeleteAfterDate *time.Time `json:"deleteAfterDate,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the IP access list to which you want to add one or more entries.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// IP address that you want to add to the project's IP access list. Your IP access list entry can be one `awsSecurityGroup`, one `cidrBlock`, or one `ipAddress`. Don't set this parameter if you set `awsSecurityGroup` or `cidrBlock`.
	IpAddress *string `json:"ipAddress,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewNetworkPermissionEntry instantiates a new NetworkPermissionEntry object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNetworkPermissionEntry() *NetworkPermissionEntry {
	this := NetworkPermissionEntry{}
	return &this
}

// NewNetworkPermissionEntryWithDefaults instantiates a new NetworkPermissionEntry object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNetworkPermissionEntryWithDefaults() *NetworkPermissionEntry {
	this := NetworkPermissionEntry{}
	return &this
}

// GetAwsSecurityGroup returns the AwsSecurityGroup field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetAwsSecurityGroup() string {
	if o == nil || IsNil(o.AwsSecurityGroup) {
		var ret string
		return ret
	}
	return *o.AwsSecurityGroup
}

// GetAwsSecurityGroupOk returns a tuple with the AwsSecurityGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetAwsSecurityGroupOk() (*string, bool) {
	if o == nil || IsNil(o.AwsSecurityGroup) {
		return nil, false
	}

	return o.AwsSecurityGroup, true
}

// HasAwsSecurityGroup returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasAwsSecurityGroup() bool {
	if o != nil && !IsNil(o.AwsSecurityGroup) {
		return true
	}

	return false
}

// SetAwsSecurityGroup gets a reference to the given string and assigns it to the AwsSecurityGroup field.
func (o *NetworkPermissionEntry) SetAwsSecurityGroup(v string) {
	o.AwsSecurityGroup = &v
}

// GetCidrBlock returns the CidrBlock field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetCidrBlock() string {
	if o == nil || IsNil(o.CidrBlock) {
		var ret string
		return ret
	}
	return *o.CidrBlock
}

// GetCidrBlockOk returns a tuple with the CidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.CidrBlock) {
		return nil, false
	}

	return o.CidrBlock, true
}

// HasCidrBlock returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasCidrBlock() bool {
	if o != nil && !IsNil(o.CidrBlock) {
		return true
	}

	return false
}

// SetCidrBlock gets a reference to the given string and assigns it to the CidrBlock field.
func (o *NetworkPermissionEntry) SetCidrBlock(v string) {
	o.CidrBlock = &v
}

// GetComment returns the Comment field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetComment() string {
	if o == nil || IsNil(o.Comment) {
		var ret string
		return ret
	}
	return *o.Comment
}

// GetCommentOk returns a tuple with the Comment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetCommentOk() (*string, bool) {
	if o == nil || IsNil(o.Comment) {
		return nil, false
	}

	return o.Comment, true
}

// HasComment returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasComment() bool {
	if o != nil && !IsNil(o.Comment) {
		return true
	}

	return false
}

// SetComment gets a reference to the given string and assigns it to the Comment field.
func (o *NetworkPermissionEntry) SetComment(v string) {
	o.Comment = &v
}

// GetDeleteAfterDate returns the DeleteAfterDate field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetDeleteAfterDate() time.Time {
	if o == nil || IsNil(o.DeleteAfterDate) {
		var ret time.Time
		return ret
	}
	return *o.DeleteAfterDate
}

// GetDeleteAfterDateOk returns a tuple with the DeleteAfterDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetDeleteAfterDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DeleteAfterDate) {
		return nil, false
	}

	return o.DeleteAfterDate, true
}

// HasDeleteAfterDate returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasDeleteAfterDate() bool {
	if o != nil && !IsNil(o.DeleteAfterDate) {
		return true
	}

	return false
}

// SetDeleteAfterDate gets a reference to the given time.Time and assigns it to the DeleteAfterDate field.
func (o *NetworkPermissionEntry) SetDeleteAfterDate(v time.Time) {
	o.DeleteAfterDate = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *NetworkPermissionEntry) SetGroupId(v string) {
	o.GroupId = &v
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetIpAddress() string {
	if o == nil || IsNil(o.IpAddress) {
		var ret string
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}

	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given string and assigns it to the IpAddress field.
func (o *NetworkPermissionEntry) SetIpAddress(v string) {
	o.IpAddress = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *NetworkPermissionEntry) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntry) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *NetworkPermissionEntry) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *NetworkPermissionEntry) SetLinks(v []Link) {
	o.Links = &v
}
