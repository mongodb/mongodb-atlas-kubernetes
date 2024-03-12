package search

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
)

const (
	TestTitle = "search"

	DeploymentTimeout = 6 * time.Minute

	// ReuseDeployed stores deployed resource references for reuse on next run
	ReuseDeployed = false
)

var (
	projectID      string
	deploymentName string
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	beforeAll(ctx)
	code := m.Run()
	afterAll(ctx)
	os.Exit(code)
}

func beforeAll(ctx context.Context) {
	if status := contract.LoadDeployedStatus(TestTitle); status.Error == nil {
		projectID = status.ProjectID
		deploymentName = status.DeploymentName
		log.Printf("Reusing project ID %s and deployment %s", projectID, deploymentName)
		return
	}
	projectID = contract.MustCreateDefaultProject(ctx, TestTitle)
	deploymentName = contract.MustCreateDefaultDeployment(ctx, projectID, TestTitle)
	contract.Must(contract.CreateDeploymentInTime(ctx, projectID, deploymentName, DeploymentTimeout))
	log.Printf("Project ID %s and deployment %s are ready", projectID, deploymentName)
}

func afterAll(ctx context.Context) {
	if ReuseDeployed {
		contract.StoreDeployedStatus(TestTitle, &contract.DeployedStatus{
			ProjectID:      projectID,
			DeploymentName: deploymentName,
		})
	} else {
		contract.RemoveDeployment(ctx, projectID, deploymentName)
		contract.Report(contract.WaitDeploymentRemoved(ctx, projectID, deploymentName, DeploymentTimeout))
		contract.RemoveProject(ctx, projectID)
		log.Printf("Project ID %s and deployment %s were removed", projectID, deploymentName)
	}
}

func TestCreateSearchIndex(t *testing.T) {
	fmt.Printf("yay!\n")
}
