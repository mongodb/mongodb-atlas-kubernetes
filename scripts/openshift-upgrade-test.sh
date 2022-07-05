#!/bin/bash

set -eou pipefail

# This test is designed to be launched from "mongodb-atlas-kubernetes/scripts" catalog
# Test conf
TEST_NAMESPACE=${TEST_NAMESPACE:-"atlas-upgrade-test"}
LATEST_RELEASE_VERSION="${LATEST_RELEASE_VERSION:-1.0.0}"
LATEST_RELEASE_REGISTRY=${LATEST_RELEASE_REGISTRY:-"quay.io/mongodb"}
REGISTRY=${REGISTRY:-"quay.io/igorkarpukhin"}

# This is used to build directory-based Openshift catalog
CATALOG_DIR="${CATALOG_DIR:-./openshift/atlas-catalog}"

if [ -z "${CURRENT_VERSION+x}" ]; then
  # opm doesn't allow 'v' prefix for versions
  CURRENT_VERSION=$(git describe --tags | awk '{gsub(/v/,"",$0); print}')
	echo "CURRENT_VERSION is not set. Setting to default: ${CURRENT_VERSION}"
fi

OPERATOR_NAME="mongodb-atlas-kubernetes-operator"
OPERATOR_IMAGE="${REGISTRY}/${OPERATOR_NAME}:${CURRENT_VERSION}"
OPERATOR_BUNDLE_IMAGE="${REGISTRY}/${OPERATOR_NAME}-bundle:${CURRENT_VERSION}"
LATEST_RELEASE_BUNDLE_IMAGE="${LATEST_RELEASE_REGISTRY}/${OPERATOR_NAME}-bundle:${LATEST_RELEASE_VERSION}"
OPERATOR_CATALOG_NAME="${OPERATOR_NAME}-catalog"
OPERATOR_CATALOG_IMAGE="${REGISTRY}/${OPERATOR_NAME}-catalog:${CURRENT_VERSION}"
OPERATOR_CATALOGSOURCE_NAME="${OPERATOR_CATALOG_NAME}-catalogsource"
OPERATOR_SUBSCRIPTION_NAME="${OPERATOR_NAME}-subscription"

millisecond=1
second=$(( 1000 * millisecond ))
DEFAULT_TIMEOUT=$((2 * second))

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
  #TODO: Added log collecting procedure
  echo "Done"
  return ${exit_code}
}
# Collect logs and remove all resources before exiting
#trap cleanup exit

try_until_success() {
  local cmd=$1
  local timeout=$2
  local interval=${3:-0.2}
  local now
  now="$(date +%s%3)"
  local expire=$(($now + $timeout))
  while [ "$now" -lt $expire ]; do
    if $cmd ; then
      echo "Passed"
      return 0
    fi
    sleep $interval
    now=$(date +%s%3)
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
  echo "Building and publishing bundle..."
  cd ../ && VERSION="${CURRENT_VERSION}" IMG="${OPERATOR_IMAGE}" OPERATOR_IMAGE="${OPERATOR_IMAGE}" BUNDLE_IMG="${OPERATOR_BUNDLE_IMAGE}" BUNDLE_METADATA_OPTS="--channels=candidate --default-channel=candidate" make image bundle bundle-build bundle-push
  echo "Bundle has been build and published"
  cd -
}

build_and_publish_catalog_with_two_channels() {
  echo "Building catalog with both bundles"
  rm -f "${CATALOG_DIR}"/*.yaml
  rm -f "$(dirname "${CATALOG_DIR}")"/atlas-catalog.Dockerfile

  mkdir -p "${CATALOG_DIR}"
  echo "Generating the dockerfile"
  opm alpha generate dockerfile "${CATALOG_DIR}"
  echo "Generating the catalog"

  # Stable - latest release, candidate - current version
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

  echo "Adding current version channel as CANDIDATE to ${CATALOG_DIR}/operator.yaml"
  echo "---
schema: olm.channel
package: mongodb-atlas-kubernetes
name: candidate
entries:
  - name: mongodb-atlas-kubernetes.v${CURRENT_VERSION}" >> "${CATALOG_DIR}"/operator.yaml

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
  channel: candidate
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

main() {
  echo "Test upgrade from ${LATEST_RELEASE_VERSION} to ${CURRENT_VERSION}"
  oc_login
  podman_version
  opm_version

  # Build and install previous version of the operator
  cleanup_previous_installation
  prepare_test_namespace
  build_and_publish_image_and_bundle
  build_and_publish_catalog_with_two_channels
  build_and_deploy_catalog_source
  build_and_deploy_operator_group
  build_and_deploy_subscription
#  wait_for_olm_to_install_operator
#
#  # Perform operator upgrade
#  patch_subscription
#  wait_for_olm_to_install_operator


  # make sure the is a PODMAN, OC and KUBECTL commands
  # Login to cluster with oc
  # Delete previous test namespace
  # Delete previous catalogsource if exists
  # Delete previous subscription if exists

  # Create test namespace
  # Build catalog from the current version
  # Create new catalog with the old bundle (e.g. 1.0.0) and a new one (e.g. 1.1.0)
  # Add catalogsource to the openshift-marketplace
  # Add subscription to the test namespace
  # Wait for the OLM to install the previous version of the operator
  # Create atlas project
  # Create atlas deployment
  # Wait for deployment

  # Patch subscription to point to the new bundle
  # Wait for the OLM to install new version
  # Check if the pod is up and running
  # Check if project and deployment are still exist
}

# Entrypoint to the test
main