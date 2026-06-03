// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LDAPVerifyConnectivityJobRequest struct for LDAPVerifyConnectivityJobRequest
type LDAPVerifyConnectivityJobRequest struct {
	// Unique 24-hexadecimal digit string that identifies the project associated with this Lightweight Directory Access Protocol (LDAP) over Transport Layer Security (TLS) configuration.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links   *[]Link                                 `json:"links,omitempty"`
	Request *LDAPVerifyConnectivityJobRequestParams `json:"request,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this request to verify an Lightweight Directory Access Protocol (LDAP) configuration.
	// Read only field.
	RequestId *string `json:"requestId,omitempty"`
	// Human-readable string that indicates the status of the Lightweight Directory Access Protocol (LDAP) over Transport Layer Security (TLS) configuration.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// List that contains the validation messages related to the verification of the provided Lightweight Directory Access Protocol (LDAP) over Transport Layer Security (TLS) configuration details. The list contains a document for each test that MongoDB Cloud runs. MongoDB Cloud stops running tests after the first failure.
	// Read only field.
	Validations *[]LDAPVerifyConnectivityJobRequestValidation `json:"validations,omitempty"`
}

// NewLDAPVerifyConnectivityJobRequest instantiates a new LDAPVerifyConnectivityJobRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLDAPVerifyConnectivityJobRequest() *LDAPVerifyConnectivityJobRequest {
	this := LDAPVerifyConnectivityJobRequest{}
	return &this
}

// NewLDAPVerifyConnectivityJobRequestWithDefaults instantiates a new LDAPVerifyConnectivityJobRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLDAPVerifyConnectivityJobRequestWithDefaults() *LDAPVerifyConnectivityJobRequest {
	this := LDAPVerifyConnectivityJobRequest{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequest) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequest) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequest) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *LDAPVerifyConnectivityJobRequest) SetGroupId(v string) {
	o.GroupId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequest) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequest) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequest) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *LDAPVerifyConnectivityJobRequest) SetLinks(v []Link) {
	o.Links = &v
}

// GetRequest returns the Request field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequest) GetRequest() LDAPVerifyConnectivityJobRequestParams {
	if o == nil || IsNil(o.Request) {
		var ret LDAPVerifyConnectivityJobRequestParams
		return ret
	}
	return *o.Request
}

// GetRequestOk returns a tuple with the Request field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequest) GetRequestOk() (*LDAPVerifyConnectivityJobRequestParams, bool) {
	if o == nil || IsNil(o.Request) {
		return nil, false
	}

	return o.Request, true
}

// HasRequest returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequest) HasRequest() bool {
	if o != nil && !IsNil(o.Request) {
		return true
	}

	return false
}

// SetRequest gets a reference to the given LDAPVerifyConnectivityJobRequestParams and assigns it to the Request field.
func (o *LDAPVerifyConnectivityJobRequest) SetRequest(v LDAPVerifyConnectivityJobRequestParams) {
	o.Request = &v
}

// GetRequestId returns the RequestId field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequest) GetRequestId() string {
	if o == nil || IsNil(o.RequestId) {
		var ret string
		return ret
	}
	return *o.RequestId
}

// GetRequestIdOk returns a tuple with the RequestId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequest) GetRequestIdOk() (*string, bool) {
	if o == nil || IsNil(o.RequestId) {
		return nil, false
	}

	return o.RequestId, true
}

// HasRequestId returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequest) HasRequestId() bool {
	if o != nil && !IsNil(o.RequestId) {
		return true
	}

	return false
}

// SetRequestId gets a reference to the given string and assigns it to the RequestId field.
func (o *LDAPVerifyConnectivityJobRequest) SetRequestId(v string) {
	o.RequestId = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequest) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequest) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequest) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *LDAPVerifyConnectivityJobRequest) SetStatus(v string) {
	o.Status = &v
}

// GetValidations returns the Validations field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequest) GetValidations() []LDAPVerifyConnectivityJobRequestValidation {
	if o == nil || IsNil(o.Validations) {
		var ret []LDAPVerifyConnectivityJobRequestValidation
		return ret
	}
	return *o.Validations
}

// GetValidationsOk returns a tuple with the Validations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequest) GetValidationsOk() (*[]LDAPVerifyConnectivityJobRequestValidation, bool) {
	if o == nil || IsNil(o.Validations) {
		return nil, false
	}

	return o.Validations, true
}

// HasValidations returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequest) HasValidations() bool {
	if o != nil && !IsNil(o.Validations) {
		return true
	}

	return false
}

// SetValidations gets a reference to the given []LDAPVerifyConnectivityJobRequestValidation and assigns it to the Validations field.
func (o *LDAPVerifyConnectivityJobRequest) SetValidations(v []LDAPVerifyConnectivityJobRequestValidation) {
	o.Validations = &v
}
