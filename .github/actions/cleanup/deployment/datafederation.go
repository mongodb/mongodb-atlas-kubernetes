package deployment

import (
	"context"
	"fmt"
	"go.mongodb.org/atlas/mongodbatlas"
	"log"
)

func DeleteAllDataFederationInstances(ctx context.Context, client mongodbatlas.DataFederationService, projectID string) error {
	dfInstances, _, err := client.List(ctx, projectID)
	if err != nil {
		return fmt.Errorf("error listing datafederation instances: %w", err)
	}

	for _, df := range dfInstances {
		log.Printf("Removing DataFederation instance: %s", df.Name)
		if _, err := client.Delete(ctx, projectID, df.Name); err != nil {
			return fmt.Errorf("unable to remove DataFederation instance: %w", err)
		}
	}
	return nil
}
