package crd2go

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	FirstVersion = ""
)

const (
	metav1 = "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateStream(w io.Writer, r io.Reader, version string) error {
	for {
		crd, err := ParseCRD(r)
		if errors.Is(err, io.EOF) {
			return err // EOF might be an error or a proper reply, so no wrapping
		}
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		stmt, err := Generate(crd, version)
		if err != nil {
			return fmt.Errorf("failed to generate CRD code: %w", err)
		}
		if _, err := w.Write(([]byte)(stmt.GoString())); err != nil {
			return fmt.Errorf("failed to write Go code: %w", err)
		}
	}
}

func Generate(crd *apiextensions.CustomResourceDefinition, version string) (*jen.File, error) {
	v := selectVersion(&crd.Spec, version)
	if v == nil {
		if version == "" {
			return nil, fmt.Errorf("no versions to generate code from")
		}
		return nil, fmt.Errorf("no version %q to generate code from", version)
	}
	f := jen.NewFile(v.Name)
	f.ImportAlias(metav1, "metav1")
	f.Func().Id("init").Params().Block(
		jen.Id("SchemeBuilder").Dot("Register").Params(
			jen.Op("&").Id("Group").Values(),
		),
	)

	if err := generateCRDRootObject(f, crd, v); err != nil {
		return nil, fmt.Errorf("failed to generate root object: %w", err)
	}
	return f, nil
}

// generateCRDRootObject generates the root object for the CRD
func generateCRDRootObject(f *jen.File, crd *apiextensions.CustomResourceDefinition, v *apiextensions.CustomResourceDefinitionVersion) error {
	f.Comment("+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object").Line()

	specType := fmt.Sprintf("%sSpec", crd.Spec.Names.Kind)
	statusType := fmt.Sprintf("%sStatus", crd.Spec.Names.Kind)

	code := f.Type().Id(crd.Spec.Names.Kind).Struct(
		jen.Qual(metav1, "TypeMeta").Tag(map[string]string{"json": ",inline"}),
		jen.Qual(metav1, "ObjectMeta").Tag(map[string]string{"json": "metadata,omitempty"}),
		jen.Line(),
		jen.Id("Spec").Id(specType).Tag(map[string]string{"json": "spec,omitempty"}),
		jen.Id("Status").Id(statusType).Tag(map[string]string{"json": "status,omitempty"}),
	)

	specSchema := v.Schema.OpenAPIV3Schema.Properties["spec"]
	specCode, err := generateOpenAPIType(specType, &specSchema)
	if err != nil {
		return fmt.Errorf("failed to generate spec code: %w", err)
	}
	code.Add(specCode)

	statusSchema := v.Schema.OpenAPIV3Schema.Properties["status"]
	statusCode, err := generateOpenAPIType(statusType, &statusSchema)
	if err != nil {
		return fmt.Errorf("failed to generate status code: %w", err)
	}
	code.Add(statusCode)
	return nil
}

// generateOpenAPIType generates a Go code type statement for the given OpenAPI schema
func generateOpenAPIType(typeName string, schema *apiextensions.JSONSchemaProps) (*jen.Statement, error) {
	switch schema.Type {
	case "object":
		return generateOpenAPIStruct(typeName, schema)
	default:
		return nil, fmt.Errorf("unsupported Open API type %q", schema.Type)
	}
}

// generateOpenAPIStruct generates a Go struct for the given OpenAPI schema
func generateOpenAPIStruct(typeName string, schema *apiextensions.JSONSchemaProps) (*jen.Statement, error) {
	subtypes := []jen.Code{}
	fields := []jen.Code{}
	for _, key := range orderedkeys(schema.Properties) {
		value := schema.Properties[key]
		id := title(fmt.Sprintf("%s%s", typeName, title(key)))
		tagValue := key
		typeSuffix, err := namedType(id, &value, slices.Contains(schema.Required, tagValue))
		if err != nil {
			return nil, fmt.Errorf("failed to parse schema type: %w", err)
		}
		entry := jen.Id(title(key)).Add(typeSuffix)
		if !slices.Contains(schema.Required, key) {
			tagValue = strings.Join([]string{tagValue, "omitempty"}, ",")
		}
		entry = entry.Tag(map[string]string{"json": tagValue})
		fields = append(fields, entry)
		for _, complexSubtype := range complexTypes(&value) {
			subtype, err := generateOpenAPIType(id, complexSubtype)
			if err != nil {
				return nil, fmt.Errorf("failed to parse schema type: %w", err)
			}
			subtypes = append(subtypes, subtype)
		}
	}
	mainType := jen.Line().Line().Type().Id(typeName).Struct(fields...)
	return mainType.Add(subtypes...), nil
}

// namedType generates a Go code statement for the given name and schema
func namedType(name string, schema *apiextensions.JSONSchemaProps, required bool) (*jen.Statement, error) {
	switch schema.Type {
	case "array":
		return requiredPrefix(required).Index().Id(title(name)), nil
	case "boolean":
		return requiredPrefix(required).Bool(), nil
	case "integer":
		return requiredPrefix(required).Int(), nil
	case "object":
		return requiredPrefix(required).Id(title(name)), nil
	case "string":
		return requiredPrefix(required).String(), nil
	default:
		return nil, fmt.Errorf("unsupported Open API type %q conversion to Go type", schema.Type)
	}
}

// requiredPrefix generates a code statement indicating whether a field is required or optional
func requiredPrefix(required bool) *jen.Statement {
	if required {
		return jen.Null()
	}
	return jen.Op("*")
}
// complexTypes returns a slice of JSONSchemaProps that represent complex types (objects or arrays) in the schema
func complexTypes(schema *apiextensions.JSONSchemaProps) []*apiextensions.JSONSchemaProps {
	switch schema.Type {
	case "object":
		return []*apiextensions.JSONSchemaProps{schema}
	case "array":
		if schema.Items.Schema != nil {
			return complexTypes(schema.Items.Schema)
		}
		schemas := []*apiextensions.JSONSchemaProps{}
		for _, schema := range schema.Items.JSONSchemas {
			schemas = append(schemas, complexTypes(&schema)...)
		}
		return schemas
	default:
		return nil
	}
}

// selectVersion returns the version from the CRD spec that matches the given version string
func selectVersion(spec *apiextensions.CustomResourceDefinitionSpec, version string) *apiextensions.CustomResourceDefinitionVersion {
	if len(spec.Versions) == 0 {
		return nil
	}
	if version == "" {
		return &spec.Versions[0]
	}
	for _, v := range spec.Versions {
		if v.Name == version {
			return &v
		}
	}
	return nil
}

// title capitalizes the first letter of a string and returns it using Go cases library
func title(s string) string {
	return cases.Upper(language.English).String(s[0:1]) + s[1:]
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
