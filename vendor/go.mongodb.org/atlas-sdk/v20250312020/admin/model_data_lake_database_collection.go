// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeDatabaseCollection A collection and data sources that map to a “stores“ data store.
type DataLakeDatabaseCollection struct {
	// Array that contains the data stores that map to a collection for this data lake.
	DataSources *[]DataLakeDatabaseDataSourceSettings `json:"dataSources,omitempty"`
	// Human-readable label that identifies the collection to which MongoDB Cloud maps the data in the data stores.
	Name *string `json:"name,omitempty"`
}

// NewDataLakeDatabaseCollection instantiates a new DataLakeDatabaseCollection object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeDatabaseCollection() *DataLakeDatabaseCollection {
	this := DataLakeDatabaseCollection{}
	return &this
}

// NewDataLakeDatabaseCollectionWithDefaults instantiates a new DataLakeDatabaseCollection object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeDatabaseCollectionWithDefaults() *DataLakeDatabaseCollection {
	this := DataLakeDatabaseCollection{}
	return &this
}

// GetDataSources returns the DataSources field value if set, zero value otherwise
func (o *DataLakeDatabaseCollection) GetDataSources() []DataLakeDatabaseDataSourceSettings {
	if o == nil || IsNil(o.DataSources) {
		var ret []DataLakeDatabaseDataSourceSettings
		return ret
	}
	return *o.DataSources
}

// GetDataSourcesOk returns a tuple with the DataSources field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseCollection) GetDataSourcesOk() (*[]DataLakeDatabaseDataSourceSettings, bool) {
	if o == nil || IsNil(o.DataSources) {
		return nil, false
	}

	return o.DataSources, true
}

// HasDataSources returns a boolean if a field has been set.
func (o *DataLakeDatabaseCollection) HasDataSources() bool {
	if o != nil && !IsNil(o.DataSources) {
		return true
	}

	return false
}

// SetDataSources gets a reference to the given []DataLakeDatabaseDataSourceSettings and assigns it to the DataSources field.
func (o *DataLakeDatabaseCollection) SetDataSources(v []DataLakeDatabaseDataSourceSettings) {
	o.DataSources = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DataLakeDatabaseCollection) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseCollection) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataLakeDatabaseCollection) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataLakeDatabaseCollection) SetName(v string) {
	o.Name = &v
}
