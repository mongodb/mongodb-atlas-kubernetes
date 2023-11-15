package v1

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
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
								Cluster:    toptr.MakePtr(false),
								Database:   toptr.MakePtr("testDB"),
								Collection: toptr.MakePtr("testCollection"),
							},
							{
								Cluster:    toptr.MakePtr(false),
								Database:   toptr.MakePtr("testDB2"),
								Collection: toptr.MakePtr("testCollection2"),
							},
						},
					},
					{
						Name: "testName2",
						Resources: []Resource{
							{
								Cluster: toptr.MakePtr(true),
							},
						},
					},
					{
						Name: "testName3",
						Resources: []Resource{
							{
								Database:   toptr.MakePtr(""),
								Collection: toptr.MakePtr(""),
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
								DB:         toptr.MakePtr("testDB"),
								Collection: toptr.MakePtr("testCollection"),
							},
							{
								Cluster:    nil,
								DB:         toptr.MakePtr("testDB2"),
								Collection: toptr.MakePtr("testCollection2"),
							},
						},
					},
					{
						Action: "testName2",
						Resources: []mongodbatlas.Resource{
							{
								Cluster: toptr.MakePtr(true),
							},
						},
					},
					{
						Action: "testName3",
						Resources: []mongodbatlas.Resource{
							{
								Cluster:    nil,
								DB:         toptr.MakePtr(""),
								Collection: toptr.MakePtr(""),
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
