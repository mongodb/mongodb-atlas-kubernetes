#!/bin/sh

target_dir="deploy"
mkdir "${target_dir}"

# Generate configuration and save it to `all-in-one`
controller-gen crd:crdVersions=v1 rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && kustomize edit set image controller="${INPUT_IMAGE_URL}"
cd - && kustomize build config/prod > "${target_dir}"/all-in-one.yaml

cat "${target_dir}"/all-in-one.yaml
