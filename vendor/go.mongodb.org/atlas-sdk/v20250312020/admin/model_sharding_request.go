// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ShardingRequest Document that configures sharding on the destination cluster when migrating from a replica set source to a sharded cluster destination on MongoDB 6.0 or higher. If you don't wish to shard any collections on the destination cluster, leave this empty.
type ShardingRequest struct {
	// Flag that lets the migration create supporting indexes for the shard keys, if none exists, as the destination cluster also needs compatible indexes for the specified shard keys.
	// Write only field.
	CreateSupportingIndexes bool `json:"createSupportingIndexes"`
	// List of shard configurations to shard destination collections. Atlas shards only those collections that you include in the sharding entries array.
	// Write only field.
	ShardingEntries []ShardEntry `json:"shardingEntries"`
}

// NewShardingRequest instantiates a new ShardingRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewShardingRequest(createSupportingIndexes bool, shardingEntries []ShardEntry) *ShardingRequest {
	this := ShardingRequest{}
	this.CreateSupportingIndexes = createSupportingIndexes
	this.ShardingEntries = shardingEntries
	return &this
}

// NewShardingRequestWithDefaults instantiates a new ShardingRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewShardingRequestWithDefaults() *ShardingRequest {
	this := ShardingRequest{}
	return &this
}

// GetCreateSupportingIndexes returns the CreateSupportingIndexes field value
func (o *ShardingRequest) GetCreateSupportingIndexes() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.CreateSupportingIndexes
}

// GetCreateSupportingIndexesOk returns a tuple with the CreateSupportingIndexes field value
// and a boolean to check if the value has been set.
func (o *ShardingRequest) GetCreateSupportingIndexesOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreateSupportingIndexes, true
}

// SetCreateSupportingIndexes sets field value
func (o *ShardingRequest) SetCreateSupportingIndexes(v bool) {
	o.CreateSupportingIndexes = v
}

// GetShardingEntries returns the ShardingEntries field value
func (o *ShardingRequest) GetShardingEntries() []ShardEntry {
	if o == nil {
		var ret []ShardEntry
		return ret
	}

	return o.ShardingEntries
}

// GetShardingEntriesOk returns a tuple with the ShardingEntries field value
// and a boolean to check if the value has been set.
func (o *ShardingRequest) GetShardingEntriesOk() (*[]ShardEntry, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ShardingEntries, true
}

// SetShardingEntries sets field value
func (o *ShardingRequest) SetShardingEntries(v []ShardEntry) {
	o.ShardingEntries = v
}
