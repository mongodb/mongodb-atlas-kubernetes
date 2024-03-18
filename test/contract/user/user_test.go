package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
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
	user, err := createUserAsNeeded(ctx, "test-user")
	require.NoError(t, err)
	err = writeData(ctx, user, "db", "collection")
	require.NoError(t, err)
}

func createUserAsNeeded(ctx context.Context, username string) (*admin.CloudDatabaseUser, error) {
	apiClient, err := contract.NewAPIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to admin API: %w", err)
	}
	user := contract.DefaultUser(username)
	user.GroupId = resources.ProjectID
	log.Printf("user to create: %s", contract.Jsonize(user))
	_, _, err = apiClient.DatabaseUsersApi.CreateDatabaseUser(ctx, resources.ProjectID, user).Execute()
	detailedError := admin.GenericOpenAPIError{}
	if errors.As(err, &detailedError) {
		apiErr := detailedError.Model()
		if (&apiErr).GetErrorCode() == "USER_ALREADY_EXISTS" {
			log.Printf("User %s already created, updating password", username)
			_, _, err = apiClient.DatabaseUsersApi.UpdateDatabaseUser(
				ctx, resources.ProjectID, user.DatabaseName, user.Username, user).Execute()
			if err != nil {
				return nil, fmt.Errorf("failed to update existing user %s: %w", username, err)
			}
			return user, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create user %s to admin API: %w", username, err)
	}
	log.Printf("Created user %s", contract.Jsonize(user))
	return user, nil
}

func writeData(ctx context.Context, user *admin.CloudDatabaseUser, dbName, collectionName string) error {
	//
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    user.DatabaseName,
		Username:      user.Username,
		Password:      *user.Password,
	}
	uri := resources.ClusterURL
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(resources.ClusterURL).SetAuth(credential))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB at %s: %w", uri, err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB at %s: %w", uri, err)
	}
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	_, err = collection.InsertOne(ctx, resources)
	if err != nil {
		return fmt.Errorf("failed to insert test data into MongoDB %s at %s.%s: %w",
			uri, dbName, collectionName, err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB at %s: %v", uri, err)
		}
	}()
	return nil
}
