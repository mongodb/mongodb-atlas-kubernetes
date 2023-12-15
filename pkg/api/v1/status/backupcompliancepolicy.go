package status

type AtlasBackupCompliancePolicyStatusOption func(s *BackupCompliancePolicyStatus)

type BackupCompliancePolicyStatus struct {
	Common `json:",inline"`
}
