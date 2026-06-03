// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterDescriptionProcessArgs20240805 struct for ClusterDescriptionProcessArgs20240805
type ClusterDescriptionProcessArgs20240805 struct {
	// The minimum pre- and post-image retention time in seconds.
	ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds *int `json:"changeStreamOptionsPreAndPostImagesExpireAfterSeconds,omitempty"`
	// Number of threads on the source shard and the receiving shard for chunk migration. The number of threads should not exceed the half the total number of CPU cores in the sharded cluster.
	ChunkMigrationConcurrency *int `json:"chunkMigrationConcurrency,omitempty"`
	// The custom OpenSSL cipher suite list for TLS 1.2. This field is only valid when `tlsCipherConfigMode` is set to `CUSTOM`.
	CustomOpensslCipherConfigTls12 *[]string `json:"customOpensslCipherConfigTls12,omitempty"`
	// The custom OpenSSL cipher suite list for TLS 1.3. This field is only valid when `tlsCipherConfigMode` is set to `CUSTOM`.
	CustomOpensslCipherConfigTls13 *[]string `json:"customOpensslCipherConfigTls13,omitempty"`
	// Default time limit in milliseconds for individual read operations to complete.
	DefaultMaxTimeMS *int `json:"defaultMaxTimeMS,omitempty"`
	// Default level of acknowledgment requested from MongoDB for write operations when none is specified by the driver.
	DefaultWriteConcern *string `json:"defaultWriteConcern,omitempty"`
	// Flag that indicates whether the cluster allows execution of operations that perform server-side executions of JavaScript. When using 8.0+, we recommend disabling server-side JavaScript and using operators of aggregation pipeline as more performant alternative.
	JavascriptEnabled *bool `json:"javascriptEnabled,omitempty"`
	// Minimum Transport Layer Security (TLS) version that the cluster accepts for incoming connections. Clusters using TLS 1.0 or 1.1 should consider setting TLS 1.2 as the minimum TLS protocol version.
	MinimumEnabledTlsProtocol *string `json:"minimumEnabledTlsProtocol,omitempty"`
	// Flag that indicates whether the cluster disables executing any query that requires a collection scan to return results.
	NoTableScan *bool `json:"noTableScan,omitempty"`
	// Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.
	OplogMinRetentionHours *float64 `json:"oplogMinRetentionHours,omitempty"`
	// Storage limit of cluster's oplog expressed in megabytes. A value of null indicates that the cluster uses the default oplog size that MongoDB Cloud calculates.
	OplogSizeMB *int `json:"oplogSizeMB,omitempty"`
	// May be set to 1 (disabled) or 3 (enabled). When set to 3, Atlas will include redacted and anonymized `$queryStats` output in MongoDB logs. `$queryStats` output does not contain literals or field values. Enabling this setting might impact the performance of your cluster.
	QueryStatsLogVerbosity *int `json:"queryStatsLogVerbosity,omitempty"`
	// Interval in seconds at which the mongosqld process re-samples data to create its relational schema.
	SampleRefreshIntervalBIConnector *int `json:"sampleRefreshIntervalBIConnector,omitempty"`
	// Number of documents per database to sample when gathering schema information.
	SampleSizeBIConnector *int `json:"sampleSizeBIConnector,omitempty"`
	// The TLS cipher suite configuration mode. The default mode uses the default cipher suites. The custom mode allows you to specify custom cipher suites for both TLS 1.2 and TLS 1.3.
	TlsCipherConfigMode *string `json:"tlsCipherConfigMode,omitempty"`
	// Lifetime, in seconds, of multi-document transactions. Atlas considers the transactions that exceed this limit as expired and so aborts them through a periodic clean-up process.
	TransactionLifetimeLimitSeconds *int64 `json:"transactionLifetimeLimitSeconds,omitempty"`
}

// NewClusterDescriptionProcessArgs20240805 instantiates a new ClusterDescriptionProcessArgs20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterDescriptionProcessArgs20240805() *ClusterDescriptionProcessArgs20240805 {
	this := ClusterDescriptionProcessArgs20240805{}
	var changeStreamOptionsPreAndPostImagesExpireAfterSeconds int = -1
	this.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds = &changeStreamOptionsPreAndPostImagesExpireAfterSeconds
	var sampleRefreshIntervalBIConnector int = 0
	this.SampleRefreshIntervalBIConnector = &sampleRefreshIntervalBIConnector
	return &this
}

// NewClusterDescriptionProcessArgs20240805WithDefaults instantiates a new ClusterDescriptionProcessArgs20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterDescriptionProcessArgs20240805WithDefaults() *ClusterDescriptionProcessArgs20240805 {
	this := ClusterDescriptionProcessArgs20240805{}
	var changeStreamOptionsPreAndPostImagesExpireAfterSeconds int = -1
	this.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds = &changeStreamOptionsPreAndPostImagesExpireAfterSeconds
	var sampleRefreshIntervalBIConnector int = 0
	this.SampleRefreshIntervalBIConnector = &sampleRefreshIntervalBIConnector
	return &this
}

// GetChangeStreamOptionsPreAndPostImagesExpireAfterSeconds returns the ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetChangeStreamOptionsPreAndPostImagesExpireAfterSeconds() int {
	if o == nil || IsNil(o.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds) {
		var ret int
		return ret
	}
	return *o.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds
}

// GetChangeStreamOptionsPreAndPostImagesExpireAfterSecondsOk returns a tuple with the ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetChangeStreamOptionsPreAndPostImagesExpireAfterSecondsOk() (*int, bool) {
	if o == nil || IsNil(o.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds) {
		return nil, false
	}

	return o.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds, true
}

// HasChangeStreamOptionsPreAndPostImagesExpireAfterSeconds returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasChangeStreamOptionsPreAndPostImagesExpireAfterSeconds() bool {
	if o != nil && !IsNil(o.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds) {
		return true
	}

	return false
}

// SetChangeStreamOptionsPreAndPostImagesExpireAfterSeconds gets a reference to the given int and assigns it to the ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds field.
func (o *ClusterDescriptionProcessArgs20240805) SetChangeStreamOptionsPreAndPostImagesExpireAfterSeconds(v int) {
	o.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds = &v
}

// GetChunkMigrationConcurrency returns the ChunkMigrationConcurrency field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetChunkMigrationConcurrency() int {
	if o == nil || IsNil(o.ChunkMigrationConcurrency) {
		var ret int
		return ret
	}
	return *o.ChunkMigrationConcurrency
}

// GetChunkMigrationConcurrencyOk returns a tuple with the ChunkMigrationConcurrency field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetChunkMigrationConcurrencyOk() (*int, bool) {
	if o == nil || IsNil(o.ChunkMigrationConcurrency) {
		return nil, false
	}

	return o.ChunkMigrationConcurrency, true
}

// HasChunkMigrationConcurrency returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasChunkMigrationConcurrency() bool {
	if o != nil && !IsNil(o.ChunkMigrationConcurrency) {
		return true
	}

	return false
}

// SetChunkMigrationConcurrency gets a reference to the given int and assigns it to the ChunkMigrationConcurrency field.
func (o *ClusterDescriptionProcessArgs20240805) SetChunkMigrationConcurrency(v int) {
	o.ChunkMigrationConcurrency = &v
}

// GetCustomOpensslCipherConfigTls12 returns the CustomOpensslCipherConfigTls12 field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetCustomOpensslCipherConfigTls12() []string {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls12) {
		var ret []string
		return ret
	}
	return *o.CustomOpensslCipherConfigTls12
}

// GetCustomOpensslCipherConfigTls12Ok returns a tuple with the CustomOpensslCipherConfigTls12 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetCustomOpensslCipherConfigTls12Ok() (*[]string, bool) {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls12) {
		return nil, false
	}

	return o.CustomOpensslCipherConfigTls12, true
}

// HasCustomOpensslCipherConfigTls12 returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasCustomOpensslCipherConfigTls12() bool {
	if o != nil && !IsNil(o.CustomOpensslCipherConfigTls12) {
		return true
	}

	return false
}

// SetCustomOpensslCipherConfigTls12 gets a reference to the given []string and assigns it to the CustomOpensslCipherConfigTls12 field.
func (o *ClusterDescriptionProcessArgs20240805) SetCustomOpensslCipherConfigTls12(v []string) {
	o.CustomOpensslCipherConfigTls12 = &v
}

// GetCustomOpensslCipherConfigTls13 returns the CustomOpensslCipherConfigTls13 field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetCustomOpensslCipherConfigTls13() []string {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls13) {
		var ret []string
		return ret
	}
	return *o.CustomOpensslCipherConfigTls13
}

// GetCustomOpensslCipherConfigTls13Ok returns a tuple with the CustomOpensslCipherConfigTls13 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetCustomOpensslCipherConfigTls13Ok() (*[]string, bool) {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls13) {
		return nil, false
	}

	return o.CustomOpensslCipherConfigTls13, true
}

// HasCustomOpensslCipherConfigTls13 returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasCustomOpensslCipherConfigTls13() bool {
	if o != nil && !IsNil(o.CustomOpensslCipherConfigTls13) {
		return true
	}

	return false
}

// SetCustomOpensslCipherConfigTls13 gets a reference to the given []string and assigns it to the CustomOpensslCipherConfigTls13 field.
func (o *ClusterDescriptionProcessArgs20240805) SetCustomOpensslCipherConfigTls13(v []string) {
	o.CustomOpensslCipherConfigTls13 = &v
}

// GetDefaultMaxTimeMS returns the DefaultMaxTimeMS field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetDefaultMaxTimeMS() int {
	if o == nil || IsNil(o.DefaultMaxTimeMS) {
		var ret int
		return ret
	}
	return *o.DefaultMaxTimeMS
}

// GetDefaultMaxTimeMSOk returns a tuple with the DefaultMaxTimeMS field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetDefaultMaxTimeMSOk() (*int, bool) {
	if o == nil || IsNil(o.DefaultMaxTimeMS) {
		return nil, false
	}

	return o.DefaultMaxTimeMS, true
}

// HasDefaultMaxTimeMS returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasDefaultMaxTimeMS() bool {
	if o != nil && !IsNil(o.DefaultMaxTimeMS) {
		return true
	}

	return false
}

// SetDefaultMaxTimeMS gets a reference to the given int and assigns it to the DefaultMaxTimeMS field.
func (o *ClusterDescriptionProcessArgs20240805) SetDefaultMaxTimeMS(v int) {
	o.DefaultMaxTimeMS = &v
}

// GetDefaultWriteConcern returns the DefaultWriteConcern field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetDefaultWriteConcern() string {
	if o == nil || IsNil(o.DefaultWriteConcern) {
		var ret string
		return ret
	}
	return *o.DefaultWriteConcern
}

// GetDefaultWriteConcernOk returns a tuple with the DefaultWriteConcern field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetDefaultWriteConcernOk() (*string, bool) {
	if o == nil || IsNil(o.DefaultWriteConcern) {
		return nil, false
	}

	return o.DefaultWriteConcern, true
}

// HasDefaultWriteConcern returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasDefaultWriteConcern() bool {
	if o != nil && !IsNil(o.DefaultWriteConcern) {
		return true
	}

	return false
}

// SetDefaultWriteConcern gets a reference to the given string and assigns it to the DefaultWriteConcern field.
func (o *ClusterDescriptionProcessArgs20240805) SetDefaultWriteConcern(v string) {
	o.DefaultWriteConcern = &v
}

// GetJavascriptEnabled returns the JavascriptEnabled field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetJavascriptEnabled() bool {
	if o == nil || IsNil(o.JavascriptEnabled) {
		var ret bool
		return ret
	}
	return *o.JavascriptEnabled
}

// GetJavascriptEnabledOk returns a tuple with the JavascriptEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetJavascriptEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.JavascriptEnabled) {
		return nil, false
	}

	return o.JavascriptEnabled, true
}

// HasJavascriptEnabled returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasJavascriptEnabled() bool {
	if o != nil && !IsNil(o.JavascriptEnabled) {
		return true
	}

	return false
}

// SetJavascriptEnabled gets a reference to the given bool and assigns it to the JavascriptEnabled field.
func (o *ClusterDescriptionProcessArgs20240805) SetJavascriptEnabled(v bool) {
	o.JavascriptEnabled = &v
}

// GetMinimumEnabledTlsProtocol returns the MinimumEnabledTlsProtocol field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetMinimumEnabledTlsProtocol() string {
	if o == nil || IsNil(o.MinimumEnabledTlsProtocol) {
		var ret string
		return ret
	}
	return *o.MinimumEnabledTlsProtocol
}

// GetMinimumEnabledTlsProtocolOk returns a tuple with the MinimumEnabledTlsProtocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetMinimumEnabledTlsProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.MinimumEnabledTlsProtocol) {
		return nil, false
	}

	return o.MinimumEnabledTlsProtocol, true
}

// HasMinimumEnabledTlsProtocol returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasMinimumEnabledTlsProtocol() bool {
	if o != nil && !IsNil(o.MinimumEnabledTlsProtocol) {
		return true
	}

	return false
}

// SetMinimumEnabledTlsProtocol gets a reference to the given string and assigns it to the MinimumEnabledTlsProtocol field.
func (o *ClusterDescriptionProcessArgs20240805) SetMinimumEnabledTlsProtocol(v string) {
	o.MinimumEnabledTlsProtocol = &v
}

// GetNoTableScan returns the NoTableScan field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetNoTableScan() bool {
	if o == nil || IsNil(o.NoTableScan) {
		var ret bool
		return ret
	}
	return *o.NoTableScan
}

// GetNoTableScanOk returns a tuple with the NoTableScan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetNoTableScanOk() (*bool, bool) {
	if o == nil || IsNil(o.NoTableScan) {
		return nil, false
	}

	return o.NoTableScan, true
}

// HasNoTableScan returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasNoTableScan() bool {
	if o != nil && !IsNil(o.NoTableScan) {
		return true
	}

	return false
}

// SetNoTableScan gets a reference to the given bool and assigns it to the NoTableScan field.
func (o *ClusterDescriptionProcessArgs20240805) SetNoTableScan(v bool) {
	o.NoTableScan = &v
}

// GetOplogMinRetentionHours returns the OplogMinRetentionHours field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetOplogMinRetentionHours() float64 {
	if o == nil || IsNil(o.OplogMinRetentionHours) {
		var ret float64
		return ret
	}
	return *o.OplogMinRetentionHours
}

// GetOplogMinRetentionHoursOk returns a tuple with the OplogMinRetentionHours field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetOplogMinRetentionHoursOk() (*float64, bool) {
	if o == nil || IsNil(o.OplogMinRetentionHours) {
		return nil, false
	}

	return o.OplogMinRetentionHours, true
}

// HasOplogMinRetentionHours returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasOplogMinRetentionHours() bool {
	if o != nil && !IsNil(o.OplogMinRetentionHours) {
		return true
	}

	return false
}

// SetOplogMinRetentionHours gets a reference to the given float64 and assigns it to the OplogMinRetentionHours field.
func (o *ClusterDescriptionProcessArgs20240805) SetOplogMinRetentionHours(v float64) {
	o.OplogMinRetentionHours = &v
}

// GetOplogSizeMB returns the OplogSizeMB field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetOplogSizeMB() int {
	if o == nil || IsNil(o.OplogSizeMB) {
		var ret int
		return ret
	}
	return *o.OplogSizeMB
}

// GetOplogSizeMBOk returns a tuple with the OplogSizeMB field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetOplogSizeMBOk() (*int, bool) {
	if o == nil || IsNil(o.OplogSizeMB) {
		return nil, false
	}

	return o.OplogSizeMB, true
}

// HasOplogSizeMB returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasOplogSizeMB() bool {
	if o != nil && !IsNil(o.OplogSizeMB) {
		return true
	}

	return false
}

// SetOplogSizeMB gets a reference to the given int and assigns it to the OplogSizeMB field.
func (o *ClusterDescriptionProcessArgs20240805) SetOplogSizeMB(v int) {
	o.OplogSizeMB = &v
}

// GetQueryStatsLogVerbosity returns the QueryStatsLogVerbosity field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetQueryStatsLogVerbosity() int {
	if o == nil || IsNil(o.QueryStatsLogVerbosity) {
		var ret int
		return ret
	}
	return *o.QueryStatsLogVerbosity
}

// GetQueryStatsLogVerbosityOk returns a tuple with the QueryStatsLogVerbosity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetQueryStatsLogVerbosityOk() (*int, bool) {
	if o == nil || IsNil(o.QueryStatsLogVerbosity) {
		return nil, false
	}

	return o.QueryStatsLogVerbosity, true
}

// HasQueryStatsLogVerbosity returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasQueryStatsLogVerbosity() bool {
	if o != nil && !IsNil(o.QueryStatsLogVerbosity) {
		return true
	}

	return false
}

// SetQueryStatsLogVerbosity gets a reference to the given int and assigns it to the QueryStatsLogVerbosity field.
func (o *ClusterDescriptionProcessArgs20240805) SetQueryStatsLogVerbosity(v int) {
	o.QueryStatsLogVerbosity = &v
}

// GetSampleRefreshIntervalBIConnector returns the SampleRefreshIntervalBIConnector field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetSampleRefreshIntervalBIConnector() int {
	if o == nil || IsNil(o.SampleRefreshIntervalBIConnector) {
		var ret int
		return ret
	}
	return *o.SampleRefreshIntervalBIConnector
}

// GetSampleRefreshIntervalBIConnectorOk returns a tuple with the SampleRefreshIntervalBIConnector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetSampleRefreshIntervalBIConnectorOk() (*int, bool) {
	if o == nil || IsNil(o.SampleRefreshIntervalBIConnector) {
		return nil, false
	}

	return o.SampleRefreshIntervalBIConnector, true
}

// HasSampleRefreshIntervalBIConnector returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasSampleRefreshIntervalBIConnector() bool {
	if o != nil && !IsNil(o.SampleRefreshIntervalBIConnector) {
		return true
	}

	return false
}

// SetSampleRefreshIntervalBIConnector gets a reference to the given int and assigns it to the SampleRefreshIntervalBIConnector field.
func (o *ClusterDescriptionProcessArgs20240805) SetSampleRefreshIntervalBIConnector(v int) {
	o.SampleRefreshIntervalBIConnector = &v
}

// GetSampleSizeBIConnector returns the SampleSizeBIConnector field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetSampleSizeBIConnector() int {
	if o == nil || IsNil(o.SampleSizeBIConnector) {
		var ret int
		return ret
	}
	return *o.SampleSizeBIConnector
}

// GetSampleSizeBIConnectorOk returns a tuple with the SampleSizeBIConnector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetSampleSizeBIConnectorOk() (*int, bool) {
	if o == nil || IsNil(o.SampleSizeBIConnector) {
		return nil, false
	}

	return o.SampleSizeBIConnector, true
}

// HasSampleSizeBIConnector returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasSampleSizeBIConnector() bool {
	if o != nil && !IsNil(o.SampleSizeBIConnector) {
		return true
	}

	return false
}

// SetSampleSizeBIConnector gets a reference to the given int and assigns it to the SampleSizeBIConnector field.
func (o *ClusterDescriptionProcessArgs20240805) SetSampleSizeBIConnector(v int) {
	o.SampleSizeBIConnector = &v
}

// GetTlsCipherConfigMode returns the TlsCipherConfigMode field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetTlsCipherConfigMode() string {
	if o == nil || IsNil(o.TlsCipherConfigMode) {
		var ret string
		return ret
	}
	return *o.TlsCipherConfigMode
}

// GetTlsCipherConfigModeOk returns a tuple with the TlsCipherConfigMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetTlsCipherConfigModeOk() (*string, bool) {
	if o == nil || IsNil(o.TlsCipherConfigMode) {
		return nil, false
	}

	return o.TlsCipherConfigMode, true
}

// HasTlsCipherConfigMode returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasTlsCipherConfigMode() bool {
	if o != nil && !IsNil(o.TlsCipherConfigMode) {
		return true
	}

	return false
}

// SetTlsCipherConfigMode gets a reference to the given string and assigns it to the TlsCipherConfigMode field.
func (o *ClusterDescriptionProcessArgs20240805) SetTlsCipherConfigMode(v string) {
	o.TlsCipherConfigMode = &v
}

// GetTransactionLifetimeLimitSeconds returns the TransactionLifetimeLimitSeconds field value if set, zero value otherwise
func (o *ClusterDescriptionProcessArgs20240805) GetTransactionLifetimeLimitSeconds() int64 {
	if o == nil || IsNil(o.TransactionLifetimeLimitSeconds) {
		var ret int64
		return ret
	}
	return *o.TransactionLifetimeLimitSeconds
}

// GetTransactionLifetimeLimitSecondsOk returns a tuple with the TransactionLifetimeLimitSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionProcessArgs20240805) GetTransactionLifetimeLimitSecondsOk() (*int64, bool) {
	if o == nil || IsNil(o.TransactionLifetimeLimitSeconds) {
		return nil, false
	}

	return o.TransactionLifetimeLimitSeconds, true
}

// HasTransactionLifetimeLimitSeconds returns a boolean if a field has been set.
func (o *ClusterDescriptionProcessArgs20240805) HasTransactionLifetimeLimitSeconds() bool {
	if o != nil && !IsNil(o.TransactionLifetimeLimitSeconds) {
		return true
	}

	return false
}

// SetTransactionLifetimeLimitSeconds gets a reference to the given int64 and assigns it to the TransactionLifetimeLimitSeconds field.
func (o *ClusterDescriptionProcessArgs20240805) SetTransactionLifetimeLimitSeconds(v int64) {
	o.TransactionLifetimeLimitSeconds = &v
}
