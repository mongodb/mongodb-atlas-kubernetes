package atlas

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"
)

type IPAccessListStatus func(ctx context.Context, projectID, entryValue string) (string, error)

func CustomIPAccessListStatus(client *mongodbatlas.Client) IPAccessListStatus {
	type ipAccessListStatus struct {
		Status string `json:"STATUS"`
	}

	return func(ctx context.Context, projectID, entryValue string) (string, error) {
		urlStr := fmt.Sprintf("/api/atlas/v1.0/groups/%s/accessList/%s/status", projectID, entryValue)
		req, err := client.NewRequest(ctx, http.MethodGet, urlStr, nil)
		if err != nil {
			return "", err
		}

		status := ipAccessListStatus{}
		_, err = client.Do(ctx, req, &status)
		if err != nil {
			return "", err
		}

		return status.Status, nil
	}
}
