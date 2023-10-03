package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

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

func TestAreIntegrationsEqual(t *testing.T) {
	atlas := aliasThirdPartyIntegration{
		Type:   "DATADOG",
		APIKey: "****************************4e6f",
		Region: "EU",
	}
	spec := aliasThirdPartyIntegration{
		Type:   "DATADOG",
		APIKey: "actual-valid-id*************4e6f",
		Region: "EU",
	}

	areEqual := AreIntegrationsEqual(&atlas, &spec)
	assert.True(t, areEqual, "Identical objects should be equal")

	spec.APIKey = "non-equal-id************1234"
	areEqual = AreIntegrationsEqual(&atlas, &spec)
	assert.False(t, areEqual, "Should fail if the last 4 characters of APIKey do not match")
}

func TestCanIntegrationsReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  mongodbatlas.Client{},
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		workflowCtx := &workflow.Context{
			Client:  mongodbatlas.Client{},
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Integrations: &atlas.ThirdPartyIntegrationsClientMock{
				ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Integrations: &atlas.ThirdPartyIntegrationsClientMock{
				ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
					return &mongodbatlas.ThirdPartyIntegrations{TotalCount: 0}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Integrations: &atlas.ThirdPartyIntegrationsClientMock{
				ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
					return &mongodbatlas.ThirdPartyIntegrations{
						Results: []*mongodbatlas.ThirdPartyIntegration{
							{
								Type:   "DATADOG",
								Region: "EU",
								APIKey: "my-api-key",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Integrations: []project.Integration{
					{
						Type: "DATADOG",
						APIKeyRef: common.ResourceRefNamespaced{
							Name:      "datadog-secret",
							Namespace: "project-namespace",
						},
						Region: "US",
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"integrations\":[{\"type\":\"DATADOG\",\"apiKeyRef\":{\"name\":\"datadog-secret\",\"namespace\":\"project-namespace\"},\"region\":\"US\"}]}",
			},
		)
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Integrations: &atlas.ThirdPartyIntegrationsClientMock{
				ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
					return &mongodbatlas.ThirdPartyIntegrations{
						Results: []*mongodbatlas.ThirdPartyIntegration{
							{
								Type:   "DATADOG",
								Region: "EU",
								APIKey: "my-api-key",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Integrations: []project.Integration{
					{
						Type: "DATADOG",
						APIKeyRef: common.ResourceRefNamespaced{
							Name:      "datadog-secret",
							Namespace: "project-namespace",
						},
						Region: "EU",
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"integrations\":[{\"type\":\"DATADOG\",\"apiKeyRef\":{\"name\":\"datadog-secret\",\"namespace\":\"project-namespace\"},\"region\":\"US\"}]}",
			},
		)
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile Integrations", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Integrations: &atlas.ThirdPartyIntegrationsClientMock{
				ListFunc: func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
					return &mongodbatlas.ThirdPartyIntegrations{
						Results: []*mongodbatlas.ThirdPartyIntegration{
							{
								Type:   "DATADOG",
								Region: "EU",
								APIKey: "my-api-key",
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Integrations: []project.Integration{
					{
						Type: "PAGER_DUTY",
						ServiceKeyRef: common.ResourceRefNamespaced{
							Name:      "pager-duty-secret",
							Namespace: "project-namespace",
						},
						Region: "EU",
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"integrations\":[{\"type\":\"PAGER_DUTY\",\"serviceKeyRef\":{\"name\":\"pager-duty-secret\",\"namespace\":\"project-namespace\"},\"region\":\"EU\"}]}",
			},
		)
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		result, err := canIntegrationsReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}
