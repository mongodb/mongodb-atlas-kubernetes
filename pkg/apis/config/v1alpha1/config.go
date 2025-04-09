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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	v1.TypeMeta `json:",inline"`
	Spec        Spec `json:"spec"`
}

type Spec struct {
	CRDConfig          []CRDConfig         `json:"crd,omitempty"`
	OpenAPIDefinitions []OpenAPIDefinition `json:"openapi,omitempty"`
}

type OpenAPIDefinition struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CRDConfig struct {
	GVK        v1.GroupVersionKind `json:"gvk,omitempty"`
	Categories []string            `json:"categories,omitempty"`
	Mappings   []Mapping           `json:"mappings,omitempty"`
}

type Mapping struct {
	OpenAPIRef      LocalObjectReference `json:"openAPIRef,omitempty"`
	MajorVersion    string               `json:"majorVersion,omitempty"`
	Path            string               `json:"path,omitempty"`
	Verb            string               `json:"verb,omitempty"`
	Transformations Transformations      `json:"transformations,omitempty"`
}

type Transformations struct {
	SkipFields      []string `json:"skipFields,omitempty"`
	SensitiveFields []string `json:"sensitiveFields,omitempty"`
}

// LocalObjectReference is a reference to an object in the same namespace as the referent
type LocalObjectReference struct {
	// Name of the resource being referred to
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name"`
}
