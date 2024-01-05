package statushandler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/kube"
)

func Test_PatchUpdateStatus(t *testing.T) {
	existingProject := &mdbv1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-project",
			Namespace: "test-ns",
		},
		Status: status.AtlasProjectStatus{
			Common: status.Common{Conditions: []status.Condition{{
				Type:   status.IPAccessListReadyType,
				Status: corev1.ConditionFalse,
			}}},
		},
	}
	// Fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(mdbv1.AddToScheme(scheme))
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existingProject).
		WithStatusSubresource(existingProject).Build()

	// Patch the existing project
	updatedProject := existingProject.DeepCopy()
	updatedProject.Status.Common.Conditions[0].Status = corev1.ConditionTrue
	updatedProject.Status.ID = "theId"
	assert.NoError(t, patchUpdateStatus(context.Background(), fakeClient, updatedProject))

	projectAfterPatch := &mdbv1.AtlasProject{}
	assert.NoError(t, fakeClient.Get(context.Background(), kube.ObjectKeyFromObject(updatedProject), projectAfterPatch))
	assert.Equal(t, []status.Condition{{Type: status.IPAccessListReadyType, Status: corev1.ConditionTrue}}, projectAfterPatch.Status.Common.Conditions)
	assert.Equal(t, "theId", projectAfterPatch.Status.ID)
}
