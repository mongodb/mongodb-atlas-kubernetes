// Code based on the AtlasAPI V2 OpenAPI file

package admin

// RateLimitEndpointSetEndpoint struct for RateLimitEndpointSetEndpoint
type RateLimitEndpointSetEndpoint struct {
	// The HTTP method of the endpoint.
	Method *string `json:"method,omitempty"`
	// The URL path of the endpoint.
	Path *string `json:"path,omitempty"`
}

// NewRateLimitEndpointSetEndpoint instantiates a new RateLimitEndpointSetEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRateLimitEndpointSetEndpoint() *RateLimitEndpointSetEndpoint {
	this := RateLimitEndpointSetEndpoint{}
	return &this
}

// NewRateLimitEndpointSetEndpointWithDefaults instantiates a new RateLimitEndpointSetEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRateLimitEndpointSetEndpointWithDefaults() *RateLimitEndpointSetEndpoint {
	this := RateLimitEndpointSetEndpoint{}
	return &this
}

// GetMethod returns the Method field value if set, zero value otherwise
func (o *RateLimitEndpointSetEndpoint) GetMethod() string {
	if o == nil || IsNil(o.Method) {
		var ret string
		return ret
	}
	return *o.Method
}

// GetMethodOk returns a tuple with the Method field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetEndpoint) GetMethodOk() (*string, bool) {
	if o == nil || IsNil(o.Method) {
		return nil, false
	}

	return o.Method, true
}

// HasMethod returns a boolean if a field has been set.
func (o *RateLimitEndpointSetEndpoint) HasMethod() bool {
	if o != nil && !IsNil(o.Method) {
		return true
	}

	return false
}

// SetMethod gets a reference to the given string and assigns it to the Method field.
func (o *RateLimitEndpointSetEndpoint) SetMethod(v string) {
	o.Method = &v
}

// GetPath returns the Path field value if set, zero value otherwise
func (o *RateLimitEndpointSetEndpoint) GetPath() string {
	if o == nil || IsNil(o.Path) {
		var ret string
		return ret
	}
	return *o.Path
}

// GetPathOk returns a tuple with the Path field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetEndpoint) GetPathOk() (*string, bool) {
	if o == nil || IsNil(o.Path) {
		return nil, false
	}

	return o.Path, true
}

// HasPath returns a boolean if a field has been set.
func (o *RateLimitEndpointSetEndpoint) HasPath() bool {
	if o != nil && !IsNil(o.Path) {
		return true
	}

	return false
}

// SetPath gets a reference to the given string and assigns it to the Path field.
func (o *RateLimitEndpointSetEndpoint) SetPath(v string) {
	o.Path = &v
}
