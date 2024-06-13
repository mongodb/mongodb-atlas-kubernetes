package atlasproject

import (
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type auditController struct {
	ctx     *workflow.Context
	project *akov2.AtlasProject
	service audit.AuditLogService
}

// reconcile dispatch state transitions
func (a *auditController) reconcile() workflow.Result {
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
func (a *auditController) configure(auditConfig *audit.AuditConfig, isUnset bool) workflow.Result {
	err := a.service.Set(a.ctx.Context, a.project.ID(), auditConfig)
	if err != nil {
		return a.terminate(workflow.ProjectAuditingReady, err)
	}

	if isUnset {
		return a.unmanage()
	}

	return a.ready()
}

// ready transitions to ready state after successfully configure audit log
func (a *auditController) ready() workflow.Result {
	result := workflow.OK()
	a.ctx.SetConditionFromResult(api.AuditingReadyType, result)

	return result
}

// terminate ends a state transition if an error occurred.
func (a *auditController) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	a.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err.Error())
	a.ctx.SetConditionFromResult(api.AuditingReadyType, result)

	return result
}

// unmanage transitions to unmanaged state if no audit config is set
func (a *auditController) unmanage() workflow.Result {
	a.ctx.UnsetCondition(api.AuditingReadyType)

	return workflow.OK()
}

// handleAudit prepare internal audit controller to handle audit log states
func handleAudit(ctx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	ctx.Log.Debug("starting audit log processing")
	defer ctx.Log.Debug("finished audit log processing")

	a := auditController{
		ctx:     ctx,
		project: project,
		service: audit.NewAuditLog(ctx.SdkClient.AuditingApi),
	}

	return a.reconcile()
}
