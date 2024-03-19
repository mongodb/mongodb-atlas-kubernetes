package contract

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const (
	InitialRetryWait = time.Second

	MaxRetries = 5
)

func DefaultUser(username string) *admin.CloudDatabaseUser {
	return &admin.CloudDatabaseUser{
		Roles: &[]admin.DatabaseUserRole{
			{
				RoleName:     "readWriteAnyDatabase",
				DatabaseName: "admin",
			},
		},
		DatabaseName: "admin",
		Username:     username,
		Password:     pointer.MakePtr(NewRandomName("some-password")),
	}
}

func WithUser(user *admin.CloudDatabaseUser) OptResourceFunc {
	return func(ctx context.Context, resources *TestResources) (*TestResources, error) {
		newUser, err := createUser(ctx, resources.ProjectID, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create username %s: %w", user.Username, err)
		}
		resources.UserDB = newUser.DatabaseName
		resources.Username = newUser.Username
		if user.Password == nil {
			return nil, fmt.Errorf("no password for username %s: %w", newUser.Username, err)
		}
		if err := waitUserPassword(ctx, resources.ClusterURL, user, InitialRetryWait, MaxRetries); err != nil {
			return nil, err
		}
		resources.Password = *user.Password
		resources.pushCleanup(func() error {
			return removeUser(ctx, resources.ProjectID, user.DatabaseName, user.Username)
		})
		return resources, nil
	}
}

func checkUser(ctx context.Context, projectID, userDB, username string) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err =
		apiClient.DatabaseUsersApi.GetDatabaseUser(ctx, projectID, userDB, username).Execute()
	return err
}

func createUser(ctx context.Context, projectID string, user *admin.CloudDatabaseUser) (*admin.CloudDatabaseUser, error) {
	log.Printf("Creating user %s...", user.Username)
	apiClient, err := NewAPIClient()
	if err != nil {
		return nil, err
	}
	newUser, _, err := apiClient.DatabaseUsersApi.CreateDatabaseUser(ctx, projectID, user).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create user %s: %w", user.Username, err)
	}
	log.Printf("Created user %s", newUser.Username)
	return newUser, nil
}

func EnsureUser(ctx context.Context, projectID, uri string, user *admin.CloudDatabaseUser) error {
	_, err := createUser(ctx, projectID, user)
	detailedError := &admin.GenericOpenAPIError{}
	if errors.As(err, &detailedError) {
		apiErr := detailedError.Model()
		if (&apiErr).GetErrorCode() == "USER_ALREADY_EXISTS" {
			log.Printf("User %s already created, updating password", user.Username)
			if err := updateUser(ctx, projectID, user); err != nil {
				return err
			}
		}
	} else if err != nil {
		return fmt.Errorf("failed to create user %s with admin API: %w", user.Username, err)
	}
	return waitUserPassword(ctx, uri, user, InitialRetryWait, MaxRetries)
}

func updateUser(ctx context.Context, projectID string, user *admin.CloudDatabaseUser) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return fmt.Errorf("failed to get an api client: %w", err)
	}
	_, _, err = apiClient.DatabaseUsersApi.UpdateDatabaseUser(
		ctx, projectID, user.DatabaseName, user.Username, user).Execute()
	if err != nil {
		return fmt.Errorf("failed to update existing user %s: %w", user.Username, err)
	}
	return nil
}

// waitUserPassword waits for the password to settle.
// When the admin API creates or updates a user password this is not applied immediately,
// if you try to access the database right away with automation code you might get auth errors.
// See https://jira.mongodb.org/browse/CLOUDP-238496 for more details
func waitUserPassword(ctx context.Context, uri string, user *admin.CloudDatabaseUser, initialRetryWait time.Duration, maxRetries int) error {
	credentials := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      user.Username,
		Password:      *user.Password,
	}
	retryWait := initialRetryWait
	var client *mongo.Client
	var err error
	defer func() {
		if client != nil {
			client.Disconnect(ctx)
		}
	}()
	for retries := maxRetries; retries > 0; retries-- {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credentials))
		if err != nil {
			return fmt.Errorf("failed to re-connect to MongoDB at %s: %w", uri, err)
		}
		if err := client.Ping(ctx, nil); err != nil {
			time.Sleep(retryWait)
			retryWait *= 2
			continue
		}
		return nil
	}
	return fmt.Errorf("timed out waiting for user %s password to be applied", user.Username)
}

func removeUser(ctx context.Context, projectID string, userDB, username string) error {
	log.Printf("Deleting user %s...", username)
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = apiClient.DatabaseUsersApi.DeleteDatabaseUser(ctx, projectID, userDB, username).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove user %s: %w", username, err)
	}
	return nil
}
