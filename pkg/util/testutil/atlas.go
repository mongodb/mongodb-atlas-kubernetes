package testutil

import (
	"context"

	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
)

// WaitForAtlasStateToNotBeReached periodically checks the given atlas cluster for a given condition. The function
// returns true after the given context timeout is exceeded.
func WaitForAtlasStateToNotBeReached(ctx context.Context, atlasClient *mongodbatlas.Client, projectName, clusterName string, fns ...func(*mongodbatlas.Cluster) bool) func() bool {
	return func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			atlasCluster, _, err := atlasClient.Clusters.Get(context.Background(), projectName, clusterName)
			if err != nil {
				return false
			}

			allTrue := true
			for _, fn := range fns {
				if !fn(atlasCluster) {
					allTrue = false
				}
			}

			Expect(allTrue).To(BeFalse())

			return allTrue
		}
	}
}
