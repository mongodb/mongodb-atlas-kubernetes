// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// MeasurementsNonIndex struct for MeasurementsNonIndex
type MeasurementsNonIndex struct {
	// Date and time that specifies when to stop retrieving measurements. If you set **end**, you must set **start**. You can't set this parameter and **period** in the same request. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	End *time.Time `json:"end,omitempty"`
	// Duration that specifies the interval between measurement data points. The parameter expresses its value in ISO 8601 timestamp format in UTC. If you set this parameter, you must set either **period** or **start** and **end**.
	// Read only field.
	Granularity *string `json:"granularity,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project. The project contains MongoDB processes that you want to return. The MongoDB process can be either the `mongod` or `mongos`.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List that contains the Atlas Search hardware measurements.
	// Read only field.
	HardwareMeasurements *[]MetricsMeasurement `json:"hardwareMeasurements,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	// Read only field.
	ProcessId *string `json:"processId,omitempty"`
	// Date and time that specifies when to start retrieving measurements. If you set **start**, you must set **end**. You can't set this parameter and **period** in the same request. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Start *time.Time `json:"start,omitempty"`
	// List that contains the Atlas Search status measurements.
	// Read only field.
	StatusMeasurements *[]MetricsMeasurement `json:"statusMeasurements,omitempty"`
}

// NewMeasurementsNonIndex instantiates a new MeasurementsNonIndex object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMeasurementsNonIndex() *MeasurementsNonIndex {
	this := MeasurementsNonIndex{}
	return &this
}

// NewMeasurementsNonIndexWithDefaults instantiates a new MeasurementsNonIndex object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMeasurementsNonIndexWithDefaults() *MeasurementsNonIndex {
	this := MeasurementsNonIndex{}
	return &this
}

// GetEnd returns the End field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetEnd() time.Time {
	if o == nil || IsNil(o.End) {
		var ret time.Time
		return ret
	}
	return *o.End
}

// GetEndOk returns a tuple with the End field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetEndOk() (*time.Time, bool) {
	if o == nil || IsNil(o.End) {
		return nil, false
	}

	return o.End, true
}

// HasEnd returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasEnd() bool {
	if o != nil && !IsNil(o.End) {
		return true
	}

	return false
}

// SetEnd gets a reference to the given time.Time and assigns it to the End field.
func (o *MeasurementsNonIndex) SetEnd(v time.Time) {
	o.End = &v
}

// GetGranularity returns the Granularity field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetGranularity() string {
	if o == nil || IsNil(o.Granularity) {
		var ret string
		return ret
	}
	return *o.Granularity
}

// GetGranularityOk returns a tuple with the Granularity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetGranularityOk() (*string, bool) {
	if o == nil || IsNil(o.Granularity) {
		return nil, false
	}

	return o.Granularity, true
}

// HasGranularity returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasGranularity() bool {
	if o != nil && !IsNil(o.Granularity) {
		return true
	}

	return false
}

// SetGranularity gets a reference to the given string and assigns it to the Granularity field.
func (o *MeasurementsNonIndex) SetGranularity(v string) {
	o.Granularity = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *MeasurementsNonIndex) SetGroupId(v string) {
	o.GroupId = &v
}

// GetHardwareMeasurements returns the HardwareMeasurements field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetHardwareMeasurements() []MetricsMeasurement {
	if o == nil || IsNil(o.HardwareMeasurements) {
		var ret []MetricsMeasurement
		return ret
	}
	return *o.HardwareMeasurements
}

// GetHardwareMeasurementsOk returns a tuple with the HardwareMeasurements field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetHardwareMeasurementsOk() (*[]MetricsMeasurement, bool) {
	if o == nil || IsNil(o.HardwareMeasurements) {
		return nil, false
	}

	return o.HardwareMeasurements, true
}

// HasHardwareMeasurements returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasHardwareMeasurements() bool {
	if o != nil && !IsNil(o.HardwareMeasurements) {
		return true
	}

	return false
}

// SetHardwareMeasurements gets a reference to the given []MetricsMeasurement and assigns it to the HardwareMeasurements field.
func (o *MeasurementsNonIndex) SetHardwareMeasurements(v []MetricsMeasurement) {
	o.HardwareMeasurements = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *MeasurementsNonIndex) SetLinks(v []Link) {
	o.Links = &v
}

// GetProcessId returns the ProcessId field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetProcessId() string {
	if o == nil || IsNil(o.ProcessId) {
		var ret string
		return ret
	}
	return *o.ProcessId
}

// GetProcessIdOk returns a tuple with the ProcessId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetProcessIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessId) {
		return nil, false
	}

	return o.ProcessId, true
}

// HasProcessId returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasProcessId() bool {
	if o != nil && !IsNil(o.ProcessId) {
		return true
	}

	return false
}

// SetProcessId gets a reference to the given string and assigns it to the ProcessId field.
func (o *MeasurementsNonIndex) SetProcessId(v string) {
	o.ProcessId = &v
}

// GetStart returns the Start field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetStart() time.Time {
	if o == nil || IsNil(o.Start) {
		var ret time.Time
		return ret
	}
	return *o.Start
}

// GetStartOk returns a tuple with the Start field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetStartOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Start) {
		return nil, false
	}

	return o.Start, true
}

// HasStart returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasStart() bool {
	if o != nil && !IsNil(o.Start) {
		return true
	}

	return false
}

// SetStart gets a reference to the given time.Time and assigns it to the Start field.
func (o *MeasurementsNonIndex) SetStart(v time.Time) {
	o.Start = &v
}

// GetStatusMeasurements returns the StatusMeasurements field value if set, zero value otherwise
func (o *MeasurementsNonIndex) GetStatusMeasurements() []MetricsMeasurement {
	if o == nil || IsNil(o.StatusMeasurements) {
		var ret []MetricsMeasurement
		return ret
	}
	return *o.StatusMeasurements
}

// GetStatusMeasurementsOk returns a tuple with the StatusMeasurements field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsNonIndex) GetStatusMeasurementsOk() (*[]MetricsMeasurement, bool) {
	if o == nil || IsNil(o.StatusMeasurements) {
		return nil, false
	}

	return o.StatusMeasurements, true
}

// HasStatusMeasurements returns a boolean if a field has been set.
func (o *MeasurementsNonIndex) HasStatusMeasurements() bool {
	if o != nil && !IsNil(o.StatusMeasurements) {
		return true
	}

	return false
}

// SetStatusMeasurements gets a reference to the given []MetricsMeasurement and assigns it to the StatusMeasurements field.
func (o *MeasurementsNonIndex) SetStatusMeasurements(v []MetricsMeasurement) {
	o.StatusMeasurements = &v
}
