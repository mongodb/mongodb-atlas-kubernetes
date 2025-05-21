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

package secret

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Secret map[string][]byte

type SecretProvider interface {
	Fetch(ctx context.Context, name string) (Secret, error)
}

type kubernetesSecrets struct {
	client client.Client
}

func NewKubernetesSecretProvider(c client.Client) SecretProvider {
	return &kubernetesSecrets{client: c}
}

func (ks *kubernetesSecrets) Fetch(ctx context.Context, secretName string) (Secret, error) {
	secret := v1.Secret{}
	err := ks.client.Get(ctx, client.ObjectKey{Name: secretName}, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret %q: %w", secretName, err)
	}
	return secret.Data, nil
}
