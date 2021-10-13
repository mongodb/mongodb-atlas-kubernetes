package config

const (
	ConfigAll = "../../deploy/" // Released generated files

	// ClusterSample      = "data/atlascluster_basic.yaml"
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
)
