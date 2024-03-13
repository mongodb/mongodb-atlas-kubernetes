package contract

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

func DefaultUser(prefix string) *admin.CloudDatabaseUser {
	password := newRandomName(fmt.Sprintf("%s-passwd", prefix))
	return &admin.CloudDatabaseUser{
		Roles: &[]admin.DatabaseUserRole{
			{
				RoleName:     "readWriteAnyDatabase",
				DatabaseName: "admin",
			},
		},
		DatabaseName: "admin",
		Username:     newRandomName(fmt.Sprintf("%s-user", prefix)),
		Password:     &password,
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
		resources.Password = *user.Password
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
