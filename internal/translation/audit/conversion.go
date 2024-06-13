package audit

import (
	"fmt"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type AuditingConfigType string

const (
	None          AuditingConfigType = "NONE"
	FilterBuilder AuditingConfigType = "FILTER_BUILDER"
	FilterJSON    AuditingConfigType = "FILTER_JSON"
)

// AuditConfig represents the Atlas Project audit log config
type AuditConfig struct {
	*akov2.Auditing

	ConfigurationType AuditingConfigType
}

func NewAuditConfig(auditConfig *akov2.Auditing) *AuditConfig {
	configType := FilterJSON

	if auditConfig == nil {
		auditConfig = &akov2.Auditing{
			AuditFilter: "{}",
		}
		configType = None
	}

	return &AuditConfig{
		Auditing:          auditConfig,
		ConfigurationType: configType,
	}
}

func toAtlas(auditing *AuditConfig) *admin.AuditLog {
	return &admin.AuditLog{
		Enabled:                   pointer.MakePtr(auditing.Enabled),
		AuditAuthorizationSuccess: pointer.MakePtr(auditing.AuditAuthorizationSuccess),
		AuditFilter:               pointer.MakePtr(auditing.AuditFilter),
		// ConfigurationType is not set on the PATCH operation to Atlas
	}
}

func fromAtlas(auditLog *admin.AuditLog) (*AuditConfig, error) {
	cfgType, err := configTypeFromAtlas(auditLog.ConfigurationType)
	if err != nil {
		return nil, err
	}
	return &AuditConfig{
		Auditing: &akov2.Auditing{
			AuditAuthorizationSuccess: pointer.GetOrDefault(auditLog.AuditAuthorizationSuccess, false),
			AuditFilter:               pointer.GetOrDefault(auditLog.AuditFilter, "{}"),
			Enabled:                   pointer.GetOrDefault(auditLog.Enabled, false),
		},
		ConfigurationType: cfgType,
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
