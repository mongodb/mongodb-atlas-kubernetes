// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ZoneMapping Human-readable label that identifies the subset of a global cluster.
type ZoneMapping struct {
	// Code that represents a location that maps to a zone in your global cluster. MongoDB Cloud represents this location with a ISO 3166-2 location and subdivision codes when possible.
	Location string `json:"location"`
	// Human-readable label that identifies the zone in your global cluster. This zone maps to a location code.
	Zone string `json:"zone"`
}

// NewZoneMapping instantiates a new ZoneMapping object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewZoneMapping(location string, zone string) *ZoneMapping {
	this := ZoneMapping{}
	this.Location = location
	this.Zone = zone
	return &this
}

// NewZoneMappingWithDefaults instantiates a new ZoneMapping object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewZoneMappingWithDefaults() *ZoneMapping {
	this := ZoneMapping{}
	return &this
}

// GetLocation returns the Location field value
func (o *ZoneMapping) GetLocation() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Location
}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
func (o *ZoneMapping) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Location, true
}

// SetLocation sets field value
func (o *ZoneMapping) SetLocation(v string) {
	o.Location = v
}

// GetZone returns the Zone field value
func (o *ZoneMapping) GetZone() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Zone
}

// GetZoneOk returns a tuple with the Zone field value
// and a boolean to check if the value has been set.
func (o *ZoneMapping) GetZoneOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Zone, true
}

// SetZone sets field value
func (o *ZoneMapping) SetZone(v string) {
	o.Zone = v
}
