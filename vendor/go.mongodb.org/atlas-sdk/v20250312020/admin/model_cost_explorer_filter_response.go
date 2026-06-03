// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CostExplorerFilterResponse Response object to give information about created query.
type CostExplorerFilterResponse struct {
	// The token used to identify the created Cost Explorer query.
	Token *string `json:"token,omitempty"`
}

// NewCostExplorerFilterResponse instantiates a new CostExplorerFilterResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCostExplorerFilterResponse() *CostExplorerFilterResponse {
	this := CostExplorerFilterResponse{}
	return &this
}

// NewCostExplorerFilterResponseWithDefaults instantiates a new CostExplorerFilterResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCostExplorerFilterResponseWithDefaults() *CostExplorerFilterResponse {
	this := CostExplorerFilterResponse{}
	return &this
}

// GetToken returns the Token field value if set, zero value otherwise
func (o *CostExplorerFilterResponse) GetToken() string {
	if o == nil || IsNil(o.Token) {
		var ret string
		return ret
	}
	return *o.Token
}

// GetTokenOk returns a tuple with the Token field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterResponse) GetTokenOk() (*string, bool) {
	if o == nil || IsNil(o.Token) {
		return nil, false
	}

	return o.Token, true
}

// HasToken returns a boolean if a field has been set.
func (o *CostExplorerFilterResponse) HasToken() bool {
	if o != nil && !IsNil(o.Token) {
		return true
	}

	return false
}

// SetToken gets a reference to the given string and assigns it to the Token field.
func (o *CostExplorerFilterResponse) SetToken(v string) {
	o.Token = &v
}
