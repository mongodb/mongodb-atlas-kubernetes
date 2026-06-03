// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiSearchDeploymentSpec struct for ApiSearchDeploymentSpec
type ApiSearchDeploymentSpec struct {
	// Hardware specification for the Search Node instance sizes.
	InstanceSize string `json:"instanceSize"`
	// Number of Search Nodes in the cluster.
	NodeCount int `json:"nodeCount"`
}

// NewApiSearchDeploymentSpec instantiates a new ApiSearchDeploymentSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiSearchDeploymentSpec(instanceSize string, nodeCount int) *ApiSearchDeploymentSpec {
	this := ApiSearchDeploymentSpec{}
	this.InstanceSize = instanceSize
	this.NodeCount = nodeCount
	return &this
}

// NewApiSearchDeploymentSpecWithDefaults instantiates a new ApiSearchDeploymentSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiSearchDeploymentSpecWithDefaults() *ApiSearchDeploymentSpec {
	this := ApiSearchDeploymentSpec{}
	return &this
}

// GetInstanceSize returns the InstanceSize field value
func (o *ApiSearchDeploymentSpec) GetInstanceSize() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.InstanceSize
}

// GetInstanceSizeOk returns a tuple with the InstanceSize field value
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentSpec) GetInstanceSizeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.InstanceSize, true
}

// SetInstanceSize sets field value
func (o *ApiSearchDeploymentSpec) SetInstanceSize(v string) {
	o.InstanceSize = v
}

// GetNodeCount returns the NodeCount field value
func (o *ApiSearchDeploymentSpec) GetNodeCount() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentSpec) GetNodeCountOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NodeCount, true
}

// SetNodeCount sets field value
func (o *ApiSearchDeploymentSpec) SetNodeCount(v int) {
	o.NodeCount = v
}
