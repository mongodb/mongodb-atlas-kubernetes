// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CollStatsLatencyNamespaceMetrics struct for CollStatsLatencyNamespaceMetrics
type CollStatsLatencyNamespaceMetrics struct {
	// Unique 24-hexadecimal digit string that identifies the project.
	// Read only field.
	GroupId string `json:"groupId"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of Coll Stats Latency metric names and their respective units.
	// Read only field.
	Metrics []CollStatsLatencyNamespaceMetric `json:"metrics"`
}

// NewCollStatsLatencyNamespaceMetrics instantiates a new CollStatsLatencyNamespaceMetrics object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCollStatsLatencyNamespaceMetrics(groupId string, metrics []CollStatsLatencyNamespaceMetric) *CollStatsLatencyNamespaceMetrics {
	this := CollStatsLatencyNamespaceMetrics{}
	this.GroupId = groupId
	this.Metrics = metrics
	return &this
}

// NewCollStatsLatencyNamespaceMetricsWithDefaults instantiates a new CollStatsLatencyNamespaceMetrics object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCollStatsLatencyNamespaceMetricsWithDefaults() *CollStatsLatencyNamespaceMetrics {
	this := CollStatsLatencyNamespaceMetrics{}
	return &this
}

// GetGroupId returns the GroupId field value
func (o *CollStatsLatencyNamespaceMetrics) GetGroupId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value
// and a boolean to check if the value has been set.
func (o *CollStatsLatencyNamespaceMetrics) GetGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GroupId, true
}

// SetGroupId sets field value
func (o *CollStatsLatencyNamespaceMetrics) SetGroupId(v string) {
	o.GroupId = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *CollStatsLatencyNamespaceMetrics) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CollStatsLatencyNamespaceMetrics) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *CollStatsLatencyNamespaceMetrics) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *CollStatsLatencyNamespaceMetrics) SetLinks(v []Link) {
	o.Links = &v
}

// GetMetrics returns the Metrics field value
func (o *CollStatsLatencyNamespaceMetrics) GetMetrics() []CollStatsLatencyNamespaceMetric {
	if o == nil {
		var ret []CollStatsLatencyNamespaceMetric
		return ret
	}

	return o.Metrics
}

// GetMetricsOk returns a tuple with the Metrics field value
// and a boolean to check if the value has been set.
func (o *CollStatsLatencyNamespaceMetrics) GetMetricsOk() (*[]CollStatsLatencyNamespaceMetric, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metrics, true
}

// SetMetrics sets field value
func (o *CollStatsLatencyNamespaceMetrics) SetMetrics(v []CollStatsLatencyNamespaceMetric) {
	o.Metrics = v
}
