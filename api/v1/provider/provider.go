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

package provider

const (
	ProviderAWS        ProviderName = "AWS"
	ProviderGCP        ProviderName = "GCP"
	ProviderAzure      ProviderName = "AZURE"
	ProviderTenant     ProviderName = "TENANT"
	ProviderServerless ProviderName = "SERVERLESS"
)

type ProviderName string
type CloudProviders map[ProviderName]struct{}

func (cp *CloudProviders) IsSupported(name ProviderName) bool {
	_, ok := (*cp)[name]

	return ok
}

func SupportedProviders() CloudProviders {
	return CloudProviders{
		ProviderAWS:   {},
		ProviderGCP:   {},
		ProviderAzure: {},
	}
}
