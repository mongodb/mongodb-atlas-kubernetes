// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gotype

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		SetAlias(MustTypeFrom(reflect.TypeOf(metav1.Condition{})), "metav1"),
	}
}
