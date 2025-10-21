// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
)

// GenerateServiceLayers generates service layers for all mappings
func GenerateServiceLayers(resourceName string, mappings []MappingWithConfig, translationOutDir string) error {
	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		if err := generateServiceLayerWithVersion(resourceName, versionSuffix, mapping.OpenAPIConfig.Package, translationOutDir); err != nil {
			return fmt.Errorf("failed to generate service layer for version %s: %w", versionSuffix, err)
		}
	}
	return nil
}

func generateServiceLayerWithVersion(resourceName, versionSuffix, sdkPackage, translationOutDir string) error {
	packageName := strings.ToLower(resourceName) + versionSuffix
	translationDir := filepath.Join(translationOutDir, packageName)

	if err := os.MkdirAll(translationDir, 0755); err != nil {
		return fmt.Errorf("failed to create translation directory: %w", err)
	}

	// Generate main translation file
	if err := generateTranslationFileWithVersion(translationDir, resourceName, packageName); err != nil {
		return fmt.Errorf("failed to generate translation file: %w", err)
	}

	// Generate service file
	if err := generateServiceFileWithVersion(translationDir, resourceName, packageName, sdkPackage); err != nil {
		return fmt.Errorf("failed to generate service file: %w", err)
	}

	return nil
}

func generateTranslationFileWithVersion(dir, resourceName, packageName string) error {
	f := jen.NewFile(packageName)
	AddLicenseHeader(f)

	// Atlas resource struct
	f.Type().Id("Atlas" + resourceName).Struct(
		jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName+"Spec"),
	)

	// ConvertFrom method
	f.Func().Params(jen.Id("a").Op("*").Id("Atlas"+resourceName)).Id("ConvertFrom").Params(
		jen.Id("k8s").Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName),
	).Error().Block(
		jen.Comment("TODO: Implement conversion from Kubernetes resource to Atlas resource"),
		jen.Return(jen.Nil()),
	)

	// Compare method
	f.Func().Params(jen.Id("a").Op("*").Id("Atlas"+resourceName)).Id("Compare").Params(
		jen.Id("other").Op("*").Id("Atlas"+resourceName),
	).Bool().Block(
		jen.Comment("TODO: Implement field-by-field comparison"),
		jen.Return(jen.True()),
	)

	fileName := filepath.Join(dir, strings.ToLower(resourceName)+".go")
	return f.Save(fileName)
}

func generateServiceFileWithVersion(dir, resourceName, packageName, sdkPackage string) error {
	atlasAPI, err := GetAtlasAPIForCRD(resourceName)
	if err != nil {
		return fmt.Errorf("failed to get Atlas API for CRD %s: %w", resourceName, err)
	}

	f := jen.NewFile(packageName)
	AddLicenseHeader(f)

	// Atlas SDK import for this version
	f.ImportAlias(sdkPackage, "admin")

	// Service interface
	f.Type().Id("Atlas"+resourceName+"Service").Interface(
		jen.Id("Get").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("orgID").String(),
			jen.Id("resourceID").String(),
		).Params(jen.Op("*").Id("Atlas"+resourceName), jen.Error()),
		jen.Id("List").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("orgID").String(),
		).Params(jen.Index().Op("*").Id("Atlas"+resourceName), jen.Error()),
		jen.Id("Update").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("orgID").String(),
			jen.Id("resourceID").String(),
			jen.Id("a"+strings.ToLower(resourceName)).Op("*").Id("Atlas"+resourceName),
		).Params(jen.Op("*").Id("Atlas"+resourceName), jen.Error()),
		jen.Id("Delete").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("orgID").String(),
			jen.Id("resourceID").String(),
		).Error(),
	)

	// Service implementation struct
	f.Type().Id("Atlas" + resourceName + "ServiceImpl").Struct(
		jen.Id(strings.ToLower(resourceName)+"API").Qual(sdkPackage, atlasAPI),
	)

	// Constructor
	f.Func().Id("NewAtlas" + resourceName + "Service").Params(
		jen.Id("api").Qual(sdkPackage, atlasAPI),
	).Id("Atlas" + resourceName + "Service").Block(
		jen.Return(jen.Op("&").Id("Atlas" + resourceName + "ServiceImpl").Values(jen.Dict{
			jen.Id(strings.ToLower(resourceName) + "API"): jen.Id("api"),
		})),
	)

	// Method implementations
	serviceVar := "s"
	methods := []struct {
		name   string
		params []jen.Code
		ret    []jen.Code
	}{
		{
			name: "Get",
			params: []jen.Code{
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("orgID").String(),
				jen.Id("resourceID").String(),
			},
			ret: []jen.Code{jen.Op("*").Id("Atlas" + resourceName), jen.Error()},
		},
		{
			name: "List",
			params: []jen.Code{
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("orgID").String(),
			},
			ret: []jen.Code{jen.Index().Op("*").Id("Atlas" + resourceName), jen.Error()},
		},
		{
			name: "Update",
			params: []jen.Code{
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("orgID").String(),
				jen.Id("resourceID").String(),
				jen.Id("a" + strings.ToLower(resourceName)).Op("*").Id("Atlas" + resourceName),
			},
			ret: []jen.Code{jen.Op("*").Id("Atlas" + resourceName), jen.Error()},
		},
		{
			name: "Delete",
			params: []jen.Code{
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("orgID").String(),
				jen.Id("resourceID").String(),
			},
			ret: []jen.Code{jen.Error()},
		},
	}

	for _, method := range methods {
		f.Func().Params(jen.Id(serviceVar).Op("*").Id("Atlas"+resourceName+"ServiceImpl")).Id(method.name).Params(method.params...).Params(method.ret...).Block(
			jen.Comment("TODO: Implement Atlas API call"),
			jen.Return(jen.Nil().Op(",").Qual("fmt", "Errorf").Call(jen.Lit("not implemented"))),
		)
	}

	fileName := filepath.Join(dir, "service.go")
	return f.Save(fileName)
}

// GetAtlasAPIForCRD maps CRD kinds to their corresponding Atlas API types
func GetAtlasAPIForCRD(crdKind string) (string, error) {
	apiMapping := map[string]string{
		"Project":                  "ProjectsApi",
		"Group":                    "ProjectsApi", // Groups are managed by ProjectsApi
		"Organization":             "OrganizationsApi",
		"DatabaseUser":             "DatabaseUsersApi",
		"Deployment":               "ClustersApi",
		"StreamInstance":           "StreamsApi",
		"PrivateEndpoint":          "PrivateEndpointServicesApi",
		"NetworkPeering":           "NetworkPeeringApi",
		"NetworkPeeringConnection": "NetworkPeeringApi",
		"NetworkContainer":         "NetworkPeeringApi",
		"IPAccessList":             "ProjectIPAccessListApi",
		"CustomRole":               "CustomDatabaseRolesApi",
		"BackupCompliancePolicy":   "CompliancePoliciesApi",
		"DataFederation":           "DataFederationApi",
		"ThirdPartyIntegrations":   "ThirdPartyIntegrationsApi",
		"FederatedAuth":            "FederatedAuthenticationApi",
		"SearchIndexConfig":        "AtlasSearchApi",
		"OrgSettings":              "OrganizationsApi",
		"StreamConnection":         "StreamsApi",
		"Team":                     "TeamsApi",
	}

	if api, exists := apiMapping[crdKind]; exists {
		return api, nil
	}
	return "", fmt.Errorf("no Atlas API mapping found for CRD kind '%s'", crdKind)
}

// AddLicenseHeader adds the standard license header to generated files
func AddLicenseHeader(f *jen.File) {
	f.HeaderComment("Copyright 2025 MongoDB Inc")
	f.HeaderComment("")
	f.HeaderComment("Licensed under the Apache License, Version 2.0 (the \"License\");")
	f.HeaderComment("you may not use this file except in compliance with the License.")
	f.HeaderComment("You may obtain a copy of the License at")
	f.HeaderComment("")
	f.HeaderComment("\thttp://www.apache.org/licenses/LICENSE-2.0")
	f.HeaderComment("")
	f.HeaderComment("Unless required by applicable law or agreed to in writing, software")
	f.HeaderComment("distributed under the License is distributed on an \"AS IS\" BASIS,")
	f.HeaderComment("WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.")
	f.HeaderComment("See the License for the specific language governing permissions and")
	f.HeaderComment("limitations under the License.")
}
