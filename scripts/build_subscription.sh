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

if [ -z "${CATALOG_DIR+x}" ]; then
	echo "CATALOG_DIR is not set"
	exit 1
fi

if [ -z "${CHANNEL+x}" ]; then
	echo "CHANNEL is not set"
	exit 1
fi

if [ -z "${TARGET_NAMESPACE+x}" ]; then
	echo "TARGET_NAMESPACE is not set"
	exit 1
fi

echo "Building subscription ${CATALOG_DIR}"

cat <<EOF> "${CATALOG_DIR}/subscription.yaml"
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: mongodb-atlas-operator-local
  namespace: ${TARGET_NAMESPACE} 
spec:
  channel: ${CHANNEL}
  name: mongodb-atlas-kubernetes
  source: mongodb-atlas-kubernetes-local
  sourceNamespace: openshift-marketplace
  installPlanApproval: Automatic
EOF
