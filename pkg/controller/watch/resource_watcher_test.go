package watch

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func TestWatchedResourcesSnapshot(t *testing.T) {
	for _, tc := range []struct {
		name      string
		dependant client.ObjectKey
		resources []WatchedObject
		want      map[WatchedObject]map[client.ObjectKey]bool
	}{
		{
			name: "empty",
		},
		{
			name:      "no watched resources",
			dependant: kube.ObjectKey("test", "project1"),
		},
		{
			name: "one watched resources",

			dependant: kube.ObjectKey("test", "project1"),
			resources: []WatchedObject{
				{
					ResourceKind: "Secret",
					Resource:     kube.ObjectKey("test", "secret"),
				},
			},

			want: map[WatchedObject]map[client.ObjectKey]bool{
				{
					ResourceKind: "Secret",
					Resource:     kube.ObjectKey("test", "secret"),
				}: {
					kube.ObjectKey("test", "project1"): true,
				},
			},
		},
		{
			name: "multiple watched resources",

			dependant: kube.ObjectKey("test", "project1"),
			resources: []WatchedObject{
				{
					ResourceKind: "Secret",
					Resource:     kube.ObjectKey("test", "secret"),
				},
				{
					ResourceKind: "ConfigMap",
					Resource:     kube.ObjectKey("test", "configmap"),
				},
			},

			want: map[WatchedObject]map[client.ObjectKey]bool{
				{
					ResourceKind: "ConfigMap",
					Resource:     kube.ObjectKey("test", "configmap"),
				}: {
					kube.ObjectKey("test", "project1"): true,
				},
				{
					ResourceKind: "Secret",
					Resource:     kube.ObjectKey("test", "secret"),
				}: {
					kube.ObjectKey("test", "project1"): true,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := NewDeprecatedResourceWatcher()
			r.EnsureMultiplesResourcesAreWatched(tc.dependant, zap.S(), tc.resources...)
			got := r.WatchedResourcesSnapshot()
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want watched resources %v, got %v", tc.want, got)
			}
		})
	}
}

func TestEnsureResourcesAreWatched(t *testing.T) {
	t.Run("One secret is watched by two resources", func(t *testing.T) {
		watcher := NewDeprecatedResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		project2 := kube.ObjectKey("test", "project2")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")

		var wg sync.WaitGroup

		wg.Add(2)
		go func() {
			defer wg.Done()
			watcher.EnsureResourcesAreWatched(project1, "Secret", zap.S(), connectionSecret)
		}()
		go func() {
			defer wg.Done()
			watcher.EnsureResourcesAreWatched(project2, "Secret", zap.S(), connectionSecret)
		}()
		wg.Wait()
		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}: {project1: true, project2: true},
		}
		assert.Equal(t, expectedWatched, watcher.WatchedResourcesSnapshot())
	})
	t.Run("One resource watches two secrets", func(t *testing.T) {
		watcher := NewDeprecatedResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")
		connectionSecret2 := kube.ObjectKey("test", "connectionSecret2")

		watcher.EnsureResourcesAreWatched(project1, "Secret", zap.S(), connectionSecret, connectionSecret2)

		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}:  {project1: true},
			{ResourceKind: "Secret", Resource: connectionSecret2}: {project1: true},
		}
		assert.Equal(t, expectedWatched, watcher.WatchedResourcesSnapshot())
	})
	t.Run("Resource stops watching one secret", func(t *testing.T) {
		watcher := NewDeprecatedResourceWatcher()
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
		assert.Equal(t, expectedWatched, watcher.WatchedResourcesSnapshot())
	})
	t.Run("Watcher to watch multiple resources of different kinds", func(t *testing.T) {
		watcher := NewDeprecatedResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")
		backupSchedule := kube.ObjectKey("test", "backupSchedule")

		watcher.EnsureMultiplesResourcesAreWatched(project1, zap.S(),
			WatchedObject{
				ResourceKind: "Secret",
				Resource:     connectionSecret,
			},
			WatchedObject{
				ResourceKind: "AtlasBackupSchedule",
				Resource:     backupSchedule,
			},
		)

		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}:            {project1: true},
			{ResourceKind: "AtlasBackupSchedule", Resource: backupSchedule}: {project1: true},
		}

		assert.Equal(t, expectedWatched, watcher.WatchedResourcesSnapshot())
	})
	t.Run("Watcher to watch multiple resources of different kinds for multiple projects", func(t *testing.T) {
		watcher := NewDeprecatedResourceWatcher()
		project1 := kube.ObjectKey("test", "project1")
		project2 := kube.ObjectKey("test", "project2")
		connectionSecret := kube.ObjectKey("test", "connectionSecret")
		backupSchedule := kube.ObjectKey("test", "backupSchedule")

		watcher.EnsureMultiplesResourcesAreWatched(project1, zap.S(),
			WatchedObject{
				ResourceKind: "Secret",
				Resource:     connectionSecret,
			},
			WatchedObject{
				ResourceKind: "AtlasBackupSchedule",
				Resource:     backupSchedule,
			},
		)

		watcher.EnsureMultiplesResourcesAreWatched(project2, zap.S(),
			WatchedObject{
				ResourceKind: "Secret",
				Resource:     connectionSecret,
			},
			WatchedObject{
				ResourceKind: "AtlasBackupSchedule",
				Resource:     backupSchedule,
			},
		)

		expectedWatched := map[WatchedObject]map[client.ObjectKey]bool{
			{ResourceKind: "Secret", Resource: connectionSecret}:            {project1: true, project2: true},
			{ResourceKind: "AtlasBackupSchedule", Resource: backupSchedule}: {project1: true, project2: true},
		}

		assert.Equal(t, expectedWatched, watcher.WatchedResourcesSnapshot())
	})
	// TODO: add test for different kind of resources
}
