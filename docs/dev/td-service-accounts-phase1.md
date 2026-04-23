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

- **Service Account (SA)** — An Atlas OAuth 2.0 client (client ID + client secret) scoped to one organisation. Used to obtain short-lived access tokens for the Atlas Administration API.
- **Programmatic API Key (PAK)** — Legacy Atlas authentication: public key + private key, HTTP Digest. Still supported; not deprecated by this project.
- **Connection Secret** — The Kubernetes Secret referenced by an Atlas custom resource (or the global operator flag) that holds Atlas credentials. Today: `orgId + publicApiKey + privateApiKey`. After this project: also accepts `orgId + clientId + clientSecret`.
- **Access Token Secret** — A Kubernetes Secret managed by the `ServiceAccountToken` controller. Holds the short-lived OAuth bearer token and its expiry timestamp. Never created or edited by users.
- **Bearer Token** — A short-lived OAuth 2.0 access token (1 hour / 3600 seconds) returned by the Atlas token endpoint. Used in `Authorization: Bearer <token>` on all Atlas API requests.
- **`ServiceAccountToken` controller** — New controller introduced in this project. Watches Connection Secrets, detects Service Account credentials, and manages the Access Token Secret lifecycle (create, refresh).

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
ServiceAccountToken Controller
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

## Shared Package: `accesstoken`

**Package:** `internal/controller/accesstoken/`

A small, dependency-free helper package used by both the `ServiceAccountToken` controller (producer) and `reconciler.GetConnectionConfig` (consumer). It holds the schema of the Access Token Secret — exported constants `AccessTokenKey`, `ExpiryKey`, `CredentialsHashKey` for the data-field names, and the two helpers operating on the shape:

- `DeriveSecretName(namespace, connectionSecretName string) (string, error)` — the deterministic `"atlas-access-token-<name>-<hash>"` derivation.
- `CredentialsHash(clientID, clientSecret string) (string, error)` — the FNV-1a-64 fingerprint used for rotation detection.

The package imports only `fmt`, `hash/fnv`, and `k8s.io/apimachinery/pkg/util/rand`. The Connection Secret field keys (`ClientIDKey`, `ClientSecretKey`) remain in the `reconciler` package because they describe a different Secret; both packages re-use them through the exported names.

---

## Component 1: ServiceAccountToken Controller

**Package:** `internal/controller/serviceaccounttoken/`

**Files:**
- `serviceaccounttoken_controller.go` — main controller. Uses `accesstoken.*` for the Access Token Secret schema and `reconciler.ClientIDKey` / `reconciler.ClientSecretKey` for the Connection Secret fields.
- `token_provider.go` — Atlas SDK wrapper.
- `serviceaccounttoken_controller_test.go` — unit tests.

### Watches

The controller watches `&corev1.Secret{}` through two registrations:

1. **`For(&corev1.Secret{}, builder.WithPredicates(credentialsLabelPredicate(), ResourceVersionChangedPredicate{}))`** — the primary watch. Filters to Secrets carrying `atlas.mongodb.com/type: credentials` and enqueues the event's own `NamespacedName`. This covers creation, data updates, and rotation of Connection Secrets.
2. **`Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(MapAccessTokenSecretToOwner), builder.WithPredicates(credentialsLabelPredicate()))`** — the secondary watch. On any event for a credentials-labelled Secret, the map function inspects `ownerReferences` and, if the Secret has a `Secret` owner, enqueues the owner's `NamespacedName`. For user-created Connection Secrets (which have no owner) it returns nothing, avoiding duplicate enqueues. For Access Token Secrets (owned by a Connection Secret) it causes the controller to re-reconcile the owner — notably recreating the token on accidental deletion.

The controller-runtime cache in `internal/operator/builder.go` additionally restricts the whole informer to credentials-labelled Secrets when the operator runs cluster-wide. The label predicate on both watches duplicates that filter so the controller behaves the same way in namespaced mode, where the cache is not label-restricted.

### Reconcile Logic

```
1. Read the Connection Secret.
2. If Secret does not contain clientId and clientSecret → skip (API key secret; not our concern).
3. Compute the Access Token Secret name:
     tokenSecretName = "atlas-access-token-" + connectionSecretName + "-" + hash(namespace + "/" + connectionSecretName)
   (The literal connectionSecretName is truncated if the total would exceed the 253-character Kubernetes DNS-subdomain limit.)
4. Compute the credentials fingerprint:
     currentHash = FNV-1a-64(clientId + "\x00" + clientSecret)
   This is stored on the Access Token Secret as data.credentialsHash and used to detect credential rotation.
5. Attempt to read the Access Token Secret by the derived name.
6. If the Access Token Secret does not exist:
     a. Call Atlas OAuth endpoint to obtain a token (see Token Provider below).
     b. Create the Access Token Secret with:
          - name: <tokenSecretName>
          - namespace: same as Connection Secret
          - label: atlas.mongodb.com/type: credentials
          - ownerReference → Connection Secret (controller: true, blockOwnerDeletion: true)
          - data.accessToken: <bearer token>
          - data.expiry: <RFC3339 timestamp>
          - data.credentialsHash: <currentHash>
     c. On AlreadyExists (concurrent event race): log and requeue at minRequeue; the next reconcile enters the "exists" branch below.
     d. Re-enqueue at 2/3 of token TTL.
7. If the Access Token Secret exists:
     a. Compare stored data.credentialsHash to currentHash. If they differ → credentials have been rotated; immediately refresh the token (call Atlas OAuth endpoint, overwrite data.accessToken + data.expiry + data.credentialsHash in place). Re-enqueue at 2/3 of the new TTL. Done.
     b. Parse data.expiry. If more than 2/3 of TTL remains → re-enqueue at that remaining time. Done.
     c. Otherwise, call Atlas OAuth endpoint to obtain a new token.
     d. Update Access Token Secret data in place (accessToken + expiry + credentialsHash).
     e. Re-enqueue at 2/3 of the new TTL.
```

Minimum re-enqueue interval: 10 seconds (guards against degenerate short-expiry tokens).

### Token Provider (`token_provider.go`)

```go
import "go.mongodb.org/atlas-sdk/v20250312018/auth/clientcredentials"

// TokenProvider abstracts the OAuth token acquisition so it can be mocked in tests.
type TokenProvider interface {
    FetchToken(ctx context.Context, clientID, clientSecret string) (token string, expiry time.Time, err error)
}

type AtlasTokenProvider struct {
    atlasDomain string
}

func (p *AtlasTokenProvider) FetchToken(ctx context.Context, clientID, clientSecret string) (string, time.Time, error) {
    cfg := clientcredentials.NewConfig(clientID, clientSecret)
    cfg.TokenURL = p.atlasDomain + clientcredentials.TokenAPIPath
    token, err := cfg.Token(ctx)
    if err != nil {
        return "", time.Time{}, fmt.Errorf("failed to acquire OAuth token: %w", err)
    }
    return token.AccessToken, token.Expiry, nil
}
```

`atlasDomain` is received at construction time from the `--atlas-domain` operator flag (same as other controllers receive it via `AtlasProvider`). Overriding `TokenURL` enables QA and Gov environment support. The interface returns `(string, time.Time, error)` — a flat tuple rather than an `oauth2.Token` — so tests can substitute a fake without depending on the SDK's types.

### RBAC Markers

```go
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
```

These markers are on the controller struct and will be picked up by `make manifests` to regenerate RBAC in `config/rbac/`.

### Rationale for Deterministic Naming vs. Annotation

The PoC branch used `generateName` to create the Access Token Secret and stored the generated name as an annotation (`atlas.mongodb.com/access-token`) on the user-provided Connection Secret. This has a GitOps compatibility problem: ArgoCD and Flux reconcile user-managed Secrets back to their declared state, which would remove the operator-written annotation and break the token lookup.

Instead, the Access Token Secret name is derived deterministically in `accesstoken.DeriveSecretName`:

```
tokenSecretName = "atlas-access-token-" + connectionSecretName + "-" + rand.SafeEncodeString(fnv64a(namespace + "/" + connectionSecretName))
```

This matches the pattern already used in the codebase (see `pkg/controller/state/tracker.go`). Any component — the `ServiceAccountToken` controller or any resource reconciler — can compute the same name independently, with no shared mutable state and no annotation reads.

### Rationale for Credentials Hash Staleness Detection

Kubernetes Secrets can be rotated in place — a user updates the `clientId` / `clientSecret` fields on the Connection Secret while keeping the same resource name. Without an explicit check, the controller would keep using the cached bearer token until the previous token's natural expiry (up to one hour), even though the new credentials could have already invalidated the old ones on Atlas's side.

To detect rotation the controller writes a non-cryptographic FNV-1a-64 fingerprint of `(clientId, clientSecret)`, computed by `accesstoken.CredentialsHash`, into `data.credentialsHash` on the Access Token Secret. On every reconcile the controller computes the same fingerprint from the current Connection Secret and compares it to the stored value. Any mismatch forces an immediate refresh — new token is fetched, hash is updated. A nul (`\x00`) separator in the hash input disambiguates `("ab", "c")` from `("a", "bc")`.

Content hash is preferred over `ResourceVersion`: the latter changes for any update (new labels, annotations, unrelated data), triggering unnecessary token fetches. A hash of the credential material alone changes exactly when rotation occurs.

---

## Component 2: GetConnectionConfig Changes

**File:** `internal/controller/reconciler/credentials.go`

### Updated Validation

`validateConnectionSecret(secret *corev1.Secret) error` inspects the raw Secret `data` map and returns `nil` when valid. Rules:

- `orgId` is always required.
- Exactly one of the following must be present:
  - `publicApiKey` AND `privateApiKey` (API key path)
  - `clientId` AND `clientSecret` (Service Account path)
- Both credential types present simultaneously → error: `"secret contains both API key and service account credentials; only one type is allowed"`.
- Neither type present, or partial (only one half of a pair) → error: `"missing required fields: [...]"` listing the missing keys.

`GetConnectionConfig` wraps the validation error with the secretRef: `"invalid connection secret <ref>: <validation error>"`.

### Updated GetConnectionConfig

After reading the Connection Secret, detect which credential type is present:

1. **API key path** (`publicApiKey` + `privateApiKey`): unchanged. Returns `ConnectionConfig` with `Credentials.APIKeys` populated.
2. **Service Account path** (`clientId` + `clientSecret`):
   - Compute `tokenSecretName` using `accesstoken.DeriveSecretName`.
   - Fetch the Access Token Secret from the informer cache.
   - If not found: return a descriptive error. The reconciler will `Terminate` and re-enqueue on `DefaultRetry` (10 seconds). Meanwhile, the `ServiceAccountToken` controller will create the token on its next reconcile. The credential Secret's `ResourceVersion` does not change when the Access Token Secret is created (different object), so the reconciler relies on the retry timer rather than a watch event for the initial token creation.
   - Compare the current credentials' `accesstoken.CredentialsHash` against the `data.credentialsHash` on the token Secret. On mismatch return `"access token secret <ref> is stale (credentials rotated); waiting for the service-account-token controller to refresh"`; the downstream reconciler retries until the controller refreshes and the hash matches again.
   - If found and hash matches: return `ConnectionConfig` with `Credentials.ServiceAccount.BearerToken` populated.

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
var baseTransport http.RoundTripper
switch {
case creds.ServiceAccount != nil:
    baseTransport = &oauth2.Transport{
        Source: oauth2.StaticTokenSource(&oauth2.Token{
            AccessToken: creds.ServiceAccount.BearerToken,
        }),
    }
case creds.APIKeys != nil:
    baseTransport = digest.NewTransport(creds.APIKeys.PublicKey, creds.APIKeys.PrivateKey)
default:
    return nil, fmt.Errorf("no credentials provided")
}
```

### Bearer Token Transport

Rather than a hand-rolled `http.RoundTripper`, the Service Account path uses `golang.org/x/oauth2`'s `oauth2.Transport` with a `StaticTokenSource`. `oauth2.Transport` injects `Authorization: Bearer <token>` on every request and delegates to the default transport. Using the upstream package keeps the surface area small and avoids re-implementing header manipulation and cloning semantics.

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
  name: atlas-access-token-my-atlas-service-account-5557f49b459c76694d  # deterministic name, never edited by user
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
  accessToken: <base64>     # bearer token
  expiry: <base64>          # RFC3339 timestamp
  credentialsHash: <base64> # FNV fingerprint of (clientId, clientSecret) used to detect rotation
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

Two controller layers watch Secrets: the `ServiceAccountToken` controller (producer of the Access Token Secret) and each downstream resource reconciler (consumer of the Connection Secret). The informer cache is shared across all controllers.

### Global informer cache

- In **cluster-wide mode** (operator watches all namespaces) the cache in `internal/operator/builder.go` restricts cached Secrets to those carrying `atlas.mongodb.com/type: credentials`. Both Connection Secrets and Access Token Secrets carry that label, so both sit in the cache.
- In **namespaced mode** the cache is not label-restricted; each controller applies the same filter via its own predicate so behaviour is identical in both modes.

### ServiceAccountToken controller — primary watch

`For(&corev1.Secret{}, builder.WithPredicates(credentialsLabelPredicate, ResourceVersionChangedPredicate))` enqueues the event's own `NamespacedName`. Triggered by:

- creation, data updates, and rotation of Connection Secrets — the Reconcile then creates or refreshes the owned Access Token Secret;
- updates to Access Token Secrets themselves (for example the controller's own refresh write) — these early-exit in Reconcile because the object has no `clientId` / `clientSecret`.

### ServiceAccountToken controller — Access Token Secret owner back-channel

`Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(MapAccessTokenSecretToOwner), builder.WithPredicates(credentialsLabelPredicate))` inspects `ownerReferences` on any credentials-labelled Secret event. When a Secret owner is present (i.e. the event came from an Access Token Secret owned by a Connection Secret) the map function enqueues a Reconcile for the **owner**. For user-authored Connection Secrets, which have no owner, the map returns nothing and the primary watch remains the only enqueue path.

The net effect:

- Connection Secret created / updated / rotated → primary watch → Reconcile the Connection Secret directly.
- Access Token Secret deleted out-of-band → owner back-channel → Reconcile the owning Connection Secret → `createToken` recreates it at the derived name within seconds.
- Access Token Secret updated (the controller's own refresh) → owner back-channel → Reconcile the owning Connection Secret → matching hash, fresh expiry → requeue, no token fetch.

### Resource reconcilers (downstream)

Unchanged by this feature. Each resource reconciler (AtlasProject, AtlasDeployment, AtlasDatabaseUser, …) has its own `Watches(&corev1.Secret{}, …)` with a map function backed by a field index (e.g. `AtlasProjectBySecretsIndex`) that matches Secrets referenced as `connectionSecretRef` on the custom resource. The Access Token Secret's derived name is never referenced by a user's custom resource, so downstream reconcilers are not re-triggered by Access Token Secret refreshes — they pick up the freshly-minted bearer on their own next reconcile.

---

## Edge Cases

- **Connection Secret deleted** — `ownerReference` causes the Access Token Secret to be garbage-collected automatically.
- **Token refresh fails (network, Atlas outage)** — Controller errors and re-enqueues with backoff. Reconcilers use the existing valid token until it expires. On expiry, Atlas returns 401 and reconcilers re-enqueue. User sees a condition/event on the resource.
- **Token expires mid-reconcile** — Atlas returns 401. Reconciler terminates and re-enqueues at `DefaultRetry`. Token should already be refreshed (controller refreshes at 40 min; tokens live 60 min).
- **Access Token Secret not yet created (first reconcile race)** — `GetConnectionConfig` returns `"access token secret <name> does not exist yet"` error. Reconciler re-enqueues at 10s. Controller creates the token on its own reconcile. Reconciler retries and finds it.
- **Access Token Secret stale after in-place credential rotation** — `GetConnectionConfig` compares the current credentials' FNV hash to `data.credentialsHash` on the token Secret. On mismatch it returns `"access token secret <name> is stale (credentials rotated); waiting for the service-account-token controller to refresh"`. The downstream reconciler re-enqueues; the service-account-token controller refreshes on its own reconcile; the next reconcile sees a matching hash and proceeds.
- **Access Token Secret unexpectedly deleted** — The controller's `Watches` on credentials-labelled Secrets runs `MapAccessTokenSecretToOwner` on the delete event, which reads the Secret's `ownerReferences` and enqueues a Reconcile for the owning Connection Secret. Reconcile sees no existing token at the derived name and calls `createToken`. Recreation happens within seconds; downstream reconcilers see a transient "does not exist yet" error in the interim.
- **Secret has both API key and SA fields** — `validateConnectionSecret` returns an error: rejected as ambiguous. Resource condition is set to failed.
- **Credentials rotated in place on the Connection Secret** — Controller detects `data.credentialsHash` mismatch on next reconcile and refreshes the token immediately, regardless of the cached token's remaining TTL. The Access Token Secret is updated in place.
- **Credentials revoked in Atlas** — Token request returns error. Controller emits event, re-enqueues. Reconcilers see errors. User must update the Connection Secret.
- **Multiple resources sharing one Connection Secret** — All share the same Access Token Secret (same deterministic name). Token is refreshed once, used by all.

---

## Known Limitations

- **First-reconcile and post-rotation transient errors.** Downstream reconcilers (AtlasProject, AtlasDeployment, etc.) watch the user-provided Connection Secret, not the derived Access Token Secret. When a Service Account Connection Secret is first created, and again briefly after in-place credential rotation, the downstream reconciler may run before the service-account-token controller has created or refreshed the Access Token Secret. `GetConnectionConfig` returns a specific error (`"does not exist yet"` or `"is stale (credentials rotated)"`) and the downstream reconciler relies on its own retry timer to recover. The window is bounded by the reconciler's retry backoff — typically seconds — and resolves as soon as the service-account-token controller catches up. No persistent failure results. Adding cross-watches from every downstream reconciler to the derived Access Token Secret name would eliminate the transient error at the cost of broader changes to each controller's `SetupWithManager`; this is deferred pending operational data on how noticeable the transient errors are in practice.

---

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
message).

Full install gates: `make bundle-validate` (OLM schema), `make
validate-manifests` (regen drift), `make validate-crds-chart` (chart CRD
parity with the bundle), `make unit-test` (chart rendering).

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

- **Configure Access to Atlas (API keys)** — Add a note at the top: "If you prefer to use Service Accounts (recommended), see [Configure Access to Atlas Using Service Accounts]."
- **Quick Start** — Add Service Account as an alternative authentication option alongside API key setup.
- **Helm chart installation page** — Document the new `atlas.serviceAccount.*` values.
- **Changelog** — Entry for the release that includes this feature.

---

## Testing

### Unit Tests

New unit tests in the new and modified packages:

**`internal/controller/serviceaccounttoken/`**
- Skip non-SA Connection Secrets (API key secrets).
- Create Access Token Secret on first reconcile (includes `credentialsHash` in data).
- Skip refresh when token has more than 2/3 TTL remaining AND `credentialsHash` matches.
- Refresh token when 2/3 TTL has elapsed (updates Secret in place).
- Refresh token when `credentialsHash` on the Access Token Secret differs from the current Connection Secret credentials, regardless of expiry.
- Handle Atlas token endpoint error (returns error, does not crash).
- Handle missing Access Token Secret (detected, recreated).
- Handle `AlreadyExists` on create (concurrent event race) — fall through to refresh path on the next reconcile.
- Idempotent duplicate reconcile event — second call does not double-fetch the token.
- Correct deterministic Secret name computation.
- Credential Secret is never mutated by the controller (no annotations written).

**`internal/controller/accesstoken/accesstoken_test.go`**
- `DeriveSecretName` pins a known input to its expected output (compatibility contract), plus namespace- and name-sensitivity, far-past-limit length bound, and exact-253-character at the truncation boundary.
- `CredentialsHash` pins a known input to its expected output, distinguishes distinct credential pairs, and disambiguates `("ab","c")` from `("a","bc")` via the nul separator.

**`internal/controller/reconciler/credentials_test.go`**
- SA Secret with no Access Token Secret yet → `"access token secret <ref> does not exist yet"`.
- SA Secret with valid Access Token Secret → returns `ConnectionConfig` with `BearerToken`.
- SA Secret with stale `data.credentialsHash` → `"is stale (credentials rotated); waiting for the service-account-token controller to refresh"`.
- `validateConnectionSecret`: orgId-only / partial API keys / partial SA / both pairs present → each returns a specific error.
- `validateConnectionSecret`: complete API keys or complete SA credentials with orgId → nil.

**`internal/controller/atlas/provider_test.go`**
- `SdkClientSet` with `ServiceAccount` credentials wraps the Atlas SDK client with `oauth2.Transport`.
- `SdkClientSet` with `APIKeys` credentials uses digest transport (regression).
- `SdkClientSet` with nil credentials returns error.

### E2E Tests

New e2e test file: `test/e2e2/serviceaccount_test.go`

Scenarios (using the `test/e2e2/` framework):

1. **Token creation**: Create Atlas SA via Admin API → create Connection Secret → verify Access Token Secret is created with correct fields, label, and ownerReference within timeout.
2. **API key passthrough**: Create API key Connection Secret → verify no Access Token Secret is created (no-op for the SA controller).
3. **Credential rotation**: Create two Atlas SAs (pre- and post-rotation) → create Connection Secret with the first → wait for minted token → rotate `clientId`/`clientSecret` on the Connection Secret in place → verify `accessToken` and `credentialsHash` on the token Secret are replaced within a bounded time.
4. **Accidental deletion recovery**: Create SA → create Connection Secret → wait for minted token → delete the Access Token Secret directly → verify it is recreated with a fresh UID, non-empty bearer token, correct label, and ownerReference within seconds.
5. **Full project lifecycle**: Create SA → create Connection Secret → create `AtlasProject` → verify project reaches `Ready` condition → delete.
6. **Full deployment lifecycle**: Create SA → project → `AtlasDeployment` (Flex) → verify both reach `Ready` → clean up.

---

## Open Questions

- **Service Account credential expiry notification** (Owner: Engineering + PM; Status: Open — investigate in TD review). The `clientSecret` stored in the Connection Secret has a configurable expiry (up to 365 days). Should the operator emit a Kubernetes Event or set a Warning condition when the secret is approaching expiry? What is the threshold (e.g., 30 days)? How does this align with Atlas's own near-expiry notifications?
- **Secret rotation flow** (Owner: PM + Engineering; Status: Open). When a user rotates their Atlas Service Account secret (creates a new secret in Atlas and updates the Kubernetes Connection Secret), what is the expected operator behaviour? The current design re-acquires a token on the next reconcile. Is zero-downtime rotation a requirement for Phase 1?
- **Deprecation of Programmatic API Keys** (Owner: PM; Status: Open — check Atlas API + Terraform Provider docs). Should the docs and Helm chart explicitly discourage new API key usage? What is MongoDB's current public deprecation stance for API keys across Atlas tooling?
- **Multiple SA secrets** (Owner: Engineering; Status: Deferred to Phase 2 / rotation investigation). Atlas Service Accounts can have multiple active secrets. The current design uses a single `clientId` + `clientSecret` pair. Should the operator support multiple secrets for rotation overlap (zero-downtime rotation)?

---

## Files Changed

- `internal/controller/accesstoken/accesstoken.go` — **New.** Shared schema (constants `AccessTokenKey`, `ExpiryKey`, `CredentialsHashKey`) and helpers (`DeriveSecretName`, `CredentialsHash`) for the Access Token Secret.
- `internal/controller/accesstoken/accesstoken_test.go` — **New.** Pinned-output and behavioural tests for `DeriveSecretName` and `CredentialsHash`.
- `internal/controller/serviceaccounttoken/serviceaccounttoken_controller.go` — **New.** Main `ServiceAccountToken` controller. Watches credential Secrets, manages Access Token Secret lifecycle. Consumes `accesstoken.*` for the Access Token Secret schema and `reconciler.ClientIDKey`/`ClientSecretKey` for the Connection Secret fields. Includes `MapAccessTokenSecretToOwner` so an accidentally deleted Access Token Secret is recreated on the next Reconcile of the owning Connection Secret.
- `internal/controller/serviceaccounttoken/token_provider.go` — **New.** Atlas SDK `clientcredentials` wrapper with `TokenURL` override for QA/Gov.
- `internal/controller/serviceaccounttoken/serviceaccounttoken_controller_test.go` — **New.** Unit tests for `ServiceAccountToken` controller, including map-to-owner tests.
- `internal/controller/atlas/provider.go` — **Modified.** Add `ServiceAccountToken` type to `Credentials`. Branch `SdkClientSet` on credential type; use `oauth2.Transport` for bearer auth.
- `internal/controller/atlas/provider_test.go` — **Modified.** Tests for bearer transport and new credential branching.
- `internal/controller/reconciler/credentials.go` — **Modified.** Detect `clientId`/`clientSecret` in `GetConnectionConfig`. Export `ClientIDKey` / `ClientSecretKey` for cross-package use. Add `validateConnectionSecret` helper and the stale-hash guard in `getServiceAccountAccessToken` (both via `accesstoken.*`).
- `internal/controller/reconciler/credentials_test.go` — **Modified.** Tests for SA credential path, including the stale-hash branch.
- `internal/controller/registry.go` — **Modified.** Register `ServiceAccountTokenReconciler`.
- `config/rbac/role.yaml` — **Regenerated** by `make manifests` — adds `create;update;patch` to Secret permissions.
- `helm-charts/atlas-operator/` — **Modified.** Add `globalConnectionSecret.{clientId,clientSecret}` values and `secret.yaml` branching to emit either API-key or Service Account Secret; README documents both install variants.
- `helm-charts/atlas-advanced/` — **Modified.** Add `secret.{clientId,clientSecret}` values and templates/secret.yaml branching.
- `helm-charts/atlas-basic/` — **Modified.** Add `secret.{clientId,clientSecret}` values and templates/secret.yaml branching.
- `helm-charts/atlas-deployment/` — **Modified.** Add `atlas.secret.{clientId,clientSecret}` values and templates/secret.yaml branching; README documents both install variants.
- `helm-charts/atlas-operator/rbac.yaml` — **Regenerated** by `make helm-crds` — picks up the enriched Secret verbs.
- `bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml` — **Regenerated** by `make bundle` — CSV reflects the enriched Secret verbs.
- `test/helm/` — **New.** Unit tests pinning `helm template` output per chart (API-key Secret, SA Secret, mutual-exclusion rejection).
- `test/e2e2/serviceaccount_test.go` — **New.** E2E test scenarios.
- `docs/dev/td-service-accounts-phase1.md` — **New.** This document.

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
