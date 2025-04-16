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

package atlasdeployment

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
)

type CMZTestData struct {
	desired             []akov2.CustomZoneMapping
	existing            map[string]string
	czmMap              map[string]string
	expectedToCreate    bool
	expectedToBeDeleted bool
	name                string
}

const (
	zone1     = "Zone 1"
	zone2     = "Zone 2"
	location1 = "CA"
	location2 = "US"
)

func runCMZTest(t *testing.T, testData *CMZTestData) {
	shouldCreate, shouldDelete := compareZoneMappingStates(testData.existing, testData.desired, testData.czmMap)
	if shouldCreate != testData.expectedToCreate {
		t.Errorf("Test: %s. expected to shouldCreate %v, got %v", testData.name, testData.expectedToCreate, shouldCreate)
	}
	if shouldDelete != testData.expectedToBeDeleted {
		t.Errorf("Test: %s. expected to shouldDelete %v, got %v", testData.name, testData.expectedToBeDeleted, shouldDelete)
	}
}

func TestCompareZoneMappingStates(t *testing.T) {
	tests := []*CMZTestData{
		{
			name: "All synced. No changes needed",
			desired: []akov2.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
				{
					Zone:     zone2,
					Location: location2,
				},
			},
			existing: map[string]string{
				location1: "1",
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    false,
			expectedToBeDeleted: false,
		},
		{
			name: "Wrong zone. Should be recreated",
			desired: []akov2.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
				{
					Zone:     zone2,
					Location: location2,
				},
			},
			existing: map[string]string{
				location1: "1",
				location2: "1",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    true,
			expectedToBeDeleted: true,
		},
		{
			name: "Exist more than needed. Should be recreated",
			desired: []akov2.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
			},
			existing: map[string]string{
				location1: "1",
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    true,
			expectedToBeDeleted: true,
		},
		{
			name:    "Empty desired. Should be deleted",
			desired: []akov2.CustomZoneMapping{},
			existing: map[string]string{
				location1: "1",
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    false,
			expectedToBeDeleted: true,
		},
		{
			name: "Exist less than needed. Should be created",
			desired: []akov2.CustomZoneMapping{
				{
					Zone:     zone1,
					Location: location1,
				},
				{
					Zone:     zone2,
					Location: location2,
				},
			},
			existing: map[string]string{
				location2: "2",
			},
			czmMap: map[string]string{
				"1": zone1,
				"2": zone2,
			},
			expectedToCreate:    true,
			expectedToBeDeleted: false,
		},
	}
	for _, test := range tests {
		runCMZTest(t, test)
	}
}

func TestEnsureCustomZoneMapping(t *testing.T) {
	projectID := "test-project"
	deploymentName := "test-deployment"

	for _, tc := range []struct {
		name               string
		customZoneMappings []akov2.CustomZoneMapping
		deploymentAPI      deployment.AtlasDeploymentsService
		isOK               bool
	}{
		{
			name: "GET errors",
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(nil, errors.New("test GET error"))
				return service
			}(),
			isOK: false,
		},
		{
			name: "No zone mappings in AKO or Atlas (no op)",
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(nil, nil)
				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(nil, nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Zone mapping in AKO but not Atlas (create)",
			customZoneMappings: []akov2.CustomZoneMapping{
				{
					Location: "test-location",
					Zone:     "test-zone",
				},
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(nil, nil)
				service.EXPECT().CreateCustomZones(context.Background(), projectID, deploymentName, mock.AnythingOfType("[]v1.CustomZoneMapping")).Return(map[string]string{"test-location": "test-zone"}, nil)

				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(map[string]string{}, nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Create errors",
			customZoneMappings: []akov2.CustomZoneMapping{
				{
					Location: "test-location",
					Zone:     "test-zone",
				},
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(nil, nil)
				service.EXPECT().CreateCustomZones(context.Background(), projectID, deploymentName, mock.AnythingOfType("[]v1.CustomZoneMapping")).Return(nil, errors.New("test POST error"))
				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(map[string]string{}, nil)

				return service
			}(),
			isOK: false,
		},
		{
			name:               "Zone mapping in Atlas but not AKO (delete)",
			customZoneMappings: []akov2.CustomZoneMapping{},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(map[string]string{"test-location": "test-zone"}, nil)
				service.EXPECT().DeleteCustomZones(context.Background(), projectID, deploymentName).Return(nil)
				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(map[string]string{"test-id": "test-zone"}, nil)

				return service
			}(),
			isOK: true,
		},
		{
			name:               "Delete errors",
			customZoneMappings: []akov2.CustomZoneMapping{},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(map[string]string{"test-location": "test-zone"}, nil)
				service.EXPECT().DeleteCustomZones(context.Background(), projectID, deploymentName).Return(errors.New("test DELETE error"))
				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(map[string]string{"test-id": "test-zone"}, nil)

				return service
			}(),
			isOK: false,
		},
		{
			name: "Zone mapping the same in Atlas and AKO (no op)",
			customZoneMappings: []akov2.CustomZoneMapping{
				{
					Location: "test-location",
					Zone:     "test-zone",
				},
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(map[string]string{"test-location": "test-id"}, nil)
				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(map[string]string{"test-id": "test-zone"}, nil)

				return service
			}(),
			isOK: true,
		},
		{
			name: "Differing zone mapping in Atlas and AKO (update)",
			customZoneMappings: []akov2.CustomZoneMapping{
				{
					Location: "new-test-location",
					Zone:     "new-test-zone",
				},
			},
			deploymentAPI: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				service.EXPECT().GetCustomZones(context.Background(), projectID, deploymentName).Return(map[string]string{"test-location": "test-id"}, nil)
				service.EXPECT().CreateCustomZones(context.Background(), projectID, deploymentName, mock.AnythingOfType("[]v1.CustomZoneMapping")).Return(map[string]string{"new-test-location": "new-test-zone"}, nil)
				service.EXPECT().DeleteCustomZones(context.Background(), projectID, deploymentName).Return(nil)
				service.EXPECT().GetZoneMapping(context.Background(), projectID, deploymentName).Return(map[string]string{"test-id": "test-zone"}, nil)

				return service
			}(),
			isOK: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &AtlasDeploymentReconciler{}
			ctx := &workflow.Context{
				Log:     zaptest.NewLogger(t).Sugar(),
				Context: context.Background(),
			}

			result := r.ensureCustomZoneMapping(
				ctx,
				tc.deploymentAPI,
				projectID,
				tc.customZoneMappings,
				deploymentName,
			)

			equal := (result.IsOk() == tc.isOK)
			if !equal {
				t.Log(result.GetMessage())
			}
			require.True(t, equal)
		})
	}
}
