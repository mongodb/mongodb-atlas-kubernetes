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

package atlasdeployment

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
)

func (r *AtlasDeploymentReconciler) ensureCustomZoneMapping(service *workflow.Context, deploymentService deployment.AtlasDeploymentsService, groupID string, customZoneMappings []akov2.CustomZoneMapping, deploymentName string) workflow.Result {
	result := r.syncCustomZoneMapping(service, deploymentService, groupID, deploymentName, customZoneMappings)
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

func (r *AtlasDeploymentReconciler) syncCustomZoneMapping(service *workflow.Context, deploymentService deployment.AtlasDeploymentsService, groupID string, deploymentName string, customZoneMappings []akov2.CustomZoneMapping) workflow.Result {
	logger := service.Log
	err := verifyZoneMapping(customZoneMappings)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, err)
	}
	existingZoneMapping, err := deploymentService.GetCustomZones(service.Context, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Errorf("failed to get zone mapping state: %w", err))
	}
	logger.Debugf("Existing zone mapping: %v", existingZoneMapping)
	var customZoneMappingStatus status.CustomZoneMapping
	zoneMappingMap, err := deploymentService.GetZoneMapping(service.Context, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Errorf("failed to get zone mapping map: %w", err))
	}

	if shouldAdd, shouldDelete := compareZoneMappingStates(existingZoneMapping, customZoneMappings, zoneMappingMap); shouldDelete || shouldAdd {
		skipAdd := false
		if shouldDelete {
			err = deploymentService.DeleteCustomZones(service.Context, groupID, deploymentName)
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
			zoneMapping, err := deploymentService.CreateCustomZones(service.Context, groupID, deploymentName, customZoneMappings)
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
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Errorf("zone mapping is not ready: %v", customZoneMapping.ZoneMappingErrMessage))
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
