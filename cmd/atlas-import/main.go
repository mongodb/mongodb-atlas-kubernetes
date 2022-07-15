package main

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
)

func main() {
	log, _ = zap.NewDevelopment()
	zap.ReplaceGlobals(log)
	log.Debug("Beginning import procedure")
	exampleConfig := generateExampleConfig()
	err := runImports(exampleConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type atlasImportConfig struct {
	orgID            string
	publicKey        string
	privateKey       string
	atlasDomain      string
	importNamespace  string
	importAll        bool
	projectsToImport []importedProject
}

type importedProject struct {
	id string
	// importAll   bool
	// deployments []string
}

func generateExampleConfig() atlasImportConfig {
	// deploymentIDS := []string{"deploymentID1", "deploymentID2"}

	exampleConfig := atlasImportConfig{
		orgID:           "62a9dbe9fb598f6e67d540c5",
		publicKey:       "usrgnqxh",
		privateKey:      "SECRET",
		importNamespace: "test-namespace",
		atlasDomain:     "https://cloud-qa.mongodb.com/",
		importAll:       true,
		//projectsToImport: []importedProject{
		//	{
		//		id:          "projectID1",
		//		importAll:   false,
		//		deployments: deploymentIDS,
		//	},
		//	{
		//		id:          "projectID2",
		//		importAll:   false,
		//		deployments: deploymentIDS,
		//	},
		//	{
		//		id:          "projectID3",
		//		importAll:   true,
		//		deployments: nil,
		//	},
		// },
	}

	return exampleConfig
}

var backgroundCtx = context.Background()
var log *zap.Logger

func getAllProjects(atlasClient *mongodbatlas.Client) ([]*mongodbatlas.Project, error) {
	// Retrieve all projects associated to credentials
	// TODO for each paginated call, make sure to retrieve every items
	allProjects, _, err := atlasClient.Projects.GetAllProjects(backgroundCtx, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}
	projects := allProjects.Results
	return projects, err
}

func getListedProjects(atlasClient *mongodbatlas.Client, importConfig []importedProject) ([]*mongodbatlas.Project, error) {
	projects := make([]*mongodbatlas.Project, 0, len(importConfig))
	for _, importProject := range importConfig {
		atlasProject, _, err := atlasClient.Projects.GetOneProject(backgroundCtx, importProject.id)
		if err != nil {
			return nil, err
		}
		projects = append(projects, atlasProject)
	}
	return projects, nil
}

// setUpAtlasClient instantiate the client to interact with the Atlas API
// Credentials are provided in the import configuration
func setUpAtlasClient(config *atlasImportConfig) (*mongodbatlas.Client, error) {
	log.Debug("Creating AtlasClient")
	credentials := atlas.Connection{
		OrgID:      config.orgID,
		PublicKey:  config.publicKey,
		PrivateKey: config.privateKey,
	}
	atlasDomain := "https://cloud.mongodb.com/"
	if config.atlasDomain != "" {
		atlasDomain = config.atlasDomain
	}
	atlasClient, err := atlas.Client(atlasDomain, credentials, log.Sugar())

	return &atlasClient, err
}

// setUpKubernetesClient instantiate the client to interact with resources in the kubernetes cluster
// It also adds the operator CRDs to the scheme
// The kubernetes configuration can be retrieved the following ways (ordered by precedence) :
// * --kubeconfig flag pointing at a file
// * KUBECONFIG environment variable pointing at a file
// * In-cluster config if running in cluster
// * $HOME/.kube/config if exists.
func setUpKubernetesClient() (client.Client, error) {
	log.Debug("Creating kube client")

	// Add CRDs definitions to client scheme
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))

	// Instantiate the client to interact with k8s cluster
	kubeConfig, err := config.GetConfig()
	if err != nil {
		log.Error("Failed to retrieve kube config")
		return nil, err
	}
	kubeClient, err := client.New(kubeConfig, client.Options{
		Scheme: scheme,
	})

	if err != nil {
		log.Error("Failed to create kube client")
		return nil, err
	}

	return kubeClient, nil
}

func runImports(importConfig atlasImportConfig) error {
	// TODO check namespace exists

	atlasClient, err := setUpAtlasClient(&importConfig)
	if err != nil {
		return err
	}

	kubeClient, err := setUpKubernetesClient()
	if err != nil {
		return err
	}

	log.Debug("Importing projects")
	// Import all project if flag is set, otherwise import the ones specified by User
	var projects []*mongodbatlas.Project
	if importConfig.importAll {
		projects, err = getAllProjects(atlasClient)
	} else {
		projects, err = getListedProjects(atlasClient, importConfig.projectsToImport)
	}

	if err != nil {
		return err
	}

	log.Debug("Populating all imported projects")
	for _, atlasProject := range projects {
		// For each atlas project, retrieve associated information and convert to K8s kubernetesProject
		kubernetesProject, err := completeAndConvertProject(atlasProject, atlasClient, kubeClient, importConfig)
		if err != nil {
			return err
		}
		// Add resource to k8s cluster
		log.Debug(fmt.Sprintf("Instantiating project %s in Cluster", atlasProject.Name))
		if err := kubeClient.Create(backgroundCtx, kubernetesProject); err != nil {
			log.Error("Error when instantiating project in Cluster")
			return err
		}

		// Retrieve deployments, db users, backup schedules/policies associated
		kubernetesDatabaseUsers, err := getAndConvertDBUsers(atlasProject, atlasClient, importConfig)
		if err != nil {
			return err
		}

		// Add these resources to k8s cluster
		projectRef := common.ResourceRefNamespaced{
			Name:      kubernetesProject.ObjectMeta.Name,
			Namespace: importConfig.importNamespace,
		}
		log.Debug("Instantiating database users")
		for _, kubernetesDatabaseUser := range kubernetesDatabaseUsers {
			kubernetesDatabaseUser.Spec.Project = projectRef
			if err := kubeClient.Create(backgroundCtx, kubernetesDatabaseUser); err != nil {
				return err
			}
		}
	}

	return nil
}

func getAndConvertDBUsers(atlasProject *mongodbatlas.Project, atlasClient *mongodbatlas.Client, importConfig atlasImportConfig) ([]*mdbv1.AtlasDatabaseUser, error) {
	atlasDatabaseUsers, _, err := atlasClient.DatabaseUsers.List(backgroundCtx, atlasProject.ID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	kubernetesDatabaseUsers := make([]*mdbv1.AtlasDatabaseUser, 0, len(atlasDatabaseUsers))
	for i := range atlasDatabaseUsers {
		convertedUser, err := mdbv1.AtlasDatabaseUserFromAtlas(&atlasDatabaseUsers[i], nil)
		// Username should already be alphanumeric according to Atlas API
		convertedUser.ObjectMeta.Name = toLowercaseAlphaNumeric(convertedUser.Spec.Username)
		convertedUser.Namespace = importConfig.importNamespace
		if err != nil {
			return nil, err
		}
		kubernetesDatabaseUsers = append(kubernetesDatabaseUsers, convertedUser)
	}
	return kubernetesDatabaseUsers, nil
}

// Retrieve the MaintenanceWindow associated to the project ID and convert it to kubernetes type.
func getWindow(atlasClient *mongodbatlas.Client, projectID string) (*project.MaintenanceWindow, error) {
	atlasWindow, _, err := atlasClient.MaintenanceWindows.Get(backgroundCtx, projectID)
	if err != nil {
		return nil, err
	}
	kubernetesWindow := project.MaintenanceWindowFromAtlas(atlasWindow)
	return kubernetesWindow, nil
}

// TODO refactor 3 methods below with Generics

// Retrieve the IpAccessLists associated to the project ID and convert them to kubernetes type.
func getAccessLists(atlasClient *mongodbatlas.Client, projectID string) ([]project.IPAccessList, error) {
	atlasAccessLists, _, err := atlasClient.ProjectIPAccessList.List(backgroundCtx, projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	kubernetesAccessLists := make([]project.IPAccessList, 0, len(atlasAccessLists.Results))
	for i := range atlasAccessLists.Results {
		convertedList, err := project.IPAccessListFromAtlas(&atlasAccessLists.Results[i])
		if err != nil {
			return nil, err
		}
		kubernetesAccessLists = append(kubernetesAccessLists, *convertedList)
	}

	return kubernetesAccessLists, nil
}

// Retrieve the Integrations associated to the project ID and convert them to kubernetes type.
func getIntegrations(atlasClient *mongodbatlas.Client, projectID string, kubeClient client.Client, importConfig atlasImportConfig) ([]project.Integration, error) {
	atlasIntegrations, _, err := atlasClient.Integrations.List(backgroundCtx, projectID)
	if err != nil {
		return nil, err
	}

	kubernetesIntegrations := make([]project.Integration, 0, len(atlasIntegrations.Results))
	for _, atlasIntegration := range atlasIntegrations.Results {
		convertedIntegration, err := project.IntegrationFromAtlas(atlasIntegration, kubeClient, importConfig.importNamespace, projectID)
		if err != nil {
			return nil, err
		}
		kubernetesIntegrations = append(kubernetesIntegrations, *convertedIntegration)
	}

	return kubernetesIntegrations, nil
}

// Retrieve the Private Endpoints associated to the project ID and convert them to kubernetes type.
func getPrivateEndpoints(atlasClient *mongodbatlas.Client, projectID string) ([]mdbv1.PrivateEndpoint, error) {
	var kubernetesPrivateEndpoints []mdbv1.PrivateEndpoint
	// Retrieving endpoints for each cloud provider
	for _, cloudProvider := range []string{"AWS", "GCP", "AZURE"} {
		atlasProviderEndpoints, _, err := atlasClient.PrivateEndpoints.List(backgroundCtx, projectID, cloudProvider, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, err
		}

		for i := range atlasProviderEndpoints {
			kubernetesPrivateEndpoint, err := mdbv1.PrivateEndpointFromAtlas(&atlasProviderEndpoints[i])
			if err != nil {
				// The endpoint is either not ready or invalid, skip it and continue to convert the others
				continue
			}
			kubernetesPrivateEndpoints = append(kubernetesPrivateEndpoints, *kubernetesPrivateEndpoint)
		}
	}
	return kubernetesPrivateEndpoints, nil
}

func completeAndConvertProject(atlasProject *mongodbatlas.Project, atlasClient *mongodbatlas.Client, kubeClient client.Client, importConfig atlasImportConfig) (*mdbv1.AtlasProject, error) {
	kubernetesWindow, err := getWindow(atlasClient, atlasProject.ID)
	if err != nil {
		return nil, err
	}
	kubernetesAccessLists, err := getAccessLists(atlasClient, atlasProject.ID)
	if err != nil {
		return nil, err
	}
	kubernetesIntegrations, err := getIntegrations(atlasClient, atlasProject.ID, kubeClient, importConfig)
	if err != nil {
		return nil, err
	}
	kubernetesPrivateEndpoints, err := getPrivateEndpoints(atlasClient, atlasProject.ID)
	if err != nil {
		return nil, err
	}

	kubernetesProject := &mdbv1.AtlasProject{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasProject",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			// Concatenate the "sanitized" human-readable name with project ID to guarantee the name to be unique
			Name:      toLowercaseAlphaNumeric(atlasProject.Name) + "-" + atlasProject.ID,
			Namespace: importConfig.importNamespace,
		},
		Spec: mdbv1.AtlasProjectSpec{
			Name:             atlasProject.Name,
			ConnectionSecret: nil, // Create a secret containing the three connection fields from atlasImportConfig and link it, as for integrations
			// TODO maybe better to just not specify connection secret (the operator's default is used in that case)
			ProjectIPAccessList:       kubernetesAccessLists,
			MaintenanceWindow:         *kubernetesWindow,
			PrivateEndpoints:          kubernetesPrivateEndpoints,
			WithDefaultAlertsSettings: false,
			X509CertRef:               nil, // Double check with anton if it can be ignored or not
			Integrations:              kubernetesIntegrations,
		},
		Status: status.AtlasProjectStatus{},
	}

	if atlasProject.WithDefaultAlertsSettings != nil {
		kubernetesProject.Spec.WithDefaultAlertsSettings = *atlasProject.WithDefaultAlertsSettings
	}

	return kubernetesProject, nil
}

// toLowercaseAlphaNumeric only keeps characters a-z, A-Z and 0-9 from a string, and turns uppercase chars to lowercase
func toLowercaseAlphaNumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return strings.ToLower(result.String())
}
