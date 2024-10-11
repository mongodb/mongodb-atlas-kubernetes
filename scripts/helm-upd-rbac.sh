#!/bin/bash

set -eou pipefail

if [[ -z "${HELM_RBAC_FILE}" ]]; then
  echo "HELM_RBAC_FILE is not set"
  exit 1
fi

yq '.spec.install.spec.clusterPermissions[0].rules' ./bundle/manifests/mongodb-atlas-kubernetes.clusterserviceversion.yaml > rbac.yaml

echo "Comparing RBAC for CSV to RBAC in AKO helm chart"
if ! diff rbac.yaml "$HELM_RBAC_FILE"; then
  echo "Copying RBAC"
  cp rbac.yaml "$HELM_RBAC_FILE"
else
  echo "No changes detected"
fi
