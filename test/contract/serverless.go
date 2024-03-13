package contract

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

const (
	ServerlessDeploymentTimeout = 5 * time.Minute
)

func DefaultServerless(prefix string) *admin.ServerlessInstanceDescriptionCreate {
	return &admin.ServerlessInstanceDescriptionCreate{
		Name: newRandomName(fmt.Sprintf("%s-serverless", prefix)),
		ProviderSettings: admin.ServerlessProviderSettings{
			BackingProviderName: DefaultProviderName(),
			RegionName:          DefaultRegion(),
		},
	}
}

func WithServerless(serverless *admin.ServerlessInstanceDescriptionCreate) OptResourceFunc {
	return func(ctx context.Context, resources *TestResources) (*TestResources, error) {
		deploymentName, err := createServerless(ctx, resources.ProjectID, serverless)
		if err != nil {
			return nil, fmt.Errorf("failed to create serverless deployment %s: %w", serverless.Name, err)
		}
		resources.ServerlessName = deploymentName
		if err := waitServerless(ctx, resources.ProjectID, resources.ServerlessName, "IDLE", ServerlessDeploymentTimeout); err != nil {
			return nil, fmt.Errorf("failed to get serverless deployment %s up and running: %w", serverless.Name, err)
		}
		return resources, nil
	}
}

func checkServerless(ctx context.Context, projectID, serverlessName string) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = apiClient.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID, serverlessName).Execute()
	return err
}

func createServerless(ctx context.Context, projectID string, serverless *admin.ServerlessInstanceDescriptionCreate) (string, error) {
	log.Printf("Creating serverless deployment %s...", serverless.Name)
	apiClient, err := NewAPIClient()
	if err != nil {
		return "", err
	}
	deployment, _, err :=
		apiClient.ServerlessInstancesApi.CreateServerlessInstance(ctx, projectID, serverless).Execute()
	if err != nil {
		return "", err
	}
	log.Printf("Created serverless deployment %s ID=%s", serverless.Name, *deployment.Id)
	return serverless.Name, nil
}

func removeServerless(ctx context.Context, projectID, serverlessName string) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = client.ServerlessInstancesApi.DeleteServerlessInstance(ctx, projectID, serverlessName).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove serverless deployment %s: %w", serverlessName, err)
	}
	log.Printf("Removed serverless deployment %s...", serverlessName)
	return nil
}

func waitServerlessRemoval(ctx context.Context, projectID, deploymentName string, timeout time.Duration) error {
	err := waitServerless(ctx, projectID, deploymentName, "", timeout)
	if strings.Contains(err.Error(), "SERVERLESS_INSTANCE_NOT_FOUND") {
		return nil
	}
	return err
}

func waitServerless(ctx context.Context, projectID, deploymentName, goal string, timeout time.Duration) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	start := time.Now()
	for time.Since(start) < timeout {
		deployment, _, err :=
			client.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID, deploymentName).Execute()
		if err != nil {
			return fmt.Errorf("failed to check deployment %s: %w", deploymentName, err)
		}
		if *deployment.StateName == goal {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("%v timeout", timeout)
}
