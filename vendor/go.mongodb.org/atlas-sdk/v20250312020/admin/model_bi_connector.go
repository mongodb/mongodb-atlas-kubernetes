// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BiConnector Settings needed to configure the MongoDB Connector for Business Intelligence for this cluster.
type BiConnector struct {
	// Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.
	Enabled *bool `json:"enabled,omitempty"`
	// Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud. The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences. Defaults to `ANALYTICS` node, or `SECONDARY` if there are no `ANALYTICS` nodes.
	ReadPreference *string `json:"readPreference,omitempty"`
}

// NewBiConnector instantiates a new BiConnector object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBiConnector() *BiConnector {
	this := BiConnector{}
	return &this
}

// NewBiConnectorWithDefaults instantiates a new BiConnector object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBiConnectorWithDefaults() *BiConnector {
	this := BiConnector{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *BiConnector) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BiConnector) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *BiConnector) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *BiConnector) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetReadPreference returns the ReadPreference field value if set, zero value otherwise
func (o *BiConnector) GetReadPreference() string {
	if o == nil || IsNil(o.ReadPreference) {
		var ret string
		return ret
	}
	return *o.ReadPreference
}

// GetReadPreferenceOk returns a tuple with the ReadPreference field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BiConnector) GetReadPreferenceOk() (*string, bool) {
	if o == nil || IsNil(o.ReadPreference) {
		return nil, false
	}

	return o.ReadPreference, true
}

// HasReadPreference returns a boolean if a field has been set.
func (o *BiConnector) HasReadPreference() bool {
	if o != nil && !IsNil(o.ReadPreference) {
		return true
	}

	return false
}

// SetReadPreference gets a reference to the given string and assigns it to the ReadPreference field.
func (o *BiConnector) SetReadPreference(v string) {
	o.ReadPreference = &v
}
