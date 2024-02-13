package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const ipAccessStatusPending = "PENDING"
const ipAccessStatusFailed = "FAILED"

// ensureIPAccessList ensures that the state of the Atlas IP Access List matches the
// state of the IP Access list specified in the project CR. Any Access Lists which exist
// in Atlas but are not specified in the CR are deleted.
func ensureIPAccessList(service *workflow.Context, statusFunc atlas.IPAccessListStatus, akoProject *mdbv1.AtlasProject, subobjectProtect bool) workflow.Result {
	canReconcile, err := canIPAccessListReconcile(service.Context, service.SdkClient, subobjectProtect, akoProject)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		service.SetConditionFromResult(status.IPAccessListReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile IP Access List due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		service.SetConditionFromResult(status.IPAccessListReadyType, result)

		return result
	}

	desiredList, expiredList := filterActiveIPAccessLists(akoProject.Spec.ProjectIPAccessList)
	service.EnsureStatusOption(status.AtlasProjectExpiredIPAccessOption(expiredList))

	list, _, err := service.SdkClient.ProjectIPAccessListApi.ListProjectIpAccessLists(service.Context, akoProject.ID()).Execute()
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("failed to retrieve IP Access list: %s", err))
		service.SetConditionFromResult(status.IPAccessListReadyType, result)

		return result
	}

	currentList := mapToOperatorSpec(list.GetResults())
	if d := cmp.Diff(currentList, akoProject.Spec.ProjectIPAccessList, cmpopts.EquateEmpty()); d != "" {
		service.Log.Infof("IP Access List differs from spec: %s", d)
		err = syncIPAccessList(service, akoProject.ID(), currentList, desiredList)
		if err != nil {
			result := workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, fmt.Sprintf("failed to sync desired state with Atlas: %s", err))
			service.SetConditionFromResult(status.IPAccessListReadyType, result)

			return result
		}
	}

	for _, ipAccessList := range desiredList {
		ipAccessStatus, err := statusFunc(service.Context, akoProject.ID(), mapToEntryValue(ipAccessList))
		if err != nil {
			result := workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, fmt.Sprintf("failed to check status in Atlas: %s", err))
			service.SetConditionFromResult(status.IPAccessListReadyType, result)

			return result
		}

		if ipAccessStatus == ipAccessStatusFailed {
			result := workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, fmt.Sprintf("configuration of %s failed in Atlas", mapToEntryValue(ipAccessList)))
			service.SetConditionFromResult(status.IPAccessListReadyType, result)

			return result
		}

		if ipAccessStatus == ipAccessStatusPending {
			result := workflow.InProgress(workflow.ProjectIPAccessListNotActive, fmt.Sprintf("waiting Atlas to configure entry %s", mapToEntryValue(ipAccessList)))
			service.SetConditionFromResult(status.IPAccessListReadyType, result)

			return result
		}
	}

	service.SetConditionTrue(status.IPAccessListReadyType)

	if len(akoProject.Spec.ProjectIPAccessList) == 0 {
		service.UnsetCondition(status.IPAccessListReadyType)
	}

	return workflow.OK()
}

func mapToOperatorSpec(projectIPAccessList []admin.NetworkPermissionEntry) []project.IPAccessList {
	ipAccessList := make([]project.IPAccessList, 0, len(projectIPAccessList))

	for _, projectIPAccess := range projectIPAccessList {
		deleteAfterDate := ""
		if deleteAfterDateFromAtlas, ok := projectIPAccess.GetDeleteAfterDateOk(); ok {
			deleteAfterDate = timeutil.FormatISO8601(*deleteAfterDateFromAtlas)
		}

		ipAccessList = append(
			ipAccessList,
			project.IPAccessList{
				AwsSecurityGroup: projectIPAccess.GetAwsSecurityGroup(),
				CIDRBlock:        projectIPAccess.GetCidrBlock(),
				Comment:          projectIPAccess.GetComment(),
				DeleteAfterDate:  deleteAfterDate,
				IPAddress:        projectIPAccess.GetIpAddress(),
			},
		)
	}

	return ipAccessList
}

func syncIPAccessList(service *workflow.Context, projectID string, current, desired []project.IPAccessList) error {
	currentMap := map[string]project.IPAccessList{}
	for _, item := range current {
		currentMap[genIPAccessListKey(item)] = item
	}

	desiredMap := map[string]project.IPAccessList{}
	for _, item := range desired {
		desiredMap[genIPAccessListKey(item)] = item
	}

	for key, ipAccessList := range currentMap {
		if _, ok := desiredMap[key]; ok {
			continue
		}

		_, _, err := service.SdkClient.ProjectIPAccessListApi.DeleteProjectIpAccessList(service.Context, projectID, mapToEntryValue(ipAccessList)).Execute()
		if err != nil {
			return err
		}
	}

	toCreate := make([]admin.NetworkPermissionEntry, 0, len(desired))
	for key, ipAccessList := range desiredMap {
		if _, ok := currentMap[key]; !ok {
			entry := admin.NetworkPermissionEntry{
				AwsSecurityGroup: pointer.SetOrNil(ipAccessList.AwsSecurityGroup, ""),
				CidrBlock:        pointer.SetOrNil(ipAccessList.CIDRBlock, ""),
				Comment:          pointer.SetOrNil(ipAccessList.Comment, ""),
				GroupId:          pointer.SetOrNil(projectID, ""),
				IpAddress:        pointer.SetOrNil(ipAccessList.IPAddress, ""),
			}
			if ipAccessList.DeleteAfterDate != "" {
				deleteAfterDate, err := timeutil.ParseISO8601(ipAccessList.DeleteAfterDate)
				if err != nil {
					return fmt.Errorf("error parsing deleteAfterDate: %w", err)
				}
				entry.SetDeleteAfterDate(deleteAfterDate)
			}
			toCreate = append(toCreate, entry)
		}
	}

	if len(toCreate) == 0 {
		return nil
	}

	_, _, err := service.SdkClient.ProjectIPAccessListApi.CreateProjectIpAccessList(service.Context, projectID, &toCreate).Execute()

	return err
}

func mapToEntryValue(ipAccessList project.IPAccessList) string {
	entry := ""

	switch {
	case ipAccessList.CIDRBlock != "":
		entry = ipAccessList.CIDRBlock
		quads := strings.Split(ipAccessList.CIDRBlock, "/")
		if quads[1] == "32" {
			entry = quads[0]
		}
	case ipAccessList.IPAddress != "":
		ip := strings.Split(ipAccessList.IPAddress, "/")

		entry = ip[0]
	case ipAccessList.AwsSecurityGroup != "":
		entry = ipAccessList.AwsSecurityGroup
	}

	return entry
}

func genIPAccessListKey(ipAccessList project.IPAccessList) string {
	entry := mapToEntryValue(ipAccessList)

	if ipAccessList.DeleteAfterDate != "" {
		entry += "." + ipAccessList.DeleteAfterDate
	}

	return entry
}

func filterActiveIPAccessLists(accessLists []project.IPAccessList) ([]project.IPAccessList, []project.IPAccessList) {
	active := make([]project.IPAccessList, 0)
	expired := make([]project.IPAccessList, 0)
	for _, list := range accessLists {
		if list.DeleteAfterDate != "" {
			// We are ignoring the error as it will never happen due to validation check before
			iso8601, _ := timeutil.ParseISO8601(list.DeleteAfterDate)
			if iso8601.Before(time.Now()) {
				expired = append(expired, list)
				continue
			}
		}
		// Either 'deleteAfterDate' field is not specified or it's higher than the current time
		active = append(active, list)
	}
	return active, expired
}

func canIPAccessListReconcile(ctx context.Context, atlasClient *admin.APIClient, protected bool, akoProject *mdbv1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &mdbv1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	list, _, err := atlasClient.ProjectIPAccessListApi.ListProjectIpAccessLists(ctx, akoProject.ID()).Execute()
	if err != nil {
		return false, err
	}

	if list.GetTotalCount() == 0 {
		return true, nil
	}

	atlasAccessLists := mapToOperatorSpec(list.GetResults())
	if cmp.Equal(atlasAccessLists, latestConfig.ProjectIPAccessList, cmpopts.EquateEmpty()) {
		return true, nil
	}

	return cmp.Equal(akoProject.Spec.ProjectIPAccessList, atlasAccessLists, cmpopts.EquateEmpty()), nil
}
