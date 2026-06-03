// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CustomZoneMappings struct for CustomZoneMappings
type CustomZoneMappings struct {
	// List that contains comma-separated key value pairs to map zones to geographic regions. These pairs map an ISO 3166-1a2 location code, with an ISO 3166-2 subdivision code when possible, to the human-readable label for the desired custom zone. MongoDB Cloud maps the ISO 3166-1a2 code to the nearest geographical zone by default. Include this parameter to override the default mappings.  This parameter returns an empty object if no custom zones exist.
	CustomZoneMappings []ZoneMapping `json:"customZoneMappings"`
}

// NewCustomZoneMappings instantiates a new CustomZoneMappings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCustomZoneMappings(customZoneMappings []ZoneMapping) *CustomZoneMappings {
	this := CustomZoneMappings{}
	this.CustomZoneMappings = customZoneMappings
	return &this
}

// NewCustomZoneMappingsWithDefaults instantiates a new CustomZoneMappings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCustomZoneMappingsWithDefaults() *CustomZoneMappings {
	this := CustomZoneMappings{}
	return &this
}

// GetCustomZoneMappings returns the CustomZoneMappings field value
func (o *CustomZoneMappings) GetCustomZoneMappings() []ZoneMapping {
	if o == nil {
		var ret []ZoneMapping
		return ret
	}

	return o.CustomZoneMappings
}

// GetCustomZoneMappingsOk returns a tuple with the CustomZoneMappings field value
// and a boolean to check if the value has been set.
func (o *CustomZoneMappings) GetCustomZoneMappingsOk() (*[]ZoneMapping, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomZoneMappings, true
}

// SetCustomZoneMappings sets field value
func (o *CustomZoneMappings) SetCustomZoneMappings(v []ZoneMapping) {
	o.CustomZoneMappings = v
}
