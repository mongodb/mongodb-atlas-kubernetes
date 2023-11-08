package featureflags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FeatureFlags(t *testing.T) {
	t.Run("Should parse feature without a value", func(t *testing.T) {
		f := NewFeatureFlags(func() []string {
			return []string{"FEATURE_TEST"}
		})
		assert.True(t, f.IsFeaturePresent("FEATURE_TEST"))
	})

	t.Run("Should parse feature with a value", func(t *testing.T) {
		f := NewFeatureFlags(func() []string {
			return []string{"FEATURE_TEST=true"}
		})
		assert.True(t, f.IsFeaturePresent("FEATURE_TEST"))
		assert.Equal(t, "true", f.GetFeatureValue("FEATURE_TEST"))
	})
}
