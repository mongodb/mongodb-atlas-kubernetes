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
operator-sdk generate kustomize manifests -q --apis-dir=pkg/api
kustomize build --load-restrictor LoadRestrictionsNone config/manifests | operator-sdk generate bundle -q --overwrite --version "${INPUT_VERSION}"
operator-sdk bundle validate ./bundle
