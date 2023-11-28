#!/bin/bash

set -eup pipefail

dir=${METRICS_DIR:-metrics}
collection=${COLLECTION:-collection}
metrics=${METRICS_JSON:-"[]"}

echo "Inserting into ${collection} metrics: ${metrics}"
mkdir -p "${dir}"
if [ ! -f "${dir}/${collection}.json" ]; then
	echo "[]" > "${dir}/${collection}.json"
fi
jq ". += ${metrics}" < "${dir}/${collection}.json" > "${dir}/${collection}.new.json"
mv "${dir}/${collection}.new.json" "${dir}/${collection}.json"
