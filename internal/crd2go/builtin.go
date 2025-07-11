package crd2go

import (
	"fmt"
	"reflect"
)

var (
	formatAliases = map[string]string{
		"date-time": "datetime",
		"datetime":  "datetime",
	}

	timeType = builtInType("Time", "metav1", "k8s.io/apimachinery/pkg/apis/meta/v1")
	jsonType = builtInType("JSON", "apiextensionsv1", "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1")

	builtInTypes = map[string]*GoType{
		timeType.signature(): timeType,
		jsonType.signature(): jsonType,
	}

	format2Builtin = map[string]*GoType{
		"datetime": builtInTypes[timeType.signature()],
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
