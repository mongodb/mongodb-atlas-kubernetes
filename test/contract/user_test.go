package contract_test

import (
	"context"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithUser(t *testing.T) {
	resources := &contract.TestResources{
		Name:           "db-test",
		ProjectID:      "65f2b65e46ee7f32c077558c",
		ServerlessName: "ServerlessInstance0",
		ClusterURL:     "mongodb+srv://serverlessinstance0.evr1a.mongodb-qa.net",
	}
	configUser := contract.WithUser(contract.DefaultUser(resources.Name))
	newResources, err := configUser(context.Background(), resources)
	require.NoError(t, err)
	assert.NotEmpty(t, newResources.Password)
}
