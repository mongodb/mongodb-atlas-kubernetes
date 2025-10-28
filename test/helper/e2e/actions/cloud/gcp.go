// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloud

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/onsi/ginkgo/v2/dsl/core"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type GCPAction struct {
	t         core.GinkgoTInterface
	projectID string
	network   *gcpNetwork

	networkClient       *compute.NetworksClient
	subnetClient        *compute.SubnetworksClient
	addressClient       *compute.AddressesClient
	forwardRuleClient   *compute.ForwardingRulesClient
	keyManagementClient *kms.KeyManagementClient
}

type gcpNetwork struct {
	VPC     string
	Subnets map[string]string
}

const (
	// TODO get from GCP
	GoogleProjectID     = "atlasoperator" // Google Cloud Project ID
	googleConnectPrefix = "ao"            // Private Service Connect Endpoint Prefix
	googleKeyName       = "projects/atlasoperator/locations/global/keyRings/atlas-operator-test-key-ring/cryptoKeys/encryption-at-rest-test-key"
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

func (a *GCPAction) CreatePrivateEndpoint(ctx context.Context, name, region, subnet, target string, index int) (string, string, error) {
	a.t.Helper()

	address := fmt.Sprintf("%s-%s-ip-%d", googleConnectPrefix, name, index)
	rule := fmt.Sprintf("%s-%s-fr-%d", googleConnectPrefix, name, index)

	ipAddress, err := a.reserveFreeVirtualAddress(ctx, address, subnet, region)
	if err != nil {
		return "", "", err
	}

	a.t.Cleanup(func() {
		err := a.deleteVirtualAddress(ctx, address, region)
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
				Name:                 pointer.MakePtr(peerName),
				Network:              pointer.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", peerProjectID, peerVPCName)),
				ExchangeSubnetRoutes: pointer.MakePtr(true),
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
			Name:                  pointer.MakePtr(vpcName),
			Description:           pointer.MakePtr("Atlas Kubernetes Operator E2E Tests VPC"),
			AutoCreateSubnetworks: pointer.MakePtr(false),
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
			Name:        pointer.MakePtr(subnetName),
			Network:     pointer.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", a.projectID, vpcName)),
			IpCidrRange: pointer.MakePtr(ipRange),
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

func (a *GCPAction) reserveFreeVirtualAddress(ctx context.Context, name, subnet, region string) (string, error) {
	backoff := wait.Backoff{
		Duration: time.Second,
		Factor:   1.5,
		Jitter:   0.7,
		Steps:    7,
		Cap:      time.Minute,
	}
	var err error
	ip := ""
	wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
		ip = a.randomIP(subnet)
		err = a.createVirtualAddress(ctx, ip, name, subnet, region)
		if err != nil {
			if strings.Contains(err.Error(), "IP_IN_USE_BY_ANOTHER_RESOURCE") {
				return false, nil
			}
		}
		return true, err
	})
	return ip, err
}

func (a *GCPAction) randomIP(subnet string) string {
	ip, network, _ := net.ParseCIDR(a.network.Subnets[subnet])

	ipParts := strings.Split(ip.String(), ".")
	if len(ipParts) != 4 {
		panic(fmt.Errorf("failed to parse IPv4 %q into 4 byte parts", ip))
	}
	const maxRandValue = 256
	for {
		randNumberBig, err := rand.Int(rand.Reader, big.NewInt(maxRandValue))
		if err != nil {
			a.t.Error(err)
		}
		randNumber := randNumberBig.String()
		ipParts[3] = randNumber
		genIP := net.ParseIP(strings.Join(ipParts, "."))

		if network.Contains(genIP) {
			return genIP.String()
		}
	}
}

func (a *GCPAction) createVirtualAddress(ctx context.Context, ip, name, subnet, region string) error {
	a.t.Helper()

	addressRequest := &computepb.InsertAddressRequest{
		Project: a.projectID,
		Region:  region,
		AddressResource: &computepb.Address{
			Name:        pointer.MakePtr(name),
			Description: pointer.MakePtr(name),
			Address:     pointer.MakePtr(ip),
			AddressType: pointer.MakePtr("INTERNAL"),
			Subnetwork: pointer.MakePtr(
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
		return err
	}

	if err = op.Wait(ctx); err != nil {
		return err
	}

	return nil
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
			Name:      pointer.MakePtr(rule),
			IPAddress: pointer.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/addresses/%s", a.projectID, region, address)),
			Network:   pointer.MakePtr(fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", a.projectID, a.network.VPC)),
			Target:    pointer.MakePtr(target),
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
			Name: pointer.MakePtr(peerName),
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

func (a *GCPAction) CreateKMS() (string, error) {
	a.t.Helper()

	ctx := context.Background()

	result, err := a.keyManagementClient.CreateCryptoKeyVersion(ctx, &kmspb.CreateCryptoKeyVersionRequest{
		Parent: googleKeyName,
	})
	if err != nil {
		return "", err
	}

	a.t.Cleanup(func() {
		err = a.deleteKMS(ctx, result.Name)
		if err != nil {
			a.t.Error(err)
		}
	})

	ver := strings.Split(result.Name, "/")
	keyVersion := ver[len(ver)-1]

	_, err = a.keyManagementClient.UpdateCryptoKeyPrimaryVersion(ctx, &kmspb.UpdateCryptoKeyPrimaryVersionRequest{
		Name:               googleKeyName,
		CryptoKeyVersionId: keyVersion,
	})
	if err != nil {
		return "", err
	}

	return result.Name, nil
}

func (a *GCPAction) deleteKMS(ctx context.Context, keyName string) error {
	a.t.Helper()

	req := &kmspb.DestroyCryptoKeyVersionRequest{
		Name: keyName,
	}

	_, err := a.keyManagementClient.DestroyCryptoKeyVersion(ctx, req)
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

	keyManagementClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCPAction{
		t:         t,
		projectID: projectID,

		networkClient:       networkClient,
		subnetClient:        subnetClient,
		addressClient:       addressClient,
		forwardRuleClient:   forwardRuleClient,
		keyManagementClient: keyManagementClient,
	}, nil
}
