package atlasdatafederation

import (
	"context"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap/zaptest"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type DataFederationMock struct {
	mongodbatlas.DataFederationService
}

func (m *DataFederationMock) Delete(context.Context, string, string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func TestDeleteConnectionSecrets(t *testing.T) {
	for _, tc := range []struct {
		name              string
		atlasProject      *akov2.AtlasProject
		dataFederation    *akov2.AtlasDataFederation
		connectionSecrets []*corev1.Secret

		wantDataFederation *akov2.AtlasDataFederation
		wantSecrets        []corev1.Secret
		wantResult         workflow.Result
	}{
		{
			name: "no finalizer",
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"},
			},
			wantResult: workflow.OK(),
		},
		{
			name: "finalizer set and deletion protection is enabled",
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers:  []string{customresource.FinalizerLabel},
					Annotations: map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep},
				},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
					Annotations: map[string]string{customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep},
				},
			},
			wantResult: workflow.OK(),
		},
		{
			name: "federation object without secrets",
			atlasProject: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "fooProject", Namespace: "bar"},
			},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.DataFederationSpec{Project: common.ResourceRefNamespaced{Name: "fooProject"}},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
				},
				Spec: akov2.DataFederationSpec{Project: common.ResourceRefNamespaced{Name: "fooProject"}},
			},
			wantResult: workflow.OK(),
		},
		{
			name: "federation object without secrets",
			atlasProject: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "fooProject", Namespace: "bar"},
			},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{Name: "fooProject"},
				},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
				},
				Spec: akov2.DataFederationSpec{Project: common.ResourceRefNamespaced{Name: "fooProject"}},
			},
			wantResult: workflow.OK(),
		},
		{
			name: "federation object with secrets",
			atlasProject: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "fooProject", Namespace: "bar"},
				Status:     status.AtlasProjectStatus{ID: "123"},
			},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.DataFederationSpec{
					Name:    "data-federation-name",
					Project: common.ResourceRefNamespaced{Name: "fooProject"},
				},
			},
			wantDataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo", Namespace: "bar",
					// Finalizers removed
				},
				Spec: akov2.DataFederationSpec{
					Name:    "data-federation-name",
					Project: common.ResourceRefNamespaced{Name: "fooProject"},
				},
			},
			connectionSecrets: []*corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "fooSecret", Namespace: "bar",
						Labels: map[string]string{
							connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
							connectionsecret.ProjectLabelKey: "123",
							connectionsecret.ClusterLabelKey: "data-federation-name",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keepSecret", Namespace: "bar",
						Labels: map[string]string{
							connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
							connectionsecret.ProjectLabelKey: "123",
							connectionsecret.ClusterLabelKey: "some-cluster",
						},
					},
				},
			},
			wantSecrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "keepSecret", Namespace: "bar",
						Labels: map[string]string{
							connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
							connectionsecret.ProjectLabelKey: "123",
							connectionsecret.ClusterLabelKey: "some-cluster",
						},
					},
				},
			},
			wantResult: workflow.OK(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			var objects []client.Object
			if tc.atlasProject != nil {
				objects = append(objects, tc.atlasProject)
			}
			if tc.dataFederation != nil {
				objects = append(objects, tc.dataFederation)
			}
			for _, s := range tc.connectionSecrets {
				objects = append(objects, s)
			}
			scheme := runtime.NewScheme()
			utilruntime.Must(corev1.AddToScheme(scheme))
			utilruntime.Must(akov2.AddToScheme(scheme))
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objects...).
				Build()
			project := &akov2.AtlasProject{}
			atlasClient := &mongodbatlas.Client{
				DataFederation: &DataFederationMock{},
			}

			r := &AtlasDataFederationReconciler{
				Client: fakeClient,
				Log:    logger,
			}
			gotResult := r.handleDelete(ctx, logger, tc.dataFederation, project, atlasClient)
			assert.Equal(t, tc.wantResult, gotResult)

			gotDataFederation := &akov2.AtlasDataFederation{}
			err := fakeClient.Get(ctx.Context, client.ObjectKeyFromObject(tc.dataFederation), gotDataFederation)
			assert.NoError(t, err)
			gotDataFederation.ResourceVersion = ""
			assert.Equal(t, tc.wantDataFederation, gotDataFederation)

			gotSecrets := &corev1.SecretList{}
			err = fakeClient.List(ctx.Context, gotSecrets)
			for i := range gotSecrets.Items {
				gotSecrets.Items[i].ResourceVersion = ""
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantSecrets, gotSecrets.Items)
		})
	}
}
