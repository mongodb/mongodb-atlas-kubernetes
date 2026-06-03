// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BaseNetworkPeeringConnectionSettings struct for BaseNetworkPeeringConnectionSettings
type BaseNetworkPeeringConnectionSettings struct {
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container that contains the specified network peering connection.
	ContainerId string `json:"containerId"`
	// Unique 24-hexadecimal digit string that identifies the network peering connection.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Cloud service provider that serves the requested network peering connection.
	ProviderName *string `json:"providerName,omitempty"`
	// Amazon Web Services (AWS) region where the Virtual Peering Connection (VPC) that you peered with the MongoDB Cloud VPC resides. The resource returns `null` if your VPC and the MongoDB Cloud VPC reside in the same region.
	AccepterRegionName *string `json:"accepterRegionName,omitempty"`
	// Unique twelve-digit string that identifies the Amazon Web Services (AWS) account that owns the VPC that you peered with the MongoDB Cloud VPC.
	AwsAccountId *string `json:"awsAccountId,omitempty"`
	// Unique string that identifies the peering connection on AWS.
	// Read only field.
	ConnectionId *string `json:"connectionId,omitempty"`
	// Type of error that can be returned when requesting an Amazon Web Services (AWS) peering connection. The resource returns `null` if the request succeeded.
	// Read only field.
	ErrorStateName *string `json:"errorStateName,omitempty"`
	// Internet Protocol (IP) addresses expressed in Classless Inter-Domain Routing (CIDR) notation of the VPC's subnet that you want to peer with the MongoDB Cloud VPC.
	RouteTableCidrBlock *string `json:"routeTableCidrBlock,omitempty"`
	// State of the network peering connection at the time you made the request.
	// Read only field.
	StatusName *string `json:"statusName,omitempty"`
	// Unique string that identifies the VPC on Amazon Web Services (AWS) that you want to peer with the MongoDB Cloud VPC.
	VpcId *string `json:"vpcId,omitempty"`
	// Unique string that identifies the Azure AD directory in which the VNet peered with the MongoDB Cloud VNet resides.
	AzureDirectoryId *string `json:"azureDirectoryId,omitempty"`
	// Unique string that identifies the Azure subscription in which the VNet you peered with the MongoDB Cloud VNet resides.
	AzureSubscriptionId *string `json:"azureSubscriptionId,omitempty"`
	// Error message returned when a requested Azure network peering resource returns `\"status\" : \"FAILED\"`. The resource returns `null` if the request succeeded.
	// Read only field.
	ErrorState *string `json:"errorState,omitempty"`
	// Human-readable label that identifies the resource group in which the VNet to peer with the MongoDB Cloud VNet resides.
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	// State of the network peering connection at the time you made the request.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Human-readable label that identifies the VNet that you want to peer with the MongoDB Cloud VNet.
	VnetName *string `json:"vnetName,omitempty"`
	// Details of the error returned when requesting a GCP network peering resource. The resource returns `null` if the request succeeded.
	// Read only field.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Human-readable label that identifies the GCP project that contains the network that you want to peer with the MongoDB Cloud VPC.
	GcpProjectId *string `json:"gcpProjectId,omitempty"`
	// Human-readable label that identifies the network to peer with the MongoDB Cloud VPC.
	NetworkName *string `json:"networkName,omitempty"`
}

// NewBaseNetworkPeeringConnectionSettings instantiates a new BaseNetworkPeeringConnectionSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBaseNetworkPeeringConnectionSettings(containerId string) *BaseNetworkPeeringConnectionSettings {
	this := BaseNetworkPeeringConnectionSettings{}
	this.ContainerId = containerId
	return &this
}

// NewBaseNetworkPeeringConnectionSettingsWithDefaults instantiates a new BaseNetworkPeeringConnectionSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBaseNetworkPeeringConnectionSettingsWithDefaults() *BaseNetworkPeeringConnectionSettings {
	this := BaseNetworkPeeringConnectionSettings{}
	return &this
}

// GetContainerId returns the ContainerId field value
func (o *BaseNetworkPeeringConnectionSettings) GetContainerId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ContainerId
}

// GetContainerIdOk returns a tuple with the ContainerId field value
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetContainerIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ContainerId, true
}

// SetContainerId sets field value
func (o *BaseNetworkPeeringConnectionSettings) SetContainerId(v string) {
	o.ContainerId = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BaseNetworkPeeringConnectionSettings) SetId(v string) {
	o.Id = &v
}

// GetProviderName returns the ProviderName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetProviderName() string {
	if o == nil || IsNil(o.ProviderName) {
		var ret string
		return ret
	}
	return *o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderName) {
		return nil, false
	}

	return o.ProviderName, true
}

// HasProviderName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasProviderName() bool {
	if o != nil && !IsNil(o.ProviderName) {
		return true
	}

	return false
}

// SetProviderName gets a reference to the given string and assigns it to the ProviderName field.
func (o *BaseNetworkPeeringConnectionSettings) SetProviderName(v string) {
	o.ProviderName = &v
}

// GetAccepterRegionName returns the AccepterRegionName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetAccepterRegionName() string {
	if o == nil || IsNil(o.AccepterRegionName) {
		var ret string
		return ret
	}
	return *o.AccepterRegionName
}

// GetAccepterRegionNameOk returns a tuple with the AccepterRegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetAccepterRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.AccepterRegionName) {
		return nil, false
	}

	return o.AccepterRegionName, true
}

// HasAccepterRegionName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasAccepterRegionName() bool {
	if o != nil && !IsNil(o.AccepterRegionName) {
		return true
	}

	return false
}

// SetAccepterRegionName gets a reference to the given string and assigns it to the AccepterRegionName field.
func (o *BaseNetworkPeeringConnectionSettings) SetAccepterRegionName(v string) {
	o.AccepterRegionName = &v
}

// GetAwsAccountId returns the AwsAccountId field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetAwsAccountId() string {
	if o == nil || IsNil(o.AwsAccountId) {
		var ret string
		return ret
	}
	return *o.AwsAccountId
}

// GetAwsAccountIdOk returns a tuple with the AwsAccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetAwsAccountIdOk() (*string, bool) {
	if o == nil || IsNil(o.AwsAccountId) {
		return nil, false
	}

	return o.AwsAccountId, true
}

// HasAwsAccountId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasAwsAccountId() bool {
	if o != nil && !IsNil(o.AwsAccountId) {
		return true
	}

	return false
}

// SetAwsAccountId gets a reference to the given string and assigns it to the AwsAccountId field.
func (o *BaseNetworkPeeringConnectionSettings) SetAwsAccountId(v string) {
	o.AwsAccountId = &v
}

// GetConnectionId returns the ConnectionId field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetConnectionId() string {
	if o == nil || IsNil(o.ConnectionId) {
		var ret string
		return ret
	}
	return *o.ConnectionId
}

// GetConnectionIdOk returns a tuple with the ConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetConnectionIdOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectionId) {
		return nil, false
	}

	return o.ConnectionId, true
}

// HasConnectionId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasConnectionId() bool {
	if o != nil && !IsNil(o.ConnectionId) {
		return true
	}

	return false
}

// SetConnectionId gets a reference to the given string and assigns it to the ConnectionId field.
func (o *BaseNetworkPeeringConnectionSettings) SetConnectionId(v string) {
	o.ConnectionId = &v
}

// GetErrorStateName returns the ErrorStateName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetErrorStateName() string {
	if o == nil || IsNil(o.ErrorStateName) {
		var ret string
		return ret
	}
	return *o.ErrorStateName
}

// GetErrorStateNameOk returns a tuple with the ErrorStateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetErrorStateNameOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorStateName) {
		return nil, false
	}

	return o.ErrorStateName, true
}

// HasErrorStateName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasErrorStateName() bool {
	if o != nil && !IsNil(o.ErrorStateName) {
		return true
	}

	return false
}

// SetErrorStateName gets a reference to the given string and assigns it to the ErrorStateName field.
func (o *BaseNetworkPeeringConnectionSettings) SetErrorStateName(v string) {
	o.ErrorStateName = &v
}

// GetRouteTableCidrBlock returns the RouteTableCidrBlock field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetRouteTableCidrBlock() string {
	if o == nil || IsNil(o.RouteTableCidrBlock) {
		var ret string
		return ret
	}
	return *o.RouteTableCidrBlock
}

// GetRouteTableCidrBlockOk returns a tuple with the RouteTableCidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetRouteTableCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.RouteTableCidrBlock) {
		return nil, false
	}

	return o.RouteTableCidrBlock, true
}

// HasRouteTableCidrBlock returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasRouteTableCidrBlock() bool {
	if o != nil && !IsNil(o.RouteTableCidrBlock) {
		return true
	}

	return false
}

// SetRouteTableCidrBlock gets a reference to the given string and assigns it to the RouteTableCidrBlock field.
func (o *BaseNetworkPeeringConnectionSettings) SetRouteTableCidrBlock(v string) {
	o.RouteTableCidrBlock = &v
}

// GetStatusName returns the StatusName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetStatusName() string {
	if o == nil || IsNil(o.StatusName) {
		var ret string
		return ret
	}
	return *o.StatusName
}

// GetStatusNameOk returns a tuple with the StatusName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetStatusNameOk() (*string, bool) {
	if o == nil || IsNil(o.StatusName) {
		return nil, false
	}

	return o.StatusName, true
}

// HasStatusName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasStatusName() bool {
	if o != nil && !IsNil(o.StatusName) {
		return true
	}

	return false
}

// SetStatusName gets a reference to the given string and assigns it to the StatusName field.
func (o *BaseNetworkPeeringConnectionSettings) SetStatusName(v string) {
	o.StatusName = &v
}

// GetVpcId returns the VpcId field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetVpcId() string {
	if o == nil || IsNil(o.VpcId) {
		var ret string
		return ret
	}
	return *o.VpcId
}

// GetVpcIdOk returns a tuple with the VpcId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetVpcIdOk() (*string, bool) {
	if o == nil || IsNil(o.VpcId) {
		return nil, false
	}

	return o.VpcId, true
}

// HasVpcId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasVpcId() bool {
	if o != nil && !IsNil(o.VpcId) {
		return true
	}

	return false
}

// SetVpcId gets a reference to the given string and assigns it to the VpcId field.
func (o *BaseNetworkPeeringConnectionSettings) SetVpcId(v string) {
	o.VpcId = &v
}

// GetAzureDirectoryId returns the AzureDirectoryId field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetAzureDirectoryId() string {
	if o == nil || IsNil(o.AzureDirectoryId) {
		var ret string
		return ret
	}
	return *o.AzureDirectoryId
}

// GetAzureDirectoryIdOk returns a tuple with the AzureDirectoryId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetAzureDirectoryIdOk() (*string, bool) {
	if o == nil || IsNil(o.AzureDirectoryId) {
		return nil, false
	}

	return o.AzureDirectoryId, true
}

// HasAzureDirectoryId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasAzureDirectoryId() bool {
	if o != nil && !IsNil(o.AzureDirectoryId) {
		return true
	}

	return false
}

// SetAzureDirectoryId gets a reference to the given string and assigns it to the AzureDirectoryId field.
func (o *BaseNetworkPeeringConnectionSettings) SetAzureDirectoryId(v string) {
	o.AzureDirectoryId = &v
}

// GetAzureSubscriptionId returns the AzureSubscriptionId field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetAzureSubscriptionId() string {
	if o == nil || IsNil(o.AzureSubscriptionId) {
		var ret string
		return ret
	}
	return *o.AzureSubscriptionId
}

// GetAzureSubscriptionIdOk returns a tuple with the AzureSubscriptionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetAzureSubscriptionIdOk() (*string, bool) {
	if o == nil || IsNil(o.AzureSubscriptionId) {
		return nil, false
	}

	return o.AzureSubscriptionId, true
}

// HasAzureSubscriptionId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasAzureSubscriptionId() bool {
	if o != nil && !IsNil(o.AzureSubscriptionId) {
		return true
	}

	return false
}

// SetAzureSubscriptionId gets a reference to the given string and assigns it to the AzureSubscriptionId field.
func (o *BaseNetworkPeeringConnectionSettings) SetAzureSubscriptionId(v string) {
	o.AzureSubscriptionId = &v
}

// GetErrorState returns the ErrorState field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetErrorState() string {
	if o == nil || IsNil(o.ErrorState) {
		var ret string
		return ret
	}
	return *o.ErrorState
}

// GetErrorStateOk returns a tuple with the ErrorState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetErrorStateOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorState) {
		return nil, false
	}

	return o.ErrorState, true
}

// HasErrorState returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasErrorState() bool {
	if o != nil && !IsNil(o.ErrorState) {
		return true
	}

	return false
}

// SetErrorState gets a reference to the given string and assigns it to the ErrorState field.
func (o *BaseNetworkPeeringConnectionSettings) SetErrorState(v string) {
	o.ErrorState = &v
}

// GetResourceGroupName returns the ResourceGroupName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetResourceGroupName() string {
	if o == nil || IsNil(o.ResourceGroupName) {
		var ret string
		return ret
	}
	return *o.ResourceGroupName
}

// GetResourceGroupNameOk returns a tuple with the ResourceGroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetResourceGroupNameOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceGroupName) {
		return nil, false
	}

	return o.ResourceGroupName, true
}

// HasResourceGroupName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasResourceGroupName() bool {
	if o != nil && !IsNil(o.ResourceGroupName) {
		return true
	}

	return false
}

// SetResourceGroupName gets a reference to the given string and assigns it to the ResourceGroupName field.
func (o *BaseNetworkPeeringConnectionSettings) SetResourceGroupName(v string) {
	o.ResourceGroupName = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *BaseNetworkPeeringConnectionSettings) SetStatus(v string) {
	o.Status = &v
}

// GetVnetName returns the VnetName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetVnetName() string {
	if o == nil || IsNil(o.VnetName) {
		var ret string
		return ret
	}
	return *o.VnetName
}

// GetVnetNameOk returns a tuple with the VnetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetVnetNameOk() (*string, bool) {
	if o == nil || IsNil(o.VnetName) {
		return nil, false
	}

	return o.VnetName, true
}

// HasVnetName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasVnetName() bool {
	if o != nil && !IsNil(o.VnetName) {
		return true
	}

	return false
}

// SetVnetName gets a reference to the given string and assigns it to the VnetName field.
func (o *BaseNetworkPeeringConnectionSettings) SetVnetName(v string) {
	o.VnetName = &v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}

	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *BaseNetworkPeeringConnectionSettings) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

// GetGcpProjectId returns the GcpProjectId field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetGcpProjectId() string {
	if o == nil || IsNil(o.GcpProjectId) {
		var ret string
		return ret
	}
	return *o.GcpProjectId
}

// GetGcpProjectIdOk returns a tuple with the GcpProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetGcpProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.GcpProjectId) {
		return nil, false
	}

	return o.GcpProjectId, true
}

// HasGcpProjectId returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasGcpProjectId() bool {
	if o != nil && !IsNil(o.GcpProjectId) {
		return true
	}

	return false
}

// SetGcpProjectId gets a reference to the given string and assigns it to the GcpProjectId field.
func (o *BaseNetworkPeeringConnectionSettings) SetGcpProjectId(v string) {
	o.GcpProjectId = &v
}

// GetNetworkName returns the NetworkName field value if set, zero value otherwise
func (o *BaseNetworkPeeringConnectionSettings) GetNetworkName() string {
	if o == nil || IsNil(o.NetworkName) {
		var ret string
		return ret
	}
	return *o.NetworkName
}

// GetNetworkNameOk returns a tuple with the NetworkName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseNetworkPeeringConnectionSettings) GetNetworkNameOk() (*string, bool) {
	if o == nil || IsNil(o.NetworkName) {
		return nil, false
	}

	return o.NetworkName, true
}

// HasNetworkName returns a boolean if a field has been set.
func (o *BaseNetworkPeeringConnectionSettings) HasNetworkName() bool {
	if o != nil && !IsNil(o.NetworkName) {
		return true
	}

	return false
}

// SetNetworkName gets a reference to the given string and assigns it to the NetworkName field.
func (o *BaseNetworkPeeringConnectionSettings) SetNetworkName(v string) {
	o.NetworkName = &v
}
