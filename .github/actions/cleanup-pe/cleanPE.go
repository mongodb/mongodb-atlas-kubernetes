package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func main() {
	cleanOnlyTaggedPE := false
	if os.Getenv("CLEAN_PE") == "true" {
		cleanOnlyTaggedPE = true
	}
	err := SetGCPCredentials()
	if err != nil {
		log.Fatal(err)
	}
	err = CleanAllPE(cleanOnlyTaggedPE)
	if err != nil {
		log.Fatal(err)
	}
}

func SetGCPCredentials() error {
	err := utils.SaveToFile(config.FileNameSAGCP, []byte(os.Getenv("GCP_SA_CRED")))
	if err != nil {
		return fmt.Errorf("error saving gcp sa cred to file: %v", err)
	}
	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.FileNameSAGCP)
	if err != nil {
		return fmt.Errorf("error setting GOOGLE_APPLICATION_CREDENTIALS: %v", err)
	}
	return nil
}

func CleanAllPE(onlyTagged bool) error {
	ctx := context.Background()
	groupNameAzure := cloud.ResourceGroup
	awsRegions := []string{
		config.AWSRegionEU,
		config.AWSRegionUS,
	}
	gcpRegion := config.GCPRegion
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	if onlyTagged {
		err := cleanAllTaggedAzurePE(ctx, config.TagForTestKey, config.TagForTestValue, groupNameAzure, subscriptionID)
		if err != nil {
			return fmt.Errorf("error while cleaning all azure pe: %v", err)
		}

		for _, awsRegion := range awsRegions {
			errClean := cleanAllTaggedAWSPE(awsRegion, config.TagForTestKey, config.TagForTestValue)
			if errClean != nil {
				return fmt.Errorf("error cleaning all aws PE. region %s. error: %v", awsRegion, errClean)
			}
		}

		err = cleanAllTaggedGCPPE(ctx, cloud.GoogleProjectID, cloud.GoogleVPC,
			gcpRegion, config.TagForTestKey, config.TagForTestValue)
		if err != nil {
			return fmt.Errorf("error while cleaning all gcp pe: %v", err)
		}
	} else {
		err := cleanAllAzurePE(ctx, groupNameAzure, subscriptionID)
		if err != nil {
			return fmt.Errorf("error while cleaning all azure pe: %v", err)
		}

		for _, awsRegion := range awsRegions {
			errClean := cleanAllAWSPE(awsRegion)
			if errClean != nil {
				return fmt.Errorf("error cleaning all aws PE. region %s. error: %v", awsRegion, errClean)
			}
		}

		err = cleanAllGCPPE(ctx, cloud.GoogleProjectID, cloud.GoogleVPC, gcpRegion)
		if err != nil {
			return fmt.Errorf("error while cleaning all gcp pe: %v", err)
		}
	}
	return nil
}
