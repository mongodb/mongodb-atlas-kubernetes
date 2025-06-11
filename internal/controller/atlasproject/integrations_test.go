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

package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
)

const (
	testProjectID = "project-id"

	testNamespace = "some-namespace"
)

var errTest = fmt.Errorf("fake test error")

func TestToAlias(t *testing.T) {
	sample := []*mongodbatlas.ThirdPartyIntegration{{
		Type:   "DATADOG",
		APIKey: "some",
		Region: "EU",
	}}
	result := toAliasThirdPartyIntegration(sample)
	assert.Equal(t, sample[0].APIKey, result[0].APIKey)
	assert.Equal(t, sample[0].Type, result[0].Type)
	assert.Equal(t, sample[0].Region, result[0].Region)
}

func TestUpdateIntegrationsAtlas(t *testing.T) {
	calls := 0
	for _, tc := range []struct {
		title          string
		toUpdate       [][]set.DeprecatedIdentifiable
		client         *mongodbatlas.Client
		expectedResult workflow.Result
		expectedCalls  int
	}{
		{
			title:          "nil list does nothing",
			expectedResult: workflow.OK(),
		},

		{
			title:          "empty list does nothing",
			toUpdate:       [][]set.DeprecatedIdentifiable{},
			expectedResult: workflow.OK(),
		},

		{
			title: "different integrations get updated",
			toUpdate: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/some-otherpath/some-othersecret",
						Enabled:                  true,
					},
				}),
			client: &mongodbatlas.Client{
				Integrations: &atlas.IntegrationsMock{
					ReplaceFunc: func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
						calls += 1
						return nil, nil, nil
					},
				},
			},
			expectedResult: workflow.OK(),
			expectedCalls:  1,
		},

		{
			title: "matching integrations get updated anyway",
			toUpdate: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				}),
			client: &mongodbatlas.Client{
				Integrations: &atlas.IntegrationsMock{
					ReplaceFunc: func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
						calls += 1
						return nil, nil, nil
					},
				},
			},
			expectedResult: workflow.OK(),
			expectedCalls:  1,
		},

		{
			title: "integrations fail to update and return error",
			toUpdate: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				}),
			client: &mongodbatlas.Client{
				Integrations: &atlas.IntegrationsMock{
					ReplaceFunc: func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
						calls += 1
						return nil, nil, errTest
					},
				},
			},
			expectedResult: workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Errorf("cannot apply integration: %w", errTest)),
			expectedCalls:  1,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
				Client:  tc.client,
			}
			r := AtlasProjectReconciler{}
			calls = 0
			result := r.updateIntegrationsAtlas(workflowCtx, testProjectID, tc.toUpdate, testNamespace)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedCalls, calls)
		})
	}
}

func TestCheckIntegrationsReady(t *testing.T) {
	for _, tc := range []struct {
		title     string
		toCheck   [][]set.DeprecatedIdentifiable
		requested []project.Integration
		expected  bool
	}{
		{
			title:    "nil list does nothing",
			expected: true,
		},

		{
			title:     "empty list does nothing",
			toCheck:   [][]set.DeprecatedIdentifiable{},
			requested: []project.Integration{},
			expected:  true,
		},

		{
			title:     "when requested list differs in length it bails early",
			toCheck:   [][]set.DeprecatedIdentifiable{},
			requested: []project.Integration{{}},
			expected:  false,
		},

		{
			title: "matching integrations are considered applied",
			toCheck: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				}),
			requested: []project.Integration{{}},
			expected:  true,
		},

		{
			title: "different integrations are considered also applied",
			toCheck: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/some-otherpath/some-othersecret",
						Enabled:                  true,
					},
				}),
			requested: []project.Integration{{}},
			expected:  true,
		},

		{
			title: "matching integrations including prometheus are considered applied",
			toCheck: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
					{
						Type:             "PROMETHEUS",
						UserName:         "prometheus",
						ServiceDiscovery: "http",
						Enabled:          true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
					{
						Type:             "PROMETHEUS",
						UserName:         "prometheus",
						ServiceDiscovery: "http",
						Enabled:          true,
					},
				}),
			requested: []project.Integration{{}, {}},
			expected:  true,
		},

		{
			title: "matching integrations with a differing prometheus are considered different",
			toCheck: set.DeprecatedIntersection(
				[]aliasThirdPartyIntegration{
					{
						Type:                     "MICROSOFT_TEAMS",
						Name:                     testNamespace,
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
					{
						Type:             "PROMETHEUS",
						UserName:         "prometheus",
						ServiceDiscovery: "http",
						Enabled:          true,
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
					{
						Type:             "PROMETHEUS",
						UserName:         "zeus",
						ServiceDiscovery: "file",
						Enabled:          true,
					},
				}),
			requested: []project.Integration{{}, {}},
			expected:  false,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
			}
			r := AtlasProjectReconciler{}
			result := r.checkIntegrationsReady(workflowCtx, tc.toCheck, tc.requested)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapLastAppliedProjectIntegrations(t *testing.T) {
	for _, tt := range []struct {
		name    string
		project *akov2.AtlasProject
		want    []project.Integration
		wantErr bool
	}{
		{
			name: "returns integrations when present",
			project: setLastApplied(
				defaultTestProject(),
				mustJSONIZE(
					appendIntegrations(&defaultTestProject().Spec, []project.Integration{
						{Type: "TYPE1"},
						{Type: "TYPE2"},
					}),
				),
			),
			want:    []project.Integration{{Type: "TYPE1"}, {Type: "TYPE2"}},
			wantErr: false,
		},
		{
			name: "returns nil when no integrations",
			project: setLastApplied(
				defaultTestProject(),
				mustJSONIZE(appendIntegrations(&defaultTestProject().Spec, nil)),
			),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "returns nil when lastApplied is nil",
			project: defaultTestProject(),
			want:    nil,
			wantErr: false,
		},
		{
			name: "returns error if lastAppliedSpecFrom fails",
			project: setLastApplied(
				defaultTestProject(),
				"broken json",
			),
			want:    nil,
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapLastAppliedProjectIntegrations(tt.project)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func defaultTestProject() *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-name",
			Namespace: "ns",
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "test-project",
		},
	}
}

func appendIntegrations(spec *akov2.AtlasProjectSpec, integrations []project.Integration) *akov2.AtlasProjectSpec {
	spec.Integrations = integrations
	return spec
}

func mustJSONIZE(obj any) string {
	js, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(js)
}

func setLastApplied(project *akov2.AtlasProject, lastApplied string) *akov2.AtlasProject {
	if project.Annotations == nil {
		project.Annotations = map[string]string{}
	}
	project.Annotations[customresource.AnnotationLastAppliedConfiguration] = lastApplied
	return project
}

func TestFilterOwnedIntegrations(t *testing.T) {
	type args struct {
		integrationIDs []project.Integration
		lastApplied    []project.Integration
	}
	for _, tc := range []struct {
		name string
		args args
		want []set.DeprecatedIdentifiable
	}{
		{
			name: "returns only owned integrations",
			args: args{
				integrationIDs: []project.Integration{
					{Type: "TYPE1"},
					{Type: "TYPE2"},
					{Type: "TYPE3"},
				},
				lastApplied: []project.Integration{
					{Type: "TYPE1"},
					{Type: "TYPE3"},
				},
			},
			want: toIdentifiableSlice([]project.Integration{
				{Type: "TYPE1"},
				{Type: "TYPE3"},
			}),
		},
		{
			name: "returns nil if no integrationIDs",
			args: args{
				integrationIDs: nil,
				lastApplied:    []project.Integration{{Type: "TYPE1"}},
			},
			want: nil,
		},
		{
			name: "returns nil if none are owned",
			args: args{
				integrationIDs: []project.Integration{
					{Type: "TYPE4"},
				},
				lastApplied: []project.Integration{
					{Type: "TYPE1"},
					{Type: "TYPE2"},
				},
			},
			want: []set.DeprecatedIdentifiable{},
		},
		{
			name: "returns all if all are owned",
			args: args{
				integrationIDs: []project.Integration{
					{Type: "TYPE1"},
					{Type: "TYPE2"},
				},
				lastApplied: []project.Integration{
					{Type: "TYPE1"},
					{Type: "TYPE2"},
				},
			},
			want: toIdentifiableSlice([]project.Integration{
				{Type: "TYPE1"},
				{Type: "TYPE2"},
			}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := filterOwnedIntegrations(
				toIdentifiableSlice(tc.args.integrationIDs),
				tc.args.lastApplied,
			)
			require.Equal(t, tc.want, got)
		})
	}
}

type identifiableIntegration struct {
	project.Integration
}

func (ii identifiableIntegration) Identifier() interface{} {
	return ii.Integration.Type
}

func toIdentifiableSlice(integrations []project.Integration) []set.DeprecatedIdentifiable {
	identifiables := make([]set.DeprecatedIdentifiable, 0, len(integrations))
	for _, integration := range integrations {
		identifiables = append(identifiables, identifiableIntegration{Integration: integration})
	}
	return identifiables
}
