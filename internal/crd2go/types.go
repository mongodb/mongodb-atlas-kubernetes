package crd2go

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type GoType struct {
	Name    string
	Kind    string
	Fields  []*GoField
	Element *GoType
}

func (g *GoType) isPrimitive(kind string) bool {
	switch kind {
	case "string", "int", "float64", "bool":
		return true
	default:
		return false
	}
}

func (g *GoType) Equal(gt *GoType) bool {
	if g.Name != gt.Name || g.Kind != gt.Kind {
		return false
	}

	if len(g.Fields) != len(gt.Fields) {
		return false
	}

	if !slices.EqualFunc(g.Fields, gt.Fields, func(f1, f2 *GoField) bool {
		return f1.Name == f2.Name && f1.GoType.Equal(f2.GoType)
	}) {
		return false
	}

	if (g.Element == nil) != (gt.Element == nil) {
		return false
	}

	if g.Element != nil && !g.Element.Equal(gt.Element) {
		return false
	}

	return true
}

func (gt *GoType) Signature() string {
	if gt == nil {
		return "nil"
	}
	if gt.Kind == "object" {
		fieldSignatures := make([]string, len(gt.Fields))
		for _, field := range gt.Fields {
			fieldSignatures = append(fieldSignatures, field.Signature())
		}
		return fmt.Sprintf("{%s}", strings.Join(fieldSignatures, ","))
	}
	if gt.Kind == "array" {
		return fmt.Sprintf("[%s]", gt.Element.Signature())
	}
	return fmt.Sprintf("%s(%s)", gt.Name, gt.Kind)
}

type GoField struct {
	Comment string
	Name    string
	GoType  *GoType
}

func NewGoField(name string, gt *GoType) *GoField {
	return &GoField{
		Name:   name,
		GoType: gt,
	}
}

func (g *GoField) Signature() string {
	if g == nil {
		return "nil"
	}
	return fmt.Sprintf("%s:%s", g.Name, g.GoType.Signature())
}

func (f *GoField) RenameType(td TypeDict, parentNames []string) error {
	if f.GoType == nil {
		return fmt.Errorf("failed to rename type for field %s: GoType is nil", f.Name)
	}
	if f.GoType.isPrimitive(f.GoType.Kind) {
		return nil // primitive types are not to be renamed
	}
	if td.Has(f.GoType) {
		existingType := td.bySiganture[f.GoType.Signature()]
		if existingType == nil {
			return fmt.Errorf("failed to find existing type for %v", f)
		}
		f.GoType.Name = existingType.Name
		return nil
	}

	typeName := title(f.GoType.Name)
	for i := len(parentNames) - 1; i >= 0; i-- {
		_, used := td.Get(typeName)
		if !used {
			break
		}
		typeName = fmt.Sprintf("%s%s", title(parentNames[i]), typeName)
	}

	_, used := td.Get(typeName)
	if used {
		return fmt.Errorf("failed to find a free type name for type %v", f)
	}
	f.GoType.Name = typeName
	td.Add(f.GoType)

	return nil
}

func NewPrimitive(name, kind string) *GoType {
	return &GoType{
		Name: name,
		Kind: kind,
	}
}

func NewArray(name string, element *GoType) *GoType {
	return &GoType{
		Name:    name,
		Kind:    "array",
		Element: element,
	}
}

func NewObject(name string, fields []*GoField) *GoType {
	return &GoType{
		Name:   name,
		Kind:   "object",
		Fields: orderFieldsByName(fields),
	}
}

type TypeDict struct {
	bySiganture map[string]*GoType
	byName      map[string]*GoType
}

func NewTypeDict(goTypes ... *GoType) TypeDict {
	td := TypeDict{
		bySiganture: make(map[string]*GoType),
		byName:      make(map[string]*GoType),
	}
	for _, gt := range goTypes {
		td.Add(gt)
	}
	return td
}

func (d TypeDict) Has(gt *GoType) bool {
	_, ok := d.bySiganture[gt.Signature()]
	return ok
}

func (d TypeDict) Get(name string) (*GoType, bool) {
	gt, ok := d.byName[name]
	return gt, ok
}

func (d TypeDict) Add(gt *GoType) {
	d.bySiganture[gt.Signature()] = gt
	d.byName[gt.Name] = gt
}

func orderFieldsByName(fields []*GoField) []*GoField {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

func FromOpenAPIType(td TypeDict, typeName string, parents []string, schema *apiextensions.JSONSchemaProps) (*GoType, error) {
	switch schema.Type {
	case "object":
		return fromOpenAPIStruct(td, typeName, parents, schema)
	case "array":
		return fromOpenAPIArray(td, typeName, schema)
	case "string", "integer", "number", "boolean":
		return fromOpenAPIPrimitive(schema.Type)
	default:
		return nil, fmt.Errorf("unsupported Open API type %q", schema.Type)
	}
}

func fromOpenAPIStruct(td TypeDict, typeName string, parents []string, schema *apiextensions.JSONSchemaProps) (*GoType, error) {
	fields := []*GoField{}
	fieldsParents := append(parents, typeName)
	for _, key := range orderedkeys(schema.Properties) {
		props := schema.Properties[key]
		fieldType, err := FromOpenAPIType(td, key, fieldsParents, &props)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s type: %w", key, err)
		}
		field := NewGoField(title(key), fieldType)
		field.Comment = props.Description
		if err := field.RenameType(td, parents); err != nil {
			return nil, fmt.Errorf("failed to rename field %v: %w", field, err)
		}
		fields = append(fields, field)
	}
	return NewObject(typeName, fields), nil
}

func fromOpenAPIArray(td TypeDict, typeName string, schema *apiextensions.JSONSchemaProps) (*GoType, error) {
	if schema.Items == nil {
		return nil, fmt.Errorf("array %s has no items", typeName)
	}
	if schema.Items.Schema == nil {
		return nil, fmt.Errorf("array %s has no items schema", typeName)
	}
	elementType, err := FromOpenAPIType(td, typeName, nil, schema.Items.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse array %s element type: %w", typeName, err)
	}
	return NewArray(typeName, elementType), nil
}

func fromOpenAPIPrimitive(kind string) (*GoType, error) {
	goTypeName, err := openAPIKindtoGoType(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI kind %s: %w", kind, err)
	}
	return NewPrimitive(goTypeName, goTypeName), nil
}

func openAPIKindtoGoType(kind string) (string, error) {
	switch kind {
	case "string":
		return "string", nil
	case "integer":
		return "int", nil
	case "number":
		return "float64", nil
	case "boolean":
		return "bool", nil
	default:
		return "", fmt.Errorf("unsuppoerted Open API kind %s", kind)
	}
}
