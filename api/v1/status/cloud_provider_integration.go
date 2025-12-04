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

package status

type CloudProviderIntegration struct {
	// Amazon Resource Name that identifies the Amazon Web Services user account that MongoDB Atlas uses when it assumes the Identity and Access Management role.
	AtlasAWSAccountArn string `json:"atlasAWSAccountArn,omitempty"`
	// Unique external ID that MongoDB Atlas uses when it assumes the IAM role in your Amazon Web Services account.
	AtlasAssumedRoleExternalID string `json:"atlasAssumedRoleExternalId"`
	// Date and time when someone authorized this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	AuthorizedDate string `json:"authorizedDate,omitempty"`
	// Date and time when someone created this role for the specified cloud service provider. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	CreatedDate string `json:"createdDate,omitempty"`
	// List that contains application features associated with this Amazon Web Services Identity and Access Management role.
	FeatureUsages []FeatureUsage `json:"featureUsages,omitempty"`
	// Amazon Resource Name that identifies the Amazon Web Services Identity and Access Management role that MongoDB Cloud assumes when it accesses resources in your AWS account.
	IamAssumedRoleArn string `json:"iamAssumedRoleArn,omitempty"`
	// Human-readable label that identifies the cloud provider of the role.
	ProviderName string `json:"providerName"`
	// Unique 24-hexadecimal digit string that identifies the role.
	RoleID string `json:"roleId,omitempty"`
	// Provision status of the service account.
	// Values are IN_PROGRESS, COMPLETE, FAILED, or NOT_INITIATED.
	Status string `json:"status,omitempty"`
	// Application error message returned.
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type FeatureUsage struct {
	// Human-readable label that describes one MongoDB Cloud feature linked to this Amazon Web Services Identity and Access Management role.
	FeatureType string `json:"featureType,omitempty"`
	// Identifying characteristics about the data lake linked to this Amazon Web Services Identity and Access Management role.
	FeatureID string `json:"featureId,omitempty"`
}

const (
	CloudProviderIntegrationStatusNew                 = "NEW"
	CloudProviderIntegrationStatusCreated             = "CREATED"
	CloudProviderIntegrationStatusAuthorized          = "AUTHORIZED"
	CloudProviderIntegrationStatusDeAuthorize         = "DEAUTHORIZE"
	CloudProviderIntegrationStatusFailedToCreate      = "FAILED_TO_CREATE"
	CloudProviderIntegrationStatusFailedToAuthorize   = "FAILED_TO_AUTHORIZE"
	CloudProviderIntegrationStatusFailedToDeAuthorize = "FAILED_TO_DEAUTHORIZE"

	StatusFailed = "FAILED"
	StatusReady  = "READY"
)

func NewCloudProviderIntegration(providerName, assumedRoleArn string) CloudProviderIntegration {
	return CloudProviderIntegration{
		ProviderName:      providerName,
		IamAssumedRoleArn: assumedRoleArn,
		Status:            CloudProviderIntegrationStatusNew,
	}
}
