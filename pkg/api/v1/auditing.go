package v1

import (
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

// Auditing represents MongoDB Maintenance Windows
type Auditing struct {
	// Indicates whether the auditing system captures successful authentication attempts for audit filters using the "atype" : "authCheck" auditing event. For more information, see auditAuthorizationSuccess
	// +optional
	AuditAuthorizationSuccess bool `json:"auditAuthorizationSuccess,omitempty"`
	// JSON-formatted audit filter used by the project
	// +optional
	AuditFilter string `json:"auditFilter,omitempty"`
	// Denotes whether or not the project associated with the {GROUP-ID} has database auditing enabled.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}

func (a Auditing) ToAtlas() *mongodbatlas.Auditing {
	return &mongodbatlas.Auditing{
		AuditAuthorizationSuccess: toptr.MakePtr(a.AuditAuthorizationSuccess),
		AuditFilter:               strings.Trim(a.AuditFilter, "\n"),
		Enabled:                   toptr.MakePtr(a.Enabled),
	}
}
