// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	for i := range uint(fuzzIterations) {
		f.Add(fmt.Appendf(nil, "seed sample %x", i), i)
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
