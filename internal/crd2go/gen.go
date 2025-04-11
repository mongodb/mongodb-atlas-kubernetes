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
			return fmt.Errorf("generation failed read input: %w", err)
		}
		stmt, err := Generate(crd, version)
		if err != nil {
			return fmt.Errorf("crd code generation failed: %w", err)
		}
		if _, err := w.Write(([]byte)(stmt.GoString())); err != nil {
			return fmt.Errorf("code writing failed: %w", err)
		}
	}
}

func Generate(crd *CRD, version string) (*jen.File, error) {
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

	if err := genRootObject(f, crd, v); err != nil {
		return nil, fmt.Errorf("failed to generate root object: %w", err)
	}
	return f, nil
}

func genRootObject(f *jen.File, crd *CRD, v *Version) error {
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
	specCode, err := generateOpenAPICode(specType, &specSchema)
	if err != nil {
		return fmt.Errorf("failed to generate spec code: %w", err)
	}
	code.Add(specCode)

	statusSchema := v.Schema.OpenAPIV3Schema.Properties["status"]
	statusCode, err := generateOpenAPICode(statusType, &statusSchema)
	if err != nil {
		return fmt.Errorf("failed to generate status code: %w", err)
	}
	code.Add(statusCode)
	return nil
}

func generateOpenAPICode(typeName string, schema *OpenAPISchema) (*jen.Statement, error) {
	switch schema.Type {
	case "object":
		return generateOpenAPIStruct(typeName, schema)
	default:
		return nil, fmt.Errorf("unsupported Open API type %q", schema.Type)
	}
}

func generateOpenAPIStruct(typeName string, schema *OpenAPISchema) (*jen.Statement, error) {
	subtypes := []jen.Code{}
	fields := []jen.Code{}
	for _, key := range orderedkeys(schema.Properties) {
		value := schema.Properties[key]
		id := title(fmt.Sprintf("%s%s", typeName, title(key)))
		tagValue := key
		typeSuffix, err := namedType(id, &value, required(tagValue, schema))
		if err != nil {
			return nil, fmt.Errorf("failed to parse schema type: %w", err)
		}
		entry := jen.Id(title(key)).Add(typeSuffix)
		if schema.Required == nil || !slices.Contains(*schema.Required, key) {
			tagValue = strings.Join([]string{tagValue, "omitempty"}, ",")
		}
		entry = entry.Tag(map[string]string{"json": tagValue})
		fields = append(fields, entry)
		if complexSubtype := complexType(&value); complexSubtype != nil {
			subtype, err := generateOpenAPICode(id, complexSubtype)
			if err != nil {
				return nil, fmt.Errorf("failed to parse schema type: %w", err)
			}
			subtypes = append(subtypes, subtype)
		}
	}
	mainType := jen.Line().Line().Type().Id(typeName).Struct(fields...)
	return mainType.Add(subtypes...), nil
}

func namedType(name string, schema *OpenAPISchema, required bool) (*jen.Statement, error) {
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

func requiredPrefix(required bool) *jen.Statement {
	if required {
		return jen.Null()
	}
	return jen.Op("*")
}

func required(field string, schema *OpenAPISchema) bool {
	if schema.Required == nil {
		return false
	}
	for _, required := range *schema.Required {
		if required == field {
			return true
		}
	}
	return false
}

func complexType(schema *OpenAPISchema) *OpenAPISchema {
	switch schema.Type {
	case "object":
		return schema
	case "array":
		return complexType(schema.Items)
	default:
		return nil
	}
}

func selectVersion(spec *Spec, version string) *Version {
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

func title(s string) string {
	return cases.Upper(language.English).String(s[0:1]) + s[1:]
}

func orderedkeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
