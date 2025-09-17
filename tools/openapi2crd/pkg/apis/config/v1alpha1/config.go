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

/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	metav1.TypeMeta `json:",inline"`
	Spec            Spec `json:"spec"`
}

type Spec struct {
	CRDConfig          []CRDConfig         `json:"crd,omitempty"`
	OpenAPIDefinitions []OpenAPIDefinition `json:"openapi,omitempty"`
}

type OpenAPIDefinition struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Package string `json:"package"`
}

type CRDConfig struct {
	GVK        metav1.GroupVersionKind `json:"gvk,omitempty"`
	Categories []string                `json:"categories,omitempty"`
	Mappings   []CRDMapping            `json:"mappings,omitempty"`
	ShortNames []string                `json:"shortNames,omitempty"`
}

type CRDMapping struct {
	OpenAPIRef        LocalObjectReference `json:"openAPIRef,omitempty"`
	MajorVersion      string               `json:"majorVersion,omitempty"`
	ParametersMapping PropertyMapping      `json:"parameters,omitempty"`
	EntryMapping      PropertyMapping      `json:"entry,omitempty"`
	StatusMapping     PropertyMapping      `json:"status,omitempty"`
	Extensions        []Extension          `json:"extensions,omitempty"`
}

type Reference struct {
	Name     string `json:"name,omitempty"`     // Name of the reference
	Property string `json:"property,omitempty"` // The OpenAPI property to map to
	Target   Target `json:"target,omitempty"`   // The target CRD to map to
}

type Target struct {
	Type       Type     `json:"type,omitempty"`       // The GroupVersionResource of the target CRD.
	Properties []string `json:"properties,omitempty"` // The target CRD properties to map to.
}

type Type struct {
	Group    string `json:"group,omitempty"`
	Version  string `json:"version,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type PropertyMapping struct {
	Schema     string       `json:"schema,omitempty"`
	Path       PropertyPath `json:"path,omitempty"`
	Filters    Filters      `json:"filters,omitempty"`
	References []Reference  `json:"references,omitempty"`
}

type Extension struct {
	Property               string                          `json:"property,omitempty"`
	XKubernetesValidations apiextensionsv1.ValidationRules `json:"x-kubernetes-validations,omitempty"`
}

type XOpenApiMapping struct {
	Property string `json:"property,omitempty"`
	Type     string `json:"type,omitempty"`
}

type XKubernetesMapping struct {
	GVR              string   `json:"gvr,omitempty"`
	PropertySelector string   `json:"property-selector,omitempty"` // Selector in the referenced resource for the property to map to.
	Properties       []string `json:"properties,omitempty"`        // List of properties to map to. First available value wins.
	Property         string   `json:"property,omitempty"`          // Single property to map to.
}

type PropertyPath struct {
	Name        string      `json:"name,omitempty"`
	Verb        string      `json:"verb,omitempty"`
	RequestBody RequestBody `json:"requestBody,omitempty"`
}

type RequestBody struct {
	MimeType string `json:"mimeType,omitempty"`
}

type Filters struct {
	ReadOnly            bool     `json:"readOnly,omitempty"`
	ReadWriteOnly       bool     `json:"readWriteOnly,omitempty"`
	SkipProperties      []string `json:"skipProperties,omitempty"`
	SensitiveProperties []string `json:"sensitiveProperties,omitempty"`
}

// LocalObjectReference is a reference to an object in the same namespace as the referent
type LocalObjectReference struct {
	// Name of the resource being referred to
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name"`
}
