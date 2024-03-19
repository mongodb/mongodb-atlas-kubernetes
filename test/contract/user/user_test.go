package user

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
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
	os.Exit(contract.RunTests(m, &resources, func(ctx context.Context) (*contract.TestResources, error) {
		return contract.DeployTestResources(ctx,
			TestName,
			WipeResources,
			contract.DefaultProject(TestName),
			contract.WithIPAccessList(contract.DefaultIPAccessList()),
			contract.WithServerless(contract.DefaultServerless(TestName)),
		)
	}))
}

func TestUser(t *testing.T) {
	ctx := context.Background()
	user := contract.DefaultUser(TestName)
	require.NoError(t, contract.EnsureUser(ctx, resources.ProjectID, resources.ClusterURL, user))

	credentials := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      user.Username,
		Password:      *user.Password,
	}
	err := writeTestData(ctx, credentials, "db", "collection")
	require.NoError(t, err)
}

func writeTestData(ctx context.Context, credentials options.Credential, dbName, collectionName string) error {
	uri := resources.ClusterURL
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(resources.ClusterURL).SetAuth(credentials))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB at %s: %w", uri, err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB at %s: %v", uri, err)
		}
	}()
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	_, err = collection.InsertOne(ctx, resources)
	if err != nil {
		return fmt.Errorf("failed to insert test data into MongoDB %s at %s.%s: %w",
			uri, dbName, collectionName, err)
	}
	return nil
}
