// Code based on the AtlasAPI V2 OpenAPI file

package admin

// TargetOrgRequest struct for TargetOrgRequest
type TargetOrgRequest struct {
	// IP address access list entries associated with the API key.
	AccessListIps *[]string `json:"accessListIps,omitempty"`
}

// NewTargetOrgRequest instantiates a new TargetOrgRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTargetOrgRequest() *TargetOrgRequest {
	this := TargetOrgRequest{}
	return &this
}

// NewTargetOrgRequestWithDefaults instantiates a new TargetOrgRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTargetOrgRequestWithDefaults() *TargetOrgRequest {
	this := TargetOrgRequest{}
	return &this
}

// GetAccessListIps returns the AccessListIps field value if set, zero value otherwise
func (o *TargetOrgRequest) GetAccessListIps() []string {
	if o == nil || IsNil(o.AccessListIps) {
		var ret []string
		return ret
	}
	return *o.AccessListIps
}

// GetAccessListIpsOk returns a tuple with the AccessListIps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TargetOrgRequest) GetAccessListIpsOk() (*[]string, bool) {
	if o == nil || IsNil(o.AccessListIps) {
		return nil, false
	}

	return o.AccessListIps, true
}

// HasAccessListIps returns a boolean if a field has been set.
func (o *TargetOrgRequest) HasAccessListIps() bool {
	if o != nil && !IsNil(o.AccessListIps) {
		return true
	}

	return false
}

// SetAccessListIps gets a reference to the given []string and assigns it to the AccessListIps field.
func (o *TargetOrgRequest) SetAccessListIps(v []string) {
	o.AccessListIps = &v
}
