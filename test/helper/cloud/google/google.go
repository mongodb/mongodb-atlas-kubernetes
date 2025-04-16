// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"context"
	"errors"
	"fmt"
	"os"

	compute "cloud.google.com/go/compute/apiv1"
)

const (
	googleSAFilename = ".googleServiceAccount.json"
)

type googleConnection struct {
	projectID string

	networkClient *compute.NetworksClient
}

func newConnection(ctx context.Context, projectID string) (*googleConnection, error) {
	if err := ensureCredentials(); err != nil {
		return nil, fmt.Errorf("failed to prepare credentials")
	}

	networkClient, err := compute.NewNetworksRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup network rest client")
	}

	return &googleConnection{
		projectID:     projectID,
		networkClient: networkClient,
	}, nil
}

func ensureCredentials() error {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		return nil
	}
	credentials := os.Getenv("GCP_SA_CRED")
	if credentials == "" {
		return errors.New("GOOGLE_APPLICATION_CREDENTIALS and GCP_SA_CRED are unset, cant setup Google credentials")
	}
	if err := os.WriteFile(googleSAFilename, ([]byte)(credentials), 0600); err != nil {
		return fmt.Errorf("failed to save credentials contents GCP_SA_CRED to %s: %w",
			googleSAFilename, err)
	}
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", googleSAFilename)
	return nil
}
