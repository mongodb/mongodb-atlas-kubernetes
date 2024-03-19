package search

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
)

const (
	TestName = "search"
)

var (
	// WipeResources removes resources at test cleanup time, no reuse will be possible afterwards
	WipeResources = contract.BoolEnv("WIPE_RESOURCES", false)
)

var resources *contract.TestResources

func TestMain(m *testing.M) {
	os.Exit(contract.RunTests(m, &resources, func(ctx context.Context) (*contract.TestResources, error) {
		return contract.DeployTestResources(ctx,
			TestName,
			WipeResources,
			contract.DefaultProject(TestName),
			contract.WithIPAccessList(contract.DefaultIPAccessList()),
			contract.WithCluster(contract.DefaultM0(TestName)),
			contract.WithUser(contract.DefaultUser(TestName)),
			contract.WithDatabase(TestName),
		)
	}))
}

func TestCreateSearchIndex(t *testing.T) {
	ctx := context.Background()
	apiClient, err := contract.NewAPIClient()
	require.NoError(t, err)
	assert.NotNil(t, apiClient)
	dynamic := true
	csi, _, err := apiClient.AtlasSearchApi.CreateAtlasSearchIndex(
		ctx,
		resources.ProjectID,
		resources.ClusterName,
		&admin.ClusterSearchIndex{
			CollectionName: resources.CollectionName,
			Database:       resources.DatabaseName,
			Name:           TestName,
			Mappings:       &admin.ApiAtlasFTSMappings{Dynamic: &dynamic},
		}).Execute()
	require.NoError(t, err)
	assert.NotNil(t, csi)

	_, _, err = apiClient.AtlasSearchApi.DeleteAtlasSearchIndex(
		ctx, resources.ProjectID, resources.ClusterName, *csi.IndexID).Execute()
	require.NoError(t, err)
}
