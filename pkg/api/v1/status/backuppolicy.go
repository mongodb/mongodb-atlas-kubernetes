package status

import "github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/collection"

// +k8s:deepcopy-gen=false

// AtlasBackupPolicyStatusOption is the option that is applied to AtlasBackupPolicy Status
type AtlasBackupPolicyStatusOption func(s *BackupPolicyStatus)

func AtlasBackupPolicySetScheduleID(ID string) AtlasBackupPolicyStatusOption {
	return func(s *BackupPolicyStatus) {
		IDs := collection.CopyWithSkip(s.BackupScheduleIDs, ID)
		IDs = append(IDs, ID)

		s.BackupScheduleIDs = IDs
	}
}

func AtlasBackupPolicyUnsetScheduleID(ID string) AtlasBackupPolicyStatusOption {
	return func(s *BackupPolicyStatus) {
		s.BackupScheduleIDs = collection.CopyWithSkip(s.BackupScheduleIDs, ID)
	}
}

type BackupPolicyStatus struct {
	Common `json:",inline"`

	// DeploymentID of the deployment using the backup policy
	BackupScheduleIDs []string `json:"backupScheduleIDs,omitempty"`
}
