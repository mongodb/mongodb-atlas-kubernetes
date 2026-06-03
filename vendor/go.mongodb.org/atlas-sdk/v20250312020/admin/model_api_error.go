// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiError struct for ApiError
type ApiError struct {
	BadRequestDetail *BadRequestDetail `json:"badRequestDetail,omitempty"`
	// Describes the specific conditions or reasons that cause each type of error.
	Detail *string `json:"detail,omitempty"`
	// HTTP status code returned with this error.
	// Read only field.
	Error int `json:"error"`
	// Application error code returned with this error.
	// Read only field.
	ErrorCode string `json:"errorCode"`
	// Parameters used to give more information about the error.
	// Read only field.
	Parameters *[]any `json:"parameters,omitempty"`
	// Application error message returned with this error.
	// Read only field.
	Reason *string `json:"reason,omitempty"`
}

// NewApiError instantiates a new ApiError object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiError(error_ int, errorCode string) *ApiError {
	this := ApiError{}
	this.Error = error_
	this.ErrorCode = errorCode
	return &this
}

// NewApiErrorWithDefaults instantiates a new ApiError object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiErrorWithDefaults() *ApiError {
	this := ApiError{}
	return &this
}

// GetBadRequestDetail returns the BadRequestDetail field value if set, zero value otherwise
func (o *ApiError) GetBadRequestDetail() BadRequestDetail {
	if o == nil || IsNil(o.BadRequestDetail) {
		var ret BadRequestDetail
		return ret
	}
	return *o.BadRequestDetail
}

// GetBadRequestDetailOk returns a tuple with the BadRequestDetail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiError) GetBadRequestDetailOk() (*BadRequestDetail, bool) {
	if o == nil || IsNil(o.BadRequestDetail) {
		return nil, false
	}

	return o.BadRequestDetail, true
}

// HasBadRequestDetail returns a boolean if a field has been set.
func (o *ApiError) HasBadRequestDetail() bool {
	if o != nil && !IsNil(o.BadRequestDetail) {
		return true
	}

	return false
}

// SetBadRequestDetail gets a reference to the given BadRequestDetail and assigns it to the BadRequestDetail field.
func (o *ApiError) SetBadRequestDetail(v BadRequestDetail) {
	o.BadRequestDetail = &v
}

// GetDetail returns the Detail field value if set, zero value otherwise
func (o *ApiError) GetDetail() string {
	if o == nil || IsNil(o.Detail) {
		var ret string
		return ret
	}
	return *o.Detail
}

// GetDetailOk returns a tuple with the Detail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiError) GetDetailOk() (*string, bool) {
	if o == nil || IsNil(o.Detail) {
		return nil, false
	}

	return o.Detail, true
}

// HasDetail returns a boolean if a field has been set.
func (o *ApiError) HasDetail() bool {
	if o != nil && !IsNil(o.Detail) {
		return true
	}

	return false
}

// SetDetail gets a reference to the given string and assigns it to the Detail field.
func (o *ApiError) SetDetail(v string) {
	o.Detail = &v
}

// GetError returns the Error field value
func (o *ApiError) GetError() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.Error
}

// GetErrorOk returns a tuple with the Error field value
// and a boolean to check if the value has been set.
func (o *ApiError) GetErrorOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Error, true
}

// SetError sets field value
func (o *ApiError) SetError(v int) {
	o.Error = v
}

// GetErrorCode returns the ErrorCode field value
func (o *ApiError) GetErrorCode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ErrorCode
}

// GetErrorCodeOk returns a tuple with the ErrorCode field value
// and a boolean to check if the value has been set.
func (o *ApiError) GetErrorCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ErrorCode, true
}

// SetErrorCode sets field value
func (o *ApiError) SetErrorCode(v string) {
	o.ErrorCode = v
}

// GetParameters returns the Parameters field value if set, zero value otherwise
func (o *ApiError) GetParameters() []any {
	if o == nil || IsNil(o.Parameters) {
		var ret []any
		return ret
	}
	return *o.Parameters
}

// GetParametersOk returns a tuple with the Parameters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiError) GetParametersOk() (*[]any, bool) {
	if o == nil || IsNil(o.Parameters) {
		return nil, false
	}

	return o.Parameters, true
}

// HasParameters returns a boolean if a field has been set.
func (o *ApiError) HasParameters() bool {
	if o != nil && !IsNil(o.Parameters) {
		return true
	}

	return false
}

// SetParameters gets a reference to the given []any and assigns it to the Parameters field.
func (o *ApiError) SetParameters(v []any) {
	o.Parameters = &v
}

// GetReason returns the Reason field value if set, zero value otherwise
func (o *ApiError) GetReason() string {
	if o == nil || IsNil(o.Reason) {
		var ret string
		return ret
	}
	return *o.Reason
}

// GetReasonOk returns a tuple with the Reason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiError) GetReasonOk() (*string, bool) {
	if o == nil || IsNil(o.Reason) {
		return nil, false
	}

	return o.Reason, true
}

// HasReason returns a boolean if a field has been set.
func (o *ApiError) HasReason() bool {
	if o != nil && !IsNil(o.Reason) {
		return true
	}

	return false
}

// SetReason gets a reference to the given string and assigns it to the Reason field.
func (o *ApiError) SetReason(v string) {
	o.Reason = &v
}
