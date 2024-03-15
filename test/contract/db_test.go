package contract_test

import (
	"context"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithDatabase(t *testing.T) {
	resources := &contract.TestResources{
		Name:           "db-test",
		ProjectID:      "65f41f0ff67c3b41f75b4c94",
		ServerlessName: "search-serverless-7e21ab",
		ClusterURL:     "mongodb+srv://search-serverless-7e21a.7ruqalf.mongodb-qa.net",
		UserDB:         "admin",
		Username:       "search-user-0b877c",
		Password:       "blasjhdlsajdldefj",
	}
	configDB := contract.WithDatabase(resources.Name)
	newResources, err := configDB(context.Background(), resources)
	require.NoError(t, err)
	assert.NotEmpty(t, newResources.DatabaseName)
}
