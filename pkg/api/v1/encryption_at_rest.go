package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"

	"go.mongodb.org/atlas/mongodbatlas"
)

// EncryptionAtRest allows to specify the Encryption at Rest for AWS, Azure and GCP providers
type EncryptionAtRest struct {
	AwsKms         AwsKms         `json:"awsKms,omitempty"`         // AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
	AzureKeyVault  AzureKeyVault  `json:"azureKeyVault,omitempty"`  // AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
	GoogleCloudKms GoogleCloudKms `json:"googleCloudKms,omitempty"` // Specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
}

// AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AwsKms struct {
	Enabled             *bool  `json:"enabled,omitempty"`             // Specifies whether Encryption at Rest is enabled for an Atlas project, To disable Encryption at Rest, pass only this parameter with a value of false, When you disable Encryption at Rest, Atlas also removes the configuration details.
	AccessKeyID         string `json:"accessKeyID,omitempty"`         // The IAM access key ID with permissions to access the customer master key specified by customerMasterKeyID.
	SecretAccessKey     string `json:"secretAccessKey,omitempty"`     // The IAM secret access key with permissions to access the customer master key specified by customerMasterKeyID.
	CustomerMasterKeyID string `json:"customerMasterKeyID,omitempty"` // The AWS customer master key used to encrypt and decrypt the MongoDB master keys.
	Region              string `json:"region,omitempty"`              // The AWS region in which the AWS customer master key exists: CA_CENTRAL_1, US_EAST_1, US_EAST_2, US_WEST_1, US_WEST_2, SA_EAST_1
	RoleID              string `json:"roleId,omitempty"`              // ID of an AWS IAM role authorized to manage an AWS customer master key.
	Valid               *bool  `json:"valid,omitempty"`               // Specifies whether the encryption key set for the provider is valid and may be used to encrypt and decrypt data.
}

// AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AzureKeyVault struct {
	Enabled           *bool  `json:"enabled,omitempty"`           // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ClientID          string `json:"clientID,omitempty"`          // The Client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.
	AzureEnvironment  string `json:"azureEnvironment,omitempty"`  // The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY
	SubscriptionID    string `json:"subscriptionID,omitempty"`    // The unique identifier associated with an Azure subscription.
	ResourceGroupName string `json:"resourceGroupName,omitempty"` // The name of the Azure Resource group that contains an Azure Key Vault.
	KeyVaultName      string `json:"keyVaultName,omitempty"`      // The name of an Azure Key Vault containing your key.
	KeyIdentifier     string `json:"keyIdentifier,omitempty"`     // The unique identifier of a key in an Azure Key Vault.
	Secret            string `json:"secret,omitempty"`            // The secret associated with the Azure Key Vault specified by azureKeyVault.tenantID.
	TenantID          string `json:"tenantID,omitempty"`          // The unique identifier for an Azure AD tenant within an Azure subscription.
}

// GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type GoogleCloudKms struct {
	Enabled              *bool  `json:"enabled,omitempty"`              // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ServiceAccountKey    string `json:"serviceAccountKey,omitempty"`    // String-formatted JSON object containing GCP KMS credentials from your GCP account.
	KeyVersionResourceID string `json:"keyVersionResourceID,omitempty"` // 	The Key Version Resource ID from your GCP account.
}

func (e EncryptionAtRest) ToAtlas(projectID string) (*mongodbatlas.EncryptionAtRest, error) {
	result := &mongodbatlas.EncryptionAtRest{
		GroupID: projectID,
	}
	err := compat.JSONCopy(result, e)
	return result, err
}
