package crd2go

import (
	"fmt"
	"path"
	"reflect"
	"slices"
	"sort"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/josvazg/crd2go/k8s"
)

const (
	UnsupportedKind = "unsupported"
	StructKind      = "struct"
	ArrayKind       = "array"
	StringKind      = "string"
	IntKind         = "int"
	Uint64Kind      = "uint64"
	FloatKind       = "float64"
	BoolKind        = "bool"
	MapKind         = "map"
	OpaqueKind      = "opaque"
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
	if gt.Kind == OpaqueKind {
		if gt.Import != nil {
			return fmt.Sprintf("%s.%s", gt.Import.Path, gt.Name)
		}
		return gt.Name
	}
	if gt.Kind == StructKind {
		if len(gt.Fields) == 0 { // de-duplicate empty structs as different types
			return fmt.Sprintf("{%s}", gt.Name)
		}
		fieldSignatures := make([]string, 0, len(gt.Fields))
		for _, field := range gt.Fields {
			fieldSignatures = append(fieldSignatures, field.signature())
		}
		return fmt.Sprintf("{%s}", strings.Join(fieldSignatures, ","))
	}
	if gt.Kind == ArrayKind {
		return fmt.Sprintf("[%s]", gt.Element.signature())
	}
	return fmt.Sprintf("%s", gt.Kind)
}

// baseType returns the base type of the GoType.
// If a type is an array, it returns the element type,
// traversing until a non-array type is found.
func (g *GoType) baseType() *GoType {
	if g == nil {
		return nil
	}
	if g.Kind == ArrayKind || g.Kind == MapKind {
		return g.Element.baseType()
	}
	return g
}

// GoField is a field in a Go struct
type GoField struct {
	Comment  string
	Required bool
	Name     string
	Key      string
	GoType   *GoType
}

// NewGoField creates a new GoField with the given name and GoType
func NewGoField(name string, gt *GoType) *GoField {
	return NewGoFieldWithKey(name, untitle(name), gt)
}

// NewGoFieldWithKey creates a new GoField with the given name, key, and GoType
func NewGoFieldWithKey(name, key string, gt *GoType) *GoField {
	return &GoField{
		Name:   name,
		Key:    key,
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
	if err := RenameType(td, parentNames, f.GoType); err != nil {
		return fmt.Errorf("failed to rename field type: %w", err)
	}
	return nil
}

// RenameType renames the given GoType to ensure it is unique within the
// TypeDict. It uses the parent names as needed to create a unique name for the
// type, if the type is not a primitive and its name is already taken.
func RenameType(td TypeDict, parentNames []string, gt *GoType) error {
	goType := gt.baseType()
	if goType.isPrimitive() {
		return nil
	}
	goType.Name = td.Rename(goType.Name)
	if td.Has(goType) {
		existingType := td.bySignature[goType.signature()]
		if existingType == nil {
			return fmt.Errorf("failed to find existing type for %v", gt)
		}
		goType.Name = existingType.Name
		goType.Import = existingType.Import
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
		return fmt.Errorf("failed to find a free type name for type %v", gt)
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

// NewOpaqueType creates a new GoType representing an opaque type with hidden internals
func NewOpaqueType(name string) *GoType {
	return &GoType{
		Name: title(name),
		Kind: OpaqueKind,
	}
}

func AddImportInfo(gt *GoType, alias, packagePath string) *GoType {
	effectiveAlias := alias
	if effectiveAlias == "" {
		effectiveAlias = path.Base(packagePath)
	}
	gt.Import = &ImportInfo{Path: packagePath, Alias: effectiveAlias}
	return gt
}

func setAlias(gt *GoType, alias string) *GoType {
	if gt.Import != nil {
		gt.Import.Alias = alias
	}
	return gt
}

// orderFieldsByName sorts the fields of a GoType by name
func orderFieldsByName(fields []*GoField) []*GoField {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

// FromOpenAPIType converts an OpenAPI schema to a GoType
func FromOpenAPIType(td TypeDict, typeName string, parents []string, schema *apiextensionsv1.JSONSchemaProps) (*GoType, error) {
	switch schema.Type {
	case OpenAPIObject:
		return fromOpenAPIStruct(td, typeName, parents, schema)
	case OpenAPIArray:
		return fromOpenAPIArray(td, typeName, parents, schema)
	case OpenAPIString, OpenAPIInteger, OpenAPINumber, OpenAPIBoolean:
		return fromOpenAPIFormattedType(schema)
	default:
		return nil, fmt.Errorf("unsupported Open API type %q", schema.Type)
	}
}

// fromOpenAPIStruct converts an OpenAPI object schema to a GoType struct
func fromOpenAPIStruct(td TypeDict, typeName string, parents []string, schema *apiextensionsv1.JSONSchemaProps) (*GoType, error) {
	if isUnstructured(schema) {
		return jsonType, nil
	}
	if isDict(schema) {
		return fromOpenAPIDict(td, typeName, parents, schema)
	}
	fields := []*GoField{}
	fieldsParents := append(parents, typeName)
	for _, key := range orderedkeys(schema.Properties) {
		props := schema.Properties[key]
		fieldType, err := FromOpenAPIType(td, key, fieldsParents, &props)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s type: %w", key, err)
		}
		field := NewGoFieldWithKey(title(key), key, fieldType)
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
func fromOpenAPIArray(td TypeDict, typeName string, parents []string, schema *apiextensionsv1.JSONSchemaProps) (*GoType, error) {
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
	if err := RenameType(td, parents, elementType); err != nil {
		return nil, fmt.Errorf("failed to rename element type under %s: %w", typeName, err)
	}
	return NewArray(elementType), nil
}

// fromOpenAPIFormattedType converts some OpenAPI formatted primitives to a hardwired GoType,
// or just fallsback to fromOpenAPIPrimitive
func fromOpenAPIFormattedType(schema *apiextensionsv1.JSONSchemaProps) (*GoType, error) {
	// - bsonobjectid: a bson object ID, i.e. a 24 characters hex string
	// - uri: an URI as parsed by Golang net/url.ParseRequestURI
	// - email: an email address as parsed by Golang net/mail.ParseAddress
	// - hostname: a valid representation for an Internet host name, as defined by RFC 1034, section 3.1 [RFC1034].
	// - ipv4: an IPv4 IP as parsed by Golang net.ParseIP
	// - ipv6: an IPv6 IP as parsed by Golang net.ParseIP
	// - cidr: a CIDR as parsed by Golang net.ParseCIDR
	// - mac: a MAC address as parsed by Golang net.ParseMAC
	// - uuid: an UUID that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$
	// - uuid3: an UUID3 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?3[0-9a-f]{3}-?[0-9a-f]{4}-?[0-9a-f]{12}$
	// - uuid4: an UUID4 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?4[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$
	// - uuid5: an UUID5 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?5[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$
	// - isbn: an ISBN10 or ISBN13 number string like "0321751043" or "978-0321751041"
	// - isbn10: an ISBN10 number string like "0321751043"
	// - isbn13: an ISBN13 number string like "978-0321751041"
	// - creditcard: a credit card number defined by the regex ^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$ with any non digit characters mixed in
	// - ssn: a U.S. social security number following the regex ^\\d{3}[- ]?\\d{2}[- ]?\\d{4}$
	// - hexcolor: an hexadecimal color code like "#FFFFFF: following the regex ^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$
	// - rgbcolor: an RGB color code like rgb like "rgb(255,255,2559"
	// - byte: base64 encoded binary data
	// - password: any kind of string
	// - date: a date string like "2006-01-02" as defined by full-date in RFC3339
	// - duration: a duration string like "22 ns" as parsed by Golang time.ParseDuration or compatible with Scala duration format
	// - datetime: a date time string like "2014-12-15T19:30:20.000Z" as defined by date-time in RFC3339.
	gt := format2Builtin[formatAliases[schema.Format]]
	if gt != nil {
		return gt, nil
	}
	return fromOpenAPIPrimitive(schema.Type)
}

// fromOpenAPIPrimitive converts an OpenAPI primitive type to a GoType
func fromOpenAPIPrimitive(kind string) (*GoType, error) {
	goTypeName, err := openAPIKindtoGoType(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI kind %s: %w", kind, err)
	}
	return NewPrimitive(goTypeName, goTypeName), nil
}

func fromOpenAPIDict(td TypeDict, typeName string, parents []string, schema *apiextensionsv1.JSONSchemaProps) (*GoType, error) {
	elemType := jsonType
	if schema.AdditionalProperties.Schema != nil {
		var err error
		elemType, err = FromOpenAPIType(td, typeName, parents, schema.AdditionalProperties.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to check map value type: %w", err)
		}
	}
	return &GoType{Name: MapKind, Kind: MapKind, Element: elemType}, nil
}

func isUnstructured(schema *apiextensionsv1.JSONSchemaProps) bool {
	return (len(schema.Properties) == 0 && schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields == true)
}

func isDict(schema *apiextensionsv1.JSONSchemaProps) bool {
	return schema.AdditionalProperties != nil
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
		setAlias(MustTypeFrom(reflect.TypeOf(metav1.Condition{})), "metav1"),
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
	builtInType := toBuiltInType(t)
	if builtInType != nil {
		return builtInType, nil
	}
	kind := GoKind(t.Kind())
	switch kind {
	case StructKind:
		return StructTypeFrom(t)
	case ArrayKind:
		return ArrayTypeFrom(t)
	case StringKind, IntKind, Uint64Kind, FloatKind, BoolKind:
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
		fields = append(fields, NewGoField(f.Name, gt))
	}
	return AddImportInfo(NewStruct(t.Name(), fields), "", t.PkgPath()), nil
}

func ArrayTypeFrom(t reflect.Type) (*GoType, error) {
	gt, err := TypeFrom(t.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to translate array element type %v: %w",
			t.Elem(), err)
	}
	return AddImportInfo(NewArray(gt), "", t.Key().PkgPath()), nil
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
