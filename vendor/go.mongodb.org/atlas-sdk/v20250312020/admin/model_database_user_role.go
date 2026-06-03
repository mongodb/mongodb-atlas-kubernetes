// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DatabaseUserRole Range of resources available to this database user.
type DatabaseUserRole struct {
	// Collection on which this role applies.
	CollectionName *string `json:"collectionName,omitempty"`
	// Database to which the user is granted access privileges.
	DatabaseName string `json:"databaseName"`
	// Human-readable label that identifies a group of privileges assigned to a database user. This value can either be a built-in role or a custom role.
	RoleName string `json:"roleName"`
}

// NewDatabaseUserRole instantiates a new DatabaseUserRole object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatabaseUserRole(databaseName string, roleName string) *DatabaseUserRole {
	this := DatabaseUserRole{}
	this.DatabaseName = databaseName
	this.RoleName = roleName
	return &this
}

// NewDatabaseUserRoleWithDefaults instantiates a new DatabaseUserRole object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatabaseUserRoleWithDefaults() *DatabaseUserRole {
	this := DatabaseUserRole{}
	return &this
}

// GetCollectionName returns the CollectionName field value if set, zero value otherwise
func (o *DatabaseUserRole) GetCollectionName() string {
	if o == nil || IsNil(o.CollectionName) {
		var ret string
		return ret
	}
	return *o.CollectionName
}

// GetCollectionNameOk returns a tuple with the CollectionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseUserRole) GetCollectionNameOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionName) {
		return nil, false
	}

	return o.CollectionName, true
}

// HasCollectionName returns a boolean if a field has been set.
func (o *DatabaseUserRole) HasCollectionName() bool {
	if o != nil && !IsNil(o.CollectionName) {
		return true
	}

	return false
}

// SetCollectionName gets a reference to the given string and assigns it to the CollectionName field.
func (o *DatabaseUserRole) SetCollectionName(v string) {
	o.CollectionName = &v
}

// GetDatabaseName returns the DatabaseName field value
func (o *DatabaseUserRole) GetDatabaseName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DatabaseName
}

// GetDatabaseNameOk returns a tuple with the DatabaseName field value
// and a boolean to check if the value has been set.
func (o *DatabaseUserRole) GetDatabaseNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DatabaseName, true
}

// SetDatabaseName sets field value
func (o *DatabaseUserRole) SetDatabaseName(v string) {
	o.DatabaseName = v
}

// GetRoleName returns the RoleName field value
func (o *DatabaseUserRole) GetRoleName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RoleName
}

// GetRoleNameOk returns a tuple with the RoleName field value
// and a boolean to check if the value has been set.
func (o *DatabaseUserRole) GetRoleNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RoleName, true
}

// SetRoleName sets field value
func (o *DatabaseUserRole) SetRoleName(v string) {
	o.RoleName = v
}
