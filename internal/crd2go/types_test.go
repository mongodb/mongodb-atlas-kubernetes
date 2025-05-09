package crd2go_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestRenameType(t *testing.T) {
	for _, tc := range []struct {
		name      string
		preloaded []*crd2go.GoType
		input     *crd2go.GoField
		parents   []string
		want      string
	}{
		{
			name: "Group Spec named Spec without preloads",
			input: crd2go.NewGoField(
				"Spec",
				crd2go.NewStruct("Spec", []*crd2go.GoField{
					{
						Name:   "V20231115",
						GoType: &crd2go.GoType{},
					},
				}),
			),
			parents: []string{"Group"},
			want:    "Spec",
		},

		{
			name: "Group Spec named GroupSpec with preloads",
			preloaded: []*crd2go.GoType{
				{
					Name: "Spec", // reserves this type name
					Kind: "object",
				},
			},
			input: crd2go.NewGoField(
				"Spec",
				crd2go.NewStruct("Spec", []*crd2go.GoField{
					{
						Name:   "V20231115",
						GoType: &crd2go.GoType{},
					},
				},
				),
			),
			parents: []string{"Group"},
			want:    "GroupSpec",
		},

		{
			name:      "Identify Cross Reference",
			preloaded: []*crd2go.GoType{CrossReference()},
			input: crd2go.NewGoField(
				"SomeRef",
				crd2go.NewStruct("SomeRef", []*crd2go.GoField{
					{
						Name: "Namespace",
						GoType: &crd2go.GoType{
							Name: "string",
							Kind: "string",
						},
					},
					{
						Name: "Name",
						GoType: &crd2go.GoType{
							Name: "string",
							Kind: "string",
						},
					},
				}),
			),
			parents: []string{"Group", "Spec"},
			want:    "K8sCrossReference",
		},

		{
			name:      "Identify Local Reference",
			preloaded: []*crd2go.GoType{LocalReference()},
			input: crd2go.NewGoField(
				"SomeRef",
				crd2go.NewStruct("SomeRef", []*crd2go.GoField{
					{
						Name: "Name",
						GoType: &crd2go.GoType{
							Name: "string",
							Kind: "string",
						},
					},
				}),
			),
			parents: []string{"Group", "Spec"},
			want:    "K8sLocalReference",
		},

		{
			name:      "Identify Local Reference behind an Array",
			preloaded: []*crd2go.GoType{LocalReference()},
			input: crd2go.NewGoField(
				"SomeRef",
				crd2go.NewArray(
					crd2go.NewStruct("SomeRef", []*crd2go.GoField{
						{
							Name: "Name",
							GoType: &crd2go.GoType{
								Name: "string",
								Kind: "string",
							},
						},
					}),
				),
			),
			parents: []string{"Group", "Spec"},
			want:    "K8sLocalReference",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			td := crd2go.NewTypeDict(tc.preloaded...)
			err := tc.input.RenameType(td, tc.parents)
			require.NoError(t, err)
			goType := tc.input.GoType
			if goType.Kind == crd2go.ArrayKind {
				goType = goType.Element
			}
			assert.Equal(t, tc.want, goType.Name)
		})
	}
}

func TestBuildOpenAPIType(t *testing.T) {
	td := crd2go.NewTypeDict(CrossReference(), LocalReference())

	schema := &apiextensions.JSONSchemaProps{
		Type: "object",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"arrayOfStrings": {
				Type: "array",
				Items: &apiextensions.JSONSchemaPropsOrArray{
					Schema: &apiextensions.JSONSchemaProps{
						Type: "string",
					},
				},
			},
			"arrayOfObjects": {
				Type: "array",
				Items: &apiextensions.JSONSchemaPropsOrArray{
					Schema: &apiextensions.JSONSchemaProps{
						Type: "object",
						Properties: map[string]apiextensions.JSONSchemaProps{
							"key":   {Type: "string"},
							"value": {Type: "integer"},
						},
					},
				},
			},
			"randomObject": {
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"field1": {Type: "string"},
					"field2": {Type: "number"},
				},
			},
			"localReference": {
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name": {Type: "string"},
				},
			},
			"crossReference": {
				Type: "object",
				Properties: map[string]apiextensions.JSONSchemaProps{
					"name":      {Type: "string"},
					"namespace": {Type: "string"},
				},
			},
			"simpleString":  {Type: "string"},
			"simpleNumber":  {Type: "number"},
			"simpleInteger": {Type: "integer"},
		},
	}

	goType, err := crd2go.FromOpenAPIType(td, "RootType", []string{}, schema)
	assert.NoError(t, err)
	assert.NotNil(t, goType)

	expectedType := crd2go.NewStruct("RootType", []*crd2go.GoField{
		crd2go.NewGoField("ArrayOfStrings", crd2go.NewArray(crd2go.NewPrimitive("string", "string"))),
		crd2go.NewGoField("ArrayOfObjects", crd2go.NewArray(
			crd2go.NewStruct("arrayOfObjects", []*crd2go.GoField{
				crd2go.NewGoField("Key", crd2go.NewPrimitive("string", "string")),
				crd2go.NewGoField("Value", crd2go.NewPrimitive("int", "int")),
			}),
		)),
		crd2go.NewGoField("RandomObject", crd2go.NewStruct("RandomObject", []*crd2go.GoField{
			crd2go.NewGoField("Field1", crd2go.NewPrimitive("string", "string")),
			crd2go.NewGoField("Field2", crd2go.NewPrimitive("float64", "float64")),
		})),
		crd2go.NewGoField("LocalReference", crd2go.NewStruct("K8sLocalReference", []*crd2go.GoField{
			crd2go.NewGoField("Name", crd2go.NewPrimitive("string", "string")),
		})),
		crd2go.NewGoField("CrossReference", crd2go.NewStruct("K8sCrossReference", []*crd2go.GoField{
			crd2go.NewGoField("Name", crd2go.NewPrimitive("string", "string")),
			crd2go.NewGoField("Namespace", crd2go.NewPrimitive("string", "string")),
		})),
		crd2go.NewGoField("SimpleString", crd2go.NewPrimitive("string", "string")),
		crd2go.NewGoField("SimpleNumber", crd2go.NewPrimitive("float64", "float64")),
		crd2go.NewGoField("SimpleInteger", crd2go.NewPrimitive("int", "int")),
	})

	fmt.Printf("Generated GoType: %s\n", jsonize(goType))
	fmt.Printf("Expected GoType: %s\n", jsonize(expectedType))

	assert.Equal(t, expectedType, goType)
}

func jsonize(obj any) string {
	js, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(js)
}

func CrossReference() *crd2go.GoType {
	return crd2go.NewStruct("K8sCrossReference", []*crd2go.GoField{
		crd2go.NewGoField("Name", crd2go.NewPrimitive("string", "string")),
		crd2go.NewGoField("Namespace", crd2go.NewPrimitive("string", "string")),
	})
}

func LocalReference() *crd2go.GoType {
	return crd2go.NewStruct("K8sLocalReference", []*crd2go.GoField{
		crd2go.NewGoField("Name", crd2go.NewPrimitive("string", "string")),
	})
}
