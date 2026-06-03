// Code based on the AtlasAPI V2 OpenAPI file

package admin

// MongoDBAccessLogsList struct for MongoDBAccessLogsList
type MongoDBAccessLogsList struct {
	// Authentication attempt, one per object, made against the cluster.
	// Read only field.
	AccessLogs *[]MongoDBAccessLogs `json:"accessLogs,omitempty"`
}

// NewMongoDBAccessLogsList instantiates a new MongoDBAccessLogsList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMongoDBAccessLogsList() *MongoDBAccessLogsList {
	this := MongoDBAccessLogsList{}
	return &this
}

// NewMongoDBAccessLogsListWithDefaults instantiates a new MongoDBAccessLogsList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMongoDBAccessLogsListWithDefaults() *MongoDBAccessLogsList {
	this := MongoDBAccessLogsList{}
	return &this
}

// GetAccessLogs returns the AccessLogs field value if set, zero value otherwise
func (o *MongoDBAccessLogsList) GetAccessLogs() []MongoDBAccessLogs {
	if o == nil || IsNil(o.AccessLogs) {
		var ret []MongoDBAccessLogs
		return ret
	}
	return *o.AccessLogs
}

// GetAccessLogsOk returns a tuple with the AccessLogs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MongoDBAccessLogsList) GetAccessLogsOk() (*[]MongoDBAccessLogs, bool) {
	if o == nil || IsNil(o.AccessLogs) {
		return nil, false
	}

	return o.AccessLogs, true
}

// HasAccessLogs returns a boolean if a field has been set.
func (o *MongoDBAccessLogsList) HasAccessLogs() bool {
	if o != nil && !IsNil(o.AccessLogs) {
		return true
	}

	return false
}

// SetAccessLogs gets a reference to the given []MongoDBAccessLogs and assigns it to the AccessLogs field.
func (o *MongoDBAccessLogsList) SetAccessLogs(v []MongoDBAccessLogs) {
	o.AccessLogs = &v
}
