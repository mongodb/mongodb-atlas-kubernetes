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
