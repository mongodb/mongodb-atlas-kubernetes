// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ActivityFeedLinkResponse Response containing a shareable activity feed link.
type ActivityFeedLinkResponse struct {
	// Shareable link to the activity feed with pre-applied filters.
	Link string `json:"link"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewActivityFeedLinkResponse instantiates a new ActivityFeedLinkResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewActivityFeedLinkResponse(link string) *ActivityFeedLinkResponse {
	this := ActivityFeedLinkResponse{}
	this.Link = link
	return &this
}

// NewActivityFeedLinkResponseWithDefaults instantiates a new ActivityFeedLinkResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewActivityFeedLinkResponseWithDefaults() *ActivityFeedLinkResponse {
	this := ActivityFeedLinkResponse{}
	return &this
}

// GetLink returns the Link field value
func (o *ActivityFeedLinkResponse) GetLink() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Link
}

// GetLinkOk returns a tuple with the Link field value
// and a boolean to check if the value has been set.
func (o *ActivityFeedLinkResponse) GetLinkOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Link, true
}

// SetLink sets field value
func (o *ActivityFeedLinkResponse) SetLink(v string) {
	o.Link = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ActivityFeedLinkResponse) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ActivityFeedLinkResponse) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ActivityFeedLinkResponse) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ActivityFeedLinkResponse) SetLinks(v []Link) {
	o.Links = &v
}
