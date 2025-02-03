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
		gofuzz.NewFromGoFuzz(data).Fuzz(&peerData)
		peerData.Provider = providerNames[index%3]
		cleanupPeer(&peerData)
		atlasConn, err := toAtlasConnection(&peerData)
		require.NoError(t, err)
		result, err := fromAtlasConnection(atlasConn)
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
			gofuzz.NewFromGoFuzz(data).Fuzz(&peerData)
			peerData.Provider = providerNames[index%3]
			cleanupPeer(&peerData)
			atlasConn, err := toAtlasConnection(&peerData)
			require.NoError(t, err)
			expectedConn, err := fromAtlasConnection(atlasConn)
			require.NoError(t, err)
			expected = append(expected, *expectedConn)
			atlasConnItem, err := toAtlasConnection(&peerData)
			require.NoError(t, err)
			conns = append(conns, *atlasConnItem)
		}
		result, err := fromAtlasConnectionList(conns)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func FuzzConvertContainer(f *testing.F) {
	for i := uint(0); i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i)
	}
	f.Fuzz(func(t *testing.T, data []byte, index uint) {
		containerData := ProviderContainer{}
		gofuzz.NewFromGoFuzz(data).Fuzz(&containerData)
		containerData.Provider = providerNames[index%3]
		cleanupContainer(&containerData)
		result := fromAtlasContainer(toAtlasContainer(&containerData))
		assert.Equal(t, &containerData, result, "failed for index=%d", index)
	})
}

func FuzzConvertListOfContainers(f *testing.F) {
	for i := uint(0); i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i, (i % 5))
	}
	f.Fuzz(func(t *testing.T, data []byte, index uint, size uint) {
		containers := []admin.CloudProviderContainer{}
		expected := []ProviderContainer{}
		for i := uint(0); i < size; i++ {
			containerData := ProviderContainer{}
			gofuzz.NewFromGoFuzz(data).Fuzz(&containerData)
			containerData.Provider = providerNames[index%3]
			cleanupContainer(&containerData)
			expectedContainer := fromAtlasContainer(toAtlasContainer(&containerData))
			expected = append(expected, *expectedContainer)
			containers = append(containers, *toAtlasContainer(&containerData))
		}
		result := fromAtlasContainerList(containers)
		assert.Equal(t, expected, result)
	})
}

func cleanupPeer(peer *NetworkPeer) {
	peer.ID = ""
	if peer.Provider != string(provider.ProviderAWS) {
		peer.AWSConfiguration = nil
	}
	if peer.Provider != string(provider.ProviderGCP) {
		peer.GCPConfiguration = nil
	}
	if peer.Provider != string(provider.ProviderAzure) {
		peer.AzureConfiguration = nil
	}
	// status fields are only populated from Atlas they do not complete a roundtrip
	peer.Status = ""
	peer.ErrorMessage = ""
	peer.AWSStatus = nil
}

func cleanupContainer(container *ProviderContainer) {
	// status fields are only populated from Atlas they do not complete a roundtrip
	container.AWSStatus = nil
	container.AzureStatus = nil
	container.GoogleStatus = nil
}
