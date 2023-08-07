//go:build e2e

package cloud

import (
	"path"
	"time"

	. "github.com/onsi/gomega"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
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
	GetAWSAccountID() string
	SetupNetwork(providerName provider.ProviderName, configs ProviderConfig) string
	SetupPrivateEndpoint(request PrivateEndpointRequest) *PrivateEndpointDetails
	ValidatePrivateEndpointStatus(providerName provider.ProviderName, endpoint, region string, gcpNumAttachments int)
	SetupNetworkPeering(providerName provider.ProviderName, peerID, peerVPC string)
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

func (a *ProviderAction) GetAWSAccountID() string {
	ID, err := a.awsProvider.GetAccountID()
	Expect(err).To(BeNil())

	return ID
}

func (a *ProviderAction) SetupNetwork(providerName provider.ProviderName, config ProviderConfig) string {
	a.t.Helper()

	if config != nil {
		config(a)
	}

	switch providerName {
	case provider.ProviderAWS:
		id, err := a.awsProvider.InitNetwork(a.awsConfig.VPC, a.awsConfig.CIDR, a.awsConfig.Region, a.awsConfig.Subnets, a.awsConfig.EnableCleanup)
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

func (a *ProviderAction) SetupPrivateEndpoint(request PrivateEndpointRequest) *PrivateEndpointDetails {
	a.t.Helper()

	switch req := request.(type) {
	case *AWSPrivateEndpointRequest:
		ID, err := a.awsProvider.CreatePrivateEndpoint(req.ServiceName, req.ID, req.Region)
		Expect(err).To(BeNil())

		return &PrivateEndpointDetails{
			ProviderName: provider.ProviderAWS,
			Region:       req.Region,
			ID:           ID,
		}
	case *GCPPrivateEndpointRequest:
		endpoints := make([]GCPPrivateEndpoint, 0, len(req.Targets))
		for index, target := range req.Targets {
			rule, ip, err := a.gcpProvider.CreatePrivateEndpoint(req.ID, req.Region, req.SubnetName, target, index)
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
		pe, err := a.azureProvider.CreatePrivateEndpoint(vpcName, req.SubnetName, req.ID, req.ServiceResourceID, req.Region)
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

func (a *ProviderAction) ValidatePrivateEndpointStatus(providerName provider.ProviderName, endpoint, region string, gcpNumAttachments int) {
	a.t.Helper()

	Eventually(func(g Gomega) bool {
		switch providerName {
		case provider.ProviderAWS:
			pe, err := a.awsProvider.GetPrivateEndpoint(endpoint, region)
			g.Expect(err).To(BeNil())

			return *pe.State == "available"
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

func (a *ProviderAction) SetupNetworkPeering(providerName provider.ProviderName, peerID, peerVPC string) {
	switch providerName {
	case provider.ProviderAWS:
		Expect(a.awsProvider.AcceptVpcPeeringConnection(peerID, a.awsConfig.Region)).To(Succeed())
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
