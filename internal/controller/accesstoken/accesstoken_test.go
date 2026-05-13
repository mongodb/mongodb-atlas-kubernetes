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

package accesstoken_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/accesstoken"
)

func TestDeriveSecretName(t *testing.T) {
	dataProvider := map[string]struct {
		namespace            string
		connectionSecretName string
		expected             string
	}{
		"pinned output for known input": {
			namespace:            "atlas-operator",
			connectionSecretName: "my-sa-creds",
			expected:             "atlas-access-token-my-sa-creds-587997bcdf678bb69ff7",
		},
		"namespace sensitivity: ns-a + creds": {
			namespace:            "ns-a",
			connectionSecretName: "creds",
			expected:             "atlas-access-token-creds-58694dcd4586644b4778",
		},
		"namespace sensitivity: ns-b + creds": {
			namespace:            "ns-b",
			connectionSecretName: "creds",
			expected:             "atlas-access-token-creds-5477fbbc7d7fb85964bc",
		},
		"connection-name sensitivity: ns + creds-a": {
			namespace:            "ns",
			connectionSecretName: "creds-a",
			expected:             "atlas-access-token-creds-a-7dbc78d46547845579d",
		},
		"connection-name sensitivity: ns + creds-b": {
			namespace:            "ns",
			connectionSecretName: "creds-b",
			expected:             "atlas-access-token-creds-b-7dbc78bf659667d758c",
		},
		"long connection name is truncated to the DNS-1123 253-char limit": {
			namespace:            "ns",
			connectionSecretName: strings.Repeat("a", 500),
			expected:             "atlas-access-token-" + strings.Repeat("a", 214) + "-cf5697c4fb6ddb7dcff",
		},
		"short connection name is preserved literally": {
			namespace:            "ns",
			connectionSecretName: "short-name",
			expected:             "atlas-access-token-short-name-ccf99675cf99c8ffb85",
		},
	}

	for desc, data := range dataProvider {
		t.Run(desc, func(t *testing.T) {
			got := accesstoken.DeriveSecretName(data.namespace, data.connectionSecretName)
			assert.Equal(t, data.expected, got)
		})
	}
}

func TestCredentialsHash(t *testing.T) {
	dataProvider := map[string]struct {
		clientID     string
		clientSecret string
		expected     string
	}{
		"pinned output for known input": {
			clientID:     "client-id",
			clientSecret: "client-secret",
			expected:     "3974328787184052522",
		},
		"distinctness: id-1 + secret-1": {
			clientID:     "id-1",
			clientSecret: "secret-1",
			expected:     "6130960229688205592",
		},
		"distinctness: different clientID (id-2 + secret-1)": {
			clientID:     "id-2",
			clientSecret: "secret-1",
			expected:     "1640858821594590263",
		},
		"distinctness: different clientSecret (id-1 + secret-2)": {
			clientID:     "id-1",
			clientSecret: "secret-2",
			expected:     "6130963528223090225",
		},
		"nul separator: ab + c": {
			clientID:     "ab",
			clientSecret: "c",
			expected:     "18258086037221804135",
		},
		"nul separator: a + bc": {
			clientID:     "a",
			clientSecret: "bc",
			expected:     "12340134017423684899",
		},
	}

	for desc, data := range dataProvider {
		t.Run(desc, func(t *testing.T) {
			got := accesstoken.CredentialsHash(data.clientID, data.clientSecret)
			assert.Equal(t, data.expected, got)
		})
	}
}
