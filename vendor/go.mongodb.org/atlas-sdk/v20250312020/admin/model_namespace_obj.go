// Code based on the AtlasAPI V2 OpenAPI file

package admin

// NamespaceObj Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
type NamespaceObj struct {
	// Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
	// Read only field.
	Namespace *string `json:"namespace,omitempty"`
	// Human-readable label that identifies the type of namespace.
	// Read only field.
	Type *string `json:"type,omitempty"`
}

// NewNamespaceObj instantiates a new NamespaceObj object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNamespaceObj() *NamespaceObj {
	this := NamespaceObj{}
	return &this
}

// NewNamespaceObjWithDefaults instantiates a new NamespaceObj object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNamespaceObjWithDefaults() *NamespaceObj {
	this := NamespaceObj{}
	return &this
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *NamespaceObj) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NamespaceObj) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *NamespaceObj) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *NamespaceObj) SetNamespace(v string) {
	o.Namespace = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *NamespaceObj) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NamespaceObj) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *NamespaceObj) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *NamespaceObj) SetType(v string) {
	o.Type = &v
}
