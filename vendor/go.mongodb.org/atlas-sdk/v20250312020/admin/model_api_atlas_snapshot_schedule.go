// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasSnapshotSchedule struct for ApiAtlasSnapshotSchedule
type ApiAtlasSnapshotSchedule struct {
	// Quantity of time expressed in minutes between successive cluster checkpoints. This parameter applies only to sharded clusters. This number determines the granularity of continuous cloud backups for sharded clusters.
	ClusterCheckpointIntervalMin int `json:"clusterCheckpointIntervalMin"`
	// Unique 24-hexadecimal digit string that identifies the cluster with the snapshot you want to return.
	ClusterId string `json:"clusterId"`
	// Quantity of time to keep daily snapshots. MongoDB Cloud expresses this value in days. Set this value to `0` to disable daily snapshot retention.
	DailySnapshotRetentionDays int `json:"dailySnapshotRetentionDays"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the cluster.
	// Read only field.
	GroupId string `json:"groupId"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Number of months that MongoDB Cloud must keep monthly snapshots. Set this value to `0` to disable monthly snapshot retention.
	MonthlySnapshotRetentionMonths int `json:"monthlySnapshotRetentionMonths"`
	// Number of hours before the current time from which MongoDB Cloud can create a Continuous Cloud Backup snapshot.
	PointInTimeWindowHours int `json:"pointInTimeWindowHours"`
	// Number of hours that must elapse before taking another snapshot.
	SnapshotIntervalHours int `json:"snapshotIntervalHours"`
	// Number of days that MongoDB Cloud must keep recent snapshots.
	SnapshotRetentionDays int `json:"snapshotRetentionDays"`
	// Number of weeks that MongoDB Cloud must keep weekly snapshots. Set this value to `0` to disable weekly snapshot retention.
	WeeklySnapshotRetentionWeeks int `json:"weeklySnapshotRetentionWeeks"`
}

// NewApiAtlasSnapshotSchedule instantiates a new ApiAtlasSnapshotSchedule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasSnapshotSchedule(clusterCheckpointIntervalMin int, clusterId string, dailySnapshotRetentionDays int, groupId string, monthlySnapshotRetentionMonths int, pointInTimeWindowHours int, snapshotIntervalHours int, snapshotRetentionDays int, weeklySnapshotRetentionWeeks int) *ApiAtlasSnapshotSchedule {
	this := ApiAtlasSnapshotSchedule{}
	this.ClusterCheckpointIntervalMin = clusterCheckpointIntervalMin
	this.ClusterId = clusterId
	this.DailySnapshotRetentionDays = dailySnapshotRetentionDays
	this.GroupId = groupId
	this.MonthlySnapshotRetentionMonths = monthlySnapshotRetentionMonths
	this.PointInTimeWindowHours = pointInTimeWindowHours
	this.SnapshotIntervalHours = snapshotIntervalHours
	this.SnapshotRetentionDays = snapshotRetentionDays
	this.WeeklySnapshotRetentionWeeks = weeklySnapshotRetentionWeeks
	return &this
}

// NewApiAtlasSnapshotScheduleWithDefaults instantiates a new ApiAtlasSnapshotSchedule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasSnapshotScheduleWithDefaults() *ApiAtlasSnapshotSchedule {
	this := ApiAtlasSnapshotSchedule{}
	return &this
}

// GetClusterCheckpointIntervalMin returns the ClusterCheckpointIntervalMin field value
func (o *ApiAtlasSnapshotSchedule) GetClusterCheckpointIntervalMin() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.ClusterCheckpointIntervalMin
}

// GetClusterCheckpointIntervalMinOk returns a tuple with the ClusterCheckpointIntervalMin field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetClusterCheckpointIntervalMinOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterCheckpointIntervalMin, true
}

// SetClusterCheckpointIntervalMin sets field value
func (o *ApiAtlasSnapshotSchedule) SetClusterCheckpointIntervalMin(v int) {
	o.ClusterCheckpointIntervalMin = v
}

// GetClusterId returns the ClusterId field value
func (o *ApiAtlasSnapshotSchedule) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *ApiAtlasSnapshotSchedule) SetClusterId(v string) {
	o.ClusterId = v
}

// GetDailySnapshotRetentionDays returns the DailySnapshotRetentionDays field value
func (o *ApiAtlasSnapshotSchedule) GetDailySnapshotRetentionDays() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.DailySnapshotRetentionDays
}

// GetDailySnapshotRetentionDaysOk returns a tuple with the DailySnapshotRetentionDays field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetDailySnapshotRetentionDaysOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DailySnapshotRetentionDays, true
}

// SetDailySnapshotRetentionDays sets field value
func (o *ApiAtlasSnapshotSchedule) SetDailySnapshotRetentionDays(v int) {
	o.DailySnapshotRetentionDays = v
}

// GetGroupId returns the GroupId field value
func (o *ApiAtlasSnapshotSchedule) GetGroupId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GroupId, true
}

// SetGroupId sets field value
func (o *ApiAtlasSnapshotSchedule) SetGroupId(v string) {
	o.GroupId = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ApiAtlasSnapshotSchedule) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ApiAtlasSnapshotSchedule) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ApiAtlasSnapshotSchedule) SetLinks(v []Link) {
	o.Links = &v
}

// GetMonthlySnapshotRetentionMonths returns the MonthlySnapshotRetentionMonths field value
func (o *ApiAtlasSnapshotSchedule) GetMonthlySnapshotRetentionMonths() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.MonthlySnapshotRetentionMonths
}

// GetMonthlySnapshotRetentionMonthsOk returns a tuple with the MonthlySnapshotRetentionMonths field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetMonthlySnapshotRetentionMonthsOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MonthlySnapshotRetentionMonths, true
}

// SetMonthlySnapshotRetentionMonths sets field value
func (o *ApiAtlasSnapshotSchedule) SetMonthlySnapshotRetentionMonths(v int) {
	o.MonthlySnapshotRetentionMonths = v
}

// GetPointInTimeWindowHours returns the PointInTimeWindowHours field value
func (o *ApiAtlasSnapshotSchedule) GetPointInTimeWindowHours() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.PointInTimeWindowHours
}

// GetPointInTimeWindowHoursOk returns a tuple with the PointInTimeWindowHours field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetPointInTimeWindowHoursOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PointInTimeWindowHours, true
}

// SetPointInTimeWindowHours sets field value
func (o *ApiAtlasSnapshotSchedule) SetPointInTimeWindowHours(v int) {
	o.PointInTimeWindowHours = v
}

// GetSnapshotIntervalHours returns the SnapshotIntervalHours field value
func (o *ApiAtlasSnapshotSchedule) GetSnapshotIntervalHours() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.SnapshotIntervalHours
}

// GetSnapshotIntervalHoursOk returns a tuple with the SnapshotIntervalHours field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetSnapshotIntervalHoursOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SnapshotIntervalHours, true
}

// SetSnapshotIntervalHours sets field value
func (o *ApiAtlasSnapshotSchedule) SetSnapshotIntervalHours(v int) {
	o.SnapshotIntervalHours = v
}

// GetSnapshotRetentionDays returns the SnapshotRetentionDays field value
func (o *ApiAtlasSnapshotSchedule) GetSnapshotRetentionDays() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.SnapshotRetentionDays
}

// GetSnapshotRetentionDaysOk returns a tuple with the SnapshotRetentionDays field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetSnapshotRetentionDaysOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SnapshotRetentionDays, true
}

// SetSnapshotRetentionDays sets field value
func (o *ApiAtlasSnapshotSchedule) SetSnapshotRetentionDays(v int) {
	o.SnapshotRetentionDays = v
}

// GetWeeklySnapshotRetentionWeeks returns the WeeklySnapshotRetentionWeeks field value
func (o *ApiAtlasSnapshotSchedule) GetWeeklySnapshotRetentionWeeks() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.WeeklySnapshotRetentionWeeks
}

// GetWeeklySnapshotRetentionWeeksOk returns a tuple with the WeeklySnapshotRetentionWeeks field value
// and a boolean to check if the value has been set.
func (o *ApiAtlasSnapshotSchedule) GetWeeklySnapshotRetentionWeeksOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WeeklySnapshotRetentionWeeks, true
}

// SetWeeklySnapshotRetentionWeeks sets field value
func (o *ApiAtlasSnapshotSchedule) SetWeeklySnapshotRetentionWeeks(v int) {
	o.WeeklySnapshotRetentionWeeks = v
}
