// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgServiceAccountUpdateRequest struct for OrgServiceAccountUpdateRequest
type OrgServiceAccountUpdateRequest struct {
	// Human readable description for the Service Account.
	Description *string `json:"description,omitempty"`
	// Human-readable name for the Service Account. The name is modifiable and does not have to be unique.
	Name *string `json:"name,omitempty"`
	// A list of organization-level roles for the Service Account.
	Roles *[]string `json:"roles,omitempty"`
}

// NewOrgServiceAccountUpdateRequest instantiates a new OrgServiceAccountUpdateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgServiceAccountUpdateRequest() *OrgServiceAccountUpdateRequest {
	this := OrgServiceAccountUpdateRequest{}
	return &this
}

// NewOrgServiceAccountUpdateRequestWithDefaults instantiates a new OrgServiceAccountUpdateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgServiceAccountUpdateRequestWithDefaults() *OrgServiceAccountUpdateRequest {
	this := OrgServiceAccountUpdateRequest{}
	return &this
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *OrgServiceAccountUpdateRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccountUpdateRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *OrgServiceAccountUpdateRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *OrgServiceAccountUpdateRequest) SetDescription(v string) {
	o.Description = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *OrgServiceAccountUpdateRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccountUpdateRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *OrgServiceAccountUpdateRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *OrgServiceAccountUpdateRequest) SetName(v string) {
	o.Name = &v
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *OrgServiceAccountUpdateRequest) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccountUpdateRequest) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *OrgServiceAccountUpdateRequest) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *OrgServiceAccountUpdateRequest) SetRoles(v []string) {
	o.Roles = &v
}
