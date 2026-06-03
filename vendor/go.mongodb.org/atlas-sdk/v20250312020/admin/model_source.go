// Code based on the AtlasAPI V2 OpenAPI file

package admin

// Source Document that describes the source of the migration.
type Source struct {
	// Path to the CA certificate that signed SSL certificates use to authenticate to the source cluster.
	CaCertificatePath *string `json:"caCertificatePath,omitempty"`
	// Label that identifies the source cluster name.
	ClusterName string `json:"clusterName"`
	// Unique 24-hexadecimal digit string that identifies the source project.
	GroupId string `json:"groupId"`
	// Flag that indicates whether MongoDB Automation manages authentication to the source cluster. If true, do not provide values for username and password.
	ManagedAuthentication bool `json:"managedAuthentication"`
	// Password that authenticates the username to the source cluster.
	// Write only field.
	Password *string `json:"password,omitempty"`
	// Flag that indicates whether you have SSL enabled.
	Ssl bool `json:"ssl"`
	// Label that identifies the SCRAM-SHA user that connects to the source cluster.
	// Write only field.
	Username *string `json:"username,omitempty"`
}

// NewSource instantiates a new Source object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSource(clusterName string, groupId string, managedAuthentication bool, ssl bool) *Source {
	this := Source{}
	this.ClusterName = clusterName
	this.GroupId = groupId
	this.ManagedAuthentication = managedAuthentication
	this.Ssl = ssl
	return &this
}

// NewSourceWithDefaults instantiates a new Source object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSourceWithDefaults() *Source {
	this := Source{}
	return &this
}

// GetCaCertificatePath returns the CaCertificatePath field value if set, zero value otherwise
func (o *Source) GetCaCertificatePath() string {
	if o == nil || IsNil(o.CaCertificatePath) {
		var ret string
		return ret
	}
	return *o.CaCertificatePath
}

// GetCaCertificatePathOk returns a tuple with the CaCertificatePath field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Source) GetCaCertificatePathOk() (*string, bool) {
	if o == nil || IsNil(o.CaCertificatePath) {
		return nil, false
	}

	return o.CaCertificatePath, true
}

// HasCaCertificatePath returns a boolean if a field has been set.
func (o *Source) HasCaCertificatePath() bool {
	if o != nil && !IsNil(o.CaCertificatePath) {
		return true
	}

	return false
}

// SetCaCertificatePath gets a reference to the given string and assigns it to the CaCertificatePath field.
func (o *Source) SetCaCertificatePath(v string) {
	o.CaCertificatePath = &v
}

// GetClusterName returns the ClusterName field value
func (o *Source) GetClusterName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value
// and a boolean to check if the value has been set.
func (o *Source) GetClusterNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterName, true
}

// SetClusterName sets field value
func (o *Source) SetClusterName(v string) {
	o.ClusterName = v
}

// GetGroupId returns the GroupId field value
func (o *Source) GetGroupId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value
// and a boolean to check if the value has been set.
func (o *Source) GetGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GroupId, true
}

// SetGroupId sets field value
func (o *Source) SetGroupId(v string) {
	o.GroupId = v
}

// GetManagedAuthentication returns the ManagedAuthentication field value
func (o *Source) GetManagedAuthentication() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.ManagedAuthentication
}

// GetManagedAuthenticationOk returns a tuple with the ManagedAuthentication field value
// and a boolean to check if the value has been set.
func (o *Source) GetManagedAuthenticationOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ManagedAuthentication, true
}

// SetManagedAuthentication sets field value
func (o *Source) SetManagedAuthentication(v bool) {
	o.ManagedAuthentication = v
}

// GetPassword returns the Password field value if set, zero value otherwise
func (o *Source) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Source) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}

	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *Source) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *Source) SetPassword(v string) {
	o.Password = &v
}

// GetSsl returns the Ssl field value
func (o *Source) GetSsl() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Ssl
}

// GetSslOk returns a tuple with the Ssl field value
// and a boolean to check if the value has been set.
func (o *Source) GetSslOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ssl, true
}

// SetSsl sets field value
func (o *Source) SetSsl(v bool) {
	o.Ssl = v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *Source) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Source) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *Source) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *Source) SetUsername(v string) {
	o.Username = &v
}
