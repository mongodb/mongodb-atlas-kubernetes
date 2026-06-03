// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AzureConnection Azure-specific configuration for the connection.
type AzureConnection struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Azure region where the storage account is located.
	// Deprecated
	Region *string `json:"region,omitempty"`
	// Unique ID of the Azure Service Principal that has access to the storage account.
	ServicePrincipalId *string `json:"servicePrincipalId,omitempty"`
	// Name of the Azure Storage Account to connect to.
	StorageAccountName *string `json:"storageAccountName,omitempty"`
}

// NewAzureConnection instantiates a new AzureConnection object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAzureConnection() *AzureConnection {
	this := AzureConnection{}
	return &this
}

// NewAzureConnectionWithDefaults instantiates a new AzureConnection object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAzureConnectionWithDefaults() *AzureConnection {
	this := AzureConnection{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *AzureConnection) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureConnection) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *AzureConnection) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *AzureConnection) SetLinks(v []Link) {
	o.Links = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
// Deprecated
func (o *AzureConnection) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *AzureConnection) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *AzureConnection) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
// Deprecated
func (o *AzureConnection) SetRegion(v string) {
	o.Region = &v
}

// GetServicePrincipalId returns the ServicePrincipalId field value if set, zero value otherwise
func (o *AzureConnection) GetServicePrincipalId() string {
	if o == nil || IsNil(o.ServicePrincipalId) {
		var ret string
		return ret
	}
	return *o.ServicePrincipalId
}

// GetServicePrincipalIdOk returns a tuple with the ServicePrincipalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureConnection) GetServicePrincipalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ServicePrincipalId) {
		return nil, false
	}

	return o.ServicePrincipalId, true
}

// HasServicePrincipalId returns a boolean if a field has been set.
func (o *AzureConnection) HasServicePrincipalId() bool {
	if o != nil && !IsNil(o.ServicePrincipalId) {
		return true
	}

	return false
}

// SetServicePrincipalId gets a reference to the given string and assigns it to the ServicePrincipalId field.
func (o *AzureConnection) SetServicePrincipalId(v string) {
	o.ServicePrincipalId = &v
}

// GetStorageAccountName returns the StorageAccountName field value if set, zero value otherwise
func (o *AzureConnection) GetStorageAccountName() string {
	if o == nil || IsNil(o.StorageAccountName) {
		var ret string
		return ret
	}
	return *o.StorageAccountName
}

// GetStorageAccountNameOk returns a tuple with the StorageAccountName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AzureConnection) GetStorageAccountNameOk() (*string, bool) {
	if o == nil || IsNil(o.StorageAccountName) {
		return nil, false
	}

	return o.StorageAccountName, true
}

// HasStorageAccountName returns a boolean if a field has been set.
func (o *AzureConnection) HasStorageAccountName() bool {
	if o != nil && !IsNil(o.StorageAccountName) {
		return true
	}

	return false
}

// SetStorageAccountName gets a reference to the given string and assigns it to the StorageAccountName field.
func (o *AzureConnection) SetStorageAccountName(v string) {
	o.StorageAccountName = &v
}
