package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func Test_getIndicesFromAnnotations(t *testing.T) {
	t.Run("Should return empty map if annotations are empty", func(t *testing.T) {
		in := map[string]string{}
		assert.Nil(t, getIndicesFromAnnotations(in))
	})

	t.Run("Should return valid IndexName:IndexID pairs", func(t *testing.T) {
		in := map[string]string{
			DeploymentIndicesAnnotation: "IndexOne:1,IndexTwo:2,IndexThree:3",
		}
		result := getIndicesFromAnnotations(in)
		assert.Len(t, result, 3)
		assert.Equal(t, "1", result["IndexOne"])
		assert.Equal(t, "2", result["IndexTwo"])
		assert.Equal(t, "3", result["IndexThree"])
	})

	t.Run("Should return ONLY valid IndexName:IndexID pairs", func(t *testing.T) {
		in := map[string]string{
			DeploymentIndicesAnnotation: "IndexOne:1,IndexTwo:2,IndexThree:3,IndexName4",
		}
		result := getIndicesFromAnnotations(in)
		assert.Len(t, result, 3)
		assert.Equal(t, "1", result["IndexOne"])
		assert.Equal(t, "2", result["IndexTwo"])
		assert.Equal(t, "3", result["IndexThree"])
	})
}

func Test_verifyAllIndicesNamesAreUnique(t *testing.T) {
	t.Run("Should return true if all indices names are unique", func(t *testing.T) {
		in := []akov2.SearchIndex{
			{
				Name: "Index-One",
			},
			{
				Name: "Index-Two",
			},
			{
				Name: "Index-Three",
			},
		}
		assert.True(t, verifyAllIndicesNamesAreUnique(in))
	})
	t.Run("Should return false if one index name appeared twice", func(t *testing.T) {
		in := []akov2.SearchIndex{
			{
				Name: "Index-One",
			},
			{
				Name: "Index-Two",
			},
			{
				Name: "Index-One",
			},
		}
		assert.False(t, verifyAllIndicesNamesAreUnique(in))
	})
}
