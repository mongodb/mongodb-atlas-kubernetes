// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PemFileInfo PEM file information for the identity provider's current certificates.
type PemFileInfo struct {
	// List of certificates in the file.
	Certificates *[]X509Certificate `json:"certificates,omitempty"`
	// Human-readable label given to the file.
	FileName *string `json:"fileName,omitempty"`
}

// NewPemFileInfo instantiates a new PemFileInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPemFileInfo() *PemFileInfo {
	this := PemFileInfo{}
	return &this
}

// NewPemFileInfoWithDefaults instantiates a new PemFileInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPemFileInfoWithDefaults() *PemFileInfo {
	this := PemFileInfo{}
	return &this
}

// GetCertificates returns the Certificates field value if set, zero value otherwise
func (o *PemFileInfo) GetCertificates() []X509Certificate {
	if o == nil || IsNil(o.Certificates) {
		var ret []X509Certificate
		return ret
	}
	return *o.Certificates
}

// GetCertificatesOk returns a tuple with the Certificates field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PemFileInfo) GetCertificatesOk() (*[]X509Certificate, bool) {
	if o == nil || IsNil(o.Certificates) {
		return nil, false
	}

	return o.Certificates, true
}

// HasCertificates returns a boolean if a field has been set.
func (o *PemFileInfo) HasCertificates() bool {
	if o != nil && !IsNil(o.Certificates) {
		return true
	}

	return false
}

// SetCertificates gets a reference to the given []X509Certificate and assigns it to the Certificates field.
func (o *PemFileInfo) SetCertificates(v []X509Certificate) {
	o.Certificates = &v
}

// GetFileName returns the FileName field value if set, zero value otherwise
func (o *PemFileInfo) GetFileName() string {
	if o == nil || IsNil(o.FileName) {
		var ret string
		return ret
	}
	return *o.FileName
}

// GetFileNameOk returns a tuple with the FileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PemFileInfo) GetFileNameOk() (*string, bool) {
	if o == nil || IsNil(o.FileName) {
		return nil, false
	}

	return o.FileName, true
}

// HasFileName returns a boolean if a field has been set.
func (o *PemFileInfo) HasFileName() bool {
	if o != nil && !IsNil(o.FileName) {
		return true
	}

	return false
}

// SetFileName gets a reference to the given string and assigns it to the FileName field.
func (o *PemFileInfo) SetFileName(v string) {
	o.FileName = &v
}
