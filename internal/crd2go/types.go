package crd2go

import (
	"fmt"
	"path"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/josvazg/crd2go/k8s"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	UnsupportedKind = "unsupported"
	StructKind      = "struct"
	ArrayKind       = "array"
	StringKind      = "string"
	IntKind         = "int"
	FloatKind       = "float64"
	BoolKind        = "bool"
)

const (
	OpenAPIObject  = "object"
	OpenAPIArray   = "array"
	OpenAPIString  = "string"
	OpenAPIInteger = "integer"
	OpenAPINumber  = "number"
	OpenAPIBoolean = "boolean"
)

const PACKAGE_BASE = "github.com/josvazg/crd2go"

// GoType represents a Go type, which can be a primitive type, a struct, or an array.
// It is used in conjunbction with TypeDict to track and ensure unique type names.
type GoType struct {
	Name    string
	Kind    string
	Fields  []*GoField
	Element *GoType
	Import  *ImportInfo
}

// ImportInfo holds the import path and alias for existing types
type ImportInfo struct {
	Alias string
	Path  string
}

// isPrimitive checks if the GoType is a primitive type
func (g *GoType) isPrimitive() bool {
	switch g.Kind {
	case StringKind, IntKind, FloatKind, BoolKind:
		return true
	default:
		return false
	}
}

// signature generates a unique signature for the GoType.
// This is leveraged by TypeDict to quickly check if a type is already registered.
func (gt *GoType) signature() string {
	if gt == nil {
		return "nil"
	}
	if gt.Kind == StructKind {
		fieldSignatures := make([]string, len(gt.Fields))
		for _, field := range gt.Fields {
			fieldSignatures = append(fieldSignatures, field.signature())
		}
		return fmt.Sprintf("{%s}", strings.Join(fieldSignatures, ","))
	}
	if gt.Kind == ArrayKind {
		return fmt.Sprintf("[%s]", gt.Element.signature())
	}
	return fmt.Sprintf("%s(%s)", gt.Name, gt.Kind)
}

// baseType returns the base type of the GoType.
// If a type is an array, it returns the element type,
// traversing until a non-array type is found.
func (g *GoType) baseType() *GoType {
	if g == nil {
		return nil
	}
	if g.Kind == ArrayKind {
		return g.Element.baseType()
	}
	return g
}

// GoField is a field in a Go struct
type GoField struct {
	Comment  string
	Required bool
	Name     string
	GoType   *GoType
}

// NewGoField creates a new GoField with the given name and GoType
func NewGoField(name string, gt *GoType) *GoField {
	return &GoField{
		Name:   name,
		GoType: gt,
	}
}

// signature generates a unique signature for the GoField leveraging the type
// signature
func (g *GoField) signature() string {
	if g == nil {
		return "nil"
	}
	return fmt.Sprintf("%s:%s", g.Name, g.GoType.signature())
}

// RenameType renames the GoType of the field to ensure it is unique within the
// TypeDict. It uses the parent names as needed to create a unique name for the
// type, if the type is not a primitive and its name is already taken.
func (f *GoField) RenameType(td TypeDict, parentNames []string) error {
	if f.GoType == nil {
		return fmt.Errorf("failed to rename type for field %s: GoType is nil", f.Name)
	}
	goType := f.GoType.baseType()
	if goType.isPrimitive() {
		return nil // primitive types are not to be renamed
	}
	if td.Has(goType) {
		existingType := td.bySignature[goType.signature()]
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

// NewPrimitive creates a new GoType representing a primitive type
func NewPrimitive(name, kind string) *GoType {
	return &GoType{
		Name: name,
		Kind: kind,
	}
}

// NewArray creates a new GoType representing an array type
func NewArray(element *GoType) *GoType {
	return &GoType{
		Name:    "",
		Kind:    ArrayKind,
		Element: element,
	}
}

// NewStruct creates a new GoType representing a struct type
func NewStruct(name string, fields []*GoField) *GoType {
	return &GoType{
		Name:   title(name),
		Kind:   StructKind,
		Fields: orderFieldsByName(fields),
	}
}

func AddImportInfo(gt *GoType, packagePath, alias string) *GoType {
	effectiveAlias := alias
	if effectiveAlias == "" {
		effectiveAlias = path.Base(packagePath)
	}
	gt.Import = &ImportInfo{Path: packagePath, Alias: effectiveAlias}
	return gt
}

// TypeDict is a dictionary of Go types, used to track and ensure unique type names.
// It also keeps track of generated types to avoid re-genrating the same type again.
type TypeDict struct {
	bySignature map[string]*GoType
	byName      map[string]*GoType
	generated   map[string]bool
}

// NewTypeDict creates a new TypeDict with the given Go types
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

// Has checks if the TypeDict contains a GoType with the same signature
func (td TypeDict) Has(gt *GoType) bool {
	_, ok := td.bySignature[gt.signature()]
	return ok
}

// Get retrieves a GoType by its name from the TypeDict
func (td TypeDict) Get(name string) (*GoType, bool) {
	gt, ok := td.byName[name]
	return gt, ok
}

// Add adds a GoType to the TypeDict, ensuring that the type name is unique
func (td TypeDict) Add(gt *GoType) {
	titledName := title(gt.Name)
	if gt.Name != titledName {
		panic(fmt.Sprintf("type name %s is not titled", gt.Name))
	}
	td.bySignature[gt.signature()] = gt
	td.byName[gt.Name] = gt
}

// MarkGenerated marks a GoType as generated
func (td TypeDict) MarkGenerated(gt *GoType) {
	if !td.Has(gt) {
		td.Add(gt)
	}
	td.generated[gt.Name] = true
}

// WasGenerated checks if a GoType was marked as generated
func (td TypeDict) WasGenerated(gt *GoType) bool {
	if td.Has(gt) {
		return td.generated[gt.Name]
	}
	return false
}

// orderFieldsByName sorts the fields of a GoType by name
func orderFieldsByName(fields []*GoField) []*GoField {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

// FromOpenAPIType converts an OpenAPI schema to a GoType
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

// fromOpenAPIStruct converts an OpenAPI object schema to a GoType struct
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

// fromOpenAPIArray converts an OpenAPI array schema to a GoType array
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

// fromOpenAPIPrimitive converts an OpenAPI primitive type to a GoType
func fromOpenAPIPrimitive(kind string) (*GoType, error) {
	goTypeName, err := openAPIKindtoGoType(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI kind %s: %w", kind, err)
	}
	return NewPrimitive(goTypeName, goTypeName), nil
}

// openAPIKindtoGoType converts an OpenAPI kind to a Go type
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

func KnownTypes() []*GoType {
	return []*GoType{
		MustTypeFrom(reflect.TypeOf(k8s.LocalReference{})),
		MustTypeFrom(reflect.TypeOf(k8s.Reference{})),
	}
}

func MustTypeFrom(t reflect.Type) *GoType {
	gt, err := TypeFrom(t)
	if err != nil {
		panic(fmt.Errorf("failed to translate type %v: %w", t.Name(), err))
	}
	return gt
}

func TypeFrom(t reflect.Type) (*GoType, error) {
	kind := GoKind(t.Kind())
	switch kind {
	case StructKind:
		return StructTypeFrom(t)
	case ArrayKind:
		return ArrayTypeFrom(t)
	case StringKind, IntKind, FloatKind, BoolKind:
		return NewPrimitive(t.Name(), kind), nil
	default:
		return nil, fmt.Errorf("unsupported kind %v", kind)
	}
}

func StructTypeFrom(t reflect.Type) (*GoType, error) {
	fields := []*GoField{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		gt, err := TypeFrom(f.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to translate field's %s type %v: %w",
				f.Name, f.Type, err)
		}
		fields = append(gt.Fields, NewGoField(f.Name, gt))
	}
	return AddImportInfo(NewStruct(t.Name(), fields), t.PkgPath(), ""), nil
}

func ArrayTypeFrom(t reflect.Type) (*GoType, error) {
	gt, err := TypeFrom(t.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to translate array element type %v: %w",
			t.Elem(), err)
	}
	return AddImportInfo(NewArray(gt), t.Key().PkgPath(), ""), nil
}

func GoKind(k reflect.Kind) string {
	switch k {
	case reflect.Array:
		return ArrayKind
	case reflect.Bool:
		return BoolKind
	case reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64:
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return IntKind
	case reflect.String:
		return StringKind
	case reflect.Struct:
		return StructKind
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
	default:
		panic(fmt.Sprintf("%s reflect.Kind: %#v", UnsupportedKind, k))
	}
	return ""
}
