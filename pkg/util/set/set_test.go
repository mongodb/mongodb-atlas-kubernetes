package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromSlice(t *testing.T) {
	t.Run("return a map(set) from a slice", func(t *testing.T) {
		type myType struct {
			Key   int
			Value string
		}
		list := []myType{
			{
				Key:   0,
				Value: "Zero",
			},
			{
				Key:   1,
				Value: "One",
			},
			{
				Key:   2,
				Value: "Two",
			},
		}

		assert.Equal(
			t,
			map[int]myType{
				0: {Key: 0, Value: "Zero"},
				1: {Key: 1, Value: "One"},
				2: {Key: 2, Value: "Two"},
			},
			FromSlice(list, func(i myType) int {
				return i.Key
			}))
	})
}
