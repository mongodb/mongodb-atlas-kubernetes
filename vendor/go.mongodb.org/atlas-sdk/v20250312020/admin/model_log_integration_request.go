// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LogIntegrationRequest Request schema for creating or updating a log integration.
type LogIntegrationRequest struct {
	// Array of log types exported by this integration.
	LogTypes []string `json:"logTypes"`
	// Human-readable label that identifies the service to which you want to integrate with Atlas. The value must match the log integration type. This value cannot be modified after the integration is created.
	Type string `json:"type"`
	// Name of the bucket to store log files.
	BucketName *string `json:"bucketName,omitempty"`
	// Unique 24-character hexadecimal string that identifies the AWS IAM role that Atlas uses to access the S3 bucket.
	IamRoleId *string `json:"iamRoleId,omitempty"`
	// AWS KMS key ID or ARN for server-side encryption (optional). If not provided, uses bucket default encryption settings.
	KmsKey *string `json:"kmsKey,omitempty"`
	// Path prefix where the log files will be stored. Atlas will add further sub-directories based on the log type.
	PrefixPath *string `json:"prefixPath,omitempty"`
	// When true, uses the legacy daily-folder path structure compatible with Push-Based Log Export: `{prefix}/{cluster}/{hostname}/{logType}/{YYYY-MM-DD}/{timestamp}-{logType}.log`. When false (default), uses the flat timestamped structure: `{prefix}/{cluster}/{hostname}/{logType}/{timestamp}-{logType}.log`.
	UseLegacyPathStructure *bool `json:"useLegacyPathStructure,omitempty"`
	// API key for authentication.
	ApiKey *string `json:"apiKey,omitempty"`
	// Datadog site/region for log ingestion. Valid values: US1, US3, US5, EU, AP1, AP2, US1_FED.
	Region *string `json:"region,omitempty"`
	// Unique 24-character hexadecimal string that identifies the Atlas Cloud Provider Access role.
	RoleId *string `json:"roleId,omitempty"`
	// OpenTelemetry collector endpoint URL. Must be HTTPS and not exceed 2048 characters.
	OtelEndpoint *string `json:"otelEndpoint,omitempty"`
	// HTTP headers for authentication and configuration. Maximum 10 headers, total size limit 2KB.
	OtelSuppliedHeaders *[]Header `json:"otelSuppliedHeaders,omitempty"`
	// HTTP Event Collector (HEC) token for authentication.
	HecToken *string `json:"hecToken,omitempty"`
	// HTTP Event Collector (HEC) endpoint URL.
	HecUrl *string `json:"hecUrl,omitempty"`
	// Storage account name where logs will be stored.
	StorageAccountName *string `json:"storageAccountName,omitempty"`
	// Storage container name for log files.
	StorageContainerName *string `json:"storageContainerName,omitempty"`
}

// NewLogIntegrationRequest instantiates a new LogIntegrationRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLogIntegrationRequest(logTypes []string, type_ string) *LogIntegrationRequest {
	this := LogIntegrationRequest{}
	this.LogTypes = logTypes
	this.Type = type_
	return &this
}

// NewLogIntegrationRequestWithDefaults instantiates a new LogIntegrationRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLogIntegrationRequestWithDefaults() *LogIntegrationRequest {
	this := LogIntegrationRequest{}
	return &this
}

// GetLogTypes returns the LogTypes field value
func (o *LogIntegrationRequest) GetLogTypes() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.LogTypes
}

// GetLogTypesOk returns a tuple with the LogTypes field value
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetLogTypesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LogTypes, true
}

// SetLogTypes sets field value
func (o *LogIntegrationRequest) SetLogTypes(v []string) {
	o.LogTypes = v
}

// GetType returns the Type field value
func (o *LogIntegrationRequest) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *LogIntegrationRequest) SetType(v string) {
	o.Type = v
}

// GetBucketName returns the BucketName field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetBucketName() string {
	if o == nil || IsNil(o.BucketName) {
		var ret string
		return ret
	}
	return *o.BucketName
}

// GetBucketNameOk returns a tuple with the BucketName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetBucketNameOk() (*string, bool) {
	if o == nil || IsNil(o.BucketName) {
		return nil, false
	}

	return o.BucketName, true
}

// HasBucketName returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasBucketName() bool {
	if o != nil && !IsNil(o.BucketName) {
		return true
	}

	return false
}

// SetBucketName gets a reference to the given string and assigns it to the BucketName field.
func (o *LogIntegrationRequest) SetBucketName(v string) {
	o.BucketName = &v
}

// GetIamRoleId returns the IamRoleId field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetIamRoleId() string {
	if o == nil || IsNil(o.IamRoleId) {
		var ret string
		return ret
	}
	return *o.IamRoleId
}

// GetIamRoleIdOk returns a tuple with the IamRoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetIamRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.IamRoleId) {
		return nil, false
	}

	return o.IamRoleId, true
}

// HasIamRoleId returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasIamRoleId() bool {
	if o != nil && !IsNil(o.IamRoleId) {
		return true
	}

	return false
}

// SetIamRoleId gets a reference to the given string and assigns it to the IamRoleId field.
func (o *LogIntegrationRequest) SetIamRoleId(v string) {
	o.IamRoleId = &v
}

// GetKmsKey returns the KmsKey field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetKmsKey() string {
	if o == nil || IsNil(o.KmsKey) {
		var ret string
		return ret
	}
	return *o.KmsKey
}

// GetKmsKeyOk returns a tuple with the KmsKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetKmsKeyOk() (*string, bool) {
	if o == nil || IsNil(o.KmsKey) {
		return nil, false
	}

	return o.KmsKey, true
}

// HasKmsKey returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasKmsKey() bool {
	if o != nil && !IsNil(o.KmsKey) {
		return true
	}

	return false
}

// SetKmsKey gets a reference to the given string and assigns it to the KmsKey field.
func (o *LogIntegrationRequest) SetKmsKey(v string) {
	o.KmsKey = &v
}

// GetPrefixPath returns the PrefixPath field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetPrefixPath() string {
	if o == nil || IsNil(o.PrefixPath) {
		var ret string
		return ret
	}
	return *o.PrefixPath
}

// GetPrefixPathOk returns a tuple with the PrefixPath field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetPrefixPathOk() (*string, bool) {
	if o == nil || IsNil(o.PrefixPath) {
		return nil, false
	}

	return o.PrefixPath, true
}

// HasPrefixPath returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasPrefixPath() bool {
	if o != nil && !IsNil(o.PrefixPath) {
		return true
	}

	return false
}

// SetPrefixPath gets a reference to the given string and assigns it to the PrefixPath field.
func (o *LogIntegrationRequest) SetPrefixPath(v string) {
	o.PrefixPath = &v
}

// GetUseLegacyPathStructure returns the UseLegacyPathStructure field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetUseLegacyPathStructure() bool {
	if o == nil || IsNil(o.UseLegacyPathStructure) {
		var ret bool
		return ret
	}
	return *o.UseLegacyPathStructure
}

// GetUseLegacyPathStructureOk returns a tuple with the UseLegacyPathStructure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetUseLegacyPathStructureOk() (*bool, bool) {
	if o == nil || IsNil(o.UseLegacyPathStructure) {
		return nil, false
	}

	return o.UseLegacyPathStructure, true
}

// HasUseLegacyPathStructure returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasUseLegacyPathStructure() bool {
	if o != nil && !IsNil(o.UseLegacyPathStructure) {
		return true
	}

	return false
}

// SetUseLegacyPathStructure gets a reference to the given bool and assigns it to the UseLegacyPathStructure field.
func (o *LogIntegrationRequest) SetUseLegacyPathStructure(v bool) {
	o.UseLegacyPathStructure = &v
}

// GetApiKey returns the ApiKey field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetApiKey() string {
	if o == nil || IsNil(o.ApiKey) {
		var ret string
		return ret
	}
	return *o.ApiKey
}

// GetApiKeyOk returns a tuple with the ApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetApiKeyOk() (*string, bool) {
	if o == nil || IsNil(o.ApiKey) {
		return nil, false
	}

	return o.ApiKey, true
}

// HasApiKey returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasApiKey() bool {
	if o != nil && !IsNil(o.ApiKey) {
		return true
	}

	return false
}

// SetApiKey gets a reference to the given string and assigns it to the ApiKey field.
func (o *LogIntegrationRequest) SetApiKey(v string) {
	o.ApiKey = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *LogIntegrationRequest) SetRegion(v string) {
	o.Region = &v
}

// GetRoleId returns the RoleId field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetRoleId() string {
	if o == nil || IsNil(o.RoleId) {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetRoleIdOk() (*string, bool) {
	if o == nil || IsNil(o.RoleId) {
		return nil, false
	}

	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasRoleId() bool {
	if o != nil && !IsNil(o.RoleId) {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *LogIntegrationRequest) SetRoleId(v string) {
	o.RoleId = &v
}

// GetOtelEndpoint returns the OtelEndpoint field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetOtelEndpoint() string {
	if o == nil || IsNil(o.OtelEndpoint) {
		var ret string
		return ret
	}
	return *o.OtelEndpoint
}

// GetOtelEndpointOk returns a tuple with the OtelEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetOtelEndpointOk() (*string, bool) {
	if o == nil || IsNil(o.OtelEndpoint) {
		return nil, false
	}

	return o.OtelEndpoint, true
}

// HasOtelEndpoint returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasOtelEndpoint() bool {
	if o != nil && !IsNil(o.OtelEndpoint) {
		return true
	}

	return false
}

// SetOtelEndpoint gets a reference to the given string and assigns it to the OtelEndpoint field.
func (o *LogIntegrationRequest) SetOtelEndpoint(v string) {
	o.OtelEndpoint = &v
}

// GetOtelSuppliedHeaders returns the OtelSuppliedHeaders field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetOtelSuppliedHeaders() []Header {
	if o == nil || IsNil(o.OtelSuppliedHeaders) {
		var ret []Header
		return ret
	}
	return *o.OtelSuppliedHeaders
}

// GetOtelSuppliedHeadersOk returns a tuple with the OtelSuppliedHeaders field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetOtelSuppliedHeadersOk() (*[]Header, bool) {
	if o == nil || IsNil(o.OtelSuppliedHeaders) {
		return nil, false
	}

	return o.OtelSuppliedHeaders, true
}

// HasOtelSuppliedHeaders returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasOtelSuppliedHeaders() bool {
	if o != nil && !IsNil(o.OtelSuppliedHeaders) {
		return true
	}

	return false
}

// SetOtelSuppliedHeaders gets a reference to the given []Header and assigns it to the OtelSuppliedHeaders field.
func (o *LogIntegrationRequest) SetOtelSuppliedHeaders(v []Header) {
	o.OtelSuppliedHeaders = &v
}

// GetHecToken returns the HecToken field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetHecToken() string {
	if o == nil || IsNil(o.HecToken) {
		var ret string
		return ret
	}
	return *o.HecToken
}

// GetHecTokenOk returns a tuple with the HecToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetHecTokenOk() (*string, bool) {
	if o == nil || IsNil(o.HecToken) {
		return nil, false
	}

	return o.HecToken, true
}

// HasHecToken returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasHecToken() bool {
	if o != nil && !IsNil(o.HecToken) {
		return true
	}

	return false
}

// SetHecToken gets a reference to the given string and assigns it to the HecToken field.
func (o *LogIntegrationRequest) SetHecToken(v string) {
	o.HecToken = &v
}

// GetHecUrl returns the HecUrl field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetHecUrl() string {
	if o == nil || IsNil(o.HecUrl) {
		var ret string
		return ret
	}
	return *o.HecUrl
}

// GetHecUrlOk returns a tuple with the HecUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetHecUrlOk() (*string, bool) {
	if o == nil || IsNil(o.HecUrl) {
		return nil, false
	}

	return o.HecUrl, true
}

// HasHecUrl returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasHecUrl() bool {
	if o != nil && !IsNil(o.HecUrl) {
		return true
	}

	return false
}

// SetHecUrl gets a reference to the given string and assigns it to the HecUrl field.
func (o *LogIntegrationRequest) SetHecUrl(v string) {
	o.HecUrl = &v
}

// GetStorageAccountName returns the StorageAccountName field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetStorageAccountName() string {
	if o == nil || IsNil(o.StorageAccountName) {
		var ret string
		return ret
	}
	return *o.StorageAccountName
}

// GetStorageAccountNameOk returns a tuple with the StorageAccountName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetStorageAccountNameOk() (*string, bool) {
	if o == nil || IsNil(o.StorageAccountName) {
		return nil, false
	}

	return o.StorageAccountName, true
}

// HasStorageAccountName returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasStorageAccountName() bool {
	if o != nil && !IsNil(o.StorageAccountName) {
		return true
	}

	return false
}

// SetStorageAccountName gets a reference to the given string and assigns it to the StorageAccountName field.
func (o *LogIntegrationRequest) SetStorageAccountName(v string) {
	o.StorageAccountName = &v
}

// GetStorageContainerName returns the StorageContainerName field value if set, zero value otherwise
func (o *LogIntegrationRequest) GetStorageContainerName() string {
	if o == nil || IsNil(o.StorageContainerName) {
		var ret string
		return ret
	}
	return *o.StorageContainerName
}

// GetStorageContainerNameOk returns a tuple with the StorageContainerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LogIntegrationRequest) GetStorageContainerNameOk() (*string, bool) {
	if o == nil || IsNil(o.StorageContainerName) {
		return nil, false
	}

	return o.StorageContainerName, true
}

// HasStorageContainerName returns a boolean if a field has been set.
func (o *LogIntegrationRequest) HasStorageContainerName() bool {
	if o != nil && !IsNil(o.StorageContainerName) {
		return true
	}

	return false
}

// SetStorageContainerName gets a reference to the given string and assigns it to the StorageContainerName field.
func (o *LogIntegrationRequest) SetStorageContainerName(v string) {
	o.StorageContainerName = &v
}
