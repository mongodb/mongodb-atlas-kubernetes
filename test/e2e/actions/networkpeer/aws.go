package networkpeer

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

type AWSNetworkPeer struct {
	ec2 *ec2.EC2
}

func NewAWSNetworkPeerService(region string) (AWSNetworkPeer, error) {
	newSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return AWSNetworkPeer{}, err
	}
	svc := ec2.New(newSession)
	return AWSNetworkPeer{svc}, nil
}

func (p *AWSNetworkPeer) CreateVPCForAWS(appCidr string) (string, string, error) {
	vpcInput := &ec2.CreateVpcInput{
		AmazonProvidedIpv6CidrBlock: aws.Bool(false),
		CidrBlock:                   aws.String(appCidr),
		TagSpecifications: []*ec2.TagSpecification{{
			ResourceType: aws.String(ec2.ResourceTypeVpc),
			Tags: []*ec2.Tag{
				{Key: aws.String("Name"), Value: aws.String(config.TagName)},
			},
		}},
	}
	vpc, err := p.ec2.CreateVpc(vpcInput)
	if err != nil {
		return "", "", err
	}
	return *vpc.Vpc.OwnerId, *vpc.Vpc.VpcId, nil
}

func EstablishAWSPeerConnection(peer status.AtlasNetworkPeer) error {
	if peer.Region == "" {
		return fmt.Errorf("region is required for %s", peer.VPC)
	}
	if peer.ProviderName == "" {
		return fmt.Errorf("providerName is required for %s", peer.VPC)
	}
	newSession, err := session.NewSession(&aws.Config{
		Region: aws.String(peer.Region)},
	)
	if err != nil {
		return fmt.Errorf("failed to create AWS session for status %v: %w", peer, err)
	}
	svc := ec2.New(newSession)
	_, err = svc.AcceptVpcPeeringConnection(&ec2.AcceptVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(peer.ConnectionID),
	})
	if err != nil {
		return err
	}
	return nil
}

func DeletePeerConnectionAndVPCForAWS(conID string, region string) error {
	newSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return err
	}
	svc := ec2.New(newSession)
	connections, err := svc.DescribeVpcPeeringConnections(&ec2.DescribeVpcPeeringConnectionsInput{
		VpcPeeringConnectionIds: []*string{aws.String(conID)},
	})
	if err != nil {
		return err
	}
	if len(connections.VpcPeeringConnections) == 0 {
		return nil
	}
	vpcID := connections.VpcPeeringConnections[0].AccepterVpcInfo.VpcId
	_, err = svc.DeleteVpcPeeringConnection(&ec2.DeleteVpcPeeringConnectionInput{
		VpcPeeringConnectionId: aws.String(conID),
	})
	if err != nil {
		return err
	}
	_, err = svc.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: vpcID,
	})
	return err
}
