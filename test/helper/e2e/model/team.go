package model

import (
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewTeam(name, namespace string) *v1.AtlasTeam {
	return &v1.AtlasTeam{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasTeam",
			APIVersion: "atlas.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.TeamSpec{
			Name:      name,
			Usernames: []v1.TeamUser{},
		},
	}
}
