// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ServerlessInstanceDescription Group of settings that configure a MongoDB serverless instance.
type ServerlessInstanceDescription struct {
	ConnectionStrings *ServerlessInstanceDescriptionConnectionStrings `json:"connectionStrings,omitempty"`
	// Date and time when MongoDB Cloud created this serverless instance. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
	// Read only field.
	CreateDate *time.Time `json:"createDate,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the serverless instance.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Version of MongoDB that the serverless instance runs.
	// Read only field.
	MongoDBVersion *string `json:"mongoDBVersion,omitempty"`
	// Human-readable label that identifies the serverless instance.
	Name                    *string                         `json:"name,omitempty"`
	ProviderSettings        ServerlessProviderSettings      `json:"providerSettings"`
	ServerlessBackupOptions *ClusterServerlessBackupOptions `json:"serverlessBackupOptions,omitempty"`
	// Human-readable label that indicates any current activity being taken on this cluster by the Atlas control plane. With the exception of CREATING and DELETING states, clusters should always be available and have a Primary node even when in states indicating ongoing activity.   - `IDLE`: Atlas is making no changes to this cluster and all changes requested via the UI or API can be assumed to have been applied.  - `CREATING`: A cluster being provisioned for the very first time returns state CREATING until it is ready for connections. Ensure IP Access List and DB Users are configured before attempting to connect.  - `UPDATING`: A change requested via the UI, API, AutoScaling, or other scheduled activity is taking place.  - `DELETING`: The cluster is in the process of deletion and will soon be deleted.  - `REPAIRING`: One or more nodes in the cluster are being returned to service by the Atlas control plane. Other nodes should continue to provide service as normal.
	// Read only field.
	StateName *string `json:"stateName,omitempty"`
	// List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the serverless instance.
	Tags *[]ResourceTag `json:"tags,omitempty"`
	// Flag that indicates whether termination protection is enabled on the serverless instance. If set to `true`, MongoDB Cloud won't delete the serverless instance. If set to `false`, MongoDB Cloud will delete the serverless instance.
	TerminationProtectionEnabled *bool `json:"terminationProtectionEnabled,omitempty"`
}

// NewServerlessInstanceDescription instantiates a new ServerlessInstanceDescription object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessInstanceDescription(providerSettings ServerlessProviderSettings) *ServerlessInstanceDescription {
	this := ServerlessInstanceDescription{}
	this.ProviderSettings = providerSettings
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	return &this
}

// NewServerlessInstanceDescriptionWithDefaults instantiates a new ServerlessInstanceDescription object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessInstanceDescriptionWithDefaults() *ServerlessInstanceDescription {
	this := ServerlessInstanceDescription{}
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	return &this
}

// GetConnectionStrings returns the ConnectionStrings field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetConnectionStrings() ServerlessInstanceDescriptionConnectionStrings {
	if o == nil || IsNil(o.ConnectionStrings) {
		var ret ServerlessInstanceDescriptionConnectionStrings
		return ret
	}
	return *o.ConnectionStrings
}

// GetConnectionStringsOk returns a tuple with the ConnectionStrings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetConnectionStringsOk() (*ServerlessInstanceDescriptionConnectionStrings, bool) {
	if o == nil || IsNil(o.ConnectionStrings) {
		return nil, false
	}

	return o.ConnectionStrings, true
}

// HasConnectionStrings returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasConnectionStrings() bool {
	if o != nil && !IsNil(o.ConnectionStrings) {
		return true
	}

	return false
}

// SetConnectionStrings gets a reference to the given ServerlessInstanceDescriptionConnectionStrings and assigns it to the ConnectionStrings field.
func (o *ServerlessInstanceDescription) SetConnectionStrings(v ServerlessInstanceDescriptionConnectionStrings) {
	o.ConnectionStrings = &v
}

// GetCreateDate returns the CreateDate field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetCreateDate() time.Time {
	if o == nil || IsNil(o.CreateDate) {
		var ret time.Time
		return ret
	}
	return *o.CreateDate
}

// GetCreateDateOk returns a tuple with the CreateDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetCreateDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateDate) {
		return nil, false
	}

	return o.CreateDate, true
}

// HasCreateDate returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasCreateDate() bool {
	if o != nil && !IsNil(o.CreateDate) {
		return true
	}

	return false
}

// SetCreateDate gets a reference to the given time.Time and assigns it to the CreateDate field.
func (o *ServerlessInstanceDescription) SetCreateDate(v time.Time) {
	o.CreateDate = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ServerlessInstanceDescription) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ServerlessInstanceDescription) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ServerlessInstanceDescription) SetLinks(v []Link) {
	o.Links = &v
}

// GetMongoDBVersion returns the MongoDBVersion field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetMongoDBVersion() string {
	if o == nil || IsNil(o.MongoDBVersion) {
		var ret string
		return ret
	}
	return *o.MongoDBVersion
}

// GetMongoDBVersionOk returns a tuple with the MongoDBVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetMongoDBVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MongoDBVersion) {
		return nil, false
	}

	return o.MongoDBVersion, true
}

// HasMongoDBVersion returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasMongoDBVersion() bool {
	if o != nil && !IsNil(o.MongoDBVersion) {
		return true
	}

	return false
}

// SetMongoDBVersion gets a reference to the given string and assigns it to the MongoDBVersion field.
func (o *ServerlessInstanceDescription) SetMongoDBVersion(v string) {
	o.MongoDBVersion = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ServerlessInstanceDescription) SetName(v string) {
	o.Name = &v
}

// GetProviderSettings returns the ProviderSettings field value
func (o *ServerlessInstanceDescription) GetProviderSettings() ServerlessProviderSettings {
	if o == nil {
		var ret ServerlessProviderSettings
		return ret
	}

	return o.ProviderSettings
}

// GetProviderSettingsOk returns a tuple with the ProviderSettings field value
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetProviderSettingsOk() (*ServerlessProviderSettings, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderSettings, true
}

// SetProviderSettings sets field value
func (o *ServerlessInstanceDescription) SetProviderSettings(v ServerlessProviderSettings) {
	o.ProviderSettings = v
}

// GetServerlessBackupOptions returns the ServerlessBackupOptions field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetServerlessBackupOptions() ClusterServerlessBackupOptions {
	if o == nil || IsNil(o.ServerlessBackupOptions) {
		var ret ClusterServerlessBackupOptions
		return ret
	}
	return *o.ServerlessBackupOptions
}

// GetServerlessBackupOptionsOk returns a tuple with the ServerlessBackupOptions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetServerlessBackupOptionsOk() (*ClusterServerlessBackupOptions, bool) {
	if o == nil || IsNil(o.ServerlessBackupOptions) {
		return nil, false
	}

	return o.ServerlessBackupOptions, true
}

// HasServerlessBackupOptions returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasServerlessBackupOptions() bool {
	if o != nil && !IsNil(o.ServerlessBackupOptions) {
		return true
	}

	return false
}

// SetServerlessBackupOptions gets a reference to the given ClusterServerlessBackupOptions and assigns it to the ServerlessBackupOptions field.
func (o *ServerlessInstanceDescription) SetServerlessBackupOptions(v ClusterServerlessBackupOptions) {
	o.ServerlessBackupOptions = &v
}

// GetStateName returns the StateName field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetStateName() string {
	if o == nil || IsNil(o.StateName) {
		var ret string
		return ret
	}
	return *o.StateName
}

// GetStateNameOk returns a tuple with the StateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetStateNameOk() (*string, bool) {
	if o == nil || IsNil(o.StateName) {
		return nil, false
	}

	return o.StateName, true
}

// HasStateName returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasStateName() bool {
	if o != nil && !IsNil(o.StateName) {
		return true
	}

	return false
}

// SetStateName gets a reference to the given string and assigns it to the StateName field.
func (o *ServerlessInstanceDescription) SetStateName(v string) {
	o.StateName = &v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetTags() []ResourceTag {
	if o == nil || IsNil(o.Tags) {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *ServerlessInstanceDescription) SetTags(v []ResourceTag) {
	o.Tags = &v
}

// GetTerminationProtectionEnabled returns the TerminationProtectionEnabled field value if set, zero value otherwise
func (o *ServerlessInstanceDescription) GetTerminationProtectionEnabled() bool {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		var ret bool
		return ret
	}
	return *o.TerminationProtectionEnabled
}

// GetTerminationProtectionEnabledOk returns a tuple with the TerminationProtectionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessInstanceDescription) GetTerminationProtectionEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		return nil, false
	}

	return o.TerminationProtectionEnabled, true
}

// HasTerminationProtectionEnabled returns a boolean if a field has been set.
func (o *ServerlessInstanceDescription) HasTerminationProtectionEnabled() bool {
	if o != nil && !IsNil(o.TerminationProtectionEnabled) {
		return true
	}

	return false
}

// SetTerminationProtectionEnabled gets a reference to the given bool and assigns it to the TerminationProtectionEnabled field.
func (o *ServerlessInstanceDescription) SetTerminationProtectionEnabled(v bool) {
	o.TerminationProtectionEnabled = &v
}
