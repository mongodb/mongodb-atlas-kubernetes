package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func int64Ref(n int64) *int64 {
	return &n
}

func intRef(n int) *int {
	return &n
}

func newOperatorSpec(diskOpts *int64, ebsVolumeType, instanceSize string, nodeCount *int) *v1.Specs {
	return &v1.Specs{
		DiskIOPS:      diskOpts,
		EbsVolumeType: ebsVolumeType,
		InstanceSize:  instanceSize,
		NodeCount:     nodeCount,
	}
}

func newAtlasSpec(diskOpts *int64, ebsVolumeType, instanceSize string, nodeCount *int) *mongodbatlas.Specs {
	return &mongodbatlas.Specs{
		DiskIOPS:      diskOpts,
		EbsVolumeType: ebsVolumeType,
		InstanceSize:  instanceSize,
		NodeCount:     nodeCount,
	}
}

func TestSpecsAreEqual(t *testing.T) {
	t.Run("Nil specs are equal", func(t *testing.T) {
		assert.True(t, specsAreEqual(nil, nil))
	})

	t.Run("One nil spec is not equal", func(t *testing.T) {
		t.Run("Operator spec is not nil", func(t *testing.T) {
			assert.False(t, specsAreEqual(nil, newOperatorSpec(int64Ref(10), "test", "M5", intRef(5))))
		})
		t.Run("Atlas spec is not nil", func(t *testing.T) {
			assert.False(t, specsAreEqual(newAtlasSpec(int64Ref(10), "test", "M5", intRef(5)), nil))
		})
	})

	t.Run("Equal specs", func(t *testing.T) {
		assert.True(t, specsAreEqual(newAtlasSpec(int64Ref(10), "test", "M5", intRef(5)), newOperatorSpec(int64Ref(10), "test", "M5", intRef(5))))
	})
	t.Run("Different specs", func(t *testing.T) {
		assert.False(t, specsAreEqual(newAtlasSpec(int64Ref(10), "test", "M5", intRef(5)), newOperatorSpec(int64Ref(10), "different", "M5", intRef(1))))
	})
}
