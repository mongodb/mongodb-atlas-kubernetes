// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// EventViewForOrg struct for EventViewForOrg
type EventViewForOrg struct {
	// Unique 24-hexadecimal digit string that identifies the API Key that triggered the event. If this resource returns this parameter, it doesn't return the `userId` parameter.
	// Read only field.
	ApiKeyId *string `json:"apiKeyId,omitempty"`
	// Date and time when this event occurred. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Unique identifier of event type.
	EventTypeName *string `json:"eventTypeName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project in which the event occurred. The `eventId` identifies the specific event.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the event.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Flag that indicates whether a MongoDB employee triggered the specified event.
	// Read only field.
	IsGlobalAdmin *bool `json:"isGlobalAdmin,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization to which these events apply.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// Public part of the API key that triggered the event. If this resource returns this parameter, it doesn't return the **username** parameter.
	// Read only field.
	PublicKey *string `json:"publicKey,omitempty"`
	Raw       *Raw    `json:"raw,omitempty"`
	// IPv4 or IPv6 address from which the user triggered this event.
	// Read only field.
	RemoteAddress *string `json:"remoteAddress,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the console user who triggered the event. If this resource returns this parameter, it doesn't return the `apiKeyId` parameter.
	// Read only field.
	UserId *string `json:"userId,omitempty"`
	// Email address for the user who triggered this event. If this resource returns this parameter, it doesn't return the `publicApiKey` parameter.
	// Read only field.
	Username *string `json:"username,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the alert associated with the event.
	// Read only field.
	AlertId *string `json:"alertId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the alert configuration associated with the `alertId`.
	// Read only field.
	AlertConfigId *string `json:"alertConfigId,omitempty"`
	// Public part of the API key that this event targets.
	// Read only field.
	TargetPublicKey *string `json:"targetPublicKey,omitempty"`
	// Entry in the list of source host addresses that the API key accepts and this event targets.
	// Read only field.
	WhitelistEntry *string `json:"whitelistEntry,omitempty"`
	// Unique 24-hexadecimal digit string that identifies of the invoice associated with the event.
	// Read only field.
	InvoiceId *string `json:"invoiceId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the invoice payment associated with this event.
	// Read only field.
	PaymentId *string `json:"paymentId,omitempty"`
	// The username of the MongoDB User that was created, deleted, or edited.
	// Read only field.
	DbUserUsername *string `json:"dbUserUsername,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization team associated with this event.
	// Read only field.
	TeamId *string `json:"teamId,omitempty"`
	// Email address for the console user that this event targets. The resource returns this parameter when `\"eventTypeName\" : \"USER\"`.
	// Read only field.
	TargetUsername *string `json:"targetUsername,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the resource associated with the event.
	// Read only field.
	ResourceId *string `json:"resourceId,omitempty"`
	// Unique identifier of resource type.
	ResourceType *string `json:"resourceType,omitempty"`
	// Unique 24-hexadecimal character string that identifies the resource policy.
	// Read only field.
	ResourcePolicyId *string `json:"resourcePolicyId,omitempty"`
}

// NewEventViewForOrg instantiates a new EventViewForOrg object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEventViewForOrg() *EventViewForOrg {
	this := EventViewForOrg{}
	return &this
}

// NewEventViewForOrgWithDefaults instantiates a new EventViewForOrg object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEventViewForOrgWithDefaults() *EventViewForOrg {
	this := EventViewForOrg{}
	return &this
}

// GetApiKeyId returns the ApiKeyId field value if set, zero value otherwise
func (o *EventViewForOrg) GetApiKeyId() string {
	if o == nil || IsNil(o.ApiKeyId) {
		var ret string
		return ret
	}
	return *o.ApiKeyId
}

// GetApiKeyIdOk returns a tuple with the ApiKeyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetApiKeyIdOk() (*string, bool) {
	if o == nil || IsNil(o.ApiKeyId) {
		return nil, false
	}

	return o.ApiKeyId, true
}

// HasApiKeyId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasApiKeyId() bool {
	if o != nil && !IsNil(o.ApiKeyId) {
		return true
	}

	return false
}

// SetApiKeyId gets a reference to the given string and assigns it to the ApiKeyId field.
func (o *EventViewForOrg) SetApiKeyId(v string) {
	o.ApiKeyId = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *EventViewForOrg) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *EventViewForOrg) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *EventViewForOrg) SetCreated(v time.Time) {
	o.Created = &v
}

// GetEventTypeName returns the EventTypeName field value if set, zero value otherwise
func (o *EventViewForOrg) GetEventTypeName() string {
	if o == nil || IsNil(o.EventTypeName) {
		var ret string
		return ret
	}
	return *o.EventTypeName
}

// GetEventTypeNameOk returns a tuple with the EventTypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetEventTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.EventTypeName) {
		return nil, false
	}

	return o.EventTypeName, true
}

// HasEventTypeName returns a boolean if a field has been set.
func (o *EventViewForOrg) HasEventTypeName() bool {
	if o != nil && !IsNil(o.EventTypeName) {
		return true
	}

	return false
}

// SetEventTypeName gets a reference to the given string and assigns it to the EventTypeName field.
func (o *EventViewForOrg) SetEventTypeName(v string) {
	o.EventTypeName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *EventViewForOrg) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *EventViewForOrg) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *EventViewForOrg) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *EventViewForOrg) SetId(v string) {
	o.Id = &v
}

// GetIsGlobalAdmin returns the IsGlobalAdmin field value if set, zero value otherwise
func (o *EventViewForOrg) GetIsGlobalAdmin() bool {
	if o == nil || IsNil(o.IsGlobalAdmin) {
		var ret bool
		return ret
	}
	return *o.IsGlobalAdmin
}

// GetIsGlobalAdminOk returns a tuple with the IsGlobalAdmin field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetIsGlobalAdminOk() (*bool, bool) {
	if o == nil || IsNil(o.IsGlobalAdmin) {
		return nil, false
	}

	return o.IsGlobalAdmin, true
}

// HasIsGlobalAdmin returns a boolean if a field has been set.
func (o *EventViewForOrg) HasIsGlobalAdmin() bool {
	if o != nil && !IsNil(o.IsGlobalAdmin) {
		return true
	}

	return false
}

// SetIsGlobalAdmin gets a reference to the given bool and assigns it to the IsGlobalAdmin field.
func (o *EventViewForOrg) SetIsGlobalAdmin(v bool) {
	o.IsGlobalAdmin = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *EventViewForOrg) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *EventViewForOrg) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *EventViewForOrg) SetLinks(v []Link) {
	o.Links = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *EventViewForOrg) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *EventViewForOrg) SetOrgId(v string) {
	o.OrgId = &v
}

// GetPublicKey returns the PublicKey field value if set, zero value otherwise
func (o *EventViewForOrg) GetPublicKey() string {
	if o == nil || IsNil(o.PublicKey) {
		var ret string
		return ret
	}
	return *o.PublicKey
}

// GetPublicKeyOk returns a tuple with the PublicKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetPublicKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PublicKey) {
		return nil, false
	}

	return o.PublicKey, true
}

// HasPublicKey returns a boolean if a field has been set.
func (o *EventViewForOrg) HasPublicKey() bool {
	if o != nil && !IsNil(o.PublicKey) {
		return true
	}

	return false
}

// SetPublicKey gets a reference to the given string and assigns it to the PublicKey field.
func (o *EventViewForOrg) SetPublicKey(v string) {
	o.PublicKey = &v
}

// GetRaw returns the Raw field value if set, zero value otherwise
func (o *EventViewForOrg) GetRaw() Raw {
	if o == nil || IsNil(o.Raw) {
		var ret Raw
		return ret
	}
	return *o.Raw
}

// GetRawOk returns a tuple with the Raw field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetRawOk() (*Raw, bool) {
	if o == nil || IsNil(o.Raw) {
		return nil, false
	}

	return o.Raw, true
}

// HasRaw returns a boolean if a field has been set.
func (o *EventViewForOrg) HasRaw() bool {
	if o != nil && !IsNil(o.Raw) {
		return true
	}

	return false
}

// SetRaw gets a reference to the given Raw and assigns it to the Raw field.
func (o *EventViewForOrg) SetRaw(v Raw) {
	o.Raw = &v
}

// GetRemoteAddress returns the RemoteAddress field value if set, zero value otherwise
func (o *EventViewForOrg) GetRemoteAddress() string {
	if o == nil || IsNil(o.RemoteAddress) {
		var ret string
		return ret
	}
	return *o.RemoteAddress
}

// GetRemoteAddressOk returns a tuple with the RemoteAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetRemoteAddressOk() (*string, bool) {
	if o == nil || IsNil(o.RemoteAddress) {
		return nil, false
	}

	return o.RemoteAddress, true
}

// HasRemoteAddress returns a boolean if a field has been set.
func (o *EventViewForOrg) HasRemoteAddress() bool {
	if o != nil && !IsNil(o.RemoteAddress) {
		return true
	}

	return false
}

// SetRemoteAddress gets a reference to the given string and assigns it to the RemoteAddress field.
func (o *EventViewForOrg) SetRemoteAddress(v string) {
	o.RemoteAddress = &v
}

// GetUserId returns the UserId field value if set, zero value otherwise
func (o *EventViewForOrg) GetUserId() string {
	if o == nil || IsNil(o.UserId) {
		var ret string
		return ret
	}
	return *o.UserId
}

// GetUserIdOk returns a tuple with the UserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.UserId) {
		return nil, false
	}

	return o.UserId, true
}

// HasUserId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasUserId() bool {
	if o != nil && !IsNil(o.UserId) {
		return true
	}

	return false
}

// SetUserId gets a reference to the given string and assigns it to the UserId field.
func (o *EventViewForOrg) SetUserId(v string) {
	o.UserId = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *EventViewForOrg) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *EventViewForOrg) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *EventViewForOrg) SetUsername(v string) {
	o.Username = &v
}

// GetAlertId returns the AlertId field value if set, zero value otherwise
func (o *EventViewForOrg) GetAlertId() string {
	if o == nil || IsNil(o.AlertId) {
		var ret string
		return ret
	}
	return *o.AlertId
}

// GetAlertIdOk returns a tuple with the AlertId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetAlertIdOk() (*string, bool) {
	if o == nil || IsNil(o.AlertId) {
		return nil, false
	}

	return o.AlertId, true
}

// HasAlertId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasAlertId() bool {
	if o != nil && !IsNil(o.AlertId) {
		return true
	}

	return false
}

// SetAlertId gets a reference to the given string and assigns it to the AlertId field.
func (o *EventViewForOrg) SetAlertId(v string) {
	o.AlertId = &v
}

// GetAlertConfigId returns the AlertConfigId field value if set, zero value otherwise
func (o *EventViewForOrg) GetAlertConfigId() string {
	if o == nil || IsNil(o.AlertConfigId) {
		var ret string
		return ret
	}
	return *o.AlertConfigId
}

// GetAlertConfigIdOk returns a tuple with the AlertConfigId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetAlertConfigIdOk() (*string, bool) {
	if o == nil || IsNil(o.AlertConfigId) {
		return nil, false
	}

	return o.AlertConfigId, true
}

// HasAlertConfigId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasAlertConfigId() bool {
	if o != nil && !IsNil(o.AlertConfigId) {
		return true
	}

	return false
}

// SetAlertConfigId gets a reference to the given string and assigns it to the AlertConfigId field.
func (o *EventViewForOrg) SetAlertConfigId(v string) {
	o.AlertConfigId = &v
}

// GetTargetPublicKey returns the TargetPublicKey field value if set, zero value otherwise
func (o *EventViewForOrg) GetTargetPublicKey() string {
	if o == nil || IsNil(o.TargetPublicKey) {
		var ret string
		return ret
	}
	return *o.TargetPublicKey
}

// GetTargetPublicKeyOk returns a tuple with the TargetPublicKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetTargetPublicKeyOk() (*string, bool) {
	if o == nil || IsNil(o.TargetPublicKey) {
		return nil, false
	}

	return o.TargetPublicKey, true
}

// HasTargetPublicKey returns a boolean if a field has been set.
func (o *EventViewForOrg) HasTargetPublicKey() bool {
	if o != nil && !IsNil(o.TargetPublicKey) {
		return true
	}

	return false
}

// SetTargetPublicKey gets a reference to the given string and assigns it to the TargetPublicKey field.
func (o *EventViewForOrg) SetTargetPublicKey(v string) {
	o.TargetPublicKey = &v
}

// GetWhitelistEntry returns the WhitelistEntry field value if set, zero value otherwise
func (o *EventViewForOrg) GetWhitelistEntry() string {
	if o == nil || IsNil(o.WhitelistEntry) {
		var ret string
		return ret
	}
	return *o.WhitelistEntry
}

// GetWhitelistEntryOk returns a tuple with the WhitelistEntry field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetWhitelistEntryOk() (*string, bool) {
	if o == nil || IsNil(o.WhitelistEntry) {
		return nil, false
	}

	return o.WhitelistEntry, true
}

// HasWhitelistEntry returns a boolean if a field has been set.
func (o *EventViewForOrg) HasWhitelistEntry() bool {
	if o != nil && !IsNil(o.WhitelistEntry) {
		return true
	}

	return false
}

// SetWhitelistEntry gets a reference to the given string and assigns it to the WhitelistEntry field.
func (o *EventViewForOrg) SetWhitelistEntry(v string) {
	o.WhitelistEntry = &v
}

// GetInvoiceId returns the InvoiceId field value if set, zero value otherwise
func (o *EventViewForOrg) GetInvoiceId() string {
	if o == nil || IsNil(o.InvoiceId) {
		var ret string
		return ret
	}
	return *o.InvoiceId
}

// GetInvoiceIdOk returns a tuple with the InvoiceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetInvoiceIdOk() (*string, bool) {
	if o == nil || IsNil(o.InvoiceId) {
		return nil, false
	}

	return o.InvoiceId, true
}

// HasInvoiceId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasInvoiceId() bool {
	if o != nil && !IsNil(o.InvoiceId) {
		return true
	}

	return false
}

// SetInvoiceId gets a reference to the given string and assigns it to the InvoiceId field.
func (o *EventViewForOrg) SetInvoiceId(v string) {
	o.InvoiceId = &v
}

// GetPaymentId returns the PaymentId field value if set, zero value otherwise
func (o *EventViewForOrg) GetPaymentId() string {
	if o == nil || IsNil(o.PaymentId) {
		var ret string
		return ret
	}
	return *o.PaymentId
}

// GetPaymentIdOk returns a tuple with the PaymentId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetPaymentIdOk() (*string, bool) {
	if o == nil || IsNil(o.PaymentId) {
		return nil, false
	}

	return o.PaymentId, true
}

// HasPaymentId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasPaymentId() bool {
	if o != nil && !IsNil(o.PaymentId) {
		return true
	}

	return false
}

// SetPaymentId gets a reference to the given string and assigns it to the PaymentId field.
func (o *EventViewForOrg) SetPaymentId(v string) {
	o.PaymentId = &v
}

// GetDbUserUsername returns the DbUserUsername field value if set, zero value otherwise
func (o *EventViewForOrg) GetDbUserUsername() string {
	if o == nil || IsNil(o.DbUserUsername) {
		var ret string
		return ret
	}
	return *o.DbUserUsername
}

// GetDbUserUsernameOk returns a tuple with the DbUserUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetDbUserUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.DbUserUsername) {
		return nil, false
	}

	return o.DbUserUsername, true
}

// HasDbUserUsername returns a boolean if a field has been set.
func (o *EventViewForOrg) HasDbUserUsername() bool {
	if o != nil && !IsNil(o.DbUserUsername) {
		return true
	}

	return false
}

// SetDbUserUsername gets a reference to the given string and assigns it to the DbUserUsername field.
func (o *EventViewForOrg) SetDbUserUsername(v string) {
	o.DbUserUsername = &v
}

// GetTeamId returns the TeamId field value if set, zero value otherwise
func (o *EventViewForOrg) GetTeamId() string {
	if o == nil || IsNil(o.TeamId) {
		var ret string
		return ret
	}
	return *o.TeamId
}

// GetTeamIdOk returns a tuple with the TeamId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetTeamIdOk() (*string, bool) {
	if o == nil || IsNil(o.TeamId) {
		return nil, false
	}

	return o.TeamId, true
}

// HasTeamId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasTeamId() bool {
	if o != nil && !IsNil(o.TeamId) {
		return true
	}

	return false
}

// SetTeamId gets a reference to the given string and assigns it to the TeamId field.
func (o *EventViewForOrg) SetTeamId(v string) {
	o.TeamId = &v
}

// GetTargetUsername returns the TargetUsername field value if set, zero value otherwise
func (o *EventViewForOrg) GetTargetUsername() string {
	if o == nil || IsNil(o.TargetUsername) {
		var ret string
		return ret
	}
	return *o.TargetUsername
}

// GetTargetUsernameOk returns a tuple with the TargetUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetTargetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.TargetUsername) {
		return nil, false
	}

	return o.TargetUsername, true
}

// HasTargetUsername returns a boolean if a field has been set.
func (o *EventViewForOrg) HasTargetUsername() bool {
	if o != nil && !IsNil(o.TargetUsername) {
		return true
	}

	return false
}

// SetTargetUsername gets a reference to the given string and assigns it to the TargetUsername field.
func (o *EventViewForOrg) SetTargetUsername(v string) {
	o.TargetUsername = &v
}

// GetResourceId returns the ResourceId field value if set, zero value otherwise
func (o *EventViewForOrg) GetResourceId() string {
	if o == nil || IsNil(o.ResourceId) {
		var ret string
		return ret
	}
	return *o.ResourceId
}

// GetResourceIdOk returns a tuple with the ResourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetResourceIdOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceId) {
		return nil, false
	}

	return o.ResourceId, true
}

// HasResourceId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasResourceId() bool {
	if o != nil && !IsNil(o.ResourceId) {
		return true
	}

	return false
}

// SetResourceId gets a reference to the given string and assigns it to the ResourceId field.
func (o *EventViewForOrg) SetResourceId(v string) {
	o.ResourceId = &v
}

// GetResourceType returns the ResourceType field value if set, zero value otherwise
func (o *EventViewForOrg) GetResourceType() string {
	if o == nil || IsNil(o.ResourceType) {
		var ret string
		return ret
	}
	return *o.ResourceType
}

// GetResourceTypeOk returns a tuple with the ResourceType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetResourceTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceType) {
		return nil, false
	}

	return o.ResourceType, true
}

// HasResourceType returns a boolean if a field has been set.
func (o *EventViewForOrg) HasResourceType() bool {
	if o != nil && !IsNil(o.ResourceType) {
		return true
	}

	return false
}

// SetResourceType gets a reference to the given string and assigns it to the ResourceType field.
func (o *EventViewForOrg) SetResourceType(v string) {
	o.ResourceType = &v
}

// GetResourcePolicyId returns the ResourcePolicyId field value if set, zero value otherwise
func (o *EventViewForOrg) GetResourcePolicyId() string {
	if o == nil || IsNil(o.ResourcePolicyId) {
		var ret string
		return ret
	}
	return *o.ResourcePolicyId
}

// GetResourcePolicyIdOk returns a tuple with the ResourcePolicyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForOrg) GetResourcePolicyIdOk() (*string, bool) {
	if o == nil || IsNil(o.ResourcePolicyId) {
		return nil, false
	}

	return o.ResourcePolicyId, true
}

// HasResourcePolicyId returns a boolean if a field has been set.
func (o *EventViewForOrg) HasResourcePolicyId() bool {
	if o != nil && !IsNil(o.ResourcePolicyId) {
		return true
	}

	return false
}

// SetResourcePolicyId gets a reference to the given string and assigns it to the ResourcePolicyId field.
func (o *EventViewForOrg) SetResourcePolicyId(v string) {
	o.ResourcePolicyId = &v
}
