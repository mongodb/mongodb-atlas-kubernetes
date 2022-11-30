package atlasdeployment

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/stringutil"

	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

const (
	SPEStatusDeleting = "DELETING"

	SPEStatusReserved             = "RESERVED"              // 1
	SPEStatusReservationRequested = "RESERVATION_REQUESTED" // 1

	SPEStatusAvailable  = "AVAILABLE"  // 2
	SPEStatusInitiating = "INITIATING" // 2
	SPEStatusFailed     = "FAILED"     // 2

)

func ensureServerlessPrivateEndpoints(service *workflow.Context, groupID string, deploymentSpec *mdbv1.ServerlessSpec, deploymentName string) workflow.Result {
	if deploymentSpec == nil {
		return workflow.Terminate(workflow.ServerlessPrivateEndpointReady, "deployment spec is empty")
	}
	if result := IsClusterSupportServerlessPE(deploymentSpec); !result.IsOk() {
		return result
	}

	if deploymentSpec.PrivateEndpoints == nil {
		service.UnsetCondition(status.ServerlessPrivateEndpointReadyType)
		service.EnsureStatusOption(status.AtlasDeploymentServerlessPrivateEndpointOption(nil))
		return workflow.OK()
	}

	result := syncServerlessPrivateEndpoints(context.Background(), service, groupID, deploymentName, provider.ProviderName(deploymentSpec.ProviderSettings.BackingProviderName), deploymentSpec.PrivateEndpoints)
	if !result.IsOk() {
		service.SetConditionFromResult(status.ServerlessPrivateEndpointReadyType, result)
		return result
	}
	service.SetConditionTrue(status.ServerlessPrivateEndpointReadyType)
	return result
}

func IsClusterSupportServerlessPE(deploymentSpec *mdbv1.ServerlessSpec) workflow.Result {
	if deploymentSpec.ProviderSettings.BackingProviderName != string(provider.ProviderAWS) &&
		deploymentSpec.ProviderSettings.BackingProviderName != string(provider.ProviderAzure) {
		return workflow.Terminate(workflow.Internal, "serverless private endpoints are only supported for AWS and Azure")
	}

	return workflow.OK()
}

func syncServerlessPrivateEndpoints(ctx context.Context, service *workflow.Context, groupID, deploymentName string, providerName provider.ProviderName, desiredPE []mdbv1.ServerlessPrivateEndpoint) workflow.Result {
	logger := service.Log
	client := service.Client.ServerlessPrivateEndpoints
	logger.Debugf("Syncing serverless private endpoints for deployment %s", deploymentName)
	existingPE, err := getAllExistingServerlessPE(ctx, client, groupID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.ServerlessPrivateEndpointReady, err.Error())
	}
	logger.Debugf("Existing serverless private endpoints: %v", existingPE)
	diff := sortServerlessPE(logger, existingPE, desiredPE)
	logger.Debugf("Serverless private endpoints diff: %v", diff)
	speStatuses := createSPE(ctx, logger, client, groupID, deploymentName, diff.PEToCreate)
	speStatuses = append(speStatuses, connectSPE(ctx, logger, client, groupID, deploymentName, providerName, diff.PEToConnect)...)
	speStatuses = append(speStatuses, getStatusFromReadySPE(diff.PEToUpdateStatus)...)
	speStatuses = append(speStatuses, handleDuplicatePE(diff)...)
	service.EnsureStatusOption(status.AtlasDeploymentServerlessPrivateEndpointOption(speStatuses))
	logger.Debugf("Serverless Private Endpoints statuses: %v", speStatuses)
	errors := deleteSPE(ctx, client, groupID, deploymentName, diff.PEToDelete)
	if len(errors) > 0 {
		return workflow.Terminate(workflow.ServerlessPrivateEndpointReady, fmt.Sprintf("failed to delete serverless private endpoints: %v", errors))
	}
	return checkStatuses(speStatuses)
}

func handleDuplicatePE(diff *SPEDiff) []status.ServerlessPrivateEndpoint {
	var result []status.ServerlessPrivateEndpoint
	for _, pe := range diff.DuplicateToCreate {
		result = append(result, status.FailedDuplicationSPE(pe.Name, pe.CloudProviderEndpointID, pe.PrivateEndpointIPAddress))
	}
	return result
}

func checkStatuses(pe []status.ServerlessPrivateEndpoint) workflow.Result {
	for _, p := range pe {
		if p.Status != SPEStatusAvailable {
			return workflow.Terminate(workflow.ServerlessPrivateEndpointReady, fmt.Sprintf("serverless private endpoint %s is not ready", p.ID))
		}
	}
	return workflow.OK()
}

func deleteSPE(ctx context.Context, client mongodbatlas.ServerlessPrivateEndpointsService, groupID, deploymentName string, peToDelete []string) []error {
	var result []error
	for _, id := range peToDelete {
		_, err := client.Delete(ctx, groupID, deploymentName, id)
		if err != nil {
			result = append(result, fmt.Errorf("failed to delete serverless private endpoint: %w", err))
		}
	}
	return result
}

func getStatusFromReadySPE(pe []mongodbatlas.ServerlessPrivateEndpointConnection) []status.ServerlessPrivateEndpoint {
	var result []status.ServerlessPrivateEndpoint
	for _, endpoint := range pe {
		result = append(result, status.SPEFromAtlas(endpoint))
	}
	return result
}

func connectSPE(ctx context.Context, logger *zap.SugaredLogger, client mongodbatlas.ServerlessPrivateEndpointsService, groupID, deploymentName string, providerName provider.ProviderName, pe []mongodbatlas.ServerlessPrivateEndpointConnection) []status.ServerlessPrivateEndpoint {
	var result []status.ServerlessPrivateEndpoint
	for _, endpoint := range pe {
		id := endpoint.ID
		req := mongodbatlas.ServerlessPrivateEndpointConnection{
			PrivateEndpointIPAddress: endpoint.PrivateEndpointIPAddress,
			CloudProviderEndpointID:  endpoint.CloudProviderEndpointID,
			ProviderName:             string(providerName),
		}
		logger.Debugf("Connecting serverless private endpoint %s", id)
		resultPE, _, err := client.Update(ctx, groupID, deploymentName, id, &req)
		if err != nil {
			logger.Errorf("Failed to connect serverless private endpoint %s: %v", id, err)
			result = append(result, status.FailedToConnectSPE(endpoint, fmt.Sprintf("failed to connect serverless private endpoint: %s", err)))
		} else {
			logger.Debugf("Serverless private endpoint %s is connected", id)
			result = append(result, status.SPEFromAtlas(*resultPE))
		}
	}
	return result
}

func createSPE(ctx context.Context, logger *zap.SugaredLogger, client mongodbatlas.ServerlessPrivateEndpointsService, groupID, deploymentName string, pe []mdbv1.ServerlessPrivateEndpoint) []status.ServerlessPrivateEndpoint {
	var result []status.ServerlessPrivateEndpoint
	for _, endpoint := range pe {
		created, _, err := client.Create(ctx, groupID, deploymentName, prepareForCreation(endpoint))
		if err != nil {
			logger.Errorf("Failed to create serverless private endpoint: %v, err: %v", prepareForCreation(endpoint), err)
			result = append(result, status.FailedToCreateSPE(endpoint.Name, fmt.Sprintf("failed to create serverless private endpoint: %s", err)))
		} else {
			result = append(result, status.SPEFromAtlas(*created))
		}
	}
	return result
}

func prepareForCreation(pe mdbv1.ServerlessPrivateEndpoint) *mongodbatlas.ServerlessPrivateEndpointConnection {
	return &mongodbatlas.ServerlessPrivateEndpointConnection{
		Comment: pe.Name,
	}
}

type SPEDiff struct {
	PEToCreate        []mdbv1.ServerlessPrivateEndpoint
	PEToConnect       []mongodbatlas.ServerlessPrivateEndpointConnection
	PEToUpdateStatus  []mongodbatlas.ServerlessPrivateEndpointConnection
	PEToDelete        []string
	DuplicateToCreate []mdbv1.ServerlessPrivateEndpoint
}

func (d *SPEDiff) appendToCreate(pe mdbv1.ServerlessPrivateEndpoint) {
	for _, p := range d.PEToCreate {
		if p.Name == pe.Name {
			d.DuplicateToCreate = append(d.DuplicateToCreate, pe)
			return
		}
	}
	d.PEToCreate = append(d.PEToCreate, pe)
}

func sortServerlessPE(logger *zap.SugaredLogger, existedPE []mongodbatlas.ServerlessPrivateEndpointConnection, desiredPE []mdbv1.ServerlessPrivateEndpoint) *SPEDiff {
	existingPEToCreate := make([]mongodbatlas.ServerlessPrivateEndpointConnection, 0)
	existingReadyPE := make([]mongodbatlas.ServerlessPrivateEndpointConnection, 0)
	existingReservedPE := make([]mongodbatlas.ServerlessPrivateEndpointConnection, 0)

	desiredPEToCreate := make([]mdbv1.ServerlessPrivateEndpoint, 0)
	desiredReadyPE := make([]mdbv1.ServerlessPrivateEndpoint, 0)

	for _, pe := range existedPE {
		switch pe.Status {
		case SPEStatusInitiating, SPEStatusReservationRequested:
			existingPEToCreate = append(existingPEToCreate, pe)
		case SPEStatusReserved:
			existingReservedPE = append(existingReservedPE, pe)
		case SPEStatusFailed, SPEStatusAvailable:
			existingReadyPE = append(existingReadyPE, pe)
		case SPEStatusDeleting:

		default:
			logger.Errorf("Unknown status %s for serverless private endpoint %s", pe.Status, pe.ID)
		}
	}

	logger.Debugf("Existing serverless private endpoints to connect: %v", existingPEToCreate)
	logger.Debugf("Existing ready serverless private endpoints: %v", existingReadyPE)

	for _, pe := range desiredPE {
		if pe.IsInitialState() {
			desiredPEToCreate = append(desiredPEToCreate, pe)
		} else {
			desiredReadyPE = append(desiredReadyPE, pe)
		}
	}

	for _, ePE := range existingReservedPE {
		ready := false
		for _, dPE := range desiredReadyPE {
			if dPE.Name == ePE.Comment {
				ready = true
				existingReadyPE = append(existingReadyPE, ePE)
			}
		}
		if !ready {
			existingPEToCreate = append(existingPEToCreate, ePE)
		}
	}

	logger.Debugf("Desired serverless private endpoints to connect: %v", desiredPEToCreate)
	logger.Debugf("Desired ready serverless private endpoints: %v", desiredReadyPE)

	secondDiff := sortReadySPE(existingReadyPE, desiredReadyPE)
	logger.Debugf("Second diff:  %v", secondDiff)

	var uniqueNotAvailableNames []string
	for _, pe := range secondDiff.PEToConnect {
		uniqueNotAvailableNames = append(uniqueNotAvailableNames, pe.Comment)
	}

	firstDiff := sortSPEToConnect(existingPEToCreate, desiredPEToCreate, uniqueNotAvailableNames)
	logger.Debugf("First diff:  %v", firstDiff)
	mergedDiff := &SPEDiff{
		PEToCreate:        append(firstDiff.PEToCreate, secondDiff.PEToCreate...),
		PEToConnect:       append(firstDiff.PEToConnect, secondDiff.PEToConnect...),
		PEToUpdateStatus:  append(firstDiff.PEToUpdateStatus, secondDiff.PEToUpdateStatus...),
		PEToDelete:        append(firstDiff.PEToDelete, secondDiff.PEToDelete...),
		DuplicateToCreate: append(firstDiff.DuplicateToCreate, secondDiff.DuplicateToCreate...),
	}

	return mergedDiff
}

func preparePEForConnection(atlasPE mongodbatlas.ServerlessPrivateEndpointConnection, pe mdbv1.ServerlessPrivateEndpoint) mongodbatlas.ServerlessPrivateEndpointConnection {
	return mongodbatlas.ServerlessPrivateEndpointConnection{
		ID:                       atlasPE.ID,
		Comment:                  pe.Name,
		CloudProviderEndpointID:  pe.CloudProviderEndpointID,
		PrivateEndpointIPAddress: pe.PrivateEndpointIPAddress,
		ProviderName:             atlasPE.ProviderName,
	}
}

func sortReadySPE(existingPEs []mongodbatlas.ServerlessPrivateEndpointConnection, desiredPEs []mdbv1.ServerlessPrivateEndpoint) *SPEDiff {
	var result SPEDiff

	for _, desiredPE := range desiredPEs {
		found := false
		for _, existingPE := range existingPEs {
			if isReadySPEEqual(existingPE, desiredPE) {
				result.PEToUpdateStatus = append(result.PEToUpdateStatus, existingPE)
				found = true
				break
			} else if existingPE.Comment == desiredPE.Name && existingPE.Status == SPEStatusReserved {
				result.PEToConnect = append(result.PEToConnect, preparePEForConnection(existingPE, desiredPE))
				found = true
				break
			}
		}
		if !found {
			result.appendToCreate(desiredPE)
		}
	}

	for _, existingPE := range existingPEs {
		toDelete := true
		for _, desiredPE := range result.PEToConnect {
			if existingPE.ID == desiredPE.ID {
				toDelete = false
				break
			}
		}
		for _, desiredPE := range result.PEToUpdateStatus {
			if existingPE.ID == desiredPE.ID {
				toDelete = false
				break
			}
		}
		if toDelete {
			result.PEToDelete = append(result.PEToDelete, existingPE.ID)
		}
	}

	return &result
}

func isReadySPEEqual(existingPE mongodbatlas.ServerlessPrivateEndpointConnection, desiredPE mdbv1.ServerlessPrivateEndpoint) bool {
	return existingPE.Comment == desiredPE.Name && desiredPE.CloudProviderEndpointID == existingPE.CloudProviderEndpointID && desiredPE.PrivateEndpointIPAddress == existingPE.PrivateEndpointIPAddress
}

func sortSPEToConnect(existingPEs []mongodbatlas.ServerlessPrivateEndpointConnection, desiredPEs []mdbv1.ServerlessPrivateEndpoint, uniqueComments []string) *SPEDiff {
	var result SPEDiff
	for _, desiredPE := range desiredPEs {
		if stringutil.Contains(uniqueComments, desiredPE.Name) {
			result.DuplicateToCreate = append(result.DuplicateToCreate, desiredPE)
			continue
		}
		
		uniqueComments = append(uniqueComments, desiredPE.Name)

		toCreate := true
		for _, existingPE := range existingPEs {
			if desiredPE.Name == existingPE.Comment {
				toCreate = false
				result.PEToUpdateStatus = append(result.PEToUpdateStatus, existingPE)
				break
			}
		}
		if toCreate {
			result.appendToCreate(desiredPE)
		}
	}

	for _, existingPE := range existingPEs {
		toDelete := true
		for _, desiredPE := range result.PEToUpdateStatus {
			if existingPE.ID == desiredPE.ID {
				toDelete = false
				break
			}
		}
		if toDelete {
			result.PEToDelete = append(result.PEToDelete, existingPE.ID)
		}
	}

	return &result
}

func getAllExistingServerlessPE(ctx context.Context, service mongodbatlas.ServerlessPrivateEndpointsService, groupID, clusterName string) ([]mongodbatlas.ServerlessPrivateEndpointConnection, error) {
	list, _, err := service.List(ctx, groupID, clusterName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list serverless private endpoints: %w", err)
	}
	return list, nil
}
