package atlasdeployment

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment/globaldeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func EnsureCustomZoneMapping(service *workflow.Context, groupID string, customZoneMappings []mdbv1.CustomZoneMapping, deploymentName string) workflow.Result {
	result := syncCustomZoneMapping(context.Background(), service, groupID, deploymentName, customZoneMappings)
	if !result.IsOk() {
		service.SetConditionFromResult(status.CustomZoneMappingReadyType, result)
		return result
	}

	if customZoneMappings == nil {
		service.UnsetCondition(status.CustomZoneMappingReadyType)
		service.EnsureStatusOption(status.AtlasDeploymentCustomZoneMappingOption(nil))
	} else {
		service.SetConditionTrue(status.CustomZoneMappingReadyType)
	}

	return result
}

func syncCustomZoneMapping(ctx context.Context, service *workflow.Context, groupID string, deploymentName string, customZoneMappings []mdbv1.CustomZoneMapping) workflow.Result {
	logger := service.Log
	_, existingZoneMapping, err := globaldeployment.GetGlobalDeploymentState(ctx, service.Client, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Sprintf("Failed to get zone mapping state: %v", err))
	}
	logger.Debugf("Existing zone mapping: %v", existingZoneMapping)
	var customZoneMapping status.CustomZoneMapping
	if shouldAdd, shouldDelete := compareZoneMappingStates(existingZoneMapping, customZoneMappings); shouldDelete || shouldAdd {
		skipAdd := false
		if shouldDelete {
			err = deleteZoneMapping(ctx, service.Client, groupID, deploymentName)
			if err != nil {
				skipAdd = true
				logger.Errorf("Failed to sync zone mapping: %v", err)
				customZoneMapping.ZoneMappingErrMessage = fmt.Sprintf("Failed to sync zone mapping: %v", err)
				customZoneMapping.ZoneMappingState = status.StatusFailed
			}
		}

		if shouldAdd && !skipAdd {
			zoneMapping, errRecreate := createZoneMapping(ctx, service.Client, groupID, deploymentName, customZoneMappings)
			if errRecreate != nil {
				logger.Errorf("Failed to sync zone mapping: %v", errRecreate)
				customZoneMapping.ZoneMappingErrMessage = fmt.Sprintf("Failed to sync zone mapping: %v", errRecreate)
				customZoneMapping.ZoneMappingState = status.StatusFailed
			} else {
				logger.Debugf("Zone mapping added: %v", zoneMapping)
				customZoneMapping.ZoneMappingState = status.StatusReady
				customZoneMapping.CustomZoneMapping = zoneMapping
			}
		}
	} else {
		customZoneMapping.ZoneMappingState = status.StatusReady
		customZoneMapping.CustomZoneMapping = existingZoneMapping
	}

	service.EnsureStatusOption(status.AtlasDeploymentCustomZoneMappingOption(&customZoneMapping))
	return checkCustomZoneMapping(customZoneMapping)
}

func deleteZoneMapping(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string) error {
	return globaldeployment.DeleteAllCustomZoneMapping(ctx, client, groupID, deploymentName)
}

func createZoneMapping(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string, mappings []mdbv1.CustomZoneMapping) (map[string]string, error) {
	var atlasMappings []mongodbatlas.CustomZoneMapping
	for _, m := range mappings {
		atlasMappings = append(atlasMappings, m.ToAtlas())
	}
	czm, err := globaldeployment.CreateCustomZoneMapping(ctx, client, groupID, deploymentName, &mongodbatlas.CustomZoneMappingsRequest{CustomZoneMappings: atlasMappings})
	if err != nil {
		return nil, fmt.Errorf("failed to create custom zone mapping: %w", err)
	}
	return czm, nil
}

func checkCustomZoneMapping(customZoneMapping status.CustomZoneMapping) workflow.Result {
	if customZoneMapping.ZoneMappingState != status.StatusReady {
		return workflow.Terminate(workflow.CustomZoneMappingReady, "Global cluster zone mapping is not ready")
	}
	return workflow.OK()
}

func compareZoneMappingStates(existing map[string]string, desired []mdbv1.CustomZoneMapping) (bool, bool) {
	shouldAdd, shouldDelete := false, false

	if len(desired) > len(existing) {
		shouldAdd = true
	} else {
		for _, d := range desired {
			if _, ok := existing[d.Location]; !ok {
				shouldAdd = true
				break
			}
		}
	}

	if len(desired) < len(existing) {
		shouldDelete = true
	} else {
		for k := range existing {
			if _, ok := existing[k]; !ok {
				shouldDelete = true
				break
			}
		}
	}

	return shouldAdd, shouldDelete
}
