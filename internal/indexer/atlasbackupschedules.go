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

//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasDeploymentByBackupScheduleIndex = "atlasdeployment.spec.backupScheduleRef"
)

type AtlasDeploymentByBackupScheduleIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasDeploymentByBackupScheduleIndexer(logger *zap.Logger) *AtlasDeploymentByBackupScheduleIndexer {
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
