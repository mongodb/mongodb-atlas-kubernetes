package gotype

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeDictHas(t *testing.T) {
	tests := map[string]struct {
		goType      *GoType
		bySignature map[string]*GoType
		expected    bool
	}{
		"empty dict": {
			goType:      NewPrimitive("string", StringKind),
			bySignature: map[string]*GoType{},
			expected:    false,
		},
		"existing item": {
			goType: NewPrimitive("string", StringKind),
			bySignature: map[string]*GoType{
				"string": NewPrimitive("string", StringKind),
			},
			expected: true,
		},
		"non-existing item": {
			goType: NewPrimitive("int", IntKind),
			bySignature: map[string]*GoType{
				"string": NewPrimitive("string", StringKind),
			},
			expected: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := TypeDict{
				bySignature: tt.bySignature,
			}
			assert.Equal(t, tt.expected, td.Has(tt.goType))
		})
	}
}

func TestTypeDictGet(t *testing.T) {
	tests := map[string]struct {
		byName       map[string]*GoType
		name         string
		expectedType *GoType
		expected     bool
	}{
		"empty dict": {
			byName:       map[string]*GoType{},
			name:         "MyType",
			expectedType: nil,
			expected:     false,
		},
		"existing item": {
			byName: map[string]*GoType{
				"MyType": NewPrimitive("MyType", StringKind),
			},
			name:         "MyType",
			expectedType: NewPrimitive("MyType", StringKind),
			expected:     true,
		},
		"non-existing item": {
			byName: map[string]*GoType{
				"MyType": NewPrimitive("MyType", StringKind),
			},
			name:         "OtherType",
			expectedType: nil,
			expected:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := TypeDict{
				byName: tt.byName,
			}
			goType, ok := td.Get(tt.name)
			assert.Equal(t, tt.expected, ok)
			assert.Equal(t, tt.expectedType, goType)
		})
	}
}

func TestTypeDictAddAll(t *testing.T) {
	typeString := NewPrimitive("String", StringKind)
	typeArrayOfString := NewArray(NewPrimitive("String", StringKind))

	tests := map[string]struct {
		goTypes      []*GoType
		expectedDict *TypeDict
	}{
		"add multiple types": {
			goTypes: []*GoType{
				typeString,
				typeArrayOfString,
			},
			expectedDict: &TypeDict{
				bySignature: map[string]*GoType{
					"string":   typeString,
					"[string]": typeArrayOfString,
				},
				byName: map[string]*GoType{
					"String": typeString,
					"":       typeArrayOfString,
				},
				generated: make(map[string]bool),
				renames:   make(map[string]string),
			},
		},
		"add duplicate types": {
			goTypes: []*GoType{
				typeString,
				typeString,
			},
			expectedDict: &TypeDict{
				bySignature: map[string]*GoType{
					"string": typeString,
				},
				byName: map[string]*GoType{
					"String": typeString,
				},
				generated: make(map[string]bool),
				renames:   make(map[string]string),
			},
		},
		"add no types": {
			goTypes: []*GoType{},
			expectedDict: &TypeDict{
				bySignature: make(map[string]*GoType),
				byName:      make(map[string]*GoType),
				generated:   make(map[string]bool),
				renames:     make(map[string]string),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := NewTypeDict(map[string]string{}, []*GoType{}...)
			td.AddAll(tt.goTypes...)
			assert.Equal(t, tt.expectedDict, td)
		})
	}
}

func TestTypeDict_MarkGenerated(t *testing.T) {
	tests := map[string]struct {
		gt *GoType
	}{
		"mark generated type": {
			gt: NewPrimitive("String", StringKind),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := NewTypeDict(map[string]string{}, []*GoType{}...)
			td.MarkGenerated(tt.gt)
			assert.True(t, td.WasGenerated(tt.gt))
		})
	}
}

func TestTypeDict_RenameField(t *testing.T) {
	tests := map[string]struct {
		goField       *GoField
		parentNames   []string
		expectedField *GoField
		expectedErr   error
	}{
		"rename field with nil GoType": {
			goField:       NewGoField("Name", nil),
			parentNames:   []string{},
			expectedField: NewGoField("Name", nil),
			expectedErr:   fmt.Errorf("failed to rename type for field Name: GoType is nil"),
		},
		"rename field with primitive GoType": {
			goField:       NewGoField("Name", NewPrimitive("string", StringKind)),
			parentNames:   []string{},
			expectedField: NewGoField("Name", NewPrimitive("string", StringKind)),
			expectedErr:   nil,
		},
		"rename field with non-primitive GoType": {
			goField: NewGoField(
				"Project",
				NewStruct(
					"ProjectObject",
					[]*GoField{},
				),
			),
			parentNames: []string{"Parent"},
			expectedField: NewGoField(
				"Project",
				NewStruct(
					"Project",
					[]*GoField{},
				),
			),
			expectedErr: nil,
		},
		"rename field with non-primitive GoType and already mapped type": {
			goField: NewGoField(
				"TeamList",
				NewStruct("TeamList", []*GoField{}),
			),
			parentNames: []string{"Company"},
			expectedField: NewGoField(
				"TeamList",
				NewStruct("CompanyTeams", []*GoField{}),
			),
			expectedErr: nil,
		},
		"reuse existing type": {
			goField: NewGoField(
				"Org",
				NewStruct("Org", []*GoField{}),
			),
			parentNames: []string{},
			expectedField: NewGoField(
				"Org",
				NewStruct("Organization", []*GoField{}),
			),
			expectedErr: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			td := NewTypeDict(
				map[string]string{
					"ProjectObject": "Project",
					"TeamList":      "Teams",
					"Org":           "Organization",
				},
				[]*GoType{
					NewPrimitive("Teams", StringKind),
					NewStruct("Organization", []*GoField{}),
				}...,
			)
			assert.Equal(t, tt.expectedErr, td.RenameField(tt.goField, tt.parentNames))
			assert.Equal(t, tt.expectedField, tt.goField)
		})
	}
}
