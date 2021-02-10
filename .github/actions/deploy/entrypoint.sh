#!/bin/bash
echo "${INPUT_KUBE_CONFIG_DATA}" >> ./kube.config
export KUBECONFIG="./kube.config"

kubectl version

#Prepare CRDs
controller-gen crd:crdVersions=v1 rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

#Installing the CRD,Operator,Role
ns=mongodb-atlas-kubernetes-system
kubectl delete deployment mongodb-atlas-kubernetes-controller-manager -n "${ns}" || true # temporary
cd config/manager && kustomize edit set image controller="${INPUT_IMAGE_URL}"
cd - && kustomize build config/dev | kubectl apply -f -

# Ensuring the Atlas credentials Secret
kubectl delete secrets my-atlas-key --ignore-not-found -n "${ns}"
kubectl create secret generic my-atlas-key --from-literal="orgId=${INPUT_ATLAS_ORG_ID}" --from-literal="publicApiKey=${INPUT_ATLAS_PUBLIC_KEY}" --from-literal="privateApiKey=${INPUT_ATLAS_PRIVATE_KEY}" -n "${ns}"

# Wait for the Operator to start
cmd="while ! kubectl -n ${ns} get pods -l control-plane=controller-manager -o jsonpath={.items[0].status.phase} 2>/dev/null | grep -q Running ; do printf .; sleep 1; done"
timeout --foreground "1m" bash -c "${cmd}" || true
if ! kubectl -n "${ns}" get pods -l "control-plane=controller-manager" -o 'jsonpath="{.items[0].status.phase}"' | grep -q "Running"; then
    echo "Operator hasn't reached RUNNING state after 1 minute. The full yaml configuration for the pod is:"
    kubectl -n "${ns}" get pods -l "control-plane=controller-manager" -o yaml

    echo "Operator failed to start, exiting"
    exit 1
fi
