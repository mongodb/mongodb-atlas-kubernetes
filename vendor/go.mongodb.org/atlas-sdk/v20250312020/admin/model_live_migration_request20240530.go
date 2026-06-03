// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LiveMigrationRequest20240530 struct for LiveMigrationRequest20240530
type LiveMigrationRequest20240530 struct {
	// Unique 24-hexadecimal digit string that identifies the migration request.
	// Read only field.
	Id          *string     `json:"_id,omitempty"`
	Destination Destination `json:"destination"`
	// Flag that indicates whether the migration process drops all collections from the destination cluster before the migration starts.
	// Write only field.
	DropDestinationData *bool `json:"dropDestinationData,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of migration hosts used for this migration.
	MigrationHosts []string         `json:"migrationHosts"`
	Sharding       *ShardingRequest `json:"sharding,omitempty"`
	Source         Source           `json:"source"`
}

// NewLiveMigrationRequest20240530 instantiates a new LiveMigrationRequest20240530 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLiveMigrationRequest20240530(destination Destination, migrationHosts []string, source Source) *LiveMigrationRequest20240530 {
	this := LiveMigrationRequest20240530{}
	this.Destination = destination
	var dropDestinationData bool = false
	this.DropDestinationData = &dropDestinationData
	this.MigrationHosts = migrationHosts
	this.Source = source
	return &this
}

// NewLiveMigrationRequest20240530WithDefaults instantiates a new LiveMigrationRequest20240530 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLiveMigrationRequest20240530WithDefaults() *LiveMigrationRequest20240530 {
	this := LiveMigrationRequest20240530{}
	var dropDestinationData bool = false
	this.DropDestinationData = &dropDestinationData
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *LiveMigrationRequest20240530) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *LiveMigrationRequest20240530) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *LiveMigrationRequest20240530) SetId(v string) {
	o.Id = &v
}

// GetDestination returns the Destination field value
func (o *LiveMigrationRequest20240530) GetDestination() Destination {
	if o == nil {
		var ret Destination
		return ret
	}

	return o.Destination
}

// GetDestinationOk returns a tuple with the Destination field value
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetDestinationOk() (*Destination, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Destination, true
}

// SetDestination sets field value
func (o *LiveMigrationRequest20240530) SetDestination(v Destination) {
	o.Destination = v
}

// GetDropDestinationData returns the DropDestinationData field value if set, zero value otherwise
func (o *LiveMigrationRequest20240530) GetDropDestinationData() bool {
	if o == nil || IsNil(o.DropDestinationData) {
		var ret bool
		return ret
	}
	return *o.DropDestinationData
}

// GetDropDestinationDataOk returns a tuple with the DropDestinationData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetDropDestinationDataOk() (*bool, bool) {
	if o == nil || IsNil(o.DropDestinationData) {
		return nil, false
	}

	return o.DropDestinationData, true
}

// HasDropDestinationData returns a boolean if a field has been set.
func (o *LiveMigrationRequest20240530) HasDropDestinationData() bool {
	if o != nil && !IsNil(o.DropDestinationData) {
		return true
	}

	return false
}

// SetDropDestinationData gets a reference to the given bool and assigns it to the DropDestinationData field.
func (o *LiveMigrationRequest20240530) SetDropDestinationData(v bool) {
	o.DropDestinationData = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *LiveMigrationRequest20240530) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *LiveMigrationRequest20240530) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *LiveMigrationRequest20240530) SetLinks(v []Link) {
	o.Links = &v
}

// GetMigrationHosts returns the MigrationHosts field value
func (o *LiveMigrationRequest20240530) GetMigrationHosts() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.MigrationHosts
}

// GetMigrationHostsOk returns a tuple with the MigrationHosts field value
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetMigrationHostsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MigrationHosts, true
}

// SetMigrationHosts sets field value
func (o *LiveMigrationRequest20240530) SetMigrationHosts(v []string) {
	o.MigrationHosts = v
}

// GetSharding returns the Sharding field value if set, zero value otherwise
func (o *LiveMigrationRequest20240530) GetSharding() ShardingRequest {
	if o == nil || IsNil(o.Sharding) {
		var ret ShardingRequest
		return ret
	}
	return *o.Sharding
}

// GetShardingOk returns a tuple with the Sharding field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetShardingOk() (*ShardingRequest, bool) {
	if o == nil || IsNil(o.Sharding) {
		return nil, false
	}

	return o.Sharding, true
}

// HasSharding returns a boolean if a field has been set.
func (o *LiveMigrationRequest20240530) HasSharding() bool {
	if o != nil && !IsNil(o.Sharding) {
		return true
	}

	return false
}

// SetSharding gets a reference to the given ShardingRequest and assigns it to the Sharding field.
func (o *LiveMigrationRequest20240530) SetSharding(v ShardingRequest) {
	o.Sharding = &v
}

// GetSource returns the Source field value
func (o *LiveMigrationRequest20240530) GetSource() Source {
	if o == nil {
		var ret Source
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *LiveMigrationRequest20240530) GetSourceOk() (*Source, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *LiveMigrationRequest20240530) SetSource(v Source) {
	o.Source = v
}
