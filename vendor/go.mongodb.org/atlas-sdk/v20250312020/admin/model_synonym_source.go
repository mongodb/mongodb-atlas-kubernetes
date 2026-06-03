// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SynonymSource Data set that stores words and their applicable synonyms.
type SynonymSource struct {
	// Label that identifies the MongoDB collection that stores words and their applicable synonyms.
	Collection string `json:"collection"`
}

// NewSynonymSource instantiates a new SynonymSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSynonymSource(collection string) *SynonymSource {
	this := SynonymSource{}
	this.Collection = collection
	return &this
}

// NewSynonymSourceWithDefaults instantiates a new SynonymSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSynonymSourceWithDefaults() *SynonymSource {
	this := SynonymSource{}
	return &this
}

// GetCollection returns the Collection field value
func (o *SynonymSource) GetCollection() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Collection
}

// GetCollectionOk returns a tuple with the Collection field value
// and a boolean to check if the value has been set.
func (o *SynonymSource) GetCollectionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Collection, true
}

// SetCollection sets field value
func (o *SynonymSource) SetCollection(v string) {
	o.Collection = v
}
