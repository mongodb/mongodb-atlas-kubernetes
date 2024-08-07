package atlasproject

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	ObjectIDRegex = "^([a-f0-9]{24})$"
)

func (r *AtlasProjectReconciler) ensureEncryptionAtRest(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	if err := readEncryptionAtRestSecrets(r.Client, workflowCtx, project.Spec.EncryptionAtRest, project.Namespace); err != nil {
		workflowCtx.UnsetCondition(api.EncryptionAtRestReadyType)
		return workflow.Terminate(workflow.ProjectEncryptionAtRestReady, err.Error())
	}

	result := createOrDeleteEncryptionAtRests(workflowCtx, project.ID(), project)
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

func readEncryptionAtRestSecrets(kubeClient client.Client, service *workflow.Context, encRest *akov2.EncryptionAtRest, parentNs string) error {
	if encRest == nil {
		return nil
	}

	if encRest.AwsKms.Enabled != nil && *encRest.AwsKms.Enabled && encRest.AwsKms.SecretRef.Name != "" {
		err := readAndFillAWSSecret(service.Context, kubeClient, parentNs, &encRest.AwsKms)
		if err != nil {
			return err
		}
	}

	if encRest.GoogleCloudKms.Enabled != nil && *encRest.GoogleCloudKms.Enabled && encRest.GoogleCloudKms.SecretRef.Name != "" {
		err := readAndFillGoogleSecret(service.Context, kubeClient, parentNs, &encRest.GoogleCloudKms)
		if err != nil {
			return err
		}
	}

	if encRest.AzureKeyVault.Enabled != nil && *encRest.AzureKeyVault.Enabled && encRest.AzureKeyVault.SecretRef.Name != "" {
		err := readAndFillAzureSecret(service.Context, kubeClient, parentNs, &encRest.AzureKeyVault)
		if err != nil {
			return err
		}
	}

	return nil
}

func readAndFillAWSSecret(ctx context.Context, kubeClient client.Client, parentNs string, awsKms *akov2.AwsKms) error {
	fieldData, err := readSecretData(ctx, kubeClient, awsKms.SecretRef, parentNs, "CustomerMasterKeyID", "RoleID")
	if err != nil {
		return err
	}

	awsKms.SetSecrets(fieldData["CustomerMasterKeyID"], fieldData["RoleID"])

	return nil
}

func readAndFillGoogleSecret(ctx context.Context, kubeClient client.Client, parentNs string, gkms *akov2.GoogleCloudKms) error {
	fieldData, err := readSecretData(ctx, kubeClient, gkms.SecretRef, parentNs, "ServiceAccountKey", "KeyVersionResourceID")
	if err != nil {
		return err
	}

	gkms.SetSecrets(fieldData["ServiceAccountKey"], fieldData["KeyVersionResourceID"])

	return nil
}

func readAndFillAzureSecret(ctx context.Context, kubeClient client.Client, parentNs string, azureVault *akov2.AzureKeyVault) error {
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

func createOrDeleteEncryptionAtRests(ctx *workflow.Context, projectID string, project *akov2.AtlasProject) workflow.Result {
	encryptionAtRestsInAtlas, err := fetchEncryptionAtRests(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	inSync, err := AtlasInSync(encryptionAtRestsInAtlas, project.Spec.EncryptionAtRest)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if inSync {
		return workflow.OK()
	}

	if err := syncEncryptionAtRestsInAtlas(ctx, projectID, project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	return workflow.OK()
}

func fetchEncryptionAtRests(ctx *workflow.Context, projectID string) (*mongodbatlas.EncryptionAtRest, error) {
	encryptionAtRestsInAtlas, _, err := ctx.Client.EncryptionsAtRest.Get(ctx.Context, projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got EncryptionAtRests From Atlas: %v", *encryptionAtRestsInAtlas)
	return encryptionAtRestsInAtlas, nil
}

func syncEncryptionAtRestsInAtlas(ctx *workflow.Context, projectID string, project *akov2.AtlasProject) error {
	requestBody := mongodbatlas.EncryptionAtRest{
		GroupID:        projectID,
		AwsKms:         getAwsKMS(project),
		AzureKeyVault:  getAzureKeyVault(project),
		GoogleCloudKms: getGoogleCloudKms(project),
	}

	if err := normalizeAwsKms(ctx, projectID, &requestBody.AwsKms); err != nil {
		return err
	}

	if _, _, err := ctx.Client.EncryptionsAtRest.Create(ctx.Context, &requestBody); err != nil { // Create() sends PATCH request
		return err
	}

	return nil
}

func normalizeAwsKms(ctx *workflow.Context, projectID string, awsKms *mongodbatlas.AwsKms) error {
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

	// assume that role ID is set as AWS ARN
	resp, _, err := ctx.Client.CloudProviderAccess.ListRoles(ctx.Context, projectID)
	if err != nil {
		return err
	}

	for _, role := range resp.AWSIAMRoles {
		if role.IAMAssumedRoleARN == awsKms.RoleID {
			awsKms.RoleID = role.RoleID
			return nil
		}
	}

	ctx.Log.Debugf("no match for provided AWS RoleID ARN: '%s'. Is the CPA configured for the project?", awsKms.RoleID)
	return fmt.Errorf("can not use '%s' aws roleID for encryption at rest. AWS ARN not configured as Cloud Provider Access", awsKms.RoleID)
}

func AtlasInSync(atlas *mongodbatlas.EncryptionAtRest, spec *akov2.EncryptionAtRest) (bool, error) {
	if IsEncryptionAtlasEmpty(atlas) && IsEncryptionSpecEmpty(spec) {
		return true, nil
	}

	if IsEncryptionAtlasEmpty(atlas) || IsEncryptionSpecEmpty(spec) {
		return false, nil
	}

	specAsAtlas, err := spec.ToAtlas(atlas.GroupID)
	if err != nil {
		return false, err
	}

	balanceAsymmetricalFields(atlas, specAsAtlas)

	return reflect.DeepEqual(atlas, specAsAtlas), nil
}

func balanceAsymmetricalFields(atlas *mongodbatlas.EncryptionAtRest, spec *mongodbatlas.EncryptionAtRest) {
	if spec.AwsKms.RoleID == "" && atlas.AwsKms.RoleID != "" {
		spec.AwsKms.RoleID = atlas.AwsKms.RoleID
	}
	if spec.AzureKeyVault.Secret == "" && atlas.AzureKeyVault.Secret != "" {
		spec.AzureKeyVault.Secret = atlas.AzureKeyVault.Secret
	}
	if spec.GoogleCloudKms.ServiceAccountKey == "" && atlas.GoogleCloudKms.ServiceAccountKey != "" {
		spec.GoogleCloudKms.ServiceAccountKey = ""
	}

	if isNotNilAndFalse(atlas.AwsKms.Enabled) {
		spec.AwsKms.Enabled = pointer.MakePtr(false)
	}
	if isNotNilAndFalse(atlas.AzureKeyVault.Enabled) {
		spec.AzureKeyVault.Enabled = pointer.MakePtr(false)
	}
	if isNotNilAndFalse(atlas.GoogleCloudKms.Enabled) {
		spec.GoogleCloudKms.Enabled = pointer.MakePtr(false)
	}

	spec.Valid = atlas.Valid
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

func IsEncryptionAtlasEmpty(atlas *mongodbatlas.EncryptionAtRest) bool {
	if atlas == nil {
		return true
	}

	awsEnabled := atlas.AwsKms.Enabled
	azureEnabled := atlas.AzureKeyVault.Enabled
	gcpEnabled := atlas.GoogleCloudKms.Enabled

	if isNotNilAndTrue(awsEnabled) || isNotNilAndTrue(azureEnabled) || isNotNilAndTrue(gcpEnabled) {
		return false
	}

	return true
}

func isNotNilAndTrue(val *bool) bool {
	return val != nil && *val
}

func isNotNilAndFalse(val *bool) bool {
	return val != nil && !*val
}

func getAwsKMS(project *akov2.AtlasProject) (result mongodbatlas.AwsKms) {
	if project.Spec.EncryptionAtRest == nil {
		return
	}

	result = project.Spec.EncryptionAtRest.AwsKms.ToAtlas()

	if (result == mongodbatlas.AwsKms{}) {
		result.Enabled = pointer.MakePtr(false)
	}

	if result.RoleID == "" {
		awsRole, foundRole := selectRole(project.Status.CloudProviderIntegrations, "AWS")
		if foundRole {
			result.RoleID = awsRole.RoleID
		}
	}

	return
}

func getAzureKeyVault(project *akov2.AtlasProject) (result mongodbatlas.AzureKeyVault) {
	if project.Spec.EncryptionAtRest == nil {
		return
	}

	result = project.Spec.EncryptionAtRest.AzureKeyVault.ToAtlas()

	if (result == mongodbatlas.AzureKeyVault{}) {
		result.Enabled = pointer.MakePtr(false)
	}

	return
}

func getGoogleCloudKms(project *akov2.AtlasProject) (result mongodbatlas.GoogleCloudKms) {
	if project.Spec.EncryptionAtRest == nil {
		return
	}

	result = project.Spec.EncryptionAtRest.GoogleCloudKms.ToAtlas()

	if (result == mongodbatlas.GoogleCloudKms{}) {
		result.Enabled = pointer.MakePtr(false)
	}

	return
}

func selectRole(accessRoles []status.CloudProviderIntegration, providerName string) (result status.CloudProviderIntegration, found bool) {
	for _, role := range accessRoles {
		if role.ProviderName == providerName {
			return role, true
		}
	}

	return
}
