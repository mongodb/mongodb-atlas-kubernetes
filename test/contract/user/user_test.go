package user

import (
	"context"
	"log"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
	"github.com/stretchr/testify/require"
)

const (
	TestName = "user"
)

var (
	// WipeResources removes resources at test cleanup time, no reuse will be possible afterwards
	WipeResources = contract.BoolEnv("WIPE_RESOURCES", false)
)

var resources *contract.TestResources

func TestMain(m *testing.M) {
	contract.TestMain(m,
		func(ctx context.Context) {
			log.Printf("WipeResources set to %v", WipeResources)
			resources = contract.MustDeployTestResources(ctx,
				TestName,
				WipeResources,
				contract.DefaultProject(TestName),
				contract.WithIPAccessList(contract.DefaultIPAccessList()),
				contract.WithServerless(contract.DefaultServerless(TestName)),
			)
		},
		func(ctx context.Context) {
			resources.MustRecycle(ctx, WipeResources)
		},
	)
}

func TestUser(t *testing.T) {
	apiClient, err := contract.NewAPIClient()
	require.NoError(t, err)
	user := contract.DefaultUser(resources.Name)
	log.Printf("user to create: %s", contract.Jsonize(user))
	ctx := context.Background()
	newUser, _, err := apiClient.DatabaseUsersApi.CreateDatabaseUser(ctx, resources.ProjectID, user).Execute()
	require.NoError(t, err)
	log.Printf("Created user %s", contract.Jsonize(newUser))
}
