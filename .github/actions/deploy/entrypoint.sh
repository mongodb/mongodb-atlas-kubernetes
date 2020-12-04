#!/bin/bash
echo "$INPUT_KUBE_CONFIG_DATA" >> ./kube.config
export KUBECONFIG="./kube.config"

go version
kustomize version
kubectl version

make install
make deploy IMG="${INPUT_IMAGE}"
