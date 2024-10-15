package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	testProjectID = "project-id"

	testNamespace = "some-namespace"
)

var errTest = fmt.Errorf("fake test error")

func TestToAlias(t *testing.T) {
	sample := []admin.ThirdPartyIntegration{{
		Type:   pointer.MakePtr("DATADOG"),
		ApiKey: pointer.MakePtr("some"),
		Region: pointer.MakePtr("EU"),
	}}
	result := toAliasThirdPartyIntegration(sample)
	assert.Equal(t, sample[0].ApiKey, result[0].ApiKey)
	assert.Equal(t, sample[0].Type, result[0].Type)
	assert.Equal(t, sample[0].Region, result[0].Region)
}

func TestUpdateIntegrationsAtlas(t *testing.T) {
	for _, tc := range []struct {
		title          string
		toUpdate       [][]set.Identifiable
		integrationAPI *mockadmin.ThirdPartyIntegrationsApi
		expectedResult workflow.Result
	}{
		{
			title:          "nil list does nothing",
			expectedResult: workflow.OK(),
		},

		{
			title:          "empty list does nothing",
			toUpdate:       [][]set.Identifiable{},
			expectedResult: workflow.OK(),
		},

		{
			title: "different integrations get updated",
			toUpdate: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/some-otherpath/some-othersecret",
						Enabled:                  true,
					},
				}),
			integrationAPI: func() *mockadmin.ThirdPartyIntegrationsApi {
				api := mockadmin.NewThirdPartyIntegrationsApi(t)
				api.EXPECT().UpdateThirdPartyIntegration(context.Background(), "MICROSOFT_TEAMS", "project-id", mock.AnythingOfType("*admin.ThirdPartyIntegration")).
					Return(admin.UpdateThirdPartyIntegrationApiRequest{ApiService: api})
				api.EXPECT().UpdateThirdPartyIntegrationExecute(mock.Anything).
					Return(&admin.PaginatedIntegration{}, &http.Response{}, nil)
				return api
			}(),
			expectedResult: workflow.OK(),
		},

		{
			title: "matching integrations get updated anyway",
			toUpdate: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				}),
			integrationAPI: func() *mockadmin.ThirdPartyIntegrationsApi {
				api := mockadmin.NewThirdPartyIntegrationsApi(t)
				api.EXPECT().UpdateThirdPartyIntegration(context.Background(), "MICROSOFT_TEAMS", "project-id", mock.AnythingOfType("*admin.ThirdPartyIntegration")).
					Return(admin.UpdateThirdPartyIntegrationApiRequest{ApiService: api}).Once()
				api.EXPECT().UpdateThirdPartyIntegrationExecute(mock.Anything).
					Return(&admin.PaginatedIntegration{}, &http.Response{}, nil).Once()
				return api
			}(),
			expectedResult: workflow.OK(),
		},

		{
			title: "integrations fail to update and return error",
			toUpdate: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
					},
				},
				[]project.Integration{
					{
						Type:                     "MICROSOFT_TEAMS",
						MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
						Enabled:                  true,
					},
				}),
			integrationAPI: func() *mockadmin.ThirdPartyIntegrationsApi {
				api := mockadmin.NewThirdPartyIntegrationsApi(t)
				api.EXPECT().UpdateThirdPartyIntegration(context.Background(), "MICROSOFT_TEAMS", "project-id", mock.AnythingOfType("*admin.ThirdPartyIntegration")).
					Return(admin.UpdateThirdPartyIntegrationApiRequest{ApiService: api}).Once()
				api.EXPECT().UpdateThirdPartyIntegrationExecute(mock.Anything).
					Return(&admin.PaginatedIntegration{}, &http.Response{}, errors.New("fake test error")).Once()
				return api
			}(),
			expectedResult: workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Can not apply integration: %v", errTest)),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
				SdkClient: &admin.APIClient{
					ThirdPartyIntegrationsApi: tc.integrationAPI,
				},
			}
			r := AtlasProjectReconciler{}
			result := r.updateIntegrationsAtlas(workflowCtx, testProjectID, tc.toUpdate, testNamespace)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCheckIntegrationsReady(t *testing.T) {
	for _, tc := range []struct {
		title     string
		toCheck   [][]set.Identifiable
		requested []project.Integration
		expected  bool
	}{
		{
			title:    "nil list does nothing",
			expected: true,
		},

		{
			title:     "empty list does nothing",
			toCheck:   [][]set.Identifiable{},
			requested: []project.Integration{},
			expected:  true,
		},

		{
			title:     "when requested list differs in length it bails early",
			toCheck:   [][]set.Identifiable{},
			requested: []project.Integration{{}},
			expected:  false,
		},

		{
			title: "matching integrations are considered applied",
			toCheck: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
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
			toCheck: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
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
			toCheck: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
					},
					{
						admin.ThirdPartyIntegration{
							Type:             pointer.MakePtr("PROMETHEUS"),
							Username:         pointer.MakePtr("prometheus"),
							ServiceDiscovery: pointer.MakePtr("http"),
							Enabled:          pointer.MakePtr(true),
						},
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
			toCheck: set.Intersection(
				[]aliasThirdPartyIntegration{
					{
						admin.ThirdPartyIntegration{
							Type:                     pointer.MakePtr("MICROSOFT_TEAMS"),
							MicrosoftTeamsWebhookUrl: pointer.MakePtr("https://somehost/somepath/somesecret"),
							Enabled:                  pointer.MakePtr(true),
						},
					},
					{
						admin.ThirdPartyIntegration{
							Type:             pointer.MakePtr("PROMETHEUS"),
							Username:         pointer.MakePtr("prometheus"),
							ServiceDiscovery: pointer.MakePtr("http"),
							Enabled:          pointer.MakePtr(true),
						},
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
			result := r.checkIntegrationsReady(workflowCtx, testNamespace, tc.toCheck, tc.requested)
			assert.Equal(t, tc.expected, result)
		})
	}
}
