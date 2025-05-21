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

func Fetch(ctx context.Context, kubeClient client.Client, secretName string) (map[string][]byte, error) {
	secret := v1.Secret{}
	err := kubeClient.Get(ctx, client.ObjectKey{Name: secretName}, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret %q: %w", secretName, err)
	}
	return secret.Data, nil
}
