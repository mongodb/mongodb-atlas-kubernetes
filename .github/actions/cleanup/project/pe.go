package project

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/atlas/mongodbatlas"
)

func deleteAllPE(ctx context.Context, client mongodbatlas.PrivateEndpointsService, projectID string) error {
	connections, err := getAllPrivateEndpoints(ctx, client, projectID)
	if err != nil {
		return fmt.Errorf("error getting private endpoints: %s", err)
	}
	err = deletePrivateEndpointsFromAtlas(ctx, client, projectID, connections)
	if err != nil {
		return fmt.Errorf("error deleting private endpoints: %s", err)
	}
	return nil
}

func deletePrivateEndpointsFromAtlas(ctx context.Context, client mongodbatlas.PrivateEndpointsService, projectID string, listsToRemove []mongodbatlas.PrivateEndpointConnection) error {
	for _, peService := range listsToRemove {
		if firstInterfaceEndpointID(peService) != "" {
			log.Printf("Deleting private endpoint %s", firstInterfaceEndpointID(peService))
			if _, err := client.DeleteOnePrivateEndpoint(ctx, projectID, peService.ProviderName, peService.ID, firstInterfaceEndpointID(peService)); err != nil {
				return fmt.Errorf("error deleting private endpoint interface: %s", err)
			}

			continue
		}

		log.Printf("Deleting private endpoint %s", peService.EndpointServiceName)
		if _, err := client.Delete(ctx, projectID, peService.ProviderName, peService.ID); err != nil {
			return fmt.Errorf("error deleting private endpoint service: %s", err)
		}
	}
	return nil
}

func firstInterfaceEndpointID(connection mongodbatlas.PrivateEndpointConnection) string {
	if len(connection.InterfaceEndpoints) != 0 {
		return connection.InterfaceEndpoints[0]
	}

	if len(connection.PrivateEndpoints) != 0 {
		return connection.PrivateEndpoints[0]
	}

	if len(connection.EndpointGroupNames) != 0 {
		return connection.EndpointGroupNames[0]
	}

	return ""
}

func getAllPrivateEndpoints(ctx context.Context, client mongodbatlas.PrivateEndpointsService, projectID string) ([]mongodbatlas.PrivateEndpointConnection, error) {
	providers := []string{ProviderAzure, ProviderGCP, ProviderAWS}
	var result []mongodbatlas.PrivateEndpointConnection
	for _, provider := range providers {
		atlasPeConnections, _, err := client.List(ctx, projectID, provider, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, err
		}

		for connIdx := range atlasPeConnections {
			atlasPeConnections[connIdx].ProviderName = provider
		}

		for _, atlasPeConnection := range atlasPeConnections {
			result = append(result, atlasPeConnection)
		}
	}
	return result, nil
}
