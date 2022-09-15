#!/bin/bash

set -eou pipefail

if [ -z "${CATALOG_DIR+x}" ]; then
	echo "CATALOG_DIR is not set"
	exit 1
fi

if [ -z "${CATALOG_IMAGE+x}" ]; then
	echo "CATALOG_IMAGE is not set"
	exit 1
fi

if [ -z "${CATALOG_DISPLAY_NAME+x}" ]; then
  CATALOG_DISPLAY_NAME="MongoDB Atlas operator local"
  echo "CATALOG_DISPLAY_NAME is not set. Setting to default: ${CATALOG_DISPLAY_NAME}"
fi

echo "Building catalog ${CATALOG_IMAGE}"

cat <<EOF> "${CATALOG_DIR}/catalogsource.yaml"
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: mongodb-atlas-kubernetes-local
spec:
  sourceType: grpc
  image: ${CATALOG_IMAGE}
  displayName: ${CATALOG_DISPLAY_NAME}
  publisher: MongoDB
  updateStrategy:
    registryPoll:
      interval: 10m
EOF
