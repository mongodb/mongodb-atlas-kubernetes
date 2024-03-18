package search

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			contract.WithServerless(contract.DefaultServerless(TestName)),
			contract.WithUser(contract.DefaultUser(TestName)),
			contract.WithDatabase(TestName),
		)
	}))
}

func TestCreateSearchIndex(t *testing.T) {
	//ctx := context.Background()
	apiClient, err := contract.NewAPIClient()
	require.NoError(t, err)
	assert.NotNil(t, apiClient)
	// apiClient.AtlasSearchApi.CreateAtlasSearchIndex(
	// 	ctx,
	// 	resources.ProjectID,
	// 	resources.ServerlessName,
	// 	&admin.ClusterSearchIndex{
	// 		CollectionName: resources.CollectionName,
	// 		Database:       resources.DatabaseName,
	// 		IndexID:        new(string),
	// 		Name:           TestName,
	// 		Status:         new(string),
	// 		Type:           new(string),
	// 		Analyzer:       new(string),
	// 		Analyzers:      &[]admin.ApiAtlasFTSAnalyzers{},
	// 		Mappings:       &admin.ApiAtlasFTSMappings{},
	// 		SearchAnalyzer: new(string),
	// 		Synonyms:       &[]admin.SearchSynonymMappingDefinition{},
	// 		Fields:         &[]map[string]interface{}{},
	// 	}).Execute()
}
