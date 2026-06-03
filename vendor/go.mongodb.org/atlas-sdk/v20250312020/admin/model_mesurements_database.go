// Code based on the AtlasAPI V2 OpenAPI file

package admin

// MesurementsDatabase struct for MesurementsDatabase
type MesurementsDatabase struct {
	// Human-readable label that identifies the database that the specified MongoDB process serves.
	DatabaseName *string `json:"databaseName,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewMesurementsDatabase instantiates a new MesurementsDatabase object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMesurementsDatabase() *MesurementsDatabase {
	this := MesurementsDatabase{}
	return &this
}

// NewMesurementsDatabaseWithDefaults instantiates a new MesurementsDatabase object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMesurementsDatabaseWithDefaults() *MesurementsDatabase {
	this := MesurementsDatabase{}
	return &this
}

// GetDatabaseName returns the DatabaseName field value if set, zero value otherwise
func (o *MesurementsDatabase) GetDatabaseName() string {
	if o == nil || IsNil(o.DatabaseName) {
		var ret string
		return ret
	}
	return *o.DatabaseName
}

// GetDatabaseNameOk returns a tuple with the DatabaseName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MesurementsDatabase) GetDatabaseNameOk() (*string, bool) {
	if o == nil || IsNil(o.DatabaseName) {
		return nil, false
	}

	return o.DatabaseName, true
}

// HasDatabaseName returns a boolean if a field has been set.
func (o *MesurementsDatabase) HasDatabaseName() bool {
	if o != nil && !IsNil(o.DatabaseName) {
		return true
	}

	return false
}

// SetDatabaseName gets a reference to the given string and assigns it to the DatabaseName field.
func (o *MesurementsDatabase) SetDatabaseName(v string) {
	o.DatabaseName = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *MesurementsDatabase) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MesurementsDatabase) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *MesurementsDatabase) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *MesurementsDatabase) SetLinks(v []Link) {
	o.Links = &v
}
