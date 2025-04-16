// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasProjectByBackupCompliancePolicyIndex = "atlasproject.spec.backupCompliancePolicyRef"
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

	if project.Spec.BackupCompliancePolicyRef == nil {
		return nil
	}

	if project.Spec.BackupCompliancePolicyRef.IsEmpty() {
		return nil
	}

	return []string{project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace).String()}
}
