// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ProjectSettingItem struct for ProjectSettingItem
type ProjectSettingItem struct {
	// Flag that indicates whether someone enabled the regionalized private endpoint setting for the specified project.  - Set this value to `true` to enable regionalized private endpoints. This allows you to create more than one private endpoint in a cloud provider region. You need to enable this setting to connect to multi-region and global MongoDB Cloud sharded clusters. Enabling regionalized private endpoints introduces the following limitations:   - Your applications must use the new connection strings for existing multi-region and global sharded clusters. This might cause downtime.   - Your MongoDB Cloud project can't contain replica sets nor can you create new replica sets in this project.    - You can't disable this setting if you have:     - more than one private endpoint in more than one region     - more than one private endpoint in one region and one private endpoint in one or more regions.  - Set this value to `false` to disable regionalized private endpoints.
	Enabled bool `json:"enabled"`
}

// NewProjectSettingItem instantiates a new ProjectSettingItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectSettingItem(enabled bool) *ProjectSettingItem {
	this := ProjectSettingItem{}
	this.Enabled = enabled
	return &this
}

// NewProjectSettingItemWithDefaults instantiates a new ProjectSettingItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectSettingItemWithDefaults() *ProjectSettingItem {
	this := ProjectSettingItem{}
	return &this
}

// GetEnabled returns the Enabled field value
func (o *ProjectSettingItem) GetEnabled() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value
// and a boolean to check if the value has been set.
func (o *ProjectSettingItem) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Enabled, true
}

// SetEnabled sets field value
func (o *ProjectSettingItem) SetEnabled(v bool) {
	o.Enabled = v
}
