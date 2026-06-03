// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiCheckpointPart Metadata contained in one document that describes the complete snapshot taken for this node.
type ApiCheckpointPart struct {
	// Human-readable label that identifies the replica set to which this checkpoint applies.
	// Read only field.
	ReplicaSetName *string `json:"replicaSetName,omitempty"`
	// Human-readable label that identifies the shard to which this checkpoint applies.
	// Read only field.
	ShardName *string `json:"shardName,omitempty"`
	// Flag that indicates whether the token exists.
	// Read only field.
	TokenDiscovered *bool             `json:"tokenDiscovered,omitempty"`
	TokenTimestamp  *ApiBSONTimestamp `json:"tokenTimestamp,omitempty"`
	// Human-readable label that identifies the type of host that the part represents.
	// Read only field.
	TypeName *string `json:"typeName,omitempty"`
}

// NewApiCheckpointPart instantiates a new ApiCheckpointPart object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiCheckpointPart() *ApiCheckpointPart {
	this := ApiCheckpointPart{}
	return &this
}

// NewApiCheckpointPartWithDefaults instantiates a new ApiCheckpointPart object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiCheckpointPartWithDefaults() *ApiCheckpointPart {
	this := ApiCheckpointPart{}
	return &this
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *ApiCheckpointPart) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCheckpointPart) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *ApiCheckpointPart) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *ApiCheckpointPart) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}

// GetShardName returns the ShardName field value if set, zero value otherwise
func (o *ApiCheckpointPart) GetShardName() string {
	if o == nil || IsNil(o.ShardName) {
		var ret string
		return ret
	}
	return *o.ShardName
}

// GetShardNameOk returns a tuple with the ShardName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCheckpointPart) GetShardNameOk() (*string, bool) {
	if o == nil || IsNil(o.ShardName) {
		return nil, false
	}

	return o.ShardName, true
}

// HasShardName returns a boolean if a field has been set.
func (o *ApiCheckpointPart) HasShardName() bool {
	if o != nil && !IsNil(o.ShardName) {
		return true
	}

	return false
}

// SetShardName gets a reference to the given string and assigns it to the ShardName field.
func (o *ApiCheckpointPart) SetShardName(v string) {
	o.ShardName = &v
}

// GetTokenDiscovered returns the TokenDiscovered field value if set, zero value otherwise
func (o *ApiCheckpointPart) GetTokenDiscovered() bool {
	if o == nil || IsNil(o.TokenDiscovered) {
		var ret bool
		return ret
	}
	return *o.TokenDiscovered
}

// GetTokenDiscoveredOk returns a tuple with the TokenDiscovered field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCheckpointPart) GetTokenDiscoveredOk() (*bool, bool) {
	if o == nil || IsNil(o.TokenDiscovered) {
		return nil, false
	}

	return o.TokenDiscovered, true
}

// HasTokenDiscovered returns a boolean if a field has been set.
func (o *ApiCheckpointPart) HasTokenDiscovered() bool {
	if o != nil && !IsNil(o.TokenDiscovered) {
		return true
	}

	return false
}

// SetTokenDiscovered gets a reference to the given bool and assigns it to the TokenDiscovered field.
func (o *ApiCheckpointPart) SetTokenDiscovered(v bool) {
	o.TokenDiscovered = &v
}

// GetTokenTimestamp returns the TokenTimestamp field value if set, zero value otherwise
func (o *ApiCheckpointPart) GetTokenTimestamp() ApiBSONTimestamp {
	if o == nil || IsNil(o.TokenTimestamp) {
		var ret ApiBSONTimestamp
		return ret
	}
	return *o.TokenTimestamp
}

// GetTokenTimestampOk returns a tuple with the TokenTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCheckpointPart) GetTokenTimestampOk() (*ApiBSONTimestamp, bool) {
	if o == nil || IsNil(o.TokenTimestamp) {
		return nil, false
	}

	return o.TokenTimestamp, true
}

// HasTokenTimestamp returns a boolean if a field has been set.
func (o *ApiCheckpointPart) HasTokenTimestamp() bool {
	if o != nil && !IsNil(o.TokenTimestamp) {
		return true
	}

	return false
}

// SetTokenTimestamp gets a reference to the given ApiBSONTimestamp and assigns it to the TokenTimestamp field.
func (o *ApiCheckpointPart) SetTokenTimestamp(v ApiBSONTimestamp) {
	o.TokenTimestamp = &v
}

// GetTypeName returns the TypeName field value if set, zero value otherwise
func (o *ApiCheckpointPart) GetTypeName() string {
	if o == nil || IsNil(o.TypeName) {
		var ret string
		return ret
	}
	return *o.TypeName
}

// GetTypeNameOk returns a tuple with the TypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCheckpointPart) GetTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.TypeName) {
		return nil, false
	}

	return o.TypeName, true
}

// HasTypeName returns a boolean if a field has been set.
func (o *ApiCheckpointPart) HasTypeName() bool {
	if o != nil && !IsNil(o.TypeName) {
		return true
	}

	return false
}

// SetTypeName gets a reference to the given string and assigns it to the TypeName field.
func (o *ApiCheckpointPart) SetTypeName(v string) {
	o.TypeName = &v
}
