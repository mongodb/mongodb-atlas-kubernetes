package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

type DeployedStatus struct {
	Error          error  `json:"-"`
	ProjectID      string `json:"projectId"`
	DeploymentName string `json:"deploymentName"`
}

func LoadDeployedStatus(name string) *DeployedStatus {
	status := DeployedStatus{}
	jsonData, err := os.ReadFile(deployedStatusFilename(name))
	status.Error = err
	if status.Error != nil {
		return &status
	}
	status.Error = json.Unmarshal(jsonData, &status)
	if status.Error == nil {
		status.Error =
			CreateDeploymentInTime(context.Background(), status.ProjectID, status.DeploymentName, time.Second)
	}
	return &status
}

func StoreDeployedStatus(name string, status *DeployedStatus) error {
	jsonData, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return os.WriteFile(deployedStatusFilename(name), jsonData, 0600)
}

func deployedStatusFilename(name string) string {
	return fmt.Sprintf(".%s-deployed.json", name)
}

func MustCreateDefaultProject(ctx context.Context, prefix string) string {
	name := newRandomName(prefix)
	log.Printf("Creating default project %s...", name)
	project, _, err := MustCreateAPIClient().ProjectsApi.CreateProject(ctx, &admin.Group{
		Name:  name,
		OrgId: OrgID(),
	}).Execute()
	if err != nil {
		panic(err)
	}
	log.Printf("Created default project %s ID=%s", name, *project.Id)
	return *project.Id
}

func MustCreateDefaultDeployment(ctx context.Context, projectID, prefix string) string {
	name := newRandomName(prefix)
	log.Printf("Creating default deployment %s...", name)
	deployment, _, err := MustCreateAPIClient().ServerlessInstancesApi.CreateServerlessInstance(
		ctx,
		projectID,
		&admin.ServerlessInstanceDescriptionCreate{
			Name: name,
			ProviderSettings: admin.ServerlessProviderSettings{
				BackingProviderName: DefaultProviderName(),
				RegionName:          DefaultRegion(),
			},
		}).Execute()
	if err != nil {
		panic(err)
	}
	log.Printf("Created default deployment %s ID=%s", name, *deployment.Id)
	return name
}

func CreateDeploymentInTime(ctx context.Context, projectID, deploymentName string, timeout time.Duration) error {
	return waitDeploymentInTime(ctx, projectID, deploymentName, "IDLE", timeout)
}

func WaitDeploymentRemoved(ctx context.Context, projectID, deploymentName string, timeout time.Duration) error {
	err := waitDeploymentInTime(ctx, projectID, deploymentName, "", timeout)
	if strings.Contains(err.Error(), "SERVERLESS_INSTANCE_NOT_FOUND") {
		return nil
	}
	return err
}

func waitDeploymentInTime(ctx context.Context, projectID, deploymentName, goal string, timeout time.Duration) error {
	client, err := NewAPIClient()
	if err != nil {
		return fmt.Errorf("failed to get API client: %w", err)
	}
	start := time.Now()
	for time.Since(start) < timeout {
		deployment, _, err :=
			client.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID, deploymentName).Execute()
		if err != nil {
			return fmt.Errorf("failed to check deployment %s: %w", deploymentName, err)
		}
		if *deployment.StateName == goal {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("%v timeout", timeout)
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Report(err error) {
	if err != nil {
		log.Print(err)
	}
}

func DefaultProviderName() string {
	return "AWS"
}

func DefaultRegion() string {
	return "US_EAST_2"
}

func RemoveProject(ctx context.Context, projectID string) {
	client, err := NewAPIClient()
	if err != nil {
		log.Printf("failed to get API client: %v", err)
		return
	}
	_, _, err = client.ProjectsApi.DeleteProject(ctx, projectID).Execute()
	if err != nil {
		log.Printf("failed to remove project %s: %v", projectID, err)
		return
	}
	log.Printf("Removed default project %s...", projectID)
}

func RemoveDeployment(ctx context.Context, projectID, deploymentName string) {
	client, err := NewAPIClient()
	if err != nil {
		log.Printf("failed to get API client: %v", err)
		return
	}
	_, _, err = client.ServerlessInstancesApi.DeleteServerlessInstance(ctx, projectID, deploymentName).Execute()
	if err != nil {
		log.Printf("failed to remove deployment %s: %v", deploymentName, err)
		return
	}
	log.Printf("Removed default deployment %s...", deploymentName)
}

func MustCreateAPIClient() *admin.APIClient {
	apiClient, err := NewAPIClient()
	if err != nil {
		panic(err)
	}
	return apiClient
}

func NewAPIClient() (*admin.APIClient, error) {
	return admin.NewClient(
		admin.UseBaseURL(Domain()),
		admin.UseDigestAuth(PublicAPIKey(), PrivateAPIKey()),
		admin.UseUserAgent(contractTestUserAgent()),
	)
}

func contractTestUserAgent() string {
	return fmt.Sprintf("%s/%s (%s;%s)", "MongoDBContractTestsAKO", version.Version, runtime.GOOS, runtime.GOARCH)
}

func OrgID() string {
	return mustGetEnv("MCLI_ORG_ID")
}

func Domain() string {
	return mustGetEnv("MCLI_OPS_MANAGER_URL")
}

func PublicAPIKey() string {
	return mustGetEnv("MCLI_PUBLIC_API_KEY")
}

func PrivateAPIKey() string {
	return mustGetEnv("MCLI_PRIVATE_API_KEY")
}

func mustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic("expected MCLI_ORG_ID was not set")
	}
	return value
}

func newRandomName(prefix string) string {
	randomSuffix := uuid.New().String()[0:6]
	return fmt.Sprintf("%s-%s", prefix, randomSuffix)
}
