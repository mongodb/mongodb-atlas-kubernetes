// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PerformanceAdvisorSlowQueryList struct for PerformanceAdvisorSlowQueryList
type PerformanceAdvisorSlowQueryList struct {
	// List of operations that the Performance Advisor detected that took longer to execute than a specified threshold.
	// Read only field.
	SlowQueries *[]PerformanceAdvisorSlowQuery `json:"slowQueries,omitempty"`
}

// NewPerformanceAdvisorSlowQueryList instantiates a new PerformanceAdvisorSlowQueryList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPerformanceAdvisorSlowQueryList() *PerformanceAdvisorSlowQueryList {
	this := PerformanceAdvisorSlowQueryList{}
	return &this
}

// NewPerformanceAdvisorSlowQueryListWithDefaults instantiates a new PerformanceAdvisorSlowQueryList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPerformanceAdvisorSlowQueryListWithDefaults() *PerformanceAdvisorSlowQueryList {
	this := PerformanceAdvisorSlowQueryList{}
	return &this
}

// GetSlowQueries returns the SlowQueries field value if set, zero value otherwise
func (o *PerformanceAdvisorSlowQueryList) GetSlowQueries() []PerformanceAdvisorSlowQuery {
	if o == nil || IsNil(o.SlowQueries) {
		var ret []PerformanceAdvisorSlowQuery
		return ret
	}
	return *o.SlowQueries
}

// GetSlowQueriesOk returns a tuple with the SlowQueries field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorSlowQueryList) GetSlowQueriesOk() (*[]PerformanceAdvisorSlowQuery, bool) {
	if o == nil || IsNil(o.SlowQueries) {
		return nil, false
	}

	return o.SlowQueries, true
}

// HasSlowQueries returns a boolean if a field has been set.
func (o *PerformanceAdvisorSlowQueryList) HasSlowQueries() bool {
	if o != nil && !IsNil(o.SlowQueries) {
		return true
	}

	return false
}

// SetSlowQueries gets a reference to the given []PerformanceAdvisorSlowQuery and assigns it to the SlowQueries field.
func (o *PerformanceAdvisorSlowQueryList) SetSlowQueries(v []PerformanceAdvisorSlowQuery) {
	o.SlowQueries = &v
}
