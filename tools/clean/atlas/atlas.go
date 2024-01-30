package atlas

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"tools/clean/provider"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
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
	NetworkPeering        []admin.BaseNetworkPeeringConnectionSettings
	AWSPrivateEndpoints   []admin.EndpointService
	GCPPrivateEndpoints   []admin.EndpointService
	AzurePrivateEndpoints []admin.EndpointService
	Clusters              []admin.AdvancedClusterDescription
	ServerlessClusters    []admin.ServerlessInstanceDescription
	FederatedDatabases    []admin.DataLakeTenant
	EncryptionAtRest      *admin.EncryptionAtRest
}

func (pd *ProjectDependencies) Length() int {
	return len(pd.NetworkPeering) + len(pd.AWSPrivateEndpoints) + len(pd.GCPPrivateEndpoints) + len(pd.AzurePrivateEndpoints) +
		len(pd.Clusters) + len(pd.ServerlessClusters) + len(pd.FederatedDatabases)
}

func (c *Cleaner) Clean(ctx context.Context, lifetime int) error {
	projects := c.listProjectsByOrg(ctx, c.orgID)

	if len(projects) > 0 {
		fmt.Println(text.FgGreen.Sprintf("Deletion Progress of %d projects ...\n", len(projects)))
	}

	for _, proj := range projects {
		p := proj

		fmt.Println(text.FgHiWhite.Sprintf("\tStarting deletion of project %s(%s) (created at %v)...", p.GetName(), p.GetId(), p.Created))

		if time.Since(p.Created) < time.Duration(lifetime)*time.Hour {
			fmt.Println(text.FgYellow.Sprintf("\tProject %s(%s) skipped once created less than %d hour ago", p.GetName(), p.GetId(), lifetime))

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

	return nil
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
