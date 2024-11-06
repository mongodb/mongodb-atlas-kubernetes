package atlascustomrole

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type roleController struct {
	ctx       *workflow.Context
	projectID string
	service   customroles.CustomRoleService
	role      *akov2.AtlasCustomRole
	deleted   bool
}

func handleCustomRole(ctx *workflow.Context, akoCustomRole *akov2.AtlasCustomRole) workflow.Result {
	ctx.Log.Debug("starting custom role processing")
	defer ctx.Log.Debug("finished custom role processing")

	r := roleController{
		ctx:       ctx,
		service:   customroles.NewCustomRoles(ctx.SdkClient.CustomDatabaseRolesApi),
		projectID: akoCustomRole.Spec.ProjectIDRef.ID,
		deleted:   !akoCustomRole.DeletionTimestamp.IsZero(),
		role:      akoCustomRole,
	}

	return r.Reconcile()
}

func (r *roleController) Reconcile() workflow.Result {
	currentCustomRoles, err := r.service.List(r.ctx.Context, r.projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}

	roleFoundInAtlas := false
	roleDeleted := r.deleted

	roleInAKO := customroles.NewCustomRole(&r.role.Spec.Role)
	var roleInAtlas customroles.CustomRole
	for _, role := range currentCustomRoles {
		if role.Name == roleInAKO.Name {
			roleFoundInAtlas = true
			roleInAtlas = role
			break
		}
	}

	switch {
	case !roleFoundInAtlas && !roleDeleted:
		return r.create(roleInAKO)
	case roleFoundInAtlas && !roleDeleted:
		return r.update(roleInAKO, roleInAtlas)
	case roleFoundInAtlas && roleDeleted:
		return r.delete(roleInAtlas)
	}

	return r.idle()
}

func (r *roleController) create(role customroles.CustomRole) workflow.Result {
	err := r.service.Create(r.ctx.Context, r.projectID, role)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotCreated, err)
	}
	return r.idle()
}

func (r *roleController) update(roleInAKO, roleInAtlas customroles.CustomRole) workflow.Result {
	if cmp.Diff(roleInAKO, roleInAtlas, cmpopts.EquateEmpty()) == "" {
		return r.idle()
	}
	err := r.service.Update(r.ctx.Context, r.projectID, roleInAKO.Name, roleInAKO)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotUpdated, err)
	}
	return r.idle()
}

func (r *roleController) delete(roleInAtlas customroles.CustomRole) workflow.Result {
	err := r.service.Delete(r.ctx.Context, r.projectID, roleInAtlas.Name)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotDeleted, err)
	}
	return r.idle()
}

func (r *roleController) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	r.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err.Error())
	r.ctx.SetConditionFromResult(api.ProjectCustomRolesReadyType, result)
	return result
}

func (r *roleController) idle() workflow.Result {
	r.ctx.SetConditionTrue(api.ProjectCustomRolesReadyType)
	return workflow.OK()
}
