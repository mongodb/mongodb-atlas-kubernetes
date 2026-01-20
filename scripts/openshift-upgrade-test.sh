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


set -eou pipefail

# This test is designed to be launched from "mongodb-atlas-kubernetes/scripts" catalog
#
# Before running, make sure you have the following tools:
# * opm 4.9+
# * kustomize
# * operator-sdk
# * controller-gen
# * yq
#
# Test conf
RANDOM_NAMESPACE_SUFFIX=${RANDOM}
TEST_NAMESPACE=${TEST_NAMESPACE:-"atlas-upgrade-test-${RANDOM_NAMESPACE_SUFFIX}"}
LATEST_RELEASE_VERSION="${LATEST_RELEASE_VERSION:-1.0.0}"
LATEST_RELEASE_REGISTRY=${LATEST_RELEASE_REGISTRY:-"quay.io/mongodb"}
REGISTRY=${REGISTRY:-"quay.io/mongodb"}

# This is used to build directory-based Openshift catalog for current version
CATALOG_DIR="${CATALOG_DIR:-./openshift/atlas-catalog}"
# This is used to build directory-based Openshift catalog for current RELEASE version
CATALOG_RELEASE_DIR="${CATALOG_RELEASE_DIR:-./openshift/atlas-catalog-release}"

if [ -z "${CURRENT_VERSION+x}" ]; then
  # opm doesn't allow 'v' prefix for versions
  CURRENT_VERSION=$(jq -r .current version.json)
	echo "CURRENT_VERSION is not set. Setting to default: ${CURRENT_VERSION}"
fi

OPERATOR_NAME="mongodb-atlas-kubernetes-operator-prerelease"
LATEST_RELEASE_OPERATOR_NAME="mongodb-atlas-kubernetes-operator"
NEW_OPERATOR_IMAGE="${REGISTRY}/${OPERATOR_NAME}:${CURRENT_VERSION}"
OPERATOR_BUNDLE_IMAGE="${REGISTRY}/${OPERATOR_NAME}-bundle:${CURRENT_VERSION}"
LATEST_RELEASE_BUNDLE_IMAGE="${LATEST_RELEASE_REGISTRY}/${LATEST_RELEASE_OPERATOR_NAME}-bundle:${LATEST_RELEASE_VERSION}"
OPERATOR_CATALOG_NAME="${OPERATOR_NAME}-catalog"
OPERATOR_CATALOG_IMAGE="${REGISTRY}/${OPERATOR_NAME}-catalog:${CURRENT_VERSION}"
OPERATOR_CATALOGSOURCE_NAME="${OPERATOR_CATALOG_NAME}"
OPERATOR_SUBSCRIPTION_NAME="${OPERATOR_NAME}-subscription"

second=1
minute=$(( 60 * second ))
DEFAULT_TIMEOUT=$((4 * minute))

if [ -z "${OC_TOKEN+x}" ]; then
	echo "OC_TOKEN is not set"
	exit 1
fi

if [ -z "${CLUSTER_API_URL+x}" ]; then
	echo "CLUSTER_API_URL is not set"
	exit 1
fi
# end test config

cleanup() {
  local exit_code="$?"
  set +e
  echo "Cleaning up..."
  # TODO: Added log collecting data if needed
  oc delete namespace "${TEST_NAMESPACE}"
  oc -n openshift-marketplace delete catalogsource "${OPERATOR_CATALOGSOURCE_NAME}" --ignore-not-found
  echo "Done"
  return ${exit_code}
}
# Collect logs and remove all resources before exiting
trap cleanup exit

try_until_success() {
  local cmd=$1
  local timeout=$2
  local interval=${3:-0.2}
  local now
  now="$(date +%s)"
  local expire=$((now + timeout))
  while [ "$now" -lt $expire ]; do
    if $cmd ; then
      echo "Passed"
      return 0
    fi
    sleep "$interval"
    now=$(date +%s)
  done
  echo "Fail"
  return 1
}

try_until_text() {
  local cmd=$1
  local expected=$2
  local timeout=$3
  local interval=${4:-1}
  local now
  now="$(date +%s)"
  local expire=$((now + timeout))
  while [ "$now" -lt $expire ]; do
    echo "Running ${cmd}"
    res=$($cmd || true)
    echo "Result: ${res}, Expected: ${expected}"
    if [[ ${res} == "${expected}" ]] ; then
        echo "Passed"
        return 0
    fi
    sleep "$interval"
    now=$(date +%s)
  done
  echo "Fail"
  return 1
}

expect_success_silent(){
  local cmd=$1
  if $cmd ; then
    return 0
  fi
  return 1
}

expect_success(){
  local cmd=$1
  echo "Running '$cmd'"
  if $cmd ; then
    return 0
  fi
  return 1
}

oc_login() {
  echo "Logging in to the cluster..."
  expect_success_silent "oc login --token=${OC_TOKEN} --server=${CLUSTER_API_URL}"
  echo "Logged in to the cluster"
  kubectl version
  oc get clusterversion/version
}

prepare_test_namespace() {
  expect_success "oc delete project ${TEST_NAMESPACE} --ignore-not-found" $((10 * second))
  try_until_success "oc create namespace ${TEST_NAMESPACE}" ${DEFAULT_TIMEOUT}
  try_until_success "oc project ${TEST_NAMESPACE}" ${DEFAULT_TIMEOUT}
}

podman_version() {
  expect_success "podman version"
}

opm_version() {
  expect_success "opm version"
}

deploy_latest_release() {
  echo "Deploying previous release"
}

cleanup_previous_installation() {
  echo "Removing previous installation"
  expect_success "oc -n openshift-marketplace delete catalogsource ${OPERATOR_CATALOGSOURCE_NAME} --ignore-not-found"
}

build_and_publish_image_and_bundle() {
  echo "Building bundle"
  cd ../
  VERSION="${CURRENT_VERSION}" IMG="${NEW_OPERATOR_IMAGE}" OPERATOR_IMAGE="${NEW_OPERATOR_IMAGE}" BUNDLE_IMG="${OPERATOR_BUNDLE_IMAGE}" BUNDLE_METADATA_OPTS="--channels=candidate --default-channel=candidate" make image bundle
  echo "Adding REPLACE parameter to the CSV"
  value="v1.0.0" yq e -i '.spec.replaces = env(value)' bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml
  echo "Publishing bundle"
  VERSION="${CURRENT_VERSION}" IMG="${NEW_OPERATOR_IMAGE}" OPERATOR_IMAGE="${NEW_OPERATOR_IMAGE}" BUNDLE_IMG="${OPERATOR_BUNDLE_IMAGE}" BUNDLE_METADATA_OPTS="--channels=candidate --default-channel=candidate" make bundle-build bundle-push
  cd -
  echo "Bundle has been build and published"
}

build_and_publish_catalog_with_two_channels() {
  echo "Building catalog with both bundles"
  rm -f "${CATALOG_DIR}"/*.yaml
  rm -f "$(dirname "${CATALOG_DIR}")"/atlas-catalog.Dockerfile

  mkdir -p "${CATALOG_DIR}"
  echo "Generating the dockerfile"
  opm alpha generate dockerfile "${CATALOG_DIR}"
  echo "Generating the catalog"

  # Stable - latest release, fast - current version
  opm init mongodb-atlas-kubernetes \
  	--default-channel="stable" \
  	--output=yaml \
  	> "${CATALOG_DIR}"/operator.yaml

  echo "Adding previous release ${LATEST_RELEASE_BUNDLE_IMAGE} to the catalog"
  opm render "${LATEST_RELEASE_BUNDLE_IMAGE}" --output=yaml >> "${CATALOG_DIR}"/operator.yaml

  echo "Adding ${OPERATOR_BUNDLE_IMAGE} to the catalog"
  opm render "${OPERATOR_BUNDLE_IMAGE}" --output=yaml >> "${CATALOG_DIR}"/operator.yaml

  echo "Adding previous release channel as STABLE to ${CATALOG_DIR}/operator.yaml"
  echo "---
schema: olm.channel
package: mongodb-atlas-kubernetes
name: stable
entries:
  - name: mongodb-atlas-kubernetes.v${LATEST_RELEASE_VERSION}" >> "${CATALOG_DIR}"/operator.yaml

  echo "Adding current version channel as FAST to ${CATALOG_DIR}/operator.yaml"
  echo "---
schema: olm.channel
package: mongodb-atlas-kubernetes
name: fast
entries:
  - name: mongodb-atlas-kubernetes.v${CURRENT_VERSION}
    skips:
      - mongodb-atlas-kubernetes.v${LATEST_RELEASE_VERSION}" >> "${CATALOG_DIR}"/operator.yaml

  echo "Validating catalog"
  expect_success "opm validate ${CATALOG_DIR}"
  echo "Catalog is valid"
  echo "Building catalog image"
  cd "$(dirname "${CATALOG_DIR}")" && docker build . -f atlas-catalog.Dockerfile -t "${OPERATOR_CATALOG_IMAGE}"
  expect_success "docker push ${OPERATOR_CATALOG_IMAGE}"
  echo "Catalog has been build and published"
  cd -
}

build_and_deploy_catalog_source() {
  echo "Creating catalog source"
  echo "
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: ${OPERATOR_CATALOGSOURCE_NAME}
  namespace: openshift-marketplace
spec:
  sourceType: grpc
  image: ${OPERATOR_CATALOG_IMAGE}
  displayName: MongoDB Atlas operator upgrade test from ${LATEST_RELEASE_VERSION} to ${CURRENT_VERSION}
  publisher: MongoDB
  updateStrategy:
    registryPoll:
      interval: 5m" > "${CATALOG_DIR}"/catalogsource.yaml
  echo "Applying catalogsource to openshift-marketplace namespace"
  expect_success "oc -n openshift-marketplace apply -f ${CATALOG_DIR}/catalogsource.yaml"
}

build_and_deploy_subscription() {
  echo "Deploying subscription ${OPERATOR_SUBSCRIPTION_NAME} to openshift-marketplace"
  echo "
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: ${OPERATOR_SUBSCRIPTION_NAME}
  namespace: ${TEST_NAMESPACE}
spec:
  channel: stable
  name: mongodb-atlas-kubernetes
  source: ${OPERATOR_CATALOGSOURCE_NAME}
  sourceNamespace: openshift-marketplace
  installPlanApproval: Automatic" > "${CATALOG_DIR}"/subscription.yaml
  expect_success "oc -n ${TEST_NAMESPACE} apply -f ${CATALOG_DIR}/subscription.yaml"
}

build_and_deploy_operator_group() {
  echo "Deploying operator group to ${TEST_NAMESPACE}"
  echo "
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: mongodb-group
  namespace: ${TEST_NAMESPACE}
spec:
  targetNamespaces:
    - ${TEST_NAMESPACE}
  " > "${CATALOG_DIR}"/operatorgroup.yaml
  expect_success "oc -n ${TEST_NAMESPACE} apply -f ${CATALOG_DIR}/operatorgroup.yaml"
}

wait_for_olm_to_install_operator() {
  try_until_text "oc -n ${TEST_NAMESPACE} get deployment mongodb-atlas-operator -o jsonpath='{.status.readyReplicas}'" "'1'" ${DEFAULT_TIMEOUT}
}

patch_subscription() {
  echo "Patching subscription to point to FAST channel"
  patch="[{\"op\": \"replace\", \"path\":\"/spec/channel\", \"value\": \"fast\"},{\"op\":\"replace\", \"path\":\"/spec/source\", \"value\":\"${OPERATOR_CATALOGSOURCE_NAME}\"}]"
  oc -n "${TEST_NAMESPACE}" patch subscription "${OPERATOR_SUBSCRIPTION_NAME}" --type json -p "${patch}"
  # Let the OLM approve the new InstallPlan
  sleep 20
}

wait_for_new_deployment() {
  echo "Waiting for OLM to create a pod with the new image..."
  try_until_success "oc -n ${TEST_NAMESPACE} rollout status deployment/mongodb-atlas-operator" ${DEFAULT_TIMEOUT}
}

verify_deployment () {
  echo "Validating deployment image to be ${NEW_OPERATOR_IMAGE}"
  try_until_text "oc -n ${TEST_NAMESPACE} get deployment mongodb-atlas-operator -o jsonpath='{.spec.template.spec.containers[0].image}'" "'${NEW_OPERATOR_IMAGE}'" ${DEFAULT_TIMEOUT}
}

main() {
  echo "Test upgrade from ${LATEST_RELEASE_VERSION} to ${CURRENT_VERSION}"
  oc_login
  opm_version

  # Build and install previous version of the operator
  cleanup_previous_installation
  prepare_test_namespace
  build_and_publish_image_and_bundle

  # Build and publish catalog with two bundles
  build_and_publish_catalog_with_two_channels

  # Build catalog and deploy CatalogSource
  build_and_deploy_catalog_source

  # Deploy OperatorGroup, subscription that points to the previous release
  build_and_deploy_operator_group
  build_and_deploy_subscription

  # Await for 1 replica to be ready
  wait_for_olm_to_install_operator

  # Perform operator upgrade
  # Patch subscription to point to the catalog with new version of the operator
  patch_subscription
  wait_for_new_deployment
  verify_deployment

  echo "Test passed!"
}

# Entrypoint to the test
main
