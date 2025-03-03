package networkcontainer

import (
	"fmt"
	"testing"

	gofuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
)

const fuzzIterations = 100

var providerNames = []string{
	string(provider.ProviderAWS),
	string(provider.ProviderAzure),
	string(provider.ProviderGCP),
}

func FuzzConvertContainer(f *testing.F) {
	for i := uint(0); i < fuzzIterations; i++ {
		f.Add(([]byte)(fmt.Sprintf("seed sample %x", i)), i)
	}
	f.Fuzz(func(t *testing.T, data []byte, index uint) {
		containerData := NetworkContainer{}
		gofuzz.NewFromGoFuzz(data).Fuzz(&containerData)
		containerData.Provider = providerNames[index%3]
		cleanupContainer(&containerData)
		result := fromAtlas(toAtlas(&containerData))
		assert.Equal(t, &containerData, result, "failed for index=%d", index)
	})
}

func cleanupContainer(container *NetworkContainer) {
	container.AtlasNetworkContainerConfig.ID = ""
	// status fields are only populated from Atlas they do not complete a roundtrip
	container.AWSStatus = nil
	container.AzureStatus = nil
	container.GCPStatus = nil
}
