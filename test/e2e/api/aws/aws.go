package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type sessionAWS struct {
	ec2 *ec2.EC2
}

func SessionAWS(region string) sessionAWS { // eu-west-2
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(session)
	return sessionAWS{svc}
}

func (s sessionAWS) GetVPCID() (string, error) {
	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String(config.TagName),
			},
		}},
	}
	result, err := s.ec2.DescribeVpcs(input)
	if err != nil {
		return "", getError(err)
	}
	if len(result.Vpcs) < 1 {
		return "", errors.New("Can not find VPC")
	}
	fmt.Println(result)
	return *result.Vpcs[0].VpcId, nil
}

func (s sessionAWS) CreateVPC(testID string) (string, error) {
	input := &ec2.CreateVpcInput{
		AmazonProvidedIpv6CidrBlock: aws.Bool(false),
		CidrBlock:                   aws.String("10.0.0.0/24"),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeVpc),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(config.TagName)},
				{Key: aws.String("Test"), Value: aws.String(testID)},
			},
		}},
	}
	result, err := s.ec2.CreateVpc(input)
	if err != nil {
		return "", getError(err)
	}
	fmt.Println(result)
	return *result.Vpc.VpcId, nil
}

func (s sessionAWS) DescribeVPCStatus(vpcID string) (string, error) {
	input := &ec2.DescribeVpcsInput{
		DryRun: aws.Bool(false),
		VpcIds: []*string{aws.String(vpcID)},
	}
	result, err := s.ec2.DescribeVpcs(input)
	if err != nil {
		return "", getError(err)
	}
	return *result.Vpcs[0].State, nil
}

func (s sessionAWS) DeleteVPC(vpcID string) error {
	input := &ec2.DeleteVpcInput{
		DryRun: aws.Bool(false),
		VpcId:  aws.String(vpcID),
	}
	_, err := s.ec2.DeleteVpc(input)
	if err != nil {
		return getError(err)
	}
	return nil
}

func (s sessionAWS) GetSubnetID() (string, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String(config.TagName),
			},
		}},
	}
	result, err := s.ec2.DescribeSubnets(input)
	if err != nil {
		return "", getError(err)
	}
	if len(result.Subnets) < 1 {
		return "", errors.New("Can not find Subnet")
	}
	fmt.Println(result)
	return *result.Subnets[0].SubnetId, nil
}

func (s sessionAWS) CreateSubnet(vpcID, cidr, testID string) (string, error) {
	input := &ec2.CreateSubnetInput{
		CidrBlock: aws.String(cidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeSubnet),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(config.TagName)},
				{Key: aws.String("Test"), Value: aws.String(testID)},
			},
		}},
		VpcId: aws.String(vpcID),
	}
	result, err := s.ec2.CreateSubnet(input)
	if err != nil {
		return "", getError(err)
	}
	fmt.Println(result)
	return *result.Subnet.SubnetId, nil
}

func (s sessionAWS) DescribeSubnetStatus(subnetID string) (string, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{aws.String(subnetID)},
	}
	result, err := s.ec2.DescribeSubnets(input)
	if err != nil {
		return "", getError(err)
	}
	return *result.Subnets[0].State, nil
}

func (s sessionAWS) DeleteSubnet(subnetID string) error {
	input := &ec2.DeleteSubnetInput{
		SubnetId: aws.String(subnetID),
	}
	_, err := s.ec2.DeleteSubnet(input)
	if err != nil {
		return getError(err)
	}
	return nil
}

func (s sessionAWS) CreatePrivateEndpoint(vpcID, subnetID, serviceName, testID string) (string, error) {
	input := &ec2.CreateVpcEndpointInput{
		ServiceName: aws.String(serviceName),
		SubnetIds:   []*string{aws.String(subnetID)},
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeVpcEndpoint),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(config.TagName)},
				{Key: aws.String("Test"), Value: aws.String(testID)},
			},
		}},
		VpcEndpointType: aws.String("Interface"),
		VpcId:           aws.String(vpcID),
	}
	result, err := s.ec2.CreateVpcEndpoint(input)
	if err != nil {
		return "", getError(err)
	}
	fmt.Println(result)
	return *result.VpcEndpoint.VpcEndpointId, nil
}

func (s sessionAWS) DescribePrivateEndpointStatus(endpointID string) (string, error) {
	input := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(endpointID)},
	}
	result, err := s.ec2.DescribeVpcEndpoints(input)
	if err != nil {
		return "", getError(err)
	}
	return *result.VpcEndpoints[0].State, nil
}

func (s sessionAWS) DeletePrivateLink(endpointID string) error {
	input := &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: []*string{aws.String(endpointID)},
	}
	_, err := s.ec2.DeleteVpcEndpoints(input)
	if err != nil {
		return getError(err)
	}
	return nil
}

func getError(err error) error {
	if aerr, ok := err.(awserr.Error); ok { // nolint
		return aerr
	}
	return err
}
