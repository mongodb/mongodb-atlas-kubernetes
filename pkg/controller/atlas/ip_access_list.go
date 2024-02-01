package atlas

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

type IPAccessListStatus func(ctx context.Context, projectID, entryValue string) (string, error)

func CustomIPAccessListStatus(client *admin.APIClient) IPAccessListStatus {
	return func(ctx context.Context, projectID, entryValue string) (string, error) {
		status, _, err := client.ProjectIPAccessListApi.GetProjectIpAccessListStatus(ctx, projectID, entryValue).Execute()
		if err != nil {
			return "", err
		}
		return status.STATUS, nil
	}
}
