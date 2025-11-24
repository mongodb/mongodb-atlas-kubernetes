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

package atlasproject

import (
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
)

type auditController struct {
	ctx     *workflow.Context
	project *akov2.AtlasProject
	service audit.AuditLogService
}

// reconcile dispatch state transitions
func (a *auditController) reconcile() workflow.DeprecatedResult {
	auditInAtlas, err := a.service.Get(a.ctx.Context, a.project.ID())
	if err != nil {
		return a.terminate(workflow.Internal, err)
	}

	isUnset := a.project.Spec.Auditing == nil
	auditInAKO := audit.NewAuditConfig(a.project.Spec.Auditing.DeepCopy())

	if !reflect.DeepEqual(auditInAKO, auditInAtlas) {
		return a.configure(auditInAKO, isUnset)
	}

	if a.project.Spec.Auditing == nil {
		return a.unmanage()
	}

	return a.ready()
}

// configure update Atlas with new audit log configuration
func (a *auditController) configure(auditConfig *audit.AuditConfig, isUnset bool) workflow.DeprecatedResult {
	err := a.service.Update(a.ctx.Context, a.project.ID(), auditConfig)
	if err != nil {
		return a.terminate(workflow.ProjectAuditingReady, err)
	}

	if isUnset {
		return a.unmanage()
	}

	return a.ready()
}

// ready transitions to ready state after successfully configure audit log
func (a *auditController) ready() workflow.DeprecatedResult {
	result := workflow.OK()
	a.ctx.SetConditionFromResult(api.AuditingReadyType, result)

	return result
}

// terminate ends a state transition if an error occurred.
func (a *auditController) terminate(reason workflow.ConditionReason, err error) workflow.DeprecatedResult {
	a.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err)
	a.ctx.SetConditionFromResult(api.AuditingReadyType, result)

	return result
}

// unmanage transitions to unmanaged state if no audit config is set
func (a *auditController) unmanage() workflow.DeprecatedResult {
	a.ctx.UnsetCondition(api.AuditingReadyType)

	return workflow.OK()
}

// handleAudit prepare internal audit controller to handle audit log states
func handleAudit(ctx *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	ctx.Log.Debug("starting audit log processing")
	defer ctx.Log.Debug("finished audit log processing")

	a := auditController{
		ctx:     ctx,
		project: project,
		service: audit.NewAuditLog(ctx.SdkClientSet.SdkClient20250312009.AuditingApi),
	}

	return a.reconcile()
}
