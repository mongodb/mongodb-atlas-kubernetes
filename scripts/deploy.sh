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

kubectl version

if [[ -z "${REGISTRY:-}" ]]; then
  REGISTRY="localhost:5000"
fi

image="${REGISTRY}/mongodb-atlas-kubernetes-operator"
docker build -f fast.Dockerfile --rm -t "${image}" .
docker push "${image}"

#Prepare CRDs
controller-gen crd:crdVersions=v1,ignoreUnexportedFields=true rbac:roleName=manager-role webhook paths="./api/..." paths="./internal/controller/..." output:crd:artifacts:config=config/crd/bases

#Installing the CRD,Operator,Role
ns=mongodb-atlas-system
kubectl delete deployment mongodb-atlas-operator -n "${ns}" || true # temporary
cd config/manager && kustomize edit set image controller="${image}"
cd - && kustomize build --load-restrictor LoadRestrictionsNone config/release/dev/allinone | kubectl apply -f -

# Ensuring the Atlas credentials Secret
# Get credentials from environment variables (preferred) or prompt if not set
if [[ -z "${ATLAS_PUBLIC_KEY:-}" ]]; then
  echo "ATLAS_PUBLIC_KEY environment variable is not set"
  exit 1
fi
if [[ -z "${ATLAS_PRIVATE_KEY:-}" ]]; then
  echo "ATLAS_PRIVATE_KEY environment variable is not set"
  exit 1
fi
if [[ -z "${ATLAS_ORG_ID:-}" ]]; then
  echo "ATLAS_ORG_ID environment variable is not set"
  exit 1
fi
public_key="${ATLAS_PUBLIC_KEY}"
private_key="${ATLAS_PRIVATE_KEY}"
org_id="${ATLAS_ORG_ID}"

# both global and project keys
kubectl delete secrets my-atlas-key --ignore-not-found -n "${ns}"
kubectl create secret generic my-atlas-key --from-literal="orgId=${org_id}" --from-literal="publicApiKey=${public_key}" --from-literal="privateApiKey=${private_key}" -n "${ns}"
kubectl label secret my-atlas-key atlas.mongodb.com/type=credentials -n "${ns}"

kubectl delete secrets mongodb-atlas-operator-api-key --ignore-not-found -n "${ns}"
kubectl create secret generic mongodb-atlas-operator-api-key --from-literal="orgId=${org_id}" --from-literal="publicApiKey=${public_key}" --from-literal="privateApiKey=${private_key}" -n "${ns}"
kubectl label secret mongodb-atlas-operator-api-key atlas.mongodb.com/type=credentials -n "${ns}"

label="app.kubernetes.io/instance=mongodb-atlas-kubernetes-operator"
# Wait for the Operator to start
cmd="while ! kubectl -n ${ns} get pods -l $label -o jsonpath={.items[0].status.phase} 2>/dev/null | grep -q Running ; do printf .; sleep 1; done"
timeout --foreground "1m" bash -c "${cmd}" || true
if ! kubectl -n "${ns}" get pods -l "$label" -o 'jsonpath="{.items[0].status.phase}"' | grep -q "Running"; then
    echo "Operator hasn't reached RUNNING state after 1 minute. The full yaml configuration for the pod is:"
    kubectl -n "${ns}" get pods -l "$label" -o yaml

    echo "Operator failed to start, exiting"
    exit 1
fi
