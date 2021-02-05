package utils

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ISpec interface {
	ProjectName(string) ISpec
	SecretRef(string) ISpec
	// TODO WhiteIP
	CompleteK8sConfig(string) []byte
}

type ProjectSpec v1.AtlasProjectSpec

type ap struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec            ProjectSpec        `json:"spec,omitempty"`
}

func NewProject() ISpec {
	return &ProjectSpec{}
}

func (s *ProjectSpec) ProjectName(name string) ISpec {
	s.Name = name
	return s
}

func (s *ProjectSpec) SecretRef(name string) ISpec {
	s.ConnectionSecret = &v1.ResourceRef{Name: name}
	return s
}

func (s ProjectSpec) CompleteK8sConfig(k8sname string) []byte {
	var t ap
	t.TypeMeta = metav1.TypeMeta{
		APIVersion: "atlas.mongodb.com/v1",
		Kind:       "AtlasProject",
	}
	t.ObjectMeta = &metav1.ObjectMeta{
		Name: k8sname,
	}
	t.Spec = s
	yamlConf, _ := JSONToYAMLConvert(t)
	return yamlConf
}
