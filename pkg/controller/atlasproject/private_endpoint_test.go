package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

func TestGetEndpointsNotInAtlas(t *testing.T) {
	const region1 = "SOME_REGION"
	const region2 = "OTHER_REGION"
	specPEs := []v1.PrivateEndpoint{
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
	result := getEndpointsNotInAtlas(specPEs, atlasPEs)
	assert.Equalf(t, len(result), 2, "getEndpointsNotInAtlas should remove a duplicate PE Service")
	assert.NotEqualf(t, result[0].Region, result[1].Region, "getEndpointsNotInAtlas should return unique PEs")

	atlasPEs = append(atlasPEs, atlasPE{
		ProviderName: string(provider.ProviderAWS),
		RegionName:   region1,
	})

	result = getEndpointsNotInAtlas(specPEs, atlasPEs)
	assert.Equalf(t, len(result), 1, "getEndpointsNotInAtlas should remove both PE Service copies if there is one in Atlas")
}
