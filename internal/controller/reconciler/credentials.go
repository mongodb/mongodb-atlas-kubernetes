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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

const (
	orgIDKey      = "orgId"
	publicAPIKey  = "publicApiKey"
	privateAPIKey = "privateApiKey"
)

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

	cfg := &atlas.ConnectionConfig{
		OrgID: string(secret.Data[orgIDKey]),
		Credentials: &atlas.Credentials{
			APIKeys: &atlas.APIKeys{
				PublicKey:  string(secret.Data[publicAPIKey]),
				PrivateKey: string(secret.Data[privateAPIKey]),
			},
		},
	}

	if missingFields, valid := validate(cfg); !valid {
		return nil, fmt.Errorf("the following fields are missing in the secret %v: %v", secretRef, missingFields)
	}

	return cfg, nil
}

func validate(cfg *atlas.ConnectionConfig) ([]string, bool) {
	missingFields := make([]string, 0, 3)

	if cfg == nil {
		return []string{orgIDKey, publicAPIKey, privateAPIKey}, false
	}

	if cfg.OrgID == "" {
		missingFields = append(missingFields, orgIDKey)
	}

	if cfg.Credentials == nil || cfg.Credentials.APIKeys == nil {
		return append(missingFields, []string{publicAPIKey, privateAPIKey}...), false
	}

	if cfg.Credentials.APIKeys.PublicKey == "" {
		missingFields = append(missingFields, publicAPIKey)
	}

	if cfg.Credentials.APIKeys.PrivateKey == "" {
		missingFields = append(missingFields, privateAPIKey)
	}

	if len(missingFields) > 0 {
		return missingFields, false
	}

	return nil, true
}
