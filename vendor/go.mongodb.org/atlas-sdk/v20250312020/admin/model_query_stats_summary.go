// Code based on the AtlasAPI V2 OpenAPI file

package admin

// QueryStatsSummary A summary of execution statistics for a given query shape.
type QueryStatsSummary struct {
	// Average total time in milliseconds spent running queries with the given query shape. If the query resulted in `getMore` commands, this metric includes the time spent processing the `getMore` requests. This metric does not include time spent waiting for the client.
	AvgWorkingMillis *float64 `json:"avgWorkingMillis,omitempty"`
	// The number of bytes read by the given query shape from the disk to the cache.
	BytesRead *float64 `json:"bytesRead,omitempty"`
	// The MongoDB command issued for this query shape.
	Command *string `json:"command,omitempty"`
	// Total CPU time in nanoseconds consumed by queries with the given query shape. Available for MDB 8.2 and higher.
	CpuTime *float64 `json:"cpuTime,omitempty"`
	// Total number of documents examined by queries with the given query shape.
	DocsExamined *float64 `json:"docsExamined,omitempty"`
	// Ratio of documents examined to documents returned by queries with the given query shape.
	DocsExaminedRatio *float64 `json:"docsExaminedRatio,omitempty"`
	// Total number of documents returned by queries with the given query shape.
	DocsReturned *float64 `json:"docsReturned,omitempty"`
	// Total number of times that queries with the given query shape have been executed.
	ExecCount *float64 `json:"execCount,omitempty"`
	// Total number of in-bounds and out-of-bounds index keys examined by queries with the given query shape.
	KeysExamined *float64 `json:"keysExamined,omitempty"`
	// Ratio of in-bounds and out-of-bounds index keys examined to indexes containing documents returned by queries with the given query shape.
	KeysExaminedRatio *float64 `json:"keysExaminedRatio,omitempty"`
	// Execution runtime in microseconds for the most recent query with the given query shape.
	LastExecMicros *float64 `json:"lastExecMicros,omitempty"`
	// Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
	Namespace *string `json:"namespace,omitempty"`
	// The 50th percentile value of execution time in microseconds.
	P50ExecMicros *float64 `json:"p50ExecMicros,omitempty"`
	// The 90th percentile value of execution time in microseconds.
	P90ExecMicros *float64 `json:"p90ExecMicros,omitempty"`
	// The 99th percentile value of execution time in microseconds.
	P99ExecMicros *float64 `json:"p99ExecMicros,omitempty"`
	// A query shape is a set of specifications that group similar queries together. Specifications can include filters, sorts, projections, aggregation pipeline stages, a namespace, and others. Queries that have similar specifications have the same query shape.
	QueryShape *string `json:"queryShape,omitempty"`
	// A hexadecimal string that represents the hash of a MongoDB query shape.
	QueryShapeHash *string `json:"queryShapeHash,omitempty"`
	// Indicates whether this query shape represents a system-initiated query.
	SystemQuery *bool `json:"systemQuery,omitempty"`
	// Time in microseconds spent from the beginning of query processing to the first server response.
	TotalTimeToResponseMicros *float64 `json:"totalTimeToResponseMicros,omitempty"`
	// Total time in milliseconds spent running queries with the given query shape. If the query resulted in `getMore` commands, this metric includes the time spent processing the `getMore` requests. This metric does not include time spent waiting for the client.
	TotalWorkingMillis *float64 `json:"totalWorkingMillis,omitempty"`
}

// NewQueryStatsSummary instantiates a new QueryStatsSummary object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQueryStatsSummary() *QueryStatsSummary {
	this := QueryStatsSummary{}
	return &this
}

// NewQueryStatsSummaryWithDefaults instantiates a new QueryStatsSummary object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQueryStatsSummaryWithDefaults() *QueryStatsSummary {
	this := QueryStatsSummary{}
	return &this
}

// GetAvgWorkingMillis returns the AvgWorkingMillis field value if set, zero value otherwise
func (o *QueryStatsSummary) GetAvgWorkingMillis() float64 {
	if o == nil || IsNil(o.AvgWorkingMillis) {
		var ret float64
		return ret
	}
	return *o.AvgWorkingMillis
}

// GetAvgWorkingMillisOk returns a tuple with the AvgWorkingMillis field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetAvgWorkingMillisOk() (*float64, bool) {
	if o == nil || IsNil(o.AvgWorkingMillis) {
		return nil, false
	}

	return o.AvgWorkingMillis, true
}

// HasAvgWorkingMillis returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasAvgWorkingMillis() bool {
	if o != nil && !IsNil(o.AvgWorkingMillis) {
		return true
	}

	return false
}

// SetAvgWorkingMillis gets a reference to the given float64 and assigns it to the AvgWorkingMillis field.
func (o *QueryStatsSummary) SetAvgWorkingMillis(v float64) {
	o.AvgWorkingMillis = &v
}

// GetBytesRead returns the BytesRead field value if set, zero value otherwise
func (o *QueryStatsSummary) GetBytesRead() float64 {
	if o == nil || IsNil(o.BytesRead) {
		var ret float64
		return ret
	}
	return *o.BytesRead
}

// GetBytesReadOk returns a tuple with the BytesRead field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetBytesReadOk() (*float64, bool) {
	if o == nil || IsNil(o.BytesRead) {
		return nil, false
	}

	return o.BytesRead, true
}

// HasBytesRead returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasBytesRead() bool {
	if o != nil && !IsNil(o.BytesRead) {
		return true
	}

	return false
}

// SetBytesRead gets a reference to the given float64 and assigns it to the BytesRead field.
func (o *QueryStatsSummary) SetBytesRead(v float64) {
	o.BytesRead = &v
}

// GetCommand returns the Command field value if set, zero value otherwise
func (o *QueryStatsSummary) GetCommand() string {
	if o == nil || IsNil(o.Command) {
		var ret string
		return ret
	}
	return *o.Command
}

// GetCommandOk returns a tuple with the Command field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetCommandOk() (*string, bool) {
	if o == nil || IsNil(o.Command) {
		return nil, false
	}

	return o.Command, true
}

// HasCommand returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasCommand() bool {
	if o != nil && !IsNil(o.Command) {
		return true
	}

	return false
}

// SetCommand gets a reference to the given string and assigns it to the Command field.
func (o *QueryStatsSummary) SetCommand(v string) {
	o.Command = &v
}

// GetCpuTime returns the CpuTime field value if set, zero value otherwise
func (o *QueryStatsSummary) GetCpuTime() float64 {
	if o == nil || IsNil(o.CpuTime) {
		var ret float64
		return ret
	}
	return *o.CpuTime
}

// GetCpuTimeOk returns a tuple with the CpuTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetCpuTimeOk() (*float64, bool) {
	if o == nil || IsNil(o.CpuTime) {
		return nil, false
	}

	return o.CpuTime, true
}

// HasCpuTime returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasCpuTime() bool {
	if o != nil && !IsNil(o.CpuTime) {
		return true
	}

	return false
}

// SetCpuTime gets a reference to the given float64 and assigns it to the CpuTime field.
func (o *QueryStatsSummary) SetCpuTime(v float64) {
	o.CpuTime = &v
}

// GetDocsExamined returns the DocsExamined field value if set, zero value otherwise
func (o *QueryStatsSummary) GetDocsExamined() float64 {
	if o == nil || IsNil(o.DocsExamined) {
		var ret float64
		return ret
	}
	return *o.DocsExamined
}

// GetDocsExaminedOk returns a tuple with the DocsExamined field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetDocsExaminedOk() (*float64, bool) {
	if o == nil || IsNil(o.DocsExamined) {
		return nil, false
	}

	return o.DocsExamined, true
}

// HasDocsExamined returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasDocsExamined() bool {
	if o != nil && !IsNil(o.DocsExamined) {
		return true
	}

	return false
}

// SetDocsExamined gets a reference to the given float64 and assigns it to the DocsExamined field.
func (o *QueryStatsSummary) SetDocsExamined(v float64) {
	o.DocsExamined = &v
}

// GetDocsExaminedRatio returns the DocsExaminedRatio field value if set, zero value otherwise
func (o *QueryStatsSummary) GetDocsExaminedRatio() float64 {
	if o == nil || IsNil(o.DocsExaminedRatio) {
		var ret float64
		return ret
	}
	return *o.DocsExaminedRatio
}

// GetDocsExaminedRatioOk returns a tuple with the DocsExaminedRatio field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetDocsExaminedRatioOk() (*float64, bool) {
	if o == nil || IsNil(o.DocsExaminedRatio) {
		return nil, false
	}

	return o.DocsExaminedRatio, true
}

// HasDocsExaminedRatio returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasDocsExaminedRatio() bool {
	if o != nil && !IsNil(o.DocsExaminedRatio) {
		return true
	}

	return false
}

// SetDocsExaminedRatio gets a reference to the given float64 and assigns it to the DocsExaminedRatio field.
func (o *QueryStatsSummary) SetDocsExaminedRatio(v float64) {
	o.DocsExaminedRatio = &v
}

// GetDocsReturned returns the DocsReturned field value if set, zero value otherwise
func (o *QueryStatsSummary) GetDocsReturned() float64 {
	if o == nil || IsNil(o.DocsReturned) {
		var ret float64
		return ret
	}
	return *o.DocsReturned
}

// GetDocsReturnedOk returns a tuple with the DocsReturned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetDocsReturnedOk() (*float64, bool) {
	if o == nil || IsNil(o.DocsReturned) {
		return nil, false
	}

	return o.DocsReturned, true
}

// HasDocsReturned returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasDocsReturned() bool {
	if o != nil && !IsNil(o.DocsReturned) {
		return true
	}

	return false
}

// SetDocsReturned gets a reference to the given float64 and assigns it to the DocsReturned field.
func (o *QueryStatsSummary) SetDocsReturned(v float64) {
	o.DocsReturned = &v
}

// GetExecCount returns the ExecCount field value if set, zero value otherwise
func (o *QueryStatsSummary) GetExecCount() float64 {
	if o == nil || IsNil(o.ExecCount) {
		var ret float64
		return ret
	}
	return *o.ExecCount
}

// GetExecCountOk returns a tuple with the ExecCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetExecCountOk() (*float64, bool) {
	if o == nil || IsNil(o.ExecCount) {
		return nil, false
	}

	return o.ExecCount, true
}

// HasExecCount returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasExecCount() bool {
	if o != nil && !IsNil(o.ExecCount) {
		return true
	}

	return false
}

// SetExecCount gets a reference to the given float64 and assigns it to the ExecCount field.
func (o *QueryStatsSummary) SetExecCount(v float64) {
	o.ExecCount = &v
}

// GetKeysExamined returns the KeysExamined field value if set, zero value otherwise
func (o *QueryStatsSummary) GetKeysExamined() float64 {
	if o == nil || IsNil(o.KeysExamined) {
		var ret float64
		return ret
	}
	return *o.KeysExamined
}

// GetKeysExaminedOk returns a tuple with the KeysExamined field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetKeysExaminedOk() (*float64, bool) {
	if o == nil || IsNil(o.KeysExamined) {
		return nil, false
	}

	return o.KeysExamined, true
}

// HasKeysExamined returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasKeysExamined() bool {
	if o != nil && !IsNil(o.KeysExamined) {
		return true
	}

	return false
}

// SetKeysExamined gets a reference to the given float64 and assigns it to the KeysExamined field.
func (o *QueryStatsSummary) SetKeysExamined(v float64) {
	o.KeysExamined = &v
}

// GetKeysExaminedRatio returns the KeysExaminedRatio field value if set, zero value otherwise
func (o *QueryStatsSummary) GetKeysExaminedRatio() float64 {
	if o == nil || IsNil(o.KeysExaminedRatio) {
		var ret float64
		return ret
	}
	return *o.KeysExaminedRatio
}

// GetKeysExaminedRatioOk returns a tuple with the KeysExaminedRatio field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetKeysExaminedRatioOk() (*float64, bool) {
	if o == nil || IsNil(o.KeysExaminedRatio) {
		return nil, false
	}

	return o.KeysExaminedRatio, true
}

// HasKeysExaminedRatio returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasKeysExaminedRatio() bool {
	if o != nil && !IsNil(o.KeysExaminedRatio) {
		return true
	}

	return false
}

// SetKeysExaminedRatio gets a reference to the given float64 and assigns it to the KeysExaminedRatio field.
func (o *QueryStatsSummary) SetKeysExaminedRatio(v float64) {
	o.KeysExaminedRatio = &v
}

// GetLastExecMicros returns the LastExecMicros field value if set, zero value otherwise
func (o *QueryStatsSummary) GetLastExecMicros() float64 {
	if o == nil || IsNil(o.LastExecMicros) {
		var ret float64
		return ret
	}
	return *o.LastExecMicros
}

// GetLastExecMicrosOk returns a tuple with the LastExecMicros field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetLastExecMicrosOk() (*float64, bool) {
	if o == nil || IsNil(o.LastExecMicros) {
		return nil, false
	}

	return o.LastExecMicros, true
}

// HasLastExecMicros returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasLastExecMicros() bool {
	if o != nil && !IsNil(o.LastExecMicros) {
		return true
	}

	return false
}

// SetLastExecMicros gets a reference to the given float64 and assigns it to the LastExecMicros field.
func (o *QueryStatsSummary) SetLastExecMicros(v float64) {
	o.LastExecMicros = &v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *QueryStatsSummary) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *QueryStatsSummary) SetNamespace(v string) {
	o.Namespace = &v
}

// GetP50ExecMicros returns the P50ExecMicros field value if set, zero value otherwise
func (o *QueryStatsSummary) GetP50ExecMicros() float64 {
	if o == nil || IsNil(o.P50ExecMicros) {
		var ret float64
		return ret
	}
	return *o.P50ExecMicros
}

// GetP50ExecMicrosOk returns a tuple with the P50ExecMicros field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetP50ExecMicrosOk() (*float64, bool) {
	if o == nil || IsNil(o.P50ExecMicros) {
		return nil, false
	}

	return o.P50ExecMicros, true
}

// HasP50ExecMicros returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasP50ExecMicros() bool {
	if o != nil && !IsNil(o.P50ExecMicros) {
		return true
	}

	return false
}

// SetP50ExecMicros gets a reference to the given float64 and assigns it to the P50ExecMicros field.
func (o *QueryStatsSummary) SetP50ExecMicros(v float64) {
	o.P50ExecMicros = &v
}

// GetP90ExecMicros returns the P90ExecMicros field value if set, zero value otherwise
func (o *QueryStatsSummary) GetP90ExecMicros() float64 {
	if o == nil || IsNil(o.P90ExecMicros) {
		var ret float64
		return ret
	}
	return *o.P90ExecMicros
}

// GetP90ExecMicrosOk returns a tuple with the P90ExecMicros field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetP90ExecMicrosOk() (*float64, bool) {
	if o == nil || IsNil(o.P90ExecMicros) {
		return nil, false
	}

	return o.P90ExecMicros, true
}

// HasP90ExecMicros returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasP90ExecMicros() bool {
	if o != nil && !IsNil(o.P90ExecMicros) {
		return true
	}

	return false
}

// SetP90ExecMicros gets a reference to the given float64 and assigns it to the P90ExecMicros field.
func (o *QueryStatsSummary) SetP90ExecMicros(v float64) {
	o.P90ExecMicros = &v
}

// GetP99ExecMicros returns the P99ExecMicros field value if set, zero value otherwise
func (o *QueryStatsSummary) GetP99ExecMicros() float64 {
	if o == nil || IsNil(o.P99ExecMicros) {
		var ret float64
		return ret
	}
	return *o.P99ExecMicros
}

// GetP99ExecMicrosOk returns a tuple with the P99ExecMicros field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetP99ExecMicrosOk() (*float64, bool) {
	if o == nil || IsNil(o.P99ExecMicros) {
		return nil, false
	}

	return o.P99ExecMicros, true
}

// HasP99ExecMicros returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasP99ExecMicros() bool {
	if o != nil && !IsNil(o.P99ExecMicros) {
		return true
	}

	return false
}

// SetP99ExecMicros gets a reference to the given float64 and assigns it to the P99ExecMicros field.
func (o *QueryStatsSummary) SetP99ExecMicros(v float64) {
	o.P99ExecMicros = &v
}

// GetQueryShape returns the QueryShape field value if set, zero value otherwise
func (o *QueryStatsSummary) GetQueryShape() string {
	if o == nil || IsNil(o.QueryShape) {
		var ret string
		return ret
	}
	return *o.QueryShape
}

// GetQueryShapeOk returns a tuple with the QueryShape field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetQueryShapeOk() (*string, bool) {
	if o == nil || IsNil(o.QueryShape) {
		return nil, false
	}

	return o.QueryShape, true
}

// HasQueryShape returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasQueryShape() bool {
	if o != nil && !IsNil(o.QueryShape) {
		return true
	}

	return false
}

// SetQueryShape gets a reference to the given string and assigns it to the QueryShape field.
func (o *QueryStatsSummary) SetQueryShape(v string) {
	o.QueryShape = &v
}

// GetQueryShapeHash returns the QueryShapeHash field value if set, zero value otherwise
func (o *QueryStatsSummary) GetQueryShapeHash() string {
	if o == nil || IsNil(o.QueryShapeHash) {
		var ret string
		return ret
	}
	return *o.QueryShapeHash
}

// GetQueryShapeHashOk returns a tuple with the QueryShapeHash field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetQueryShapeHashOk() (*string, bool) {
	if o == nil || IsNil(o.QueryShapeHash) {
		return nil, false
	}

	return o.QueryShapeHash, true
}

// HasQueryShapeHash returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasQueryShapeHash() bool {
	if o != nil && !IsNil(o.QueryShapeHash) {
		return true
	}

	return false
}

// SetQueryShapeHash gets a reference to the given string and assigns it to the QueryShapeHash field.
func (o *QueryStatsSummary) SetQueryShapeHash(v string) {
	o.QueryShapeHash = &v
}

// GetSystemQuery returns the SystemQuery field value if set, zero value otherwise
func (o *QueryStatsSummary) GetSystemQuery() bool {
	if o == nil || IsNil(o.SystemQuery) {
		var ret bool
		return ret
	}
	return *o.SystemQuery
}

// GetSystemQueryOk returns a tuple with the SystemQuery field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetSystemQueryOk() (*bool, bool) {
	if o == nil || IsNil(o.SystemQuery) {
		return nil, false
	}

	return o.SystemQuery, true
}

// HasSystemQuery returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasSystemQuery() bool {
	if o != nil && !IsNil(o.SystemQuery) {
		return true
	}

	return false
}

// SetSystemQuery gets a reference to the given bool and assigns it to the SystemQuery field.
func (o *QueryStatsSummary) SetSystemQuery(v bool) {
	o.SystemQuery = &v
}

// GetTotalTimeToResponseMicros returns the TotalTimeToResponseMicros field value if set, zero value otherwise
func (o *QueryStatsSummary) GetTotalTimeToResponseMicros() float64 {
	if o == nil || IsNil(o.TotalTimeToResponseMicros) {
		var ret float64
		return ret
	}
	return *o.TotalTimeToResponseMicros
}

// GetTotalTimeToResponseMicrosOk returns a tuple with the TotalTimeToResponseMicros field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetTotalTimeToResponseMicrosOk() (*float64, bool) {
	if o == nil || IsNil(o.TotalTimeToResponseMicros) {
		return nil, false
	}

	return o.TotalTimeToResponseMicros, true
}

// HasTotalTimeToResponseMicros returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasTotalTimeToResponseMicros() bool {
	if o != nil && !IsNil(o.TotalTimeToResponseMicros) {
		return true
	}

	return false
}

// SetTotalTimeToResponseMicros gets a reference to the given float64 and assigns it to the TotalTimeToResponseMicros field.
func (o *QueryStatsSummary) SetTotalTimeToResponseMicros(v float64) {
	o.TotalTimeToResponseMicros = &v
}

// GetTotalWorkingMillis returns the TotalWorkingMillis field value if set, zero value otherwise
func (o *QueryStatsSummary) GetTotalWorkingMillis() float64 {
	if o == nil || IsNil(o.TotalWorkingMillis) {
		var ret float64
		return ret
	}
	return *o.TotalWorkingMillis
}

// GetTotalWorkingMillisOk returns a tuple with the TotalWorkingMillis field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummary) GetTotalWorkingMillisOk() (*float64, bool) {
	if o == nil || IsNil(o.TotalWorkingMillis) {
		return nil, false
	}

	return o.TotalWorkingMillis, true
}

// HasTotalWorkingMillis returns a boolean if a field has been set.
func (o *QueryStatsSummary) HasTotalWorkingMillis() bool {
	if o != nil && !IsNil(o.TotalWorkingMillis) {
		return true
	}

	return false
}

// SetTotalWorkingMillis gets a reference to the given float64 and assigns it to the TotalWorkingMillis field.
func (o *QueryStatsSummary) SetTotalWorkingMillis(v float64) {
	o.TotalWorkingMillis = &v
}
