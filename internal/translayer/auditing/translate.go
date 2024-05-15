package auditing

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
)

func toAtlas(spec *v1alpha1.AtlasAuditingSpec) *admin.AuditLog {
	return &admin.AuditLog{
		Enabled:                   pointer.MakePtr(spec.Enabled),
		AuditAuthorizationSuccess: pointer.MakePtr(spec.AuditAuthorizationSuccess),
		AuditFilter:               pointer.MakePtr(jsonToAtlas(spec.AuditFilter)),
		// ConfigurationType is not set on the PATCH operation to Atlas
	}
}

func fromAtlas(auditLog *admin.AuditLog) (*v1alpha1.AtlasAuditingSpec, error) {
	cfgType, err := configTypeFromAtlas(auditLog.ConfigurationType)
	if err != nil {
		return nil, err
	}
	return &v1alpha1.AtlasAuditingSpec{
		Enabled:                   pointer.GetOrDefault(auditLog.Enabled, false),
		AuditAuthorizationSuccess: pointer.GetOrDefault(auditLog.AuditAuthorizationSuccess, false),
		ConfigurationType:         cfgType,
		AuditFilter:               jsonFromAtlas(auditLog.AuditFilter),
	}, nil
}

func jsonToAtlas(js *apiextensionsv1.JSON) string {
	if js == nil {
		return ""
	}
	return string(js.Raw)
}

func jsonFromAtlas(js *string) *apiextensionsv1.JSON {
	if js == nil {
		return nil
	}
	return &apiextensionsv1.JSON{Raw: ([]byte)(*js)}
}

func configTypeFromAtlas(configType *string) (v1alpha1.AuditingConfigTypes, error) {
	ct := pointer.GetOrDefault(configType, string(v1alpha1.None))
	switch ct {
	case string(v1alpha1.None), string(v1alpha1.FilterBuilder), string(v1alpha1.FilterJSON):
		return v1alpha1.AuditingConfigTypes(ct), nil
	default:
		return v1alpha1.AuditingConfigTypes(ct), fmt.Errorf("unsupported Auditing Config Type %q", ct)
	}
}
