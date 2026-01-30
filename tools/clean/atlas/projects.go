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
	"fmt"
	"sync"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
)

func (c *Cleaner) listProjectsByOrg(ctx context.Context, orgID string) []admin.Group {
	projectsList, _, err := c.client.OrganizationsApi.
		GetOrgGroups(ctx, orgID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("Failed to list projects of organization %s: %s", orgID, err))

		return nil
	}

	if projectsList.GetTotalCount() == 0 {
		fmt.Println(text.FgYellow.Sprintf("No projects found in organization %s", orgID))

		return nil
	}

	return *projectsList.Results
}

func (c *Cleaner) deleteProject(ctx context.Context, p *admin.Group) {
	_, err := c.client.ProjectsApi.DeleteGroup(ctx, p.GetId()).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to request deletion of project %s(%s): %s", p.GetName(), p.GetId(), err))

		return
	}

	fmt.Println(text.FgGreen.Sprintf("\tRequested deletion of project %s(%s)", p.GetName(), p.GetId()))
}

func (c *Cleaner) GetProjectDependencies(ctx context.Context, projectID string, isGov bool) *ProjectDependencies {
	deps := &ProjectDependencies{}

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()

		providers := SupportedProviders
		if isGov {
			providers = GovSupportedProviders
		}
		deps.NetworkPeering = c.listNetworkPeering(ctx, projectID, providers)
	}()

	go func() {
		defer wg.Done()

		deps.AWSPrivateEndpoints = c.listPrivateEndpoints(ctx, projectID, CloudProviderAWS)
	}()

	go func() {
		defer wg.Done()

		deps.Clusters = c.listClusters(ctx, projectID)
	}()

	go func() {
		defer wg.Done()

		deps.EncryptionAtRest = c.getEncryptionAtRest(ctx, projectID)
	}()

	if !isGov {
		wg.Add(7)

		go func() {
			defer wg.Done()

			deps.Streams = c.listStreams(ctx, projectID)
		}()

		go func() {
			defer wg.Done()

			deps.GCPPrivateEndpoints = c.listPrivateEndpoints(ctx, projectID, CloudProviderGCP)
		}()

		go func() {
			defer wg.Done()

			deps.AzurePrivateEndpoints = c.listPrivateEndpoints(ctx, projectID, CloudProviderAZURE)
		}()

		go func() {
			defer wg.Done()

			deps.ServerlessClusters = c.listServerlessClusters(ctx, projectID)
		}()

		go func() {
			defer wg.Done()

			deps.FlexClusters = c.listFlexClusters(ctx, projectID)
		}()

		go func() {
			defer wg.Done()

			deps.FederatedDBPrivateEndpoints = c.listFederatedDBPrivateEndpoints(ctx, projectID)
		}()

		go func() {
			defer wg.Done()

			deps.FederatedDatabases = c.listFederatedDatabases(ctx, projectID)
		}()
	}

	wg.Wait()

	return deps
}

func (c *Cleaner) DeleteProjectDependencies(ctx context.Context, projectID string, deps *ProjectDependencies) {
	c.deleteClusters(ctx, projectID, deps.Clusters)

	c.deleteServerlessClusters(ctx, projectID, deps.ServerlessClusters)

	c.deleteFlexClusters(ctx, projectID, deps.FlexClusters)

	c.deleteFederatedDBPrivateEndpoints(ctx, projectID, deps.FederatedDBPrivateEndpoints)

	c.deleteFederatedDatabases(ctx, projectID, deps.FederatedDatabases)

	c.deleteEncryptionAtRest(ctx, projectID, deps.EncryptionAtRest)

	c.deleteNetworkPeering(ctx, projectID, deps.NetworkPeering)

	c.deletePrivateEndpoints(ctx, projectID, CloudProviderAWS, deps.AWSPrivateEndpoints)

	c.deletePrivateEndpoints(ctx, projectID, CloudProviderGCP, deps.GCPPrivateEndpoints)

	c.deletePrivateEndpoints(ctx, projectID, CloudProviderAZURE, deps.AzurePrivateEndpoints)

	c.deleteStreams(ctx, projectID, deps.Streams)
}
