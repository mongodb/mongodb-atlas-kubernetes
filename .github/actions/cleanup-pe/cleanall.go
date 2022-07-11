package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"google.golang.org/api/compute/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

func cleanAllAWSPE(region string) error {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return fmt.Errorf("error creating awsSession: %v", err)
	}
	svc := ec2.New(awsSession)

	subnetInput := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String(config.TagName),
			},
		}},
	}
	subnetOutput, err := svc.DescribeSubnets(subnetInput)
	if err != nil {
		return fmt.Errorf("error while listing subnets: %v", err)
	}
	if len(subnetOutput.Subnets) == 0 {
		return fmt.Errorf("no subnets found")
	}
	subnetID := subnetOutput.Subnets[0].SubnetId

	endpoints, err := svc.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{})
	if err != nil {
		return fmt.Errorf("error fething all vpcEP: %v", err)
	}
	var endpointIDs []*string
	for _, endpoint := range endpoints.VpcEndpoints {
		if containsPtr(endpoint.SubnetIds, subnetID) {
			endpointIDs = append(endpointIDs, endpoint.VpcEndpointId)
		}
	}
	if len(endpointIDs) > 0 {
		input := &ec2.DeleteVpcEndpointsInput{
			VpcEndpointIds: endpointIDs,
		}
		_, err = svc.DeleteVpcEndpoints(input)
		if err != nil {
			return fmt.Errorf("error deleting vpcEP: %v", err)
		}
	}
	log.Printf("deleted %d AWS PEs in region %s", len(endpointIDs), region)
	return nil
}

func cleanAllAzurePE(ctx context.Context, resourceGroupName, azureSubscriptionID string) error {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return fmt.Errorf("error creating authorizer: %v", err)
	}
	peClient := network.NewPrivateEndpointsClient(azureSubscriptionID)
	peClient.Authorizer = authorizer

	peList, err := peClient.List(ctx, resourceGroupName)
	if err != nil {
		return fmt.Errorf("error fething all PE: %v", err)
	}
	var endpointNames []string
	for _, endpoint := range peList.Values() {
		if endpoint.Subnet.Name != nil {
			if *endpoint.Subnet.Name == cloud.SubnetName {
				endpointNames = append(endpointNames, *endpoint.Name)
			}
		}

	}

	for _, peName := range endpointNames {
		_, errDelete := peClient.Delete(ctx, resourceGroupName, peName)
		if errDelete != nil {
			return errDelete
		}
		log.Printf("successfully deleted Azure PE %s", peName)
	}
	return nil
}

func cleanAllGCPPE(ctx context.Context, projectID, vpc, region, subnet string) error {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("error while creating new compute service: %v", err)
	}

	networkURL := gcp.FormNetworkURL(vpc, projectID)
	subnetURL := gcp.FormSubnetURL(region, subnet, projectID)

	forwardRules, err := computeService.ForwardingRules.List(projectID, region).Do()
	if err != nil {
		return fmt.Errorf("error while listing forwarding rules: %v", err)
	}

	for _, forwardRule := range forwardRules.Items {
		log.Printf("deleting forwarding rule %s. subnet %s. network %s", forwardRule.Name, forwardRule.Subnetwork, forwardRule.Network) // TODO: remove this line
		if forwardRule.Network == networkURL {
			_, err = computeService.ForwardingRules.Delete(projectID, region, forwardRule.Name).Do()
			if err != nil {
				return fmt.Errorf("error while deleting forwarding rule: %v", err)
			}
			ruleName := forwardRule.Name
			log.Printf("successfully deleted GCP forward rule: %s", ruleName)
		}
	}

	err = deleteGCPAddressBySubnet(computeService, projectID, region, subnetURL)
	if err != nil {
		return fmt.Errorf("error while deleting GCP address: %v", err)
	}

	return nil
}

func deleteGCPAddressBySubnet(service *compute.Service, projectID, region, subnetURL string) error {
	addressList, err := service.Addresses.List(projectID, region).Do()
	if err != nil {
		return fmt.Errorf("error while listing addresses: %v", err)
	}
	for _, address := range addressList.Items {
		if address.Subnetwork == subnetURL {
			_, err = service.Addresses.Delete(projectID, region, address.Name).Do()
			if err != nil {
				return fmt.Errorf("error while deleting address: %v", err)
			}
			log.Printf("successfully deleted GCP address: %s", address.Name)
		}
	}
	return nil
}

func containsPtr(slice []*string, elem *string) bool {
	for _, s := range slice {
		if s != nil && elem != nil {
			if *s == *elem {
				return true
			}
		}
	}
	return false
}
