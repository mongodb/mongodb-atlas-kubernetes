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
	"errors"
	"fmt"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

const ipAccessStatusPending = "PENDING"
const ipAccessStatusFailed = "FAILED"

type ipAccessListController struct {
	ctx         *workflow.Context
	project     *akov2.AtlasProject
	service     ipaccesslist.IPAccessListService
	lastApplied ipaccesslist.IPAccessEntries
}

// reconcile dispatch state transitions
func (i *ipAccessListController) reconcile() workflow.DeprecatedResult {
	ialInAtlas, err := i.service.List(i.ctx.Context, i.project.ID())
	if err != nil {
		return i.terminate(workflow.Internal, err)
	}

	isUnset := len(i.project.Spec.ProjectIPAccessList) == 0
	ialInAKO, err := ipaccesslist.NewIPAccessEntries(i.project.Spec.ProjectIPAccessList)
	if err != nil {
		return i.terminate(workflow.Internal, err)
	}

	if !reflect.DeepEqual(ialInAKO.GetByStatus(false), ialInAtlas) {
		return i.configure(ialInAtlas, ialInAKO, isUnset)
	}

	if isUnset {
		return i.unmanage()
	}

	return i.progress(ialInAKO)
}

// configure update Atlas with new ip access list
func (i *ipAccessListController) configure(current, desired ipaccesslist.IPAccessEntries, isUnset bool) workflow.DeprecatedResult {
	err := i.service.Add(i.ctx.Context, i.project.ID(), desired.GetByStatus(false))
	if err != nil {
		return i.terminate(workflow.ProjectIPNotCreatedInAtlas, err)
	}

	managedEntries := len(i.lastApplied) > 0
	if managedEntries {
		for key, entry := range current {
			if _, ok := desired[key]; !ok {
				err = i.service.Delete(i.ctx.Context, i.project.ID(), entry)
				if err != nil {
					return i.terminate(workflow.ProjectIPNotCreatedInAtlas, err)
				}
			}
		}
	}

	if isUnset {
		return i.unmanage()
	}

	return i.progress(desired)
}

// progress transitions to pending while ip access list are not active
func (i *ipAccessListController) progress(ipAccessEntries ipaccesslist.IPAccessEntries) workflow.DeprecatedResult {
	for _, entry := range ipAccessEntries.GetByStatus(false) {
		stat, err := i.service.Status(i.ctx.Context, i.project.ID(), entry)
		if err != nil {
			return i.terminate(workflow.Internal, err)
		}

		switch stat {
		case ipAccessStatusPending:
			result := workflow.InProgress(
				workflow.ProjectIPAccessListNotActive,
				"atlas is adding access. this entry may not apply to all cloud providers at the time of this request",
			)
			i.ctx.SetConditionFromResult(api.IPAccessListReadyType, result)

			return result
		case ipAccessStatusFailed:
			return i.terminate(workflow.ProjectIPNotCreatedInAtlas, errors.New("atlas didn't succeed in adding this access entry"))
		}
	}

	return i.ready(ipAccessEntries)
}

// ready transitions to ready state after successfully configure ip access list
func (i *ipAccessListController) ready(ipAccessEntries ipaccesslist.IPAccessEntries) workflow.DeprecatedResult {
	i.ctx.EnsureStatusOption(status.AtlasProjectExpiredIPAccessOption(ipaccesslist.FromInternal(ipAccessEntries.GetByStatus(true))))

	result := workflow.OK()
	i.ctx.SetConditionFromResult(api.IPAccessListReadyType, result)

	return result
}

// terminate ends a state transition if an error occurred.
func (i *ipAccessListController) terminate(reason workflow.ConditionReason, err error) workflow.DeprecatedResult {
	i.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err)
	i.ctx.SetConditionFromResult(api.IPAccessListReadyType, result)

	return result
}

// unmanage transitions to unmanaged state if no ip access list config is set
func (i *ipAccessListController) unmanage() workflow.DeprecatedResult {
	i.ctx.UnsetCondition(api.IPAccessListReadyType)

	return workflow.OK()
}

// handleIPAccessList prepare internal ip access list controller to handle states
func handleIPAccessList(ctx *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	ctx.Log.Debug("starting ip access list processing")
	defer ctx.Log.Debug("finished ip access list processing")

	lastApplied, err := mapLastAppliedIPAccessList(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Errorf("failed to get last applied configuration: %w", err))
	}

	c := ipAccessListController{
		ctx:         ctx,
		project:     project,
		service:     ipaccesslist.NewIPAccessList(ctx.SdkClientSet.SdkClient20250312011.ProjectIPAccessListApi),
		lastApplied: lastApplied,
	}

	return c.reconcile()
}

func mapLastAppliedIPAccessList(atlasProject *akov2.AtlasProject) (ipaccesslist.IPAccessEntries, error) {
	lastApplied, err := lastAppliedSpecFrom(atlasProject)
	if err != nil {
		return nil, err
	}

	if lastApplied == nil || len(lastApplied.ProjectIPAccessList) == 0 {
		return nil, nil
	}

	entries, err := ipaccesslist.NewIPAccessEntries(lastApplied.ProjectIPAccessList)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
