// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AccountDetails Account details for the group, region, and provider.
type AccountDetails struct {
	// The AWS Account ID.
	AwsAccountId *string `json:"awsAccountId,omitempty"`
	// The VPC CIDR Block.
	CidrBlock *string `json:"cidrBlock,omitempty"`
	// Cloud provider.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// The VPC ID.
	VpcId *string `json:"vpcId,omitempty"`
	// The Azure Subscription ID.
	AzureSubscriptionId *string `json:"azureSubscriptionId,omitempty"`
	// The name of the virtual network.
	VirtualNetworkName *string `json:"virtualNetworkName,omitempty"`
	// The GCP Project ID.
	GcpProjectId *string `json:"gcpProjectId,omitempty"`
	// The name of the VPC network.
	VpcNetworkName *string `json:"vpcNetworkName,omitempty"`
}

// NewAccountDetails instantiates a new AccountDetails object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccountDetails() *AccountDetails {
	this := AccountDetails{}
	return &this
}

// NewAccountDetailsWithDefaults instantiates a new AccountDetails object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccountDetailsWithDefaults() *AccountDetails {
	this := AccountDetails{}
	return &this
}

// GetAwsAccountId returns the AwsAccountId field value if set, zero value otherwise
func (o *AccountDetails) GetAwsAccountId() string {
	if o == nil || IsNil(o.AwsAccountId) {
		var ret string
		return ret
	}
	return *o.AwsAccountId
}

// GetAwsAccountIdOk returns a tuple with the AwsAccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetAwsAccountIdOk() (*string, bool) {
	if o == nil || IsNil(o.AwsAccountId) {
		return nil, false
	}

	return o.AwsAccountId, true
}

// HasAwsAccountId returns a boolean if a field has been set.
func (o *AccountDetails) HasAwsAccountId() bool {
	if o != nil && !IsNil(o.AwsAccountId) {
		return true
	}

	return false
}

// SetAwsAccountId gets a reference to the given string and assigns it to the AwsAccountId field.
func (o *AccountDetails) SetAwsAccountId(v string) {
	o.AwsAccountId = &v
}

// GetCidrBlock returns the CidrBlock field value if set, zero value otherwise
func (o *AccountDetails) GetCidrBlock() string {
	if o == nil || IsNil(o.CidrBlock) {
		var ret string
		return ret
	}
	return *o.CidrBlock
}

// GetCidrBlockOk returns a tuple with the CidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.CidrBlock) {
		return nil, false
	}

	return o.CidrBlock, true
}

// HasCidrBlock returns a boolean if a field has been set.
func (o *AccountDetails) HasCidrBlock() bool {
	if o != nil && !IsNil(o.CidrBlock) {
		return true
	}

	return false
}

// SetCidrBlock gets a reference to the given string and assigns it to the CidrBlock field.
func (o *AccountDetails) SetCidrBlock(v string) {
	o.CidrBlock = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *AccountDetails) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *AccountDetails) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *AccountDetails) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *AccountDetails) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *AccountDetails) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *AccountDetails) SetLinks(v []Link) {
	o.Links = &v
}

// GetVpcId returns the VpcId field value if set, zero value otherwise
func (o *AccountDetails) GetVpcId() string {
	if o == nil || IsNil(o.VpcId) {
		var ret string
		return ret
	}
	return *o.VpcId
}

// GetVpcIdOk returns a tuple with the VpcId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetVpcIdOk() (*string, bool) {
	if o == nil || IsNil(o.VpcId) {
		return nil, false
	}

	return o.VpcId, true
}

// HasVpcId returns a boolean if a field has been set.
func (o *AccountDetails) HasVpcId() bool {
	if o != nil && !IsNil(o.VpcId) {
		return true
	}

	return false
}

// SetVpcId gets a reference to the given string and assigns it to the VpcId field.
func (o *AccountDetails) SetVpcId(v string) {
	o.VpcId = &v
}

// GetAzureSubscriptionId returns the AzureSubscriptionId field value if set, zero value otherwise
func (o *AccountDetails) GetAzureSubscriptionId() string {
	if o == nil || IsNil(o.AzureSubscriptionId) {
		var ret string
		return ret
	}
	return *o.AzureSubscriptionId
}

// GetAzureSubscriptionIdOk returns a tuple with the AzureSubscriptionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetAzureSubscriptionIdOk() (*string, bool) {
	if o == nil || IsNil(o.AzureSubscriptionId) {
		return nil, false
	}

	return o.AzureSubscriptionId, true
}

// HasAzureSubscriptionId returns a boolean if a field has been set.
func (o *AccountDetails) HasAzureSubscriptionId() bool {
	if o != nil && !IsNil(o.AzureSubscriptionId) {
		return true
	}

	return false
}

// SetAzureSubscriptionId gets a reference to the given string and assigns it to the AzureSubscriptionId field.
func (o *AccountDetails) SetAzureSubscriptionId(v string) {
	o.AzureSubscriptionId = &v
}

// GetVirtualNetworkName returns the VirtualNetworkName field value if set, zero value otherwise
func (o *AccountDetails) GetVirtualNetworkName() string {
	if o == nil || IsNil(o.VirtualNetworkName) {
		var ret string
		return ret
	}
	return *o.VirtualNetworkName
}

// GetVirtualNetworkNameOk returns a tuple with the VirtualNetworkName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetVirtualNetworkNameOk() (*string, bool) {
	if o == nil || IsNil(o.VirtualNetworkName) {
		return nil, false
	}

	return o.VirtualNetworkName, true
}

// HasVirtualNetworkName returns a boolean if a field has been set.
func (o *AccountDetails) HasVirtualNetworkName() bool {
	if o != nil && !IsNil(o.VirtualNetworkName) {
		return true
	}

	return false
}

// SetVirtualNetworkName gets a reference to the given string and assigns it to the VirtualNetworkName field.
func (o *AccountDetails) SetVirtualNetworkName(v string) {
	o.VirtualNetworkName = &v
}

// GetGcpProjectId returns the GcpProjectId field value if set, zero value otherwise
func (o *AccountDetails) GetGcpProjectId() string {
	if o == nil || IsNil(o.GcpProjectId) {
		var ret string
		return ret
	}
	return *o.GcpProjectId
}

// GetGcpProjectIdOk returns a tuple with the GcpProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetGcpProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.GcpProjectId) {
		return nil, false
	}

	return o.GcpProjectId, true
}

// HasGcpProjectId returns a boolean if a field has been set.
func (o *AccountDetails) HasGcpProjectId() bool {
	if o != nil && !IsNil(o.GcpProjectId) {
		return true
	}

	return false
}

// SetGcpProjectId gets a reference to the given string and assigns it to the GcpProjectId field.
func (o *AccountDetails) SetGcpProjectId(v string) {
	o.GcpProjectId = &v
}

// GetVpcNetworkName returns the VpcNetworkName field value if set, zero value otherwise
func (o *AccountDetails) GetVpcNetworkName() string {
	if o == nil || IsNil(o.VpcNetworkName) {
		var ret string
		return ret
	}
	return *o.VpcNetworkName
}

// GetVpcNetworkNameOk returns a tuple with the VpcNetworkName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AccountDetails) GetVpcNetworkNameOk() (*string, bool) {
	if o == nil || IsNil(o.VpcNetworkName) {
		return nil, false
	}

	return o.VpcNetworkName, true
}

// HasVpcNetworkName returns a boolean if a field has been set.
func (o *AccountDetails) HasVpcNetworkName() bool {
	if o != nil && !IsNil(o.VpcNetworkName) {
		return true
	}

	return false
}

// SetVpcNetworkName gets a reference to the given string and assigns it to the VpcNetworkName field.
func (o *AccountDetails) SetVpcNetworkName(v string) {
	o.VpcNetworkName = &v
}
