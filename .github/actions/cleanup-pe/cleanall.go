package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"google.golang.org/api/compute/v1"

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
	err = deleteAWSPEsByID(svc, endpointIDs)
	if err != nil {
		return err
	}
	log.Printf("deleted %d AWS PEs in region %s", len(endpointIDs), region)
	return nil
}

func deleteAWSPEsByID(svc *ec2.EC2, endpointIDs []*string) error {
	if len(endpointIDs) > 0 {
		endpointsIDByPortion := chunkSlice(endpointIDs, 25) // aws has a limit of 25 endpointIDs per request
		for _, endpointsIDPortion := range endpointsIDByPortion {
			input := &ec2.DeleteVpcEndpointsInput{
				VpcEndpointIds: endpointsIDPortion,
			}
			_, err := svc.DeleteVpcEndpoints(input)
			if err != nil {
				return fmt.Errorf("error deleting vpcEP: %v", err)
			}
		}
	}
	return nil
}

func cleanAllAzurePE(ctx context.Context, resourceGroupName, azureSubscriptionID, subnet string) error {
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
			if *endpoint.Subnet.Name == subnet {
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
	log.Printf("deleted %d Azure PEs", len(endpointNames))
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

	counter := 0
	for _, forwardRule := range forwardRules.Items {
		if forwardRule.Network == networkURL {
			_, err = computeService.ForwardingRules.Delete(projectID, region, forwardRule.Name).Do()
			if err != nil {
				return fmt.Errorf("error while deleting forwarding rule: %v", err)
			}

			counter++
			log.Printf("successfully deleted GCP forward rule: %s. network:  %s",
				forwardRule.Name, forwardRule.Network)
		}
	}
	log.Printf("deleted %d GCP Forfard rules", counter)

	time.Sleep(time.Second * 20) // need to wait for GCP to delete the forwarding rule
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
	counter := 0
	for _, address := range addressList.Items {
		if address.Subnetwork == subnetURL {
			_, err = service.Addresses.Delete(projectID, region, address.Name).Do()
			if err != nil {
				return fmt.Errorf("error while deleting address: %v", err)
			}
			counter++
			log.Printf("successfully deleted GCP address: %s. subnet: %s", address.Name, address.Subnetwork)
		}
	}
	log.Printf("deleted %d GCP addresses", counter)
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
