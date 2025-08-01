package gotype

import (
	"fmt"
	"reflect"
)

var (
	FormatAliases = map[string]string{
		"date-time": "datetime",
		"datetime":  "datetime",
	}

	Format2Builtin = map[string]*GoType{
		"datetime": builtInTypes[timeType.Signature()],
	}

	JSONType = builtInType("JSON", "apiextensionsv1", "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1")
)

var (
	timeType = builtInType("Time", "metav1", "k8s.io/apimachinery/pkg/apis/meta/v1")

	builtInTypes = map[string]*GoType{
		timeType.Signature(): timeType,
		JSONType.Signature(): JSONType,
	}
)

func builtInType(name, alias, path string) *GoType {
	return AddImportInfo(NewOpaqueType(name), alias, path)
}

func toBuiltInType(t reflect.Type) *GoType {
	builtInKey := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	gt, ok := builtInTypes[builtInKey]
	if ok {
		return gt
	}
	return nil
}
