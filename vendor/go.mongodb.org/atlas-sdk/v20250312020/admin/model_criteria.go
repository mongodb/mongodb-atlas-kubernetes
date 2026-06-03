// Code based on the AtlasAPI V2 OpenAPI file

package admin

// Criteria Rules by which MongoDB Cloud archives data.  Use the `criteria.type` field to choose how MongoDB Cloud selects data to archive. Choose data using the age of the data or a MongoDB query. `\"criteria.type\": \"DATE\"` selects documents to archive based on a date. `\"criteria.type\": \"CUSTOM\"` selects documents to archive based on a custom JSON query. MongoDB Cloud doesn't support `\"criteria.type\": \"CUSTOM\"` when `\"collectionType\": \"TIMESERIES\"`.
type Criteria struct {
	// Means by which MongoDB Cloud selects data to archive. Data can be chosen using the age of the data or a MongoDB query. `DATE` selects documents to archive based on a date. `CUSTOM` selects documents to archive based on a custom JSON query. MongoDB Cloud doesn't support `CUSTOM` when `\"collectionType\": \"TIMESERIES\"`.
	Type *string `json:"type,omitempty"`
	// MongoDB find query that selects documents to archive. The specified query follows the syntax of the `db.collection.find(query)` command. This query can't use the empty document (`{}`) to return all documents. Set this parameter when `\"criteria.type\" : \"CUSTOM\"`.
	Query *string `json:"query,omitempty"`
	// Indexed database parameter that stores the date that determines when data moves to the online archive. MongoDB Cloud archives the data when the current date exceeds the date in this database parameter plus the number of days specified through the `expireAfterDays` parameter. Set this parameter when you set `\"criteria.type\" : \"DATE\"`.
	DateField *string `json:"dateField,omitempty"`
	// Syntax used to write the date after which data moves to the online archive. Date can be expressed as ISO 8601, Epoch timestamps, or Object ID. The Epoch timestamp can be expressed as nanoseconds, milliseconds, or seconds. Set this parameter when `criteria.type` : `DATE`. You must set `criteria.type` : `DATE` if `collectionType`: `TIMESERIES`.
	DateFormat *string `json:"dateFormat,omitempty"`
	// Number of days after the value in the `criteria.dateField` when MongoDB Cloud archives data in the specified cluster. Set this parameter when you set `\"criteria.type\" : \"DATE\"`.
	ExpireAfterDays *int `json:"expireAfterDays,omitempty"`
}

// NewCriteria instantiates a new Criteria object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCriteria() *Criteria {
	this := Criteria{}
	var dateFormat string = "ISODATE"
	this.DateFormat = &dateFormat
	return &this
}

// NewCriteriaWithDefaults instantiates a new Criteria object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCriteriaWithDefaults() *Criteria {
	this := Criteria{}
	var dateFormat string = "ISODATE"
	this.DateFormat = &dateFormat
	return &this
}

// GetType returns the Type field value if set, zero value otherwise
func (o *Criteria) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Criteria) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *Criteria) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *Criteria) SetType(v string) {
	o.Type = &v
}

// GetQuery returns the Query field value if set, zero value otherwise
func (o *Criteria) GetQuery() string {
	if o == nil || IsNil(o.Query) {
		var ret string
		return ret
	}
	return *o.Query
}

// GetQueryOk returns a tuple with the Query field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Criteria) GetQueryOk() (*string, bool) {
	if o == nil || IsNil(o.Query) {
		return nil, false
	}

	return o.Query, true
}

// HasQuery returns a boolean if a field has been set.
func (o *Criteria) HasQuery() bool {
	if o != nil && !IsNil(o.Query) {
		return true
	}

	return false
}

// SetQuery gets a reference to the given string and assigns it to the Query field.
func (o *Criteria) SetQuery(v string) {
	o.Query = &v
}

// GetDateField returns the DateField field value if set, zero value otherwise
func (o *Criteria) GetDateField() string {
	if o == nil || IsNil(o.DateField) {
		var ret string
		return ret
	}
	return *o.DateField
}

// GetDateFieldOk returns a tuple with the DateField field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Criteria) GetDateFieldOk() (*string, bool) {
	if o == nil || IsNil(o.DateField) {
		return nil, false
	}

	return o.DateField, true
}

// HasDateField returns a boolean if a field has been set.
func (o *Criteria) HasDateField() bool {
	if o != nil && !IsNil(o.DateField) {
		return true
	}

	return false
}

// SetDateField gets a reference to the given string and assigns it to the DateField field.
func (o *Criteria) SetDateField(v string) {
	o.DateField = &v
}

// GetDateFormat returns the DateFormat field value if set, zero value otherwise
func (o *Criteria) GetDateFormat() string {
	if o == nil || IsNil(o.DateFormat) {
		var ret string
		return ret
	}
	return *o.DateFormat
}

// GetDateFormatOk returns a tuple with the DateFormat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Criteria) GetDateFormatOk() (*string, bool) {
	if o == nil || IsNil(o.DateFormat) {
		return nil, false
	}

	return o.DateFormat, true
}

// HasDateFormat returns a boolean if a field has been set.
func (o *Criteria) HasDateFormat() bool {
	if o != nil && !IsNil(o.DateFormat) {
		return true
	}

	return false
}

// SetDateFormat gets a reference to the given string and assigns it to the DateFormat field.
func (o *Criteria) SetDateFormat(v string) {
	o.DateFormat = &v
}

// GetExpireAfterDays returns the ExpireAfterDays field value if set, zero value otherwise
func (o *Criteria) GetExpireAfterDays() int {
	if o == nil || IsNil(o.ExpireAfterDays) {
		var ret int
		return ret
	}
	return *o.ExpireAfterDays
}

// GetExpireAfterDaysOk returns a tuple with the ExpireAfterDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Criteria) GetExpireAfterDaysOk() (*int, bool) {
	if o == nil || IsNil(o.ExpireAfterDays) {
		return nil, false
	}

	return o.ExpireAfterDays, true
}

// HasExpireAfterDays returns a boolean if a field has been set.
func (o *Criteria) HasExpireAfterDays() bool {
	if o != nil && !IsNil(o.ExpireAfterDays) {
		return true
	}

	return false
}

// SetExpireAfterDays gets a reference to the given int and assigns it to the ExpireAfterDays field.
func (o *Criteria) SetExpireAfterDays(v int) {
	o.ExpireAfterDays = &v
}
