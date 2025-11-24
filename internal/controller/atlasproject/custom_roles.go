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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
)

type roleController struct {
	ctx     *workflow.Context
	project *akov2.AtlasProject
	service customroles.CustomRoleService
}

func getLastAppliedCustomRoles(atlasProject *akov2.AtlasProject) ([]akov2.CustomRole, error) {
	lastAppliedSpec := akov2.AtlasProjectSpec{}
	lastAppliedSpecStr, ok := atlasProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if !ok {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(lastAppliedSpecStr), &lastAppliedSpec); err != nil {
		return nil, fmt.Errorf("failed to parse last applied configuration: %w", err)
	}

	return lastAppliedSpec.CustomRoles, nil
}

func findRolesToDelete(prevSpec, akoRoles, atlasRoles []customroles.CustomRole) map[string]customroles.CustomRole {
	result := map[string]customroles.CustomRole{}
	lastAppliedRolesMap := mapCustomRolesByName(prevSpec)
	akoRolesMap := mapCustomRolesByName(akoRoles)
	atlasRolesMap := mapCustomRolesByName(atlasRoles)

	for atlasName, atlasRole := range atlasRolesMap {
		_, inAKO := akoRolesMap[atlasName]
		_, inLastApplied := lastAppliedRolesMap[atlasName]
		if !inAKO && inLastApplied {
			result[atlasName] = atlasRole
		}
	}

	return result
}

func convertToInternalRoles(roles []akov2.CustomRole) []customroles.CustomRole {
	result := make([]customroles.CustomRole, 0, len(roles))
	for i := range roles {
		result = append(result, customroles.NewCustomRole(&roles[i]))
	}
	return result
}

func ensureCustomRoles(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	lastAppliedCustomRoles, err := getLastAppliedCustomRoles(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	r := roleController{
		ctx:     workflowCtx,
		project: project,
		service: customroles.NewCustomRoles(workflowCtx.SdkClientSet.SdkClient20250312009.CustomDatabaseRolesApi),
	}

	currentAtlasCustomRoles, err := r.service.List(r.ctx.Context, r.project.ID())
	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, err)
	}

	akoRoles := make([]customroles.CustomRole, len(project.Spec.CustomRoles))
	for i := range project.Spec.CustomRoles {
		akoRoles[i] = customroles.NewCustomRole(&project.Spec.CustomRoles[i])
	}

	ops := calculateChanges(currentAtlasCustomRoles, akoRoles)

	var deleteStatus map[string]status.CustomRole
	if len(lastAppliedCustomRoles) > 0 {
		deleteStatus = r.deleteCustomRoles(workflowCtx, project.ID(),
			findRolesToDelete(convertToInternalRoles(lastAppliedCustomRoles), akoRoles, currentAtlasCustomRoles))
	}
	updateStatus := r.updateCustomRoles(workflowCtx, project.ID(), ops.Update)
	createStatus := r.createCustomRoles(workflowCtx, project.ID(), ops.Create)

	result := syncCustomRolesStatus(workflowCtx, akoRoles, createStatus, updateStatus, deleteStatus)

	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.ProjectCustomRolesReadyType, result)

		return result
	}

	workflowCtx.SetConditionTrue(api.ProjectCustomRolesReadyType)

	if len(project.Spec.CustomRoles) == 0 {
		workflowCtx.UnsetCondition(api.ProjectCustomRolesReadyType)
	}

	return result
}

func (r *roleController) deleteCustomRoles(ctx *workflow.Context, projectID string, toDelete map[string]customroles.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be deleted", "NumItems", len(toDelete))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toDelete {
		err := r.service.Delete(ctx.Context, projectID, customRole.Name)

		opStatus, errorMsg := evaluateOperation(err)
		statuses[customRole.Name] = status.CustomRole{
			Name:   customRole.Name,
			Status: opStatus,
			Error:  errorMsg,
		}

		if errorMsg != "" {
			ctx.Log.Warnf("Failed to delete custom role \"%s\": %s", customRole.Name, errorMsg)
		} else {
			ctx.Log.Debugw("Removed Custom Role in current AtlasProject", "custom role:", customRole.Name)
		}
	}

	return statuses
}

func (r *roleController) updateCustomRoles(ctx *workflow.Context, projectID string, toUpdate map[string]customroles.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be updated", "NumItems", len(toUpdate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toUpdate {
		err := r.service.Update(ctx.Context, projectID, customRole.Name, customRole)

		opStatus, errorMsg := evaluateOperation(err)

		statuses[customRole.Name] = status.CustomRole{
			Name:   customRole.Name,
			Status: opStatus,
			Error:  errorMsg,
		}

		if errorMsg != "" {
			ctx.Log.Warnf("Failed to update custom role \"%s\": %s", customRole.Name, errorMsg)
		} else {
			ctx.Log.Debugw("Updated Custom Role in current AtlasProject", "custom role:", customRole.Name)
		}
	}

	return statuses
}

func (r *roleController) createCustomRoles(ctx *workflow.Context, projectID string, toCreate map[string]customroles.CustomRole) map[string]status.CustomRole {
	ctx.Log.Debugw("Custom Roles to be added", "NumItems", len(toCreate))

	statuses := map[string]status.CustomRole{}
	for _, customRole := range toCreate {
		err := r.service.Create(ctx.Context, projectID, customRole)

		opStatus, errorMsg := evaluateOperation(err)

		statuses[customRole.Name] = status.CustomRole{
			Name:   customRole.Name,
			Status: opStatus,
			Error:  errorMsg,
		}

		if errorMsg != "" {
			ctx.Log.Warnf("Failed to create custom role \"%s\": %s", customRole.Name, errorMsg)
		} else {
			ctx.Log.Debugw("Created Custom Role in current AtlasProject", "custom role:", customRole.Name)
		}
	}

	return statuses
}

func mapCustomRolesByName(customRoles []customroles.CustomRole) map[string]customroles.CustomRole {
	customRolesByName := map[string]customroles.CustomRole{}

	for _, customRole := range customRoles {
		customRolesByName[customRole.Name] = customRole
	}

	return customRolesByName
}

type CustomRolesOperations struct {
	Create map[string]customroles.CustomRole
	Update map[string]customroles.CustomRole
}

func calculateChanges(currentCustomRoles []customroles.CustomRole, desiredCustomRoles []customroles.CustomRole) CustomRolesOperations {
	currentCustomRolesByName := mapCustomRolesByName(currentCustomRoles)
	ops := CustomRolesOperations{
		Create: map[string]customroles.CustomRole{},
		Update: map[string]customroles.CustomRole{},
	}

	for _, desiredCustomRole := range desiredCustomRoles {
		customRole, ok := currentCustomRolesByName[desiredCustomRole.Name]

		if !ok {
			ops.Create[desiredCustomRole.Name] = desiredCustomRole

			continue
		}

		if d := cmp.Diff(desiredCustomRole, customRole, cmpopts.EquateEmpty()); d != "" {
			ops.Update[desiredCustomRole.Name] = desiredCustomRole
		}
	}

	return ops
}

func evaluateOperation(err error) (status.CustomRoleStatus, string) {
	if err != nil {
		return status.CustomRoleStatusFailed, err.Error()
	}

	return status.CustomRoleStatusOK, ""
}

func syncCustomRolesStatus(ctx *workflow.Context, desiredCustomRoles []customroles.CustomRole, created, updated, deleted map[string]status.CustomRole) workflow.DeprecatedResult {
	statuses := make([]status.CustomRole, 0, len(desiredCustomRoles))
	var err error

	for _, customRole := range desiredCustomRoles {
		stat, ok := deleted[customRole.Name]

		if ok {
			if stat.Status == status.CustomRoleStatusFailed {
				statuses = append(statuses, stat)
				err = errors.Join(err, fmt.Errorf("%s", stat.Error))
			}

			continue
		}

		if stat, ok = updated[customRole.Name]; ok {
			statuses = append(statuses, stat)

			if stat.Status == status.CustomRoleStatusFailed {
				err = errors.Join(err, fmt.Errorf("%s", stat.Error))
			}

			continue
		}

		if stat, ok = created[customRole.Name]; ok {
			statuses = append(statuses, stat)

			if stat.Status == status.CustomRoleStatusFailed {
				err = errors.Join(err, fmt.Errorf("%s", stat.Error))
			}

			continue
		}

		statuses = append(statuses, status.CustomRole{
			Name:   customRole.Name,
			Status: status.CustomRoleStatusOK,
		})
	}

	ctx.EnsureStatusOption(status.AtlasProjectSetCustomRolesOption(&statuses))

	if err != nil {
		return workflow.Terminate(workflow.ProjectCustomRolesReady, fmt.Errorf("failed to apply changes to custom roles: %w", err))
	}

	return workflow.OK()
}
