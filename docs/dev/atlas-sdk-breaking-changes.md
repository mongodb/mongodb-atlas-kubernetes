# Handling Breaking Changes in Atlas SDK Updates

When the [Weekly Atlas SDK Update](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/update-sdk.yml) PR fails to merge cleanly, classify the breakage and apply the right fix.

## Background

The weekly workflow runs `scripts/update-sdk.sh` then `make clean gen-crds` and opens a PR titled "Atlas SDK Dependency Update". The script bumps `go.mongodb.org/atlas-sdk/<release>` in `go.mod`/imports via `gomajor`, then blindly replaces the old SDK major release string with the new one across `crd2go/openapi2crd.yaml`, `crd2go/openapi2crd.experimental.yaml`, `test/scaffolder/testdata/crds.yaml`, and `test/scaffolder/testdata/atlas-crds.yaml`.

Most bumps merge cleanly. When they don't (compile errors from renamed/removed SDK types, changed signatures), manual intervention is needed.

## Step 1: Classify the Breakage

| Category | Definition | Action |
|---|---|---|
| **Trivial / internal-only** | Break is in internal helpers, mocks, or SDK types not exposed in any CRD spec/status. | Fix in place (Step 2). No versioning needed. |
| **Non-breaking / additive** | New optional fields, new enum values, new endpoints, docs-only changes. | Fix in place (Step 2). |
| **Breaking (public API)** | A field exposed in the operator's public customer-facing API surface is removed, renamed, or its behavior/contract changed (validation, defaults, enum meanings). | Add new SDK version alongside existing (Step 3). File upstream design bug for behavior changes. |

The operator has **three relevant Go packages**, and they must not be confused:

| Package | Backing config | Status | Breaking-change procedure? |
|---|---|---|---|
| `api/v1/` | — (hand-written) | **Maintenance mode**: legacy hand-written CRD types (`AtlasProject`, `AtlasDeployment`, `AtlasDatabaseUser`, etc.). Do the bare minimum to keep it compiling; do not add elaborate versioning schemes here. | Minimal fix only. No version splitting. |
| `generated/v1/` | `crd2go/openapi2crd.yaml` | **PUBLIC, customer-facing generated API surface** — produced by `crd2go`, backing `atlas.generated.mongodb.com` CRDs (Cluster, Group, DatabaseUser, FlexCluster, IPAccessListEntry, etc.) that real customers deploy. **Breaking changes here directly break customers.** | Yes — Step 3 version-splitting procedure applies. |
| `internal/nextapi/generated/v1/` | `crd2go/openapi2crd.experimental.yaml` | **Experimental, non-public, internal-only code** (gated behind `EXPERIMENTAL=1`). Changes here do NOT break customers. | No — treat as trivial/internal (Step 2) even if fields are removed/renamed. |

### Decision Checklist

- `go build ./...` / `make unit-test` fails only in non-API, non-generated-CRD code? → **trivial fix**
- A field removed/renamed that exists in `api/v1/` → **minimal fix only** (see Step 3 note about `api/v1`); do not add versioning for `api/v1`.
- A field removed/renamed in `crd2go/openapi2crd.yaml` (feeding `generated/v1/`)? → **breaking, customer-facing**, use Step 3.
- A field removed/renamed in `crd2go/openapi2crd.experimental.yaml` (feeding `internal/nextapi/generated/v1/`)? → **trivial/internal**, Step 2, no customer impact.
- A behavior/semantic change to an existing field? → **breaking** (report upstream)
- Otherwise (additive only)? → **non-breaking**, fix in place

## Step 2: Handling Trivial / Non-Breaking Fix

1. Fix the compile/test failure in place.
2. Run `make fmt && make lint && make unit-test` (or `make ci`) before pushing.
3. No versioning changes needed.

This includes all breakages in `internal/nextapi/generated/v1/` (backed by `crd2go/openapi2crd.experimental.yaml`), since that package is experimental/non-public and no customers are affected.

## Step 3: Handling a Breaking Change in the Public Generated API (`generated/v1/`, driven by `crd2go/openapi2crd.yaml`)

This procedure applies **only** to the public customer-facing generated API surface (`generated/v1/`, backed by `crd2go/openapi2crd.yaml`).

When a field is removed or renamed in the Atlas SDK, we split the changed fields off into a **new SDK version entry** in spec/status, while the **old SDK version entry** is left untouched so existing users on the old spec/status version keep working unaffected.

Key structure of `crd2go/openapi2crd.yaml`:

- `openapi:` list keyed by `name:` (e.g. `v20250312`), each with a `package:` SDK import path.
- `crd.mappings` entries reference `majorVersion:` + `openAPIRef.name:` pointing at an `openapi` entry.
- CRD status nests fields per version, e.g. `$.status.v20250312.id`.

This work might require explicit effort allocation and re-planing, when unexpected.

### Procedure

1. **Keep the existing spec|status/sdk-version entry completely untouched/working.** Do not let the blind `sed` in `scripts/update-sdk.sh` rename it in place for a breaking bump. If the rename already happened, revert it in `crd2go/openapi2crd.yaml`, `crd2go/openapi2crd.experimental.yaml`, and `test/scaffolder/testdata/*.yaml` — the old SDK version entry must keep working for users of that spec/status version.

2. **Split off the removed/renamed/changed fields into a NEW spec|status/sdk-version entry** (e.g. `v20260101`) alongside the existing one(s). Add a new item in the `openapi:` list (both `name:` and `package:` fields) and new `crd.mappings` entries (with `majorVersion` and `openAPIRef` blocks) pointing at this new version for the affected CRD kinds. Both old and new versions generate side-by-side in `generated/v1/`'s spec and status, e.g. both `$.status.v20250312.id` and `$.status.v20260101.id` exist after regeneration. The old version keeps working unchanged for existing customers; the new version carries the breaking change.

3. **Regenerate:** run `make clean gen-crds` (and `make gen-all` if Go types/controllers need scaffolding). Verify both versions appear in `config/generated/crd/bases/crds.yaml` and in the regenerated `generated/v1/` Go types.

4. **If `api/v1/` is also affected** (compile break in the legacy hand-written types): fix it with the bare minimum change — rename the Go field to match the new SDK name, update references in the legacy controller code, and stop there. Do **not** extend `api/v1/` with elaborate new-version schemes; this package is in maintenance mode. Do not replicate the version-splitting scheme from Step 3 there.

5. **Evaluate pruning the oldest unsupported version.** The moment a new spec|status/sdk-version entry is split off is exactly the trigger to check whether the OLDEST entry in `crd2go/openapi2crd.yaml` (its `openapi:` entry, `crd.mappings`, generated status blocks, and any test fixtures in `test/scaffolder/testdata/`) can now be removed entirely. Coordinate with the team on the current deprecation/support policy. Either:
    - Prune it now (in this PR or a closely tracked follow-up), or
    - Explicitly document why it cannot be pruned yet (e.g. users still on that version per telemetry, or the deprecation window hasn't elapsed).

    This evaluation applies only to the oldest supported version when it is already unsupported. In that case, it is a required step of the breaking-change workflow — it must not be deferred indefinitely.

6. **File an upstream issue** with the Atlas API/SDK team for any contract behavior change on an existing field. Behavior changes must only ship under new version names, never silently on an existing one.

7. Run `make fmt && make lint && make unit-test` (or `make ci`), plus `make manifests` if CRDs under `config/crd/bases` changed, before pushing.

## Related Files

- `.github/workflows/update-sdk.yml` — Weekly workflow that bumps the Atlas SDK and opens a PR
- `scripts/update-sdk.sh` — SDK bump script
- `crd2go/openapi2crd.yaml` — CRD generation config that generates the PUBLIC `generated/v1/` API (breaking-change procedure applies)
- `generated/v1/` — Public, customer-facing generated Go types (Cluster, Group, DatabaseUser, FlexCluster, IPAccessListEntry, etc.); do not edit by hand
- `crd2go/openapi2crd.experimental.yaml` — CRD generation config that generates the experimental, non-public `internal/nextapi/generated/v1/` API (no customer impact, trivial-fix path)
- `internal/nextapi/generated/v1/` — Experimental/internal generated Go types; do not edit by hand
- `test/scaffolder/testdata/crds.yaml` — Scaffolder test fixture
- `test/scaffolder/testdata/atlas-crds.yaml` — Scaffolder test fixture
- `api/v1/` — Legacy hand-written public API types, **maintenance mode** — do not introduce elaborate versioning here

