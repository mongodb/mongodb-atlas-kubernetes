// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlas

import (
	"context"

	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

// WaitForAtlasDeploymentStateToNotBeReached periodically checks the given atlas deployment for a given condition. The function
// returns true after the given context timeout is exceeded.
func WaitForAtlasDeploymentStateToNotBeReached(ctx context.Context, atlasClient *admin.APIClient, projectName, deploymentName string, fns ...func(description *admin.ClusterDescription20240805) bool) func() bool {
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
