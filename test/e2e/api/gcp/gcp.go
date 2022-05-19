package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type sessionGCP struct {
	computeService *compute.Service
	gProjectID     string
}

func SessionGCP(gProjectID string) (sessionGCP, error) {
	computeService, err := compute.NewService(context.Background(), option.WithCredentialsFile("my-atlasoperator-ba1b0d70afc5.json")) // TODO
	if err != nil {
		return sessionGCP{}, fmt.Errorf("compute.NewClient: %v", err)
	}
	return sessionGCP{computeService, gProjectID}, nil
}

func (s *sessionGCP) AddIPAdress(region, addressName, subnet string) (string, error) {
	address := &compute.Address{
		AddressType:    "INTERNAL",
		Description:    addressName,
		Name:           addressName,
		Network:        "",
		Region:         region,
		Subnetwork:     s.formSubnetURL(region, subnet),
		ServerResponse: googleapi.ServerResponse{},
	}
	_, err := s.computeService.Addresses.Insert(s.gProjectID, region, address).Context(context.Background()).Do()
	if err != nil {
		return "", fmt.Errorf("computeService.Addresses.Insert: %v", err)
	}
	// TODO add get IP and return
	return "", nil
}

func (s *sessionGCP) DescribeIPStatus(region, addressName string) (string, error) {
	resp, err := s.computeService.Addresses.Get(s.gProjectID, region, addressName).Context(context.Background()).Do()
	if err != nil {
		return "", fmt.Errorf("computeService.Addresses.Get: %v", err)
	}
	return resp.Status, nil
}

func (s *sessionGCP) DeleteIPAdress(region, addressName string) error {
	_, err := s.computeService.Addresses.Delete(s.gProjectID, region, addressName).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("computeService.Addresses.Delete: %v", err)
	}
	return nil
}

func (s *sessionGCP) AddForwardRule(region, ruleName, addressName, network, subnet, target string) ( error) {
	rules := &compute.ForwardingRule{
		IPAddress:                     s.formAddressURL(region, addressName),
		IPProtocol:                    "",
		AllPorts:                      false,
		AllowGlobalAccess:             false,
		BackendService:                "",
		Description:                   "",
		Fingerprint:                   "",
		IpVersion:                     "",
		IsMirroringCollector:          false,
		Kind:                          "",
		LabelFingerprint:              "",
		Labels:                        map[string]string{},
		LoadBalancingScheme:           "",
		MetadataFilters:               []*compute.MetadataFilter{},
		Name:                          ruleName,
		Network:                       s.formNetworkURL(network),
		NetworkTier:                   "",
		PortRange:                     "",
		Ports:                         []string{},
		PscConnectionId:               0,
		PscConnectionStatus:           "",
		Region:                        region,
		SelfLink:                      "",
		ServiceDirectoryRegistrations: []*compute.ForwardingRuleServiceDirectoryRegistration{},
		ServiceLabel:                  "",
		Subnetwork:                    "",
		Target:                        target,
		ServerResponse:                googleapi.ServerResponse{},
	}
	_, err := s.computeService.ForwardingRules.Insert(s.gProjectID, region, rules).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("computeService.ForwardingRules.Insert: %v", err)
	}
	return  nil
}

func (s *sessionGCP) DeleteForwardRule(region, ruleName string) error {
	_, err := s.computeService.ForwardingRules.Delete(s.gProjectID, region, ruleName).Do()
	if err != nil {
		return fmt.Errorf("computeService.ForwardingRules.Insert: %v", err)
	}
	return nil
}

func (s *sessionGCP) formNetworkURL(network string) string {
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s",
		s.gProjectID, network,
	)
}

func (s *sessionGCP) formSubnetURL(region, subnet string) string {
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks/%s",
		s.gProjectID, region, subnet,
	)
}

func (s *sessionGCP) formAddressURL(region, addressName string) string {
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/addresses/%s",
		s.gProjectID, region, addressName,
	)
}
