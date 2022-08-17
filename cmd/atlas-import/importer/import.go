package importer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap/zapcore"

	v1 "k8s.io/api/core/v1"
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

// TODO improve credentials mgmt, see discussions on Spec document
// TODO import project by names instead of IDs
// TODO depending on the team discussion outcome about using Advanced API only for clusters, import advanced only
// TODO pre-requisite : having applied CRDs to the k8s cluster

// AtlasImportConfig contains the full import configuration
type AtlasImportConfig struct {
	OrgID            string            `yaml:"OrgID"`
	PublicKey        string            `yaml:"PublicKey"`
	PrivateKey       string            `yaml:"PrivateKey"`
	AtlasDomain      string            `yaml:"AtlasDomain"`
	ImportNamespace  string            `yaml:"ImportNamespace"`
	ImportAll        bool              `yaml:"ImportAll"`
	ImportedProjects []ImportedProject `yaml:"ImportedProjects"`
	Verbose          bool              `yaml:"Verbose"`
}

type ImportedProject struct {
	ID          string   `yaml:"Id"`
	ImportAll   bool     `yaml:"ImportAll"`
	Deployments []string `yaml:"Deployments"`
}

// Global variables
var backgroundCtx = context.Background()
var Log *zap.Logger

const kubeAPIVersion = "v1"
const maxItemPerPageAtlas = 500

// setUpAtlasClient instantiate the client to interact with the Atlas API
// Credentials are provided in the import configuration
func setUpAtlasClient(config *AtlasImportConfig) (*mongodbatlas.Client, error) {
	Log.Debug("Creating AtlasClient")
	credentials := atlas.Connection{
		OrgID:      config.OrgID,
		PublicKey:  config.PublicKey,
		PrivateKey: config.PrivateKey,
	}
	atlasDomain := "https://cloud.mongodb.com/"
	if config.AtlasDomain != "" {
		atlasDomain = config.AtlasDomain
	}
	atlasClient, err := atlas.Client(atlasDomain, credentials, Log.Sugar())
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
	Log.Debug("Creating kube client")

	// Add CRDs definitions to client scheme
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))

	// Instantiate the client to interact with k8s cluster
	kubeConfig, err := config.GetConfig()
	if err != nil {
		Log.Error("Failed to retrieve kube config")
		return nil, err
	}
	kubeClient, err := client.New(kubeConfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		Log.Error("Failed to create kube client")
		return nil, err
	}
	return kubeClient, nil
}

func ensureNamespaceExists(kubeClient client.Client, namespace string) error {
	// The namespace object should be part of global namespace
	objKey := client.ObjectKey{
		Namespace: "global",
		Name:      namespace,
	}
	nameSpaceObject := v1.Namespace{}
	if err := kubeClient.Get(backgroundCtx, objKey, &nameSpaceObject); err != nil {
		return err
	}
	if nameSpaceObject.Name != namespace {
		Log.Error("the specified import namespace couldn't be retrieved from kubernetes cluster")
		return errors.New("namespace " + namespace + " doesn't exist")
	}
	return nil
}

func initCustomZapLogger(level, encoding string) (*zap.Logger, error) {
	lv := zap.AtomicLevel{}
	err := lv.UnmarshalText([]byte(strings.ToLower(level)))
	if err != nil {
		return nil, err
	}

	enc := strings.ToLower(encoding)
	if enc != "json" && enc != "console" {
		return nil, errors.New("'encoding' parameter can only by either 'json' or 'console'")
	}

	cfg := zap.Config{
		Level:             lv,
		OutputPaths:       []string{"stdout"},
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          enc,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	return cfg.Build()
}

func fullSetUp(importConfig AtlasImportConfig) (*mongodbatlas.Client, client.Client, error) {
	var err error
	logLevel := "info"
	if importConfig.Verbose {
		logLevel = "debug"
	}
	Log, err = initCustomZapLogger(logLevel, "console")
	if err != nil {
		println("Impossible to initialize logger : " + err.Error())
		os.Exit(1)
	}
	zap.ReplaceGlobals(Log)

	Log.Info("Instantiating Atlas and Kubernetes Clients")

	atlasClient, err := setUpAtlasClient(&importConfig)
	if err != nil {
		return nil, nil, err
	}

	kubeClient, err := setUpKubernetesClient()
	if err != nil {
		return nil, nil, err
	}

	// Verifying that import namespace exists
	if err = ensureNamespaceExists(kubeClient, importConfig.ImportNamespace); err != nil {
		return nil, nil, err
	}

	return atlasClient, kubeClient, nil
}

func RunImports(importConfig AtlasImportConfig) error {
	atlasClient, kubeClient, err := fullSetUp(importConfig)
	if err != nil {
		return err
	}

	Log.Info("Listing projects to import")
	// Import all project if flag is set, otherwise import the ones specified by User
	var projects []*mongodbatlas.Project
	if importConfig.ImportAll {
		projects, err = getAllProjects(atlasClient)
	} else {
		projects, err = getListedProjects(atlasClient, importConfig.ImportedProjects)
	}
	if err != nil {
		return err
	}

	Log.Info("Populating projects")
	for _, atlasProject := range projects {
		// For each atlas project, retrieve associated information and convert to K8s kubernetesProject
		kubernetesProject, err := completeAndConvertProject(atlasProject, atlasClient, kubeClient, importConfig)
		if err != nil {
			return err
		}
		// Add resource to k8s cluster
		Log.Info(fmt.Sprintf("Instantiating Project : %s in Cluster", atlasProject.Name))
		instantiateKubernetesObject(kubeClient, kubernetesProject)

		projectRef := &common.ResourceRefNamespaced{
			Name:      kubernetesProject.ObjectMeta.Name,
			Namespace: importConfig.ImportNamespace,
		}

		// Retrieve and instantiate associated db users
		kubernetesDatabaseUsers, err := getAndConvertDBUsers(atlasProject, atlasClient, importConfig)
		if err != nil {
			return err
		}

		Log.Info("Instantiating project's Database Users")
		for _, kubernetesDatabaseUser := range kubernetesDatabaseUsers {
			kubernetesDatabaseUser.Spec.Project = *projectRef
			instantiateKubernetesObject(kubeClient, kubernetesDatabaseUser)
		}

		// Retrieve and instantiate associated Deployments
		kubernetesDeployments, err := getAndConvertDeployments(atlasProject, atlasClient, importConfig, projectRef)
		if err != nil {
			return err
		}

		Log.Info("Instantiating project's Deployments and their Backup Policies")
		for _, kubernetesDeployment := range kubernetesDeployments {
			deploymentName, err := getDeploymentName(kubernetesDeployment)
			if err != nil {
				return err
			}
			schedule, policy, err := retrieveBackupSchedule(atlasClient, atlasProject.ID, deploymentName, importConfig)
			// TODO check that the error is indeed a 404 (meaning there's no schedule) and can be ignored
			if err == nil {
				// Linking deployment to its schedule (policy is already linked to the schedule)
				kubernetesDeployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
					Name:      schedule.Name,
					Namespace: importConfig.ImportNamespace,
				}
				instantiateKubernetesObject(kubeClient, schedule)
				instantiateKubernetesObject(kubeClient, policy)
			}
			instantiateKubernetesObject(kubeClient, kubernetesDeployment)
		}
	}
	return nil
}

// ======================= ATLAS DEPLOYMENTS =======================

func getDeploymentName(deployment *mdbv1.AtlasDeployment) (string, error) {
	switch {
	case deployment.Spec.DeploymentSpec != nil:
		return deployment.Spec.DeploymentSpec.Name, nil
	case deployment.Spec.AdvancedDeploymentSpec != nil:
		return deployment.Spec.AdvancedDeploymentSpec.Name, nil
	case deployment.Spec.ServerlessSpec != nil:
		return deployment.Spec.ServerlessSpec.Name, nil
	default:
		return "", errors.New("deployment resource contains no specification")
	}
}

func getAllPaginatedResources[resource any](paginatedCall func(options *mongodbatlas.ListOptions) ([]resource, *mongodbatlas.Response, error)) []resource {
	listOptions := &mongodbatlas.ListOptions{
		PageNum:      0,
		ItemsPerPage: maxItemPerPageAtlas,
		IncludeCount: true,
	}
	shouldContinue := true
	resources := make([]resource, 0)
	for currPageNum := 1; shouldContinue; currPageNum++ {
		listOptions.PageNum = currPageNum
		newResources, res, err := paginatedCall(listOptions)
		if err != nil {
			Log.Fatal("Impossible to fetch resource : " + err.Error())
		}
		resources = append(resources, newResources...)
		if res.IsLastPage() {
			shouldContinue = false
		}
	}
	return resources
}

func getAllDeployments(atlasClient *mongodbatlas.Client, projectID string) ([]mongodbatlas.Cluster, []*mongodbatlas.AdvancedCluster, []*mongodbatlas.Cluster) {
	atlasDeployments := getAllPaginatedResources(
		func(options *mongodbatlas.ListOptions) ([]mongodbatlas.Cluster, *mongodbatlas.Response, error) {
			return atlasClient.Clusters.List(backgroundCtx, projectID, options)
		},
	)
	atlasAdvancedDeployments := getAllPaginatedResources(
		func(options *mongodbatlas.ListOptions) ([]*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
			rep, res, err := atlasClient.AdvancedClusters.List(backgroundCtx, projectID, options)
			return rep.Results, res, err
		},
	)
	atlasServerlessDeployments := getAllPaginatedResources(
		func(options *mongodbatlas.ListOptions) ([]*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
			rep, res, err := atlasClient.ServerlessInstances.List(backgroundCtx, projectID, options)
			return rep.Results, res, err
		},
	)
	return atlasDeployments, atlasAdvancedDeployments, atlasServerlessDeployments
}

func getListedDeployments(atlasClient *mongodbatlas.Client, projectID string, deploymentList []string) (
	[]mongodbatlas.Cluster, []*mongodbatlas.AdvancedCluster, []*mongodbatlas.Cluster, error) {
	var atlasDeployments []mongodbatlas.Cluster
	var atlasAdvancedDeployments []*mongodbatlas.AdvancedCluster
	var atlasServerlessDeployments []*mongodbatlas.Cluster
	for _, deploymentName := range deploymentList {
		// Here we need to find which type of deployment we need to import
		deployment, res, err := atlasClient.Clusters.Get(backgroundCtx, projectID, deploymentName)
		if err != nil && res.StatusCode != 400 && res.StatusCode != 404 {
			// Atlas returns a 400 code when trying to import a Multi-Cloud cluster through the normal API
			return nil, nil, nil, err
		} else if err == nil {
			atlasDeployments = append(atlasDeployments, *deployment)
			continue
		}
		advancedDeployment, res, err := atlasClient.AdvancedClusters.Get(backgroundCtx, projectID, deploymentName)
		if err != nil && res.StatusCode != 400 && res.StatusCode != 404 {
			return nil, nil, nil, err
		} else if err == nil {
			atlasAdvancedDeployments = append(atlasAdvancedDeployments, advancedDeployment)
			continue
		}
		serverlessDeployment, res, err := atlasClient.ServerlessInstances.Get(backgroundCtx, projectID, deploymentName)
		if err != nil && res.StatusCode != 404 {
			return nil, nil, nil, err
		} else if err == nil {
			atlasServerlessDeployments = append(atlasServerlessDeployments, serverlessDeployment)
		} else {
			Log.Fatal("Couldn't find any deployment with the name " + deploymentName + ". Error message : " + err.Error())
		}
	}
	return atlasDeployments, atlasAdvancedDeployments, atlasServerlessDeployments, nil
}

func getProjectConfig(importConfig AtlasImportConfig, projectID string) (bool, []string) {
	var deploymentList []string
	importAllDeployments := true
	if !importConfig.ImportAll {
		for _, configProject := range importConfig.ImportedProjects {
			if configProject.ID == projectID {
				importAllDeployments = configProject.ImportAll
				if !importAllDeployments {
					deploymentList = configProject.Deployments
				}
			}
		}
	}
	return importAllDeployments, deploymentList
}

func getAndConvertDeployments(atlasProject *mongodbatlas.Project, atlasClient *mongodbatlas.Client,
	importConfig AtlasImportConfig, projectRef *common.ResourceRefNamespaced) ([]*mdbv1.AtlasDeployment, error) {
	/*
		Atlas separates Deployments in 3 types : normal, advanced and serverless
		Normal and Serverless are returned as type "Cluster", Advanced is returned as "AdvancedCluster"
		But the API call for advanced clusters returns the normal ones as well
		Under the hood, they are the same thing in Atlas, normal is a legacy version which doesn't
		support multi-cloud Deployments
	*/

	// Normal and serverless are both of type cluster in Atlas API, but are returned by different API calls
	var atlasDeployments []mongodbatlas.Cluster
	var atlasAdvancedDeployments []*mongodbatlas.AdvancedCluster
	var atlasServerlessDeployments []*mongodbatlas.Cluster

	importAllDeployments, deploymentList := getProjectConfig(importConfig, atlasProject.ID)

	if importAllDeployments {
		atlasDeployments, atlasAdvancedDeployments, atlasServerlessDeployments = getAllDeployments(atlasClient, atlasProject.ID)
	} else {
		var err error
		atlasDeployments, atlasAdvancedDeployments, atlasServerlessDeployments, err = getListedDeployments(atlasClient, atlasProject.ID, deploymentList)
		if err != nil {
			return nil, err
		}
	}

	kubernetesDeployments := make([]*mdbv1.AtlasDeployment, 0, len(atlasAdvancedDeployments)+len(atlasServerlessDeployments))

	normalDeploymentSet := make(map[string]bool)

	for i := range atlasDeployments {
		kubernetesDeploymentSpec, err := mdbv1.DeploymentFromAtlas(&atlasDeployments[i])
		if err != nil {
			return nil, err
		}
		normalDeploymentSet[kubernetesDeploymentSpec.Name] = true
		kubernetesDeployment := generateKubeDeploymentCRDFromSpecs(kubernetesDeploymentSpec, nil,
			nil, importConfig, kubernetesDeploymentSpec.Name, projectRef)
		kubernetesDeployments = append(kubernetesDeployments, kubernetesDeployment)
	}

	for i := range atlasAdvancedDeployments {
		_, isNormal := normalDeploymentSet[atlasAdvancedDeployments[i].Name]
		if isNormal {
			Log.Debug("Skipping " + atlasAdvancedDeployments[i].Name + " because already imported as normal")
			// Already managed previously, skip
			continue
		}
		kubernetesAdvancedDeploymentSpec, err := mdbv1.AdvancedDeploymentFromAtlas(atlasAdvancedDeployments[i])
		if err != nil {
			return nil, err
		}
		kubernetesDeployment := generateKubeDeploymentCRDFromSpecs(nil, kubernetesAdvancedDeploymentSpec,
			nil, importConfig, kubernetesAdvancedDeploymentSpec.Name, projectRef)
		kubernetesDeployments = append(kubernetesDeployments, kubernetesDeployment)
	}

	for i := range atlasServerlessDeployments {
		kubernetesServerlessDeploymentSpec, err := mdbv1.ServerlessDeploymentFromAtlas(atlasServerlessDeployments[i])
		if err != nil {
			return nil, err
		}
		kubernetesDeployment := generateKubeDeploymentCRDFromSpecs(nil, nil,
			kubernetesServerlessDeploymentSpec, importConfig, kubernetesServerlessDeploymentSpec.Name, projectRef)
		kubernetesDeployments = append(kubernetesDeployments, kubernetesDeployment)
	}

	return kubernetesDeployments, nil
}

func generateKubeDeploymentCRDFromSpecs(normalSpec *mdbv1.DeploymentSpec,
	advancedSpec *mdbv1.AdvancedDeploymentSpec, serverlessSpec *mdbv1.ServerlessSpec, importConfig AtlasImportConfig,
	deploymentName string, projectRef *common.ResourceRefNamespaced) *mdbv1.AtlasDeployment {
	kubernetesDeployment := mdbv1.AtlasDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasDeployment",
			APIVersion: kubeAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			// Deployment names in Atlas are not case-sensitive, so we cannot have collisions even when simplifying the
			// name with toLowercaseAlphaNumeric
			Name:      stringNameConcatenation(projectRef.Name, toLowercaseAlphaNumeric(deploymentName)),
			Namespace: importConfig.ImportNamespace,
		},
		Spec: mdbv1.AtlasDeploymentSpec{
			Project:                *projectRef,
			DeploymentSpec:         normalSpec,
			AdvancedDeploymentSpec: advancedSpec,
			// Schedule ref will be updated before instantiation of the resources in the k8s cluster
			BackupScheduleRef: common.ResourceRefNamespaced{},
			ServerlessSpec:    serverlessSpec,
			ProcessArgs:       nil,
		},
		Status: status.AtlasDeploymentStatus{},
	}
	return &kubernetesDeployment
}

func retrieveBackupSchedule(atlasClient *mongodbatlas.Client, projectID string, deploymentName string,
	importConfig AtlasImportConfig) (*mdbv1.AtlasBackupSchedule, *mdbv1.AtlasBackupPolicy, error) {
	atlasBackupPolicy, _, err := atlasClient.CloudProviderSnapshotBackupPolicies.Get(backgroundCtx, projectID, deploymentName)
	if err != nil {
		return nil, nil, err
	}

	scheduleSpec, policySpec, err := mdbv1.BackupScheduleFromAtlas(atlasBackupPolicy)
	if err != nil {
		return nil, nil, err
	}

	prefix := stringNameConcatenation(toLowercaseAlphaNumeric(projectID), toLowercaseAlphaNumeric(deploymentName))
	backupScheduleName := stringNameConcatenation(prefix, "backup-schedule")
	backupPolicyName := stringNameConcatenation(prefix, "backup-policy")

	backupPolicy := &mdbv1.AtlasBackupPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasBackupPolicy",
			APIVersion: kubeAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupPolicyName,
			Namespace: importConfig.ImportNamespace,
		},
		Spec:   *policySpec,
		Status: mdbv1.AtlasBackupPolicyStatus{},
	}

	backupSchedule := &mdbv1.AtlasBackupSchedule{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasBackupSchedule",
			APIVersion: kubeAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupScheduleName,
			Namespace: importConfig.ImportNamespace,
		},
		Spec:   *scheduleSpec,
		Status: mdbv1.AtlasBackupScheduleStatus{},
	}

	backupSchedule.Spec.PolicyRef = common.ResourceRefNamespaced{
		Name:      backupPolicyName,
		Namespace: importConfig.ImportNamespace,
	}

	return backupSchedule, backupPolicy, nil
}

// ======================= ATLAS PROJECTS =======================

func getAllProjects(atlasClient *mongodbatlas.Client) ([]*mongodbatlas.Project, error) {
	// Retrieve all projects associated to credentials
	allProjects := getAllPaginatedResources(
		func(options *mongodbatlas.ListOptions) ([]*mongodbatlas.Project, *mongodbatlas.Response, error) {
			rep, res, err := atlasClient.Projects.GetAllProjects(backgroundCtx, options)
			if err != nil {
				return nil, res, err
			}
			return rep.Results, res, err
		},
	)
	projects := allProjects
	return projects, nil
}

func getListedProjects(atlasClient *mongodbatlas.Client, importConfig []ImportedProject) ([]*mongodbatlas.Project, error) {
	projects := make([]*mongodbatlas.Project, 0, len(importConfig))
	for _, importProject := range importConfig {
		atlasProject, _, err := atlasClient.Projects.GetOneProject(backgroundCtx, importProject.ID)
		if err != nil {
			return nil, err
		}
		projects = append(projects, atlasProject)
	}
	return projects, nil
}

func getAndConvertDBUsers(atlasProject *mongodbatlas.Project, atlasClient *mongodbatlas.Client,
	importConfig AtlasImportConfig) ([]*mdbv1.AtlasDatabaseUser, error) {
	atlasDatabaseUsers := getAllPaginatedResources(
		func(options *mongodbatlas.ListOptions) ([]mongodbatlas.DatabaseUser, *mongodbatlas.Response, error) {
			rep, res, err := atlasClient.DatabaseUsers.List(backgroundCtx, atlasProject.ID, options)
			return rep, res, err
		},
	)
	kubernetesDatabaseUsers := make([]*mdbv1.AtlasDatabaseUser, 0, len(atlasDatabaseUsers))
	for i := range atlasDatabaseUsers {
		convertedUser, err := mdbv1.AtlasDatabaseUserFromAtlas(&atlasDatabaseUsers[i], nil)
		// Username should already be alphanumeric according to Atlas API
		convertedUser.ObjectMeta.Name = toLowercaseAlphaNumeric(convertedUser.Spec.Username)
		convertedUser.Namespace = importConfig.ImportNamespace
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

func getAndConvertAssociatedResource[kubernetesResource any, atlasResource any](
	projectID string, conversionMethod func(*atlasResource) (*kubernetesResource, error),
	getMethod func(context.Context, string) ([]atlasResource, error)) (
	[]kubernetesResource, error) {
	atlasResourceList, err := getMethod(backgroundCtx, projectID)
	if err != nil {
		return nil, err
	}

	resourceList := make([]kubernetesResource, 0, len(atlasResourceList))
	for i := range atlasResourceList {
		convertedResource, err := conversionMethod(&atlasResourceList[i])
		if err != nil {
			return nil, err
		}
		resourceList = append(resourceList, *convertedResource)
	}

	return resourceList, nil
}

// Retrieve the IpAccessLists associated to the project ID and convert them to kubernetes type.
func getAccessLists(atlasClient *mongodbatlas.Client, projectID string) ([]project.IPAccessList, error) {
	getMethod := func(context.Context, string) ([]mongodbatlas.ProjectIPAccessList, error) {
		atlasAccessLists := getAllPaginatedResources(
			func(options *mongodbatlas.ListOptions) ([]mongodbatlas.ProjectIPAccessList, *mongodbatlas.Response, error) {
				atlasAccessLists, res, err := atlasClient.ProjectIPAccessList.List(backgroundCtx, projectID, options)
				return atlasAccessLists.Results, res, err
			})
		return atlasAccessLists, nil
	}
	ipAccessList, err := getAndConvertAssociatedResource[project.IPAccessList, mongodbatlas.ProjectIPAccessList](
		projectID, project.IPAccessListFromAtlas, getMethod)
	return ipAccessList, err
}

// Retrieve the Integrations associated to the project ID and convert them to kubernetes type.
func getIntegrations(atlasClient *mongodbatlas.Client, projectID string, kubeClient client.Client,
	importConfig AtlasImportConfig) ([]project.Integration, error) {
	atlasIntegrations, _, err := atlasClient.Integrations.List(backgroundCtx, projectID)
	if err != nil {
		return nil, err
	}

	kubernetesIntegrations := make([]project.Integration, 0, len(atlasIntegrations.Results))
	for _, atlasIntegration := range atlasIntegrations.Results {
		convertedIntegration, err := project.IntegrationFromAtlas(atlasIntegration, kubeClient, importConfig.ImportNamespace, projectID)
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
		atlasProviderEndpoints := getAllPaginatedResources(
			func(options *mongodbatlas.ListOptions) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
				return atlasClient.PrivateEndpoints.List(backgroundCtx, projectID, cloudProvider, options)
			})

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

func completeAndConvertProject(atlasProject *mongodbatlas.Project, atlasClient *mongodbatlas.Client,
	kubeClient client.Client, importConfig AtlasImportConfig) (*mdbv1.AtlasProject, error) {
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

	// Concatenate the "sanitized" human-readable name with project ID to guarantee the name to be unique
	projectName := stringNameConcatenation(toLowercaseAlphaNumeric(atlasProject.Name), atlasProject.ID)

	connectionSecret := storeAtlasSecret(projectName, importConfig, kubeClient)

	kubernetesProject := &mdbv1.AtlasProject{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasProject",
			APIVersion: kubeAPIVersion,
		},

		ObjectMeta: metav1.ObjectMeta{

			Name:      projectName,
			Namespace: importConfig.ImportNamespace,
		},
		Spec: mdbv1.AtlasProjectSpec{
			Name:                      atlasProject.Name,
			ConnectionSecret:          connectionSecret,
			ProjectIPAccessList:       kubernetesAccessLists,
			MaintenanceWindow:         *kubernetesWindow,
			PrivateEndpoints:          kubernetesPrivateEndpoints,
			WithDefaultAlertsSettings: false,
			X509CertRef:               nil, // TODO import certificate for Atlas
			Integrations:              kubernetesIntegrations,
		},
		Status: status.AtlasProjectStatus{},
	}

	if atlasProject.WithDefaultAlertsSettings != nil {
		kubernetesProject.Spec.WithDefaultAlertsSettings = *atlasProject.WithDefaultAlertsSettings
	}

	return kubernetesProject, nil
}

func storeAtlasSecret(projectName string, importConfig AtlasImportConfig, kubeClient client.Client) *common.ResourceRef {
	secretName := stringNameConcatenation(projectName, "secret")

	connectionSecretRef := common.ResourceRef{
		Name: secretName,
	}

	//TODO : import constants from pkg/api/controller/atlas/connection.go for field names ? Need to export them first
	data := map[string][]byte{
		"orgId":         []byte(importConfig.OrgID),
		"publicApiKey":  []byte(importConfig.PublicKey),
		"privateApiKey": []byte(importConfig.PrivateKey),
	}
	object := metav1.ObjectMeta{Name: secretName, Namespace: importConfig.ImportNamespace}
	secret := &v1.Secret{Data: data, ObjectMeta: object}

	instantiateKubernetesObject(kubeClient, secret)

	return &connectionSecretRef
}

// ======================= HELPER METHODS =======================

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

// Instantiates a kubernetes object in the cluster, using default ctx
// Terminates the program if an error occurs
func instantiateKubernetesObject(kubeClient client.Client, object client.Object) {
	if err := kubeClient.Create(backgroundCtx, object); err != nil {
		Log.Fatal("Failed to instantiate object " + object.GetName() + " in kube cluster, error is : " + err.Error())
	}
}

// Defines how we concatenate different fields to generate unique resource names
func stringNameConcatenation(str1 string, str2 string) string {
	return str1 + "-" + str2
}
