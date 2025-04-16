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

package indexer

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestAtlasProjectByConnectionSecretIndexer(t *testing.T) {
	for _, tc := range []struct {
		name     string
		object   client.Object
		wantKeys []string
	}{
		{
			name:     "should return nil on wrong type",
			object:   &akov2.AtlasDatabaseUser{},
			wantKeys: nil,
		},
		{
			name:     "should return nil when there are no references",
			object:   &akov2.AtlasProject{},
			wantKeys: []string{},
		},
		{
			name: "should return nil when there is an empty reference",
			object: &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					ConnectionSecret: &common.ResourceRefNamespaced{},
				},
			},
			wantKeys: []string{},
		},
		{
			name: "should return project namespace if name is set only",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "projectName", Namespace: "projectNamespace"},
				Spec: akov2.AtlasProjectSpec{
					ConnectionSecret: &common.ResourceRefNamespaced{Name: "someSecret"},
				},
			},
			wantKeys: []string{"projectNamespace/someSecret"},
		},
		{
			name: "should return secret namespace and name if set",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "projectName", Namespace: "projectNamespace"},
				Spec: akov2.AtlasProjectSpec{
					ConnectionSecret: &common.ResourceRefNamespaced{Name: "someSecret", Namespace: "secretNamespace"},
				},
			},
			wantKeys: []string{"secretNamespace/someSecret"},
		},
		{
			name: "should also return secrets in encryption at rest",
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "projectName", Namespace: "projectNamespace"},
				Spec: akov2.AtlasProjectSpec{
					ConnectionSecret: &common.ResourceRefNamespaced{Name: "ConnectionSecret", Namespace: "secretNamespace"},
					AlertConfigurations: []akov2.AlertConfiguration{
						{
							Notifications: []akov2.Notification{
								{APITokenRef: common.ResourceRefNamespaced{Name: "APITokenRef"}},
								{APITokenRef: common.ResourceRefNamespaced{Name: "APITokenRef"}}, // double entry
								{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "DatadogAPIKeyRef"}},
								{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "FlowdockAPITokenRef"}},
								{OpsGenieAPIKeyRef: common.ResourceRefNamespaced{Name: "OpsGenieAPIKeyRef"}},
								{ServiceKeyRef: common.ResourceRefNamespaced{Name: "ServiceKeyRef"}},
								{VictorOpsSecretRef: common.ResourceRefNamespaced{Name: "VictorOpsSecretRef"}},
							},
						},
						{ // duplicate some entries from the previous alertconfig
							Notifications: []akov2.Notification{
								{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "DatadogAPIKeyRef"}},
								{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "DatadogAPIKeyRef"}},
								{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "FlowdockAPITokenRef"}},
								{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "FlowdockAPITokenRef"}},
								{VictorOpsSecretRef: common.ResourceRefNamespaced{Name: "VictorOpsSecretRef"}},
								{VictorOpsSecretRef: common.ResourceRefNamespaced{Name: "VictorOpsSecretRef"}},
							},
						},
					},
					EncryptionAtRest: &akov2.EncryptionAtRest{
						AwsKms:         akov2.AwsKms{SecretRef: common.ResourceRefNamespaced{Name: "AwsKms"}},
						AzureKeyVault:  akov2.AzureKeyVault{SecretRef: common.ResourceRefNamespaced{Name: "AzureKeyVault"}},
						GoogleCloudKms: akov2.GoogleCloudKms{SecretRef: common.ResourceRefNamespaced{Name: "GoogleCloudKms"}},
					},
				},
			},
			wantKeys: []string{
				"projectNamespace/APITokenRef",
				"projectNamespace/AwsKms",
				"projectNamespace/AzureKeyVault",
				"projectNamespace/DatadogAPIKeyRef",
				"projectNamespace/FlowdockAPITokenRef",
				"projectNamespace/GoogleCloudKms",
				"projectNamespace/OpsGenieAPIKeyRef",
				"projectNamespace/ServiceKeyRef",
				"projectNamespace/VictorOpsSecretRef",
				"secretNamespace/ConnectionSecret",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			indexer := NewAtlasProjectByConnectionSecretIndexer(zaptest.NewLogger(t))
			keys := indexer.Keys(tc.object)
			sort.Strings(keys)
			assert.Equal(t, tc.wantKeys, keys)
		})
	}
}
