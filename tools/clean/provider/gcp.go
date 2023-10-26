package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

type GCP struct {
	projectID string

	networkClient       *compute.NetworksClient
	forwardRuleClient   *compute.ForwardingRulesClient
	addressClient       *compute.AddressesClient
	keyManagementClient *kms.KeyManagementClient
}

func (gcp *GCP) DeleteVpc(ctx context.Context, vpcName string) error {
	vpcRequest := &computepb.DeleteNetworkRequest{
		Project: gcp.projectID,
		Network: vpcName,
	}

	op, err := gcp.networkClient.Delete(ctx, vpcRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (gcp *GCP) DeletePrivateEndpoint(ctx context.Context, groupName, sAttachment string) error {
	region := extractRegion(sAttachment)

	err := gcp.deleteForwardRules(ctx, region, groupName)
	if err != nil {
		return err
	}

	err = gcp.deleteIPAddresses(ctx, region, groupName)
	if err != nil {
		return err
	}

	return nil
}

func (gcp *GCP) DeleteCryptoKey(ctx context.Context, keyName string) error {
	_, err := gcp.keyManagementClient.GetCryptoKeyVersion(ctx, &kmspb.GetCryptoKeyVersionRequest{
		Name: keyName,
	})
	if err != nil {
		var respErr *googleapi.Error
		if ok := errors.As(err, &respErr); ok && respErr.Code == 404 {
			return nil
		}
		return err
	}

	req := &kmspb.DestroyCryptoKeyVersionRequest{
		Name: keyName,
	}

	_, err = gcp.keyManagementClient.DestroyCryptoKeyVersion(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (gcp *GCP) deleteForwardRules(ctx context.Context, region, groupName string) error {
	filter := fmt.Sprintf("name=%s*", groupName)

	request := &computepb.ListForwardingRulesRequest{
		Filter:  &filter,
		Project: gcp.projectID,
		Region:  region,
	}

	it := gcp.forwardRuleClient.List(ctx, request)
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		frRequest := &computepb.DeleteForwardingRuleRequest{
			ForwardingRule: resp.GetName(),
			Project:        gcp.projectID,
			Region:         resp.GetRegion(),
		}

		op, err := gcp.forwardRuleClient.Delete(ctx, frRequest)
		if err != nil {
			return err
		}

		err = op.Wait(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gcp *GCP) deleteIPAddresses(ctx context.Context, region, groupName string) error {
	filter := fmt.Sprintf("name=%s*", groupName)

	request := &computepb.ListAddressesRequest{
		Filter:  &filter,
		Project: gcp.projectID,
		Region:  region,
	}

	it := gcp.addressClient.List(ctx, request)
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		aRequest := &computepb.DeleteAddressRequest{
			Address: resp.GetName(),
			Project: gcp.projectID,
			Region:  resp.GetRegion(),
		}

		op, err := gcp.addressClient.Delete(ctx, aRequest)
		if err != nil {
			return err
		}

		err = op.Wait(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGCPCleaner(ctx context.Context) (*GCP, error) {
	_, defined := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if !defined {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS must be set")
	}

	projectID, defined := os.LookupEnv("GOOGLE_PROJECT_ID")
	if !defined {
		return nil, fmt.Errorf("GOOGLE_PROJECT_ID must be set")
	}

	networkClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	forwardRuleClient, err := compute.NewForwardingRulesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	addressClient, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	keyManagementClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCP{
		projectID:           projectID,
		networkClient:       networkClient,
		forwardRuleClient:   forwardRuleClient,
		addressClient:       addressClient,
		keyManagementClient: keyManagementClient,
	}, nil
}

func extractRegion(serviceAttachment string) string {
	parts := strings.Split(serviceAttachment, "/")

	return parts[3]
}
