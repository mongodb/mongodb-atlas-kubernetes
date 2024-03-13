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

const (
	ServerlessDeploymentTimeout = 5 * time.Minute
)

type Resources struct {
	Name           string `json:"name"`
	ProjectID      string `json:"projectId"`
	ServerlessName string `json:"deploymentName"`
}

func LoadResources(name string) (*Resources, error) {
	resources := Resources{Name: name}
	jsonData, err := os.ReadFile(resources.Filename())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &resources)
	if err == nil {
		if err := resources.check(); err != nil {
			return nil, err
		}
	}
	return &resources, nil
}

func (resources *Resources) String() string {
	return fmt.Sprintf("Name: %s\nProjectID: %s\nServerlessName: %s",
		resources.Name, resources.ProjectID, resources.ServerlessName)
}

func (resources *Resources) Filename() string {
	return fmt.Sprintf(".%s-resources.json", resources.Name)
}

func (resources *Resources) MustRecycle(ctx context.Context, wipe bool) {
	err := resources.Recycle(ctx, wipe)
	if err != nil {
		panic(err)
	}
}

func (resources *Resources) Recycle(ctx context.Context, wipe bool) error {
	if !wipe {
		log.Printf("Trying to stored %s test resource references for reuse...", resources.Name)
		err := resources.Store()
		if err != nil {
			return fmt.Errorf("failed to wipe store resources: %w", err)
		}
		log.Printf("Stored %s test resource references for reuse", resources.Name)
		return nil
	}
	log.Printf("Wiping %s test resources...", resources.Name)
	resources.wipe()
	if resources.ProjectID != "" && resources.ServerlessName != "" {
		err := removeServerless(ctx, resources.ProjectID, resources.ServerlessName)
		if err != nil {
			return err
		}
		err = waitServerlessRemoval(ctx, resources.ProjectID, resources.ServerlessName, ServerlessDeploymentTimeout)
		if err != nil {
			return fmt.Errorf("failed to get serverless deployment %s removed: %w", resources.ServerlessName, err)
		}
	}
	if resources.ProjectID != "" {
		err := removeProject(ctx, resources.ProjectID)
		if err != nil {
			return err
		}
	}
	log.Printf("Removed %s test resources", resources.Name)
	return nil
}

func (resources *Resources) Store() error {
	jsonData, err := json.Marshal(resources)
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return os.WriteFile(resources.Filename(), jsonData, 0600)
}

func (resources *Resources) wipe() error {
	return os.Remove(resources.Filename())
}

func (resources *Resources) check() error {
	ctx := context.Background()
	if resources.ProjectID != "" {
		if err := checkProject(ctx, resources.ProjectID); err != nil {
			return fmt.Errorf("failed to check project %s: %w", resources.ProjectID, err)
		}
	}
	if resources.ProjectID != "" && resources.ServerlessName != "" {
		if err := checkServerless(ctx, resources.ProjectID, resources.ServerlessName); err != nil {
			return fmt.Errorf("failed to check deployment %s: %w", resources.ServerlessName, err)
		}
	}
	return nil
}

func checkProject(ctx context.Context, projectID string) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = apiClient.ProjectsApi.GetProject(ctx, projectID).Execute()
	return err
}

func checkServerless(ctx context.Context, projectID, serverlessName string) error {
	apiClient, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = apiClient.ServerlessInstancesApi.GetServerlessInstance(ctx, projectID, serverlessName).Execute()
	return err
}

type OptResourceFunc func(ctx context.Context, resources *Resources) (*Resources, error)

func MustDeployTestResources(ctx context.Context, name string, wipe bool, project *admin.Group, optResourcesFn ...OptResourceFunc) *Resources {
	resources, err := DeployTestResources(ctx, name, wipe, project, optResourcesFn...)
	if err != nil {
		panic(err)
	}
	return resources
}

func DeployTestResources(ctx context.Context, name string, wipe bool, project *admin.Group, optResourcesFn ...OptResourceFunc) (*Resources, error) {
	existing, err := LoadResources(name)
	if err == nil {
		log.Printf("Reusing existing resources:\n%v", existing)
		return existing, nil
	}

	log.Printf("Cannot reuse resources: %v", err)
	resources := &Resources{Name: name}

	id, err := CreateProject(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to create test project for %s: %w", name, err)
	}
	resources.ProjectID = id

	for _, optResourceFn := range optResourcesFn {
		resources, err = optResourceFn(ctx, resources)
		if err != nil {
			return nil, fmt.Errorf("failed to create optional test resource for %s: %w", name, err)
		}
	}

	return resources, nil
}

func DefaultProject(prefix string) *admin.Group {
	return &admin.Group{
		Name:  newRandomName(prefix),
		OrgId: OrgID(),
	}
}

func DefaultServerless(prefix string) *admin.ServerlessInstanceDescriptionCreate {
	return &admin.ServerlessInstanceDescriptionCreate{
		Name: newRandomName(fmt.Sprintf("%s-serverless", prefix)),
		ProviderSettings: admin.ServerlessProviderSettings{
			BackingProviderName: DefaultProviderName(),
			RegionName:          DefaultRegion(),
		},
	}
}

func WithServerless(serverless *admin.ServerlessInstanceDescriptionCreate) OptResourceFunc {
	return func(ctx context.Context, resources *Resources) (*Resources, error) {
		deploymentName, err := CreateServerless(ctx, resources.ProjectID, serverless)
		if err != nil {
			return nil, fmt.Errorf("failed to create serverless deployment %s: %w", serverless.Name, err)
		}
		resources.ServerlessName = deploymentName
		if err := waitServerless(ctx, resources.ProjectID, resources.ServerlessName, "IDLE", ServerlessDeploymentTimeout); err != nil {
			return nil, fmt.Errorf("failed to get serverless deployment %s up and running: %w", serverless.Name, err)
		}
		return resources, nil
	}
}

func CreateProject(ctx context.Context, project *admin.Group) (string, error) {
	log.Printf("Creating project %s...", project.Name)
	apiClient, err := NewAPIClient()
	if err != nil {
		return "", err
	}
	newProject, _, err := apiClient.ProjectsApi.CreateProject(ctx, project).Execute()
	if err != nil {
		panic(err)
	}
	log.Printf("Created project %s ID=%s", newProject.Name, *newProject.Id)
	return *newProject.Id, nil
}

func CreateServerless(ctx context.Context, projectID string, serverless *admin.ServerlessInstanceDescriptionCreate) (string, error) {
	log.Printf("Creating serverless deployment %s...", serverless.Name)
	apiClient, err := NewAPIClient()
	if err != nil {
		return "", err
	}
	deployment, _, err :=
		apiClient.ServerlessInstancesApi.CreateServerlessInstance(ctx, projectID, serverless).Execute()
	if err != nil {
		return "", err
	}
	log.Printf("Created serverless deployment %s ID=%s", serverless.Name, *deployment.Id)
	return serverless.Name, nil
}

func waitServerlessRemoval(ctx context.Context, projectID, deploymentName string, timeout time.Duration) error {
	err := waitServerless(ctx, projectID, deploymentName, "", timeout)
	if strings.Contains(err.Error(), "SERVERLESS_INSTANCE_NOT_FOUND") {
		return nil
	}
	return err
}

func waitServerless(ctx context.Context, projectID, deploymentName, goal string, timeout time.Duration) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
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

func DefaultProviderName() string {
	return "AWS"
}

func DefaultRegion() string {
	return "US_EAST_2"
}

func removeProject(ctx context.Context, projectID string) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = client.ProjectsApi.DeleteProject(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove project %s: %w", projectID, err)
	}
	log.Printf("Removed default project %s...", projectID)
	return nil
}

func removeServerless(ctx context.Context, projectID, serverlessName string) error {
	client, err := NewAPIClient()
	if err != nil {
		return err
	}
	_, _, err = client.ServerlessInstancesApi.DeleteServerlessInstance(ctx, projectID, serverlessName).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove serverless deployment %s: %w", serverlessName, err)
	}
	log.Printf("Removed serverless deployment %s...", serverlessName)
	return nil
}

func NewAPIClient() (*admin.APIClient, error) {
	client, err := admin.NewClient(
		admin.UseBaseURL(Domain()),
		admin.UseDigestAuth(PublicAPIKey(), PrivateAPIKey()),
		admin.UseUserAgent(contractTestUserAgent()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get API client: %w", err)
	}
	return client, nil
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

func BoolEnv(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value == strings.ToLower("true")
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
