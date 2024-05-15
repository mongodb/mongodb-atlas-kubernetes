package indexer

import (
	"go.uber.org/zap"

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasDeploymentByBackupScheduleIndex = ".spec.backupScheduleRef"
)

type AtlasDeploymentByBackupScheduleIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasBackupScheduleToDeploymentIndex(logger *zap.Logger) *AtlasDeploymentByBackupScheduleIndexer {
	return &AtlasDeploymentByBackupScheduleIndexer{
		logger: logger.Named(AtlasDeploymentByBackupScheduleIndex).Sugar(),
	}
}

func (*AtlasDeploymentByBackupScheduleIndexer) Object() client.Object {
	return &akov2.AtlasDeployment{}
}

func (*AtlasDeploymentByBackupScheduleIndexer) Name() string {
	return AtlasDeploymentByBackupScheduleIndex
}

func (a *AtlasDeploymentByBackupScheduleIndexer) Keys(object client.Object) []string {
	deployment, ok := object.(*akov2.AtlasDeployment)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasDeployment but got %T", object)
		return nil
	}

	if deployment.Spec.BackupScheduleRef.IsEmpty() {
		return nil
	}

	return []string{deployment.Spec.BackupScheduleRef.GetObject(deployment.Namespace).String()}
}
