package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
)

type sessionAWS struct {
	ec2 *ec2.EC2
	kms *kms.KMS
}

func SessionAWS(region string) sessionAWS {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(session)
	kms := kms.New(session)
	return sessionAWS{svc, kms}
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
		return "", errors.New("can not find VPC")
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
		return "", errors.New("can not find Subnet")
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
				{Key: aws.String(config.TagForTestKey), Value: aws.String(config.TagForTestValue)},
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

func (s sessionAWS) GetFuncPrivateEndpointStatus(privateEndpointID string) func() string {
	return func() string {
		r, err := s.DescribePrivateEndpointStatus(privateEndpointID)
		if err != nil {
			return ""
		}
		return r
	}
}

func (s sessionAWS) GetCustomerMasterKeyID(atlasAccountArn, assumedRoleArn string) (keyId string, err error) {
	keyId, adminARNs, err := getKeyIDAndAdminARNs()
	if err != nil {
		return "", err
	}

	policyString, err := RolePolicyString(atlasAccountArn, assumedRoleArn, adminARNs)
	if err != nil {
		return "", err
	}

	policyInput := &kms.PutKeyPolicyInput{
		KeyId:      &keyId,
		PolicyName: aws.String("default"),
		Policy:     aws.String(policyString),
	}
	_, err = s.kms.PutKeyPolicy(policyInput)
	if err != nil {
		return "", err
	}

	return keyId, nil
}

func getKeyIDAndAdminARNs() (keyID string, adminARNs []string, err error) {
	keyID = os.Getenv("AWS_KMS_KEY_ID")
	if keyID == "" {
		err = errors.New("AWS_KMS_KEY_ID secret is empty")
		return
	}
	adminArnString := os.Getenv("AWS_ACCOUNT_ARN_LIST")
	if adminArnString == "" {
		err = errors.New("AWS_ACCOUNT_ARN_LIST secret is empty")
		return
	}

	adminARNs = strings.Split(adminArnString, ",")
	if len(adminARNs) == 0 {
		err = errors.New("AWS_ACCOUNT_ARN_LIST wasn't parsed properly, please separate accounts via a comma")
		return
	}

	return keyID, adminARNs, nil
}

func RolePolicyString(atlasAccountARN, assumedRoleARN string, adminARNs []string) (string, error) {
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

type principal struct {
	AWS []string `json:"AWS,omitempty"`
}
