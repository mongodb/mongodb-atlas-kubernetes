// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ClusterOutageSimulation struct for ClusterOutageSimulation
type ClusterOutageSimulation struct {
	// Human-readable label that identifies the cluster that undergoes outage simulation.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Date and time when MongoDB Cloud expires the outage simulation. This parameter expresses its value in the ISO 8601 timestamp format in UTC. If not provided, defaults to 3 days from the start date.
	// Read only field.
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project that contains the cluster to undergo outage simulation.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal character string that identifies the outage simulation.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of settings that specify the type of cluster outage simulation.
	OutageFilters *[]AtlasClusterOutageSimulationOutageFilter `json:"outageFilters,omitempty"`
	// Date and time when MongoDB Cloud started the regional outage simulation. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	StartRequestDate *time.Time `json:"startRequestDate,omitempty"`
	// Phase of the outage simulation.  | State       | Indication | |-------------|------------| | `START_REQUESTED`    | User has requested cluster outage simulation.| | `STARTING`           | MongoDB Cloud is starting cluster outage simulation.| | `SIMULATING`         | MongoDB Cloud is simulating cluster outage.| | `RECOVERY_REQUESTED` | User has requested recovery from the simulated outage.| | `RECOVERING`         | MongoDB Cloud is recovering the cluster from the simulated outage.| | `COMPLETE`           | MongoDB Cloud has completed the cluster outage simulation.|
	// Read only field.
	State *string `json:"state,omitempty"`
}

// NewClusterOutageSimulation instantiates a new ClusterOutageSimulation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterOutageSimulation() *ClusterOutageSimulation {
	this := ClusterOutageSimulation{}
	return &this
}

// NewClusterOutageSimulationWithDefaults instantiates a new ClusterOutageSimulation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterOutageSimulationWithDefaults() *ClusterOutageSimulation {
	this := ClusterOutageSimulation{}
	return &this
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *ClusterOutageSimulation) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetExpirationDate returns the ExpirationDate field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetExpirationDate() time.Time {
	if o == nil || IsNil(o.ExpirationDate) {
		var ret time.Time
		return ret
	}
	return *o.ExpirationDate
}

// GetExpirationDateOk returns a tuple with the ExpirationDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetExpirationDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpirationDate) {
		return nil, false
	}

	return o.ExpirationDate, true
}

// HasExpirationDate returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasExpirationDate() bool {
	if o != nil && !IsNil(o.ExpirationDate) {
		return true
	}

	return false
}

// SetExpirationDate gets a reference to the given time.Time and assigns it to the ExpirationDate field.
func (o *ClusterOutageSimulation) SetExpirationDate(v time.Time) {
	o.ExpirationDate = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ClusterOutageSimulation) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ClusterOutageSimulation) SetId(v string) {
	o.Id = &v
}

// GetOutageFilters returns the OutageFilters field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetOutageFilters() []AtlasClusterOutageSimulationOutageFilter {
	if o == nil || IsNil(o.OutageFilters) {
		var ret []AtlasClusterOutageSimulationOutageFilter
		return ret
	}
	return *o.OutageFilters
}

// GetOutageFiltersOk returns a tuple with the OutageFilters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetOutageFiltersOk() (*[]AtlasClusterOutageSimulationOutageFilter, bool) {
	if o == nil || IsNil(o.OutageFilters) {
		return nil, false
	}

	return o.OutageFilters, true
}

// HasOutageFilters returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasOutageFilters() bool {
	if o != nil && !IsNil(o.OutageFilters) {
		return true
	}

	return false
}

// SetOutageFilters gets a reference to the given []AtlasClusterOutageSimulationOutageFilter and assigns it to the OutageFilters field.
func (o *ClusterOutageSimulation) SetOutageFilters(v []AtlasClusterOutageSimulationOutageFilter) {
	o.OutageFilters = &v
}

// GetStartRequestDate returns the StartRequestDate field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetStartRequestDate() time.Time {
	if o == nil || IsNil(o.StartRequestDate) {
		var ret time.Time
		return ret
	}
	return *o.StartRequestDate
}

// GetStartRequestDateOk returns a tuple with the StartRequestDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetStartRequestDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartRequestDate) {
		return nil, false
	}

	return o.StartRequestDate, true
}

// HasStartRequestDate returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasStartRequestDate() bool {
	if o != nil && !IsNil(o.StartRequestDate) {
		return true
	}

	return false
}

// SetStartRequestDate gets a reference to the given time.Time and assigns it to the StartRequestDate field.
func (o *ClusterOutageSimulation) SetStartRequestDate(v time.Time) {
	o.StartRequestDate = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *ClusterOutageSimulation) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterOutageSimulation) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *ClusterOutageSimulation) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *ClusterOutageSimulation) SetState(v string) {
	o.State = &v
}
