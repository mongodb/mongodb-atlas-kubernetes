package project

import "go.mongodb.org/atlas/mongodbatlas"

// Auditing represents MongoDB Maintenance Windows
type Auditing struct {
	// Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event. For more information, see auditAuthorizationSuccess
	// +optional
	AuditAuthorizationSuccess *bool `json:"auditAuthorizationSuccess,omitempty"`
	// JSON-formatted audit filter used by the project
	// +optional
	AuditFilter string `json:"auditFilter,omitempty"`
	// Denotes whether or not the project associated with the {GROUP-ID} has database auditing enabled.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

func (a Auditing) ToAtlas() *mongodbatlas.Auditing {
	return &mongodbatlas.Auditing{
		AuditAuthorizationSuccess: a.AuditAuthorizationSuccess,
		AuditFilter:               a.AuditFilter,
		Enabled:                   a.Enabled,
	}
}
