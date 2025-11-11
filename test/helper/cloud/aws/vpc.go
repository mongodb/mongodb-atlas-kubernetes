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

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	taghelper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper"
)

func CreateVPC(ctx context.Context, name, cidr, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create an AWS config: %w", err)
	}
	ec2Client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = region
	})
	result, err := ec2Client.CreateVpc(
		ctx,
		&ec2.CreateVpcInput{
			AmazonProvidedIpv6CidrBlock: aws.Bool(false),
			CidrBlock:                   aws.String(cidr),
			TagSpecifications: []types.TagSpecification{{
				ResourceType: types.ResourceTypeVpc,
				Tags: []types.Tag{
					{Key: aws.String("Name"), Value: aws.String(name)},
					{Key: aws.String(taghelper.OwnerEmailTag), Value: aws.String(taghelper.AKOEmail)},
					{Key: aws.String(taghelper.CostCenterTag), Value: aws.String(taghelper.AKOCostCenter)},
					{Key: aws.String(taghelper.EnvironmentTag), Value: aws.String(taghelper.AKOEnvTest)},
				},
			}},
		})
	if err != nil {
		return "", fmt.Errorf("failed to create an AWS VPC: %w", err)
	}

	_, err = ec2Client.ModifyVpcAttribute(
		ctx,
		&ec2.ModifyVpcAttributeInput{
			EnableDnsHostnames: &types.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
			VpcId: result.Vpc.VpcId,
		})
	if err != nil {
		return "", fmt.Errorf("failed to configure AWS VPC: %w", err)
	}

	return *result.Vpc.VpcId, nil
}

func DeleteVPC(ctx context.Context, vpcID, region string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to create an AWS config: %w", err)
	}
	ec2Client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = region
	})
	_, err = ec2Client.DeleteVpc(
		ctx,
		&ec2.DeleteVpcInput{
			DryRun: aws.Bool(false),
			VpcId:  aws.String(vpcID),
		})
	return err
}
