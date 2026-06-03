// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupIPAddresses List of IP addresses in a project.
type GroupIPAddresses struct {
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud project.
	// Read only field.
	GroupId  *string       `json:"groupId,omitempty"`
	Services *GroupService `json:"services,omitempty"`
}

// NewGroupIPAddresses instantiates a new GroupIPAddresses object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupIPAddresses() *GroupIPAddresses {
	this := GroupIPAddresses{}
	return &this
}

// NewGroupIPAddressesWithDefaults instantiates a new GroupIPAddresses object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupIPAddressesWithDefaults() *GroupIPAddresses {
	this := GroupIPAddresses{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *GroupIPAddresses) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupIPAddresses) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *GroupIPAddresses) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *GroupIPAddresses) SetGroupId(v string) {
	o.GroupId = &v
}

// GetServices returns the Services field value if set, zero value otherwise
func (o *GroupIPAddresses) GetServices() GroupService {
	if o == nil || IsNil(o.Services) {
		var ret GroupService
		return ret
	}
	return *o.Services
}

// GetServicesOk returns a tuple with the Services field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupIPAddresses) GetServicesOk() (*GroupService, bool) {
	if o == nil || IsNil(o.Services) {
		return nil, false
	}

	return o.Services, true
}

// HasServices returns a boolean if a field has been set.
func (o *GroupIPAddresses) HasServices() bool {
	if o != nil && !IsNil(o.Services) {
		return true
	}

	return false
}

// SetServices gets a reference to the given GroupService and assigns it to the Services field.
func (o *GroupIPAddresses) SetServices(v GroupService) {
	o.Services = &v
}
