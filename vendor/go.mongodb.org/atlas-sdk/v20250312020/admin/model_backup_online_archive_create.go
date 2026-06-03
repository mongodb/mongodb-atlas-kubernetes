// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BackupOnlineArchiveCreate struct for BackupOnlineArchiveCreate
type BackupOnlineArchiveCreate struct {
	// Unique 24-hexadecimal digit string that identifies the online archive.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Human-readable label that identifies the cluster that contains the collection for which you want to create an online archive.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Human-readable label that identifies the collection for which you created the online archive.
	// Write only field.
	CollName string `json:"collName"`
	// Classification of MongoDB database collection that you want to return.  If you set this parameter to `TIMESERIES`, set `\"criteria.type\" : \"date\"` and `\"criteria.dateFormat\" : \"ISODATE\"`.
	// Write only field.
	CollectionType     *string                  `json:"collectionType,omitempty"`
	Criteria           Criteria                 `json:"criteria"`
	DataExpirationRule *DataExpirationRule      `json:"dataExpirationRule,omitempty"`
	DataProcessRegion  *CreateDataProcessRegion `json:"dataProcessRegion,omitempty"`
	// Human-readable label that identifies the dataset that Atlas generates for this online archive.
	// Read only field.
	DataSetName *string `json:"dataSetName,omitempty"`
	// Human-readable label of the database that contains the collection that contains the online archive.
	// Write only field.
	DbName string `json:"dbName"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the specified cluster. The specified cluster contains the collection for which to create the online archive.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List that contains document parameters to use to logically divide data within a collection. Partitions provide a coarse level of filtering of the underlying collection data. To divide your data, specify parameters that you frequently query. If you specified `criteria.type`: `DATE` in the Create One Online Archive endpoint, then you can specify up to three parameters by which to query. One of these parameters must be the `DATE` value, which is required in this case. If you specified `criteria.type`: `CUSTOM` in the Create One Online Archive endpoint, then you can specify up to two parameters by which to query. Queries that don't use `criteria.type`: `DATE` or `criteria.type`: `CUSTOM` parameters cause MongoDB to scan a full collection of all archived documents. This takes more time and increases your costs.
	// Write only field.
	PartitionFields *[]PartitionField `json:"partitionFields,omitempty"`
	// Flag that indicates whether this online archive exists in the paused state. A request to resume fails if the collection has another active online archive. To pause an active online archive or resume a paused online archive, you must include this parameter. To pause an active archive, set this to **true**. To resume a paused archive, set this to **false**.
	Paused   *bool                  `json:"paused,omitempty"`
	Schedule *OnlineArchiveSchedule `json:"schedule,omitempty"`
	// Phase of the process to create this online archive when you made this request.  | State       | Indication | |-------------|------------| | `PENDING`   | MongoDB Cloud has queued documents for archive. Archiving hasn't started. | | `ARCHIVING` | MongoDB Cloud started archiving documents that meet the archival criteria. | | `IDLE`      | MongoDB Cloud waits to start the next archival job. | | `PAUSING`   | Someone chose to stop archiving. MongoDB Cloud finishes the running archival job then changes the state to `PAUSED` when that job completes. | | `PAUSED`    | MongoDB Cloud has stopped archiving. Archived documents can be queried. The specified archiving operation on the active cluster cannot archive additional documents. You can resume archiving for paused archives at any time. | | `ORPHANED`  | Someone has deleted the collection associated with an active or paused archive. MongoDB Cloud doesn't delete the archived data. You must manually delete the online archives associated with the deleted collection. | | `DELETED`   | Someone has deleted the archive was deleted. When someone deletes an online archive, MongoDB Cloud removes all associated archived documents from the cloud object storage. |
	// Read only field.
	State *string `json:"state,omitempty"`
}

// NewBackupOnlineArchiveCreate instantiates a new BackupOnlineArchiveCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupOnlineArchiveCreate(collName string, criteria Criteria, dbName string) *BackupOnlineArchiveCreate {
	this := BackupOnlineArchiveCreate{}
	this.CollName = collName
	var collectionType string = "STANDARD"
	this.CollectionType = &collectionType
	this.Criteria = criteria
	this.DbName = dbName
	return &this
}

// NewBackupOnlineArchiveCreateWithDefaults instantiates a new BackupOnlineArchiveCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupOnlineArchiveCreateWithDefaults() *BackupOnlineArchiveCreate {
	this := BackupOnlineArchiveCreate{}
	var collectionType string = "STANDARD"
	this.CollectionType = &collectionType
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BackupOnlineArchiveCreate) SetId(v string) {
	o.Id = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *BackupOnlineArchiveCreate) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCollName returns the CollName field value
func (o *BackupOnlineArchiveCreate) GetCollName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CollName
}

// GetCollNameOk returns a tuple with the CollName field value
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetCollNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CollName, true
}

// SetCollName sets field value
func (o *BackupOnlineArchiveCreate) SetCollName(v string) {
	o.CollName = v
}

// GetCollectionType returns the CollectionType field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetCollectionType() string {
	if o == nil || IsNil(o.CollectionType) {
		var ret string
		return ret
	}
	return *o.CollectionType
}

// GetCollectionTypeOk returns a tuple with the CollectionType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetCollectionTypeOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionType) {
		return nil, false
	}

	return o.CollectionType, true
}

// HasCollectionType returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasCollectionType() bool {
	if o != nil && !IsNil(o.CollectionType) {
		return true
	}

	return false
}

// SetCollectionType gets a reference to the given string and assigns it to the CollectionType field.
func (o *BackupOnlineArchiveCreate) SetCollectionType(v string) {
	o.CollectionType = &v
}

// GetCriteria returns the Criteria field value
func (o *BackupOnlineArchiveCreate) GetCriteria() Criteria {
	if o == nil {
		var ret Criteria
		return ret
	}

	return o.Criteria
}

// GetCriteriaOk returns a tuple with the Criteria field value
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetCriteriaOk() (*Criteria, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Criteria, true
}

// SetCriteria sets field value
func (o *BackupOnlineArchiveCreate) SetCriteria(v Criteria) {
	o.Criteria = v
}

// GetDataExpirationRule returns the DataExpirationRule field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetDataExpirationRule() DataExpirationRule {
	if o == nil || IsNil(o.DataExpirationRule) {
		var ret DataExpirationRule
		return ret
	}
	return *o.DataExpirationRule
}

// GetDataExpirationRuleOk returns a tuple with the DataExpirationRule field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetDataExpirationRuleOk() (*DataExpirationRule, bool) {
	if o == nil || IsNil(o.DataExpirationRule) {
		return nil, false
	}

	return o.DataExpirationRule, true
}

// HasDataExpirationRule returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasDataExpirationRule() bool {
	if o != nil && !IsNil(o.DataExpirationRule) {
		return true
	}

	return false
}

// SetDataExpirationRule gets a reference to the given DataExpirationRule and assigns it to the DataExpirationRule field.
func (o *BackupOnlineArchiveCreate) SetDataExpirationRule(v DataExpirationRule) {
	o.DataExpirationRule = &v
}

// GetDataProcessRegion returns the DataProcessRegion field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetDataProcessRegion() CreateDataProcessRegion {
	if o == nil || IsNil(o.DataProcessRegion) {
		var ret CreateDataProcessRegion
		return ret
	}
	return *o.DataProcessRegion
}

// GetDataProcessRegionOk returns a tuple with the DataProcessRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetDataProcessRegionOk() (*CreateDataProcessRegion, bool) {
	if o == nil || IsNil(o.DataProcessRegion) {
		return nil, false
	}

	return o.DataProcessRegion, true
}

// HasDataProcessRegion returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasDataProcessRegion() bool {
	if o != nil && !IsNil(o.DataProcessRegion) {
		return true
	}

	return false
}

// SetDataProcessRegion gets a reference to the given CreateDataProcessRegion and assigns it to the DataProcessRegion field.
func (o *BackupOnlineArchiveCreate) SetDataProcessRegion(v CreateDataProcessRegion) {
	o.DataProcessRegion = &v
}

// GetDataSetName returns the DataSetName field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetDataSetName() string {
	if o == nil || IsNil(o.DataSetName) {
		var ret string
		return ret
	}
	return *o.DataSetName
}

// GetDataSetNameOk returns a tuple with the DataSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetDataSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.DataSetName) {
		return nil, false
	}

	return o.DataSetName, true
}

// HasDataSetName returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasDataSetName() bool {
	if o != nil && !IsNil(o.DataSetName) {
		return true
	}

	return false
}

// SetDataSetName gets a reference to the given string and assigns it to the DataSetName field.
func (o *BackupOnlineArchiveCreate) SetDataSetName(v string) {
	o.DataSetName = &v
}

// GetDbName returns the DbName field value
func (o *BackupOnlineArchiveCreate) GetDbName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DbName
}

// GetDbNameOk returns a tuple with the DbName field value
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetDbNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DbName, true
}

// SetDbName sets field value
func (o *BackupOnlineArchiveCreate) SetDbName(v string) {
	o.DbName = v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *BackupOnlineArchiveCreate) SetGroupId(v string) {
	o.GroupId = &v
}

// GetPartitionFields returns the PartitionFields field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetPartitionFields() []PartitionField {
	if o == nil || IsNil(o.PartitionFields) {
		var ret []PartitionField
		return ret
	}
	return *o.PartitionFields
}

// GetPartitionFieldsOk returns a tuple with the PartitionFields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetPartitionFieldsOk() (*[]PartitionField, bool) {
	if o == nil || IsNil(o.PartitionFields) {
		return nil, false
	}

	return o.PartitionFields, true
}

// HasPartitionFields returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasPartitionFields() bool {
	if o != nil && !IsNil(o.PartitionFields) {
		return true
	}

	return false
}

// SetPartitionFields gets a reference to the given []PartitionField and assigns it to the PartitionFields field.
func (o *BackupOnlineArchiveCreate) SetPartitionFields(v []PartitionField) {
	o.PartitionFields = &v
}

// GetPaused returns the Paused field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetPaused() bool {
	if o == nil || IsNil(o.Paused) {
		var ret bool
		return ret
	}
	return *o.Paused
}

// GetPausedOk returns a tuple with the Paused field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetPausedOk() (*bool, bool) {
	if o == nil || IsNil(o.Paused) {
		return nil, false
	}

	return o.Paused, true
}

// HasPaused returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasPaused() bool {
	if o != nil && !IsNil(o.Paused) {
		return true
	}

	return false
}

// SetPaused gets a reference to the given bool and assigns it to the Paused field.
func (o *BackupOnlineArchiveCreate) SetPaused(v bool) {
	o.Paused = &v
}

// GetSchedule returns the Schedule field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetSchedule() OnlineArchiveSchedule {
	if o == nil || IsNil(o.Schedule) {
		var ret OnlineArchiveSchedule
		return ret
	}
	return *o.Schedule
}

// GetScheduleOk returns a tuple with the Schedule field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetScheduleOk() (*OnlineArchiveSchedule, bool) {
	if o == nil || IsNil(o.Schedule) {
		return nil, false
	}

	return o.Schedule, true
}

// HasSchedule returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasSchedule() bool {
	if o != nil && !IsNil(o.Schedule) {
		return true
	}

	return false
}

// SetSchedule gets a reference to the given OnlineArchiveSchedule and assigns it to the Schedule field.
func (o *BackupOnlineArchiveCreate) SetSchedule(v OnlineArchiveSchedule) {
	o.Schedule = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *BackupOnlineArchiveCreate) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupOnlineArchiveCreate) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *BackupOnlineArchiveCreate) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *BackupOnlineArchiveCreate) SetState(v string) {
	o.State = &v
}
