// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StateReason State reason of the Job. This is set when the job state is \"Failed\".
type StateReason struct {
	// Error code relating to state.
	ErrorCode *string `json:"errorCode,omitempty"`
	// Message describing error or state.
	Message *string `json:"message,omitempty"`
}

// NewStateReason instantiates a new StateReason object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStateReason() *StateReason {
	this := StateReason{}
	return &this
}

// NewStateReasonWithDefaults instantiates a new StateReason object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStateReasonWithDefaults() *StateReason {
	this := StateReason{}
	return &this
}

// GetErrorCode returns the ErrorCode field value if set, zero value otherwise
func (o *StateReason) GetErrorCode() string {
	if o == nil || IsNil(o.ErrorCode) {
		var ret string
		return ret
	}
	return *o.ErrorCode
}

// GetErrorCodeOk returns a tuple with the ErrorCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StateReason) GetErrorCodeOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorCode) {
		return nil, false
	}

	return o.ErrorCode, true
}

// HasErrorCode returns a boolean if a field has been set.
func (o *StateReason) HasErrorCode() bool {
	if o != nil && !IsNil(o.ErrorCode) {
		return true
	}

	return false
}

// SetErrorCode gets a reference to the given string and assigns it to the ErrorCode field.
func (o *StateReason) SetErrorCode(v string) {
	o.ErrorCode = &v
}

// GetMessage returns the Message field value if set, zero value otherwise
func (o *StateReason) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StateReason) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}

	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *StateReason) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *StateReason) SetMessage(v string) {
	o.Message = &v
}
