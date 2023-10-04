package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/api/compute/v1"
)

const (
	fileNameSAGCP = "gcp_service_account.json"
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
	var allErr error
	for _, network := range networks.Items {
		if strings.HasPrefix(network.Name, nameSubstr) {
			log.Printf("deleting network %s", network.Name)
			_, err = networkService.Delete(gcpProjectID, network.Name).Do()
			if err != nil {
				allErr = errors.Join(allErr, fmt.Errorf("failed to delete network %s: %s", network.Name, err))
				ok = false
			}
		}
	}
	return ok, allErr
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
