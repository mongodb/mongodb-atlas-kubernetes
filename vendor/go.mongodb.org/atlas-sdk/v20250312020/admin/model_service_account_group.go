// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServiceAccountGroup struct for ServiceAccountGroup
type ServiceAccountGroup struct {
	// Unique 24-hexadecimal digit string that identifies your project. **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
}

// NewServiceAccountGroup instantiates a new ServiceAccountGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServiceAccountGroup() *ServiceAccountGroup {
	this := ServiceAccountGroup{}
	return &this
}

// NewServiceAccountGroupWithDefaults instantiates a new ServiceAccountGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServiceAccountGroupWithDefaults() *ServiceAccountGroup {
	this := ServiceAccountGroup{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ServiceAccountGroup) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountGroup) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ServiceAccountGroup) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ServiceAccountGroup) SetGroupId(v string) {
	o.GroupId = &v
}
