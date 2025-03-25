#!/bin/bash

set -euo pipefail

helm version
go version
cd test/e2e

# no `long-run`, no `broken` tests. `Long-run` tests run as a separate job
if [[ $TEST_NAME == "long-run" ]]; then
	filter="long-run && !broken";
else 
	filter="$TEST_NAME && !long-run && !broken";
fi

AKO_E2E_TEST=1 ginkgo --output-interceptor-mode=none --label-filter="${filter}" --timeout 120m --nodes=10 \
  --flake-attempts=1 --race --cover --v --coverpkg=github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/... \
  --coverprofile=coverprofile.out
