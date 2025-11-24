package generate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
)

func generateMainHandlerFile(dir, resourceName, typesPath string, mappings []MappingWithConfig, refsByKind map[string][]ReferenceField, config *ParsedConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

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
				// Capitalize first letter of version (e.g., v20250312 -> V20250312)
				capitalizedVersion := strings.ToUpper(string(versionSuffix[0])) + versionSuffix[1:]
				// Construct CRD name: {plural}.{group}
				crdName := config.PluralName + "." + config.CRDGroup
				sdkImportPathSplit := strings.Split(mapping.OpenAPIConfig.Package, "/")
				sdkVersionSuffix := strings.TrimPrefix(sdkImportPathSplit[len(sdkImportPathSplit)-2], "v")

				g.If(jen.Id(strings.ToLower(resourceName)).Dot("Spec").Dot(capitalizedVersion).Op("!=").Nil()).Block(

					jen.List(jen.Id("translationReq"), jen.Id("err")).Op(":=").Id("getTranslationRequest").Call(
						jen.Id("ctx"),
						jen.Id("h").Dot("Client"),
						jen.Lit(crdName),
						jen.Lit(config.StorageVersion),
						jen.Lit(versionSuffix),
					),
					jen.If(jen.Id("err").Op("!=").Nil()).Block(
						jen.Return(jen.Nil(), jen.Id("err")),
					),
					jen.Id("versionCount").Op("++"),
					jen.Id("selectedHandler").
						Op("=").Id("h").
						Dot("handler"+versionSuffix).
						Call(
							jen.Id("h").
								Dot("Client"),
							jen.Id("atlasClients").Dot("SdkClient"+sdkVersionSuffix),
							jen.Id("translationReq"),
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

	generateDelegatingStateHandlers(f, resourceName, apiPkg, refsByKind)
	// ClientSet and translation request helpers
	generateSDKClientSetMethod(f, resourceName, apiPkg)

	// Generate package-level helper function attached to the handler
	generatePackageLevelTranslationHelper(f)

	fileName := filepath.Join(dir, "handler.go")
	return f.Save(fileName)
}

func generateDelegatingStateHandlers(f *jen.File, resourceName, apiPkg string, refsByKind map[string][]ReferenceField) {
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

	generateMapperFunctions(f, resourceName, apiPkg, refsByKind)

	generateSetupWithManager(f, resourceName, refsByKind)
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

func generatePackageLevelTranslationHelper(f *jen.File) {
	f.Comment("getTranslationRequest creates a translation request for converting entities between API and AKO.")
	f.Comment("This is a package-level function that can be called from any handler.")
	f.Func().Id("getTranslationRequest").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("crdName").String(),
		jen.Id("storageVersion").String(),
		jen.Id("targetVersion").String(),
	).Params(
		jen.Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "Request"),
		jen.Error(),
	).Block(
		jen.Id("crd").Op(":=").Op("&").Qual("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1", "CustomResourceDefinition").Values(),
		jen.Id("err").Op(":=").Id("kubeClient").Dot("Get").Call(
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
