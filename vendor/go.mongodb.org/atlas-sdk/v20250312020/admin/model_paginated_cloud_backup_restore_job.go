// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PaginatedCloudBackupRestoreJob struct for PaginatedCloudBackupRestoreJob
type PaginatedCloudBackupRestoreJob struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of returned documents that MongoDB Cloud provides when completing this request.
	// Read only field.
	Results []DiskBackupSnapshotRestoreJob `json:"results"`
	// Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`. The total number is an estimate and may not be exact.
	// Read only field.
	TotalCount *int `json:"totalCount,omitempty"`
}

// NewPaginatedCloudBackupRestoreJob instantiates a new PaginatedCloudBackupRestoreJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginatedCloudBackupRestoreJob(results []DiskBackupSnapshotRestoreJob) *PaginatedCloudBackupRestoreJob {
	this := PaginatedCloudBackupRestoreJob{}
	this.Results = results
	return &this
}

// NewPaginatedCloudBackupRestoreJobWithDefaults instantiates a new PaginatedCloudBackupRestoreJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginatedCloudBackupRestoreJobWithDefaults() *PaginatedCloudBackupRestoreJob {
	this := PaginatedCloudBackupRestoreJob{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *PaginatedCloudBackupRestoreJob) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedCloudBackupRestoreJob) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PaginatedCloudBackupRestoreJob) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *PaginatedCloudBackupRestoreJob) SetLinks(v []Link) {
	o.Links = &v
}

// GetResults returns the Results field value
func (o *PaginatedCloudBackupRestoreJob) GetResults() []DiskBackupSnapshotRestoreJob {
	if o == nil {
		var ret []DiskBackupSnapshotRestoreJob
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *PaginatedCloudBackupRestoreJob) GetResultsOk() (*[]DiskBackupSnapshotRestoreJob, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Results, true
}

// SetResults sets field value
func (o *PaginatedCloudBackupRestoreJob) SetResults(v []DiskBackupSnapshotRestoreJob) {
	o.Results = v
}

// GetTotalCount returns the TotalCount field value if set, zero value otherwise
func (o *PaginatedCloudBackupRestoreJob) GetTotalCount() int {
	if o == nil || IsNil(o.TotalCount) {
		var ret int
		return ret
	}
	return *o.TotalCount
}

// GetTotalCountOk returns a tuple with the TotalCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedCloudBackupRestoreJob) GetTotalCountOk() (*int, bool) {
	if o == nil || IsNil(o.TotalCount) {
		return nil, false
	}

	return o.TotalCount, true
}

// HasTotalCount returns a boolean if a field has been set.
func (o *PaginatedCloudBackupRestoreJob) HasTotalCount() bool {
	if o != nil && !IsNil(o.TotalCount) {
		return true
	}

	return false
}

// SetTotalCount gets a reference to the given int and assigns it to the TotalCount field.
func (o *PaginatedCloudBackupRestoreJob) SetTotalCount(v int) {
	o.TotalCount = &v
}
