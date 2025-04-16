// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/encryptionatrest"
)

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

		encRest := &encryptionatrest.EncryptionAtRest{
			AWS: encryptionatrest.AwsKms{
				AwsKms: akov2.AwsKms{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name:      "aws-secret",
						Namespace: "test",
					},
					Region: "testRegion",
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["CustomerMasterKeyID"]), encRest.AWS.CustomerMasterKeyID)
		assert.Equal(t, string(secretData["RoleID"]), encRest.AWS.RoleID)
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

		encRest := &encryptionatrest.EncryptionAtRest{
			AWS: encryptionatrest.AwsKms{
				AwsKms: akov2.AwsKms{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name: "aws-secret",
					},
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test-fallback-ns")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["CustomerMasterKeyID"]), encRest.AWS.CustomerMasterKeyID)
		assert.Equal(t, string(secretData["RoleID"]), encRest.AWS.RoleID)
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

		encRest := &encryptionatrest.EncryptionAtRest{
			AWS: encryptionatrest.AwsKms{
				AwsKms: akov2.AwsKms{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name:      "aws-secret",
						Namespace: "test",
					},
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

		encRest := &encryptionatrest.EncryptionAtRest{
			GCP: encryptionatrest.GoogleCloudKms{
				GoogleCloudKms: akov2.GoogleCloudKms{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name: "gcp-secret",
					},
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["ServiceAccountKey"]), encRest.GCP.ServiceAccountKey)
		assert.Equal(t, string(secretData["KeyVersionResourceID"]), encRest.GCP.KeyVersionResourceID)
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

		encRest := &encryptionatrest.EncryptionAtRest{
			GCP: encryptionatrest.GoogleCloudKms{
				GoogleCloudKms: akov2.GoogleCloudKms{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name: "gcp-secret",
					},
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

		encRest := &encryptionatrest.EncryptionAtRest{
			Azure: encryptionatrest.AzureKeyVault{
				AzureKeyVault: akov2.AzureKeyVault{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name: "azure-secret",
					},
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.Nil(t, err)

		assert.Equal(t, string(secretData["Secret"]), encRest.Azure.Secret)
		assert.Equal(t, string(secretData["SubscriptionID"]), encRest.Azure.SubscriptionID)
		assert.Equal(t, string(secretData["KeyVaultName"]), encRest.Azure.KeyVaultName)
		assert.Equal(t, string(secretData["KeyIdentifier"]), encRest.Azure.KeyIdentifier)
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

		encRest := &encryptionatrest.EncryptionAtRest{
			Azure: encryptionatrest.AzureKeyVault{
				AzureKeyVault: akov2.AzureKeyVault{
					Enabled: pointer.MakePtr(true),
					SecretRef: common.ResourceRefNamespaced{
						Name: "gcp-secret",
					},
				},
			},
		}

		err := readEncryptionAtRestSecrets(kk, service, encRest, "test")
		assert.NotNil(t, err)
	})
}

func TestIsEncryptionAtlasEmpty(t *testing.T) {
	spec := &akov2.EncryptionAtRest{}
	isEmpty := IsEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Empty spec should be empty")

	spec.AwsKms.Enabled = pointer.MakePtr(true)
	isEmpty = IsEncryptionSpecEmpty(spec)
	assert.False(t, isEmpty, "Non-empty spec")

	spec.AwsKms.Enabled = pointer.MakePtr(false)
	isEmpty = IsEncryptionSpecEmpty(spec)
	assert.True(t, isEmpty, "Enabled flag set to false is same as empty")
}
