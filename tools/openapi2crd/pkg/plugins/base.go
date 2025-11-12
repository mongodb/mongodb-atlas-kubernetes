// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package plugins

import (
	"fmt"
	"strings"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

// Base is a plugin that add minimum required configuration to a CRD
type Base struct{}

func (p *Base) Name() string {
	return "base"
}

func (p *Base) Process(req *CRDProcessorRequest) error {
	pluralGvk, singularGvk := guessKindToResource(req.CRDConfig.GVK)

	req.CRD.ObjectMeta = v1.ObjectMeta{
		Name: fmt.Sprintf("%s.%s", pluralGvk.Resource, pluralGvk.Group),
	}

	req.CRD.Spec = apiextensions.CustomResourceDefinitionSpec{
		Group: pluralGvk.Group,
		Scope: apiextensions.NamespaceScoped,
		Names: apiextensions.CustomResourceDefinitionNames{
			Kind:     req.CRDConfig.GVK.Kind,
			ListKind: fmt.Sprintf("%sList", req.CRDConfig.GVK.Kind),
			Plural:   pluralGvk.Resource,
			Singular: singularGvk.Resource,
		},
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    req.CRDConfig.GVK.Version,
				Served:  true,
				Storage: true,
			},
		},
		PreserveUnknownFields: ptr.To(false),
		Validation: &apiextensions.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: fmt.Sprintf("A %v, managed by the MongoDB Kubernetes Atlas Operator.", singularGvk.Resource),
				Properties: map[string]apiextensions.JSONSchemaProps{
					"spec": {
						Type: "object",
						Description: fmt.Sprintf(`Specification of the %v supporting the following versions:

%v

At most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource, strings.Join(majorVersions(req.CRDConfig), "\n")),
						Properties: map[string]apiextensions.JSONSchemaProps{},
					},
					"status": {
						Type:        "object",
						Description: fmt.Sprintf(`Most recently observed read-only status of the %v for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource),
						Properties:  map[string]apiextensions.JSONSchemaProps{},
					},
				},
			},
		},
	}

	req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties["conditions"] = apiextensions.JSONSchemaProps{
		Type:        "array",
		Description: "Represents the latest available observations of a resource's current state.",
		Items: &apiextensions.JSONSchemaPropsOrArray{
			Schema: &apiextensions.JSONSchemaProps{
				Type:     "object",
				Required: []string{"type", "status"},
				Properties: map[string]apiextensions.JSONSchemaProps{
					"type":               {Type: "string", Description: "Type of condition."},
					"status":             {Type: "string", Description: "Status of the condition, one of True, False, Unknown."},
					"observedGeneration": {Type: "integer", Description: "observedGeneration represents the .metadata.generation that the condition was set based upon."},
					"message":            {Type: "string", Description: "A human readable message indicating details about the transition."},
					"reason":             {Type: "string", Description: "The reason for the condition's last transition."},
					"lastTransitionTime": {Type: "string", Format: "date-time", Description: "Last time the condition transitioned from one status to another."},
				},
			},
		},
		XListMapKeys: []string{
			"type",
		},
		XListType: ptr.To("map"),
	}

	req.CRD.Status.StoredVersions = []string{}

	// enable status subresource
	req.CRD.Spec.Subresources = &apiextensions.CustomResourceSubresources{
		Status: &apiextensions.CustomResourceSubresourceStatus{},
	}

	req.CRD.Spec.Names.Categories = req.CRDConfig.Categories
	req.CRD.Spec.Names.ShortNames = req.CRDConfig.ShortNames

	for _, version := range req.CRD.Spec.Versions {
		if version.Storage {
			req.CRD.Status.StoredVersions = append(req.CRD.Status.StoredVersions, version.Name)
		}
	}

	return nil
}

func guessKindToResource(gvk v1.GroupVersionKind) ( /*plural*/ runtimeschema.GroupVersionResource /*singular*/, runtimeschema.GroupVersionResource) {
	runtimeGVK := runtimeschema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}
	kindName := runtimeGVK.Kind
	if len(kindName) == 0 {
		return runtimeschema.GroupVersionResource{}, runtimeschema.GroupVersionResource{}
	}
	singularName := strings.ToLower(kindName)
	singular := runtimeGVK.GroupVersion().WithResource(singularName)

	switch string(singularName[len(singularName)-1]) {
	case "s":
		return runtimeGVK.GroupVersion().WithResource(singularName + "es"), singular
	case "x":
		return runtimeGVK.GroupVersion().WithResource(singularName + "es"), singular
	case "y":
		return runtimeGVK.GroupVersion().WithResource(strings.TrimSuffix(singularName, "y") + "ies"), singular
	}

	return runtimeGVK.GroupVersion().WithResource(singularName + "s"), singular
}

func majorVersions(crdConfig *configv1alpha1.CRDConfig) []string {
	result := make([]string, 0, len(crdConfig.Mappings))
	for _, m := range crdConfig.Mappings {
		result = append(result, "- "+m.MajorVersion)
	}
	return result
}
