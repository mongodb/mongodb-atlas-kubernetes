package v1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestAtlasCustomRoles_ToAtlas(t *testing.T) {
	tests := []struct {
		name string
		spec CustomRole
		want *mongodbatlas.CustomDBRole
	}{
		{
			name: "Should convert all fields",
			spec: CustomRole{
				Name: "testName",
				InheritedRoles: []Role{
					{
						Name:     "testName",
						Database: "testDB",
					},
				},
				Actions: []Action{
					{
						Name: "testName",
						Resources: []Resource{
							{
								Cluster:    pointer.MakePtr(false),
								Database:   pointer.MakePtr("testDB"),
								Collection: pointer.MakePtr("testCollection"),
							},
							{
								Cluster:    pointer.MakePtr(false),
								Database:   pointer.MakePtr("testDB2"),
								Collection: pointer.MakePtr("testCollection2"),
							},
						},
					},
					{
						Name: "testName2",
						Resources: []Resource{
							{
								Cluster: pointer.MakePtr(true),
							},
						},
					},
					{
						Name: "testName3",
						Resources: []Resource{
							{
								Database:   pointer.MakePtr(""),
								Collection: pointer.MakePtr(""),
							},
						},
					},
				},
			},
			want: &mongodbatlas.CustomDBRole{
				RoleName: "testName",
				InheritedRoles: []mongodbatlas.InheritedRole{
					{
						Role: "testName",
						Db:   "testDB",
					},
				},
				Actions: []mongodbatlas.Action{
					{
						Action: "testName",
						Resources: []mongodbatlas.Resource{
							{
								Cluster:    nil,
								DB:         pointer.MakePtr("testDB"),
								Collection: pointer.MakePtr("testCollection"),
							},
							{
								Cluster:    nil,
								DB:         pointer.MakePtr("testDB2"),
								Collection: pointer.MakePtr("testCollection2"),
							},
						},
					},
					{
						Action: "testName2",
						Resources: []mongodbatlas.Resource{
							{
								Cluster: pointer.MakePtr(true),
							},
						},
					},
					{
						Action: "testName3",
						Resources: []mongodbatlas.Resource{
							{
								Cluster:    nil,
								DB:         pointer.MakePtr(""),
								Collection: pointer.MakePtr(""),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.ToAtlas()
			if !reflect.DeepEqual(got, tt.want) {
				g, _ := json.MarshalIndent(got, "", " ")
				w, _ := json.MarshalIndent(tt.want, "", " ")
				fmt.Println("GOT", string(g))
				fmt.Println("WANT", string(w))
			}

			assert.Equalf(t, tt.want, got, "ToAtlas()")
		})
	}
}
