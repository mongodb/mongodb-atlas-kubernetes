package utils

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AC struct {
	metav1.TypeMeta   `json:",inline"`
	ObjectMeta *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec ClusterSpec `json:"spec,omitempty"`
}

type ClusterSpec v1.AtlasClusterSpec
