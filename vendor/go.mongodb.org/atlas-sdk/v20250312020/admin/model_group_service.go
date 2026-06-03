// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupService List of IP addresses in a project categorized by services.
type GroupService struct {
	// IP addresses of clusters.
	// Read only field.
	Clusters *[]ClusterIPAddresses `json:"clusters,omitempty"`
}

// NewGroupService instantiates a new GroupService object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupService() *GroupService {
	this := GroupService{}
	return &this
}

// NewGroupServiceWithDefaults instantiates a new GroupService object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupServiceWithDefaults() *GroupService {
	this := GroupService{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise
func (o *GroupService) GetClusters() []ClusterIPAddresses {
	if o == nil || IsNil(o.Clusters) {
		var ret []ClusterIPAddresses
		return ret
	}
	return *o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupService) GetClustersOk() (*[]ClusterIPAddresses, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}

	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *GroupService) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []ClusterIPAddresses and assigns it to the Clusters field.
func (o *GroupService) SetClusters(v []ClusterIPAddresses) {
	o.Clusters = &v
}
