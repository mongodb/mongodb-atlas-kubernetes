package cloud

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/api/googleapi"

	"google.golang.org/api/iterator"

	"github.com/onsi/ginkgo/v2/dsl/core"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
)

type GCPAction struct {
	t         core.GinkgoTInterface
	projectID string
	network   *gcpNetwork

	networkClient     *compute.NetworksClient
	subnetClient      *compute.SubnetworksClient
	addressClient     *compute.AddressesClient
	forwardRuleClient *compute.ForwardingRulesClient
}

type gcpNetwork struct {
	VPC     string
	Subnets map[string]string
}

const (
	// TODO get from GCP
	GoogleProjectID     = "atlasoperator"             // Google Cloud Project ID
	GoogleVPC           = "atlas-operator-test"       // VPC Name
	GoogleSubnetName    = "atlas-operator-subnet-leo" // Subnet Name
	googleConnectPrefix = "ao"                        // Private Service Connect Endpoint Prefix

	gcpSubnetIPMask = "10.0.0.%d"
)

func (a *GCPAction) InitNetwork(vpcName, region string, subnets map[string]string, cleanup bool) (string, error) {
	a.t.Helper()
	ctx := context.Background()

	vpc, err := a.findVPC(ctx, vpcName)
	if err != nil {
		return "", err
	}

	if vpc == nil {
		err = a.createVPC(ctx, vpcName)
		if err != nil {
			return "", err
		}
	}

	if cleanup {
		a.t.Cleanup(func() {
			err = a.deleteVPC(ctx, vpcName)
			if err != nil {
				a.t.Error(err)
			}
		})
	}

	existingSubnets, err := a.getSubnets(ctx, region)
	if err != nil {
		return "", err
	}

	for name, ipRange := range subnets {
		if _, ok := existingSubnets[name]; !ok {
			if err = a.createSubnet(ctx, vpcName, name, ipRange, region); err != nil {
				return "", err
			}

			if cleanup {
				a.t.Cleanup(func() {
					err = a.deleteSubnet(ctx, name, region)
					if err != nil {
						a.t.Error(err)
					}
				})
			}
		}
	}

	a.network = &gcpNetwork{
		VPC:     vpcName,
		Subnets: subnets,
	}

	return vpcName, nil
}
func (a *GCPAction) CreatePrivateEndpoint(name, region, subnet, target string, index int) (string, string, error) {
	a.t.Helper()
	ctx := context.Background()

	address := fmt.Sprintf("%s-%s-ip-%d", googleConnectPrefix, name, index)
	rule := fmt.Sprintf("%s-%s-fr-%d", googleConnectPrefix, name, index)

	ipAddress, err := a.createVirtualAddress(ctx, address, subnet, region)
	if err != nil {
		return "", "", err
	}

	a.t.Cleanup(func() {
		err = a.deleteVirtualAddress(ctx, address, region)
		if err != nil {
			a.t.Error(err)
		}
	})

	err = a.createForwardRule(ctx, rule, address, region, target)
	if err != nil {
		return "", "", err
	}

	a.t.Cleanup(func() {
		err = a.deleteForwardRule(ctx, rule, region)
		if err != nil {
			a.t.Error(err)
		}
	})

	return rule, ipAddress, err
}

func (a *GCPAction) GetForwardingRule(name, region string, suffixIndex int) (*computepb.ForwardingRule, error) {
	a.t.Helper()

	ruleRequest := &computepb.GetForwardingRuleRequest{
		Project:        a.projectID,
		ForwardingRule: fmt.Sprintf("%s-%s-fr-%d", googleConnectPrefix, name, suffixIndex),
		Region:         region,
	}

	rule, err := a.forwardRuleClient.Get(context.Background(), ruleRequest)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (a *GCPAction) CreateNetworkPeering(vpcName, peerProjectID, peerVPCName string) error {
	a.t.Helper()
	ctx := context.Background()

	peerName := "atlas-networking-peering"
	peerRequest := &computepb.AddPeeringNetworkRequest{
		Project: a.projectID,
		Network: vpcName,
		NetworksAddPeeringRequestResource: &computepb.NetworksAddPeeringRequest{
			NetworkPeering: &computepb.NetworkPeering{
				Name:                 toptr.MakePtr(peerName),
				Network:              toptr.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", peerProjectID, peerVPCName)),
				ExchangeSubnetRoutes: toptr.MakePtr(true),
			},
		},
	}

	op, err := a.networkClient.AddPeering(ctx, peerRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	a.t.Cleanup(func() {
		err = a.deleteVPCPeering(ctx, vpcName, peerName)
		if err != nil {
			a.t.Error(err)
		}
	})

	return nil
}

func (a *GCPAction) findVPC(ctx context.Context, vpcName string) (*computepb.Network, error) {
	a.t.Helper()

	vpcRequest := &computepb.GetNetworkRequest{
		Project: a.projectID,
		Network: vpcName,
	}

	vpc, err := a.networkClient.Get(ctx, vpcRequest)
	if err != nil {
		var respErr *googleapi.Error
		if ok := errors.As(err, &respErr); ok && respErr.Code == 404 {
			return nil, nil
		}
		return nil, err
	}

	return vpc, nil
}

func (a *GCPAction) createVPC(ctx context.Context, vpcName string) error {
	a.t.Helper()

	vpcRequest := &computepb.InsertNetworkRequest{
		Project: a.projectID,
		NetworkResource: &computepb.Network{
			Name:                  toptr.MakePtr(vpcName),
			Description:           toptr.MakePtr("Atlas Kubernetes Operator E2E Tests VPC"),
			AutoCreateSubnetworks: toptr.MakePtr(false),
		},
	}

	op, err := a.networkClient.Insert(ctx, vpcRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) deleteVPC(ctx context.Context, vpcName string) error {
	a.t.Helper()

	vpcRequest := &computepb.DeleteNetworkRequest{
		Project: a.projectID,
		Network: vpcName,
	}

	op, err := a.networkClient.Delete(ctx, vpcRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) getSubnets(ctx context.Context, region string) (map[string]string, error) {
	a.t.Helper()

	subnetRequest := &computepb.ListSubnetworksRequest{
		Project: a.projectID,
		Region:  region,
	}

	subnets := map[string]string{}
	list := a.subnetClient.List(ctx, subnetRequest)
	for {
		subnet, err := list.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		subnets[*subnet.Name] = *subnet.IpCidrRange
	}

	return subnets, nil
}

func (a *GCPAction) createSubnet(ctx context.Context, vpcName, subnetName, ipRange, region string) error {
	a.t.Helper()

	subnetRequest := &computepb.InsertSubnetworkRequest{
		Project: a.projectID,
		Region:  region,
		SubnetworkResource: &computepb.Subnetwork{
			Name:        toptr.MakePtr(subnetName),
			Network:     toptr.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", a.projectID, vpcName)),
			IpCidrRange: toptr.MakePtr(ipRange),
		},
	}

	op, err := a.subnetClient.Insert(ctx, subnetRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) deleteSubnet(ctx context.Context, subnetName, region string) error {
	a.t.Helper()

	subnetRequest := &computepb.DeleteSubnetworkRequest{
		Subnetwork: subnetName,
		Project:    a.projectID,
		Region:     region,
	}

	op, err := a.subnetClient.Delete(ctx, subnetRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) createVirtualAddress(ctx context.Context, name, subnet, region string) (string, error) {
	a.t.Helper()

	ip := fmt.Sprintf(gcpSubnetIPMask, rand.IntnRange(2, 200))
	addressRequest := &computepb.InsertAddressRequest{
		Project: a.projectID,
		Region:  region,
		AddressResource: &computepb.Address{
			Name:        toptr.MakePtr(name),
			Description: toptr.MakePtr(name),
			Address:     toptr.MakePtr(ip),
			AddressType: toptr.MakePtr("INTERNAL"),
			Subnetwork: toptr.MakePtr(
				fmt.Sprintf(
					"https://www.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks/%s",
					a.projectID,
					region,
					subnet,
				),
			),
		},
	}

	op, err := a.addressClient.Insert(ctx, addressRequest)
	if err != nil {
		return "", err
	}

	err = op.Wait(ctx)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func (a *GCPAction) deleteVirtualAddress(ctx context.Context, name, region string) error {
	a.t.Helper()

	addressRequest := &computepb.DeleteAddressRequest{
		Address: name,
		Project: a.projectID,
		Region:  region,
	}

	op, err := a.addressClient.Delete(ctx, addressRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) createForwardRule(ctx context.Context, rule, address, region, target string) error {
	a.t.Helper()

	ruleRequest := &computepb.InsertForwardingRuleRequest{
		Project: a.projectID,
		Region:  region,
		ForwardingRuleResource: &computepb.ForwardingRule{
			Name:      toptr.MakePtr(rule),
			IPAddress: toptr.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/addresses/%s", a.projectID, region, address)),
			Network:   toptr.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", a.projectID, a.network.VPC)),
			Target:    toptr.MakePtr(target),
		},
	}

	op, err := a.forwardRuleClient.Insert(ctx, ruleRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) deleteForwardRule(ctx context.Context, rule, region string) error {
	a.t.Helper()

	addressRequest := &computepb.DeleteForwardingRuleRequest{
		ForwardingRule: rule,
		Project:        a.projectID,
		Region:         region,
	}

	op, err := a.forwardRuleClient.Delete(ctx, addressRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *GCPAction) deleteVPCPeering(ctx context.Context, vpcName, peerName string) error {
	a.t.Helper()

	peerRequest := &computepb.RemovePeeringNetworkRequest{
		Project: a.projectID,
		Network: vpcName,
		NetworksRemovePeeringRequestResource: &computepb.NetworksRemovePeeringRequest{
			Name: toptr.MakePtr(peerName),
		},
	}

	op, err := a.networkClient.RemovePeering(ctx, peerRequest)
	if err != nil {
		return err
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGCPAction(t core.GinkgoTInterface, projectID string) (*GCPAction, error) {
	t.Helper()

	ctx := context.Background()

	networkClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	subnetClient, err := compute.NewSubnetworksRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	addressClient, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	forwardRuleClient, err := compute.NewForwardingRulesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCPAction{
		t:         t,
		projectID: projectID,

		networkClient:     networkClient,
		subnetClient:      subnetClient,
		addressClient:     addressClient,
		forwardRuleClient: forwardRuleClient,
	}, nil
}
