package audit

import (
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

func fromAtlas(auditLog *admin.AuditLog) *v1alpha1.AtlasAuditingSpec {
	return &v1alpha1.AtlasAuditingSpec{
		Enabled:                   pointer.GetOrDefault(auditLog.Enabled, false),
		AuditAuthorizationSuccess: pointer.GetOrDefault(auditLog.AuditAuthorizationSuccess, false),
		AuditFilter:               jsonFromAtlas(auditLog.AuditFilter),
	}
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
