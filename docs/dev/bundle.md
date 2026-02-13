# OLM Bundle Generation

The operator is distributed on OpenShift via [OLM](https://olm.operatorframework.io/) (Operator Lifecycle Manager). OLM uses **bundles** — container images that package the operator's CRDs, RBAC, and a ClusterServiceVersion (CSV) manifest describing the operator to the OLM catalog.

## Key concepts

- **ClusterServiceVersion (CSV)**: The central manifest that describes the operator to OLM — its name, version, description, install strategy, owned CRDs, required permissions, and upgrade path (`replaces` field).
- **Bundle**: A directory (and corresponding container image) containing the CSV, CRDs, and OLM metadata. Lives in `bundle/` when generated locally.
- **Catalog**: A registry of bundles that OLM uses to discover and install operators. Built with `make catalog-build`.

## Directory layout

| Path | Purpose |
|---|---|
| `config/manifests-template/bases/` | CSV template — the hand-maintained base for the ClusterServiceVersion. Edit this to change operator metadata, description, install modes, or annotations. |
| `config/manifests/` | Generated kustomization that combines the CSV template with CRDs, RBAC, samples, and scorecard config. |
| `config/release/{dev,prod}/` | Kustomize overlays per environment. Each contains sub-overlays for `allinone`, `clusterwide`, `namespaced`, and `openshift` deployment modes. |
| `config/release/dev/dev_patch.json` | Dev-specific patches (e.g. QA API endpoint). |
| `config/release/prod/prod_patch.json` | Prod-specific patches. |
| `bundle/` | Generated output directory. Contains `manifests/` (CSV + CRDs) and `metadata/` (annotations). Created by `make bundle`, not checked in. |
| `bundle.Dockerfile` | Generated Dockerfile for building the bundle image. Created by `make bundle`, not checked in. |
| `version.json` | Source of truth for versioning. Contains `current` (the released version) and `next` (the upcoming version). |

## How `make bundle` works

The `bundle` target orchestrates several steps:

1. **`prepare-dirs`** — creates the `bundle/manifests/`, `deploy/` and sub-directories.

2. **Copy CSV template** — copies the base CSV from `config/manifests-template/bases/` into `bundle/manifests/`.

3. **`release`** — builds all deployment manifests (all-in-one, clusterwide, namespaced, OpenShift) via Kustomize from `config/release/$(ENV)/` overlays. Also copies CRDs into `deploy/crds/`.

4. **Generate kustomize manifests** — runs `operator-sdk generate kustomize manifests` to produce the `config/manifests/` kustomization from the template, CRDs, and API types.

5. **Generate bundle** — pipes the Kustomize output through `operator-sdk generate bundle` which produces the final `bundle/` directory with the CSV, CRDs, and metadata for the given `VERSION`.

6. **Patch CSV** — a series of post-processing steps:
   - Adds a `replaces: mongodb-atlas-kubernetes.v<CURRENT_VERSION>` field to enable OLM upgrades from the previous version.
   - Injects the `WATCH_NAMESPACE` environment variable sourced from the OLM `olm.targetNamespaces` annotation.
   - Sets the `containerImage` annotation to the operator image.

7. **Patch bundle.Dockerfile** — adds Red Hat delivery labels (`com.redhat.openshift.versions`, `com.redhat.delivery.backport`, `com.redhat.delivery.operator.bundle`) required for OpenShift certification.

## Makefile targets

| Target | Description |
|---|---|
| `make bundle` | Generate the full OLM bundle (CSV, CRDs, metadata, Dockerfile). Depends on `prepare-dirs`, CSV template copy, and `release`. |
| `make bundle-dev` | Run `make bundle` then patch the CSV to use the QA Atlas endpoint (`cloud-qa.mongodb.com`). |
| `make bundle-build` | Build the bundle container image from `bundle.Dockerfile`. |
| `make bundle-push` | Push the bundle image to the registry. |
| `make bundle-validate` | Validate the bundle with `operator-sdk bundle validate`. |
| `make clean-bundle` | Remove all generated bundle and deploy artifacts. |
| `make catalog-build` | Build an OLM catalog image containing the bundle. |
| `make catalog-push` | Push the catalog image. |
| `make deploy-olm` | Full OLM deployment: build and push bundle + catalog, then deploy to OpenShift via CatalogSource and Subscription. |

## Environment variable: `ENV`

The `ENV` variable (default: `dev`) selects which Kustomize overlay set is used under `config/release/`. Set `ENV=prod` for production bundles.

## Versioning

Versions are managed in `version.json`:

```json
{
  "current": "2.13.1",
  "next": "2.13.2"
}
```

- `current` — the last released version. Used in the CSV `replaces` field for OLM upgrade chains.
- `next` — the upcoming version. Used when `USE_NEXT_VERSION` is set.
- `make bump-version-file` — advances `current` to `next` and increments the minor version of `next`.

The effective `VERSION` defaults to `current` (plus `-dirty` if the working tree has uncommitted changes). Override with `VERSION=x.y.z make bundle`.

## Deploying a bundle locally

Deploy to a bare Kind cluster (not one started with `make run-kind`):

```shell
make bundle-build BUNDLE_IMG=$USERNAME/test-bundle:v$VERSION
make bundle-push BUNDLE_IMG=$USERNAME/test-bundle:v$VERSION
operator-sdk run bundle docker.io/$USERNAME/test-bundle:v$VERSION
```

Remove with:

```shell
operator-sdk cleanup mongodb-atlas-kubernetes
```

## Deploying via OLM on OpenShift

```shell
make deploy-olm
```

This builds and pushes the bundle and catalog images, then creates the CatalogSource, OperatorGroup, and Subscription in the target namespace.
