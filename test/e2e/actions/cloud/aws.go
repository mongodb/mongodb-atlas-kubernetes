package cloud

import (
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	aws "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
)

type awsAction struct{}

func (awsAction *awsAction) createPrivateEndpoint(pe status.ProjectPrivateEndpoint, privatelinkName string) (string, error) {
	fmt.Print("create AWS LINK")
	session := aws.SessionAWS(pe.Region)
	vpcID, err := session.GetVPCID()
	if err != nil {
		return "", err
	}
	subnetID, err := session.GetSubnetID()
	if err != nil {
		return "", err
	}

	privateEndpointID, err := session.CreatePrivateEndpoint(vpcID, subnetID, pe.ServiceName, privatelinkName)
	if err != nil {
		return "", err
	}

	return privateEndpointID, nil
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
	if status != "deleting" {
		return errors.New("AWS PE status is NOT 'deleting'. Actual status: " + status)
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
