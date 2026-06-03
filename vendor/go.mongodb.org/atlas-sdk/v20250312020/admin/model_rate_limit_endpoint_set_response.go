// Code based on the AtlasAPI V2 OpenAPI file

package admin

// RateLimitEndpointSetResponse struct for RateLimitEndpointSetResponse
type RateLimitEndpointSetResponse struct {
	Capacity *RateLimitEndpointSetCapacity `json:"capacity,omitempty"`
	// The ID of the endpoint set.
	EndpointSetId *string `json:"endpointSetId,omitempty"`
	// The endpoint set name.
	EndpointSetName *string `json:"endpointSetName,omitempty"`
	// A list of endpoints associated with the specified endpoint set.
	Endpoints             *[]RateLimitEndpointSetEndpoint            `json:"endpoints,omitempty"`
	RefillDurationSeconds *RateLimitEndpointSetRefillDurationSeconds `json:"refillDurationSeconds,omitempty"`
	RefillRate            *RateLimitEndpointSetRefillRate            `json:"refillRate,omitempty"`
	// The scope of the endpoint set.
	Scope *string `json:"scope,omitempty"`
}

// NewRateLimitEndpointSetResponse instantiates a new RateLimitEndpointSetResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRateLimitEndpointSetResponse() *RateLimitEndpointSetResponse {
	this := RateLimitEndpointSetResponse{}
	return &this
}

// NewRateLimitEndpointSetResponseWithDefaults instantiates a new RateLimitEndpointSetResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRateLimitEndpointSetResponseWithDefaults() *RateLimitEndpointSetResponse {
	this := RateLimitEndpointSetResponse{}
	return &this
}

// GetCapacity returns the Capacity field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetCapacity() RateLimitEndpointSetCapacity {
	if o == nil || IsNil(o.Capacity) {
		var ret RateLimitEndpointSetCapacity
		return ret
	}
	return *o.Capacity
}

// GetCapacityOk returns a tuple with the Capacity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetCapacityOk() (*RateLimitEndpointSetCapacity, bool) {
	if o == nil || IsNil(o.Capacity) {
		return nil, false
	}

	return o.Capacity, true
}

// HasCapacity returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasCapacity() bool {
	if o != nil && !IsNil(o.Capacity) {
		return true
	}

	return false
}

// SetCapacity gets a reference to the given RateLimitEndpointSetCapacity and assigns it to the Capacity field.
func (o *RateLimitEndpointSetResponse) SetCapacity(v RateLimitEndpointSetCapacity) {
	o.Capacity = &v
}

// GetEndpointSetId returns the EndpointSetId field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetEndpointSetId() string {
	if o == nil || IsNil(o.EndpointSetId) {
		var ret string
		return ret
	}
	return *o.EndpointSetId
}

// GetEndpointSetIdOk returns a tuple with the EndpointSetId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetEndpointSetIdOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetId) {
		return nil, false
	}

	return o.EndpointSetId, true
}

// HasEndpointSetId returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasEndpointSetId() bool {
	if o != nil && !IsNil(o.EndpointSetId) {
		return true
	}

	return false
}

// SetEndpointSetId gets a reference to the given string and assigns it to the EndpointSetId field.
func (o *RateLimitEndpointSetResponse) SetEndpointSetId(v string) {
	o.EndpointSetId = &v
}

// GetEndpointSetName returns the EndpointSetName field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetEndpointSetName() string {
	if o == nil || IsNil(o.EndpointSetName) {
		var ret string
		return ret
	}
	return *o.EndpointSetName
}

// GetEndpointSetNameOk returns a tuple with the EndpointSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetEndpointSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetName) {
		return nil, false
	}

	return o.EndpointSetName, true
}

// HasEndpointSetName returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasEndpointSetName() bool {
	if o != nil && !IsNil(o.EndpointSetName) {
		return true
	}

	return false
}

// SetEndpointSetName gets a reference to the given string and assigns it to the EndpointSetName field.
func (o *RateLimitEndpointSetResponse) SetEndpointSetName(v string) {
	o.EndpointSetName = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetEndpoints() []RateLimitEndpointSetEndpoint {
	if o == nil || IsNil(o.Endpoints) {
		var ret []RateLimitEndpointSetEndpoint
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetEndpointsOk() (*[]RateLimitEndpointSetEndpoint, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}

	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []RateLimitEndpointSetEndpoint and assigns it to the Endpoints field.
func (o *RateLimitEndpointSetResponse) SetEndpoints(v []RateLimitEndpointSetEndpoint) {
	o.Endpoints = &v
}

// GetRefillDurationSeconds returns the RefillDurationSeconds field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetRefillDurationSeconds() RateLimitEndpointSetRefillDurationSeconds {
	if o == nil || IsNil(o.RefillDurationSeconds) {
		var ret RateLimitEndpointSetRefillDurationSeconds
		return ret
	}
	return *o.RefillDurationSeconds
}

// GetRefillDurationSecondsOk returns a tuple with the RefillDurationSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetRefillDurationSecondsOk() (*RateLimitEndpointSetRefillDurationSeconds, bool) {
	if o == nil || IsNil(o.RefillDurationSeconds) {
		return nil, false
	}

	return o.RefillDurationSeconds, true
}

// HasRefillDurationSeconds returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasRefillDurationSeconds() bool {
	if o != nil && !IsNil(o.RefillDurationSeconds) {
		return true
	}

	return false
}

// SetRefillDurationSeconds gets a reference to the given RateLimitEndpointSetRefillDurationSeconds and assigns it to the RefillDurationSeconds field.
func (o *RateLimitEndpointSetResponse) SetRefillDurationSeconds(v RateLimitEndpointSetRefillDurationSeconds) {
	o.RefillDurationSeconds = &v
}

// GetRefillRate returns the RefillRate field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetRefillRate() RateLimitEndpointSetRefillRate {
	if o == nil || IsNil(o.RefillRate) {
		var ret RateLimitEndpointSetRefillRate
		return ret
	}
	return *o.RefillRate
}

// GetRefillRateOk returns a tuple with the RefillRate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetRefillRateOk() (*RateLimitEndpointSetRefillRate, bool) {
	if o == nil || IsNil(o.RefillRate) {
		return nil, false
	}

	return o.RefillRate, true
}

// HasRefillRate returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasRefillRate() bool {
	if o != nil && !IsNil(o.RefillRate) {
		return true
	}

	return false
}

// SetRefillRate gets a reference to the given RateLimitEndpointSetRefillRate and assigns it to the RefillRate field.
func (o *RateLimitEndpointSetResponse) SetRefillRate(v RateLimitEndpointSetRefillRate) {
	o.RefillRate = &v
}

// GetScope returns the Scope field value if set, zero value otherwise
func (o *RateLimitEndpointSetResponse) GetScope() string {
	if o == nil || IsNil(o.Scope) {
		var ret string
		return ret
	}
	return *o.Scope
}

// GetScopeOk returns a tuple with the Scope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RateLimitEndpointSetResponse) GetScopeOk() (*string, bool) {
	if o == nil || IsNil(o.Scope) {
		return nil, false
	}

	return o.Scope, true
}

// HasScope returns a boolean if a field has been set.
func (o *RateLimitEndpointSetResponse) HasScope() bool {
	if o != nil && !IsNil(o.Scope) {
		return true
	}

	return false
}

// SetScope gets a reference to the given string and assigns it to the Scope field.
func (o *RateLimitEndpointSetResponse) SetScope(v string) {
	o.Scope = &v
}
