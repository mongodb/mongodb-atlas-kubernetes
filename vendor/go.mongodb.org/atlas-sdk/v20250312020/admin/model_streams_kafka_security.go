// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsKafkaSecurity Properties for the secure transport connection to Kafka. For SSL, this can include the trusted certificate to use.
type StreamsKafkaSecurity struct {
	// A trusted, public x509 certificate for connecting to Kafka over SSL.
	BrokerPublicCertificate *string `json:"brokerPublicCertificate,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Describes the transport type. Can be either `SASL_PLAINTEXT`, `SASL_SSL`, or `SSL`.
	Protocol *string `json:"protocol,omitempty"`
}

// NewStreamsKafkaSecurity instantiates a new StreamsKafkaSecurity object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsKafkaSecurity() *StreamsKafkaSecurity {
	this := StreamsKafkaSecurity{}
	return &this
}

// NewStreamsKafkaSecurityWithDefaults instantiates a new StreamsKafkaSecurity object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsKafkaSecurityWithDefaults() *StreamsKafkaSecurity {
	this := StreamsKafkaSecurity{}
	return &this
}

// GetBrokerPublicCertificate returns the BrokerPublicCertificate field value if set, zero value otherwise
func (o *StreamsKafkaSecurity) GetBrokerPublicCertificate() string {
	if o == nil || IsNil(o.BrokerPublicCertificate) {
		var ret string
		return ret
	}
	return *o.BrokerPublicCertificate
}

// GetBrokerPublicCertificateOk returns a tuple with the BrokerPublicCertificate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaSecurity) GetBrokerPublicCertificateOk() (*string, bool) {
	if o == nil || IsNil(o.BrokerPublicCertificate) {
		return nil, false
	}

	return o.BrokerPublicCertificate, true
}

// HasBrokerPublicCertificate returns a boolean if a field has been set.
func (o *StreamsKafkaSecurity) HasBrokerPublicCertificate() bool {
	if o != nil && !IsNil(o.BrokerPublicCertificate) {
		return true
	}

	return false
}

// SetBrokerPublicCertificate gets a reference to the given string and assigns it to the BrokerPublicCertificate field.
func (o *StreamsKafkaSecurity) SetBrokerPublicCertificate(v string) {
	o.BrokerPublicCertificate = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsKafkaSecurity) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaSecurity) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsKafkaSecurity) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsKafkaSecurity) SetLinks(v []Link) {
	o.Links = &v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise
func (o *StreamsKafkaSecurity) GetProtocol() string {
	if o == nil || IsNil(o.Protocol) {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsKafkaSecurity) GetProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}

	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *StreamsKafkaSecurity) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *StreamsKafkaSecurity) SetProtocol(v string) {
	o.Protocol = &v
}
