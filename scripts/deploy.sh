#!/bin/bash

set -eou pipefail

kubectl version

if [[ -z "${REGISTRY:-}" ]]; then
  REGISTRY="localhost:5000"
fi

image="${REGISTRY}/mongodb-atlas-kubernetes-operator"
docker build -t "${image}" .
docker push "${image}"

#Prepare CRDs
controller-gen crd:crdVersions=v1 rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

#Installing the CRD,Operator,Role
ns=mongodb-atlas-system
kubectl delete deployment mongodb-atlas-operator -n "${ns}" || true # temporary
cd config/manager && kustomize edit set image controller="${image}"
cd - && kustomize build --load-restrictor LoadRestrictionsNone config/release/dev/allinone | kubectl apply -f -

# Ensuring the Atlas credentials Secret
public_key=$(grep "ATLAS_PUBLIC_KEY" .actrc | cut -d "=" -f 2)
private_key=$(grep "ATLAS_PRIVATE_KEY" .actrc | cut -d "=" -f 2)
org_id=$(grep "ATLAS_ORG_ID" .actrc | cut -d "=" -f 2)

# both global and project keys
kubectl delete secrets my-atlas-key --ignore-not-found -n "${ns}"
kubectl create secret generic my-atlas-key --from-literal="orgId=${org_id}" --from-literal="publicApiKey=${public_key}" --from-literal="privateApiKey=${private_key}" -n "${ns}"
kubectl label secret my-atlas-key atlas.mongodb.com/type=credentials

kubectl delete secrets mongodb-atlas-operator-api-key --ignore-not-found -n "${ns}"
kubectl create secret generic mongodb-atlas-operator-api-key --from-literal="orgId=${org_id}" --from-literal="publicApiKey=${public_key}" --from-literal="privateApiKey=${private_key}" -n "${ns}"
kubectl label secret mongodb-atlas-operator-api-key atlas.mongodb.com/type=credentials

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
