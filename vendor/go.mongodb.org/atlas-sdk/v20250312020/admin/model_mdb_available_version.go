// Code based on the AtlasAPI V2 OpenAPI file

package admin

// MdbAvailableVersion struct for MdbAvailableVersion
type MdbAvailableVersion struct {
	// Cloud service provider on which MongoDB Cloud provisions the hosts. Set dedicated clusters to `AWS`, `GCP`, `AZURE` or `TENANT`.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// Whether the version is the current default for the Instance Size and Cloud Provider.
	DefaultStatus *string `json:"defaultStatus,omitempty"`
	// Instance size boundary to which your cluster can automatically scale.
	// Read only field.
	InstanceSize *string `json:"instanceSize,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// The MongoDB Major Version in question.
	Version *string `json:"version,omitempty"`
}

// NewMdbAvailableVersion instantiates a new MdbAvailableVersion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMdbAvailableVersion() *MdbAvailableVersion {
	this := MdbAvailableVersion{}
	return &this
}

// NewMdbAvailableVersionWithDefaults instantiates a new MdbAvailableVersion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMdbAvailableVersionWithDefaults() *MdbAvailableVersion {
	this := MdbAvailableVersion{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *MdbAvailableVersion) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MdbAvailableVersion) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *MdbAvailableVersion) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *MdbAvailableVersion) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetDefaultStatus returns the DefaultStatus field value if set, zero value otherwise
func (o *MdbAvailableVersion) GetDefaultStatus() string {
	if o == nil || IsNil(o.DefaultStatus) {
		var ret string
		return ret
	}
	return *o.DefaultStatus
}

// GetDefaultStatusOk returns a tuple with the DefaultStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MdbAvailableVersion) GetDefaultStatusOk() (*string, bool) {
	if o == nil || IsNil(o.DefaultStatus) {
		return nil, false
	}

	return o.DefaultStatus, true
}

// HasDefaultStatus returns a boolean if a field has been set.
func (o *MdbAvailableVersion) HasDefaultStatus() bool {
	if o != nil && !IsNil(o.DefaultStatus) {
		return true
	}

	return false
}

// SetDefaultStatus gets a reference to the given string and assigns it to the DefaultStatus field.
func (o *MdbAvailableVersion) SetDefaultStatus(v string) {
	o.DefaultStatus = &v
}

// GetInstanceSize returns the InstanceSize field value if set, zero value otherwise
func (o *MdbAvailableVersion) GetInstanceSize() string {
	if o == nil || IsNil(o.InstanceSize) {
		var ret string
		return ret
	}
	return *o.InstanceSize
}

// GetInstanceSizeOk returns a tuple with the InstanceSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MdbAvailableVersion) GetInstanceSizeOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceSize) {
		return nil, false
	}

	return o.InstanceSize, true
}

// HasInstanceSize returns a boolean if a field has been set.
func (o *MdbAvailableVersion) HasInstanceSize() bool {
	if o != nil && !IsNil(o.InstanceSize) {
		return true
	}

	return false
}

// SetInstanceSize gets a reference to the given string and assigns it to the InstanceSize field.
func (o *MdbAvailableVersion) SetInstanceSize(v string) {
	o.InstanceSize = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *MdbAvailableVersion) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MdbAvailableVersion) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *MdbAvailableVersion) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *MdbAvailableVersion) SetLinks(v []Link) {
	o.Links = &v
}

// GetVersion returns the Version field value if set, zero value otherwise
func (o *MdbAvailableVersion) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MdbAvailableVersion) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}

	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *MdbAvailableVersion) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *MdbAvailableVersion) SetVersion(v string) {
	o.Version = &v
}
