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

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/smithy-go"
)

type AWS struct{}

func (a *AWS) DeleteVpc(ctx context.Context, ID, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create AWS config: %w", err)
	}
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = normalizeAwsRegion(region)
	})

	vpc, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{VpcIds: []string{ID}})
	if err != nil {
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &apiErr) && apiErr.Code == "InvalidVpcID.NotFound" {
			return nil
		}
		return err
	}

	if len(vpc.Vpcs) == 0 {
		return nil
	}

	err = a.deletePeeringConnection(ctx, client, ID)
	if err != nil {
		return err
	}

	_, err = client.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(ID)})

	return err
}

func (a *AWS) DeleteEndpoint(ctx context.Context, ID, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create AWS config: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	endpoints, err := client.DescribeVpcEndpointServices(
		ctx,
		&ec2.DescribeVpcEndpointServicesInput{ServiceNames: []string{ID}},
	)
	if err != nil {
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &apiErr) && apiErr.Code == "InvalidVpcEndpointID.NotFound" {
			return nil
		}
		return err
	}

	if len(endpoints.ServiceNames) == 0 {
		return nil
	}

	endpointIDs := make([]string, 0, len(endpoints.ServiceNames))

	for _, serviceName := range endpoints.ServiceNames {
		endpointID := strings.TrimPrefix(serviceName, fmt.Sprintf("com.amazonaws.vpce.%s.", normalizeAwsRegion(region)))
		endpointIDs = append(endpointIDs, endpointID)
	}

	_, err = client.DeleteVpcEndpointServiceConfigurations(
		ctx, &ec2.DeleteVpcEndpointServiceConfigurationsInput{ServiceIds: endpointIDs})

	return err
}

func (a *AWS) DeleteKMS(ctx context.Context, ID, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create AWS config: %w", err)
	}

	kmsClient := kms.NewFromConfig(cfg)

	_, err = kmsClient.DescribeKey(ctx, &kms.DescribeKeyInput{
		GrantTokens: nil,
		KeyId:       aws.String(ID),
	})
	if err != nil {
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &apiErr) && apiErr.Code == "NotFoundException" {
			return nil
		}
		return err
	}

	_, err = kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
		KeyId:               aws.String(ID),
		PendingWindowInDays: aws.Int32(7), // this is the minimum possible and can be up to 24h longer than value set
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AWS) deletePeeringConnection(ctx context.Context, ec2Client *ec2.Client, vpcID string) error {
	input := ec2.DescribeVpcPeeringConnectionsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("accepter-vpc-info.vpc-id"),
				Values: []string{vpcID},
			},
		},
	}

	peers, err := ec2Client.DescribeVpcPeeringConnections(ctx, &input)
	if err != nil {
		return err
	}

	if len(peers.VpcPeeringConnections) == 0 {
		return nil
	}

	for _, peer := range peers.VpcPeeringConnections {
		_, err = ec2Client.DeleteVpcPeeringConnection(ctx, &ec2.DeleteVpcPeeringConnectionInput{
			VpcPeeringConnectionId: peer.VpcPeeringConnectionId,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func NewAWSCleaner() *AWS {
	return &AWS{}
}

func normalizeAwsRegion(region string) string {
	return strings.ToLower(strings.ReplaceAll(region, "_", "-"))
}
