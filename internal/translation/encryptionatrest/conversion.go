package encryptionatrest

import (
	"reflect"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type EncryptionAtRest struct {
	AWS   AwsKms
	Azure AzureKeyVault
	GCP   GoogleCloudKms
}

type AwsKms struct {
	akov2.AwsKms
	AccessKeyID                  string
	SecretAccessKey              string
	CustomerMasterKeyID          string
	RoleID                       string
	CloudProviderIntegrationRole string
}

func (a *AwsKms) SetSecrets(customerMasterKeyID, roleID string) {
	a.CustomerMasterKeyID = customerMasterKeyID
	a.RoleID = roleID
}

func (a *AwsKms) ToAtlas() *admin.AWSKMSConfiguration {
	if a == nil {
		return nil
	}

	result := &admin.AWSKMSConfiguration{
		Enabled: a.AwsKms.Enabled,
		Region:  &a.AwsKms.Region,
		Valid:   a.AwsKms.Valid,

		RoleId:              &a.RoleID,
		CustomerMasterKeyID: &a.CustomerMasterKeyID,
	}

	if result.RoleId == nil && a.CloudProviderIntegrationRole != "" {
		result.RoleId = &a.CloudProviderIntegrationRole
	}

	return result
}

type AzureKeyVault struct {
	akov2.AzureKeyVault
	SubscriptionID string
	KeyVaultName   string
	KeyIdentifier  string
	Secret         string
}

func (az *AzureKeyVault) SetSecrets(subscriptionID, keyVaultName, keyIdentifier, secret string) {
	az.SubscriptionID = subscriptionID
	az.KeyVaultName = keyVaultName
	az.KeyIdentifier = keyIdentifier
	az.Secret = secret
}

func (az *AzureKeyVault) ToAtlas() *admin.AzureKeyVault {
	if az == nil {
		return nil
	}

	return &admin.AzureKeyVault{
		Enabled:           az.AzureKeyVault.Enabled,
		ClientID:          &az.AzureKeyVault.ClientID,
		AzureEnvironment:  &az.AzureKeyVault.AzureEnvironment,
		ResourceGroupName: &az.AzureKeyVault.ResourceGroupName,
		TenantID:          &az.AzureKeyVault.TenantID,

		SubscriptionID: &az.SubscriptionID,
		KeyVaultName:   &az.KeyVaultName,
		KeyIdentifier:  &az.KeyIdentifier,
		Secret:         &az.Secret,
	}
}

type GoogleCloudKms struct {
	akov2.GoogleCloudKms
	ServiceAccountKey    string
	KeyVersionResourceID string
}

func (g *GoogleCloudKms) SetSecrets(serviceAccountKey, keyVersionResourceID string) {
	g.ServiceAccountKey = serviceAccountKey
	g.KeyVersionResourceID = keyVersionResourceID
}

func (g *GoogleCloudKms) ToAtlas() *admin.GoogleCloudKMS {
	if g == nil {
		return nil
	}

	return &admin.GoogleCloudKMS{
		Enabled: g.GoogleCloudKms.Enabled,

		ServiceAccountKey:    &g.ServiceAccountKey,
		KeyVersionResourceID: &g.KeyVersionResourceID,
	}
}

func NewEncryptionAtRest(project *akov2.AtlasProject) *EncryptionAtRest {
	spec := project.Spec.EncryptionAtRest
	if spec == nil {
		return nil
	}

	ear := &EncryptionAtRest{}
	if spec.AwsKms.IsEnabled() {
		ear.AWS.AwsKms = spec.AwsKms
		for _, role := range project.Status.CloudProviderIntegrations {
			if role.ProviderName == "AWS" {
				ear.AWS.CloudProviderIntegrationRole = role.RoleID
			}
		}
	}
	if spec.AzureKeyVault.IsEnabled() {
		ear.Azure.AzureKeyVault = spec.AzureKeyVault
	}
	if spec.GoogleCloudKms.IsEnabled() {
		ear.GCP.GoogleCloudKms = spec.GoogleCloudKms
	}
	return ear
}

func toAtlas(spec *EncryptionAtRest) *admin.EncryptionAtRest {
	if spec == nil {
		return nil
	}

	return &admin.EncryptionAtRest{
		AwsKms:         spec.AWS.ToAtlas(),
		AzureKeyVault:  spec.Azure.ToAtlas(),
		GoogleCloudKms: spec.GCP.ToAtlas(),
	}
}

func fromAtlas(ear *admin.EncryptionAtRest) *EncryptionAtRest {
	out := &EncryptionAtRest{}
	if ear.HasAwsKms() {
		out.AWS.AwsKms = akov2.AwsKms{
			Enabled: ear.AwsKms.Enabled,
			Region:  ear.AwsKms.GetRegion(),
			Valid:   ear.AwsKms.Valid,
		}
	}
	if ear.HasAzureKeyVault() {
		out.Azure.AzureKeyVault = akov2.AzureKeyVault{
			Enabled:           ear.AzureKeyVault.Enabled,
			AzureEnvironment:  ear.AzureKeyVault.GetAzureEnvironment(),
			ClientID:          ear.AzureKeyVault.GetClientID(),
			ResourceGroupName: ear.AzureKeyVault.GetResourceGroupName(),
			TenantID:          ear.AzureKeyVault.GetTenantID(),
		}
	}
	if ear.HasGoogleCloudKms() {
		out.GCP.GoogleCloudKms = akov2.GoogleCloudKms{
			Enabled: ear.GoogleCloudKms.Enabled,
		}
	}
	return out
}

func EqualSpecs(spec, atlas *EncryptionAtRest) bool {
	// Retracted secrets mean that this will never be equal
	return reflect.DeepEqual(atlas, spec)
}
