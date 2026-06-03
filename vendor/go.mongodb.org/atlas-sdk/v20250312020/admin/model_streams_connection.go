// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsConnection Settings that define a connection to an external data store.
type StreamsConnection struct {
	// Unique identifier of the connection.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the stream connection. In the case of the Sample type, this is the name of the sample source.
	Name *string `json:"name,omitempty"`
	// The connection's region.
	Region *string `json:"region,omitempty"`
	// The state of the connection.
	// Read only field.
	State *string `json:"state,omitempty"`
	// Type of the connection.
	Type *string `json:"type,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the configured cluster. Required if the ID does not match the project containing the streams workspace. You must first enable the organization setting.
	ClusterGroupId *string `json:"clusterGroupId,omitempty"`
	// Name of the cluster configured for this connection.
	ClusterName     *string                     `json:"clusterName,omitempty"`
	DbRoleToExecute *DBRoleToExecute            `json:"dbRoleToExecute,omitempty"`
	Authentication  *StreamsKafkaAuthentication `json:"authentication,omitempty"`
	// Comma separated list of server addresses.
	BootstrapServers *string `json:"bootstrapServers,omitempty"`
	// A map of Kafka key-value pairs for optional configuration. This is a flat object, and keys can have '.' characters.
	Config     *map[string]string      `json:"config,omitempty"`
	Networking *StreamsKafkaNetworking `json:"networking,omitempty"`
	Security   *StreamsKafkaSecurity   `json:"security,omitempty"`
	// A map of key-value pairs that will be passed as headers for the request.
	Headers *map[string]string `json:"headers,omitempty"`
	// The URL to be used for the request.
	Url *string                     `json:"url,omitempty"`
	Aws *StreamsAWSConnectionConfig `json:"aws,omitempty"`
	// The Schema Registry provider.
	Provider                     *string                       `json:"provider,omitempty"`
	SchemaRegistryAuthentication *SchemaRegistryAuthentication `json:"schemaRegistryAuthentication,omitempty"`
	// List of Schema Registry endpoint URLs used by this connection. Each URL must use the http or https scheme and specify a valid host and optional port.
	SchemaRegistryUrls      *[]string                           `json:"schemaRegistryUrls,omitempty"`
	Azure                   *AzureConnection                    `json:"azure,omitempty"`
	PublicPrivateNetworking *StreamsPublicPrivateLinkNetworking `json:"publicPrivateNetworking,omitempty"`
	Gcp                     *StreamsGCPConnectionConfig         `json:"gcp,omitempty"`
}

// NewStreamsConnection instantiates a new StreamsConnection object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsConnection() *StreamsConnection {
	this := StreamsConnection{}
	return &this
}

// NewStreamsConnectionWithDefaults instantiates a new StreamsConnection object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsConnectionWithDefaults() *StreamsConnection {
	this := StreamsConnection{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *StreamsConnection) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *StreamsConnection) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *StreamsConnection) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsConnection) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsConnection) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsConnection) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *StreamsConnection) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *StreamsConnection) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *StreamsConnection) SetName(v string) {
	o.Name = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *StreamsConnection) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *StreamsConnection) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *StreamsConnection) SetRegion(v string) {
	o.Region = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *StreamsConnection) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *StreamsConnection) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *StreamsConnection) SetState(v string) {
	o.State = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *StreamsConnection) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *StreamsConnection) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *StreamsConnection) SetType(v string) {
	o.Type = &v
}

// GetClusterGroupId returns the ClusterGroupId field value if set, zero value otherwise
func (o *StreamsConnection) GetClusterGroupId() string {
	if o == nil || IsNil(o.ClusterGroupId) {
		var ret string
		return ret
	}
	return *o.ClusterGroupId
}

// GetClusterGroupIdOk returns a tuple with the ClusterGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetClusterGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterGroupId) {
		return nil, false
	}

	return o.ClusterGroupId, true
}

// HasClusterGroupId returns a boolean if a field has been set.
func (o *StreamsConnection) HasClusterGroupId() bool {
	if o != nil && !IsNil(o.ClusterGroupId) {
		return true
	}

	return false
}

// SetClusterGroupId gets a reference to the given string and assigns it to the ClusterGroupId field.
func (o *StreamsConnection) SetClusterGroupId(v string) {
	o.ClusterGroupId = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *StreamsConnection) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *StreamsConnection) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *StreamsConnection) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetDbRoleToExecute returns the DbRoleToExecute field value if set, zero value otherwise
func (o *StreamsConnection) GetDbRoleToExecute() DBRoleToExecute {
	if o == nil || IsNil(o.DbRoleToExecute) {
		var ret DBRoleToExecute
		return ret
	}
	return *o.DbRoleToExecute
}

// GetDbRoleToExecuteOk returns a tuple with the DbRoleToExecute field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetDbRoleToExecuteOk() (*DBRoleToExecute, bool) {
	if o == nil || IsNil(o.DbRoleToExecute) {
		return nil, false
	}

	return o.DbRoleToExecute, true
}

// HasDbRoleToExecute returns a boolean if a field has been set.
func (o *StreamsConnection) HasDbRoleToExecute() bool {
	if o != nil && !IsNil(o.DbRoleToExecute) {
		return true
	}

	return false
}

// SetDbRoleToExecute gets a reference to the given DBRoleToExecute and assigns it to the DbRoleToExecute field.
func (o *StreamsConnection) SetDbRoleToExecute(v DBRoleToExecute) {
	o.DbRoleToExecute = &v
}

// GetAuthentication returns the Authentication field value if set, zero value otherwise
func (o *StreamsConnection) GetAuthentication() StreamsKafkaAuthentication {
	if o == nil || IsNil(o.Authentication) {
		var ret StreamsKafkaAuthentication
		return ret
	}
	return *o.Authentication
}

// GetAuthenticationOk returns a tuple with the Authentication field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetAuthenticationOk() (*StreamsKafkaAuthentication, bool) {
	if o == nil || IsNil(o.Authentication) {
		return nil, false
	}

	return o.Authentication, true
}

// HasAuthentication returns a boolean if a field has been set.
func (o *StreamsConnection) HasAuthentication() bool {
	if o != nil && !IsNil(o.Authentication) {
		return true
	}

	return false
}

// SetAuthentication gets a reference to the given StreamsKafkaAuthentication and assigns it to the Authentication field.
func (o *StreamsConnection) SetAuthentication(v StreamsKafkaAuthentication) {
	o.Authentication = &v
}

// GetBootstrapServers returns the BootstrapServers field value if set, zero value otherwise
func (o *StreamsConnection) GetBootstrapServers() string {
	if o == nil || IsNil(o.BootstrapServers) {
		var ret string
		return ret
	}
	return *o.BootstrapServers
}

// GetBootstrapServersOk returns a tuple with the BootstrapServers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetBootstrapServersOk() (*string, bool) {
	if o == nil || IsNil(o.BootstrapServers) {
		return nil, false
	}

	return o.BootstrapServers, true
}

// HasBootstrapServers returns a boolean if a field has been set.
func (o *StreamsConnection) HasBootstrapServers() bool {
	if o != nil && !IsNil(o.BootstrapServers) {
		return true
	}

	return false
}

// SetBootstrapServers gets a reference to the given string and assigns it to the BootstrapServers field.
func (o *StreamsConnection) SetBootstrapServers(v string) {
	o.BootstrapServers = &v
}

// GetConfig returns the Config field value if set, zero value otherwise
func (o *StreamsConnection) GetConfig() map[string]string {
	if o == nil || IsNil(o.Config) {
		var ret map[string]string
		return ret
	}
	return *o.Config
}

// GetConfigOk returns a tuple with the Config field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetConfigOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Config) {
		return nil, false
	}

	return o.Config, true
}

// HasConfig returns a boolean if a field has been set.
func (o *StreamsConnection) HasConfig() bool {
	if o != nil && !IsNil(o.Config) {
		return true
	}

	return false
}

// SetConfig gets a reference to the given map[string]string and assigns it to the Config field.
func (o *StreamsConnection) SetConfig(v map[string]string) {
	o.Config = &v
}

// GetNetworking returns the Networking field value if set, zero value otherwise
func (o *StreamsConnection) GetNetworking() StreamsKafkaNetworking {
	if o == nil || IsNil(o.Networking) {
		var ret StreamsKafkaNetworking
		return ret
	}
	return *o.Networking
}

// GetNetworkingOk returns a tuple with the Networking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetNetworkingOk() (*StreamsKafkaNetworking, bool) {
	if o == nil || IsNil(o.Networking) {
		return nil, false
	}

	return o.Networking, true
}

// HasNetworking returns a boolean if a field has been set.
func (o *StreamsConnection) HasNetworking() bool {
	if o != nil && !IsNil(o.Networking) {
		return true
	}

	return false
}

// SetNetworking gets a reference to the given StreamsKafkaNetworking and assigns it to the Networking field.
func (o *StreamsConnection) SetNetworking(v StreamsKafkaNetworking) {
	o.Networking = &v
}

// GetSecurity returns the Security field value if set, zero value otherwise
func (o *StreamsConnection) GetSecurity() StreamsKafkaSecurity {
	if o == nil || IsNil(o.Security) {
		var ret StreamsKafkaSecurity
		return ret
	}
	return *o.Security
}

// GetSecurityOk returns a tuple with the Security field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetSecurityOk() (*StreamsKafkaSecurity, bool) {
	if o == nil || IsNil(o.Security) {
		return nil, false
	}

	return o.Security, true
}

// HasSecurity returns a boolean if a field has been set.
func (o *StreamsConnection) HasSecurity() bool {
	if o != nil && !IsNil(o.Security) {
		return true
	}

	return false
}

// SetSecurity gets a reference to the given StreamsKafkaSecurity and assigns it to the Security field.
func (o *StreamsConnection) SetSecurity(v StreamsKafkaSecurity) {
	o.Security = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise
func (o *StreamsConnection) GetHeaders() map[string]string {
	if o == nil || IsNil(o.Headers) {
		var ret map[string]string
		return ret
	}
	return *o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetHeadersOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Headers) {
		return nil, false
	}

	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *StreamsConnection) HasHeaders() bool {
	if o != nil && !IsNil(o.Headers) {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given map[string]string and assigns it to the Headers field.
func (o *StreamsConnection) SetHeaders(v map[string]string) {
	o.Headers = &v
}

// GetUrl returns the Url field value if set, zero value otherwise
func (o *StreamsConnection) GetUrl() string {
	if o == nil || IsNil(o.Url) {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetUrlOk() (*string, bool) {
	if o == nil || IsNil(o.Url) {
		return nil, false
	}

	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *StreamsConnection) HasUrl() bool {
	if o != nil && !IsNil(o.Url) {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
func (o *StreamsConnection) SetUrl(v string) {
	o.Url = &v
}

// GetAws returns the Aws field value if set, zero value otherwise
func (o *StreamsConnection) GetAws() StreamsAWSConnectionConfig {
	if o == nil || IsNil(o.Aws) {
		var ret StreamsAWSConnectionConfig
		return ret
	}
	return *o.Aws
}

// GetAwsOk returns a tuple with the Aws field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetAwsOk() (*StreamsAWSConnectionConfig, bool) {
	if o == nil || IsNil(o.Aws) {
		return nil, false
	}

	return o.Aws, true
}

// HasAws returns a boolean if a field has been set.
func (o *StreamsConnection) HasAws() bool {
	if o != nil && !IsNil(o.Aws) {
		return true
	}

	return false
}

// SetAws gets a reference to the given StreamsAWSConnectionConfig and assigns it to the Aws field.
func (o *StreamsConnection) SetAws(v StreamsAWSConnectionConfig) {
	o.Aws = &v
}

// GetProvider returns the Provider field value if set, zero value otherwise
func (o *StreamsConnection) GetProvider() string {
	if o == nil || IsNil(o.Provider) {
		var ret string
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetProviderOk() (*string, bool) {
	if o == nil || IsNil(o.Provider) {
		return nil, false
	}

	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *StreamsConnection) HasProvider() bool {
	if o != nil && !IsNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given string and assigns it to the Provider field.
func (o *StreamsConnection) SetProvider(v string) {
	o.Provider = &v
}

// GetSchemaRegistryAuthentication returns the SchemaRegistryAuthentication field value if set, zero value otherwise
func (o *StreamsConnection) GetSchemaRegistryAuthentication() SchemaRegistryAuthentication {
	if o == nil || IsNil(o.SchemaRegistryAuthentication) {
		var ret SchemaRegistryAuthentication
		return ret
	}
	return *o.SchemaRegistryAuthentication
}

// GetSchemaRegistryAuthenticationOk returns a tuple with the SchemaRegistryAuthentication field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetSchemaRegistryAuthenticationOk() (*SchemaRegistryAuthentication, bool) {
	if o == nil || IsNil(o.SchemaRegistryAuthentication) {
		return nil, false
	}

	return o.SchemaRegistryAuthentication, true
}

// HasSchemaRegistryAuthentication returns a boolean if a field has been set.
func (o *StreamsConnection) HasSchemaRegistryAuthentication() bool {
	if o != nil && !IsNil(o.SchemaRegistryAuthentication) {
		return true
	}

	return false
}

// SetSchemaRegistryAuthentication gets a reference to the given SchemaRegistryAuthentication and assigns it to the SchemaRegistryAuthentication field.
func (o *StreamsConnection) SetSchemaRegistryAuthentication(v SchemaRegistryAuthentication) {
	o.SchemaRegistryAuthentication = &v
}

// GetSchemaRegistryUrls returns the SchemaRegistryUrls field value if set, zero value otherwise
func (o *StreamsConnection) GetSchemaRegistryUrls() []string {
	if o == nil || IsNil(o.SchemaRegistryUrls) {
		var ret []string
		return ret
	}
	return *o.SchemaRegistryUrls
}

// GetSchemaRegistryUrlsOk returns a tuple with the SchemaRegistryUrls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetSchemaRegistryUrlsOk() (*[]string, bool) {
	if o == nil || IsNil(o.SchemaRegistryUrls) {
		return nil, false
	}

	return o.SchemaRegistryUrls, true
}

// HasSchemaRegistryUrls returns a boolean if a field has been set.
func (o *StreamsConnection) HasSchemaRegistryUrls() bool {
	if o != nil && !IsNil(o.SchemaRegistryUrls) {
		return true
	}

	return false
}

// SetSchemaRegistryUrls gets a reference to the given []string and assigns it to the SchemaRegistryUrls field.
func (o *StreamsConnection) SetSchemaRegistryUrls(v []string) {
	o.SchemaRegistryUrls = &v
}

// GetAzure returns the Azure field value if set, zero value otherwise
func (o *StreamsConnection) GetAzure() AzureConnection {
	if o == nil || IsNil(o.Azure) {
		var ret AzureConnection
		return ret
	}
	return *o.Azure
}

// GetAzureOk returns a tuple with the Azure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetAzureOk() (*AzureConnection, bool) {
	if o == nil || IsNil(o.Azure) {
		return nil, false
	}

	return o.Azure, true
}

// HasAzure returns a boolean if a field has been set.
func (o *StreamsConnection) HasAzure() bool {
	if o != nil && !IsNil(o.Azure) {
		return true
	}

	return false
}

// SetAzure gets a reference to the given AzureConnection and assigns it to the Azure field.
func (o *StreamsConnection) SetAzure(v AzureConnection) {
	o.Azure = &v
}

// GetPublicPrivateNetworking returns the PublicPrivateNetworking field value if set, zero value otherwise
func (o *StreamsConnection) GetPublicPrivateNetworking() StreamsPublicPrivateLinkNetworking {
	if o == nil || IsNil(o.PublicPrivateNetworking) {
		var ret StreamsPublicPrivateLinkNetworking
		return ret
	}
	return *o.PublicPrivateNetworking
}

// GetPublicPrivateNetworkingOk returns a tuple with the PublicPrivateNetworking field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetPublicPrivateNetworkingOk() (*StreamsPublicPrivateLinkNetworking, bool) {
	if o == nil || IsNil(o.PublicPrivateNetworking) {
		return nil, false
	}

	return o.PublicPrivateNetworking, true
}

// HasPublicPrivateNetworking returns a boolean if a field has been set.
func (o *StreamsConnection) HasPublicPrivateNetworking() bool {
	if o != nil && !IsNil(o.PublicPrivateNetworking) {
		return true
	}

	return false
}

// SetPublicPrivateNetworking gets a reference to the given StreamsPublicPrivateLinkNetworking and assigns it to the PublicPrivateNetworking field.
func (o *StreamsConnection) SetPublicPrivateNetworking(v StreamsPublicPrivateLinkNetworking) {
	o.PublicPrivateNetworking = &v
}

// GetGcp returns the Gcp field value if set, zero value otherwise
func (o *StreamsConnection) GetGcp() StreamsGCPConnectionConfig {
	if o == nil || IsNil(o.Gcp) {
		var ret StreamsGCPConnectionConfig
		return ret
	}
	return *o.Gcp
}

// GetGcpOk returns a tuple with the Gcp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsConnection) GetGcpOk() (*StreamsGCPConnectionConfig, bool) {
	if o == nil || IsNil(o.Gcp) {
		return nil, false
	}

	return o.Gcp, true
}

// HasGcp returns a boolean if a field has been set.
func (o *StreamsConnection) HasGcp() bool {
	if o != nil && !IsNil(o.Gcp) {
		return true
	}

	return false
}

// SetGcp gets a reference to the given StreamsGCPConnectionConfig and assigns it to the Gcp field.
func (o *StreamsConnection) SetGcp(v StreamsGCPConnectionConfig) {
	o.Gcp = &v
}
