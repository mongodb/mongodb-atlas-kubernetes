package status

import (
	"context"
	"testing"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_PatchUpdateStatus(t *testing.T) {
	existingProject := &mdbv1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-project",
			Namespace: "test-ns",
		},
		Status: mdbv1.AtlasProjectStatus{
			Common: status.Common{Phase: status.PhaseReconciling},
		},
	}
	// Fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(mdbv1.AddToScheme(scheme))
	fakeClient := fake.NewFakeClientWithScheme(scheme, existingProject)

	// Patch the existing project
	updatedProject := existingProject.DeepCopy()
	updatedProject.Status.Common.Phase = status.PhasePending
	updatedProject.Status.ID = "theId"
	assert.NoError(t, PatchUpdateStatus(fakeClient, updatedProject))

	projectAfterPatch := &mdbv1.AtlasProject{}
	assert.NoError(t, fakeClient.Get(context.Background(), kube.ObjectKeyFromObject(updatedProject), projectAfterPatch))
	assert.Equal(t, status.PhasePending, projectAfterPatch.Status.Common.Phase)
	assert.Equal(t, "theId", projectAfterPatch.Status.ID)
}
