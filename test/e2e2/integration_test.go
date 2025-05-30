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

package e2e2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e2"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
)

const (
	AtlasThirdPartyIntegrationsCRD = "atlasthirdpartyintegrations.atlas.nextapi.mongodb.com"
)

func TestAtlasThirdPartyIntegrationsCreate(t *testing.T) {
	control.SkipTestUnless(t, "AKO_E2E2_TEST")
	ctx := context.Background()
	for _, tc := range []struct {
		name      string
		objs      []client.Object
		wantReady string
	}{
		{},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := e2e2.InitK8sTest(ctx, AtlasThirdPartyIntegrationsCRD)
			require.NoError(t, err, "Kubernetes test env is not available")
		})
	}
}
