// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ApiBSONTimestamp BSON timestamp that indicates when the checkpoint token entry in the oplog occurred.
type ApiBSONTimestamp struct {
	// Date and time when the oplog recorded this database operation. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Date *time.Time `json:"date,omitempty"`
	// Order of the database operation that the oplog recorded at specific date and time.
	// Read only field.
	Increment *int `json:"increment,omitempty"`
}

// NewApiBSONTimestamp instantiates a new ApiBSONTimestamp object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiBSONTimestamp() *ApiBSONTimestamp {
	this := ApiBSONTimestamp{}
	return &this
}

// NewApiBSONTimestampWithDefaults instantiates a new ApiBSONTimestamp object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiBSONTimestampWithDefaults() *ApiBSONTimestamp {
	this := ApiBSONTimestamp{}
	return &this
}

// GetDate returns the Date field value if set, zero value otherwise
func (o *ApiBSONTimestamp) GetDate() time.Time {
	if o == nil || IsNil(o.Date) {
		var ret time.Time
		return ret
	}
	return *o.Date
}

// GetDateOk returns a tuple with the Date field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiBSONTimestamp) GetDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Date) {
		return nil, false
	}

	return o.Date, true
}

// HasDate returns a boolean if a field has been set.
func (o *ApiBSONTimestamp) HasDate() bool {
	if o != nil && !IsNil(o.Date) {
		return true
	}

	return false
}

// SetDate gets a reference to the given time.Time and assigns it to the Date field.
func (o *ApiBSONTimestamp) SetDate(v time.Time) {
	o.Date = &v
}

// GetIncrement returns the Increment field value if set, zero value otherwise
func (o *ApiBSONTimestamp) GetIncrement() int {
	if o == nil || IsNil(o.Increment) {
		var ret int
		return ret
	}
	return *o.Increment
}

// GetIncrementOk returns a tuple with the Increment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiBSONTimestamp) GetIncrementOk() (*int, bool) {
	if o == nil || IsNil(o.Increment) {
		return nil, false
	}

	return o.Increment, true
}

// HasIncrement returns a boolean if a field has been set.
func (o *ApiBSONTimestamp) HasIncrement() bool {
	if o != nil && !IsNil(o.Increment) {
		return true
	}

	return false
}

// SetIncrement gets a reference to the given int and assigns it to the Increment field.
func (o *ApiBSONTimestamp) SetIncrement(v int) {
	o.Increment = &v
}
