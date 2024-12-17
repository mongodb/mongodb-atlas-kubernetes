package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
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
