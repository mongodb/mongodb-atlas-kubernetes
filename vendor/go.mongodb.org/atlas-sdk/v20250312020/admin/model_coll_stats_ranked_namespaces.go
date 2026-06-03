// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CollStatsRankedNamespaces struct for CollStatsRankedNamespaces
type CollStatsRankedNamespaces struct {
	// Unique 24-hexadecimal digit string that identifies the request project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the request process.
	// Read only field.
	IdentifierId *string `json:"identifierId,omitempty"`
	// Ordered list of the hottest namespaces, highest value first.
	// Read only field.
	RankedNamespaces []string `json:"rankedNamespaces"`
}

// NewCollStatsRankedNamespaces instantiates a new CollStatsRankedNamespaces object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCollStatsRankedNamespaces(rankedNamespaces []string) *CollStatsRankedNamespaces {
	this := CollStatsRankedNamespaces{}
	this.RankedNamespaces = rankedNamespaces
	return &this
}

// NewCollStatsRankedNamespacesWithDefaults instantiates a new CollStatsRankedNamespaces object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCollStatsRankedNamespacesWithDefaults() *CollStatsRankedNamespaces {
	this := CollStatsRankedNamespaces{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *CollStatsRankedNamespaces) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CollStatsRankedNamespaces) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *CollStatsRankedNamespaces) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *CollStatsRankedNamespaces) SetGroupId(v string) {
	o.GroupId = &v
}

// GetIdentifierId returns the IdentifierId field value if set, zero value otherwise
func (o *CollStatsRankedNamespaces) GetIdentifierId() string {
	if o == nil || IsNil(o.IdentifierId) {
		var ret string
		return ret
	}
	return *o.IdentifierId
}

// GetIdentifierIdOk returns a tuple with the IdentifierId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CollStatsRankedNamespaces) GetIdentifierIdOk() (*string, bool) {
	if o == nil || IsNil(o.IdentifierId) {
		return nil, false
	}

	return o.IdentifierId, true
}

// HasIdentifierId returns a boolean if a field has been set.
func (o *CollStatsRankedNamespaces) HasIdentifierId() bool {
	if o != nil && !IsNil(o.IdentifierId) {
		return true
	}

	return false
}

// SetIdentifierId gets a reference to the given string and assigns it to the IdentifierId field.
func (o *CollStatsRankedNamespaces) SetIdentifierId(v string) {
	o.IdentifierId = &v
}

// GetRankedNamespaces returns the RankedNamespaces field value
func (o *CollStatsRankedNamespaces) GetRankedNamespaces() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.RankedNamespaces
}

// GetRankedNamespacesOk returns a tuple with the RankedNamespaces field value
// and a boolean to check if the value has been set.
func (o *CollStatsRankedNamespaces) GetRankedNamespacesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RankedNamespaces, true
}

// SetRankedNamespaces sets field value
func (o *CollStatsRankedNamespaces) SetRankedNamespaces(v []string) {
	o.RankedNamespaces = v
}
