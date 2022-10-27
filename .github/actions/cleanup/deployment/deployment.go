package deployment

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

func DeleteAllDeployments(ctx context.Context, client mongodbatlas.ClustersService, projectID string) error {
	deployments, err := GetAllDeployments(ctx, client, projectID)
	if err != nil {
		return fmt.Errorf("error getting deployments: %s", err)
	}
	for _, deployment := range deployments {
		if _, err = client.Delete(ctx, projectID, deployment.Name); err != nil {
			return fmt.Errorf("error deleting deployment: %s", err)
		}
	}
	return nil
}

func GetAllDeployments(ctx context.Context, client mongodbatlas.ClustersService, projectID string) ([]mongodbatlas.Cluster, error) {
	deployments, _, err := client.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting deployments: %s", err)
	}
	return deployments, nil
}
