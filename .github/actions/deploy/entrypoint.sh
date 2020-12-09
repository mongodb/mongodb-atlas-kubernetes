#!/bin/bash
echo "${INPUT_KUBE_CONFIG_DATA}" >> ./kube.config
export KUBECONFIG="./kube.config"

kubectl version

#Install CRDs
controller-gen crd:crdVersions=v1 rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -

#Installing the Operator
cd config/manager && kustomize edit set image controller="${INPUT_IMAGE}"
cd - && kustomize build config/default | kubectl apply -f -
