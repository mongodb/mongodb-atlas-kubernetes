package networkpeering

import (
	"fmt"
	"testing"

	gofuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
)

const fuzzIterations = 100

var providerNames = []string{
	string(provider.ProviderAWS),
	string(provider.ProviderAzure),
	string(provider.ProviderGCP),
}

func FuzzConvertConnection(f *testing.F) {
	for i := uint(0); i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i)
	}
	f.Fuzz(func(t *testing.T, data []byte, index uint) {
		peerData := NetworkPeer{}
		fuzzPeer(gofuzz.NewFromGoFuzz(data), index, &peerData)
		atlasConn, err := toAtlas(&peerData)
		require.NoError(t, err)
		result, err := fromAtlas(atlasConn)
		require.NoError(t, err)
		assert.Equal(t, &peerData, result, "failed for index=%d", index)
	})
}

func FuzzConvertListOfConnections(f *testing.F) {
	for i := uint(0); i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i, (i % 5))
	}
	f.Fuzz(func(t *testing.T, data []byte, index uint, size uint) {
		conns := []admin.BaseNetworkPeeringConnectionSettings{}
		expected := []NetworkPeer{}
		for i := uint(0); i < size; i++ {
			peerData := NetworkPeer{}
			fuzzPeer(gofuzz.NewFromGoFuzz(data), index, &peerData)
			atlasConn, err := toAtlas(&peerData)
			require.NoError(t, err)
			expectedConn, err := fromAtlas(atlasConn)
			require.NoError(t, err)
			expected = append(expected, *expectedConn)
			atlasConnItem, err := toAtlas(&peerData)
			require.NoError(t, err)
			conns = append(conns, *atlasConnItem)
		}
		result, err := fromAtlasConnectionList(conns)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func fuzzPeer(fuzzer *gofuzz.Fuzzer, index uint, peer *NetworkPeer) {
	fuzzer.NilChance(0).Fuzz(peer)
	peer.ID = ""                           // ID is provided by Atlas, cannoy complete a roundtrip
	peer.Provider = providerNames[index%3] // provider can only be one of 3 AWS, AZURE or GCP
	switch peer.Provider {                 // only the selected provider config is expected
	case string(provider.ProviderAWS):
		peer.AzureConfiguration = nil
		peer.GCPConfiguration = nil
	case string(provider.ProviderAzure):
		peer.AWSConfiguration = nil
		peer.GCPConfiguration = nil
	case string(provider.ProviderGCP):
		peer.AWSConfiguration = nil
		peer.AzureConfiguration = nil
	}
	// status fields are only populated from Atlas they do not complete a roundtrip
	peer.Status = ""
	peer.ErrorMessage = ""
	peer.AWSStatus = nil
}
