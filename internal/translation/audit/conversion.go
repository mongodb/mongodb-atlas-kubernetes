package audit

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

type AuditingConfigType string

const (
	FilterDefault = "{}"
)

// AuditConfig represents the Atlas Project audit log config
type AuditConfig struct {
	*akov2.Auditing
}

func NewAuditConfig(auditConfig *akov2.Auditing) *AuditConfig {
	if auditConfig == nil {
		auditConfig = &akov2.Auditing{}
	}

	if auditConfig.AuditFilter == "" {
		auditConfig.AuditFilter = FilterDefault
	}

	return &AuditConfig{
		Auditing: auditConfig,
	}
}

func toAtlas(auditing *AuditConfig) *admin.AuditLog {
	auditLog := admin.NewAuditLogWithDefaults()
	auditLog.SetEnabled(auditing.Enabled)
	auditLog.SetAuditAuthorizationSuccess(auditing.AuditAuthorizationSuccess)
	auditLog.SetAuditFilter(auditing.AuditFilter)
	// ConfigurationType is not set on the PATCH operation to Atlas

	return auditLog
}

func fromAtlas(auditLog *admin.AuditLog) *AuditConfig {
	auditFilter := FilterDefault

	if auditLog.GetAuditFilter() != "" {
		auditFilter = auditLog.GetAuditFilter()
	}

	return &AuditConfig{
		Auditing: &akov2.Auditing{
			Enabled:                   auditLog.GetEnabled(),
			AuditAuthorizationSuccess: auditLog.GetAuditAuthorizationSuccess(),
			AuditFilter:               auditFilter,
		},
	}
}
