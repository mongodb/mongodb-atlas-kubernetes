package deployment

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/atlas/mongodbatlas"
)

func DeleteAllServerless(ctx context.Context, client mongodbatlas.ServerlessInstancesService, projectID string) error {
	serverless, err := GetAllServerless(ctx, client, projectID)
	if err != nil {
		return fmt.Errorf("error getting serverless: %s", err)
	}
	var allErr error
	for _, s := range serverless {
		log.Printf("Deleting serverless %s", s.Name)
		if _, err = client.Delete(ctx, projectID, s.Name); err != nil {
			allErr = errors.Join(allErr, fmt.Errorf("error deleting serverless: %s", err))
		}
	}
	return allErr
}

func GetAllServerless(ctx context.Context, client mongodbatlas.ServerlessInstancesService, projectID string) ([]*mongodbatlas.Cluster, error) {
	serverless, _, err := client.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting serverless: %s", err)
	}
	return serverless.Results, nil
}
