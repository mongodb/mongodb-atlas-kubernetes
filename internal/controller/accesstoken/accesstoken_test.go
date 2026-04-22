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
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/accesstoken"
)

func TestDeriveSecretName_PinsOutput(t *testing.T) {
	// Pin the exact output for a known input pair. Any change to the derivation
	// scheme (algorithm, separator, prefix, encoding) flips this value and
	// every existing Access Token Secret on a live cluster would be orphaned
	// behind a new name. Treat this assertion as a compatibility contract.
	got, err := accesstoken.DeriveSecretName("atlas-operator", "my-sa-creds")
	require.NoError(t, err)
	assert.Equal(t, "atlas-access-token-my-sa-creds-587997bcdf678bb69ff7", got)
}

func TestDeriveSecretName_NamespaceSensitive(t *testing.T) {
	a, err := accesstoken.DeriveSecretName("ns-a", "creds")
	require.NoError(t, err)
	b, err := accesstoken.DeriveSecretName("ns-b", "creds")
	require.NoError(t, err)
	assert.NotEqual(t, a, b, "same name in different namespaces must yield different outputs")
}

func TestDeriveSecretName_NameSensitive(t *testing.T) {
	a, err := accesstoken.DeriveSecretName("ns", "creds-a")
	require.NoError(t, err)
	b, err := accesstoken.DeriveSecretName("ns", "creds-b")
	require.NoError(t, err)
	assert.NotEqual(t, a, b, "different names in same namespace must yield different outputs")
}

func TestDeriveSecretName_LengthFarPastLimit(t *testing.T) {
	longName := strings.Repeat("x", 500)
	result, err := accesstoken.DeriveSecretName("ns", longName)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(result), 253, "result must fit in DNS-1123 subdomain limit")
	assert.True(t, strings.HasPrefix(result, "atlas-access-token-"))
}

func TestDeriveSecretName_LengthAtBoundary(t *testing.T) {
	const ns = "ns"

	// A very long input forces truncation; the result must be exactly 253.
	veryLong, err := accesstoken.DeriveSecretName(ns, strings.Repeat("a", 500))
	require.NoError(t, err)
	assert.Equal(t, 253, len(veryLong),
		"when truncation is forced, result length must be exactly 253 — guards off-by-one in maxNameLen")
	assert.True(t, strings.HasPrefix(veryLong, "atlas-access-token-"),
		"prefix must be preserved even under maximal truncation")

	// Short names are preserved literally.
	short, err := accesstoken.DeriveSecretName(ns, "short-name")
	require.NoError(t, err)
	assert.LessOrEqual(t, len(short), 253)
	assert.Contains(t, short, "short-name")
}

func TestCredentialsHash_PinsOutput(t *testing.T) {
	// Pin the exact FNV-1a-64 output for a known input pair. Any accidental
	// change to the hashing scheme (algorithm, separator, encoding) flips this
	// value — which means every previously issued Access Token Secret on a
	// live cluster would look "stale" and be needlessly refreshed. Treat this
	// assertion as a compatibility contract, not a self-test.
	got, err := accesstoken.CredentialsHash("client-id", "client-secret")
	require.NoError(t, err)
	assert.Equal(t, "3974328787184052522", got)
}

func TestCredentialsHash_DistinguishesInputs(t *testing.T) {
	a, err := accesstoken.CredentialsHash("id-1", "secret-1")
	require.NoError(t, err)
	b, err := accesstoken.CredentialsHash("id-2", "secret-1")
	require.NoError(t, err)
	c, err := accesstoken.CredentialsHash("id-1", "secret-2")
	require.NoError(t, err)
	assert.Equal(t, "6130960229688205592", a)
	assert.Equal(t, "1640858821594590263", b, "different clientIDs must produce different hashes")
	assert.Equal(t, "6130963528223090225", c, "different clientSecrets must produce different hashes")
}

func TestCredentialsHash_NulSeparatorDisambiguation(t *testing.T) {
	// Without a nul separator "ab"+"c" would equal "a"+"bc" after concat,
	// producing a hash collision. The nul separator must prevent this.
	a, err := accesstoken.CredentialsHash("ab", "c")
	require.NoError(t, err)
	b, err := accesstoken.CredentialsHash("a", "bc")
	require.NoError(t, err)
	assert.Equal(t, "18258086037221804135", a)
	assert.Equal(t, "12340134017423684899", b, "nul separator must disambiguate credential pairs with shared concatenation")
}
