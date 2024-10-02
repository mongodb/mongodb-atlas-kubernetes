package datafederation

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akocmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
)

func NonEmptyString(s *string, c fuzz.Continue) {
	for {
		fuzz.UnicodeRange{First: 'a', Last: 'z'}.CustomStringFuzzFunc()(s, c)
		if *s != "" {
			return
		}
	}
}

// EnsureEmptySliceOf ensures that slices in the form of *[]typ are initialized with empty slices.
func EnsureEmptySliceOf[typ any](f *fuzz.Fuzzer) func(**[]typ, fuzz.Continue) {
	return func(slice **[]typ, c fuzz.Continue) {
		var empty []typ
		*slice = &empty
		f.FuzzNoCustom(*slice)
	}
}

func TestRoundtrip_DataFederation(t *testing.T) {
	f := fuzz.New().NilChance(0.0).NumElements(1, 10)
	f.Funcs(
		NonEmptyString,
		EnsureEmptySliceOf[admin.DataLakeDatabaseCollection](f),
		EnsureEmptySliceOf[admin.DataLakeDatabaseDataSourceSettings](f),
		EnsureEmptySliceOf[admin.DataLakeDatabaseInstance](f),
		EnsureEmptySliceOf[string](f),
	)

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

		require.NoError(t, akocmp.Normalize(&atlas))
		require.NoError(t, akocmp.Normalize(toAtlasResult))

		equals := reflect.DeepEqual(&atlas, toAtlasResult)
		if !equals {
			t.Log(cmp.Diff(&atlas, toAtlasResult))
		}
		require.True(t, equals)
	}
}
