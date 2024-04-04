package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
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
	}
	atlasPEs := []atlasPE{
		{
			EndpointService: admin.EndpointService{
				CloudProvider: string(provider.ProviderAWS),
				RegionName:    admin.PtrString(region1),
			},
		},
		{
			EndpointService: admin.EndpointService{
				CloudProvider: string(provider.ProviderAWS),
				RegionName:    admin.PtrString(region1),
			},
		},
	}

	uniqueItems := getEndpointsNotInSpec(specPEs, atlasPEs)
	assert.Equalf(t, 0, len(uniqueItems), "getEndpointsNotInSpec should not return anything if PEs are in spec")

	atlasPEs = append(atlasPEs, atlasPE{
		EndpointService: admin.EndpointService{
			CloudProvider: string(provider.ProviderAWS),
			RegionName:    admin.PtrString(region2),
		},
	})
	uniqueItems = getEndpointsNotInSpec(specPEs, atlasPEs)
	assert.Equalf(t, 1, len(uniqueItems), "getEndpointsNotInSpec should get a spec item")
}
