package config

const (

	// Kubernetes configuration samples for users in deploy directory
	DefaultDeployConfig              = "../../deploy/" // Released generated files
	DefaultClusterWideCRDConfig      = "../../deploy/clusterwide/crds.yaml"
	DefaultClusterWideOperatorConfig = "../../deploy/clusterwide/clusterwide-config.yaml"
	DefaultNamespacedCRDConfig       = "../../deploy/namespaced/crds.yaml"
	DefaultNamespacedOperatorConfig  = "../../deploy/namespaced/namespaced-config.yaml"

	// Default names/path for tests coordinates
	DataGenFolder            = "output" // for generated configs
	DefaultOperatorNS        = "mongodb-atlas-system"
	DefaultOperatorName      = "mongodb-atlas-operator"
	DefaultOperatorGlobalKey = "mongodb-atlas-operator-api-key"
	AtlasHost                = "https://cloud-qa.mongodb.com/"
	AtlasAPIURL              = AtlasHost + "api/atlas/v1.0/"
	TestAppLabelPrefix       = "app=test-app-"
	ActrcPath                = "../../.actrc"

	// HELM relative path
	TestAppHelmChartPath          = "../app/helm/"
	AtlasOperatorHelmChartPath    = "../../helm-charts/atlas-operator"
	AtlasOperatorCRDHelmChartPath = "../../helm-charts/atlas-operator-crds"
	AtlasDeploymentHelmChartPath  = "../../helm-charts/atlas-deployment"
	HelmChartDirectory            = "../../helm-charts"
	MajorVersionFile              = "../../major-version"

	// Tags for test
	TagName         = "atlas-operator-test"
	TagForTestValue = "atlas-operator-e2e-test"
	TagForTestKey   = "atlas-operator-e2e-key"
	TagBusy         = "busy"

	// Regions for tests

	GCPRegion     = "europe-west1"
	AWSRegionUS   = "us-east-1"
	AWSRegionEU   = "eu-west-2"
	AzureRegionEU = "northeurope"

	// GCP
	FileNameSAGCP = "output/gcp_service_account.json"

	// X509 auth test PEM key
	PEMCertFileName = "output/x509cert.pem"
)
