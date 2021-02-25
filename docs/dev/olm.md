## Operator sdk support for OLM and OLM bundles

One of the features that operator-sdk provides is support for [OLM](https://olm.operatorframework.io/docs/getting-started/) and bundles that it uses.

OLM is an "Operator Lifecycle Manager" developed by RedHat. It's shipped together with Openshift but can be installed 
to "vanilla" clusters as well. When installed it provides the way to see, install (both from CLI and UI), upgrade Operators from the catalog.

### Installing OLM to K8s

```
operator-sdk olm install
```

### Creating/Updating OLM bundles

operator-sdk provides support for generating bundles:

```
make bundle VERSION=0.3.0 IMG=mongodb/mongodb-atlas-kubernetes-operator:0.3.0
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

As per https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/:
```
make bundle-build BUNDLE_IMG=antonlisovenko/test-bundle:v0.3.0
make docker-push IMG=antonlisovenko/test-bundle:v0.3.0
operator-sdk run bundle docker.io/antonlisovenko/test-bundle:v0.3.0
```
(so far the last command fails)