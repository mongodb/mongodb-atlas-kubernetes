// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CollStatsLatencyNamespaceMetric Coll Stats Latency metric name and its unit of measurement.
type CollStatsLatencyNamespaceMetric struct {
	// Human-readable label that identifies this metric.
	// Read only field.
	MetricName string `json:"metricName"`
	// Unit of measurement that applies to this metric.
	// Read only field.
	Units string `json:"units"`
}

// NewCollStatsLatencyNamespaceMetric instantiates a new CollStatsLatencyNamespaceMetric object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCollStatsLatencyNamespaceMetric(metricName string, units string) *CollStatsLatencyNamespaceMetric {
	this := CollStatsLatencyNamespaceMetric{}
	this.MetricName = metricName
	this.Units = units
	return &this
}

// NewCollStatsLatencyNamespaceMetricWithDefaults instantiates a new CollStatsLatencyNamespaceMetric object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCollStatsLatencyNamespaceMetricWithDefaults() *CollStatsLatencyNamespaceMetric {
	this := CollStatsLatencyNamespaceMetric{}
	return &this
}

// GetMetricName returns the MetricName field value
func (o *CollStatsLatencyNamespaceMetric) GetMetricName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MetricName
}

// GetMetricNameOk returns a tuple with the MetricName field value
// and a boolean to check if the value has been set.
func (o *CollStatsLatencyNamespaceMetric) GetMetricNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MetricName, true
}

// SetMetricName sets field value
func (o *CollStatsLatencyNamespaceMetric) SetMetricName(v string) {
	o.MetricName = v
}

// GetUnits returns the Units field value
func (o *CollStatsLatencyNamespaceMetric) GetUnits() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Units
}

// GetUnitsOk returns a tuple with the Units field value
// and a boolean to check if the value has been set.
func (o *CollStatsLatencyNamespaceMetric) GetUnitsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Units, true
}

// SetUnits sets field value
func (o *CollStatsLatencyNamespaceMetric) SetUnits(v string) {
	o.Units = v
}
