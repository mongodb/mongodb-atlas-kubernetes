// Code based on the AtlasAPI V2 OpenAPI file

package admin

// QueryShapeResponse Response containing the details and status of a query shape. The query shape field may be null if the user lacks PII view access.
type QueryShapeResponse struct {
	// The MongoDB command type issued for a query shape.
	// Read only field.
	Command *string `json:"command,omitempty"`
	// Human-readable label that identifies the namespace on the specified host. The resource expresses this parameter value as `<database>.<collection>`.
	// Read only field.
	Namespace *string `json:"namespace,omitempty"`
	// A query shape is a set of specifications that group similar queries together. Specifications can include filters, sorts, projections, aggregation pipeline stages, a namespace, and others. Queries that have similar specifications have the same query shape. This field may be null if the user lacks PII view access.
	// Read only field.
	QueryShape *string `json:"queryShape,omitempty"`
	// A hexadecimal string that represents the hash of a MongoDB query shape.
	// Read only field.
	QueryShapeHash string `json:"queryShapeHash"`
	// The rejection status of a query shape. Use REJECTED to prevent the query shape from executing on the cluster, or UNREJECTED to allow it to execute.
	Status string `json:"status"`
}

// NewQueryShapeResponse instantiates a new QueryShapeResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQueryShapeResponse(queryShapeHash string, status string) *QueryShapeResponse {
	this := QueryShapeResponse{}
	this.QueryShapeHash = queryShapeHash
	this.Status = status
	return &this
}

// NewQueryShapeResponseWithDefaults instantiates a new QueryShapeResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQueryShapeResponseWithDefaults() *QueryShapeResponse {
	this := QueryShapeResponse{}
	return &this
}

// GetCommand returns the Command field value if set, zero value otherwise
func (o *QueryShapeResponse) GetCommand() string {
	if o == nil || IsNil(o.Command) {
		var ret string
		return ret
	}
	return *o.Command
}

// GetCommandOk returns a tuple with the Command field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryShapeResponse) GetCommandOk() (*string, bool) {
	if o == nil || IsNil(o.Command) {
		return nil, false
	}

	return o.Command, true
}

// HasCommand returns a boolean if a field has been set.
func (o *QueryShapeResponse) HasCommand() bool {
	if o != nil && !IsNil(o.Command) {
		return true
	}

	return false
}

// SetCommand gets a reference to the given string and assigns it to the Command field.
func (o *QueryShapeResponse) SetCommand(v string) {
	o.Command = &v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *QueryShapeResponse) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryShapeResponse) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *QueryShapeResponse) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *QueryShapeResponse) SetNamespace(v string) {
	o.Namespace = &v
}

// GetQueryShape returns the QueryShape field value if set, zero value otherwise
func (o *QueryShapeResponse) GetQueryShape() string {
	if o == nil || IsNil(o.QueryShape) {
		var ret string
		return ret
	}
	return *o.QueryShape
}

// GetQueryShapeOk returns a tuple with the QueryShape field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *QueryShapeResponse) GetQueryShapeOk() (*string, bool) {
	if o == nil || IsNil(o.QueryShape) {
		return nil, false
	}

	return o.QueryShape, true
}

// HasQueryShape returns a boolean if a field has been set.
func (o *QueryShapeResponse) HasQueryShape() bool {
	if o != nil && !IsNil(o.QueryShape) {
		return true
	}

	return false
}

// SetQueryShape gets a reference to the given string and assigns it to the QueryShape field.
func (o *QueryShapeResponse) SetQueryShape(v string) {
	o.QueryShape = &v
}

// GetQueryShapeHash returns the QueryShapeHash field value
func (o *QueryShapeResponse) GetQueryShapeHash() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.QueryShapeHash
}

// GetQueryShapeHashOk returns a tuple with the QueryShapeHash field value
// and a boolean to check if the value has been set.
func (o *QueryShapeResponse) GetQueryShapeHashOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.QueryShapeHash, true
}

// SetQueryShapeHash sets field value
func (o *QueryShapeResponse) SetQueryShapeHash(v string) {
	o.QueryShapeHash = v
}

// GetStatus returns the Status field value
func (o *QueryShapeResponse) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *QueryShapeResponse) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *QueryShapeResponse) SetStatus(v string) {
	o.Status = v
}
