package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	// "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
)

func TestFromAtlas(t *testing.T) {
	t.Run("FromAtlas", func(t *testing.T) {
		t.Log("test")
		atlasSide := []*mongodbatlas.ThirdPartyIntegration{{
			Type:        "DATADOG",
			LicenseKey:  "",
			AccountID:   "",
			WriteToken:  "",
			ReadToken:   "",
			APIKey:      "somekey",
			Region:      "EU",
			ServiceKey:  "",
			APIToken:    "",
			TeamName:    "",
			ChannelName: "",
			RoutingKey:  "",
			FlowName:    "",
			OrgName:     "",
			URL:         "",
			Secret:      "",
		}}
		converted := fromAtlas(atlasSide)

		assert.Equal(t, atlasSide[0].Type, converted[0].Type)
		assert.Equal(t, atlasSide[0].APIKey, converted[0].APIKey)
		assert.Equal(t, atlasSide[0].URL, converted[0].URL)
	})

}
