// Code based on the AtlasAPI V2 OpenAPI file

package admin

// IngestionSource Ingestion Source of a Data Lake Pipeline.
type IngestionSource struct {
	// Type of ingestion source of this Data Lake Pipeline.
	Type *string `json:"type,omitempty"`
	// Human-readable name that identifies the cluster.
	ClusterName *string `json:"clusterName,omitempty"`
	// Human-readable name that identifies the collection.
	CollectionName *string `json:"collectionName,omitempty"`
	// Human-readable name that identifies the database.
	DatabaseName *string `json:"databaseName,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal character string that identifies a policy item.
	PolicyItemId *string `json:"policyItemId,omitempty"`
}

// NewIngestionSource instantiates a new IngestionSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIngestionSource() *IngestionSource {
	this := IngestionSource{}
	return &this
}

// NewIngestionSourceWithDefaults instantiates a new IngestionSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIngestionSourceWithDefaults() *IngestionSource {
	this := IngestionSource{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise
func (o *IngestionSource) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSource) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *IngestionSource) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *IngestionSource) SetType(v string) {
	o.Type = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *IngestionSource) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSource) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *IngestionSource) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *IngestionSource) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCollectionName returns the CollectionName field value if set, zero value otherwise
func (o *IngestionSource) GetCollectionName() string {
	if o == nil || IsNil(o.CollectionName) {
		var ret string
		return ret
	}
	return *o.CollectionName
}

// GetCollectionNameOk returns a tuple with the CollectionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSource) GetCollectionNameOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionName) {
		return nil, false
	}

	return o.CollectionName, true
}

// HasCollectionName returns a boolean if a field has been set.
func (o *IngestionSource) HasCollectionName() bool {
	if o != nil && !IsNil(o.CollectionName) {
		return true
	}

	return false
}

// SetCollectionName gets a reference to the given string and assigns it to the CollectionName field.
func (o *IngestionSource) SetCollectionName(v string) {
	o.CollectionName = &v
}

// GetDatabaseName returns the DatabaseName field value if set, zero value otherwise
func (o *IngestionSource) GetDatabaseName() string {
	if o == nil || IsNil(o.DatabaseName) {
		var ret string
		return ret
	}
	return *o.DatabaseName
}

// GetDatabaseNameOk returns a tuple with the DatabaseName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSource) GetDatabaseNameOk() (*string, bool) {
	if o == nil || IsNil(o.DatabaseName) {
		return nil, false
	}

	return o.DatabaseName, true
}

// HasDatabaseName returns a boolean if a field has been set.
func (o *IngestionSource) HasDatabaseName() bool {
	if o != nil && !IsNil(o.DatabaseName) {
		return true
	}

	return false
}

// SetDatabaseName gets a reference to the given string and assigns it to the DatabaseName field.
func (o *IngestionSource) SetDatabaseName(v string) {
	o.DatabaseName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *IngestionSource) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSource) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *IngestionSource) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *IngestionSource) SetGroupId(v string) {
	o.GroupId = &v
}

// GetPolicyItemId returns the PolicyItemId field value if set, zero value otherwise
func (o *IngestionSource) GetPolicyItemId() string {
	if o == nil || IsNil(o.PolicyItemId) {
		var ret string
		return ret
	}
	return *o.PolicyItemId
}

// GetPolicyItemIdOk returns a tuple with the PolicyItemId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IngestionSource) GetPolicyItemIdOk() (*string, bool) {
	if o == nil || IsNil(o.PolicyItemId) {
		return nil, false
	}

	return o.PolicyItemId, true
}

// HasPolicyItemId returns a boolean if a field has been set.
func (o *IngestionSource) HasPolicyItemId() bool {
	if o != nil && !IsNil(o.PolicyItemId) {
		return true
	}

	return false
}

// SetPolicyItemId gets a reference to the given string and assigns it to the PolicyItemId field.
func (o *IngestionSource) SetPolicyItemId(v string) {
	o.PolicyItemId = &v
}
