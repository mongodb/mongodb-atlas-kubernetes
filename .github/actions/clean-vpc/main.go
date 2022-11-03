package main

import (
	"context"
	"log"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/networkpeer"
)

func main() {
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
	ctx := context.Background()
	azureOk, err := deleteAzureVPCBySubstr(ctx, subID, networkpeer.AzureResourceGroupName, networkpeer.AzureVPCName)
	if err != nil {
		log.Fatal(err)
	}
	if !azureOk {
		log.Println("Not all Azure VPC was deleted")
	}
	if !azureOk || !gcpOk {
		os.Exit(1)
	}
}
