package crd_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/josvazg/crd2go/internal/crd"
	"github.com/josvazg/crd2go/internal/crd/hooks"
	"github.com/josvazg/crd2go/internal/debug"
	"github.com/josvazg/crd2go/internal/gotype"
	"github.com/josvazg/crd2go/k8s"
	"github.com/josvazg/crd2go/pkg/config"
)

func TestRenameType(t *testing.T) {
	for _, tc := range []struct {
		name           string
		preloaded      []*gotype.GoType
		input          *gotype.GoField
		parents        []string
		want           string
		wantImportInfo bool
	}{
		{
			name: "Group Spec named Spec without preloads",
			input: gotype.NewGoField(
				"Spec",
				gotype.NewStruct("Spec", []*gotype.GoField{
					{
						Name:   "V20231115",
						GoType: &gotype.GoType{},
					},
				}),
			),
			parents: []string{"Group"},
			want:    "Spec",
		},

		{
			name: "Group Spec named GroupSpec with preloads",
			preloaded: []*gotype.GoType{
				{
					Name: "Spec", // reserves this type name
					Kind: "object",
				},
			},
			input: gotype.NewGoField(
				"Spec",
				gotype.NewStruct("Spec", []*gotype.GoField{
					{
						Name:   "V20231115",
						GoType: &gotype.GoType{},
					},
				},
				),
			),
			parents: []string{"Group"},
			want:    "GroupSpec",
		},

		{
			name:      "Identify Cross Reference",
			preloaded: []*gotype.GoType{CrossReference()},
			input: gotype.NewGoField(
				"SomeRef",
				gotype.NewStruct("SomeRef", []*gotype.GoField{
					{
						Name: "Namespace",
						GoType: &gotype.GoType{
							Name: "string",
							Kind: "string",
						},
					},
					{
						Name: "Name",
						GoType: &gotype.GoType{
							Name: "string",
							Kind: "string",
						},
					},
				}),
			),
			parents:        []string{"Group", "Spec"},
			want:           "Reference",
			wantImportInfo: true,
		},

		{
			name:      "Identify Local Reference",
			preloaded: []*gotype.GoType{LocalReference()},
			input: gotype.NewGoField(
				"SomeRef",
				gotype.NewStruct("SomeRef", []*gotype.GoField{
					{
						Name: "Name",
						GoType: &gotype.GoType{
							Name: "string",
							Kind: "string",
						},
					},
				}),
			),
			parents:        []string{"Group", "Spec"},
			want:           "LocalReference",
			wantImportInfo: true,
		},

		{
			name:      "Identify Local Reference behind an Array",
			preloaded: []*gotype.GoType{LocalReference()},
			input: gotype.NewGoField(
				"SomeRef",
				gotype.NewArray(
					gotype.NewStruct("SomeRef", []*gotype.GoField{
						{
							Name: "Name",
							GoType: &gotype.GoType{
								Name: "string",
								Kind: "string",
							},
						},
					}),
				),
			),
			parents:        []string{"Group", "Spec"},
			want:           "LocalReference",
			wantImportInfo: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			td := gotype.NewTypeDict(nil, tc.preloaded...)
			err := td.RenameField(tc.input, tc.parents)
			require.NoError(t, err)
			goType := tc.input.GoType
			if goType.Kind == gotype.ArrayKind {
				goType = goType.Element
			}
			assert.Equal(t, tc.want, goType.Name)
			if tc.wantImportInfo {
				assert.NotNil(t, goType.Import)
			}
		})
	}
}

func TestBuildOpenAPIType(t *testing.T) {
	td := gotype.NewTypeDict(nil, CrossReference(), LocalReference())
	crdRootType := &crd.CRDType{
		Name:    "RootType",
		Parents: []string{},
		Schema: &apiextensionsv1.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextensionsv1.JSONSchemaProps{
				"arrayOfStrings": {
					Type: "array",
					Items: &apiextensionsv1.JSONSchemaPropsOrArray{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: "string",
						},
					},
				},
				"arrayOfObjects": {
					Type: "array",
					Items: &apiextensionsv1.JSONSchemaPropsOrArray{
						Schema: &apiextensionsv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"key":   {Type: "string"},
								"value": {Type: "integer"},
							},
						},
					},
				},
				"randomObject": {
					Type: "object",
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"field1": {Type: "string"},
						"field2": {Type: "number"},
					},
				},
				"localReference": {
					Type: "object",
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"name": {Type: "string"},
					},
				},
				"crossReference": {
					Type: "object",
					Properties: map[string]apiextensionsv1.JSONSchemaProps{
						"name":      {Type: "string"},
						"namespace": {Type: "string"},
					},
				},
				"simpleString":  {Type: "string"},
				"simpleNumber":  {Type: "number"},
				"simpleInteger": {Type: "integer"},
			},
		},
	}
	goType, err := crd.FromOpenAPIType(td, hooks.Hooks, crdRootType)
	assert.NoError(t, err)
	assert.NotNil(t, goType)

	expectedType := gotype.NewStruct("RootType", []*gotype.GoField{
		gotype.NewGoField("ArrayOfStrings", gotype.NewArray(gotype.NewPrimitive("string", "string"))),
		gotype.NewGoField("ArrayOfObjects", gotype.NewArray(
			gotype.NewStruct("arrayOfObjects", []*gotype.GoField{
				gotype.NewGoField("Key", gotype.NewPrimitive("string", "string")),
				gotype.NewGoField("Value", gotype.NewPrimitive("int", "int")),
			}),
		)),
		gotype.NewGoField("RandomObject", gotype.NewStruct("RandomObject", []*gotype.GoField{
			gotype.NewGoField("Field1", gotype.NewPrimitive("string", "string")),
			gotype.NewGoField("Field2", gotype.NewPrimitive("float64", "float64")),
		})),
		gotype.NewGoField("LocalReference", gotype.AddImportInfo(gotype.NewStruct("LocalReference", []*gotype.GoField{
			gotype.NewGoField("Name", gotype.NewPrimitive("string", "string")),
		}), "k8s", "github.com/josvazg/crd2go/k8s")),
		gotype.NewGoField("CrossReference", gotype.AddImportInfo(gotype.NewStruct("Reference", []*gotype.GoField{
			gotype.NewGoField("Name", gotype.NewPrimitive("string", "string")),
			gotype.NewGoField("Namespace", gotype.NewPrimitive("string", "string")),
		}), "k8s", "github.com/josvazg/crd2go/k8s")),
		gotype.NewGoField("SimpleString", gotype.NewPrimitive("string", "string")),
		gotype.NewGoField("SimpleNumber", gotype.NewPrimitive("float64", "float64")),
		gotype.NewGoField("SimpleInteger", gotype.NewPrimitive("int", "int")),
	})

	fmt.Printf("Generated GoType: %s\n", debug.JSONize(goType))
	fmt.Printf("Expected GoType: %s\n", debug.JSONize(expectedType))

	assert.Equal(t, expectedType, goType)
}

func TestBuiltInFormat2Type(t *testing.T) {
	td := gotype.NewTypeDict(nil, gotype.KnownTypes()...)
	crdTimeType := &crd.CRDType{
		Name:    "time",
		Parents: []string{},
		Schema: &apiextensionsv1.JSONSchemaProps{
			Type:   "string",
			Format: "date-time",
		},
	}
	got, err := crd.FromOpenAPIType(td, hooks.Hooks, crdTimeType)
	require.NoError(t, err)
	want := &gotype.GoType{
		Name: "Time",
		Kind: "opaque",
		Import: &config.ImportInfo{
			Alias: "metav1",
			Path:  "k8s.io/apimachinery/pkg/apis/meta/v1",
		},
	}
	assert.Equal(t, want, got)
}

func TestConditionsMatch(t *testing.T) {
	for _, tc := range []struct {
		title string
		td    *gotype.TypeDict
	}{
		{
			title: "match conditions with a known type",
			td:    gotype.NewTypeDict(nil, gotype.KnownTypes()...),
		},
		{
			title: "match conditions with renames and imports",
			td: gotype.NewTypeDict(
				map[string]string{
					"Cond": "Condition",
				},
				gotype.NewAutoImportType(&config.ImportedTypeConfig{
					Name: "Condition",
					ImportInfo: config.ImportInfo{
						Alias: "metav1",
						Path:  "k8s.io/apimachinery/pkg/apis/meta/v1",
					},
				}),
			),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			input := &gotype.GoType{
				Name: "Cond",
				Kind: "struct",
				Fields: []*gotype.GoField{
					{
						Comment: "Last time the condition transitioned from one status to another.",
						Name:    "LastTransitionTime",
						GoType: &gotype.GoType{
							Name: "Time",
							Kind: "opaque",
							Import: &config.ImportInfo{
								Alias: "metav1",
								Path:  "k8s.io/apimachinery/pkg/apis/meta/v1",
							},
						},
					},
					{
						Comment: "A human readable message indicating details about the transition.",
						Name:    "Message",
						GoType:  &gotype.GoType{Name: "string", Kind: gotype.StringKind},
					},
					{
						Comment: "observedGeneration represents the .metadata.generation that the condition was set based upon.",
						Name:    "ObservedGeneration",
						GoType:  &gotype.GoType{Name: "int64", Kind: gotype.IntKind},
					},
					{
						Comment: "The reason for the condition's last transition.",
						Name:    "Reason",
						GoType:  &gotype.GoType{Name: "string", Kind: gotype.StringKind},
					},
					{
						Comment: "Status of the condition, one of True, False, Unknown.",
						Name:    "Status",
						GoType:  &gotype.GoType{Name: "ConditionStatus", Kind: gotype.StringKind},
					},
					{
						Comment: "Type of condition.",
						Name:    "Type",
						GoType:  &gotype.GoType{Name: "string", Kind: gotype.StringKind},
					},
				},
				Import: &config.ImportInfo{},
			}
			require.NoError(t, tc.td.RenameType([]string{"conditions"}, input))
			want := &gotype.GoType{
				Name: "Condition",
				Kind: "struct",
				Fields: []*gotype.GoField{
					{
						Comment: "Last time the condition transitioned from one status to another.",
						Name:    "LastTransitionTime",
						GoType: &gotype.GoType{
							Name: "Time",
							Kind: "opaque",
							Import: &config.ImportInfo{
								Alias: "metav1",
								Path:  "k8s.io/apimachinery/pkg/apis/meta/v1",
							},
						},
					},
					{
						Comment: "A human readable message indicating details about the transition.",
						Name:    "Message",
						GoType:  &gotype.GoType{Name: "string", Kind: gotype.StringKind},
					},
					{
						Comment: "observedGeneration represents the .metadata.generation that the condition was set based upon.",
						Name:    "ObservedGeneration",
						GoType:  &gotype.GoType{Name: "int64", Kind: gotype.IntKind},
					},
					{
						Comment: "The reason for the condition's last transition.",
						Name:    "Reason",
						GoType:  &gotype.GoType{Name: "string", Kind: gotype.StringKind},
					},
					{
						Comment: "Status of the condition, one of True, False, Unknown.",
						Name:    "Status",
						GoType:  &gotype.GoType{Name: "ConditionStatus", Kind: gotype.StringKind},
					},
					{
						Comment: "Type of condition.",
						Name:    "Type",
						GoType:  &gotype.GoType{Name: "string", Kind: gotype.StringKind},
					},
				},
				Import: &config.ImportInfo{
					Alias: "metav1",
					Path:  "k8s.io/apimachinery/pkg/apis/meta/v1",
				},
			}
			assert.Equal(t, want, input)
		})
	}
}

func CrossReference() *gotype.GoType {
	return gotype.MustTypeFrom(reflect.TypeOf(k8s.Reference{}))
}

func LocalReference() *gotype.GoType {
	return gotype.MustTypeFrom(reflect.TypeOf(k8s.LocalReference{}))
}
