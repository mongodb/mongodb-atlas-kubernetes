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

const (
	pkgCtrlState = "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

// FromConfig generates controllers and handlers based on the parsed CRD result file
func FromConfig(resultPath, crdKind, controllerOutDir, translationOutDir string) error {
	parsedConfig, err := ParseCRDConfig(resultPath, crdKind)
	if err != nil {
		return err
	}

	resourceName := parsedConfig.ResourceName
	
	// Set default directories if not provided
	if controllerOutDir == "" {
		controllerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "controller")
	}
	if translationOutDir == "" {
		translationOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "translation")
	}
	
	baseControllerDir := filepath.Join(controllerOutDir, strings.ToLower(resourceName))

	// Generate translation layers for all mappings
	if err := GenerateTranslationLayers(resourceName, parsedConfig.Mappings, translationOutDir); err != nil {
		return fmt.Errorf("failed to generate translation layers: %w", err)
	}

	controllerName := resourceName
	controllerDir := baseControllerDir

	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("failed to create controller directory: %w", err)
	}

	if err := generateControllerFileWithMultipleVersions(controllerDir, controllerName, resourceName, parsedConfig.Mappings); err != nil {
		return fmt.Errorf("failed to generate controller file: %w", err)
	}

	if err := generateHandlerFileWithMultipleVersions(controllerDir, controllerName, resourceName, parsedConfig.Mappings); err != nil {
		return fmt.Errorf("failed to generate handler file: %w", err)
	}

	fmt.Printf("Successfully generated controller %s for resource %s with %d SDK versions at %s\n",
		controllerName, resourceName, len(parsedConfig.Mappings), controllerDir)

	return nil
}

// PrintCRDs displays available CRDs from the result file
func PrintCRDs(resultPath string) error {
	crdInfos, err := ListCRDs(resultPath)
	if err != nil {
		return err
	}

	fmt.Printf("Available CRDs in %s:\n\n", resultPath)
	for _, crd := range crdInfos {
		fmt.Printf("Kind: %s\n", crd.Kind)
		fmt.Printf("  Group: %s\n", crd.Group)
		fmt.Printf("  Version: %s\n", crd.Version)
		if len(crd.ShortNames) > 0 {
			fmt.Printf("  Short Names: %s\n", strings.Join(crd.ShortNames, ", "))
		}
		if len(crd.Categories) > 0 {
			fmt.Printf("  Categories: %s\n", strings.Join(crd.Categories, ", "))
		}
		if len(crd.Versions) > 0 {
			fmt.Printf("  SDK Versions:\n")
			for _, version := range crd.Versions {
				fmt.Printf("    %s: %s\n", version.Version, version.AtlasSDKVersion)
			}
		}
		fmt.Println()
	}
	return nil
}

func generateControllerFileWithMultipleVersions(dir, controllerName, resourceName string, mappings []MappingWithConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	atlasAPI, err := GetAtlasAPIForCRD(resourceName)
	if err != nil {
		return fmt.Errorf("failed to get Atlas API for CRD %s: %w", resourceName, err)
	}

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")

	// RBAC
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,resources=%s,verbs=get;list;watch;create;update;patch;delete", strings.ToLower("atlas"+resourceName+"s")))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,resources=%s/status,verbs=get;update;patch", strings.ToLower("atlas"+resourceName+"s")))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,resources=%s/finalizers,verbs=update", strings.ToLower("atlas"+resourceName+"s")))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=%s,verbs=get;list;watch;create;update;patch;delete", strings.ToLower("atlas"+resourceName+"s")))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=%s/status,verbs=get;update;patch", strings.ToLower("atlas"+resourceName+"s")))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=%s/finalizers,verbs=update", strings.ToLower("atlas"+resourceName+"s")))

	// Service builder types for each version
	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		f.Type().Id("serviceBuilderFunc"+versionSuffix).Func().Params(
			jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "ClientSet"),
		).Qual(fmt.Sprintf("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/%s%s", atlasResourceName, versionSuffix), "Atlas"+resourceName+"Service")
	}

	// Handler struct for ALL CRD versions
	handlerFields := []jen.Code{
		jen.Qual(pkgCtrlState, "StateHandler").Types(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName)),
		jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler"),
	}

	// Service builder for each version
	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		handlerFields = append(handlerFields, jen.Id("serviceBuilder"+versionSuffix).Id("serviceBuilderFunc"+versionSuffix))
	}

	f.Type().Id(controllerName + "Handler").Struct(handlerFields...)

	// NewReconciler method with all service builders
	reconcilerParams := []jen.Code{
		jen.Id("c").Qual("sigs.k8s.io/controller-runtime/pkg/cluster", "Cluster"),
		jen.Id("atlasProvider").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "Provider"),
		jen.Id("logger").Op("*").Qual("go.uber.org/zap", "Logger"),
		jen.Id("globalSecretRef").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey"),
		jen.Id("reapplySupport").Bool(),
	}

	serviceBuilderValues := jen.Dict{
		jen.Id("AtlasReconciler"): jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler").Values(jen.Dict{
			jen.Id("Client"):          jen.Id("c").Dot("GetClient").Call(),
			jen.Id("AtlasProvider"):   jen.Id("atlasProvider"),
			jen.Id("Log"):             jen.Id("logger").Dot("Named").Call(jen.Lit("controllers")).Dot("Named").Call(jen.Lit("Atlas" + resourceName)).Dot("Sugar").Call(),
			jen.Id("GlobalSecretRef"): jen.Id("globalSecretRef"),
		}),
	}

	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		serviceBuilderValues[jen.Id("serviceBuilder"+versionSuffix)] = jen.Func().Params(jen.Id("clientSet").Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "ClientSet")).Qual(fmt.Sprintf("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/%s%s", atlasResourceName, versionSuffix), "Atlas"+resourceName+"Service").Block(
			jen.Return(jen.Qual(fmt.Sprintf("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/%s%s", atlasResourceName, versionSuffix), "NewAtlas"+resourceName+"Service").Call(jen.Id("clientSet").Dot("SdkClient" + versionSuffix + "006").Dot(atlasAPI))),
		)
	}

	f.Func().Id("New"+controllerName+"Reconciler").Params(reconcilerParams...).Op("*").Qual(pkgCtrlState, "Reconciler").Types(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName)).Block(
		jen.Id(strings.ToLower(controllerName)+"Handler").Op(":=").Op("&").Id(controllerName+"Handler").Values(serviceBuilderValues),
		jen.Return(jen.Qual(pkgCtrlState, "NewStateReconciler").Call(
			jen.Id(strings.ToLower(controllerName)+"Handler"),
			jen.Qual(pkgCtrlState, "WithCluster").Types(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName)).Call(jen.Id("c")),
			jen.Qual(pkgCtrlState, "WithReapplySupport").Types(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName)).Call(jen.Id("reapplySupport")),
		)),
	)

	fileName := filepath.Join(dir, atlasResourceName+"_controller.go")
	return f.Save(fileName)
}

func generateHandlerFileWithMultipleVersions(dir, controllerName, resourceName string, mappings []MappingWithConfig) error {
	atlasResourceName := strings.ToLower(resourceName)

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")

	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		translationPkg := fmt.Sprintf("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/%s%s", atlasResourceName, versionSuffix)
		f.ImportAlias(translationPkg, atlasResourceName+versionSuffix)
	}

	f.Type().Id("reconcileRequest").Struct(
		jen.Id("version").String(),
		jen.Id(strings.ToLower("a"+resourceName)).Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName),
	)

	// Helper method to get service for resource (merged getSDKVersion + getServiceForVersion)
	f.Comment("getServiceForResource determines the SDK version from the resource spec and returns the appropriate service")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("getServiceForResource").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(strings.ToLower("a"+resourceName)).Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName),
	).Params(jen.Interface(), jen.Error()).Block(
		jen.Comment("Determine which SDK version to use from resource spec"),
		jen.Var().Id("version").String(),
		jen.BlockFunc(func(g *jen.Group) {
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				// Capitalize first letter of version (e.g., v20250312 -> V20250312)
				capitalizedVersion := strings.ToUpper(string(versionSuffix[0])) + versionSuffix[1:]
				g.If(jen.Id(strings.ToLower("a" + resourceName)).Dot("Spec").Dot(capitalizedVersion).Op("!=").Nil()).Block(
					jen.Id("version").Op("=").Lit(versionSuffix),
				)
			}
		}),
		jen.Comment("Ensure a version was specified"),
		jen.If(jen.Id("version").Op("==").Lit("")).Block(
			jen.Return(jen.Nil().Op(",").Qual("fmt", "Errorf").Call(jen.Lit("no SDK version specified in resource spec - please specify one of the supported versions"))),
		),
		jen.Comment("Get client set"),
		jen.Id("clientSet").Op(",").Id("err").Op(":=").Id("h").Dot("AtlasProvider").Dot("SdkClient").Call(jen.Id("ctx").Op(",").Id("h").Dot("GlobalSecretRef").Op(",").Id("h").Dot("Log")),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil().Op(",").Id("err")),
		),
		jen.Comment("Return appropriate service for version"),
		jen.Switch(jen.Id("version")).BlockFunc(func(g *jen.Group) {
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				g.Case(jen.Lit(versionSuffix)).Block(
					jen.Return(jen.Id("h").Dot("serviceBuilder" + versionSuffix).Call(jen.Id("clientSet")).Op(",").Nil()),
				)
			}
			g.Default().Block(
				jen.Return(jen.Nil().Op(",").Qual("fmt", "Errorf").Call(jen.Lit("unsupported SDK version: %s"), jen.Id("version"))),
			)
		}),
	)

	generateVersionAwareStateHandlers(f, controllerName, resourceName, mappings)

	fileName := filepath.Join(dir, "handler.go")
	return f.Save(fileName)
}

func generateVersionAwareStateHandlers(f *jen.File, controllerName, resourceName string, mappings []MappingWithConfig) {
	// HandleInitial method
	f.Comment("HandleInitial handles the initial state")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("HandleInitial").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(strings.ToLower("a"+resourceName)).Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName),
	).Params(
		jen.Qual(pkgCtrlState, "Result"),
		jen.Error(),
	).Block(
		jen.Comment("TODO: Implement initial state logic"),
		jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "NextState").Call(
			jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", "StateUpdated"),
			jen.Lit("Updated Atlas"+resourceName+"."),
		)),
	)

	// HandleUpdated method
	f.Comment("HandleUpdated handles the updated state")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("HandleUpdated").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(strings.ToLower("a"+resourceName)).Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName),
	).Params(
		jen.Qual(pkgCtrlState, "Result"),
		jen.Error(),
	).Block(
		jen.Comment("Get the appropriate service for this resource"),
		jen.List(jen.Id("svc"), jen.Id("err")).Op(":=").Id("h").Dot("getServiceForResource").Call(jen.Id("ctx"), jen.Id(strings.ToLower("a"+resourceName))),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "Error").Call(
				jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", "StateUpdated"),
				jen.Id("err"),
			)),
		),
		jen.Comment("TODO: Use the service to implement updated state logic with Atlas API calls"),
		jen.Id("_").Op("=").Id("svc").Comment("Use service in implementation"),
		jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "NextState").Call(
			jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", "StateUpdated"),
			jen.Lit("Ready"),
		)),
	)

	// Generate all other state handler methods
	handlers := []struct {
		name      string
		nextState string
		message   string
	}{
		{"HandleImportRequested", "StateImported", "Import completed"},
		{"HandleImported", "StateUpdated", "Ready"},
		{"HandleCreating", "StateCreated", "Resource created"},
		{"HandleCreated", "StateUpdated", "Ready"},
		{"HandleUpdating", "StateUpdated", "Update completed"},
		{"HandleDeletionRequested", "StateDeleting", "Deletion started"},
		{"HandleDeleting", "StateDeleted", "Deleted"},
	}

	for _, handler := range handlers {
		f.Comment(fmt.Sprintf("%s handles the %s state", handler.name, strings.ToLower(strings.TrimPrefix(handler.name, "Handle"))))
		f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id(handler.name).Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower("a"+resourceName)).Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName),
		).Params(
			jen.Qual(pkgCtrlState, "Result"),
			jen.Error(),
		).Block(
			jen.Comment("TODO: Implement "+strings.ToLower(strings.TrimPrefix(handler.name, "Handle"))+" state logic"),
			jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "NextState").Call(
				jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", handler.nextState),
				jen.Lit(handler.message),
			)),
		)
	}

	// For method
	f.Comment("For returns the resource and predicates for the controller")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("For").Params().Params(
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "Predicates"),
	).Block(
		jen.Id("obj").Op(":=").Op("&").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1", "Atlas"+resourceName).Values(),
		jen.Comment("TODO: Add appropriate predicates"),
		jen.Return(jen.Id("obj"), jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call()),
	)

	// SetupWithManager method
	f.Comment("SetupWithManager sets up the controller with the Manager")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("SetupWithManager").Params(
		jen.Id("mgr").Qual("sigs.k8s.io/controller-runtime", "Manager"),
		jen.Id("rec").Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Reconciler"),
		jen.Id("defaultOptions").Qual("sigs.k8s.io/controller-runtime/pkg/controller", "Options"),
	).Error().Block(
		jen.Id("h").Dot("Client").Op("=").Id("mgr").Dot("GetClient").Call(),
		jen.Return(jen.Qual("sigs.k8s.io/controller-runtime", "NewControllerManagedBy").Call(jen.Id("mgr")).
			Dot("Named").Call(jen.Lit("Atlas"+resourceName)).
			Dot("For").Call(jen.Id("h").Dot("For").Call()).
			Dot("WithOptions").Call(jen.Id("defaultOptions")).
			Dot("Complete").Call(jen.Id("rec"))),
	)
}

