// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FTSMetric Measurement of one Atlas Search status when MongoDB Atlas received this request.
type FTSMetric struct {
	// Human-readable label that identifies this Atlas Search hardware, status, or index measurement.
	// Read only field.
	MetricName string `json:"metricName"`
	// Unit of measurement that applies to this Atlas Search metric.
	// Read only field.
	Units string `json:"units"`
}

// NewFTSMetric instantiates a new FTSMetric object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFTSMetric(metricName string, units string) *FTSMetric {
	this := FTSMetric{}
	this.MetricName = metricName
	this.Units = units
	return &this
}

// NewFTSMetricWithDefaults instantiates a new FTSMetric object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFTSMetricWithDefaults() *FTSMetric {
	this := FTSMetric{}
	return &this
}

// GetMetricName returns the MetricName field value
func (o *FTSMetric) GetMetricName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MetricName
}

// GetMetricNameOk returns a tuple with the MetricName field value
// and a boolean to check if the value has been set.
func (o *FTSMetric) GetMetricNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MetricName, true
}

// SetMetricName sets field value
func (o *FTSMetric) SetMetricName(v string) {
	o.MetricName = v
}

// GetUnits returns the Units field value
func (o *FTSMetric) GetUnits() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Units
}

// GetUnitsOk returns a tuple with the Units field value
// and a boolean to check if the value has been set.
func (o *FTSMetric) GetUnitsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Units, true
}

// SetUnits sets field value
func (o *FTSMetric) SetUnits(v string) {
	o.Units = v
}
