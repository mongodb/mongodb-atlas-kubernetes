package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPluginCatalog(t *testing.T) {
	t.Run("should create new plugin catalog", func(t *testing.T) {
		catalog := NewPluginCatalog(nil)
		assert.Len(t, catalog, 12)
	})
}
