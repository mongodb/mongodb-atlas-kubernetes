# Design: Service Account Credentials in Installation Methods

**Author:** Maciej Karaś
**Date:** 2026-04-22
**Related:** [`docs/dev/td-service-accounts-phase1.md`](../../dev/td-service-accounts-phase1.md), Service Accounts Phase 1

## Context

Service Account support at runtime (Phase 1) is complete: the operator consumes Secrets carrying `orgId + clientId + clientSecret`, the `ServiceAccountToken` controller mints bearer tokens, and all resource reconcilers read them through `GetConnectionConfig`. None of the installation methods, however, know how to produce a Secret in the SA shape, and the generated RBAC artefacts on disk have not been refreshed since the new controller's markers were added.

Specifically:

- The four operator-adjacent Helm charts (`atlas-operator`, `atlas-advanced`, `atlas-basic`, `atlas-deployment`) only render API-key Secrets from their `values.yaml`. A user who wants to install the operator using Service Account credentials has to bypass the chart and apply a Secret by hand.
- `config/rbac/role.yaml` (split into cluster-scoped and namespaced variants) is stale — it pre-dates the `ServiceAccountTokenReconciler` markers that grant `create/update/patch` on Secrets, so a fresh `make manifests` would produce a diff.
- The OLM bundle CSV at `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml` still embeds the old role. A fresh `make bundle` would update it, and `make helm-crds` would in turn propagate the refreshed `clusterPermissions[0].rules` into `helm-charts/atlas-operator/rbac.yaml`, which is sourced by the chart's ClusterRole / Role templates.
- No test currently exercises the chart rendering for different credential types. `test/helm/flex_test.go` is the only helm-template-based unit test and it covers the `atlas-deployment` Flex spec specifically.

This spec records the concrete changes needed to bring all four installation methods up to parity with the runtime feature, plus a helm-template-level unit-test suite.

## Goals

- Every Helm chart that currently creates a Connection Secret from its `values.yaml` can also create one in Service Account shape. Both shapes are mutually exclusive at chart render time.
- The generated RBAC artefacts on disk (`config/rbac/`, bundle CSV, `helm-charts/atlas-operator/rbac.yaml`) reflect the `ServiceAccountToken` controller's markers.
- OLM bundle remains valid (`make bundle-validate`).
- Every chart's rendering contract is pinned by a unit test that runs under `make unit-test` — no live cluster required.
- Existing users (API-key paths) see no change in their values layout or resulting manifests.

## Non-Goals

- Changes to `atlas-advanced`, `atlas-deployment`, `atlas-basic` field-naming inconsistencies (`orgID` vs `orgId`, `publicKey` vs `publicApiKey`). Each chart keeps its established API-key naming; the new SA keys use a uniform `clientId`/`clientSecret` across all charts.
- `scripts/openshift/*` and the OpenShift upgrade test — they consume the bundle image, so the regenerated CSV flows through without per-script edits.
- An end-to-end install-via-helm test that runs against a live Atlas org. The runtime SA path is already covered by `test/e2e2/serviceaccount_test.go`; chart-level unit tests are sufficient evidence that the install manifests are correct.
- Chart-level Secret rotation. Rotation remains the runtime concern of the controllers.

## Design

### 1. Helm chart edits

Four charts render a credentials Secret from `values.yaml`. Each follows the same pattern: add SA credential fields to `values.yaml`, update the template to select between API-key and SA data fields, fail at render time if both are configured, and document the new option in the README.

**Per-chart paths** (the names are inconsistent by existing convention — leave them alone for backwards compatibility):

- `atlas-operator` — `values.yaml` key path `.Values.globalConnectionSecret`. Template `templates/global-secret.yaml`.
- `atlas-advanced` — `values.yaml` key path `.Values.secret`. Template `templates/atlas-secret.yaml`.
- `atlas-basic` — `values.yaml` key path `.Values.secret`. Template `templates/atlas-secret.yaml`.
- `atlas-deployment` — `values.yaml` key path `.Values.atlas.secret`. Template `templates/atlas-secret.yaml`.

**`values.yaml` additions (per chart, under the corresponding key path):**

```yaml
# ...existing API-key fields stay unchanged...

# Service Account credentials (alternative to API key).
# Set EITHER (publicApiKey,privateApiKey) OR (clientId,clientSecret), not both.
clientId: ""
clientSecret: ""
```

The `orgId` field is shared; already present in every chart.

**Template rules (apply uniformly):** pseudocode in terms of `<path>` (each chart's values key) and `<publicKeyField>` / `<privateKeyField>` (each chart's legacy API-key field names — `publicApiKey`/`privateApiKey` in `atlas-operator` and `atlas-deployment`, but `publicKey`/`privateKey` in `atlas-advanced` and `atlas-basic`; do NOT rewrite them).

```
{{- /* Mutual exclusion — fail at render time, matching the operator's
       validateConnectionSecret behaviour at runtime. */ -}}
{{- if and <path>.<publicKeyField> <path>.clientSecret }}
  {{- fail "<chart>: set either (<publicKeyField>,<privateKeyField>) or (clientId,clientSecret), not both" }}
{{- end }}

{{- if or <path>.<publicKeyField> <path>.clientSecret }}
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: ...
  namespace: ...
  labels:
    atlas.mongodb.com/type: "credentials"
    # ...other labels retained as today...
data:
    orgId: {{ <path>.orgId | b64enc }}       # or <path>.orgID where that is the existing name
    {{- if <path>.<publicKeyField> }}
    publicApiKey: {{ <path>.<publicKeyField> | b64enc }}
    privateApiKey: {{ <path>.<privateKeyField> | b64enc }}
    {{- else }}
    clientId: {{ <path>.clientId | b64enc }}
    clientSecret: {{ <path>.clientSecret | b64enc }}
    {{- end }}
{{- end }}
```

Note that the **Kubernetes Secret data keys are always `orgId`, `publicApiKey`, `privateApiKey`, `clientId`, `clientSecret`** — that is the contract the operator reads. Only the values-file key names (`publicKey` vs `publicApiKey`, `orgID` vs `orgId`) differ across charts.

Per-chart substitution:

| Chart | `<path>` | `<publicKeyField>` | `<privateKeyField>` | values-file `orgId` key |
|---|---|---|---|---|
| `atlas-operator` | `.Values.globalConnectionSecret` | `publicApiKey` | `privateApiKey` | `orgId` |
| `atlas-deployment` | `.Values.atlas.secret` | `publicApiKey` | `privateApiKey` | `orgId` |
| `atlas-advanced` | `.Values.secret` | `publicKey` | `privateKey` | `orgID` |
| `atlas-basic` | `.Values.secret` | `publicKey` | `privateKey` | `orgID` |

Notes:

- The trigger condition (`<publicKeyField> || clientSecret`) is the test for "should we render a Secret at all?". The existing field is the legacy trigger for the API-key path; `clientSecret` is the trigger for the SA path. `clientId` alone is not sufficient because the operator's validation rejects a partial SA pair anyway.
- `atlas-deployment`'s existing `.Values.atlas.secret.global`/`.existing`/`.setCustomName` outer conditionals remain in place — only the inner `data:` block gains the SA branch.
- All four charts' non-`data` blocks (metadata, labels, annotations) stay unchanged beyond what is shown above.

**README updates** — each chart's `README.md` gets a second install example alongside the existing API-key example, showing Service Account credentials being set through `--set`.

### 2. Regenerated RBAC and OLM artefacts

The `ServiceAccountToken` controller already carries the markers (`serviceaccounttoken_controller.go`, lines 53–54):

```go
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch;create;update;patch
```

Running the existing `make` chain regenerates every downstream artefact from that single source of truth:

1. `make manifests` — controller-gen re-reads markers → writes `config/rbac/role.yaml`, then `scripts/split_roles_yaml.sh` emits cluster-scoped + namespaced variants. controller-gen emits the union of verbs across all matching markers per (group, resource), so the prior `get;list;watch`-only rule for Secrets is replaced by a rule carrying `get;list;watch;create;update;patch`.
2. `make bundle` — `operator-sdk generate bundle` consumes the refreshed `config/rbac/` and emits a new `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml`. The CSV's `spec.install.spec.clusterPermissions[0].rules` now contains the enriched Secret verbs.
3. `make helm-crds` — `yq` extracts `clusterPermissions[0].rules` from the CSV into `helm-charts/atlas-operator/rbac.yaml`. The chart's `cluster-roles.yaml` and `roles.yaml` templates already source that file via `.Files.Lines "rbac.yaml"`.
4. `make bundle-validate` — `operator-sdk bundle validate ./bundle` succeeds.

No hand-edits to generated artefacts. If the pipeline produces other incidental diffs (timestamps, creation dates), keep them; they are expected side effects of regen and are already exercised by CI's `make validate-manifests` gate.

### 3. Tests

Extend `test/helm/` with one file per chart, following the existing `flex_test.go` pattern: `helm template` → decode YAML → assert against Go structs.

**New test files:**

- `test/helm/atlas_operator_secret_test.go`
- `test/helm/atlas_advanced_secret_test.go`
- `test/helm/atlas_basic_secret_test.go`
- `test/helm/atlas_deployment_secret_test.go`

Each file covers three scenarios:

1. **API key only** — set the chart's legacy fields. Assert the rendered `corev1.Secret` has data keys `orgId`, `publicApiKey`, `privateApiKey`; no `clientId`/`clientSecret`. Labels include `atlas.mongodb.com/type: credentials`.
2. **Service Account only** — set `clientId`/`clientSecret`. Assert the rendered Secret has data keys `orgId`, `clientId`, `clientSecret`; no `publicApiKey`/`privateApiKey`. Labels include `atlas.mongodb.com/type: credentials`.
3. **Both set → render fails** — invoke `helm template`, assert non-zero exit, assert stderr contains the mutual-exclusion message.

**Additional RBAC regression test** — in `test/helm/atlas_operator_rbac_test.go`, `helm template` the `atlas-operator` chart and decode the ClusterRole. Assert the Secret rule lists verbs `get`, `list`, `watch`, `create`, `update`, `patch`. This pins the regeneration chain so any future change to a kubebuilder marker that doesn't flow through to the helm chart trips CI.

**Test helper** — if stderr-capture for the `{{ fail }}` cases is not already supported by `test/helper/cmd.RunCommand`, add a minimal helper in `test/helm` that runs `helm template` via `exec.Command`, captures stderr, and returns `(stdout, stderr, exitErr)`. Reuse across all four "both set" cases.

**Tooling assumption** — `make unit-test` runs inside `devbox` (per `CLAUDE.md`), which provisions `helm` on `$PATH`. The tests therefore require no additional install step locally or in CI when the `unit-test` job runs under `devbox run`, as it does today.

**Values fixtures** — colocate input fixtures as `*_values.yaml` next to the test files, in the style of the existing `flex_values.yaml`. One fixture per scenario per chart, for readability.

### Files Changed

- `helm-charts/atlas-operator/values.yaml` — add SA fields to `globalConnectionSecret`.
- `helm-charts/atlas-operator/templates/global-secret.yaml` — add mutual-exclusion guard + SA branch.
- `helm-charts/atlas-operator/README.md` — SA install example.
- `helm-charts/atlas-advanced/values.yaml` — add SA fields to `secret`.
- `helm-charts/atlas-advanced/templates/atlas-secret.yaml` — guard + SA branch.
- `helm-charts/atlas-advanced/README.md` — SA install example.
- `helm-charts/atlas-basic/values.yaml` — add SA fields to `secret`.
- `helm-charts/atlas-basic/templates/atlas-secret.yaml` — guard + SA branch.
- `helm-charts/atlas-basic/README.md` — SA install example.
- `helm-charts/atlas-deployment/values.yaml` — add SA fields to `atlas.secret`.
- `helm-charts/atlas-deployment/templates/atlas-secret.yaml` — guard + SA branch.
- `helm-charts/atlas-deployment/README.md` — SA install example.
- `config/rbac/role.yaml` + split cluster/namespaced variants — regenerated by `make manifests`.
- `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml` and any collateral bundle files — regenerated by `make bundle`.
- `helm-charts/atlas-operator/rbac.yaml` — regenerated by `make helm-crds`.
- `test/helm/atlas_operator_secret_test.go` — **new.**
- `test/helm/atlas_advanced_secret_test.go` — **new.**
- `test/helm/atlas_basic_secret_test.go` — **new.**
- `test/helm/atlas_deployment_secret_test.go` — **new.**
- `test/helm/atlas_operator_rbac_test.go` — **new.**
- `test/helm/*_values.yaml` — **new** fixture files for each test scenario.
- `docs/dev/td-service-accounts-phase1.md` — Installation section rewritten to reflect actual state (helm + OLM paths both now support SA) and new test plan entries in the Testing section.

## Edge Cases

- **Empty `orgId`** — not a new case; existing templates already happily render a Secret with an empty `orgId` b64 value. The operator rejects it at runtime via `validateConnectionSecret`. No new guard added.
- **Only `clientId` set, `clientSecret` empty** — our trigger is `clientSecret`, so the Secret is not rendered. Consistent with today's behaviour for the API-key path where `publicApiKey` is the trigger.
- **Only `clientSecret` set, `clientId` empty** — Secret IS rendered with an empty `clientId`. The operator's `validateConnectionSecret` rejects it at runtime. Matches the parallel case for the API-key path (an API-key Secret with only `publicApiKey` set also renders and fails at runtime). No new guard added.
- **Chart upgrades where a prior release had `publicApiKey` set and a new release has `clientSecret` set** — the rendered Secret's name is unchanged; its `data` contents change. Helm's apply behaviour updates the Secret in place. The operator's `credentialsHash` guard detects the change and refreshes the token. No chart-level migration code needed.

## Testing

- New helm-template unit tests above, run by `make unit-test`.
- Existing `make bundle-validate` continues to pass.
- Existing `make validate-manifests` gate catches any regen drift.
- Existing `make validate-crds-chart` gate catches helm chart vs bundle CRD drift.

## Backward Compatibility

- Users whose `values.yaml` has `publicApiKey` + `privateApiKey` set continue to render an identical Secret. No value renames. No field moves.
- Users whose `values.yaml` has neither credential pair set continue to render nothing (trigger unchanged).
- Upgrading a release from API-key values to SA values: the Secret in-cluster is updated in place; operator picks up the rotation via its existing `credentialsHash` detection.

## Open Questions

None.
