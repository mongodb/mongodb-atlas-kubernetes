package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"

	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	ObjectIDRegex = "^([a-f0-9]{24})$"
)

func (r *AtlasProjectReconciler) ensureEncryptionAtRest(workflowCtx *workflow.Context, project *mdbv1.AtlasProject, protected bool) workflow.Result {
	if err := readEncryptionAtRestSecrets(r.Client, workflowCtx, project.Spec.EncryptionAtRest, project.Namespace); err != nil {
		workflowCtx.UnsetCondition(status.EncryptionAtRestReadyType)
		return workflow.Terminate(workflow.ProjectEncryptionAtRestReady, err.Error())
	}

	canReconcile, err := canEncryptionAtRestReconcile(workflowCtx, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.EncryptionAtRestReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.EncryptionAtRestReadyType, result)

		return result
	}

	result := createOrDeleteEncryptionAtRests(workflowCtx, project.ID(), project)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.EncryptionAtRestReadyType, result)
		return result
	}

	if IsEncryptionSpecEmpty(project.Spec.EncryptionAtRest) {
		workflowCtx.UnsetCondition(status.EncryptionAtRestReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(status.EncryptionAtRestReadyType)
	return workflow.OK()
}

func readEncryptionAtRestSecrets(kubeClient client.Client, service *workflow.Context, encRest *mdbv1.EncryptionAtRest, parentNs string) error {
	if encRest == nil {
		return nil
	}

	if encRest.AwsKms.Enabled != nil && *encRest.AwsKms.Enabled && encRest.AwsKms.SecretRef.Name != "" {
		watchObj, err := readAndFillAWSSecret(kubeClient, parentNs, &encRest.AwsKms)
		service.AddResourcesToWatch(*watchObj)
		if err != nil {
			return err
		}
	}

	if encRest.GoogleCloudKms.Enabled != nil && *encRest.GoogleCloudKms.Enabled && encRest.GoogleCloudKms.SecretRef.Name != "" {
		watchObj, err := readAndFillGoogleSecret(kubeClient, parentNs, &encRest.GoogleCloudKms)
		service.AddResourcesToWatch(*watchObj)
		if err != nil {
			return err
		}
	}

	if encRest.AzureKeyVault.Enabled != nil && *encRest.AzureKeyVault.Enabled && encRest.AzureKeyVault.SecretRef.Name != "" {
		watchObj, err := readAndFillAzureSecret(kubeClient, parentNs, &encRest.AzureKeyVault)
		service.AddResourcesToWatch(*watchObj)
		if err != nil {
			return err
		}
	}

	return nil
}

func readAndFillAWSSecret(kubeClient client.Client, parentNs string, awsKms *mdbv1.AwsKms) (*watch.WatchedObject, error) {
	fieldData, watchObj, err := readSecretData(kubeClient, awsKms.SecretRef, parentNs, "CustomerMasterKeyID", "RoleID")
	if err != nil {
		return watchObj, err
	}

	awsKms.SetSecrets(fieldData["CustomerMasterKeyID"], fieldData["RoleID"])

	return watchObj, nil
}

func readAndFillGoogleSecret(kubeClient client.Client, parentNs string, gkms *mdbv1.GoogleCloudKms) (*watch.WatchedObject, error) {
	fieldData, watchObj, err := readSecretData(kubeClient, gkms.SecretRef, parentNs, "ServiceAccountKey", "KeyVersionResourceID")
	if err != nil {
		return watchObj, err
	}

	gkms.SetSecrets(fieldData["ServiceAccountKey"], fieldData["KeyVersionResourceID"])

	return watchObj, nil
}

func readAndFillAzureSecret(kubeClient client.Client, parentNs string, azureVault *mdbv1.AzureKeyVault) (*watch.WatchedObject, error) {
	fieldData, watchObj, err := readSecretData(kubeClient, azureVault.SecretRef, parentNs, "Secret", "SubscriptionID", "KeyVaultName", "KeyIdentifier")
	if err != nil {
		return watchObj, err
	}

	azureVault.SetSecrets(fieldData["SubscriptionID"], fieldData["KeyVaultName"], fieldData["KeyIdentifier"], fieldData["Secret"])

	return watchObj, nil
}

// Return all requested field from a secret
func readSecretData(kubeClient client.Client, res common.ResourceRefNamespaced, parentNamespace string, fieldNames ...string) (map[string]string, *watch.WatchedObject, error) {
	secret := &v1.Secret{}
	var ns string
	if res.Namespace == "" {
		ns = parentNamespace
	} else {
		ns = res.Namespace
	}

	result := map[string]string{}

	secretObj := client.ObjectKey{Name: res.Name, Namespace: ns}
	obj := &watch.WatchedObject{ResourceKind: "Secret", Resource: secretObj}

	if err := kubeClient.Get(context.Background(), secretObj, secret); err != nil {
		return result, obj, err
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
		return result, obj, fmt.Errorf("the following fields are either missing or their values are empty: %s", strings.Join(missingFields, ", "))
	}

	return result, obj, nil
}

func createOrDeleteEncryptionAtRests(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
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
	encryptionAtRestsInAtlas, _, err := ctx.Client.EncryptionsAtRest.Get(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got EncryptionAtRests From Atlas: %v", *encryptionAtRestsInAtlas)
	return encryptionAtRestsInAtlas, nil
}

func syncEncryptionAtRestsInAtlas(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) error {
	requestBody := mongodbatlas.EncryptionAtRest{
		GroupID:        projectID,
		AwsKms:         getAwsKMS(project),
		AzureKeyVault:  getAzureKeyVault(project),
		GoogleCloudKms: getGoogleCloudKms(project),
	}

	if err := normalizeAwsKms(ctx, projectID, &requestBody.AwsKms); err != nil {
		return err
	}

	if _, _, err := ctx.Client.EncryptionsAtRest.Create(context.Background(), &requestBody); err != nil { // Create() sends PATCH request
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
	resp, _, err := ctx.Client.CloudProviderAccess.ListRoles(context.Background(), projectID)
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

func AtlasInSync(atlas *mongodbatlas.EncryptionAtRest, spec *mdbv1.EncryptionAtRest) (bool, error) {
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
		spec.AwsKms.Enabled = toptr.MakePtr(false)
	}
	if isNotNilAndFalse(atlas.AzureKeyVault.Enabled) {
		spec.AzureKeyVault.Enabled = toptr.MakePtr(false)
	}
	if isNotNilAndFalse(atlas.GoogleCloudKms.Enabled) {
		spec.GoogleCloudKms.Enabled = toptr.MakePtr(false)
	}

	spec.Valid = atlas.Valid
}

func IsEncryptionSpecEmpty(spec *mdbv1.EncryptionAtRest) bool {
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

func getAwsKMS(project *mdbv1.AtlasProject) (result mongodbatlas.AwsKms) {
	if project.Spec.EncryptionAtRest == nil {
		return
	}

	result = project.Spec.EncryptionAtRest.AwsKms.ToAtlas()

	if (result == mongodbatlas.AwsKms{}) {
		result.Enabled = toptr.MakePtr(false)
	}

	if result.RoleID == "" {
		awsRole, foundRole := selectRole(project.Status.CloudProviderIntegrations, "AWS")
		if foundRole {
			result.RoleID = awsRole.RoleID
		}
	}

	return
}

func getAzureKeyVault(project *mdbv1.AtlasProject) (result mongodbatlas.AzureKeyVault) {
	if project.Spec.EncryptionAtRest == nil {
		return
	}

	result = project.Spec.EncryptionAtRest.AzureKeyVault.ToAtlas()

	if (result == mongodbatlas.AzureKeyVault{}) {
		result.Enabled = toptr.MakePtr(false)
	}

	return
}

func getGoogleCloudKms(project *mdbv1.AtlasProject) (result mongodbatlas.GoogleCloudKms) {
	if project.Spec.EncryptionAtRest == nil {
		return
	}

	result = project.Spec.EncryptionAtRest.GoogleCloudKms.ToAtlas()

	if (result == mongodbatlas.GoogleCloudKms{}) {
		result.Enabled = toptr.MakePtr(false)
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

func canEncryptionAtRestReconcile(workflowCtx *workflow.Context, protected bool, akoProject *mdbv1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &mdbv1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	ear, _, err := workflowCtx.Client.EncryptionsAtRest.Get(workflowCtx.Context, akoProject.ID())
	if err != nil {
		return false, err
	}

	if IsEncryptionAtlasEmpty(ear) {
		return true, nil
	}

	return areEaRConfigEqual(*latestConfig.EncryptionAtRest, ear, true) ||
		areEaRConfigEqual(*akoProject.Spec.EncryptionAtRest, ear, false), nil
}

func areEaRConfigEqual(operator mdbv1.EncryptionAtRest, atlas *mongodbatlas.EncryptionAtRest, lastApplied bool) bool {
	return areAWSConfigEqual(operator.AwsKms, atlas.AwsKms, lastApplied) &&
		areGCPConfigEqual(operator.GoogleCloudKms, atlas.GoogleCloudKms, lastApplied) &&
		areAzureConfigEqual(operator.AzureKeyVault, atlas.AzureKeyVault, lastApplied)
}

func areAWSConfigEqual(operator mdbv1.AwsKms, atlas mongodbatlas.AwsKms, lastApplied bool) bool {
	if operator.Enabled == nil {
		operator.Enabled = toptr.MakePtr(false)
	}

	if lastApplied {
		return *operator.Enabled == *atlas.Enabled &&
			operator.Region == atlas.Region
	}
	// Atlas API does not return SecretAccessKey or RoleID
	return *operator.Enabled == *atlas.Enabled &&
		operator.CustomerMasterKeyID() == atlas.CustomerMasterKeyID &&
		operator.Region == atlas.Region
}

func areGCPConfigEqual(operator mdbv1.GoogleCloudKms, atlas mongodbatlas.GoogleCloudKms, lastApplied bool) bool {
	if operator.Enabled == nil {
		operator.Enabled = toptr.MakePtr(false)
	}

	if lastApplied {
		return *operator.Enabled == *atlas.Enabled
	}

	// Atlas API does not return service account key
	return *operator.Enabled == *atlas.Enabled &&
		operator.KeyVersionResourceID() == atlas.KeyVersionResourceID
}

func areAzureConfigEqual(operator mdbv1.AzureKeyVault, atlas mongodbatlas.AzureKeyVault, lastApplied bool) bool {
	if operator.Enabled == nil {
		operator.Enabled = toptr.MakePtr(false)
	}

	if lastApplied {
		return *operator.Enabled == *atlas.Enabled &&
			operator.AzureEnvironment == atlas.AzureEnvironment &&
			operator.ClientID == atlas.ClientID &&
			operator.ResourceGroupName == atlas.ResourceGroupName &&
			operator.TenantID == atlas.TenantID
	}

	return *operator.Enabled == *atlas.Enabled &&
		operator.AzureEnvironment == atlas.AzureEnvironment &&
		operator.ClientID == atlas.ClientID &&
		operator.KeyIdentifier() == atlas.KeyIdentifier &&
		operator.KeyVaultName() == atlas.KeyVaultName &&
		operator.ResourceGroupName == atlas.ResourceGroupName &&
		operator.SubscriptionID() == atlas.SubscriptionID &&
		operator.TenantID == atlas.TenantID
}
