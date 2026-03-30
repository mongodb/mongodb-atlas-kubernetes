# Technical Design: Atlas Kubernetes Operator — Service Accounts Support, Phase 1

**Author:** Maciej Karaś
**Date:** 2026-03-30
**Jira Epic:** [CLOUDP-383935](https://jira.mongodb.org/browse/CLOUDP-383935)
**Scope Document:** [PD+Scope: Service Accounts support — Phase 1](https://docs.google.com/document/d/1Kn8bwUWXx4Piojzd2eVh7QXgpymPciivi1h78pTSa7I)
**Spike:** [Spike: Service Accounts support in AKO](https://docs.google.com/document/d/1FwV7GfC9i__dRGP6m3cExj7t3hc0wsKuB0dWDwYK0Mc)

Please add your LGTM to this document.

---

## Overview

MongoDB recommends Atlas API Service Accounts (OAuth 2.0 Client Credentials) as the preferred method for programmatic access to the Atlas Administration API. API keys (HTTP Digest) are labelled as legacy in public documentation.

Today, the Atlas Kubernetes Operator only supports API key–based authentication. This means customers who have standardised on Service Accounts for Terraform, Atlas CLI, or other tooling cannot use them with the operator. This project adds Phase 1 Service Account support: the operator can consume Service Account credentials (`clientId` / `clientSecret`) from a Kubernetes Secret, obtain short-lived OAuth access tokens from Atlas, and use those tokens for all Atlas API calls. No new CRDs, no changes to the custom resource API surface, no changes to how existing API key users operate.

Phase 1 is consumption-only. The operator does not create or manage Service Accounts or their secrets in Atlas. That is explicitly deferred to Phase 2, pending product alignment and a separate security investigation (Separation of Duties).

---

## Terminology

| Term                                 | Definition                                                                                                                                                                                                                                  |
|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Service Account (SA)**             | An Atlas OAuth 2.0 client (client ID + client secret) scoped to one organisation. Used to obtain short-lived access tokens for the Atlas Administration API.                                                                                |
| **Programmatic API Key (PAK)**       | Legacy Atlas authentication: public key + private key, HTTP Digest. Still supported; not deprecated by this project.                                                                                                                        |
| **Connection Secret**                | The Kubernetes Secret referenced by an Atlas custom resource (or the global operator flag) that holds Atlas credentials. Today: `orgId + publicApiKey + privateApiKey`. After this project: also accepts `orgId + clientId + clientSecret`. |
| **Access Token Secret**              | A Kubernetes Secret managed by the `ServiceAccountToken` controller. Holds the short-lived OAuth bearer token and its expiry timestamp. Never created or edited by users.                                                                   |
| **Bearer Token**                     | A short-lived OAuth 2.0 access token (1 hour / 3600 seconds) returned by the Atlas token endpoint. Used in `Authorization: Bearer <token>` on all Atlas API requests.                                                                       |
| **`ServiceAccountToken` controller** | New controller introduced in this project. Watches Connection Secrets, detects Service Account credentials, and manages the Access Token Secret lifecycle (create, refresh).                                                                |

---

## Background

### Current Authentication Flow

All resource reconcilers (AtlasProject, AtlasDeployment, AtlasDatabaseUser, and generated controllers) share a common credential-reading path:

1. Read the Connection Secret referenced by the custom resource (or the global secret).
2. Call `GetConnectionConfig()` in `internal/controller/reconciler/credentials.go`, which returns a `ConnectionConfig` containing `Credentials.APIKeys`.
3. Call `provider.SdkClientSet()` in `internal/controller/atlas/provider.go`, which wraps the credentials in `digest.NewTransport` (HTTP Digest) and creates an Atlas SDK client.

The credential Secret shape today:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mongodb-atlas-operator-api-key
  labels:
    atlas.mongodb.com/type: credentials
type: Opaque
stringData:
  orgId: "66d9b0a20167996f36ac694e"
  publicApiKey: "66d9198f-72a3-46ab-b741-18ccc6aa835e"
  privateApiKey: "cwcnxswp"
```

### Problem

- Service Account credentials cannot be supplied to the operator today.
- Every Atlas API call uses a fresh Digest nonce derived from the long-lived key pair. There is no token to manage.
- Service accounts produce short-lived OAuth access tokens (1 hour) that must be proactively refreshed. If each of the operator's many reconcilers requested a fresh token on every reconcile, this would generate thousands of `POST /api/oauth/token` requests per hour and risk rate-limiting (`HTTP 429`).

---

## Goals

- Support `clientId` / `clientSecret` in Connection Secrets. When these fields are present, the operator uses OAuth 2.0 for Atlas API calls.
- Centralise token lifecycle (acquire, cache, refresh) in a single dedicated controller. All reconcilers share one token per credential Secret; no reconciler performs token acquisition itself.
- Token naming is deterministic and requires no mutation of user-provided Secrets. GitOps-safe: ArgoCD / Flux workflows that reconcile user Secrets cannot break the token lookup.
- Preserve existing API key behavior unchanged. A Secret with API key fields continues to work exactly as before.
- Mixed Secrets (containing both API key and Service Account fields) are rejected with a clear validation error.
- No CRD or API surface changes. Custom resources continue to reference Connection Secrets the same way.
- Enabled by default. No feature flag or operator configuration required.
- Support all Atlas environments: commercial (`cloud.mongodb.com`), QA (`cloud-qa.mongodb.com`), and Gov (`cloud.mongodbgov.com`).

## Non-Goals

- Creating or managing Service Accounts or their secrets in Atlas (Phase 2).
- Migration tooling from API keys to Service Accounts.
- A new Custom Resource Definition for Service Account credential lifecycle.
- Deprecating Programmatic API Keys.
- Zero-downtime credential secret rotation (may be addressed in a follow-up).

---

## Architecture Overview

```
User creates Connection Secret (clientId/clientSecret, label: atlas.mongodb.com/type: credentials)
     │
     ▼
ServiceAccountToken Controller  (internal/controller/serviceaccount/)
  ├── Watches Secrets with label atlas.mongodb.com/type: credentials
  ├── Skips Secrets without clientId/clientSecret (API key secrets)
  ├── Derives Access Token Secret name deterministically (FNV hash)
  ├── If Access Token Secret missing or expired → calls Atlas OAuth endpoint
  │     POST <atlasDomain>/api/oauth/token  (via Atlas SDK clientcredentials)
  ├── Creates / updates Access Token Secret in place
  │     - label: atlas.mongodb.com/type: credentials  (required for informer cache)
  │     - ownerReference → Connection Secret  (for garbage collection)
  └── Re-enqueues at 2/3 of token TTL (~40 min)
     │
     ▼
Resource Reconcilers (AtlasProject, AtlasDeployment, AtlasDatabaseUser, generated controllers)
  ├── Call GetConnectionConfig() as today
  ├── Detect clientId/clientSecret → derive Access Token Secret name (same hash)
  ├── Read Access Token Secret from informer cache
  └── Return ConnectionConfig with Credentials.ServiceAccount.BearerToken
     │
     ▼
provider.SdkClientSet()
  └── bearerTokenTransport  →  Authorization: Bearer <token>  on all Atlas API requests
```

---

## Component 1: ServiceAccountToken Controller

**Package:** `internal/controller/serviceaccount/`

**Files:**
- `serviceaccounttoken_controller.go` — main controller
- `token_provider.go` — Atlas SDK wrapper
- `serviceaccounttoken_controller_test.go` — unit tests

### Watches

The controller-runtime cache in `internal/operator/builder.go` is already globally configured to only cache Secrets with the label `atlas.mongodb.com/type: credentials`. The new controller uses a simple `For(&corev1.Secret{})` — no additional label predicate needed. The cache restriction ensures only credential-labelled Secrets are visible.

### Reconcile Logic

```
1. Read the Connection Secret.
2. If Secret does not contain clientId and clientSecret → skip (API key secret; not our concern).
3. Compute the Access Token Secret name:
     tokenSecretName = "atlas-access-token-" + FNVhash(namespace + "/" + connectionSecretName)
4. Attempt to read the Access Token Secret by that name.
5. If the Access Token Secret does not exist:
     a. Call Atlas OAuth endpoint to obtain a token (see Token Provider below).
     b. Create the Access Token Secret with:
          - name: <tokenSecretName>
          - namespace: same as Connection Secret
          - label: atlas.mongodb.com/type: credentials
          - ownerReference → Connection Secret (controller: true, blockOwnerDeletion: true)
          - data.accessToken: <bearer token>
          - data.expiry: <RFC3339 timestamp>
     c. Re-enqueue at 2/3 of token TTL.
6. If the Access Token Secret exists:
     a. Parse expiry. If more than 2/3 of TTL remains → re-enqueue at that remaining time. Done.
     b. Otherwise, call Atlas OAuth endpoint to obtain a new token.
     c. Update Access Token Secret data in place (accessToken + expiry).
     d. Re-enqueue at 2/3 of the new TTL.
```

Minimum re-enqueue interval: 10 seconds (guards against degenerate short-expiry tokens).

### Token Provider (`token_provider.go`)

```go
import "go.mongodb.org/atlas-sdk/v20250312013/auth/clientcredentials"

type AtlasTokenProvider struct {
    AtlasDomain string
}

func (p *AtlasTokenProvider) Token(ctx context.Context, clientID, clientSecret string) (*oauth2.Token, error) {
    cfg := clientcredentials.NewConfig(clientID, clientSecret)
    cfg.TokenURL = p.AtlasDomain + clientcredentials.TokenAPIPath
    return cfg.Token(ctx)
}
```

`AtlasDomain` is received at construction time from the `--atlas-domain` operator flag (same as other controllers receive it via `AtlasProvider`). Overriding `TokenURL` enables QA and Gov environment support.

### RBAC Markers

```go
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
```

These markers are on the controller struct and will be picked up by `make manifests` to regenerate RBAC in `config/rbac/`.

### Rationale for Deterministic Naming vs. Annotation

The PoC branch used `generateName` to create the Access Token Secret and stored the generated name as an annotation (`atlas.mongodb.com/access-token`) on the user-provided Connection Secret. This has a GitOps compatibility problem: ArgoCD and Flux reconcile user-managed Secrets back to their declared state, which would remove the operator-written annotation and break the token lookup.

Instead, the Access Token Secret name is derived deterministically:

```
tokenSecretName = "atlas-access-token-" + base16(FNV-1a-64(namespace + "/" + connectionSecretName))
```

This matches the pattern already used in the codebase (see `pkg/controller/state/tracker.go`). Any component — the `ServiceAccountToken` controller or any resource reconciler — can compute the same name independently, with no shared mutable state and no annotation reads.

---

## Component 2: GetConnectionConfig Changes

**File:** `internal/controller/reconciler/credentials.go`

### Updated Validation

`validate()` is extended to accept either credential type:

- `orgId` is always required.
- Exactly one of the following must be present:
  - `publicApiKey` AND `privateApiKey` (API key path)
  - `clientId` AND `clientSecret` (Service Account path)
- Both present simultaneously → validation error: `"secret must not contain both API key and service account fields"`.
- Neither present → existing validation error.

### Updated GetConnectionConfig

After reading the Connection Secret, detect which credential type is present:

1. **API key path** (`publicApiKey` + `privateApiKey`): unchanged. Returns `ConnectionConfig` with `Credentials.APIKeys` populated.
2. **Service Account path** (`clientId` + `clientSecret`):
   - Compute `tokenSecretName` using the same deterministic hash as the controller.
   - Fetch the Access Token Secret from the informer cache.
   - If not found: return a descriptive error. The reconciler will `Terminate` and re-enqueue on `DefaultRetry` (10 seconds). Meanwhile, the `ServiceAccountToken` controller will create the token on its next reconcile. The credential Secret's `ResourceVersion` does not change when the Access Token Secret is created (different object), so the reconciler relies on the retry timer rather than a watch event for the initial token creation.
   - If found: return `ConnectionConfig` with `Credentials.ServiceAccount.BearerToken` populated.

---

## Component 3: Provider and SDK Client

**File:** `internal/controller/atlas/provider.go`

### New Types

```go
type Credentials struct {
    APIKeys        *APIKeys
    ServiceAccount *ServiceAccountToken
}

type ServiceAccountToken struct {
    BearerToken string
}
```

### SdkClientSet Branching

```go
switch {
case creds.ServiceAccount != nil:
    httpClient = bearerTokenHTTPClient(creds.ServiceAccount.BearerToken)
case creds.APIKeys != nil:
    httpClient = digestHTTPClient(creds.APIKeys.PublicKey, creds.APIKeys.PrivateKey)
default:
    return nil, errors.New("no credentials provided")
}
```

### bearerTokenTransport

A minimal `http.RoundTripper` that injects `Authorization: Bearer <token>` on every request:

```go
type bearerTokenTransport struct {
    token string
    base  http.RoundTripper
}

func (t *bearerTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    req = req.Clone(req.Context())
    req.Header.Set("Authorization", "Bearer "+t.token)
    return t.base.RoundTrip(req)
}
```

This is used with the Atlas SDK via `admin.UseHTTPClient(httpClient)`.

---

## Component 4: Controller Registration

**File:** `internal/controller/registry.go`

```go
reconcilers = append(reconcilers,
    serviceaccounttoken.NewServiceAccountTokenReconciler(c, r.logger, r.atlasDomain, r.maxConcurrentReconciles),
)
```

The `ServiceAccountToken` controller is registered alongside all other reconcilers. It requires the configured `atlasDomain` to override the OAuth token URL for non-production environments.

---

## User-Facing Contract

No CRD changes. No changes to how custom resources reference Connection Secrets. The only user action is to create a Connection Secret with a different set of fields.

### Service Account Connection Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-atlas-service-account
  namespace: atlas-operator
  labels:
    atlas.mongodb.com/type: credentials
type: Opaque
stringData:
  orgId: "64abc123def456"
  clientId: "mdb_sa_id_xxxxxxxxxxxxxxxxxxxxxxxxxx"
  clientSecret: "mdb_sa_sk_xxxxxxxxxxxxxxxxxxxxxxxxxx"
```

This Secret is referenced from custom resources exactly as today:

```yaml
spec:
  connectionSecretRef:
    name: my-atlas-service-account
```

Or, when used as the global operator secret, via `--global-api-secret-name my-atlas-service-account` (unchanged).

### Access Token Secret (operator-managed, not user-created)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: atlas-access-token-a3f7c2d1  # deterministic name, never edited by user
  namespace: atlas-operator
  labels:
    atlas.mongodb.com/type: credentials  # required: must be in informer cache
  ownerReferences:
    - apiVersion: v1
      kind: Secret
      name: my-atlas-service-account
      uid: <uid>
      controller: true
      blockOwnerDeletion: true
type: Opaque
data:
  accessToken: <base64>  # bearer token
  expiry: <base64>       # RFC3339 timestamp
```

Users should not create, modify, or reference this Secret. It is automatically created, updated, and garbage-collected by the operator.

---

## Token Lifecycle

```
t=0                   t=40min              t=1h
|------ valid ---------|--- refresh zone ---|--- expired ---|
                       ↑
                 Controller re-enqueues here (2/3 of TTL)
                 and fetches a new token before the old one expires
```

- Token TTL: 3600 seconds (fixed by Atlas OAuth endpoint).
- Refresh at: 2/3 of TTL = ~2400 seconds (~40 minutes).
- On refresh, the Access Token Secret is **updated in place** — `data.accessToken` and `data.expiry` are overwritten. The Secret name does not change.
- If refresh fails (Atlas unreachable, credentials revoked), the controller re-enqueues with exponential backoff. Resource reconcilers continue using the existing token until it expires, then fail with a 401 and re-enqueue. Clear error conditions and events are emitted (see Testing section).

---

## Watch Infrastructure

No changes to existing watch configurations. The behaviour follows from existing infrastructure:

1. The global cache in `internal/operator/builder.go` restricts cached Secrets to those with label `atlas.mongodb.com/type: credentials`. Both Connection Secrets and Access Token Secrets carry this label, so both are in the informer cache.
2. Existing resource reconcilers already watch `&corev1.Secret{}` with `ResourceVersionChangedPredicate`. When the `ServiceAccountToken` controller updates the Access Token Secret (on refresh), the `ResourceVersion` changes → all reconcilers referencing the corresponding Connection Secret are re-triggered via the existing map function. **No changes to existing watch configurations or predicates.**
3. The `ServiceAccountToken` controller uses `For(&corev1.Secret{})`. No own predicate needed — the global cache restriction already filters to credential Secrets only.

---

## Edge Cases

| Scenario                                                   | Behavior                                                                                                                                                                                                         |
|------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Connection Secret deleted                                  | `ownerReference` causes Access Token Secret to be garbage-collected automatically.                                                                                                                               |
| Token refresh fails (network, Atlas outage)                | Controller errors and re-enqueues with backoff. Reconcilers use the existing valid token until it expires. On expiry, Atlas returns 401 and reconcilers re-enqueue. User sees a condition/event on the resource. |
| Token expires mid-reconcile                                | Atlas returns 401. Reconciler terminates and re-enqueues at `DefaultRetry`. Token should already be refreshed (controller refreshes at 40 min; tokens live 60 min).                                              |
| Access Token Secret not yet created (first reconcile race) | `GetConnectionConfig` returns `"access token secret not found"` error. Reconciler re-enqueues at 10s. Controller creates the token on its own reconcile. Reconciler retries and finds it.                        |
| Access Token Secret unexpectedly deleted                   | Controller detects missing Secret on next reconcile and recreates it. Reconcilers get errors until recreation (seconds).                                                                                         |
| Secret has both API key and SA fields                      | `validate()` returns a clear error: rejected as ambiguous. Resource condition is set to failed.                                                                                                                  |
| Credentials revoked in Atlas                               | Token request returns error. Controller emits event, re-enqueues. Reconcilers see errors. User must update the Connection Secret.                                                                                |
| Multiple resources sharing one Connection Secret           | All share the same Access Token Secret (same deterministic name). Token is refreshed once, used by all.                                                                                                          |

---

## Installation

### Helm Chart Changes

**File:** `charts/atlas-operator/values.yaml`

A new `credentials` section is added under the existing `atlas` block to document the Service Account option. The chart already generates a Secret from these values — the only change is adding `clientId` and `clientSecret` as optional alternatives to `apiKey.publicKey` / `apiKey.privateKey`.

```yaml
atlas:
  # Option 1 (legacy): Programmatic API Key
  apiKey:
    orgId: ""
    publicKey: ""
    privateKey: ""

  # Option 2 (recommended): Service Account
  serviceAccount:
    orgId: ""
    clientId: ""
    clientSecret: ""
```

The chart template generates the Connection Secret from whichever option is populated. The generated Secret carries the `atlas.mongodb.com/type: credentials` label as today.

Helm chart validation (`_helpers.tpl` or a named template) ensures exactly one of `apiKey` or `serviceAccount` is configured, not both.

**Quick-start guide updates:** The default `values.yaml` comment block is updated to explain both options, with a note that Service Account is recommended.

### OLM / CSV Changes

The `ServiceAccountToken` controller requires the following RBAC on Secrets (in addition to what already exists):

```
get, list, watch, create, update, patch
```

The existing Secret permissions in the CSV cover `get`, `list`, `watch`. The new `create`, `update`, `patch` verbs are added.

**How these get applied:** The `// +kubebuilder:rbac:...` marker on the controller struct is picked up by `make manifests`, which regenerates `config/rbac/role.yaml`. The updated role is then included in the next `make bundle` run, which regenerates the CSV. No manual CSV edits required.

**OLM-specific note:** The CSV's `installModes` and `WATCH_NAMESPACE` injection are unaffected. No new env vars are introduced.

---

## Documentation Plan

The public documentation at `https://www.mongodb.com/docs/atlas/operator/current/` has a top-level **Atlas Access** section. The current page **Configure Access to Atlas** covers API key setup. The following documentation changes are required:

### New Page: Configure Access to Atlas Using Service Accounts

**Location:** Under **Atlas Access**, as a new sibling page to the existing "Configure Access to Atlas" (API keys) page.

**Suggested title:** `Configure Access to Atlas Using Service Accounts`

**Content outline:**
1. Overview: what Atlas API Service Accounts are; why they are recommended over API keys.
2. Prerequisites: how to create a Service Account in Atlas (link to Atlas docs); required roles (Organisation Owner for org-level SA).
3. Step-by-step guide:
   - Create the Service Account in Atlas (UI or API), copy `clientId` and `clientSecret`.
   - Create the Kubernetes Secret with `orgId`, `clientId`, `clientSecret`.
   - Reference the Secret from `AtlasProject` or use it as the global operator Secret.
   - Verify the operator is using the Service Account (check resource conditions).
4. YAML example (Connection Secret).
5. Note on token management: tokens are managed automatically by the operator; do not edit the Access Token Secret.
6. Note on credential rotation: when the Atlas Service Account secret is rotated, update the Kubernetes Connection Secret; the operator will re-acquire a new token on the next reconcile.
7. Note on IP access list: the Service Account's API access list must include the operator pod's egress IP(s), same as for API keys.

### Updates to Existing Pages

| Page                                 | Change                                                                                                                                |
|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| Configure Access to Atlas (API keys) | Add a note at the top: "If you prefer to use Service Accounts (recommended), see [Configure Access to Atlas Using Service Accounts]." |
| Quick Start                          | Add Service Account as an alternative authentication option alongside API key setup.                                                  |
| Helm chart installation page         | Document the new `atlas.serviceAccount.*` values.                                                                                     |
| Changelog                            | Entry for the release that includes this feature.                                                                                     |

---

## Testing

### Unit Tests

New unit tests in the new and modified packages:

**`internal/controller/serviceaccount/`**
- Skip non-SA Connection Secrets (API key secrets).
- Create Access Token Secret on first reconcile.
- Skip refresh when token has more than 2/3 TTL remaining.
- Refresh token when 2/3 TTL has elapsed (updates Secret in place).
- Handle Atlas token endpoint error (returns error, does not crash).
- Handle missing Access Token Secret (detected, recreated).
- Handle deleted Connection Secret (ownerReference GC — validated by checking ownerReference is set).
- Correct deterministic Secret name computation.
- Correct `TokenURL` override for QA and Gov domains.

**`internal/controller/reconciler/credentials_test.go`**
- SA Secret with no Access Token Secret yet → returns descriptive error.
- SA Secret with valid Access Token Secret → returns `ConnectionConfig` with `BearerToken`.
- SA Secret with both API key and SA fields → validation error.
- SA Secret with `orgId` missing → validation error.

**`internal/controller/atlas/provider_test.go`**
- `SdkClientSet` with `ServiceAccount` credentials uses `bearerTokenTransport`.
- `bearerTokenTransport` injects correct `Authorization: Bearer` header.
- `SdkClientSet` with `APIKeys` credentials uses digest transport (regression).
- `SdkClientSet` with nil credentials returns error.

### E2E Tests

New e2e test file: `test/e2e2/serviceaccount_test.go`

Scenarios (using the `test/e2e2/` framework):

1. **Token creation**: Create Atlas SA via Admin API → create Connection Secret → verify Access Token Secret is created with correct fields, label, and ownerReference within timeout.
2. **API key passthrough**: Create API key Connection Secret → verify no Access Token Secret is created (no-op for the SA controller).
3. **Full project lifecycle**: Create SA → create Connection Secret → create `AtlasProject` → verify project reaches `Ready` condition → delete.
4. **Full deployment lifecycle**: Create SA → project → `AtlasDeployment` (Flex) → verify both reach `Ready` → clean up.

---

## Open Questions

| # | Question                                                                                                                                                                                                                                                                                                                                                                | Owner            | Status                                           |
|---|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------|--------------------------------------------------|
| 1 | **Service Account credential expiry notification**: The `clientSecret` stored in the Connection Secret has a configurable expiry (up to 365 days). Should the operator emit a Kubernetes Event or set a Warning condition when the secret is approaching expiry? What is the threshold (e.g., 30 days)? How does this align with Atlas's own near-expiry notifications? | Engineering + PM | Open — investigate in TD review                  |
| 2 | **Secret rotation flow**: When a user rotates their Atlas Service Account secret (creates a new secret in Atlas and updates the Kubernetes Connection Secret), what is the expected operator behaviour? The current design re-acquires a token on the next reconcile. Is zero-downtime rotation a requirement for Phase 1?                                              | PM + Engineering | Open                                             |
| 3 | **Deprecation of Programmatic API Keys**: Should the docs and Helm chart explicitly discourage new API key usage? What is MongoDB's current public deprecation stance for API keys across Atlas tooling?                                                                                                                                                                | PM               | Open — check Atlas API + Terraform Provider docs |
| 4 | **Multiple SA secrets**: Atlas Service Accounts can have multiple active secrets. The current design uses a single `clientId` + `clientSecret` pair. Should the operator support multiple secrets for rotation overlap (zero-downtime rotation)?                                                                                                                        | Engineering      | Deferred to Phase 2 / rotation investigation     |

---

## Files Changed

| File                                                                   | Change                                                                                                                               |
|------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| `internal/controller/serviceaccount/serviceaccounttoken_controller.go`      | **New.** Main `ServiceAccountToken` controller. Watches credential Secrets, manages Access Token Secret lifecycle.                   |
| `internal/controller/serviceaccount/token_provider.go`                      | **New.** Atlas SDK `clientcredentials` wrapper with `TokenURL` override for QA/Gov.                                                  |
| `internal/controller/serviceaccount/serviceaccounttoken_controller_test.go` | **New.** Unit tests for `ServiceAccountToken` controller.                                                                            |
| `internal/controller/atlas/provider.go`                                     | **Modified.** Add `ServiceAccountToken` type to `Credentials`. Branch `SdkClientSet` on credential type. Add `bearerTokenTransport`. |
| `internal/controller/atlas/provider_test.go`                                | **Modified.** Tests for bearer transport and new credential branching.                                                               |
| `internal/controller/reconciler/credentials.go`                             | **Modified.** Detect `clientId`/`clientSecret` in `GetConnectionConfig`. Derive Access Token Secret name. Update `validate()`.       |
| `internal/controller/reconciler/credentials_test.go`                        | **Modified.** Tests for SA credential path.                                                                                          |
| `internal/controller/registry.go`                                           | **Modified.** Register `ServiceAccountTokenReconciler`.                                                                              |
| `config/rbac/role.yaml`                                                | **Regenerated** by `make manifests` — adds `create;update;patch` to Secret permissions.                                              |
| `charts/atlas-operator/values.yaml`                                    | **Modified.** Add `atlas.serviceAccount.{orgId,clientId,clientSecret}` values.                                                       |
| `charts/atlas-operator/templates/secret.yaml`                          | **Modified.** Generate Connection Secret from either API key or Service Account values. Add chart-level validation.                  |
| `test/e2e2/serviceaccount_test.go`                                     | **New.** E2E test scenarios.                                                                                                         |
| `docs/dev/td-service-accounts-phase1.md`                               | **New.** This document.                                                                                                              |

> **Note:** `service-account-support.md` (checked into the spike branch `CLOUDP-373260/spike-branch`) is superseded by this TD. The PoC in that branch uses annotation-based token naming (Option A); the implementation must be updated to use deterministic hashing (Option B) before the branch is merged.

---

## References

- [Spike: Service Accounts support in AKO](https://docs.google.com/document/d/1FwV7GfC9i__dRGP6m3cExj7t3hc0wsKuB0dWDwYK0Mc)
- [PD+Scope: Service Accounts support — Phase 1](https://docs.google.com/document/d/1Kn8bwUWXx4Piojzd2eVh7QXgpymPciivi1h78pTSa7I)
- [Atlas API Authentication (public docs)](https://www.mongodb.com/docs/atlas/api/api-authentication/)
- [Atlas Service Accounts Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/)
- [Generate OAuth2 Token](https://www.mongodb.com/docs/atlas/api/service-accounts/generate-oauth2-token/)
- [Atlas Go SDK — clientcredentials package](https://pkg.go.dev/go.mongodb.org/atlas-sdk/v20250312013/auth/clientcredentials)
- Terraform Provider reference: [Service Account (recommended)](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/provider-configuration#service-account-recommended)
- Jira: [CLOUDP-373260](https://jira.mongodb.org/browse/CLOUDP-373260) (Spike), [CLOUDP-383935](https://jira.mongodb.org/browse/CLOUDP-383935) (Epic)
