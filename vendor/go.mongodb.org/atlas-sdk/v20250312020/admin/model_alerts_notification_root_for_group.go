// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AlertsNotificationRootForGroup One target that MongoDB Cloud sends notifications when an alert triggers.
type AlertsNotificationRootForGroup struct {
	// Datadog API Key that MongoDB Cloud needs to send alert notifications to Datadog. You can find this API key in the Datadog dashboard. The resource requires this parameter when `\"notifications.[n].typeName\" : \"DATADOG\"`.  **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:  * View or edit the alert through the Atlas UI.  * Query the alert for the notification through the Atlas Administration API.
	DatadogApiKey *string `json:"datadogApiKey,omitempty"`
	// Datadog region that indicates which API Uniform Resource Locator (URL) to use. The resource requires this parameter when `\"notifications.[n].typeName\" : \"DATADOG\"`.
	DatadogRegion *string `json:"datadogRegion,omitempty"`
	// Number of minutes that MongoDB Cloud waits after detecting an alert condition before it sends out the first notification.
	DelayMin *int `json:"delayMin,omitempty"`
	// The id of the associated integration, the credentials of which to use for requests.
	IntegrationId *string `json:"integrationId,omitempty"`
	// Number of minutes to wait between successive notifications. MongoDB Cloud sends notifications until someone acknowledges the unacknowledged alert.  PagerDuty, VictorOps, and OpsGenie notifications don't return this element. Configure and manage the notification interval within each of those services.
	IntervalMin *int `json:"intervalMin,omitempty"`
	// The `notifierId` is a system-generated unique identifier assigned to each notification method. This is needed when updating third-party notifications without requiring explicit authentication credentials.
	NotifierId *string `json:"notifierId,omitempty"`
	// Human-readable label that displays the alert notification type.
	TypeName *string `json:"typeName,omitempty"`
	// Email address to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `\"notifications.[n].typeName\" : \"EMAIL\"`. You don't need to set this value to send emails to individual or groups of MongoDB Cloud users including:  - specific MongoDB Cloud users (`\"notifications.[n].typeName\" : \"USER\"`) - MongoDB Cloud users with specific project roles (`\"notifications.[n].typeName\" : \"GROUP\"`) - MongoDB Cloud users with specific organization roles (`\"notifications.[n].typeName\" : \"ORG\"`) - MongoDB Cloud teams (`\"notifications.[n].typeName\" : \"TEAM\"`)  To send emails to one MongoDB Cloud user or grouping of users, set the `notifications.[n].emailEnabled` parameter.
	EmailAddress *string `json:"emailAddress,omitempty"`
	// Flag that indicates whether MongoDB Cloud should send email notifications. The resource requires this parameter when one of the following values have been set:  - `\"notifications.[n].typeName\" : \"ORG\"` - `\"notifications.[n].typeName\" : \"GROUP\"` - `\"notifications.[n].typeName\" : \"USER\"`
	EmailEnabled *bool `json:"emailEnabled,omitempty"`
	// List that contains the one or more organization roles that receive the configured alert. This parameter is available when `\"notifications.[n].typeName\" : \"GROUP\"` or `\"notifications.[n].typeName\" : \"ORG\"`. If you include this parameter, MongoDB Cloud sends alerts only to users assigned the roles you specify in the array. If you omit this parameter, MongoDB Cloud sends alerts to users assigned any role.
	Roles *[]string `json:"roles,omitempty"`
	// Flag that indicates whether MongoDB Cloud should send text message notifications. The resource requires this parameter when one of the following values have been set:  - `\"notifications.[n].typeName\" : \"ORG\"` - `\"notifications.[n].typeName\" : \"GROUP\"` - `\"notifications.[n].typeName\" : \"USER\"`
	SmsEnabled *bool `json:"smsEnabled,omitempty"`
	// HipChat API token that MongoDB Cloud needs to send alert notifications to HipChat. The resource requires this parameter when `\"notifications.[n].typeName\" : \"HIP_CHAT\"`\". If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes it.  **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:  * View or edit the alert through the Atlas UI.  * Query the alert for the notification through the Atlas Administration API.
	NotificationToken *string `json:"notificationToken,omitempty"`
	// HipChat API room name to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `\"notifications.[n].typeName\" : \"HIP_CHAT\"`\".
	RoomName *string `json:"roomName,omitempty"`
	// Microsoft Teams Webhook Uniform Resource Locator (URL) that MongoDB Cloud needs to send this notification via Microsoft Teams. The resource requires this parameter when `\"notifications.[n].typeName\" : \"MICROSOFT_TEAMS\"`. If the URL later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.  **NOTE**: When you view or edit the alert for a Microsoft Teams notification, the URL appears partially redacted.
	MicrosoftTeamsWebhookUrl *string `json:"microsoftTeamsWebhookUrl,omitempty"`
	// API Key that MongoDB Cloud needs to send this notification via OpsGenie. The resource requires this parameter when `\"notifications.[n].typeName\" : \"OPS_GENIE\"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.  **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:  * View or edit the alert through the Atlas UI.  * Query the alert for the notification through the Atlas Administration API.
	OpsGenieApiKey *string `json:"opsGenieApiKey,omitempty"`
	// OpsGenie region that indicates which API Uniform Resource Locator (URL) to use.
	OpsGenieRegion *string `json:"opsGenieRegion,omitempty"`
	// PagerDuty region that indicates which API Uniform Resource Locator (URL) to use.
	Region *string `json:"region,omitempty"`
	// PagerDuty service key that MongoDB Cloud needs to send notifications via PagerDuty. The resource requires this parameter when `\"notifications.[n].typeName\" : \"PAGER_DUTY\"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.  **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:  * View or edit the alert through the Atlas UI.  * Query the alert for the notification through the Atlas Administration API.
	ServiceKey *string `json:"serviceKey,omitempty"`
	// Slack API token or Bot token that MongoDB Cloud needs to send alert notifications via Slack. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`. If the token later becomes invalid, MongoDB Cloud sends an email to the project owners. If the token remains invalid, MongoDB Cloud removes the token.   **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:  * View or edit the alert through the Atlas UI.  * Query the alert for the notification through the Atlas Administration API.
	ApiToken *string `json:"apiToken,omitempty"`
	// Name of the Slack channel to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SLACK\"`.
	ChannelName *string `json:"channelName,omitempty"`
	// Mobile phone number to which MongoDB Cloud sends alert notifications. The resource requires this parameter when `\"notifications.[n].typeName\" : \"SMS\"`.
	MobileNumber *string `json:"mobileNumber,omitempty"`
	// Unique 24-hexadecimal digit string that identifies one MongoDB Cloud team. The resource requires this parameter when `\"notifications.[n].typeName\" : \"TEAM\"`.
	TeamId *string `json:"teamId,omitempty"`
	// Name of the MongoDB Cloud team that receives this notification. The resource requires this parameter when `\"notifications.[n].typeName\" : \"TEAM\"`.
	TeamName *string `json:"teamName,omitempty"`
	// MongoDB Cloud username of the person to whom MongoDB Cloud sends notifications. Specify only MongoDB Cloud users who belong to the project that owns the alert configuration. The resource requires this parameter when `\"notifications.[n].typeName\" : \"USER\"`.
	Username *string `json:"username,omitempty"`
	// API key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `\"notifications.[n].typeName\" : \"VICTOR_OPS\"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.  **NOTE**: After you create a notification which requires an API or integration key, the key appears partially redacted when you:  * View or edit the alert through the Atlas UI.  * Query the alert for the notification through the Atlas Administration API.
	VictorOpsApiKey *string `json:"victorOpsApiKey,omitempty"`
	// Routing key that MongoDB Cloud needs to send alert notifications to Splunk On-Call. The resource requires this parameter when `\"notifications.[n].typeName\" : \"VICTOR_OPS\"`. If the key later becomes invalid, MongoDB Cloud sends an email to the project owners. If the key remains invalid, MongoDB Cloud removes it.
	VictorOpsRoutingKey *string `json:"victorOpsRoutingKey,omitempty"`
	// Authentication secret for a webhook-based alert.  Atlas returns this value if you set `notifications.[n].typeName` :`WEBHOOK` and either: * You set `notification.[n].webhookSecret` to a non-empty string * You set a default webhook secret either on the Integrations page, or with the Integrations API  **NOTE**: When you view or edit the alert for a webhook notification, the secret appears completely redacted.
	WebhookSecret *string `json:"webhookSecret,omitempty"`
	// Target URL for a webhook-based alert.  Atlas returns this value if you set `\"notifications.[n].typeName\" :\"WEBHOOK\"` and either: * You set `notification.[n].webhookURL` to a non-empty string * You set a default webhook URL either on the Integrations page, or with the Integrations API  **NOTE**: When you view or edit the alert for a Webhook URL notification, the URL appears partially redacted.
	WebhookUrl *string `json:"webhookUrl,omitempty"`
}

// NewAlertsNotificationRootForGroup instantiates a new AlertsNotificationRootForGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAlertsNotificationRootForGroup() *AlertsNotificationRootForGroup {
	this := AlertsNotificationRootForGroup{}
	var datadogRegion string = "US"
	this.DatadogRegion = &datadogRegion
	var opsGenieRegion string = "US"
	this.OpsGenieRegion = &opsGenieRegion
	var region string = "US"
	this.Region = &region
	return &this
}

// NewAlertsNotificationRootForGroupWithDefaults instantiates a new AlertsNotificationRootForGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAlertsNotificationRootForGroupWithDefaults() *AlertsNotificationRootForGroup {
	this := AlertsNotificationRootForGroup{}
	var datadogRegion string = "US"
	this.DatadogRegion = &datadogRegion
	var opsGenieRegion string = "US"
	this.OpsGenieRegion = &opsGenieRegion
	var region string = "US"
	this.Region = &region
	return &this
}

// GetDatadogApiKey returns the DatadogApiKey field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetDatadogApiKey() string {
	if o == nil || IsNil(o.DatadogApiKey) {
		var ret string
		return ret
	}
	return *o.DatadogApiKey
}

// GetDatadogApiKeyOk returns a tuple with the DatadogApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetDatadogApiKeyOk() (*string, bool) {
	if o == nil || IsNil(o.DatadogApiKey) {
		return nil, false
	}

	return o.DatadogApiKey, true
}

// HasDatadogApiKey returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasDatadogApiKey() bool {
	if o != nil && !IsNil(o.DatadogApiKey) {
		return true
	}

	return false
}

// SetDatadogApiKey gets a reference to the given string and assigns it to the DatadogApiKey field.
func (o *AlertsNotificationRootForGroup) SetDatadogApiKey(v string) {
	o.DatadogApiKey = &v
}

// GetDatadogRegion returns the DatadogRegion field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetDatadogRegion() string {
	if o == nil || IsNil(o.DatadogRegion) {
		var ret string
		return ret
	}
	return *o.DatadogRegion
}

// GetDatadogRegionOk returns a tuple with the DatadogRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetDatadogRegionOk() (*string, bool) {
	if o == nil || IsNil(o.DatadogRegion) {
		return nil, false
	}

	return o.DatadogRegion, true
}

// HasDatadogRegion returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasDatadogRegion() bool {
	if o != nil && !IsNil(o.DatadogRegion) {
		return true
	}

	return false
}

// SetDatadogRegion gets a reference to the given string and assigns it to the DatadogRegion field.
func (o *AlertsNotificationRootForGroup) SetDatadogRegion(v string) {
	o.DatadogRegion = &v
}

// GetDelayMin returns the DelayMin field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetDelayMin() int {
	if o == nil || IsNil(o.DelayMin) {
		var ret int
		return ret
	}
	return *o.DelayMin
}

// GetDelayMinOk returns a tuple with the DelayMin field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetDelayMinOk() (*int, bool) {
	if o == nil || IsNil(o.DelayMin) {
		return nil, false
	}

	return o.DelayMin, true
}

// HasDelayMin returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasDelayMin() bool {
	if o != nil && !IsNil(o.DelayMin) {
		return true
	}

	return false
}

// SetDelayMin gets a reference to the given int and assigns it to the DelayMin field.
func (o *AlertsNotificationRootForGroup) SetDelayMin(v int) {
	o.DelayMin = &v
}

// GetIntegrationId returns the IntegrationId field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetIntegrationId() string {
	if o == nil || IsNil(o.IntegrationId) {
		var ret string
		return ret
	}
	return *o.IntegrationId
}

// GetIntegrationIdOk returns a tuple with the IntegrationId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetIntegrationIdOk() (*string, bool) {
	if o == nil || IsNil(o.IntegrationId) {
		return nil, false
	}

	return o.IntegrationId, true
}

// HasIntegrationId returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasIntegrationId() bool {
	if o != nil && !IsNil(o.IntegrationId) {
		return true
	}

	return false
}

// SetIntegrationId gets a reference to the given string and assigns it to the IntegrationId field.
func (o *AlertsNotificationRootForGroup) SetIntegrationId(v string) {
	o.IntegrationId = &v
}

// GetIntervalMin returns the IntervalMin field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetIntervalMin() int {
	if o == nil || IsNil(o.IntervalMin) {
		var ret int
		return ret
	}
	return *o.IntervalMin
}

// GetIntervalMinOk returns a tuple with the IntervalMin field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetIntervalMinOk() (*int, bool) {
	if o == nil || IsNil(o.IntervalMin) {
		return nil, false
	}

	return o.IntervalMin, true
}

// HasIntervalMin returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasIntervalMin() bool {
	if o != nil && !IsNil(o.IntervalMin) {
		return true
	}

	return false
}

// SetIntervalMin gets a reference to the given int and assigns it to the IntervalMin field.
func (o *AlertsNotificationRootForGroup) SetIntervalMin(v int) {
	o.IntervalMin = &v
}

// GetNotifierId returns the NotifierId field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetNotifierId() string {
	if o == nil || IsNil(o.NotifierId) {
		var ret string
		return ret
	}
	return *o.NotifierId
}

// GetNotifierIdOk returns a tuple with the NotifierId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetNotifierIdOk() (*string, bool) {
	if o == nil || IsNil(o.NotifierId) {
		return nil, false
	}

	return o.NotifierId, true
}

// HasNotifierId returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasNotifierId() bool {
	if o != nil && !IsNil(o.NotifierId) {
		return true
	}

	return false
}

// SetNotifierId gets a reference to the given string and assigns it to the NotifierId field.
func (o *AlertsNotificationRootForGroup) SetNotifierId(v string) {
	o.NotifierId = &v
}

// GetTypeName returns the TypeName field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetTypeName() string {
	if o == nil || IsNil(o.TypeName) {
		var ret string
		return ret
	}
	return *o.TypeName
}

// GetTypeNameOk returns a tuple with the TypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.TypeName) {
		return nil, false
	}

	return o.TypeName, true
}

// HasTypeName returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasTypeName() bool {
	if o != nil && !IsNil(o.TypeName) {
		return true
	}

	return false
}

// SetTypeName gets a reference to the given string and assigns it to the TypeName field.
func (o *AlertsNotificationRootForGroup) SetTypeName(v string) {
	o.TypeName = &v
}

// GetEmailAddress returns the EmailAddress field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetEmailAddress() string {
	if o == nil || IsNil(o.EmailAddress) {
		var ret string
		return ret
	}
	return *o.EmailAddress
}

// GetEmailAddressOk returns a tuple with the EmailAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetEmailAddressOk() (*string, bool) {
	if o == nil || IsNil(o.EmailAddress) {
		return nil, false
	}

	return o.EmailAddress, true
}

// HasEmailAddress returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasEmailAddress() bool {
	if o != nil && !IsNil(o.EmailAddress) {
		return true
	}

	return false
}

// SetEmailAddress gets a reference to the given string and assigns it to the EmailAddress field.
func (o *AlertsNotificationRootForGroup) SetEmailAddress(v string) {
	o.EmailAddress = &v
}

// GetEmailEnabled returns the EmailEnabled field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetEmailEnabled() bool {
	if o == nil || IsNil(o.EmailEnabled) {
		var ret bool
		return ret
	}
	return *o.EmailEnabled
}

// GetEmailEnabledOk returns a tuple with the EmailEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetEmailEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.EmailEnabled) {
		return nil, false
	}

	return o.EmailEnabled, true
}

// HasEmailEnabled returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasEmailEnabled() bool {
	if o != nil && !IsNil(o.EmailEnabled) {
		return true
	}

	return false
}

// SetEmailEnabled gets a reference to the given bool and assigns it to the EmailEnabled field.
func (o *AlertsNotificationRootForGroup) SetEmailEnabled(v bool) {
	o.EmailEnabled = &v
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *AlertsNotificationRootForGroup) SetRoles(v []string) {
	o.Roles = &v
}

// GetSmsEnabled returns the SmsEnabled field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetSmsEnabled() bool {
	if o == nil || IsNil(o.SmsEnabled) {
		var ret bool
		return ret
	}
	return *o.SmsEnabled
}

// GetSmsEnabledOk returns a tuple with the SmsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetSmsEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.SmsEnabled) {
		return nil, false
	}

	return o.SmsEnabled, true
}

// HasSmsEnabled returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasSmsEnabled() bool {
	if o != nil && !IsNil(o.SmsEnabled) {
		return true
	}

	return false
}

// SetSmsEnabled gets a reference to the given bool and assigns it to the SmsEnabled field.
func (o *AlertsNotificationRootForGroup) SetSmsEnabled(v bool) {
	o.SmsEnabled = &v
}

// GetNotificationToken returns the NotificationToken field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetNotificationToken() string {
	if o == nil || IsNil(o.NotificationToken) {
		var ret string
		return ret
	}
	return *o.NotificationToken
}

// GetNotificationTokenOk returns a tuple with the NotificationToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetNotificationTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NotificationToken) {
		return nil, false
	}

	return o.NotificationToken, true
}

// HasNotificationToken returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasNotificationToken() bool {
	if o != nil && !IsNil(o.NotificationToken) {
		return true
	}

	return false
}

// SetNotificationToken gets a reference to the given string and assigns it to the NotificationToken field.
func (o *AlertsNotificationRootForGroup) SetNotificationToken(v string) {
	o.NotificationToken = &v
}

// GetRoomName returns the RoomName field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetRoomName() string {
	if o == nil || IsNil(o.RoomName) {
		var ret string
		return ret
	}
	return *o.RoomName
}

// GetRoomNameOk returns a tuple with the RoomName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetRoomNameOk() (*string, bool) {
	if o == nil || IsNil(o.RoomName) {
		return nil, false
	}

	return o.RoomName, true
}

// HasRoomName returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasRoomName() bool {
	if o != nil && !IsNil(o.RoomName) {
		return true
	}

	return false
}

// SetRoomName gets a reference to the given string and assigns it to the RoomName field.
func (o *AlertsNotificationRootForGroup) SetRoomName(v string) {
	o.RoomName = &v
}

// GetMicrosoftTeamsWebhookUrl returns the MicrosoftTeamsWebhookUrl field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetMicrosoftTeamsWebhookUrl() string {
	if o == nil || IsNil(o.MicrosoftTeamsWebhookUrl) {
		var ret string
		return ret
	}
	return *o.MicrosoftTeamsWebhookUrl
}

// GetMicrosoftTeamsWebhookUrlOk returns a tuple with the MicrosoftTeamsWebhookUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetMicrosoftTeamsWebhookUrlOk() (*string, bool) {
	if o == nil || IsNil(o.MicrosoftTeamsWebhookUrl) {
		return nil, false
	}

	return o.MicrosoftTeamsWebhookUrl, true
}

// HasMicrosoftTeamsWebhookUrl returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasMicrosoftTeamsWebhookUrl() bool {
	if o != nil && !IsNil(o.MicrosoftTeamsWebhookUrl) {
		return true
	}

	return false
}

// SetMicrosoftTeamsWebhookUrl gets a reference to the given string and assigns it to the MicrosoftTeamsWebhookUrl field.
func (o *AlertsNotificationRootForGroup) SetMicrosoftTeamsWebhookUrl(v string) {
	o.MicrosoftTeamsWebhookUrl = &v
}

// GetOpsGenieApiKey returns the OpsGenieApiKey field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetOpsGenieApiKey() string {
	if o == nil || IsNil(o.OpsGenieApiKey) {
		var ret string
		return ret
	}
	return *o.OpsGenieApiKey
}

// GetOpsGenieApiKeyOk returns a tuple with the OpsGenieApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetOpsGenieApiKeyOk() (*string, bool) {
	if o == nil || IsNil(o.OpsGenieApiKey) {
		return nil, false
	}

	return o.OpsGenieApiKey, true
}

// HasOpsGenieApiKey returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasOpsGenieApiKey() bool {
	if o != nil && !IsNil(o.OpsGenieApiKey) {
		return true
	}

	return false
}

// SetOpsGenieApiKey gets a reference to the given string and assigns it to the OpsGenieApiKey field.
func (o *AlertsNotificationRootForGroup) SetOpsGenieApiKey(v string) {
	o.OpsGenieApiKey = &v
}

// GetOpsGenieRegion returns the OpsGenieRegion field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetOpsGenieRegion() string {
	if o == nil || IsNil(o.OpsGenieRegion) {
		var ret string
		return ret
	}
	return *o.OpsGenieRegion
}

// GetOpsGenieRegionOk returns a tuple with the OpsGenieRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetOpsGenieRegionOk() (*string, bool) {
	if o == nil || IsNil(o.OpsGenieRegion) {
		return nil, false
	}

	return o.OpsGenieRegion, true
}

// HasOpsGenieRegion returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasOpsGenieRegion() bool {
	if o != nil && !IsNil(o.OpsGenieRegion) {
		return true
	}

	return false
}

// SetOpsGenieRegion gets a reference to the given string and assigns it to the OpsGenieRegion field.
func (o *AlertsNotificationRootForGroup) SetOpsGenieRegion(v string) {
	o.OpsGenieRegion = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *AlertsNotificationRootForGroup) SetRegion(v string) {
	o.Region = &v
}

// GetServiceKey returns the ServiceKey field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetServiceKey() string {
	if o == nil || IsNil(o.ServiceKey) {
		var ret string
		return ret
	}
	return *o.ServiceKey
}

// GetServiceKeyOk returns a tuple with the ServiceKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetServiceKeyOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceKey) {
		return nil, false
	}

	return o.ServiceKey, true
}

// HasServiceKey returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasServiceKey() bool {
	if o != nil && !IsNil(o.ServiceKey) {
		return true
	}

	return false
}

// SetServiceKey gets a reference to the given string and assigns it to the ServiceKey field.
func (o *AlertsNotificationRootForGroup) SetServiceKey(v string) {
	o.ServiceKey = &v
}

// GetApiToken returns the ApiToken field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetApiToken() string {
	if o == nil || IsNil(o.ApiToken) {
		var ret string
		return ret
	}
	return *o.ApiToken
}

// GetApiTokenOk returns a tuple with the ApiToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetApiTokenOk() (*string, bool) {
	if o == nil || IsNil(o.ApiToken) {
		return nil, false
	}

	return o.ApiToken, true
}

// HasApiToken returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasApiToken() bool {
	if o != nil && !IsNil(o.ApiToken) {
		return true
	}

	return false
}

// SetApiToken gets a reference to the given string and assigns it to the ApiToken field.
func (o *AlertsNotificationRootForGroup) SetApiToken(v string) {
	o.ApiToken = &v
}

// GetChannelName returns the ChannelName field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetChannelName() string {
	if o == nil || IsNil(o.ChannelName) {
		var ret string
		return ret
	}
	return *o.ChannelName
}

// GetChannelNameOk returns a tuple with the ChannelName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetChannelNameOk() (*string, bool) {
	if o == nil || IsNil(o.ChannelName) {
		return nil, false
	}

	return o.ChannelName, true
}

// HasChannelName returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasChannelName() bool {
	if o != nil && !IsNil(o.ChannelName) {
		return true
	}

	return false
}

// SetChannelName gets a reference to the given string and assigns it to the ChannelName field.
func (o *AlertsNotificationRootForGroup) SetChannelName(v string) {
	o.ChannelName = &v
}

// GetMobileNumber returns the MobileNumber field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetMobileNumber() string {
	if o == nil || IsNil(o.MobileNumber) {
		var ret string
		return ret
	}
	return *o.MobileNumber
}

// GetMobileNumberOk returns a tuple with the MobileNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetMobileNumberOk() (*string, bool) {
	if o == nil || IsNil(o.MobileNumber) {
		return nil, false
	}

	return o.MobileNumber, true
}

// HasMobileNumber returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasMobileNumber() bool {
	if o != nil && !IsNil(o.MobileNumber) {
		return true
	}

	return false
}

// SetMobileNumber gets a reference to the given string and assigns it to the MobileNumber field.
func (o *AlertsNotificationRootForGroup) SetMobileNumber(v string) {
	o.MobileNumber = &v
}

// GetTeamId returns the TeamId field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetTeamId() string {
	if o == nil || IsNil(o.TeamId) {
		var ret string
		return ret
	}
	return *o.TeamId
}

// GetTeamIdOk returns a tuple with the TeamId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetTeamIdOk() (*string, bool) {
	if o == nil || IsNil(o.TeamId) {
		return nil, false
	}

	return o.TeamId, true
}

// HasTeamId returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasTeamId() bool {
	if o != nil && !IsNil(o.TeamId) {
		return true
	}

	return false
}

// SetTeamId gets a reference to the given string and assigns it to the TeamId field.
func (o *AlertsNotificationRootForGroup) SetTeamId(v string) {
	o.TeamId = &v
}

// GetTeamName returns the TeamName field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetTeamName() string {
	if o == nil || IsNil(o.TeamName) {
		var ret string
		return ret
	}
	return *o.TeamName
}

// GetTeamNameOk returns a tuple with the TeamName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetTeamNameOk() (*string, bool) {
	if o == nil || IsNil(o.TeamName) {
		return nil, false
	}

	return o.TeamName, true
}

// HasTeamName returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasTeamName() bool {
	if o != nil && !IsNil(o.TeamName) {
		return true
	}

	return false
}

// SetTeamName gets a reference to the given string and assigns it to the TeamName field.
func (o *AlertsNotificationRootForGroup) SetTeamName(v string) {
	o.TeamName = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *AlertsNotificationRootForGroup) SetUsername(v string) {
	o.Username = &v
}

// GetVictorOpsApiKey returns the VictorOpsApiKey field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetVictorOpsApiKey() string {
	if o == nil || IsNil(o.VictorOpsApiKey) {
		var ret string
		return ret
	}
	return *o.VictorOpsApiKey
}

// GetVictorOpsApiKeyOk returns a tuple with the VictorOpsApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetVictorOpsApiKeyOk() (*string, bool) {
	if o == nil || IsNil(o.VictorOpsApiKey) {
		return nil, false
	}

	return o.VictorOpsApiKey, true
}

// HasVictorOpsApiKey returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasVictorOpsApiKey() bool {
	if o != nil && !IsNil(o.VictorOpsApiKey) {
		return true
	}

	return false
}

// SetVictorOpsApiKey gets a reference to the given string and assigns it to the VictorOpsApiKey field.
func (o *AlertsNotificationRootForGroup) SetVictorOpsApiKey(v string) {
	o.VictorOpsApiKey = &v
}

// GetVictorOpsRoutingKey returns the VictorOpsRoutingKey field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetVictorOpsRoutingKey() string {
	if o == nil || IsNil(o.VictorOpsRoutingKey) {
		var ret string
		return ret
	}
	return *o.VictorOpsRoutingKey
}

// GetVictorOpsRoutingKeyOk returns a tuple with the VictorOpsRoutingKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetVictorOpsRoutingKeyOk() (*string, bool) {
	if o == nil || IsNil(o.VictorOpsRoutingKey) {
		return nil, false
	}

	return o.VictorOpsRoutingKey, true
}

// HasVictorOpsRoutingKey returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasVictorOpsRoutingKey() bool {
	if o != nil && !IsNil(o.VictorOpsRoutingKey) {
		return true
	}

	return false
}

// SetVictorOpsRoutingKey gets a reference to the given string and assigns it to the VictorOpsRoutingKey field.
func (o *AlertsNotificationRootForGroup) SetVictorOpsRoutingKey(v string) {
	o.VictorOpsRoutingKey = &v
}

// GetWebhookSecret returns the WebhookSecret field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetWebhookSecret() string {
	if o == nil || IsNil(o.WebhookSecret) {
		var ret string
		return ret
	}
	return *o.WebhookSecret
}

// GetWebhookSecretOk returns a tuple with the WebhookSecret field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetWebhookSecretOk() (*string, bool) {
	if o == nil || IsNil(o.WebhookSecret) {
		return nil, false
	}

	return o.WebhookSecret, true
}

// HasWebhookSecret returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasWebhookSecret() bool {
	if o != nil && !IsNil(o.WebhookSecret) {
		return true
	}

	return false
}

// SetWebhookSecret gets a reference to the given string and assigns it to the WebhookSecret field.
func (o *AlertsNotificationRootForGroup) SetWebhookSecret(v string) {
	o.WebhookSecret = &v
}

// GetWebhookUrl returns the WebhookUrl field value if set, zero value otherwise
func (o *AlertsNotificationRootForGroup) GetWebhookUrl() string {
	if o == nil || IsNil(o.WebhookUrl) {
		var ret string
		return ret
	}
	return *o.WebhookUrl
}

// GetWebhookUrlOk returns a tuple with the WebhookUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertsNotificationRootForGroup) GetWebhookUrlOk() (*string, bool) {
	if o == nil || IsNil(o.WebhookUrl) {
		return nil, false
	}

	return o.WebhookUrl, true
}

// HasWebhookUrl returns a boolean if a field has been set.
func (o *AlertsNotificationRootForGroup) HasWebhookUrl() bool {
	if o != nil && !IsNil(o.WebhookUrl) {
		return true
	}

	return false
}

// SetWebhookUrl gets a reference to the given string and assigns it to the WebhookUrl field.
func (o *AlertsNotificationRootForGroup) SetWebhookUrl(v string) {
	o.WebhookUrl = &v
}
