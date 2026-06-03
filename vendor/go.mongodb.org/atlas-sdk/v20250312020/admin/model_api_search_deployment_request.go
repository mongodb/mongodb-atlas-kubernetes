// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiSearchDeploymentRequest struct for ApiSearchDeploymentRequest
type ApiSearchDeploymentRequest struct {
	// List of settings that configure the Search Nodes for your cluster. Provide one element per region when configuring asymmetric deployments; a single element applies to all regions.
	Specs []ApiSearchDeploymentRequestSpec `json:"specs"`
}

// NewApiSearchDeploymentRequest instantiates a new ApiSearchDeploymentRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiSearchDeploymentRequest(specs []ApiSearchDeploymentRequestSpec) *ApiSearchDeploymentRequest {
	this := ApiSearchDeploymentRequest{}
	this.Specs = specs
	return &this
}

// NewApiSearchDeploymentRequestWithDefaults instantiates a new ApiSearchDeploymentRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiSearchDeploymentRequestWithDefaults() *ApiSearchDeploymentRequest {
	this := ApiSearchDeploymentRequest{}
	return &this
}

// GetSpecs returns the Specs field value
func (o *ApiSearchDeploymentRequest) GetSpecs() []ApiSearchDeploymentRequestSpec {
	if o == nil {
		var ret []ApiSearchDeploymentRequestSpec
		return ret
	}

	return o.Specs
}

// GetSpecsOk returns a tuple with the Specs field value
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentRequest) GetSpecsOk() (*[]ApiSearchDeploymentRequestSpec, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Specs, true
}

// SetSpecs sets field value
func (o *ApiSearchDeploymentRequest) SetSpecs(v []ApiSearchDeploymentRequestSpec) {
	o.Specs = v
}
