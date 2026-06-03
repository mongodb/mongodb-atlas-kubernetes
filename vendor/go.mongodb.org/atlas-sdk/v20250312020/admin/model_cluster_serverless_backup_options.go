// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterServerlessBackupOptions Group of settings that configure serverless backup.
type ClusterServerlessBackupOptions struct {
	// Flag that indicates whether the serverless instance uses **Serverless Continuous Backup**.  If this parameter is `false`, the serverless instance uses **Basic Backup**.   | Option | Description |  |---|---|  | Serverless Continuous Backup | Atlas takes incremental [snapshots](https://www.mongodb.com/docs/atlas/backup/cloud-backup/overview/#std-label-serverless-snapshots) of the data in your serverless instance every six hours and lets you restore the data from a selected point in time within the last 72 hours. Atlas also takes daily snapshots and retains these daily snapshots for 35 days. To learn more, see [Serverless Instance Costs](https://www.mongodb.com/docs/atlas/billing/serverless-instance-costs/#std-label-serverless-instance-costs). |  | Basic Backup | Atlas takes incremental [snapshots](https://www.mongodb.com/docs/atlas/backup/cloud-backup/overview/#std-label-serverless-snapshots) of the data in your serverless instance every six hours and retains only the two most recent snapshots. You can use this option for free. |
	ServerlessContinuousBackupEnabled *bool `json:"serverlessContinuousBackupEnabled,omitempty"`
}

// NewClusterServerlessBackupOptions instantiates a new ClusterServerlessBackupOptions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterServerlessBackupOptions() *ClusterServerlessBackupOptions {
	this := ClusterServerlessBackupOptions{}
	var serverlessContinuousBackupEnabled bool = true
	this.ServerlessContinuousBackupEnabled = &serverlessContinuousBackupEnabled
	return &this
}

// NewClusterServerlessBackupOptionsWithDefaults instantiates a new ClusterServerlessBackupOptions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterServerlessBackupOptionsWithDefaults() *ClusterServerlessBackupOptions {
	this := ClusterServerlessBackupOptions{}
	var serverlessContinuousBackupEnabled bool = true
	this.ServerlessContinuousBackupEnabled = &serverlessContinuousBackupEnabled
	return &this
}

// GetServerlessContinuousBackupEnabled returns the ServerlessContinuousBackupEnabled field value if set, zero value otherwise
func (o *ClusterServerlessBackupOptions) GetServerlessContinuousBackupEnabled() bool {
	if o == nil || IsNil(o.ServerlessContinuousBackupEnabled) {
		var ret bool
		return ret
	}
	return *o.ServerlessContinuousBackupEnabled
}

// GetServerlessContinuousBackupEnabledOk returns a tuple with the ServerlessContinuousBackupEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterServerlessBackupOptions) GetServerlessContinuousBackupEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.ServerlessContinuousBackupEnabled) {
		return nil, false
	}

	return o.ServerlessContinuousBackupEnabled, true
}

// HasServerlessContinuousBackupEnabled returns a boolean if a field has been set.
func (o *ClusterServerlessBackupOptions) HasServerlessContinuousBackupEnabled() bool {
	if o != nil && !IsNil(o.ServerlessContinuousBackupEnabled) {
		return true
	}

	return false
}

// SetServerlessContinuousBackupEnabled gets a reference to the given bool and assigns it to the ServerlessContinuousBackupEnabled field.
func (o *ClusterServerlessBackupOptions) SetServerlessContinuousBackupEnabled(v bool) {
	o.ServerlessContinuousBackupEnabled = &v
}
