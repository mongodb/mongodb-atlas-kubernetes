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

package reconciler

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

const (
	orgIDKey      = "orgId"
	publicAPIKey  = "publicApiKey"
	privateAPIKey = "privateApiKey"

	clientIDKey     = "clientId"
	clientSecretKey = "clientSecret"

	accessTokenKey          = "accessToken"
	credentialsHashKey      = "credentialsHash"
	accessTokenSecretPrefix = "atlas-access-token-"
)

// CredentialsHash returns a deterministic, non-cryptographic fingerprint of
// the (clientID, clientSecret) pair. The service-account-token controller
// stores this fingerprint on the Access Token Secret under credentialsHashKey
// so any component reading the Secret can detect that the source credentials
// have been rotated since the cached bearer token was issued. The nul
// separator disambiguates ("ab","c") from ("a","bc").
func CredentialsHash(clientID, clientSecret string) (string, error) {
	h := fnv.New64a()
	if _, err := h.Write([]byte(clientID + "\x00" + clientSecret)); err != nil {
		return "", fmt.Errorf("failed to compute credentials hash: %w", err)
	}
	return fmt.Sprint(h.Sum64()), nil
}

func (r *AtlasReconciler) ResolveConnectionConfig(ctx context.Context, referrer project.ProjectReferrerObject) (*atlas.ConnectionConfig, error) {
	connectionSecret := r.connectionSecretRef(referrer)
	if connectionSecret != nil && connectionSecret.Name != "" {
		cfg, err := GetConnectionConfig(ctx, r.Client, connectionSecret, &r.GlobalSecretRef)
		if err != nil {
			return nil, fmt.Errorf("error getting credentials from connection secret: %w", err)
		}
		return cfg, nil
	}

	prj, err := r.fetchProject(ctx, referrer)
	if err != nil {
		return nil, fmt.Errorf("error resolving project reference: %w", err)
	}

	var projectSecret *client.ObjectKey
	if prj != nil {
		projectSecret = prj.ConnectionSecretObjectKey()
	}

	cfg, err := GetConnectionConfig(ctx, r.Client, projectSecret, &r.GlobalSecretRef)
	if err != nil {
		return nil, fmt.Errorf("error getting credentials from project reference: %w", err)
	}
	return cfg, nil
}

func (r *AtlasReconciler) connectionSecretRef(pro project.ProjectReferrerObject) *client.ObjectKey {
	key := client.ObjectKeyFromObject(pro)
	pdr := pro.ProjectDualRef()
	if pdr.ConnectionSecret == nil {
		return nil
	}
	key.Name = pdr.ConnectionSecret.Name
	return &key
}

func GetConnectionConfig(ctx context.Context, k8sClient client.Client, secretRef, fallbackRef *client.ObjectKey) (*atlas.ConnectionConfig, error) {
	if secretRef == nil {
		secretRef = fallbackRef
	}

	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, *secretRef, secret); err != nil {
		return nil, fmt.Errorf("failed to read Atlas API credentials from the secret %s: %w", secretRef.String(), err)
	}

	if err := validateConnectionSecret(secret); err != nil {
		return nil, fmt.Errorf("invalid connection secret %s: %w", secretRef, err)
	}

	if isServiceAccountCredentials(secret) {
		bearerToken, err := getServiceAccountAccessToken(ctx, k8sClient, secret)
		if err != nil {
			return nil, err
		}

		return &atlas.ConnectionConfig{
			OrgID: string(secret.Data[orgIDKey]),
			Credentials: &atlas.Credentials{
				ServiceAccount: &atlas.ServiceAccountToken{
					BearerToken: bearerToken,
				},
			},
		}, nil
	}

	return &atlas.ConnectionConfig{
		OrgID: string(secret.Data[orgIDKey]),
		Credentials: &atlas.Credentials{
			APIKeys: &atlas.APIKeys{
				PublicKey:  string(secret.Data[publicAPIKey]),
				PrivateKey: string(secret.Data[privateAPIKey]),
			},
		},
	}, nil
}

// DeriveAccessTokenSecretName returns the deterministic name of the Access Token Secret for a given Connection Secret.
// The Connection Secret name is included literally for operator debuggability; it is truncated when the total
// exceeds the Kubernetes 253-character DNS-subdomain limit.
func DeriveAccessTokenSecretName(namespace, connectionSecretName string) (string, error) {
	hasher := fnv.New64a()
	_, err := hasher.Write([]byte(namespace + "/" + connectionSecretName))
	if err != nil {
		return "", fmt.Errorf("failed to compute hash for access token secret name: %w", err)
	}
	hash := rand.SafeEncodeString(fmt.Sprint(hasher.Sum64()))

	const k8sNameLimit = 253
	maxNameLen := k8sNameLimit - len(accessTokenSecretPrefix) - 1 - len(hash)
	name := connectionSecretName
	if len(name) > maxNameLen {
		name = name[:maxNameLen]
	}

	return accessTokenSecretPrefix + name + "-" + hash, nil
}

func getServiceAccountAccessToken(ctx context.Context, k8sClient client.Client, secret *corev1.Secret) (string, error) {
	tokenSecretName, err := DeriveAccessTokenSecretName(secret.Namespace, secret.Name)
	if err != nil {
		return "", err
	}
	tokenRef := client.ObjectKey{Namespace: secret.Namespace, Name: tokenSecretName}

	tokenSecret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, tokenRef, tokenSecret); err != nil {
		if apierrors.IsNotFound(err) {
			return "", fmt.Errorf("access token secret %s does not exist yet", tokenRef.String())
		}
		return "", fmt.Errorf("failed to read access token secret %s: %w", tokenRef.String(), err)
	}

	// Guard against a stale cached token — if the credential Secret was
	// rotated since the token was issued, the service-account-token controller
	// may not have caught up yet. Returning an error prompts the downstream
	// reconciler to retry rather than hitting Atlas with revoked credentials.
	currentHash, err := CredentialsHash(string(secret.Data[clientIDKey]), string(secret.Data[clientSecretKey]))
	if err != nil {
		return "", err
	}
	if string(tokenSecret.Data[credentialsHashKey]) != currentHash {
		return "", fmt.Errorf("access token secret %s is stale (credentials rotated); waiting for the service-account-token controller to refresh", tokenRef.String())
	}

	bearerToken := string(tokenSecret.Data[accessTokenKey])
	if bearerToken == "" {
		return "", fmt.Errorf("access token secret %s has an empty accessToken field", tokenRef.String())
	}

	return bearerToken, nil
}

func isServiceAccountCredentials(credentials *corev1.Secret) bool {
	clientID := credentials.Data[clientIDKey]
	clientSecret := credentials.Data[clientSecretKey]

	return len(clientID) > 0 && len(clientSecret) > 0
}

func validateConnectionSecret(secret *corev1.Secret) error {
	hasAnyAPIKey := len(secret.Data[publicAPIKey]) > 0 || len(secret.Data[privateAPIKey]) > 0
	hasAnySA := len(secret.Data[clientIDKey]) > 0 || len(secret.Data[clientSecretKey]) > 0

	if hasAnyAPIKey && hasAnySA {
		return errors.New("secret contains both API key and service account credentials; only one type is allowed")
	}

	var missingFields []string

	if len(secret.Data[orgIDKey]) == 0 {
		missingFields = append(missingFields, orgIDKey)
	}

	if hasAnyAPIKey {
		if len(secret.Data[publicAPIKey]) == 0 {
			missingFields = append(missingFields, publicAPIKey)
		}
		if len(secret.Data[privateAPIKey]) == 0 {
			missingFields = append(missingFields, privateAPIKey)
		}
	} else if hasAnySA {
		if len(secret.Data[clientIDKey]) == 0 {
			missingFields = append(missingFields, clientIDKey)
		}
		if len(secret.Data[clientSecretKey]) == 0 {
			missingFields = append(missingFields, clientSecretKey)
		}
	} else {
		//By default, we are expecting API keys
		missingFields = append(missingFields, publicAPIKey, privateAPIKey)
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %v", missingFields)
	}

	return nil
}
