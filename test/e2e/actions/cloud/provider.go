package cloud

import (
	"errors"

	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

const (
	gcpRegion   = "europe-west1"
	awsRegion   = "eu-west-2"
	azureRegion = "northeurope"
	vpcName     = "atlas-operator-e2e-test-vpc"
	vpcCIDR     = "10.0.0.0/24"
	subnetName  = "atlas-operator-e2e-test-subnet"
	subnetCIDR  = "10.0.0.0/24"
)

type Provider interface {
	SetupNetwork(providerName provider.ProviderName, configs ...ProviderConfig) error
	SetupPrivateEndpoint(request PrivateEndpointRequest) (*PrivateEndpointDetails, error)
	IsPrivateEndpointAvailable(providerName provider.ProviderName, endpoint, region string) (bool, error)
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
	Region  string
	VPC     string
	CIDR    string
	Subnets []string
}

type GCPConfig struct {
	Region  string
	VPC     string
	Subnets map[string]string
}

type AzureConfig struct {
	Region  string
	VPC     string
	CIDR    string
	Subnets map[string]string
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
	}
}

func (a *ProviderAction) SetupNetwork(providerName provider.ProviderName, configs ...ProviderConfig) error {
	a.t.Helper()

	for _, config := range configs {
		config(a)
	}

	switch providerName {
	case provider.ProviderAWS:
		err := a.awsProvider.InitNetwork(a.awsConfig.VPC, a.awsConfig.CIDR, a.awsConfig.Region, a.awsConfig.Subnets)
		if err != nil {
			return err
		}
	case provider.ProviderGCP:
		err := a.gcpProvider.InitNetwork(a.gcpConfig.VPC, a.gcpConfig.Region, a.gcpConfig.Subnets)
		if err != nil {
			return err
		}
	case provider.ProviderAzure:
		err := a.azureProvider.InitNetwork(a.azureConfig.VPC, a.azureConfig.CIDR, a.azureConfig.Region, a.azureConfig.Subnets)
		if err != nil {
			return err
		}
	}

	return nil
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
}

func (r *AzurePrivateEndpointRequest) isPrivateEndpointRequest() {}

type GCPPrivateEndpointRequest struct {
	ID      string
	Region  string
	Targets []string
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

func (a *ProviderAction) SetupPrivateEndpoint(request PrivateEndpointRequest) (*PrivateEndpointDetails, error) {
	a.t.Helper()

	switch req := request.(type) {
	case *AWSPrivateEndpointRequest:
		ID, err := a.awsProvider.CreatePrivateEndpoint(req.ServiceName, req.ID, req.Region)
		if err != nil {
			return nil, err
		}

		return &PrivateEndpointDetails{
			ProviderName: provider.ProviderAWS,
			Region:       req.Region,
			ID:           ID,
		}, nil
	case *GCPPrivateEndpointRequest:
		endpoints := make([]GCPPrivateEndpoint, 0, len(req.Targets))
		for index, target := range req.Targets {
			rule, ip, err := a.gcpProvider.CreatePrivateEndpoint(req.ID, req.Region, subnetName, target, index)
			if err != nil {
				return nil, err
			}

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
		}, nil
	case *AzurePrivateEndpointRequest:
		ID, ip, err := a.azureProvider.CreatePrivateEndpoint(vpcName, subnetName, req.ID, req.ServiceResourceID, req.Region)
		if err != nil {
			return nil, err
		}

		return &PrivateEndpointDetails{
			ProviderName: provider.ProviderAzure,
			Region:       req.Region,
			ID:           ID,
			IP:           ip,
		}, nil
	}

	return nil, nil
}

func (a *ProviderAction) IsPrivateEndpointAvailable(providerName provider.ProviderName, endpoint, region string) (bool, error) {
	a.t.Helper()

	switch providerName {
	case provider.ProviderAWS:
		pe, err := a.awsProvider.GetPrivateEndpoint(endpoint, region)
		if err != nil {
			return false, err
		}

		return *pe.State == "available", nil
	case provider.ProviderGCP:
		rule, err := a.gcpProvider.GetForwardingRule(endpoint, region)
		if err != nil {
			return false, err
		}

		return *rule.PscConnectionStatus == "ACCEPTED", nil
	case provider.ProviderAzure:
		pe, err := a.azureProvider.GetPrivateEndpoint(endpoint)
		if err != nil {
			return false, err
		}

		return *pe.Properties.ManualPrivateLinkServiceConnections[0].Properties.PrivateLinkServiceConnectionState.Status == "Approved", nil
	}

	return false, errors.New("invalid data")
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
		Region:  awsRegion,
		VPC:     vpcName,
		CIDR:    vpcCIDR,
		Subnets: []string{subnetCIDR},
	}
}

func getGCPConfigDefaults() *GCPConfig {
	return &GCPConfig{
		Region:  gcpRegion,
		VPC:     vpcName,
		Subnets: map[string]string{subnetName: subnetCIDR},
	}
}

func getAzureConfigDefaults() *AzureConfig {
	return &AzureConfig{
		Region:  azureRegion,
		VPC:     vpcName,
		CIDR:    vpcCIDR,
		Subnets: map[string]string{subnetName: subnetCIDR},
	}
}
