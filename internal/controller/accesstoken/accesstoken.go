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

// Package accesstoken holds the schema of the Access Token Secret produced
// by the service-account-token controller and consumed by
// reconciler.GetConnectionConfig, plus the two helpers that operate on it:
// DeriveSecretName (namespace+name → deterministic Secret name) and
// CredentialsHash (clientId+clientSecret → staleness fingerprint).
//
// This package has no operator-specific imports by design — it is a pure
// data-model/helper package that both the producer and the consumer depend on.
package accesstoken

import (
	"fmt"
	"hash/fnv"

	"k8s.io/apimachinery/pkg/util/rand"
)

// Data-field keys on an Access Token Secret.
const (
	// AccessTokenKey holds the OAuth bearer token.
	AccessTokenKey = "accessToken"
	// ExpiryKey holds the RFC3339 timestamp at which the bearer token expires.
	ExpiryKey = "expiry"
	// CredentialsHashKey holds the FNV fingerprint of the (clientId, clientSecret)
	// pair used to mint the current bearer token. Used to detect credential
	// rotation before the cached bearer is used.
	CredentialsHashKey = "credentialsHash"

	// secretNamePrefix is the name prefix of every Access Token Secret.
	// Changing this would invalidate every existing Access Token Secret across
	// all deployments.
	secretNamePrefix = "atlas-access-token-"
)

// DeriveSecretName returns the deterministic Access Token Secret name for a
// given Connection Secret. The connection secret name is included literally
// in the result for operator debuggability; it is truncated when the total
// would exceed the Kubernetes 253-character DNS-subdomain limit. Uniqueness
// is guaranteed by the hash suffix, which always uses the full (untruncated)
// name as input.
func DeriveSecretName(namespace, connectionSecretName string) string {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(namespace + "/" + connectionSecretName))
	hash := rand.SafeEncodeString(fmt.Sprint(hasher.Sum64()))

	const k8sNameLimit = 253
	maxNameLen := k8sNameLimit - len(secretNamePrefix) - 1 - len(hash)
	name := connectionSecretName
	if len(name) > maxNameLen {
		name = name[:maxNameLen]
	}
	return secretNamePrefix + name + "-" + hash
}

// CredentialsHash returns a non-cryptographic fingerprint of the credential
// pair. The nul separator disambiguates ("ab","c") from ("a","bc") — without
// it, both would concatenate to the same input and collide.
func CredentialsHash(clientID, clientSecret string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(clientID + "\x00" + clientSecret))
	return fmt.Sprint(h.Sum64())
}
