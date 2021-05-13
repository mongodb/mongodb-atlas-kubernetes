#!/bin/bash
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

if [ -z "${CONTAINER_ENGINE+x}" ]; then
	echo "CONTAINER_ENGINE is not set"
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
cd "$(dirname "${CATALOG_DIR}")" && ${CONTAINER_ENGINE} build . -f atlas-catalog.Dockerfile -t "${CATALOG_IMAGE}"
