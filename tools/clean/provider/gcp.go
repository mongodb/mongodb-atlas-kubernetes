package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

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
	subnetworksClient   *compute.SubnetworksClient
	forwardRuleClient   *compute.ForwardingRulesClient
	addressClient       *compute.AddressesClient
	keyManagementClient *kms.KeyManagementClient
}

func (gcp *GCP) DeleteVpc(ctx context.Context, vpcName string) error {
	vpcGetRequest := &computepb.GetNetworkRequest{
		Project: gcp.projectID,
		Network: vpcName,
	}
	net, err := gcp.networkClient.Get(ctx, vpcGetRequest)
	if err != nil {
		return fmt.Errorf("failed to get VPC %q: %v", vpcName, err)
	}
	for _, subnetURL := range net.Subnetworks {
		subnet, region := decodeSubnetURL(subnetURL)
		if subnet == "" {
			return fmt.Errorf("failed to decode subnet URL %q", subnetURL)
		}
		subnetDeleteRequest := &computepb.DeleteSubnetworkRequest{
			Project:    gcp.projectID,
			Subnetwork: subnet,
			Region:     region,
		}
		op, err := gcp.subnetworksClient.Delete(ctx, subnetDeleteRequest)
		if err := waitOrFailOp(ctx, op, err); err != nil {
			return fmt.Errorf("failed to delete subnet %q: %v", subnet, err)
		}
	}
	vpcRequest := &computepb.DeleteNetworkRequest{
		Project: gcp.projectID,
		Network: vpcName,
	}

	op, err := gcp.networkClient.Delete(ctx, vpcRequest)
	if err := waitOrFailOp(ctx, op, err); err != nil {
		return fmt.Errorf("failed to delete VPC %q: %v", vpcName, err)
	}

	return nil
}

func decodeSubnetURL(subnetURL string) (string, string) {
	parts := strings.Split(subnetURL, "/")
	if len(parts) < 11 {
		return "", ""
	}
	region := parts[8]
	subnet := parts[10]
	return subnet, region
}

func waitOrFailOp(ctx context.Context, op *compute.Operation, err error) error {
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

func (gcp *GCP) DeleteOrphanPrivateEndpoints(ctx context.Context, lifetime int, region string, subnet string) ([]string, []string, error) {
	addresses := gcp.addressClient.List(ctx, &computepb.ListAddressesRequest{
		Project: gcp.projectID,
		Region:  region,
	})
	done := []string{}
	skipped := []string{}
	for {
		addr, err := addresses.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed iterating addresses in project %v region %v: %w",
				gcp.projectID, region, err)
		}
		suffix := fmt.Sprintf("subnetworks/%s", subnet)
		if !strings.HasSuffix(addr.GetSubnetwork(), suffix) {
			skipped = append(skipped,
				fmt.Sprintf("Address %s(%s) skipped, not in %s\n", addr.GetName(), addr.GetAddress(), subnet))
			continue
		}
		createdAt, err := asTime(addr.GetCreationTimestamp())
		if err != nil {
			return nil, nil, fmt.Errorf("failed parsing Address creation timestamp %q: %w",
				addr.GetCreationTimestamp(), err)
		}
		if time.Since(createdAt) < time.Duration(lifetime)*time.Hour {
			skipped = append(skipped, fmt.Sprintf("Address %s(%s) skipped once created less than %d hours ago\n",
				addr.GetName(), addr.GetAddress(), lifetime))
		}
		frName, err := expectForwardingRule(addr.GetUsers())
		if err != nil {
			return nil, nil, err
		}

		if frName != "" {
			if err := gcp.deleteForwardingRule(ctx, frName, region); err != nil {
				return nil, nil, fmt.Errorf("failed deleting Forwarding Rule %q in region %q: %w", region, frName, err)
			}
			done = append(done, fmt.Sprintf("Deleted Forwarding Rule %s for %s\n",
				frName, addr.GetAddress()))
		} else {
			skipped = append(skipped,
				fmt.Sprintf("No forwarding rule using Address %s(%s)", addr.GetName(), addr.GetAddress()))
		}
		if err := gcp.deleteIPAddress(ctx, addr.GetName(), region); err != nil {
			return nil, nil, fmt.Errorf("error deleting Address %s(%s) in region %q: %w",
				region, addr.GetName(), addr.GetAddress(), err)
		}
		done = append(done, fmt.Sprintf("Released orphan Address %s(%s)\n", addr.GetName(), addr.GetAddress()))
	}

	return done, skipped, nil
}

func asTime(rfc3339time string) (time.Time, error) {
	return time.Parse(time.RFC3339, rfc3339time)
}

func expectForwardingRule(usersOfEndpointAddress []string) (string, error) {
	if len(usersOfEndpointAddress) == 0 {
		return "", nil
	}
	if len(usersOfEndpointAddress) > 1 {
		return "", fmt.Errorf("expected a single user of an Endpoint Address, but got %v", usersOfEndpointAddress)
	}
	if strings.Contains(usersOfEndpointAddress[0], "/forwardingRules/") {
		return path.Base(usersOfEndpointAddress[0]), nil
	}
	return "", fmt.Errorf("expected a Forwarding Rule user for Endpoint Address but got %s", usersOfEndpointAddress[0])
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
		fwr, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		if err := gcp.deleteForwardingRule(ctx, fwr.GetName(), fwr.GetRegion()); err != nil {
			return os.ErrClosed
		}
	}

	return nil
}

func (gcp *GCP) deleteForwardingRule(ctx context.Context, name, region string) error {
	op, err := gcp.forwardRuleClient.Delete(ctx, &computepb.DeleteForwardingRuleRequest{
		ForwardingRule: name,
		Project:        gcp.projectID,
		Region:         region,
	})
	if err != nil {
		return err
	}

	return op.Wait(ctx)
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
		addr, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		if err := gcp.deleteIPAddress(ctx, addr.GetName(), addr.GetRegion()); err != nil {
			return err
		}
	}

	return nil
}

func (gcp *GCP) deleteIPAddress(ctx context.Context, name, region string) error {
	op, err := gcp.addressClient.Delete(ctx, &computepb.DeleteAddressRequest{
		Address: name,
		Project: gcp.projectID,
		Region:  region,
	})
	if err != nil {
		return fmt.Errorf("failed to delete address %s at %s: %v", name, region, err)
	}

	return op.Wait(ctx)
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

	subnetworksClient, err := compute.NewSubnetworksRESTClient(ctx)
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
		subnetworksClient:   subnetworksClient,
		forwardRuleClient:   forwardRuleClient,
		addressClient:       addressClient,
		keyManagementClient: keyManagementClient,
	}, nil
}

func extractRegion(serviceAttachment string) string {
	parts := strings.Split(serviceAttachment, "/")

	return parts[3]
}
