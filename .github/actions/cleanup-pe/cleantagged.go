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

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"
)

func cleanAllTaggedAWSPE(region, tagName, tagValue string) error {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return fmt.Errorf("error creating awsSession: %v", err)
	}
	svc := ec2.New(awsSession)
	endpoints, err := svc.DescribeVpcEndpoints(&ec2.DescribeVpcEndpointsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(tagName),
				},
			},
			{
				Name: aws.String("tag:Value"),
				Values: []*string{
					aws.String(tagValue),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error fething all vpcEP with tag %s: %v", tagName, err)
	}
	var endpointIDs []*string
	for _, endpoint := range endpoints.VpcEndpoints {
		endpointIDs = append(endpointIDs, endpoint.VpcEndpointId)
	}
	if len(endpointIDs) > 0 {
		endpointsIDByPortion := chunkSlice(endpointIDs, 25) // aws has a limit of 25 endpointIDs per request
		for _, endpointsIDPortion := range endpointsIDByPortion {
			input := &ec2.DeleteVpcEndpointsInput{
				VpcEndpointIds: endpointsIDPortion,
			}
			_, err = svc.DeleteVpcEndpoints(input)
			if err != nil {
				return fmt.Errorf("error deleting vpcEP: %v", err)
			}
		}
	}
	log.Printf("deleted %d AWS PEs in region %s", len(endpointIDs), region)
	return nil
}

func cleanAllTaggedAzurePE(ctx context.Context, tagName, tagValue, resourceGroupName, azureSubscriptionID string) error {
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
		tags := endpoint.Tags
		if peTagValue, ok := tags[tagName]; ok {
			if peTagValue != nil {
				if *peTagValue == tagValue {
					endpointNames = append(endpointNames, *endpoint.Name)
				}
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

func cleanAllTaggedGCPPE(ctx context.Context, projectID, vpc, region, subnet string) error {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("error while creating new compute service: %v", err)
	}

	networkURL := gcp.FormNetworkURL(vpc, projectID)
	subnetURL := gcp.FormSubnetURL(region, subnet, projectID)

	addressList, err := computeService.Addresses.List(projectID, region).Do()
	if err != nil {
		return fmt.Errorf("error while listing addresses: %v", err)
	}
	var addressNamesToDelete []string
	for _, address := range addressList.Items {
		if address.Network == networkURL && address.Subnetwork == subnetURL {
			addressNamesToDelete = append(addressNamesToDelete, address.Name)
		}
	}

	for _, addressName := range addressNamesToDelete {
		err = deleteGCPForwardRuleByAddressName(computeService, projectID, region, addressName)
		if err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 20) // need to wait for GCP to delete the forwarding rule
	for _, addressName := range addressNamesToDelete {
		_, err = computeService.Addresses.Delete(projectID, region, addressName).Do()
		if err != nil {
			return fmt.Errorf("error while deleting address: %v", err)
		}
		log.Printf("successfully deleted GCP PE %s", addressName)
	}

	return nil
}

func deleteGCPForwardRuleByAddressName(service *compute.Service, projectID, region, addressName string) error {
	ruleName, err := cloud.FormForwardRuleNameByAddressName(addressName)
	if err != nil {
		return fmt.Errorf("error while forming forward rule name: %v", err)
	}
	_, err = service.ForwardingRules.Delete(projectID, region, ruleName).Do()
	if err != nil {
		return fmt.Errorf("error while deleting forwarding rule: %v", err)
	}
	log.Printf("successfully deleted forwarding rule %s", ruleName)
	return nil
}

func chunkSlice(slice []*string, chunkSize int) [][]*string {
	var chunks [][]*string
	for {
		if len(slice) == 0 {
			break
		}
		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}
		chunk := slice[:chunkSize]
		chunks = append(chunks, chunk)
		slice = slice[chunkSize:]
	}

	return chunks
}
