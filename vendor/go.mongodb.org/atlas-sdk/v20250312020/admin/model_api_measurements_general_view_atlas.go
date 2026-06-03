// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ApiMeasurementsGeneralViewAtlas struct for ApiMeasurementsGeneralViewAtlas
type ApiMeasurementsGeneralViewAtlas struct {
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
	// Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	// Read only field.
	HostId *string `json:"hostId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]LinkAtlas `json:"links,omitempty"`
	// List that contains measurements and their data points.
	// Read only field.
	Measurements *[]MetricsMeasurementAtlas `json:"measurements,omitempty"`
	// Human-readable label of the disk or partition to which the measurements apply.
	// Read only field.
	PartitionName *string `json:"partitionName,omitempty"`
	// Combination of hostname and Internet Assigned Numbers Authority (IANA) port that serves the MongoDB process. The host must be the hostname, fully qualified domain name (FQDN), or Internet Protocol address (IPv4 or IPv6) of the host that runs the MongoDB process (`mongod` or `mongos`). The port must be the IANA port on which the MongoDB process listens for requests.
	// Read only field.
	ProcessId *string `json:"processId,omitempty"`
	// Date and time that specifies when to start retrieving measurements. If you set **start**, you must set **end**. You can't set this parameter and **period** in the same request. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Start *time.Time `json:"start,omitempty"`
}

// NewApiMeasurementsGeneralViewAtlas instantiates a new ApiMeasurementsGeneralViewAtlas object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiMeasurementsGeneralViewAtlas() *ApiMeasurementsGeneralViewAtlas {
	this := ApiMeasurementsGeneralViewAtlas{}
	return &this
}

// NewApiMeasurementsGeneralViewAtlasWithDefaults instantiates a new ApiMeasurementsGeneralViewAtlas object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiMeasurementsGeneralViewAtlasWithDefaults() *ApiMeasurementsGeneralViewAtlas {
	this := ApiMeasurementsGeneralViewAtlas{}
	return &this
}

// GetDatabaseName returns the DatabaseName field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetDatabaseName() string {
	if o == nil || IsNil(o.DatabaseName) {
		var ret string
		return ret
	}
	return *o.DatabaseName
}

// GetDatabaseNameOk returns a tuple with the DatabaseName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetDatabaseNameOk() (*string, bool) {
	if o == nil || IsNil(o.DatabaseName) {
		return nil, false
	}

	return o.DatabaseName, true
}

// HasDatabaseName returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasDatabaseName() bool {
	if o != nil && !IsNil(o.DatabaseName) {
		return true
	}

	return false
}

// SetDatabaseName gets a reference to the given string and assigns it to the DatabaseName field.
func (o *ApiMeasurementsGeneralViewAtlas) SetDatabaseName(v string) {
	o.DatabaseName = &v
}

// GetEnd returns the End field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetEnd() time.Time {
	if o == nil || IsNil(o.End) {
		var ret time.Time
		return ret
	}
	return *o.End
}

// GetEndOk returns a tuple with the End field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetEndOk() (*time.Time, bool) {
	if o == nil || IsNil(o.End) {
		return nil, false
	}

	return o.End, true
}

// HasEnd returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasEnd() bool {
	if o != nil && !IsNil(o.End) {
		return true
	}

	return false
}

// SetEnd gets a reference to the given time.Time and assigns it to the End field.
func (o *ApiMeasurementsGeneralViewAtlas) SetEnd(v time.Time) {
	o.End = &v
}

// GetGranularity returns the Granularity field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetGranularity() string {
	if o == nil || IsNil(o.Granularity) {
		var ret string
		return ret
	}
	return *o.Granularity
}

// GetGranularityOk returns a tuple with the Granularity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetGranularityOk() (*string, bool) {
	if o == nil || IsNil(o.Granularity) {
		return nil, false
	}

	return o.Granularity, true
}

// HasGranularity returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasGranularity() bool {
	if o != nil && !IsNil(o.Granularity) {
		return true
	}

	return false
}

// SetGranularity gets a reference to the given string and assigns it to the Granularity field.
func (o *ApiMeasurementsGeneralViewAtlas) SetGranularity(v string) {
	o.Granularity = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ApiMeasurementsGeneralViewAtlas) SetGroupId(v string) {
	o.GroupId = &v
}

// GetHostId returns the HostId field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetHostId() string {
	if o == nil || IsNil(o.HostId) {
		var ret string
		return ret
	}
	return *o.HostId
}

// GetHostIdOk returns a tuple with the HostId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetHostIdOk() (*string, bool) {
	if o == nil || IsNil(o.HostId) {
		return nil, false
	}

	return o.HostId, true
}

// HasHostId returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasHostId() bool {
	if o != nil && !IsNil(o.HostId) {
		return true
	}

	return false
}

// SetHostId gets a reference to the given string and assigns it to the HostId field.
func (o *ApiMeasurementsGeneralViewAtlas) SetHostId(v string) {
	o.HostId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetLinks() []LinkAtlas {
	if o == nil || IsNil(o.Links) {
		var ret []LinkAtlas
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetLinksOk() (*[]LinkAtlas, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []LinkAtlas and assigns it to the Links field.
func (o *ApiMeasurementsGeneralViewAtlas) SetLinks(v []LinkAtlas) {
	o.Links = &v
}

// GetMeasurements returns the Measurements field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetMeasurements() []MetricsMeasurementAtlas {
	if o == nil || IsNil(o.Measurements) {
		var ret []MetricsMeasurementAtlas
		return ret
	}
	return *o.Measurements
}

// GetMeasurementsOk returns a tuple with the Measurements field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetMeasurementsOk() (*[]MetricsMeasurementAtlas, bool) {
	if o == nil || IsNil(o.Measurements) {
		return nil, false
	}

	return o.Measurements, true
}

// HasMeasurements returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasMeasurements() bool {
	if o != nil && !IsNil(o.Measurements) {
		return true
	}

	return false
}

// SetMeasurements gets a reference to the given []MetricsMeasurementAtlas and assigns it to the Measurements field.
func (o *ApiMeasurementsGeneralViewAtlas) SetMeasurements(v []MetricsMeasurementAtlas) {
	o.Measurements = &v
}

// GetPartitionName returns the PartitionName field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetPartitionName() string {
	if o == nil || IsNil(o.PartitionName) {
		var ret string
		return ret
	}
	return *o.PartitionName
}

// GetPartitionNameOk returns a tuple with the PartitionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetPartitionNameOk() (*string, bool) {
	if o == nil || IsNil(o.PartitionName) {
		return nil, false
	}

	return o.PartitionName, true
}

// HasPartitionName returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasPartitionName() bool {
	if o != nil && !IsNil(o.PartitionName) {
		return true
	}

	return false
}

// SetPartitionName gets a reference to the given string and assigns it to the PartitionName field.
func (o *ApiMeasurementsGeneralViewAtlas) SetPartitionName(v string) {
	o.PartitionName = &v
}

// GetProcessId returns the ProcessId field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetProcessId() string {
	if o == nil || IsNil(o.ProcessId) {
		var ret string
		return ret
	}
	return *o.ProcessId
}

// GetProcessIdOk returns a tuple with the ProcessId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetProcessIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessId) {
		return nil, false
	}

	return o.ProcessId, true
}

// HasProcessId returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasProcessId() bool {
	if o != nil && !IsNil(o.ProcessId) {
		return true
	}

	return false
}

// SetProcessId gets a reference to the given string and assigns it to the ProcessId field.
func (o *ApiMeasurementsGeneralViewAtlas) SetProcessId(v string) {
	o.ProcessId = &v
}

// GetStart returns the Start field value if set, zero value otherwise
func (o *ApiMeasurementsGeneralViewAtlas) GetStart() time.Time {
	if o == nil || IsNil(o.Start) {
		var ret time.Time
		return ret
	}
	return *o.Start
}

// GetStartOk returns a tuple with the Start field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiMeasurementsGeneralViewAtlas) GetStartOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Start) {
		return nil, false
	}

	return o.Start, true
}

// HasStart returns a boolean if a field has been set.
func (o *ApiMeasurementsGeneralViewAtlas) HasStart() bool {
	if o != nil && !IsNil(o.Start) {
		return true
	}

	return false
}

// SetStart gets a reference to the given time.Time and assigns it to the Start field.
func (o *ApiMeasurementsGeneralViewAtlas) SetStart(v time.Time) {
	o.Start = &v
}
