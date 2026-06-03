// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServerlessInstanceDescriptionConnectionStrings Collection of Uniform Resource Locators that point to the MongoDB database.
type ServerlessInstanceDescriptionConnectionStrings struct {
	// List of private endpoint-aware connection strings that you can use to connect to this serverless instance through a private endpoint. This parameter returns only if you created a private endpoint for this serverless instance and it is AVAILABLE.
	// Read only field.
	PrivateEndpoint *[]ServerlessConnectionStringsPrivateEndpointList `json:"privateEndpoint,omitempty"`
	// Public connection string that you can use to connect to this serverless instance. This connection string uses the `mongodb+srv://` protocol.
	// Read only field.
	StandardSrv *string `json:"standardSrv,omitempty"`
}

// NewServerlessInstanceDescriptionConnectionStrings instantiates a new ServerlessInstanceDescriptionConnectionStrings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessInstanceDescriptionConnectionStrings() *ServerlessInstanceDescriptionConnectionStrings {
	this := ServerlessInstanceDescriptionConnectionStrings{}
	return &this
}

// NewServerlessInstanceDescriptionConnectionStringsWithDefaults instantiates a new ServerlessInstanceDescriptionConnectionStrings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessInstanceDescriptionConnectionStringsWithDefaults() *ServerlessInstanceDescriptionConnectionStrings {
	this := ServerlessInstanceDescriptionConnectionStrings{}
	return &this
}

// GetPrivateEndpoint returns the PrivateEndpoint field value if set, zero value otherwise
func (o *ServerlessInstanceDescriptionConnectionStrings) GetPrivateEndpoint() []ServerlessConnectionStringsPrivateEndpointList {
	if o == nil || IsNil(o.PrivateEndpoint) {
		var ret []ServerlessConnectionStringsPrivateEndpointList
		return ret
	}
	return *o.PrivateEndpoint
}

// GetPrivateEndpointOk returns a tuple with the PrivateEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescriptionConnectionStrings) GetPrivateEndpointOk() (*[]ServerlessConnectionStringsPrivateEndpointList, bool) {
	if o == nil || IsNil(o.PrivateEndpoint) {
		return nil, false
	}

	return o.PrivateEndpoint, true
}

// HasPrivateEndpoint returns a boolean if a field has been set.
func (o *ServerlessInstanceDescriptionConnectionStrings) HasPrivateEndpoint() bool {
	if o != nil && !IsNil(o.PrivateEndpoint) {
		return true
	}

	return false
}

// SetPrivateEndpoint gets a reference to the given []ServerlessConnectionStringsPrivateEndpointList and assigns it to the PrivateEndpoint field.
func (o *ServerlessInstanceDescriptionConnectionStrings) SetPrivateEndpoint(v []ServerlessConnectionStringsPrivateEndpointList) {
	o.PrivateEndpoint = &v
}

// GetStandardSrv returns the StandardSrv field value if set, zero value otherwise
func (o *ServerlessInstanceDescriptionConnectionStrings) GetStandardSrv() string {
	if o == nil || IsNil(o.StandardSrv) {
		var ret string
		return ret
	}
	return *o.StandardSrv
}

// GetStandardSrvOk returns a tuple with the StandardSrv field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescriptionConnectionStrings) GetStandardSrvOk() (*string, bool) {
	if o == nil || IsNil(o.StandardSrv) {
		return nil, false
	}

	return o.StandardSrv, true
}

// HasStandardSrv returns a boolean if a field has been set.
func (o *ServerlessInstanceDescriptionConnectionStrings) HasStandardSrv() bool {
	if o != nil && !IsNil(o.StandardSrv) {
		return true
	}

	return false
}

// SetStandardSrv gets a reference to the given string and assigns it to the StandardSrv field.
func (o *ServerlessInstanceDescriptionConnectionStrings) SetStandardSrv(v string) {
	o.StandardSrv = &v
}
