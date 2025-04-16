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

# sample from https://kind.sigs.k8s.io/docs/user/local-registry/
set -o errexit

if [[ "${USE_KIND:-}" = "false" ]]; then
  echo "skipping creation of kind cluster"
  exit 0
fi

if [[ "$(kind get clusters | grep -x kind)" = "kind" ]]; then
  echo "kind cluster already exists"
  exit 0
fi

# create registry container unless it already exists
reg_name='kind-registry'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
nodes:
- role: control-plane
  # port forward 80 on the host to 80 on this node
  extraPortMappings:
  - containerPort: 30000
    hostPort: 30000
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30001
    hostPort: 30001
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30002
    hostPort: 30002
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30003
    hostPort: 30003
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30004
    hostPort: 30004
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30005
    hostPort: 30005
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30006
    hostPort: 30006
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30007
    hostPort: 30007
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30008
    hostPort: 30008
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30009
    hostPort: 30009
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30010
    hostPort: 30010
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30011
    hostPort: 30011
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30012
    hostPort: 30012
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30013
    hostPort: 30013
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30014
    hostPort: 30014
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30015
    hostPort: 30015
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30016
    hostPort: 30016
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30017
    hostPort: 30017
    listenAddress: "0.0.0.0"
    protocol: TCP
  - containerPort: 30018
    hostPort: 30018
    listenAddress: "0.0.0.0"
    protocol: TCP
EOF

# connect the registry to the cluster network
# (the network may already be connected)
docker network connect "kind" "${reg_name}" || true

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF
