package atlasproject

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func ensurePrivateEndpoint(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	specPEs := project.Spec.DeepCopy().PrivateEndpoints

	atlasPEs, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	result, conditionType := syncPrivateEndpointsWithAtlas(ctx, projectID, specPEs, atlasPEs)
	if !result.IsOk() {
		if conditionType == status.PrivateEndpointServiceReadyType {
			ctx.UnsetCondition(status.PrivateEndpointReadyType)
		}
		ctx.SetConditionFromResult(conditionType, result)
		return result
	}

	if len(specPEs) == 0 && len(atlasPEs) == 0 {
		ctx.UnsetCondition(status.PrivateEndpointServiceReadyType)
		ctx.UnsetCondition(status.PrivateEndpointReadyType)
		return workflow.OK()
	}

	serviceStatus, endpointsAreFullyDefined := getStatusForServices(ctx, atlasPEs)
	ctx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, serviceStatus)
	if !endpointsAreFullyDefined {
		ctx.UnsetCondition(status.PrivateEndpointReadyType)
		return serviceStatus
	}

	interfaceStatus := getStatusForInterfaces(ctx, atlasPEs, projectID)
	ctx.SetConditionFromResult(status.PrivateEndpointReadyType, interfaceStatus)

	return interfaceStatus
}

func syncPrivateEndpointsWithAtlas(ctx *workflow.Context, projectID string, specPEs []mdbv1.PrivateEndpoint, atlasPEs []atlasPE) (workflow.Result, status.ConditionType) {
	log := ctx.Log

	log.Debugw("PE Connections", "atlasPEs", atlasPEs, "specPEs", specPEs)

	endpointsToDelete := set.Difference(atlasPEs, specPEs)
	log.Debugw("Private Endpoints to delete", "difference", endpointsToDelete)
	if result := deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete); !result.IsOk() {
		return result, status.PrivateEndpointServiceReadyType
	}

	endpointsToCreate := set.Difference(specPEs, atlasPEs)
	log.Debugw("Private Endpoints to create", "difference", endpointsToCreate)
	newConnections, err := createPeServiceInAtlas(ctx, projectID, endpointsToCreate)
	if err != nil {
		return terminateWithError(ctx, status.PrivateEndpointServiceReadyType, "Failed to create PE Service in Atlas", err)
	}

	endpointsToSync := set.Intersection(specPEs, atlasPEs)
	log.Debugw("Private Endpoints to sync", "difference", endpointsToSync)
	syncedConnections, err := syncPeInterfaceInAtlas(ctx, projectID, endpointsToSync)
	if err != nil {
		return terminateWithError(ctx, status.PrivateEndpointReadyType, "Failed to sync PE Interface in Atlas", err)
	}

	log.Debugw("PE Changes", "newConnections", newConnections, "syncedConnections", syncedConnections)
	updatePEStatusOption(ctx, projectID, newConnections, syncedConnections)

	return determineResult(newConnections, syncedConnections, atlasPEs)
}

func determineResult(newConnections, syncedConnections, atlasPEs []atlasPE) (workflow.Result, status.ConditionType) {
	if len(newConnections) != 0 {
		return notReadyServiceResult, status.PrivateEndpointServiceReadyType
	}

	if len(syncedConnections) != len(atlasPEs) {
		return notReadyInterfaceResult, status.PrivateEndpointReadyType
	}

	return workflow.OK(), status.PrivateEndpointReadyType
}

func getStatusForServices(ctx *workflow.Context, atlasPEs []atlasPE) (result workflow.Result, allDefined bool) {
	allAvailable, allDefined, failureMessage := areServicesAvailableOrFailed(atlasPEs)
	ctx.Log.Debugw("Get Status for Services", "allAvailable", allAvailable, "failureMessage", failureMessage)
	if failureMessage != "" {
		result = workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, failureMessage)
		return
	}
	if !allAvailable {
		result = notReadyServiceResult
		return
	}
	if !allDefined {
		result = workflow.OK().WithMessage("Awaiting additional user configuration")
		return
	}

	return
}

func getStatusForInterfaces(ctx *workflow.Context, atlasPEs []atlasPE, projectID string) workflow.Result {
	for _, statusPeService := range convertAllToStatus(ctx, projectID, atlasPEs) {
		if statusPeService.InterfaceEndpointID == "" {
			return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "InterfaceEndpointID is empty")
		}

		interfaceEndpoint, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, string(statusPeService.Provider), statusPeService.ID, statusPeService.InterfaceEndpointID)
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}

		interfaceIsAvailable, interfaceFailureMessage := checkIfInterfaceIsAvailable(interfaceEndpoint)
		if interfaceFailureMessage != "" {
			return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, interfaceFailureMessage)
		}
		if !interfaceIsAvailable {
			return notReadyInterfaceResult
		}
	}

	return workflow.OK()
}

func areServicesAvailableOrFailed(atlasPeConnections []atlasPE) (allAvailable, allDefined bool, failureMessage string) {
	allAvailable = true
	allDefined = true

	for _, conn := range atlasPeConnections {
		if conn.InterfaceEndpointID() == "" {
			allDefined = false
		}
		if isFailed(conn.Status) {
			failureMessage = conn.ErrorMessage
			return
		}
		if !isAvailable(conn.Status) {
			allAvailable = false
		}
	}

	return
}

func updatePEStatusOption(ctx *workflow.Context, projectID string, newConnections, syncedConnections []atlasPE) {
	setPEStatusOption(ctx, projectID, syncedConnections)
	addPEStatusOption(ctx, projectID, newConnections)
}

func addPEStatusOption(ctx *workflow.Context, projectID string, newPEs []atlasPE) {
	statusPEs := convertAllToStatus(ctx, projectID, newPEs)
	ctx.EnsureStatusOption(status.AtlasProjectAddPrivateEnpointsOption(statusPEs))
}

func setPEStatusOption(ctx *workflow.Context, projectID string, atlasPeConnections []atlasPE) {
	statusPEs := convertAllToStatus(ctx, projectID, atlasPeConnections)
	ctx.EnsureStatusOption(status.AtlasProjectSetPrivateEnpointsOption(statusPEs))
}

type atlasPE mongodbatlas.PrivateEndpointConnection

func (a atlasPE) Identifier() interface{} {
	return a.ProviderName + status.TransformRegionToID(a.RegionName)
}

func (a atlasPE) InterfaceEndpointID() string {
	if len(a.InterfaceEndpoints) != 0 {
		return a.InterfaceEndpoints[0]
	}

	if len(a.PrivateEndpoints) != 0 {
		return a.PrivateEndpoints[0]
	}

	return ""
}

func (a atlasPE) EndpointGroupName() string {
	if len(a.EndpointGroupNames) == 0 {
		return ""
	}

	return a.EndpointGroupNames[0]
}

func getAllPrivateEndpoints(client mongodbatlas.Client, projectID string) (result []atlasPE, err error) {
	providers := []string{"AWS", "AZURE", "GCP"}
	for _, provider := range providers {
		atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, err
		}

		for connIdx := range atlasPeConnections {
			atlasPeConnections[connIdx].ProviderName = provider
		}

		for _, atlasPeConnection := range atlasPeConnections {
			result = append(result, atlasPE(atlasPeConnection))
		}
	}

	return
}

func createPeServiceInAtlas(ctx *workflow.Context, projectID string, endpointsToCreate []set.Identifiable) (newConnections []atlasPE, err error) {
	newConnections = make([]atlasPE, 0)
	for _, item := range endpointsToCreate {
		pe := item.(mdbv1.PrivateEndpoint)

		conn, _, err := ctx.Client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		})
		if err != nil {
			return newConnections, err
		}

		conn.ProviderName = string(pe.Provider)
		conn.Region = pe.Region
		newConnections = append(newConnections, atlasPE(*conn))
	}

	return newConnections, nil
}

func syncPeInterfaceInAtlas(ctx *workflow.Context, projectID string, endpointsToUpdate [][]set.Identifiable) (syncedEndpoints []atlasPE, err error) {
	syncedEndpoints = make([]atlasPE, 0)
	for _, pair := range endpointsToUpdate {
		specPeService := pair[0].(mdbv1.PrivateEndpoint)
		atlasPeService := pair[1].(atlasPE)

		ctx.Log.Debugw("endpointNeedsUpdating", "output", endpointNeedsUpdating(specPeService, atlasPeService))
		if endpointNeedsUpdating(specPeService, atlasPeService) {
			interfaceConn := &mongodbatlas.InterfaceEndpointConnection{
				ID:                       specPeService.ID,
				PrivateEndpointIPAddress: specPeService.IP,
				EndpointGroupName:        specPeService.EndpointGroupName,
				GCPProjectID:             specPeService.GCPProjectID,
			}
			if gcpEndpoints, err := specPeService.Endpoints.ConvertToAtlas(); err == nil {
				interfaceConn.Endpoints = gcpEndpoints
			}

			interfaceConn, response, err := ctx.Client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(specPeService.Provider), atlasPeService.ID, interfaceConn)
			ctx.Log.Debugw("AddOnePrivateEndpoint Reply", "interfaceConn", interfaceConn, "err", err)
			if err != nil {
				if response.StatusCode == 400 || response.StatusCode == 409 {
					ctx.Log.Debugw("failed to create PE Interface", "error", err)
					continue
				}

				return syncedEndpoints, err
			}
		}

		atlasPeService.ProviderName = string(specPeService.Provider)
		atlasPeService.Region = specPeService.Region
		syncedEndpoints = append(syncedEndpoints, atlasPeService)
	}

	return
}

func endpointNeedsUpdating(specPeService mdbv1.PrivateEndpoint, atlasPeService atlasPE) bool {
	if isAvailable(atlasPeService.Status) {
		switch specPeService.Provider {
		case provider.ProviderAWS, provider.ProviderAzure:
			return specPeService.ID != "" && specPeService.ID != atlasPeService.InterfaceEndpointID()
		case provider.ProviderGCP:
			return specPeService.EndpointGroupName != "" && specPeService.EndpointGroupName != atlasPeService.EndpointGroupName() && len(atlasPeService.EndpointServiceName) != 0
		}
	}

	return false
}

func DeleteAllPrivateEndpoints(ctx *workflow.Context, projectID string) workflow.Result {
	atlasPEs, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	endpointsToDelete := set.Difference(atlasPEs, []atlasPE{})
	ctx.Log.Debugw("List of endpoints to delete", "endpointsToDelete", endpointsToDelete)
	return deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete)
}

func deletePrivateEndpointsFromAtlas(ctx *workflow.Context, projectID string, listsToRemove []set.Identifiable) workflow.Result {
	if len(listsToRemove) == 0 {
		return workflow.OK()
	}

	for _, item := range listsToRemove {
		peService := item.(atlasPE)
		if isDeleting(peService.Status) {
			ctx.Log.Debugf("%s Private Endpoint Service for the region %s is being deleted", peService.ProviderName, peService.RegionName)
			continue
		}

		if peService.InterfaceEndpointID() != "" {
			if _, err := ctx.Client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, peService.ProviderName, peService.ID, peService.InterfaceEndpointID()); err != nil {
				return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "failed to delete Private Endpoint")
			}

			continue
		}

		if _, err := ctx.Client.PrivateEndpoints.Delete(context.Background(), projectID, peService.ProviderName, peService.ID); err != nil {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, "failed to delete Private Endpoint Service")
		}

		ctx.Log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "provider", peService.ProviderName, "regionName", peService.RegionName)
	}

	return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
}

func convertAllToStatus(ctx *workflow.Context, projectID string, peList []atlasPE) (result []status.ProjectPrivateEndpoint) {
	for _, endpoint := range peList {
		result = append(result, convertOneToStatus(ctx, projectID, endpoint))
	}

	return result
}

func convertOneToStatus(ctx *workflow.Context, projectID string, conn atlasPE) status.ProjectPrivateEndpoint {
	pe := status.ProjectPrivateEndpoint{
		ID:       conn.ID,
		Provider: provider.ProviderName(conn.ProviderName),
		Region:   conn.Region,
	}

	switch pe.Provider {
	case provider.ProviderAWS:
		pe.ServiceName = conn.EndpointServiceName
		pe.ServiceResourceID = conn.ID
		if len(conn.InterfaceEndpoints) != 0 {
			pe.InterfaceEndpointID = conn.InterfaceEndpoints[0]
		}
	case provider.ProviderAzure:
		pe.ServiceName = conn.PrivateLinkServiceName
		pe.ServiceResourceID = conn.PrivateLinkServiceResourceID
		if len(conn.PrivateEndpoints) != 0 {
			pe.InterfaceEndpointID = conn.PrivateEndpoints[0]
		}
	case provider.ProviderGCP:
		pe.ServiceAttachmentNames = conn.ServiceAttachmentNames
		if len(conn.EndpointGroupNames) != 0 {
			var err error
			pe.InterfaceEndpointID = conn.EndpointGroupNames[0]
			pe.Endpoints, err = getGCPInterfaceEndpoint(ctx, projectID, pe)
			if err != nil {
				ctx.Log.Warnw("failed to get Interface Endpoint Data for GCP", "err", err, "pe", pe)
			}
		}
	}
	ctx.Log.Debugw("Converted Status", "status", pe, "connection", conn)

	return pe
}

// getGCPInterfaceEndpoint returns an InterfaceEndpointID and a list of GCP endpoints
func getGCPInterfaceEndpoint(ctx *workflow.Context, projectID string, endpoint status.ProjectPrivateEndpoint) ([]status.GCPEndpoint, error) {
	log := ctx.Log
	if endpoint.InterfaceEndpointID == "" {
		return nil, errors.New("InterfaceEndpointID is empty")
	}
	interfaceEndpointConn, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, string(provider.ProviderGCP), endpoint.ID, endpoint.InterfaceEndpointID)
	if err != nil {
		return nil, err
	}

	interfaceConns := interfaceEndpointConn.Endpoints
	listOfInterfaces := make([]status.GCPEndpoint, 0)
	for _, e := range interfaceConns {
		endpoint := status.GCPEndpoint{
			Status:       e.Status,
			EndpointName: e.EndpointName,
			IPAddress:    e.IPAddress,
		}
		listOfInterfaces = append(listOfInterfaces, endpoint)
	}
	log.Debugw("Result of getGCPEndpointData", "endpoint.ID", endpoint.ID, "listOfInterfaces", listOfInterfaces)

	return listOfInterfaces, nil
}

// checkIfInterfaceIsAvailable checks if an interface and all of its nested endpoints are available and also returns an error message
func checkIfInterfaceIsAvailable(interfaceEndpointConn *mongodbatlas.InterfaceEndpointConnection) (allAvailable bool, failureMessage string) {
	allAvailable = true

	if isFailed(interfaceEndpointConn.Status) {
		return false, interfaceEndpointConn.ErrorMessage
	}
	if !isAvailable(interfaceEndpointConn.Status) && !isAvailable(interfaceEndpointConn.AWSConnectionStatus) {
		allAvailable = false
	}

	for _, endpoint := range interfaceEndpointConn.Endpoints {
		if isFailed(endpoint.Status) {
			return false, interfaceEndpointConn.ErrorMessage
		}
		if !isAvailable(endpoint.Status) {
			allAvailable = false
		}
	}

	return
}

func isAvailable(status string) bool {
	return status == "AVAILABLE"
}

func isDeleting(status string) bool {
	return status == "DELETING"
}

func isFailed(status string) bool {
	return status == "FAILED"
}

func terminateWithError(ctx *workflow.Context, conditionType status.ConditionType, message string, err error) (workflow.Result, status.ConditionType) {
	ctx.Log.Debugw(message, "error", err)
	result := workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, err.Error()).WithoutRetry()
	return result, conditionType
}

var notReadyServiceResult = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is not ready")
var notReadyInterfaceResult = workflow.InProgress(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "Interface Private Endpoint is not ready")
