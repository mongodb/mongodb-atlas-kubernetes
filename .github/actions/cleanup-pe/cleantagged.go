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

func cleanAllTaggedGCPPE(ctx context.Context, projectID, vpc, region, tagName, tagValue string) error {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("error while creating new compute service: %v", err)
	}

	networkURL := gcp.FormNetworkURL(vpc, projectID)

	forwardRules, err := computeService.ForwardingRules.List(projectID, region).Do()
	if err != nil {
		return fmt.Errorf("error while listing forwarding rules: %v", err)
	}

	for _, forwardRule := range forwardRules.Items {
		forwardRuleLabels := forwardRule.Labels
		if forwardRuleLabels[tagName] == tagValue && forwardRule.Network == networkURL {
			_, err = computeService.ForwardingRules.Delete(projectID, region, forwardRule.Name).Do()
			if err != nil {
				return fmt.Errorf("error while deleting forwarding rule: %v", err)
			}
			ruleName := forwardRule.Name
			log.Printf("successfully deleted GCP forward rule: %s", ruleName)
			err = deleteGCPAddressByForwardRuleName(computeService, projectID, region, ruleName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteGCPAddressByForwardRuleName(service *compute.Service, projectID, region, ruleName string) error {
	addressName, err := cloud.FormAddressNameByRuleName(ruleName)
	if err != nil {
		return fmt.Errorf("unexpected forvard rule name pattern: %v", err)
	}
	_, err = service.Addresses.Delete(projectID, region, addressName).Do()
	if err != nil {
		return fmt.Errorf("error while deleting address: %v", err)
	}
	log.Printf("successfully deleted GCP address: %s", addressName)
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
