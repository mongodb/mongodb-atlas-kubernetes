package contract

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

func DefaultProject(prefix string) *admin.Group {
	return &admin.Group{
		Name:  newRandomName(prefix),
		OrgId: OrgID(),
	}
}

func checkProject(ctx context.Context, projectID string) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = apiClient.ProjectsApi.GetProject(ctx, projectID).Execute()
	return err
}

func createProject(ctx context.Context, project *admin.Group) (string, error) {
	log.Printf("Creating project %s...", project.Name)
	apiClient, err := NewAPIClient()
	if err != nil {
		return "", err
	}
	newProject, _, err := apiClient.ProjectsApi.CreateProject(ctx, project).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to create project %s: %w", project.Name, err)
	}
	log.Printf("Created project %s ID=%s", newProject.Name, *newProject.Id)
	return *newProject.Id, nil
}

func removeProject(ctx context.Context, projectID string) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = client.ProjectsApi.DeleteProject(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove project %s: %w", projectID, err)
	}
	log.Printf("Removed default project %s...", projectID)
	return nil
}
