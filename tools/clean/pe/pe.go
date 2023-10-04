package pe

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	subnet1Name       = "atlas-operator-e2e-test-subnet1"
	subnet2Name       = "atlas-operator-e2e-test-subnet2"
	fileNameSAGCP     = "gcp_service_account.json"
	googleProjectID   = "atlasoperator"
	gcpVPCName        = "network-peering-gcp-1-vpc"
	resourceGroupName = "svet-test"
)

func CleanPEs() {
	err := setGCPCredentials()
	if err != nil {
		log.Fatal(err)
	}

	err = cleanAllPE()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PE cleanups completed OK")
}

func setGCPCredentials() error {
	err := os.MkdirAll(filepath.Dir(fileNameSAGCP), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileNameSAGCP, []byte(os.Getenv("GCP_SA_CRED")), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", fileNameSAGCP)
	if err != nil {
		return fmt.Errorf("error setting GOOGLE_APPLICATION_CREDENTIALS: %v", err)
	}

	return nil
}

func cleanAllPE() error {
	ctx := context.Background()
	groupNameAzure := resourceGroupName
	awsRegions := []string{
		"eu-west-2",
		"us-east-1",
	}
	gcpRegion := "europe-west1"
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	var allErr error
	err := cleanAllAzurePE(ctx, groupNameAzure, subscriptionID, []string{subnet1Name, subnet2Name})
	if err != nil {
		allErr = errors.Join(allErr, fmt.Errorf("error while cleaning all azure pe: %v", err))
	}

	for _, awsRegion := range awsRegions {
		errClean := cleanAllAWSPE(awsRegion, []string{subnet1Name, subnet2Name})
		if errClean != nil {
			allErr = errors.Join(allErr, fmt.Errorf("error cleaning all aws PE. region %s. error: %v", awsRegion, errClean))
		}
	}

	err = cleanAllGCPPE(ctx, googleProjectID, gcpVPCName, gcpRegion, []string{subnet1Name, subnet2Name})
	if err != nil {
		allErr = errors.Join(allErr, fmt.Errorf("error while cleaning all gcp pe: %v", err))
	}

	return allErr
}
