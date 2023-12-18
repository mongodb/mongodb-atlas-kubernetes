#!/bin/bash

set -euo pipefail

BIN_VERSION=$("${BINARY}" -v)

if [ "${BIN_VERSION}" == "unknown" ]; then
    echo "${BINARY} version ${BIN_VERSION}: was not set"
    exit 1
elif [[ "${BIN_VERSION}" =~ .*-dirty ]]; then
    echo "${BINARY} version ${BIN_VERSION}: is dirty"
    exit 1
elif [ "${BIN_VERSION}" != "${VERSION}" ]; then
    echo "${BINARY} version ${BIN_VERSION}: does not match expected ${VERSION}"
    exit 1
fi

echo "${BINARY} version ${BIN_VERSION}: OK"
exit 0
