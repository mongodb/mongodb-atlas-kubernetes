// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServerlessConnectionStringsPrivateEndpointList Private endpoint connection string that you can use to connect to this serverless instance through a private endpoint.
type ServerlessConnectionStringsPrivateEndpointList struct {
	// List that contains the private endpoints through which you connect to MongoDB Cloud when you use `connectionStrings.privateEndpoint[n].srvConnectionString`.
	// Read only field.
	Endpoints *[]ServerlessConnectionStringsPrivateEndpointItem `json:"endpoints,omitempty"`
	// Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS).
	// Read only field.
	SrvConnectionString *string `json:"srvConnectionString,omitempty"`
	// MongoDB process type to which your application connects.
	// Read only field.
	Type *string `json:"type,omitempty"`
}

// NewServerlessConnectionStringsPrivateEndpointList instantiates a new ServerlessConnectionStringsPrivateEndpointList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessConnectionStringsPrivateEndpointList() *ServerlessConnectionStringsPrivateEndpointList {
	this := ServerlessConnectionStringsPrivateEndpointList{}
	return &this
}

// NewServerlessConnectionStringsPrivateEndpointListWithDefaults instantiates a new ServerlessConnectionStringsPrivateEndpointList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessConnectionStringsPrivateEndpointListWithDefaults() *ServerlessConnectionStringsPrivateEndpointList {
	this := ServerlessConnectionStringsPrivateEndpointList{}
	return &this
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise
func (o *ServerlessConnectionStringsPrivateEndpointList) GetEndpoints() []ServerlessConnectionStringsPrivateEndpointItem {
	if o == nil || IsNil(o.Endpoints) {
		var ret []ServerlessConnectionStringsPrivateEndpointItem
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessConnectionStringsPrivateEndpointList) GetEndpointsOk() (*[]ServerlessConnectionStringsPrivateEndpointItem, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}

	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *ServerlessConnectionStringsPrivateEndpointList) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []ServerlessConnectionStringsPrivateEndpointItem and assigns it to the Endpoints field.
func (o *ServerlessConnectionStringsPrivateEndpointList) SetEndpoints(v []ServerlessConnectionStringsPrivateEndpointItem) {
	o.Endpoints = &v
}

// GetSrvConnectionString returns the SrvConnectionString field value if set, zero value otherwise
func (o *ServerlessConnectionStringsPrivateEndpointList) GetSrvConnectionString() string {
	if o == nil || IsNil(o.SrvConnectionString) {
		var ret string
		return ret
	}
	return *o.SrvConnectionString
}

// GetSrvConnectionStringOk returns a tuple with the SrvConnectionString field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessConnectionStringsPrivateEndpointList) GetSrvConnectionStringOk() (*string, bool) {
	if o == nil || IsNil(o.SrvConnectionString) {
		return nil, false
	}

	return o.SrvConnectionString, true
}

// HasSrvConnectionString returns a boolean if a field has been set.
func (o *ServerlessConnectionStringsPrivateEndpointList) HasSrvConnectionString() bool {
	if o != nil && !IsNil(o.SrvConnectionString) {
		return true
	}

	return false
}

// SetSrvConnectionString gets a reference to the given string and assigns it to the SrvConnectionString field.
func (o *ServerlessConnectionStringsPrivateEndpointList) SetSrvConnectionString(v string) {
	o.SrvConnectionString = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *ServerlessConnectionStringsPrivateEndpointList) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessConnectionStringsPrivateEndpointList) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ServerlessConnectionStringsPrivateEndpointList) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ServerlessConnectionStringsPrivateEndpointList) SetType(v string) {
	o.Type = &v
}
