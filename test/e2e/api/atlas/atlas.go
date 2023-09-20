package atlas

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils/debug"

	"github.com/mongodb-forks/digest"
	"go.mongodb.org/atlas/mongodbatlas"
)

var globalAtlas *Atlas

type Atlas struct {
	OrgID   string
	Public  string
	Private string
	Client  *mongodbatlas.Client
}

const (
	keyDescription = "created from the AO test"
)

func AClient() (Atlas, error) {
	var A Atlas
	A.OrgID = os.Getenv("MCLI_ORG_ID")
	A.Public = os.Getenv("MCLI_PUBLIC_API_KEY")
	A.Private = os.Getenv("MCLI_PRIVATE_API_KEY")
	t := digest.NewTransport(A.Public, A.Private)
	tc, err := t.Client()
	if err != nil {
		return A, err
	}
	A.Client = mongodbatlas.NewClient(tc)
	u, _ := url.Parse(os.Getenv("MCLI_OPS_MANAGER_URL"))
	A.Client.BaseURL = u
	return A, nil
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

func (a *Atlas) AddKeyWithAccessList(projectID string, roles, access []string) (*mongodbatlas.APIKey, error) {
	createKeyRequest := &mongodbatlas.APIKeyInput{
		Desc:  keyDescription,
		Roles: roles,
	}
	newKey, _, err := a.Client.ProjectAPIKeys.Create(context.Background(), projectID, createKeyRequest)
	if err != nil {
		return nil, err
	}
	createAccessRequest := formAccessRequest(access)
	_, _, err = a.Client.AccessListAPIKeys.Create(context.Background(), a.OrgID, newKey.ID, createAccessRequest)
	if err != nil {
		return nil, err
	}
	return newKey, nil
}

func formAccessRequest(access []string) []*mongodbatlas.AccessListAPIKeysReq {
	createRequest := make([]*mongodbatlas.AccessListAPIKeysReq, 0)
	var req *mongodbatlas.AccessListAPIKeysReq
	for _, item := range access {
		if strings.Contains(item, "/") {
			req = &mongodbatlas.AccessListAPIKeysReq{CidrBlock: item}
		} else {
			req = &mongodbatlas.AccessListAPIKeysReq{IPAddress: item}
		}
		createRequest = append(createRequest, req)
	}
	return createRequest
}

func (a *Atlas) GetPrivateEndpoint(projectID, provider string) ([]mongodbatlas.PrivateEndpointConnection, error) {
	endpointsList, _, err := a.Client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}
	ginkgoPrettyPrintf(endpointsList, "listing private endpoints in project %s", projectID)
	return endpointsList, nil
}

func (a *Atlas) IsDeploymentExist(projectID string, name string) bool {
	for _, c := range a.GetDeployments(projectID) {
		if c.Name == name {
			return true
		}
	}
	return false
}

func (a *Atlas) IsProjectExists(g Gomega, projectID string) bool {
	project, _, err := a.Client.Projects.GetOneProject(context.Background(), projectID)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && (apiError.ErrorCode == "GROUP_NOT_FOUND" || apiError.ErrorCode == "RESOURCE_NOT_FOUND") {
			return false
		}
		g.Expect(err).NotTo(HaveOccurred())
	}
	return project != nil
}

func (a *Atlas) GetDeployments(projectID string) []*mongodbatlas.AdvancedCluster {
	reply, _, err := a.Client.AdvancedClusters.List(context.Background(), projectID, nil)
	Expect(err).NotTo(HaveOccurred())
	deployments := reply.Results
	ginkgoPrettyPrintf(deployments, "listing legacy deployments in project %s", projectID)
	return deployments
}

func (a *Atlas) GetDeployment(projectId, deploymentName string) (*mongodbatlas.AdvancedCluster, error) {
	advancedDeployment, _, err := a.Client.AdvancedClusters.Get(context.Background(), projectId, deploymentName)
	if err != nil {
		return nil, err
	}
	ginkgoPrettyPrintf(advancedDeployment, "getting advanced deployment %s in project %s", deploymentName, projectId)
	return advancedDeployment, nil
}

func (a *Atlas) GetServerlessInstance(projectId, deploymentName string) (*mongodbatlas.Cluster, error) {
	serverlessInstance, _, err := a.Client.ServerlessInstances.Get(context.Background(), projectId, deploymentName)
	if err != nil {
		return nil, err
	}
	ginkgoPrettyPrintf(serverlessInstance, "getting serverless instance %s in project %s", deploymentName, projectId)
	return serverlessInstance, nil
}

func (a *Atlas) GetDBUser(database, userName, projectID string) (*mongodbatlas.DatabaseUser, error) {
	user, _, err := a.Client.DatabaseUsers.Get(context.Background(), database, projectID, userName)
	if err != nil {
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) &&
				(apiError.ErrorCode == "USERNAME_NOT_FOUND" || apiError.ErrorCode == "RESOURCE_NOT_FOUND" || apiError.ErrorCode == "USER_NOT_IN_GROUP" || apiError.Response.StatusCode == 400) {
				return nil, nil
			}
			return nil, err
		}
	}
	return user, nil
}

// ginkgoPrettyPrintf displays a message and a formatted json object through the Ginkgo Writer.
func ginkgoPrettyPrintf(obj interface{}, msg string, formatArgs ...interface{}) {
	ginkgo.GinkgoWriter.Println(fmt.Sprintf(msg, formatArgs...))
	ginkgo.GinkgoWriter.Println(debug.PrettyString(obj))
}

func (a *Atlas) GetIntegrationByType(projectId, iType string) (*mongodbatlas.ThirdPartyIntegration, error) {
	integration, _, err := a.Client.Integrations.Get(context.Background(), projectId, iType)
	if err != nil {
		return nil, err
	}
	return integration, nil
}

func (a *Atlas) GetUserByName(database, projectID, username string) (*mongodbatlas.DatabaseUser, error) {
	dbUser, _, err := a.Client.DatabaseUsers.Get(context.Background(), database, projectID, username)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

func (a *Atlas) DeleteGlobalKey(key mongodbatlas.APIKey) error {
	_, err := a.Client.APIKeys.Delete(context.Background(), a.OrgID, key.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Atlas) GetEncryptionAtRest(projectID string) (*mongodbatlas.EncryptionAtRest, error) {
	encryptionAtRest, _, err := a.Client.EncryptionsAtRest.Get(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	return encryptionAtRest, nil
}

func (a *Atlas) GetOrgUsers(projectID string) ([]mongodbatlas.AtlasUser, error) {
	users, _, err := a.Client.AtlasUsers.List(context.Background(), projectID, nil)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a *Atlas) CreateExportBucket(projectID, bucketName, roleID string) (*mongodbatlas.CloudProviderSnapshotExportBucket, error) {
	r, _, err := a.Client.CloudProviderSnapshotExportBuckets.Create(
		context.Background(),
		projectID,
		&mongodbatlas.CloudProviderSnapshotExportBucket{
			BucketName:    bucketName,
			CloudProvider: "AWS",
			IAMRoleID:     roleID,
		},
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}
