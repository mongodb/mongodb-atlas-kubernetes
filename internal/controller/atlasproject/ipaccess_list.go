package atlasproject

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
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
func (i *ipAccessListController) reconcile() workflow.Result {
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
func (i *ipAccessListController) configure(current, desired ipaccesslist.IPAccessEntries, isUnset bool) workflow.Result {
	err := i.service.Add(i.ctx.Context, i.project.ID(), desired.GetByStatus(false))
	if err != nil {
		return i.terminate(workflow.ProjectIPNotCreatedInAtlas, err)
	}

	for key, entry := range current {
		if _, ok := i.lastApplied[key]; !ok {
			continue
		}

		if _, ok := desired[key]; !ok {
			err = i.service.Delete(i.ctx.Context, i.project.ID(), entry)
			if err != nil {
				return i.terminate(workflow.ProjectIPNotCreatedInAtlas, err)
			}
		}
	}

	if isUnset {
		return i.unmanage()
	}

	return i.progress(desired)
}

// progress transitions to pending while ip access list are not active
func (i *ipAccessListController) progress(ipAccessEntries ipaccesslist.IPAccessEntries) workflow.Result {
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
func (i *ipAccessListController) ready(ipAccessEntries ipaccesslist.IPAccessEntries) workflow.Result {
	i.ctx.EnsureStatusOption(status.AtlasProjectExpiredIPAccessOption(ipaccesslist.FromInternal(ipAccessEntries.GetByStatus(true))))

	result := workflow.OK()
	i.ctx.SetConditionFromResult(api.IPAccessListReadyType, result)

	return result
}

// terminate ends a state transition if an error occurred.
func (i *ipAccessListController) terminate(reason workflow.ConditionReason, err error) workflow.Result {
	i.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err)
	i.ctx.SetConditionFromResult(api.IPAccessListReadyType, result)

	return result
}

// unmanage transitions to unmanaged state if no ip access list config is set
func (i *ipAccessListController) unmanage() workflow.Result {
	i.ctx.UnsetCondition(api.IPAccessListReadyType)

	return workflow.OK()
}

// handleIPAccessList prepare internal ip access list controller to handle states
func handleIPAccessList(ctx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	ctx.Log.Debug("starting ip access list processing")
	defer ctx.Log.Debug("finished ip access list processing")

	skipped, err := hasSkippedIPAccessListConfiguration(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if skipped {
		ctx.EnsureStatusOption(status.AtlasProjectExpiredIPAccessOption(nil))
		ctx.UnsetCondition(api.IPAccessListReadyType)

		return workflow.OK()
	}

	lastApplied, err := mapLastAppliedIPAccessList(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Sprintf("failed to get last applied configuration: %s", err))
	}

	c := ipAccessListController{
		ctx:         ctx,
		project:     project,
		service:     ipaccesslist.NewIPAccessList(ctx.SdkClient.ProjectIPAccessListApi),
		lastApplied: lastApplied,
	}

	return c.reconcile()
}

func hasSkippedIPAccessListConfiguration(atlasProject *akov2.AtlasProject) (bool, error) {
	lastSkippedSpec := akov2.AtlasProjectSpec{}
	lastSkippedSpecString, ok := atlasProject.Annotations[customresource.AnnotationLastSkippedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(lastSkippedSpecString), &lastSkippedSpec); err != nil {
			return false, fmt.Errorf("failed to parse last skipped configuration: %w", err)
		}

		return len(lastSkippedSpec.ProjectIPAccessList) > 0, nil
	}

	return false, nil
}

func mapLastAppliedIPAccessList(atlasProject *akov2.AtlasProject) (ipaccesslist.IPAccessEntries, error) {
	lastApplied, err := lastSpecFrom(atlasProject, customresource.AnnotationLastAppliedConfiguration)
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
