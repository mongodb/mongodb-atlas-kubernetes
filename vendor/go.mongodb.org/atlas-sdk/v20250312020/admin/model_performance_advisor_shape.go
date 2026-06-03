// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PerformanceAdvisorShape struct for PerformanceAdvisorShape
type PerformanceAdvisorShape struct {
	// Average duration in milliseconds for the queries examined that match this shape.
	// Read only field.
	AvgMs *int64 `json:"avgMs,omitempty"`
	// Number of queries examined that match this shape.
	// Read only field.
	Count *int64 `json:"count,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this shape. This string exists only for the duration of this API request.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Average number of documents read for every document that the query returns.
	// Read only field.
	InefficiencyScore *int64 `json:"inefficiencyScore,omitempty"`
	// Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
	// Read only field.
	Namespace *string `json:"namespace,omitempty"`
	// List that contains specific about individual queries.
	// Read only field.
	Operations *[]PerformanceAdvisorOperation `json:"operations,omitempty"`
}

// NewPerformanceAdvisorShape instantiates a new PerformanceAdvisorShape object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPerformanceAdvisorShape() *PerformanceAdvisorShape {
	this := PerformanceAdvisorShape{}
	return &this
}

// NewPerformanceAdvisorShapeWithDefaults instantiates a new PerformanceAdvisorShape object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPerformanceAdvisorShapeWithDefaults() *PerformanceAdvisorShape {
	this := PerformanceAdvisorShape{}
	return &this
}

// GetAvgMs returns the AvgMs field value if set, zero value otherwise
func (o *PerformanceAdvisorShape) GetAvgMs() int64 {
	if o == nil || IsNil(o.AvgMs) {
		var ret int64
		return ret
	}
	return *o.AvgMs
}

// GetAvgMsOk returns a tuple with the AvgMs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorShape) GetAvgMsOk() (*int64, bool) {
	if o == nil || IsNil(o.AvgMs) {
		return nil, false
	}

	return o.AvgMs, true
}

// HasAvgMs returns a boolean if a field has been set.
func (o *PerformanceAdvisorShape) HasAvgMs() bool {
	if o != nil && !IsNil(o.AvgMs) {
		return true
	}

	return false
}

// SetAvgMs gets a reference to the given int64 and assigns it to the AvgMs field.
func (o *PerformanceAdvisorShape) SetAvgMs(v int64) {
	o.AvgMs = &v
}

// GetCount returns the Count field value if set, zero value otherwise
func (o *PerformanceAdvisorShape) GetCount() int64 {
	if o == nil || IsNil(o.Count) {
		var ret int64
		return ret
	}
	return *o.Count
}

// GetCountOk returns a tuple with the Count field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorShape) GetCountOk() (*int64, bool) {
	if o == nil || IsNil(o.Count) {
		return nil, false
	}

	return o.Count, true
}

// HasCount returns a boolean if a field has been set.
func (o *PerformanceAdvisorShape) HasCount() bool {
	if o != nil && !IsNil(o.Count) {
		return true
	}

	return false
}

// SetCount gets a reference to the given int64 and assigns it to the Count field.
func (o *PerformanceAdvisorShape) SetCount(v int64) {
	o.Count = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *PerformanceAdvisorShape) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorShape) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *PerformanceAdvisorShape) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *PerformanceAdvisorShape) SetId(v string) {
	o.Id = &v
}

// GetInefficiencyScore returns the InefficiencyScore field value if set, zero value otherwise
func (o *PerformanceAdvisorShape) GetInefficiencyScore() int64 {
	if o == nil || IsNil(o.InefficiencyScore) {
		var ret int64
		return ret
	}
	return *o.InefficiencyScore
}

// GetInefficiencyScoreOk returns a tuple with the InefficiencyScore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorShape) GetInefficiencyScoreOk() (*int64, bool) {
	if o == nil || IsNil(o.InefficiencyScore) {
		return nil, false
	}

	return o.InefficiencyScore, true
}

// HasInefficiencyScore returns a boolean if a field has been set.
func (o *PerformanceAdvisorShape) HasInefficiencyScore() bool {
	if o != nil && !IsNil(o.InefficiencyScore) {
		return true
	}

	return false
}

// SetInefficiencyScore gets a reference to the given int64 and assigns it to the InefficiencyScore field.
func (o *PerformanceAdvisorShape) SetInefficiencyScore(v int64) {
	o.InefficiencyScore = &v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *PerformanceAdvisorShape) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorShape) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *PerformanceAdvisorShape) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *PerformanceAdvisorShape) SetNamespace(v string) {
	o.Namespace = &v
}

// GetOperations returns the Operations field value if set, zero value otherwise
func (o *PerformanceAdvisorShape) GetOperations() []PerformanceAdvisorOperation {
	if o == nil || IsNil(o.Operations) {
		var ret []PerformanceAdvisorOperation
		return ret
	}
	return *o.Operations
}

// GetOperationsOk returns a tuple with the Operations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorShape) GetOperationsOk() (*[]PerformanceAdvisorOperation, bool) {
	if o == nil || IsNil(o.Operations) {
		return nil, false
	}

	return o.Operations, true
}

// HasOperations returns a boolean if a field has been set.
func (o *PerformanceAdvisorShape) HasOperations() bool {
	if o != nil && !IsNil(o.Operations) {
		return true
	}

	return false
}

// SetOperations gets a reference to the given []PerformanceAdvisorOperation and assigns it to the Operations field.
func (o *PerformanceAdvisorShape) SetOperations(v []PerformanceAdvisorOperation) {
	o.Operations = &v
}
