// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BackupOnlineArchive struct for BackupOnlineArchive
type BackupOnlineArchive struct {
	// Unique 24-hexadecimal digit string that identifies the online archive.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Human-readable label that identifies the cluster that contains the collection for which you want to create an online archive.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Human-readable label that identifies the collection for which you created the online archive.
	// Read only field.
	CollName *string `json:"collName,omitempty"`
	// Classification of MongoDB database collection that you want to return.  If you set this parameter to `TIMESERIES`, set `\"criteria.type\" : \"date\"` and `\"criteria.dateFormat\" : \"ISODATE\"`.
	// Read only field.
	CollectionType     *string             `json:"collectionType,omitempty"`
	Criteria           *Criteria           `json:"criteria,omitempty"`
	DataExpirationRule *DataExpirationRule `json:"dataExpirationRule,omitempty"`
	DataProcessRegion  *DataProcessRegion  `json:"dataProcessRegion,omitempty"`
	// Human-readable label that identifies the dataset that Atlas generates for this online archive.
	// Read only field.
	DataSetName *string `json:"dataSetName,omitempty"`
	// Human-readable label of the database that contains the collection that contains the online archive.
	// Read only field.
	DbName *string `json:"dbName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the specified cluster. The specified cluster contains the collection for which to create the online archive.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List that contains document parameters to use to logically divide data within a collection. Partitions provide a coarse level of filtering of the underlying collection data. To divide your data, specify parameters that you frequently query. If you specified `criteria.type`: `DATE` in the Create One Online Archive endpoint, then you can specify up to three parameters by which to query. One of these parameters must be the `DATE` value, which is required in this case. If you specified `criteria.type`: `CUSTOM` in the Create One Online Archive endpoint, then you can specify up to two parameters by which to query. Queries that don't use `criteria.type`: `DATE` or `criteria.type`: `CUSTOM` parameters cause MongoDB to scan a full collection of all archived documents. This takes more time and increases your costs.
	// Read only field.
	PartitionFields *[]PartitionField `json:"partitionFields,omitempty"`
	// Flag that indicates whether this online archive exists in the paused state. A request to resume fails if the collection has another active online archive. To pause an active online archive or resume a paused online archive, you must include this parameter. To pause an active archive, set this to **true**. To resume a paused archive, set this to **false**.
	Paused   *bool                  `json:"paused,omitempty"`
	Schedule *OnlineArchiveSchedule `json:"schedule,omitempty"`
	// Phase of the process to create this online archive when you made this request.  | State       | Indication | |-------------|------------| | `PENDING`   | MongoDB Cloud has queued documents for archive. Archiving hasn't started. | | `ARCHIVING` | MongoDB Cloud started archiving documents that meet the archival criteria. | | `IDLE`      | MongoDB Cloud waits to start the next archival job. | | `PAUSING`   | Someone chose to stop archiving. MongoDB Cloud finishes the running archival job then changes the state to `PAUSED` when that job completes. | | `PAUSED`    | MongoDB Cloud has stopped archiving. Archived documents can be queried. The specified archiving operation on the active cluster cannot archive additional documents. You can resume archiving for paused archives at any time. | | `ORPHANED`  | Someone has deleted the collection associated with an active or paused archive. MongoDB Cloud doesn't delete the archived data. You must manually delete the online archives associated with the deleted collection. | | `DELETED`   | Someone has deleted the archive was deleted. When someone deletes an online archive, MongoDB Cloud removes all associated archived documents from the cloud object storage. |
	// Read only field.
	State *string `json:"state,omitempty"`
}

// NewBackupOnlineArchive instantiates a new BackupOnlineArchive object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupOnlineArchive() *BackupOnlineArchive {
	this := BackupOnlineArchive{}
	return &this
}

// NewBackupOnlineArchiveWithDefaults instantiates a new BackupOnlineArchive object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupOnlineArchiveWithDefaults() *BackupOnlineArchive {
	this := BackupOnlineArchive{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BackupOnlineArchive) SetId(v string) {
	o.Id = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *BackupOnlineArchive) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCollName returns the CollName field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetCollName() string {
	if o == nil || IsNil(o.CollName) {
		var ret string
		return ret
	}
	return *o.CollName
}

// GetCollNameOk returns a tuple with the CollName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetCollNameOk() (*string, bool) {
	if o == nil || IsNil(o.CollName) {
		return nil, false
	}

	return o.CollName, true
}

// HasCollName returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasCollName() bool {
	if o != nil && !IsNil(o.CollName) {
		return true
	}

	return false
}

// SetCollName gets a reference to the given string and assigns it to the CollName field.
func (o *BackupOnlineArchive) SetCollName(v string) {
	o.CollName = &v
}

// GetCollectionType returns the CollectionType field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetCollectionType() string {
	if o == nil || IsNil(o.CollectionType) {
		var ret string
		return ret
	}
	return *o.CollectionType
}

// GetCollectionTypeOk returns a tuple with the CollectionType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetCollectionTypeOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionType) {
		return nil, false
	}

	return o.CollectionType, true
}

// HasCollectionType returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasCollectionType() bool {
	if o != nil && !IsNil(o.CollectionType) {
		return true
	}

	return false
}

// SetCollectionType gets a reference to the given string and assigns it to the CollectionType field.
func (o *BackupOnlineArchive) SetCollectionType(v string) {
	o.CollectionType = &v
}

// GetCriteria returns the Criteria field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetCriteria() Criteria {
	if o == nil || IsNil(o.Criteria) {
		var ret Criteria
		return ret
	}
	return *o.Criteria
}

// GetCriteriaOk returns a tuple with the Criteria field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetCriteriaOk() (*Criteria, bool) {
	if o == nil || IsNil(o.Criteria) {
		return nil, false
	}

	return o.Criteria, true
}

// HasCriteria returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasCriteria() bool {
	if o != nil && !IsNil(o.Criteria) {
		return true
	}

	return false
}

// SetCriteria gets a reference to the given Criteria and assigns it to the Criteria field.
func (o *BackupOnlineArchive) SetCriteria(v Criteria) {
	o.Criteria = &v
}

// GetDataExpirationRule returns the DataExpirationRule field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetDataExpirationRule() DataExpirationRule {
	if o == nil || IsNil(o.DataExpirationRule) {
		var ret DataExpirationRule
		return ret
	}
	return *o.DataExpirationRule
}

// GetDataExpirationRuleOk returns a tuple with the DataExpirationRule field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetDataExpirationRuleOk() (*DataExpirationRule, bool) {
	if o == nil || IsNil(o.DataExpirationRule) {
		return nil, false
	}

	return o.DataExpirationRule, true
}

// HasDataExpirationRule returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasDataExpirationRule() bool {
	if o != nil && !IsNil(o.DataExpirationRule) {
		return true
	}

	return false
}

// SetDataExpirationRule gets a reference to the given DataExpirationRule and assigns it to the DataExpirationRule field.
func (o *BackupOnlineArchive) SetDataExpirationRule(v DataExpirationRule) {
	o.DataExpirationRule = &v
}

// GetDataProcessRegion returns the DataProcessRegion field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetDataProcessRegion() DataProcessRegion {
	if o == nil || IsNil(o.DataProcessRegion) {
		var ret DataProcessRegion
		return ret
	}
	return *o.DataProcessRegion
}

// GetDataProcessRegionOk returns a tuple with the DataProcessRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetDataProcessRegionOk() (*DataProcessRegion, bool) {
	if o == nil || IsNil(o.DataProcessRegion) {
		return nil, false
	}

	return o.DataProcessRegion, true
}

// HasDataProcessRegion returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasDataProcessRegion() bool {
	if o != nil && !IsNil(o.DataProcessRegion) {
		return true
	}

	return false
}

// SetDataProcessRegion gets a reference to the given DataProcessRegion and assigns it to the DataProcessRegion field.
func (o *BackupOnlineArchive) SetDataProcessRegion(v DataProcessRegion) {
	o.DataProcessRegion = &v
}

// GetDataSetName returns the DataSetName field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetDataSetName() string {
	if o == nil || IsNil(o.DataSetName) {
		var ret string
		return ret
	}
	return *o.DataSetName
}

// GetDataSetNameOk returns a tuple with the DataSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetDataSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.DataSetName) {
		return nil, false
	}

	return o.DataSetName, true
}

// HasDataSetName returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasDataSetName() bool {
	if o != nil && !IsNil(o.DataSetName) {
		return true
	}

	return false
}

// SetDataSetName gets a reference to the given string and assigns it to the DataSetName field.
func (o *BackupOnlineArchive) SetDataSetName(v string) {
	o.DataSetName = &v
}

// GetDbName returns the DbName field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetDbName() string {
	if o == nil || IsNil(o.DbName) {
		var ret string
		return ret
	}
	return *o.DbName
}

// GetDbNameOk returns a tuple with the DbName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetDbNameOk() (*string, bool) {
	if o == nil || IsNil(o.DbName) {
		return nil, false
	}

	return o.DbName, true
}

// HasDbName returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasDbName() bool {
	if o != nil && !IsNil(o.DbName) {
		return true
	}

	return false
}

// SetDbName gets a reference to the given string and assigns it to the DbName field.
func (o *BackupOnlineArchive) SetDbName(v string) {
	o.DbName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *BackupOnlineArchive) SetGroupId(v string) {
	o.GroupId = &v
}

// GetPartitionFields returns the PartitionFields field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetPartitionFields() []PartitionField {
	if o == nil || IsNil(o.PartitionFields) {
		var ret []PartitionField
		return ret
	}
	return *o.PartitionFields
}

// GetPartitionFieldsOk returns a tuple with the PartitionFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetPartitionFieldsOk() (*[]PartitionField, bool) {
	if o == nil || IsNil(o.PartitionFields) {
		return nil, false
	}

	return o.PartitionFields, true
}

// HasPartitionFields returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasPartitionFields() bool {
	if o != nil && !IsNil(o.PartitionFields) {
		return true
	}

	return false
}

// SetPartitionFields gets a reference to the given []PartitionField and assigns it to the PartitionFields field.
func (o *BackupOnlineArchive) SetPartitionFields(v []PartitionField) {
	o.PartitionFields = &v
}

// GetPaused returns the Paused field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetPaused() bool {
	if o == nil || IsNil(o.Paused) {
		var ret bool
		return ret
	}
	return *o.Paused
}

// GetPausedOk returns a tuple with the Paused field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetPausedOk() (*bool, bool) {
	if o == nil || IsNil(o.Paused) {
		return nil, false
	}

	return o.Paused, true
}

// HasPaused returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasPaused() bool {
	if o != nil && !IsNil(o.Paused) {
		return true
	}

	return false
}

// SetPaused gets a reference to the given bool and assigns it to the Paused field.
func (o *BackupOnlineArchive) SetPaused(v bool) {
	o.Paused = &v
}

// GetSchedule returns the Schedule field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetSchedule() OnlineArchiveSchedule {
	if o == nil || IsNil(o.Schedule) {
		var ret OnlineArchiveSchedule
		return ret
	}
	return *o.Schedule
}

// GetScheduleOk returns a tuple with the Schedule field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetScheduleOk() (*OnlineArchiveSchedule, bool) {
	if o == nil || IsNil(o.Schedule) {
		return nil, false
	}

	return o.Schedule, true
}

// HasSchedule returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasSchedule() bool {
	if o != nil && !IsNil(o.Schedule) {
		return true
	}

	return false
}

// SetSchedule gets a reference to the given OnlineArchiveSchedule and assigns it to the Schedule field.
func (o *BackupOnlineArchive) SetSchedule(v OnlineArchiveSchedule) {
	o.Schedule = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *BackupOnlineArchive) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchive) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *BackupOnlineArchive) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *BackupOnlineArchive) SetState(v string) {
	o.State = &v
}
