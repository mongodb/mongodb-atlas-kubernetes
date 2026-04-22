# Service Account Credentials in Installation Methods — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.
>
> **Branch:** `CLOUDP-373260/service-accounts-phase-1-helm-charts` (branched off `CLOUDP-373260/service-accounts-phase-1`).

**Goal:** Bring every install method (four Helm charts + OLM bundle + generated RBAC) up to parity with the Service-Account runtime feature, and pin the rendering contract with helm-template unit tests.

**Architecture:** Each chart's credentials-Secret template gains a second path that emits `orgId + clientId + clientSecret` when SA values are populated, while preserving the legacy `orgId + publicApiKey + privateApiKey` path. The `ServiceAccountToken` controller already carries the kubebuilder markers that expand the Secret RBAC; running `make manifests && make bundle && make helm-crds` propagates the refreshed rules through `config/rbac/` → the OLM CSV → `helm-charts/atlas-operator/rbac.yaml`.

**Tech Stack:** Go, Helm templating, controller-gen, operator-sdk (for bundle), devbox-provisioned tools (`helm`, `yq`, `kustomize`, `kubectl`).

**Spec:** [`docs/superpowers/specs/2026-04-22-sa-installation-methods-design.md`](../specs/2026-04-22-sa-installation-methods-design.md)

---

## File Map

**Charts (hand-edited per chunk 2–5):**

- `helm-charts/atlas-operator/values.yaml` — add `clientId`/`clientSecret` to `globalConnectionSecret`.
- `helm-charts/atlas-operator/templates/global-secret.yaml` — add mutual-exclusion guard + SA data-branch.
- `helm-charts/atlas-operator/README.md` — SA install example.
- `helm-charts/atlas-advanced/values.yaml` — add `clientId`/`clientSecret` to `secret`.
- `helm-charts/atlas-advanced/templates/atlas-secret.yaml` — guard + SA branch.
- `helm-charts/atlas-advanced/README.md` — SA install example.
- `helm-charts/atlas-basic/values.yaml` — add `clientId`/`clientSecret` to `secret`.
- `helm-charts/atlas-basic/templates/atlas-secret.yaml` — guard + SA branch.
- `helm-charts/atlas-basic/README.md` — SA install example.
- `helm-charts/atlas-deployment/values.yaml` — add `clientId`/`clientSecret` to `atlas.secret`.
- `helm-charts/atlas-deployment/templates/atlas-secret.yaml` — guard + SA branch inside the existing outer conditionals.
- `helm-charts/atlas-deployment/README.md` — SA install example.

**Regenerated artefacts (chunk 1):**

- `config/rbac/role.yaml` + any split `clusterwide/` and `namespaced/` variants.
- `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml` + related bundle files.
- `helm-charts/atlas-operator/rbac.yaml`.

**Tests (new, chunks 1–5):**

- `test/helm/atlas_operator_rbac_test.go` — pins Secret verbs in rendered ClusterRole.
- `test/helm/atlas_operator_secret_test.go` + `*_values.yaml` fixtures.
- `test/helm/atlas_advanced_secret_test.go` + fixtures.
- `test/helm/atlas_basic_secret_test.go` + fixtures.
- `test/helm/atlas_deployment_secret_test.go` + fixtures.
- `test/helm/helm_cli.go` — new test-local helper that runs `helm template` via `exec.Command` and returns `(stdout, stderr, err)`, for the `{{ fail }}` scenarios.

**TD update (chunk 6):**

- `docs/dev/td-service-accounts-phase1.md` — rewrite the Installation section to reflect what actually ships.

---

## Chunk 1: Regenerate RBAC, bundle, and helm-charts/atlas-operator/rbac.yaml

Runs the existing `make` chain to bring the on-disk RBAC/CSV artefacts in sync with the `ServiceAccountToken` controller's kubebuilder markers, then pins the result with a unit test.

### Task 1: Run the regeneration chain

**Files (regenerated, not hand-edited):**

- `config/rbac/role.yaml` (+ split variants under `config/rbac/clusterwide/`, `config/rbac/namespaced/` as emitted by `scripts/split_roles_yaml.sh`)
- `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml`
- `helm-charts/atlas-operator/rbac.yaml`

- [ ] **Step 1.1: Regenerate RBAC from the kubebuilder markers**

```bash
cd /Users/maciej.karas/mongodb/mongodb-atlas-kubernetes
devbox run -- make manifests
```

Expected: clean exit. Compare `git diff config/rbac/` — the Secret rule should now list verbs `get`, `list`, `watch`, `create`, `update`, `patch` (up from `get`, `list`, `watch` before regeneration).

- [ ] **Step 1.2: Regenerate the OLM bundle**

```bash
devbox run -- make bundle
```

Expected: clean exit. `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml` updates its `spec.install.spec.clusterPermissions[0].rules` Secret entry with the expanded verbs.

- [ ] **Step 1.3: Regenerate the chart's rbac.yaml from the fresh CSV**

```bash
devbox run -- make helm-crds
```

Expected: clean exit. `helm-charts/atlas-operator/rbac.yaml` now contains the expanded Secret rule. `helm-charts/atlas-operator-crds/templates/atlas.*` CRDs also refresh to match the bundle.

- [ ] **Step 1.4: Run validation gates**

```bash
devbox run -- make bundle-validate
devbox run -- make validate-manifests
devbox run -- make validate-crds-chart
```

Expected: all three succeed. `validate-manifests` would fail if regen produced uncommitted drift from what `make manifests` expects; this is the CI gate that catches forgotten regens.

### Task 2: Pin the regenerated ClusterRole with a unit test

**Files:**

- Create: `test/helm/atlas_operator_rbac_test.go`

- [ ] **Step 2.1: Write the failing test**

Write this content to `test/helm/atlas_operator_rbac_test.go`:

```go
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cmd"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/decoder"
)

// TestOperatorClusterRole_SecretVerbs pins the verbs on the Secret rule in
// the rendered ClusterRole. If anyone changes a kubebuilder marker that
// affects Secrets but forgets to run `make manifests && make bundle &&
// make helm-crds`, this test fails and points them at the regen chain.
func TestOperatorClusterRole_SecretVerbs(t *testing.T) {
	output := cmd.RunCommand(t,
		"helm", "template",
		"--namespace=default",
		"../../helm-charts/atlas-operator")
	objects := decoder.DecodeAll(t, output)

	var clusterRoles []*rbacv1.ClusterRole
	for _, obj := range objects {
		if cr, ok := obj.(*rbacv1.ClusterRole); ok {
			clusterRoles = append(clusterRoles, cr)
		}
	}
	require.NotEmpty(t, clusterRoles, "expected at least one ClusterRole in rendered chart")

	var secretVerbs []string
	for _, cr := range clusterRoles {
		for _, rule := range cr.Rules {
			isCoreGroup := false
			for _, g := range rule.APIGroups {
				if g == "" {
					isCoreGroup = true
					break
				}
			}
			if !isCoreGroup {
				continue
			}
			for _, res := range rule.Resources {
				if res == "secrets" {
					secretVerbs = append(secretVerbs, rule.Verbs...)
				}
			}
		}
	}

	require.NotEmpty(t, secretVerbs,
		"no Secret rule in any rendered ClusterRole; regen chain likely not run")
	for _, wanted := range []string{"get", "list", "watch", "create", "update", "patch"} {
		assert.Containsf(t, secretVerbs, wanted,
			"Secret verbs missing %q — run `make manifests && make bundle && make helm-crds`",
			wanted)
	}
}
```

- [ ] **Step 2.2: Run the test to confirm it passes**

```bash
cd /Users/maciej.karas/mongodb/mongodb-atlas-kubernetes
go test ./test/helm/ -run TestOperatorClusterRole_SecretVerbs -v
```

Expected: PASS. If it fails on missing verbs, Step 1.1–1.3 didn't actually regenerate (re-run). If it fails with "no ClusterRole in rendered chart", the chart default values don't render a ClusterRole (unlikely — check `crossNamespaceRoles: true` default in `values.yaml`).

- [ ] **Step 2.3: Run the full helm unit-test package**

```bash
go test ./test/helm/ -v
```

Expected: `TestFlexSpec` passes (existing); `TestOperatorClusterRole_SecretVerbs` passes. No other tests in the package yet.

- [ ] **Step 2.4: Commit**

```bash
git add config/rbac/ bundle/ helm-charts/atlas-operator/rbac.yaml \
        helm-charts/atlas-operator-crds/templates/ \
        test/helm/atlas_operator_rbac_test.go
git commit -m "Regenerate RBAC, bundle, and chart rbac.yaml; pin Secret verbs in unit test"
```

(The `helm-charts/atlas-operator-crds/templates/` CRD files may or may not show diff — include them if they do; they're also regenerated by `make helm-crds`.)

---

## Chunk 2: Test helper + `atlas-operator` chart SA path

### Task 3: Add the test-local `helm template` helper

**Files:**

- Create: `test/helm/helm_cli.go`

This is used by every secret test in chunks 2–5. Add it once, share across.

- [ ] **Step 3.1: Write the helper**

Write this content to `test/helm/helm_cli.go`:

```go
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
	"bytes"
	"os/exec"
	"testing"
)

// helmTemplate runs `helm template <args>` and returns (stdout, stderr, err).
// Unlike test/helper/cmd.RunCommand, this does not fail the test on non-zero
// exit — callers that expect `{{ fail "..." }}` need to inspect stderr.
func helmTemplate(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("helm", append([]string{"template"}, args...)...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
```

- [ ] **Step 3.2: Confirm it compiles**

```bash
go build ./test/helm/...
```

Expected: clean exit.

### Task 4: `atlas-operator` chart — API-key rendering test and fixture

**Files:**

- Create: `test/helm/atlas_operator_apikey_values.yaml`
- Create: `test/helm/atlas_operator_secret_test.go` (one file, grows across tasks 4–6)

- [ ] **Step 4.1: Write the API-key values fixture**

Write `test/helm/atlas_operator_apikey_values.yaml`:

```yaml
globalConnectionSecret:
  orgId: "6500000000000000000000aa"
  publicApiKey: "abcdefgh"
  privateApiKey: "12345678-1234-1234-1234-1234567890ab"
```

- [ ] **Step 4.2: Write the API-key rendering test**

Write `test/helm/atlas_operator_secret_test.go` (initial content — tasks 5 and 6 append to the same file):

```go
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

const atlasOperatorChartPath = "../../helm-charts/atlas-operator"

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
		require.Nilf(t, found, "more than one credentials Secret rendered: %v and %v", found, s)
		found = s
	}
	return found
}

func TestAtlasOperator_RendersAPIKeySecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_operator_apikey_values.yaml",
		atlasOperatorChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "credentials", secret.Labels["atlas.mongodb.com/type"])
	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "abcdefgh", string(secret.Data["publicApiKey"]))
	assert.Equal(t, "12345678-1234-1234-1234-1234567890ab", string(secret.Data["privateApiKey"]))
	assert.NotContains(t, secret.Data, "clientId",
		"API-key path must not render clientId")
	assert.NotContains(t, secret.Data, "clientSecret",
		"API-key path must not render clientSecret")
}
```

- [ ] **Step 4.3: Run the test — expect PASS against the CURRENT template**

```bash
go test ./test/helm/ -run TestAtlasOperator_RendersAPIKeySecret -v
```

Expected: PASS. The current (unmodified) template renders an API-key Secret matching the assertions. This is the pre-change baseline.

### Task 5: `atlas-operator` chart — Service-Account rendering test (failing) and template change

**Files:**

- Create: `test/helm/atlas_operator_sa_values.yaml`
- Modify: `test/helm/atlas_operator_secret_test.go` (append)
- Modify: `helm-charts/atlas-operator/values.yaml`
- Modify: `helm-charts/atlas-operator/templates/global-secret.yaml`

- [ ] **Step 5.1: Write the SA values fixture**

Write `test/helm/atlas_operator_sa_values.yaml`:

```yaml
globalConnectionSecret:
  orgId: "6500000000000000000000aa"
  clientId: "mdb_sa_id_01234567890abcdef"
  clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"
```

- [ ] **Step 5.2: Append the SA rendering test**

Append this to `test/helm/atlas_operator_secret_test.go`:

```go
func TestAtlasOperator_RendersServiceAccountSecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_operator_sa_values.yaml",
		atlasOperatorChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "credentials", secret.Labels["atlas.mongodb.com/type"])
	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "mdb_sa_id_01234567890abcdef", string(secret.Data["clientId"]))
	assert.Equal(t, "mdb_sa_sk_01234567890abcdefghijklmnop", string(secret.Data["clientSecret"]))
	assert.NotContains(t, secret.Data, "publicApiKey",
		"Service-Account path must not render publicApiKey")
	assert.NotContains(t, secret.Data, "privateApiKey",
		"Service-Account path must not render privateApiKey")
}
```

- [ ] **Step 5.3: Run the test — expect FAIL**

```bash
go test ./test/helm/ -run TestAtlasOperator_RendersServiceAccountSecret -v
```

Expected: FAIL — `findCredentialsSecret` returns nil because the unmodified template only renders when `publicApiKey` is set, and the SA fixture doesn't set it.

- [ ] **Step 5.4: Update `values.yaml`**

Modify `helm-charts/atlas-operator/values.yaml`. Find the `globalConnectionSecret:` block (around lines 33–36) and replace it with:

```yaml
# globalConnectionSecret is a default "global" Secret containing Atlas
# authentication information.
#
# Configure EITHER the API-key fields (publicApiKey + privateApiKey) OR the
# Service Account fields (clientId + clientSecret). Setting both at the same
# time is rejected at chart-render time.
#
# You should never check-in these values as part of values.yaml file on your
# CVS. Instead set these values with `--set`.
globalConnectionSecret:
  orgId: ""
  # API key (legacy)
  publicApiKey: ""
  privateApiKey: ""
  # Service Account (recommended)
  clientId: ""
  clientSecret: ""
```

(Preserve the multi-line comment block immediately above the block that was already there — the snippet above shows the replacement comment. Do not introduce a blank line between the comment and the key.)

- [ ] **Step 5.5: Update `templates/global-secret.yaml`**

Replace the full content of `helm-charts/atlas-operator/templates/global-secret.yaml` with:

```
{{- if and .Values.globalConnectionSecret.publicApiKey .Values.globalConnectionSecret.clientSecret }}
{{- fail "globalConnectionSecret: set either (publicApiKey,privateApiKey) or (clientId,clientSecret), not both" }}
{{- end }}
{{- if or .Values.globalConnectionSecret.publicApiKey .Values.globalConnectionSecret.clientSecret }}
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: "{{ include "mongodb-atlas-operator.name" . }}-api-key"
  labels:
    atlas.mongodb.com/type: "credentials"
    {{- include "mongodb-atlas-operator.labels" . | nindent 4 }}
data:
    orgId: {{ .Values.globalConnectionSecret.orgId | b64enc }}
    {{- if .Values.globalConnectionSecret.publicApiKey }}
    publicApiKey: {{ .Values.globalConnectionSecret.publicApiKey | b64enc }}
    privateApiKey: {{ .Values.globalConnectionSecret.privateApiKey | b64enc }}
    {{- else }}
    clientId: {{ .Values.globalConnectionSecret.clientId | b64enc }}
    clientSecret: {{ .Values.globalConnectionSecret.clientSecret | b64enc }}
    {{- end }}
{{- end }}
```

- [ ] **Step 5.6: Run both rendering tests — expect PASS**

```bash
go test ./test/helm/ -run "TestAtlasOperator_Renders" -v
```

Expected: both `TestAtlasOperator_RendersAPIKeySecret` and `TestAtlasOperator_RendersServiceAccountSecret` PASS.

### Task 6: `atlas-operator` chart — both-set rejection test

**Files:**

- Create: `test/helm/atlas_operator_both_values.yaml`
- Modify: `test/helm/atlas_operator_secret_test.go` (append)

- [ ] **Step 6.1: Write the both-set values fixture**

Write `test/helm/atlas_operator_both_values.yaml`:

```yaml
globalConnectionSecret:
  orgId: "6500000000000000000000aa"
  publicApiKey: "abcdefgh"
  privateApiKey: "12345678-1234-1234-1234-1234567890ab"
  clientId: "mdb_sa_id_01234567890abcdef"
  clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"
```

- [ ] **Step 6.2: Append the rejection test**

Append to `test/helm/atlas_operator_secret_test.go`:

```go
func TestAtlasOperator_RejectsBothCredentialTypes(t *testing.T) {
	_, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_operator_both_values.yaml",
		atlasOperatorChartPath,
	)
	require.Error(t, err, "expected helm template to fail when both credential types are set")
	assert.Contains(t, stderr, "set either (publicApiKey,privateApiKey) or (clientId,clientSecret), not both",
		"stderr did not include the mutual-exclusion message; got: %s", stderr)
}
```

- [ ] **Step 6.3: Run all atlas-operator tests**

```bash
go test ./test/helm/ -run "TestAtlasOperator_" -v
```

Expected: three PASS (`APIKeySecret`, `ServiceAccountSecret`, `RejectsBothCredentialTypes`).

### Task 7: `atlas-operator` chart — README

**Files:**

- Modify: `helm-charts/atlas-operator/README.md`

- [ ] **Step 7.1: Add an SA install example below the existing API-key example**

Find the existing `--set globalConnectionSecret.publicApiKey=...` block (around lines 67–69) and add a second install-example block immediately after, separated by a blank line and a heading:

```markdown
Alternatively, you can install the Operator using Atlas Service Account
credentials (recommended) instead of an API key:

    helm install mongodb-atlas-operator mongodb/mongodb-atlas-operator \
    --set globalConnectionSecret.clientId=<the_client_id> \
    --set globalConnectionSecret.clientSecret=<the_client_secret> \
    --set globalConnectionSecret.orgId=<the_org_id>
```

(Match the indentation / style of the existing snippet — same 4-space indent if that's what's used, same command-name formatting.)

- [ ] **Step 7.2: Re-run tests and commit the whole chunk**

```bash
go test ./test/helm/ -v
git add helm-charts/atlas-operator/values.yaml \
        helm-charts/atlas-operator/templates/global-secret.yaml \
        helm-charts/atlas-operator/README.md \
        test/helm/helm_cli.go \
        test/helm/atlas_operator_apikey_values.yaml \
        test/helm/atlas_operator_sa_values.yaml \
        test/helm/atlas_operator_both_values.yaml \
        test/helm/atlas_operator_secret_test.go
git commit -m "atlas-operator chart: support Service Account credentials"
```

---

## Chunk 3: `atlas-advanced` chart

Same pattern as chunk 2 but with field names `publicKey`/`privateKey` (no `Api` suffix) and `orgID` (capital ID), values path `.Values.secret`, template always renders (no `{{- if .Values.secret.publicKey }}` guard today).

### Task 8: `atlas-advanced` — API-key baseline test and fixture

**Files:**

- Create: `test/helm/atlas_advanced_apikey_values.yaml`
- Create: `test/helm/atlas_advanced_secret_test.go`

- [ ] **Step 8.1: Fixture**

`test/helm/atlas_advanced_apikey_values.yaml`:

```yaml
secret:
  orgID: "6500000000000000000000aa"
  publicKey: "abcdefgh"
  privateKey: "12345678-1234-1234-1234-1234567890ab"

project:
  name: "test-project"
```

- [ ] **Step 8.2: Write the baseline test**

`test/helm/atlas_advanced_secret_test.go`:

```go
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const atlasAdvancedChartPath = "../../helm-charts/atlas-advanced"

func TestAtlasAdvanced_RendersAPIKeySecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_advanced_apikey_values.yaml",
		atlasAdvancedChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "abcdefgh", string(secret.Data["publicApiKey"]))
	assert.Equal(t, "12345678-1234-1234-1234-1234567890ab", string(secret.Data["privateApiKey"]))
	assert.NotContains(t, secret.Data, "clientId")
	assert.NotContains(t, secret.Data, "clientSecret")
}
```

- [ ] **Step 8.3: Verify baseline PASS**

```bash
go test ./test/helm/ -run "TestAtlasAdvanced_RendersAPIKeySecret" -v
```

Expected: PASS (template unchanged so far; baseline matches).

### Task 9: `atlas-advanced` — SA rendering test + template change

**Files:**

- Create: `test/helm/atlas_advanced_sa_values.yaml`
- Modify: `test/helm/atlas_advanced_secret_test.go` (append)
- Modify: `helm-charts/atlas-advanced/values.yaml`
- Modify: `helm-charts/atlas-advanced/templates/atlas-secret.yaml`

- [ ] **Step 9.1: SA fixture**

`test/helm/atlas_advanced_sa_values.yaml`:

```yaml
secret:
  orgID: "6500000000000000000000aa"
  clientId: "mdb_sa_id_01234567890abcdef"
  clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"

project:
  name: "test-project"
```

- [ ] **Step 9.2: Append SA test**

Append to `test/helm/atlas_advanced_secret_test.go`:

```go
func TestAtlasAdvanced_RendersServiceAccountSecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_advanced_sa_values.yaml",
		atlasAdvancedChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "mdb_sa_id_01234567890abcdef", string(secret.Data["clientId"]))
	assert.Equal(t, "mdb_sa_sk_01234567890abcdefghijklmnop", string(secret.Data["clientSecret"]))
	assert.NotContains(t, secret.Data, "publicApiKey")
	assert.NotContains(t, secret.Data, "privateApiKey")
}
```

- [ ] **Step 9.3: Run — expect FAIL**

```bash
go test ./test/helm/ -run "TestAtlasAdvanced_RendersServiceAccountSecret" -v
```

Expected: FAIL — the current template renders `publicApiKey`/`privateApiKey` from empty strings, and does not emit `clientId`/`clientSecret`.

- [ ] **Step 9.4: Update `values.yaml`**

Modify `helm-charts/atlas-advanced/values.yaml`. Find the `secret:` block at the top (lines 1–4) and replace with:

```yaml
# Configure EITHER the API-key fields (publicKey + privateKey) OR the Service
# Account fields (clientId + clientSecret). Setting both at the same time is
# rejected at chart-render time.
secret:
  orgID: ""
  # API key (legacy)
  privateKey: ""
  publicKey: ""
  # Service Account (recommended)
  clientId: ""
  clientSecret: ""
```

- [ ] **Step 9.5: Update `templates/atlas-secret.yaml`**

Replace the full content of `helm-charts/atlas-advanced/templates/atlas-secret.yaml` with:

```
{{- if and .Values.secret.publicKey .Values.secret.clientSecret }}
{{- fail "secret: set either (publicKey,privateKey) or (clientId,clientSecret), not both" }}
{{- end }}
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ include "atlas-advanced.fullname" . }}-secret
  namespace: {{ .Release.Namespace }}
  labels:
    atlas.mongodb.com/type: "credentials"
data:
  orgId: {{ .Values.secret.orgID | b64enc }}
  {{- if .Values.secret.clientSecret }}
  clientId: {{ .Values.secret.clientId | b64enc }}
  clientSecret: {{ .Values.secret.clientSecret | b64enc }}
  {{- else }}
  publicApiKey: {{ .Values.secret.publicKey | b64enc }}
  privateApiKey: {{ .Values.secret.privateKey | b64enc }}
  {{- end }}
```

Note: unlike `atlas-operator`, the `atlas-advanced` chart always renders a Secret (no outer `{{- if .Values.secret.publicKey }}` guard today — the template just assumes the user filled in values). Preserve that behaviour. The inner branch picks data fields based on `.Values.secret.clientSecret`. This preserves the current default-rendering behaviour for users who set no values at all (they get an empty-orgId/empty-keys Secret, which fails at operator runtime — same as today).

- [ ] **Step 9.6: Run both tests — expect PASS**

```bash
go test ./test/helm/ -run "TestAtlasAdvanced_Renders" -v
```

Expected: PASS.

### Task 10: `atlas-advanced` — both-set rejection test

**Files:**

- Create: `test/helm/atlas_advanced_both_values.yaml`
- Modify: `test/helm/atlas_advanced_secret_test.go` (append)

- [ ] **Step 10.1: Fixture**

`test/helm/atlas_advanced_both_values.yaml`:

```yaml
secret:
  orgID: "6500000000000000000000aa"
  publicKey: "abcdefgh"
  privateKey: "12345678-1234-1234-1234-1234567890ab"
  clientId: "mdb_sa_id_01234567890abcdef"
  clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"

project:
  name: "test-project"
```

- [ ] **Step 10.2: Append rejection test**

Append to `test/helm/atlas_advanced_secret_test.go`:

```go
func TestAtlasAdvanced_RejectsBothCredentialTypes(t *testing.T) {
	_, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_advanced_both_values.yaml",
		atlasAdvancedChartPath,
	)
	require.Error(t, err, "expected helm template to fail when both credential types are set")
	assert.Contains(t, stderr, "set either (publicKey,privateKey) or (clientId,clientSecret), not both",
		"stderr did not include the mutual-exclusion message; got: %s", stderr)
}
```

- [ ] **Step 10.3: Run all atlas-advanced tests**

```bash
go test ./test/helm/ -run "TestAtlasAdvanced_" -v
```

Expected: three PASS.

### Task 11: `atlas-advanced` — README + commit

**Files:**

- Modify: `helm-charts/atlas-advanced/README.md`

- [ ] **Step 11.1: Add SA install example**

Find the existing install snippet in the README and append an SA variant. Use the same `--set secret.clientId=...` / `--set secret.clientSecret=...` / `--set secret.orgID=...` form.

- [ ] **Step 11.2: Commit**

```bash
go test ./test/helm/ -v
git add helm-charts/atlas-advanced/values.yaml \
        helm-charts/atlas-advanced/templates/atlas-secret.yaml \
        helm-charts/atlas-advanced/README.md \
        test/helm/atlas_advanced_apikey_values.yaml \
        test/helm/atlas_advanced_sa_values.yaml \
        test/helm/atlas_advanced_both_values.yaml \
        test/helm/atlas_advanced_secret_test.go
git commit -m "atlas-advanced chart: support Service Account credentials"
```

---

## Chunk 4: `atlas-basic` chart

Identical shape to `atlas-advanced`: same values path (`.Values.secret`), same field names, same "always renders" template behaviour.

### Task 12: `atlas-basic` — API-key baseline

**Files:**

- Create: `test/helm/atlas_basic_apikey_values.yaml`
- Create: `test/helm/atlas_basic_secret_test.go`

- [ ] **Step 12.1: Fixture and baseline test**

Write `test/helm/atlas_basic_apikey_values.yaml`:

```yaml
secret:
  orgID: "6500000000000000000000aa"
  publicKey: "abcdefgh"
  privateKey: "12345678-1234-1234-1234-1234567890ab"

project:
  name: "test-project"

deployment:
  name: "test-deployment"
```

Write `test/helm/atlas_basic_secret_test.go` (full content):

```go
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const atlasBasicChartPath = "../../helm-charts/atlas-basic"

func TestAtlasBasic_RendersAPIKeySecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_basic_apikey_values.yaml",
		atlasBasicChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret)

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "abcdefgh", string(secret.Data["publicApiKey"]))
	assert.Equal(t, "12345678-1234-1234-1234-1234567890ab", string(secret.Data["privateApiKey"]))
	assert.NotContains(t, secret.Data, "clientId")
	assert.NotContains(t, secret.Data, "clientSecret")
}
```

- [ ] **Step 12.2: Verify baseline**

```bash
go test ./test/helm/ -run TestAtlasBasic_RendersAPIKeySecret -v
```

Expected: PASS.

### Task 13: `atlas-basic` — SA test + template change

**Files:**

- Create: `test/helm/atlas_basic_sa_values.yaml`
- Modify: `test/helm/atlas_basic_secret_test.go` (append)
- Modify: `helm-charts/atlas-basic/values.yaml`
- Modify: `helm-charts/atlas-basic/templates/atlas-secret.yaml`

- [ ] **Step 13.1: SA fixture**

`test/helm/atlas_basic_sa_values.yaml`:

```yaml
secret:
  orgID: "6500000000000000000000aa"
  clientId: "mdb_sa_id_01234567890abcdef"
  clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"

project:
  name: "test-project"

deployment:
  name: "test-deployment"
```

- [ ] **Step 13.2: Append SA test**

Append to `test/helm/atlas_basic_secret_test.go`:

```go
func TestAtlasBasic_RendersServiceAccountSecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_basic_sa_values.yaml",
		atlasBasicChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "mdb_sa_id_01234567890abcdef", string(secret.Data["clientId"]))
	assert.Equal(t, "mdb_sa_sk_01234567890abcdefghijklmnop", string(secret.Data["clientSecret"]))
	assert.NotContains(t, secret.Data, "publicApiKey")
	assert.NotContains(t, secret.Data, "privateApiKey")
}
```

- [ ] **Step 13.3: Run — expect FAIL**

```bash
go test ./test/helm/ -run "TestAtlasBasic_RendersServiceAccountSecret" -v
```

Expected: FAIL — the current template renders `publicApiKey`/`privateApiKey` from empty strings and does not emit `clientId`/`clientSecret`.

- [ ] **Step 13.4: Update `helm-charts/atlas-basic/values.yaml`**

Find the `secret:` block at the top (lines 1–4) and replace with:

```yaml
# Configure EITHER the API-key fields (publicKey + privateKey) OR the Service
# Account fields (clientId + clientSecret). Setting both at the same time is
# rejected at chart-render time.
secret:
  orgID: ""
  # API key (legacy)
  privateKey: ""
  publicKey: ""
  # Service Account (recommended)
  clientId: ""
  clientSecret: ""
```

- [ ] **Step 13.5: Verify starting-state parity with `atlas-advanced`, then update `helm-charts/atlas-basic/templates/atlas-secret.yaml`**

Pre-flight sanity check — before editing, confirm the existing `atlas-basic` template differs from the `atlas-advanced` template only on the include helper name:

```bash
diff -u helm-charts/atlas-advanced/templates/atlas-secret.yaml \
        helm-charts/atlas-basic/templates/atlas-secret.yaml
```

Expected: a single-line diff on the `include` helper. If anything else differs, stop and investigate before continuing.

Then replace the full content of `helm-charts/atlas-basic/templates/atlas-secret.yaml` with:

```
{{- if and .Values.secret.publicKey .Values.secret.clientSecret }}
{{- fail "secret: set either (publicKey,privateKey) or (clientId,clientSecret), not both" }}
{{- end }}
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ include "atlas-basic.fullname" . }}-secret
  namespace: {{ .Release.Namespace }}
  labels:
    atlas.mongodb.com/type: "credentials"
data:
  orgId: {{ .Values.secret.orgID | b64enc }}
  {{- if .Values.secret.clientSecret }}
  clientId: {{ .Values.secret.clientId | b64enc }}
  clientSecret: {{ .Values.secret.clientSecret | b64enc }}
  {{- else }}
  publicApiKey: {{ .Values.secret.publicKey | b64enc }}
  privateApiKey: {{ .Values.secret.privateKey | b64enc }}
  {{- end }}
```

This mirrors the `atlas-advanced` template exactly except for the include helper name (`atlas-basic.fullname`). The chunk preamble's "same 'always renders' behaviour" note applies — the spec's outer `{{- if or <path>.<publicKeyField> <path>.clientSecret }}` guard is intentionally omitted to preserve the existing unconditional-render contract.

- [ ] **Step 13.6: Run both rendering tests — expect PASS**

```bash
go test ./test/helm/ -run "TestAtlasBasic_Renders" -v
```

Expected: PASS.

### Task 14: `atlas-basic` — both-set rejection + commit

**Files:**

- Create: `test/helm/atlas_basic_both_values.yaml`
- Modify: `test/helm/atlas_basic_secret_test.go` (append)

Note: `helm-charts/atlas-basic/` does not ship a README.md. No documentation update is required for this chart.

- [ ] **Step 14.1: Both-set fixture**

`test/helm/atlas_basic_both_values.yaml`:

```yaml
secret:
  orgID: "6500000000000000000000aa"
  publicKey: "abcdefgh"
  privateKey: "12345678-1234-1234-1234-1234567890ab"
  clientId: "mdb_sa_id_01234567890abcdef"
  clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"

project:
  name: "test-project"

deployment:
  name: "test-deployment"
```

- [ ] **Step 14.2: Append rejection test**

Append to `test/helm/atlas_basic_secret_test.go`:

```go
func TestAtlasBasic_RejectsBothCredentialTypes(t *testing.T) {
	_, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_basic_both_values.yaml",
		atlasBasicChartPath,
	)
	require.Error(t, err, "expected helm template to fail when both credential types are set")
	assert.Contains(t, stderr, "set either (publicKey,privateKey) or (clientId,clientSecret), not both",
		"stderr did not include the mutual-exclusion message; got: %s", stderr)
}
```

- [ ] **Step 14.3: Run all atlas-basic tests**

```bash
go test ./test/helm/ -run "TestAtlasBasic_" -v
```

Expected: three PASS.

- [ ] **Step 14.4: Commit**

```bash
go test ./test/helm/ -v
git add helm-charts/atlas-basic/values.yaml \
        helm-charts/atlas-basic/templates/atlas-secret.yaml \
        test/helm/atlas_basic_apikey_values.yaml \
        test/helm/atlas_basic_sa_values.yaml \
        test/helm/atlas_basic_both_values.yaml \
        test/helm/atlas_basic_secret_test.go
git commit -m "atlas-basic chart: support Service Account credentials"
```

---

## Chunk 5: `atlas-deployment` chart

Differences from the other charts:

- Values path is `.Values.atlas.secret` (nested).
- Field names `publicApiKey`/`privateApiKey` (matching `atlas-operator`).
- Template has outer conditionals `{{- if and (not .Values.atlas.secret.global) (not .Values.atlas.secret.existing) }}` — preserve those; only the inner `data:` block gains the SA branch.

### Task 15: `atlas-deployment` — baseline + SA test

**Files:**

- Create: `test/helm/atlas_deployment_apikey_values.yaml`
- Create: `test/helm/atlas_deployment_sa_values.yaml`
- Create: `test/helm/atlas_deployment_secret_test.go`

- [ ] **Step 15.1: API-key fixture**

`test/helm/atlas_deployment_apikey_values.yaml`:

```yaml
atlas:
  secret:
    global: false
    existing: ""
    orgId: "6500000000000000000000aa"
    publicApiKey: "abcdefgh"
    privateApiKey: "12345678-1234-1234-1234-1234567890ab"
    setCustomName: ""
```

- [ ] **Step 15.2: SA fixture**

`test/helm/atlas_deployment_sa_values.yaml`:

```yaml
atlas:
  secret:
    global: false
    existing: ""
    orgId: "6500000000000000000000aa"
    clientId: "mdb_sa_id_01234567890abcdef"
    clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"
    setCustomName: ""
```

- [ ] **Step 15.3: Write baseline + SA tests**

Create `test/helm/atlas_deployment_secret_test.go`:

```go
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const atlasDeploymentChartPath = "../../helm-charts/atlas-deployment"

func TestAtlasDeployment_RendersAPIKeySecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_deployment_apikey_values.yaml",
		atlasDeploymentChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "abcdefgh", string(secret.Data["publicApiKey"]))
	assert.Equal(t, "12345678-1234-1234-1234-1234567890ab", string(secret.Data["privateApiKey"]))
	assert.NotContains(t, secret.Data, "clientId")
	assert.NotContains(t, secret.Data, "clientSecret")
}

func TestAtlasDeployment_RendersServiceAccountSecret(t *testing.T) {
	stdout, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_deployment_sa_values.yaml",
		atlasDeploymentChartPath,
	)
	require.NoError(t, err, "stderr: %s", stderr)

	secret := findCredentialsSecret(t, stdout)
	require.NotNil(t, secret, "expected a credentials Secret in rendered output")

	assert.Equal(t, "6500000000000000000000aa", string(secret.Data["orgId"]))
	assert.Equal(t, "mdb_sa_id_01234567890abcdef", string(secret.Data["clientId"]))
	assert.Equal(t, "mdb_sa_sk_01234567890abcdefghijklmnop", string(secret.Data["clientSecret"]))
	assert.NotContains(t, secret.Data, "publicApiKey")
	assert.NotContains(t, secret.Data, "privateApiKey")
}
```

- [ ] **Step 15.4: Run both tests — expect baseline PASS, SA FAIL**

```bash
go test ./test/helm/ -run "TestAtlasDeployment_Renders" -v
```

Expected:
- `TestAtlasDeployment_RendersAPIKeySecret`: PASS (template unchanged so far).
- `TestAtlasDeployment_RendersServiceAccountSecret`: FAIL — the current template writes `publicApiKey`/`privateApiKey` from empty strings and does not emit `clientId`/`clientSecret`.

- [ ] **Step 15.5: Update `helm-charts/atlas-deployment/values.yaml`**

Find the `atlas.secret:` block (lines 13–24) and add `clientId` and `clientSecret` fields directly below `privateApiKey`:

```yaml
# Please provide Atlas API credentials and Organization.
# Configure EITHER the API-key fields (publicApiKey + privateApiKey) OR the
# Service Account fields (clientId + clientSecret). Setting both at the same
# time is rejected at chart-render time.
atlas:
  secret:
    # project uses Global Key (highest priority)
    global: false
    # secret already exist in the same namespace
    existing: ""

    orgId: "<Atlas Organization ID>"
    # API key (legacy)
    publicApiKey: "<Atlas_api_public_key>"
    privateApiKey: "<Atlas_api_private_key>"
    # Service Account (recommended)
    clientId: ""
    clientSecret: ""
    # use custom secret name during new secret creation
    setCustomName: ""
```

- [ ] **Step 15.6: Update `helm-charts/atlas-deployment/templates/atlas-secret.yaml`**

Replace the full content of the template with:

```
{{- if and (not .Values.atlas.secret.global) (not .Values.atlas.secret.existing) }}
{{- if and .Values.atlas.secret.publicApiKey .Values.atlas.secret.clientSecret }}
{{- fail "atlas.secret: set either (publicApiKey,privateApiKey) or (clientId,clientSecret), not both" }}
{{- end }}
apiVersion: v1
kind: Secret
type: Opaque
metadata:
{{- if .Values.atlas.secret.setCustomName }}
  name: {{ .Values.atlas.secret.setCustomName}}
{{- else }}
  name: {{ include "atlas-deployment.fullname" . }}-secret
{{- end }}
  namespace: {{ .Release.Namespace }}
  labels:
    atlas.mongodb.com/type: "credentials"
    {{- include "atlas-deployment.labels" . | nindent 4 }}
  annotations:
    'helm.sh/hook': post-delete,pre-install,pre-upgrade
data:
    orgId: {{ .Values.atlas.secret.orgId | b64enc }}
    {{- if .Values.atlas.secret.clientSecret }}
    clientId: {{ .Values.atlas.secret.clientId | b64enc }}
    clientSecret: {{ .Values.atlas.secret.clientSecret | b64enc }}
    {{- else }}
    publicApiKey: {{ .Values.atlas.secret.publicApiKey | b64enc }}
    privateApiKey: {{ .Values.atlas.secret.privateApiKey | b64enc }}
    {{- end }}
{{- end }}
```

Note: the outer `{{- if and (not .global) (not .existing) }}` preserves the existing "only render when not reusing an existing Secret" behaviour. The mutual-exclusion `fail` lives inside that outer guard so that consumers who use `global` or `existing` to point at an out-of-chart Secret are never blocked by a stray combination of unused API-key and SA fields. The inner branch picks API-key vs SA data based on `.clientSecret`.

- [ ] **Step 15.7: Run both tests — expect PASS**

```bash
go test ./test/helm/ -run "TestAtlasDeployment_Renders" -v
```

Expected: both PASS.

### Task 16: `atlas-deployment` — both-set + README + commit

**Files:**

- Create: `test/helm/atlas_deployment_both_values.yaml`
- Modify: `test/helm/atlas_deployment_secret_test.go` (append)
- Modify: `helm-charts/atlas-deployment/README.md`

- [ ] **Step 16.1: Both-set fixture**

`test/helm/atlas_deployment_both_values.yaml`:

```yaml
atlas:
  secret:
    global: false
    existing: ""
    orgId: "6500000000000000000000aa"
    publicApiKey: "abcdefgh"
    privateApiKey: "12345678-1234-1234-1234-1234567890ab"
    clientId: "mdb_sa_id_01234567890abcdef"
    clientSecret: "mdb_sa_sk_01234567890abcdefghijklmnop"
    setCustomName: ""
```

- [ ] **Step 16.2: Append rejection test**

Append to `test/helm/atlas_deployment_secret_test.go`:

```go
func TestAtlasDeployment_RejectsBothCredentialTypes(t *testing.T) {
	_, stderr, err := helmTemplate(t,
		"--namespace=default",
		"--values=atlas_deployment_both_values.yaml",
		atlasDeploymentChartPath,
	)
	require.Error(t, err, "expected helm template to fail when both credential types are set")
	assert.Contains(t, stderr, "set either (publicApiKey,privateApiKey) or (clientId,clientSecret), not both",
		"stderr did not include the mutual-exclusion message; got: %s", stderr)
}
```

- [ ] **Step 16.3: Run all atlas-deployment tests (and TestFlexSpec)**

```bash
go test ./test/helm/ -run "TestAtlasDeployment_|TestFlexSpec" -v
```

Expected: three `TestAtlasDeployment_*` PASS, and `TestFlexSpec` still PASS (unchanged).

- [ ] **Step 16.4: README update**

Open `helm-charts/atlas-deployment/README.md`. Find the existing install block (around lines 44–50) — the real block includes `--namespace=...` and `--create-namespace` flags that are elided below for brevity; keep them untouched, we are only appending below this block, not replacing it:

````markdown
```shell
helm install atlas-deployment mongodb/atlas-deployment\
    --set project.atlasProjectName='My Project' \
    --set atlas.secret.orgId='<orgid>' \
    --set atlas.secret.publicApiKey='<publicKey>' \
    --set atlas.secret.privateApiKey='<privateApiKey>'
```
````

Immediately below it, append a Service Account variant:

````markdown
Or with a MongoDB Atlas Service Account (recommended):

```shell
helm install atlas-deployment mongodb/atlas-deployment \
    --set project.atlasProjectName='My Project' \
    --set atlas.secret.orgId='<orgid>' \
    --set atlas.secret.clientId='<clientId>' \
    --set atlas.secret.clientSecret='<clientSecret>'
```

Setting both API-key and Service Account fields at the same time is rejected at chart-render time.
````

- [ ] **Step 16.5: Commit**

```bash
go test ./test/helm/ -v
git add helm-charts/atlas-deployment/values.yaml \
        helm-charts/atlas-deployment/templates/atlas-secret.yaml \
        helm-charts/atlas-deployment/README.md \
        test/helm/atlas_deployment_apikey_values.yaml \
        test/helm/atlas_deployment_sa_values.yaml \
        test/helm/atlas_deployment_both_values.yaml \
        test/helm/atlas_deployment_secret_test.go
git commit -m "atlas-deployment chart: support Service Account credentials"
```

---

## Chunk 6: TD update + final verification

### Task 17: Rewrite the TD Installation section

**Files:**

- Modify: `docs/dev/td-service-accounts-phase1.md` — Installation section.

- [ ] **Step 17.1: Replace the Installation section**

Find `## Installation` and replace the body (Helm Chart Changes + OLM / CSV Changes subsections) with actual-state content:

```markdown
## Installation

### Helm Charts

Every operator-adjacent chart that creates a credentials Secret from its
`values.yaml` now renders either an API-key Secret or a Service Account
Secret based on which fields are populated:

- `atlas-operator` — `.Values.globalConnectionSecret.clientId` / `.clientSecret`.
- `atlas-advanced` — `.Values.secret.clientId` / `.clientSecret`.
- `atlas-basic` — `.Values.secret.clientId` / `.clientSecret`.
- `atlas-deployment` — `.Values.atlas.secret.clientId` / `.clientSecret`.

Each chart keeps its existing API-key values keys unchanged and adds
`clientId` / `clientSecret` next to them. The chart template emits the
Secret with data keys `orgId + clientId + clientSecret` when SA fields
are populated, otherwise `orgId + publicApiKey + privateApiKey`.
Setting both pairs fails chart rendering with a clear message, matching
the operator's `validateConnectionSecret` behaviour at runtime.

Each chart's README now documents both install variants with `--set`
examples.

### OLM / CSV

The `ServiceAccountToken` controller's `// +kubebuilder:rbac:...` markers
grant `create`, `update`, `patch` on Secrets in addition to the existing
`get`, `list`, `watch`. Running `make manifests && make bundle &&
make helm-crds` propagates the enriched verbs through:

1. `config/rbac/role.yaml` (and split cluster-scoped / namespaced variants)
   — via `controller-gen` + `scripts/split_roles_yaml.sh`.
2. `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml`
   — via `operator-sdk generate bundle`.
3. `helm-charts/atlas-operator/rbac.yaml` — via `yq` in `make helm-crds`.

The `atlas-operator` chart's `cluster-roles.yaml` and `roles.yaml` already
`.Files.Lines "rbac.yaml"`, so the enriched verbs flow through to both
deployment modes (cluster-wide and namespaced) without hand-edits. No
new env vars, `installModes`, or watch-namespace semantics.

### Testing

Install-time rendering is pinned by unit tests under `test/helm/` that run
`helm template` against each chart and assert the resulting Secret (or,
for the mutual-exclusion case, assert rendering fails with the expected
message). An additional unit test pins the Secret verbs in the rendered
`atlas-operator` ClusterRole, catching any regen-chain drift.

Full install gates: `make bundle-validate` (OLM schema), `make
validate-manifests` (regen drift), `make validate-crds-chart` (chart CRD
parity with the bundle), `make unit-test` (chart rendering + RBAC).
```

- [ ] **Step 17.2: Commit the TD**

```bash
git add docs/dev/td-service-accounts-phase1.md
git commit -m "TD: document Helm chart + OLM updates for Service Account support"
```

### Task 18: Full-branch verification

- [ ] **Step 18.1: Run the full helm test package**

```bash
go test ./test/helm/ -v
```

Expected: 13 tests PASS (1 flex + 1 rbac + 4 charts × 3 scenarios). No failures, no skipped.

- [ ] **Step 18.2: Run the full repo gate**

```bash
devbox run -- make ci
```

Expected: PASS. If `make ci` isn't wired to run every gate we touched, instead run individually:

```bash
devbox run -- make fmt lint unit-test validate-manifests validate-crds-chart bundle-validate
```

Expected: each target exits 0.

- [ ] **Step 18.3: Grep for stale references**

```bash
# Any accidental "charts/" (the old path from the first TD draft)?
grep -rn 'charts/atlas-operator\b' docs/ 2>&1 | grep -v 'helm-charts/' | head -5
# Any leftover references to the old "API key only" assumption in the chart READMEs?
grep -rn 'Service Account' helm-charts/*/README.md 2>&1 | head -10
```

Expected from the first grep: no matches. Expected from the second: at least one match per chart's README (i.e. every README got its SA example).

- [ ] **Step 18.4: Branch ready**

```bash
git log --oneline origin/CLOUDP-373260/service-accounts-phase-1..HEAD
```

Expected to show 6 commits (one per chunk — RBAC regen, atlas-operator, atlas-advanced, atlas-basic, atlas-deployment, TD).
