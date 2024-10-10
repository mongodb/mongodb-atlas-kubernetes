package connectionsecret_test

import (
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
)

const (
	testProjectID = "123456"

	testNamespace = "some-namespace"
)

func TestReapOrphanConnectionSecrets(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))

	for _, tc := range []struct {
		title            string
		deployments      []string
		objects          []client.Object
		expectedErr      error
		expectedRemovals []string
	}{
		{
			title:            "Empty list of secrets returns empty list of removals",
			expectedRemovals: []string{},
		},

		{
			title:            "Matching secrets do not get removed",
			deployments:      sampleDeployments(),
			objects:          matchingSecrets(),
			expectedRemovals: []string{},
		},

		{
			title:            "Secrets to non existing clusters get removed",
			deployments:      sampleDeployments(),
			objects:          merge(matchingSecrets(), nonMatchingSecrets()),
			expectedRemovals: namesOf(nonMatchingSecrets()),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.objects...).Build()
			removedOrphans, err := connectionsecret.ReapOrphanConnectionSecrets(
				context.Background(),
				fakeClient,
				testProjectID,
				testNamespace,
				tc.deployments,
			)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRemovals, removedOrphans)
		})
	}
}

func sampleDeployments() []string {
	return []string{"cluster1", "serverless2"}
}

func matchingSecrets() []client.Object {
	return []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret1",
				Namespace: testNamespace,
				Labels: map[string]string{
					connectionsecret.ClusterLabelKey: "cluster1",
					connectionsecret.ProjectLabelKey: testProjectID,
					connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
				},
			},
		},

		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret2",
				Namespace: testNamespace,
				Labels: map[string]string{
					connectionsecret.ClusterLabelKey: "serverless2",
					connectionsecret.ProjectLabelKey: testProjectID,
					connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
				},
			},
		},
	}
}

func nonMatchingSecrets() []client.Object {
	return []client.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret3",
				Namespace: testNamespace,
				Labels: map[string]string{
					connectionsecret.ClusterLabelKey: "cluster3",
					connectionsecret.ProjectLabelKey: testProjectID,
					connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
				},
			},
		},

		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret4",
				Namespace: testNamespace,
				Labels: map[string]string{
					connectionsecret.ClusterLabelKey: "serverless4",
					connectionsecret.ProjectLabelKey: testProjectID,
					connectionsecret.TypeLabelKey:    connectionsecret.CredLabelVal,
				},
			},
		},
	}
}

func namesOf(objs []client.Object) []string {
	names := make([]string, 0, len(objs))
	for _, obj := range objs {
		names = append(names, fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()))
	}
	return names
}

func merge(objs ...[]client.Object) []client.Object {
	if len(objs) == 0 {
		return []client.Object{}
	}
	result := objs[0]
	for i := 1; i < len(objs); i++ {
		result = append(result, objs[i]...)
	}
	return result
}
