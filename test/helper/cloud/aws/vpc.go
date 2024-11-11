package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateVPC(name, cidr, region string) (string, error) {
	awsSession, err := newSession(region)
	if err != nil {
		return "", fmt.Errorf("failed to create an AWS session: %w", err)
	}
	ec2Client := ec2.New(awsSession)
	result, err := ec2Client.CreateVpc(&ec2.CreateVpcInput{
		AmazonProvidedIpv6CidrBlock: aws.Bool(false),
		CidrBlock:                   aws.String(cidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeVpc),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(name)},
			},
		}},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create an AWS VPC: %w", err)
	}

	_, err = ec2Client.ModifyVpcAttribute(&ec2.ModifyVpcAttributeInput{
		EnableDnsHostnames: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
		VpcId: result.Vpc.VpcId,
	})
	if err != nil {
		return "", fmt.Errorf("failed to configure AWS VPC: %w", err)
	}

	return *result.Vpc.VpcId, nil
}

func DeleteVPC(vpcID, region string) error {
	awsSession, err := newSession(region)
	if err != nil {
		return fmt.Errorf("failed to create an AWS session: %w", err)
	}
	ec2Client := ec2.New(awsSession)
	_, err = ec2Client.DeleteVpc(&ec2.DeleteVpcInput{
		DryRun: aws.Bool(false),
		VpcId:  aws.String(vpcID),
	})
	return err
}
