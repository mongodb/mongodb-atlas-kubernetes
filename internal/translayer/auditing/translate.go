package auditing

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type AuditingConfigType string

const (
	None          AuditingConfigType = "NONE"
	FilterBuilder AuditingConfigType = "FILTER_BUILDER"
	FilterJSON    AuditingConfigType = "FILTER_JSON"
)

type Auditing struct {
	Enabled                   bool
	AuditAuthorizationSuccess bool
	ConfigurationType         AuditingConfigType
	AuditFilter               string
}

func toAtlas(auditing *Auditing) *admin.AuditLog {
	return &admin.AuditLog{
		Enabled:                   pointer.MakePtr(auditing.Enabled),
		AuditAuthorizationSuccess: pointer.MakePtr(auditing.AuditAuthorizationSuccess),
		AuditFilter:               pointer.MakePtr(auditing.AuditFilter),
		// ConfigurationType is not set on the PATCH operation to Atlas
	}
}

func fromAtlas(auditLog *admin.AuditLog) (*Auditing, error) {
	cfgType, err := configTypeFromAtlas(auditLog.ConfigurationType)
	if err != nil {
		return nil, err
	}
	return &Auditing{
		Enabled:                   pointer.GetOrDefault(auditLog.Enabled, false),
		AuditAuthorizationSuccess: pointer.GetOrDefault(auditLog.AuditAuthorizationSuccess, false),
		ConfigurationType:         cfgType,
		AuditFilter:               pointer.GetOrDefault(auditLog.AuditFilter, ""),
	}, nil
}

func configTypeFromAtlas(configType *string) (AuditingConfigType, error) {
	ct := pointer.GetOrDefault(configType, string(None))
	switch ct {
	case string(None), string(FilterBuilder), string(FilterJSON):
		return AuditingConfigType(ct), nil
	default:
		return AuditingConfigType(ct), fmt.Errorf("unsupported Auditing Config type %q", ct)
	}
}
