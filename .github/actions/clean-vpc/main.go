package main

import (
	"log"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/networkpeer"
)

func main() {
	if os.Getenv("DELETE_VPC") == "true" {
		err := setGCPCredentials()
		if err != nil {
			log.Fatal(err)
		}
		gcpOk, err := deleteGCPVPCBySubstr(cloud.GoogleProjectID, networkpeer.GCPVPCName)
		if err != nil {
			log.Fatal(err)
		}
		if !gcpOk {
			log.Println("Not all GCP VPC was deleted")
		}
		subID := os.Getenv("AZURE_SUBSCRIPTION_ID")
		if subID == "" {
			log.Fatal("AZURE_SUBSCRIPTION_ID is not set")
		}
		awsOk, err := deleteAzureVPCBySubstr(subID, networkpeer.AzureResourceGroupName, networkpeer.AzureVPCName)
		if err != nil {
			log.Fatal(err)
		}
		if !awsOk {
			log.Println("Not all Azure VPC was deleted")
		}
	}
}
