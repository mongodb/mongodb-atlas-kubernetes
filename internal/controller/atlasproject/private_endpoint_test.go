package atlasproject

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestGetEndpointsNotInAtlas(t *testing.T) {
	const region1 = "SOME_REGION"
	const region2 = "OTHER_REGION"
	specPEs := []akov2.PrivateEndpoint{
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
		EndpointService: admin.EndpointService{
			CloudProvider: string(provider.ProviderAWS),
			RegionName:    admin.PtrString(region1),
		},
	})

	uniqueItems, _ = getEndpointsNotInAtlas(specPEs, atlasPEs)
	assert.Equalf(t, len(uniqueItems), 1, "getEndpointsNotInAtlas should remove both PE Service copies if there is one in Atlas")
}

func TestGetEndpointsNotInSpec(t *testing.T) {
	tests := map[string]struct {
		specPEs       []akov2.PrivateEndpoint
		atlasPEs      []atlasPE
		lastPEs       map[string]akov2.PrivateEndpoint
		expectedItems []atlasPE
	}{
		"should return no items when spec and atlas are the same": {
			specPEs: []akov2.PrivateEndpoint{
				{
					Provider: provider.ProviderAWS,
					Region:   "us_east1",
				},
				{
					Provider: provider.ProviderAWS,
					Region:   "us_east2",
				},
			},
			atlasPEs: []atlasPE{
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_east1"),
					},
				},
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_east2"),
					},
				},
			},
			expectedItems: []atlasPE{},
		},
		"should return no items when spec and atlas are different but not previously managed by the operator": {
			specPEs: []akov2.PrivateEndpoint{
				{
					Provider: provider.ProviderAWS,
					Region:   "us_east1",
				},
				{
					Provider: provider.ProviderAWS,
					Region:   "us_east2",
				},
			},
			atlasPEs: []atlasPE{
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_east1"),
					},
				},
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_west1"),
					},
				},
			},
			expectedItems: []atlasPE{},
		},
		"should return items when spec and atlas are different but previously managed by the operator": {
			specPEs: []akov2.PrivateEndpoint{
				{
					Provider: provider.ProviderAWS,
					Region:   "us_east1",
				},
			},
			atlasPEs: []atlasPE{
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_east1"),
					},
				},
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_east2"),
					},
				},
			},
			lastPEs: map[string]akov2.PrivateEndpoint{
				"AWSaesstu": {
					Provider: provider.ProviderAWS,
					Region:   "us_east1",
				},
				"AWS2aesstu": {
					Provider: provider.ProviderAWS,
					Region:   "us_east2",
				},
			},
			expectedItems: []atlasPE{
				{
					EndpointService: admin.EndpointService{
						CloudProvider: string(provider.ProviderAWS),
						RegionName:    admin.PtrString("us_east2"),
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			uniqueItems := getEndpointsNotInSpec(tt.specPEs, tt.atlasPEs, tt.lastPEs)
			assert.Equal(t, tt.expectedItems, uniqueItems)
		})
	}
}

func TestMapLastAppliedPrivateEndpoint(t *testing.T) {
	tests := map[string]struct {
		annotations   map[string]string
		expectedPEs   map[string]akov2.PrivateEndpoint
		expectedError string
	}{
		"should return error when last spec annotation is wrong": {
			annotations: map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"},
			expectedError: "error reading AtlasProject Spec from annotation [mongodb.com/last-applied-configuration]:" +
				" invalid character 'w' looking for beginning of object key string",
		},
		"should return nil when there is no last spec": {},
		"should return map of last private endpoints": {
			annotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"privateEndpoints\": [{\"provider\":\"AWS\",\"region\":\"us_east1\"},{\"provider\":\"AWS\",\"region\":\"us_east2\"}]}"},
			expectedPEs: map[string]akov2.PrivateEndpoint{
				"AWSaesstu": {
					Provider: provider.ProviderAWS,
					Region:   "us_east1",
				},
				"AWS2aesstu": {
					Provider: provider.ProviderAWS,
					Region:   "us_east2",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &akov2.AtlasProject{}
			p.WithAnnotations(tt.annotations)

			result, err := mapLastAppliedPrivateEndpoint(p)
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expectedPEs, result)
		})
	}
}

func TestPrivateEndpointsNonGreedyBehaviour(t *testing.T) {
	for _, tc := range []struct {
		title            string
		lastAppliedPEids []string
		specPEids        []string
		atlasPEids       []string
		wantRemoved      []string
		wantResult       workflow.Result
	}{
		{
			title:            "no last applied no removal in Atlas",
			lastAppliedPEids: []string{},
			specPEids:        []string{},
			atlasPEids:       []string{"pe1", "pe2"},
			wantRemoved:      []string{},
			wantResult:       workflow.OK(),
		},
		{
			title:            "removed from last applied removes from Atlas",
			lastAppliedPEids: []string{"pe1", "pe2"},
			specPEids:        []string{"pe1"},
			atlasPEids:       []string{"pe1", "pe2"},
			wantRemoved:      []string{"pe2"},
			wantResult: workflow.InProgress(
				workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting"),
		},
		{
			title:            "removed all from last applied removes all from Atlas",
			lastAppliedPEids: []string{"pe1", "pe2"},
			specPEids:        []string{},
			atlasPEids:       []string{"pe1", "pe2"},
			wantRemoved:      []string{"pe1", "pe2"},
			wantResult: workflow.InProgress(
				workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting"),
		},
		{
			title:            "not in last applied not removed from Atlas",
			lastAppliedPEids: []string{"pe1"},
			specPEids:        []string{"pe1"},
			atlasPEids:       []string{"pe1", "pe2"},
			wantRemoved:      []string{},
			wantResult: workflow.InProgress(
				workflow.ProjectPrivateEndpointIsNotReadyInAtlas, "Interface Private Endpoint is not ready"),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			prj := newTestPEProject(tc.specPEids)
			lastPrj := newTestPEProject(tc.lastAppliedPEids)
			prj.Annotations[customresource.AnnotationLastAppliedConfiguration] = jsonize(t, lastPrj.Spec)

			privateEndpointsAPI := mockadmin.NewPrivateEndpointServicesApi(t)
			privateEndpointsAPI.EXPECT().ListPrivateEndpointServices(mock.Anything, mock.Anything, "AWS").
				Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: privateEndpointsAPI}).Once()
			privateEndpointsAPI.EXPECT().ListPrivateEndpointServicesExecute(
				mock.AnythingOfType("admin.ListPrivateEndpointServicesApiRequest")).Return(
				synthesizeAtlasPEs(tc.atlasPEids), nil, nil,
			).Once()
			privateEndpointsAPI.EXPECT().ListPrivateEndpointServices(mock.Anything, mock.Anything, "AZURE").
				Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: privateEndpointsAPI}).Once()
			privateEndpointsAPI.EXPECT().ListPrivateEndpointServicesExecute(
				mock.AnythingOfType("admin.ListPrivateEndpointServicesApiRequest")).Return(
				nil, nil, nil,
			).Once()
			privateEndpointsAPI.EXPECT().ListPrivateEndpointServices(mock.Anything, mock.Anything, "GCP").
				Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: privateEndpointsAPI}).Once()
			privateEndpointsAPI.EXPECT().ListPrivateEndpointServicesExecute(
				mock.AnythingOfType("admin.ListPrivateEndpointServicesApiRequest")).Return(
				nil, nil, nil,
			).Once()

			removals := len(tc.wantRemoved)
			if removals > 0 {
				privateEndpointsAPI.EXPECT().DeletePrivateEndpointServiceWithParams(
					mock.Anything, mock.Anything,
				).Return(admin.DeletePrivateEndpointServiceApiRequest{ApiService: privateEndpointsAPI}).Times(removals)
				privateEndpointsAPI.EXPECT().DeletePrivateEndpointServiceExecute(
					mock.AnythingOfType("admin.DeletePrivateEndpointServiceApiRequest")).Return(
					nil, nil, nil,
				).Times(removals)
			}
			privateEndpointsAPI.EXPECT().CreatePrivateEndpointWithParams(
				mock.Anything, mock.Anything,
			).Return(admin.CreatePrivateEndpointApiRequest{ApiService: privateEndpointsAPI}).Maybe()
			privateEndpointsAPI.EXPECT().CreatePrivateEndpointExecute(
				mock.AnythingOfType("admin.CreatePrivateEndpointApiRequest")).Return(
				nil, nil, nil,
			).Maybe()

			workflowCtx := workflow.Context{
				Log:     zaptest.NewLogger(t).Sugar(),
				Context: context.Background(),
				SdkClient: &admin.APIClient{
					PrivateEndpointServicesApi: privateEndpointsAPI,
				},
			}

			result := ensurePrivateEndpoint(&workflowCtx, prj)
			require.Equal(t, tc.wantResult, result)
		})
	}
}

func newTestPEProject(peIDs []string) *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name:             "test-project",
			PrivateEndpoints: synthesizePEs(peIDs),
		},
	}
}

func synthesizePEs(peIDs []string) []akov2.PrivateEndpoint {
	pes := make([]akov2.PrivateEndpoint, 0, len(peIDs))
	for _, id := range peIDs {
		pes = append(pes, akov2.PrivateEndpoint{
			Provider: "AWS",
			ID:       id,
			Region:   fmt.Sprintf("fake-region-%s", id),
		})
	}
	return pes
}

func synthesizeAtlasPEs(peIDs []string) []admin.EndpointService {
	atlasPEs := make([]admin.EndpointService, 0, len(peIDs))
	for _, id := range peIDs {
		atlasPEs = append(atlasPEs, admin.EndpointService{
			CloudProvider: "AWS",
			Id:            pointer.MakePtr(id),
			RegionName:    pointer.MakePtr(fmt.Sprintf("fake-region-%s", id)),
			Status:        pointer.MakePtr("AVAILABLE"),
		})
	}
	return atlasPEs
}
