package atlasbackupcompliancepolicy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func TestReconcile(t *testing.T) {
	t.Run("should terminate silently when resource is not found", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "bcp",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})
}

func TestEnsureAtlasBackupCompliancePolicy(t *testing.T) {
	t.Run("should skip reconciliation when annotation is set", func(t *testing.T) {
		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bcp",
				Namespace: "default",
				Annotations: map[string]string{
					customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			WithStatusSubresource(bcp).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}
		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-bcp",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})

	t.Run("should transition to error state when resource version is invalid", func(t *testing.T) {
		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bcp",
				Namespace: "default",
				Labels: map[string]string{
					customresource.ResourceVersion: "blah",
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			WithStatusSubresource(bcp).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}
		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-bcp",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			result,
		)
	})

	t.Run("should transition to error state when resource is unsupported", func(t *testing.T) {
		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bcp",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			WithStatusSubresource(bcp).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return false
				},
			},
		}
		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-bcp",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{},
			result,
		)
	})

	t.Run("should lock when there are references", func(t *testing.T) {
		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bcp",
				Namespace: "default",
			},
			Spec: akov2.AtlasBackupCompliancePolicySpec{
				AuthorizedEmail:         "test@example.com",
				AuthorizedUserFirstName: "John",
				AuthorizedUserLastName:  "Doe",
				CopyProtectionEnabled:   false,
				EncryptionAtRestEnabled: false,
				PITEnabled:              false,
				RestoreWindowDays:       42,
				ScheduledPolicyItems: []akov2.AtlasBackupPolicyItem{
					{
						FrequencyType:     "monthly",
						FrequencyInterval: 4,
						RetentionUnit:     "months",
						RetentionValue:    1,
					},
				},
				OnDemandPolicy: akov2.AtlasOnDemandPolicy{
					RetentionUnit:  "weeks",
					RetentionValue: 3,
				},
			},
		}

		project := akov2.DefaultProject("default", "connection-secret").WithBackupCompliancePolicyNamespaced("my-bcp", "default")
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		bcpIndexer := indexer.NewAtlasProjectByBackupCompliancePolicyIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp, project).
			WithStatusSubresource(bcp).
			WithIndex(
				bcpIndexer.Object(),
				bcpIndexer.Name(),
				bcpIndexer.Keys,
			).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}
		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-bcp",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{},
			result,
		)
	})

	t.Run("should release when there are no references", func(t *testing.T) {
		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "my-bcp",
				Namespace:  "default",
				Finalizers: []string{customresource.FinalizerLabel},
			},
			Spec: akov2.AtlasBackupCompliancePolicySpec{
				AuthorizedEmail:         "test@example.com",
				AuthorizedUserFirstName: "John",
				AuthorizedUserLastName:  "Doe",
				CopyProtectionEnabled:   false,
				EncryptionAtRestEnabled: false,
				PITEnabled:              false,
				RestoreWindowDays:       42,
				ScheduledPolicyItems: []akov2.AtlasBackupPolicyItem{
					{
						FrequencyType:     "monthly",
						FrequencyInterval: 4,
						RetentionUnit:     "months",
						RetentionValue:    1,
					},
				},
				OnDemandPolicy: akov2.AtlasOnDemandPolicy{
					RetentionUnit:  "weeks",
					RetentionValue: 3,
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		bcpIndexer := indexer.NewAtlasProjectByBackupCompliancePolicyIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			WithStatusSubresource(bcp).
			WithIndex(
				bcpIndexer.Object(),
				bcpIndexer.Name(),
				bcpIndexer.Keys,
			).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}
		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-bcp",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{},
			result,
		)
	})
}

func TestFindBCPForProjects(t *testing.T) {
	t.Run("should return a slice of requests for BCP", func(t *testing.T) {
		project := akov2.NewProject("project1", "default", "connection-secret").WithBackupCompliancePolicyNamespaced("bcp1", "other-ns")

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		reqs := reconciler.findBCPForProjects(context.Background(), project)
		assert.Equal(
			t,
			[]ctrl.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      "bcp1",
						Namespace: "other-ns",
					},
				},
			},
			reqs,
		)
	})
	t.Run("should return nil when no bcp specified", func(t *testing.T) {
		project := akov2.NewProject("project1", "default", "connection-secret")

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			Build()

		reconciler := &AtlasBackupCompliancePolicyReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		reqs := reconciler.findBCPForProjects(context.Background(), project)

		assert.Equal(t, nil, reqs)
	})
}
