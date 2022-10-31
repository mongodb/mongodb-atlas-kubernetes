package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"

	"google.golang.org/api/compute/v1"
)

func deleteGCPVPCBySubstr(gcpProjectID, nameSubstr string) (bool, error) {
	ok := true
	computeService, err := compute.NewService(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to create compute service: %s", err)
	}
	networkService := compute.NewNetworksService(computeService)
	networks, err := networkService.List(gcpProjectID).Do()
	if err != nil {
		return false, fmt.Errorf("failed to list networks: %s", err)
	}
	for _, network := range networks.Items {
		if strings.HasPrefix(network.Name, nameSubstr) {
			log.Printf("deleting network %s", network.Name)
			_, err = networkService.Delete(gcpProjectID, network.Name).Do()
			if err != nil {
				log.Printf(fmt.Sprintf("failed to delete network %s: %s", network.Name, err))
				ok = false
			}
		}
	}
	return ok, nil
}

func setGCPCredentials() error {
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
