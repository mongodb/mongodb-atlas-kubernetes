package atlasdeployment

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	fakeInstanceName = "fake-instance-name"
)

func TestCanReconcileServerlessPrivateEndpoints(t *testing.T) {
	t.Run("when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canServerlessPrivateEndpointsReconcile(
			&workflow.Context{},
			false,
			"fake-project-id-wont-be-checked",
			&akov2.AtlasDeployment{})

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("when protected but there is no Atlas Serverless Endpoint configured", func(t *testing.T) {
		ctx := context.Background()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return []mongodbatlas.ServerlessPrivateEndpointConnection{}, nil, nil
				},
			},
		}
		deployment := sampleServerlessDeployment()
		workflowCtx := workflow.Context{Client: &client, Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("when protected but configs match", func(t *testing.T) {
		ctx := context.Background()
		endpointsConfig := sampleAtlasSPEConfig()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return endpointsConfig, nil, nil
				},
			},
		}
		deployment := sampleAnnotatedServerlessDeployment(endpointsFrom(endpointsConfig))
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("when protected but configs match, even with different order", func(t *testing.T) {
		ctx := context.Background()
		endpointsConfig := sampleAtlasSPEConfig()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return endpointsConfig, nil, nil
				},
			},
		}
		deployment := sampleAnnotatedServerlessDeployment(reverse(endpointsFrom(endpointsConfig)))
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("when protected but only old configs matches", func(t *testing.T) {
		ctx := context.Background()
		endpointsConfig := sampleAtlasSPEConfig()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return endpointsConfig, nil, nil
				},
			},
		}
		deployment := sampleAnnotatedServerlessDeployment(endpointsFrom(endpointsConfig))
		// remove all PEs in the current desired setup
		deployment.Spec.ServerlessSpec.PrivateEndpoints = []akov2.ServerlessPrivateEndpoint{}
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.True(t, result)
	})
}

func TestCannotReconcileServerlessPrivateEndpoints(t *testing.T) {
	t.Run("when configs do not match", func(t *testing.T) {
		ctx := context.Background()
		endpointsConfig := sampleAtlasSPEConfig()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return endpointsConfig, nil, nil
				},
			},
		}
		endpoints := endpointsFrom(endpointsConfig)
		endpoints[0].Name = "non-matching-fake-name"
		deployment := sampleAnnotatedServerlessDeployment(endpoints)
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("when ownership cannot be assured (empty prior config)", func(t *testing.T) {
		ctx := context.Background()
		endpointsConfig := sampleAtlasSPEConfig()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return endpointsConfig, nil, nil
				},
			},
		}
		deployment := sampleServerlessDeployment()
		deployment.Annotations = map[string]string{
			customresource.AnnotationLastAppliedConfiguration: "{}",
		}
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("when ownership cannot be assured (unset prior config)", func(t *testing.T) {
		ctx := context.Background()
		endpointsConfig := sampleAtlasSPEConfig()
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return endpointsConfig, nil, nil
				},
			},
		}
		deployment := sampleServerlessDeployment()
		deployment.Annotations = map[string]string{}
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.NoError(t, err)
		assert.False(t, result)
	})
}

func TestCanReconcileServerlessPrivateEndpointsFail(t *testing.T) {
	t.Run("when the old config is not a proper JSON", func(t *testing.T) {
		ctx := context.Background()
		client := mongodbatlas.Client{}
		deployment := sampleServerlessDeployment()
		deployment.Annotations = map[string]string{
			customresource.AnnotationLastAppliedConfiguration: "{",
		}
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.False(t, result)
		var aJSONError *json.SyntaxError
		assert.ErrorAs(t, err, &aJSONError)
	})

	t.Run("when list fails in Atlas", func(t *testing.T) {
		ctx := context.Background()
		fakeError := fmt.Errorf("fake error from Atlas")
		client := mongodbatlas.Client{
			ServerlessPrivateEndpoints: ServerlessPrivateEndpointClientMock{
				ListFn: func(groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
					return nil, nil, fakeError
				},
			},
		}
		deployment := sampleServerlessDeployment()
		workflowCtx := workflow.Context{Client: &client, Log: debugLogger(t), Context: ctx}

		result, err := canServerlessPrivateEndpointsReconcile(&workflowCtx, true, fakeProjectID, deployment)

		require.False(t, result)
		assert.ErrorIs(t, err, fakeError)
	})
}

func sampleServerlessDeployment() *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		Spec: akov2.AtlasDeploymentSpec{
			ServerlessSpec: &akov2.ServerlessSpec{Name: fakeInstanceName},
		},
	}
}

func sampleAnnotatedServerlessDeployment(endpoints []akov2.ServerlessPrivateEndpoint) *akov2.AtlasDeployment {
	deployment := &akov2.AtlasDeployment{
		Spec: akov2.AtlasDeploymentSpec{ServerlessSpec: &akov2.ServerlessSpec{
			Name:             fakeInstanceName,
			PrivateEndpoints: endpoints,
		}},
	}
	deployment.Annotations = map[string]string{
		customresource.AnnotationLastAppliedConfiguration: jsonize(deployment.Spec),
	}
	return deployment
}

func sampleAtlasSPEConfig() []mongodbatlas.ServerlessPrivateEndpointConnection {
	return []mongodbatlas.ServerlessPrivateEndpointConnection{
		{
			ID:                      "fake-id-1",
			CloudProviderEndpointID: "opaque-cloud-fake-id-1",
			Comment:                 "fake-name-1",
			Status:                  SPEStatusAvailable,
			ProviderName:            "AWS",
		},
		{
			ID:                       "fake-id-2",
			CloudProviderEndpointID:  "opaque-cloud-fake-id-2",
			Comment:                  "fake-name-2",
			Status:                   SPEStatusAvailable,
			ProviderName:             "Azure",
			PrivateEndpointIPAddress: "11.11.10.0",
		},
	}
}

func endpointsFrom(configs []mongodbatlas.ServerlessPrivateEndpointConnection) []akov2.ServerlessPrivateEndpoint {
	endpoints := []akov2.ServerlessPrivateEndpoint{}
	for _, cfg := range configs {
		endpoints = append(endpoints, akov2.ServerlessPrivateEndpoint{
			Name:                     cfg.Comment,
			CloudProviderEndpointID:  cfg.CloudProviderEndpointID,
			PrivateEndpointIPAddress: cfg.PrivateEndpointIPAddress,
		})
	}
	return endpoints
}

func reverse(endpoints []akov2.ServerlessPrivateEndpoint) []akov2.ServerlessPrivateEndpoint {
	reversed := make([]akov2.ServerlessPrivateEndpoint, 0, len(endpoints))
	for i := len(endpoints) - 1; i >= 0; i-- {
		reversed = append(reversed, endpoints[i])
	}
	return reversed
}

func jsonize(obj any) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}

func debugLogger(t *testing.T) *zap.SugaredLogger {
	t.Helper()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	return logger.Sugar()
}
