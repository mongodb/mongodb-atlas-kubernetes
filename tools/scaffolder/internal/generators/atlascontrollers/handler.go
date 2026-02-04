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

package atlascontrollers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/boilerplate"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/indexers"
)

func generateMainHandlerFile(dir, resourceName, typesPath, indexerImportPath string, mappings []config.MappingWithConfig, refsByKind map[string][]indexers.ReferenceField, _ *config.ParsedConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	boilerplate.AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")
	f.ImportAlias("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1", "apiextensionsv1")
	f.ImportAlias(apiPkg, "akov2generated")

	f.Comment("getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set")
	f.Func().Params(jen.Id("h").Op("*").Id("Handler")).Id("getHandlerForResource").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(atlasResourceName).Op("*").Qual(apiPkg, resourceName),
	).Params(jen.Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)), jen.Error()).Block(
		jen.List(jen.Id("atlasClients"), jen.Id("err")).Op(":=").Id("h").Dot("getSDKClientSet").Call(
			jen.Id("ctx"),
			jen.Id(atlasResourceName),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Id("err")),
		),

		jen.Comment("Check which resource spec version is set and validate that only one is specified"),
		jen.Var().Id("versionCount").Int(),
		jen.Var().Id("selectedHandler").Qual(pkgCtrlState, "StateHandler").Types(jen.Qual(apiPkg, resourceName)),
		jen.Line(),
		jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
			for _, mapping := range mappings {
				versionSuffix := mapping.Version
				capitalizedVersion := strings.ToUpper(string(versionSuffix[0])) + versionSuffix[1:]

				sdkImportPathSplit := strings.Split(mapping.OpenAPIConfig.Package, "/")
				sdkVersionSuffix := strings.TrimPrefix(sdkImportPathSplit[len(sdkImportPathSplit)-2], "v")

				g.If(jen.Id(strings.ToLower(resourceName)).Dot("Spec").Dot(capitalizedVersion).Op("!=").Nil()).Block(

					jen.List(jen.Id("translator"), jen.Id("ok")).Op(":=").Id("h").Dot("translators").Op("[").Lit(versionSuffix).Op("]"),
					jen.If(jen.Id("ok").Op("!=").True()).Block(
						jen.Return(
							jen.Nil(),
							jen.Qual("errors", "New").Call(
								jen.Lit(fmt.Sprintf("unsupported version %s set in CR", versionSuffix)),
							)),
					),
					jen.Id("versionCount").Op("++"),
					jen.Id("selectedHandler").
						Op("=").Id("h").
						Dot("handler"+versionSuffix).
						Call(
							jen.Id("h").
								Dot("Client"),
							jen.Id("atlasClients").Dot("SdkClient"+sdkVersionSuffix),
							jen.Id("translator"),
							jen.Id("h").Dot("deletionProtection"),
						),
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

	generateDelegatingStateHandlers(f, resourceName, apiPkg, indexerImportPath, refsByKind)
	generateSDKClientSetMethod(f, resourceName, apiPkg)

	fileName := filepath.Join(dir, "handler.go")
	return f.Save(fileName)
}

func generateDelegatingStateHandlers(f *jen.File, resourceName, apiPkg, indexerImportPath string, refsByKind map[string][]indexers.ReferenceField) {
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
		f.Func().Params(jen.Id("h").Op("*").Id("Handler")).Id(handlerName).Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(resourceName)).Op("*").Qual(apiPkg, resourceName),
		).Params(
			jen.Qual(pkgCtrlState, "Result"),
			jen.Error(),
		).Block(
			jen.List(jen.Id("handler"), jen.Id("err")).
				Op(":=").
				Id("h").
				Dot("getHandlerForResource").
				Call(jen.Id("ctx"), jen.Id(strings.ToLower(resourceName))),
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
	f.Func().Params(jen.Id("h").Op("*").Id("Handler")).Id("For").Params().Params(
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "Predicates"),
	).Block(
		jen.Id("obj").Op(":=").Op("&").Qual(apiPkg, resourceName).Values(),
		jen.Return(
			jen.Id("obj"),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call(jen.Id("h").Dot("predicates").Op("...")),
		),
	)

	generateSetupWithManager(f, resourceName, apiPkg, indexerImportPath, refsByKind)
}

func generateSetupWithManager(f *jen.File, resourceName, typesPath, indexerImportPath string, refsByKind map[string][]indexers.ReferenceField) {
	newControllerManagedBy := jen.Qual("sigs.k8s.io/controller-runtime", "NewControllerManagedBy").Call(jen.Id("mgr")).
		Dot("Named").Call(jen.Lit(resourceName)).
		Dot("For").Call(jen.Id("h").Dot("For").Call())

	for referencedKind := range refsByKind {
		newControllerManagedBy = newControllerManagedBy.Dot("Watches").Call(
			getWatchedTypeInstance(referencedKind, typesPath),
			jen.Qual("sigs.k8s.io/controller-runtime/pkg/handler", "EnqueueRequestsFromMapFunc").
				Call(
					jen.Qual(indexerImportPath, fmt.Sprintf("New%sBy%sMapFunc", resourceName, referencedKind)).
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

func getWatchedTypeInstance(kind, typesPath string) *jen.Statement {
	switch kind {
	case "Secret":
		return jen.Op("&").Qual("k8s.io/api/core/v1", "Secret").Values()
	default:
		return jen.Op("&").Qual(typesPath, kind).Values()
	}
}

func generateSDKClientSetMethod(f *jen.File, resourceName, apiPkg string) {
	resourceLower := strings.ToLower(resourceName)

	f.Comment("getSDKClientSet creates an Atlas SDK client set using credentials from the resource's connection secret")
	f.Func().Params(
		jen.Id("h").Op("*").Id("Handler"),
	).Id("getSDKClientSet").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(resourceLower).Op("*").Qual(apiPkg, resourceName),
	).Params(
		jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas", "ClientSet"),
		jen.Error(),
	).Block(
		jen.Var().Id("connectionSecretRef").Op("*").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey"),
		jen.If(jen.Id(resourceLower).Dot("Spec").Dot("ConnectionSecretRef").Op("!=").Nil()).Block(
			jen.Id("connectionSecretRef").Op("=").Op("&").Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey").Values(jen.Dict{
				jen.Id("Name"):      jen.Id(resourceLower).Dot("Spec").Dot("ConnectionSecretRef").Dot("Name"),
				jen.Id("Namespace"): jen.Id(resourceLower).Dot("Namespace"),
			}),
		),
		jen.Line(),

		jen.List(jen.Id("connectionConfig"), jen.Id("err")).Op(":=").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler", "GetConnectionConfig").Call(
			jen.Id("ctx"),
			jen.Id("h").Dot("Client"),
			jen.Id("connectionSecretRef"),
			jen.Op("&").Id("h").Dot("GlobalSecretRef"),
		),
		jen.If(jen.Id("err").Op("!=").Nil()).Block(
			jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(
				jen.Lit("failed to resolve Atlas credentials: %w"),
				jen.Id("err"),
			)),
		),
		jen.Line(),
		jen.List(jen.Id("clientSet"), jen.Id("err")).Op(":=").Id("h").Dot("AtlasProvider").Dot("SdkClientSet").Call(
			jen.Id("ctx"),
			jen.Id("connectionConfig").Dot("Credentials"),
			jen.Id("h").Dot("Log"),
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
