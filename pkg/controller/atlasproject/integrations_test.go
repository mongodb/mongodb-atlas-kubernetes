package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
)

func TestToAlias(t *testing.T) {
	sample := []*mongodbatlas.ThirdPartyIntegration{{
		Type:   "DATADOG",
		APIKey: "some",
		Region: "EU",
	}}
	result := toAliasThirdPartyIntegration(sample)
	assert.Equal(t, sample[0].APIKey, result[0].APIKey)
	assert.Equal(t, sample[0].Type, result[0].Type)
	assert.Equal(t, sample[0].Region, result[0].Region)
}
