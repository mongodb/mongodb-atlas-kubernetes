package watch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllertest"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func TestHandleCreate(t *testing.T) {
	t.Run("Create event is not handled", func(t *testing.T) {
		secret := secretForTesting("testSecret")
		watcher := NewDeprecatedResourceWatcher()
		handler := NewSecretHandler(&watcher)
		createEvent := event.CreateEvent{Object: secret}
		queue := controllertest.Queue{TypedInterface: workqueue.NewTyped[reconcile.Request]()}

		handler.Create(context.Background(), createEvent, &queue)
		assert.Zero(t, queue.Len())
	})
	t.Run("Create event is handled", func(t *testing.T) {
		secret := secretForTesting("testSecret")
		dependentResourceKey := kube.ObjectKey("ns", "testAtlasProject")
		watcher := NewDeprecatedResourceWatcher()
		watcher.EnsureMultiplesResourcesAreWatched(dependentResourceKey, zap.S(), WatchedObject{
			ResourceKind: secret.GetObjectKind().GroupVersionKind().Kind,
			Resource:     kube.ObjectKeyFromObject(secret),
		})
		handler := NewSecretHandler(&watcher)

		createEvent := event.CreateEvent{Object: secret}
		queue := controllertest.Queue{TypedInterface: workqueue.NewTyped[reconcile.Request]()}

		handler.Create(context.Background(), createEvent, &queue)
		assert.Equal(t, queue.Len(), 1)

		enqueued, _ := queue.Get()

		// We expect the "dependent" resource to appear in the queue
		assert.Equal(t, reconcile.Request{NamespacedName: dependentResourceKey}, enqueued)
	})
}

func TestHandleUpdate(t *testing.T) {
	t.Run("Update event is not handled", func(t *testing.T) {
		// Update event is not handled as the Secret that triggered the update event is not a watched one
		watchedSecret := secretForTesting("someOtherSecret")
		dependentResourceKey := kube.ObjectKey("ns", "testAtlasProject")

		oldSecret := secretForTesting("testSecret")
		newSecret := oldSecret.DeepCopy()
		newSecret.Data["secondKey"] = []byte("secondValue")
		watcher := NewDeprecatedResourceWatcher()
		watcher.EnsureMultiplesResourcesAreWatched(dependentResourceKey, zap.S(), WatchedObject{
			ResourceKind: watchedSecret.GetObjectKind().GroupVersionKind().Kind,
			Resource:     kube.ObjectKeyFromObject(watchedSecret),
		})
		handler := NewSecretHandler(&watcher)
		updateEvent := event.UpdateEvent{ObjectOld: oldSecret, ObjectNew: newSecret}
		queue := controllertest.Queue{TypedInterface: workqueue.NewTyped[reconcile.Request]()}

		handler.Update(context.Background(), updateEvent, &queue)
		assert.Zero(t, queue.Len())
	})
	t.Run("Update event is handled", func(t *testing.T) {
		secret := secretForTesting("testSecret")
		dependentResourceKey := kube.ObjectKey("ns", "testAtlasProject")

		oldSecret := secretForTesting("testSecret")
		newSecret := oldSecret.DeepCopy()
		newSecret.Data["secondKey"] = []byte("secondValue")

		watcher := NewDeprecatedResourceWatcher()
		watcher.EnsureMultiplesResourcesAreWatched(dependentResourceKey, zap.S(), WatchedObject{
			ResourceKind: secret.GetObjectKind().GroupVersionKind().Kind,
			Resource:     kube.ObjectKeyFromObject(secret),
		})
		handler := NewSecretHandler(&watcher)

		updateEvent := event.UpdateEvent{ObjectOld: oldSecret, ObjectNew: newSecret}
		queue := controllertest.Queue{TypedInterface: workqueue.NewTyped[reconcile.Request]()}

		handler.Update(context.Background(), updateEvent, &queue)
		assert.Equal(t, queue.Len(), 1)

		enqueued, _ := queue.Get()

		assert.Equal(t, reconcile.Request{NamespacedName: dependentResourceKey}, enqueued)
	})
}

func TestShouldHandleUpdate(t *testing.T) {
	t.Run("Update shouldn't happen if Secrets data hasn't changed", func(t *testing.T) {
		oldObj := secretForTesting("testValue")
		newObj := oldObj.DeepCopy()
		newObj.ObjectMeta.ResourceVersion = "4243"

		assert.False(t, shouldHandleUpdate(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj}))
	})
	t.Run("Update should happen if the data has changed for Secret", func(t *testing.T) {
		oldObj := secretForTesting("testValue")
		newObj := oldObj.DeepCopy()
		newObj.ObjectMeta.ResourceVersion = "4243"
		newObj.Data["secondKey"] = []byte("secondValue")

		assert.True(t, shouldHandleUpdate(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj}))
	})
}

func secretForTesting(name string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "ns",
		},
		Data: map[string][]byte{"testKey": []byte("testValue")},
	}
}
