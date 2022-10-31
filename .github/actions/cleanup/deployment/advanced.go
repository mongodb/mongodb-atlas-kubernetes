package deployment

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/atlas/mongodbatlas"
)

func getAllAdvancedClusters(ctx context.Context, client mongodbatlas.AdvancedClustersService, projectID string) ([]*mongodbatlas.AdvancedCluster, error) {
	advancedClusters, _, err := client.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting advanced clusters: %w", err)
	}
	return advancedClusters.Results, nil
}

func DeleteAllAdvancedClusters(ctx context.Context, client mongodbatlas.AdvancedClustersService, projectID string) error {
	advancedClusters, err := getAllAdvancedClusters(ctx, client, projectID)
	if err != nil {
		return fmt.Errorf("error getting advanced clusters: %w", err)
	}
	for _, cluster := range advancedClusters {
		log.Printf("Deleting advanced cluster %s", cluster.Name)
		if _, err = client.Delete(ctx, projectID, cluster.Name); err != nil {
			return fmt.Errorf("error deleting advanced cluster: %w", err)
		}
	}
	return nil
}
