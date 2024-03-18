package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

const (
	ProtocolPrefix = "mongodb+srv://"
)

// TestResources keeps track of all resourced for a given named test to allow
// reuse or consistent cleanup
type TestResources struct {
	Name string `json:"name"`

	WipeResources bool `json:""`

	ProjectID      string `json:"projectId"`
	ServerlessName string `json:"deploymentName"`
	ClusterURL     string `json:"clusterURL"`
	UserDB         string `json:"userDB"`
	Username       string `json:"username"`
	Password       string `json:"password"`

	DatabaseName   string `json:"databaseName"`
	CollectionName string `json:"collectionName"`
}

// OptResourceFunc allows to add and setup optional resources
type OptResourceFunc func(ctx context.Context, resources *TestResources) (*TestResources, error)

func DeployTestResources(ctx context.Context, name string, wipe bool, project *admin.Group, optResourcesFn ...OptResourceFunc) (*TestResources, error) {
	log.Printf("Wipe resources set to %v", wipe)
	existing, err := LoadResources(name)
	if err == nil {
		log.Printf("Reusing existing resources\n")
		return existing, nil
	}

	log.Printf("Cannot reuse resources: %v", err)
	resources := &TestResources{Name: name, WipeResources: wipe}

	id, err := createProject(ctx, project)
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
	log.Printf("Created new resources\n")
	log.Printf("Trying to stored %s test resource references for reuse...", resources.Name)
	if err := resources.store(); err != nil {
		return nil, fmt.Errorf("failed to wipe store resources: %w", err)
	}
	log.Printf("Stored %s test resource references for reuse", resources.Name)
	return resources, nil
}

func LoadResources(name string) (*TestResources, error) {
	resources := TestResources{Name: name}
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

func (resources *TestResources) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "Name: %s\n", resources.Name)
	fmt.Fprintf(&buf, "  ProjectID       %s\n", resources.ProjectID)
	fmt.Fprintf(&buf, "  ServerlessName  %s\n", resources.ServerlessName)
	fmt.Fprintf(&buf, "  ClusterURL      %s\n", resources.ClusterURL)
	fmt.Fprintf(&buf, "  UserDB          %s\n", resources.UserDB)
	fmt.Fprintf(&buf, "  Username        %s\n", resources.Username)
	fmt.Fprintf(&buf, "  Password        *******\n")
	fmt.Fprintf(&buf, "  DatabaseName    %s\n", resources.DatabaseName)
	fmt.Fprintf(&buf, "  CollectionName  %s\n", resources.CollectionName)
	return buf.String()
}

func (resources *TestResources) Filename() string {
	return fmt.Sprintf(".%s-resources.json", resources.Name)
}

func (resources *TestResources) Recycle(ctx context.Context) error {
	if !resources.WipeResources {
		log.Printf("Trying to stored %s test resource references for reuse...", resources.Name)
		err := resources.store()
		if err != nil {
			return fmt.Errorf("failed to wipe store resources: %w", err)
		}
		log.Printf("Stored %s test resource references for reuse", resources.Name)
		return nil
	}
	log.Printf("Wiping %s test resources...", resources.Name)
	resources.wipe()
	if resources.ProjectID != "" && resources.Username != "" {
		err := removeUser(ctx, resources.ProjectID, resources.UserDB, resources.Username)
		if err != nil {
			return err
		}
	}
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
		err := clearIPAccessList(ctx, resources.ProjectID)
		if err != nil {
			return err
		}
		err = removeProject(ctx, resources.ProjectID)
		if err != nil {
			return err
		}
	}
	log.Printf("Removed %s test resources", resources.Name)
	return nil
}

func (resources *TestResources) store() error {
	jsonData, err := json.Marshal(resources)
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	return os.WriteFile(resources.Filename(), jsonData, 0600)
}

func (resources *TestResources) wipe() error {
	return os.Remove(resources.Filename())
}

func (resources *TestResources) check() error {
	ctx := context.Background()
	if resources.ProjectID != "" {
		if err := checkProject(ctx, resources.ProjectID); err != nil {
			return fmt.Errorf("failed to check project %s: %w", resources.ProjectID, err)
		}
	}
	if resources.ProjectID != "" && resources.ServerlessName != "" {
		if err := checkServerless(ctx, resources.ProjectID, resources.ServerlessName); err != nil {
			return fmt.Errorf("failed to check serverless deployment %s: %w", resources.ServerlessName, err)
		}
	}
	if resources.ProjectID != "" && resources.Username != "" {
		if err := checkUser(ctx, resources.ProjectID, resources.UserDB, resources.Username); err != nil {
			return fmt.Errorf("failed to check user %s: %w", resources.ServerlessName, err)
		}
	}
	return nil
}
