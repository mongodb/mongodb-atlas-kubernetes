// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// CloudDatabaseUser struct for CloudDatabaseUser
type CloudDatabaseUser struct {
	// Human-readable label that indicates whether the new database user authenticates with the Amazon Web Services (AWS) Identity and Access Management (IAM) credentials associated with the user or the user's role.
	AwsIAMType *string `json:"awsIAMType,omitempty"`
	// The database against which the database user authenticates. Database users must provide both a username and authentication database to log into MongoDB. If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be `$external`. If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be `admin`.
	DatabaseName string `json:"databaseName"`
	// Date and time when MongoDB Cloud deletes the user. This parameter expresses its value in the ISO 8601 timestamp format in UTC and can include the time zone designation. You must specify a future date that falls within one week of making the Application Programming Interface (API) request.
	DeleteAfterDate *time.Time `json:"deleteAfterDate,omitempty"`
	// Description of this database user.
	Description *string `json:"description,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project.
	// Write only field.
	GroupId string `json:"groupId"`
	// List that contains the key-value pairs for tagging and categorizing the MongoDB database user. The labels that you define do not appear in the console.
	Labels *[]ComponentLabel `json:"labels,omitempty"`
	// Part of the Lightweight Directory Access Protocol (LDAP) record that the database uses to authenticate this database user on the LDAP host.
	LdapAuthType *string `json:"ldapAuthType,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that indicates whether the new database user or group authenticates with OIDC federated authentication. To create a federated authentication user, specify the value of USER in this field. To create a federated authentication group, specify the value of `IDP_GROUP` in this field.
	OidcAuthType *string `json:"oidcAuthType,omitempty"`
	// Alphanumeric string that authenticates this database user against the database specified in `databaseName`. To authenticate with SCRAM-SHA, you must specify this parameter. This parameter doesn't appear in this response.
	// Write only field.
	Password *string `json:"password,omitempty"`
	// List that provides the pairings of one role with one applicable database.
	Roles []DatabaseUserRole `json:"roles"`
	// List that contains clusters, MongoDB Atlas Data Lakes, and MongoDB Atlas Streams Workspaces that this database user can access. If omitted, MongoDB Cloud grants the database user access to all the clusters, MongoDB Atlas Data Lakes, and MongoDB Atlas Streams Workspaces in the project.
	Scopes *[]UserScope `json:"scopes,omitempty"`
	// Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:  | Authentication Method | Parameter Needed | Parameter Value | username Format | |---|---|---|---| | AWS IAM | `awsIAMType` | `ROLE` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | AWS IAM | `awsIAMType` | `USER` | <abbr title=\"Amazon Resource Name\">ARN</abbr> | | x.509 | `x509Type` | `CUSTOMER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | x.509 | `x509Type` | `MANAGED` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `USER` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | LDAP | `ldapAuthType` | `GROUP` | [RFC 2253](https://tools.ietf.org/html/2253) Distinguished Name | | OIDC Workforce | `oidcAuthType` | `IDP_GROUP` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP group name | | OIDC Workload | `oidcAuthType` | `USER` | Atlas OIDC IdP ID (found in federation settings), followed by a '/', followed by the IdP user name | | SCRAM-SHA | `awsIAMType`, `x509Type`, `ldapAuthType`, `oidcAuthType` | `NONE` | Alphanumeric string |
	Username string `json:"username"`
	// X.509 method that MongoDB Cloud uses to authenticate the database user.  - For application-managed X.509, specify `MANAGED`. - For self-managed X.509, specify `CUSTOMER`.  Users created with the `CUSTOMER` method require a Common Name (CN) in the **username** parameter. You must create externally authenticated users on the `$external` database.
	X509Type *string `json:"x509Type,omitempty"`
}

// NewCloudDatabaseUser instantiates a new CloudDatabaseUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudDatabaseUser(databaseName string, groupId string, roles []DatabaseUserRole, username string) *CloudDatabaseUser {
	this := CloudDatabaseUser{}
	var awsIAMType string = "NONE"
	this.AwsIAMType = &awsIAMType
	this.DatabaseName = databaseName
	this.GroupId = groupId
	var ldapAuthType string = "NONE"
	this.LdapAuthType = &ldapAuthType
	var oidcAuthType string = "NONE"
	this.OidcAuthType = &oidcAuthType
	this.Roles = roles
	this.Username = username
	var x509Type string = "NONE"
	this.X509Type = &x509Type
	return &this
}

// NewCloudDatabaseUserWithDefaults instantiates a new CloudDatabaseUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudDatabaseUserWithDefaults() *CloudDatabaseUser {
	this := CloudDatabaseUser{}
	var awsIAMType string = "NONE"
	this.AwsIAMType = &awsIAMType
	var databaseName string = "admin"
	this.DatabaseName = databaseName
	var ldapAuthType string = "NONE"
	this.LdapAuthType = &ldapAuthType
	var oidcAuthType string = "NONE"
	this.OidcAuthType = &oidcAuthType
	var x509Type string = "NONE"
	this.X509Type = &x509Type
	return &this
}

// GetAwsIAMType returns the AwsIAMType field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetAwsIAMType() string {
	if o == nil || IsNil(o.AwsIAMType) {
		var ret string
		return ret
	}
	return *o.AwsIAMType
}

// GetAwsIAMTypeOk returns a tuple with the AwsIAMType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetAwsIAMTypeOk() (*string, bool) {
	if o == nil || IsNil(o.AwsIAMType) {
		return nil, false
	}

	return o.AwsIAMType, true
}

// HasAwsIAMType returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasAwsIAMType() bool {
	if o != nil && !IsNil(o.AwsIAMType) {
		return true
	}

	return false
}

// SetAwsIAMType gets a reference to the given string and assigns it to the AwsIAMType field.
func (o *CloudDatabaseUser) SetAwsIAMType(v string) {
	o.AwsIAMType = &v
}

// GetDatabaseName returns the DatabaseName field value
func (o *CloudDatabaseUser) GetDatabaseName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DatabaseName
}

// GetDatabaseNameOk returns a tuple with the DatabaseName field value
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetDatabaseNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DatabaseName, true
}

// SetDatabaseName sets field value
func (o *CloudDatabaseUser) SetDatabaseName(v string) {
	o.DatabaseName = v
}

// GetDeleteAfterDate returns the DeleteAfterDate field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetDeleteAfterDate() time.Time {
	if o == nil || IsNil(o.DeleteAfterDate) {
		var ret time.Time
		return ret
	}
	return *o.DeleteAfterDate
}

// GetDeleteAfterDateOk returns a tuple with the DeleteAfterDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetDeleteAfterDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.DeleteAfterDate) {
		return nil, false
	}

	return o.DeleteAfterDate, true
}

// HasDeleteAfterDate returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasDeleteAfterDate() bool {
	if o != nil && !IsNil(o.DeleteAfterDate) {
		return true
	}

	return false
}

// SetDeleteAfterDate gets a reference to the given time.Time and assigns it to the DeleteAfterDate field.
func (o *CloudDatabaseUser) SetDeleteAfterDate(v time.Time) {
	o.DeleteAfterDate = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *CloudDatabaseUser) SetDescription(v string) {
	o.Description = &v
}

// GetGroupId returns the GroupId field value
func (o *CloudDatabaseUser) GetGroupId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GroupId, true
}

// SetGroupId sets field value
func (o *CloudDatabaseUser) SetGroupId(v string) {
	o.GroupId = v
}

// GetLabels returns the Labels field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetLabels() []ComponentLabel {
	if o == nil || IsNil(o.Labels) {
		var ret []ComponentLabel
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetLabelsOk() (*[]ComponentLabel, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}

	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given []ComponentLabel and assigns it to the Labels field.
func (o *CloudDatabaseUser) SetLabels(v []ComponentLabel) {
	o.Labels = &v
}

// GetLdapAuthType returns the LdapAuthType field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetLdapAuthType() string {
	if o == nil || IsNil(o.LdapAuthType) {
		var ret string
		return ret
	}
	return *o.LdapAuthType
}

// GetLdapAuthTypeOk returns a tuple with the LdapAuthType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetLdapAuthTypeOk() (*string, bool) {
	if o == nil || IsNil(o.LdapAuthType) {
		return nil, false
	}

	return o.LdapAuthType, true
}

// HasLdapAuthType returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasLdapAuthType() bool {
	if o != nil && !IsNil(o.LdapAuthType) {
		return true
	}

	return false
}

// SetLdapAuthType gets a reference to the given string and assigns it to the LdapAuthType field.
func (o *CloudDatabaseUser) SetLdapAuthType(v string) {
	o.LdapAuthType = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *CloudDatabaseUser) SetLinks(v []Link) {
	o.Links = &v
}

// GetOidcAuthType returns the OidcAuthType field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetOidcAuthType() string {
	if o == nil || IsNil(o.OidcAuthType) {
		var ret string
		return ret
	}
	return *o.OidcAuthType
}

// GetOidcAuthTypeOk returns a tuple with the OidcAuthType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetOidcAuthTypeOk() (*string, bool) {
	if o == nil || IsNil(o.OidcAuthType) {
		return nil, false
	}

	return o.OidcAuthType, true
}

// HasOidcAuthType returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasOidcAuthType() bool {
	if o != nil && !IsNil(o.OidcAuthType) {
		return true
	}

	return false
}

// SetOidcAuthType gets a reference to the given string and assigns it to the OidcAuthType field.
func (o *CloudDatabaseUser) SetOidcAuthType(v string) {
	o.OidcAuthType = &v
}

// GetPassword returns the Password field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}

	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *CloudDatabaseUser) SetPassword(v string) {
	o.Password = &v
}

// GetRoles returns the Roles field value
func (o *CloudDatabaseUser) GetRoles() []DatabaseUserRole {
	if o == nil {
		var ret []DatabaseUserRole
		return ret
	}

	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetRolesOk() (*[]DatabaseUserRole, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Roles, true
}

// SetRoles sets field value
func (o *CloudDatabaseUser) SetRoles(v []DatabaseUserRole) {
	o.Roles = v
}

// GetScopes returns the Scopes field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetScopes() []UserScope {
	if o == nil || IsNil(o.Scopes) {
		var ret []UserScope
		return ret
	}
	return *o.Scopes
}

// GetScopesOk returns a tuple with the Scopes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetScopesOk() (*[]UserScope, bool) {
	if o == nil || IsNil(o.Scopes) {
		return nil, false
	}

	return o.Scopes, true
}

// HasScopes returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasScopes() bool {
	if o != nil && !IsNil(o.Scopes) {
		return true
	}

	return false
}

// SetScopes gets a reference to the given []UserScope and assigns it to the Scopes field.
func (o *CloudDatabaseUser) SetScopes(v []UserScope) {
	o.Scopes = &v
}

// GetUsername returns the Username field value
func (o *CloudDatabaseUser) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *CloudDatabaseUser) SetUsername(v string) {
	o.Username = v
}

// GetX509Type returns the X509Type field value if set, zero value otherwise
func (o *CloudDatabaseUser) GetX509Type() string {
	if o == nil || IsNil(o.X509Type) {
		var ret string
		return ret
	}
	return *o.X509Type
}

// GetX509TypeOk returns a tuple with the X509Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudDatabaseUser) GetX509TypeOk() (*string, bool) {
	if o == nil || IsNil(o.X509Type) {
		return nil, false
	}

	return o.X509Type, true
}

// HasX509Type returns a boolean if a field has been set.
func (o *CloudDatabaseUser) HasX509Type() bool {
	if o != nil && !IsNil(o.X509Type) {
		return true
	}

	return false
}

// SetX509Type gets a reference to the given string and assigns it to the X509Type field.
func (o *CloudDatabaseUser) SetX509Type(v string) {
	o.X509Type = &v
}
