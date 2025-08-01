package crd

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/josvazg/crd2go/internal/gotype"
	"github.com/josvazg/crd2go/k8s"
)

const (
	OpenAPIObject  = "object"
	OpenAPIArray   = "array"
	OpenAPIString  = "string"
	OpenAPIInteger = "integer"
	OpenAPINumber  = "number"
	OpenAPIBoolean = "boolean"
)

type CRDType struct {
	Name    string
	Parents []string
	Schema  *apiextensionsv1.JSONSchemaProps
}

// FromOpenAPIType converts an OpenAPI schema to a GoType
func FromOpenAPIType(td *gotype.TypeDict, crdType *CRDType) (*gotype.GoType, error) {
	switch crdType.Schema.Type {
	case OpenAPIObject:
		return fromOpenAPIObject(td, crdType)
	case OpenAPIArray:
		return fromOpenAPIArray(td, crdType)
	case OpenAPIString, OpenAPIInteger, OpenAPINumber, OpenAPIBoolean:
		return fromOpenAPIFormattedType(crdType.Schema)
	default:
		return nil, fmt.Errorf("unsupported Open API type %q", crdType.Name)
	}
}

// fromOpenAPIObject converts an OpenAPI object schema into  unstructured JSON,
// a map or a GoType struct
func fromOpenAPIObject(td *gotype.TypeDict, crdType *CRDType) (*gotype.GoType, error) {
	if isUnstructured(crdType.Schema) {
		return gotype.JSONType, nil
	}
	if isDict(crdType.Schema) {
		return fromOpenAPIDict(td, crdType)
	}
	return fromOpenAPIStruct(td, crdType)
}

// fromOpenAPIStruct converts and OpenAPI object to a GoType struct
func fromOpenAPIStruct(td *gotype.TypeDict, crdType *CRDType) (*gotype.GoType, error) {
	fields := []*gotype.GoField{}
	fieldsParents := append(crdType.Parents, crdType.Name)
	for _, key := range orderedkeys(crdType.Schema.Properties) {
		props := crdType.Schema.Properties[key]
		fieldType, err := FromOpenAPIType(td, &CRDType{
			Name:    key,
			Parents: fieldsParents,
			Schema:  &props,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s type: %w", key, err)
		}
		field := gotype.NewGoFieldWithKey(key, key, fieldType)
		field.Comment = props.Description
		field.Required = slices.Contains(crdType.Schema.Required, key)
		if err := td.RenameField(field, fieldsParents); err != nil {
			return nil, fmt.Errorf("failed to rename field %v: %w", field, err)
		}
		fields = append(fields, field)
	}
	return gotype.NewStruct(crdType.Name, fields), nil
}

// fromOpenAPIArray converts an OpenAPI array schema to a GoType array
func fromOpenAPIArray(td *gotype.TypeDict, crdType *CRDType) (*gotype.GoType, error) {
	if crdType.Schema.Items == nil {
		return nil, fmt.Errorf("array %s has no items", crdType.Name)
	}
	if crdType.Schema.Items.Schema == nil {
		return nil, fmt.Errorf("array %s has no items schema", crdType.Name)
	}
	elementType, err := FromOpenAPIType(td, &CRDType{
		Name:   crdType.Name,
		Schema: crdType.Schema.Items.Schema,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse array %s element type: %w", crdType.Name, err)
	}
	if err := td.RenameType(crdType.Parents, elementType); err != nil {
		return nil, fmt.Errorf("failed to rename element type under %s: %w", crdType.Name, err)
	}
	return gotype.NewArray(elementType), nil
}

// fromOpenAPIDict converts an OpenAPI dictionary to a GoType map
func fromOpenAPIDict(td *gotype.TypeDict, crdType *CRDType) (*gotype.GoType, error) {
	elemType := gotype.JSONType
	if crdType.Schema.AdditionalProperties.Schema != nil {
		var err error
		elemType, err = FromOpenAPIType(td, &CRDType{
			Name:    crdType.Name,
			Parents: crdType.Parents,
			Schema:  crdType.Schema.AdditionalProperties.Schema,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to check map value type: %w", err)
		}
	}
	return &gotype.GoType{Name: gotype.MapKind, Kind: gotype.MapKind, Element: elemType}, nil
}

// fromOpenAPIFormattedType converts some OpenAPI formatted primitives to a hardwired GoType,
// or just fallsback to fromOpenAPIPrimitive
func fromOpenAPIFormattedType(schema *apiextensionsv1.JSONSchemaProps) (*gotype.GoType, error) {
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
	gt := gotype.Format2Builtin[gotype.FormatAliases[schema.Format]]
	if gt != nil {
		return gt, nil
	}
	return fromOpenAPIPrimitive(schema.Type)
}

// fromOpenAPIPrimitive converts an OpenAPI primitive type to a GoType
func fromOpenAPIPrimitive(kind string) (*gotype.GoType, error) {
	goTypeName, err := openAPIKindtoGoType(kind)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI kind %s: %w", kind, err)
	}
	return gotype.NewPrimitive(goTypeName, goTypeName), nil
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
		return gotype.StringKind, nil
	case OpenAPIInteger:
		return gotype.IntKind, nil
	case OpenAPINumber:
		return gotype.FloatKind, nil
	case OpenAPIBoolean:
		return gotype.BoolKind, nil
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

func KnownTypes() []*gotype.GoType {
	return []*gotype.GoType{
		gotype.MustTypeFrom(reflect.TypeOf(k8s.LocalReference{})),
		gotype.MustTypeFrom(reflect.TypeOf(k8s.Reference{})),
		gotype.SetAlias(gotype.MustTypeFrom(reflect.TypeOf(metav1.Condition{})), "metav1"),
	}
}

func CRD2Filename(crd *apiextensionsv1.CustomResourceDefinition) string {
	return fmt.Sprintf("%s.go", strings.ToLower(crd.Spec.Names.Kind))
}
