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

package atlascustomrole

import (
	"fmt"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

type roleController struct {
	ctx       *workflow.Context
	project   *project.Project
	service   customroles.CustomRoleService
	role      *akov2.AtlasCustomRole
	dpEnabled bool
	k8sClient client.Client
}

func handleCustomRole(ctx *workflow.Context, k8sClient client.Client, project *project.Project, service customroles.CustomRoleService, akoCustomRole *akov2.AtlasCustomRole, dpEnabled bool) workflow.Result {
	ctx.Log.Debug("starting custom role processing")
	defer ctx.Log.Debug("finished custom role processing")

	r := roleController{
		ctx:       ctx,
		project:   project,
		service:   service,
		role:      akoCustomRole,
		dpEnabled: dpEnabled,
		k8sClient: k8sClient,
	}

	result := r.Reconcile()
	ctx.SetConditionFromResult(api.ReadyType, result)
	return result
}

func (r *roleController) Reconcile() workflow.Result {
	if r.project.ID == "" {
		return workflow.Terminate(workflow.ProjectCustomRolesReady,
			fmt.Errorf("the referenced AtlasProject resource '%s' doesn't have ID (status.ID is empty)", r.project.Name))
	}
	roleFoundInAtlas := false
	roleInAtlas, err := r.service.Get(r.ctx.Context, r.project.ID, r.role.Spec.Role.Name)
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err)
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
	err := r.service.Create(r.ctx.Context, r.project.ID, role)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotCreated, err)
	}
	return r.managed()
}

func (r *roleController) update(roleInAKO, roleInAtlas customroles.CustomRole) workflow.Result {
	if reflect.DeepEqual(roleInAKO, roleInAtlas) {
		return r.idle()
	}
	err := r.service.Update(r.ctx.Context, r.project.ID, roleInAKO.Name, roleInAKO)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotUpdated, err)
	}
	return r.managed()
}

func (r *roleController) delete(roleInAtlas customroles.CustomRole) workflow.Result {
	err := r.service.Delete(r.ctx.Context, r.project.ID, roleInAtlas.Name)
	if err != nil {
		return r.terminate(workflow.AtlasCustomRoleNotDeleted, err)
	}
	return r.unmanaged()
}

func (r *roleController) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	r.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err)
	r.ctx.SetConditionFromResult(api.ProjectCustomRolesReadyType, result)
	return result
}

func (r *roleController) idle() workflow.Result {
	r.ctx.SetConditionTrue(api.ProjectCustomRolesReadyType)
	return workflow.OK()
}
