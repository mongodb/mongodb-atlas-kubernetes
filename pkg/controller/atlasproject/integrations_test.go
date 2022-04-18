package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	// "github.com/stretchr/testify/assert"
	// "go.mongodb.org/atlas/mongodbatlas"
	// "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
)

func TestFromAtlas(t *testing.T) {
	t.Run("FromAtlas", func(t *testing.T) {
		t.Log("test")
		// atlasSide := []*mongodbatlas.ThirdPartyIntegration{{
		// 	Type:   "DATADOG",
		// 	APIKey: "somekey",
		// 	Region: "EU",
		// }}
		// converted := fromAtlas(atlasSide)

		// assert.Equal(t, atlasSide[0].Type, converted[0].Type)
		// assert.Equal(t, atlasSide[0].APIKey, converted[0].APIKey)
		// assert.Equal(t, atlasSide[0].URL, converted[0].URL)
	})
}

func TestToAlias(t *testing.T) {
	//toAliasThirdPartyIntegration(list []*mongodbatlas.ThirdPartyIntegration) []aliasThirdPartyIntegration
	sample := []*mongodbatlas.ThirdPartyIntegration{{
		Type:   "DATADOG",
		APIKey: "some",
		Region: "EU",
	}}
	result := toAliasThirdPartyIntegration(sample)
	assert.Equal(t, sample[0].APIKey, result[0].APIKey)
}
