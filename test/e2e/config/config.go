package config

const (

	// Kubernetes configuration samples for users in deploy directory
	DefaultDeployConfig              = "../../deploy/" // Released generated files
	DefaultClusterWideCRDConfig      = "../../deploy/clusterwide/crds.yaml"
	DefaultClusterWideOperatorConfig = "../../deploy/clusterwide/clusterwide-config.yaml"
	DefaultNamespacedCRDConfig       = "../../deploy/namespaced/crds.yaml"
	DefaultNamespacedOperatorConfig  = "../../deploy/namespaced/namespaced-config.yaml"

	// Default names/path for tests coordinates
	DataFolder               = "data"
	DefaultOperatorNS        = "mongodb-atlas-system"
	DefaultOperatorName      = "mongodb-atlas-operator"
	DefaultOperatorGlobalKey = "mongodb-atlas-operator-api-key"
	AtlasHost                = "https://cloud-qa.mongodb.com"
	AtlasAPIURL              = AtlasHost + "/api/atlas/v1.0/"
	TestAppLabelPrefix       = "app=test-app-"
	ActrcPath                = "../../.actrc"

	// HELM related path
	HelmTestAppPath             = "../app/helm/"
	HelmOperatorChartPath       = "../../helm-charts/charts/atlas-operator"
	HelmCRDChartPath            = "../../helm-charts/charts/atlas-operator-crds"
	HelmAtlasResourcesChartPath = "../../helm-charts/charts/atlas-cluster"

	// AWS Tags for test
	TagName = "atlas-operator-test"
	TagBusy = "busy"
)
