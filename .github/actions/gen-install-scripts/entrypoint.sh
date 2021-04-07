#!/bin/bash

set -xeou pipefail

target_dir="deploy"
clusterwide_dir="${target_dir}/clusterwide"
namespaced_dir="${target_dir}/namespaced"
crds_dir="${target_dir}/crds"

mkdir -p "${clusterwide_dir}"
mkdir -p "${namespaced_dir}"
mkdir -p "${crds_dir}"

# Generate configuration and save it to `all-in-one`
controller-gen crd:crdVersions=v1 rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && kustomize edit set image controller="${INPUT_IMAGE_URL}"
cd -


which kustomize
kustomize version

# all-in-one
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/allinone" > "${target_dir}/all-in-one.yaml"
echo "Created all-in-one config"

# clusterwide
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/clusterwide" > "${clusterwide_dir}/clusterwide-config.yaml"
kustomize build "config/crd" > "${clusterwide_dir}/crds.yaml"
echo "Created clusterwide config"

# namespaced
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/namespaced" > "${namespaced_dir}/namespaced-config.yaml"
kustomize build "config/crd" > "${namespaced_dir}/crds.yaml"
echo "Created namespaced config"

# crds
cp config/crd/bases/* "${crds_dir}"

# CSV bundle
# get the current version so we could put it into the "replaces:"
# TODO current_version=$(grep "name: mongodb-atlas-kubernetes.v0.5.0")
operator-sdk generate kustomize manifests -q --apis-dir=pkg/api

# We pass the version only for non-dev deployments (it's ok to have "0.0.0" for dev)
channel="beta"
if [[ "${INPUT_ENV}" == "dev" ]]; then
  kustomize build --load-restrictor LoadRestrictionsNone config/manifests | operator-sdk generate bundle -q --overwrite --default-channel="${channel}" --channels="${channel}"
else
  kustomize build --load-restrictor LoadRestrictionsNone config/manifests | operator-sdk generate bundle -q --overwrite --version "${INPUT_VERSION}" --default-channel="${channel}" --channels="${channel}"
fi

# temporary - should be done by composing kustomize
sed -i .bak '/runAsNonRoot: true/d' "bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml"
sed -i .bak '/runAsUser: 2000/d' "bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml"
# TODO
sed -i .bak '/replaces: mongodb-atlas-kubernetes.v0.4.0/d' "bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml"

rm "bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml.bak"

operator-sdk bundle validate ./bundle
