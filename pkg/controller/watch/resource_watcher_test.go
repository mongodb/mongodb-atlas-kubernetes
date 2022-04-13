package watch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

func TestEnsureResourcesAreWatched(t *testing.T) {
	t.Run("One secret is watched by two resources", func(t *testing.T) {
		watcher := NewResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		project2 := kube.ObjectKey("test", "project2")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")

		watcher.EnsureResourcesAreWatched(project1, "Secret", zap.S(), connectionSecret)
		watcher.EnsureResourcesAreWatched(project2, "Secret", zap.S(), connectionSecret)

		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}: {project1: true, project2: true},
		}
		assert.Equal(t, expectedWatched, watcher.WatchedResources)
	})
	t.Run("One resource watches two secrets", func(t *testing.T) {
		watcher := NewResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")
		connectionSecret2 := kube.ObjectKey("test", "connectionSecret2")

		watcher.EnsureResourcesAreWatched(project1, "Secret", zap.S(), connectionSecret, connectionSecret2)

		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}:  {project1: true},
			{ResourceKind: "Secret", Resource: connectionSecret2}: {project1: true},
		}
		assert.Equal(t, expectedWatched, watcher.WatchedResources)
	})
	t.Run("Resource stops watching one secret", func(t *testing.T) {
		watcher := NewResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		project2 := kube.ObjectKey("test", "project2")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")
		connectionSecret2 := kube.ObjectKey("test", "connectionSecret2")
		connectionSecret3 := kube.ObjectKey("test", "connectionSecret3")

		watcher.EnsureResourcesAreWatched(project1, "Secret", zap.S(), connectionSecret, connectionSecret2)
		watcher.EnsureResourcesAreWatched(project2, "Secret", zap.S(), connectionSecret, connectionSecret2)

		// The second secret is not watched any more
		watcher.EnsureResourcesAreWatched(project1, "Secret", zap.S(), connectionSecret, connectionSecret3)

		// We expect that the watching state stays consistent and the project 1 doesn't watch secret2 anymore

		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}:  {project1: true, project2: true},
			{ResourceKind: "Secret", Resource: connectionSecret2}: {project2: true},
			{ResourceKind: "Secret", Resource: connectionSecret3}: {project1: true},
		}
		assert.Equal(t, expectedWatched, watcher.WatchedResources)
	})
	// TODO: add test for different kind of resources
}
