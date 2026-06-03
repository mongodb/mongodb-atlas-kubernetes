// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ShardKeys Document that configures the shard key on the destination cluster.
type ShardKeys struct {
	// List of fields to use for the shard key.
	// Write only field.
	Key []any `json:"key"`
}

// NewShardKeys instantiates a new ShardKeys object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewShardKeys(key []any) *ShardKeys {
	this := ShardKeys{}
	this.Key = key
	return &this
}

// NewShardKeysWithDefaults instantiates a new ShardKeys object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewShardKeysWithDefaults() *ShardKeys {
	this := ShardKeys{}
	return &this
}

// GetKey returns the Key field value
func (o *ShardKeys) GetKey() []any {
	if o == nil {
		var ret []any
		return ret
	}

	return o.Key
}

// GetKeyOk returns a tuple with the Key field value
// and a boolean to check if the value has been set.
func (o *ShardKeys) GetKeyOk() (*[]any, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Key, true
}

// SetKey sets field value
func (o *ShardKeys) SetKey(v []any) {
	o.Key = v
}
