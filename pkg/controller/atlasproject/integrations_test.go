package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/integrations"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	testProjectID = "project-id"

	testNamespace = "some-namespace"
)

var errTest = fmt.Errorf("fake test error")

func TestUpdateIntegrationsAtlas(t *testing.T) {
	for _, tc := range []struct {
		title              string
		toUpdate           [][]set.Identifiable
		integrationService integrations.AtlasIntegrationsService
		secret             *corev1.Secret
		expectedResult     workflow.Result
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
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/some-otherpath/some-othersecret",
							Enabled:                  true,
							APIKeyRef: common.ResourceRefNamespaced{
								Name:      "test-secret",
								Namespace: "test-ns",
							},
						},
					},
				}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Update(context.Background(), "project-id", "MICROSOFT_TEAMS", mock.AnythingOfType("integrations.Integration"), map[string]string{"apiKey": "Passw0rd!"}).Return(nil)
				return service
			}(),
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "test-ns",
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			expectedResult: workflow.OK(),
		},

		{
			title: "matching integrations get updated anyway",
			toUpdate: set.Intersection(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Update(context.Background(), "project-id", "MICROSOFT_TEAMS", mock.AnythingOfType("integrations.Integration"), map[string]string{}).Return(nil)
				return service
			}(),
			expectedResult: workflow.OK(),
		},

		{
			title: "integrations fail to update and return error",
			toUpdate: set.Intersection(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Update(context.Background(), "project-id", "MICROSOFT_TEAMS", mock.AnythingOfType("integrations.Integration"), map[string]string{}).Return(errors.New("fake test error"))
				return service
			}(),
			expectedResult: workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Cannot apply integration %v: %v", "MICROSOFT_TEAMS", errTest)),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme)

			if tc.secret != nil {
				k8sClient.WithObjects(tc.secret)
			}

			r := AtlasProjectReconciler{
				Client:              k8sClient.Build(),
				integrationsService: tc.integrationService,
			}
			result := r.updateIntegrationsAtlas(workflowCtx, testProjectID, tc.toUpdate, testNamespace)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestDeleteIntegrationsAtlas(t *testing.T) {
	testScheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(testScheme))
	assert.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		Build()

	for _, tc := range []struct {
		title              string
		toDelete           []set.Identifiable
		integrationService integrations.AtlasIntegrationsService
		expectedResult     workflow.Result
	}{
		{
			title:          "nil list does nothing",
			expectedResult: workflow.OK(),
		},
		{
			title: "successfully deletes",
			toDelete: set.Difference(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Delete(context.Background(), "project-id", "MICROSOFT_TEAMS").Return(nil)
				return service
			}(),
			expectedResult: workflow.OK(),
		},
		{
			title: "attempt to delete errors",
			toDelete: set.Difference(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Delete(context.Background(), "project-id", "MICROSOFT_TEAMS").Return(errors.New("fake test error"))
				return service
			}(),
			expectedResult: workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Cannot delete integration %v: %v", "MICROSOFT_TEAMS", errTest)),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
			}

			r := AtlasProjectReconciler{
				Client:              k8sClient,
				integrationsService: tc.integrationService,
			}
			result := r.deleteIntegrationsFromAtlas(workflowCtx, testProjectID, tc.toDelete)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCreateIntegrationsAtlas(t *testing.T) {
	testScheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(testScheme))
	assert.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		Build()

	for _, tc := range []struct {
		title              string
		toCreate           []set.Identifiable
		integrationService integrations.AtlasIntegrationsService
		expectedResult     workflow.Result
	}{
		{
			title:          "nil list does nothing",
			expectedResult: workflow.OK(),
		},
		{
			title: "successfully creates",
			toCreate: set.Difference(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Create(context.Background(), "project-id", "MICROSOFT_TEAMS", mock.AnythingOfType("integrations.Integration"), map[string]string{}).Return(nil)
				return service
			}(),
			expectedResult: workflow.OK(),
		},
		{
			title: "attempt to create errors",
			toCreate: set.Difference(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{}),
			integrationService: func() integrations.AtlasIntegrationsService {
				service := translation.NewAtlasIntegrationsServiceMock(t)
				service.EXPECT().Create(context.Background(), "project-id", "MICROSOFT_TEAMS", mock.AnythingOfType("integrations.Integration"), map[string]string{}).Return(errors.New("fake test error"))
				return service
			}(),
			expectedResult: workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Cannot create integration %v: %v", "MICROSOFT_TEAMS", errTest)),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
			}

			r := AtlasProjectReconciler{
				Client:              k8sClient,
				integrationsService: tc.integrationService,
			}
			result := r.createIntegrationsInAtlas(workflowCtx, testProjectID, tc.toCreate, testNamespace)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCheckIntegrationsReady(t *testing.T) {
	for _, tc := range []struct {
		title     string
		toCheck   [][]set.Identifiable
		requested []integrations.Integration
		expected  bool
	}{
		{
			title:    "nil list does nothing",
			expected: true,
		},

		{
			title:     "empty list does nothing",
			toCheck:   [][]set.Identifiable{},
			requested: []integrations.Integration{},
			expected:  true,
		},

		{
			title:     "when requested list differs in length it bails early",
			toCheck:   [][]set.Identifiable{},
			requested: []integrations.Integration{{}},
			expected:  false,
		},

		{
			title: "matching integrations are considered applied",
			toCheck: set.Intersection(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				}),
			requested: []integrations.Integration{{}},
			expected:  true,
		},

		{
			title: "different integrations are considered also applied",
			toCheck: set.Intersection(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/some-otherpath/some-othersecret",
							Enabled:                  true,
						},
					},
				}),
			requested: []integrations.Integration{{}},
			expected:  true,
		},

		{
			title: "matching integrations including prometheus are considered applied",
			toCheck: set.Intersection(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
					{
						Integration: project.Integration{
							Type:             "PROMETHEUS",
							UserName:         "prometheus",
							ServiceDiscovery: "http",
							Enabled:          true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
					{
						Integration: project.Integration{
							Type:             "PROMETHEUS",
							UserName:         "prometheus",
							ServiceDiscovery: "http",
							Enabled:          true,
						},
					},
				}),
			requested: []integrations.Integration{{}, {}},
			expected:  true,
		},

		{
			title: "matching integrations with a differing prometheus are considered different",
			toCheck: set.Intersection(
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							Name:                     testNamespace,
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
					{
						Integration: project.Integration{
							Type:             "PROMETHEUS",
							UserName:         "prometheus",
							ServiceDiscovery: "http",
							Enabled:          true,
						},
					},
				},
				[]integrations.Integration{
					{
						Integration: project.Integration{
							Type:                     "MICROSOFT_TEAMS",
							MicrosoftTeamsWebhookURL: "https://somehost/somepath/somesecret",
							Enabled:                  true,
						},
					},
					{
						Integration: project.Integration{
							Type:             "PROMETHEUS",
							UserName:         "zeus",
							ServiceDiscovery: "file",
							Enabled:          true,
						},
					},
				}),
			requested: []integrations.Integration{{}, {}},
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

func TestSyncPrometheusStatus(t *testing.T) {
	workflowCtx := &workflow.Context{
		Context:   context.Background(),
		Log:       zap.S(),
		SdkClient: admin.NewAPIClient(admin.NewConfiguration()),
	}

	proj := &akov2.AtlasProject{
		Status: status.AtlasProjectStatus{
			ID: "testid123",
		},
	}

	int := set.Intersection(
		[]integrations.Integration{
			{
				Integration: project.Integration{
					Type: "PROMETHEUS",
				},
			},
		},
		[]integrations.Integration{
			{
				Integration: project.Integration{
					Type: "PROMETHEUS",
				},
			},
		},
	)

	syncPrometheusStatus(workflowCtx, proj, int)

	for _, opt := range workflowCtx.StatusOptions() {
		if o, ok := opt.(status.AtlasProjectStatusOption); ok {
			o(&proj.Status)
		}
	}

	assert.NotNil(t, proj.Status.Prometheus)
	assert.Equal(t, `https://cloud.mongodb.com/prometheus/v1.0/groups/testid123/discovery`, proj.Status.Prometheus.DiscoveryURL)
}
