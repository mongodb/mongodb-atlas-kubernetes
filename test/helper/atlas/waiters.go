package atlas

import (
	"context"

	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
)

// WaitForAtlasDeploymentStateToNotBeReached periodically checks the given atlas deployment for a given condition. The function
// returns true after the given context timeout is exceeded.
func WaitForAtlasDeploymentStateToNotBeReached(ctx context.Context, atlasClient *admin.APIClient, projectName, deploymentName string, fns ...func(description *admin.AdvancedClusterDescription) bool) func() bool {
	return func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			atlasDeployment, _, err := atlasClient.ClustersApi.GetCluster(ctx, projectName, deploymentName).Execute()
			if err != nil {
				return false
			}

			allTrue := true
			for _, fn := range fns {
				if !fn(atlasDeployment) {
					allTrue = false
				}
			}

			Expect(allTrue).To(BeFalse())

			return allTrue
		}
	}
}

// WaitForAtlasDatabaseUserStateToNotBeReached periodically checks the given atlas database user for a given condition.
// The function returns true after the given context timeout is exceeded.
func WaitForAtlasDatabaseUserStateToNotBeReached(ctx context.Context, atlasClient *admin.APIClient, authDb, groupId, userName string, fns ...func(user *admin.CloudDatabaseUser) bool) func() bool {
	return func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			atlasDatabaseUser, _, err := atlasClient.DatabaseUsersApi.
				GetDatabaseUser(ctx, groupId, authDb, userName).
				Execute()
			if err != nil {
				return false
			}

			allTrue := true
			for _, fn := range fns {
				if !fn(atlasDatabaseUser) {
					allTrue = false
				}
			}

			Expect(allTrue).To(BeFalse())

			return allTrue
		}
	}
}
