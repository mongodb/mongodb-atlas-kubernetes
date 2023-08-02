package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap/zaptest"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
)

type alertConfigurationClient struct {
	ListFunc func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
}

func (c *alertConfigurationClient) Create(_ context.Context, _ string, _ *mongodbatlas.AlertConfiguration) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *alertConfigurationClient) EnableAnAlertConfig(_ context.Context, _ string, _ string, _ *bool) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *alertConfigurationClient) GetAnAlertConfig(_ context.Context, _ string, _ string) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *alertConfigurationClient) GetOpenAlertsConfig(_ context.Context, _ string, _ string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *alertConfigurationClient) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	return c.ListFunc(projectID)
}

func (c *alertConfigurationClient) ListMatcherFields(_ context.Context) ([]string, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *alertConfigurationClient) Update(_ context.Context, _ string, _ string, _ *mongodbatlas.AlertConfiguration) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *alertConfigurationClient) Delete(_ context.Context, _ string, _ string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func TestCanAlertConfigurationReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canAlertConfigurationReconcile(context.TODO(), mongodbatlas.Client{}, false, &mdbv1.AtlasProject{}, zaptest.NewLogger(t).Sugar())
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canAlertConfigurationReconcile(context.TODO(), mongodbatlas.Client{}, true, akoProject, zaptest.NewLogger(t).Sugar())
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAlertConfigurationReconcile(context.TODO(), atlasClient, true, akoProject, zaptest.NewLogger(t).Sugar())

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when there are no items in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return []mongodbatlas.AlertConfiguration{}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAlertConfigurationReconcile(context.TODO(), atlasClient, true, akoProject, zaptest.NewLogger(t).Sugar())

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return []mongodbatlas.AlertConfiguration{
						{
							EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
							Enabled:       toptr.MakePtr(true),
							Threshold: &mongodbatlas.Threshold{
								Operator:  "LESS_THAN",
								Threshold: 1,
								Units:     "HOURS",
							},
							Notifications: []mongodbatlas.Notification{
								{
									IntervalMin:  5,
									DelayMin:     toptr.MakePtr(5),
									EmailEnabled: toptr.MakePtr(true),
									SMSEnabled:   toptr.MakePtr(false),
									Roles: []string{
										"GROUP_OWNER",
									},
									TypeName: "GROUP",
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
					{
						EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
						Enabled:       true,
						Threshold: &mdbv1.Threshold{
							Operator:  "LESS_THAN",
							Threshold: "1",
							Units:     "HOURS",
						},
						Notifications: []mdbv1.Notification{
							{
								IntervalMin:  6,
								DelayMin:     toptr.MakePtr(5),
								EmailEnabled: toptr.MakePtr(true),
								SMSEnabled:   toptr.MakePtr(false),
								Roles: []string{
									"GROUP_OWNER",
								},
								TypeName: "GROUP",
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"alertConfigurations\":[{\"enabled\":true,\"eventTypeName\":\"REPLICATION_OPLOG_WINDOW_RUNNING_OUT\",\"threshold\":{\"operator\":\"LESS_THAN\",\"units\":\"HOURS\",\"threshold\":\"1\"},\"notifications\":[{\"apiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"datadogAPIKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"delayMin\":5,\"emailEnabled\":true,\"flowdockApiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"intervalMin\":5,\"opsGenieApiKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"serviceKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"smsEnabled\":false,\"typeName\":\"GROUP\",\"victorOpsSecretRef\":{\"name\":\"\",\"namespace\":\"\"},\"roles\":[\"GROUP_OWNER\"]}]}]}",
			},
		)
		result, err := canAlertConfigurationReconcile(context.TODO(), atlasClient, true, akoProject, zaptest.NewLogger(t).Sugar())

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return []mongodbatlas.AlertConfiguration{
						{
							EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
							Enabled:       toptr.MakePtr(true),
							Threshold: &mongodbatlas.Threshold{
								Operator:  "LESS_THAN",
								Threshold: 1,
								Units:     "HOURS",
							},
							Notifications: []mongodbatlas.Notification{
								{
									IntervalMin:  5,
									DelayMin:     toptr.MakePtr(5),
									EmailEnabled: toptr.MakePtr(true),
									SMSEnabled:   toptr.MakePtr(false),
									Roles: []string{
										"GROUP_OWNER",
									},
									TypeName: "GROUP",
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
					{
						EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
						Enabled:       true,
						Threshold: &mdbv1.Threshold{
							Operator:  "LESS_THAN",
							Threshold: "1",
							Units:     "HOURS",
						},
						Notifications: []mdbv1.Notification{
							{
								IntervalMin:  5,
								DelayMin:     toptr.MakePtr(5),
								EmailEnabled: toptr.MakePtr(true),
								SMSEnabled:   toptr.MakePtr(false),
								Roles: []string{
									"GROUP_OWNER",
								},
								TypeName: "GROUP",
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"alertConfigurations\":[{\"enabled\":true,\"eventTypeName\":\"REPLICATION_OPLOG_WINDOW_RUNNING_OUT\",\"threshold\":{\"operator\":\"LESS_THAN\",\"units\":\"HOURS\",\"threshold\":\"1\"},\"notifications\":[{\"apiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"datadogAPIKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"delayMin\":5,\"emailEnabled\":true,\"flowdockApiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"intervalMin\":6,\"opsGenieApiKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"serviceKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"smsEnabled\":false,\"typeName\":\"GROUP\",\"victorOpsSecretRef\":{\"name\":\"\",\"namespace\":\"\"},\"roles\":[\"GROUP_OWNER\"]}]}]}",
			},
		)
		result, err := canAlertConfigurationReconcile(context.TODO(), atlasClient, true, akoProject, zaptest.NewLogger(t).Sugar())

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return []mongodbatlas.AlertConfiguration{
						{
							EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
							Enabled:       toptr.MakePtr(true),
							Threshold: &mongodbatlas.Threshold{
								Operator:  "LESS_THAN",
								Threshold: 1,
								Units:     "HOURS",
							},
							Notifications: []mongodbatlas.Notification{
								{
									IntervalMin:  5,
									DelayMin:     toptr.MakePtr(5),
									EmailEnabled: toptr.MakePtr(true),
									SMSEnabled:   toptr.MakePtr(false),
									Roles: []string{
										"GROUP_OWNER",
									},
									TypeName: "GROUP",
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
					{
						EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
						Enabled:       true,
						Threshold: &mdbv1.Threshold{
							Operator:  "LESS_THAN",
							Threshold: "1",
							Units:     "HOURS",
						},
						Notifications: []mdbv1.Notification{
							{
								IntervalMin:  6,
								DelayMin:     toptr.MakePtr(5),
								EmailEnabled: toptr.MakePtr(true),
								SMSEnabled:   toptr.MakePtr(false),
								Roles: []string{
									"GROUP_OWNER",
								},
								TypeName: "GROUP",
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"alertConfigurations\":[{\"enabled\":true,\"eventTypeName\":\"REPLICATION_OPLOG_WINDOW_RUNNING_OUT\",\"threshold\":{\"operator\":\"LESS_THAN\",\"units\":\"HOURS\",\"threshold\":\"1\"},\"notifications\":[{\"apiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"datadogAPIKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"delayMin\":5,\"emailEnabled\":true,\"flowdockApiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"intervalMin\":4,\"opsGenieApiKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"serviceKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"smsEnabled\":false,\"typeName\":\"GROUP\",\"victorOpsSecretRef\":{\"name\":\"\",\"namespace\":\"\"},\"roles\":[\"GROUP_OWNER\"]}]}]}",
			},
		)
		result, err := canAlertConfigurationReconcile(context.TODO(), atlasClient, true, akoProject, zaptest.NewLogger(t).Sugar())

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsureAlertConfigurations(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client: atlasClient,
		}
		reconciler := &AtlasProjectReconciler{}
		result := reconciler.ensureAlertConfigurations(context.TODO(), workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			AlertConfigurations: &alertConfigurationClient{
				ListFunc: func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
					return []mongodbatlas.AlertConfiguration{
						{
							EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
							Enabled:       toptr.MakePtr(true),
							Threshold: &mongodbatlas.Threshold{
								Operator:  "LESS_THAN",
								Threshold: 1,
								Units:     "HOURS",
							},
							Notifications: []mongodbatlas.Notification{
								{
									IntervalMin:  5,
									DelayMin:     toptr.MakePtr(5),
									EmailEnabled: toptr.MakePtr(true),
									SMSEnabled:   toptr.MakePtr(false),
									Roles: []string{
										"GROUP_OWNER",
									},
									TypeName: "GROUP",
								},
							},
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
					{
						EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
						Enabled:       true,
						Threshold: &mdbv1.Threshold{
							Operator:  "LESS_THAN",
							Threshold: "1",
							Units:     "HOURS",
						},
						Notifications: []mdbv1.Notification{
							{
								IntervalMin:  6,
								DelayMin:     toptr.MakePtr(5),
								EmailEnabled: toptr.MakePtr(true),
								SMSEnabled:   toptr.MakePtr(false),
								Roles: []string{
									"GROUP_OWNER",
								},
								TypeName: "GROUP",
							},
						},
					},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"alertConfigurations\":[{\"enabled\":true,\"eventTypeName\":\"REPLICATION_OPLOG_WINDOW_RUNNING_OUT\",\"threshold\":{\"operator\":\"LESS_THAN\",\"units\":\"HOURS\",\"threshold\":\"1\"},\"notifications\":[{\"apiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"datadogAPIKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"delayMin\":5,\"emailEnabled\":true,\"flowdockApiTokenRef\":{\"name\":\"\",\"namespace\":\"\"},\"intervalMin\":4,\"opsGenieApiKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"serviceKeyRef\":{\"name\":\"\",\"namespace\":\"\"},\"smsEnabled\":false,\"typeName\":\"GROUP\",\"victorOpsSecretRef\":{\"name\":\"\",\"namespace\":\"\"},\"roles\":[\"GROUP_OWNER\"]}]}]}",
			},
		)
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		reconciler := &AtlasProjectReconciler{}
		result := reconciler.ensureAlertConfigurations(context.TODO(), workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Alert Configuration due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}
