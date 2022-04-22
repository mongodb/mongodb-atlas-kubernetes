package config

const (

	// Kubernetes configuration samples for users in deploy directory
	DefaultDeployConfig              = "../../deploy/" // Released generated files
	DefaultClusterWideCRDConfig      = "../../deploy/clusterwide/crds.yaml"
	DefaultClusterWideOperatorConfig = "../../deploy/clusterwide/clusterwide-config.yaml"
	DefaultNamespacedCRDConfig       = "../../deploy/namespaced/crds.yaml"
	DefaultNamespacedOperatorConfig  = "../../deploy/namespaced/namespaced-config.yaml"

	// Default names/path for tests coordinates
	DataGenFolder            = "data/gen" // for generated configs
	DefaultOperatorNS        = "mongodb-atlas-system"
	DefaultOperatorName      = "mongodb-atlas-operator"
	DefaultOperatorGlobalKey = "mongodb-atlas-operator-api-key"
	AtlasHost                = "https://cloud-qa.mongodb.com/"
	AtlasAPIURL              = AtlasHost + "api/atlas/v1.0/"
	TestAppLabelPrefix       = "app=test-app-"
	ActrcPath                = "../../.actrc"

	// HELM relative path
	TestAppHelmChartPath          = "../app/helm/"
	AtlasOperatorHelmChartPath    = "../../helm-charts/charts/atlas-operator"
	AtlasOperatorCRDHelmChartPath = "../../helm-charts/charts/atlas-operator-crds"
	AtlasClusterHelmChartPath     = "../../helm-charts/charts/atlas-cluster"
	HelmChartDirectory            = "../../helm-charts/charts"

	// AWS Tags for test
	TagName = "atlas-operator-test"
	TagBusy = "busy"
)
