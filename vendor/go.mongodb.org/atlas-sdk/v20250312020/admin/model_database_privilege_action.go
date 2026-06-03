// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DatabasePrivilegeAction Privilege action that the role grants.
type DatabasePrivilegeAction struct {
	// Human-readable label that identifies the privilege action.
	Action string `json:"action"`
	// List of resources on which you grant the action.
	Resources []DatabasePermittedNamespaceResource `json:"resources"`
}

// NewDatabasePrivilegeAction instantiates a new DatabasePrivilegeAction object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatabasePrivilegeAction(action string, resources []DatabasePermittedNamespaceResource) *DatabasePrivilegeAction {
	this := DatabasePrivilegeAction{}
	this.Action = action
	this.Resources = resources
	return &this
}

// NewDatabasePrivilegeActionWithDefaults instantiates a new DatabasePrivilegeAction object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatabasePrivilegeActionWithDefaults() *DatabasePrivilegeAction {
	this := DatabasePrivilegeAction{}
	return &this
}

// GetAction returns the Action field value
func (o *DatabasePrivilegeAction) GetAction() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Action
}

// GetActionOk returns a tuple with the Action field value
// and a boolean to check if the value has been set.
func (o *DatabasePrivilegeAction) GetActionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Action, true
}

// SetAction sets field value
func (o *DatabasePrivilegeAction) SetAction(v string) {
	o.Action = v
}

// GetResources returns the Resources field value
func (o *DatabasePrivilegeAction) GetResources() []DatabasePermittedNamespaceResource {
	if o == nil {
		var ret []DatabasePermittedNamespaceResource
		return ret
	}

	return o.Resources
}

// GetResourcesOk returns a tuple with the Resources field value
// and a boolean to check if the value has been set.
func (o *DatabasePrivilegeAction) GetResourcesOk() (*[]DatabasePermittedNamespaceResource, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Resources, true
}

// SetResources sets field value
func (o *DatabasePrivilegeAction) SetResources(v []DatabasePermittedNamespaceResource) {
	o.Resources = v
}
