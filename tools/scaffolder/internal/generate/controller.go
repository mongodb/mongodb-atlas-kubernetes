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

	if err := generateMainHandlerFile(controllerDir, resourceName, typesPath, parsedConfig.Mappings, refsByKind, parsedConfig); err != nil {
		return fmt.Errorf("failed to generate main handler file: %w", err)
	}

	// Generate version-specific handlers
	for _, mapping := range parsedConfig.Mappings {
		if err := generateVersionHandlerFile(controllerDir, resourceName, typesPath, resultPath, mapping, override); err != nil {
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
	f.ImportAlias(apiPkg, "akov2generated")

	f.Const().Defs(
		jen.Comment("crdVersion of the controler"),
		jen.Id("crdVersion").Op("=").Lit("v1"),
	)

	supportedSDKVersions := []jen.Code{}
	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		supportedSDKVersions = append(supportedSDKVersions, jen.Lit(versionSuffix))
	}
	f.Var().Defs(
		jen.Comment("sdkVersions supported by this controller"),
		jen.Id("sdkVersions").Op("=").Op("[]").String().Values(supportedSDKVersions...),
	)

	resourcePlural := strings.ToLower(resourceName) + "s"
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=%s,verbs=get;list;watch;create;update;patch;delete", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=%s/status,verbs=get;update;patch", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,resources=%s/finalizers,verbs=update", resourcePlural))
	f.Comment("+kubebuilder:rbac:groups=\"\",resources=secrets,verbs=get;list;watch")
	f.Comment("+kubebuilder:rbac:groups=\"events.k8s.io\",resources=events,verbs=create;patch")
	f.Line()
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=%s,verbs=get;list;watch;create;update;patch;delete", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=%s/status,verbs=get;update;patch", resourcePlural))
	f.Comment(fmt.Sprintf("+kubebuilder:rbac:groups=atlas.generated.mongodb.com,namespace=default,resources=%s/finalizers,verbs=update", resourcePlural))
	f.Comment("+kubebuilder:rbac:groups=\"\",namespace=default,resources=secrets,verbs=get;list;watch")
	f.Comment("+kubebuilder:rbac:groups=\"events.k8s.io\",namespace=default,resources=events,verbs=create;patch")
	f.Line()

	// Handler struct for ALL CRD versions
	handlerFields := []jen.Code{
		jen.Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)),
		jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler"),
		jen.Id("deletionProtection").Bool(),
		jen.Id("predicates").Index().Qual("sigs.k8s.io/controller-runtime/pkg/predicate", "Predicate"),
	}

	// Version-specific handler for each version
	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		sdkImportPath := mapping.OpenAPIConfig.Package

		f.ImportAlias(sdkImportPath, versionSuffix+"sdk")

		handlerFields = append(
			handlerFields,
			jen.Id("translators").Map(jen.String()).Qual(
				"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi", "Translator",
			),
			jen.Id("handler"+versionSuffix).
				Qual(pkgCtrlState, "VersionedHandlerFunc").
				Types(
					jen.Qual(sdkImportPath, "APIClient"),
					jen.Qual(apiPkg, resourceName),
				),
		)
	}

	f.Type().Id("Handler").Struct(handlerFields...)

	// NewReconciler method with all service builders
	reconcilerParams := []jen.Code{
		jen.Line().Id("c").Qual("sigs.k8s.io/controller-runtime/pkg/cluster", "Cluster"),
		jen.Line().Id("atlasProvider").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "Provider"),
		jen.Line().Id("logger").Op("*").Qual("go.uber.org/zap", "Logger"),
		jen.Line().Id("globalSecretRef").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey"),
		jen.Line().Id("deletionProtection").Bool(),
		jen.Line().Id("reapplySupport").Bool(),
		jen.Line().Id("predicates").Index().Qual("sigs.k8s.io/controller-runtime/pkg/predicate", "Predicate"),
	}

	atlasReconcilerBase := jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "AtlasReconciler").Values(jen.Dict{
		jen.Id("Client"):          jen.Id("c").Dot("GetClient").Call(),
		jen.Id("AtlasProvider"):   jen.Id("atlasProvider"),
		jen.Id("Log"):             jen.Id("logger").Dot("Named").Call(jen.Lit("controllers")).Dot("Named").Call(jen.Lit("Atlas" + resourceName)).Dot("Sugar").Call(),
		jen.Id("GlobalSecretRef"): jen.Id("globalSecretRef"),
	})

	f.Func().Id("New"+controllerName+"Reconciler").Params(
		reconcilerParams...,
	).Params(
		jen.Op("*").Qual(pkgCtrlState, "Reconciler").Types(jen.Qual(apiPkg, resourceName)),
		jen.Error(),
	).Block(
		jen.List(jen.Id("crd"), jen.Id("err")).Op(":=").Qual(
			"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crds", "EmbeddedCRD",
		).Call(jen.Lit(resourceName)),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit(fmt.Sprintf("failed to read CRD for %s: %%w", resourceName)),
					jen.Id("err"),
				),
			),
		),
		jen.List(jen.Id("translators"), jen.Id("err")).Op(":=").Qual(
			"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi", "NewPerVersionTranslators",
		).Call(
			jen.Id("c").Dot("GetScheme").Call(),
			jen.Id("crd"),
			jen.Id("crdVersion"),
			jen.Id("sdkVersions").Op("..."),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(
				jen.Nil(),
				jen.Qual("fmt", "Errorf").Call(
					jen.Lit(fmt.Sprintf("failed to get translator set for %s: %%w", resourceName)),
					jen.Id("err"),
				),
			),
		),
		jen.Comment("Create main handler dispatcher"),
		jen.Id(strings.ToLower(controllerName)+"Handler").Op(":=").Op("&").Id("Handler").Values(jen.DictFunc(func(d jen.Dict) {
			d[jen.Id("AtlasReconciler")] = atlasReconcilerBase
			d[jen.Id("deletionProtection")] = jen.Id("deletionProtection")
			d[jen.Id("translators")] = jen.Id("translators")
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				d[jen.Id("handler"+versionSuffix)] = jen.Id("handler" + versionSuffix + "Func")
			}
			d[jen.Id("predicates")] = jen.Id("predicates")
		})),
		jen.Line(),
		jen.Return(jen.Qual(pkgCtrlState, "NewStateReconciler").Call(
			jen.Id(strings.ToLower(controllerName)+"Handler"),
			jen.Qual(pkgCtrlState, "WithCluster").Types(jen.Qual(apiPkg, resourceName)).Call(jen.Id("c")),
			jen.Qual(pkgCtrlState, "WithReapplySupport").Types(jen.Qual(apiPkg, resourceName)).Call(jen.Id("reapplySupport")),
		), jen.Nil()),
	)

	for _, mapping := range mappings {
		versionSuffix := mapping.Version
		sdkImportPath := mapping.OpenAPIConfig.Package

		handlerFuncParams := []jen.Code{
			jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
			jen.Id("atlasClient").Op("*").Qual(sdkImportPath, "APIClient"),
			jen.Id("translator").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi", "Translator"),
			jen.Id("deletionProtection").Bool(),
		}

		f.Func().
			Id("handler"+versionSuffix+"Func").
			Params(handlerFuncParams...).
			Qual(pkgCtrlState, "StateHandler").
			Types(jen.Qual(apiPkg, resourceName)).Block(
			jen.Return(
				jen.Id("NewHandler"+versionSuffix).Call(
					jen.Id("kubeClient"),
					jen.Id("atlasClient"),
					jen.Id("translator"),
					jen.Id("deletionProtection"),
				),
			),
		)
	}

	fileName := filepath.Join(dir, atlasResourceName+"_controller.go")
	return f.Save(fileName)
}

func generateSetupWithManager(f *jen.File, resourceName string, refsByKind map[string][]ReferenceField) {
	newControllerManagedBy := jen.Qual("sigs.k8s.io/controller-runtime", "NewControllerManagedBy").Call(jen.Id("mgr")).
		Dot("Named").Call(jen.Lit(resourceName)).
		Dot("For").Call(jen.Id("h").Dot("For").Call())

	for referencedKind := range refsByKind {
		newControllerManagedBy = newControllerManagedBy.Dot("Watches").Call(
			getWatchedTypeInstance(referencedKind),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/handler", "EnqueueRequestsFromMapFunc").
				Call(
					jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers", fmt.Sprintf("New%sBy%sMapFunc", resourceName, referencedKind)).
						Call(jen.Id("h").Dot("Client")),
				),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").
				Call(
					jen.Qual("sigs.k8s.io/controller-runtime/pkg/predicate", "ResourceVersionChangedPredicate").Values(),
				),
		)
	}

	newControllerManagedBy = newControllerManagedBy.
		Dot("WithOptions").Call(jen.Id("defaultOptions")).
		Dot("Complete").Call(jen.Id("rec"))

	f.Func().
		Params(
			jen.Id("h").Op("*").Id("Handler"),
		).
		Id("SetupWithManager").
		Params(
			jen.Id("mgr").Qual("sigs.k8s.io/controller-runtime", "Manager"),
			jen.Id("rec").Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Reconciler"),
			jen.Id("defaultOptions").Qual("sigs.k8s.io/controller-runtime/pkg/controller", "Options"),
		).
		Params(
			jen.Error(),
		).
		Block(
			jen.Id("h").Dot("Client").Op("=").Id("mgr").Dot("GetClient").Call(),
			jen.Return(newControllerManagedBy),
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
