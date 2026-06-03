// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LiveImportValidation struct for LiveImportValidation
type LiveImportValidation struct {
	// Unique 24-hexadecimal digit string that identifies the validation.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Reason why the validation job failed.
	// Read only field.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project to validate.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the source project.
	SourceGroupId *string `json:"sourceGroupId,omitempty"`
	// State of the specified validation job returned at the time of the request.
	// Read only field.
	Status *string `json:"status,omitempty"`
}

// NewLiveImportValidation instantiates a new LiveImportValidation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLiveImportValidation() *LiveImportValidation {
	this := LiveImportValidation{}
	return &this
}

// NewLiveImportValidationWithDefaults instantiates a new LiveImportValidation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLiveImportValidationWithDefaults() *LiveImportValidation {
	this := LiveImportValidation{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *LiveImportValidation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveImportValidation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *LiveImportValidation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *LiveImportValidation) SetId(v string) {
	o.Id = &v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise
func (o *LiveImportValidation) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveImportValidation) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}

	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *LiveImportValidation) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *LiveImportValidation) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *LiveImportValidation) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveImportValidation) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *LiveImportValidation) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *LiveImportValidation) SetGroupId(v string) {
	o.GroupId = &v
}

// GetSourceGroupId returns the SourceGroupId field value if set, zero value otherwise
func (o *LiveImportValidation) GetSourceGroupId() string {
	if o == nil || IsNil(o.SourceGroupId) {
		var ret string
		return ret
	}
	return *o.SourceGroupId
}

// GetSourceGroupIdOk returns a tuple with the SourceGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveImportValidation) GetSourceGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.SourceGroupId) {
		return nil, false
	}

	return o.SourceGroupId, true
}

// HasSourceGroupId returns a boolean if a field has been set.
func (o *LiveImportValidation) HasSourceGroupId() bool {
	if o != nil && !IsNil(o.SourceGroupId) {
		return true
	}

	return false
}

// SetSourceGroupId gets a reference to the given string and assigns it to the SourceGroupId field.
func (o *LiveImportValidation) SetSourceGroupId(v string) {
	o.SourceGroupId = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *LiveImportValidation) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveImportValidation) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *LiveImportValidation) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *LiveImportValidation) SetStatus(v string) {
	o.Status = &v
}
