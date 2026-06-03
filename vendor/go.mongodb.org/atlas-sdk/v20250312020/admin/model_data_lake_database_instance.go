// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeDatabaseInstance Database associated with this data lake. Databases contain collections and views.
type DataLakeDatabaseInstance struct {
	// Array of collections and data sources that map to a ``stores`` data store.
	Collections *[]DataLakeDatabaseCollection `json:"collections,omitempty"`
	// Maximum number of wildcard collections in the database. This only applies to S3 data sources.
	MaxWildcardCollections *int `json:"maxWildcardCollections,omitempty"`
	// Human-readable label that identifies the database to which the data lake maps data.
	Name *string `json:"name,omitempty"`
	// Array of aggregation pipelines that apply to the collection. This only applies to S3 data sources.
	Views *[]DataLakeApiBase `json:"views,omitempty"`
}

// NewDataLakeDatabaseInstance instantiates a new DataLakeDatabaseInstance object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeDatabaseInstance() *DataLakeDatabaseInstance {
	this := DataLakeDatabaseInstance{}
	var maxWildcardCollections int = 100
	this.MaxWildcardCollections = &maxWildcardCollections
	return &this
}

// NewDataLakeDatabaseInstanceWithDefaults instantiates a new DataLakeDatabaseInstance object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeDatabaseInstanceWithDefaults() *DataLakeDatabaseInstance {
	this := DataLakeDatabaseInstance{}
	var maxWildcardCollections int = 100
	this.MaxWildcardCollections = &maxWildcardCollections
	return &this
}

// GetCollections returns the Collections field value if set, zero value otherwise
func (o *DataLakeDatabaseInstance) GetCollections() []DataLakeDatabaseCollection {
	if o == nil || IsNil(o.Collections) {
		var ret []DataLakeDatabaseCollection
		return ret
	}
	return *o.Collections
}

// GetCollectionsOk returns a tuple with the Collections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseInstance) GetCollectionsOk() (*[]DataLakeDatabaseCollection, bool) {
	if o == nil || IsNil(o.Collections) {
		return nil, false
	}

	return o.Collections, true
}

// HasCollections returns a boolean if a field has been set.
func (o *DataLakeDatabaseInstance) HasCollections() bool {
	if o != nil && !IsNil(o.Collections) {
		return true
	}

	return false
}

// SetCollections gets a reference to the given []DataLakeDatabaseCollection and assigns it to the Collections field.
func (o *DataLakeDatabaseInstance) SetCollections(v []DataLakeDatabaseCollection) {
	o.Collections = &v
}

// GetMaxWildcardCollections returns the MaxWildcardCollections field value if set, zero value otherwise
func (o *DataLakeDatabaseInstance) GetMaxWildcardCollections() int {
	if o == nil || IsNil(o.MaxWildcardCollections) {
		var ret int
		return ret
	}
	return *o.MaxWildcardCollections
}

// GetMaxWildcardCollectionsOk returns a tuple with the MaxWildcardCollections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseInstance) GetMaxWildcardCollectionsOk() (*int, bool) {
	if o == nil || IsNil(o.MaxWildcardCollections) {
		return nil, false
	}

	return o.MaxWildcardCollections, true
}

// HasMaxWildcardCollections returns a boolean if a field has been set.
func (o *DataLakeDatabaseInstance) HasMaxWildcardCollections() bool {
	if o != nil && !IsNil(o.MaxWildcardCollections) {
		return true
	}

	return false
}

// SetMaxWildcardCollections gets a reference to the given int and assigns it to the MaxWildcardCollections field.
func (o *DataLakeDatabaseInstance) SetMaxWildcardCollections(v int) {
	o.MaxWildcardCollections = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DataLakeDatabaseInstance) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseInstance) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataLakeDatabaseInstance) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataLakeDatabaseInstance) SetName(v string) {
	o.Name = &v
}

// GetViews returns the Views field value if set, zero value otherwise
func (o *DataLakeDatabaseInstance) GetViews() []DataLakeApiBase {
	if o == nil || IsNil(o.Views) {
		var ret []DataLakeApiBase
		return ret
	}
	return *o.Views
}

// GetViewsOk returns a tuple with the Views field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseInstance) GetViewsOk() (*[]DataLakeApiBase, bool) {
	if o == nil || IsNil(o.Views) {
		return nil, false
	}

	return o.Views, true
}

// HasViews returns a boolean if a field has been set.
func (o *DataLakeDatabaseInstance) HasViews() bool {
	if o != nil && !IsNil(o.Views) {
		return true
	}

	return false
}

// SetViews gets a reference to the given []DataLakeApiBase and assigns it to the Views field.
func (o *DataLakeDatabaseInstance) SetViews(v []DataLakeApiBase) {
	o.Views = &v
}
