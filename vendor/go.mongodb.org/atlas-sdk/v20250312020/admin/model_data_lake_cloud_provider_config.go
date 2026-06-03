// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeCloudProviderConfig Cloud provider where this Federated Database Instance is hosted.
type DataLakeCloudProviderConfig struct {
	Aws   *DataLakeAWSCloudProviderConfig         `json:"aws,omitempty"`
	Azure *DataFederationAzureCloudProviderConfig `json:"azure,omitempty"`
	Gcp   *DataFederationGCPCloudProviderConfig   `json:"gcp,omitempty"`
}

// NewDataLakeCloudProviderConfig instantiates a new DataLakeCloudProviderConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeCloudProviderConfig() *DataLakeCloudProviderConfig {
	this := DataLakeCloudProviderConfig{}
	return &this
}

// NewDataLakeCloudProviderConfigWithDefaults instantiates a new DataLakeCloudProviderConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeCloudProviderConfigWithDefaults() *DataLakeCloudProviderConfig {
	this := DataLakeCloudProviderConfig{}
	return &this
}

// GetAws returns the Aws field value if set, zero value otherwise
func (o *DataLakeCloudProviderConfig) GetAws() DataLakeAWSCloudProviderConfig {
	if o == nil || IsNil(o.Aws) {
		var ret DataLakeAWSCloudProviderConfig
		return ret
	}
	return *o.Aws
}

// GetAwsOk returns a tuple with the Aws field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeCloudProviderConfig) GetAwsOk() (*DataLakeAWSCloudProviderConfig, bool) {
	if o == nil || IsNil(o.Aws) {
		return nil, false
	}

	return o.Aws, true
}

// HasAws returns a boolean if a field has been set.
func (o *DataLakeCloudProviderConfig) HasAws() bool {
	if o != nil && !IsNil(o.Aws) {
		return true
	}

	return false
}

// SetAws gets a reference to the given DataLakeAWSCloudProviderConfig and assigns it to the Aws field.
func (o *DataLakeCloudProviderConfig) SetAws(v DataLakeAWSCloudProviderConfig) {
	o.Aws = &v
}

// GetAzure returns the Azure field value if set, zero value otherwise
func (o *DataLakeCloudProviderConfig) GetAzure() DataFederationAzureCloudProviderConfig {
	if o == nil || IsNil(o.Azure) {
		var ret DataFederationAzureCloudProviderConfig
		return ret
	}
	return *o.Azure
}

// GetAzureOk returns a tuple with the Azure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeCloudProviderConfig) GetAzureOk() (*DataFederationAzureCloudProviderConfig, bool) {
	if o == nil || IsNil(o.Azure) {
		return nil, false
	}

	return o.Azure, true
}

// HasAzure returns a boolean if a field has been set.
func (o *DataLakeCloudProviderConfig) HasAzure() bool {
	if o != nil && !IsNil(o.Azure) {
		return true
	}

	return false
}

// SetAzure gets a reference to the given DataFederationAzureCloudProviderConfig and assigns it to the Azure field.
func (o *DataLakeCloudProviderConfig) SetAzure(v DataFederationAzureCloudProviderConfig) {
	o.Azure = &v
}

// GetGcp returns the Gcp field value if set, zero value otherwise
func (o *DataLakeCloudProviderConfig) GetGcp() DataFederationGCPCloudProviderConfig {
	if o == nil || IsNil(o.Gcp) {
		var ret DataFederationGCPCloudProviderConfig
		return ret
	}
	return *o.Gcp
}

// GetGcpOk returns a tuple with the Gcp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeCloudProviderConfig) GetGcpOk() (*DataFederationGCPCloudProviderConfig, bool) {
	if o == nil || IsNil(o.Gcp) {
		return nil, false
	}

	return o.Gcp, true
}

// HasGcp returns a boolean if a field has been set.
func (o *DataLakeCloudProviderConfig) HasGcp() bool {
	if o != nil && !IsNil(o.Gcp) {
		return true
	}

	return false
}

// SetGcp gets a reference to the given DataFederationGCPCloudProviderConfig and assigns it to the Gcp field.
func (o *DataLakeCloudProviderConfig) SetGcp(v DataFederationGCPCloudProviderConfig) {
	o.Gcp = &v
}
