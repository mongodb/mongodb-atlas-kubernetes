// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiSearchDeploymentResponse struct for ApiSearchDeploymentResponse
type ApiSearchDeploymentResponse struct {
	// Cloud service provider that manages your customer keys to provide an additional layer of Encryption At Rest for the cluster.
	// Read only field.
	EncryptionAtRestProvider *string `json:"encryptionAtRestProvider,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the search deployment.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of settings that configure the Search Nodes for your cluster. The configuration will be returned for each region and shard.
	// Read only field.
	Specs *[]ApiSearchDeploymentSpec `json:"specs,omitempty"`
	// Human-readable label that indicates the current operating condition of this search deployment.
	// Read only field.
	StateName *string `json:"stateName,omitempty"`
}

// NewApiSearchDeploymentResponse instantiates a new ApiSearchDeploymentResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiSearchDeploymentResponse() *ApiSearchDeploymentResponse {
	this := ApiSearchDeploymentResponse{}
	return &this
}

// NewApiSearchDeploymentResponseWithDefaults instantiates a new ApiSearchDeploymentResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiSearchDeploymentResponseWithDefaults() *ApiSearchDeploymentResponse {
	this := ApiSearchDeploymentResponse{}
	return &this
}

// GetEncryptionAtRestProvider returns the EncryptionAtRestProvider field value if set, zero value otherwise
func (o *ApiSearchDeploymentResponse) GetEncryptionAtRestProvider() string {
	if o == nil || IsNil(o.EncryptionAtRestProvider) {
		var ret string
		return ret
	}
	return *o.EncryptionAtRestProvider
}

// GetEncryptionAtRestProviderOk returns a tuple with the EncryptionAtRestProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentResponse) GetEncryptionAtRestProviderOk() (*string, bool) {
	if o == nil || IsNil(o.EncryptionAtRestProvider) {
		return nil, false
	}

	return o.EncryptionAtRestProvider, true
}

// HasEncryptionAtRestProvider returns a boolean if a field has been set.
func (o *ApiSearchDeploymentResponse) HasEncryptionAtRestProvider() bool {
	if o != nil && !IsNil(o.EncryptionAtRestProvider) {
		return true
	}

	return false
}

// SetEncryptionAtRestProvider gets a reference to the given string and assigns it to the EncryptionAtRestProvider field.
func (o *ApiSearchDeploymentResponse) SetEncryptionAtRestProvider(v string) {
	o.EncryptionAtRestProvider = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ApiSearchDeploymentResponse) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentResponse) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ApiSearchDeploymentResponse) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ApiSearchDeploymentResponse) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiSearchDeploymentResponse) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentResponse) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiSearchDeploymentResponse) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiSearchDeploymentResponse) SetId(v string) {
	o.Id = &v
}

// GetSpecs returns the Specs field value if set, zero value otherwise
func (o *ApiSearchDeploymentResponse) GetSpecs() []ApiSearchDeploymentSpec {
	if o == nil || IsNil(o.Specs) {
		var ret []ApiSearchDeploymentSpec
		return ret
	}
	return *o.Specs
}

// GetSpecsOk returns a tuple with the Specs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentResponse) GetSpecsOk() (*[]ApiSearchDeploymentSpec, bool) {
	if o == nil || IsNil(o.Specs) {
		return nil, false
	}

	return o.Specs, true
}

// HasSpecs returns a boolean if a field has been set.
func (o *ApiSearchDeploymentResponse) HasSpecs() bool {
	if o != nil && !IsNil(o.Specs) {
		return true
	}

	return false
}

// SetSpecs gets a reference to the given []ApiSearchDeploymentSpec and assigns it to the Specs field.
func (o *ApiSearchDeploymentResponse) SetSpecs(v []ApiSearchDeploymentSpec) {
	o.Specs = &v
}

// GetStateName returns the StateName field value if set, zero value otherwise
func (o *ApiSearchDeploymentResponse) GetStateName() string {
	if o == nil || IsNil(o.StateName) {
		var ret string
		return ret
	}
	return *o.StateName
}

// GetStateNameOk returns a tuple with the StateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiSearchDeploymentResponse) GetStateNameOk() (*string, bool) {
	if o == nil || IsNil(o.StateName) {
		return nil, false
	}

	return o.StateName, true
}

// HasStateName returns a boolean if a field has been set.
func (o *ApiSearchDeploymentResponse) HasStateName() bool {
	if o != nil && !IsNil(o.StateName) {
		return true
	}

	return false
}

// SetStateName gets a reference to the given string and assigns it to the StateName field.
func (o *ApiSearchDeploymentResponse) SetStateName(v string) {
	o.StateName = &v
}
