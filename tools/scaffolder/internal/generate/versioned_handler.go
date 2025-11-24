package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
)

func generateVersionHandlerFile(dir, resourceName, typesPath string, mapping MappingWithConfig, override bool) error {
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

	f := jen.NewFile(atlasResourceName)
	AddLicenseHeader(f)

	f.ImportAlias(pkgCtrlState, "ctrlstate")
	f.ImportAlias(apiPkg, "akov2generated")
	f.ImportAlias(sdkImportPath, versionSuffix+"sdk")

	f.Type().Id("Handler"+versionSuffix).Struct(
		jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("atlasClient").Op("*").Qual(sdkImportPath, "APIClient"),
		jen.Id("translationRequest").Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "Request"),
		jen.Id("deletionProtection").Bool(),
	)

	f.Func().Id("NewHandler"+versionSuffix).Params(
		jen.Id("kubeClient").Qual("sigs.k8s.io/controller-runtime/pkg/client", "Client"),
		jen.Id("atlasClient").Op("*").Qual(sdkImportPath, "APIClient"),
		jen.Id("translationRequest").Op("*").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate", "Request"),
		jen.Id("deletionProtection").Bool(),
	).Op("*").Id("Handler" + versionSuffix).Block(
		jen.Return(jen.Op("&").Id("Handler" + versionSuffix).Values(jen.Dict{
			jen.Id("kubeClient"):         jen.Id("kubeClient"),
			jen.Id("atlasClient"):        jen.Id("atlasClient"),
			jen.Id("translationRequest"): jen.Id("translationRequest"),
			jen.Id("deletionProtection"): jen.Id("deletionProtection"),
		})),
	)

	generateVersionStateHandlers(f, resourceName, apiPkg, versionSuffix)

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

	for _, handler := range handlers {
		f.Comment(fmt.Sprintf("%s handles the %s state for version %s", handler.name, strings.ToLower(strings.TrimPrefix(handler.name, "Handle")), versionSuffix))
		f.Func().Params(jen.Id("h").Op("*").Id("Handler"+versionSuffix)).Id(handler.name).Params(
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
