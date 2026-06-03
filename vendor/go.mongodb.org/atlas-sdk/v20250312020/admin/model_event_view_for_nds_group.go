// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// EventViewForNdsGroup struct for EventViewForNdsGroup
type EventViewForNdsGroup struct {
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
	// Human-readable label of the shard associated with the event.
	// Read only field.
	ShardName *string `json:"shardName,omitempty"`
	// Human-readable label of the collection on which the event occurred. The resource returns this parameter when the `eventTypeName` includes `DATA_EXPLORER`.
	// Read only field.
	Collection *string `json:"collection,omitempty"`
	// Human-readable label of the database on which this incident occurred. The resource returns this parameter when `\"eventTypeName\" : \"DATA_EXPLORER\"` or `\"eventTypeName\" : \"DATA_EXPLORER_CRUD\"`.
	// Read only field.
	Database *string `json:"database,omitempty"`
	// Action that the database attempted to execute when the event triggered. The response returns this parameter when `eventTypeName\" : \"DATA_EXPLORER\"`.
	// Read only field.
	OpType *string `json:"opType,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the Data Explorer session associated with the event.
	// Read only field.
	SessionId *string `json:"sessionId,omitempty"`
	// Desk location of MongoDB employee associated with the event.
	// Read only field.
	DeskLocation *string `json:"deskLocation,omitempty"`
	// Identifier of MongoDB employee associated with the event.
	// Read only field.
	EmployeeIdentifier *string `json:"employeeIdentifier,omitempty"`
	// IANA port on which the MongoDB process listens for requests.
	// Read only field.
	Port *int `json:"port,omitempty"`
	// Human-readable label of the replica set associated with the event.
	// Read only field.
	ReplicaSetName *string            `json:"replicaSetName,omitempty"`
	CurrentValue   *NumberMetricValue `json:"currentValue,omitempty"`
	// Human-readable label of the metric associated with the `alertId`. This field may change type of `currentValue` field.
	// Read only field.
	MetricName *string `json:"metricName,omitempty"`
	// The username of the MongoDB User that was created, deleted, or edited.
	// Read only field.
	DbUserUsername *string `json:"dbUserUsername,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the endpoint associated with this event.
	// Read only field.
	EndpointId *string `json:"endpointId,omitempty"`
	// Unique identification string that the cloud provider uses to identify the private endpoint.
	// Read only field.
	ProviderEndpointId *string `json:"providerEndpointId,omitempty"`
	// Fully qualified domain name (FQDN) of the host associated with the event.
	// Read only field.
	Hostname *string `json:"hostname,omitempty"`
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
	// Name of the stream processing workspace associated with the event.
	// Read only field.
	InstanceName *string `json:"instanceName,omitempty"`
	// Error message linked to the stream processor associated with the event.
	// Read only field.
	ProcessorErrorMsg *string `json:"processorErrorMsg,omitempty"`
	// Name of the stream processor associated with the event.
	// Read only field.
	ProcessorName *string `json:"processorName,omitempty"`
	// State of the stream processor associated with the event.
	// Read only field.
	ProcessorState *string `json:"processorState,omitempty"`
	// Unique 24-hexadecimal character string that identifies the resource policy.
	// Read only field.
	ResourcePolicyId *string `json:"resourcePolicyId,omitempty"`
	// String representation of the violated resource policy ids.
	ViolatedPolicies *[]string `json:"violatedPolicies,omitempty"`
	// Resource policy action taken by the user and evaluated against the currently active policies.
	// Read only field.
	ViolationAction *string `json:"violationAction,omitempty"`
}

// NewEventViewForNdsGroup instantiates a new EventViewForNdsGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEventViewForNdsGroup() *EventViewForNdsGroup {
	this := EventViewForNdsGroup{}
	return &this
}

// NewEventViewForNdsGroupWithDefaults instantiates a new EventViewForNdsGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEventViewForNdsGroupWithDefaults() *EventViewForNdsGroup {
	this := EventViewForNdsGroup{}
	return &this
}

// GetApiKeyId returns the ApiKeyId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetApiKeyId() string {
	if o == nil || IsNil(o.ApiKeyId) {
		var ret string
		return ret
	}
	return *o.ApiKeyId
}

// GetApiKeyIdOk returns a tuple with the ApiKeyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetApiKeyIdOk() (*string, bool) {
	if o == nil || IsNil(o.ApiKeyId) {
		return nil, false
	}

	return o.ApiKeyId, true
}

// HasApiKeyId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasApiKeyId() bool {
	if o != nil && !IsNil(o.ApiKeyId) {
		return true
	}

	return false
}

// SetApiKeyId gets a reference to the given string and assigns it to the ApiKeyId field.
func (o *EventViewForNdsGroup) SetApiKeyId(v string) {
	o.ApiKeyId = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *EventViewForNdsGroup) SetCreated(v time.Time) {
	o.Created = &v
}

// GetEventTypeName returns the EventTypeName field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetEventTypeName() string {
	if o == nil || IsNil(o.EventTypeName) {
		var ret string
		return ret
	}
	return *o.EventTypeName
}

// GetEventTypeNameOk returns a tuple with the EventTypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetEventTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.EventTypeName) {
		return nil, false
	}

	return o.EventTypeName, true
}

// HasEventTypeName returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasEventTypeName() bool {
	if o != nil && !IsNil(o.EventTypeName) {
		return true
	}

	return false
}

// SetEventTypeName gets a reference to the given string and assigns it to the EventTypeName field.
func (o *EventViewForNdsGroup) SetEventTypeName(v string) {
	o.EventTypeName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *EventViewForNdsGroup) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *EventViewForNdsGroup) SetId(v string) {
	o.Id = &v
}

// GetIsGlobalAdmin returns the IsGlobalAdmin field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetIsGlobalAdmin() bool {
	if o == nil || IsNil(o.IsGlobalAdmin) {
		var ret bool
		return ret
	}
	return *o.IsGlobalAdmin
}

// GetIsGlobalAdminOk returns a tuple with the IsGlobalAdmin field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetIsGlobalAdminOk() (*bool, bool) {
	if o == nil || IsNil(o.IsGlobalAdmin) {
		return nil, false
	}

	return o.IsGlobalAdmin, true
}

// HasIsGlobalAdmin returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasIsGlobalAdmin() bool {
	if o != nil && !IsNil(o.IsGlobalAdmin) {
		return true
	}

	return false
}

// SetIsGlobalAdmin gets a reference to the given bool and assigns it to the IsGlobalAdmin field.
func (o *EventViewForNdsGroup) SetIsGlobalAdmin(v bool) {
	o.IsGlobalAdmin = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *EventViewForNdsGroup) SetLinks(v []Link) {
	o.Links = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *EventViewForNdsGroup) SetOrgId(v string) {
	o.OrgId = &v
}

// GetPublicKey returns the PublicKey field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetPublicKey() string {
	if o == nil || IsNil(o.PublicKey) {
		var ret string
		return ret
	}
	return *o.PublicKey
}

// GetPublicKeyOk returns a tuple with the PublicKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetPublicKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PublicKey) {
		return nil, false
	}

	return o.PublicKey, true
}

// HasPublicKey returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasPublicKey() bool {
	if o != nil && !IsNil(o.PublicKey) {
		return true
	}

	return false
}

// SetPublicKey gets a reference to the given string and assigns it to the PublicKey field.
func (o *EventViewForNdsGroup) SetPublicKey(v string) {
	o.PublicKey = &v
}

// GetRaw returns the Raw field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetRaw() Raw {
	if o == nil || IsNil(o.Raw) {
		var ret Raw
		return ret
	}
	return *o.Raw
}

// GetRawOk returns a tuple with the Raw field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetRawOk() (*Raw, bool) {
	if o == nil || IsNil(o.Raw) {
		return nil, false
	}

	return o.Raw, true
}

// HasRaw returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasRaw() bool {
	if o != nil && !IsNil(o.Raw) {
		return true
	}

	return false
}

// SetRaw gets a reference to the given Raw and assigns it to the Raw field.
func (o *EventViewForNdsGroup) SetRaw(v Raw) {
	o.Raw = &v
}

// GetRemoteAddress returns the RemoteAddress field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetRemoteAddress() string {
	if o == nil || IsNil(o.RemoteAddress) {
		var ret string
		return ret
	}
	return *o.RemoteAddress
}

// GetRemoteAddressOk returns a tuple with the RemoteAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetRemoteAddressOk() (*string, bool) {
	if o == nil || IsNil(o.RemoteAddress) {
		return nil, false
	}

	return o.RemoteAddress, true
}

// HasRemoteAddress returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasRemoteAddress() bool {
	if o != nil && !IsNil(o.RemoteAddress) {
		return true
	}

	return false
}

// SetRemoteAddress gets a reference to the given string and assigns it to the RemoteAddress field.
func (o *EventViewForNdsGroup) SetRemoteAddress(v string) {
	o.RemoteAddress = &v
}

// GetUserId returns the UserId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetUserId() string {
	if o == nil || IsNil(o.UserId) {
		var ret string
		return ret
	}
	return *o.UserId
}

// GetUserIdOk returns a tuple with the UserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.UserId) {
		return nil, false
	}

	return o.UserId, true
}

// HasUserId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasUserId() bool {
	if o != nil && !IsNil(o.UserId) {
		return true
	}

	return false
}

// SetUserId gets a reference to the given string and assigns it to the UserId field.
func (o *EventViewForNdsGroup) SetUserId(v string) {
	o.UserId = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *EventViewForNdsGroup) SetUsername(v string) {
	o.Username = &v
}

// GetAlertId returns the AlertId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetAlertId() string {
	if o == nil || IsNil(o.AlertId) {
		var ret string
		return ret
	}
	return *o.AlertId
}

// GetAlertIdOk returns a tuple with the AlertId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetAlertIdOk() (*string, bool) {
	if o == nil || IsNil(o.AlertId) {
		return nil, false
	}

	return o.AlertId, true
}

// HasAlertId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasAlertId() bool {
	if o != nil && !IsNil(o.AlertId) {
		return true
	}

	return false
}

// SetAlertId gets a reference to the given string and assigns it to the AlertId field.
func (o *EventViewForNdsGroup) SetAlertId(v string) {
	o.AlertId = &v
}

// GetAlertConfigId returns the AlertConfigId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetAlertConfigId() string {
	if o == nil || IsNil(o.AlertConfigId) {
		var ret string
		return ret
	}
	return *o.AlertConfigId
}

// GetAlertConfigIdOk returns a tuple with the AlertConfigId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetAlertConfigIdOk() (*string, bool) {
	if o == nil || IsNil(o.AlertConfigId) {
		return nil, false
	}

	return o.AlertConfigId, true
}

// HasAlertConfigId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasAlertConfigId() bool {
	if o != nil && !IsNil(o.AlertConfigId) {
		return true
	}

	return false
}

// SetAlertConfigId gets a reference to the given string and assigns it to the AlertConfigId field.
func (o *EventViewForNdsGroup) SetAlertConfigId(v string) {
	o.AlertConfigId = &v
}

// GetTargetPublicKey returns the TargetPublicKey field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetTargetPublicKey() string {
	if o == nil || IsNil(o.TargetPublicKey) {
		var ret string
		return ret
	}
	return *o.TargetPublicKey
}

// GetTargetPublicKeyOk returns a tuple with the TargetPublicKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetTargetPublicKeyOk() (*string, bool) {
	if o == nil || IsNil(o.TargetPublicKey) {
		return nil, false
	}

	return o.TargetPublicKey, true
}

// HasTargetPublicKey returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasTargetPublicKey() bool {
	if o != nil && !IsNil(o.TargetPublicKey) {
		return true
	}

	return false
}

// SetTargetPublicKey gets a reference to the given string and assigns it to the TargetPublicKey field.
func (o *EventViewForNdsGroup) SetTargetPublicKey(v string) {
	o.TargetPublicKey = &v
}

// GetWhitelistEntry returns the WhitelistEntry field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetWhitelistEntry() string {
	if o == nil || IsNil(o.WhitelistEntry) {
		var ret string
		return ret
	}
	return *o.WhitelistEntry
}

// GetWhitelistEntryOk returns a tuple with the WhitelistEntry field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetWhitelistEntryOk() (*string, bool) {
	if o == nil || IsNil(o.WhitelistEntry) {
		return nil, false
	}

	return o.WhitelistEntry, true
}

// HasWhitelistEntry returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasWhitelistEntry() bool {
	if o != nil && !IsNil(o.WhitelistEntry) {
		return true
	}

	return false
}

// SetWhitelistEntry gets a reference to the given string and assigns it to the WhitelistEntry field.
func (o *EventViewForNdsGroup) SetWhitelistEntry(v string) {
	o.WhitelistEntry = &v
}

// GetInvoiceId returns the InvoiceId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetInvoiceId() string {
	if o == nil || IsNil(o.InvoiceId) {
		var ret string
		return ret
	}
	return *o.InvoiceId
}

// GetInvoiceIdOk returns a tuple with the InvoiceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetInvoiceIdOk() (*string, bool) {
	if o == nil || IsNil(o.InvoiceId) {
		return nil, false
	}

	return o.InvoiceId, true
}

// HasInvoiceId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasInvoiceId() bool {
	if o != nil && !IsNil(o.InvoiceId) {
		return true
	}

	return false
}

// SetInvoiceId gets a reference to the given string and assigns it to the InvoiceId field.
func (o *EventViewForNdsGroup) SetInvoiceId(v string) {
	o.InvoiceId = &v
}

// GetPaymentId returns the PaymentId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetPaymentId() string {
	if o == nil || IsNil(o.PaymentId) {
		var ret string
		return ret
	}
	return *o.PaymentId
}

// GetPaymentIdOk returns a tuple with the PaymentId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetPaymentIdOk() (*string, bool) {
	if o == nil || IsNil(o.PaymentId) {
		return nil, false
	}

	return o.PaymentId, true
}

// HasPaymentId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasPaymentId() bool {
	if o != nil && !IsNil(o.PaymentId) {
		return true
	}

	return false
}

// SetPaymentId gets a reference to the given string and assigns it to the PaymentId field.
func (o *EventViewForNdsGroup) SetPaymentId(v string) {
	o.PaymentId = &v
}

// GetShardName returns the ShardName field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetShardName() string {
	if o == nil || IsNil(o.ShardName) {
		var ret string
		return ret
	}
	return *o.ShardName
}

// GetShardNameOk returns a tuple with the ShardName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetShardNameOk() (*string, bool) {
	if o == nil || IsNil(o.ShardName) {
		return nil, false
	}

	return o.ShardName, true
}

// HasShardName returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasShardName() bool {
	if o != nil && !IsNil(o.ShardName) {
		return true
	}

	return false
}

// SetShardName gets a reference to the given string and assigns it to the ShardName field.
func (o *EventViewForNdsGroup) SetShardName(v string) {
	o.ShardName = &v
}

// GetCollection returns the Collection field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetCollection() string {
	if o == nil || IsNil(o.Collection) {
		var ret string
		return ret
	}
	return *o.Collection
}

// GetCollectionOk returns a tuple with the Collection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetCollectionOk() (*string, bool) {
	if o == nil || IsNil(o.Collection) {
		return nil, false
	}

	return o.Collection, true
}

// HasCollection returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasCollection() bool {
	if o != nil && !IsNil(o.Collection) {
		return true
	}

	return false
}

// SetCollection gets a reference to the given string and assigns it to the Collection field.
func (o *EventViewForNdsGroup) SetCollection(v string) {
	o.Collection = &v
}

// GetDatabase returns the Database field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetDatabase() string {
	if o == nil || IsNil(o.Database) {
		var ret string
		return ret
	}
	return *o.Database
}

// GetDatabaseOk returns a tuple with the Database field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetDatabaseOk() (*string, bool) {
	if o == nil || IsNil(o.Database) {
		return nil, false
	}

	return o.Database, true
}

// HasDatabase returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasDatabase() bool {
	if o != nil && !IsNil(o.Database) {
		return true
	}

	return false
}

// SetDatabase gets a reference to the given string and assigns it to the Database field.
func (o *EventViewForNdsGroup) SetDatabase(v string) {
	o.Database = &v
}

// GetOpType returns the OpType field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetOpType() string {
	if o == nil || IsNil(o.OpType) {
		var ret string
		return ret
	}
	return *o.OpType
}

// GetOpTypeOk returns a tuple with the OpType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetOpTypeOk() (*string, bool) {
	if o == nil || IsNil(o.OpType) {
		return nil, false
	}

	return o.OpType, true
}

// HasOpType returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasOpType() bool {
	if o != nil && !IsNil(o.OpType) {
		return true
	}

	return false
}

// SetOpType gets a reference to the given string and assigns it to the OpType field.
func (o *EventViewForNdsGroup) SetOpType(v string) {
	o.OpType = &v
}

// GetSessionId returns the SessionId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetSessionId() string {
	if o == nil || IsNil(o.SessionId) {
		var ret string
		return ret
	}
	return *o.SessionId
}

// GetSessionIdOk returns a tuple with the SessionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetSessionIdOk() (*string, bool) {
	if o == nil || IsNil(o.SessionId) {
		return nil, false
	}

	return o.SessionId, true
}

// HasSessionId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasSessionId() bool {
	if o != nil && !IsNil(o.SessionId) {
		return true
	}

	return false
}

// SetSessionId gets a reference to the given string and assigns it to the SessionId field.
func (o *EventViewForNdsGroup) SetSessionId(v string) {
	o.SessionId = &v
}

// GetDeskLocation returns the DeskLocation field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetDeskLocation() string {
	if o == nil || IsNil(o.DeskLocation) {
		var ret string
		return ret
	}
	return *o.DeskLocation
}

// GetDeskLocationOk returns a tuple with the DeskLocation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetDeskLocationOk() (*string, bool) {
	if o == nil || IsNil(o.DeskLocation) {
		return nil, false
	}

	return o.DeskLocation, true
}

// HasDeskLocation returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasDeskLocation() bool {
	if o != nil && !IsNil(o.DeskLocation) {
		return true
	}

	return false
}

// SetDeskLocation gets a reference to the given string and assigns it to the DeskLocation field.
func (o *EventViewForNdsGroup) SetDeskLocation(v string) {
	o.DeskLocation = &v
}

// GetEmployeeIdentifier returns the EmployeeIdentifier field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetEmployeeIdentifier() string {
	if o == nil || IsNil(o.EmployeeIdentifier) {
		var ret string
		return ret
	}
	return *o.EmployeeIdentifier
}

// GetEmployeeIdentifierOk returns a tuple with the EmployeeIdentifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetEmployeeIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.EmployeeIdentifier) {
		return nil, false
	}

	return o.EmployeeIdentifier, true
}

// HasEmployeeIdentifier returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasEmployeeIdentifier() bool {
	if o != nil && !IsNil(o.EmployeeIdentifier) {
		return true
	}

	return false
}

// SetEmployeeIdentifier gets a reference to the given string and assigns it to the EmployeeIdentifier field.
func (o *EventViewForNdsGroup) SetEmployeeIdentifier(v string) {
	o.EmployeeIdentifier = &v
}

// GetPort returns the Port field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetPort() int {
	if o == nil || IsNil(o.Port) {
		var ret int
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetPortOk() (*int, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}

	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int and assigns it to the Port field.
func (o *EventViewForNdsGroup) SetPort(v int) {
	o.Port = &v
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *EventViewForNdsGroup) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}

// GetCurrentValue returns the CurrentValue field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetCurrentValue() NumberMetricValue {
	if o == nil || IsNil(o.CurrentValue) {
		var ret NumberMetricValue
		return ret
	}
	return *o.CurrentValue
}

// GetCurrentValueOk returns a tuple with the CurrentValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetCurrentValueOk() (*NumberMetricValue, bool) {
	if o == nil || IsNil(o.CurrentValue) {
		return nil, false
	}

	return o.CurrentValue, true
}

// HasCurrentValue returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasCurrentValue() bool {
	if o != nil && !IsNil(o.CurrentValue) {
		return true
	}

	return false
}

// SetCurrentValue gets a reference to the given NumberMetricValue and assigns it to the CurrentValue field.
func (o *EventViewForNdsGroup) SetCurrentValue(v NumberMetricValue) {
	o.CurrentValue = &v
}

// GetMetricName returns the MetricName field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetMetricName() string {
	if o == nil || IsNil(o.MetricName) {
		var ret string
		return ret
	}
	return *o.MetricName
}

// GetMetricNameOk returns a tuple with the MetricName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetMetricNameOk() (*string, bool) {
	if o == nil || IsNil(o.MetricName) {
		return nil, false
	}

	return o.MetricName, true
}

// HasMetricName returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasMetricName() bool {
	if o != nil && !IsNil(o.MetricName) {
		return true
	}

	return false
}

// SetMetricName gets a reference to the given string and assigns it to the MetricName field.
func (o *EventViewForNdsGroup) SetMetricName(v string) {
	o.MetricName = &v
}

// GetDbUserUsername returns the DbUserUsername field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetDbUserUsername() string {
	if o == nil || IsNil(o.DbUserUsername) {
		var ret string
		return ret
	}
	return *o.DbUserUsername
}

// GetDbUserUsernameOk returns a tuple with the DbUserUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetDbUserUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.DbUserUsername) {
		return nil, false
	}

	return o.DbUserUsername, true
}

// HasDbUserUsername returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasDbUserUsername() bool {
	if o != nil && !IsNil(o.DbUserUsername) {
		return true
	}

	return false
}

// SetDbUserUsername gets a reference to the given string and assigns it to the DbUserUsername field.
func (o *EventViewForNdsGroup) SetDbUserUsername(v string) {
	o.DbUserUsername = &v
}

// GetEndpointId returns the EndpointId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetEndpointId() string {
	if o == nil || IsNil(o.EndpointId) {
		var ret string
		return ret
	}
	return *o.EndpointId
}

// GetEndpointIdOk returns a tuple with the EndpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetEndpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointId) {
		return nil, false
	}

	return o.EndpointId, true
}

// HasEndpointId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasEndpointId() bool {
	if o != nil && !IsNil(o.EndpointId) {
		return true
	}

	return false
}

// SetEndpointId gets a reference to the given string and assigns it to the EndpointId field.
func (o *EventViewForNdsGroup) SetEndpointId(v string) {
	o.EndpointId = &v
}

// GetProviderEndpointId returns the ProviderEndpointId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetProviderEndpointId() string {
	if o == nil || IsNil(o.ProviderEndpointId) {
		var ret string
		return ret
	}
	return *o.ProviderEndpointId
}

// GetProviderEndpointIdOk returns a tuple with the ProviderEndpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetProviderEndpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderEndpointId) {
		return nil, false
	}

	return o.ProviderEndpointId, true
}

// HasProviderEndpointId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasProviderEndpointId() bool {
	if o != nil && !IsNil(o.ProviderEndpointId) {
		return true
	}

	return false
}

// SetProviderEndpointId gets a reference to the given string and assigns it to the ProviderEndpointId field.
func (o *EventViewForNdsGroup) SetProviderEndpointId(v string) {
	o.ProviderEndpointId = &v
}

// GetHostname returns the Hostname field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetHostname() string {
	if o == nil || IsNil(o.Hostname) {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetHostnameOk() (*string, bool) {
	if o == nil || IsNil(o.Hostname) {
		return nil, false
	}

	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasHostname() bool {
	if o != nil && !IsNil(o.Hostname) {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *EventViewForNdsGroup) SetHostname(v string) {
	o.Hostname = &v
}

// GetTeamId returns the TeamId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetTeamId() string {
	if o == nil || IsNil(o.TeamId) {
		var ret string
		return ret
	}
	return *o.TeamId
}

// GetTeamIdOk returns a tuple with the TeamId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetTeamIdOk() (*string, bool) {
	if o == nil || IsNil(o.TeamId) {
		return nil, false
	}

	return o.TeamId, true
}

// HasTeamId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasTeamId() bool {
	if o != nil && !IsNil(o.TeamId) {
		return true
	}

	return false
}

// SetTeamId gets a reference to the given string and assigns it to the TeamId field.
func (o *EventViewForNdsGroup) SetTeamId(v string) {
	o.TeamId = &v
}

// GetTargetUsername returns the TargetUsername field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetTargetUsername() string {
	if o == nil || IsNil(o.TargetUsername) {
		var ret string
		return ret
	}
	return *o.TargetUsername
}

// GetTargetUsernameOk returns a tuple with the TargetUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetTargetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.TargetUsername) {
		return nil, false
	}

	return o.TargetUsername, true
}

// HasTargetUsername returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasTargetUsername() bool {
	if o != nil && !IsNil(o.TargetUsername) {
		return true
	}

	return false
}

// SetTargetUsername gets a reference to the given string and assigns it to the TargetUsername field.
func (o *EventViewForNdsGroup) SetTargetUsername(v string) {
	o.TargetUsername = &v
}

// GetResourceId returns the ResourceId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetResourceId() string {
	if o == nil || IsNil(o.ResourceId) {
		var ret string
		return ret
	}
	return *o.ResourceId
}

// GetResourceIdOk returns a tuple with the ResourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetResourceIdOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceId) {
		return nil, false
	}

	return o.ResourceId, true
}

// HasResourceId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasResourceId() bool {
	if o != nil && !IsNil(o.ResourceId) {
		return true
	}

	return false
}

// SetResourceId gets a reference to the given string and assigns it to the ResourceId field.
func (o *EventViewForNdsGroup) SetResourceId(v string) {
	o.ResourceId = &v
}

// GetResourceType returns the ResourceType field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetResourceType() string {
	if o == nil || IsNil(o.ResourceType) {
		var ret string
		return ret
	}
	return *o.ResourceType
}

// GetResourceTypeOk returns a tuple with the ResourceType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetResourceTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceType) {
		return nil, false
	}

	return o.ResourceType, true
}

// HasResourceType returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasResourceType() bool {
	if o != nil && !IsNil(o.ResourceType) {
		return true
	}

	return false
}

// SetResourceType gets a reference to the given string and assigns it to the ResourceType field.
func (o *EventViewForNdsGroup) SetResourceType(v string) {
	o.ResourceType = &v
}

// GetInstanceName returns the InstanceName field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetInstanceName() string {
	if o == nil || IsNil(o.InstanceName) {
		var ret string
		return ret
	}
	return *o.InstanceName
}

// GetInstanceNameOk returns a tuple with the InstanceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetInstanceNameOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceName) {
		return nil, false
	}

	return o.InstanceName, true
}

// HasInstanceName returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasInstanceName() bool {
	if o != nil && !IsNil(o.InstanceName) {
		return true
	}

	return false
}

// SetInstanceName gets a reference to the given string and assigns it to the InstanceName field.
func (o *EventViewForNdsGroup) SetInstanceName(v string) {
	o.InstanceName = &v
}

// GetProcessorErrorMsg returns the ProcessorErrorMsg field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetProcessorErrorMsg() string {
	if o == nil || IsNil(o.ProcessorErrorMsg) {
		var ret string
		return ret
	}
	return *o.ProcessorErrorMsg
}

// GetProcessorErrorMsgOk returns a tuple with the ProcessorErrorMsg field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetProcessorErrorMsgOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessorErrorMsg) {
		return nil, false
	}

	return o.ProcessorErrorMsg, true
}

// HasProcessorErrorMsg returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasProcessorErrorMsg() bool {
	if o != nil && !IsNil(o.ProcessorErrorMsg) {
		return true
	}

	return false
}

// SetProcessorErrorMsg gets a reference to the given string and assigns it to the ProcessorErrorMsg field.
func (o *EventViewForNdsGroup) SetProcessorErrorMsg(v string) {
	o.ProcessorErrorMsg = &v
}

// GetProcessorName returns the ProcessorName field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetProcessorName() string {
	if o == nil || IsNil(o.ProcessorName) {
		var ret string
		return ret
	}
	return *o.ProcessorName
}

// GetProcessorNameOk returns a tuple with the ProcessorName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetProcessorNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessorName) {
		return nil, false
	}

	return o.ProcessorName, true
}

// HasProcessorName returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasProcessorName() bool {
	if o != nil && !IsNil(o.ProcessorName) {
		return true
	}

	return false
}

// SetProcessorName gets a reference to the given string and assigns it to the ProcessorName field.
func (o *EventViewForNdsGroup) SetProcessorName(v string) {
	o.ProcessorName = &v
}

// GetProcessorState returns the ProcessorState field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetProcessorState() string {
	if o == nil || IsNil(o.ProcessorState) {
		var ret string
		return ret
	}
	return *o.ProcessorState
}

// GetProcessorStateOk returns a tuple with the ProcessorState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetProcessorStateOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessorState) {
		return nil, false
	}

	return o.ProcessorState, true
}

// HasProcessorState returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasProcessorState() bool {
	if o != nil && !IsNil(o.ProcessorState) {
		return true
	}

	return false
}

// SetProcessorState gets a reference to the given string and assigns it to the ProcessorState field.
func (o *EventViewForNdsGroup) SetProcessorState(v string) {
	o.ProcessorState = &v
}

// GetResourcePolicyId returns the ResourcePolicyId field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetResourcePolicyId() string {
	if o == nil || IsNil(o.ResourcePolicyId) {
		var ret string
		return ret
	}
	return *o.ResourcePolicyId
}

// GetResourcePolicyIdOk returns a tuple with the ResourcePolicyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetResourcePolicyIdOk() (*string, bool) {
	if o == nil || IsNil(o.ResourcePolicyId) {
		return nil, false
	}

	return o.ResourcePolicyId, true
}

// HasResourcePolicyId returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasResourcePolicyId() bool {
	if o != nil && !IsNil(o.ResourcePolicyId) {
		return true
	}

	return false
}

// SetResourcePolicyId gets a reference to the given string and assigns it to the ResourcePolicyId field.
func (o *EventViewForNdsGroup) SetResourcePolicyId(v string) {
	o.ResourcePolicyId = &v
}

// GetViolatedPolicies returns the ViolatedPolicies field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetViolatedPolicies() []string {
	if o == nil || IsNil(o.ViolatedPolicies) {
		var ret []string
		return ret
	}
	return *o.ViolatedPolicies
}

// GetViolatedPoliciesOk returns a tuple with the ViolatedPolicies field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetViolatedPoliciesOk() (*[]string, bool) {
	if o == nil || IsNil(o.ViolatedPolicies) {
		return nil, false
	}

	return o.ViolatedPolicies, true
}

// HasViolatedPolicies returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasViolatedPolicies() bool {
	if o != nil && !IsNil(o.ViolatedPolicies) {
		return true
	}

	return false
}

// SetViolatedPolicies gets a reference to the given []string and assigns it to the ViolatedPolicies field.
func (o *EventViewForNdsGroup) SetViolatedPolicies(v []string) {
	o.ViolatedPolicies = &v
}

// GetViolationAction returns the ViolationAction field value if set, zero value otherwise
func (o *EventViewForNdsGroup) GetViolationAction() string {
	if o == nil || IsNil(o.ViolationAction) {
		var ret string
		return ret
	}
	return *o.ViolationAction
}

// GetViolationActionOk returns a tuple with the ViolationAction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventViewForNdsGroup) GetViolationActionOk() (*string, bool) {
	if o == nil || IsNil(o.ViolationAction) {
		return nil, false
	}

	return o.ViolationAction, true
}

// HasViolationAction returns a boolean if a field has been set.
func (o *EventViewForNdsGroup) HasViolationAction() bool {
	if o != nil && !IsNil(o.ViolationAction) {
		return true
	}

	return false
}

// SetViolationAction gets a reference to the given string and assigns it to the ViolationAction field.
func (o *EventViewForNdsGroup) SetViolationAction(v string) {
	o.ViolationAction = &v
}
