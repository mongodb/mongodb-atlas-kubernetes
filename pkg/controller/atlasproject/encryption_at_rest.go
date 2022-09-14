package atlasproject

import (
	"context"
	"reflect"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"go.mongodb.org/atlas/mongodbatlas"
)

func ensureEncryptionAtRest(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	result := createOrDeleteEncryptionAtRests(ctx, projectID, project)
	if !result.IsOk() {
		ctx.SetConditionFromResult(status.EncryptionAtRestReadyType, result)
		return result
	}

	if isEncryptionSpecEmpty(project.Spec.EncryptionAtRest) {
		ctx.UnsetCondition(status.EncryptionAtRestReadyType)
		return workflow.OK()
	}

	ctx.SetConditionTrue(status.EncryptionAtRestReadyType)
	return workflow.OK()
}

func createOrDeleteEncryptionAtRests(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	encryptionAtRestsInAtlas, err := fetchEncryptionAtRests(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	inSync, err := atlasInSync(encryptionAtRestsInAtlas, project.Spec.EncryptionAtRest)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	if inSync {
		return workflow.OK()
	}

	if err := syncEncryptionAtRestsInAtlas(ctx, projectID, project, encryptionAtRestsInAtlas); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	return workflow.InProgress(workflow.ProjectEncryptionAtRestReady, "Encryption at Rest is being synced")
}

func fetchEncryptionAtRests(ctx *workflow.Context, projectID string) (*mongodbatlas.EncryptionAtRest, error) {
	encryptionAtRestsInAtlas, _, err := ctx.Client.EncryptionsAtRest.Get(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got EncryptionAtRests From Atlas: %v", *encryptionAtRestsInAtlas)
	return encryptionAtRestsInAtlas, nil
}

func syncEncryptionAtRestsInAtlas(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject, atlas *mongodbatlas.EncryptionAtRest) error {
	requestBody := mongodbatlas.EncryptionAtRest{
		GroupID:        projectID,
		AwsKms:         getAwsKMS(project, atlas),
		AzureKeyVault:  getAzureKeyVault(project, atlas),
		GoogleCloudKms: getGoogleCloudKms(project, atlas),
	}

	if _, _, err := ctx.Client.EncryptionsAtRest.Create(context.Background(), &requestBody); err != nil { // Create() sends PATCH request
		return err
	}

	return nil
}

func atlasInSync(atlas *mongodbatlas.EncryptionAtRest, spec *mdbv1.EncryptionAtRest) (bool, error) {
	if IsEncryptionAtlasEmpty(atlas) && isEncryptionSpecEmpty(spec) {
		return true, nil
	}

	if IsEncryptionAtlasEmpty(atlas) || isEncryptionSpecEmpty(spec) {
		return false, nil
	}

	specAsAtlas, err := spec.ToAtlas(atlas.GroupID)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(atlas, specAsAtlas), nil
}

func isEncryptionSpecEmpty(spec *mdbv1.EncryptionAtRest) bool {
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

func getAwsKMS(project *mdbv1.AtlasProject, atlas *mongodbatlas.EncryptionAtRest) (result mongodbatlas.AwsKms) {
	if project.Spec.EncryptionAtRest != nil {
		result = mongodbatlas.AwsKms(project.Spec.EncryptionAtRest.AwsKms)
	}

	if (atlas == nil || atlas.AwsKms == mongodbatlas.AwsKms{}) {
		result.Enabled = toptr.MakePtr(false)
	}

	if result.RoleID == "" {
		awsRole, foundRole := selectRole(project.Status.CloudProviderAccessRoles, "AWS")
		if foundRole {
			result.RoleID = awsRole.RoleID
		}
	}

	return
}

func getAzureKeyVault(project *mdbv1.AtlasProject, atlas *mongodbatlas.EncryptionAtRest) (result mongodbatlas.AzureKeyVault) {
	if project.Spec.EncryptionAtRest != nil {
		result = mongodbatlas.AzureKeyVault(project.Spec.EncryptionAtRest.AzureKeyVault)
	}

	if (atlas == nil || atlas.GoogleCloudKms == mongodbatlas.GoogleCloudKms{}) {
		result.Enabled = toptr.MakePtr(false)
	}

	return
}

func getGoogleCloudKms(project *mdbv1.AtlasProject, atlas *mongodbatlas.EncryptionAtRest) (result mongodbatlas.GoogleCloudKms) {
	if project.Spec.EncryptionAtRest != nil {
		result = mongodbatlas.GoogleCloudKms(project.Spec.EncryptionAtRest.GoogleCloudKms)
	}

	if (atlas == nil || atlas.GoogleCloudKms == mongodbatlas.GoogleCloudKms{}) {
		result.Enabled = toptr.MakePtr(false)
	}

	return
}

func selectRole(accessRoles []status.CloudProviderAccessRole, providerName string) (result status.CloudProviderAccessRole, found bool) {
	for _, role := range accessRoles {
		if role.ProviderName == providerName {
			return role, true
		}
	}

	return
}
