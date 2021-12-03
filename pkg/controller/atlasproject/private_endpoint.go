package atlasproject

import (
	"context"
	"regexp"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func (r *AtlasProjectReconciler) ensurePrivateEndpoint(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	peList := project.Spec.DeepCopy().PrivateEndpoints
	result, privateEndpoints := createOrDeletePEInAtlas(ctx.Client, projectID, peList, ctx.Log)
	if !result.IsOk() {
		return result
	}

	peStatus, err := convertToStatus(privateEndpoints)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	ctx.EnsureStatusOption(status.AtlasProjectPrivateEnpointOption(peStatus))
	return workflow.OK()
}

func createOrDeletePEInAtlas(client mongodbatlas.Client, projectID string, operatorPrivateEndpoints []project.PrivateEndpoint, log *zap.SugaredLogger) (result workflow.Result, privateEndpoints []atlasProjectPrivateEndpoint) {
	provider := "AWS"
	atlasPrivateEndpoints, err := getAtlasPrivateEndpoint(client, projectID, provider)
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

	if err := deletePrivateEndpointsFromAtlas(client, projectID, provider, endpointsToDelete, log); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error()), nil
	}

	if result := createPeServiceInAtlas(client, projectID, endpointsToCreate); !result.IsOk() {
		return result, nil
	}

	if result := createPrivateEndpointInAtlas(client, projectID, endpointsToUpdate); !result.IsOk() {
		return result, nil
	}

	return workflow.OK(), atlasPrivateEndpoints
}

type atlasProjectPrivateEndpoint mongodbatlas.PrivateEndpointConnection

func (pe atlasProjectPrivateEndpoint) Identifier() interface{} {
	return pe.ProviderName + pe.Region
}

func getAtlasPrivateEndpoint(client mongodbatlas.Client, projectID string, provider string) (result []atlasProjectPrivateEndpoint, err error) {
	atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Making a new slice with synonyms as Atlas IP Access list to enable usage of 'Identifiable'
	result = make([]atlasProjectPrivateEndpoint, len(atlasPeConnections))
	for i, r := range atlasPeConnections {
		result[i] = atlasProjectPrivateEndpoint(r)
	}

	return fillProviderAndRegion(result), nil
}

func createPeServiceInAtlas(client mongodbatlas.Client, projectID string, endpointsToCreate []set.Identifiable) workflow.Result {
	for _, item := range endpointsToCreate {
		pe := item.(project.PrivateEndpoint)

		if _, _, err := client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		}); err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}
	}

	return workflow.OK()
}

func createPrivateEndpointInAtlas(client mongodbatlas.Client, projectID string, endpointsToUpdate [][]set.Identifiable) workflow.Result {
	for _, pair := range endpointsToUpdate {
		operatorPeService := pair[0].(project.PrivateEndpoint)
		atlasPeService := pair[1].(atlasProjectPrivateEndpoint)

		if operatorPeService.ID != "" {
			_, _, err := client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(operatorPeService.Provider), atlasPeService.ID, &mongodbatlas.InterfaceEndpointConnection{
				ID: operatorPeService.ID,
			})
			if err != nil {
				return workflow.InProgress(workflow.ProjectPrivateEndpointNotCreatedInAtlas, err.Error())
			}

			return workflow.OK()
		}
	}

	return workflow.OK()
}

func DeleteAllPrivateEndpoints(client mongodbatlas.Client, projectID string, log *zap.SugaredLogger) error {
	provider := "AWS"
	atlasPrivateEndpoints, err := getAtlasPrivateEndpoint(client, projectID, provider)
	if err != nil {
		return err
	}

	allEndpoints := set.Difference(atlasPrivateEndpoints, []atlasProjectPrivateEndpoint{})

	if err := deletePrivateEndpointsFromAtlas(client, projectID, provider, allEndpoints, log); err != nil {
		return err
	}

	return nil
}

func deletePrivateEndpointsFromAtlas(client mongodbatlas.Client, projectID string, provider string, listsToRemove []set.Identifiable, log *zap.SugaredLogger) error {
	for _, item := range listsToRemove {
		log.Debugw("item to delete", "l", item)
		peService := item.(atlasProjectPrivateEndpoint)
		if len(peService.PrivateEndpoints) != 0 {
			for _, endpointID := range peService.PrivateEndpoints {
				if _, err := client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, provider, peService.ID, endpointID); err != nil {
					return err
				}
			}
		}

		if _, err := client.PrivateEndpoints.Delete(context.Background(), projectID, provider, peService.ID); err != nil {
			return err
		}
		log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "id", item.Identifier())
	}
	return nil
}

func fillProviderAndRegion(atlasPeConnections []atlasProjectPrivateEndpoint) (result []atlasProjectPrivateEndpoint) {
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

func filterOutCreatingEndpoints(endpointsToCreate []set.Identifiable, atlasPrivateEndpoints []atlasProjectPrivateEndpoint, log *zap.SugaredLogger) (result []set.Identifiable) {
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

func matchAtlasPE(operatorServices []set.Identifiable, atlasPrivateEndpoints []atlasProjectPrivateEndpoint) (result []operatorAtlasServicePair) {
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
	atlas    *atlasProjectPrivateEndpoint
}

func endpointsMatch(endpoint1 set.Identifiable, endpoint2 set.Identifiable) bool {
	return endpoint1.Identifier() == endpoint2.Identifier()
}

func convertToStatus(atlasPEs []atlasProjectPrivateEndpoint) (result []status.PrivateEndpoint, err error) {
	err = compat.JSONCopy(&result, atlasPEs)
	return result, err
}
