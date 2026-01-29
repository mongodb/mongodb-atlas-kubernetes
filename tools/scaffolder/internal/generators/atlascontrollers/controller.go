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
)

const (
	pkgCtrlState = "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

func generateControllerFile(dir, resourceName, typesPath string, mappings []config.MappingWithConfig) error {
	atlasResourceName := strings.ToLower(resourceName)
	apiPkg := typesPath

	f := jen.NewFile(atlasResourceName)
	boilerplate.AddLicenseHeader(f)

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
				"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi", "Translator",
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

	f.Func().Id("New"+resourceName+"Reconciler").Params(
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
			"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi", "NewPerVersionTranslators",
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
		jen.Id(strings.ToLower(resourceName)+"Handler").Op(":=").Op("&").Id("Handler").Values(jen.DictFunc(func(d jen.Dict) {
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
			jen.Id(strings.ToLower(resourceName)+"Handler"),
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
			jen.Id("translator").Qual("github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi", "Translator"),
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
