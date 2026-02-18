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

package model

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CatalogSource struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            CatalogSourceSpec  `json:"spec"`
}

type CatalogSourceSpec struct {
	SourceType  string `json:"sourceType"`
	Image       string `json:"image"`
	DisplayName string `json:"displayName"`
	Publisher   string `json:"publisher"`
}

func NewCatalogSource(imageURL string) CatalogSource {
	name := strings.Split(imageURL, ":")[1]
	name = strings.ToLower(name)
	return CatalogSource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "operators.coreos.com/v1alpha1",
			Kind:       "CatalogSource",
		},
		ObjectMeta: &metav1.ObjectMeta{
			Name:      name,
			Namespace: "openshift-marketplace",
		},
		Spec: CatalogSourceSpec{
			SourceType:  "grpc",
			Image:       imageURL,
			DisplayName: name,
			Publisher:   name,
		},
	}
}
