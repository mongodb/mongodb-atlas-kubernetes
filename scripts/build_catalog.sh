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

echo "Building catalog"

if [ -z "${CATALOG_DIR+x}" ]; then
	echo "CATALOG_DIR is not set"
	exit 1
fi

if [ -z "${CHANNEL+x}" ]; then
	echo "CHANNEL is not set"
	exit 1
fi

if [ -z "${CATALOG_IMAGE+x}" ]; then
	echo "CATALOG_IMAGE is not set"
	exit 1
fi

if [ -z "${BUNDLE_IMAGE+x}" ]; then
	echo "BUNDLE_IMAGE is not set"
	exit 1
fi

if [ -z "${VERSION+x}" ]; then
	echo "VERSION is not set"
	exit 1
fi

rm -f "${CATALOG_DIR}"/operator.yaml
rm -f "${CATALOG_DIR}"/channel.yaml
rm -f "$(dirname "${CATALOG_DIR}")"/atlas-catalog.Dockerfile

mkdir -p "${CATALOG_DIR}"

echo "Generating the dockerfile"
opm alpha generate dockerfile "${CATALOG_DIR}"
echo "Generating the catalog"
opm init mongodb-atlas-kubernetes \
	--default-channel="${CHANNEL}" \
	--output=yaml \
	> "${CATALOG_DIR}"/operator.yaml
echo "Adding ${BUNDLE_IMAGE} to the catalog"
opm render "${BUNDLE_IMAGE}" --output=yaml \
	>> "${CATALOG_DIR}"/operator.yaml

echo "Adding ${CHANNEL} channel to the catalog. Version: ${VERSION}"
cat <<EOF> "${CATALOG_DIR}"/channel.yaml
---
schema: olm.channel
package: mongodb-atlas-kubernetes
name: ${CHANNEL}
entries:
  - name: mongodb-atlas-kubernetes.v${VERSION}
EOF

echo "Validating catalog..."
opm validate "${CATALOG_DIR}"
echo "Catalog is valid"
echo "Building catalog image"
cd "$(dirname "${CATALOG_DIR}")" && docker build . -f atlas-catalog.Dockerfile -t "${CATALOG_IMAGE}"
