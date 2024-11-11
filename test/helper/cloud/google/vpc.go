package google

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/compute/apiv1/computepb"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func CreateVPC(ctx context.Context, vpcName string) error {
	gce, err := newConnection(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return fmt.Errorf("failed to get Google Cloud connection: %w", err)
	}

	op, err := gce.networkClient.Insert(ctx, &computepb.InsertNetworkRequest{
		Project: gce.projectID,
		NetworkResource: &computepb.Network{
			Name:                  pointer.MakePtr(vpcName),
			Description:           pointer.MakePtr("Atlas Kubernetes Operator E2E Tests VPC"),
			AutoCreateSubnetworks: pointer.MakePtr(false),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to request creation of Google VPC: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Google VPC: %w", err)
	}

	return nil
}

func DeleteVPC(ctx context.Context, vpcName string) error {
	gce, err := newConnection(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		return fmt.Errorf("failed to get Google Cloud connection: %w", err)
	}
	op, err := gce.networkClient.Delete(ctx, &computepb.DeleteNetworkRequest{
		Project: gce.projectID,
		Network: vpcName,
	})
	if err != nil {
		return fmt.Errorf("failed to request deletion of Google VPC: %w", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete Google VPC: %w", err)
	}

	return nil
}
