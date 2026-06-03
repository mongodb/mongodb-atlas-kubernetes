// Code based on the AtlasAPI V2 OpenAPI file

package admin

// RestoreJobFileHash Key and value pair that map one restore file to one hashed checksum. This parameter applies after you download the corresponding `delivery.url`.
type RestoreJobFileHash struct {
	// Human-readable label that identifies the hashed file.
	// Read only field.
	FileName *string `json:"fileName,omitempty"`
	// Hashed checksum that maps to the restore file.
	// Read only field.
	Hash *string `json:"hash,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the hashing algorithm used to compute the hash value.
	// Read only field.
	TypeName *string `json:"typeName,omitempty"`
}

// NewRestoreJobFileHash instantiates a new RestoreJobFileHash object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRestoreJobFileHash() *RestoreJobFileHash {
	this := RestoreJobFileHash{}
	return &this
}

// NewRestoreJobFileHashWithDefaults instantiates a new RestoreJobFileHash object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRestoreJobFileHashWithDefaults() *RestoreJobFileHash {
	this := RestoreJobFileHash{}
	return &this
}

// GetFileName returns the FileName field value if set, zero value otherwise
func (o *RestoreJobFileHash) GetFileName() string {
	if o == nil || IsNil(o.FileName) {
		var ret string
		return ret
	}
	return *o.FileName
}

// GetFileNameOk returns a tuple with the FileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RestoreJobFileHash) GetFileNameOk() (*string, bool) {
	if o == nil || IsNil(o.FileName) {
		return nil, false
	}

	return o.FileName, true
}

// HasFileName returns a boolean if a field has been set.
func (o *RestoreJobFileHash) HasFileName() bool {
	if o != nil && !IsNil(o.FileName) {
		return true
	}

	return false
}

// SetFileName gets a reference to the given string and assigns it to the FileName field.
func (o *RestoreJobFileHash) SetFileName(v string) {
	o.FileName = &v
}

// GetHash returns the Hash field value if set, zero value otherwise
func (o *RestoreJobFileHash) GetHash() string {
	if o == nil || IsNil(o.Hash) {
		var ret string
		return ret
	}
	return *o.Hash
}

// GetHashOk returns a tuple with the Hash field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RestoreJobFileHash) GetHashOk() (*string, bool) {
	if o == nil || IsNil(o.Hash) {
		return nil, false
	}

	return o.Hash, true
}

// HasHash returns a boolean if a field has been set.
func (o *RestoreJobFileHash) HasHash() bool {
	if o != nil && !IsNil(o.Hash) {
		return true
	}

	return false
}

// SetHash gets a reference to the given string and assigns it to the Hash field.
func (o *RestoreJobFileHash) SetHash(v string) {
	o.Hash = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *RestoreJobFileHash) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RestoreJobFileHash) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *RestoreJobFileHash) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *RestoreJobFileHash) SetLinks(v []Link) {
	o.Links = &v
}

// GetTypeName returns the TypeName field value if set, zero value otherwise
func (o *RestoreJobFileHash) GetTypeName() string {
	if o == nil || IsNil(o.TypeName) {
		var ret string
		return ret
	}
	return *o.TypeName
}

// GetTypeNameOk returns a tuple with the TypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *RestoreJobFileHash) GetTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.TypeName) {
		return nil, false
	}

	return o.TypeName, true
}

// HasTypeName returns a boolean if a field has been set.
func (o *RestoreJobFileHash) HasTypeName() bool {
	if o != nil && !IsNil(o.TypeName) {
		return true
	}

	return false
}

// SetTypeName gets a reference to the given string and assigns it to the TypeName field.
func (o *RestoreJobFileHash) SetTypeName(v string) {
	o.TypeName = &v
}
