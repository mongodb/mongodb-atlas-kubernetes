// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeDatabaseDataSourceSettings Data store that maps to a collection for this data lake.
type DataLakeDatabaseDataSourceSettings struct {
	// Flag that validates the scheme in the specified URLs. If `true`, allows insecure `HTTP` scheme, doesn't verify the server's certificate chain and hostname, and accepts any certificate with any hostname presented by the server. If `false`, allows secure `HTTPS` scheme only.
	AllowInsecure *bool `json:"allowInsecure,omitempty"`
	// Human-readable label that identifies the collection in the database. For creating a wildcard (`*`) collection, you must omit this parameter.
	Collection *string `json:"collection,omitempty"`
	// Regex pattern to use for creating the wildcard (*) collection. To learn more about the regex syntax, see [Go programming language](https://pkg.go.dev/regexp).
	CollectionRegex *string `json:"collectionRegex,omitempty"`
	// Human-readable label that identifies the database, which contains the collection in the cluster. You must omit this parameter to generate wildcard (`*`) collections for dynamically generated databases.
	Database *string `json:"database,omitempty"`
	// Regex pattern to use for creating the wildcard (*) database. To learn more about the regex syntax, see [Go programming language](https://pkg.go.dev/regexp).
	DatabaseRegex *string `json:"databaseRegex,omitempty"`
	// Human-readable label that identifies the dataset that Atlas generates for an ingestion pipeline run or Online Archive.
	DatasetName *string `json:"datasetName,omitempty"`
	// Human-readable label that matches against the dataset names for ingestion pipeline runs or Online Archives.
	DatasetPrefix *string `json:"datasetPrefix,omitempty"`
	// File format that MongoDB Cloud uses if it encounters a file without a file extension while searching **storeName**.
	DefaultFormat *string `json:"defaultFormat,omitempty"`
	// File path that controls how MongoDB Cloud searches for and parses files in the **storeName** before mapping them to a collection.Specify ``/`` to capture all files and folders from the ``prefix`` path.
	Path *string `json:"path,omitempty"`
	// Name for the field that includes the provenance of the documents in the results. MongoDB Cloud returns different fields in the results for each supported provider.
	ProvenanceFieldName *string `json:"provenanceFieldName,omitempty"`
	// Human-readable label that identifies the data store that MongoDB Cloud maps to the collection.
	StoreName *string `json:"storeName,omitempty"`
	// Unsigned integer that specifies how many fields of the dataset name to trim from the left of the dataset name before mapping the remaining fields to a wildcard collection name.
	TrimLevel *int `json:"trimLevel,omitempty"`
	// URLs of the publicly accessible data files. You can't specify URLs that require authentication. Atlas Data Lake creates a partition for each URL. If empty or omitted, Data Lake uses the URLs from the store specified in the **dataSources.storeName** parameter.
	Urls *[]string `json:"urls,omitempty"`
}

// NewDataLakeDatabaseDataSourceSettings instantiates a new DataLakeDatabaseDataSourceSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeDatabaseDataSourceSettings() *DataLakeDatabaseDataSourceSettings {
	this := DataLakeDatabaseDataSourceSettings{}
	var allowInsecure bool = false
	this.AllowInsecure = &allowInsecure
	return &this
}

// NewDataLakeDatabaseDataSourceSettingsWithDefaults instantiates a new DataLakeDatabaseDataSourceSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeDatabaseDataSourceSettingsWithDefaults() *DataLakeDatabaseDataSourceSettings {
	this := DataLakeDatabaseDataSourceSettings{}
	var allowInsecure bool = false
	this.AllowInsecure = &allowInsecure
	return &this
}

// GetAllowInsecure returns the AllowInsecure field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetAllowInsecure() bool {
	if o == nil || IsNil(o.AllowInsecure) {
		var ret bool
		return ret
	}
	return *o.AllowInsecure
}

// GetAllowInsecureOk returns a tuple with the AllowInsecure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetAllowInsecureOk() (*bool, bool) {
	if o == nil || IsNil(o.AllowInsecure) {
		return nil, false
	}

	return o.AllowInsecure, true
}

// HasAllowInsecure returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasAllowInsecure() bool {
	if o != nil && !IsNil(o.AllowInsecure) {
		return true
	}

	return false
}

// SetAllowInsecure gets a reference to the given bool and assigns it to the AllowInsecure field.
func (o *DataLakeDatabaseDataSourceSettings) SetAllowInsecure(v bool) {
	o.AllowInsecure = &v
}

// GetCollection returns the Collection field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetCollection() string {
	if o == nil || IsNil(o.Collection) {
		var ret string
		return ret
	}
	return *o.Collection
}

// GetCollectionOk returns a tuple with the Collection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetCollectionOk() (*string, bool) {
	if o == nil || IsNil(o.Collection) {
		return nil, false
	}

	return o.Collection, true
}

// HasCollection returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasCollection() bool {
	if o != nil && !IsNil(o.Collection) {
		return true
	}

	return false
}

// SetCollection gets a reference to the given string and assigns it to the Collection field.
func (o *DataLakeDatabaseDataSourceSettings) SetCollection(v string) {
	o.Collection = &v
}

// GetCollectionRegex returns the CollectionRegex field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetCollectionRegex() string {
	if o == nil || IsNil(o.CollectionRegex) {
		var ret string
		return ret
	}
	return *o.CollectionRegex
}

// GetCollectionRegexOk returns a tuple with the CollectionRegex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetCollectionRegexOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionRegex) {
		return nil, false
	}

	return o.CollectionRegex, true
}

// HasCollectionRegex returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasCollectionRegex() bool {
	if o != nil && !IsNil(o.CollectionRegex) {
		return true
	}

	return false
}

// SetCollectionRegex gets a reference to the given string and assigns it to the CollectionRegex field.
func (o *DataLakeDatabaseDataSourceSettings) SetCollectionRegex(v string) {
	o.CollectionRegex = &v
}

// GetDatabase returns the Database field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetDatabase() string {
	if o == nil || IsNil(o.Database) {
		var ret string
		return ret
	}
	return *o.Database
}

// GetDatabaseOk returns a tuple with the Database field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetDatabaseOk() (*string, bool) {
	if o == nil || IsNil(o.Database) {
		return nil, false
	}

	return o.Database, true
}

// HasDatabase returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasDatabase() bool {
	if o != nil && !IsNil(o.Database) {
		return true
	}

	return false
}

// SetDatabase gets a reference to the given string and assigns it to the Database field.
func (o *DataLakeDatabaseDataSourceSettings) SetDatabase(v string) {
	o.Database = &v
}

// GetDatabaseRegex returns the DatabaseRegex field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetDatabaseRegex() string {
	if o == nil || IsNil(o.DatabaseRegex) {
		var ret string
		return ret
	}
	return *o.DatabaseRegex
}

// GetDatabaseRegexOk returns a tuple with the DatabaseRegex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetDatabaseRegexOk() (*string, bool) {
	if o == nil || IsNil(o.DatabaseRegex) {
		return nil, false
	}

	return o.DatabaseRegex, true
}

// HasDatabaseRegex returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasDatabaseRegex() bool {
	if o != nil && !IsNil(o.DatabaseRegex) {
		return true
	}

	return false
}

// SetDatabaseRegex gets a reference to the given string and assigns it to the DatabaseRegex field.
func (o *DataLakeDatabaseDataSourceSettings) SetDatabaseRegex(v string) {
	o.DatabaseRegex = &v
}

// GetDatasetName returns the DatasetName field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetDatasetName() string {
	if o == nil || IsNil(o.DatasetName) {
		var ret string
		return ret
	}
	return *o.DatasetName
}

// GetDatasetNameOk returns a tuple with the DatasetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetDatasetNameOk() (*string, bool) {
	if o == nil || IsNil(o.DatasetName) {
		return nil, false
	}

	return o.DatasetName, true
}

// HasDatasetName returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasDatasetName() bool {
	if o != nil && !IsNil(o.DatasetName) {
		return true
	}

	return false
}

// SetDatasetName gets a reference to the given string and assigns it to the DatasetName field.
func (o *DataLakeDatabaseDataSourceSettings) SetDatasetName(v string) {
	o.DatasetName = &v
}

// GetDatasetPrefix returns the DatasetPrefix field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetDatasetPrefix() string {
	if o == nil || IsNil(o.DatasetPrefix) {
		var ret string
		return ret
	}
	return *o.DatasetPrefix
}

// GetDatasetPrefixOk returns a tuple with the DatasetPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetDatasetPrefixOk() (*string, bool) {
	if o == nil || IsNil(o.DatasetPrefix) {
		return nil, false
	}

	return o.DatasetPrefix, true
}

// HasDatasetPrefix returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasDatasetPrefix() bool {
	if o != nil && !IsNil(o.DatasetPrefix) {
		return true
	}

	return false
}

// SetDatasetPrefix gets a reference to the given string and assigns it to the DatasetPrefix field.
func (o *DataLakeDatabaseDataSourceSettings) SetDatasetPrefix(v string) {
	o.DatasetPrefix = &v
}

// GetDefaultFormat returns the DefaultFormat field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetDefaultFormat() string {
	if o == nil || IsNil(o.DefaultFormat) {
		var ret string
		return ret
	}
	return *o.DefaultFormat
}

// GetDefaultFormatOk returns a tuple with the DefaultFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetDefaultFormatOk() (*string, bool) {
	if o == nil || IsNil(o.DefaultFormat) {
		return nil, false
	}

	return o.DefaultFormat, true
}

// HasDefaultFormat returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasDefaultFormat() bool {
	if o != nil && !IsNil(o.DefaultFormat) {
		return true
	}

	return false
}

// SetDefaultFormat gets a reference to the given string and assigns it to the DefaultFormat field.
func (o *DataLakeDatabaseDataSourceSettings) SetDefaultFormat(v string) {
	o.DefaultFormat = &v
}

// GetPath returns the Path field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetPath() string {
	if o == nil || IsNil(o.Path) {
		var ret string
		return ret
	}
	return *o.Path
}

// GetPathOk returns a tuple with the Path field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetPathOk() (*string, bool) {
	if o == nil || IsNil(o.Path) {
		return nil, false
	}

	return o.Path, true
}

// HasPath returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasPath() bool {
	if o != nil && !IsNil(o.Path) {
		return true
	}

	return false
}

// SetPath gets a reference to the given string and assigns it to the Path field.
func (o *DataLakeDatabaseDataSourceSettings) SetPath(v string) {
	o.Path = &v
}

// GetProvenanceFieldName returns the ProvenanceFieldName field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetProvenanceFieldName() string {
	if o == nil || IsNil(o.ProvenanceFieldName) {
		var ret string
		return ret
	}
	return *o.ProvenanceFieldName
}

// GetProvenanceFieldNameOk returns a tuple with the ProvenanceFieldName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetProvenanceFieldNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProvenanceFieldName) {
		return nil, false
	}

	return o.ProvenanceFieldName, true
}

// HasProvenanceFieldName returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasProvenanceFieldName() bool {
	if o != nil && !IsNil(o.ProvenanceFieldName) {
		return true
	}

	return false
}

// SetProvenanceFieldName gets a reference to the given string and assigns it to the ProvenanceFieldName field.
func (o *DataLakeDatabaseDataSourceSettings) SetProvenanceFieldName(v string) {
	o.ProvenanceFieldName = &v
}

// GetStoreName returns the StoreName field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetStoreName() string {
	if o == nil || IsNil(o.StoreName) {
		var ret string
		return ret
	}
	return *o.StoreName
}

// GetStoreNameOk returns a tuple with the StoreName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetStoreNameOk() (*string, bool) {
	if o == nil || IsNil(o.StoreName) {
		return nil, false
	}

	return o.StoreName, true
}

// HasStoreName returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasStoreName() bool {
	if o != nil && !IsNil(o.StoreName) {
		return true
	}

	return false
}

// SetStoreName gets a reference to the given string and assigns it to the StoreName field.
func (o *DataLakeDatabaseDataSourceSettings) SetStoreName(v string) {
	o.StoreName = &v
}

// GetTrimLevel returns the TrimLevel field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetTrimLevel() int {
	if o == nil || IsNil(o.TrimLevel) {
		var ret int
		return ret
	}
	return *o.TrimLevel
}

// GetTrimLevelOk returns a tuple with the TrimLevel field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetTrimLevelOk() (*int, bool) {
	if o == nil || IsNil(o.TrimLevel) {
		return nil, false
	}

	return o.TrimLevel, true
}

// HasTrimLevel returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasTrimLevel() bool {
	if o != nil && !IsNil(o.TrimLevel) {
		return true
	}

	return false
}

// SetTrimLevel gets a reference to the given int and assigns it to the TrimLevel field.
func (o *DataLakeDatabaseDataSourceSettings) SetTrimLevel(v int) {
	o.TrimLevel = &v
}

// GetUrls returns the Urls field value if set, zero value otherwise
func (o *DataLakeDatabaseDataSourceSettings) GetUrls() []string {
	if o == nil || IsNil(o.Urls) {
		var ret []string
		return ret
	}
	return *o.Urls
}

// GetUrlsOk returns a tuple with the Urls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeDatabaseDataSourceSettings) GetUrlsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Urls) {
		return nil, false
	}

	return o.Urls, true
}

// HasUrls returns a boolean if a field has been set.
func (o *DataLakeDatabaseDataSourceSettings) HasUrls() bool {
	if o != nil && !IsNil(o.Urls) {
		return true
	}

	return false
}

// SetUrls gets a reference to the given []string and assigns it to the Urls field.
func (o *DataLakeDatabaseDataSourceSettings) SetUrls(v []string) {
	o.Urls = &v
}
