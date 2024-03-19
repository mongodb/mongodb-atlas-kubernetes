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
	ServerlessDeploymentTimeout = 10 * time.Minute
)

func DefaultServerless(prefix string) *admin.ServerlessInstanceDescriptionCreate {
	return &admin.ServerlessInstanceDescriptionCreate{
		Name: NewRandomName(fmt.Sprintf("%s-serverless", prefix)),
		ProviderSettings: admin.ServerlessProviderSettings{
			BackingProviderName: DefaultProviderName(),
			RegionName:          DefaultRegion(),
		},
	}
}

func WithServerless(serverless *admin.ServerlessInstanceDescriptionCreate) OptResourceFunc {
	return func(ctx context.Context, resources *TestResources) (*TestResources, error) {
		deployment, err := createServerless(ctx, resources.ProjectID, serverless)
		if err != nil {
			return nil, fmt.Errorf("failed to create serverless deployment %s: %w", serverless.Name, err)
		}
		resources.ServerlessName = *deployment.Name
		if err := waitServerless(ctx, resources.ProjectID, resources.ServerlessName, "IDLE", ServerlessDeploymentTimeout); err != nil {
			return nil, fmt.Errorf("failed to get serverless deployment %s up and running: %w", serverless.Name, err)
		}
		readyDeployment, err := getServerless(ctx, resources.ProjectID, resources.ServerlessName)
		if err != nil {
			return nil, fmt.Errorf("failed to query ready serverless deployment %s: %w", serverless.Name, err)
		}
		if readyDeployment.ConnectionStrings == nil || readyDeployment.ConnectionStrings.StandardSrv == nil {
			return nil, fmt.Errorf("missing connection string for serverless %s: %w", serverless.Name, err)
		}
		resources.ClusterURL = *readyDeployment.ConnectionStrings.StandardSrv

		resources.pushCleanup(func() error {
			if err := removeServerless(ctx, resources.ProjectID, resources.ServerlessName); err != nil {
				return err
			}
			if err = waitServerlessRemoval(ctx, resources.ProjectID, resources.ServerlessName, ServerlessDeploymentTimeout); err != nil {
				return fmt.Errorf("failed to get serverless deployment %s removed: %w", resources.ServerlessName, err)
			}
			return nil
		})
		return resources, nil
	}
}

func checkServerless(ctx context.Context, projectID, serverlessName string) error {
	_, err := getServerless(ctx, projectID, serverlessName)
	return err
}

func getServerless(ctx context.Context, projectID, serverlessName string) (*admin.ServerlessInstanceDescription, error) {
	apiClient, err := NewAPIClient()
	if err != nil {
		return nil, err
	}
	serverless, _, err := apiClient.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID, serverlessName).Execute()
	return serverless, err
}

func createServerless(ctx context.Context, projectID string, serverless *admin.ServerlessInstanceDescriptionCreate) (*admin.ServerlessInstanceDescription, error) {
	log.Printf("Creating serverless deployment %s...", serverless.Name)
	apiClient, err := NewAPIClient()
	if err != nil {
		return nil, err
	}
	deployment, _, err :=
		apiClient.ServerlessInstancesApi.CreateServerlessInstance(ctx, projectID, serverless).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create serverless instance %s: %w", serverless.Name, err)
	}
	log.Printf("Created serverless deployment %s ID=%s", serverless.Name, *deployment.Id)
	return deployment, nil
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
