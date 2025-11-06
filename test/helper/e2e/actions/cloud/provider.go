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
	"path"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
)

const (
	GCPRegion         = "europe-west1"
	AWSRegion         = "eu-west-2"
	AzureRegion       = "northeurope"
	ResourceGroupName = "svet-test"
	Subnet1Name       = "atlas-operator-e2e-test-subnet1"
	Subnet2Name       = "atlas-operator-e2e-test-subnet2"
	Subnet1CIDR       = "10.0.0.0/25"
	Subnet2CIDR       = "10.0.0.128/25"
	vpcName           = "atlas-operator-e2e-test-vpc"
	vpcCIDR           = "10.0.0.0/24"
)

type Provider interface {
	GetAWSAccountID(ctx context.Context) string
	SetupNetwork(ctx context.Context, providerName provider.ProviderName, configs ProviderConfig) string
	SetupPrivateEndpoint(ctx context.Context, request PrivateEndpointRequest) *PrivateEndpointDetails
	ValidatePrivateEndpointStatus(ctx context.Context, providerName provider.ProviderName, endpoint, region string, gcpNumAttachments int)
	SetupNetworkPeering(ctx context.Context, providerName provider.ProviderName, peerID, peerVPC string)
}

type ProviderAction struct {
	t core.GinkgoTInterface

	awsConfig   *AWSConfig
	gcpConfig   *GCPConfig
	azureConfig *AzureConfig

	awsProvider   *AwsAction
	gcpProvider   *GCPAction
	azureProvider *AzureAction
}

type CloudConfig interface {
	AWSConfig | AzureConfig | GCPConfig
}

type AWSConfig struct {
	Region        string
	VPC           string
	CIDR          string
	Subnets       map[string]string
	EnableCleanup bool
}

type GCPConfig struct {
	Region        string
	VPC           string
	Subnets       map[string]string
	EnableCleanup bool
}

type AzureConfig struct {
	Region        string
	VPC           string
	CIDR          string
	Subnets       map[string]string
	EnableCleanup bool
}

type ProviderConfig func(action *ProviderAction)

func WithAWSConfig(config *AWSConfig) ProviderConfig {
	return func(action *ProviderAction) {
		if config.Region != "" {
			action.awsConfig.Region = config.Region
		}

		if config.VPC != "" {
			action.awsConfig.VPC = config.VPC
		}

		if config.CIDR != "" {
			action.awsConfig.CIDR = config.CIDR
		}

		if len(config.Subnets) > 0 {
			action.awsConfig.Subnets = config.Subnets
		}

		action.awsConfig.EnableCleanup = config.EnableCleanup
	}
}

func WithGCPConfig(config *GCPConfig) ProviderConfig {
	return func(action *ProviderAction) {
		if config.Region != "" {
			action.gcpConfig.Region = config.Region
		}

		if config.VPC != "" {
			action.gcpConfig.VPC = config.VPC
		}

		if len(config.Subnets) > 0 {
			action.gcpConfig.Subnets = config.Subnets
		}

		action.gcpConfig.EnableCleanup = config.EnableCleanup
	}
}

func WithAzureConfig(config *AzureConfig) ProviderConfig {
	return func(action *ProviderAction) {
		if config.Region != "" {
			action.azureConfig.Region = config.Region
		}

		if config.VPC != "" {
			action.azureConfig.VPC = config.VPC
		}

		if config.CIDR != "" {
			action.azureConfig.CIDR = config.CIDR
		}

		if len(config.Subnets) > 0 {
			action.azureConfig.Subnets = config.Subnets
		}

		action.azureConfig.EnableCleanup = config.EnableCleanup
	}
}

func (a *ProviderAction) GetAWSAccountID(ctx context.Context) string {
	ID, err := a.awsProvider.GetAccountID(ctx)
	Expect(err).To(BeNil())

	return ID
}

func (a *ProviderAction) SetupNetwork(ctx context.Context, providerName provider.ProviderName, config ProviderConfig) string {
	a.t.Helper()

	if config != nil {
		config(a)
	}

	switch providerName {
	case provider.ProviderAWS:
		id, err := a.awsProvider.InitNetwork(ctx, a.awsConfig.VPC, a.awsConfig.CIDR, a.awsConfig.Region, a.awsConfig.Subnets, a.awsConfig.EnableCleanup)
		Expect(err).To(BeNil())
		return id
	case provider.ProviderGCP:
		id, err := a.gcpProvider.InitNetwork(a.gcpConfig.VPC, a.gcpConfig.Region, a.gcpConfig.Subnets, a.gcpConfig.EnableCleanup)
		Expect(err).To(BeNil())
		return id
	case provider.ProviderAzure:
		id, err := a.azureProvider.InitNetwork(a.azureConfig.VPC, a.azureConfig.CIDR, a.azureConfig.Region, a.azureConfig.Subnets, a.azureConfig.EnableCleanup)
		Expect(err).To(BeNil())
		return id
	}

	return ""
}

type PrivateEndpointRequest interface {
	isPrivateEndpointRequest()
}

type AWSPrivateEndpointRequest struct {
	ID          string
	Region      string
	ServiceName string
}

func (r *AWSPrivateEndpointRequest) isPrivateEndpointRequest() {}

type AzurePrivateEndpointRequest struct {
	ID                string
	Region            string
	ServiceResourceID string
	SubnetName        string
}

func (r *AzurePrivateEndpointRequest) isPrivateEndpointRequest() {}

type GCPPrivateEndpointRequest struct {
	ID         string
	Region     string
	Targets    []string
	SubnetName string
}

func (r *GCPPrivateEndpointRequest) isPrivateEndpointRequest() {}

type PrivateEndpointDetails struct {
	ProviderName      provider.ProviderName
	Region            string
	ID                string
	IP                string
	GCPProjectID      string
	EndpointGroupName string
	Endpoints         []GCPPrivateEndpoint
}

type GCPPrivateEndpoint struct {
	Name string
	IP   string
}

func (a *ProviderAction) SetupPrivateEndpoint(ctx context.Context, request PrivateEndpointRequest) *PrivateEndpointDetails {
	a.t.Helper()

	switch req := request.(type) {
	case *AWSPrivateEndpointRequest:
		ID, err := a.awsProvider.CreatePrivateEndpoint(ctx, req.ServiceName, req.ID, req.Region)
		Expect(err).To(BeNil())

		return &PrivateEndpointDetails{
			ProviderName: provider.ProviderAWS,
			Region:       req.Region,
			ID:           ID,
		}
	case *GCPPrivateEndpointRequest:
		endpoints := make([]GCPPrivateEndpoint, 0, len(req.Targets))
		for index, target := range req.Targets {
			rule, ip, err := a.gcpProvider.CreatePrivateEndpoint(ctx, req.ID, req.Region, req.SubnetName, target, index)
			Expect(err).To(BeNil())

			endpoints = append(
				endpoints,
				GCPPrivateEndpoint{
					Name: rule,
					IP:   ip,
				},
			)
		}

		return &PrivateEndpointDetails{
			ProviderName:      provider.ProviderGCP,
			Region:            req.Region,
			GCPProjectID:      a.gcpProvider.projectID,
			EndpointGroupName: a.gcpProvider.network.VPC,
			Endpoints:         endpoints,
		}
	case *AzurePrivateEndpointRequest:
		pe, err := a.azureProvider.CreatePrivateEndpoint(a.azureConfig.VPC, req.SubnetName, req.ID, req.ServiceResourceID, req.Region)
		Expect(err).To(BeNil())
		Expect(pe).ShouldNot(BeNil())
		Expect(pe.Properties).ShouldNot(BeNil())
		Expect(pe.Properties.NetworkInterfaces).ShouldNot(BeNil())
		Expect(pe.Properties.NetworkInterfaces).ShouldNot(BeEmpty())

		var itf *armnetwork.Interface
		Eventually(func(g Gomega) bool {
			itf, err = a.azureProvider.GetInterface(path.Base(*pe.Properties.NetworkInterfaces[0].ID))
			g.Expect(err).To(BeNil())
			g.Expect(itf).ShouldNot(BeNil())
			g.Expect(itf.Properties).ShouldNot(BeNil())
			g.Expect(itf.Properties.IPConfigurations).ShouldNot(BeNil())
			g.Expect(itf.Properties.IPConfigurations).ShouldNot(BeEmpty())
			g.Expect(itf.Properties.IPConfigurations[0].Properties).ShouldNot(BeNil())

			return true
		}).WithTimeout(5 * time.Minute).WithPolling(15 * time.Second).Should(BeTrue())

		return &PrivateEndpointDetails{
			ProviderName: provider.ProviderAzure,
			Region:       req.Region,
			ID:           *pe.ID,
			IP:           *itf.Properties.IPConfigurations[0].Properties.PrivateIPAddress,
		}
	}

	return nil
}

func (a *ProviderAction) ValidatePrivateEndpointStatus(ctx context.Context, providerName provider.ProviderName, endpoint, region string, gcpNumAttachments int) {
	a.t.Helper()

	Eventually(func(g Gomega) bool {
		switch providerName {
		case provider.ProviderAWS:
			pe, err := a.awsProvider.GetPrivateEndpoint(ctx, endpoint, region)
			g.Expect(err).To(BeNil())

			return pe.State == "available"
		case provider.ProviderGCP:
			res := true
			for i := 0; i < gcpNumAttachments; i++ {
				rule, err := a.gcpProvider.GetForwardingRule(endpoint, region, i)
				g.Expect(err).To(BeNil())

				res = res && (*rule.PscConnectionStatus == "ACCEPTED")
			}

			return res
		case provider.ProviderAzure:
			pe, err := a.azureProvider.GetPrivateEndpoint(endpoint)
			g.Expect(err).To(BeNil())

			return *pe.Properties.ManualPrivateLinkServiceConnections[0].Properties.PrivateLinkServiceConnectionState.Status == "Approved"
		}

		return false
	}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
}

func (a *ProviderAction) SetupNetworkPeering(ctx context.Context, providerName provider.ProviderName, peerID, peerVPC string) {
	switch providerName {
	case provider.ProviderAWS:
		Expect(a.awsProvider.AcceptVpcPeeringConnection(ctx, peerID, a.awsConfig.Region)).To(Succeed())
	case provider.ProviderGCP:
		Expect(a.gcpProvider.CreateNetworkPeering(a.gcpConfig.VPC, peerID, peerVPC)).To(Succeed())
	}
}

func NewProviderAction(t core.GinkgoTInterface, aws *AwsAction, gcp *GCPAction, azure *AzureAction) *ProviderAction {
	return &ProviderAction{
		t: t,

		awsConfig:   getAWSConfigDefaults(),
		gcpConfig:   getGCPConfigDefaults(),
		azureConfig: getAzureConfigDefaults(),

		awsProvider:   aws,
		gcpProvider:   gcp,
		azureProvider: azure,
	}
}

func getAWSConfigDefaults() *AWSConfig {
	return &AWSConfig{
		Region:  AWSRegion,
		VPC:     vpcName,
		CIDR:    vpcCIDR,
		Subnets: map[string]string{Subnet1Name: Subnet1CIDR, Subnet2Name: Subnet2CIDR},
	}
}

func getGCPConfigDefaults() *GCPConfig {
	return &GCPConfig{
		Region:  GCPRegion,
		VPC:     vpcName,
		Subnets: map[string]string{Subnet1Name: Subnet1CIDR, Subnet2Name: Subnet2CIDR},
	}
}

func getAzureConfigDefaults() *AzureConfig {
	return &AzureConfig{
		Region:  AzureRegion,
		VPC:     vpcName,
		CIDR:    vpcCIDR,
		Subnets: map[string]string{Subnet1Name: Subnet1CIDR, Subnet2Name: Subnet2CIDR},
	}
}

func GenerateCloudConfig[T CloudConfig](cloudProvider, region, prefixName string) (*T, error) {
	vpc, subnets, err := generateVPCWithSubnets()
	if err != nil {
		return nil, err
	}

	uniqueID := strings.ToLower(uuid.New().String()[0:6])

	switch cloudProvider {
	case "AWS":
		return any(&AWSConfig{
			Region: region,
			VPC:    fmt.Sprintf("%s-aws-vpc-%s", prefixName, uniqueID),
			CIDR:   vpc,
			Subnets: map[string]string{
				fmt.Sprintf("%s-aws-sn1-%s", prefixName, uniqueID): subnets[0],
				fmt.Sprintf("%s-aws-sn2-%s", prefixName, uniqueID): subnets[1],
			},
			EnableCleanup: true,
		}).(*T), nil
	case "AZURE":
		return any(&AzureConfig{
			Region: region,
			VPC:    fmt.Sprintf("%s-azure-vpc-%s", prefixName, uniqueID),
			CIDR:   vpc,
			Subnets: map[string]string{
				fmt.Sprintf("%s-azure-sn1-%s", prefixName, uniqueID): subnets[0],
				fmt.Sprintf("%s-azure-sn2-%s", prefixName, uniqueID): subnets[1],
			},
			EnableCleanup: true,
		}).(*T), nil
	case "GCP":
		return any(&GCPConfig{
			Region: region,
			VPC:    fmt.Sprintf("%s-gcp-vpc-%s", prefixName, uniqueID),
			Subnets: map[string]string{
				fmt.Sprintf("%s-gcp-sn1-%s", prefixName, uniqueID): subnets[0],
				fmt.Sprintf("%s-gcp-sn2-%s", prefixName, uniqueID): subnets[1],
			},
			EnableCleanup: true,
		}).(*T), nil
	}

	return nil, errors.New("unsupported provider, valid options are: AWS, Azure, GCP")
}

func generateVPCWithSubnets() (string, []string, error) {
	privateRanges := []struct {
		base       string
		subnetMask int
	}{
		{"10.0.0.0", 24},
		{"172.16.0.0", 24},
		{"192.168.0.0", 24},
	}

	// Pick a random range
	r, err := rand.Int(rand.Reader, big.NewInt(int64(len(privateRanges))))
	if err != nil {
		return "", nil, err
	}
	privateRange := privateRanges[r.Int64()]

	_, network, err := net.ParseCIDR(fmt.Sprintf("%s/%d", privateRange.base, privateRange.subnetMask))
	if err != nil {
		return "", nil, err
	}

	// Generate 2 subnets
	ip := network.IP.To4()
	r, err = rand.Int(rand.Reader, big.NewInt(int64(255)))
	if err != nil {
		return "", nil, err
	}
	ip[2] = byte(r.Int64())

	vpcCIDR := fmt.Sprintf("%s/%d", ip.String(), privateRange.subnetMask)
	subnet1CIDR := fmt.Sprintf("%s/%d", ip.String(), privateRange.subnetMask+1)
	ip[3] = 128
	subnet2CIDR := fmt.Sprintf("%s/%d", ip.String(), privateRange.subnetMask+1)

	return vpcCIDR, []string{subnet1CIDR, subnet2CIDR}, nil
}

func GetAtlasRegionByProvider(cloudProvider string) (string, error) {
	regionMap := map[string][]string{
		"AWS": {
			"US_WEST_2",      // North America - Oregon
			"CA_CENTRAL_1",   // North America - Canada Central
			"SA_EAST_1",      // South America - São Paulo
			"EU_NORTH_1",     // Europe - Stockholm
			"EU_WEST_3",      // Europe - Paris
			"ME_SOUTH_1",     // Middle East - Bahrain
			"AP_SOUTH_2",     // Asia - Hyderabad
			"AP_NORTHEAST_2", // Asia - Seoul
			"AP_SOUTHEAST_2", // Oceania - Sydney
		},
		"AZURE": {
			"US_WEST_3",           // North America - West US
			"US_EAST_2",           // North America - East US
			"BRAZIL_SOUTHEAST",    // South America - Brazil Southeast
			"EUROPE_NORTH",        // Europe - North Europe (Ireland)
			"NORWAY_EAST",         // Europe - Norway East
			"FRANCE_SOUTH",        // Europe - France South
			"UAE_CENTRAL",         // Middle East - UAE Central
			"KOREA_CENTRAL",       // Asia - Korea Central
			"INDIA_CENTRAL",       // Asia - India Central
			"AUSTRALIA_CENTRAL_2", // Oceania - Australia Central
			"SOUTH_AFRICA_WEST",   // Africa - South Africa West
		},
		"GCP": {
			"US_WEST_3",             // North America - Salt Lake City
			"US_EAST_5",             // North America - Columbus
			"SOUTH_AMERICA_EAST_1",  // South America - São Paulo
			"EUROPE_WEST_3",         // Europe - Frankfurt
			"EUROPE_NORTH_1",        // Europe - Finland
			"EUROPE_WEST_6",         // Europe - Zurich
			"ASIA_EAST_2",           // Asia - Hong Kong
			"ASIA_NORTHEAST_2",      // Asia - Osaka
			"AUSTRALIA_SOUTHEAST_2", // Oceania - Melbourne
		},
	}

	// Validate the input provider
	regions := regionMap[cloudProvider]
	r, err := rand.Int(rand.Reader, big.NewInt(int64(len(regions))))
	if err != nil {
		return "", err
	}

	return regions[r.Int64()], nil
}

func MapCloudProviderRegion(cloudProvider, atlasRegion string) string {
	regionMap := map[string]map[string]string{
		"AWS": {
			"US_WEST_2":      "us-west-2",      // North America - Oregon
			"CA_CENTRAL_1":   "ca-central-1",   // North America - Canada Central
			"SA_EAST_1":      "sa-east-1",      // South America - São Paulo
			"EU_NORTH_1":     "eu-north-1",     // Europe - Stockholm
			"EU_WEST_3":      "eu-west-3",      // Europe - Paris
			"ME_SOUTH_1":     "me-south-1",     // Middle East - Bahrain
			"AP_SOUTH_2":     "ap-south-2",     // Asia - Hyderabad
			"AP_NORTHEAST_2": "ap-northeast-2", // Asia - Seoul
			"AP_SOUTHEAST_2": "ap-southeast-2", // Oceania - Sydney
		},
		"AZURE": {
			"US_WEST_3":           "westus3",           // North America - West US
			"US_EAST_2":           "eastus2",           // North America - East US
			"BRAZIL_SOUTHEAST":    "brazilsoutheast",   // South America - Brazil Southeast
			"EUROPE_NORTH":        "northeurope",       // Europe - North Europe (Ireland)
			"NORWAY_EAST":         "norwayeast",        // Europe - Norway East
			"FRANCE_SOUTH":        "francesouth",       // Europe - France South
			"UAE_CENTRAL":         "uaecentral",        // Middle East - UAE Central
			"KOREA_CENTRAL":       "koreacentral",      // Asia - Korea Central
			"INDIA_CENTRAL":       "centralindia",      // Asia - India Central
			"AUSTRALIA_CENTRAL_2": "australiacentral2", // Oceania - Australia Central
			"SOUTH_AFRICA_WEST":   "southafricawest",   // Africa - South Africa West
		},
		"GCP": {
			"US_WEST_3":             "us-west3",             // North America - Salt Lake City
			"US_EAST_5":             "us-east5",             // North America - Columbus
			"SOUTH_AMERICA_EAST_1":  "southamerica-east1",   // South America - São Paulo
			"EUROPE_WEST_3":         "europe-west3",         // Europe - Frankfurt
			"EUROPE_NORTH_1":        "europe-north1",        // Europe - Finland
			"EUROPE_WEST_6":         "europe-west6",         // Europe - Zurich
			"ASIA_EAST_2":           "asia-east2",           // Asia - Hong Kong
			"ASIA_NORTHEAST_2":      "asia-northeast2",      // Asia - Osaka
			"AUSTRALIA_SOUTHEAST_2": "australia-southeast2", // Oceania - Melbourne
		},
	}

	if _, ok := regionMap[cloudProvider]; !ok {
		return ""
	}

	if _, ok := regionMap[cloudProvider][atlasRegion]; !ok {
		return ""
	}

	return regionMap[cloudProvider][atlasRegion]
}
