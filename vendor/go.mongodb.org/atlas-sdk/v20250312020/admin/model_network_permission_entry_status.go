// Code based on the AtlasAPI V2 OpenAPI file

package admin

// NetworkPermissionEntryStatus struct for NetworkPermissionEntryStatus
type NetworkPermissionEntryStatus struct {
	// State of the access list entry when MongoDB Cloud made this request.  `ACTIVE`: This access list entry applies to all relevant cloud providers.  `PENDING`: MongoDB Cloud has started to add access list entry. This access list entry may not apply to all cloud providers at the time of this request.  `FAILED`: MongoDB Cloud didn't succeed in adding this access list entry.
	// Read only field.
	STATUS string `json:"STATUS"`
}

// NewNetworkPermissionEntryStatus instantiates a new NetworkPermissionEntryStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNetworkPermissionEntryStatus(sTATUS string) *NetworkPermissionEntryStatus {
	this := NetworkPermissionEntryStatus{}
	this.STATUS = sTATUS
	return &this
}

// NewNetworkPermissionEntryStatusWithDefaults instantiates a new NetworkPermissionEntryStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNetworkPermissionEntryStatusWithDefaults() *NetworkPermissionEntryStatus {
	this := NetworkPermissionEntryStatus{}
	return &this
}

// GetSTATUS returns the STATUS field value
func (o *NetworkPermissionEntryStatus) GetSTATUS() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.STATUS
}

// GetSTATUSOk returns a tuple with the STATUS field value
// and a boolean to check if the value has been set.
func (o *NetworkPermissionEntryStatus) GetSTATUSOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.STATUS, true
}

// SetSTATUS sets field value
func (o *NetworkPermissionEntryStatus) SetSTATUS(v string) {
	o.STATUS = v
}
