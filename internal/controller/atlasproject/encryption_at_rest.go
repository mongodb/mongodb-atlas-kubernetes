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

package atlasproject

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/encryptionatrest"
)

const (
	ObjectIDRegex = "^([a-f0-9]{24})$"
)

func (r *AtlasProjectReconciler) ensureEncryptionAtRest(workflowCtx *workflow.Context, project *akov2.AtlasProject, encryptionAtRestService encryptionatrest.EncryptionAtRestService) workflow.DeprecatedResult {
	encRest := encryptionatrest.NewEncryptionAtRest(project)

	if err := readEncryptionAtRestSecrets(r.Client, workflowCtx, encRest, project.Namespace); err != nil {
		workflowCtx.UnsetCondition(api.EncryptionAtRestReadyType)
		return workflow.Terminate(workflow.ProjectEncryptionAtRestReady, err)
	}

	result := createOrDeleteEncryptionAtRests(workflowCtx, encryptionAtRestService, project.ID(), encRest)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.EncryptionAtRestReadyType, result)
		return result
	}

	if IsEncryptionSpecEmpty(project.Spec.EncryptionAtRest) {
		workflowCtx.UnsetCondition(api.EncryptionAtRestReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(api.EncryptionAtRestReadyType)
	return workflow.OK()
}

func readEncryptionAtRestSecrets(kubeClient client.Client, service *workflow.Context, encRest *encryptionatrest.EncryptionAtRest, parentNs string) error {
	if encRest == nil {
		return nil
	}

	if encRest.AWS.IsEnabled() && encRest.AWS.SecretRef.Name != "" {
		err := readAndFillAWSSecret(service.Context, kubeClient, parentNs, &encRest.AWS)
		if err != nil {
			return err
		}
	}

	if encRest.GCP.IsEnabled() && encRest.GCP.SecretRef.Name != "" {
		err := readAndFillGoogleSecret(service.Context, kubeClient, parentNs, &encRest.GCP)
		if err != nil {
			return err
		}
	}

	if encRest.Azure.IsEnabled() && encRest.Azure.SecretRef.Name != "" {
		err := readAndFillAzureSecret(service.Context, kubeClient, parentNs, &encRest.Azure)
		if err != nil {
			return err
		}
	}

	return nil
}

func readAndFillAWSSecret(ctx context.Context, kubeClient client.Client, parentNs string, awsKms *encryptionatrest.AwsKms) error {
	fieldData, err := readSecretData(ctx, kubeClient, awsKms.SecretRef, parentNs, "CustomerMasterKeyID", "RoleID")
	if err != nil {
		return err
	}

	awsKms.SetSecrets(fieldData["CustomerMasterKeyID"], fieldData["RoleID"])

	return nil
}

func readAndFillGoogleSecret(ctx context.Context, kubeClient client.Client, parentNs string, gkms *encryptionatrest.GoogleCloudKms) error {
	fieldData, err := readSecretData(ctx, kubeClient, gkms.SecretRef, parentNs, "ServiceAccountKey", "KeyVersionResourceID")
	if err != nil {
		return err
	}

	gkms.SetSecrets(fieldData["ServiceAccountKey"], fieldData["KeyVersionResourceID"])

	return nil
}

func readAndFillAzureSecret(ctx context.Context, kubeClient client.Client, parentNs string, azureVault *encryptionatrest.AzureKeyVault) error {
	fieldData, err := readSecretData(ctx, kubeClient, azureVault.SecretRef, parentNs, "Secret", "SubscriptionID", "KeyVaultName", "KeyIdentifier")
	if err != nil {
		return err
	}

	azureVault.SetSecrets(fieldData["SubscriptionID"], fieldData["KeyVaultName"], fieldData["KeyIdentifier"], fieldData["Secret"])

	return nil
}

// Return all requested field from a secret
func readSecretData(ctx context.Context, kubeClient client.Client, res common.ResourceRefNamespaced, parentNamespace string, fieldNames ...string) (map[string]string, error) {
	secret := &v1.Secret{}
	var ns string
	if res.Namespace == "" {
		ns = parentNamespace
	} else {
		ns = res.Namespace
	}

	result := map[string]string{}

	secretObj := client.ObjectKey{Name: res.Name, Namespace: ns}

	if err := kubeClient.Get(ctx, secretObj, secret); err != nil {
		return result, err
	}

	missingFields := []string{}
	for i := range fieldNames {
		val, exists := secret.Data[fieldNames[i]]
		if !exists || len(val) == 0 {
			missingFields = append(missingFields, fieldNames[i])
		}
		result[fieldNames[i]] = string(val)
	}

	if len(missingFields) != 0 {
		return result, fmt.Errorf("the following fields are either missing or their values are empty: %s", strings.Join(missingFields, ", "))
	}

	return result, nil
}

func createOrDeleteEncryptionAtRests(ctx *workflow.Context, service encryptionatrest.EncryptionAtRestService, projectID string, encRest *encryptionatrest.EncryptionAtRest) workflow.DeprecatedResult {
	encryptionAtRestsInAtlas, err := service.Get(ctx.Context, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	inSync := encryptionatrest.EqualSpecs(encRest, encryptionAtRestsInAtlas)

	if inSync {
		return workflow.OK()
	}

	if encRest != nil {
		if err := normalizeAwsKms(ctx, projectID, &encRest.AWS); err != nil {
			return workflow.Terminate(workflow.Internal, err)
		}
		if err := service.Update(ctx.Context, projectID, *encRest); err != nil {
			return workflow.Terminate(workflow.Internal, err)
		}
	}

	return workflow.OK()
}

func normalizeAwsKms(ctx *workflow.Context, projectID string, awsKms *encryptionatrest.AwsKms) error {
	if awsKms == nil || awsKms.Enabled == nil || !*awsKms.Enabled {
		return nil
	}

	// verify if role ID is set as AtlasObjectID
	matched, err := regexp.MatchString(ObjectIDRegex, awsKms.RoleID)
	if err != nil {
		ctx.Log.Debugf("normalizing aws kms roleID failed: %v", err)
		return err
	}
	if matched {
		return nil
	}

	// This endpoint does not offer paginated responses.
	// Assume that role ID is set as AWS ARN.
	resp, _, err := ctx.SdkClientSet.SdkClient20250312009.CloudProviderAccessApi.ListCloudProviderAccessRoles(ctx.Context, projectID).Execute()
	if err != nil {
		return err
	}

	for _, role := range resp.GetAwsIamRoles() {
		if role.GetIamAssumedRoleArn() == awsKms.RoleID {
			awsKms.RoleID = role.GetRoleId()
			return nil
		}
	}

	ctx.Log.Debugf("no match for provided AWS RoleID ARN: '%s'. Is the CPA configured for the project?", awsKms.RoleID)
	return fmt.Errorf("can not use '%s' aws roleID for encryption at rest. AWS ARN not configured as Cloud Provider Access", awsKms.RoleID)
}

func IsEncryptionSpecEmpty(spec *akov2.EncryptionAtRest) bool {
	if spec == nil {
		return true
	}

	awsEnabled := spec.AwsKms.Enabled
	azureEnabled := spec.AzureKeyVault.Enabled
	gcpEnabled := spec.GoogleCloudKms.Enabled

	if isNotNilAndTrue(awsEnabled) || isNotNilAndTrue(azureEnabled) || isNotNilAndTrue(gcpEnabled) {
		return false
	}

	return true
}

func isNotNilAndTrue(val *bool) bool {
	return val != nil && *val
}
