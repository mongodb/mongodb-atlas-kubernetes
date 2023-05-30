package cloud

import (
	aws_sdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
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

func NewAwsAction() *AwsAction {
	return new(AwsAction)
}

func (awsAction *AwsAction) CreateKMS(region, atlasAccountArn, assumedRoleArn string) (key string, err error) {
	session := aws.SessionAWS(region)
	return session.GetCustomerMasterKeyID(atlasAccountArn, assumedRoleArn)
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

func (a *AwsAction) InitNetwork(vpcName, cidr, region string, subnets []string, cleanup bool) (string, error) {
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

	for _, subnet := range subnets {
		subnetID, ok := subnetsMap[subnet]
		if !ok {
			subnetID, err = a.createSubnet(vpc, subnet, region)
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
	}

	a.network = &awsNetwork{
		VPC:     vpc,
		Subnets: subnetsIDs,
	}

	return vpc, nil
}

func (a *AwsAction) CreatePrivateEndpoint(serviceName, privateEndpointName, region string) (string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	createInput := &ec2.CreateVpcEndpointInput{
		ServiceName: aws_sdk.String(serviceName),
		SubnetIds:   a.network.Subnets,
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws_sdk.String(ec2.ResourceTypeVpcEndpoint),
			Tags: []*ec2.Tag{
				{Key: aws_sdk.String("Name"), Value: aws_sdk.String(config.TagName)},
				{Key: aws_sdk.String("PrivateEndpointName"), Value: aws_sdk.String(privateEndpointName)},
			},
		}},
		VpcEndpointType: aws_sdk.String("Interface"),
		VpcId:           aws_sdk.String(a.network.VPC),
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
	})

	return *result.VpcEndpoint.VpcEndpointId, nil
}

func (a *AwsAction) GetPrivateEndpoint(endpointID, region string) (*ec2.VpcEndpoint, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.DescribeVpcEndpointsInput{
		VpcEndpointIds: []*string{aws_sdk.String(endpointID)},
	}

	result, err := ec2Client.DescribeVpcEndpoints(input)
	if err != nil {
		return nil, err
	}

	return result.VpcEndpoints[0], nil
}

func (a *AwsAction) AcceptVpcPeeringConnection(connectionID, region string) error {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))
	_, err := ec2Client.AcceptVpcPeeringConnection(
		&ec2.AcceptVpcPeeringConnectionInput{
			VpcPeeringConnectionId: aws_sdk.String(connectionID),
		},
	)

	a.t.Cleanup(func() {
		_, err = ec2Client.DeleteVpcPeeringConnection(
			&ec2.DeleteVpcPeeringConnectionInput{
				VpcPeeringConnectionId: aws_sdk.String(connectionID),
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

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name: aws_sdk.String("tag:Name"),
			Values: []*string{
				aws_sdk.String(name),
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

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.CreateVpcInput{
		AmazonProvidedIpv6CidrBlock: aws_sdk.Bool(false),
		CidrBlock:                   aws_sdk.String(cidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws_sdk.String(ec2.ResourceTypeVpc),
			Tags: []*ec2.Tag{
				{Key: aws_sdk.String("Name"), Value: aws_sdk.String(name)},
			},
		}},
	}

	result, err := ec2Client.CreateVpc(input)
	if err != nil {
		return "", err
	}

	return *result.Vpc.VpcId, nil
}

func (a *AwsAction) deleteVPC(id, region string) error {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.DeleteVpcInput{
		DryRun: aws_sdk.Bool(false),
		VpcId:  aws_sdk.String(id),
	}

	_, err := ec2Client.DeleteVpc(input)

	return err
}

func (a *AwsAction) getSubnets(vpcID, region string) (map[string]*string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws_sdk.String("vpc-id"),
				Values: []*string{
					aws_sdk.String(vpcID),
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

func (a *AwsAction) createSubnet(vpcID, cidr, region string) (*string, error) {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.CreateSubnetInput{
		CidrBlock: aws_sdk.String(cidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws_sdk.String(ec2.ResourceTypeSubnet),
			Tags: []*ec2.Tag{
				{Key: aws_sdk.String("Name"), Value: aws_sdk.String(config.TagName)},
			},
		}},
		VpcId: aws_sdk.String(vpcID),
	}
	result, err := ec2Client.CreateSubnet(input)
	if err != nil {
		return nil, err
	}

	return result.Subnet.SubnetId, nil
}

func (a *AwsAction) deleteSubnet(subnetID, region string) error {
	a.t.Helper()

	ec2Client := ec2.New(a.session, aws_sdk.NewConfig().WithRegion(region))

	input := &ec2.DeleteSubnetInput{
		SubnetId: aws_sdk.String(subnetID),
	}
	_, err := ec2Client.DeleteSubnet(input)
	if err != nil {
		return err
	}

	return nil
}

func NewAWSAction(t core.GinkgoTInterface) (*AwsAction, error) {
	t.Helper()

	awsSession, err := session.NewSession(aws_sdk.NewConfig())
	if err != nil {
		return nil, err
	}

	return &AwsAction{
		t: t,

		session: awsSession,
	}, nil
}
