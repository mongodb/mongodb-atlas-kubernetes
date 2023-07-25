package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

type privateEndpointClient struct {
	ListFunc func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error)
}

func (c *privateEndpointClient) Create(_ context.Context, _ string, _ *mongodbatlas.PrivateEndpointConnection) (*mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *privateEndpointClient) Get(_ context.Context, _ string, _ string, _ string) (*mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *privateEndpointClient) List(_ context.Context, projectID string, providerName string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
	return c.ListFunc(projectID, providerName)
}

func (c *privateEndpointClient) Delete(_ context.Context, _ string, _ string, _ string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func (c *privateEndpointClient) AddOnePrivateEndpoint(_ context.Context, _ string, _ string, _ string, _ *mongodbatlas.InterfaceEndpointConnection) (*mongodbatlas.InterfaceEndpointConnection, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *privateEndpointClient) GetOnePrivateEndpoint(_ context.Context, _ string, _ string, _ string, _ string) (*mongodbatlas.InterfaceEndpointConnection, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *privateEndpointClient) DeleteOnePrivateEndpoint(_ context.Context, _ string, _ string, _ string, _ string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func (c *privateEndpointClient) UpdateRegionalizedPrivateEndpointSetting(_ context.Context, _ string, _ bool) (*mongodbatlas.RegionalizedPrivateEndpointSetting, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *privateEndpointClient) GetRegionalizedPrivateEndpointSetting(_ context.Context, _ string) (*mongodbatlas.RegionalizedPrivateEndpointSetting, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func TestGetEndpointsNotInAtlas(t *testing.T) {
	const region1 = "SOME_REGION"
	const region2 = "OTHER_REGION"
	specPEs := []mdbv1.PrivateEndpoint{
		{
			Provider: provider.ProviderAWS,
			Region:   region1,
		},
		{
			Provider: provider.ProviderAWS,
			Region:   region1,
		},
		{
			Provider: provider.ProviderAWS,
			Region:   region2,
		},
	}
	atlasPEs := []atlasPE{}
	uniqueItems, itemCounts := getEndpointsNotInAtlas(specPEs, atlasPEs)
	assert.Equalf(t, 2, len(uniqueItems), "getEndpointsNotInAtlas should remove a duplicate PE Service")
	assert.NotEqualf(t, uniqueItems[0].Region, uniqueItems[1].Region, "getEndpointsNotInAtlas should return unique PEs")
	assert.Equalf(t, len(uniqueItems), len(itemCounts), "item counts should have the same length as items")
	assert.Equalf(t, 3, itemCounts[0]+itemCounts[1], "item counts should sum up to the actual value of spec endpoints")

	atlasPEs = append(atlasPEs, atlasPE{
		ProviderName: string(provider.ProviderAWS),
		RegionName:   region1,
	})

	uniqueItems, _ = getEndpointsNotInAtlas(specPEs, atlasPEs)
	assert.Equalf(t, len(uniqueItems), 1, "getEndpointsNotInAtlas should remove both PE Service copies if there is one in Atlas")
}

func TestGetEndpointsNotInSpec(t *testing.T) {
	const region1 = "SOME_REGION"
	const region2 = "OTHER_REGION"
	specPEs := []mdbv1.PrivateEndpoint{
		{
			Provider: provider.ProviderAWS,
			Region:   region1,
		},
		{
			Provider: provider.ProviderAWS,
			Region:   region1,
		},
	}
	atlasPEs := []atlasPE{
		{
			ProviderName: string(provider.ProviderAWS),
			RegionName:   region1,
		},
		{
			ProviderName: string(provider.ProviderAWS),
			RegionName:   region1,
		},
	}

	uniqueItems := getEndpointsNotInSpec(specPEs, atlasPEs)
	assert.Equalf(t, 0, len(uniqueItems), "getEndpointsNotInSpec should not return anything if PEs are in spec")

	atlasPEs = append(atlasPEs, atlasPE{
		ProviderName: string(provider.ProviderAWS),
		Region:       region2,
	})
	uniqueItems = getEndpointsNotInSpec(specPEs, atlasPEs)
	assert.Equalf(t, 1, len(uniqueItems), "getEndpointsNotInSpec should get a spec item")
}

func TestCanPrivateEndpointReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canPrivateEndpointReconcile(mongodbatlas.Client{}, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canPrivateEndpointReconcile(mongodbatlas.Client{}, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canPrivateEndpointReconcile(atlasClient, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					return []mongodbatlas.PrivateEndpointConnection{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canPrivateEndpointReconcile(atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					if providerName == "AWS" {
						return []mongodbatlas.PrivateEndpointConnection{
							{
								ID:           "123456",
								ProviderName: "AWS",
								Region:       "eu-west-2",
								RegionName:   "eu-west-2",
							},
						}, nil, nil
					}

					return []mongodbatlas.PrivateEndpointConnection{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-2",
					},
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"privateEndpoints\":[{\"provider\":\"AWS\",\"region\":\"eu-west-2\"}]}"})
		result, err := canPrivateEndpointReconcile(atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					if providerName == "AWS" {
						return []mongodbatlas.PrivateEndpointConnection{
							{
								ID:           "123456",
								ProviderName: "AWS",
								Region:       "eu-west-2",
								RegionName:   "eu-west-2",
							}, {
								ID:           "654321",
								ProviderName: "AWS",
								Region:       "eu-west-1",
								RegionName:   "eu-west-1",
							},
						}, nil, nil
					}

					return []mongodbatlas.PrivateEndpointConnection{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-2",
					},
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"privateEndpoints\":[{\"provider\":\"AWS\",\"region\":\"eu-west-2\"}]}"})
		result, err := canPrivateEndpointReconcile(atlasClient, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile private endpoints", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					if providerName == "AWS" {
						return []mongodbatlas.PrivateEndpointConnection{
							{
								ID:           "123456",
								ProviderName: "AWS",
								Region:       "eu-west-2",
								RegionName:   "eu-west-2",
							}, {
								ID:           "654321",
								ProviderName: "AWS",
								Region:       "us-west-1",
								RegionName:   "us-west-1",
							},
						}, nil, nil
					}

					return []mongodbatlas.PrivateEndpointConnection{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-2",
					},
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"privateEndpoints\":[{\"provider\":\"AWS\",\"region\":\"eu-west-2\"}]}"})
		result, err := canPrivateEndpointReconcile(atlasClient, true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsurePrivateEndpoint(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensurePrivateEndpoint(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "the operator could not validate ownership for deletion protection"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			PrivateEndpoints: &privateEndpointClient{
				ListFunc: func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
					if providerName == "AWS" {
						return []mongodbatlas.PrivateEndpointConnection{
							{
								ID:           "123456",
								ProviderName: "AWS",
								Region:       "eu-west-2",
								RegionName:   "eu-west-2",
							}, {
								ID:           "654321",
								ProviderName: "AWS",
								Region:       "us-west-1",
								RegionName:   "us-west-1",
							},
						}, nil, nil
					}

					return []mongodbatlas.PrivateEndpointConnection{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-2",
					},
					{
						Provider: provider.ProviderAWS,
						Region:   "eu-west-1",
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{\"privateEndpoints\":[{\"provider\":\"AWS\",\"region\":\"eu-west-2\"}]}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		result := ensurePrivateEndpoint(workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Private Endpoint(s) due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}
