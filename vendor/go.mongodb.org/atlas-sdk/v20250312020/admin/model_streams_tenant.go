// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsTenant struct for StreamsTenant
type StreamsTenant struct {
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// List of connections configured in the stream workspace.
	// Read only field.
	Connections       *[]StreamsConnection      `json:"connections,omitempty"`
	DataProcessRegion *StreamsDataProcessRegion `json:"dataProcessRegion,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List that contains the hostnames assigned to the stream workspace.
	// Read only field.
	Hostnames *[]string `json:"hostnames,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Label that identifies the stream workspace.
	Name              *string                   `json:"name,omitempty"`
	SampleConnections *StreamsSampleConnections `json:"sampleConnections,omitempty"`
	StreamConfig      *StreamConfig             `json:"streamConfig,omitempty"`
}

// NewStreamsTenant instantiates a new StreamsTenant object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsTenant() *StreamsTenant {
	this := StreamsTenant{}
	return &this
}

// NewStreamsTenantWithDefaults instantiates a new StreamsTenant object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsTenantWithDefaults() *StreamsTenant {
	this := StreamsTenant{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *StreamsTenant) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *StreamsTenant) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *StreamsTenant) SetId(v string) {
	o.Id = &v
}

// GetConnections returns the Connections field value if set, zero value otherwise
func (o *StreamsTenant) GetConnections() []StreamsConnection {
	if o == nil || IsNil(o.Connections) {
		var ret []StreamsConnection
		return ret
	}
	return *o.Connections
}

// GetConnectionsOk returns a tuple with the Connections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetConnectionsOk() (*[]StreamsConnection, bool) {
	if o == nil || IsNil(o.Connections) {
		return nil, false
	}

	return o.Connections, true
}

// HasConnections returns a boolean if a field has been set.
func (o *StreamsTenant) HasConnections() bool {
	if o != nil && !IsNil(o.Connections) {
		return true
	}

	return false
}

// SetConnections gets a reference to the given []StreamsConnection and assigns it to the Connections field.
func (o *StreamsTenant) SetConnections(v []StreamsConnection) {
	o.Connections = &v
}

// GetDataProcessRegion returns the DataProcessRegion field value if set, zero value otherwise
func (o *StreamsTenant) GetDataProcessRegion() StreamsDataProcessRegion {
	if o == nil || IsNil(o.DataProcessRegion) {
		var ret StreamsDataProcessRegion
		return ret
	}
	return *o.DataProcessRegion
}

// GetDataProcessRegionOk returns a tuple with the DataProcessRegion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetDataProcessRegionOk() (*StreamsDataProcessRegion, bool) {
	if o == nil || IsNil(o.DataProcessRegion) {
		return nil, false
	}

	return o.DataProcessRegion, true
}

// HasDataProcessRegion returns a boolean if a field has been set.
func (o *StreamsTenant) HasDataProcessRegion() bool {
	if o != nil && !IsNil(o.DataProcessRegion) {
		return true
	}

	return false
}

// SetDataProcessRegion gets a reference to the given StreamsDataProcessRegion and assigns it to the DataProcessRegion field.
func (o *StreamsTenant) SetDataProcessRegion(v StreamsDataProcessRegion) {
	o.DataProcessRegion = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *StreamsTenant) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *StreamsTenant) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *StreamsTenant) SetGroupId(v string) {
	o.GroupId = &v
}

// GetHostnames returns the Hostnames field value if set, zero value otherwise
func (o *StreamsTenant) GetHostnames() []string {
	if o == nil || IsNil(o.Hostnames) {
		var ret []string
		return ret
	}
	return *o.Hostnames
}

// GetHostnamesOk returns a tuple with the Hostnames field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetHostnamesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Hostnames) {
		return nil, false
	}

	return o.Hostnames, true
}

// HasHostnames returns a boolean if a field has been set.
func (o *StreamsTenant) HasHostnames() bool {
	if o != nil && !IsNil(o.Hostnames) {
		return true
	}

	return false
}

// SetHostnames gets a reference to the given []string and assigns it to the Hostnames field.
func (o *StreamsTenant) SetHostnames(v []string) {
	o.Hostnames = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsTenant) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsTenant) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsTenant) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *StreamsTenant) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *StreamsTenant) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *StreamsTenant) SetName(v string) {
	o.Name = &v
}

// GetSampleConnections returns the SampleConnections field value if set, zero value otherwise
func (o *StreamsTenant) GetSampleConnections() StreamsSampleConnections {
	if o == nil || IsNil(o.SampleConnections) {
		var ret StreamsSampleConnections
		return ret
	}
	return *o.SampleConnections
}

// GetSampleConnectionsOk returns a tuple with the SampleConnections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetSampleConnectionsOk() (*StreamsSampleConnections, bool) {
	if o == nil || IsNil(o.SampleConnections) {
		return nil, false
	}

	return o.SampleConnections, true
}

// HasSampleConnections returns a boolean if a field has been set.
func (o *StreamsTenant) HasSampleConnections() bool {
	if o != nil && !IsNil(o.SampleConnections) {
		return true
	}

	return false
}

// SetSampleConnections gets a reference to the given StreamsSampleConnections and assigns it to the SampleConnections field.
func (o *StreamsTenant) SetSampleConnections(v StreamsSampleConnections) {
	o.SampleConnections = &v
}

// GetStreamConfig returns the StreamConfig field value if set, zero value otherwise
func (o *StreamsTenant) GetStreamConfig() StreamConfig {
	if o == nil || IsNil(o.StreamConfig) {
		var ret StreamConfig
		return ret
	}
	return *o.StreamConfig
}

// GetStreamConfigOk returns a tuple with the StreamConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenant) GetStreamConfigOk() (*StreamConfig, bool) {
	if o == nil || IsNil(o.StreamConfig) {
		return nil, false
	}

	return o.StreamConfig, true
}

// HasStreamConfig returns a boolean if a field has been set.
func (o *StreamsTenant) HasStreamConfig() bool {
	if o != nil && !IsNil(o.StreamConfig) {
		return true
	}

	return false
}

// SetStreamConfig gets a reference to the given StreamConfig and assigns it to the StreamConfig field.
func (o *StreamsTenant) SetStreamConfig(v StreamConfig) {
	o.StreamConfig = &v
}
