// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
