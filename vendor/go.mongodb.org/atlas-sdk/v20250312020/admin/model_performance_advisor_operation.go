// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PerformanceAdvisorOperation struct for PerformanceAdvisorOperation
type PerformanceAdvisorOperation struct {
	// List that contains the search criteria that the query uses. To use the values in key-value pairs in these predicates requires **Project Data Access Read Only** permissions or greater. Otherwise, MongoDB Cloud redacts these values.
	// Read only field.
	Predicates *[]any                     `json:"predicates,omitempty"`
	Stats      *PerformanceAdvisorOpStats `json:"stats,omitempty"`
}

// NewPerformanceAdvisorOperation instantiates a new PerformanceAdvisorOperation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPerformanceAdvisorOperation() *PerformanceAdvisorOperation {
	this := PerformanceAdvisorOperation{}
	return &this
}

// NewPerformanceAdvisorOperationWithDefaults instantiates a new PerformanceAdvisorOperation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPerformanceAdvisorOperationWithDefaults() *PerformanceAdvisorOperation {
	this := PerformanceAdvisorOperation{}
	return &this
}

// GetPredicates returns the Predicates field value if set, zero value otherwise
func (o *PerformanceAdvisorOperation) GetPredicates() []any {
	if o == nil || IsNil(o.Predicates) {
		var ret []any
		return ret
	}
	return *o.Predicates
}

// GetPredicatesOk returns a tuple with the Predicates field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorOperation) GetPredicatesOk() (*[]any, bool) {
	if o == nil || IsNil(o.Predicates) {
		return nil, false
	}

	return o.Predicates, true
}

// HasPredicates returns a boolean if a field has been set.
func (o *PerformanceAdvisorOperation) HasPredicates() bool {
	if o != nil && !IsNil(o.Predicates) {
		return true
	}

	return false
}

// SetPredicates gets a reference to the given []any and assigns it to the Predicates field.
func (o *PerformanceAdvisorOperation) SetPredicates(v []any) {
	o.Predicates = &v
}

// GetStats returns the Stats field value if set, zero value otherwise
func (o *PerformanceAdvisorOperation) GetStats() PerformanceAdvisorOpStats {
	if o == nil || IsNil(o.Stats) {
		var ret PerformanceAdvisorOpStats
		return ret
	}
	return *o.Stats
}

// GetStatsOk returns a tuple with the Stats field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorOperation) GetStatsOk() (*PerformanceAdvisorOpStats, bool) {
	if o == nil || IsNil(o.Stats) {
		return nil, false
	}

	return o.Stats, true
}

// HasStats returns a boolean if a field has been set.
func (o *PerformanceAdvisorOperation) HasStats() bool {
	if o != nil && !IsNil(o.Stats) {
		return true
	}

	return false
}

// SetStats gets a reference to the given PerformanceAdvisorOpStats and assigns it to the Stats field.
func (o *PerformanceAdvisorOperation) SetStats(v PerformanceAdvisorOpStats) {
	o.Stats = &v
}
