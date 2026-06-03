// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudProviderContainer Collection of settings that configures the network container for a virtual private connection on Amazon Web Services.
type CloudProviderContainer struct {
	// Unique 24-hexadecimal digit string that identifies the network peering container.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Cloud service provider that serves the requested network peering containers.
	ProviderName *string `json:"providerName,omitempty"`
	// Flag that indicates whether MongoDB Cloud clusters exist in the specified network peering container.
	// Read only field.
	Provisioned *bool `json:"provisioned,omitempty"`
	// IP addresses expressed in Classless Inter-Domain Routing (CIDR) notation that MongoDB Cloud uses for the network peering containers in your project. MongoDB Cloud assigns all of the project's clusters deployed to this cloud provider an IP address from this range. MongoDB Cloud locks this value if an M10 or greater cluster or a network peering connection exists in this project.  These CIDR blocks must fall within the ranges reserved per RFC 1918. AWS and Azure further limit the block to between the `/24` and  `/21` ranges.  To modify the CIDR block, the target project cannot have:  - Any M10 or greater clusters - Any other VPC peering connections   You can also create a new project and create a network peering connection to set the desired MongoDB Cloud network peering container CIDR block for that project. MongoDB Cloud limits the number of MongoDB nodes per network peering connection based on the CIDR block and the region selected for the project.   **Example:** A project in an Amazon Web Services (AWS) region supporting three availability zones and an MongoDB CIDR network peering container block of limit of `/24` equals 27 three-node replica sets.
	AtlasCidrBlock *string `json:"atlasCidrBlock,omitempty"`
	// Unique string that identifies the Azure subscription in which the MongoDB Cloud VNet resides.
	// Read only field.
	AzureSubscriptionId *string `json:"azureSubscriptionId,omitempty"`
	// Azure region to which MongoDB Cloud deployed this network peering container.
	Region *string `json:"region,omitempty"`
	// Unique string that identifies the Azure VNet in which MongoDB Cloud clusters in this network peering container exist. The response returns **null** if no clusters exist in this network peering container.
	// Read only field.
	VnetName *string `json:"vnetName,omitempty"`
	// Unique string that identifies the GCP project in which MongoDB Cloud clusters in this network peering container exist. The response returns **null** if no clusters exist in this network peering container.
	// Read only field.
	GcpProjectId *string `json:"gcpProjectId,omitempty"`
	// Human-readable label that identifies the network in which MongoDB Cloud clusters in this network peering container exist. MongoDB Cloud returns **null** if no clusters exist in this network peering container.
	// Read only field.
	NetworkName *string `json:"networkName,omitempty"`
	// List of GCP regions to which you want to deploy this MongoDB Cloud network peering container.  In this MongoDB Cloud project, you can deploy clusters only to the GCP regions in this list. To deploy MongoDB Cloud clusters to other GCP regions, create additional projects.
	Regions *[]string `json:"regions,omitempty"`
	// Geographic area that Amazon Web Services (AWS) defines to which MongoDB Cloud deployed this network peering container.
	RegionName *string `json:"regionName,omitempty"`
	// Unique string that identifies the MongoDB Cloud VPC on AWS.
	// Read only field.
	VpcId *string `json:"vpcId,omitempty"`
}

// NewCloudProviderContainer instantiates a new CloudProviderContainer object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderContainer() *CloudProviderContainer {
	this := CloudProviderContainer{}
	return &this
}

// NewCloudProviderContainerWithDefaults instantiates a new CloudProviderContainer object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderContainerWithDefaults() *CloudProviderContainer {
	this := CloudProviderContainer{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *CloudProviderContainer) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CloudProviderContainer) SetId(v string) {
	o.Id = &v
}

// GetProviderName returns the ProviderName field value if set, zero value otherwise
func (o *CloudProviderContainer) GetProviderName() string {
	if o == nil || IsNil(o.ProviderName) {
		var ret string
		return ret
	}
	return *o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderName) {
		return nil, false
	}

	return o.ProviderName, true
}

// HasProviderName returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasProviderName() bool {
	if o != nil && !IsNil(o.ProviderName) {
		return true
	}

	return false
}

// SetProviderName gets a reference to the given string and assigns it to the ProviderName field.
func (o *CloudProviderContainer) SetProviderName(v string) {
	o.ProviderName = &v
}

// GetProvisioned returns the Provisioned field value if set, zero value otherwise
func (o *CloudProviderContainer) GetProvisioned() bool {
	if o == nil || IsNil(o.Provisioned) {
		var ret bool
		return ret
	}
	return *o.Provisioned
}

// GetProvisionedOk returns a tuple with the Provisioned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetProvisionedOk() (*bool, bool) {
	if o == nil || IsNil(o.Provisioned) {
		return nil, false
	}

	return o.Provisioned, true
}

// HasProvisioned returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasProvisioned() bool {
	if o != nil && !IsNil(o.Provisioned) {
		return true
	}

	return false
}

// SetProvisioned gets a reference to the given bool and assigns it to the Provisioned field.
func (o *CloudProviderContainer) SetProvisioned(v bool) {
	o.Provisioned = &v
}

// GetAtlasCidrBlock returns the AtlasCidrBlock field value if set, zero value otherwise
func (o *CloudProviderContainer) GetAtlasCidrBlock() string {
	if o == nil || IsNil(o.AtlasCidrBlock) {
		var ret string
		return ret
	}
	return *o.AtlasCidrBlock
}

// GetAtlasCidrBlockOk returns a tuple with the AtlasCidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetAtlasCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.AtlasCidrBlock) {
		return nil, false
	}

	return o.AtlasCidrBlock, true
}

// HasAtlasCidrBlock returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasAtlasCidrBlock() bool {
	if o != nil && !IsNil(o.AtlasCidrBlock) {
		return true
	}

	return false
}

// SetAtlasCidrBlock gets a reference to the given string and assigns it to the AtlasCidrBlock field.
func (o *CloudProviderContainer) SetAtlasCidrBlock(v string) {
	o.AtlasCidrBlock = &v
}

// GetAzureSubscriptionId returns the AzureSubscriptionId field value if set, zero value otherwise
func (o *CloudProviderContainer) GetAzureSubscriptionId() string {
	if o == nil || IsNil(o.AzureSubscriptionId) {
		var ret string
		return ret
	}
	return *o.AzureSubscriptionId
}

// GetAzureSubscriptionIdOk returns a tuple with the AzureSubscriptionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetAzureSubscriptionIdOk() (*string, bool) {
	if o == nil || IsNil(o.AzureSubscriptionId) {
		return nil, false
	}

	return o.AzureSubscriptionId, true
}

// HasAzureSubscriptionId returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasAzureSubscriptionId() bool {
	if o != nil && !IsNil(o.AzureSubscriptionId) {
		return true
	}

	return false
}

// SetAzureSubscriptionId gets a reference to the given string and assigns it to the AzureSubscriptionId field.
func (o *CloudProviderContainer) SetAzureSubscriptionId(v string) {
	o.AzureSubscriptionId = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *CloudProviderContainer) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *CloudProviderContainer) SetRegion(v string) {
	o.Region = &v
}

// GetVnetName returns the VnetName field value if set, zero value otherwise
func (o *CloudProviderContainer) GetVnetName() string {
	if o == nil || IsNil(o.VnetName) {
		var ret string
		return ret
	}
	return *o.VnetName
}

// GetVnetNameOk returns a tuple with the VnetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetVnetNameOk() (*string, bool) {
	if o == nil || IsNil(o.VnetName) {
		return nil, false
	}

	return o.VnetName, true
}

// HasVnetName returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasVnetName() bool {
	if o != nil && !IsNil(o.VnetName) {
		return true
	}

	return false
}

// SetVnetName gets a reference to the given string and assigns it to the VnetName field.
func (o *CloudProviderContainer) SetVnetName(v string) {
	o.VnetName = &v
}

// GetGcpProjectId returns the GcpProjectId field value if set, zero value otherwise
func (o *CloudProviderContainer) GetGcpProjectId() string {
	if o == nil || IsNil(o.GcpProjectId) {
		var ret string
		return ret
	}
	return *o.GcpProjectId
}

// GetGcpProjectIdOk returns a tuple with the GcpProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetGcpProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.GcpProjectId) {
		return nil, false
	}

	return o.GcpProjectId, true
}

// HasGcpProjectId returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasGcpProjectId() bool {
	if o != nil && !IsNil(o.GcpProjectId) {
		return true
	}

	return false
}

// SetGcpProjectId gets a reference to the given string and assigns it to the GcpProjectId field.
func (o *CloudProviderContainer) SetGcpProjectId(v string) {
	o.GcpProjectId = &v
}

// GetNetworkName returns the NetworkName field value if set, zero value otherwise
func (o *CloudProviderContainer) GetNetworkName() string {
	if o == nil || IsNil(o.NetworkName) {
		var ret string
		return ret
	}
	return *o.NetworkName
}

// GetNetworkNameOk returns a tuple with the NetworkName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetNetworkNameOk() (*string, bool) {
	if o == nil || IsNil(o.NetworkName) {
		return nil, false
	}

	return o.NetworkName, true
}

// HasNetworkName returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasNetworkName() bool {
	if o != nil && !IsNil(o.NetworkName) {
		return true
	}

	return false
}

// SetNetworkName gets a reference to the given string and assigns it to the NetworkName field.
func (o *CloudProviderContainer) SetNetworkName(v string) {
	o.NetworkName = &v
}

// GetRegions returns the Regions field value if set, zero value otherwise
func (o *CloudProviderContainer) GetRegions() []string {
	if o == nil || IsNil(o.Regions) {
		var ret []string
		return ret
	}
	return *o.Regions
}

// GetRegionsOk returns a tuple with the Regions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetRegionsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Regions) {
		return nil, false
	}

	return o.Regions, true
}

// HasRegions returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasRegions() bool {
	if o != nil && !IsNil(o.Regions) {
		return true
	}

	return false
}

// SetRegions gets a reference to the given []string and assigns it to the Regions field.
func (o *CloudProviderContainer) SetRegions(v []string) {
	o.Regions = &v
}

// GetRegionName returns the RegionName field value if set, zero value otherwise
func (o *CloudProviderContainer) GetRegionName() string {
	if o == nil || IsNil(o.RegionName) {
		var ret string
		return ret
	}
	return *o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionName) {
		return nil, false
	}

	return o.RegionName, true
}

// HasRegionName returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasRegionName() bool {
	if o != nil && !IsNil(o.RegionName) {
		return true
	}

	return false
}

// SetRegionName gets a reference to the given string and assigns it to the RegionName field.
func (o *CloudProviderContainer) SetRegionName(v string) {
	o.RegionName = &v
}

// GetVpcId returns the VpcId field value if set, zero value otherwise
func (o *CloudProviderContainer) GetVpcId() string {
	if o == nil || IsNil(o.VpcId) {
		var ret string
		return ret
	}
	return *o.VpcId
}

// GetVpcIdOk returns a tuple with the VpcId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderContainer) GetVpcIdOk() (*string, bool) {
	if o == nil || IsNil(o.VpcId) {
		return nil, false
	}

	return o.VpcId, true
}

// HasVpcId returns a boolean if a field has been set.
func (o *CloudProviderContainer) HasVpcId() bool {
	if o != nil && !IsNil(o.VpcId) {
		return true
	}

	return false
}

// SetVpcId gets a reference to the given string and assigns it to the VpcId field.
func (o *CloudProviderContainer) SetVpcId(v string) {
	o.VpcId = &v
}
