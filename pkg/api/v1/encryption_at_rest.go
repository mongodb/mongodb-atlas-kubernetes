package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
)

// EncryptionAtRest allows to specify the Encryption at Rest for AWS, Azure and GCP providers
type EncryptionAtRest struct {
	AwsKms         AwsKms         `json:"awsKms,omitempty"`         // AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
	AzureKeyVault  AzureKeyVault  `json:"azureKeyVault,omitempty"`  // AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
	GoogleCloudKms GoogleCloudKms `json:"googleCloudKms,omitempty"` // Specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
}

// AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AwsKms struct {
	Enabled             *bool  `json:"enabled,omitempty"` // Specifies whether Encryption at Rest is enabled for an Atlas project, To disable Encryption at Rest, pass only this parameter with a value of false, When you disable Encryption at Rest, Atlas also removes the configuration details.
	accessKeyID         string // The IAM access key ID with permissions to access the customer master key specified by customerMasterKeyID.
	secretAccessKey     string // The IAM secret access key with permissions to access the customer master key specified by customerMasterKeyID.
	customerMasterKeyID string // The AWS customer master key used to encrypt and decrypt the MongoDB master keys.
	Region              string `json:"region,omitempty"` // The AWS region in which the AWS customer master key exists: CA_CENTRAL_1, US_EAST_1, US_EAST_2, US_WEST_1, US_WEST_2, SA_EAST_1
	roleID              string // ID of an AWS IAM role authorized to manage an AWS customer master key.
	Valid               *bool  `json:"valid,omitempty"` // Specifies whether the encryption key set for the provider is valid and may be used to encrypt and decrypt data.
	// A reference to as Secret containing the AccessKeyID, SecretAccessKey, CustomerMasterKey and RoleID fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

// AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AzureKeyVault struct {
	Enabled           *bool  `json:"enabled,omitempty"`          // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ClientID          string `json:"clientID,omitempty"`         // The Client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.
	AzureEnvironment  string `json:"azureEnvironment,omitempty"` // The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY
	subscriptionID    string // The unique identifier associated with an Azure subscription.
	ResourceGroupName string `json:"resourceGroupName,omitempty"` // The name of the Azure Resource group that contains an Azure Key Vault.
	keyVaultName      string // The name of an Azure Key Vault containing your key.
	keyIdentifier     string // The unique identifier of a key in an Azure Key Vault.
	secret            string // The secret associated with the Azure Key Vault specified by azureKeyVault.tenantID.
	TenantID          string `json:"tenantID,omitempty"` // The unique identifier for an Azure AD tenant within an Azure subscription.
	// A reference to as Secret containing the SubscriptionID, KeyVaultName, KeyIdentifier, Secret fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

// GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type GoogleCloudKms struct {
	Enabled              *bool  `json:"enabled,omitempty"`              // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ServiceAccountKey    string `json:"serviceAccountKey,omitempty"`    // String-formatted JSON object containing GCP KMS credentials from your GCP account.
	KeyVersionResourceID string `json:"keyVersionResourceID,omitempty"` // 	The Key Version Resource ID from your GCP account.
	// A reference to as Secret containing the ServiceAccountKey, KeyVersionResourceID fields
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
}

func (e EncryptionAtRest) ToAtlas(projectID string) (*mongodbatlas.EncryptionAtRest, error) {
	result := &mongodbatlas.EncryptionAtRest{
		GroupID: projectID,
	}

	err := compat.JSONCopy(result, e)
	return result, err
}

func (a AwsKms) ToAtlas() mongodbatlas.AwsKms {
	return mongodbatlas.AwsKms{
		Enabled:             a.Enabled,
		AccessKeyID:         a.accessKeyID,
		SecretAccessKey:     a.secretAccessKey,
		RoleID:              a.roleID,
		CustomerMasterKeyID: a.customerMasterKeyID,
		Region:              a.Region,
		Valid:               a.Valid,
	}
}

func (g GoogleCloudKms) ToAtlas() mongodbatlas.GoogleCloudKms {
	return mongodbatlas.GoogleCloudKms{
		Enabled:              g.Enabled,
		ServiceAccountKey:    g.ServiceAccountKey,
		KeyVersionResourceID: g.KeyVersionResourceID,
	}
}

func (az AzureKeyVault) ToAtlas() mongodbatlas.AzureKeyVault {
	return mongodbatlas.AzureKeyVault{
		Enabled:           az.Enabled,
		ClientID:          az.ClientID,
		AzureEnvironment:  az.AzureEnvironment,
		SubscriptionID:    az.subscriptionID,
		ResourceGroupName: az.ResourceGroupName,
		KeyVaultName:      az.keyVaultName,
		KeyIdentifier:     az.keyIdentifier,
		TenantID:          az.TenantID,
		Secret:            az.secret,
	}
}
