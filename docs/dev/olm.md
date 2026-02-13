## Operator sdk support for OLM and OLM bundles

> For a comprehensive reference on the bundle pipeline, directory layout, versioning, and Makefile targets, see [bundle.md](bundle.md).

One of the features that operator-sdk provides is support for [OLM](https://olm.operatorframework.io/docs/getting-started/) and bundles that it uses.

OLM is an "Operator Lifecycle Manager" developed by RedHat. It's shipped together with Openshift but can be installed
to "vanilla" clusters as well. When installed it provides the way to see, install (both from CLI and UI), upgrade Operators from the catalog.

### Installing OLM to K8s

```
operator-sdk olm install
```

### Creating/Updating OLM bundles

operator-sdk provides support for generating bundles (you can use the dev image here - copy the image from
Github action `Prepare E2E configuration and image`):

```
make bundle VERSION=0.5.0 IMG=mongodb/mongodb-atlas-kubernetes-operator-prerelease:CLOUDP-82589-ensure-cluster-secrets-ae16708
```

The following happens "under the hood":

1. `operator-sdk generate kustomize manifests -q --apis-dir=pkg/api`  - generates the `config/manifests` folder that contains
   the "base" for `ClusterServiceVersion` that will be generated in the end bundle. This contains some data like "description",
   "keywords", "installModes" which describe the Operator. Also the "customresourcedefinitions" is filled based on the APIs
   used (`--apis-dir=pkg/api` overrides this location from the default `api`). This base file can be edited manually if necessary.
2. `$(KUSTOMIZE) build --load-restrictor LoadRestrictionsNone config/manifests` - takes the `default` (which points to the
   "/release/prod/allinone" configuration), `samples` and `scorecard` configs and passed them all to
3. `operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)` which generates the final `bundle` directory
    with CRDs, and CSV

### Deploying the bundle

*Important! Deploying the bundle must happen on "bare" Kind cluster, not the one started with local registry support
(`make run-kind`) - otherwise the bundle is not installed*

As per https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/:
```
make bundle-build BUNDLE_IMG=$USERNAME/test-bundle:v0.5.0
make docker-push IMG=$USERNAME/test-bundle:v0.5.0
operator-sdk run bundle docker.io/$USERNAME/test-bundle:v0.5.0
```

### Removing the bundle

```
operator-sdk cleanup mongodb-atlas-kubernetes
```

### Full Output

```
operator-sdk olm install

➜  mongodb-atlas-kubernetes git:(CLOUDP-81621_olm_bundle) ✗ make bundle VERSION=0.4.0 IMG=mongodb/mongodb-atlas-kubernetes-operator-prerelease:main-2ba4112
/Users/alisovenko/workspace/mongodb-atlas-kubernetes/bin/controller-gen "crd:crdVersions=v1" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
operator-sdk generate kustomize manifests -q --apis-dir=pkg/api
cd config/manager && /Users/alisovenko/workspace/mongodb-atlas-kubernetes/bin/kustomize edit set image controller=mongodb/mongodb-atlas-kubernetes-operator:0.3.0
/Users/alisovenko/workspace/mongodb-atlas-kubernetes/bin/kustomize build --load-restrictor LoadRestrictionsNone config/manifests | operator-sdk generate bundle -q --overwrite --version 0.3.0
INFO[0000] Building annotations.yaml
INFO[0000] Writing annotations.yaml in /Users/alisovenko/workspace/mongodb-atlas-kubernetes/bundle/metadata
INFO[0000] Building Dockerfile
INFO[0000] Writing bundle.Dockerfile in /Users/alisovenko/workspace/mongodb-atlas-kubernetes
operator-sdk bundle validate ./bundle
INFO[0000] Found annotations file                        bundle-dir=bundle container-tool=docker
INFO[0000] Could not find optional dependencies file     bundle-dir=bundle container-tool=docker
INFO[0000] All validation tests have completed successfully


➜  mongodb-atlas-kubernetes git:(CLOUDP-81621_olm_bundle) ✗ make bundle-build BUNDLE_IMG=$USERNAME/test-bundle:v0.4.0


➜  mongodb-atlas-kubernetes git:(CLOUDP-81621_olm_bundle) ✗ make docker-push IMG=$USERNAME/test-bundle:v0.4.0

➜  mongodb-atlas-kubernetes git:(CLOUDP-81621_olm_bundle) ✗ operator-sdk run bundle docker.io/$USERNAME/test-bundle:v0.4.0
INFO[0010] Successfully created registry pod: docker-io-$USERNAME-test-bundle-v0-4-0
INFO[0010] Created CatalogSource: mongodb-atlas-kubernetes-catalog
INFO[0010] OperatorGroup "operator-sdk-og" created
INFO[0010] Created Subscription: mongodb-atlas-kubernetes-v0-4-0-sub
INFO[0016] Approved InstallPlan install-665lv for the Subscription: mongodb-atlas-kubernetes-v0-4-0-sub
INFO[0016] Waiting for ClusterServiceVersion "default/mongodb-atlas-kubernetes.v0.4.0" to reach 'Succeeded' phase
INFO[0016]   Waiting for ClusterServiceVersion "default/mongodb-atlas-kubernetes.v0.4.0" to appear
INFO[0026]   Found ClusterServiceVersion "default/mongodb-atlas-kubernetes.v0.4.0" phase: Pending
INFO[0028]   Found ClusterServiceVersion "default/mongodb-atlas-kubernetes.v0.4.0" phase: Installing
INFO[0055]   Found ClusterServiceVersion "default/mongodb-atlas-kubernetes.v0.4.0" phase: Succeeded
INFO[0055] OLM has successfully installed "mongodb-atlas-kubernetes.v0.4.0"

```
