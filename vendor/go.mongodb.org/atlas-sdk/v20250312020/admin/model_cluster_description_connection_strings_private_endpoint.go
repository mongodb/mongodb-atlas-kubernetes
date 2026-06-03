// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterDescriptionConnectionStringsPrivateEndpoint Private endpoint-aware connection string that you can use to connect to this cluster through a private endpoint.
type ClusterDescriptionConnectionStringsPrivateEndpoint struct {
	// Private endpoint-aware connection string that uses the `mongodb://` protocol to connect to MongoDB Cloud through a private endpoint.
	// Read only field.
	ConnectionString *string `json:"connectionString,omitempty"`
	// List that contains the private endpoints through which you connect to MongoDB Cloud when you use `connectionStrings.privateEndpoint[n].connectionString` or `connectionStrings.privateEndpoint[n].srvConnectionString`.
	// Read only field.
	Endpoints *[]ClusterDescriptionConnectionStringsPrivateEndpointEndpoint `json:"endpoints,omitempty"`
	// Private endpoint-aware connection string that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. The `mongodb+srv` protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application supports it. If it doesn't, use `connectionStrings.privateEndpoint[n].connectionString`.
	// Read only field.
	SrvConnectionString *string `json:"srvConnectionString,omitempty"`
	// Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your application and Atlas cluster supports it. If it doesn't, use and consult the documentation for `connectionStrings.privateEndpoint[n].srvConnectionString`.
	// Read only field.
	SrvShardOptimizedConnectionString *string `json:"srvShardOptimizedConnectionString,omitempty"`
	// MongoDB process type to which your application connects. Use `MONGOD` for replica sets and `MONGOS` for sharded clusters.
	// Read only field.
	Type *string `json:"type,omitempty"`
}

// NewClusterDescriptionConnectionStringsPrivateEndpoint instantiates a new ClusterDescriptionConnectionStringsPrivateEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterDescriptionConnectionStringsPrivateEndpoint() *ClusterDescriptionConnectionStringsPrivateEndpoint {
	this := ClusterDescriptionConnectionStringsPrivateEndpoint{}
	return &this
}

// NewClusterDescriptionConnectionStringsPrivateEndpointWithDefaults instantiates a new ClusterDescriptionConnectionStringsPrivateEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterDescriptionConnectionStringsPrivateEndpointWithDefaults() *ClusterDescriptionConnectionStringsPrivateEndpoint {
	this := ClusterDescriptionConnectionStringsPrivateEndpoint{}
	return &this
}

// GetConnectionString returns the ConnectionString field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetConnectionString() string {
	if o == nil || IsNil(o.ConnectionString) {
		var ret string
		return ret
	}
	return *o.ConnectionString
}

// GetConnectionStringOk returns a tuple with the ConnectionString field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetConnectionStringOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectionString) {
		return nil, false
	}

	return o.ConnectionString, true
}

// HasConnectionString returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) HasConnectionString() bool {
	if o != nil && !IsNil(o.ConnectionString) {
		return true
	}

	return false
}

// SetConnectionString gets a reference to the given string and assigns it to the ConnectionString field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) SetConnectionString(v string) {
	o.ConnectionString = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetEndpoints() []ClusterDescriptionConnectionStringsPrivateEndpointEndpoint {
	if o == nil || IsNil(o.Endpoints) {
		var ret []ClusterDescriptionConnectionStringsPrivateEndpointEndpoint
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetEndpointsOk() (*[]ClusterDescriptionConnectionStringsPrivateEndpointEndpoint, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}

	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []ClusterDescriptionConnectionStringsPrivateEndpointEndpoint and assigns it to the Endpoints field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) SetEndpoints(v []ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) {
	o.Endpoints = &v
}

// GetSrvConnectionString returns the SrvConnectionString field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetSrvConnectionString() string {
	if o == nil || IsNil(o.SrvConnectionString) {
		var ret string
		return ret
	}
	return *o.SrvConnectionString
}

// GetSrvConnectionStringOk returns a tuple with the SrvConnectionString field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetSrvConnectionStringOk() (*string, bool) {
	if o == nil || IsNil(o.SrvConnectionString) {
		return nil, false
	}

	return o.SrvConnectionString, true
}

// HasSrvConnectionString returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) HasSrvConnectionString() bool {
	if o != nil && !IsNil(o.SrvConnectionString) {
		return true
	}

	return false
}

// SetSrvConnectionString gets a reference to the given string and assigns it to the SrvConnectionString field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) SetSrvConnectionString(v string) {
	o.SrvConnectionString = &v
}

// GetSrvShardOptimizedConnectionString returns the SrvShardOptimizedConnectionString field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetSrvShardOptimizedConnectionString() string {
	if o == nil || IsNil(o.SrvShardOptimizedConnectionString) {
		var ret string
		return ret
	}
	return *o.SrvShardOptimizedConnectionString
}

// GetSrvShardOptimizedConnectionStringOk returns a tuple with the SrvShardOptimizedConnectionString field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetSrvShardOptimizedConnectionStringOk() (*string, bool) {
	if o == nil || IsNil(o.SrvShardOptimizedConnectionString) {
		return nil, false
	}

	return o.SrvShardOptimizedConnectionString, true
}

// HasSrvShardOptimizedConnectionString returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) HasSrvShardOptimizedConnectionString() bool {
	if o != nil && !IsNil(o.SrvShardOptimizedConnectionString) {
		return true
	}

	return false
}

// SetSrvShardOptimizedConnectionString gets a reference to the given string and assigns it to the SrvShardOptimizedConnectionString field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) SetSrvShardOptimizedConnectionString(v string) {
	o.SrvShardOptimizedConnectionString = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpoint) SetType(v string) {
	o.Type = &v
}
