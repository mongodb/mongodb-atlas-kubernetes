package atlasproject

import (
	"context"
	"errors"
	"regexp"

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
	peList := project.Spec.DeepCopy().PrivateEndpoints
	result, atlasPEs := createOrDeletePEInAtlas(ctx, projectID, peList, ctx.Log)
	if !result.IsOk() {
		return result
	}

	peStatus := convertToStatus(atlasPEs)
	ctx.Log.Debugw("Private Endpoints Status", "status", peStatus)

	ctx.EnsureStatusOption(status.AtlasProjectPrivateEnpointOption(peStatus))
	return workflow.OK()
}

func createOrDeletePEInAtlas(ctx *workflow.Context, projectID string, operatorPrivateEndpoints []project.PrivateEndpoint, log *zap.SugaredLogger) (result workflow.Result, atlasPrivateEndpoints []AtlasProjectPrivateEndpoint) {
	provider := "AWS"
	atlasPrivateEndpoints, err := getAtlasPrivateEndpoint(ctx.Client, projectID, provider)
	if err != nil {
		workflow.Terminate(workflow.ProjectPrivateEndpointNotCreatedInAtlas, err.Error())
	}

	endpointsToCreate := set.Difference(operatorPrivateEndpoints, atlasPrivateEndpoints)
	endpointsToUpdate := set.Intersection(operatorPrivateEndpoints, atlasPrivateEndpoints)
	endpointsToDelete := set.Difference(atlasPrivateEndpoints, operatorPrivateEndpoints)

	endpointsToCreate = filterOutCreatingEndpoints(endpointsToCreate, atlasPrivateEndpoints, log)

	log.Debugw("Items to create", "difference", endpointsToCreate)
	log.Debugw("Items to update", "difference", endpointsToUpdate)
	log.Debugw("Items to delete", "difference", endpointsToDelete)

	if result := deletePrivateEndpointsFromAtlas(ctx.Client, projectID, provider, endpointsToDelete, log); !result.IsOk() {
		return result, nil
	}

	if result := createPeServiceInAtlas(ctx.Client, projectID, endpointsToCreate); !result.IsOk() {
		ctx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, result)
		return result, nil
	}
	ctx.SetConditionTrue(status.PrivateEndpointServiceReadyType)

	if len(endpointsToUpdate) != 0 {
		if result := createPrivateEndpointInAtlas(ctx.Client, projectID, endpointsToUpdate); !result.IsOk() {
			ctx.SetConditionFromResult(status.PrivateEndpointReadyType, result)
			return result, nil
		}

		ctx.SetConditionTrue(status.PrivateEndpointReadyType)
	}

	return workflow.OK(), atlasPrivateEndpoints
}

type AtlasProjectPrivateEndpoint mongodbatlas.PrivateEndpointConnection

func (pe AtlasProjectPrivateEndpoint) Identifier() interface{} {
	return pe.ProviderName + pe.Region
}

func getAtlasPrivateEndpoint(client mongodbatlas.Client, projectID string, provider string) (result []AtlasProjectPrivateEndpoint, err error) {
	atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Making a new slice with synonyms as Atlas IP Access list to enable usage of 'Identifiable'
	result = make([]AtlasProjectPrivateEndpoint, len(atlasPeConnections))
	for i, r := range atlasPeConnections {
		result[i] = AtlasProjectPrivateEndpoint(r)
	}

	return fillProviderAndRegion(result), nil
}

func createPeServiceInAtlas(client mongodbatlas.Client, projectID string, endpointsToCreate []set.Identifiable) workflow.Result {
	result := workflow.OK()
	for _, item := range endpointsToCreate {
		pe := item.(project.PrivateEndpoint)

		if _, _, err := client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		}); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}

		result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is creating")
	}

	return result
}

func createPrivateEndpointInAtlas(client mongodbatlas.Client, projectID string, endpointsToUpdate [][]set.Identifiable) workflow.Result {
	result := workflow.OK()
	for _, pair := range endpointsToUpdate {
		operatorPeService := pair[0].(project.PrivateEndpoint)
		atlasPeService := pair[1].(AtlasProjectPrivateEndpoint)

		if operatorPeService.ID != "" {
			_, _, err := client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(operatorPeService.Provider), atlasPeService.ID, &mongodbatlas.InterfaceEndpointConnection{
				ID: operatorPeService.ID,
			})
			if err != nil {
				return workflow.InProgress(workflow.ProjectPrivateEndpointNotCreatedInAtlas, err.Error())
			}

			result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is creating")
		}
	}

	return result
}

func DeleteAllPrivateEndpoints(client mongodbatlas.Client, projectID string, log *zap.SugaredLogger) error {
	provider := "AWS"
	atlasPrivateEndpoints, err := getAtlasPrivateEndpoint(client, projectID, provider)
	if err != nil {
		return err
	}

	allEndpoints := set.Difference(atlasPrivateEndpoints, []AtlasProjectPrivateEndpoint{})

	if result := deletePrivateEndpointsFromAtlas(client, projectID, provider, allEndpoints, log); !result.IsOk() {
		return errors.New("failed to delete Private Endpoints from Atlas")
	}

	return nil
}

func deletePrivateEndpointsFromAtlas(client mongodbatlas.Client, projectID string, provider string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) workflow.Result {
	result := workflow.OK()
	for _, item := range listsToRemove {
		log.Debugw("item to delete", "l", item)
		peService := item.(AtlasProjectPrivateEndpoint)
		if len(peService.PrivateEndpoints) != 0 {
			for _, endpointID := range peService.PrivateEndpoints {
				if _, err := client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, provider, peService.ID, endpointID); err != nil {
					return workflow.Terminate(workflow.ProjectPrivateEndpointNotCreatedInAtlas, "failed to delete Private Endpoint")
				}
			}
		}

		if _, err := client.PrivateEndpoints.Delete(context.Background(), projectID, provider, peService.ID); err != nil {
			return workflow.Terminate(workflow.ProjectPrivateEndpointNotCreatedInAtlas, "failed to delete Private Endpoint Service")
		}
		log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "id", item.Identifier())
		result = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
	}
	return result
}

func fillProviderAndRegion(atlasPeConnections []AtlasProjectPrivateEndpoint) (result []AtlasProjectPrivateEndpoint) {
	awsRegex := regexp.MustCompile(`^com\.(\w+).vpce\.([\w-]+)\.`)
	for _, conn := range atlasPeConnections {
		// try to parse AWS
		matches := awsRegex.FindStringSubmatch(conn.EndpointServiceName)
		if len(matches) != 3 {
			continue
		}

		if matches[1] == "amazonaws" {
			conn.ProviderName = "AWS"
			conn.Region = matches[2]
			result = append(result, conn)
		}
	}

	return result
}

func filterOutCreatingEndpoints(endpointsToCreate []set.Identifiable, atlasPrivateEndpoints []AtlasProjectPrivateEndpoint, log *zap.SugaredLogger) (result []set.Identifiable) {
	log.Debugw("Endpoints to match", "endpointsToCreate", endpointsToCreate, "atlasPrivateEndpoints", atlasPrivateEndpoints)
	matchedPeServicePairs := matchAtlasPE(endpointsToCreate, atlasPrivateEndpoints)
	for _, item := range matchedPeServicePairs {
		log.Debugw("Matched PE service pair", "operator", item.operator, "atlas", item.atlas)
		if item.atlas == nil {
			result = append(result, *item.operator)
		}
	}

	return result
}

func matchAtlasPE(operatorServices []set.Identifiable, atlasPrivateEndpoints []AtlasProjectPrivateEndpoint) (result []operatorAtlasServicePair) {
	for i := range operatorServices {
		operatorServices := operatorServices[i]
		found := false
		for j := range atlasPrivateEndpoints {
			atlasService := atlasPrivateEndpoints[j]
			if endpointsMatch(operatorServices, atlasService) {
				found = true
				result = append(result, operatorAtlasServicePair{
					operator: &operatorServices,
					atlas:    &atlasService,
				})

				break
			}
		}

		if !found {
			result = append(result, operatorAtlasServicePair{
				operator: &operatorServices,
				atlas:    nil,
			})
		}
	}

	return result
}

type operatorAtlasServicePair struct {
	operator *set.Identifiable
	atlas    *AtlasProjectPrivateEndpoint
}

func endpointsMatch(endpoint1 set.Identifiable, endpoint2 set.Identifiable) bool {
	return endpoint1.Identifier() == endpoint2.Identifier()
}

func convertToStatus(peList []AtlasProjectPrivateEndpoint) (result []status.ProjectPrivateEndpoint) {
	for _, endpoint := range peList {
		result = append(result, status.ProjectPrivateEndpoint{
			Provider:          provider.ProviderName(endpoint.ProviderName),
			Region:            endpoint.Region,
			ServiceName:       endpoint.EndpointServiceName,
			ServiceResourceID: endpoint.ID,
		})
	}

	return result
}
