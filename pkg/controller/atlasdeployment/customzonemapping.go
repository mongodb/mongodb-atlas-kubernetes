package atlasdeployment

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

	"go.mongodb.org/atlas/mongodbatlas"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func EnsureCustomZoneMapping(service *workflow.Context, groupID string, customZoneMappings []akov2.CustomZoneMapping, deploymentName string) workflow.Result {
	result := syncCustomZoneMapping(service, groupID, deploymentName, customZoneMappings)
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

func syncCustomZoneMapping(service *workflow.Context, groupID string, deploymentName string, customZoneMappings []akov2.CustomZoneMapping) workflow.Result {
	logger := service.Log
	err := verifyZoneMapping(customZoneMappings)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, err.Error())
	}
	_, existingZoneMapping, err := GetGlobalDeploymentState(service.Context, service.Client.GlobalClusters, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Sprintf("Failed to get zone mapping state: %v", err))
	}
	logger.Debugf("Existing zone mapping: %v", existingZoneMapping)
	var customZoneMappingStatus status.CustomZoneMapping
	zoneMappingMap, err := getZoneMappingMap(service.Context, service.Client, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.CustomZoneMappingReady, fmt.Sprintf("Failed to get zone mapping map: %v", err))
	}

	if shouldAdd, shouldDelete := compareZoneMappingStates(existingZoneMapping, customZoneMappings, zoneMappingMap); shouldDelete || shouldAdd {
		skipAdd := false
		if shouldDelete {
			err = deleteZoneMapping(service.Context, service.Client.GlobalClusters, groupID, deploymentName)
			if err != nil {
				skipAdd = true
				logger.Errorf("failed to sync zone mapping: %v", err)
				customZoneMappingStatus.ZoneMappingErrMessage = fmt.Sprintf("Failed to sync zone mapping: %v", err)
				customZoneMappingStatus.ZoneMappingState = status.StatusFailed
			}
		}

		if shouldAdd && !skipAdd {
			zoneMapping, errRecreate := createZoneMapping(service.Context, service.Client.GlobalClusters, groupID, deploymentName, customZoneMappings)
			if errRecreate != nil {
				logger.Errorf("failed to sync zone mapping: %v", errRecreate)
				customZoneMappingStatus.ZoneMappingErrMessage = fmt.Sprintf("Failed to sync zone mapping: %v", errRecreate)
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

func deleteZoneMapping(ctx context.Context, client mongodbatlas.GlobalClustersService, groupID string, deploymentName string) error {
	_, _, err := client.DeleteCustomZoneMappings(ctx, groupID, deploymentName)
	if err != nil {
		return fmt.Errorf("failed to delete custom zone mapping: %w", err)
	}
	return nil
}

func createZoneMapping(ctx context.Context, client mongodbatlas.GlobalClustersService, groupID string, deploymentName string, mappings []akov2.CustomZoneMapping) (map[string]string, error) {
	var atlasMappings []mongodbatlas.CustomZoneMapping
	for _, m := range mappings {
		atlasMappings = append(atlasMappings, m.ToAtlas())
	}
	gc, _, err := client.AddCustomZoneMappings(ctx, groupID, deploymentName, &mongodbatlas.CustomZoneMappingsRequest{CustomZoneMappings: atlasMappings})
	if err != nil {
		return nil, fmt.Errorf("failed to create custom zone mapping: %w", err)
	}
	return gc.CustomZoneMapping, nil
}

func checkCustomZoneMapping(customZoneMapping status.CustomZoneMapping) workflow.Result {
	if customZoneMapping.ZoneMappingState != status.StatusReady {
		return workflow.Terminate(workflow.CustomZoneMappingReady, "Zone mapping is not ready")
	}
	return workflow.OK()
}

func getZoneMappingMap(ctx context.Context, client *mongodbatlas.Client, groupID, clusterName string) (map[string]string, error) {
	cluster, _, err := client.AdvancedClusters.Get(ctx, groupID, clusterName)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(cluster.ReplicationSpecs))
	for _, rc := range cluster.ReplicationSpecs {
		result[rc.ID] = rc.ZoneName
	}
	return result, nil
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

func GetGlobalDeploymentState(ctx context.Context, client mongodbatlas.GlobalClustersService, groupID string, deploymentName string) ([]mongodbatlas.ManagedNamespace, map[string]string, error) {
	deployment, _, err := client.Get(ctx, groupID, deploymentName)
	if err != nil {
		return nil, nil, err
	}
	return deployment.ManagedNamespaces, deployment.CustomZoneMapping, nil
}
