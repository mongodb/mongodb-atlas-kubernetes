// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsKafkaNetworkingAccess Information about networking access.
type StreamsKafkaNetworkingAccess struct {
	// Reserved. Will be used by `PRIVATE_LINK` connection type.
	ConnectionId *string `json:"connectionId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Reserved. Will be used by `PRIVATE_LINK` connection type.
	Name *string `json:"name,omitempty"`
	// Reserved. Will be used by `TRANSIT_GATEWAY` connection type.
	TgwRouteId *string `json:"tgwRouteId,omitempty"`
	// Selected networking type. Either `PUBLIC`, `VPC`, `PRIVATE_LINK`, or `TRANSIT_GATEWAY`. Defaults to `PUBLIC`. For VPC, ensure that VPC peering exists and connectivity has been established between Atlas VPC and the VPC where Kafka cluster is hosted for the connection to function properly. `TRANSIT_GATEWAY` support is coming soon.
	Type *string `json:"type,omitempty"`
}

// NewStreamsKafkaNetworkingAccess instantiates a new StreamsKafkaNetworkingAccess object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsKafkaNetworkingAccess() *StreamsKafkaNetworkingAccess {
	this := StreamsKafkaNetworkingAccess{}
	return &this
}

// NewStreamsKafkaNetworkingAccessWithDefaults instantiates a new StreamsKafkaNetworkingAccess object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsKafkaNetworkingAccessWithDefaults() *StreamsKafkaNetworkingAccess {
	this := StreamsKafkaNetworkingAccess{}
	return &this
}

// GetConnectionId returns the ConnectionId field value if set, zero value otherwise
func (o *StreamsKafkaNetworkingAccess) GetConnectionId() string {
	if o == nil || IsNil(o.ConnectionId) {
		var ret string
		return ret
	}
	return *o.ConnectionId
}

// GetConnectionIdOk returns a tuple with the ConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaNetworkingAccess) GetConnectionIdOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectionId) {
		return nil, false
	}

	return o.ConnectionId, true
}

// HasConnectionId returns a boolean if a field has been set.
func (o *StreamsKafkaNetworkingAccess) HasConnectionId() bool {
	if o != nil && !IsNil(o.ConnectionId) {
		return true
	}

	return false
}

// SetConnectionId gets a reference to the given string and assigns it to the ConnectionId field.
func (o *StreamsKafkaNetworkingAccess) SetConnectionId(v string) {
	o.ConnectionId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsKafkaNetworkingAccess) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaNetworkingAccess) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsKafkaNetworkingAccess) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsKafkaNetworkingAccess) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *StreamsKafkaNetworkingAccess) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaNetworkingAccess) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *StreamsKafkaNetworkingAccess) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *StreamsKafkaNetworkingAccess) SetName(v string) {
	o.Name = &v
}

// GetTgwRouteId returns the TgwRouteId field value if set, zero value otherwise
func (o *StreamsKafkaNetworkingAccess) GetTgwRouteId() string {
	if o == nil || IsNil(o.TgwRouteId) {
		var ret string
		return ret
	}
	return *o.TgwRouteId
}

// GetTgwRouteIdOk returns a tuple with the TgwRouteId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaNetworkingAccess) GetTgwRouteIdOk() (*string, bool) {
	if o == nil || IsNil(o.TgwRouteId) {
		return nil, false
	}

	return o.TgwRouteId, true
}

// HasTgwRouteId returns a boolean if a field has been set.
func (o *StreamsKafkaNetworkingAccess) HasTgwRouteId() bool {
	if o != nil && !IsNil(o.TgwRouteId) {
		return true
	}

	return false
}

// SetTgwRouteId gets a reference to the given string and assigns it to the TgwRouteId field.
func (o *StreamsKafkaNetworkingAccess) SetTgwRouteId(v string) {
	o.TgwRouteId = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *StreamsKafkaNetworkingAccess) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaNetworkingAccess) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *StreamsKafkaNetworkingAccess) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *StreamsKafkaNetworkingAccess) SetType(v string) {
	o.Type = &v
}
