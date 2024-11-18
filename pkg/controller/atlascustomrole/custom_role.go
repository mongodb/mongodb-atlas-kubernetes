package atlascustomrole

import (
	"fmt"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"

	"sigs.k8s.io/controller-runtime/pkg/client"

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
	dpEnabled bool
	k8sClient client.Client
}

func handleCustomRole(ctx *workflow.Context, k8sClient client.Client, akoCustomRole *akov2.AtlasCustomRole,
	dpEnabled bool) workflow.Result {
	ctx.Log.Debug("starting custom role processing")
	defer ctx.Log.Debug("finished custom role processing")

	r := roleController{
		ctx:       ctx,
		service:   customroles.NewCustomRoles(ctx.SdkClient.CustomDatabaseRolesApi),
		role:      akoCustomRole,
		dpEnabled: dpEnabled,
		k8sClient: k8sClient,
	}

	return r.Reconcile()
}

func (r *roleController) Reconcile() workflow.Result {
	if r.role.Spec.ProjectRef != nil {
		project := &akov2.AtlasProject{}
		err := r.k8sClient.Get(r.ctx.Context, *(r.role.Spec.ProjectRef.GetObject(r.role.GetNamespace())), project)
		if err != nil {
			return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
		}
		if project.Status.ID == "" {
			return workflow.Terminate(workflow.ProjectCustomRolesReady,
				fmt.Sprintf("the referenced AtlasProject resource '%s' doesn't have ID (status.ID is empty)", project.GetName()))
		}
		r.projectID = project.Status.ID
	} else if r.role.Spec.ExternalProjectIDRef != nil {
		r.projectID = r.role.Spec.ExternalProjectIDRef.ID
	}

	roleFoundInAtlas := false
	roleInAtlas, err := r.service.Get(r.ctx.Context, r.projectID, r.role.Spec.Role.Name)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err.Error())
	}
	roleFoundInAtlas = roleInAtlas != customroles.CustomRole{}

	roleDeleted := !r.role.GetDeletionTimestamp().IsZero()

	roleInAKO := customroles.NewCustomRole(&r.role.Spec.Role)

	switch {
	case !roleFoundInAtlas && !roleDeleted:
		return r.create(roleInAKO)
	case roleFoundInAtlas && !roleDeleted:
		return r.update(roleInAKO, roleInAtlas)
	case roleFoundInAtlas && roleDeleted && !r.dpEnabled:
		return r.delete(roleInAtlas)
	}

	return r.unmanaged()
}

func (r *roleController) unmanaged() workflow.Result {
	if err := customresource.ManageFinalizer(r.ctx.Context, r.k8sClient, r.role, customresource.UnsetFinalizer); err != nil {
		return r.terminate(workflow.AtlasFinalizerNotRemoved, err)
	}
	return workflow.Deleted()
}

func (r *roleController) managed() workflow.Result {
	if err := customresource.ManageFinalizer(r.ctx.Context, r.k8sClient, r.role, customresource.SetFinalizer); err != nil {
		return r.terminate(workflow.AtlasFinalizerNotSet, err)
	}
	return r.idle()
}

func (r *roleController) create(role customroles.CustomRole) workflow.Result {
	err := r.service.Create(r.ctx.Context, r.projectID, role)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotCreated, err)
	}
	return r.managed()
}

func (r *roleController) update(roleInAKO, roleInAtlas customroles.CustomRole) workflow.Result {
	if reflect.DeepEqual(roleInAKO, roleInAtlas) {
		return r.idle()
	}
	err := r.service.Update(r.ctx.Context, r.projectID, roleInAKO.Name, roleInAKO)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotUpdated, err)
	}
	return r.managed()
}

func (r *roleController) delete(roleInAtlas customroles.CustomRole) workflow.Result {
	err := r.service.Delete(r.ctx.Context, r.projectID, roleInAtlas.Name)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotDeleted, err)
	}
	return r.unmanaged()
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
