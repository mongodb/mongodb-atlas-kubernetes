package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
)

func TestCanEncryptionAtRestReconcile(t *testing.T) {
	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, false, &mdbv1.AtlasProject{})
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		workflowCtx := &workflow.Context{
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, true, akoProject)
		require.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		require.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, true, akoProject)

		require.EqualError(t, err, "failed to retrieve data")
		require.False(t, result)
	})

	t.Run("should return true when all providers are disabled in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return &mongodbatlas.EncryptionAtRest{
						AwsKms: mongodbatlas.AwsKms{
							Enabled: toptr.MakePtr(false),
						},
						AzureKeyVault: mongodbatlas.AzureKeyVault{
							Enabled: toptr.MakePtr(false),
						},
						GoogleCloudKms: mongodbatlas.GoogleCloudKms{
							Enabled: toptr.MakePtr(false),
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return &mongodbatlas.EncryptionAtRest{
						AwsKms: mongodbatlas.AwsKms{
							Enabled:             toptr.MakePtr(true),
							CustomerMasterKeyID: "aws-kms-master-key",
							Region:              "eu-west-1",
							Valid:               toptr.MakePtr(true),
						},
						AzureKeyVault: mongodbatlas.AzureKeyVault{
							Enabled: toptr.MakePtr(false),
						},
						GoogleCloudKms: mongodbatlas.GoogleCloudKms{
							Enabled: toptr.MakePtr(false),
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AwsKms: mdbv1.AwsKms{
						Enabled: toptr.MakePtr(true),
						Region:  "eu-west-2",
						SecretRef: common.ResourceRefNamespaced{
							Name:      "test-aws-secret",
							Namespace: "test-namespace",
						},
					},
					AzureKeyVault:  mdbv1.AzureKeyVault{},
					GoogleCloudKms: mdbv1.GoogleCloudKms{},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"encryptionAtRest":{"awsKms":{"enabled":true,"region":"eu-west-1","secretRef":{"name":"test-aws-secret","namespace":"test-namespace"}},"azureKeyVault":{},"googleCloudKms":{}}}`,
			},
		)
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return &mongodbatlas.EncryptionAtRest{
						AwsKms: mongodbatlas.AwsKms{
							Enabled:             toptr.MakePtr(true),
							CustomerMasterKeyID: "aws-kms-master-key",
							Region:              "eu-west-1",
							Valid:               toptr.MakePtr(true),
						},
						AzureKeyVault: mongodbatlas.AzureKeyVault{
							Enabled: toptr.MakePtr(false),
						},
						GoogleCloudKms: mongodbatlas.GoogleCloudKms{
							Enabled: toptr.MakePtr(false),
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AwsKms: mdbv1.AwsKms{
						Enabled: toptr.MakePtr(true),
						Region:  "eu-west-1",
						SecretRef: common.ResourceRefNamespaced{
							Name:      "test-aws-secret",
							Namespace: "test-namespace",
						},
					},
					AzureKeyVault:  mdbv1.AzureKeyVault{},
					GoogleCloudKms: mdbv1.GoogleCloudKms{},
				},
			},
		}
		akoProject.Spec.EncryptionAtRest.AwsKms.SetSecrets("aws-kms-master-key", "aws:id:arn/my-role")
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"encryptionAtRest":{"awsKms":{"enabled":true,"region":"eu-west-2","secretRef":{"name":"test-aws-secret","namespace":"test-namespace"}},"azureKeyVault":{},"googleCloudKms":{}}}`,
			},
		)
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("should return false when unable to reconcile Encryption at Rest", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return &mongodbatlas.EncryptionAtRest{
						AwsKms: mongodbatlas.AwsKms{
							Enabled:             toptr.MakePtr(true),
							CustomerMasterKeyID: "aws-kms-master-key",
							Region:              "eu-west-1",
							Valid:               toptr.MakePtr(true),
						},
						AzureKeyVault: mongodbatlas.AzureKeyVault{
							Enabled: toptr.MakePtr(false),
						},
						GoogleCloudKms: mongodbatlas.GoogleCloudKms{
							Enabled: toptr.MakePtr(false),
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AwsKms: mdbv1.AwsKms{
						Enabled: toptr.MakePtr(true),
						Region:  "eu-central-1",
						SecretRef: common.ResourceRefNamespaced{
							Name:      "test-aws-secret",
							Namespace: "test-namespace",
						},
					},
					AzureKeyVault:  mdbv1.AzureKeyVault{},
					GoogleCloudKms: mdbv1.GoogleCloudKms{},
				},
			},
		}
		akoProject.Spec.EncryptionAtRest.AwsKms.SetSecrets("aws-kms-master-key", "aws:id:arn/my-role")
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"encryptionAtRest":{"awsKms":{"enabled":true,"region":"eu-west-2","secretRef":{"name":"test-aws-secret","namespace":"test-namespace"}},"azureKeyVault":{},"googleCloudKms":{}}}`,
			},
		)
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canEncryptionAtRestReconcile(workflowCtx, true, akoProject)

		require.NoError(t, err)
		require.False(t, result)
	})
}

func TestEnsureEncryptionAtRest(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		reconciler := &AtlasProjectReconciler{
			SubObjectDeletionProtection: true,
		}
		result := reconciler.ensureEncryptionAtRest(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			EncryptionsAtRest: &atlas.EncryptionAtRestClientMock{
				GetFunc: func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
					return &mongodbatlas.EncryptionAtRest{
						AwsKms: mongodbatlas.AwsKms{
							Enabled:             toptr.MakePtr(true),
							CustomerMasterKeyID: "aws-kms-master-key",
							Region:              "eu-west-1",
							Valid:               toptr.MakePtr(true),
						},
						AzureKeyVault: mongodbatlas.AzureKeyVault{
							Enabled: toptr.MakePtr(false),
						},
						GoogleCloudKms: mongodbatlas.GoogleCloudKms{
							Enabled: toptr.MakePtr(false),
						},
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AwsKms: mdbv1.AwsKms{
						Enabled: toptr.MakePtr(true),
						Region:  "eu-central-1",
						SecretRef: common.ResourceRefNamespaced{
							Name:      "test-aws-secret",
							Namespace: "test-namespace",
						},
					},
					AzureKeyVault:  mdbv1.AzureKeyVault{},
					GoogleCloudKms: mdbv1.GoogleCloudKms{},
				},
			},
		}
		akoProject.WithAnnotations(
			map[string]string{
				customresource.AnnotationLastAppliedConfiguration: `{"encryptionAtRest":{"awsKms":{"enabled":true,"region":"eu-west-2","secretRef":{"name":"test-aws-secret","namespace":"test-namespace"}},"azureKeyVault":{},"googleCloudKms":{}}}`,
			},
		)
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		secretData := map[string][]byte{
			"AccessKeyID":         []byte("testAccessKeyID"),
			"SecretAccessKey":     []byte("testSecretAccessKey"),
			"CustomerMasterKeyID": []byte("aws-kms-master-key"),
			"RoleID":              []byte("aws:id:arn/my-role"),
		}

		fakeClient := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-aws-secret",
					Namespace: "test-namespace",
				},
			},
		}...).Build()

		reconciler := &AtlasProjectReconciler{
			SubObjectDeletionProtection: true,
			Client:                      fakeClient,
		}
		result := reconciler.ensureEncryptionAtRest(workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}

func TestReadEncryptionAtRestSecrets(t *testing.T) {
	t.Run("AWS with correct secret data", func(t *testing.T) {
		secretData := map[string][]byte{
			"AccessKeyID":         []byte("testAccessKeyID"),
			"SecretAccessKey":     []byte("testSecretAccessKey"),
			"CustomerMasterKeyID": []byte("testCustomerMasterKeyID"),
			"RoleID":              []byte("testRoleID"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-secret",
					Namespace: "test",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			AwsKms: mdbv1.AwsKms{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name:      "aws-secret",
					Namespace: "test",
				},
				Region: "testRegion",
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["CustomerMasterKeyID"]), encRest.AwsKms.CustomerMasterKeyID())
		assert.Equal(t, string(secretData["RoleID"]), encRest.AwsKms.RoleID())
	})

	t.Run("AWS with correct secret data (fallback namespace)", func(t *testing.T) {
		secretData := map[string][]byte{
			"AccessKeyID":         []byte("testAccessKeyID"),
			"SecretAccessKey":     []byte("testSecretAccessKey"),
			"CustomerMasterKeyID": []byte("testCustomerMasterKeyID"),
			"RoleID":              []byte("testRoleID"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-secret",
					Namespace: "test-fallback-ns",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			AwsKms: mdbv1.AwsKms{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name: "aws-secret",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test-fallback-ns")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["CustomerMasterKeyID"]), encRest.AwsKms.CustomerMasterKeyID())
		assert.Equal(t, string(secretData["RoleID"]), encRest.AwsKms.RoleID())
	})

	t.Run("AWS with missing fields", func(t *testing.T) {
		secretData := map[string][]byte{
			"AccessKeyID":         []byte("testKeyID"),
			"SecretAccessKey":     []byte("testSecretAccessKey"),
			"CustomerMasterKeyID": []byte("testCustomerMasterKeyID"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-secret",
					Namespace: "test",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			AwsKms: mdbv1.AwsKms{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name:      "aws-secret",
					Namespace: "test",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.NotNil(t, err)
	})

	t.Run("GCP with correct secret data", func(t *testing.T) {
		secretData := map[string][]byte{
			"ServiceAccountKey":    []byte("testServiceAccountKey"),
			"KeyVersionResourceID": []byte("testKeyVersionResourceID"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gcp-secret",
					Namespace: "test",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			GoogleCloudKms: mdbv1.GoogleCloudKms{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name: "gcp-secret",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["ServiceAccountKey"]), encRest.GoogleCloudKms.ServiceAccountKey())
		assert.Equal(t, string(secretData["KeyVersionResourceID"]), encRest.GoogleCloudKms.KeyVersionResourceID())
	})

	t.Run("GCP with missing fields", func(t *testing.T) {
		secretData := map[string][]byte{
			"ServiceAccountKey": []byte("testServiceAccountKey"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gcp-secret",
					Namespace: "test",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			GoogleCloudKms: mdbv1.GoogleCloudKms{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name: "gcp-secret",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.NotNil(t, err)
	})

	t.Run("Azure with correct secret data", func(t *testing.T) {
		secretData := map[string][]byte{
			"Secret":         []byte("testClientSecret"),
			"SubscriptionID": []byte("testSubscriptionID"),
			"KeyVaultName":   []byte("testKeyVaultName"),
			"KeyIdentifier":  []byte("testKeyIdentifier"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "azure-secret",
					Namespace: "test",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			AzureKeyVault: mdbv1.AzureKeyVault{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name: "azure-secret",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["Secret"]), encRest.AzureKeyVault.Secret())
		assert.Equal(t, string(secretData["SubscriptionID"]), encRest.AzureKeyVault.SubscriptionID())
		assert.Equal(t, string(secretData["KeyVaultName"]), encRest.AzureKeyVault.KeyVaultName())
		assert.Equal(t, string(secretData["KeyIdentifier"]), encRest.AzureKeyVault.KeyIdentifier())
	})

	t.Run("Azure with missing fields", func(t *testing.T) {
		secretData := map[string][]byte{
			"ClientID":          []byte("testClientID"),
			"AzureEnvironment":  []byte("testAzureEnvironment"),
			"SubscriptionID":    []byte("testSubscriptionID"),
			"ResourceGroupName": []byte("testResourceGroupName"),
		}

		kk := fake.NewClientBuilder().WithRuntimeObjects([]runtime.Object{
			&v1.Secret{
				Data: secretData,
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gcp-secret",
					Namespace: "test",
				},
			},
		}...).Build()

		service := &workflow.Context{}

		encRest := &mdbv1.EncryptionAtRest{
			AzureKeyVault: mdbv1.AzureKeyVault{
				Enabled: toptr.MakePtr(true),
				SecretRef: common.ResourceRefNamespaced{
					Name: "gcp-secret",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.NotNil(t, err)
	})
}

func TestIsEncryptionAtlasEmpty(t *testing.T) {
	spec := &mdbv1.EncryptionAtRest{}
	isEmpty := IsEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Empty spec should be empty")

	spec.AwsKms.Enabled = toptr.MakePtr(true)
	isEmpty = IsEncryptionSpecEmpty(spec)
	assert.False(t, isEmpty, "Non-empty spec")

	spec.AwsKms.Enabled = toptr.MakePtr(false)
	isEmpty = IsEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Enabled flag set to false is same as empty")
}

func TestAtlasInSync(t *testing.T) {
	areInSync, err := AtlasInSync(nil, nil)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both atlas and spec are nil")

	groupID := "0"
	atlas := mongodbatlas.EncryptionAtRest{
		GroupID: groupID,
		AwsKms: mongodbatlas.AwsKms{
			Enabled: toptr.MakePtr(true),
		},
	}
	spec := mdbv1.EncryptionAtRest{
		AwsKms: mdbv1.AwsKms{
			Enabled: toptr.MakePtr(true),
		},
	}

	areInSync, err = AtlasInSync(nil, &spec)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Nil atlas")

	areInSync, err = AtlasInSync(&atlas, nil)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Nil spec")

	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are the same")

	spec.AwsKms.Enabled = toptr.MakePtr(false)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Atlas is disabled")

	atlas.AwsKms.Enabled = toptr.MakePtr(false)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are disabled")

	atlas.AwsKms.RoleID = "example"
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are disabled but atlas RoleID field")

	spec.AwsKms.Enabled = toptr.MakePtr(true)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.False(t, areInSync, "Spec is re-enabled")

	atlas.AwsKms.Enabled = toptr.MakePtr(true)
	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Both are re-enabled and only RoleID is different")

	atlas = mongodbatlas.EncryptionAtRest{
		AwsKms: mongodbatlas.AwsKms{
			Enabled:             toptr.MakePtr(true),
			CustomerMasterKeyID: "testCustomerMasterKeyID",
			Region:              "US_EAST_1",
			RoleID:              "testRoleID",
			Valid:               toptr.MakePtr(true),
		},
		AzureKeyVault: mongodbatlas.AzureKeyVault{
			Enabled: toptr.MakePtr(false),
		},
		GoogleCloudKms: mongodbatlas.GoogleCloudKms{
			Enabled: toptr.MakePtr(false),
		},
	}
	spec = mdbv1.EncryptionAtRest{
		AwsKms: mdbv1.AwsKms{
			Enabled: toptr.MakePtr(true),
			Region:  "US_EAST_1",
			Valid:   toptr.MakePtr(true),
		},
		AzureKeyVault:  mdbv1.AzureKeyVault{},
		GoogleCloudKms: mdbv1.GoogleCloudKms{},
	}
	spec.AwsKms.SetSecrets("testCustomerMasterKeyID", "testRoleID")

	areInSync, err = AtlasInSync(&atlas, &spec)
	assert.NoError(t, err)
	assert.True(t, areInSync, "Realistic example. should be equal")
}

func TestAreAzureConfigEqual(t *testing.T) {
	type args struct {
		operator mdbv1.AzureKeyVault
		atlas    mongodbatlas.AzureKeyVault
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Azure configuration are equal",
			args: args{
				operator: mdbv1.AzureKeyVault{
					Enabled:           toptr.MakePtr(true),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					ResourceGroupName: "resource group",
					TenantID:          "tenant id",
				},
				atlas: mongodbatlas.AzureKeyVault{
					Enabled:           toptr.MakePtr(true),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					SubscriptionID:    "sub id",
					ResourceGroupName: "resource group",
					KeyVaultName:      "vault name",
					KeyIdentifier:     "key id",
					TenantID:          "tenant id",
				},
			},
			want: true,
		},
		{
			name: "Azure configuration are equal when disabled and nullable",
			args: args{
				operator: mdbv1.AzureKeyVault{
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					ResourceGroupName: "resource group",
					TenantID:          "tenant id",
				},
				atlas: mongodbatlas.AzureKeyVault{
					Enabled:           toptr.MakePtr(false),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					SubscriptionID:    "sub id",
					ResourceGroupName: "resource group",
					KeyVaultName:      "vault name",
					KeyIdentifier:     "key id",
					TenantID:          "tenant id",
				},
			},
			want: true,
		},
		{
			name: "Azure configuration differ by enabled field",
			args: args{
				operator: mdbv1.AzureKeyVault{
					Enabled:           toptr.MakePtr(false),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					ResourceGroupName: "resource group",
					TenantID:          "tenant id",
				},
				atlas: mongodbatlas.AzureKeyVault{
					Enabled:           toptr.MakePtr(true),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					SubscriptionID:    "sub id",
					ResourceGroupName: "resource group",
					KeyVaultName:      "vault name",
					KeyIdentifier:     "key id",
					TenantID:          "tenant id",
				},
			},
			want: false,
		},
		{
			name: "Azure configuration differ by other field",
			args: args{
				operator: mdbv1.AzureKeyVault{
					Enabled:           toptr.MakePtr(true),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					ResourceGroupName: "resource group",
					TenantID:          "tenant id",
				},
				atlas: mongodbatlas.AzureKeyVault{
					Enabled:           toptr.MakePtr(true),
					ClientID:          "client id",
					AzureEnvironment:  "azure env",
					SubscriptionID:    "sub id",
					ResourceGroupName: "resource group name",
					KeyVaultName:      "vault name",
					KeyIdentifier:     "key id",
					TenantID:          "tenant id",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.operator.SetSecrets("sub id", "vault name", "key id", "")
			assert.Equalf(t, tt.want, areAzureConfigEqual(tt.args.operator, tt.args.atlas, false), "areAzureConfigEqual(%v, %v)", tt.args.operator, tt.args.atlas)
		})
	}
}

func TestAreGCPConfigEqual(t *testing.T) {
	type args struct {
		operator mdbv1.GoogleCloudKms
		atlas    mongodbatlas.GoogleCloudKms
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "GCP configuration are equal",
			args: args{
				operator: mdbv1.GoogleCloudKms{
					Enabled: toptr.MakePtr(true),
				},
				atlas: mongodbatlas.GoogleCloudKms{
					Enabled:              toptr.MakePtr(true),
					KeyVersionResourceID: "key version id",
				},
			},
			want: true,
		},
		{
			name: "GCP configuration are equal when disabled and nullable",
			args: args{
				operator: mdbv1.GoogleCloudKms{},
				atlas: mongodbatlas.GoogleCloudKms{
					Enabled:              toptr.MakePtr(false),
					KeyVersionResourceID: "key version id",
				},
			},
			want: true,
		},
		{
			name: "GCP configuration are different by enable field",
			args: args{
				operator: mdbv1.GoogleCloudKms{
					Enabled: toptr.MakePtr(true),
				},
				atlas: mongodbatlas.GoogleCloudKms{
					Enabled:              toptr.MakePtr(false),
					KeyVersionResourceID: "key version id",
				},
			},
			want: false,
		},
		{
			name: "GCP configuration are different by another field",
			args: args{
				operator: mdbv1.GoogleCloudKms{
					Enabled: toptr.MakePtr(true),
				},
				atlas: mongodbatlas.GoogleCloudKms{
					Enabled:              toptr.MakePtr(true),
					KeyVersionResourceID: "key version resource id",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.operator.SetSecrets("", "key version id")
			assert.Equalf(t, tt.want, areGCPConfigEqual(tt.args.operator, tt.args.atlas, false), "areGCPConfigEqual(%v, %v)", tt.args.operator, tt.args.atlas)
		})
	}
}

func TestAreAWSConfigEqual(t *testing.T) {
	type args struct {
		operator mdbv1.AwsKms
		atlas    mongodbatlas.AwsKms
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "AWS configuration are equal",
			args: args{
				operator: mdbv1.AwsKms{
					Enabled: toptr.MakePtr(true),
				},
				atlas: mongodbatlas.AwsKms{
					Enabled:             toptr.MakePtr(true),
					CustomerMasterKeyID: "customer master key",
				},
			},
			want: true,
		},
		{
			name: "AWS configuration are equal when disabled and nullable",
			args: args{
				operator: mdbv1.AwsKms{},
				atlas: mongodbatlas.AwsKms{
					Enabled:             toptr.MakePtr(false),
					CustomerMasterKeyID: "customer master key",
				},
			},
			want: true,
		},
		{
			name: "AWS configuration are different by enable field",
			args: args{
				operator: mdbv1.AwsKms{
					Enabled: toptr.MakePtr(true),
				},
				atlas: mongodbatlas.AwsKms{
					Enabled:             toptr.MakePtr(false),
					CustomerMasterKeyID: "customer master key",
				},
			},
			want: false,
		},
		{
			name: "AWS configuration are different by another field",
			args: args{
				operator: mdbv1.AwsKms{
					Enabled: toptr.MakePtr(true),
				},
				atlas: mongodbatlas.AwsKms{
					Enabled:             toptr.MakePtr(true),
					CustomerMasterKeyID: "customer master key id",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.operator.SetSecrets("customer master key", "")
			assert.Equalf(t, tt.want, areAWSConfigEqual(tt.args.operator, tt.args.atlas, false), "areGCPConfigEqual(%v, %v)", tt.args.operator, tt.args.atlas)
		})
	}
}
