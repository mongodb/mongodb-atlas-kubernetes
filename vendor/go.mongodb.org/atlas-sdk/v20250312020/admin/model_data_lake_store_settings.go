// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeStoreSettings Group of settings that define where the data is stored.
type DataLakeStoreSettings struct {
	// Human-readable label that identifies the data store. The **databases.[n].collections.[n].dataSources.[n].storeName** field references this values as part of the mapping configuration. To use MongoDB Cloud as a data store, the data lake requires a serverless instance or an `M10` or higher cluster.
	Name     *string `json:"name,omitempty"`
	Provider string  `json:"provider"`
	// Collection of AWS S3 [storage classes](https://aws.amazon.com/s3/storage-classes/). Atlas Data Lake includes the files in these storage classes in the query results.
	AdditionalStorageClasses *[]string `json:"additionalStorageClasses,omitempty"`
	// Human-readable label that identifies the Google Cloud Storage bucket.
	Bucket *string `json:"bucket,omitempty"`
	// Delimiter.
	Delimiter *string `json:"delimiter,omitempty"`
	// Flag that indicates whether to use S3 tags on the files in the given path as additional partition attributes. If set to `true`, data lake adds the S3 tags as additional partition attributes and adds new top-level BSON elements associating each tag to each document.
	IncludeTags *bool `json:"includeTags,omitempty"`
	// Prefix.
	Prefix *string `json:"prefix,omitempty"`
	// Flag that indicates whether the bucket is public. If set to `true`, MongoDB Cloud doesn't use the configured GCP service account to access the bucket. If set to `false`, the configured GCP service acccount must include permissions to access the bucket.
	Public *bool `json:"public,omitempty"`
	// Google Cloud Platform Regions.
	Region *string `json:"region,omitempty"`
	// Human-readable label of the MongoDB Cloud cluster on which the store is based.
	ClusterName *string `json:"clusterName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project.
	// Read only field.
	ProjectId      *string                           `json:"projectId,omitempty"`
	ReadConcern    *DataLakeAtlasStoreReadConcern    `json:"readConcern,omitempty"`
	ReadPreference *DataLakeAtlasStoreReadPreference `json:"readPreference,omitempty"`
	// Flag that validates the scheme in the specified URLs. If `true`, allows insecure `HTTP` scheme, doesn't verify the server's certificate chain and hostname, and accepts any certificate with any hostname presented by the server. If `false`, allows secure `HTTPS` scheme only.
	AllowInsecure *bool `json:"allowInsecure,omitempty"`
	// Default format that Data Lake assumes if it encounters a file without an extension while searching the `storeName`. If omitted, Data Lake attempts to detect the file type by processing a few bytes of the file. The specified format only applies to the URLs specified in the **databases.[n].collections.[n].dataSources** object.
	DefaultFormat *string `json:"defaultFormat,omitempty"`
	// Comma-separated list of publicly accessible HTTP URLs where data is stored. You can't specify URLs that require authentication.
	Urls *[]string `json:"urls,omitempty"`
	// Human-readable label that identifies the name of the container.
	ContainerName *string `json:"containerName,omitempty"`
	// Replacement Delimiter.
	ReplacementDelimiter *string `json:"replacementDelimiter,omitempty"`
	// Service URL.
	ServiceURL *string `json:"serviceURL,omitempty"`
}

// NewDataLakeStoreSettings instantiates a new DataLakeStoreSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeStoreSettings(provider string) *DataLakeStoreSettings {
	this := DataLakeStoreSettings{}
	this.Provider = provider
	var includeTags bool = false
	this.IncludeTags = &includeTags
	var public bool = false
	this.Public = &public
	var allowInsecure bool = false
	this.AllowInsecure = &allowInsecure
	return &this
}

// NewDataLakeStoreSettingsWithDefaults instantiates a new DataLakeStoreSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeStoreSettingsWithDefaults() *DataLakeStoreSettings {
	this := DataLakeStoreSettings{}
	var includeTags bool = false
	this.IncludeTags = &includeTags
	var public bool = false
	this.Public = &public
	var allowInsecure bool = false
	this.AllowInsecure = &allowInsecure
	return &this
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataLakeStoreSettings) SetName(v string) {
	o.Name = &v
}

// GetProvider returns the Provider field value
func (o *DataLakeStoreSettings) GetProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Provider
}

// GetProviderOk returns a tuple with the Provider field value
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Provider, true
}

// SetProvider sets field value
func (o *DataLakeStoreSettings) SetProvider(v string) {
	o.Provider = v
}

// GetAdditionalStorageClasses returns the AdditionalStorageClasses field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetAdditionalStorageClasses() []string {
	if o == nil || IsNil(o.AdditionalStorageClasses) {
		var ret []string
		return ret
	}
	return *o.AdditionalStorageClasses
}

// GetAdditionalStorageClassesOk returns a tuple with the AdditionalStorageClasses field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetAdditionalStorageClassesOk() (*[]string, bool) {
	if o == nil || IsNil(o.AdditionalStorageClasses) {
		return nil, false
	}

	return o.AdditionalStorageClasses, true
}

// HasAdditionalStorageClasses returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasAdditionalStorageClasses() bool {
	if o != nil && !IsNil(o.AdditionalStorageClasses) {
		return true
	}

	return false
}

// SetAdditionalStorageClasses gets a reference to the given []string and assigns it to the AdditionalStorageClasses field.
func (o *DataLakeStoreSettings) SetAdditionalStorageClasses(v []string) {
	o.AdditionalStorageClasses = &v
}

// GetBucket returns the Bucket field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetBucket() string {
	if o == nil || IsNil(o.Bucket) {
		var ret string
		return ret
	}
	return *o.Bucket
}

// GetBucketOk returns a tuple with the Bucket field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetBucketOk() (*string, bool) {
	if o == nil || IsNil(o.Bucket) {
		return nil, false
	}

	return o.Bucket, true
}

// HasBucket returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasBucket() bool {
	if o != nil && !IsNil(o.Bucket) {
		return true
	}

	return false
}

// SetBucket gets a reference to the given string and assigns it to the Bucket field.
func (o *DataLakeStoreSettings) SetBucket(v string) {
	o.Bucket = &v
}

// GetDelimiter returns the Delimiter field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetDelimiter() string {
	if o == nil || IsNil(o.Delimiter) {
		var ret string
		return ret
	}
	return *o.Delimiter
}

// GetDelimiterOk returns a tuple with the Delimiter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetDelimiterOk() (*string, bool) {
	if o == nil || IsNil(o.Delimiter) {
		return nil, false
	}

	return o.Delimiter, true
}

// HasDelimiter returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasDelimiter() bool {
	if o != nil && !IsNil(o.Delimiter) {
		return true
	}

	return false
}

// SetDelimiter gets a reference to the given string and assigns it to the Delimiter field.
func (o *DataLakeStoreSettings) SetDelimiter(v string) {
	o.Delimiter = &v
}

// GetIncludeTags returns the IncludeTags field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetIncludeTags() bool {
	if o == nil || IsNil(o.IncludeTags) {
		var ret bool
		return ret
	}
	return *o.IncludeTags
}

// GetIncludeTagsOk returns a tuple with the IncludeTags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetIncludeTagsOk() (*bool, bool) {
	if o == nil || IsNil(o.IncludeTags) {
		return nil, false
	}

	return o.IncludeTags, true
}

// HasIncludeTags returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasIncludeTags() bool {
	if o != nil && !IsNil(o.IncludeTags) {
		return true
	}

	return false
}

// SetIncludeTags gets a reference to the given bool and assigns it to the IncludeTags field.
func (o *DataLakeStoreSettings) SetIncludeTags(v bool) {
	o.IncludeTags = &v
}

// GetPrefix returns the Prefix field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetPrefix() string {
	if o == nil || IsNil(o.Prefix) {
		var ret string
		return ret
	}
	return *o.Prefix
}

// GetPrefixOk returns a tuple with the Prefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetPrefixOk() (*string, bool) {
	if o == nil || IsNil(o.Prefix) {
		return nil, false
	}

	return o.Prefix, true
}

// HasPrefix returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasPrefix() bool {
	if o != nil && !IsNil(o.Prefix) {
		return true
	}

	return false
}

// SetPrefix gets a reference to the given string and assigns it to the Prefix field.
func (o *DataLakeStoreSettings) SetPrefix(v string) {
	o.Prefix = &v
}

// GetPublic returns the Public field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetPublic() bool {
	if o == nil || IsNil(o.Public) {
		var ret bool
		return ret
	}
	return *o.Public
}

// GetPublicOk returns a tuple with the Public field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetPublicOk() (*bool, bool) {
	if o == nil || IsNil(o.Public) {
		return nil, false
	}

	return o.Public, true
}

// HasPublic returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasPublic() bool {
	if o != nil && !IsNil(o.Public) {
		return true
	}

	return false
}

// SetPublic gets a reference to the given bool and assigns it to the Public field.
func (o *DataLakeStoreSettings) SetPublic(v bool) {
	o.Public = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *DataLakeStoreSettings) SetRegion(v string) {
	o.Region = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *DataLakeStoreSettings) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetProjectId returns the ProjectId field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetProjectId() string {
	if o == nil || IsNil(o.ProjectId) {
		var ret string
		return ret
	}
	return *o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProjectId) {
		return nil, false
	}

	return o.ProjectId, true
}

// HasProjectId returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasProjectId() bool {
	if o != nil && !IsNil(o.ProjectId) {
		return true
	}

	return false
}

// SetProjectId gets a reference to the given string and assigns it to the ProjectId field.
func (o *DataLakeStoreSettings) SetProjectId(v string) {
	o.ProjectId = &v
}

// GetReadConcern returns the ReadConcern field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetReadConcern() DataLakeAtlasStoreReadConcern {
	if o == nil || IsNil(o.ReadConcern) {
		var ret DataLakeAtlasStoreReadConcern
		return ret
	}
	return *o.ReadConcern
}

// GetReadConcernOk returns a tuple with the ReadConcern field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetReadConcernOk() (*DataLakeAtlasStoreReadConcern, bool) {
	if o == nil || IsNil(o.ReadConcern) {
		return nil, false
	}

	return o.ReadConcern, true
}

// HasReadConcern returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasReadConcern() bool {
	if o != nil && !IsNil(o.ReadConcern) {
		return true
	}

	return false
}

// SetReadConcern gets a reference to the given DataLakeAtlasStoreReadConcern and assigns it to the ReadConcern field.
func (o *DataLakeStoreSettings) SetReadConcern(v DataLakeAtlasStoreReadConcern) {
	o.ReadConcern = &v
}

// GetReadPreference returns the ReadPreference field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetReadPreference() DataLakeAtlasStoreReadPreference {
	if o == nil || IsNil(o.ReadPreference) {
		var ret DataLakeAtlasStoreReadPreference
		return ret
	}
	return *o.ReadPreference
}

// GetReadPreferenceOk returns a tuple with the ReadPreference field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetReadPreferenceOk() (*DataLakeAtlasStoreReadPreference, bool) {
	if o == nil || IsNil(o.ReadPreference) {
		return nil, false
	}

	return o.ReadPreference, true
}

// HasReadPreference returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasReadPreference() bool {
	if o != nil && !IsNil(o.ReadPreference) {
		return true
	}

	return false
}

// SetReadPreference gets a reference to the given DataLakeAtlasStoreReadPreference and assigns it to the ReadPreference field.
func (o *DataLakeStoreSettings) SetReadPreference(v DataLakeAtlasStoreReadPreference) {
	o.ReadPreference = &v
}

// GetAllowInsecure returns the AllowInsecure field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetAllowInsecure() bool {
	if o == nil || IsNil(o.AllowInsecure) {
		var ret bool
		return ret
	}
	return *o.AllowInsecure
}

// GetAllowInsecureOk returns a tuple with the AllowInsecure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetAllowInsecureOk() (*bool, bool) {
	if o == nil || IsNil(o.AllowInsecure) {
		return nil, false
	}

	return o.AllowInsecure, true
}

// HasAllowInsecure returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasAllowInsecure() bool {
	if o != nil && !IsNil(o.AllowInsecure) {
		return true
	}

	return false
}

// SetAllowInsecure gets a reference to the given bool and assigns it to the AllowInsecure field.
func (o *DataLakeStoreSettings) SetAllowInsecure(v bool) {
	o.AllowInsecure = &v
}

// GetDefaultFormat returns the DefaultFormat field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetDefaultFormat() string {
	if o == nil || IsNil(o.DefaultFormat) {
		var ret string
		return ret
	}
	return *o.DefaultFormat
}

// GetDefaultFormatOk returns a tuple with the DefaultFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetDefaultFormatOk() (*string, bool) {
	if o == nil || IsNil(o.DefaultFormat) {
		return nil, false
	}

	return o.DefaultFormat, true
}

// HasDefaultFormat returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasDefaultFormat() bool {
	if o != nil && !IsNil(o.DefaultFormat) {
		return true
	}

	return false
}

// SetDefaultFormat gets a reference to the given string and assigns it to the DefaultFormat field.
func (o *DataLakeStoreSettings) SetDefaultFormat(v string) {
	o.DefaultFormat = &v
}

// GetUrls returns the Urls field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetUrls() []string {
	if o == nil || IsNil(o.Urls) {
		var ret []string
		return ret
	}
	return *o.Urls
}

// GetUrlsOk returns a tuple with the Urls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetUrlsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Urls) {
		return nil, false
	}

	return o.Urls, true
}

// HasUrls returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasUrls() bool {
	if o != nil && !IsNil(o.Urls) {
		return true
	}

	return false
}

// SetUrls gets a reference to the given []string and assigns it to the Urls field.
func (o *DataLakeStoreSettings) SetUrls(v []string) {
	o.Urls = &v
}

// GetContainerName returns the ContainerName field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetContainerName() string {
	if o == nil || IsNil(o.ContainerName) {
		var ret string
		return ret
	}
	return *o.ContainerName
}

// GetContainerNameOk returns a tuple with the ContainerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetContainerNameOk() (*string, bool) {
	if o == nil || IsNil(o.ContainerName) {
		return nil, false
	}

	return o.ContainerName, true
}

// HasContainerName returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasContainerName() bool {
	if o != nil && !IsNil(o.ContainerName) {
		return true
	}

	return false
}

// SetContainerName gets a reference to the given string and assigns it to the ContainerName field.
func (o *DataLakeStoreSettings) SetContainerName(v string) {
	o.ContainerName = &v
}

// GetReplacementDelimiter returns the ReplacementDelimiter field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetReplacementDelimiter() string {
	if o == nil || IsNil(o.ReplacementDelimiter) {
		var ret string
		return ret
	}
	return *o.ReplacementDelimiter
}

// GetReplacementDelimiterOk returns a tuple with the ReplacementDelimiter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetReplacementDelimiterOk() (*string, bool) {
	if o == nil || IsNil(o.ReplacementDelimiter) {
		return nil, false
	}

	return o.ReplacementDelimiter, true
}

// HasReplacementDelimiter returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasReplacementDelimiter() bool {
	if o != nil && !IsNil(o.ReplacementDelimiter) {
		return true
	}

	return false
}

// SetReplacementDelimiter gets a reference to the given string and assigns it to the ReplacementDelimiter field.
func (o *DataLakeStoreSettings) SetReplacementDelimiter(v string) {
	o.ReplacementDelimiter = &v
}

// GetServiceURL returns the ServiceURL field value if set, zero value otherwise
func (o *DataLakeStoreSettings) GetServiceURL() string {
	if o == nil || IsNil(o.ServiceURL) {
		var ret string
		return ret
	}
	return *o.ServiceURL
}

// GetServiceURLOk returns a tuple with the ServiceURL field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeStoreSettings) GetServiceURLOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceURL) {
		return nil, false
	}

	return o.ServiceURL, true
}

// HasServiceURL returns a boolean if a field has been set.
func (o *DataLakeStoreSettings) HasServiceURL() bool {
	if o != nil && !IsNil(o.ServiceURL) {
		return true
	}

	return false
}

// SetServiceURL gets a reference to the given string and assigns it to the ServiceURL field.
func (o *DataLakeStoreSettings) SetServiceURL(v string) {
	o.ServiceURL = &v
}
