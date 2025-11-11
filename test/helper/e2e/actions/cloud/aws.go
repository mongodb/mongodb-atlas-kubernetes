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

package cloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	taghelper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper"
)

type AwsAction struct {
	network *awsNetwork

	cfg aws.Config
}

type awsNetwork struct {
	VPC     string
	Subnets []string
}

type principal struct {
	AWS []string `json:"AWS,omitempty"`
}

type kmsPolicy struct {
	Version   string      `json:"Version"`
	Statement []statement `json:"Statement"`
}

type statement struct {
	Sid       string    `json:"Sid"`
	Effect    string    `json:"Effect"`
	Principal principal `json:"Principal"`
	Action    []string  `json:"Action"`
	Resource  string    `json:"Resource"`
}

func (a *AwsAction) CreateKMS(ctx context.Context, alias, region, atlasAccountArn, assumedRoleArn string) (string, error) {
	//kmsClient := kms.New(a.session, aws.NewConfig().WithRegion(region))
	kmsClient := kms.NewFromConfig(a.cfg, func(o *kms.Options) {
		o.Region = region
	})

	adminARNs, err := getAdminARNs()
	if err != nil {
		return "", err
	}

	policyString, err := rolePolicyString(atlasAccountArn, assumedRoleArn, adminARNs)
	if err != nil {
		return "", err
	}

	key, err := kmsClient.CreateKey(ctx,
		&kms.CreateKeyInput{
			Description: aws.String("Key for E2E test"),
			KeySpec:     kmstypes.KeySpecSymmetricDefault,
			KeyUsage:    kmstypes.KeyUsageTypeEncryptDecrypt,
			MultiRegion: aws.Bool(false),
			Origin:      kmstypes.OriginTypeAwsKms,
			Policy:      aws.String(policyString),
			Tags: []kmstypes.Tag{
				{TagKey: aws.String(taghelper.OwnerTag), TagValue: aws.String(taghelper.AKOTeam)},
				{TagKey: aws.String(taghelper.OwnerEmailTag), TagValue: aws.String(taghelper.AKOEmail)},
				{TagKey: aws.String(taghelper.CostCenterTag), TagValue: aws.String(taghelper.AKOCostCenter)},
				{TagKey: aws.String(taghelper.EnvironmentTag), TagValue: aws.String(taghelper.AKOEnvTest)},
			},
		})

	if err != nil {
		return "", err
	}

	_, err = kmsClient.CreateAlias(ctx, &kms.CreateAliasInput{
		AliasName:   aws.String("alias/" + strings.ToLower(strings.ReplaceAll(alias, " ", "-"))),
		TargetKeyId: key.KeyMetadata.KeyId,
	})

	Expect(err).NotTo(HaveOccurred())

	DeferCleanup(func(ctx SpecContext) error {
		_, err = kmsClient.ScheduleKeyDeletion(ctx, &kms.ScheduleKeyDeletionInput{
			KeyId:               key.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int32(7), // this is the minimum possible and can be up to 24h longer than value set
		})
		return err
	})

	return *key.KeyMetadata.KeyId, nil
}

func getAdminARNs() ([]string, error) {
	adminArnString := os.Getenv("AWS_ACCOUNT_ARN_LIST")
	if adminArnString == "" {
		return nil, errors.New("AWS_ACCOUNT_ARN_LIST secret is empty")
	}

	adminARNs := strings.Split(adminArnString, ",")
	if len(adminARNs) == 0 {
		return nil, errors.New("AWS_ACCOUNT_ARN_LIST wasn't parsed properly, please separate accounts via a comma")
	}

	return adminARNs, nil
}

func rolePolicyString(atlasAccountARN, assumedRoleARN string, adminARNs []string) (string, error) {
	policy := defaultKMSPolicy(atlasAccountARN, assumedRoleARN, adminARNs)
	byteStr, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	return string(byteStr), nil
}

func defaultKMSPolicy(atlasAccountArn, assumedRoleArn string, adminARNs []string) kmsPolicy {
	return kmsPolicy{
		Version: "2012-10-17",
		Statement: []statement{
			{
				Sid:    "Enable IAM User Permissions",
				Effect: "Allow",
				Principal: principal{
					AWS: []string{atlasAccountArn},
				},
				Action:   []string{"kms:*"},
				Resource: "*",
			},
			{
				Sid:    "Allow access for Key Administrators",
				Effect: "Allow",
				Principal: principal{
					AWS: adminARNs,
				},
				Action: []string{
					"kms:Create*",
					"kms:Describe*",
					"kms:Enable*",
					"kms:List*",
					"kms:Put*",
					"kms:Update*",
					"kms:Revoke*",
					"kms:Disable*",
					"kms:Get*",
					"kms:Delete*",
					"kms:TagResource",
					"kms:UntagResource",
					"kms:ScheduleKeyDeletion",
					"kms:CancelKeyDeletion",
				},
				Resource: "*",
			},
			{
				Sid:    "Allow use of the key",
				Effect: "Allow",
				Principal: principal{
					AWS: []string{assumedRoleArn},
				},
				Action: []string{
					"kms:Encrypt",
					"kms:Decrypt",
					"kms:ReEncrypt*",
					"kms:GenerateDataKey*",
					"kms:DescribeKey",
				},
				Resource: "*",
			},
		},
	}
}

func (a *AwsAction) GetAccountID(ctx context.Context) (string, error) {
	stsClient := sts.NewFromConfig(a.cfg)
	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *identity.Account, nil
}

func (a *AwsAction) InitNetwork(ctx context.Context, vpcName, cidr, region string, subnets map[string]string, cleanup bool) (string, error) {
	vpc, err := a.findVPC(ctx, vpcName, region)
	if err != nil {
		return "", err
	}

	if vpc == "" {
		vpc, err = a.createVPC(ctx, vpcName, cidr, region)
		if err != nil {
			return "", err
		}
	}

	if cleanup {
		DeferCleanup(func(ctx SpecContext) error {
			return a.deleteVPC(ctx, vpc, region)
		})
	}

	subnetsMap, err := a.getSubnets(ctx, vpc, region)
	if err != nil {
		return "", err
	}

	subnetsIDs := make([]string, 0, len(subnets))
	azs := []string{"a", "b", "c"}
	counter := 0

	for subnetName, subnetCidr := range subnets {
		subnetID, ok := subnetsMap[subnetCidr]
		if !ok {
			subnetID, err = a.createSubnet(ctx, vpc, subnetName, subnetCidr, region, azs[counter%len(azs)])
			if err != nil {
				return "", err
			}

			if cleanup {
				DeferCleanup(func(ctx SpecContext) error {
					return a.deleteSubnet(ctx, subnetID, region)
				})
			}
		}

		subnetsIDs = append(subnetsIDs, subnetID)
		counter++
	}

	a.network = &awsNetwork{
		VPC:     vpc,
		Subnets: subnetsIDs,
	}

	return vpc, nil
}

func (a *AwsAction) CreatePrivateEndpoint(ctx context.Context, serviceName, privateEndpointName, region string) (string, error) {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	createInput := &ec2.CreateVpcEndpointInput{
		ServiceName: aws.String(serviceName),
		SubnetIds:   a.network.Subnets,
		TagSpecifications: []ec2types.TagSpecification{{
			ResourceType: ec2types.ResourceTypeVpcEndpoint,
			Tags: []ec2types.Tag{
				{Key: aws.String("PrivateEndpointName"), Value: aws.String(privateEndpointName)},
			},
		}},
		VpcEndpointType: ec2types.VpcEndpointTypeInterface,
		VpcId:           aws.String(a.network.VPC),
	}
	result, err := ec2Client.CreateVpcEndpoint(ctx, createInput)
	if err != nil {
		return "", err
	}
	vpcEndpointId := aws.ToString(result.VpcEndpoint.VpcEndpointId)

	DeferCleanup(func(ctx SpecContext) error {
		deleteInput := &ec2.DeleteVpcEndpointsInput{
			VpcEndpointIds: []string{vpcEndpointId},
		}
		_, err = ec2Client.DeleteVpcEndpoints(ctx, deleteInput)
		if err != nil {
			return err
		}

		timeout := 10 * time.Minute
		start := time.Now()
		for {
			fmt.Fprintf(GinkgoWriter, "deleting VPC ID %q since %v\n", vpcEndpointId, time.Since(start))

			output, err := ec2Client.DescribeVpcEndpoints(ctx, &ec2.DescribeVpcEndpointsInput{
				VpcEndpointIds: []string{vpcEndpointId},
			})

			if err != nil {
				var apiErr *smithy.GenericAPIError
				if errors.As(err, &apiErr) && apiErr.Code == "InvalidVpcEndpointId.NotFound" {
					return nil
				}
				return err
			}

			if len(output.VpcEndpoints) == 0 {
				return nil
			}

			if time.Since(start) > timeout {
				return errors.New("timeout waiting for deletion of vpc endpoints")
			}

			// we do know that deletion of VPC endpoints takes time
			time.Sleep(3 * time.Second)
		}
	})

	return *result.VpcEndpoint.VpcEndpointId, nil
}

func (a *AwsAction) GetPrivateEndpoint(ctx context.Context, endpointID, region string) (ec2types.VpcEndpoint, error) {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []string{endpointID},
	}

	result, err := ec2Client.DescribeVpcEndpoints(ctx, input)
	if err != nil {
		return ec2types.VpcEndpoint{}, err
	}

	return result.VpcEndpoints[0], nil
}

func (a *AwsAction) AcceptVpcPeeringConnection(ctx context.Context, connectionID, region string) error {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	_, err := ec2Client.AcceptVpcPeeringConnection(
		ctx,
		&ec2.AcceptVpcPeeringConnectionInput{
			VpcPeeringConnectionId: aws.String(connectionID),
		},
	)

	DeferCleanup(func(ctx SpecContext) error {
		_, err = ec2Client.DeleteVpcPeeringConnection(
			ctx,
			&ec2.DeleteVpcPeeringConnectionInput{
				VpcPeeringConnectionId: aws.String(connectionID),
			},
		)
		return err
	})

	return err
}

func (a *AwsAction) findVPC(ctx context.Context, name, region string) (string, error) {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.DescribeVpcsInput{
		Filters: []ec2types.Filter{{
			Name:   aws.String("tag:Name"),
			Values: []string{name},
		}},
	}
	result, err := ec2Client.DescribeVpcs(ctx, input)
	if err != nil {
		return "", err
	}

	if len(result.Vpcs) < 1 {
		return "", nil
	}

	return *result.Vpcs[0].VpcId, nil
}

func (a *AwsAction) createVPC(ctx context.Context, name, cidr, region string) (string, error) {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.CreateVpcInput{
		AmazonProvidedIpv6CidrBlock: aws.Bool(false),
		CidrBlock:                   aws.String(cidr),
		TagSpecifications: []ec2types.TagSpecification{{
			ResourceType: ec2types.ResourceTypeVpc,
			Tags: []ec2types.Tag{
				{Key: aws.String("Name"), Value: aws.String(name)},
				{Key: aws.String(taghelper.OwnerTag), Value: aws.String(taghelper.AKOTeam)},
				{Key: aws.String(taghelper.OwnerEmailTag), Value: aws.String(taghelper.AKOEmail)},
				{Key: aws.String(taghelper.CostCenterTag), Value: aws.String(taghelper.AKOCostCenter)},
				{Key: aws.String(taghelper.EnvironmentTag), Value: aws.String(taghelper.AKOEnvTest)},
			},
		}},
	}

	result, err := ec2Client.CreateVpc(ctx, input)
	if err != nil {
		return "", err
	}

	_, err = ec2Client.ModifyVpcAttribute(ctx, &ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2types.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: result.Vpc.VpcId,
	})
	if err != nil {
		return "", err
	}

	return *result.Vpc.VpcId, nil
}

func (a *AwsAction) deleteVPC(ctx context.Context, id, region string) error {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.DeleteVpcInput{
		DryRun: aws.Bool(false),
		VpcId:  aws.String(id),
	}

	_, err := ec2Client.DeleteVpc(ctx, input)

	return err
}

func (a *AwsAction) getSubnets(ctx context.Context, vpcID, region string) (map[string]string, error) {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	}

	result, err := ec2Client.DescribeSubnets(ctx, input)
	if err != nil {
		return nil, err
	}

	subnetsMap := map[string]string{}

	for _, subnet := range result.Subnets {
		subnetsMap[*subnet.CidrBlock] = aws.ToString(subnet.SubnetId)
	}

	return subnetsMap, nil
}

func (a *AwsAction) createSubnet(ctx context.Context, vpcID, name, cidr, region, az string) (string, error) {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidr),
		TagSpecifications: []ec2types.TagSpecification{{
			ResourceType: ec2types.ResourceTypeSubnet,
			Tags: []ec2types.Tag{
				{Key: aws.String("Name"), Value: aws.String(name)},
				{Key: aws.String(taghelper.OwnerTag), Value: aws.String(taghelper.AKOTeam)},
				{Key: aws.String(taghelper.OwnerEmailTag), Value: aws.String(taghelper.AKOEmail)},
				{Key: aws.String(taghelper.CostCenterTag), Value: aws.String(taghelper.AKOCostCenter)},
				{Key: aws.String(taghelper.EnvironmentTag), Value: aws.String(taghelper.AKOEnvTest)},
			},
		}},
		VpcId:            aws.String(vpcID),
		AvailabilityZone: pointer.MakePtr(fmt.Sprintf("%s%s", region, az)),
	}
	result, err := ec2Client.CreateSubnet(ctx, input)
	if err != nil {
		return "", err
	}

	return aws.ToString(result.Subnet.SubnetId), nil
}

func (a *AwsAction) deleteSubnet(ctx context.Context, subnetID string, region string) error {
	ec2Client := ec2.NewFromConfig(a.cfg, func(o *ec2.Options) {
		o.Region = region
	})

	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
	}
	_, err := ec2Client.DeleteSubnet(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func NewAWSAction(ctx context.Context) (*AwsAction, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &AwsAction{
		cfg: cfg,
	}, nil
}
