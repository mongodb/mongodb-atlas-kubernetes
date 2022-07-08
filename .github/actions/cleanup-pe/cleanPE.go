package main

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"google.golang.org/api/compute/v1"
	"log"
	"os"
)

func main() {
	err := CleanAllPE()
	if err != nil {
		log.Fatal(err)
	}
}

func CleanAllPE() error {
	ctx := context.Background()
	groupNameAzure := cloud.ResourceGroup
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	err := cleanAllAzurePE(ctx, config.TagForTestKey, config.TagForTestValue, groupNameAzure, subscriptionID)
	if err != nil {
		return fmt.Errorf("error while cleaning all azure pe: %v", err)
	}

	awsRegions := []string{
		config.AWSRegionEU,
		config.AWSRegionUS,
	}
	for _, awsRegion := range awsRegions {
		errClean := cleanAllAWSPE(awsRegion, config.TagForTestKey, config.TagForTestValue)
		if errClean != nil {
			return fmt.Errorf("error cleaning all aws PE. region %s. error: %v", awsRegion, errClean)
		}
	}

	gcpRegion := config.GCPRegion
	err = cleanAllGCPPE(ctx, cloud.GoogleProjectID, cloud.GoogleVPC, cloud.GoogleSubnetName,
		gcpRegion, config.TagForTestKey, config.TagForTestValue)
	if err != nil {
		return fmt.Errorf("error while cleaning all gcp pe: %v", err)
	}
	return nil
}

func cleanAllAWSPE(region, tagName, tagValue string) error {
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
	input := &ec2.DeleteVpcEndpointsInput{
		VpcEndpointIds: endpointIDs,
	}
	_, err = svc.DeleteVpcEndpoints(input)
	if err != nil {
		return fmt.Errorf("error deleting vpcEP: %v", err)
	}
	return nil
}

func cleanAllAzurePE(ctx context.Context, tagName, tagValue, resourceGroupName, azureSubscriptionID string) error {
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
	}
	return nil
}

func cleanAllGCPPE(ctx context.Context, projectID, vpc, subnetName, region, tagName, tagValue string) error {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("error while creating new compute service: %v", err)
	}
	subnet := gcp.FormSubnetURL(region, subnetName, projectID)

	addressFilter := fmt.Sprintf("subnet=%s", subnet)
	networkURL := gcp.FormNetworkURL(vpc, projectID)
	forwardRuleFilter := fmt.Sprintf("(%s) AND (network=%s)", addressFilter, networkURL)

	forwardRules, err := computeService.ForwardingRules.List(projectID, region).Filter(forwardRuleFilter).Do()
	if err != nil {
		return fmt.Errorf("error while listing forwarding rules: %v", err)
	}

	for _, forwardRule := range forwardRules.Items {
		forwardRuleLabels := forwardRule.Labels
		if forwardRuleLabels[tagName] == tagValue {
			_, err = computeService.ForwardingRules.Delete(projectID, region, forwardRule.Name).Do()
			if err != nil {
				return fmt.Errorf("error while deleting forwarding rule: %v", err)
			}
		}
	}

	addresses, err := computeService.Addresses.List(projectID, region).Filter(addressFilter).Do()
	if err != nil {
		return fmt.Errorf("error while listing addresses: %v", err)
	}

	for _, address := range addresses.Items {
		_, errDelete := computeService.Addresses.Delete(projectID, region, address.Name).Do()
		if errDelete != nil {
			return errDelete
		}
	}
	return nil
}
