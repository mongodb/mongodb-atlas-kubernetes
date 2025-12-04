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

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

// EncryptionAtRest configures the Encryption at Rest for the AWS, Azure and GCP providers.
type EncryptionAtRest struct {
	AwsKms         AwsKms         `json:"awsKms,omitempty"`         // AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
	AzureKeyVault  AzureKeyVault  `json:"azureKeyVault,omitempty"`  // AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
	GoogleCloudKms GoogleCloudKms `json:"googleCloudKms,omitempty"` // Specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
}

// AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AwsKms struct {
	// Specifies whether Encryption at Rest is enabled for an Atlas project.
	// To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	Enabled *bool `json:"enabled,omitempty"`
	// The AWS region in which the AWS customer master key exists.
	Region string `json:"region,omitempty"`
	// Specifies whether the encryption key set for the provider is valid and may be used to encrypt and decrypt data.
	Valid *bool `json:"valid,omitempty"`
	// A reference to as Secret containing the AccessKeyID, SecretAccessKey, CustomerMasterKeyID and RoleID fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

func (a *AwsKms) IsEnabled() bool {
	return a != nil && a.Enabled != nil && *a.Enabled
}

// AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AzureKeyVault struct {
	// Specifies whether Encryption at Rest is enabled for an Atlas project.
	// To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	Enabled *bool `json:"enabled,omitempty"`
	// The Client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.
	ClientID string `json:"clientID,omitempty"`
	// The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY
	AzureEnvironment string `json:"azureEnvironment,omitempty"`
	// The name of the Azure Resource group that contains an Azure Key Vault.
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// The unique identifier for an Azure AD tenant within an Azure subscription.
	TenantID string `json:"tenantID,omitempty"`
	// A reference to as Secret containing the SubscriptionID, KeyVaultName, KeyIdentifier, Secret fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

func (a *AzureKeyVault) IsEnabled() bool {
	return a != nil && a.Enabled != nil && *a.Enabled
}

// GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type GoogleCloudKms struct {
	// Specifies whether Encryption at Rest is enabled for an Atlas project.
	// To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	Enabled *bool `json:"enabled,omitempty"`
	// A reference to as Secret containing the ServiceAccountKey, KeyVersionResourceID fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

func (g *GoogleCloudKms) IsEnabled() bool {
	return g != nil && g.Enabled != nil && *g.Enabled
}
