// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServerlessProviderSettings Group of cloud provider settings that configure the provisioned MongoDB serverless instance.
type ServerlessProviderSettings struct {
	// Cloud service provider on which MongoDB Cloud provisioned the serverless instance.
	BackingProviderName string `json:"backingProviderName"`
	// Storage capacity of instance data volumes expressed in gigabytes. This value is not configurable for Serverless or effectively Flex clusters.
	// Read only field.
	EffectiveDiskSizeGBLimit *int `json:"effectiveDiskSizeGBLimit,omitempty"`
	// Instance size boundary to which your cluster can automatically scale.
	// Read only field.
	EffectiveInstanceSizeName *string `json:"effectiveInstanceSizeName,omitempty"`
	// Cloud service provider on which MongoDB Cloud effectively provisioned the serverless instance.
	// Read only field.
	EffectiveProviderName *string `json:"effectiveProviderName,omitempty"`
	// Human-readable label that identifies the cloud service provider.
	ProviderName *string `json:"providerName,omitempty"`
	// Human-readable label that identifies the geographic location of your MongoDB serverless instance. The region you choose can affect network latency for clients accessing your databases. For a complete list of region names, see [AWS](https://docs.atlas.mongodb.com/reference/amazon-aws/#std-label-amazon-aws), [GCP](https://docs.atlas.mongodb.com/reference/google-gcp/), and [Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
	RegionName string `json:"regionName"`
}

// NewServerlessProviderSettings instantiates a new ServerlessProviderSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessProviderSettings(backingProviderName string, regionName string) *ServerlessProviderSettings {
	this := ServerlessProviderSettings{}
	this.BackingProviderName = backingProviderName
	var providerName string = "SERVERLESS"
	this.ProviderName = &providerName
	this.RegionName = regionName
	return &this
}

// NewServerlessProviderSettingsWithDefaults instantiates a new ServerlessProviderSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessProviderSettingsWithDefaults() *ServerlessProviderSettings {
	this := ServerlessProviderSettings{}
	var providerName string = "SERVERLESS"
	this.ProviderName = &providerName
	return &this
}

// GetBackingProviderName returns the BackingProviderName field value
func (o *ServerlessProviderSettings) GetBackingProviderName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BackingProviderName
}

// GetBackingProviderNameOk returns a tuple with the BackingProviderName field value
// and a boolean to check if the value has been set.
func (o *ServerlessProviderSettings) GetBackingProviderNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BackingProviderName, true
}

// SetBackingProviderName sets field value
func (o *ServerlessProviderSettings) SetBackingProviderName(v string) {
	o.BackingProviderName = v
}

// GetEffectiveDiskSizeGBLimit returns the EffectiveDiskSizeGBLimit field value if set, zero value otherwise
func (o *ServerlessProviderSettings) GetEffectiveDiskSizeGBLimit() int {
	if o == nil || IsNil(o.EffectiveDiskSizeGBLimit) {
		var ret int
		return ret
	}
	return *o.EffectiveDiskSizeGBLimit
}

// GetEffectiveDiskSizeGBLimitOk returns a tuple with the EffectiveDiskSizeGBLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessProviderSettings) GetEffectiveDiskSizeGBLimitOk() (*int, bool) {
	if o == nil || IsNil(o.EffectiveDiskSizeGBLimit) {
		return nil, false
	}

	return o.EffectiveDiskSizeGBLimit, true
}

// HasEffectiveDiskSizeGBLimit returns a boolean if a field has been set.
func (o *ServerlessProviderSettings) HasEffectiveDiskSizeGBLimit() bool {
	if o != nil && !IsNil(o.EffectiveDiskSizeGBLimit) {
		return true
	}

	return false
}

// SetEffectiveDiskSizeGBLimit gets a reference to the given int and assigns it to the EffectiveDiskSizeGBLimit field.
func (o *ServerlessProviderSettings) SetEffectiveDiskSizeGBLimit(v int) {
	o.EffectiveDiskSizeGBLimit = &v
}

// GetEffectiveInstanceSizeName returns the EffectiveInstanceSizeName field value if set, zero value otherwise
func (o *ServerlessProviderSettings) GetEffectiveInstanceSizeName() string {
	if o == nil || IsNil(o.EffectiveInstanceSizeName) {
		var ret string
		return ret
	}
	return *o.EffectiveInstanceSizeName
}

// GetEffectiveInstanceSizeNameOk returns a tuple with the EffectiveInstanceSizeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessProviderSettings) GetEffectiveInstanceSizeNameOk() (*string, bool) {
	if o == nil || IsNil(o.EffectiveInstanceSizeName) {
		return nil, false
	}

	return o.EffectiveInstanceSizeName, true
}

// HasEffectiveInstanceSizeName returns a boolean if a field has been set.
func (o *ServerlessProviderSettings) HasEffectiveInstanceSizeName() bool {
	if o != nil && !IsNil(o.EffectiveInstanceSizeName) {
		return true
	}

	return false
}

// SetEffectiveInstanceSizeName gets a reference to the given string and assigns it to the EffectiveInstanceSizeName field.
func (o *ServerlessProviderSettings) SetEffectiveInstanceSizeName(v string) {
	o.EffectiveInstanceSizeName = &v
}

// GetEffectiveProviderName returns the EffectiveProviderName field value if set, zero value otherwise
func (o *ServerlessProviderSettings) GetEffectiveProviderName() string {
	if o == nil || IsNil(o.EffectiveProviderName) {
		var ret string
		return ret
	}
	return *o.EffectiveProviderName
}

// GetEffectiveProviderNameOk returns a tuple with the EffectiveProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessProviderSettings) GetEffectiveProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.EffectiveProviderName) {
		return nil, false
	}

	return o.EffectiveProviderName, true
}

// HasEffectiveProviderName returns a boolean if a field has been set.
func (o *ServerlessProviderSettings) HasEffectiveProviderName() bool {
	if o != nil && !IsNil(o.EffectiveProviderName) {
		return true
	}

	return false
}

// SetEffectiveProviderName gets a reference to the given string and assigns it to the EffectiveProviderName field.
func (o *ServerlessProviderSettings) SetEffectiveProviderName(v string) {
	o.EffectiveProviderName = &v
}

// GetProviderName returns the ProviderName field value if set, zero value otherwise
func (o *ServerlessProviderSettings) GetProviderName() string {
	if o == nil || IsNil(o.ProviderName) {
		var ret string
		return ret
	}
	return *o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessProviderSettings) GetProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderName) {
		return nil, false
	}

	return o.ProviderName, true
}

// HasProviderName returns a boolean if a field has been set.
func (o *ServerlessProviderSettings) HasProviderName() bool {
	if o != nil && !IsNil(o.ProviderName) {
		return true
	}

	return false
}

// SetProviderName gets a reference to the given string and assigns it to the ProviderName field.
func (o *ServerlessProviderSettings) SetProviderName(v string) {
	o.ProviderName = &v
}

// GetRegionName returns the RegionName field value
func (o *ServerlessProviderSettings) GetRegionName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value
// and a boolean to check if the value has been set.
func (o *ServerlessProviderSettings) GetRegionNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegionName, true
}

// SetRegionName sets field value
func (o *ServerlessProviderSettings) SetRegionName(v string) {
	o.RegionName = v
}
