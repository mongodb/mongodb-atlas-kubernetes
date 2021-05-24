package config

const (
	ConfigAll          = "../../deploy/" // Released generated files
	ClusterSample      = "data/atlascluster_basic.yaml"
	DataFolder         = "data"
	DefaultOperatorNS  = "mongodb-atlas-system"
	AtlasURL           = "https://cloud-qa.mongodb.com"
	TestAppLabelPrefix = "app=test-app-"
	// HELM related path
	HelmTestAppPath             = "../app/helm/"
	HelmOperatorChartPath       = "../../my-charts/charts/atlas-operator"
	HelmCRDChartPath            = "../../my-charts/charts/atlas-operator-crds"
	HelmAtlasResourcesChartPath = "../../my-charts/charts/atlas-cluster"
)
