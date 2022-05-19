package atlasproject

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func (r *AtlasProjectReconciler) ensurePrivateEndpoint(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	specPEs := project.Spec.DeepCopy().PrivateEndpoints
	statusPEs := project.Status.DeepCopy().PrivateEndpoints

	result := createOrDeletePEInAtlas(ctx, projectID, specPEs, statusPEs)
	if !result.IsOk() {
		return result
	}

	return workflow.OK()
}

func createOrDeletePEInAtlas(ctx *workflow.Context, projectID string, specPEs []mdbv1.PrivateEndpoint, statusPEs []status.ProjectPrivateEndpoint) (result workflow.Result) {
	log := ctx.Log

	atlasPeConnections, err := syncPEConnections(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	log.Debugw("Updated PE Connections", "atlasPeConnections", atlasPeConnections, "statusPEs", statusPEs)

	if result := clearOutNotLinkedPEs(ctx, projectID, atlasPeConnections, statusPEs); !result.IsOk() {
		return result
	}

	endpointsToDelete := set.Difference(statusPEs, specPEs)
	log.Debugw("Private Endpoints to delete", "difference", endpointsToDelete)
	if result := deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete); !result.IsOk() {
		return result
	}

	endpointsToCreate := set.Difference(specPEs, statusPEs)
	log.Debugw("Private Endpoints to create", "difference", endpointsToCreate)
	newConnections, err := createPeServiceInAtlas(ctx.Client, projectID, endpointsToCreate)
	if err != nil {
		log.Debugw("Failed to create PE Service in Atlas", "error", err)
	}
	ctx.EnsureStatusOption(status.AtlasProjectAddPrivateEnpointsOption(convertAllToStatus(ctx, projectID, newConnections)))

	endpointsToUpdate := set.Intersection(specPEs, statusPEs)
	log.Debugw("Private Endpoints to update", "difference", endpointsToUpdate)
	if err = createPrivateEndpointInAtlas(ctx.Client, projectID, endpointsToUpdate, log); err != nil {
		log.Debugw("Failed to create PE Interface in Atlas", "error", err)
	}

	return getStatusForInterfaceConnections(ctx, projectID)
}

func getStatusForInterfaceConnections(ctx *workflow.Context, projectID string) workflow.Result {
	atlasPeConnections, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if len(atlasPeConnections) != 0 {
		if allEnpointsAreAvailable(atlasPeConnections) {
			ctx.SetConditionTrue(status.PrivateEndpointServiceReadyType)
		} else {
			result := workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is not ready")
			ctx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, result)
			return result
		}
	}

	for _, statusPeService := range convertAllToStatus(ctx, projectID, atlasPeConnections) {
		if statusPeService.InterfaceEndpointID == "" {
			continue
		}

		interfaceEndpoint, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, string(statusPeService.Provider), statusPeService.ID, statusPeService.InterfaceEndpointID)
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}

		// interfaceEndpoint.Status is for the AZURE and GCP interface endpoints
		if !(isAvailable(interfaceEndpoint.AWSConnectionStatus) || isAvailable(interfaceEndpoint.Status)) {
			result := workflow.InProgress(workflow.ProjectPrivateEndpointIsNotReadyInAtlas, "Interface Private Endpoint is not ready")
			ctx.SetConditionFromResult(status.PrivateEndpointReadyType, result)
			return result
		}

		ctx.SetConditionTrue(status.PrivateEndpointReadyType)
	}

	return workflow.OK()
}

func allEnpointsAreAvailable(atlasPeConnections []mongodbatlas.PrivateEndpointConnection) bool {
	for _, conn := range atlasPeConnections {
		if !isAvailable(conn.Status) {
			return false
		}
	}

	return true
}

func syncPEConnections(ctx *workflow.Context, projectID string) ([]mongodbatlas.PrivateEndpointConnection, error) {
	atlasPeConnections, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return nil, err
	}

	ctx.EnsureStatusOption(status.AtlasProjectUpdatePrivateEnpointsOption(convertAllToStatus(ctx, projectID, atlasPeConnections)))

	return atlasPeConnections, nil
}

func getAllPrivateEndpoints(client mongodbatlas.Client, projectID string) (result []mongodbatlas.PrivateEndpointConnection, err error) {
	providers := []string{"AWS", "AZURE", "GCP"}
	for _, provider := range providers {
		atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, err
		}

		for connIdx := range atlasPeConnections {
			atlasPeConnections[connIdx].ProviderName = provider
		}

		result = append(result, atlasPeConnections...)
	}

	return
}

func createPeServiceInAtlas(client mongodbatlas.Client, projectID string, endpointsToCreate []set.Identifiable) ([]mongodbatlas.PrivateEndpointConnection, error) {
	newConnections := []mongodbatlas.PrivateEndpointConnection{}
	for _, item := range endpointsToCreate {
		pe := item.(mdbv1.PrivateEndpoint)

		conn, _, err := client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		})
		if err != nil {
			return nil, err
		}

		conn.ProviderName = string(pe.Provider)
		conn.Region = pe.Region
		newConnections = append(newConnections, *conn)
	}

	return newConnections, nil
}

func createPrivateEndpointInAtlas(client mongodbatlas.Client, projectID string, endpointsToUpdate [][]set.Identifiable, log *zap.SugaredLogger) error {
	for _, pair := range endpointsToUpdate {
		specPeService := pair[0].(mdbv1.PrivateEndpoint)
		statusPeService := pair[1].(status.ProjectPrivateEndpoint)

		log.Debugw("endpointsAreNotFullyConfigured", "output", endpointsAreNotFullyConfigured(specPeService, statusPeService))
		if endpointsAreNotFullyConfigured(specPeService, statusPeService) {
			interfaceConn := &mongodbatlas.InterfaceEndpointConnection{
				ID:                       specPeService.ID,
				PrivateEndpointIPAddress: specPeService.IP,
				EndpointGroupName:        specPeService.EndpointGroupName,
				GCPProjectID:             specPeService.GCPProjectID,
			}
			if gcpEndpoints, err := specPeService.Endpoints.ConvertToAtlas(); err == nil {
				interfaceConn.Endpoints = gcpEndpoints
			}
			interfaceConn, _, err := client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(specPeService.Provider), statusPeService.ID, interfaceConn)
			log.Debugw("AddOnePrivateEndpoint Reply", "interfaceConn", interfaceConn, "err", err)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func endpointsAreNotFullyConfigured(specPeService mdbv1.PrivateEndpoint, statusPeService status.ProjectPrivateEndpoint) bool {
	awsOrAzureCondition := specPeService.ID != "" && statusPeService.InterfaceEndpointID == ""
	gcpCondition := specPeService.GCPProjectID != "" && specPeService.EndpointGroupName != "" && len(specPeService.Endpoints) != 0 && len(statusPeService.Endpoints) != len(specPeService.Endpoints)
	return awsOrAzureCondition || gcpCondition
}

func DeleteAllPrivateEndpoints(ctx *workflow.Context, client mongodbatlas.Client, projectID string, statusPE []status.ProjectPrivateEndpoint, log *zap.SugaredLogger) workflow.Result {
	atlasPeConnections, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	endpointsToDelete := set.Difference(convertAllToStatus(ctx, projectID, atlasPeConnections), []status.ProjectPrivateEndpoint{})
	log.Debugw("List of endpoints to delete", "endpointsToDelete", endpointsToDelete)
	return deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete)
}

func clearOutNotLinkedPEs(ctx *workflow.Context, projectID string, atlasConns []mongodbatlas.PrivateEndpointConnection, statusPEs []status.ProjectPrivateEndpoint) workflow.Result {
	log := ctx.Log
	endpointsWithoutPair := []status.ProjectPrivateEndpoint{}
	endpointsAreDeleting := false
	for _, atlasConn := range atlasConns {
		if !isAvailable(atlasConn.Status) {
			continue
		}

		atlasPE := convertOneToStatus(ctx, projectID, atlasConn)
		found := false
		for _, statusPE := range statusPEs {
			if atlasPE.ID == statusPE.ID {
				found = true
			}
		}

		if !found {
			endpointsWithoutPair = append(endpointsWithoutPair, atlasPE)
			endpointsAreDeleting = true
		}
	}

	endpointsToDelete := set.Difference(endpointsWithoutPair, []status.ProjectPrivateEndpoint{})
	log.Debugw("Outdated endpoints to delete", "endpointsToDelete", endpointsToDelete)
	result := deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete)

	if endpointsAreDeleting {
		return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Endpoints are being deleted")
	}

	return result
}

func deletePrivateEndpointsFromAtlas(ctx *workflow.Context, projectID string, listsToRemove []set.Identifiable) workflow.Result {
	log := ctx.Log
	result := workflow.OK()
	for _, item := range listsToRemove {
		peService := item.(status.ProjectPrivateEndpoint)
		provider := string(peService.Provider)
		if peService.InterfaceEndpointID != "" {
			if _, err := ctx.Client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, provider, peService.ID, peService.InterfaceEndpointID); err != nil {
				return workflow.Terminate(workflow.ProjectPrivateEndpointIsNotReadyInAtlas, "failed to delete Private Endpoint")
			}

			return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
		}

		if _, err := ctx.Client.PrivateEndpoints.Delete(context.Background(), projectID, provider, peService.ID); err != nil {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, "failed to delete Private Endpoint Service")
		}
		log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "id", item.Identifier())
		result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
	}
	return result
}

func convertAllToStatus(ctx *workflow.Context, projectID string, peList []mongodbatlas.PrivateEndpointConnection) (result []status.ProjectPrivateEndpoint) {
	for _, endpoint := range peList {
		result = append(result, convertOneToStatus(ctx, projectID, endpoint))
	}

	return result
}

func convertOneToStatus(ctx *workflow.Context, projectID string, conn mongodbatlas.PrivateEndpointConnection) status.ProjectPrivateEndpoint {
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

func isAvailable(status string) bool {
	return status == "AVAILABLE"
}
