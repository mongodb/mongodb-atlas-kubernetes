#!/bin/bash

set -eou pipefail

: ${TARGET_NAMESPACE:="mongodb-altas-system"}

if [[ -z $(oc get project ${TARGET_NAMESPACE} 2> /dev/null) ]]; then
	oc create project ${TARGET_NAMESPACE}
fi
