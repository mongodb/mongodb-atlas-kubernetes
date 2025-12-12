package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
)

func generateVersionHandlerFile(dir, resourceName, typesPath, resultPath string, mapping MappingWithConfig, override bool) error {
	atlasResourceName := strings.ToLower(resourceName)
	versionSuffix := mapping.Version
	apiPkg := typesPath
	sdkImportPath := mapping.OpenAPIConfig.Package

	fileName := filepath.Join(dir, "handler_"+versionSuffix+".go")

	// Check if a versioned handler file exists
	if !override {
		if _, err := os.Stat(fileName); err == nil {
			fmt.Printf("Skipping versioned handler %s (already exists, use --override to overwrite)\n", fileName)
			return nil
		}
	}

	referenceFields, err := ParseReferenceFields(resultPath, resourceName)
	if err != nil {
		return fmt.Errorf("failed to parse reference fields: %w", err)
	}

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")
	f.ImportAlias(apiPkg, "akov2generated")
	f.ImportAlias(sdkImportPath, versionSuffix+"sdk")

	f.Type().Id("Handler"+versionSuffix).Struct(
		jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("atlasClient").Op("*").Qual(sdkImportPath, "APIClient"),
		jen.Id("translator").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi", "Translator"),
		jen.Id("deletionProtection").Bool(),
	)

	f.Func().Id("NewHandler"+versionSuffix).Params(
		jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("atlasClient").Op("*").Qual(sdkImportPath, "APIClient"),
		jen.Id("translator").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi", "Translator"),
		jen.Id("deletionProtection").Bool(),
	).Op("*").Id("Handler" + versionSuffix).Block(
		jen.Return(jen.Op("&").Id("Handler" + versionSuffix).Values(jen.Dict{
			jen.Id("kubeClient"):         jen.Id("kubeClient"),
			jen.Id("atlasClient"):        jen.Id("atlasClient"),
			jen.Id("translator"):         jen.Id("translator"),
			jen.Id("deletionProtection"): jen.Id("deletionProtection"),
		})),
	)

	generateVersionStateHandlers(f, resourceName, apiPkg, versionSuffix)

	// Generate getDependencies method to be used in for api translation calls
	generateGetDependenciesMethod(f, resourceName, apiPkg, versionSuffix, referenceFields)

	// Generate For and SetupWithManager methods to satisfy StateHandler interface
	generateVersionInterfaceMethods(f, resourceName, apiPkg, versionSuffix)

	return f.Save(fileName)
}

func generateVersionStateHandlers(f *jen.File, resourceName, apiPkg, versionSuffix string) {
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

	resourceVarName := strings.ToLower(resourceName)

	for _, handler := range handlers {
		f.Comment(fmt.Sprintf("%s handles the %s state for version %s", handler.name, strings.ToLower(strings.TrimPrefix(handler.name, "Handle")), versionSuffix))

		// Build the method body with getDependencies call
		methodBody := []jen.Code{
			jen.List(jen.Id("_"), jen.Err()).Op(":=").Id("h").Dot("getDependencies").Call(
				jen.Id("ctx"),
				jen.Id(resourceVarName),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(
					jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "Error").Call(
						jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", "State"+strings.TrimPrefix(handler.name, "Handle")),
						jen.Qual("fmt", "Errorf").Call(
							jen.Lit(fmt.Sprintf("failed to resolve %s dependencies: %%w", resourceName)),
							jen.Err(),
						),
					),
				),
			),
			jen.Line(),
			jen.Comment("TODO: Implement " + strings.ToLower(strings.TrimPrefix(handler.name, "Handle")) + " state logic"),
			jen.Comment("TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client"),
			jen.Comment("TODO: Replace _ with deps and use deps variable when calling h.translator.ToAPI() methods"),
			jen.Return(jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result", "NextState").Call(
				jen.Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state", handler.nextState),
				jen.Lit(handler.message),
			)),
		}

		f.Func().Params(jen.Id("h").Op("*").Id("Handler"+versionSuffix)).Id(handler.name).Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(resourceVarName).Op("*").Qual(apiPkg, resourceName),
		).Params(
			jen.Qual(pkgCtrlState, "Result"),
			jen.Error(),
		).Block(methodBody...)
	}
}

// generateGetDependenciesMethod generates the getDependencies method for the version-specific handler
func generateGetDependenciesMethod(f *jen.File, resourceName, apiPkg, versionSuffix string, referenceFields []ReferenceField) {
	resourceVarName := strings.ToLower(resourceName)

	blockStatements := []jen.Code{
		jen.Var().Id("deps").Index().Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Line(),
	}

	if len(referenceFields) == 0 {
		blockStatements = append(blockStatements, jen.Return(jen.Id("deps"), jen.Nil()))
	} else {
		for _, ref := range referenceFields {
			// No array-based refs
			if strings.Contains(ref.FieldPath, ".items.") {
				continue
			}

			fieldAccessPath := strings.Replace(buildFieldAccessPath(ref.FieldPath), "resource", resourceVarName, 1)

			// Build nil check conditions using RequiredSegments to handle optional pointer fields
			nilCheckCondition := buildNilCheckConditions(fieldAccessPath, ref.RequiredSegments)

			refKind := ref.ReferencedKind
			refVarName := strings.ToLower(refKind)

			// TODO: simplify?
			var refPkgQual *jen.Statement
			if refKind == "Secret" {
				refPkgQual = jen.Qual("k8s.io/api/core/v1", refKind)
			} else {
				refPkgQual = jen.Qual(apiPkg, refKind)
			}

			blockStatements = append(blockStatements,
				jen.Comment(fmt.Sprintf("Check if %s is present", ref.FieldName)),
				jen.If(nilCheckCondition).Block(
					jen.Id(refVarName).Op(":=").Op("&").Add(refPkgQual).Values(),
					jen.Err().Op(":=").Id("h").Dot("kubeClient").Dot("Get").Call(
						jen.Id("ctx"),
						jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "ObjectKey").Values(jen.Dict{
							jen.Id("Name"):      jen.Id(fieldAccessPath).Dot("Name"),
							jen.Id("Namespace"): jen.Id(resourceVarName).Dot("GetNamespace").Call(),
						}),
						jen.Id(refVarName),
					),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Id("deps"), jen.Qual("fmt", "Errorf").Call(
							jen.Lit(fmt.Sprintf("failed to get %s %%s/%%s: %%w", refKind)),
							jen.Id(resourceVarName).Dot("GetNamespace").Call(),
							jen.Id(fieldAccessPath).Dot("Name"),
							jen.Err(),
						)),
					),
					jen.Line(),
					jen.Id("deps").Op("=").Append(jen.Id("deps"), jen.Id(refVarName)),
				),
				jen.Line(),
			)
		}

		blockStatements = append(blockStatements, jen.Return(jen.Id("deps"), jen.Nil()))
	}

	f.Func().Params(jen.Id("h").Op("*").Id("Handler"+versionSuffix)).Id("getDependencies").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id(resourceVarName).Op("*").Qual(apiPkg, resourceName),
	).Params(
		jen.Index().Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Error(),
	).Block(blockStatements...)
}

// generateVersionInterfaceMethods generates For and SetupWithManager methods for version-specific handlers
func generateVersionInterfaceMethods(f *jen.File, resourceName, apiPkg, versionSuffix string) {
	// For method
	f.Comment("For returns the resource and predicates for the controller")
	f.Func().Params(jen.Id("h").Op("*").Id("Handler"+versionSuffix)).Id("For").Params().Params(
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/client", "Object"),
		jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "Predicates"),
	).Block(
		jen.Return(jen.Op("&").Qual(apiPkg, resourceName).Values(), jen.Qual("sigs.k8s.io/controller-runtime/pkg/builder", "WithPredicates").Call()),
	)

	// SetupWithManager method
	f.Comment("SetupWithManager sets up the controller with the Manager")
	f.Func().Params(jen.Id("h").Op("*").Id("Handler"+versionSuffix)).Id("SetupWithManager").Params(
		jen.Id("mgr").Qual("sigs.k8s.io/controller-runtime", "Manager"),
		jen.Id("rec").Qual("sigs.k8s.io/controller-runtime/pkg/reconcile", "Reconciler"),
		jen.Id("defaultOptions").Qual("sigs.k8s.io/controller-runtime/pkg/controller", "Options"),
	).Error().Block(
		jen.Comment("This method is not used for version-specific handlers but required by StateHandler interface"),
		jen.Return(jen.Nil()),
	)
}
