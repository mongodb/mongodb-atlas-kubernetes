package atlas

import (
	"context"
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/debug"
)

var globalAtlas *Atlas

type Atlas struct {
	OrgID   string
	Public  string
	Private string
	Client  *admin.APIClient
}

func AClient() (Atlas, error) {
	a := Atlas{
		OrgID:   os.Getenv("MCLI_ORG_ID"),
		Public:  os.Getenv("MCLI_PUBLIC_API_KEY"),
		Private: os.Getenv("MCLI_PRIVATE_API_KEY"),
	}

	c, err := atlas.NewClient(os.Getenv("MCLI_OPS_MANAGER_URL"), a.Public, a.Private)
	a.Client = c

	return a, err
}

func GetClientOrFail() *Atlas {
	if globalAtlas != nil {
		return globalAtlas
	}

	c, err := AClient()
	Expect(err).NotTo(HaveOccurred())

	globalAtlas = &c

	return globalAtlas
}

func (a *Atlas) IsDeploymentExist(projectID string, name string) bool {
	for _, c := range a.GetDeployments(projectID) {
		if c.GetName() == name {
			return true
		}
	}

	return false
}

func (a *Atlas) IsProjectExists(g Gomega, projectID string) bool {
	project, _, err := a.Client.ProjectsApi.GetProject(context.Background(), projectID).Execute()
	if admin.IsErrorCode(err, "GROUP_NOT_FOUND") || admin.IsErrorCode(err, "RESOURCE_NOT_FOUND") {
		return false
	}

	g.Expect(err).NotTo(HaveOccurred())

	return project != nil
}

func (a *Atlas) GetDeployments(projectID string) []admin.AdvancedClusterDescription {
	reply, _, err := a.Client.ClustersApi.ListClusters(context.Background(), projectID).Execute()
	Expect(err).NotTo(HaveOccurred())
	ginkgoPrettyPrintf(reply.Results, "listing legacy deployments in project %s", projectID)

	return *reply.Results
}

func (a *Atlas) GetDeployment(projectId, deploymentName string) (*admin.AdvancedClusterDescription, error) {
	advancedDeployment, _, err := a.Client.ClustersApi.
		GetCluster(context.Background(), projectId, deploymentName).
		Execute()

	ginkgoPrettyPrintf(advancedDeployment, "getting advanced deployment %s in project %s", deploymentName, projectId)

	return advancedDeployment, err
}

func (a *Atlas) GetServerlessInstance(projectId, deploymentName string) (*admin.ServerlessInstanceDescription, error) {
	serverlessInstance, _, err := a.Client.ServerlessInstancesApi.
		GetServerlessInstance(context.Background(), projectId, deploymentName).
		Execute()

	ginkgoPrettyPrintf(serverlessInstance, "getting serverless instance %s in project %s", deploymentName, projectId)

	return serverlessInstance, err
}

func (a *Atlas) GetDBUser(database, userName, projectID string) (*admin.CloudDatabaseUser, error) {
	user, _, err := a.Client.DatabaseUsersApi.
		GetDatabaseUser(context.Background(), projectID, database, userName).
		Execute()
	if admin.IsErrorCode(err, "USERNAME_NOT_FOUND") || admin.IsErrorCode(err, "RESOURCE_NOT_FOUND") || admin.IsErrorCode(err, "USER_NOT_IN_GROUP") {
		return nil, nil
	}

	return user, err
}

// ginkgoPrettyPrintf displays a message and a formatted json object through the Ginkgo Writer.
func ginkgoPrettyPrintf(obj interface{}, msg string, formatArgs ...interface{}) {
	ginkgo.GinkgoWriter.Println(fmt.Sprintf(msg, formatArgs...))
	ginkgo.GinkgoWriter.Println(debug.PrettyString(obj))
}

func (a *Atlas) GetIntegrationByType(projectId, iType string) (*admin.ThridPartyIntegration, error) {
	integration, _, err := a.Client.ThirdPartyIntegrationsApi.
		GetThirdPartyIntegration(context.Background(), projectId, iType).
		Execute()

	return integration, err
}

func (a *Atlas) GetUserByName(database, projectID, username string) (*admin.CloudDatabaseUser, error) {
	dbUser, _, err := a.Client.DatabaseUsersApi.
		GetDatabaseUser(context.Background(), projectID, database, username).
		Execute()
	if err != nil {
		return nil, err
	}

	return dbUser, nil
}

func (a *Atlas) DeleteGlobalKey(key admin.ApiKeyUserDetails) error {
	_, _, err := a.Client.ProgrammaticAPIKeysApi.DeleteApiKey(context.Background(), a.OrgID, key.GetId()).Execute()

	return err
}

func (a *Atlas) GetEncryptionAtRest(projectID string) (*admin.EncryptionAtRest, error) {
	encryptionAtRest, _, err := a.Client.EncryptionAtRestUsingCustomerKeyManagementApi.
		GetEncryptionAtRest(context.Background(), projectID).
		Execute()

	return encryptionAtRest, err
}

func (a *Atlas) GetOrgUsers() ([]admin.CloudAppUser, error) {
	users, _, err := a.Client.OrganizationsApi.ListOrganizationUsers(context.Background(), a.OrgID).Execute()

	return *users.Results, err
}

func (a *Atlas) CreateExportBucket(projectID, bucketName, roleID string) (*admin.DiskBackupSnapshotAWSExportBucket, error) {
	r, _, err := a.Client.CloudBackupsApi.
		CreateExportBucket(
			context.Background(),
			projectID,
			&admin.DiskBackupSnapshotAWSExportBucket{
				BucketName:    &bucketName,
				CloudProvider: pointer.MakePtr("AWS"),
				IamRoleId:     &roleID,
			},
		).Execute()

	return r, err
}
