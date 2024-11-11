package atlasdeployment

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDeploymentReconciler) ensureCustomZoneMapping(service *workflow.Context, groupID string, customZoneMappings []akov2.CustomZoneMapping, deploymentName string) workflow.Result {
	result := r.syncCustomZoneMapping(service, groupID, deploymentName, customZoneMappings)
	if !result.IsOk() {
		service.SetConditionFromResult(api.CustomZoneMappingReadyType, result)
		return result
	}

	if customZoneMappings == nil {
		service.UnsetCondition(api.CustomZoneMappingReadyType)
		service.EnsureStatusOption(status.AtlasDeploymentCustomZoneMappingOption(nil))
	} else {
		service.SetConditionTrue(api.CustomZoneMappingReadyType)
	}

	return result
}

func (r *AtlasDeploymentReconciler) syncCustomZoneMapping(service *workflow.Context, groupID string, deploymentName string, customZoneMappings []akov2.CustomZoneMapping) workflow.Result {
	logger := service.Log
	err := verifyZoneMapping(customZoneMappings)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, err.Error())
	}
	existingZoneMapping, err := r.deploymentService.GetCustomZones(service.Context, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Sprintf("Failed to get zone mapping state: %v", err))
	}
	logger.Debugf("Existing zone mapping: %v", existingZoneMapping)
	var customZoneMappingStatus status.CustomZoneMapping
	zoneMappingMap, err := r.deploymentService.GetZoneMapping(service.Context, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Sprintf("Failed to get zone mapping map: %v", err))
	}

	if shouldAdd, shouldDelete := compareZoneMappingStates(existingZoneMapping, customZoneMappings, zoneMappingMap); shouldDelete || shouldAdd {
		skipAdd := false
		if shouldDelete {
			err = r.deploymentService.DeleteCustomZones(service.Context, groupID, deploymentName)
			if err != nil {
				skipAdd = true
				logger.Errorf("failed to sync zone mapping: %v", err)
				customZoneMappingStatus.ZoneMappingErrMessage = fmt.Sprintf("Failed to sync zone mapping: %v", err)
				customZoneMappingStatus.ZoneMappingState = status.StatusFailed
			} else {
				logger.Debug("Zone mapping deleted")
				customZoneMappingStatus.ZoneMappingState = status.StatusReady
			}
		}

		if shouldAdd && !skipAdd {
			zoneMapping, err := r.deploymentService.CreateCustomZones(service.Context, groupID, deploymentName, customZoneMappings)
			if err != nil {
				logger.Errorf("failed to sync zone mapping: %v", err)
				customZoneMappingStatus.ZoneMappingErrMessage = fmt.Sprintf("Failed to sync zone mapping: %v", err)
				customZoneMappingStatus.ZoneMappingState = status.StatusFailed
			} else {
				logger.Debugf("Zone mapping added: %v", zoneMapping)
				customZoneMappingStatus.ZoneMappingState = status.StatusReady
				customZoneMappingStatus.CustomZoneMapping = zoneMapping
			}
		}
	} else {
		customZoneMappingStatus.ZoneMappingState = status.StatusReady
		customZoneMappingStatus.CustomZoneMapping = existingZoneMapping
	}

	service.EnsureStatusOption(status.AtlasDeploymentCustomZoneMappingOption(&customZoneMappingStatus))
	return checkCustomZoneMapping(customZoneMappingStatus)
}

func verifyZoneMapping(desired []akov2.CustomZoneMapping) error {
	locations := make(map[string]bool)
	for _, m := range desired {
		if _, ok := locations[m.Location]; ok {
			return fmt.Errorf("duplicate location %v", m.Location)
		} else {
			locations[m.Location] = true
		}
	}
	return nil
}

func checkCustomZoneMapping(customZoneMapping status.CustomZoneMapping) workflow.Result {
	if customZoneMapping.ZoneMappingState != status.StatusReady {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Sprintf("Zone mapping is not ready: %v", customZoneMapping.ZoneMappingErrMessage))
	}
	return workflow.OK()
}

func compareZoneMappingStates(existing map[string]string, desired []akov2.CustomZoneMapping, zoneMappingMap map[string]string) (bool, bool) {
	shouldAdd, shouldDelete := false, false

	if len(desired) < len(existing) {
		shouldDelete = true
	} else {
		for loc, id := range existing {
			found := false
			for _, d := range desired {
				if d.Location == loc && d.Zone == zoneMappingMap[id] {
					found = true
					break
				}
			}
			if !found {
				shouldDelete = true
				break
			}
		}
	}

	if len(desired) > len(existing) || (len(desired) > 0 && shouldDelete) {
		shouldAdd = true
	} else {
		for _, d := range desired {
			zoneID, ok := existing[d.Location]
			if !ok {
				shouldAdd = true
				break
			}
			if zoneName, ok2 := zoneMappingMap[zoneID]; ok2 && zoneName != d.Zone {
				shouldAdd = true
				break
			}
		}
	}

	return shouldAdd, shouldDelete
}
