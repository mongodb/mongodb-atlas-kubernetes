#!/bin/bash

set -eou pipefail

if [ -z "${CATALOG_DIR+x}" ]; then
	echo "CATALOG_DIR is not set"
	exit 1
fi

if [ -z "${CHANNEL+x}" ]; then
	echo "CHANNEL is not set"
	exit 1
fi

if [ -z "${TARGET_NAMESPACE+x}" ]; then
	echo "TARGET_NAMESPACE is not set"
	exit 1
fi

echo "Building subscription ${CATALOG_DIR}"

cat <<EOF> "${CATALOG_DIR}/subscription.yaml"
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mongodb-atlas-operator-local
  namespace: ${TARGET_NAMESPACE} 
spec:
  channel: ${CHANNEL}
  name: mongodb-atlas-kubernetes
  source: mongodb-atlas-kubernetes-local
  sourceNamespace: openshift-marketplace
  installPlanApproval: Automatic
EOF
