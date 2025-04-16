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
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
)

type AWS struct{}

func (a *AWS) DeleteVpc(ID, region string) error {
	awsSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String(normalizeAwsRegion(region)),
		},
	)
	if err != nil {
		return err
	}

	client := ec2.New(awsSession)

	vpc, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(ID)},
	})
	if err != nil {
		var e awserr.Error
		if errors.As(err, &e) && e.Code() == "InvalidVpcID.NotFound" {
			return nil
		}
		return err
	}

	if len(vpc.Vpcs) == 0 {
		return nil
	}

	err = a.deletePeeringConnection(client, ID)
	if err != nil {
		return err
	}

	_, err = client.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: aws.String(ID),
	})

	return err
}

func (a *AWS) DeleteEndpoint(ID, region string) error {
	awsSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String(normalizeAwsRegion(region)),
		},
	)
	if err != nil {
		return err
	}

	client := ec2.New(awsSession)

	endpoints, err := client.DescribeVpcEndpointServices(&ec2.DescribeVpcEndpointServicesInput{
		ServiceNames: []*string{aws.String(ID)},
	})
	if err != nil {
		var e awserr.Error
		if errors.As(err, &e) && e.Code() == "InvalidVpcEndpointID.NotFound" {
			return nil
		}
		return err
	}

	if len(endpoints.ServiceNames) == 0 {
		return nil
	}

	endpointIDs := make([]*string, 0, len(endpoints.ServiceNames))

	for _, serviceName := range endpoints.ServiceNames {
		endpointID := strings.TrimPrefix(*serviceName, fmt.Sprintf("com.amazonaws.vpce.%s.", normalizeAwsRegion(region)))
		endpointIDs = append(endpointIDs, &endpointID)
	}

	_, err = client.DeleteVpcEndpointServiceConfigurations(&ec2.DeleteVpcEndpointServiceConfigurationsInput{
		ServiceIds: endpointIDs,
	})

	return err
}

func (a *AWS) DeleteKMS(ID, region string) error {
	awsSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String(normalizeAwsRegion(region)),
		},
	)
	if err != nil {
		return err
	}

	kmsClient := kms.New(awsSession)

	_, err = kmsClient.DescribeKey(&kms.DescribeKeyInput{
		GrantTokens: nil,
		KeyId:       aws.String(ID),
	})
	if err != nil {
		var e awserr.Error
		if errors.As(err, &e) && e.Code() == "NotFoundException" {
			return nil
		}
		return err
	}

	_, err = kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
		KeyId:               aws.String(ID),
		PendingWindowInDays: aws.Int64(7), // this is the minimum possible and can be up to 24h longer than value set
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *AWS) deletePeeringConnection(ec2Client *ec2.EC2, vpcID string) error {
	input := ec2.DescribeVpcPeeringConnectionsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("accepter-vpc-info.vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	peers, err := ec2Client.DescribeVpcPeeringConnections(&input)
	if err != nil {
		return err
	}

	if len(peers.VpcPeeringConnections) == 0 {
		return nil
	}

	for _, peer := range peers.VpcPeeringConnections {
		_, err = ec2Client.DeleteVpcPeeringConnection(&ec2.DeleteVpcPeeringConnectionInput{
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
