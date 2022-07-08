package cloud

import (
	"fmt"
	"time"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/gcp"
)

type gcpAction struct{}

const (
	// TODO get from GCP
	GoogleProjectID     = "atlasoperator"             // Google Cloud Project ID
	GoogleVPC           = "atlas-operator-test"       // VPC Name
	GoogleSubnetName    = "atlas-operator-subnet-leo" // Subnet Name
	googleConnectPrefix = "ao"                        // Private Service Connect Endpoint Prefix
)

func (gcpAction *gcpAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (v1.PrivateEndpoint, error) {
	session, err := gcp.SessionGCP(GoogleProjectID)
	if err != nil {
		return v1.PrivateEndpoint{}, err
	}
	var cResponse v1.PrivateEndpoint
	for i, target := range pe.ServiceAttachmentNames {
		addressName := formAddressName(privatelinkName, i)
		ruleName := formRuleName(privatelinkName, i)
		ip, err := session.AddIPAddress(pe.Region, addressName, GoogleSubnetName)
		if err != nil {
			return v1.PrivateEndpoint{}, err
		}

		cResponse.Endpoints = append(cResponse.Endpoints, v1.GCPEndpoint{
			EndpointName: ruleName,
			IPAddress:    ip,
		})
		cResponse.EndpointGroupName = GoogleVPC
		cResponse.Region = pe.Region
		cResponse.Provider = pe.Provider
		cResponse.GCPProjectID = GoogleProjectID

		session.AddForwardRule(pe.Region, ruleName, addressName, GoogleVPC, GoogleSubnetName, target)
	}
	return cResponse, nil
}

func (gcpAction *gcpAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) error {
	session, err := gcp.SessionGCP(GoogleProjectID)
	if err != nil {
		return err
	}
	for i := range pe.Endpoints {
		session.DeleteForwardRule(pe.Region, formRuleName(privatelinkName, i), 10, 20*time.Second)
		session.DeleteIPAdress(pe.Region, formAddressName(privatelinkName, i))
	}
	return nil
}

func (gcpAction *gcpAction) statusPrivateEndpointPending(region, privateID string) bool {
	session, err := gcp.SessionGCP(GoogleProjectID)
	if err != nil {
		fmt.Print(err)
		return false
	}
	ruleName := formRuleName(privateID, 1)
	result, err := session.DescribePrivateLinkStatus(region, ruleName)
	if err != nil {
		fmt.Print(err)
		return false
	}
	return (result == "PENDING")
}

func (gcpAction *gcpAction) statusPrivateEndpointAvailable(region, privateID string) bool {
	session, err := gcp.SessionGCP(GoogleProjectID)
	if err != nil {
		fmt.Print(err)
		return false
	}
	ruleName := formRuleName(privateID, 1)
	result, err := session.DescribePrivateLinkStatus(region, ruleName)
	if err != nil {
		fmt.Print(err)
		return false
	}
	return (result == "ACCEPTED")
}

func formAddressName(name string, i int) string {
	return fmt.Sprintf("%s%s-ip-%d", googleConnectPrefix, name, i)
}

func formRuleName(name string, i int) string {
	return fmt.Sprintf("%s%s-%d", googleConnectPrefix, name, i)
}
