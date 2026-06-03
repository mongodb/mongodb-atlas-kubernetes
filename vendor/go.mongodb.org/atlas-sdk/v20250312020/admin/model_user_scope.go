// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UserScope Range of resources available to this database user.
type UserScope struct {
	// Human-readable label that identifies the cluster or MongoDB Atlas Data Lake that this database user can access.
	Name string `json:"name"`
	// Category of resource that this database user can access.
	Type string `json:"type"`
}

// NewUserScope instantiates a new UserScope object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserScope(name string, type_ string) *UserScope {
	this := UserScope{}
	this.Name = name
	this.Type = type_
	return &this
}

// NewUserScopeWithDefaults instantiates a new UserScope object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserScopeWithDefaults() *UserScope {
	this := UserScope{}
	return &this
}

// GetName returns the Name field value
func (o *UserScope) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *UserScope) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *UserScope) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *UserScope) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *UserScope) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *UserScope) SetType(v string) {
	o.Type = v
}
