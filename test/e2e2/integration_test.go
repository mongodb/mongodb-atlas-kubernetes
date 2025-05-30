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
	"embed"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

//go:embed configs/*
var configs embed.FS

const (
	AtlasThirdPartyIntegrationsCRD = "atlasthirdpartyintegrations.atlas.nextapi.mongodb.com"
)

func TestAtlasThirdPartyIntegrationsCreate(t *testing.T) {
	control.SkipTestUnless(t, "AKO_E2E2_TEST")
	ns := control.MustEnvVar("OPERATOR_NAMESPACE")
	ctx := context.Background()
	for _, tc := range []struct {
		name      string
		objs      []client.Object
		wantReady string
	}{
		{
			name:      "simple datadog sample",
			objs:      yml.MustParseCRs(yml.MustOpen(configs, "configs/datadog.sample.yml")),
			wantReady: "atlas-datadog-integ",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			kubeClient, err := kube.NewK8sTest(ctx, AtlasThirdPartyIntegrationsCRD)
			require.NoError(t, err, "Kubernetes test env is not available")
			ako := runTestAKO()
			ako.Start(t)
			defer ako.Stop(t)

			require.NoError(t, kube.Apply(ctx, kubeClient, ns, tc.objs...))

			integration := akov2next.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{
					Name:      tc.wantReady,
					Namespace: ns,
				},
			}
			key := client.ObjectKeyFromObject(&integration)
			assert.NoError(t, kube.WaitConditionOrFailure(time.Minute, func() (bool, error) {
				return kube.AssertObjReady(ctx, kubeClient, key, &integration)
			}))
		})
	}
}

func runTestAKO() *operator.Operator {
	return operator.NewOperator(control.MustEnvVar("OPERATOR_NAMESPACE"), os.Stdout, os.Stderr,
		"--log-level=-9",
		"--global-api-secret-name=mongodb-atlas-operator-api-key",
		`--atlas-domain=https://cloud-qa.mongodb.com`,
	)
}
