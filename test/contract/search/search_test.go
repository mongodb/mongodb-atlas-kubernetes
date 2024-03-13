package search

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
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
	if !control.Enabled("AKO_CONTRACT_TEST") {
		log.Print("Skipping e2e tests, AKO_CONTRACT_TEST is not set")
		return
	}
	ctx := context.Background()
	beforeAll(ctx)
	code := m.Run()
	afterAll(ctx)
	os.Exit(code)
}

func beforeAll(ctx context.Context) {
	log.Printf("WipeResources set to %v", WipeResources)
	resources = contract.MustDeployTestResources(ctx,
		TestName,
		WipeResources,
		contract.DefaultProject(TestName),
		contract.WithServerless(contract.DefaultServerless(TestName)),
		contract.WithUser(contract.DefaultUser(TestName)),
		contract.WithDatabase(TestName),
	)
	log.Printf("Resources ready:\n%v", resources)
}

func afterAll(ctx context.Context) {
	resources.MustRecycle(ctx, WipeResources)
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
