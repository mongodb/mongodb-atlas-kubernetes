package atlasproject

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
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

func createOrDeletePEInAtlas(ctx *workflow.Context, projectID string, specPEs []project.PrivateEndpoint, statusPEs []status.ProjectPrivateEndpoint) (result workflow.Result) {
	log := ctx.Log

	atlasPeConnections, err := syncPEConnections(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	log.Debugw("Updated PE Connections", "atlasPeConnections", atlasPeConnections, "statusPEs", statusPEs)

	if result := clearOutNotLinkedPEs(ctx.Client, projectID, atlasPeConnections, statusPEs, log); !result.IsOk() {
		return result
	}

	endpointsToCreate := set.Difference(specPEs, statusPEs)
	endpointsToUpdate := set.Intersection(specPEs, statusPEs)
	endpointsToDelete := set.Difference(statusPEs, specPEs)

	log.Debugw("Items to create", "difference", endpointsToCreate)
	log.Debugw("Items to update", "difference", endpointsToUpdate)
	log.Debugw("Items to delete", "difference", endpointsToDelete)

	if result := deletePrivateEndpointsFromAtlas(ctx.Client, projectID, endpointsToDelete, log); !result.IsOk() {
		return result
	}

	newConnections, err := createPeServiceInAtlas(ctx.Client, projectID, endpointsToCreate)
	if err != nil {
		log.Debugw("Failed to create PE Service in Atlas", "error", err)
	}
	ctx.EnsureStatusOption(status.AtlasProjectAddPrivateEnpointsOption(convertAllToStatus(newConnections)))

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

	for _, statusPeService := range convertAllToStatus(atlasPeConnections) {
		if statusPeService.InterfaceEndpointID == "" {
			continue
		}

		interfaceEndpoint, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, string(statusPeService.Provider), statusPeService.ID, statusPeService.InterfaceEndpointID)
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}

		// interfaceEndpoint.Status is for the AZURE and GCP interface endpoints
		if !(interfaceEndpoint.AWSConnectionStatus == "AVAILABLE" || interfaceEndpoint.Status == "AVAILABLE") {
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
		if conn.Status != "AVAILABLE" {
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

	ctx.EnsureStatusOption(status.AtlasProjectUpdatePrivateEnpointsOption(convertAllToStatus(atlasPeConnections)))

	return atlasPeConnections, nil
}

func getAllPrivateEndpoints(client mongodbatlas.Client, projectID string) (result []mongodbatlas.PrivateEndpointConnection, err error) {
	providers := []string{"AWS", "AZURE"}
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
		pe := item.(project.PrivateEndpoint)

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
		operatorPeService := pair[0].(project.PrivateEndpoint)
		statusPeService := pair[1].(status.ProjectPrivateEndpoint)

		if operatorPeService.ID != "" && statusPeService.InterfaceEndpointID == "" {
			interfaceConn, _, err := client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(operatorPeService.Provider), statusPeService.ID, &mongodbatlas.InterfaceEndpointConnection{
				ID:                       operatorPeService.ID,
				PrivateEndpointIPAddress: operatorPeService.IP,
			})
			log.Debugw("AddOnePrivateEndpoint Reply", "interfaceConn", interfaceConn, "err", err)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DeleteAllPrivateEndpoints(ctx *workflow.Context, client mongodbatlas.Client, projectID string, statusPE []status.ProjectPrivateEndpoint, log *zap.SugaredLogger) workflow.Result {
	atlasPeConnections, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	endpointsToDelete := set.Difference(convertAllToStatus(atlasPeConnections), []status.ProjectPrivateEndpoint{})
	log.Debugw("List of endpoints to delete", "endpointsToDelete", endpointsToDelete)
	return deletePrivateEndpointsFromAtlas(client, projectID, endpointsToDelete, log)
}

func clearOutNotLinkedPEs(client mongodbatlas.Client, projectID string, atlasConns []mongodbatlas.PrivateEndpointConnection, statusPEs []status.ProjectPrivateEndpoint, log *zap.SugaredLogger) workflow.Result {
	endpointsWithoutPair := []status.ProjectPrivateEndpoint{}
	endpointsAreDeleting := false
	for _, atlasConn := range atlasConns {
		if atlasConn.Status == "DELETING" {
			endpointsAreDeleting = true
		}

		atlasPE := convertOneToStatus(atlasConn)
		found := false
		for _, statusPE := range statusPEs {
			if atlasPE.ID == statusPE.ID {
				found = true
			}
		}

		if !found {
			endpointsWithoutPair = append(endpointsWithoutPair, atlasPE)
		}
	}

	endpointsToDelete := set.Difference(endpointsWithoutPair, []status.ProjectPrivateEndpoint{})
	log.Debugw("Outdated endpoints to delete", "endpointsToDelete", endpointsToDelete)
	result := deletePrivateEndpointsFromAtlas(client, projectID, endpointsToDelete, log)

	if endpointsAreDeleting {
		return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Endpoints are being deleted")
	}

	return result
}

func deletePrivateEndpointsFromAtlas(client mongodbatlas.Client, projectID string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) workflow.Result {
	result := workflow.OK()
	for _, item := range listsToRemove {
		peService := item.(status.ProjectPrivateEndpoint)
		provider := string(peService.Provider)
		if peService.InterfaceEndpointID != "" {
			if _, err := client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, provider, peService.ID, peService.InterfaceEndpointID); err != nil {
				return workflow.Terminate(workflow.ProjectPrivateEndpointIsNotReadyInAtlas, "failed to delete Private Endpoint")
			}

			return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
		}

		if _, err := client.PrivateEndpoints.Delete(context.Background(), projectID, provider, peService.ID); err != nil {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, "failed to delete Private Endpoint Service")
		}
		log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "id", item.Identifier())
		result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
	}
	return result
}

func convertAllToStatus(peList []mongodbatlas.PrivateEndpointConnection) (result []status.ProjectPrivateEndpoint) {
	for _, endpoint := range peList {
		result = append(result, convertOneToStatus(endpoint))
	}

	return result
}

func convertOneToStatus(endpoint mongodbatlas.PrivateEndpointConnection) status.ProjectPrivateEndpoint {
	pe := status.ProjectPrivateEndpoint{
		ID:       endpoint.ID,
		Provider: provider.ProviderName(endpoint.ProviderName),
		Region:   endpoint.Region,
	}

	switch pe.Provider {
	case provider.ProviderAWS:
		pe.ServiceName = endpoint.EndpointServiceName
		pe.ServiceResourceID = endpoint.ID
		if len(endpoint.InterfaceEndpoints) != 0 {
			pe.InterfaceEndpointID = endpoint.InterfaceEndpoints[0]
		}
	case provider.ProviderAzure:
		pe.ServiceName = endpoint.PrivateLinkServiceName
		pe.ServiceResourceID = endpoint.PrivateLinkServiceResourceID
		if len(endpoint.PrivateEndpoints) != 0 {
			pe.InterfaceEndpointID = endpoint.PrivateEndpoints[0]
		}
	}

	return pe
}
