package atlas

import (
	"context"
	"fmt"
	"sync"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

func (c *Cleaner) listProjectsByOrg(ctx context.Context, orgID string) []admin.Group {
	projectsList, _, err := c.client.OrganizationsApi.
		ListOrganizationProjects(ctx, orgID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("Failed to list projects of organization %s: %s", orgID, err))

		return nil
	}

	if projectsList.GetTotalCount() == 0 {
		fmt.Println(text.FgYellow.Sprintf("No projects found in organization %s", orgID))

		return nil
	}

	return projectsList.Results
}

func (c *Cleaner) deleteProject(ctx context.Context, p *admin.Group) {
	_, _, err := c.client.ProjectsApi.DeleteProject(ctx, p.GetId()).Execute()
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
		wg.Add(4)

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

			deps.FederatedDatabases = c.listFederatedDatabases(ctx, projectID)
		}()
	}

	wg.Wait()

	return deps
}

func (c *Cleaner) DeleteProjectDependencies(ctx context.Context, projectID string, deps *ProjectDependencies) {
	c.deleteClusters(ctx, projectID, deps.Clusters)

	c.deleteServerlessClusters(ctx, projectID, deps.ServerlessClusters)

	c.deleteFederatedDatabases(ctx, projectID, deps.FederatedDatabases)

	c.deleteEncryptionAtRest(ctx, projectID, deps.EncryptionAtRest)

	c.deleteNetworkPeering(ctx, projectID, deps.NetworkPeering)

	c.deletePrivateEndpoints(ctx, projectID, CloudProviderAWS, deps.AWSPrivateEndpoints)

	c.deletePrivateEndpoints(ctx, projectID, CloudProviderGCP, deps.GCPPrivateEndpoints)

	c.deletePrivateEndpoints(ctx, projectID, CloudProviderAZURE, deps.AzurePrivateEndpoints)
}
