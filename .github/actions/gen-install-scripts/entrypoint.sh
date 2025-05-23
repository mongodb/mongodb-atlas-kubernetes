#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -xeou pipefail

target_dir="deploy"
clusterwide_dir="${target_dir}/clusterwide"
namespaced_dir="${target_dir}/namespaced"
crds_dir="${target_dir}/crds"
openshift="${target_dir}/openshift"

mkdir -p "${clusterwide_dir}"
mkdir -p "${namespaced_dir}"
mkdir -p "${crds_dir}"
mkdir -p "${openshift}"

# Generate configuration and save it to `all-in-one`
controller-gen crd:crdVersions=v1,ignoreUnexportedFields=true rbac:roleName=manager-role webhook paths="./api/..." paths="./internal/controller/..." output:crd:artifacts:config=config/crd/bases
cd config/manager && kustomize edit set image controller="${INPUT_IMAGE_URL}"
cd -
./scripts/split_roles_yaml.sh

which kustomize
kustomize version

# all-in-one
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/allinone" >"${target_dir}/all-in-one.yaml"
echo "Created all-in-one config"

# clusterwide
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/clusterwide" >"${clusterwide_dir}/clusterwide-config.yaml"
kustomize build "config/crd" >"${clusterwide_dir}/crds.yaml"
echo "Created clusterwide config"

# base-openshift-namespace-scoped
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/openshift" >"${openshift}/openshift.yaml"
kustomize build "config/crd" >"${openshift}/crds.yaml"
echo "Created openshift namespaced config"

# namespaced
kustomize build --load-restrictor LoadRestrictionsNone "config/release/${INPUT_ENV}/namespaced" >"${namespaced_dir}/namespaced-config.yaml"
kustomize build "config/crd" >"${namespaced_dir}/crds.yaml"
echo "Created namespaced config"

# crds
cp config/crd/bases/* "${crds_dir}"

# CSV bundle
operator-sdk generate kustomize manifests -q --apis-dir=api
# get the current version so we could put it into the "replaces:"
current_version="$(yq e '.metadata.name' bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml)"

# We pass the version only for non-dev deployments (it's ok to have "0.0.0" for dev)
channel="stable"
if [[ "${INPUT_ENV}" == "dev" ]]; then
  echo "build dev purpose"
  kustomize build --load-restrictor LoadRestrictionsNone config/manifests |
    operator-sdk generate bundle -q --overwrite --default-channel="${channel}" --channels="${channel}"
else
  echo "build release version"
  echo "${INPUT_IMAGE_URL}"
  kustomize build --load-restrictor LoadRestrictionsNone config/manifests |
    operator-sdk generate bundle -q --overwrite --version "${INPUT_VERSION}" --default-channel="${channel}" --channels="${channel}"
  # add replaces
  awk '!/replaces:/' bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml >tmp && mv tmp bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml
  echo "  replaces: $current_version" >>bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml
  # Add WATCH_NAMESPACE env parameter
  value="metadata.annotations['olm.targetNamespaces']" yq e -i '.spec.install.spec.deployments[0].spec.template.spec.containers[0].env[2] |= {"name": "WATCH_NAMESPACE", "valueFrom": {"fieldRef": {"fieldPath": env(value)}}}' bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml
  # Add containerImage to bundle/manifests/ csv. containerImage - The full location (registry, repository, name and tag) of the operator image
  yq e -i ".metadata.annotations.containerImage=\"${INPUT_IMAGE_URL}\"" bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml
fi

# add additional LABELs to bundle.Docker file
label="LABEL com.redhat.openshift.versions=\"v4.8-v4.18\"\nLABEL com.redhat.delivery.backport=true\nLABEL com.redhat.delivery.operator.bundle=true"
awk -v rep="FROM scratch\n\n$label" '{sub(/FROM scratch/, rep); print}' bundle.Dockerfile >tmp && mv tmp bundle.Dockerfile

operator-sdk bundle validate ./bundle
