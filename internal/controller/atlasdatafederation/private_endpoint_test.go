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

package atlasdatafederation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
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
			name: "failed when last applied configuration annotation has wrong data",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return(nil, nil)

				return m
			},
			project: &akov2.AtlasProject{},
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"mongodb.com/last-applied-configuration": "{wrongJson",
					},
				},
			},
			wantOK: false,
			wantConditions: []api.Condition{
				{
					Type:    "DataFederationPrivateEndpointsReady",
					Status:  "False",
					Reason:  "InternalError",
					Message: "error reading data federation from last applied annotation: invalid character 'w' looking for beginning of object key string",
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
			name: "failed to create entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return(nil, nil)
				m.EXPECT().Create(mock.Anything, mock.Anything).Return(errors.New("failed to create entry"))
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
			wantOK: false,
			wantConditions: []api.Condition{
				{
					Type:    "DataFederationPrivateEndpointsReady",
					Status:  "False",
					Reason:  "InternalError",
					Message: "error creating private endpoint: failed to create entry",
				},
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
			name: "nothing to update in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"}},
				}, nil)
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
			name: "failed to delete when updating entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"}},
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "456", Provider: "bar", Type: "some"}},
				}, nil)
				m.EXPECT().Delete(mock.Anything, &datafederation.DatafederationPrivateEndpointEntry{
					DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"},
				}).Return(errors.New("failed to delete entry"))
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
			wantOK: false,
			wantConditions: []api.Condition{
				{
					Type:    "DataFederationPrivateEndpointsReady",
					Status:  "False",
					Reason:  "InternalError",
					Message: "error deleting private endpoint: failed to delete entry",
				},
			},
		},
		{
			name: "failed to create when updating entry in atlas",
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
				}).Return(errors.New("failed to create entry"))
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
			wantOK: false,
			wantConditions: []api.Condition{
				{
					Type:    "DataFederationPrivateEndpointsReady",
					Status:  "False",
					Reason:  "InternalError",
					Message: "error creating private endpoint: failed to create entry",
				},
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
		{
			name: "failed to delete tracked entry in atlas",
			service: func(m *translation.DatafederationPrivateEndpointServiceMock) datafederation.DatafederationPrivateEndpointService {
				m.EXPECT().List(mock.Anything, mock.Anything).Return([]*datafederation.DatafederationPrivateEndpointEntry{
					{DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"}},
				}, nil)
				m.EXPECT().Delete(mock.Anything, &datafederation.DatafederationPrivateEndpointEntry{
					DataFederationPE: &akov2.DataFederationPE{EndpointID: "123", Provider: "foo", Type: "some"},
				}).Return(errors.New("failed to delete entry"))
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
			wantOK: false,
			wantConditions: []api.Condition{
				{
					Type:    "DataFederationPrivateEndpointsReady",
					Status:  "False",
					Reason:  "InternalError",
					Message: "error deleting private endpoint: failed to delete entry",
				},
			},
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
