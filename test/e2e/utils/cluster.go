package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

type AC struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ClusterSpec        `json:"spec,omitempty"`
}

type ClusterSpec v1.AtlasClusterSpec

func (ac *AC) ClusterFileName() string {
	return "data/cluster-" + ac.ObjectMeta.Name + ".yaml"
}

func (ac *AC) GetClusterNameResource() string {
	return "atlascluster.atlas.mongodb.com/" + ac.ObjectMeta.Name
}
