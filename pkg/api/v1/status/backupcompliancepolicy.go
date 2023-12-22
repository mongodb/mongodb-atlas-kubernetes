package status

// +kubebuilder:object:generate=false

type AtlasBackupCompliancePolicyStatusOption func(s *BackupCompliancePolicyStatus)

type BackupCompliancePolicyStatus struct {
	Common `json:",inline"`

	ProjectIDs []string `json:"projectID,omitempty"`
}
