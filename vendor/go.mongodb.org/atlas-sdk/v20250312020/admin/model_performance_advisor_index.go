// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PerformanceAdvisorIndex struct for PerformanceAdvisorIndex
type PerformanceAdvisorIndex struct {
	// The average size of an object in the collection of this index.
	// Read only field.
	AvgObjSize *float64 `json:"avgObjSize,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this index.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List that contains unique 24-hexadecimal character string that identifies the query shapes in this response that the Performance Advisor suggests.
	// Read only field.
	Impact *[]string `json:"impact,omitempty"`
	// List that contains documents that specify a key in the index and its sort order.
	// Read only field.
	Index *[]map[string]int `json:"index,omitempty"`
	// Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
	// Read only field.
	Namespace *string `json:"namespace,omitempty"`
	// Estimated performance improvement that the suggested index provides. This value corresponds to **Impact** in the Performance Advisor user interface.
	// Read only field.
	Weight *float64 `json:"weight,omitempty"`
}

// NewPerformanceAdvisorIndex instantiates a new PerformanceAdvisorIndex object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPerformanceAdvisorIndex() *PerformanceAdvisorIndex {
	this := PerformanceAdvisorIndex{}
	return &this
}

// NewPerformanceAdvisorIndexWithDefaults instantiates a new PerformanceAdvisorIndex object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPerformanceAdvisorIndexWithDefaults() *PerformanceAdvisorIndex {
	this := PerformanceAdvisorIndex{}
	return &this
}

// GetAvgObjSize returns the AvgObjSize field value if set, zero value otherwise
func (o *PerformanceAdvisorIndex) GetAvgObjSize() float64 {
	if o == nil || IsNil(o.AvgObjSize) {
		var ret float64
		return ret
	}
	return *o.AvgObjSize
}

// GetAvgObjSizeOk returns a tuple with the AvgObjSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorIndex) GetAvgObjSizeOk() (*float64, bool) {
	if o == nil || IsNil(o.AvgObjSize) {
		return nil, false
	}

	return o.AvgObjSize, true
}

// HasAvgObjSize returns a boolean if a field has been set.
func (o *PerformanceAdvisorIndex) HasAvgObjSize() bool {
	if o != nil && !IsNil(o.AvgObjSize) {
		return true
	}

	return false
}

// SetAvgObjSize gets a reference to the given float64 and assigns it to the AvgObjSize field.
func (o *PerformanceAdvisorIndex) SetAvgObjSize(v float64) {
	o.AvgObjSize = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *PerformanceAdvisorIndex) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorIndex) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *PerformanceAdvisorIndex) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *PerformanceAdvisorIndex) SetId(v string) {
	o.Id = &v
}

// GetImpact returns the Impact field value if set, zero value otherwise
func (o *PerformanceAdvisorIndex) GetImpact() []string {
	if o == nil || IsNil(o.Impact) {
		var ret []string
		return ret
	}
	return *o.Impact
}

// GetImpactOk returns a tuple with the Impact field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorIndex) GetImpactOk() (*[]string, bool) {
	if o == nil || IsNil(o.Impact) {
		return nil, false
	}

	return o.Impact, true
}

// HasImpact returns a boolean if a field has been set.
func (o *PerformanceAdvisorIndex) HasImpact() bool {
	if o != nil && !IsNil(o.Impact) {
		return true
	}

	return false
}

// SetImpact gets a reference to the given []string and assigns it to the Impact field.
func (o *PerformanceAdvisorIndex) SetImpact(v []string) {
	o.Impact = &v
}

// GetIndex returns the Index field value if set, zero value otherwise
func (o *PerformanceAdvisorIndex) GetIndex() []map[string]int {
	if o == nil || IsNil(o.Index) {
		var ret []map[string]int
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorIndex) GetIndexOk() (*[]map[string]int, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}

	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *PerformanceAdvisorIndex) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given []map[string]int and assigns it to the Index field.
func (o *PerformanceAdvisorIndex) SetIndex(v []map[string]int) {
	o.Index = &v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *PerformanceAdvisorIndex) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorIndex) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *PerformanceAdvisorIndex) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *PerformanceAdvisorIndex) SetNamespace(v string) {
	o.Namespace = &v
}

// GetWeight returns the Weight field value if set, zero value otherwise
func (o *PerformanceAdvisorIndex) GetWeight() float64 {
	if o == nil || IsNil(o.Weight) {
		var ret float64
		return ret
	}
	return *o.Weight
}

// GetWeightOk returns a tuple with the Weight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PerformanceAdvisorIndex) GetWeightOk() (*float64, bool) {
	if o == nil || IsNil(o.Weight) {
		return nil, false
	}

	return o.Weight, true
}

// HasWeight returns a boolean if a field has been set.
func (o *PerformanceAdvisorIndex) HasWeight() bool {
	if o != nil && !IsNil(o.Weight) {
		return true
	}

	return false
}

// SetWeight gets a reference to the given float64 and assigns it to the Weight field.
func (o *PerformanceAdvisorIndex) SetWeight(v float64) {
	o.Weight = &v
}
