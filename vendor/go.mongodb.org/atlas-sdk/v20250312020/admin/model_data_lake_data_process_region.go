// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeDataProcessRegion Information about the cloud provider region to which the Federated Database Instance routes client connections.
type DataLakeDataProcessRegion struct {
	// Name of the cloud service that hosts the Federated Database Instance's infrastructure.
	CloudProvider string `json:"cloudProvider"`
	// Name of the region to which the data lake routes client connections.
	Region string `json:"region"`
}

// NewDataLakeDataProcessRegion instantiates a new DataLakeDataProcessRegion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeDataProcessRegion(cloudProvider string, region string) *DataLakeDataProcessRegion {
	this := DataLakeDataProcessRegion{}
	this.CloudProvider = cloudProvider
	this.Region = region
	return &this
}

// NewDataLakeDataProcessRegionWithDefaults instantiates a new DataLakeDataProcessRegion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeDataProcessRegionWithDefaults() *DataLakeDataProcessRegion {
	this := DataLakeDataProcessRegion{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value
func (o *DataLakeDataProcessRegion) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *DataLakeDataProcessRegion) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *DataLakeDataProcessRegion) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetRegion returns the Region field value
func (o *DataLakeDataProcessRegion) GetRegion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Region
}

// GetRegionOk returns a tuple with the Region field value
// and a boolean to check if the value has been set.
func (o *DataLakeDataProcessRegion) GetRegionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Region, true
}

// SetRegion sets field value
func (o *DataLakeDataProcessRegion) SetRegion(v string) {
	o.Region = v
}
