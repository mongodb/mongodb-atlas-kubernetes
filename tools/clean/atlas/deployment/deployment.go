package deployment

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/atlas/mongodbatlas"
)

func DeleteAllDeployments(ctx context.Context, client mongodbatlas.ClustersService, projectID string) error {
	deployments, err := GetAllDeployments(ctx, client, projectID)
	if err != nil {
		return fmt.Errorf("error getting deployments: %s", err)
	}
	var allErr error
	for _, deployment := range deployments {
		log.Printf("Deleting deployment %s", deployment.Name)
		if _, err = client.Delete(ctx, projectID, deployment.Name, &mongodbatlas.DeleteAdvanceClusterOptions{}); err != nil {
			allErr = errors.Join(allErr, fmt.Errorf("error deleting deployment: %s", err))
		}
	}
	return allErr
}

func GetAllDeployments(ctx context.Context, client mongodbatlas.ClustersService, projectID string) ([]mongodbatlas.Cluster, error) {
	deployments, _, err := client.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting deployments: %s", err)
	}
	return deployments, nil
}
