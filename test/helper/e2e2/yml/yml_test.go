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

package yml_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

//go:embed samples/*
var samples embed.FS

func TestParseCRs(t *testing.T) {
	in, err := samples.Open("samples/sample.yml")
	require.NoError(t, err)
	defer in.Close()

	objs, err := yml.ParseCRs(in)
	require.NoError(t, err)
	assert.Len(t, objs, 2)
	assert.IsType(t, &corev1.Secret{}, objs[0])
	assert.Equal(t, &akov2next.AtlasThirdPartyIntegration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasThirdPartyIntegration",
			APIVersion: "atlas.nextapi.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-atlas-integ",
		},
		Spec: akov2next.AtlasThirdPartyIntegrationSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{
					ID: "68359c51ce672533a751117e",
				},
				ConnectionSecret: &api.LocalObjectReference{
					Name: "mongodb-atlas-operator-api-key",
				},
			},
			Type: "DATADOG",
			Datadog: &akov2next.DatadogIntegration{
				APIKeySecretRef: api.LocalObjectReference{
					Name: "datadog-secret",
				},
				Region:                       "US",
				SendCollectionLatencyMetrics: pointer.MakePtr("enabled"),
				SendDatabaseMetrics:          pointer.MakePtr("enabled"),
			},
		},
	}, objs[1])
}
