// Code based on the AtlasAPI V2 OpenAPI file

package admin

// QueryStatsDetailsResponse Metadata and summary statistics for a given query shape.
type QueryStatsDetailsResponse struct {
	FirstSeen  *QueryShapeSeenMetadata `json:"firstSeen,omitempty"`
	LastSeen   *QueryShapeSeenMetadata `json:"lastSeen,omitempty"`
	QueryStats *QueryStatsSummary      `json:"queryStats,omitempty"`
}

// NewQueryStatsDetailsResponse instantiates a new QueryStatsDetailsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQueryStatsDetailsResponse() *QueryStatsDetailsResponse {
	this := QueryStatsDetailsResponse{}
	return &this
}

// NewQueryStatsDetailsResponseWithDefaults instantiates a new QueryStatsDetailsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQueryStatsDetailsResponseWithDefaults() *QueryStatsDetailsResponse {
	this := QueryStatsDetailsResponse{}
	return &this
}

// GetFirstSeen returns the FirstSeen field value if set, zero value otherwise
func (o *QueryStatsDetailsResponse) GetFirstSeen() QueryShapeSeenMetadata {
	if o == nil || IsNil(o.FirstSeen) {
		var ret QueryShapeSeenMetadata
		return ret
	}
	return *o.FirstSeen
}

// GetFirstSeenOk returns a tuple with the FirstSeen field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsDetailsResponse) GetFirstSeenOk() (*QueryShapeSeenMetadata, bool) {
	if o == nil || IsNil(o.FirstSeen) {
		return nil, false
	}

	return o.FirstSeen, true
}

// HasFirstSeen returns a boolean if a field has been set.
func (o *QueryStatsDetailsResponse) HasFirstSeen() bool {
	if o != nil && !IsNil(o.FirstSeen) {
		return true
	}

	return false
}

// SetFirstSeen gets a reference to the given QueryShapeSeenMetadata and assigns it to the FirstSeen field.
func (o *QueryStatsDetailsResponse) SetFirstSeen(v QueryShapeSeenMetadata) {
	o.FirstSeen = &v
}

// GetLastSeen returns the LastSeen field value if set, zero value otherwise
func (o *QueryStatsDetailsResponse) GetLastSeen() QueryShapeSeenMetadata {
	if o == nil || IsNil(o.LastSeen) {
		var ret QueryShapeSeenMetadata
		return ret
	}
	return *o.LastSeen
}

// GetLastSeenOk returns a tuple with the LastSeen field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsDetailsResponse) GetLastSeenOk() (*QueryShapeSeenMetadata, bool) {
	if o == nil || IsNil(o.LastSeen) {
		return nil, false
	}

	return o.LastSeen, true
}

// HasLastSeen returns a boolean if a field has been set.
func (o *QueryStatsDetailsResponse) HasLastSeen() bool {
	if o != nil && !IsNil(o.LastSeen) {
		return true
	}

	return false
}

// SetLastSeen gets a reference to the given QueryShapeSeenMetadata and assigns it to the LastSeen field.
func (o *QueryStatsDetailsResponse) SetLastSeen(v QueryShapeSeenMetadata) {
	o.LastSeen = &v
}

// GetQueryStats returns the QueryStats field value if set, zero value otherwise
func (o *QueryStatsDetailsResponse) GetQueryStats() QueryStatsSummary {
	if o == nil || IsNil(o.QueryStats) {
		var ret QueryStatsSummary
		return ret
	}
	return *o.QueryStats
}

// GetQueryStatsOk returns a tuple with the QueryStats field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsDetailsResponse) GetQueryStatsOk() (*QueryStatsSummary, bool) {
	if o == nil || IsNil(o.QueryStats) {
		return nil, false
	}

	return o.QueryStats, true
}

// HasQueryStats returns a boolean if a field has been set.
func (o *QueryStatsDetailsResponse) HasQueryStats() bool {
	if o != nil && !IsNil(o.QueryStats) {
		return true
	}

	return false
}

// SetQueryStats gets a reference to the given QueryStatsSummary and assigns it to the QueryStats field.
func (o *QueryStatsDetailsResponse) SetQueryStats(v QueryStatsSummary) {
	o.QueryStats = &v
}
