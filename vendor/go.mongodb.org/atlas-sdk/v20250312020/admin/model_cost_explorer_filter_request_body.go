// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CostExplorerFilterRequestBody Request body for a cost explorer query.
type CostExplorerFilterRequestBody struct {
	// The list of clusters to be included in the Cost Explorer Query.
	Clusters *[]string `json:"clusters,omitempty"`
	// The exclusive ending date for the Cost Explorer query. The date must be the start of a month.
	EndDate string `json:"endDate"`
	// The dimension to group the returned usage results by. At least one filter value needs to be provided for a dimension to be used.
	GroupBy *string `json:"groupBy,omitempty"`
	// Flag to control whether usage that matches the filter criteria, but does not have values for all filter criteria is included in response. Default is false, which excludes the partially matching data.
	IncludePartialMatches *bool `json:"includePartialMatches,omitempty"`
	// The list of organizations to be included in the Cost Explorer Query.
	Organizations *[]string `json:"organizations,omitempty"`
	// The list of projects to be included in the Cost Explorer Query.
	Projects *[]string `json:"projects,omitempty"`
	// The list of SKU services to be included in the Cost Explorer Query.
	Services *[]string `json:"services,omitempty"`
	// The inclusive starting date for the Cost Explorer query. The date must be the start of a month.
	StartDate string `json:"startDate"`
}

// NewCostExplorerFilterRequestBody instantiates a new CostExplorerFilterRequestBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCostExplorerFilterRequestBody(endDate string, startDate string) *CostExplorerFilterRequestBody {
	this := CostExplorerFilterRequestBody{}
	this.EndDate = endDate
	this.StartDate = startDate
	return &this
}

// NewCostExplorerFilterRequestBodyWithDefaults instantiates a new CostExplorerFilterRequestBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCostExplorerFilterRequestBodyWithDefaults() *CostExplorerFilterRequestBody {
	this := CostExplorerFilterRequestBody{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise
func (o *CostExplorerFilterRequestBody) GetClusters() []string {
	if o == nil || IsNil(o.Clusters) {
		var ret []string
		return ret
	}
	return *o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetClustersOk() (*[]string, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}

	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *CostExplorerFilterRequestBody) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []string and assigns it to the Clusters field.
func (o *CostExplorerFilterRequestBody) SetClusters(v []string) {
	o.Clusters = &v
}

// GetEndDate returns the EndDate field value
func (o *CostExplorerFilterRequestBody) GetEndDate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EndDate
}

// GetEndDateOk returns a tuple with the EndDate field value
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetEndDateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EndDate, true
}

// SetEndDate sets field value
func (o *CostExplorerFilterRequestBody) SetEndDate(v string) {
	o.EndDate = v
}

// GetGroupBy returns the GroupBy field value if set, zero value otherwise
func (o *CostExplorerFilterRequestBody) GetGroupBy() string {
	if o == nil || IsNil(o.GroupBy) {
		var ret string
		return ret
	}
	return *o.GroupBy
}

// GetGroupByOk returns a tuple with the GroupBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetGroupByOk() (*string, bool) {
	if o == nil || IsNil(o.GroupBy) {
		return nil, false
	}

	return o.GroupBy, true
}

// HasGroupBy returns a boolean if a field has been set.
func (o *CostExplorerFilterRequestBody) HasGroupBy() bool {
	if o != nil && !IsNil(o.GroupBy) {
		return true
	}

	return false
}

// SetGroupBy gets a reference to the given string and assigns it to the GroupBy field.
func (o *CostExplorerFilterRequestBody) SetGroupBy(v string) {
	o.GroupBy = &v
}

// GetIncludePartialMatches returns the IncludePartialMatches field value if set, zero value otherwise
func (o *CostExplorerFilterRequestBody) GetIncludePartialMatches() bool {
	if o == nil || IsNil(o.IncludePartialMatches) {
		var ret bool
		return ret
	}
	return *o.IncludePartialMatches
}

// GetIncludePartialMatchesOk returns a tuple with the IncludePartialMatches field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetIncludePartialMatchesOk() (*bool, bool) {
	if o == nil || IsNil(o.IncludePartialMatches) {
		return nil, false
	}

	return o.IncludePartialMatches, true
}

// HasIncludePartialMatches returns a boolean if a field has been set.
func (o *CostExplorerFilterRequestBody) HasIncludePartialMatches() bool {
	if o != nil && !IsNil(o.IncludePartialMatches) {
		return true
	}

	return false
}

// SetIncludePartialMatches gets a reference to the given bool and assigns it to the IncludePartialMatches field.
func (o *CostExplorerFilterRequestBody) SetIncludePartialMatches(v bool) {
	o.IncludePartialMatches = &v
}

// GetOrganizations returns the Organizations field value if set, zero value otherwise
func (o *CostExplorerFilterRequestBody) GetOrganizations() []string {
	if o == nil || IsNil(o.Organizations) {
		var ret []string
		return ret
	}
	return *o.Organizations
}

// GetOrganizationsOk returns a tuple with the Organizations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetOrganizationsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Organizations) {
		return nil, false
	}

	return o.Organizations, true
}

// HasOrganizations returns a boolean if a field has been set.
func (o *CostExplorerFilterRequestBody) HasOrganizations() bool {
	if o != nil && !IsNil(o.Organizations) {
		return true
	}

	return false
}

// SetOrganizations gets a reference to the given []string and assigns it to the Organizations field.
func (o *CostExplorerFilterRequestBody) SetOrganizations(v []string) {
	o.Organizations = &v
}

// GetProjects returns the Projects field value if set, zero value otherwise
func (o *CostExplorerFilterRequestBody) GetProjects() []string {
	if o == nil || IsNil(o.Projects) {
		var ret []string
		return ret
	}
	return *o.Projects
}

// GetProjectsOk returns a tuple with the Projects field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetProjectsOk() (*[]string, bool) {
	if o == nil || IsNil(o.Projects) {
		return nil, false
	}

	return o.Projects, true
}

// HasProjects returns a boolean if a field has been set.
func (o *CostExplorerFilterRequestBody) HasProjects() bool {
	if o != nil && !IsNil(o.Projects) {
		return true
	}

	return false
}

// SetProjects gets a reference to the given []string and assigns it to the Projects field.
func (o *CostExplorerFilterRequestBody) SetProjects(v []string) {
	o.Projects = &v
}

// GetServices returns the Services field value if set, zero value otherwise
func (o *CostExplorerFilterRequestBody) GetServices() []string {
	if o == nil || IsNil(o.Services) {
		var ret []string
		return ret
	}
	return *o.Services
}

// GetServicesOk returns a tuple with the Services field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetServicesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Services) {
		return nil, false
	}

	return o.Services, true
}

// HasServices returns a boolean if a field has been set.
func (o *CostExplorerFilterRequestBody) HasServices() bool {
	if o != nil && !IsNil(o.Services) {
		return true
	}

	return false
}

// SetServices gets a reference to the given []string and assigns it to the Services field.
func (o *CostExplorerFilterRequestBody) SetServices(v []string) {
	o.Services = &v
}

// GetStartDate returns the StartDate field value
func (o *CostExplorerFilterRequestBody) GetStartDate() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.StartDate
}

// GetStartDateOk returns a tuple with the StartDate field value
// and a boolean to check if the value has been set.
func (o *CostExplorerFilterRequestBody) GetStartDateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StartDate, true
}

// SetStartDate sets field value
func (o *CostExplorerFilterRequestBody) SetStartDate(v string) {
	o.StartDate = v
}
