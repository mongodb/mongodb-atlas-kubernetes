// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"testing"
)

const atlasOperatorChartPath = "../../helm-charts/atlas-operator"

func TestAtlasOperator_RendersAPIKeySecret(t *testing.T) {
	assertAPIKeySecret(t, atlasOperatorChartPath, "atlas_operator_apikey_values.yaml")
}

func TestAtlasOperator_RendersServiceAccountSecret(t *testing.T) {
	assertServiceAccountSecret(t, atlasOperatorChartPath, "atlas_operator_sa_values.yaml")
}

func TestAtlasOperator_RejectsBothCredentialTypes(t *testing.T) {
	assertRejectsBothCredentialTypes(t, atlasOperatorChartPath, "atlas_operator_sa_and_pka_values.yaml",
		"set either (publicApiKey,privateApiKey) or (clientId,clientSecret), not both")
}
