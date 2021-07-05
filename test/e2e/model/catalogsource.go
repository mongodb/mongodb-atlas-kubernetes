package model

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CatalogSource struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            CatalogSourceSpec  `json:"spec,omitempty"`
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
