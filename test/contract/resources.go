package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

// TestResources keeps track of all resourced for a given named test to allow
// reuse or consistent cleanup
type TestResources struct {
	Name           string `json:"name"`
	ProjectID      string `json:"projectId"`
	ServerlessName string `json:"deploymentName"`
}

// OptResourceFunc allows to add and setup optional resources
type OptResourceFunc func(ctx context.Context, resources *TestResources) (*TestResources, error)

func MustDeployTestResources(ctx context.Context, name string, wipe bool, project *admin.Group, optResourcesFn ...OptResourceFunc) *TestResources {
	resources, err := DeployTestResources(ctx, name, wipe, project, optResourcesFn...)
	if err != nil {
		panic(err)
	}
	return resources
}

func DeployTestResources(ctx context.Context, name string, wipe bool, project *admin.Group, optResourcesFn ...OptResourceFunc) (*TestResources, error) {
	existing, err := LoadResources(name)
	if err == nil {
		log.Printf("Reusing existing resources:\n%v", existing)
		return existing, nil
	}

	log.Printf("Cannot reuse resources: %v", err)
	resources := &TestResources{Name: name}

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
	return fmt.Sprintf("Name: %s\nProjectID: %s\nServerlessName: %s",
		resources.Name, resources.ProjectID, resources.ServerlessName)
}

func (resources *TestResources) Filename() string {
	return fmt.Sprintf(".%s-resources.json", resources.Name)
}

func (resources *TestResources) MustRecycle(ctx context.Context, wipe bool) {
	err := resources.Recycle(ctx, wipe)
	if err != nil {
		panic(err)
	}
}

func (resources *TestResources) Recycle(ctx context.Context, wipe bool) error {
	if !wipe {
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
			return fmt.Errorf("failed to check deployment %s: %w", resources.ServerlessName, err)
		}
	}
	return nil
}
