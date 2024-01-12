package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyWithSkip(t *testing.T) {
	t.Run("should return the collection without the skip item", func(t *testing.T) {
		c := []string{"a", "b", "c", "d", "e"}
		assert.Equal(t, []string{"a", "b", "d", "e"}, CopyWithSkip(c, "c"))
	})

	t.Run("should return the same collection when the skip item is not present", func(t *testing.T) {
		c := []string{"a", "b", "c", "d", "e"}
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, CopyWithSkip(c, "f"))
	})
}
