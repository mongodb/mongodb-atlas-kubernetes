package networkpeering

import (
	"fmt"
	"testing"

	gofuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

const fuzzIterations = 100

var providerNames = []provider.ProviderName{
	provider.ProviderAWS,
	provider.ProviderAzure,
	provider.ProviderGCP,
}

func FuzzConvertConnection(f *testing.F) {
	for i := 0; i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i)
	}
	f.Fuzz(func(t *testing.T, data []byte, index int) {
		peerData := akov2.NetworkPeer{}
		gofuzz.NewFromGoFuzz(data).Fuzz(&peerData)
		peerData.ProviderName = providerNames[index%3]
		cleanupPeer(&peerData)
		result := fromAtlasConnection(toAtlasConnection(&peerData))
		assert.Equal(t, &peerData, result, "failed for index=%d", index)
	})
}

func FuzzConvertConnectionList(f *testing.F) {
	for i := 0; i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i, (i%5)-1)
	}
	f.Fuzz(func(t *testing.T, data []byte, index int, size int) {
		conns := []admin.BaseNetworkPeeringConnectionSettings{}
		expected := []akov2.NetworkPeer{}
		if size < 0 {
			conns = nil
			expected = nil
		} else {
			for i := 0; i < size; i++ {
				peerData := akov2.NetworkPeer{}
				gofuzz.NewFromGoFuzz(data).Fuzz(&peerData)
				peerData.ProviderName = providerNames[index%3]
				cleanupPeer(&peerData)
				expectedConn := fromAtlasConnection(toAtlasConnection(&peerData))
				expected = append(expected, *expectedConn)
				conns = append(conns, *toAtlasConnection(&peerData))
			}
		}
		result := fromAtlasConnectionList(conns)
		assert.Equal(t, expected, result)
	})
}

func FuzzConvertContainer(f *testing.F) {
	for i := 0; i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i)
	}
	f.Fuzz(func(t *testing.T, data []byte, index int) {
		containerData := ProviderContainer{}
		gofuzz.NewFromGoFuzz(data).Fuzz(&containerData)
		containerData.ProviderName = providerNames[index%3]
		result := fromAtlasContainer(toAtlasContainer(&containerData))
		assert.Equal(t, &containerData, result, "failed for index=%d", index)
	})
}

func FuzzConvertContainerList(f *testing.F) {
	for i := 0; i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i, (i%5)-1)
	}
	f.Fuzz(func(t *testing.T, data []byte, index int, size int) {
		containers := []admin.CloudProviderContainer{}
		expected := []ProviderContainer{}
		if size < 0 {
			containers = nil
			expected = nil
		} else {
			for i := 0; i < size; i++ {
				containerData := ProviderContainer{}
				gofuzz.NewFromGoFuzz(data).Fuzz(&containerData)
				containerData.ProviderName = providerNames[index%3]
				expectedContainer := fromAtlasContainer(toAtlasContainer(&containerData))
				expected = append(expected, *expectedContainer)
				containers = append(containers, *toAtlasContainer(&containerData))
			}
		}
		result := fromAtlasContainerList(containers)
		assert.Equal(t, expected, result)
	})
}

func cleanupPeer(peer *akov2.NetworkPeer) {
	peer.ContainerRegion = ""
	peer.AtlasCIDRBlock = ""
	if peer.ProviderName != provider.ProviderAWS {
		peer.AccepterRegionName = ""
		peer.AWSAccountID = ""
		peer.RouteTableCIDRBlock = ""
		peer.VpcID = ""
	}
	if peer.ProviderName != provider.ProviderGCP {
		peer.GCPProjectID = ""
		peer.NetworkName = ""
	}
	if peer.ProviderName != provider.ProviderAzure {
		peer.AzureDirectoryID = ""
		peer.AzureSubscriptionID = ""
		peer.ResourceGroupName = ""
		peer.VNetName = ""
	}
}
