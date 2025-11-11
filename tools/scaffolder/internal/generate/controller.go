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
func FromConfig(resultPath, crdKind, controllerOutDir, indexerOutDir, typesPath string, override bool) error {
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

	// Parse reference fields for watch generation
	referenceFields, err := ParseReferenceFields(resultPath, crdKind)
	if err != nil {
		return fmt.Errorf("failed to parse reference fields: %w", err)
	}

	// Group references by target kind
	refsByKind := make(map[string][]ReferenceField)
	for _, ref := range referenceFields {
		refsByKind[ref.ReferencedKind] = append(refsByKind[ref.ReferencedKind], ref)
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

	if err := generateMainHandlerFile(controllerDir, controllerName, resourceName, typesPath, parsedConfig.Mappings, refsByKind, parsedConfig); err != nil {
		return fmt.Errorf("failed to generate main handler file: %w", err)
	}

	// Generate version-specific handlers
	for _, mapping := range parsedConfig.Mappings {
		if err := generateVersionHandlerFile(controllerDir, controllerName, resourceName, typesPath, mapping, override, parsedConfig); err != nil {
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
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=%s,verbs=get;list;watch;create;update;patch;delete", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=%s/status,verbs=get;update;patch", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=%s/finalizers,verbs=update", resourcePlural))
	f.Comment("+kubebuilder:rbac:groups=\"\",resources=secrets,verbs=get;list;watch")
	f.Comment("+kubebuilder:rbac:groups=\"\",resources=events,verbs=create;patch")
	f.Line()
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=%s,verbs=get;list;watch;create;update;patch;delete", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=%s/status,verbs=get;update;patch", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=%s/finalizers,verbs=update", resourcePlural))
	f.Comment("+kubebuilder:rbac:groups=\"\",namespace=default,resources=secrets,verbs=get;list;watch")
	f.Comment("+kubebuilder:rbac:groups=\"\",namespace=default,resources=events,verbs=create;patch")
	f.Line()

	// Handler struct for ALL CRD versions
	handlerFields := []jen.Code{
		jen.Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)),
		jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler"),
		jen.Id("predicates").Index().Qual("sigs.k8s.io/controller-runtime/pkg/predicate", "Predicate"),
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
		jen.Id("predicates").Index().Qual("sigs.k8s.io/controller-runtime/pkg/predicate", "Predicate"),
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
			d[jen.Id("predicates")] = jen.Id("predicates")
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

func generatePackageLevelTranslationHelper(f *jen.File) {
	f.Comment("getTranslationRequest creates a translation request for converting entities between API and AKO.")
	f.Comment("This is a package-level function that can be called from any handler.")
	f.Func().Id("getTranslationRequest").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("k8sClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("crdName").String(),
		jen.Id("storageVersion").String(),
		jen.Id("targetVersion").String(),
	).Params(
		jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "Request"),
		jen.Error(),
	).Block(
		jen.Id("crd").Op(":=").Op("&").Qual("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1", "CustomResourceDefinition").Values(),
		jen.Id("err").Op(":=").Id("k8sClient").Dot("Get").Call(
			jen.Id("ctx"),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey").Values(jen.Dict{
				jen.Id("Name"): jen.Id("crdName"),
			}),
			jen.Id("crd"),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("failed to resolve CRD %s: %w"),
				jen.Id("crdName"),
				jen.Id("err"),
			)),
		),
		jen.Line(),
		jen.List(jen.Id("translator"), jen.Id("err")).Op(":=").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "NewTranslator").Call(
			jen.Id("crd"),
			jen.Id("storageVersion"),
			jen.Id("targetVersion"),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("failed to setup translator: %w"),
				jen.Id("err"),
			)),
		),
		jen.Line(),
		jen.Return(jen.Op("&").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "Request").Values(jen.Dict{
			jen.Id("Translator"):   jen.Id("translator"),
			jen.Id("Dependencies"): jen.Nil(),
		}), jen.Nil()),
	)
}

func generateMainHandlerFile(dir, controllerName, resourceName, typesPath string, mappings []MappingWithConfig, refsByKind map[string][]ReferenceField, config *ParsedConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")

	// Generate package-level helper function attached to the handler
	generatePackageLevelTranslationHelper(f)

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

	generateDelegatingStateHandlers(f, controllerName, resourceName, apiPkg, refsByKind)

	fileName := filepath.Join(dir, "handler.go")
	return f.Save(fileName)
}

func generateDelegatingStateHandlers(f *jen.File, controllerName, resourceName, apiPkg string, refsByKind map[string][]ReferenceField) {
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
	startStateMap := map[string]string{
		"HandleInitial":           "StateInitial",
		"HandleImportRequested":   "StateImportRequested",
		"HandleImported":          "StateImported",
		"HandleCreating":          "StateCreating",
		"HandleCreated":           "StateCreated",
		"HandleUpdating":          "StateUpdating",
		"HandleUpdated":           "StateUpdated",
		"HandleDeletionRequested": "StateDeletionRequested",
		"HandleDeleting":          "StateDeleting",
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
					jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", startStateMap[handlerName]),
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
		jen.Return(
			jen.Id("obj"),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call(jen.Id("h").Dot("predicates").Op("...")),
		),
	)

	generateMapperFunctions(f, controllerName, resourceName, apiPkg, refsByKind)

	generateSetupWithManager(f, controllerName, resourceName, refsByKind)
}

func generateTranslationRequestWrapper(f *jen.File, controllerName, versionSuffix string, config *ParsedConfig) {
	// Construct CRD name: {plural}.{group}
	crdName := fmt.Sprintf("%s.%s", config.PluralName, config.CRDGroup)

	f.Comment("getTranslationRequest is a convenience wrapper for the package-level getTranslationRequest function")
	f.Func().Params(
		jen.Id("h").Op("*").Id(controllerName+"Handler"+versionSuffix),
	).Id("getTranslationRequest").Params(
		jen.Id("ctx").Qual("context", "Context"),
	).Params(
		jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "Request"),
		jen.Error(),
	).Block(
		jen.Return(jen.Id("getTranslationRequest").Call(
			jen.Id("ctx"),
			jen.Id("h").Dot("client"),
			jen.Lit(crdName),
			jen.Lit(config.StorageVersion),
			jen.Lit(versionSuffix),
		)),
	)
}

func generateSDKClientSetMethod(f *jen.File, controllerName, resourceName, apiPkg, versionSuffix string) {
	resourceLower := strings.ToLower(resourceName)

	f.Comment("getSDKClientSet creates an Atlas SDK client set using credentials from the resource's connection secret")
	f.Func().Params(
		jen.Id("h").Op("*").Id(controllerName+"Handler"+versionSuffix),
	).Id("getSDKClientSet").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(resourceLower).Op("*").Qual(apiPkg, resourceName),
	).Params(
		jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "ClientSet"),
		jen.Error(),
	).Block(
		jen.List(jen.Id("connectionConfig"), jen.Id("err")).Op(":=").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "GetConnectionConfig").Call(
			jen.Id("ctx"),
			jen.Id("h").Dot("client"),
			jen.Op("&").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey").Values(jen.Dict{
				jen.Id("Namespace"): jen.Id(resourceLower).Dot("Namespace"),
				jen.Id("Name"):      jen.Id(resourceLower).Dot("Spec").Dot("ConnectionSecretRef").Dot("Name"),
			}),
			jen.Op("&").Id("h").Dot("globalSecretRef"),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("failed to resolve Atlas credentials: %w"),
				jen.Id("err"),
			)),
		),
		jen.Line(),
		jen.List(jen.Id("clientSet"), jen.Id("err")).Op(":=").Id("h").Dot("atlasProvider").Dot("SdkClientSet").Call(
			jen.Id("ctx"),
			jen.Id("connectionConfig").Dot("Credentials"),
			jen.Id("h").Dot("log"),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("failed to setup Atlas SDK client: %w"),
				jen.Id("err"),
			)),
		),
		jen.Line(),
		jen.Return(jen.Id("clientSet"), jen.Nil()),
	)
}

func generateVersionHandlerFile(dir, controllerName, resourceName, typesPath string, mapping MappingWithConfig, override bool, config *ParsedConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	versionSuffix := mapping.Version
	apiPkg := typesPath

	fileName := filepath.Join(dir, "handler_"+versionSuffix+".go")

	// Check if a versioned handler file exists
	if !override {
		if _, err := os.Stat(fileName); err == nil {
			fmt.Printf("Skipping versioned handler %s (already exists, use --override to overwrite)\n", fileName)
			return nil
		}
	}

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

	// ClientSet and translation request helpers
	generateSDKClientSetMethod(f, controllerName, resourceName, apiPkg, versionSuffix)
	generateTranslationRequestWrapper(f, controllerName, versionSuffix, config)

	generateVersionStateHandlers(f, controllerName, resourceName, apiPkg, versionSuffix)

	// Generate For and SetupWithManager methods to satisfy StateHandler interface
	generateVersionInterfaceMethods(f, controllerName, resourceName, apiPkg, versionSuffix)

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

func generateSetupWithManager(f *jen.File, controllerName, resourceName string, refsByKind map[string][]ReferenceField) {
	f.Comment("SetupWithManager sets up the controller with the Manager")

	setupChain := jen.Qual("sigs.k8s.io/controller-runtime", "NewControllerManagedBy").Call(jen.Id("mgr")).
		Dot("Named").Call(jen.Lit(resourceName)).
		Dot("For").Call(jen.Id("h").Dot("For").Call())

	// Add Watches()
	for kind := range refsByKind {
		watchedTypeInstance := getWatchedTypeInstance(kind)
		mapperFuncName := fmt.Sprintf("%sFor%sMapFunc", strings.ToLower(resourceName), kind)

		setupChain = setupChain.
			Dot("Watches").Call(
			watchedTypeInstance,
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/handler", "EnqueueRequestsFromMapFunc").Call(
				jen.Id("h").Dot(mapperFuncName).Call(),
			),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call(
				jen.Qual("sigs.k8s.io/controller-runtime/pkg/predicate", "ResourceVersionChangedPredicate").Values(),
			),
		)
	}

	setupChain = setupChain.
		Dot("WithOptions").Call(jen.Id("defaultOptions")).
		Dot("Complete").Call(jen.Id("rec"))

	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id("SetupWithManager").Params(
		jen.Id("mgr").Qual("sigs.k8s.io/controller-runtime", "Manager"),
		jen.Id("rec").Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Reconciler"),
		jen.Id("defaultOptions").Qual("sigs.k8s.io/controller-runtime/pkg/controller", "Options"),
	).Error().Block(
		jen.Id("h").Dot("Client").Op("=").Id("mgr").Dot("GetClient").Call(),
		jen.Return(setupChain),
	)
}

func generateMapperFunctions(f *jen.File, controllerName, resourceName, apiPkg string, refsByKind map[string][]ReferenceField) {
	for kind, refs := range refsByKind {
		if len(refs) == 0 {
			continue
		}

		indexerType := refs[0].IndexerType
		mapperFuncName := fmt.Sprintf("%sFor%sMapFunc", strings.ToLower(resourceName), kind)

		switch indexerType {
		case "project":
			generateIndexerBasedMapperFunction(f, controllerName, resourceName, apiPkg, kind, mapperFuncName, "ProjectsIndexMapperFunc")
		case "credentials":
			generateIndexerBasedMapperFunction(f, controllerName, resourceName, apiPkg, kind, mapperFuncName, "CredentialsIndexMapperFunc")
		case "resource":
			generateResourceMapperFunction(f, controllerName, resourceName, apiPkg, kind, mapperFuncName, refs)
		}
	}
}

// For Group or Secrets
func generateIndexerBasedMapperFunction(f *jen.File, controllerName, resourceName, apiPkg, referencedKind, mapperFuncName, indexerHelperFunc string) {
	indexName := fmt.Sprintf("%sBy%sIndex", resourceName, referencedKind)
	listTypeName := fmt.Sprintf("%sList", resourceName)
	requestsFuncName := fmt.Sprintf("%sRequestsFrom%s", resourceName, referencedKind)

	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id(mapperFuncName).Params().Qual("sigs.k8s.io/controller-runtime/pkg/handler", "MapFunc").Block(
		jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer", indexerHelperFunc).Call(
			jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexer", indexName),
			jen.Func().Params().Op("*").Qual(apiPkg, listTypeName).Block(
				jen.Return(jen.Op("&").Qual(apiPkg, listTypeName).Values()),
			),
			jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexer", requestsFuncName),
			jen.Id("h").Dot("Client"),
			jen.Id("h").Dot("Log"),
		)),
	)
}

func generateResourceMapperFunction(f *jen.File, controllerName, resourceName, apiPkg, referencedKind, mapperFuncName string, refs []ReferenceField) {
	indexName := fmt.Sprintf("%sBy%sIndex", resourceName, referencedKind)
	listTypeName := fmt.Sprintf("%sList", resourceName)
	watchedType := getWatchedType(referencedKind)

	f.Func().Params(jen.Id("h").Op("*").Id(controllerName+"Handler")).Id(mapperFuncName).Params().Qual("sigs.k8s.io/controller-runtime/pkg/handler", "MapFunc").Block(
		jen.Return(jen.Func().Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("obj").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		).Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").Block(

			jen.List(jen.Id("refObj"), jen.Id("ok")).Op(":=").Id("obj").Assert(jen.Op("*").Add(watchedType)),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Id("h").Dot("Log").Dot("Warnf").Call(
					jen.Lit(fmt.Sprintf("watching %s but got %%T", referencedKind)),
					jen.Id("obj"),
				),
				jen.Return(jen.Nil()),
			),

			jen.Id("listOpts").Op(":=").Op("&").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ListOptions").Values(jen.Dict{
				jen.Id("FieldSelector"): jen.Qual("k8s.io/apimachinery/pkg/fields", "OneTermEqualSelector").Call(
					jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexer", indexName),
					jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKeyFromObject").Call(jen.Id("refObj")).Dot("String").Call(),
				),
			}),
			jen.Id("list").Op(":=").Op("&").Qual(apiPkg, listTypeName).Values(),
			jen.Id("err").Op(":=").Id("h").Dot("Client").Dot("List").Call(jen.Id("ctx"), jen.Id("list"), jen.Id("listOpts")),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Id("h").Dot("Log").Dot("Errorf").Call(
					jen.Lit(fmt.Sprintf("failed to list from indexer %s: %%v", indexName)),
					jen.Id("err"),
				),
				jen.Return(jen.Nil()),
			),

			jen.Id("requests").Op(":=").Make(jen.Index().Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request"), jen.Lit(0), jen.Len(jen.Id("list").Dot("Items"))),
			jen.For(jen.List(jen.Id("_"), jen.Id("item")).Op(":=").Range().Id("list").Dot("Items")).Block(
				jen.Id("requests").Op("=").Append(jen.Id("requests"), jen.Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Request").Values(jen.Dict{
					jen.Id("NamespacedName"): jen.Qual("k8s.io/apimachinery/pkg/types", "NamespacedName").Values(jen.Dict{
						jen.Id("Name"):      jen.Id("item").Dot("Name"),
						jen.Id("Namespace"): jen.Id("item").Dot("Namespace"),
					}),
				})),
			),
			jen.Return(jen.Id("requests")),
		)),
	)
}

func getWatchedTypeInstance(kind string) *jen.Statement {
	switch kind {
	case "Secret":
		return jen.Op("&").Qual("k8s.io/api/core/v1", "Secret").Values()
	case "Group":
		return jen.Op("&").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", "Group").Values()
	default:
		return jen.Op("&").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", kind).Values()
	}
}

func getWatchedType(kind string) *jen.Statement {
	switch kind {
	case "Secret":
		return jen.Qual("k8s.io/api/core/v1", "Secret")
	case "Group":
		return jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", "Group")
	default:
		return jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", kind)
	}
}
