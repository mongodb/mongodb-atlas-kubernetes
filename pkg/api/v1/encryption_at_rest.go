package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

// EncryptionAtRest allows to specify the Encryption at Rest for AWS, Azure and GCP providers
type EncryptionAtRest struct {
	AwsKms         AwsKms         `json:"awsKms,omitempty"`         // AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
	AzureKeyVault  AzureKeyVault  `json:"azureKeyVault,omitempty"`  // AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
	GoogleCloudKms GoogleCloudKms `json:"googleCloudKms,omitempty"` // Specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
}

// AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AwsKms struct {
	Enabled *bool  `json:"enabled,omitempty"` // Specifies whether Encryption at Rest is enabled for an Atlas project, To disable Encryption at Rest, pass only this parameter with a value of false, When you disable Encryption at Rest, Atlas also removes the configuration details.
	Region  string `json:"region,omitempty"`  // The AWS region in which the AWS customer master key exists: CA_CENTRAL_1, US_EAST_1, US_EAST_2, US_WEST_1, US_WEST_2, SA_EAST_1
	Valid   *bool  `json:"valid,omitempty"`   // Specifies whether the encryption key set for the provider is valid and may be used to encrypt and decrypt data.
	// A reference to as Secret containing the AccessKeyID, SecretAccessKey, CustomerMasterKeyID and RoleID fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

// AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AzureKeyVault struct {
	Enabled           *bool  `json:"enabled,omitempty"`           // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ClientID          string `json:"clientID,omitempty"`          // The Client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.
	AzureEnvironment  string `json:"azureEnvironment,omitempty"`  // The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY
	ResourceGroupName string `json:"resourceGroupName,omitempty"` // The name of the Azure Resource group that contains an Azure Key Vault.
	TenantID          string `json:"tenantID,omitempty"`          // The unique identifier for an Azure AD tenant within an Azure subscription.
	// A reference to as Secret containing the SubscriptionID, KeyVaultName, KeyIdentifier, Secret fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

// GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type GoogleCloudKms struct {
	Enabled *bool `json:"enabled,omitempty"` // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	// A reference to as Secret containing the ServiceAccountKey, KeyVersionResourceID fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}
