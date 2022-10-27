package deployment

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

func DeleteAllServerless(ctx context.Context, client mongodbatlas.ServerlessInstancesService, projectID string) error {
	serverless, err := GetAllServerless(ctx, client, projectID)
	if err != nil {
		return fmt.Errorf("error getting serverless: %s", err)
	}
	for _, s := range serverless {
		if _, err = client.Delete(ctx, projectID, s.Name); err != nil {
			return fmt.Errorf("error deleting serverless: %s", err)
		}
	}
	return nil
}

func GetAllServerless(ctx context.Context, client mongodbatlas.ServerlessInstancesService, projectID string) ([]*mongodbatlas.Cluster, error) {
	serverless, _, err := client.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting serverless: %s", err)
	}
	return serverless.Results, nil
}
