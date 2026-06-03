// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DatabaseRollingIndexRequest struct for DatabaseRollingIndexRequest
type DatabaseRollingIndexRequest struct {
	Collation *Collation `json:"collation,omitempty"`
	// Human-readable label of the collection for which MongoDB Cloud creates an index.
	// Write only field.
	Collection string `json:"collection"`
	// Human-readable label of the database that holds the collection on which MongoDB Cloud creates an index.
	// Write only field.
	Db string `json:"db"`
	// List that contains one or more objects that describe the parameters that you want to index.
	// Write only field.
	Keys    []map[string]string `json:"keys"`
	Options *IndexOptions       `json:"options,omitempty"`
}

// NewDatabaseRollingIndexRequest instantiates a new DatabaseRollingIndexRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatabaseRollingIndexRequest(collection string, db string, keys []map[string]string) *DatabaseRollingIndexRequest {
	this := DatabaseRollingIndexRequest{}
	this.Collection = collection
	this.Db = db
	this.Keys = keys
	return &this
}

// NewDatabaseRollingIndexRequestWithDefaults instantiates a new DatabaseRollingIndexRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatabaseRollingIndexRequestWithDefaults() *DatabaseRollingIndexRequest {
	this := DatabaseRollingIndexRequest{}
	return &this
}

// GetCollation returns the Collation field value if set, zero value otherwise
func (o *DatabaseRollingIndexRequest) GetCollation() Collation {
	if o == nil || IsNil(o.Collation) {
		var ret Collation
		return ret
	}
	return *o.Collation
}

// GetCollationOk returns a tuple with the Collation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseRollingIndexRequest) GetCollationOk() (*Collation, bool) {
	if o == nil || IsNil(o.Collation) {
		return nil, false
	}

	return o.Collation, true
}

// HasCollation returns a boolean if a field has been set.
func (o *DatabaseRollingIndexRequest) HasCollation() bool {
	if o != nil && !IsNil(o.Collation) {
		return true
	}

	return false
}

// SetCollation gets a reference to the given Collation and assigns it to the Collation field.
func (o *DatabaseRollingIndexRequest) SetCollation(v Collation) {
	o.Collation = &v
}

// GetCollection returns the Collection field value
func (o *DatabaseRollingIndexRequest) GetCollection() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Collection
}

// GetCollectionOk returns a tuple with the Collection field value
// and a boolean to check if the value has been set.
func (o *DatabaseRollingIndexRequest) GetCollectionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Collection, true
}

// SetCollection sets field value
func (o *DatabaseRollingIndexRequest) SetCollection(v string) {
	o.Collection = v
}

// GetDb returns the Db field value
func (o *DatabaseRollingIndexRequest) GetDb() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Db
}

// GetDbOk returns a tuple with the Db field value
// and a boolean to check if the value has been set.
func (o *DatabaseRollingIndexRequest) GetDbOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Db, true
}

// SetDb sets field value
func (o *DatabaseRollingIndexRequest) SetDb(v string) {
	o.Db = v
}

// GetKeys returns the Keys field value
func (o *DatabaseRollingIndexRequest) GetKeys() []map[string]string {
	if o == nil {
		var ret []map[string]string
		return ret
	}

	return o.Keys
}

// GetKeysOk returns a tuple with the Keys field value
// and a boolean to check if the value has been set.
func (o *DatabaseRollingIndexRequest) GetKeysOk() (*[]map[string]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Keys, true
}

// SetKeys sets field value
func (o *DatabaseRollingIndexRequest) SetKeys(v []map[string]string) {
	o.Keys = v
}

// GetOptions returns the Options field value if set, zero value otherwise
func (o *DatabaseRollingIndexRequest) GetOptions() IndexOptions {
	if o == nil || IsNil(o.Options) {
		var ret IndexOptions
		return ret
	}
	return *o.Options
}

// GetOptionsOk returns a tuple with the Options field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseRollingIndexRequest) GetOptionsOk() (*IndexOptions, bool) {
	if o == nil || IsNil(o.Options) {
		return nil, false
	}

	return o.Options, true
}

// HasOptions returns a boolean if a field has been set.
func (o *DatabaseRollingIndexRequest) HasOptions() bool {
	if o != nil && !IsNil(o.Options) {
		return true
	}

	return false
}

// SetOptions gets a reference to the given IndexOptions and assigns it to the Options field.
func (o *DatabaseRollingIndexRequest) SetOptions(v IndexOptions) {
	o.Options = &v
}
