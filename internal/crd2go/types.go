package crd2go

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	StructKind = "struct"
	ArrayKind  = "array"
	StringKind = "string"
	IntKind    = "int"
	FloatKind  = "float64"
	BoolKind   = "bool"
)

const (
	OpenAPIObject  = "object"
	OpenAPIArray   = "array"
	OpenAPIString  = "string"
	OpenAPIInteger = "integer"
	OpenAPINumber  = "number"
	OpenAPIBoolean = "boolean"
)

type GoType struct {
	Name    string
	Kind    string
	Fields  []*GoField
	Element *GoType
}

func (g *GoType) isPrimitive() bool {
	switch g.Kind {
	case StringKind, IntKind, FloatKind, BoolKind:
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
	if gt.Kind == StructKind {
		fieldSignatures := make([]string, len(gt.Fields))
		for _, field := range gt.Fields {
			fieldSignatures = append(fieldSignatures, field.Signature())
		}
		return fmt.Sprintf("{%s}", strings.Join(fieldSignatures, ","))
	}
	if gt.Kind == ArrayKind {
		return fmt.Sprintf("[%s]", gt.Element.Signature())
	}
	return fmt.Sprintf("%s(%s)", gt.Name, gt.Kind)
}

func (g *GoType) BaseType() *GoType {
    if g == nil {
        return nil
    }
    if g.Kind == ArrayKind {
        return g.Element.BaseType()
    }
    return g
}

type GoField struct {
	Comment  string
	Required bool
	Name     string
	GoType   *GoType
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
	goType := f.GoType.BaseType()
	if goType.isPrimitive() {
		return nil // primitive types are not to be renamed
	}
	if td.Has(goType) {
		existingType := td.bySignature[goType.Signature()]
		if existingType == nil {
			return fmt.Errorf("failed to find existing type for %v", f)
		}
		goType.Name = existingType.Name
		return nil
	}

	typeName := goType.Name
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
	goType.Name = typeName
	td.Add(goType)

	return nil
}

func NewPrimitive(name, kind string) *GoType {
	return &GoType{
		Name: name,
		Kind: kind,
	}
}

func NewArray(element *GoType) *GoType {
	return &GoType{
		Name:    "",
		Kind:    ArrayKind,
		Element: element,
	}
}

func NewStruct(name string, fields []*GoField) *GoType {
	return &GoType{
		Name:   title(name),
		Kind:   StructKind,
		Fields: orderFieldsByName(fields),
	}
}

type TypeDict struct {
	bySignature map[string]*GoType
	byName      map[string]*GoType
	generated   map[string]bool
}

func NewTypeDict(goTypes ...*GoType) TypeDict {
	td := TypeDict{
		bySignature: make(map[string]*GoType),
		byName:      make(map[string]*GoType),
		generated:   make(map[string]bool),
	}
	for _, gt := range goTypes {
		td.Add(gt)
	}
	return td
}

func (td TypeDict) Has(gt *GoType) bool {
	_, ok := td.bySignature[gt.Signature()]
	return ok
}

func (td TypeDict) Get(name string) (*GoType, bool) {
	gt, ok := td.byName[name]
	return gt, ok
}

func (td TypeDict) Add(gt *GoType) {
	titledName := title(gt.Name)
	if gt.Name != titledName {
		panic(fmt.Sprintf("type name %s is not titled", gt.Name))
	}
	td.bySignature[gt.Signature()] = gt
	td.byName[gt.Name] = gt
}

func (td TypeDict) MarkGenerated(gt *GoType) {
	if !td.Has(gt) {
		td.Add(gt)
	}
	td.generated[gt.Name] = true
}

func (td TypeDict) WasGenerated(gt *GoType) bool {
	if td.Has(gt) {
		return td.generated[gt.Name]
	}
	return false
}

func orderFieldsByName(fields []*GoField) []*GoField {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

func FromOpenAPIType(td TypeDict, typeName string, parents []string, schema *apiextensions.JSONSchemaProps) (*GoType, error) {
	switch schema.Type {
	case OpenAPIObject:
		return fromOpenAPIStruct(td, typeName, parents, schema)
	case OpenAPIArray:
		return fromOpenAPIArray(td, typeName, schema)
	case OpenAPIString, OpenAPIInteger, OpenAPINumber, OpenAPIBoolean:
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
		field.Required = slices.Contains(schema.Required, key)
		if err := field.RenameType(td, fieldsParents); err != nil {
			return nil, fmt.Errorf("failed to rename field %v: %w", field, err)
		}
		fields = append(fields, field)
	}
	return NewStruct(typeName, fields), nil
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
	return NewArray(elementType), nil
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
	case OpenAPIString:
		return StringKind, nil
	case OpenAPIInteger:
		return IntKind, nil
	case OpenAPINumber:
		return FloatKind, nil
	case OpenAPIBoolean:
		return BoolKind, nil
	default:
		return "", fmt.Errorf("unsupported Open API kind %s", kind)
	}
}

// orderedkeys returns a sorted slice of keys from the given map
func orderedkeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
