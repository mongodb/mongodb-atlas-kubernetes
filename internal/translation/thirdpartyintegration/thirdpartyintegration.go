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

package thirdpartyintegration

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

var (
	// ErrNotFound is returned when the expected integration is not found
	ErrNotFound = errors.New("integration not found")
)

type ThirdPartyIntegrationService interface {
	List(ctx context.Context, projectID string) ([]*ThirdPartyIntegration, error)
	Create(ctx context.Context, projectID string, integration *ThirdPartyIntegration) (*ThirdPartyIntegration, error)
	Get(ctx context.Context, projectID, integrationType string) (*ThirdPartyIntegration, error)
	Update(ctx context.Context, projectID string, integration *ThirdPartyIntegration) (*ThirdPartyIntegration, error)
	Delete(ctx context.Context, projectID, integrationType string) error
}

func NewThirdPartyIntegrationServiceFromClientSet(clientSet *atlas.ClientSet) ThirdPartyIntegrationService {
	return NewThirdPartyIntegrationService(clientSet.SdkClient20250312011.ThirdPartyIntegrationsApi)
}

func NewThirdPartyIntegrationService(integrationsAPI admin.ThirdPartyIntegrationsApi) ThirdPartyIntegrationService {
	return &thirdPartyIntegration{integrationsAPI: integrationsAPI}
}

type thirdPartyIntegration struct {
	integrationsAPI admin.ThirdPartyIntegrationsApi
}

func (tpi *thirdPartyIntegration) List(ctx context.Context, projectID string) ([]*ThirdPartyIntegration, error) {
	list, _, err := tpi.integrationsAPI.ListGroupIntegrations(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list integrations for project %v: %w", projectID, err)
	}

	result := make([]*ThirdPartyIntegration, 0, len(list.GetResults()))
	for _, i := range list.GetResults() {
		integration, err := fromAtlas(&i)
		if err != nil {
			return nil, fmt.Errorf("failed to convert integration from Atlas: %w", err)
		}

		result = append(result, integration)
	}

	return result, nil
}

func (tpi *thirdPartyIntegration) Create(ctx context.Context, projectID string, integration *ThirdPartyIntegration) (*ThirdPartyIntegration, error) {
	atlasIntegration, err := toAtlas(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to convert integration to Atlas: %w", err)
	}
	integrationPages, _, err := tpi.integrationsAPI.CreateGroupIntegration(
		ctx, atlasIntegration.GetType(), projectID, atlasIntegration).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create integration from config: %w", err)
	}
	newIntegration, err := getResultOfType(integrationPages.GetResults(), integration.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert integration from Atlas: %w", err)
	}
	return newIntegration, nil
}

func (tpi *thirdPartyIntegration) Get(ctx context.Context, projectID, integrationType string) (*ThirdPartyIntegration, error) {
	atlasIntegration, _, err := tpi.integrationsAPI.GetGroupIntegration(ctx, projectID, integrationType).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "INTEGRATION_NOT_CONFIGURED") {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get integration type %v for project %v: %w", integrationType, projectID, err)
	}
	peer, err := fromAtlas(atlasIntegration)
	if err != nil {
		return nil, fmt.Errorf("failed to convert integration from Atlas: %w", err)
	}
	return peer, nil
}

func (tpi *thirdPartyIntegration) Update(ctx context.Context, projectID string, integration *ThirdPartyIntegration) (*ThirdPartyIntegration, error) {
	atlasIntegration, err := toAtlas(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to convert integration to Atlas: %w", err)
	}
	integrationPages, _, err := tpi.integrationsAPI.UpdateGroupIntegration(
		ctx, atlasIntegration.GetType(), projectID, atlasIntegration).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update integration with config %v: %w", integration, err)
	}
	updatedIntegration, err := getResultOfType(integrationPages.GetResults(), integration.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert integration from Atlas: %w", err)
	}
	return updatedIntegration, nil
}

func (tpi *thirdPartyIntegration) Delete(ctx context.Context, projectID, integrationType string) error {
	_, err := tpi.integrationsAPI.DeleteGroupIntegration(ctx, integrationType, projectID).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "INTEGRATION_NOT_CONFIGURED") {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete integration type %s: %w", integrationType, err)
	}
	return nil
}

func getResultOfType(integrations []admin.ThirdPartyIntegration, typeName string) (*ThirdPartyIntegration, error) {
	if err := assertType(typeName); err != nil {
		return nil, fmt.Errorf("wrong target type: %w", err)
	}
	for _, integration := range integrations {
		tn := integration.GetType()
		if err := assertType(tn); err != nil {
			return nil, fmt.Errorf("wrong result type: %w", err)
		}
		if tn == typeName {
			return fromAtlas(&integration)
		}
	}
	return nil, fmt.Errorf("integration %q %w in API reply", typeName, ErrNotFound)
}
