// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UpdateRequirePrivateNetworkingRequest Request body to toggle the private-networking requirement for an existing export bucket.
type UpdateRequirePrivateNetworkingRequest struct {
	// True to require private networking; false to disable it.
	RequirePrivateNetworking bool `json:"requirePrivateNetworking"`
}

// NewUpdateRequirePrivateNetworkingRequest instantiates a new UpdateRequirePrivateNetworkingRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateRequirePrivateNetworkingRequest(requirePrivateNetworking bool) *UpdateRequirePrivateNetworkingRequest {
	this := UpdateRequirePrivateNetworkingRequest{}
	this.RequirePrivateNetworking = requirePrivateNetworking
	return &this
}

// NewUpdateRequirePrivateNetworkingRequestWithDefaults instantiates a new UpdateRequirePrivateNetworkingRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateRequirePrivateNetworkingRequestWithDefaults() *UpdateRequirePrivateNetworkingRequest {
	this := UpdateRequirePrivateNetworkingRequest{}
	return &this
}

// GetRequirePrivateNetworking returns the RequirePrivateNetworking field value
func (o *UpdateRequirePrivateNetworkingRequest) GetRequirePrivateNetworking() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.RequirePrivateNetworking
}

// GetRequirePrivateNetworkingOk returns a tuple with the RequirePrivateNetworking field value
// and a boolean to check if the value has been set.
func (o *UpdateRequirePrivateNetworkingRequest) GetRequirePrivateNetworkingOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RequirePrivateNetworking, true
}

// SetRequirePrivateNetworking sets field value
func (o *UpdateRequirePrivateNetworkingRequest) SetRequirePrivateNetworking(v bool) {
	o.RequirePrivateNetworking = v
}
