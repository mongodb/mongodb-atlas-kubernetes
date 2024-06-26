package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

// +kubebuilder:object:generate=false

type AtlasBackupCompliancePolicyStatusOption func(s *BackupCompliancePolicyStatus)

type BackupCompliancePolicyStatus struct {
	api.Common `json:",inline"`
}
