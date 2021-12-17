package atlasproject

import (
	"context"
	"errors"

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

	atlasPeConnections, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	ctx.EnsureStatusOption(status.AtlasProjectRemoveOldPrivateEnpointOption(convertToStatus(atlasPeConnections)))

	endpointsToCreate := set.Difference(specPEs, statusPEs)
	endpointsToUpdate := set.Intersection(specPEs, statusPEs)
	endpointsToDelete := set.Difference(statusPEs, specPEs)

	log.Debugw("Items to create", "difference", endpointsToCreate)
	log.Debugw("Items to update", "difference", endpointsToUpdate)
	log.Debugw("Items to delete", "difference", endpointsToDelete)

	if result := deletePrivateEndpointsFromAtlas(ctx.Client, projectID, endpointsToDelete, log); !result.IsOk() {
		return result
	}

	result, newConnections := createPeServiceInAtlas(ctx.Client, projectID, endpointsToCreate, log)
	ctx.EnsureStatusOption(status.AtlasProjectAddPrivateEnpointsOption(convertToStatus(newConnections)))

	atlasPeConnections, err = getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	ctx.EnsureStatusOption(status.AtlasProjectUpdatePrivateEnpointsOption(convertToStatus(atlasPeConnections)))
	if len(atlasPeConnections) != 0 {
		if allEnpointsAreAvailable(atlasPeConnections) {
			ctx.SetConditionTrue(status.PrivateEndpointServiceReadyType)
		} else {
			result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is not ready")
			ctx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, result)
			return result
		}
	}

	if !result.IsOk() {
		ctx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, result)
		return result
	}

	result, updatedConnections := createPrivateEndpointInAtlas(ctx.Client, projectID, endpointsToUpdate, log)
	ctx.EnsureStatusOption(status.AtlasProjectUpdatePrivateEnpointsOption(convertToStatus(updatedConnections)))
	log.Debugw("Updated Private Enpoints", "updatedConnections", updatedConnections)
	for _, pair := range endpointsToUpdate {
		operatorPeService := pair[0].(project.PrivateEndpoint)
		statusPeService := pair[1].(status.ProjectPrivateEndpoint)
		log.Debugw("Pair Vals", "operatorPeService.ID", operatorPeService.ID)
		if operatorPeService.ID == "" {
			continue
		}

		interfaceEndpoint, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, string(operatorPeService.Provider), statusPeService.ServiceResourceID, operatorPeService.ID)
		if err != nil {
			return workflow.InProgress(workflow.Internal, err.Error())
		}

		log.Debugw("Interface Endpoint", "interfaceEndpoint", interfaceEndpoint)
		if interfaceEndpoint.AWSConnectionStatus == "AVAILABLE" || interfaceEndpoint.AzureStatus == "AVAILABLE" {
			ctx.SetConditionTrue(status.PrivateEndpointReadyType)
			continue
		}

		result = workflow.InProgress(workflow.ProjectPrivateEndpointIsNotReadyInAtlas, "Interface Private Endpoint is not ready")
		ctx.SetConditionFromResult(status.PrivateEndpointReadyType, result)
	}

	return result
}

func allEnpointsAreAvailable(atlasPeConnections []mongodbatlas.PrivateEndpointConnection) bool {
	for _, conn := range atlasPeConnections {
		if conn.Status != "AVAILABLE" {
			return false
		}
	}

	return true
}

func getAllPrivateEndpoints(client mongodbatlas.Client, projectID string) (result []mongodbatlas.PrivateEndpointConnection, err error) {
	providers := []string{"AWS", "AZURE"}
	for _, provider := range providers {
		atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, err
		}

		result = append(result, atlasPeConnections...)
	}

	return
}

func createPeServiceInAtlas(client mongodbatlas.Client, projectID string, endpointsToCreate []set.Identifiable, log *zap.SugaredLogger) (workflow.Result, []mongodbatlas.PrivateEndpointConnection) {
	result := workflow.OK()
	newConnections := []mongodbatlas.PrivateEndpointConnection{}
	for _, item := range endpointsToCreate {
		pe := item.(project.PrivateEndpoint)

		conn, resp, err := client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		})
		log.Debug("Reply from Atlas", "resp", resp, "err", err, "conn", conn)
		if err != nil {
			conn = &mongodbatlas.PrivateEndpointConnection{}
		}
		conn.ProviderName = string(pe.Provider)
		conn.Region = pe.Region
		newConnections = append(newConnections, *conn)
		result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is creating")
	}

	return result, newConnections
}

func createPrivateEndpointInAtlas(client mongodbatlas.Client, projectID string, endpointsToUpdate [][]set.Identifiable, log *zap.SugaredLogger) (workflow.Result, []mongodbatlas.PrivateEndpointConnection) {
	result := workflow.OK()
	updatedConnections := []mongodbatlas.PrivateEndpointConnection{}
	for _, pair := range endpointsToUpdate {
		operatorPeService := pair[0].(project.PrivateEndpoint)
		statusPeService := pair[1].(status.ProjectPrivateEndpoint)

		log.Debugw("Pair CREATE", "operatorPeService", operatorPeService, "statusPeService", statusPeService)
		if operatorPeService.ID != "" && statusPeService.InterfaceEndpointID == "" {
			interfaceConn, _, err := client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(operatorPeService.Provider), statusPeService.ServiceResourceID, &mongodbatlas.InterfaceEndpointConnection{
				ID: operatorPeService.ID,
			})
			updatedConnections = append(updatedConnections, mongodbatlas.PrivateEndpointConnection{
				ID:               statusPeService.ServiceResourceID,
				PrivateEndpoints: []string{operatorPeService.ID},
			})
			log.Debugw("AddOnePrivateEndpoint Reply", "interfaceConn", interfaceConn, "err", err)
			if interfaceConn == nil {
				return workflow.InProgress(workflow.ProjectPrivateEndpointIsNotReadyInAtlas, err.Error()), nil
			}

			result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is creating")
		}
	}

	return result, updatedConnections
}

func DeleteAllPrivateEndpoints(client mongodbatlas.Client, projectID string, statusPEs []status.ProjectPrivateEndpoint, log *zap.SugaredLogger) error {
	endpointsToDelete := set.Difference(statusPEs, []status.ProjectPrivateEndpoint{})
	if result := deletePrivateEndpointsFromAtlas(client, projectID, endpointsToDelete, log); !result.IsOk() {
		return errors.New("failed to delete Private Endpoints from Atlas")
	}

	return nil
}

func deletePrivateEndpointsFromAtlas(client mongodbatlas.Client, projectID string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) workflow.Result {
	result := workflow.OK()
	for _, item := range listsToRemove {
		peService := item.(status.ProjectPrivateEndpoint)
		provider := string(peService.Provider)
		if peService.InterfaceEndpointID != "" {
			if _, err := client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, provider, peService.ServiceResourceID, peService.InterfaceEndpointID); err != nil {
				return workflow.Terminate(workflow.ProjectPrivateEndpointIsNotReadyInAtlas, "failed to delete Private Endpoint")
			}

			return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
		}

		if _, err := client.PrivateEndpoints.Delete(context.Background(), projectID, provider, peService.ServiceResourceID); err != nil {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, "failed to delete Private Endpoint Service")
		}
		log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "id", item.Identifier())
		result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
	}
	return result
}

func convertToStatus(peList []mongodbatlas.PrivateEndpointConnection) (result []status.ProjectPrivateEndpoint) {
	for _, endpoint := range peList {
		pe := status.ProjectPrivateEndpoint{
			Provider:          provider.ProviderName(endpoint.ProviderName),
			Region:            endpoint.Region,
			ServiceName:       endpoint.EndpointServiceName,
			ServiceResourceID: endpoint.ID,
		}

		if len(endpoint.InterfaceEndpoints) != 0 {
			pe.InterfaceEndpointID = endpoint.InterfaceEndpoints[0] // AWS
		}
		if len(endpoint.PrivateEndpoints) != 0 {
			pe.InterfaceEndpointID = endpoint.PrivateEndpoints[0] // AZURE
		}

		result = append(result, pe)
	}

	return result
}
