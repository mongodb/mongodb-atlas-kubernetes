// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ShardEntry Sharding configuration for a collection to be sharded on the destination cluster.
type ShardEntry struct {
	// Human-readable label that identifies the collection to be sharded on the destination cluster.
	// Write only field.
	Collection string `json:"collection"`
	// Human-readable label that identifies the database that contains the collection to be sharded on the destination cluster.
	// Write only field.
	Database        string    `json:"database"`
	ShardCollection ShardKeys `json:"shardCollection"`
}

// NewShardEntry instantiates a new ShardEntry object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewShardEntry(collection string, database string, shardCollection ShardKeys) *ShardEntry {
	this := ShardEntry{}
	this.Collection = collection
	this.Database = database
	this.ShardCollection = shardCollection
	return &this
}

// NewShardEntryWithDefaults instantiates a new ShardEntry object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewShardEntryWithDefaults() *ShardEntry {
	this := ShardEntry{}
	return &this
}

// GetCollection returns the Collection field value
func (o *ShardEntry) GetCollection() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Collection
}

// GetCollectionOk returns a tuple with the Collection field value
// and a boolean to check if the value has been set.
func (o *ShardEntry) GetCollectionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Collection, true
}

// SetCollection sets field value
func (o *ShardEntry) SetCollection(v string) {
	o.Collection = v
}

// GetDatabase returns the Database field value
func (o *ShardEntry) GetDatabase() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Database
}

// GetDatabaseOk returns a tuple with the Database field value
// and a boolean to check if the value has been set.
func (o *ShardEntry) GetDatabaseOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Database, true
}

// SetDatabase sets field value
func (o *ShardEntry) SetDatabase(v string) {
	o.Database = v
}

// GetShardCollection returns the ShardCollection field value
func (o *ShardEntry) GetShardCollection() ShardKeys {
	if o == nil {
		var ret ShardKeys
		return ret
	}

	return o.ShardCollection
}

// GetShardCollectionOk returns a tuple with the ShardCollection field value
// and a boolean to check if the value has been set.
func (o *ShardEntry) GetShardCollectionOk() (*ShardKeys, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ShardCollection, true
}

// SetShardCollection sets field value
func (o *ShardEntry) SetShardCollection(v ShardKeys) {
	o.ShardCollection = v
}
