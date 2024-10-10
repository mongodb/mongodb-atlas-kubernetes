package datafederation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func TestRoundtrip_DataFederation(t *testing.T) {
	f := fuzz.New()

	for i := 0; i < 100; i++ {
		var atlas admin.DataLakeTenant
		f.Fuzz(&atlas)

		// ignore read-only fields
		if atlas.CloudProviderConfig != nil {
			atlas.CloudProviderConfig.Aws.ExternalId = nil
			atlas.CloudProviderConfig.Aws.IamAssumedRoleARN = nil
			atlas.CloudProviderConfig.Aws.IamUserARN = nil
		}
		atlas.Hostnames = nil
		atlas.PrivateEndpointHostnames = nil
		atlas.State = nil
		if atlas.Storage != nil && atlas.Storage.Stores != nil {
			for i := range *atlas.Storage.Stores {
				(*atlas.Storage.Stores)[i].ProjectId = nil
			}
		}

		fromAtlasResult, err := fromAtlas(&atlas)
		require.NoError(t, err)
		toAtlasResult := toAtlas(fromAtlasResult)
		roundtripAtlasResult, err := fromAtlas(toAtlasResult)
		require.NoError(t, err)

		equals := fromAtlasResult.SpecEqualsTo(roundtripAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fromAtlasResult, roundtripAtlasResult))
		}
		require.True(t, equals)
	}
}
