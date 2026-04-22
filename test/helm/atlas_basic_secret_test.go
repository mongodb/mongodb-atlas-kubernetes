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

const atlasBasicChartPath = "../../helm-charts/atlas-basic"

func TestAtlasBasic_RendersAPIKeySecret(t *testing.T) {
	assertAPIKeySecret(t, atlasBasicChartPath, "atlas_basic_apikey_values.yaml")
}

func TestAtlasBasic_RendersServiceAccountSecret(t *testing.T) {
	assertServiceAccountSecret(t, atlasBasicChartPath, "atlas_basic_sa_values.yaml")
}

func TestAtlasBasic_RejectsBothCredentialTypes(t *testing.T) {
	assertRejectsBothCredentialTypes(t, atlasBasicChartPath, "atlas_basic_both_values.yaml",
		"set either (publicKey,privateKey) or (clientId,clientSecret), not both")
}
