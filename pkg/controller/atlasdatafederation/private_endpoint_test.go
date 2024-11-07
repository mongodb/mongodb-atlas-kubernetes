package atlasdatafederation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestEnsurePrivateEndpoints(t *testing.T) {
	for _, tc := range []struct {
		name           string
		project        *akov2.AtlasProject
		dataFederation *akov2.AtlasDataFederation
		service        func(mock *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService

		wantOK         bool
		wantConditions []api.Condition
	}{
		{
			name: "empty data federation",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{}, nil)
				return m
			},
			project:        &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{},
			wantOK:         true,
			wantConditions: []api.Condition{},
		},
		{
			name: "error in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return(nil, errors.New("atlas error"))
				return m
			},
			project:        &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{},
			wantOK:         false,
			wantConditions: []api.Condition{
				{
					Type:   "DataFederationPrivateEndpointsReady",
					Status: "False", Reason: "InternalError", Message: "atlas error",
				},
			},
		},
		{
			name: "create entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return(nil, nil)
				m.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
				return m
			},
			project: &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{
						{EndpointID: "123", Provider: "foo", Type: "some"},
					},
				},
			},
			wantOK: true,
			wantConditions: []api.Condition{
				{Type: "DataFederationPrivateEndpointsReady", Status: "True"},
			},
		},
		{
			name: "delete and update entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"}},
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "456", Provider: "bar", Type: "some"}},
				}, nil)
				m.EXPECT().Delete(mock.Anything, &datafederation.DatafederationPrivateEndpointEntry{
					DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"},
				}).Return(nil)
				m.EXPECT().Create(mock.Anything, &datafederation.DatafederationPrivateEndpointEntry{
					DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "CHANGE", Type: "some"},
				}).Return(nil)
				return m
			},
			project: &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{
						{EndpointID: "123", Provider: "CHANGE", Type: "some"},
					},
				},
			},
			wantOK: true,
			wantConditions: []api.Condition{
				{Type: "DataFederationPrivateEndpointsReady", Status: "True"},
			},
		},
		{
			name: "do not delete untracked entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"}},
				}, nil)
				return m
			},
			project:        &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{},
			wantOK:         true,
			wantConditions: []api.Condition{},
		},
		{
			name: "delete tracked entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"}},
				}, nil)
				m.EXPECT().Delete(mock.Anything, &datafederation.DatafederationPrivateEndpointEntry{
					DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"},
				}).Return(nil)
				return m
			},
			project: &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"mongodb.com/last-applied-configuration": `
{
  "privateEndpoints": [
   {
    "endpointId": "123",
    "provider": "foo",
    "type": "some"
   }
  ]
 }`,
					},
				},
			},
			wantOK:         true,
			wantConditions: []api.Condition{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			reconciler := &AtlasDataFederationReconciler{}
			result := reconciler.ensurePrivateEndpoints(ctx, tc.service(translation.NewDatafederationPrivateEndpointServiceMock(t)), tc.project, tc.dataFederation)
			require.Equal(t, tc.wantOK, result.IsOk())

			gotConditions := ctx.Conditions()
			for i := range ctx.Conditions() {
				gotConditions[i].LastTransitionTime = metav1.Time{}
			}
			require.Equal(t, tc.wantConditions, gotConditions)
		})
	}
}
