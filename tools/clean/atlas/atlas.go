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
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
	"tools/clean/provider"
)

const (
	CloudProviderAWS   = "AWS"
	CloudProviderGCP   = "GCP"
	CloudProviderAZURE = "AZURE"
)

var (
	SupportedProviders = []string{CloudProviderAWS, CloudProviderGCP, CloudProviderAZURE}

	GovSupportedProviders = []string{CloudProviderAWS}
)

type Cleaner struct {
	client *admin.APIClient
	aws    *provider.AWS
	gcp    *provider.GCP
	azure  *provider.Azure
	orgID  string
}

type ProjectDependencies struct {
	NetworkPeering              []admin.BaseNetworkPeeringConnectionSettings
	AWSPrivateEndpoints         []admin.EndpointService
	GCPPrivateEndpoints         []admin.EndpointService
	AzurePrivateEndpoints       []admin.EndpointService
	Clusters                    []admin.ClusterDescription20240805
	ServerlessClusters          []admin.ServerlessInstanceDescription
	FederatedDatabases          []admin.DataLakeTenant
	FederatedDBPrivateEndpoints []admin.PrivateNetworkEndpointIdEntry
	EncryptionAtRest            *admin.EncryptionAtRest
	FlexClusters                []admin.FlexClusterDescription20241113
	Streams                     []admin.StreamsTenant
}

func (pd *ProjectDependencies) Length() int {
	return len(pd.NetworkPeering) + len(pd.AWSPrivateEndpoints) + len(pd.GCPPrivateEndpoints) + len(pd.AzurePrivateEndpoints) +
		len(pd.Clusters) + len(pd.ServerlessClusters) + len(pd.FederatedDBPrivateEndpoints) + len(pd.FederatedDatabases) + len(pd.FlexClusters) + len(pd.Streams)
}

func (c *Cleaner) Clean(ctx context.Context, lifetimeHours int) error {
	projects := c.listProjectsByOrg(ctx, c.orgID)

	if len(projects) > 0 {
		fmt.Println(text.FgGreen.Sprintf("Deletion Progress of %d projects ...\n", len(projects)))
	}

	for _, proj := range projects {
		p := proj

		fmt.Println(text.FgHiWhite.Sprintf("\tStarting deletion of project %s(%s) (created at %v)...", p.GetName(), p.GetId(), p.Created))

		if time.Since(p.Created) < time.Duration(lifetimeHours)*time.Hour {
			fmt.Println(text.FgYellow.Sprintf("\tProject %s(%s) skipped once created less than %d hour ago", p.GetName(), p.GetId(), lifetimeHours))
			continue
		}

		deps := c.GetProjectDependencies(ctx, p.GetId(), isGov(c.client.GetConfig().Host))

		if deps.Length() > 0 {
			fmt.Println(text.FgWhite.Sprintf("\t\tDeleting dependencies of project %s(%s) ...", p.GetName(), p.GetId()))
			c.DeleteProjectDependencies(ctx, p.GetId(), deps)
			fmt.Println(text.FgYellow.Sprintf("\t\tProject %s(%s) should be ready for deletion on next run", p.GetName(), p.GetId()))

			continue
		}

		c.deleteProject(ctx, &p)
	}

	fmt.Println()

	teams := c.listTeamsByOrg(ctx, c.orgID)

	if len(teams) > 0 {
		fmt.Println(text.FgGreen.Sprintf("Deletion Progress of %d teams ...\n", len(teams)))
	}

	for _, team := range teams {
		t := team
		c.deleteTeam(ctx, c.orgID, &t)
	}

	c.cleanOrphanResources(ctx, lifetimeHours)
	return nil
}

func (c *Cleaner) cleanOrphanResources(ctx context.Context, lifetimeHours int) {
	gcpVpcPrefixes := []string{"gcp-pe", "ao-vpc", "migrate-private-endpoint"}
	gcpRegions := []string{
		"us-west3",
		"us-east5",
		"southamerica-east1",
		"europe-west3",
		"europe-north1",
		"europe-west6",
		"europe-west1",
		"asia-east2",
		"asia-northeast2",
		"australia-southeast2",
	}

	var done, skipped []string
	var errs []error

	addResults := func(f func() ([]string, []string, []error)) {
		d, s, e := f()
		done = append(done, d...)
		skipped = append(skipped, s...)
		errs = append(errs, e...)
	}

	for _, region := range gcpRegions {
		addResults(func() ([]string, []string, []error) {
			return c.gcp.DeleteOrphanPrivateEndpoints(ctx, region, lifetimeHours)
		})
	}

	addResults(func() ([]string, []string, []error) {
		return c.gcp.DeleteOrphanVPCs(ctx, gcpVpcPrefixes, gcpRegions, lifetimeHours)
	})

	for _, doneMsg := range done {
		fmt.Println(text.FgGreen.Sprintf("%s", doneMsg))
	}
	for _, skippedMsg := range skipped {
		fmt.Println(text.FgYellow.Sprintf("\t%s", skippedMsg))
	}
	for _, err := range errs {
		fmt.Println(text.FgRed.Sprintf("%v", err.Error()))
	}
}

func NewCleaner(aws *provider.AWS, gcp *provider.GCP, azure *provider.Azure) (*Cleaner, error) {
	apiURL, defined := os.LookupEnv("MCLI_OPS_MANAGER_URL")
	if !defined {
		return nil, fmt.Errorf("API URL must be set")
	}

	apiKey, defined := os.LookupEnv("MCLI_PUBLIC_API_KEY")
	if !defined {
		return nil, fmt.Errorf("API public key must be set")
	}

	apiSecret, defined := os.LookupEnv("MCLI_PRIVATE_API_KEY")
	if !defined {
		return nil, fmt.Errorf("API private key must be set")
	}

	orgID, defined := os.LookupEnv("MCLI_ORG_ID")
	if !defined {
		return nil, fmt.Errorf("organization ID must be set")
	}

	adminClient, err := admin.NewClient(admin.UseBaseURL(apiURL), admin.UseDigestAuth(apiKey, apiSecret))
	if err != nil {
		return nil, err
	}

	return &Cleaner{
		client: adminClient,
		aws:    aws,
		gcp:    gcp,
		azure:  azure,
		orgID:  orgID,
	}, nil
}

func isGov(url string) bool {
	return strings.HasSuffix(url, "mongodbgov.com")
}
