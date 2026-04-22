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

package helm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/decoder"
)

const atlasAdvancedChartPath = "../../helm-charts/atlas-advanced"

// findAtlasAdvancedAPICredentialsSecret locates the Atlas API-credentials
// Secret rendered by the atlas-advanced chart. The chart also renders a
// separate Secret for the database user's password that carries the same
// `atlas.mongodb.com/type: credentials` label, so we filter by the
// "-secret"-suffixed name produced by the template.
func findAtlasAdvancedAPICredentialsSecret(t *testing.T, output string) *corev1.Secret {
	t.Helper()
	objects := decoder.DecodeAll(t, strings.NewReader(output))
	var found *corev1.Secret
	for _, obj := range objects {
		s, ok := obj.(*corev1.Secret)
		if !ok {
			continue
		}
		if s.Labels["atlas.mongodb.com/type"] != "credentials" {
			continue
		}
		if !strings.HasSuffix(s.Name, "-secret") {
			continue
		}
		require.Nilf(t, found, "more than one API-credentials Secret rendered: %v and %v", found, s)
		found = s
	}
	return found
}

func TestAtlasAdvanced_RendersAPIKeySecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_advanced_apikey_values.yaml",
		atlasAdvancedChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findAtlasAdvancedAPICredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "abcdefgh", string(secret.Data["publicApiKey"]))
	assert.Equal(t, "12345678-1234-1234-1234-1234567890ab", string(secret.Data["privateApiKey"]))
	assert.NotContains(t, secret.Data, "clientId")
	assert.NotContains(t, secret.Data, "clientSecret")
}

func TestAtlasAdvanced_RendersServiceAccountSecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_advanced_sa_values.yaml",
		atlasAdvancedChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findAtlasAdvancedAPICredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "mdb_sa_id_01234567890abcdef", string(secret.Data["clientId"]))
	assert.Equal(t, "mdb_sa_sk_01234567890abcdefghijklmnop", string(secret.Data["clientSecret"]))
	assert.NotContains(t, secret.Data, "publicApiKey")
	assert.NotContains(t, secret.Data, "privateApiKey")
}

func TestAtlasAdvanced_RejectsBothCredentialTypes(t *testing.T) {
	_, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_advanced_both_values.yaml",
		atlasAdvancedChartPath,
	)
	require.Error(t, err, "expected helm template to fail when both credential types are set")
	assert.Contains(t, stderr, "set either (publicKey,privateKey) or (clientId,clientSecret), not both",
		"stderr did not include the mutual-exclusion message; got: %s", stderr)
}
