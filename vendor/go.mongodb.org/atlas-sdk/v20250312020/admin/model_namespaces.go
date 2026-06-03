// Code based on the AtlasAPI V2 OpenAPI file

package admin

// Namespaces struct for Namespaces
type Namespaces struct {
	// List that contains each combination of database, collection, and type on the specified host.
	// Read only field.
	Namespaces *[]NamespaceObj `json:"namespaces,omitempty"`
}

// NewNamespaces instantiates a new Namespaces object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNamespaces() *Namespaces {
	this := Namespaces{}
	return &this
}

// NewNamespacesWithDefaults instantiates a new Namespaces object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNamespacesWithDefaults() *Namespaces {
	this := Namespaces{}
	return &this
}

// GetNamespaces returns the Namespaces field value if set, zero value otherwise
func (o *Namespaces) GetNamespaces() []NamespaceObj {
	if o == nil || IsNil(o.Namespaces) {
		var ret []NamespaceObj
		return ret
	}
	return *o.Namespaces
}

// GetNamespacesOk returns a tuple with the Namespaces field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Namespaces) GetNamespacesOk() (*[]NamespaceObj, bool) {
	if o == nil || IsNil(o.Namespaces) {
		return nil, false
	}

	return o.Namespaces, true
}

// HasNamespaces returns a boolean if a field has been set.
func (o *Namespaces) HasNamespaces() bool {
	if o != nil && !IsNil(o.Namespaces) {
		return true
	}

	return false
}

// SetNamespaces gets a reference to the given []NamespaceObj and assigns it to the Namespaces field.
func (o *Namespaces) SetNamespaces(v []NamespaceObj) {
	o.Namespaces = &v
}
