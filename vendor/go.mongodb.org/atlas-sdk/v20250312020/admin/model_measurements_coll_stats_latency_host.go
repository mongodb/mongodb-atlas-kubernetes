// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// MeasurementsCollStatsLatencyHost struct for MeasurementsCollStatsLatencyHost
type MeasurementsCollStatsLatencyHost struct {
	// Human-readable label that identifies the collection.
	// Read only field.
	CollectionName *string `json:"collectionName,omitempty"`
	// Human-readable label that identifies the database that the specified MongoDB process serves.
	// Read only field.
	DatabaseName *string `json:"databaseName,omitempty"`
	// Date and time that specifies when to stop retrieving measurements. If you set **end**, you must set **start**. You can't set this parameter and **period** in the same request. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	End *time.Time `json:"end,omitempty"`
	// Duration that specifies the interval between measurement data points. The parameter expresses its value in ISO 8601 timestamp format in UTC. If you set this parameter, you must set either **period** or **start** and **end**.
	// Read only field.
	Granularity *string `json:"granularity,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project. The project contains MongoDB processes that you want to return. The MongoDB process can be either the `mongod` or `mongos`.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List that contains measurements and their data points.
	// Read only field.
	Measurements *[]MetricsMeasurement `json:"measurements,omitempty"`
	// Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	// Read only field.
	ProcessId *string `json:"processId,omitempty"`
	// Date and time that specifies when to start retrieving measurements. If you set **start**, you must set **end**. You can't set this parameter and **period** in the same request. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Start *time.Time `json:"start,omitempty"`
}

// NewMeasurementsCollStatsLatencyHost instantiates a new MeasurementsCollStatsLatencyHost object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMeasurementsCollStatsLatencyHost() *MeasurementsCollStatsLatencyHost {
	this := MeasurementsCollStatsLatencyHost{}
	return &this
}

// NewMeasurementsCollStatsLatencyHostWithDefaults instantiates a new MeasurementsCollStatsLatencyHost object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMeasurementsCollStatsLatencyHostWithDefaults() *MeasurementsCollStatsLatencyHost {
	this := MeasurementsCollStatsLatencyHost{}
	return &this
}

// GetCollectionName returns the CollectionName field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetCollectionName() string {
	if o == nil || IsNil(o.CollectionName) {
		var ret string
		return ret
	}
	return *o.CollectionName
}

// GetCollectionNameOk returns a tuple with the CollectionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetCollectionNameOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionName) {
		return nil, false
	}

	return o.CollectionName, true
}

// HasCollectionName returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasCollectionName() bool {
	if o != nil && !IsNil(o.CollectionName) {
		return true
	}

	return false
}

// SetCollectionName gets a reference to the given string and assigns it to the CollectionName field.
func (o *MeasurementsCollStatsLatencyHost) SetCollectionName(v string) {
	o.CollectionName = &v
}

// GetDatabaseName returns the DatabaseName field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetDatabaseName() string {
	if o == nil || IsNil(o.DatabaseName) {
		var ret string
		return ret
	}
	return *o.DatabaseName
}

// GetDatabaseNameOk returns a tuple with the DatabaseName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetDatabaseNameOk() (*string, bool) {
	if o == nil || IsNil(o.DatabaseName) {
		return nil, false
	}

	return o.DatabaseName, true
}

// HasDatabaseName returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasDatabaseName() bool {
	if o != nil && !IsNil(o.DatabaseName) {
		return true
	}

	return false
}

// SetDatabaseName gets a reference to the given string and assigns it to the DatabaseName field.
func (o *MeasurementsCollStatsLatencyHost) SetDatabaseName(v string) {
	o.DatabaseName = &v
}

// GetEnd returns the End field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetEnd() time.Time {
	if o == nil || IsNil(o.End) {
		var ret time.Time
		return ret
	}
	return *o.End
}

// GetEndOk returns a tuple with the End field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetEndOk() (*time.Time, bool) {
	if o == nil || IsNil(o.End) {
		return nil, false
	}

	return o.End, true
}

// HasEnd returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasEnd() bool {
	if o != nil && !IsNil(o.End) {
		return true
	}

	return false
}

// SetEnd gets a reference to the given time.Time and assigns it to the End field.
func (o *MeasurementsCollStatsLatencyHost) SetEnd(v time.Time) {
	o.End = &v
}

// GetGranularity returns the Granularity field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetGranularity() string {
	if o == nil || IsNil(o.Granularity) {
		var ret string
		return ret
	}
	return *o.Granularity
}

// GetGranularityOk returns a tuple with the Granularity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetGranularityOk() (*string, bool) {
	if o == nil || IsNil(o.Granularity) {
		return nil, false
	}

	return o.Granularity, true
}

// HasGranularity returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasGranularity() bool {
	if o != nil && !IsNil(o.Granularity) {
		return true
	}

	return false
}

// SetGranularity gets a reference to the given string and assigns it to the Granularity field.
func (o *MeasurementsCollStatsLatencyHost) SetGranularity(v string) {
	o.Granularity = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *MeasurementsCollStatsLatencyHost) SetGroupId(v string) {
	o.GroupId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *MeasurementsCollStatsLatencyHost) SetLinks(v []Link) {
	o.Links = &v
}

// GetMeasurements returns the Measurements field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetMeasurements() []MetricsMeasurement {
	if o == nil || IsNil(o.Measurements) {
		var ret []MetricsMeasurement
		return ret
	}
	return *o.Measurements
}

// GetMeasurementsOk returns a tuple with the Measurements field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetMeasurementsOk() (*[]MetricsMeasurement, bool) {
	if o == nil || IsNil(o.Measurements) {
		return nil, false
	}

	return o.Measurements, true
}

// HasMeasurements returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasMeasurements() bool {
	if o != nil && !IsNil(o.Measurements) {
		return true
	}

	return false
}

// SetMeasurements gets a reference to the given []MetricsMeasurement and assigns it to the Measurements field.
func (o *MeasurementsCollStatsLatencyHost) SetMeasurements(v []MetricsMeasurement) {
	o.Measurements = &v
}

// GetProcessId returns the ProcessId field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetProcessId() string {
	if o == nil || IsNil(o.ProcessId) {
		var ret string
		return ret
	}
	return *o.ProcessId
}

// GetProcessIdOk returns a tuple with the ProcessId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetProcessIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessId) {
		return nil, false
	}

	return o.ProcessId, true
}

// HasProcessId returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasProcessId() bool {
	if o != nil && !IsNil(o.ProcessId) {
		return true
	}

	return false
}

// SetProcessId gets a reference to the given string and assigns it to the ProcessId field.
func (o *MeasurementsCollStatsLatencyHost) SetProcessId(v string) {
	o.ProcessId = &v
}

// GetStart returns the Start field value if set, zero value otherwise
func (o *MeasurementsCollStatsLatencyHost) GetStart() time.Time {
	if o == nil || IsNil(o.Start) {
		var ret time.Time
		return ret
	}
	return *o.Start
}

// GetStartOk returns a tuple with the Start field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MeasurementsCollStatsLatencyHost) GetStartOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Start) {
		return nil, false
	}

	return o.Start, true
}

// HasStart returns a boolean if a field has been set.
func (o *MeasurementsCollStatsLatencyHost) HasStart() bool {
	if o != nil && !IsNil(o.Start) {
		return true
	}

	return false
}

// SetStart gets a reference to the given time.Time and assigns it to the Start field.
func (o *MeasurementsCollStatsLatencyHost) SetStart(v time.Time) {
	o.Start = &v
}
