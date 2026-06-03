// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PerformanceAdvisorSlowQueryMetrics Metrics from a slow query log.
type PerformanceAdvisorSlowQueryMetrics struct {
	// The number of documents in the collection that MongoDB scanned in order to carry out the operation.
	// Read only field.
	DocsExamined *int64 `json:"docsExamined,omitempty"`
	// Ratio of documents examined to documents returned.
	// Read only field.
	DocsExaminedReturnedRatio *float64 `json:"docsExaminedReturnedRatio,omitempty"`
	// The number of documents returned by the operation.
	// Read only field.
	DocsReturned *int64 `json:"docsReturned,omitempty"`
	// This boolean will be true when the server can identify the query source as non-server. This field is only available for MDB 8.0+.
	// Read only field.
	FromUserConnection *bool `json:"fromUserConnection,omitempty"`
	// Indicates if the query has index coverage.
	// Read only field.
	HasIndexCoverage *bool `json:"hasIndexCoverage,omitempty"`
	// This boolean will be true when a query cannot use the ordering in the index to return the requested sorted results; i.e. MongoDB must sort the documents after it receives the documents from a cursor.
	// Read only field.
	HasSort *bool `json:"hasSort,omitempty"`
	// The number of index keys that MongoDB scanned in order to carry out the operation.
	// Read only field.
	KeysExamined *int64 `json:"keysExamined,omitempty"`
	// Ratio of keys examined to documents returned.
	// Read only field.
	KeysExaminedReturnedRatio *float64 `json:"keysExaminedReturnedRatio,omitempty"`
	// The number of times the operation yielded to allow other operations to complete.
	// Read only field.
	NumYields *int64 `json:"numYields,omitempty"`
	// Total execution time of a query in milliseconds.
	// Read only field.
	OperationExecutionTime *int64 `json:"operationExecutionTime,omitempty"`
	// The length in bytes of the operation's result document.
	// Read only field.
	ResponseLength *int64 `json:"responseLength,omitempty"`
}

// NewPerformanceAdvisorSlowQueryMetrics instantiates a new PerformanceAdvisorSlowQueryMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPerformanceAdvisorSlowQueryMetrics() *PerformanceAdvisorSlowQueryMetrics {
	this := PerformanceAdvisorSlowQueryMetrics{}
	return &this
}

// NewPerformanceAdvisorSlowQueryMetricsWithDefaults instantiates a new PerformanceAdvisorSlowQueryMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPerformanceAdvisorSlowQueryMetricsWithDefaults() *PerformanceAdvisorSlowQueryMetrics {
	this := PerformanceAdvisorSlowQueryMetrics{}
	return &this
}

// GetDocsExamined returns the DocsExamined field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetDocsExamined() int64 {
	if o == nil || IsNil(o.DocsExamined) {
		var ret int64
		return ret
	}
	return *o.DocsExamined
}

// GetDocsExaminedOk returns a tuple with the DocsExamined field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetDocsExaminedOk() (*int64, bool) {
	if o == nil || IsNil(o.DocsExamined) {
		return nil, false
	}

	return o.DocsExamined, true
}

// HasDocsExamined returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasDocsExamined() bool {
	if o != nil && !IsNil(o.DocsExamined) {
		return true
	}

	return false
}

// SetDocsExamined gets a reference to the given int64 and assigns it to the DocsExamined field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetDocsExamined(v int64) {
	o.DocsExamined = &v
}

// GetDocsExaminedReturnedRatio returns the DocsExaminedReturnedRatio field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetDocsExaminedReturnedRatio() float64 {
	if o == nil || IsNil(o.DocsExaminedReturnedRatio) {
		var ret float64
		return ret
	}
	return *o.DocsExaminedReturnedRatio
}

// GetDocsExaminedReturnedRatioOk returns a tuple with the DocsExaminedReturnedRatio field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetDocsExaminedReturnedRatioOk() (*float64, bool) {
	if o == nil || IsNil(o.DocsExaminedReturnedRatio) {
		return nil, false
	}

	return o.DocsExaminedReturnedRatio, true
}

// HasDocsExaminedReturnedRatio returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasDocsExaminedReturnedRatio() bool {
	if o != nil && !IsNil(o.DocsExaminedReturnedRatio) {
		return true
	}

	return false
}

// SetDocsExaminedReturnedRatio gets a reference to the given float64 and assigns it to the DocsExaminedReturnedRatio field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetDocsExaminedReturnedRatio(v float64) {
	o.DocsExaminedReturnedRatio = &v
}

// GetDocsReturned returns the DocsReturned field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetDocsReturned() int64 {
	if o == nil || IsNil(o.DocsReturned) {
		var ret int64
		return ret
	}
	return *o.DocsReturned
}

// GetDocsReturnedOk returns a tuple with the DocsReturned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetDocsReturnedOk() (*int64, bool) {
	if o == nil || IsNil(o.DocsReturned) {
		return nil, false
	}

	return o.DocsReturned, true
}

// HasDocsReturned returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasDocsReturned() bool {
	if o != nil && !IsNil(o.DocsReturned) {
		return true
	}

	return false
}

// SetDocsReturned gets a reference to the given int64 and assigns it to the DocsReturned field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetDocsReturned(v int64) {
	o.DocsReturned = &v
}

// GetFromUserConnection returns the FromUserConnection field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetFromUserConnection() bool {
	if o == nil || IsNil(o.FromUserConnection) {
		var ret bool
		return ret
	}
	return *o.FromUserConnection
}

// GetFromUserConnectionOk returns a tuple with the FromUserConnection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetFromUserConnectionOk() (*bool, bool) {
	if o == nil || IsNil(o.FromUserConnection) {
		return nil, false
	}

	return o.FromUserConnection, true
}

// HasFromUserConnection returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasFromUserConnection() bool {
	if o != nil && !IsNil(o.FromUserConnection) {
		return true
	}

	return false
}

// SetFromUserConnection gets a reference to the given bool and assigns it to the FromUserConnection field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetFromUserConnection(v bool) {
	o.FromUserConnection = &v
}

// GetHasIndexCoverage returns the HasIndexCoverage field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetHasIndexCoverage() bool {
	if o == nil || IsNil(o.HasIndexCoverage) {
		var ret bool
		return ret
	}
	return *o.HasIndexCoverage
}

// GetHasIndexCoverageOk returns a tuple with the HasIndexCoverage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetHasIndexCoverageOk() (*bool, bool) {
	if o == nil || IsNil(o.HasIndexCoverage) {
		return nil, false
	}

	return o.HasIndexCoverage, true
}

// HasHasIndexCoverage returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasHasIndexCoverage() bool {
	if o != nil && !IsNil(o.HasIndexCoverage) {
		return true
	}

	return false
}

// SetHasIndexCoverage gets a reference to the given bool and assigns it to the HasIndexCoverage field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetHasIndexCoverage(v bool) {
	o.HasIndexCoverage = &v
}

// GetHasSort returns the HasSort field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetHasSort() bool {
	if o == nil || IsNil(o.HasSort) {
		var ret bool
		return ret
	}
	return *o.HasSort
}

// GetHasSortOk returns a tuple with the HasSort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetHasSortOk() (*bool, bool) {
	if o == nil || IsNil(o.HasSort) {
		return nil, false
	}

	return o.HasSort, true
}

// HasHasSort returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasHasSort() bool {
	if o != nil && !IsNil(o.HasSort) {
		return true
	}

	return false
}

// SetHasSort gets a reference to the given bool and assigns it to the HasSort field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetHasSort(v bool) {
	o.HasSort = &v
}

// GetKeysExamined returns the KeysExamined field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetKeysExamined() int64 {
	if o == nil || IsNil(o.KeysExamined) {
		var ret int64
		return ret
	}
	return *o.KeysExamined
}

// GetKeysExaminedOk returns a tuple with the KeysExamined field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetKeysExaminedOk() (*int64, bool) {
	if o == nil || IsNil(o.KeysExamined) {
		return nil, false
	}

	return o.KeysExamined, true
}

// HasKeysExamined returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasKeysExamined() bool {
	if o != nil && !IsNil(o.KeysExamined) {
		return true
	}

	return false
}

// SetKeysExamined gets a reference to the given int64 and assigns it to the KeysExamined field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetKeysExamined(v int64) {
	o.KeysExamined = &v
}

// GetKeysExaminedReturnedRatio returns the KeysExaminedReturnedRatio field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetKeysExaminedReturnedRatio() float64 {
	if o == nil || IsNil(o.KeysExaminedReturnedRatio) {
		var ret float64
		return ret
	}
	return *o.KeysExaminedReturnedRatio
}

// GetKeysExaminedReturnedRatioOk returns a tuple with the KeysExaminedReturnedRatio field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetKeysExaminedReturnedRatioOk() (*float64, bool) {
	if o == nil || IsNil(o.KeysExaminedReturnedRatio) {
		return nil, false
	}

	return o.KeysExaminedReturnedRatio, true
}

// HasKeysExaminedReturnedRatio returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasKeysExaminedReturnedRatio() bool {
	if o != nil && !IsNil(o.KeysExaminedReturnedRatio) {
		return true
	}

	return false
}

// SetKeysExaminedReturnedRatio gets a reference to the given float64 and assigns it to the KeysExaminedReturnedRatio field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetKeysExaminedReturnedRatio(v float64) {
	o.KeysExaminedReturnedRatio = &v
}

// GetNumYields returns the NumYields field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetNumYields() int64 {
	if o == nil || IsNil(o.NumYields) {
		var ret int64
		return ret
	}
	return *o.NumYields
}

// GetNumYieldsOk returns a tuple with the NumYields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetNumYieldsOk() (*int64, bool) {
	if o == nil || IsNil(o.NumYields) {
		return nil, false
	}

	return o.NumYields, true
}

// HasNumYields returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasNumYields() bool {
	if o != nil && !IsNil(o.NumYields) {
		return true
	}

	return false
}

// SetNumYields gets a reference to the given int64 and assigns it to the NumYields field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetNumYields(v int64) {
	o.NumYields = &v
}

// GetOperationExecutionTime returns the OperationExecutionTime field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetOperationExecutionTime() int64 {
	if o == nil || IsNil(o.OperationExecutionTime) {
		var ret int64
		return ret
	}
	return *o.OperationExecutionTime
}

// GetOperationExecutionTimeOk returns a tuple with the OperationExecutionTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetOperationExecutionTimeOk() (*int64, bool) {
	if o == nil || IsNil(o.OperationExecutionTime) {
		return nil, false
	}

	return o.OperationExecutionTime, true
}

// HasOperationExecutionTime returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasOperationExecutionTime() bool {
	if o != nil && !IsNil(o.OperationExecutionTime) {
		return true
	}

	return false
}

// SetOperationExecutionTime gets a reference to the given int64 and assigns it to the OperationExecutionTime field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetOperationExecutionTime(v int64) {
	o.OperationExecutionTime = &v
}

// GetResponseLength returns the ResponseLength field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryMetrics) GetResponseLength() int64 {
	if o == nil || IsNil(o.ResponseLength) {
		var ret int64
		return ret
	}
	return *o.ResponseLength
}

// GetResponseLengthOk returns a tuple with the ResponseLength field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) GetResponseLengthOk() (*int64, bool) {
	if o == nil || IsNil(o.ResponseLength) {
		return nil, false
	}

	return o.ResponseLength, true
}

// HasResponseLength returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryMetrics) HasResponseLength() bool {
	if o != nil && !IsNil(o.ResponseLength) {
		return true
	}

	return false
}

// SetResponseLength gets a reference to the given int64 and assigns it to the ResponseLength field.
func (o *PerformanceAdvisorSlowQueryMetrics) SetResponseLength(v int64) {
	o.ResponseLength = &v
}
