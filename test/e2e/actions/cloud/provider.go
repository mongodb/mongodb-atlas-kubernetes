package cloud

import (
	"errors"
	"fmt"

	"github.com/onsi/ginkgo/v2/dsl/core"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

const (
	vpcName    = "atlas-operator-e2e-test-vpc"
	vpcCIDR    = "10.0.0.0/24"
	subnetName = "atlas-operator-e2e-test-subnet"
	subnetCIDR = "10.0.0.0/24"
)

type Provider interface {
	SetupNetwork(configs ...ProviderConfig) error
	SetupPrivateEndpoint(request PrivateEndpointRequest) (*PrivateEndpointDetails, error)
	IsPrivateEndpointAvailable(providerName provider.ProviderName, endpoint, region string) (bool, error)
}

type ProviderAction struct {
	t core.GinkgoTInterface

	awsRegion   string
	gcpRegion   string
	azureRegion string

	awsProvider   *AwsAction
	gcpProvider   *GCPAction
	azureProvider *AzureAction
}

type ProviderConfig func(action *ProviderAction)

func WithAWSConfig(region string) ProviderConfig {
	return func(action *ProviderAction) {
		action.awsRegion = region
	}
}

func WithGCPConfig(region string) ProviderConfig {
	return func(action *ProviderAction) {
		action.gcpRegion = region
	}
}

func WithAzureConfig(region string) ProviderConfig {
	return func(action *ProviderAction) {
		action.azureRegion = region
	}
}

func (a *ProviderAction) SetupNetwork(configs ...ProviderConfig) error {
	a.t.Helper()

	for _, config := range configs {
		config(a)
	}

	providers := []provider.ProviderName{provider.ProviderAWS, provider.ProviderGCP, provider.ProviderAzure}
	for _, p := range providers {
		switch p {
		case provider.ProviderAWS:
			err := a.awsProvider.InitNetwork(vpcName, vpcCIDR, a.awsRegion, []string{subnetCIDR})
			if err != nil {
				return err
			}
		case provider.ProviderGCP:
			err := a.gcpProvider.InitNetwork(vpcName, a.gcpRegion, map[string]string{subnetName: subnetCIDR})
			if err != nil {
				return err
			}
		case provider.ProviderAzure:
			err := a.azureProvider.InitNetwork(vpcName, vpcCIDR, a.azureRegion, map[string]string{subnetName: subnetCIDR})
			if err != nil {
				return err
			}
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
			rule, ip, err := a.gcpProvider.CreatePrivateEndpoint(fmt.Sprintf("%s-%d", req.ID, index), req.Region, subnetName, target)
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
		t:             t,
		awsProvider:   aws,
		gcpProvider:   gcp,
		azureProvider: azure,
	}
}
