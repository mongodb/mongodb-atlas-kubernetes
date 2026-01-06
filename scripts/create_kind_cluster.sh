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

# Sourced from https://kind.sigs.k8s.io/examples/kind-with-registry.sh

# 1. Create registry container unless it already exists
reg_name='kind-registry'
reg_port='5000'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --network bridge --name "${reg_name}" \
    registry:2
fi

# 2. Create kind cluster with containerd registry config dir enabled
#
# NOTE: the containerd config patch is not necessary with images from kind v0.27.0+
# It may enable some older images to work similarly.
# If you're only supporting newer relases, you can just use `kind create cluster` here.
#
# See:
# https://github.com/kubernetes-sigs/kind/issues/2875
# https://github.com/containerd/containerd/blob/main/docs/cri/config.md#registry-configuration
# See: https://github.com/containerd/containerd/blob/main/docs/hosts.md
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "/etc/containerd/certs.d"
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

# 3. Add the registry config to the nodes
#
# This is necessary because localhost resolves to loopback addresses that are
# network-namespace local.
# In other words: localhost in the container is not localhost on the host.
#
# We want a consistent name that works from both ends, so we tell containerd to
# alias localhost:${reg_port} to the registry container when pulling images
REGISTRY_DIR="/etc/containerd/certs.d/localhost:${reg_port}"
for node in $(kind get nodes); do
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${reg_name}:5000"]
EOF
done

# 4. Connect the registry to the cluster network if not already connected
# This allows kind to bootstrap the network but ensures they're on the same network
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi

# 5. Document the local registry
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
