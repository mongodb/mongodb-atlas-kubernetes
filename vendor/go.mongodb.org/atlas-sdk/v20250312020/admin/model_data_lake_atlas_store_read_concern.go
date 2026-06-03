// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeAtlasStoreReadConcern MongoDB Cloud cluster read concern, which determines the consistency and isolation properties of the data read from an Atlas cluster.
type DataLakeAtlasStoreReadConcern struct {
	// Read Concern level that specifies the consistency and availability of the data read.
	Level *string `json:"level,omitempty"`
}

// NewDataLakeAtlasStoreReadConcern instantiates a new DataLakeAtlasStoreReadConcern object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeAtlasStoreReadConcern() *DataLakeAtlasStoreReadConcern {
	this := DataLakeAtlasStoreReadConcern{}
	return &this
}

// NewDataLakeAtlasStoreReadConcernWithDefaults instantiates a new DataLakeAtlasStoreReadConcern object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeAtlasStoreReadConcernWithDefaults() *DataLakeAtlasStoreReadConcern {
	this := DataLakeAtlasStoreReadConcern{}
	return &this
}

// GetLevel returns the Level field value if set, zero value otherwise
func (o *DataLakeAtlasStoreReadConcern) GetLevel() string {
	if o == nil || IsNil(o.Level) {
		var ret string
		return ret
	}
	return *o.Level
}

// GetLevelOk returns a tuple with the Level field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeAtlasStoreReadConcern) GetLevelOk() (*string, bool) {
	if o == nil || IsNil(o.Level) {
		return nil, false
	}

	return o.Level, true
}

// HasLevel returns a boolean if a field has been set.
func (o *DataLakeAtlasStoreReadConcern) HasLevel() bool {
	if o != nil && !IsNil(o.Level) {
		return true
	}

	return false
}

// SetLevel gets a reference to the given string and assigns it to the Level field.
func (o *DataLakeAtlasStoreReadConcern) SetLevel(v string) {
	o.Level = &v
}
