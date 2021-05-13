package atlasinventory

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

// discoverInstances query atlas and return list of instances found
func discoverInstances(atlasClient *mongodbatlas.Client) ([]dbaasv1alpha1.Instance, workflow.Result) {
	// Try to find the service
	projects, response, err := atlasClient.Projects.GetAllProjects(context.Background(), &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, workflow.Terminate(getReasonFromResponse(response), err.Error())
	}
	instanceList := []dbaasv1alpha1.Instance{}
	for _, p := range projects.Results {
		clusters, response, err := atlasClient.Clusters.List(context.Background(), p.ID, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, workflow.Terminate(getReasonFromResponse(response), err.Error())
		}
		for _, cluster := range clusters {
			clusterSvc := dbaasv1alpha1.Instance{
				InstanceID: cluster.ID,
				Name:       cluster.Name,
				InstanceInfo: map[string]string{
					dbaas.InstanceSizeNameKey:             cluster.ProviderSettings.InstanceSizeName,
					dbaas.CloudProviderKey:                cluster.ProviderSettings.ProviderName,
					dbaas.CloudRegionKey:                  cluster.ProviderSettings.RegionName,
					dbaas.ProjectIDKey:                    p.ID,
					dbaas.ProjectNameKey:                  p.Name,
					dbaas.ConnectionStringsStandardSrvKey: cluster.ConnectionStrings.StandardSrv,
				},
			}
			instanceList = append(instanceList, clusterSvc)
		}
	}
	return instanceList, workflow.OK()
}

func getReasonFromResponse(response *mongodbatlas.Response) workflow.ConditionReason {
	reason := workflow.MongoDBAtlasInventoryBackendError
	if response == nil {
		return reason
	}
	if response.StatusCode == http.StatusUnauthorized {
		reason = workflow.MongoDBAtlasInventoryAuthenticationError
	} else if response.StatusCode == http.StatusBadGateway || response.StatusCode == http.StatusServiceUnavailable {
		reason = workflow.MongoDBAtlasInventoryEndpointUnreachable
	}
	return reason
}
