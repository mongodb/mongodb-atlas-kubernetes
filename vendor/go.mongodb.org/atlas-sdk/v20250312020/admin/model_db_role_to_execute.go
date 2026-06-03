// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DBRoleToExecute The name of a Built in or Custom DB Role to connect to an Atlas Cluster.
type DBRoleToExecute struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// The name of the role to use. Can be a built in role or a custom role.
	Role *string `json:"role,omitempty"`
	// Type of the DB role. Can be either Built In or Custom.
	Type *string `json:"type,omitempty"`
}

// NewDBRoleToExecute instantiates a new DBRoleToExecute object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDBRoleToExecute() *DBRoleToExecute {
	this := DBRoleToExecute{}
	return &this
}

// NewDBRoleToExecuteWithDefaults instantiates a new DBRoleToExecute object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDBRoleToExecuteWithDefaults() *DBRoleToExecute {
	this := DBRoleToExecute{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DBRoleToExecute) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DBRoleToExecute) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DBRoleToExecute) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DBRoleToExecute) SetLinks(v []Link) {
	o.Links = &v
}

// GetRole returns the Role field value if set, zero value otherwise
func (o *DBRoleToExecute) GetRole() string {
	if o == nil || IsNil(o.Role) {
		var ret string
		return ret
	}
	return *o.Role
}

// GetRoleOk returns a tuple with the Role field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DBRoleToExecute) GetRoleOk() (*string, bool) {
	if o == nil || IsNil(o.Role) {
		return nil, false
	}

	return o.Role, true
}

// HasRole returns a boolean if a field has been set.
func (o *DBRoleToExecute) HasRole() bool {
	if o != nil && !IsNil(o.Role) {
		return true
	}

	return false
}

// SetRole gets a reference to the given string and assigns it to the Role field.
func (o *DBRoleToExecute) SetRole(v string) {
	o.Role = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *DBRoleToExecute) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DBRoleToExecute) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *DBRoleToExecute) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *DBRoleToExecute) SetType(v string) {
	o.Type = &v
}
