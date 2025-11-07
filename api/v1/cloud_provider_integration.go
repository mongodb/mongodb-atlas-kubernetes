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

package v1

// CloudProviderIntegration define an integration to a cloud provider
type CloudProviderIntegration struct {
	// ProviderName is the name of the cloud provider. Currently only AWS is supported.
	ProviderName string `json:"providerName"`
	// IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.
	// +optional
	IamAssumedRoleArn string `json:"iamAssumedRoleArn"`
}

// CloudProviderAccessRole define an integration to a cloud provider
// DEPRECATED: This type is deprecated in favor of CloudProviderIntegration
type CloudProviderAccessRole struct {
	// ProviderName is the name of the cloud provider. Currently only AWS is supported.
	ProviderName string `json:"providerName"`
	// IamAssumedRoleArn is the ARN of the IAM role that is assumed by the Atlas cluster.
	// +optional
	IamAssumedRoleArn string `json:"iamAssumedRoleArn"`
}
