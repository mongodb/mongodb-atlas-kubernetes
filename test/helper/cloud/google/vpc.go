// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
