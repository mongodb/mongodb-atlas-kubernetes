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

package datafederation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
)

func TestRoundtrip_DataFederationPE(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &DatafederationPrivateEndpointEntry{}
		f.Fuzz(fuzzed)
		// ignore non-Atlas fields
		fuzzed.ProjectID = ""

		toAtlasResult := endpointToAtlas(fuzzed)
		fromAtlasResult := endpointFromAtlas(toAtlasResult, "")

		equals := fuzzed.EqualsTo(fromAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fuzzed, fromAtlasResult))
		}
		require.True(t, equals)
	}
}

func TestMapDatafederationPrivateEndpoints(t *testing.T) {
	tests := map[string]struct {
		dataFederation *akov2.AtlasDataFederation
		endpoints      []*DatafederationPrivateEndpointEntry
		expectedResult map[string]*DataFederationPrivateEndpoint
		expectedErr    string
	}{
		"failed to parse last config applied annotation": {
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: "wrong,",
					},
				},
			},
			endpoints:      []*DatafederationPrivateEndpointEntry{},
			expectedResult: nil,
			expectedErr:    "error reading data federation from last applied annotation: invalid character 'w' looking for beginning of value",
		},
		"map without last applied configuration": {
			dataFederation: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{
						{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						{
							Provider:   "AZURE",
							Type:       "DATA_LAKE",
							EndpointID: "azure/resource/id",
						},
					},
				},
			},
			endpoints: []*DatafederationPrivateEndpointEntry{
				{
					DataFederationPE: &akov2.DataFederationPE{
						Provider:   "AWS",
						Type:       "DATA_LAKE",
						EndpointID: "vpcpe-123456",
					},
					ProjectID: "project-id",
				},
				{
					DataFederationPE: &akov2.DataFederationPE{
						Provider:   "AZURE",
						Type:       "DATA_LAKE",
						EndpointID: "azure/resource/id",
					},
					ProjectID: "project-id",
				},
			},
			expectedResult: map[string]*DataFederationPrivateEndpoint{
				"vpcpe-123456": {
					AKO: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						ProjectID: "project-id",
					},
					Atlas: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						ProjectID: "project-id",
					},
					LastApplied: nil,
				},
				"azure/resource/id": {
					AKO: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AZURE",
							Type:       "DATA_LAKE",
							EndpointID: "azure/resource/id",
						},
						ProjectID: "project-id",
					},
					Atlas: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AZURE",
							Type:       "DATA_LAKE",
							EndpointID: "azure/resource/id",
						},
						ProjectID: "project-id",
					},
					LastApplied: nil,
				},
			},
		},
		"map with last applied configuration": {
			dataFederation: &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						customresource.AnnotationLastAppliedConfiguration: "{\"name\":\"\",\"privateEndpoints\":[{\"endpointId\":\"vpcpe-123456\"," +
							"\"provider\":\"AWS\",\"type\":\"DATA_LAKE\"}]}",
					},
				},
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{
						{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						{
							Provider:   "AZURE",
							Type:       "DATA_LAKE",
							EndpointID: "azure/resource/id",
						},
					},
				},
			},
			endpoints: []*DatafederationPrivateEndpointEntry{
				{
					DataFederationPE: &akov2.DataFederationPE{
						Provider:   "AWS",
						Type:       "DATA_LAKE",
						EndpointID: "vpcpe-123456",
					},
					ProjectID: "project-id",
				},
				{
					DataFederationPE: &akov2.DataFederationPE{
						Provider:   "AZURE",
						Type:       "DATA_LAKE",
						EndpointID: "azure/resource/id",
					},
					ProjectID: "project-id",
				},
			},
			expectedResult: map[string]*DataFederationPrivateEndpoint{
				"vpcpe-123456": {
					AKO: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						ProjectID: "project-id",
					},
					Atlas: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						ProjectID: "project-id",
					},
					LastApplied: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
						ProjectID: "project-id",
					},
				},
				"azure/resource/id": {
					AKO: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AZURE",
							Type:       "DATA_LAKE",
							EndpointID: "azure/resource/id",
						},
						ProjectID: "project-id",
					},
					Atlas: &DatafederationPrivateEndpointEntry{
						DataFederationPE: &akov2.DataFederationPE{
							Provider:   "AZURE",
							Type:       "DATA_LAKE",
							EndpointID: "azure/resource/id",
						},
						ProjectID: "project-id",
					},
					LastApplied: nil,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m, err := MapDatafederationPrivateEndpoints("project-id", tt.dataFederation, tt.endpoints)
			if err != nil {
				assert.EqualError(t, err, tt.expectedErr)
			}
			assert.Equal(t, tt.expectedResult, m)
		})
	}
}
