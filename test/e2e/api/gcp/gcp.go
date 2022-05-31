package gcp

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

type sessionGCP struct {
	computeService *compute.Service
	gProjectID     string
}

func SessionGCP(gProjectID string) (sessionGCP, error) {
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return sessionGCP{}, fmt.Errorf("compute.NewClient: %w", err)
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
		return "", fmt.Errorf("computeService.Addresses.Insert: %w", err)
	}
	ip, err := s.GetIP(region, addressName, 20, 10)
	if err != nil {
		return "", fmt.Errorf("computeService.Addresses.Get: %w", err)
	}
	return ip, nil
}

func (s *sessionGCP) GetIP(region, addressName string, try, interval int) (string, error) {
	for i := 0; i < try; i++ {
		r, err := s.computeService.Addresses.Get(s.gProjectID, region, addressName).Do()
		if err != nil {
			return "", err
		}
		if r.Address != "" {
			return r.Address, nil
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return "", fmt.Errorf("timeout computeService.Addresses.Get")
}

func (s *sessionGCP) DeleteIPAdress(region, addressName string) error {
	_, err := s.computeService.Addresses.Delete(s.gProjectID, region, addressName).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("computeService.Addresses.Delete: %w", err)
	}
	return nil
}

func (s *sessionGCP) AddForwardRule(region, ruleName, addressName, network, subnet, target string) error {
	rules := &compute.ForwardingRule{
		IPAddress:                     s.formAddressURL(region, addressName),
		Labels:                        map[string]string{},
		Name:                          ruleName,
		Network:                       s.formNetworkURL(network),
		Ports:                         []string{},
		Region:                        region,
		ServiceDirectoryRegistrations: []*compute.ForwardingRuleServiceDirectoryRegistration{},
		Subnetwork:                    "",
		Target:                        target,
		ServerResponse:                googleapi.ServerResponse{},
	}
	_, err := s.computeService.ForwardingRules.Insert(s.gProjectID, region, rules).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("computeService.ForwardingRules.Insert: %w", err)
	}
	return nil
}

func (s *sessionGCP) DeleteForwardRule(region, ruleName string, try int, interval time.Duration) error {
	_, err := s.computeService.ForwardingRules.Delete(s.gProjectID, region, ruleName).Do()
	if err != nil {
		return fmt.Errorf("computeService.ForwardingRules.Delete: %w", err)
	}

	contain := func(list []*compute.ForwardingRule, name string) bool {
		for _, item := range list {
			if item.Name == name {
				return true
			}
		}
		return false
	}

	deleted := false
	for i := 0; i < try; i++ {
		r, err := s.computeService.ForwardingRules.List(s.gProjectID, region).Do()
		if err != nil {
			return fmt.Errorf("computeService.ForwardingRule.List: %w", err)
		}
		if !contain(r.Items, ruleName) {
			deleted = true
			break
		}
		time.Sleep(interval)
	}
	if !deleted {
		return fmt.Errorf("computeService.ForwardingRules.Delete. Could not delete forward rule after %d retries", try)
	}

	return nil
}

// Possible values:
// "ACCEPTED" - The connection has been accepted by the producer.
// "CLOSED" - The connection has been closed by the producer and will
// not serve traffic going forward.
// "PENDING" - The connection is pending acceptance by the producer.
// "REJECTED" - The connection has been rejected by the producer.
// "STATUS_UNSPECIFIED"
func (s *sessionGCP) DescribePrivateLinkStatus(region, ruleName string) (string, error) {
	resp, err := s.computeService.ForwardingRules.Get(s.gProjectID, region, ruleName).Context(context.Background()).Do()
	if err != nil {
		return "", fmt.Errorf("computeService.Addresses.Get: %w", err)
	}
	return resp.PscConnectionStatus, nil
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
