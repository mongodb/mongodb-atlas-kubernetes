//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasBackupScheduleByBackupPolicyIndex = "atlasbackupschedule.spec.policyRef"
)

type AtlasBackupScheduleByBackupPolicyIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasBackupScheduleByBackupPolicyIndexer(logger *zap.Logger) *AtlasBackupScheduleByBackupPolicyIndexer {
	return &AtlasBackupScheduleByBackupPolicyIndexer{
		logger: logger.Named(AtlasBackupScheduleByBackupPolicyIndex).Sugar(),
	}
}

func (*AtlasBackupScheduleByBackupPolicyIndexer) Object() client.Object {
	return &akov2.AtlasBackupSchedule{}
}

func (*AtlasBackupScheduleByBackupPolicyIndexer) Name() string {
	return AtlasBackupScheduleByBackupPolicyIndex
}

func (a *AtlasBackupScheduleByBackupPolicyIndexer) Keys(object client.Object) []string {
	schedule, ok := object.(*akov2.AtlasBackupSchedule)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasBackupSchedule but got %T", object)
		return nil
	}

	if schedule.Spec.PolicyRef.IsEmpty() {
		return nil
	}

	return []string{schedule.Spec.PolicyRef.GetObject(schedule.Namespace).String()}
}
