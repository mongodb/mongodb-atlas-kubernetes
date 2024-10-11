#!/bin/bash

set -eou pipefail

if [[ -z "${HELM_CRDS_PATH}" ]]; then
  echo "HELM_CRDS_PATH is not set"
  exit 1
fi

filesToCopy=()
for filename in ./bundle/manifests/atlas.mongodb.com_*.yaml; do
  absName="$(basename "$filename")"
  echo "Veryfing file: ${filename}"
  if ! diff "$filename" "${HELM_CRDS_PATH}"/"$absName"; then
    filesToCopy+=("$filename")
  fi
done

fLen=${#filesToCopy[@]}
if [ "$fLen" -eq 0 ]; then
  echo "No CRD changes detected"
  exit 0
fi

echo "The following CRD changes detected:"
for (( i=0; i < fLen; i++ )); do
  echo "${filesToCopy[$i]}"
done

for (( i=0; i < fLen; i++ )); do
  echo "Copying ${filesToCopy[$i]} to ${HELM_CRDS_PATH}/"
  cp "${filesToCopy[$i]}" "${HELM_CRDS_PATH}"/
done
