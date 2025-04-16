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
