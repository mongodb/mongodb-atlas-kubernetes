// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataExpirationRule Rule for specifying when data should be deleted from the archive.
type DataExpirationRule struct {
	// Number of days used in the date criteria for nominating documents for deletion.
	ExpireAfterDays *int `json:"expireAfterDays,omitempty"`
}

// NewDataExpirationRule instantiates a new DataExpirationRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataExpirationRule() *DataExpirationRule {
	this := DataExpirationRule{}
	return &this
}

// NewDataExpirationRuleWithDefaults instantiates a new DataExpirationRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataExpirationRuleWithDefaults() *DataExpirationRule {
	this := DataExpirationRule{}
	return &this
}

// GetExpireAfterDays returns the ExpireAfterDays field value if set, zero value otherwise
func (o *DataExpirationRule) GetExpireAfterDays() int {
	if o == nil || IsNil(o.ExpireAfterDays) {
		var ret int
		return ret
	}
	return *o.ExpireAfterDays
}

// GetExpireAfterDaysOk returns a tuple with the ExpireAfterDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataExpirationRule) GetExpireAfterDaysOk() (*int, bool) {
	if o == nil || IsNil(o.ExpireAfterDays) {
		return nil, false
	}

	return o.ExpireAfterDays, true
}

// HasExpireAfterDays returns a boolean if a field has been set.
func (o *DataExpirationRule) HasExpireAfterDays() bool {
	if o != nil && !IsNil(o.ExpireAfterDays) {
		return true
	}

	return false
}

// SetExpireAfterDays gets a reference to the given int and assigns it to the ExpireAfterDays field.
func (o *DataExpirationRule) SetExpireAfterDays(v int) {
	o.ExpireAfterDays = &v
}
