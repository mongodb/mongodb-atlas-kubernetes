package deployment_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/contract"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	depedenciesTimeout = 5 * time.Minute
)

func TestListDeployments(t *testing.T) {
	ctx := context.Background()
	contract.RunGoContractTest(ctx, t, "test list deployments", func(ch contract.ContractHelper) {
		projectName := utils.RandomName("deployments-list-test-project")
		clusterName := "cluster0"
		serverlessName := "serverless-name"
		log.Printf("Creating project with a cluster and serveless deployment...")
		require.NoError(t, ch.AddResources(
			ctx, depedenciesTimeout,
			contract.DefaultAtlasProject(projectName),
			akov2.DefaultAWSDeployment("", projectName).Lightweight().
				WithInstanceSize("M0").WithName(clusterName).WithAtlasName(clusterName),
			akov2.NewDefaultAWSServerlessInstance("", projectName).
				WithName(serverlessName).WithAtlasName(serverlessName)))
		testProjectID, err := ch.ProjectID(ctx, projectName)
		require.NoError(t, err)
		ds := deployment.NewAtlasDeployments(ch.AtlasClient().ClustersApi, ch.AtlasClient().ServerlessInstancesApi, nil, ch.AtlasClientSet().SdkClient20241113001.FlexClustersApi, false)

		names, err := ds.ListDeploymentNames(ctx, testProjectID)
		require.NoError(t, err)
		expected := []string{clusterName, serverlessName}
		assert.Equal(t, expected, names)
	})
}
