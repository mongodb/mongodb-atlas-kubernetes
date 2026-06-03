// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// Raw Additional meta information captured about this event. The response returns this parameter as a JSON object when the query parameter `includeRaw=true`. The list of fields in the raw document may change. Don't rely on raw values for formal monitoring.
type Raw struct {
	// Unique identifier of event type.
	T *string `json:"_t,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the alert configuration related to the event.
	// Read only field.
	AlertConfigId *string `json:"alertConfigId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project in which the event occurred.
	// Read only field.
	Cid *string `json:"cid,omitempty"`
	// Date and time when this event occurred. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Cre *time.Time `json:"cre,omitempty"`
	// Description of the event.
	Description *string `json:"description,omitempty"`
	// Human-readable label that identifies the project.
	Gn *string `json:"gn,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the event.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization to which these events apply.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// Human-readable label that identifies the organization that contains the project.
	OrgName *string `json:"orgName,omitempty"`
	// Severity of the event.
	Severity *string `json:"severity,omitempty"`
}

// NewRaw instantiates a new Raw object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRaw() *Raw {
	this := Raw{}
	return &this
}

// NewRawWithDefaults instantiates a new Raw object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRawWithDefaults() *Raw {
	this := Raw{}
	return &this
}

// GetT returns the T field value if set, zero value otherwise
func (o *Raw) GetT() string {
	if o == nil || IsNil(o.T) {
		var ret string
		return ret
	}
	return *o.T
}

// GetTOk returns a tuple with the T field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetTOk() (*string, bool) {
	if o == nil || IsNil(o.T) {
		return nil, false
	}

	return o.T, true
}

// HasT returns a boolean if a field has been set.
func (o *Raw) HasT() bool {
	if o != nil && !IsNil(o.T) {
		return true
	}

	return false
}

// SetT gets a reference to the given string and assigns it to the T field.
func (o *Raw) SetT(v string) {
	o.T = &v
}

// GetAlertConfigId returns the AlertConfigId field value if set, zero value otherwise
func (o *Raw) GetAlertConfigId() string {
	if o == nil || IsNil(o.AlertConfigId) {
		var ret string
		return ret
	}
	return *o.AlertConfigId
}

// GetAlertConfigIdOk returns a tuple with the AlertConfigId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetAlertConfigIdOk() (*string, bool) {
	if o == nil || IsNil(o.AlertConfigId) {
		return nil, false
	}

	return o.AlertConfigId, true
}

// HasAlertConfigId returns a boolean if a field has been set.
func (o *Raw) HasAlertConfigId() bool {
	if o != nil && !IsNil(o.AlertConfigId) {
		return true
	}

	return false
}

// SetAlertConfigId gets a reference to the given string and assigns it to the AlertConfigId field.
func (o *Raw) SetAlertConfigId(v string) {
	o.AlertConfigId = &v
}

// GetCid returns the Cid field value if set, zero value otherwise
func (o *Raw) GetCid() string {
	if o == nil || IsNil(o.Cid) {
		var ret string
		return ret
	}
	return *o.Cid
}

// GetCidOk returns a tuple with the Cid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetCidOk() (*string, bool) {
	if o == nil || IsNil(o.Cid) {
		return nil, false
	}

	return o.Cid, true
}

// HasCid returns a boolean if a field has been set.
func (o *Raw) HasCid() bool {
	if o != nil && !IsNil(o.Cid) {
		return true
	}

	return false
}

// SetCid gets a reference to the given string and assigns it to the Cid field.
func (o *Raw) SetCid(v string) {
	o.Cid = &v
}

// GetCre returns the Cre field value if set, zero value otherwise
func (o *Raw) GetCre() time.Time {
	if o == nil || IsNil(o.Cre) {
		var ret time.Time
		return ret
	}
	return *o.Cre
}

// GetCreOk returns a tuple with the Cre field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetCreOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Cre) {
		return nil, false
	}

	return o.Cre, true
}

// HasCre returns a boolean if a field has been set.
func (o *Raw) HasCre() bool {
	if o != nil && !IsNil(o.Cre) {
		return true
	}

	return false
}

// SetCre gets a reference to the given time.Time and assigns it to the Cre field.
func (o *Raw) SetCre(v time.Time) {
	o.Cre = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *Raw) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Raw) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Raw) SetDescription(v string) {
	o.Description = &v
}

// GetGn returns the Gn field value if set, zero value otherwise
func (o *Raw) GetGn() string {
	if o == nil || IsNil(o.Gn) {
		var ret string
		return ret
	}
	return *o.Gn
}

// GetGnOk returns a tuple with the Gn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetGnOk() (*string, bool) {
	if o == nil || IsNil(o.Gn) {
		return nil, false
	}

	return o.Gn, true
}

// HasGn returns a boolean if a field has been set.
func (o *Raw) HasGn() bool {
	if o != nil && !IsNil(o.Gn) {
		return true
	}

	return false
}

// SetGn gets a reference to the given string and assigns it to the Gn field.
func (o *Raw) SetGn(v string) {
	o.Gn = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *Raw) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Raw) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Raw) SetId(v string) {
	o.Id = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *Raw) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *Raw) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *Raw) SetOrgId(v string) {
	o.OrgId = &v
}

// GetOrgName returns the OrgName field value if set, zero value otherwise
func (o *Raw) GetOrgName() string {
	if o == nil || IsNil(o.OrgName) {
		var ret string
		return ret
	}
	return *o.OrgName
}

// GetOrgNameOk returns a tuple with the OrgName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetOrgNameOk() (*string, bool) {
	if o == nil || IsNil(o.OrgName) {
		return nil, false
	}

	return o.OrgName, true
}

// HasOrgName returns a boolean if a field has been set.
func (o *Raw) HasOrgName() bool {
	if o != nil && !IsNil(o.OrgName) {
		return true
	}

	return false
}

// SetOrgName gets a reference to the given string and assigns it to the OrgName field.
func (o *Raw) SetOrgName(v string) {
	o.OrgName = &v
}

// GetSeverity returns the Severity field value if set, zero value otherwise
func (o *Raw) GetSeverity() string {
	if o == nil || IsNil(o.Severity) {
		var ret string
		return ret
	}
	return *o.Severity
}

// GetSeverityOk returns a tuple with the Severity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Raw) GetSeverityOk() (*string, bool) {
	if o == nil || IsNil(o.Severity) {
		return nil, false
	}

	return o.Severity, true
}

// HasSeverity returns a boolean if a field has been set.
func (o *Raw) HasSeverity() bool {
	if o != nil && !IsNil(o.Severity) {
		return true
	}

	return false
}

// SetSeverity gets a reference to the given string and assigns it to the Severity field.
func (o *Raw) SetSeverity(v string) {
	o.Severity = &v
}
