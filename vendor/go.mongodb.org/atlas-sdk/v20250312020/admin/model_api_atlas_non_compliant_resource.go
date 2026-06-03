// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasNonCompliantResource struct for ApiAtlasNonCompliantResource
type ApiAtlasNonCompliantResource struct {
	// Unique 24-hexadecimal character string that identifies the organization the resource belongs to.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// Unique 24-hexadecimal character string that identifies the non-compliant resource.
	// Read only field.
	ResourceId *string `json:"resourceId,omitempty"`
	// Unique human readable string that identifies the non-compliant resource.
	// Read only field.
	ResourceName *string `json:"resourceName,omitempty"`
	// List of resource policies causing the resource to be considered non-compliant.
	// Read only field.
	ResourcePoliciesCausingNonCompliance *[]ApiAtlasResourcePolicyMetadata `json:"resourcePoliciesCausingNonCompliance,omitempty"`
	// Human-readable label that displays the type of a resource.
	// Read only field.
	ResourceType *string `json:"resourceType,omitempty"`
}

// NewApiAtlasNonCompliantResource instantiates a new ApiAtlasNonCompliantResource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasNonCompliantResource() *ApiAtlasNonCompliantResource {
	this := ApiAtlasNonCompliantResource{}
	return &this
}

// NewApiAtlasNonCompliantResourceWithDefaults instantiates a new ApiAtlasNonCompliantResource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasNonCompliantResourceWithDefaults() *ApiAtlasNonCompliantResource {
	this := ApiAtlasNonCompliantResource{}
	return &this
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *ApiAtlasNonCompliantResource) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasNonCompliantResource) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *ApiAtlasNonCompliantResource) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *ApiAtlasNonCompliantResource) SetOrgId(v string) {
	o.OrgId = &v
}

// GetResourceId returns the ResourceId field value if set, zero value otherwise
func (o *ApiAtlasNonCompliantResource) GetResourceId() string {
	if o == nil || IsNil(o.ResourceId) {
		var ret string
		return ret
	}
	return *o.ResourceId
}

// GetResourceIdOk returns a tuple with the ResourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasNonCompliantResource) GetResourceIdOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceId) {
		return nil, false
	}

	return o.ResourceId, true
}

// HasResourceId returns a boolean if a field has been set.
func (o *ApiAtlasNonCompliantResource) HasResourceId() bool {
	if o != nil && !IsNil(o.ResourceId) {
		return true
	}

	return false
}

// SetResourceId gets a reference to the given string and assigns it to the ResourceId field.
func (o *ApiAtlasNonCompliantResource) SetResourceId(v string) {
	o.ResourceId = &v
}

// GetResourceName returns the ResourceName field value if set, zero value otherwise
func (o *ApiAtlasNonCompliantResource) GetResourceName() string {
	if o == nil || IsNil(o.ResourceName) {
		var ret string
		return ret
	}
	return *o.ResourceName
}

// GetResourceNameOk returns a tuple with the ResourceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasNonCompliantResource) GetResourceNameOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceName) {
		return nil, false
	}

	return o.ResourceName, true
}

// HasResourceName returns a boolean if a field has been set.
func (o *ApiAtlasNonCompliantResource) HasResourceName() bool {
	if o != nil && !IsNil(o.ResourceName) {
		return true
	}

	return false
}

// SetResourceName gets a reference to the given string and assigns it to the ResourceName field.
func (o *ApiAtlasNonCompliantResource) SetResourceName(v string) {
	o.ResourceName = &v
}

// GetResourcePoliciesCausingNonCompliance returns the ResourcePoliciesCausingNonCompliance field value if set, zero value otherwise
func (o *ApiAtlasNonCompliantResource) GetResourcePoliciesCausingNonCompliance() []ApiAtlasResourcePolicyMetadata {
	if o == nil || IsNil(o.ResourcePoliciesCausingNonCompliance) {
		var ret []ApiAtlasResourcePolicyMetadata
		return ret
	}
	return *o.ResourcePoliciesCausingNonCompliance
}

// GetResourcePoliciesCausingNonComplianceOk returns a tuple with the ResourcePoliciesCausingNonCompliance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasNonCompliantResource) GetResourcePoliciesCausingNonComplianceOk() (*[]ApiAtlasResourcePolicyMetadata, bool) {
	if o == nil || IsNil(o.ResourcePoliciesCausingNonCompliance) {
		return nil, false
	}

	return o.ResourcePoliciesCausingNonCompliance, true
}

// HasResourcePoliciesCausingNonCompliance returns a boolean if a field has been set.
func (o *ApiAtlasNonCompliantResource) HasResourcePoliciesCausingNonCompliance() bool {
	if o != nil && !IsNil(o.ResourcePoliciesCausingNonCompliance) {
		return true
	}

	return false
}

// SetResourcePoliciesCausingNonCompliance gets a reference to the given []ApiAtlasResourcePolicyMetadata and assigns it to the ResourcePoliciesCausingNonCompliance field.
func (o *ApiAtlasNonCompliantResource) SetResourcePoliciesCausingNonCompliance(v []ApiAtlasResourcePolicyMetadata) {
	o.ResourcePoliciesCausingNonCompliance = &v
}

// GetResourceType returns the ResourceType field value if set, zero value otherwise
func (o *ApiAtlasNonCompliantResource) GetResourceType() string {
	if o == nil || IsNil(o.ResourceType) {
		var ret string
		return ret
	}
	return *o.ResourceType
}

// GetResourceTypeOk returns a tuple with the ResourceType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasNonCompliantResource) GetResourceTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceType) {
		return nil, false
	}

	return o.ResourceType, true
}

// HasResourceType returns a boolean if a field has been set.
func (o *ApiAtlasNonCompliantResource) HasResourceType() bool {
	if o != nil && !IsNil(o.ResourceType) {
		return true
	}

	return false
}

// SetResourceType gets a reference to the given string and assigns it to the ResourceType field.
func (o *ApiAtlasNonCompliantResource) SetResourceType(v string) {
	o.ResourceType = &v
}
