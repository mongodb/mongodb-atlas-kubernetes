// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterConnectionStrings Collection of Uniform Resource Locators that point to the MongoDB database.
type ClusterConnectionStrings struct {
	// Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.
	// Read only field.
	AwsPrivateLink *map[string]string `json:"awsPrivateLink,omitempty"`
	// Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink. Each key identifies an Amazon Web Services (AWS) interface endpoint. Each value identifies the related `mongodb://` connection string that you use to connect to Atlas through the interface endpoint that the key names. If the cluster uses an optimized connection string, `awsPrivateLinkSrv` contains the optimized connection string. If the cluster has the non-optimized (legacy) connection string, `awsPrivateLinkSrv` contains the non-optimized connection string even if an optimized connection string is also present.
	// Read only field.
	AwsPrivateLinkSrv *map[string]string `json:"awsPrivateLinkSrv,omitempty"`
	// Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter once someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change. Use this URI format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.
	// Read only field.
	Private *string `json:"private,omitempty"`
	// List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint. This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.
	// Read only field.
	PrivateEndpoint *[]ClusterDescriptionConnectionStringsPrivateEndpoint `json:"privateEndpoint,omitempty"`
	// Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster. This connection string uses the `mongodb+srv://` protocol. The resource returns this parameter when someone creates a network peering connection to this cluster. This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change. Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use `connectionStrings.private`. For Amazon Web Services (AWS) clusters, this parameter returns only if you [enable custom DNS](https://docs.atlas.mongodb.com/reference/api/aws-custom-dns-update/).
	// Read only field.
	PrivateSrv *string `json:"privateSrv,omitempty"`
	// Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb://` protocol.
	// Read only field.
	Standard *string `json:"standard,omitempty"`
	// Public connection string that you can use to connect to this cluster. This connection string uses the `mongodb+srv://` protocol.
	// Read only field.
	StandardSrv *string `json:"standardSrv,omitempty"`
}

// NewClusterConnectionStrings instantiates a new ClusterConnectionStrings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterConnectionStrings() *ClusterConnectionStrings {
	this := ClusterConnectionStrings{}
	return &this
}

// NewClusterConnectionStringsWithDefaults instantiates a new ClusterConnectionStrings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterConnectionStringsWithDefaults() *ClusterConnectionStrings {
	this := ClusterConnectionStrings{}
	return &this
}

// GetAwsPrivateLink returns the AwsPrivateLink field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetAwsPrivateLink() map[string]string {
	if o == nil || IsNil(o.AwsPrivateLink) {
		var ret map[string]string
		return ret
	}
	return *o.AwsPrivateLink
}

// GetAwsPrivateLinkOk returns a tuple with the AwsPrivateLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetAwsPrivateLinkOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.AwsPrivateLink) {
		return nil, false
	}

	return o.AwsPrivateLink, true
}

// HasAwsPrivateLink returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasAwsPrivateLink() bool {
	if o != nil && !IsNil(o.AwsPrivateLink) {
		return true
	}

	return false
}

// SetAwsPrivateLink gets a reference to the given map[string]string and assigns it to the AwsPrivateLink field.
func (o *ClusterConnectionStrings) SetAwsPrivateLink(v map[string]string) {
	o.AwsPrivateLink = &v
}

// GetAwsPrivateLinkSrv returns the AwsPrivateLinkSrv field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetAwsPrivateLinkSrv() map[string]string {
	if o == nil || IsNil(o.AwsPrivateLinkSrv) {
		var ret map[string]string
		return ret
	}
	return *o.AwsPrivateLinkSrv
}

// GetAwsPrivateLinkSrvOk returns a tuple with the AwsPrivateLinkSrv field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetAwsPrivateLinkSrvOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.AwsPrivateLinkSrv) {
		return nil, false
	}

	return o.AwsPrivateLinkSrv, true
}

// HasAwsPrivateLinkSrv returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasAwsPrivateLinkSrv() bool {
	if o != nil && !IsNil(o.AwsPrivateLinkSrv) {
		return true
	}

	return false
}

// SetAwsPrivateLinkSrv gets a reference to the given map[string]string and assigns it to the AwsPrivateLinkSrv field.
func (o *ClusterConnectionStrings) SetAwsPrivateLinkSrv(v map[string]string) {
	o.AwsPrivateLinkSrv = &v
}

// GetPrivate returns the Private field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetPrivate() string {
	if o == nil || IsNil(o.Private) {
		var ret string
		return ret
	}
	return *o.Private
}

// GetPrivateOk returns a tuple with the Private field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetPrivateOk() (*string, bool) {
	if o == nil || IsNil(o.Private) {
		return nil, false
	}

	return o.Private, true
}

// HasPrivate returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasPrivate() bool {
	if o != nil && !IsNil(o.Private) {
		return true
	}

	return false
}

// SetPrivate gets a reference to the given string and assigns it to the Private field.
func (o *ClusterConnectionStrings) SetPrivate(v string) {
	o.Private = &v
}

// GetPrivateEndpoint returns the PrivateEndpoint field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetPrivateEndpoint() []ClusterDescriptionConnectionStringsPrivateEndpoint {
	if o == nil || IsNil(o.PrivateEndpoint) {
		var ret []ClusterDescriptionConnectionStringsPrivateEndpoint
		return ret
	}
	return *o.PrivateEndpoint
}

// GetPrivateEndpointOk returns a tuple with the PrivateEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetPrivateEndpointOk() (*[]ClusterDescriptionConnectionStringsPrivateEndpoint, bool) {
	if o == nil || IsNil(o.PrivateEndpoint) {
		return nil, false
	}

	return o.PrivateEndpoint, true
}

// HasPrivateEndpoint returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasPrivateEndpoint() bool {
	if o != nil && !IsNil(o.PrivateEndpoint) {
		return true
	}

	return false
}

// SetPrivateEndpoint gets a reference to the given []ClusterDescriptionConnectionStringsPrivateEndpoint and assigns it to the PrivateEndpoint field.
func (o *ClusterConnectionStrings) SetPrivateEndpoint(v []ClusterDescriptionConnectionStringsPrivateEndpoint) {
	o.PrivateEndpoint = &v
}

// GetPrivateSrv returns the PrivateSrv field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetPrivateSrv() string {
	if o == nil || IsNil(o.PrivateSrv) {
		var ret string
		return ret
	}
	return *o.PrivateSrv
}

// GetPrivateSrvOk returns a tuple with the PrivateSrv field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetPrivateSrvOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateSrv) {
		return nil, false
	}

	return o.PrivateSrv, true
}

// HasPrivateSrv returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasPrivateSrv() bool {
	if o != nil && !IsNil(o.PrivateSrv) {
		return true
	}

	return false
}

// SetPrivateSrv gets a reference to the given string and assigns it to the PrivateSrv field.
func (o *ClusterConnectionStrings) SetPrivateSrv(v string) {
	o.PrivateSrv = &v
}

// GetStandard returns the Standard field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetStandard() string {
	if o == nil || IsNil(o.Standard) {
		var ret string
		return ret
	}
	return *o.Standard
}

// GetStandardOk returns a tuple with the Standard field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetStandardOk() (*string, bool) {
	if o == nil || IsNil(o.Standard) {
		return nil, false
	}

	return o.Standard, true
}

// HasStandard returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasStandard() bool {
	if o != nil && !IsNil(o.Standard) {
		return true
	}

	return false
}

// SetStandard gets a reference to the given string and assigns it to the Standard field.
func (o *ClusterConnectionStrings) SetStandard(v string) {
	o.Standard = &v
}

// GetStandardSrv returns the StandardSrv field value if set, zero value otherwise
func (o *ClusterConnectionStrings) GetStandardSrv() string {
	if o == nil || IsNil(o.StandardSrv) {
		var ret string
		return ret
	}
	return *o.StandardSrv
}

// GetStandardSrvOk returns a tuple with the StandardSrv field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterConnectionStrings) GetStandardSrvOk() (*string, bool) {
	if o == nil || IsNil(o.StandardSrv) {
		return nil, false
	}

	return o.StandardSrv, true
}

// HasStandardSrv returns a boolean if a field has been set.
func (o *ClusterConnectionStrings) HasStandardSrv() bool {
	if o != nil && !IsNil(o.StandardSrv) {
		return true
	}

	return false
}

// SetStandardSrv gets a reference to the given string and assigns it to the StandardSrv field.
func (o *ClusterConnectionStrings) SetStandardSrv(v string) {
	o.StandardSrv = &v
}
