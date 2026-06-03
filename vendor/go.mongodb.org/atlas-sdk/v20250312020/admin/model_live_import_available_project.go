// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LiveImportAvailableProject struct for LiveImportAvailableProject
type LiveImportAvailableProject struct {
	// List of clusters that can be migrated to MongoDB Cloud.
	Deployments []AvailableClustersDeployment `json:"deployments"`
	// Hostname of MongoDB Agent list that you configured to perform a migration.
	MigrationHosts []string `json:"migrationHosts"`
	// Human-readable label that identifies this project.
	// Read only field.
	Name string `json:"name"`
	// Unique 24-hexadecimal digit string that identifies the project to be migrated.
	// Read only field.
	ProjectId string `json:"projectId"`
}

// NewLiveImportAvailableProject instantiates a new LiveImportAvailableProject object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLiveImportAvailableProject(deployments []AvailableClustersDeployment, migrationHosts []string, name string, projectId string) *LiveImportAvailableProject {
	this := LiveImportAvailableProject{}
	this.Deployments = deployments
	this.MigrationHosts = migrationHosts
	this.Name = name
	this.ProjectId = projectId
	return &this
}

// NewLiveImportAvailableProjectWithDefaults instantiates a new LiveImportAvailableProject object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLiveImportAvailableProjectWithDefaults() *LiveImportAvailableProject {
	this := LiveImportAvailableProject{}
	return &this
}

// GetDeployments returns the Deployments field value
func (o *LiveImportAvailableProject) GetDeployments() []AvailableClustersDeployment {
	if o == nil {
		var ret []AvailableClustersDeployment
		return ret
	}

	return o.Deployments
}

// GetDeploymentsOk returns a tuple with the Deployments field value
// and a boolean to check if the value has been set.
func (o *LiveImportAvailableProject) GetDeploymentsOk() (*[]AvailableClustersDeployment, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Deployments, true
}

// SetDeployments sets field value
func (o *LiveImportAvailableProject) SetDeployments(v []AvailableClustersDeployment) {
	o.Deployments = v
}

// GetMigrationHosts returns the MigrationHosts field value
func (o *LiveImportAvailableProject) GetMigrationHosts() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.MigrationHosts
}

// GetMigrationHostsOk returns a tuple with the MigrationHosts field value
// and a boolean to check if the value has been set.
func (o *LiveImportAvailableProject) GetMigrationHostsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MigrationHosts, true
}

// SetMigrationHosts sets field value
func (o *LiveImportAvailableProject) SetMigrationHosts(v []string) {
	o.MigrationHosts = v
}

// GetName returns the Name field value
func (o *LiveImportAvailableProject) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *LiveImportAvailableProject) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *LiveImportAvailableProject) SetName(v string) {
	o.Name = v
}

// GetProjectId returns the ProjectId field value
func (o *LiveImportAvailableProject) GetProjectId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value
// and a boolean to check if the value has been set.
func (o *LiveImportAvailableProject) GetProjectIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProjectId, true
}

// SetProjectId sets field value
func (o *LiveImportAvailableProject) SetProjectId(v string) {
	o.ProjectId = v
}
