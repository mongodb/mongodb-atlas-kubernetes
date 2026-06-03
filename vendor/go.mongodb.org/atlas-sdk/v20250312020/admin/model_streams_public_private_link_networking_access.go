// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsPublicPrivateLinkNetworkingAccess Information about networking access.
type StreamsPublicPrivateLinkNetworkingAccess struct {
	// The ID of the Private Link connection. Required for `PRIVATE_LINK` type. For GCP connections using Private Service Connect (PSC), this is the PSC connection ID.
	ConnectionId *string `json:"connectionId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Selected networking type. Either `PUBLIC` or `PRIVATE_LINK`. Defaults to `PUBLIC`. For AWS, Azure, and GCP connections, use `PRIVATE_LINK` for AWS PrivateLink, Azure Private Link, or GCP Private Service Connect (PSC) respectively.
	Type *string `json:"type,omitempty"`
}

// NewStreamsPublicPrivateLinkNetworkingAccess instantiates a new StreamsPublicPrivateLinkNetworkingAccess object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsPublicPrivateLinkNetworkingAccess() *StreamsPublicPrivateLinkNetworkingAccess {
	this := StreamsPublicPrivateLinkNetworkingAccess{}
	return &this
}

// NewStreamsPublicPrivateLinkNetworkingAccessWithDefaults instantiates a new StreamsPublicPrivateLinkNetworkingAccess object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsPublicPrivateLinkNetworkingAccessWithDefaults() *StreamsPublicPrivateLinkNetworkingAccess {
	this := StreamsPublicPrivateLinkNetworkingAccess{}
	return &this
}

// GetConnectionId returns the ConnectionId field value if set, zero value otherwise
func (o *StreamsPublicPrivateLinkNetworkingAccess) GetConnectionId() string {
	if o == nil || IsNil(o.ConnectionId) {
		var ret string
		return ret
	}
	return *o.ConnectionId
}

// GetConnectionIdOk returns a tuple with the ConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsPublicPrivateLinkNetworkingAccess) GetConnectionIdOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectionId) {
		return nil, false
	}

	return o.ConnectionId, true
}

// HasConnectionId returns a boolean if a field has been set.
func (o *StreamsPublicPrivateLinkNetworkingAccess) HasConnectionId() bool {
	if o != nil && !IsNil(o.ConnectionId) {
		return true
	}

	return false
}

// SetConnectionId gets a reference to the given string and assigns it to the ConnectionId field.
func (o *StreamsPublicPrivateLinkNetworkingAccess) SetConnectionId(v string) {
	o.ConnectionId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsPublicPrivateLinkNetworkingAccess) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsPublicPrivateLinkNetworkingAccess) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsPublicPrivateLinkNetworkingAccess) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsPublicPrivateLinkNetworkingAccess) SetLinks(v []Link) {
	o.Links = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *StreamsPublicPrivateLinkNetworkingAccess) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsPublicPrivateLinkNetworkingAccess) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *StreamsPublicPrivateLinkNetworkingAccess) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *StreamsPublicPrivateLinkNetworkingAccess) SetType(v string) {
	o.Type = &v
}
