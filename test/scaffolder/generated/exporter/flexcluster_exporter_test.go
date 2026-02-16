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

package exporter

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admin "go.mongodb.org/atlas-sdk/v20250312013/admin"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/refs"
)

// mockFlexClustersApi implements admin.FlexClustersApi for testing.
type mockFlexClustersApi struct {
	admin.FlexClustersApi
	listResponse *admin.PaginatedFlexClusters20241113
	listErr      error
	// For pagination testing: map pageNum to response
	paginatedResponses map[int]*admin.PaginatedFlexClusters20241113
}

func (m *mockFlexClustersApi) ListFlexClusters(ctx context.Context, groupId string) admin.ListFlexClustersApiRequest {
	return admin.ListFlexClustersApiRequest{ApiService: m}
}

func (m *mockFlexClustersApi) ListFlexClustersExecute(r admin.ListFlexClustersApiRequest) (*admin.PaginatedFlexClusters20241113, *http.Response, error) {
	if m.listErr != nil {
		return nil, nil, m.listErr
	}
	// For paginated responses, we need to check the page number
	// Since we can't easily extract it from the request, use the single response
	if m.paginatedResponses != nil {
		// This is a simplified pagination mock - in real tests we'd need more sophisticated tracking
		return m.listResponse, &http.Response{StatusCode: http.StatusOK}, nil
	}
	return m.listResponse, &http.Response{StatusCode: http.StatusOK}, nil
}

// mockTranslator implements crapi.Translator for testing purposes.
type mockTranslator struct {
	scheme       *runtime.Scheme
	fromAPIFunc  func(target client.Object, source any, objs ...client.Object) ([]client.Object, error)
	fromAPIError error
}

func (m *mockTranslator) Scheme() *runtime.Scheme {
	return m.scheme
}

func (m *mockTranslator) MajorVersion() string {
	return "v20250312013"
}

func (m *mockTranslator) Mappings() ([]*refs.Mapping, error) {
	return nil, nil
}

func (m *mockTranslator) ToAPI(target any, source client.Object, objs ...client.Object) error {
	return nil
}

func (m *mockTranslator) FromAPI(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
	if m.fromAPIError != nil {
		return nil, m.fromAPIError
	}
	if m.fromAPIFunc != nil {
		return m.fromAPIFunc(target, source, objs...)
	}
	return []client.Object{target}, nil
}

var _ crapi.Translator = (*mockTranslator)(nil)

// mockClientObject implements client.Object for testing.
type mockClientObject struct {
	client.Object
	name string
}

func (m *mockClientObject) GetName() string {
	return m.name
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func TestExport(t *testing.T) {
	tests := []struct {
		name              string
		mockApi           *mockFlexClustersApi
		translator        *mockTranslator
		identifiers       []string
		referencedObjects []client.Object
		wantResultsCount  int
		wantErr           bool
		wantErrContains   string
	}{
		{
			name: "exports single cluster successfully",
			mockApi: &mockFlexClustersApi{
				listResponse: &admin.PaginatedFlexClusters20241113{
					Results: &[]admin.FlexClusterDescription20241113{
						{Name: stringPtr("cluster-1")},
					},
					TotalCount: intPtr(1),
				},
			},
			translator: &mockTranslator{
				fromAPIFunc: func(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
					return []client.Object{&mockClientObject{name: "cluster-1"}}, nil
				},
			},
			identifiers:      []string{"project-id"},
			wantResultsCount: 2, // primary resource + 1 translated object
			wantErr:          false,
		},
		{
			name: "exports multiple clusters successfully",
			mockApi: &mockFlexClustersApi{
				listResponse: &admin.PaginatedFlexClusters20241113{
					Results: &[]admin.FlexClusterDescription20241113{
						{Name: stringPtr("cluster-1")},
						{Name: stringPtr("cluster-2")},
						{Name: stringPtr("cluster-3")},
					},
					TotalCount: intPtr(3),
				},
			},
			translator: &mockTranslator{
				fromAPIFunc: func(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
					return []client.Object{target}, nil
				},
			},
			identifiers:      []string{"project-id"},
			wantResultsCount: 6, // 3 × (primary resource + 1 translated object)
			wantErr:          false,
		},
		{
			name: "returns empty when no clusters exist",
			mockApi: &mockFlexClustersApi{
				listResponse: &admin.PaginatedFlexClusters20241113{
					Results:    &[]admin.FlexClusterDescription20241113{},
					TotalCount: intPtr(0),
				},
			},
			translator:       &mockTranslator{},
			identifiers:      []string{"project-id"},
			wantResultsCount: 0,
			wantErr:          false,
		},
		{
			name: "returns error when API call fails",
			mockApi: &mockFlexClustersApi{
				listErr: errors.New("API connection failed"),
			},
			translator:      &mockTranslator{},
			identifiers:     []string{"project-id"},
			wantErr:         true,
			wantErrContains: "failed to list FlexClusters from Atlas",
		},
		{
			name: "returns error when API returns nil response",
			mockApi: &mockFlexClustersApi{
				listResponse: nil,
			},
			translator:      &mockTranslator{},
			identifiers:     []string{"project-id"},
			wantErr:         true,
			wantErrContains: "no response",
		},
		{
			name: "returns error when translator fails",
			mockApi: &mockFlexClustersApi{
				listResponse: &admin.PaginatedFlexClusters20241113{
					Results: &[]admin.FlexClusterDescription20241113{
						{Name: stringPtr("cluster-1")},
					},
					TotalCount: intPtr(1),
				},
			},
			translator: &mockTranslator{
				fromAPIError: errors.New("translation failed"),
			},
			identifiers:     []string{"project-id"},
			wantErr:         true,
			wantErrContains: "failed to translate FlexCluster",
		},
		{
			name: "translator returns multiple objects per resource",
			mockApi: &mockFlexClustersApi{
				listResponse: &admin.PaginatedFlexClusters20241113{
					Results: &[]admin.FlexClusterDescription20241113{
						{Name: stringPtr("cluster-1")},
					},
					TotalCount: intPtr(1),
				},
			},
			translator: &mockTranslator{
				fromAPIFunc: func(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
					// Simulate translator returning multiple objects (e.g., cluster + related secrets)
					return []client.Object{
						&mockClientObject{name: "cluster-1"},
						&mockClientObject{name: "cluster-1-secret"},
					}, nil
				},
			},
			identifiers:      []string{"project-id"},
			wantResultsCount: 3, // primary resource + 2 translated objects
			wantErr:          false,
		},
		{
			name: "passes referenced objects to translator",
			mockApi: &mockFlexClustersApi{
				listResponse: &admin.PaginatedFlexClusters20241113{
					Results: &[]admin.FlexClusterDescription20241113{
						{Name: stringPtr("cluster-1")},
					},
					TotalCount: intPtr(1),
				},
			},
			translator: &mockTranslator{
				fromAPIFunc: func(target client.Object, source any, objs ...client.Object) ([]client.Object, error) {
					// Verify that referenced objects are forwarded to the translator
					assert.Len(t, objs, 2)
					return []client.Object{target}, nil
				},
			},
			referencedObjects: []client.Object{
				&mockClientObject{name: "ref-obj-1"},
				&mockClientObject{name: "ref-obj-2"},
			},
			identifiers:      []string{"project-id"},
			wantResultsCount: 2, // primary resource + 1 translated object
			wantErr:          false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// Create API client with mock
			apiClient := &admin.APIClient{
				FlexClustersApi: tc.mockApi,
			}

			exporter := NewFlexClusterExporter(apiClient, tc.translator, tc.identifiers)
			results, err := exporter.Export(ctx, tc.referencedObjects)

			if tc.wantErr {
				require.Error(t, err)
				if tc.wantErrContains != "" {
					assert.Contains(t, err.Error(), tc.wantErrContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Len(t, results, tc.wantResultsCount)
		})
	}
}

func TestNewFlexClusterExporter(t *testing.T) {
	tests := []struct {
		name        string
		identifiers []string
	}{
		{
			name:        "creates exporter with single identifier",
			identifiers: []string{"project-id-123"},
		},
		{
			name:        "creates exporter with multiple identifiers",
			identifiers: []string{"project-id-123", "cluster-name"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			translator := &mockTranslator{}
			exp := NewFlexClusterExporter(nil, translator, tc.identifiers)

			require.NotNil(t, exp)
			assert.Equal(t, tc.identifiers, exp.identifiers)
			assert.Equal(t, translator, exp.translator)
		})
	}
}
