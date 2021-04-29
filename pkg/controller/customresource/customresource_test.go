package customresource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestResourceShouldBeLeftInAtlas(t *testing.T) {
	t.Run("Empty annotations", func(t *testing.T) {
		assert.False(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{}))
	})

	t.Run("Other annotations", func(t *testing.T) {
		assert.False(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{"foo": "bar"},
			},
		}))
	})

	t.Run("Annotation present, resources should be removed", func(t *testing.T) {
		assert.False(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				// Any other value except for "keep" is considered as "purge"
				Annotations: map[string]string{ResourcePolicyAnnotation: "foobar"},
			},
		}))
	})

	t.Run("Annotation present, resources should be kept", func(t *testing.T) {
		assert.True(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{ResourcePolicyAnnotation: ResourcePolicyKeep},
			},
		}))
	})
}
