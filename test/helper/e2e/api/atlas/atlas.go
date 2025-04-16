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

package atlas

import (
	"context"
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/debug"
)

var globalAtlas *Atlas

type Atlas struct {
	OrgID              string
	Public             string
	Private            string
	Client             *admin.APIClient
	ClientV20241113001 *adminv20241113001.APIClient
}

func AClient() (Atlas, error) {
	a := Atlas{
		OrgID:   os.Getenv("MCLI_ORG_ID"),
		Public:  os.Getenv("MCLI_PUBLIC_API_KEY"),
		Private: os.Getenv("MCLI_PRIVATE_API_KEY"),
	}

	c, err := atlas.NewClient(os.Getenv("MCLI_OPS_MANAGER_URL"), a.Public, a.Private)
	if err != nil {
		return Atlas{}, err
	}
	a.Client = c

	clientV20241113001, err := adminv20241113001.NewClient(
		adminv20241113001.UseBaseURL(os.Getenv("MCLI_OPS_MANAGER_URL")),
		adminv20241113001.UseDigestAuth(a.Public, a.Private),
	)
	if err != nil {
		return Atlas{}, err
	}
	a.ClientV20241113001 = clientV20241113001

	return a, nil
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
	for _, deploymentName := range a.GetDeploymentNames(projectID) {
		if deploymentName == name {
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

func (a *Atlas) GetDeploymentNames(projectID string) []string {
	ctx := context.Background()

	clusters, _, err := a.Client.ClustersApi.ListClusters(ctx, projectID).Execute()
	Expect(err).NotTo(HaveOccurred())
	ginkgoPrettyPrintf(clusters.GetResults(), "listing legacy deployments in project %s", projectID)
	names := []string{}
	for _, cluster := range clusters.GetResults() {
		names = append(names, cluster.GetName())
	}
	flex, _, err := a.ClientV20241113001.FlexClustersApi.ListFlexClusters(ctx, projectID).Execute()
	Expect(err).NotTo(HaveOccurred())
	ginkgoPrettyPrintf(flex.GetResults(), "listing flex deployments in project %s", projectID)
	for _, cluster := range flex.GetResults() {
		names = append(names, cluster.GetName())
	}
	return names
}

func (a *Atlas) GetDeployment(projectId, deploymentName string) (*admin.AdvancedClusterDescription, error) {
	advancedDeployment, _, err := a.Client.ClustersApi.
		GetCluster(context.Background(), projectId, deploymentName).
		Execute()

	ginkgoPrettyPrintf(advancedDeployment, "getting advanced deployment %s in project %s", deploymentName, projectId)

	return advancedDeployment, err
}

func (a *Atlas) GetFlexInstance(projectId, deploymentName string) (*adminv20241113001.FlexClusterDescription20241113, error) {
	flexInstance, _, err := a.ClientV20241113001.FlexClustersApi.
		GetFlexCluster(context.Background(), projectId, deploymentName).
		Execute()

	ginkgoPrettyPrintf(flexInstance, "getting flex instance %s in project %s", deploymentName, projectId)

	return flexInstance, err
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

func (a *Atlas) GetIntegrationByType(projectId, iType string) (*admin.ThirdPartyIntegration, error) {
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
