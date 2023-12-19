#!/bin/bash

set -euo pipefail

# TODO: remove the need for this script ASAP
MODULE_VERSION=""
if [[ ! "${VERSION}" =~ 1\..* ]]; then
    if [[ "${VERSION}" =~ [0-9]+ ]]; then
        matched_digits=${BASH_REMATCH[0]}
        MODULE_VERSION="/v${matched_digits}"
    else
        echo "No matching digits found"
        exit 1
    fi
fi
echo "github.com/mongodb/mongodb-atlas-kubernetes${MODULE_VERSION}"
