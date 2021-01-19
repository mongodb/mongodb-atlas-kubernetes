#!/bin/bash
echo "${INPUT_KUBE_CONFIG_DATA}" >> ./kube.config
export KUBECONFIG="./kube.config"

kubectl version

#Install CRDs
controller-gen crd:crdVersions=v1 rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -

#Installing the Operator
kubectl delete deployment mongodb-atlas-kubernetes-controller-manager -n mongodb-atlas-kubernetes-system || true # temporary
cd config/manager && kustomize edit set image controller="${INPUT_IMAGE_URL}"
cd - && kustomize build config/default | kubectl apply -f -

# Ensuring the Atlas credentials Secret
kubectl delete secrets my-atlas-key --ignore-not-found -n mongodb-atlas-kubernetes-system
kubectl create secret generic my-atlas-key --from-literal="orgId=${INPUT_ATLAS_ORG_ID}" --from-literal="publicApiKey=${INPUT_ATLAS_PUBLIC_KEY}" --from-literal="privateApiKey=${INPUT_ATLAS_PRIVATE_KEY}" -n mongodb-atlas-kubernetes-system
