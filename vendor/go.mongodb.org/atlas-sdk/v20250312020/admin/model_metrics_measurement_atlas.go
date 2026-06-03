// Code based on the AtlasAPI V2 OpenAPI file

package admin

// MetricsMeasurementAtlas struct for MetricsMeasurementAtlas
type MetricsMeasurementAtlas struct {
	// List that contains the value of, and metadata provided for, one data point generated at a particular moment in time. If no data point exists for a particular moment in time, the `value` parameter returns `null`.
	// Read only field.
	DataPoints *[]MetricDataPointAtlas `json:"dataPoints,omitempty"`
	// Human-readable label of the measurement that this data point covers.
	// Read only field.
	Name *string `json:"name,omitempty"`
	// Element used to quantify the measurement. The resource returns units of throughput, storage, and time.
	// Read only field.
	Units *string `json:"units,omitempty"`
}

// NewMetricsMeasurementAtlas instantiates a new MetricsMeasurementAtlas object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetricsMeasurementAtlas() *MetricsMeasurementAtlas {
	this := MetricsMeasurementAtlas{}
	return &this
}

// NewMetricsMeasurementAtlasWithDefaults instantiates a new MetricsMeasurementAtlas object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetricsMeasurementAtlasWithDefaults() *MetricsMeasurementAtlas {
	this := MetricsMeasurementAtlas{}
	return &this
}

// GetDataPoints returns the DataPoints field value if set, zero value otherwise
func (o *MetricsMeasurementAtlas) GetDataPoints() []MetricDataPointAtlas {
	if o == nil || IsNil(o.DataPoints) {
		var ret []MetricDataPointAtlas
		return ret
	}
	return *o.DataPoints
}

// GetDataPointsOk returns a tuple with the DataPoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetricsMeasurementAtlas) GetDataPointsOk() (*[]MetricDataPointAtlas, bool) {
	if o == nil || IsNil(o.DataPoints) {
		return nil, false
	}

	return o.DataPoints, true
}

// HasDataPoints returns a boolean if a field has been set.
func (o *MetricsMeasurementAtlas) HasDataPoints() bool {
	if o != nil && !IsNil(o.DataPoints) {
		return true
	}

	return false
}

// SetDataPoints gets a reference to the given []MetricDataPointAtlas and assigns it to the DataPoints field.
func (o *MetricsMeasurementAtlas) SetDataPoints(v []MetricDataPointAtlas) {
	o.DataPoints = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *MetricsMeasurementAtlas) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetricsMeasurementAtlas) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *MetricsMeasurementAtlas) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *MetricsMeasurementAtlas) SetName(v string) {
	o.Name = &v
}

// GetUnits returns the Units field value if set, zero value otherwise
func (o *MetricsMeasurementAtlas) GetUnits() string {
	if o == nil || IsNil(o.Units) {
		var ret string
		return ret
	}
	return *o.Units
}

// GetUnitsOk returns a tuple with the Units field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetricsMeasurementAtlas) GetUnitsOk() (*string, bool) {
	if o == nil || IsNil(o.Units) {
		return nil, false
	}

	return o.Units, true
}

// HasUnits returns a boolean if a field has been set.
func (o *MetricsMeasurementAtlas) HasUnits() bool {
	if o != nil && !IsNil(o.Units) {
		return true
	}

	return false
}

// SetUnits gets a reference to the given string and assigns it to the Units field.
func (o *MetricsMeasurementAtlas) SetUnits(v string) {
	o.Units = &v
}
