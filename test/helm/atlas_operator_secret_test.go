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

const (
	atlasOperatorChartPath = "../../helm-charts/atlas-operator"

	// The API-key and Service-Account credential fixtures share identical
	// sample values across all four charts (atlas-operator, atlas-advanced,
	// atlas-basic, atlas-deployment). Pin them as package constants so the
	// per-chart helpers stay parameter-light and the assertions remain
	// consistent.
	testOrgID              = "6500000000000000000000aa"
	testPublicAPIKey       = "abcdefgh"
	testPrivateAPIKey      = "12345678-1234-1234-1234-1234567890ab"
	testSAClientID         = "mdb_sa_id_01234567890abcdef"
	testSAClientSecretData = "mdb_sa_sk_01234567890abcdefghijklmnop"
)

func findCredentialsSecret(t *testing.T, output string) *corev1.Secret {
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
		if _, ok := s.Data["orgId"]; !ok {
			// Skip DB-user password Secrets, which also carry the
			// credentials type label but have no orgId.
			continue
		}
		require.Nilf(t, found, "more than one Atlas API credentials Secret rendered: %v and %v", found, s)
		found = s
	}
	return found
}

// assertAPIKeySecret renders chartPath with the given values fixture and
// asserts the resulting credentials Secret contains the expected API-key
// fields (orgId / publicApiKey / privateApiKey) and no SA fields. All four
// charts share identical fixture values, pinned as package constants.
func assertAPIKeySecret(t *testing.T, chartPath, valuesFile string) {
	t.Helper()
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values="+valuesFile,
		chartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, testOrgID, string(secret.Data["orgId"]))
	assert.Equal(t, testPublicAPIKey, string(secret.Data["publicApiKey"]))
	assert.Equal(t, testPrivateAPIKey, string(secret.Data["privateApiKey"]))
	assert.NotContains(t, secret.Data, "clientId")
	assert.NotContains(t, secret.Data, "clientSecret")
}

// assertServiceAccountSecret renders chartPath with the given values fixture
// and asserts the resulting credentials Secret contains the expected SA
// fields (orgId / clientId / clientSecret) and no API-key fields. All four
// charts share identical fixture values, pinned as package constants.
func assertServiceAccountSecret(t *testing.T, chartPath, valuesFile string) {
	t.Helper()
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values="+valuesFile,
		chartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, testOrgID, string(secret.Data["orgId"]))
	assert.Equal(t, testSAClientID, string(secret.Data["clientId"]))
	assert.Equal(t, testSAClientSecretData, string(secret.Data["clientSecret"]))
	assert.NotContains(t, secret.Data, "publicApiKey")
	assert.NotContains(t, secret.Data, "privateApiKey")
}

// assertRejectsBothCredentialTypes runs helm template against a both-set
// fixture and asserts the rejection message appears in stderr. The expected
// substring varies per chart (atlas-operator/-deployment use 'publicApiKey',
// atlas-advanced/-basic use 'publicKey').
func assertRejectsBothCredentialTypes(t *testing.T, chartPath, valuesFile, expectedSubstring string) {
	t.Helper()
	_, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values="+valuesFile,
		chartPath,
	)
	require.Error(t, err, "expected helm template to fail when both credential types are set")
	assert.Contains(t, stderr, expectedSubstring,
		"stderr did not include the mutual-exclusion message; got: %s", stderr)
}

func TestAtlasOperator_RendersAPIKeySecret(t *testing.T) {
	assertAPIKeySecret(t, atlasOperatorChartPath, "atlas_operator_apikey_values.yaml")
}

func TestAtlasOperator_RendersServiceAccountSecret(t *testing.T) {
	assertServiceAccountSecret(t, atlasOperatorChartPath, "atlas_operator_sa_values.yaml")
}

func TestAtlasOperator_RejectsBothCredentialTypes(t *testing.T) {
	assertRejectsBothCredentialTypes(t, atlasOperatorChartPath, "atlas_operator_both_values.yaml",
		"set either (publicApiKey,privateApiKey) or (clientId,clientSecret), not both")
}
