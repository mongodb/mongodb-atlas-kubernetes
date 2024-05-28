package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasProjectByBackupCompliancePolicyIndex = ""
)

type AtlasProjectByBackupCompliancePolicyIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasProjectByBackupCompliancePolicyIndexer(logger *zap.Logger) *AtlasProjectByBackupCompliancePolicyIndexer {
	return &AtlasProjectByBackupCompliancePolicyIndexer{
		logger: logger.Named(AtlasProjectByBackupCompliancePolicyIndex).Sugar(),
	}
}

func (AtlasProjectByBackupCompliancePolicyIndexer) Object() client.Object {
	return &akov2.AtlasProject{}
}

func (*AtlasProjectByBackupCompliancePolicyIndexer) Name() string {
	return AtlasProjectByBackupCompliancePolicyIndex
}

func (a *AtlasProjectByBackupCompliancePolicyIndexer) Keys(object client.Object) []string {
	project, ok := object.(*akov2.AtlasProject)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasProject but got %T", object)
		return nil
	}

	if project.Spec.BackupCompliancePolicyRef.IsEmpty() {
		return nil
	}

	return []string{project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace).String()}
}
