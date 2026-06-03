// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ApiHostViewAtlas struct for ApiHostViewAtlas
type ApiHostViewAtlas struct {
	// Date and time when MongoDB Cloud created this MongoDB process. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project. The project contains MongoDB processes that you want to return. The MongoDB process can be either the `mongod` or `mongos`.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`).
	// Read only field.
	Hostname *string `json:"hostname,omitempty"`
	// Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Date and time when MongoDB Cloud received the last ping for this MongoDB process. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastPing *time.Time `json:"lastPing,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]LinkAtlas `json:"links,omitempty"`
	// Internet Assigned Numbers Authority (IANA) port on which the MongoDB process listens for requests.
	// Read only field.
	Port *int `json:"port,omitempty"`
	// Human-readable label that identifies the replica set that contains this process. This resource returns this parameter if this process belongs to a replica set.
	// Read only field.
	ReplicaSetName *string `json:"replicaSetName,omitempty"`
	// Human-readable label that identifies the shard that contains this process. This resource returns this value only if this process belongs to a sharded cluster.
	// Read only field.
	ShardName *string `json:"shardName,omitempty"`
	// Type of MongoDB process that MongoDB Cloud tracks. MongoDB Cloud returns new processes as `NO_DATA` until MongoDB Cloud completes deploying the process.
	// Read only field.
	TypeName *string `json:"typeName,omitempty"`
	// Human-readable label that identifies the cluster node. MongoDB Cloud sets this hostname usually to the standard hostname for the cluster node. It appears in the connection string for a cluster instead of the value of the hostname parameter.
	// Read only field.
	UserAlias *string `json:"userAlias,omitempty"`
	// Version of MongoDB that this process runs.
	// Read only field.
	Version *string `json:"version,omitempty"`
}

// NewApiHostViewAtlas instantiates a new ApiHostViewAtlas object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiHostViewAtlas() *ApiHostViewAtlas {
	this := ApiHostViewAtlas{}
	return &this
}

// NewApiHostViewAtlasWithDefaults instantiates a new ApiHostViewAtlas object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiHostViewAtlasWithDefaults() *ApiHostViewAtlas {
	this := ApiHostViewAtlas{}
	return &this
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *ApiHostViewAtlas) SetCreated(v time.Time) {
	o.Created = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ApiHostViewAtlas) SetGroupId(v string) {
	o.GroupId = &v
}

// GetHostname returns the Hostname field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetHostname() string {
	if o == nil || IsNil(o.Hostname) {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetHostnameOk() (*string, bool) {
	if o == nil || IsNil(o.Hostname) {
		return nil, false
	}

	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasHostname() bool {
	if o != nil && !IsNil(o.Hostname) {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *ApiHostViewAtlas) SetHostname(v string) {
	o.Hostname = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiHostViewAtlas) SetId(v string) {
	o.Id = &v
}

// GetLastPing returns the LastPing field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetLastPing() time.Time {
	if o == nil || IsNil(o.LastPing) {
		var ret time.Time
		return ret
	}
	return *o.LastPing
}

// GetLastPingOk returns a tuple with the LastPing field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetLastPingOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastPing) {
		return nil, false
	}

	return o.LastPing, true
}

// HasLastPing returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasLastPing() bool {
	if o != nil && !IsNil(o.LastPing) {
		return true
	}

	return false
}

// SetLastPing gets a reference to the given time.Time and assigns it to the LastPing field.
func (o *ApiHostViewAtlas) SetLastPing(v time.Time) {
	o.LastPing = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetLinks() []LinkAtlas {
	if o == nil || IsNil(o.Links) {
		var ret []LinkAtlas
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetLinksOk() (*[]LinkAtlas, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []LinkAtlas and assigns it to the Links field.
func (o *ApiHostViewAtlas) SetLinks(v []LinkAtlas) {
	o.Links = &v
}

// GetPort returns the Port field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetPort() int {
	if o == nil || IsNil(o.Port) {
		var ret int
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetPortOk() (*int, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}

	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int and assigns it to the Port field.
func (o *ApiHostViewAtlas) SetPort(v int) {
	o.Port = &v
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *ApiHostViewAtlas) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}

// GetShardName returns the ShardName field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetShardName() string {
	if o == nil || IsNil(o.ShardName) {
		var ret string
		return ret
	}
	return *o.ShardName
}

// GetShardNameOk returns a tuple with the ShardName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetShardNameOk() (*string, bool) {
	if o == nil || IsNil(o.ShardName) {
		return nil, false
	}

	return o.ShardName, true
}

// HasShardName returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasShardName() bool {
	if o != nil && !IsNil(o.ShardName) {
		return true
	}

	return false
}

// SetShardName gets a reference to the given string and assigns it to the ShardName field.
func (o *ApiHostViewAtlas) SetShardName(v string) {
	o.ShardName = &v
}

// GetTypeName returns the TypeName field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetTypeName() string {
	if o == nil || IsNil(o.TypeName) {
		var ret string
		return ret
	}
	return *o.TypeName
}

// GetTypeNameOk returns a tuple with the TypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.TypeName) {
		return nil, false
	}

	return o.TypeName, true
}

// HasTypeName returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasTypeName() bool {
	if o != nil && !IsNil(o.TypeName) {
		return true
	}

	return false
}

// SetTypeName gets a reference to the given string and assigns it to the TypeName field.
func (o *ApiHostViewAtlas) SetTypeName(v string) {
	o.TypeName = &v
}

// GetUserAlias returns the UserAlias field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetUserAlias() string {
	if o == nil || IsNil(o.UserAlias) {
		var ret string
		return ret
	}
	return *o.UserAlias
}

// GetUserAliasOk returns a tuple with the UserAlias field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetUserAliasOk() (*string, bool) {
	if o == nil || IsNil(o.UserAlias) {
		return nil, false
	}

	return o.UserAlias, true
}

// HasUserAlias returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasUserAlias() bool {
	if o != nil && !IsNil(o.UserAlias) {
		return true
	}

	return false
}

// SetUserAlias gets a reference to the given string and assigns it to the UserAlias field.
func (o *ApiHostViewAtlas) SetUserAlias(v string) {
	o.UserAlias = &v
}

// GetVersion returns the Version field value if set, zero value otherwise
func (o *ApiHostViewAtlas) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiHostViewAtlas) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}

	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *ApiHostViewAtlas) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *ApiHostViewAtlas) SetVersion(v string) {
	o.Version = &v
}
