// Code based on the AtlasAPI V2 OpenAPI file

package admin

// InboundControlPlaneCloudProviderIPAddresses List of inbound IP addresses to the Atlas control plane, categorized by cloud provider. If your application allows outbound HTTP requests only to specific IP addresses, you must allow access to the following IP addresses so that your API requests can reach the Atlas control plane.
type InboundControlPlaneCloudProviderIPAddresses struct {
	// Control plane IP addresses in AWS. Each key identifies an Amazon Web Services (AWS) region. Each value identifies control plane IP addresses in the AWS region.
	// Read only field.
	Aws *map[string][]string `json:"aws,omitempty"`
	// Control plane IP addresses in Azure. Each key identifies an Azure region. Each value identifies control plane IP addresses in the Azure region.
	// Read only field.
	Azure *map[string][]string `json:"azure,omitempty"`
	// Control plane IP addresses in GCP. Each key identifies a Google Cloud (GCP) region. Each value identifies control plane IP addresses in the GCP region.
	// Read only field.
	Gcp *map[string][]string `json:"gcp,omitempty"`
}

// NewInboundControlPlaneCloudProviderIPAddresses instantiates a new InboundControlPlaneCloudProviderIPAddresses object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInboundControlPlaneCloudProviderIPAddresses() *InboundControlPlaneCloudProviderIPAddresses {
	this := InboundControlPlaneCloudProviderIPAddresses{}
	return &this
}

// NewInboundControlPlaneCloudProviderIPAddressesWithDefaults instantiates a new InboundControlPlaneCloudProviderIPAddresses object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInboundControlPlaneCloudProviderIPAddressesWithDefaults() *InboundControlPlaneCloudProviderIPAddresses {
	this := InboundControlPlaneCloudProviderIPAddresses{}
	return &this
}

// GetAws returns the Aws field value if set, zero value otherwise
func (o *InboundControlPlaneCloudProviderIPAddresses) GetAws() map[string][]string {
	if o == nil || IsNil(o.Aws) {
		var ret map[string][]string
		return ret
	}
	return *o.Aws
}

// GetAwsOk returns a tuple with the Aws field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InboundControlPlaneCloudProviderIPAddresses) GetAwsOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.Aws) {
		return nil, false
	}

	return o.Aws, true
}

// HasAws returns a boolean if a field has been set.
func (o *InboundControlPlaneCloudProviderIPAddresses) HasAws() bool {
	if o != nil && !IsNil(o.Aws) {
		return true
	}

	return false
}

// SetAws gets a reference to the given map[string][]string and assigns it to the Aws field.
func (o *InboundControlPlaneCloudProviderIPAddresses) SetAws(v map[string][]string) {
	o.Aws = &v
}

// GetAzure returns the Azure field value if set, zero value otherwise
func (o *InboundControlPlaneCloudProviderIPAddresses) GetAzure() map[string][]string {
	if o == nil || IsNil(o.Azure) {
		var ret map[string][]string
		return ret
	}
	return *o.Azure
}

// GetAzureOk returns a tuple with the Azure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InboundControlPlaneCloudProviderIPAddresses) GetAzureOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.Azure) {
		return nil, false
	}

	return o.Azure, true
}

// HasAzure returns a boolean if a field has been set.
func (o *InboundControlPlaneCloudProviderIPAddresses) HasAzure() bool {
	if o != nil && !IsNil(o.Azure) {
		return true
	}

	return false
}

// SetAzure gets a reference to the given map[string][]string and assigns it to the Azure field.
func (o *InboundControlPlaneCloudProviderIPAddresses) SetAzure(v map[string][]string) {
	o.Azure = &v
}

// GetGcp returns the Gcp field value if set, zero value otherwise
func (o *InboundControlPlaneCloudProviderIPAddresses) GetGcp() map[string][]string {
	if o == nil || IsNil(o.Gcp) {
		var ret map[string][]string
		return ret
	}
	return *o.Gcp
}

// GetGcpOk returns a tuple with the Gcp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InboundControlPlaneCloudProviderIPAddresses) GetGcpOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.Gcp) {
		return nil, false
	}

	return o.Gcp, true
}

// HasGcp returns a boolean if a field has been set.
func (o *InboundControlPlaneCloudProviderIPAddresses) HasGcp() bool {
	if o != nil && !IsNil(o.Gcp) {
		return true
	}

	return false
}

// SetGcp gets a reference to the given map[string][]string and assigns it to the Gcp field.
func (o *InboundControlPlaneCloudProviderIPAddresses) SetGcp(v map[string][]string) {
	o.Gcp = &v
}
