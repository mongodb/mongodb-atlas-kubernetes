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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type AwsAction struct {
	t       core.GinkgoTInterface
	network *awsNetwork

	session *session.Session
}

type awsNetwork struct {
	VPC     string
	Subnets []*string
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

func (a *AwsAction) CreateKMS(alias, region, atlasAccountArn, assumedRoleArn string) (string, error) {
	a.t.Helper()

	kmsClient := kms.New(a.session, aws.NewConfig().WithRegion(region))

	adminARNs, err := getAdminARNs()
	if err != nil {
		return "", err
	}

	policyString, err := rolePolicyString(atlasAccountArn, assumedRoleArn, adminARNs)
	if err != nil {
		return "", err
	}

	key, err := kmsClient.CreateKey(&kms.CreateKeyInput{
		Description: aws.String("Key for E2E test"),
		KeySpec:     aws.String("SYMMETRIC_DEFAULT"),
		KeyUsage:    aws.String("ENCRYPT_DECRYPT"),
		MultiRegion: aws.Bool(false),
		Origin:      aws.String("AWS_KMS"),
		Policy:      aws.String(policyString),
	})

	if err != nil {
		return "", err
	}

	_, err = kmsClient.CreateAlias(&kms.CreateAliasInput{
		AliasName:   aws.String("alias/" + strings.ToLower(strings.ReplaceAll(alias, " ", "-"))),
		TargetKeyId: key.KeyMetadata.KeyId,
	})

	if err != nil {
		a.t.Log(fmt.Sprintf("failed to create alias to key %s(%s): %s", alias, *key.KeyMetadata.KeyId, err))
	}

	a.t.Cleanup(func() {
		_, err = kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
			KeyId:               key.KeyMetadata.KeyId,
			PendingWindowInDays: aws.Int64(7), // this is the minimum possible and can be up to 24h longer than value set
		})
		if err != nil {
			a.t.Error(err)
		}
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

func (a *AwsAction) GetAccountID() (string, error) {
	a.t.Helper()

	stsClient := sts.New(a.session)
	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *identity.Account, nil
}

func (a *AwsAction) InitNetwork(vpcName, cidr, region string, subnets map[string]string, cleanup bool) (string, error) {
	a.t.Helper()

	vpc, err := a.findVPC(vpcName, region)
	if err != nil {
		return "", err
	}

	if vpc == "" {
		vpc, err = a.createVPC(vpcName, cidr, region)
		if err != nil {
			return "", err
		}
	}

	if cleanup {
		a.t.Cleanup(func() {
			err = a.deleteVPC(vpc, region)
			if err != nil {
				a.t.Error(err)
			}
		})
	}

	subnetsMap, err := a.getSubnets(vpc, region)
	if err != nil {
		return "", err
	}

	subnetsIDs := make([]*string, 0, len(subnets))
	azs := []string{"a", "b", "c"}
	counter := 0

	for subnetName, subnetCidr := range subnets {
		subnetID, ok := subnetsMap[subnetCidr]
		if !ok {
			subnetID, err = a.createSubnet(vpc, subnetName, subnetCidr, region, azs[counter%len(azs)])
			if err != nil {
				return "", err
			}

			if cleanup {
				a.t.Cleanup(func() {
					err = a.deleteSubnet(*subnetID, region)
					if err != nil {
						a.t.Error(err)
					}
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

func (a *AwsAction) CreatePrivateEndpoint(serviceName, privateEndpointName, region string) (string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	createInput := &ec2.CreateVpcEndpointInput{
		ServiceName: aws.String(serviceName),
		SubnetIds:   a.network.Subnets,
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeVpcEndpoint),
			Tags: []*ec2.Tag{
				{Key: aws.String("PrivateEndpointName"), Value: aws.String(privateEndpointName)},
			},
		}},
		VpcEndpointType: aws.String("Interface"),
		VpcId:           aws.String(a.network.VPC),
	}
	result, err := ec2Client.CreateVpcEndpoint(createInput)
	if err != nil {
		return "", err
	}

	a.t.Cleanup(func() {
		deleteInput := &ec2.DeleteVpcEndpointsInput{
			VpcEndpointIds: []*string{result.VpcEndpoint.VpcEndpointId},
		}
		_, err = ec2Client.DeleteVpcEndpoints(deleteInput)
		if err != nil {
			a.t.Error(err)
		}

		timeout := 10 * time.Minute
		start := time.Now()
		for {
			a.t.Log(fmt.Sprintf("deleting VPC ID %q since %v", aws.StringValue(result.VpcEndpoint.VpcEndpointId), time.Since(start)))

			output, err := ec2Client.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{
				VpcEndpointIds: []*string{result.VpcEndpoint.VpcEndpointId},
			})

			var e awserr.Error
			if (errors.As(err, &e) && e.Code() == "InvalidVpcEndpointId.NotFound") || len(output.VpcEndpoints) == 0 {
				return
			}

			if err != nil {
				a.t.Error(err)
				return
			}

			if time.Since(start) > timeout {
				a.t.Error(errors.New("timeout waiting for deletion of vpc endpoints"))
				return
			}

			// we do know that deletion of VPC endpoints takes time
			time.Sleep(3 * time.Second)
		}
	})

	return *result.VpcEndpoint.VpcEndpointId, nil
}

func (a *AwsAction) GetPrivateEndpoint(endpointID, region string) (*ec2.VpcEndpoint, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(endpointID)},
	}

	result, err := ec2Client.DescribeVpcEndpoints(input)
	if err != nil {
		return nil, err
	}

	return result.VpcEndpoints[0], nil
}

func (a *AwsAction) AcceptVpcPeeringConnection(connectionID, region string) error {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))
	_, err := ec2Client.AcceptVpcPeeringConnection(
		&ec2.AcceptVpcPeeringConnectionInput{
			VpcPeeringConnectionId: aws.String(connectionID),
		},
	)

	a.t.Cleanup(func() {
		_, err = ec2Client.DeleteVpcPeeringConnection(
			&ec2.DeleteVpcPeeringConnectionInput{
				VpcPeeringConnectionId: aws.String(connectionID),
			},
		)
		if err != nil {
			a.t.Error(err)
		}
	})

	return err
}

func (a *AwsAction) findVPC(name, region string) (string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String(name),
			},
		}},
	}
	result, err := ec2Client.DescribeVpcs(input)
	if err != nil {
		return "", err
	}

	if len(result.Vpcs) < 1 {
		return "", nil
	}

	return *result.Vpcs[0].VpcId, nil
}

func (a *AwsAction) createVPC(name, cidr, region string) (string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.CreateVpcInput{
		AmazonProvidedIpv6CidrBlock: aws.Bool(false),
		CidrBlock:                   aws.String(cidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeVpc),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(name)},
			},
		}},
	}

	result, err := ec2Client.CreateVpc(input)
	if err != nil {
		return "", err
	}

	_, err = ec2Client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: result.Vpc.VpcId,
	})
	if err != nil {
		return "", err
	}

	return *result.Vpc.VpcId, nil
}

func (a *AwsAction) deleteVPC(id, region string) error {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.DeleteVpcInput{
		DryRun: aws.Bool(false),
		VpcId:  aws.String(id),
	}

	_, err := ec2Client.DeleteVpc(input)

	return err
}

func (a *AwsAction) getSubnets(vpcID, region string) (map[string]*string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcID),
				},
			},
		},
	}

	result, err := ec2Client.DescribeSubnets(input)
	if err != nil {
		return nil, err
	}

	subnetsMap := map[string]*string{}

	for _, subnet := range result.Subnets {
		subnetsMap[*subnet.CidrBlock] = subnet.SubnetId
	}

	return subnetsMap, nil
}

func (a *AwsAction) createSubnet(vpcID, name, cidr, region, az string) (*string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeSubnet),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(name)},
			},
		}},
		VpcId:            aws.String(vpcID),
		AvailabilityZone: pointer.MakePtr(fmt.Sprintf("%s%s", region, az)),
	}
	result, err := ec2Client.CreateSubnet(input)
	if err != nil {
		return nil, err
	}

	return result.Subnet.SubnetId, nil
}

func (a *AwsAction) deleteSubnet(subnetID, region string) error {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws.NewConfig().WithRegion(region))

	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
	}
	_, err := ec2Client.DeleteSubnet(input)
	if err != nil {
		return err
	}

	return nil
}

func NewAWSAction(t core.GinkgoTInterface) (*AwsAction, error) {
	t.Helper()

	awsSession, err := session.NewSession(aws.NewConfig())
	if err != nil {
		return nil, err
	}

	return &AwsAction{
		t: t,

		session: awsSession,
	}, nil
}
