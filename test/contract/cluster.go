package contract

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const (
	ClusterDeploymentTimeout = 15 * time.Minute
)

func DefaultM0(prefix string) *admin.AdvancedClusterDescription {
	return &admin.AdvancedClusterDescription{
		AcceptDataRisksAndForceReplicaSetReconfig: &time.Time{},
		ClusterType: pointer.MakePtr("REPLICASET"),
		Name:        pointer.MakePtr(NewRandomName(prefix)),
		ReplicationSpecs: &[]admin.ReplicationSpec{
			{
				RegionConfigs: &[]admin.CloudRegionConfig{
					{
						ElectableSpecs: &admin.HardwareSpec{
							InstanceSize: pointer.MakePtr("M0"),
							NodeCount:    pointer.MakePtr(1),
						},
						Priority:            pointer.MakePtr(7),
						ProviderName:        pointer.MakePtr("TENANT"),
						RegionName:          pointer.MakePtr("US_EAST_1"),
						BackingProviderName: pointer.MakePtr("AWS"),
					},
				},
			},
		},
	}
}

func WithCluster(cluster *admin.AdvancedClusterDescription) OptResourceFunc {
	return func(ctx context.Context, resources *TestResources) (*TestResources, error) {
		if resources.ClusterName == "" || resources.ClusterURL == "" {
			_, err := createCluster(ctx, resources.ProjectID, cluster)
			if err != nil {
				return nil, fmt.Errorf("failed to create cluster %s: %w", *cluster.Name, err)
			}
			resources.ClusterName = *cluster.Name
			if err := waitCluster(ctx, resources.ProjectID, resources.ClusterName, "IDLE", ClusterDeploymentTimeout); err != nil {
				return nil, fmt.Errorf("failed to get cluster %s up and running: %w", *cluster.Name, err)
			}
			readyDeployment, err := getCluster(ctx, resources.ProjectID, resources.ClusterName)
			if err != nil {
				return nil, fmt.Errorf("failed to query ready cluster %s: %w", *cluster.Name, err)
			}
			if readyDeployment.ConnectionStrings == nil || readyDeployment.ConnectionStrings.StandardSrv == nil {
				return nil, fmt.Errorf("missing connection string for cluster %s: %w", *cluster.Name, err)
			}
			resources.ClusterURL = *readyDeployment.ConnectionStrings.StandardSrv
		} else {
			if err := checkCluster(ctx, resources.ProjectID, resources.ClusterName); err != nil {
				return nil, err
			}
		}

		resources.pushCleanup(func() error {
			if err := removeCluster(ctx, resources.ProjectID, resources.ClusterName); err != nil {
				return err
			}
			if err := waitClusterRemoval(ctx, resources.ProjectID, resources.ClusterName, ClusterDeploymentTimeout); err != nil {
				return fmt.Errorf("failed to get cluster %s removed: %w", resources.ClusterName, err)
			}
			return nil
		})
		return resources, nil
	}
}

func checkCluster(ctx context.Context, projectID, clusterName string) error {
	_, err := getCluster(ctx, projectID, clusterName)
	if err != nil {
		return fmt.Errorf("failed to check cluster %s: %w", clusterName, err)
	}
	return nil
}

func getCluster(ctx context.Context, projectID, clusterName string) (*admin.AdvancedClusterDescription, error) {
	apiClient, err := NewAPIClient()
	if err != nil {
		return nil, err
	}
	cluster, _, err := apiClient.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	return cluster, err
}

func createCluster(ctx context.Context, projectID string, cluster *admin.AdvancedClusterDescription) (*admin.AdvancedClusterDescription, error) {
	log.Printf("Creating cluster %s...", *cluster.Name)
	apiClient, err := NewAPIClient()
	if err != nil {
		return nil, err
	}
	deployment, _, err :=
		apiClient.ClustersApi.CreateCluster(ctx, projectID, cluster).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster %s: %w", *cluster.Name, err)
	}
	log.Printf("Created cluster %s ID=%s", *cluster.Name, *deployment.Id)
	return deployment, nil
}

func removeCluster(ctx context.Context, projectID, clusterName string) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = client.ServerlessInstancesApi.DeleteServerlessInstance(ctx, projectID, clusterName).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove cluster %s: %w", clusterName, err)
	}
	log.Printf("Removed cluster %s...", clusterName)
	return nil
}

func waitClusterRemoval(ctx context.Context, projectID, clusterName string, timeout time.Duration) error {
	err := waitCluster(ctx, projectID, clusterName, "", timeout)
	if strings.Contains(err.Error(), "CLUSTER_NOT_FOUND") {
		return nil
	}
	return err
}

func waitCluster(ctx context.Context, projectID, clusterName, goal string, timeout time.Duration) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	start := time.Now()
	for time.Since(start) < timeout {
		deployment, _, err :=
			client.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
		if err != nil {
			return fmt.Errorf("failed to check cluster %s: %w", clusterName, err)
		}
		if *deployment.StateName == goal {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("%v timeout", timeout)
}

// waitChanges waits for changes to the cluster (or serverless) to complete.
// This includes creating or changing a user password:
// When the admin API creates or updates a user password this is not applied immediately,
// if you try to access the database right away with automation code you might get auth errors.
//
// See https://jira.mongodb.org/browse/CLOUDP-238496 for more details
func waitChanges(ctx context.Context, projectID, clusterName string, timeout time.Duration) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	start := time.Now()
	for time.Since(start) < timeout {
		status, _, err :=
			client.ClustersApi.GetClusterStatus(ctx, projectID, clusterName).Execute()
		if err != nil {
			return fmt.Errorf("failed to check cluster %s: %w", clusterName, err)
		}
		if status.GetChangeStatus() == "APPLIED" {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("%v timeout", timeout)
}
