// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeStorage Configuration information for each data store and its mapping to MongoDB Cloud databases.
type DataLakeStorage struct {
	// Array that contains the queryable databases and collections for this data lake.
	Databases *[]DataLakeDatabaseInstance `json:"databases,omitempty"`
	// Array that contains the data stores for the data lake.
	Stores *[]DataLakeStoreSettings `json:"stores,omitempty"`
}

// NewDataLakeStorage instantiates a new DataLakeStorage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeStorage() *DataLakeStorage {
	this := DataLakeStorage{}
	return &this
}

// NewDataLakeStorageWithDefaults instantiates a new DataLakeStorage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeStorageWithDefaults() *DataLakeStorage {
	this := DataLakeStorage{}
	return &this
}

// GetDatabases returns the Databases field value if set, zero value otherwise
func (o *DataLakeStorage) GetDatabases() []DataLakeDatabaseInstance {
	if o == nil || IsNil(o.Databases) {
		var ret []DataLakeDatabaseInstance
		return ret
	}
	return *o.Databases
}

// GetDatabasesOk returns a tuple with the Databases field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStorage) GetDatabasesOk() (*[]DataLakeDatabaseInstance, bool) {
	if o == nil || IsNil(o.Databases) {
		return nil, false
	}

	return o.Databases, true
}

// HasDatabases returns a boolean if a field has been set.
func (o *DataLakeStorage) HasDatabases() bool {
	if o != nil && !IsNil(o.Databases) {
		return true
	}

	return false
}

// SetDatabases gets a reference to the given []DataLakeDatabaseInstance and assigns it to the Databases field.
func (o *DataLakeStorage) SetDatabases(v []DataLakeDatabaseInstance) {
	o.Databases = &v
}

// GetStores returns the Stores field value if set, zero value otherwise
func (o *DataLakeStorage) GetStores() []DataLakeStoreSettings {
	if o == nil || IsNil(o.Stores) {
		var ret []DataLakeStoreSettings
		return ret
	}
	return *o.Stores
}

// GetStoresOk returns a tuple with the Stores field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStorage) GetStoresOk() (*[]DataLakeStoreSettings, bool) {
	if o == nil || IsNil(o.Stores) {
		return nil, false
	}

	return o.Stores, true
}

// HasStores returns a boolean if a field has been set.
func (o *DataLakeStorage) HasStores() bool {
	if o != nil && !IsNil(o.Stores) {
		return true
	}

	return false
}

// SetStores gets a reference to the given []DataLakeStoreSettings and assigns it to the Stores field.
func (o *DataLakeStorage) SetStores(v []DataLakeStoreSettings) {
	o.Stores = &v
}
