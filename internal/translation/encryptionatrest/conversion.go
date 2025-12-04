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

package encryptionatrest

import (
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
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
		Region:  pointer.MakePtrOrNil(a.AwsKms.Region),
		Valid:   a.AwsKms.Valid,

		RoleId:              pointer.MakePtrOrNil(a.RoleID),
		CustomerMasterKeyID: pointer.MakePtrOrNil(a.CustomerMasterKeyID),
	}

	if result.RoleId == nil && a.CloudProviderIntegrationRole != "" {
		result.RoleId = pointer.MakePtrOrNil(a.CloudProviderIntegrationRole)
	}

	if result.Enabled == nil {
		result.Enabled = pointer.MakePtr(false)
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

	result := &admin.AzureKeyVault{
		Enabled:           az.AzureKeyVault.Enabled,
		ClientID:          pointer.MakePtrOrNil(az.AzureKeyVault.ClientID),
		AzureEnvironment:  pointer.MakePtrOrNil(az.AzureKeyVault.AzureEnvironment),
		ResourceGroupName: pointer.MakePtrOrNil(az.AzureKeyVault.ResourceGroupName),
		TenantID:          pointer.MakePtrOrNil(az.AzureKeyVault.TenantID),

		SubscriptionID: pointer.MakePtrOrNil(az.SubscriptionID),
		KeyVaultName:   pointer.MakePtrOrNil(az.KeyVaultName),
		KeyIdentifier:  pointer.MakePtrOrNil(az.KeyIdentifier),
		Secret:         pointer.MakePtrOrNil(az.Secret),
	}

	if result.Enabled == nil {
		result.Enabled = pointer.MakePtr(false)
	}

	return result
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

	result := &admin.GoogleCloudKMS{
		Enabled:              g.GoogleCloudKms.Enabled,
		ServiceAccountKey:    pointer.MakePtrOrNil(g.ServiceAccountKey),
		KeyVersionResourceID: pointer.MakePtrOrNil(g.KeyVersionResourceID),
	}

	if result.Enabled == nil {
		result.Enabled = pointer.MakePtr(false)
	}

	return result
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
	result := &admin.EncryptionAtRest{
		AwsKms:         spec.AWS.ToAtlas(),
		AzureKeyVault:  spec.Azure.ToAtlas(),
		GoogleCloudKms: spec.GCP.ToAtlas(),
	}
	return result
}

func fromAtlas(ear *admin.EncryptionAtRest) *EncryptionAtRest {
	out := &EncryptionAtRest{}
	out.AWS.AwsKms.Enabled = pointer.MakePtr(false)
	if ear.HasAwsKms() {
		out.AWS.AwsKms = akov2.AwsKms{
			Enabled: ear.AwsKms.Enabled,
			Region:  ear.AwsKms.GetRegion(),
			Valid:   ear.AwsKms.Valid,
		}
	}

	out.Azure.AzureKeyVault.Enabled = pointer.MakePtr(false)
	if ear.HasAzureKeyVault() {
		out.Azure.AzureKeyVault = akov2.AzureKeyVault{
			Enabled:           ear.AzureKeyVault.Enabled,
			AzureEnvironment:  ear.AzureKeyVault.GetAzureEnvironment(),
			ClientID:          ear.AzureKeyVault.GetClientID(),
			ResourceGroupName: ear.AzureKeyVault.GetResourceGroupName(),
			TenantID:          ear.AzureKeyVault.GetTenantID(),
		}
	}

	out.GCP.GoogleCloudKms.Enabled = pointer.MakePtr(false)
	if ear.HasGoogleCloudKms() {
		out.GCP.GoogleCloudKms = akov2.GoogleCloudKms{
			Enabled: ear.GoogleCloudKms.Enabled,
		}
	}
	return out
}

func EqualSpecs(spec, atlas *EncryptionAtRest) bool {
	var specCopy, atlasCopy EncryptionAtRest

	specCopy.AWS = prunedAWSSpecCopy(spec)
	specCopy.GCP = prunedGCPSpecCopy(spec)
	specCopy.Azure = prunedAzureSpecCopy(spec)

	atlasCopy.AWS = prunedAWSSpecCopy(atlas)
	atlasCopy.GCP = prunedGCPSpecCopy(atlas)
	atlasCopy.Azure = prunedAzureSpecCopy(atlas)

	setDefaultsFromAtlas(&specCopy, &atlasCopy)
	return reflect.DeepEqual(specCopy, atlasCopy)
}

func prunedAWSSpecCopy(source *EncryptionAtRest) AwsKms {
	var result AwsKms

	if source == nil {
		return result
	}

	result.AwsKms = *source.AWS.AwsKms.DeepCopy()
	result.AwsKms.SecretRef = common.ResourceRefNamespaced{}
	if result.AwsKms.Enabled == nil {
		result.AwsKms.Enabled = pointer.MakePtr(false)
	}
	return result
}

func prunedGCPSpecCopy(source *EncryptionAtRest) GoogleCloudKms {
	var result GoogleCloudKms

	if source == nil {
		return result
	}

	result.GoogleCloudKms = *source.GCP.GoogleCloudKms.DeepCopy()
	result.GoogleCloudKms.SecretRef = common.ResourceRefNamespaced{}
	if result.GoogleCloudKms.Enabled == nil {
		result.GoogleCloudKms.Enabled = pointer.MakePtr(false)
	}
	return result
}

func prunedAzureSpecCopy(source *EncryptionAtRest) AzureKeyVault {
	var result AzureKeyVault

	if source == nil {
		return result
	}

	result.AzureKeyVault = *source.Azure.AzureKeyVault.DeepCopy()
	result.AzureKeyVault.SecretRef = common.ResourceRefNamespaced{}
	if result.AzureKeyVault.Enabled == nil {
		result.AzureKeyVault.Enabled = pointer.MakePtr(false)
	}
	return result
}

func setDefaultsFromAtlas(spec, atlas *EncryptionAtRest) {
	if spec.AWS.Valid == nil && atlas.AWS.Valid != nil {
		spec.AWS.Valid = atlas.AWS.Valid
	}

	if spec.AWS.Enabled == nil && atlas.AWS.Enabled != nil {
		spec.AWS.Enabled = atlas.AWS.Enabled
	}

	if spec.GCP.Enabled == nil && atlas.GCP.Enabled != nil {
		spec.GCP.Enabled = atlas.GCP.Enabled
	}

	if spec.Azure.Enabled == nil && atlas.Azure.Enabled != nil {
		spec.Azure.Enabled = atlas.Azure.Enabled
	}
}
