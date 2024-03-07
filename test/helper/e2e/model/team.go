package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func NewTeam(name, namespace string) *akov2.AtlasTeam {
	return &akov2.AtlasTeam{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AtlasTeam",
			APIVersion: "atlas.mongodb.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: akov2.TeamSpec{
			Name:      name,
			Usernames: []akov2.TeamUser{},
		},
	}
}
