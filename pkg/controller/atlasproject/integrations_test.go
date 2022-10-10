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

func TestAreIntegrationsEqual(t *testing.T) {
	atlas := aliasThirdPartyIntegration{
		Type:   "DATADOG",
		APIKey: "****************************4e6f",
		Region: "EU",
	}
	spec := aliasThirdPartyIntegration{
		Type:   "DATADOG",
		APIKey: "actual-valid-id*************4e6f",
		Region: "EU",
	}

	areEqual := AreIntegrationsEqual(&atlas, &spec)
	assert.True(t, areEqual, "Identical objects should be equal")

	spec.APIKey = "non-equal-id************1234"
	areEqual = AreIntegrationsEqual(&atlas, &spec)
	assert.False(t, areEqual, "Should fail if the last 4 characters of APIKey do not match")
}
