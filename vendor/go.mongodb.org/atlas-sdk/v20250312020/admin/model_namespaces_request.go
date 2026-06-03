// Code based on the AtlasAPI V2 OpenAPI file

package admin

// NamespacesRequest struct for NamespacesRequest
type NamespacesRequest struct {
	// List of namespace strings (combination of database and collection) on the specified host or cluster.
	// Write only field.
	Namespaces *[]string `json:"namespaces,omitempty"`
}

// NewNamespacesRequest instantiates a new NamespacesRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNamespacesRequest() *NamespacesRequest {
	this := NamespacesRequest{}
	return &this
}

// NewNamespacesRequestWithDefaults instantiates a new NamespacesRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNamespacesRequestWithDefaults() *NamespacesRequest {
	this := NamespacesRequest{}
	return &this
}

// GetNamespaces returns the Namespaces field value if set, zero value otherwise
func (o *NamespacesRequest) GetNamespaces() []string {
	if o == nil || IsNil(o.Namespaces) {
		var ret []string
		return ret
	}
	return *o.Namespaces
}

// GetNamespacesOk returns a tuple with the Namespaces field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NamespacesRequest) GetNamespacesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Namespaces) {
		return nil, false
	}

	return o.Namespaces, true
}

// HasNamespaces returns a boolean if a field has been set.
func (o *NamespacesRequest) HasNamespaces() bool {
	if o != nil && !IsNil(o.Namespaces) {
		return true
	}

	return false
}

// SetNamespaces gets a reference to the given []string and assigns it to the Namespaces field.
func (o *NamespacesRequest) SetNamespaces(v []string) {
	o.Namespaces = &v
}
