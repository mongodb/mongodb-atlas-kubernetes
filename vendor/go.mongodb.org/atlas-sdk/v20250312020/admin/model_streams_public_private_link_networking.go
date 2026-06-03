// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsPublicPrivateLinkNetworking Networking configuration for connections that support `PUBLIC` and `PRIVATE_LINK` access types. For GCP connections, use `PRIVATE_LINK` for GCP Private Service Connect (PSC).
type StreamsPublicPrivateLinkNetworking struct {
	Access *StreamsPublicPrivateLinkNetworkingAccess `json:"access,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewStreamsPublicPrivateLinkNetworking instantiates a new StreamsPublicPrivateLinkNetworking object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsPublicPrivateLinkNetworking() *StreamsPublicPrivateLinkNetworking {
	this := StreamsPublicPrivateLinkNetworking{}
	return &this
}

// NewStreamsPublicPrivateLinkNetworkingWithDefaults instantiates a new StreamsPublicPrivateLinkNetworking object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsPublicPrivateLinkNetworkingWithDefaults() *StreamsPublicPrivateLinkNetworking {
	this := StreamsPublicPrivateLinkNetworking{}
	return &this
}

// GetAccess returns the Access field value if set, zero value otherwise
func (o *StreamsPublicPrivateLinkNetworking) GetAccess() StreamsPublicPrivateLinkNetworkingAccess {
	if o == nil || IsNil(o.Access) {
		var ret StreamsPublicPrivateLinkNetworkingAccess
		return ret
	}
	return *o.Access
}

// GetAccessOk returns a tuple with the Access field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsPublicPrivateLinkNetworking) GetAccessOk() (*StreamsPublicPrivateLinkNetworkingAccess, bool) {
	if o == nil || IsNil(o.Access) {
		return nil, false
	}

	return o.Access, true
}

// HasAccess returns a boolean if a field has been set.
func (o *StreamsPublicPrivateLinkNetworking) HasAccess() bool {
	if o != nil && !IsNil(o.Access) {
		return true
	}

	return false
}

// SetAccess gets a reference to the given StreamsPublicPrivateLinkNetworkingAccess and assigns it to the Access field.
func (o *StreamsPublicPrivateLinkNetworking) SetAccess(v StreamsPublicPrivateLinkNetworkingAccess) {
	o.Access = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsPublicPrivateLinkNetworking) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsPublicPrivateLinkNetworking) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsPublicPrivateLinkNetworking) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsPublicPrivateLinkNetworking) SetLinks(v []Link) {
	o.Links = &v
}
