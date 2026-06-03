// Code based on the AtlasAPI V2 OpenAPI file

package admin

// QueryStatsSummaryListResponse struct for QueryStatsSummaryListResponse
type QueryStatsSummaryListResponse struct {
	// List of query shape statistic summaries from Query Shape Insights.
	Summaries *[]QueryStatsSummary `json:"summaries,omitempty"`
}

// NewQueryStatsSummaryListResponse instantiates a new QueryStatsSummaryListResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQueryStatsSummaryListResponse() *QueryStatsSummaryListResponse {
	this := QueryStatsSummaryListResponse{}
	return &this
}

// NewQueryStatsSummaryListResponseWithDefaults instantiates a new QueryStatsSummaryListResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQueryStatsSummaryListResponseWithDefaults() *QueryStatsSummaryListResponse {
	this := QueryStatsSummaryListResponse{}
	return &this
}

// GetSummaries returns the Summaries field value if set, zero value otherwise
func (o *QueryStatsSummaryListResponse) GetSummaries() []QueryStatsSummary {
	if o == nil || IsNil(o.Summaries) {
		var ret []QueryStatsSummary
		return ret
	}
	return *o.Summaries
}

// GetSummariesOk returns a tuple with the Summaries field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryStatsSummaryListResponse) GetSummariesOk() (*[]QueryStatsSummary, bool) {
	if o == nil || IsNil(o.Summaries) {
		return nil, false
	}

	return o.Summaries, true
}

// HasSummaries returns a boolean if a field has been set.
func (o *QueryStatsSummaryListResponse) HasSummaries() bool {
	if o != nil && !IsNil(o.Summaries) {
		return true
	}

	return false
}

// SetSummaries gets a reference to the given []QueryStatsSummary and assigns it to the Summaries field.
func (o *QueryStatsSummaryListResponse) SetSummaries(v []QueryStatsSummary) {
	o.Summaries = &v
}
