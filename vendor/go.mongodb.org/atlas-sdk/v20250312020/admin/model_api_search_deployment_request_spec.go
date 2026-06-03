// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiSearchDeploymentRequestSpec struct for ApiSearchDeploymentRequestSpec
type ApiSearchDeploymentRequestSpec struct {
	// Hardware specification for the Search Node instance sizes.
	InstanceSize string `json:"instanceSize"`
	// Number of Search Nodes in this region. Optional; falls back to the request-level default when omitted.
	NodeCount *int `json:"nodeCount,omitempty"`
}

// NewApiSearchDeploymentRequestSpec instantiates a new ApiSearchDeploymentRequestSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiSearchDeploymentRequestSpec(instanceSize string) *ApiSearchDeploymentRequestSpec {
	this := ApiSearchDeploymentRequestSpec{}
	this.InstanceSize = instanceSize
	return &this
}

// NewApiSearchDeploymentRequestSpecWithDefaults instantiates a new ApiSearchDeploymentRequestSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiSearchDeploymentRequestSpecWithDefaults() *ApiSearchDeploymentRequestSpec {
	this := ApiSearchDeploymentRequestSpec{}
	return &this
}

// GetInstanceSize returns the InstanceSize field value
func (o *ApiSearchDeploymentRequestSpec) GetInstanceSize() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.InstanceSize
}

// GetInstanceSizeOk returns a tuple with the InstanceSize field value
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentRequestSpec) GetInstanceSizeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InstanceSize, true
}

// SetInstanceSize sets field value
func (o *ApiSearchDeploymentRequestSpec) SetInstanceSize(v string) {
	o.InstanceSize = v
}

// GetNodeCount returns the NodeCount field value if set, zero value otherwise
func (o *ApiSearchDeploymentRequestSpec) GetNodeCount() int {
	if o == nil || IsNil(o.NodeCount) {
		var ret int
		return ret
	}
	return *o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentRequestSpec) GetNodeCountOk() (*int, bool) {
	if o == nil || IsNil(o.NodeCount) {
		return nil, false
	}

	return o.NodeCount, true
}

// HasNodeCount returns a boolean if a field has been set.
func (o *ApiSearchDeploymentRequestSpec) HasNodeCount() bool {
	if o != nil && !IsNil(o.NodeCount) {
		return true
	}

	return false
}

// SetNodeCount gets a reference to the given int and assigns it to the NodeCount field.
func (o *ApiSearchDeploymentRequestSpec) SetNodeCount(v int) {
	o.NodeCount = &v
}
