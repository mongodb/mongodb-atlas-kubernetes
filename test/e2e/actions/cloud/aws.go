package cloud

import (
	"errors"
	"fmt"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	aws "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
)

type awsAction struct{}

func (awsAction *awsAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (v1.PrivateEndpoint, error) {
	fmt.Print("create AWS LINK")
	session := aws.SessionAWS(pe.Region)
	vpcID, err := session.GetVPCID()
	if err != nil {
		return v1.PrivateEndpoint{}, err
	}
	subnetID, err := session.GetSubnetID()
	if err != nil {
		return v1.PrivateEndpoint{}, err
	}

	privateEndpointID, err := session.CreatePrivateEndpoint(vpcID, subnetID, pe.ServiceName, privatelinkName)
	if err != nil {
		return v1.PrivateEndpoint{}, err
	}
	cResponse := v1.PrivateEndpoint{
		ID:       privateEndpointID,
		IP:       "",
		Provider: provider.ProviderAWS,
		Region:   pe.Region,
	}
	return cResponse, nil
}

func (awsAction *awsAction) deletePrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkID string) error {
	fmt.Print("DELETE AWS LINK")
	session := aws.SessionAWS(pe.Region)
	err := session.DeletePrivateLink(privatelinkID)
	if err != nil {
		return err
	}
	status, err := session.DescribePrivateEndpointStatus(privatelinkID)
	if err != nil {
		return err
	}
	if status != "deleting" && status != "deleted" {
		return errors.New("AWS PE status is NOT 'deleting'/'deleted'. Actual status: " + status)
	}
	return nil
}

func (awsAction *awsAction) statusPrivateEndpointPending(region, privateID string) bool {
	fmt.Print("check AWS LINK status: ")
	session := aws.SessionAWS(region)
	status, err := session.DescribePrivateEndpointStatus(privateID)
	if err != nil {
		fmt.Print(err)
		return false
	}
	fmt.Println(status)
	return (status == "pendingAcceptance")
}

func (awsAction *awsAction) statusPrivateEndpointAvailable(region, privateID string) bool {
	fmt.Print("check AWS LINK status: ")
	session := aws.SessionAWS(region)
	status, err := session.DescribePrivateEndpointStatus(privateID)
	if err != nil {
		fmt.Print(err)
		return false
	}
	fmt.Println(status)
	return (status == "available")
}
