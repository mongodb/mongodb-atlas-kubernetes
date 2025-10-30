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

func getAPIPackage(apiVersion string) string {
	return fmt.Sprintf("github.com/mongodb/mongodb-atlas-kubernetes/v2/api/%s", apiVersion)
}

// FromConfig generates controllers and handlers based on the parsed CRD result file
func FromConfig(resultPath, crdKind, controllerOutDir, indexerOutDir, typesPath string) error {
	parsedConfig, err := ParseCRDConfig(resultPath, crdKind)
	if err != nil {
		return err
	}

	resourceName := parsedConfig.ResourceName

	// Set default directories if not provided
	if controllerOutDir == "" {
		controllerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "controller")
	}

	// Generate indexers
	if indexerOutDir == "" {
		indexerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "indexer")
	}
	if err := GenerateIndexers(resultPath, crdKind, indexerOutDir); err != nil {
		return fmt.Errorf("failed to generate indexers: %w", err)
	}

	baseControllerDir := filepath.Join(controllerOutDir, strings.ToLower(resourceName))

	controllerName := resourceName
	controllerDir := baseControllerDir

	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("failed to create controller directory: %w", err)
	}

	if err := generateControllerFileWithMultipleVersions(controllerDir, controllerName, resourceName, typesPath, parsedConfig.Mappings); err != nil {
		return fmt.Errorf("failed to generate controller file: %w", err)
	}

	if err := generateMainHandlerFile(controllerDir, controllerName, resourceName, typesPath, parsedConfig.Mappings); err != nil {
		return fmt.Errorf("failed to generate main handler file: %w", err)
	}

	// Generate version-specific handlers
	for _, mapping := range parsedConfig.Mappings {
		if err := generateVersionHandlerFile(controllerDir, controllerName, resourceName, typesPath, mapping); err != nil {
			return fmt.Errorf("failed to generate handler for version %s: %w", mapping.Version, err)
		}
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

func generateControllerFileWithMultipleVersions(dir, controllerName, resourceName, typesPath string, mappings []MappingWithConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")

	resourcePlural := strings.ToLower(resourceName) + "s"
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,resources=%s,verbs=get;list;watch;create;update;patch;delete", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,resources=%s/status,verbs=get;update;patch", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,resources=%s/finalizers,verbs=update", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=%s,verbs=get;list;watch;create;update;patch;delete", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=%s/status,verbs=get;update;patch", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=%s/finalizers,verbs=update", resourcePlural))

	// Handler struct for ALL CRD versions
	handlerFields := []jen.Code{
		jen.Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)),
		jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler"),
	}

	// Version-specific handler for each version
	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		handlerFields = append(handlerFields, jen.Id("handler"+versionSuffix).Op("*").Id(controllerName+"Handler"+versionSuffix))
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

	atlasReconcilerBase := jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler").Values(jen.Dict{
		jen.Id("Client"):          jen.Id("c").Dot("GetClient").Call(),
		jen.Id("AtlasProvider"):   jen.Id("atlasProvider"),
		jen.Id("Log"):             jen.Id("logger").Dot("Named").Call(jen.Lit("controllers")).Dot("Named").Call(jen.Lit("Atlas" + resourceName)).Dot("Sugar").Call(),
		jen.Id("GlobalSecretRef"): jen.Id("globalSecretRef"),
	})

	f.Func().Id("New"+controllerName+"Reconciler").Params(reconcilerParams...).Op("*").Qual(pkgCtrlState, "Reconciler").Types(jen.Qual(apiPkg, resourceName)).Block(
		jen.Comment("Create version-specific handlers"),
		jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				g.Id("handler"+versionSuffix).Op(":=").Id("New"+controllerName+"Handler"+versionSuffix).Call(
					jen.Id("atlasProvider"),
					jen.Id("c").Dot("GetClient").Call(),
					jen.Id("logger").Dot("Named").Call(jen.Lit("controllers")).Dot("Named").Call(jen.Lit(resourceName+"-"+versionSuffix)).Dot("Sugar").Call(),
					jen.Id("globalSecretRef"),
				)
			}
		}),
		jen.Line(),
		jen.Comment("Create main handler dispatcher"),
		jen.Id(strings.ToLower(controllerName)+"Handler").Op(":=").Op("&").Id(controllerName+"Handler").Values(jen.DictFunc(func(d jen.Dict) {
			d[jen.Id("AtlasReconciler")] = atlasReconcilerBase
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				d[jen.Id("handler"+versionSuffix)] = jen.Id("handler" + versionSuffix)
			}
		})),
		jen.Line(),
		jen.Return(jen.Qual(pkgCtrlState, "NewStateReconciler").Call(
			jen.Id(strings.ToLower(controllerName)+"Handler"),
			jen.Qual(pkgCtrlState, "WithCluster").Types(jen.Qual(apiPkg, resourceName)).Call(jen.Id("c")),
			jen.Qual(pkgCtrlState, "WithReapplySupport").Types(jen.Qual(apiPkg, resourceName)).Call(jen.Id("reapplySupport")),
		)),
	)

	fileName := filepath.Join(dir, atlasResourceName+"_controller.go")
	return f.Save(fileName)
}

func generateMainHandlerFile(dir, controllerName, resourceName, typesPath string, mappings []MappingWithConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")

	f.Comment("getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("getHandlerForResource").Params(
		jen.Id(strings.ToLower(resourceName)).Op("*").Qual(apiPkg, resourceName),
	).Params(jen.Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)), jen.Error()).Block(
		jen.Comment("Check which resource spec version is set and validate that only one is specified"),
		jen.Var().Id("versionCount").Int(),
		jen.Var().Id("selectedHandler").Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)),
		jen.Line(),
		jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				// Capitalize first letter of version (e.g., v20250312 -> V20250312)
				capitalizedVersion := strings.ToUpper(string(versionSuffix[0])) + versionSuffix[1:]
				g.If(jen.Id(strings.ToLower(resourceName)).Dot("Spec").Dot(capitalizedVersion).Op("!=").Nil()).Block(
					jen.Id("versionCount").Op("++"),
					jen.Id("selectedHandler").Op("=").Id("h").Dot("handler"+versionSuffix),
				)
			}
		}),
		jen.Line(),
		jen.If(jen.Id("versionCount").Op("==").Lit(0)).Block(
			jen.Return(jen.Nil().Op(",").Qual("fmt", "Errorf").Call(jen.Lit("no resource spec version specified - please set one of the available spec versions"))),
		),
		jen.If(jen.Id("versionCount").Op(">").Lit(1)).Block(
			jen.Return(jen.Nil().Op(",").Qual("fmt", "Errorf").Call(jen.Lit("multiple resource spec versions specified - please set only one spec version"))),
		),
		jen.Return(jen.Id("selectedHandler").Op(",").Nil()),
	)

	generateDelegatingStateHandlers(f, controllerName, resourceName, apiPkg)

	fileName := filepath.Join(dir, "handler.go")
	return f.Save(fileName)
}

func generateDelegatingStateHandlers(f *jen.File, controllerName, resourceName, apiPkg string) {
	handlers := []string{
		"HandleInitial",
		"HandleImportRequested",
		"HandleImported",
		"HandleCreating",
		"HandleCreated",
		"HandleUpdating",
		"HandleUpdated",
		"HandleDeletionRequested",
		"HandleDeleting",
	}

	for _, handlerName := range handlers {
		f.Comment(fmt.Sprintf("%s delegates to the version-specific handler", handlerName))
		f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id(handlerName).Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(resourceName)).Op("*").Qual(apiPkg, resourceName),
		).Params(
			jen.Qual(pkgCtrlState, "Result"),
			jen.Error(),
		).Block(
			jen.List(jen.Id("handler"), jen.Id("err")).Op(":=").Id("h").Dot("getHandlerForResource").Call(jen.Id(strings.ToLower(resourceName))),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "Error").Call(
					jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", "StateInitial"),
					jen.Id("err"),
				)),
			),
			jen.Return(jen.Id("handler").Dot(handlerName).Call(jen.Id("ctx"), jen.Id(strings.ToLower(resourceName)))),
		)
	}

	f.Comment("For returns the resource and predicates for the controller")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("For").Params().Params(
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "Predicates"),
	).Block(
		jen.Id("obj").Op(":=").Op("&").Qual(apiPkg, resourceName).Values(),
		jen.Comment("TODO: Add appropriate predicates"),
		jen.Return(jen.Id("obj"), jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call()),
	)

	f.Comment("SetupWithManager sets up the controller with the Manager")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("SetupWithManager").Params(
		jen.Id("mgr").Qual("sigs.k8s.io/controller-runtime", "Manager"),
		jen.Id("rec").Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Reconciler"),
		jen.Id("defaultOptions").Qual("sigs.k8s.io/controller-runtime/pkg/controller", "Options"),
	).Error().Block(
		jen.Id("h").Dot("Client").Op("=").Id("mgr").Dot("GetClient").Call(),
		jen.Return(jen.Qual("sigs.k8s.io/controller-runtime", "NewControllerManagedBy").Call(jen.Id("mgr")).
			Dot("Named").Call(jen.Lit(resourceName)).
			Dot("For").Call(jen.Id("h").Dot("For").Call()).
			Dot("WithOptions").Call(jen.Id("defaultOptions")).
			Dot("Complete").Call(jen.Id("rec"))),
	)
}

func generateVersionHandlerFile(dir, controllerName, resourceName, typesPath string, mapping MappingWithConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	versionSuffix := mapping.Version
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")

	f.Type().Id(controllerName+"Handler"+versionSuffix).Struct(
		jen.Id("atlasProvider").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "Provider"),
		jen.Id("client").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("log").Op("*").Qual("go.uber.org/zap", "SugaredLogger"),
		jen.Id("globalSecretRef").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey"),
	)

	f.Func().Id("New"+controllerName+"Handler"+versionSuffix).Params(
		jen.Id("atlasProvider").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "Provider"),
		jen.Id("client").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("log").Op("*").Qual("go.uber.org/zap", "SugaredLogger"),
		jen.Id("globalSecretRef").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey"),
	).Op("*").Id(controllerName + "Handler" + versionSuffix).Block(
		jen.Return(jen.Op("&").Id(controllerName + "Handler" + versionSuffix).Values(jen.Dict{
			jen.Id("atlasProvider"):   jen.Id("atlasProvider"),
			jen.Id("client"):          jen.Id("client"),
			jen.Id("log"):             jen.Id("log"),
			jen.Id("globalSecretRef"): jen.Id("globalSecretRef"),
		})),
	)

	generateVersionStateHandlers(f, controllerName, resourceName, apiPkg, versionSuffix)

	// Generate For and SetupWithManager methods to satisfy StateHandler interface
	generateVersionInterfaceMethods(f, controllerName, resourceName, apiPkg, versionSuffix)

	fileName := filepath.Join(dir, "handler_"+versionSuffix+".go")
	return f.Save(fileName)
}

func generateVersionStateHandlers(f *jen.File, controllerName, resourceName, apiPkg, versionSuffix string) {
	handlers := []struct {
		name      string
		nextState string
		message   string
	}{
		{"HandleInitial", "StateUpdated", "Updated Atlas" + resourceName + "."},
		{"HandleImportRequested", "StateImported", "Import completed"},
		{"HandleImported", "StateUpdated", "Ready"},
		{"HandleCreating", "StateCreated", "Resource created"},
		{"HandleCreated", "StateUpdated", "Ready"},
		{"HandleUpdating", "StateUpdated", "Update completed"},
		{"HandleUpdated", "StateUpdated", "Ready"},
		{"HandleDeletionRequested", "StateDeleting", "Deletion started"},
		{"HandleDeleting", "StateDeleted", "Deleted"},
	}

	for _, handler := range handlers {
		f.Comment(fmt.Sprintf("%s handles the %s state for version %s", handler.name, strings.ToLower(strings.TrimPrefix(handler.name, "Handle")), versionSuffix))
		f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler"+versionSuffix)).Id(handler.name).Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(resourceName)).Op("*").Qual(apiPkg, resourceName),
		).Params(
			jen.Qual(pkgCtrlState, "Result"),
			jen.Error(),
		).Block(
			jen.Comment("TODO: Implement "+strings.ToLower(strings.TrimPrefix(handler.name, "Handle"))+" state logic"),
			jen.Comment("TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client"),
			jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "NextState").Call(
				jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", handler.nextState),
				jen.Lit(handler.message),
			)),
		)
	}
}

// generateVersionInterfaceMethods generates For and SetupWithManager methods for version-specific handlers
func generateVersionInterfaceMethods(f *jen.File, controllerName, resourceName, apiPkg, versionSuffix string) {
	// For method
	f.Comment("For returns the resource and predicates for the controller")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler"+versionSuffix)).Id("For").Params().Params(
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "Predicates"),
	).Block(
		jen.Return(jen.Op("&").Qual(apiPkg, resourceName).Values(), jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call()),
	)

	// SetupWithManager method
	f.Comment("SetupWithManager sets up the controller with the Manager")
	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler"+versionSuffix)).Id("SetupWithManager").Params(
		jen.Id("mgr").Qual("sigs.k8s.io/controller-runtime", "Manager"),
		jen.Id("rec").Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Reconciler"),
		jen.Id("defaultOptions").Qual("sigs.k8s.io/controller-runtime/pkg/controller", "Options"),
	).Error().Block(
		jen.Comment("This method is not used for version-specific handlers but required by StateHandler interface"),
		jen.Return(jen.Nil()),
	)
}
