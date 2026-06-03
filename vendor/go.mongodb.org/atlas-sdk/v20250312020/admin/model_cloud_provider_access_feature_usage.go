// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudProviderAccessFeatureUsage MongoDB Cloud features associated with this Amazon Web Services (AWS) Identity and Access Management (IAM) role.
type CloudProviderAccessFeatureUsage struct {
	// Human-readable label that describes one MongoDB Cloud feature linked to this Amazon Web Services (AWS) Identity and Access Management (IAM) role.
	// Read only field.
	FeatureType *string                                                     `json:"featureType,omitempty"`
	FeatureId   *CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId `json:"featureId,omitempty"`
}

// NewCloudProviderAccessFeatureUsage instantiates a new CloudProviderAccessFeatureUsage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderAccessFeatureUsage() *CloudProviderAccessFeatureUsage {
	this := CloudProviderAccessFeatureUsage{}
	return &this
}

// NewCloudProviderAccessFeatureUsageWithDefaults instantiates a new CloudProviderAccessFeatureUsage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderAccessFeatureUsageWithDefaults() *CloudProviderAccessFeatureUsage {
	this := CloudProviderAccessFeatureUsage{}
	return &this
}

// GetFeatureType returns the FeatureType field value if set, zero value otherwise
func (o *CloudProviderAccessFeatureUsage) GetFeatureType() string {
	if o == nil || IsNil(o.FeatureType) {
		var ret string
		return ret
	}
	return *o.FeatureType
}

// GetFeatureTypeOk returns a tuple with the FeatureType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessFeatureUsage) GetFeatureTypeOk() (*string, bool) {
	if o == nil || IsNil(o.FeatureType) {
		return nil, false
	}

	return o.FeatureType, true
}

// HasFeatureType returns a boolean if a field has been set.
func (o *CloudProviderAccessFeatureUsage) HasFeatureType() bool {
	if o != nil && !IsNil(o.FeatureType) {
		return true
	}

	return false
}

// SetFeatureType gets a reference to the given string and assigns it to the FeatureType field.
func (o *CloudProviderAccessFeatureUsage) SetFeatureType(v string) {
	o.FeatureType = &v
}

// GetFeatureId returns the FeatureId field value if set, zero value otherwise
func (o *CloudProviderAccessFeatureUsage) GetFeatureId() CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId {
	if o == nil || IsNil(o.FeatureId) {
		var ret CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId
		return ret
	}
	return *o.FeatureId
}

// GetFeatureIdOk returns a tuple with the FeatureId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderAccessFeatureUsage) GetFeatureIdOk() (*CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId, bool) {
	if o == nil || IsNil(o.FeatureId) {
		return nil, false
	}

	return o.FeatureId, true
}

// HasFeatureId returns a boolean if a field has been set.
func (o *CloudProviderAccessFeatureUsage) HasFeatureId() bool {
	if o != nil && !IsNil(o.FeatureId) {
		return true
	}

	return false
}

// SetFeatureId gets a reference to the given CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId and assigns it to the FeatureId field.
func (o *CloudProviderAccessFeatureUsage) SetFeatureId(v CloudProviderAccessFeatureUsagePushBasedLogExportFeatureId) {
	o.FeatureId = &v
}
