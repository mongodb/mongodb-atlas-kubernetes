package gotype

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/josvazg/crd2go/k8s"
)

var (
	JSONType = BuiltInType("JSON", "apiextensionsv1", "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1")
	TimeType = BuiltInType("Time", "metav1", "k8s.io/apimachinery/pkg/apis/meta/v1")

	builtInTypes = map[string]*GoType{
		TimeType.Signature(): TimeType,
		JSONType.Signature(): JSONType,
	}
)

func BuiltInType(name, alias, path string) *GoType {
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

func KnownTypes() []*GoType {
	return []*GoType{
		MustTypeFrom(reflect.TypeOf(k8s.LocalReference{})),
		MustTypeFrom(reflect.TypeOf(k8s.Reference{})),
		SetAlias(MustTypeFrom(reflect.TypeOf(metav1.Condition{})), "metav1"),
	}
}
