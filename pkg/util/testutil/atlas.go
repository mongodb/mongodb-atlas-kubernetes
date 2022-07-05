package testutil

import (
	"context"

	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
)

// WaitForAtlasDeploymentStateToNotBeReached periodically checks the given atlas deployment for a given condition. The function
// returns true after the given context timeout is exceeded.
func WaitForAtlasDeploymentStateToNotBeReached(ctx context.Context, atlasClient *mongodbatlas.Client, projectName, deploymentName string, fns ...func(*mongodbatlas.Cluster) bool) func() bool {
	return func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			atlasDeployment, _, err := atlasClient.Clusters.Get(context.Background(), projectName, deploymentName)
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
func WaitForAtlasDatabaseUserStateToNotBeReached(ctx context.Context, atlasClient *mongodbatlas.Client, authDb, groupId, userName string, fns ...func(user *mongodbatlas.DatabaseUser) bool) func() bool {
	return func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			atlasDatabaseUser, _, err := atlasClient.DatabaseUsers.Get(context.Background(), authDb, groupId, userName)
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

// WaitForAtlasProjectStateToNotBeReached periodically checks the given atlas project for a given condition.
// The function returns true after the given context timeout is exceeded.
func WaitForAtlasProjectStateToNotBeReached(ctx context.Context, atlasClient *mongodbatlas.Client, projectName string, fns ...func(project *mongodbatlas.Project) bool) func() bool {
	return func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			project, _, err := atlasClient.Projects.GetOneProjectByName(context.Background(), projectName)
			if err != nil {
				return false
			}

			allTrue := true
			for _, fn := range fns {
				if !fn(project) {
					allTrue = false
				}
			}

			Expect(allTrue).To(BeFalse())

			return allTrue
		}
	}
}
