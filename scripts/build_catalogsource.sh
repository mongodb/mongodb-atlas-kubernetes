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
