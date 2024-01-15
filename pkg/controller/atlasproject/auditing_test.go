package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestCanAuditingReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canAuditingReconcile(testWorkFlowContext(mongodbatlas.Client{}), false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canAuditingReconcile(testWorkFlowContext(mongodbatlas.Client{}), true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAuditingReconcile(testWorkFlowContext(atlasClient), true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when configuration is empty in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return nil, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAuditingReconcile(testWorkFlowContext(atlasClient), true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return &mongodbatlas.Auditing{
						Enabled:                   toptr.MakePtr(true),
						AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
						AuditAuthorizationSuccess: toptr.MakePtr(false),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Auditing: &mdbv1.Auditing{
					Enabled:                   true,
					AuditFilter:               `{"atype":"authenticate","param":{"db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					AuditAuthorizationSuccess: false,
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"auditing":{"auditFilter":"{\"atype\":\"authenticate\",\"param\":{\"user\":\"auditReadOnly\",\"db\":\"admin\",\"mechanism\":\"SCRAM-SHA-1\"}}","enabled":true,"auditAuthorizationSuccess":false}}`,
			},
		)
		result, err := canAuditingReconcile(testWorkFlowContext(atlasClient), true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return &mongodbatlas.Auditing{
						Enabled:                   toptr.MakePtr(true),
						AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
						AuditAuthorizationSuccess: toptr.MakePtr(false),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Auditing: &mdbv1.Auditing{
					Enabled:                   true,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					AuditAuthorizationSuccess: false,
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"auditing":{"auditFilter":"{\"atype\":\"authenticate\",\"param\":{\"user\":\"auditReadOnly\",\"db\":\"admin\",\"mechanism\":\"SCRAM-SHA-1\"}}","enabled":true,"auditAuthorizationSuccess":true}}`,
			},
		)
		result, err := canAuditingReconcile(testWorkFlowContext(atlasClient), true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile Auditing", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return &mongodbatlas.Auditing{
						Enabled:                   toptr.MakePtr(true),
						AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
						AuditAuthorizationSuccess: toptr.MakePtr(false),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Auditing: &mdbv1.Auditing{
					Enabled:                   true,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					AuditAuthorizationSuccess: true,
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"auditing":{"auditFilter":"{\"atype\":\"authenticate\",\"param\":{\"db\":\"admin\",\"mechanism\":\"SCRAM-SHA-1\"}}","enabled":true,"auditAuthorizationSuccess":true}}`,
			},
		)
		result, err := canAuditingReconcile(testWorkFlowContext(atlasClient), true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsureAuditing(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result := ensureAuditing(testWorkFlowContext(atlasClient), akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Auditing: &atlas.AuditingClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
					return &mongodbatlas.Auditing{
						Enabled:                   toptr.MakePtr(true),
						AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
						AuditAuthorizationSuccess: toptr.MakePtr(false),
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Auditing: &mdbv1.Auditing{
					Enabled:                   true,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					AuditAuthorizationSuccess: true,
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"auditing":{"auditFilter":"{\"atype\":\"authenticate\",\"param\":{\"db\":\"admin\",\"mechanism\":\"SCRAM-SHA-1\"}}","enabled":true,"auditAuthorizationSuccess":true}}`,
			},
		)
		result := ensureAuditing(testWorkFlowContext(atlasClient), akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}

func TestAuditingInSync(t *testing.T) {
	type args struct {
		atlas *mongodbatlas.Auditing
		spec  *mdbv1.Auditing
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Atlas and Operator Auditing are empty",
			args: args{
				atlas: nil,
				spec:  nil,
			},
			want: true,
		},
		{
			name: "Atlas Auditing is empty and Operator doesn't",
			args: args{
				atlas: nil,
				spec:  &mdbv1.Auditing{Enabled: true},
			},
			want: false,
		},
		{
			name: "Operator Auditing is empty and Atlas doesn't",
			args: args{
				atlas: &mongodbatlas.Auditing{Enabled: toptr.MakePtr(true)},
				spec:  nil,
			},
			want: false,
		},
		{
			name: "Operator Auditing has different config from Atlas",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &mdbv1.Auditing{
					AuditAuthorizationSuccess: true,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					Enabled:                   true,
				},
			},
			want: false,
		},
		{
			name: "Operator Auditing has different config filter from Atlas",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &mdbv1.Auditing{
					AuditAuthorizationSuccess: false,
					AuditFilter:               `{"atype":"authenticate","param":{"db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					Enabled:                   true,
				},
			},
			want: false,
		},
		{
			name: "Operator Auditing are Equal",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &mdbv1.Auditing{
					AuditAuthorizationSuccess: false,
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					Enabled:                   true,
				},
			},
			want: true,
		},
		{
			name: "Operator Auditing are Equal when filter has newline in the end",
			args: args{
				atlas: &mongodbatlas.Auditing{
					AuditAuthorizationSuccess: toptr.MakePtr(false),
					AuditFilter:               `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
					ConfigurationType:         "ReadOnly",
					Enabled:                   toptr.MakePtr(true),
				},
				spec: &mdbv1.Auditing{
					AuditAuthorizationSuccess: false,
					AuditFilter:               "{\"atype\":\"authenticate\",\"param\":{\"user\":\"auditReadOnly\",\"db\":\"admin\",\"mechanism\":\"SCRAM-SHA-1\"}}\n",
					Enabled:                   true,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, auditingInSync(tt.args.atlas, tt.args.spec), "auditingInSync(%v, %v)", tt.args.atlas, tt.args.spec)
		})
	}
}

func testWorkFlowContext(client mongodbatlas.Client) *workflow.Context {
	return &workflow.Context{Client: &client, Context: context.Background()}
}
