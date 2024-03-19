package contract

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const (
	MaxWait = 20 * time.Second
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
		if resources.UserDB == "" || resources.Username == "" || resources.Password == "" {
			newUser, err := createUser(ctx, resources.ProjectID, user)
			if err != nil {
				return nil, fmt.Errorf("failed to create username %s: %w", user.Username, err)
			}
			resources.UserDB = newUser.DatabaseName
			resources.Username = newUser.Username
			if user.Password == nil {
				return nil, fmt.Errorf("no password for username %s: %w", newUser.Username, err)
			}
			if err := waitChanges(ctx, resources.ProjectID, resources.ClusterName, MaxWait); err != nil {
				return nil, err
			}
			resources.Password = *user.Password
		} else {
			if err := checkUser(ctx, resources.ProjectID, resources.UserDB, resources.Username); err != nil {
				return nil, err
			}
		}
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
	if err != nil {
		return fmt.Errorf("failed to check user %s: %w", username, err)
	}
	return nil
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

func EnsureUser(ctx context.Context, projectID, clusterName string, user *admin.CloudDatabaseUser) error {
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
	return waitChanges(ctx, projectID, clusterName, MaxWait)
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
