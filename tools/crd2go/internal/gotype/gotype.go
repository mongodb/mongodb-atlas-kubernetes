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
	"path"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/pkg/config"
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
	AutoImportKind  = "autoImport"
)

const PACKAGE_BASE = "github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go"

// GoType represents a Go type, which can be a primitive type, a struct, or an array.
// It is used in conjunbction with TypeDict to track and ensure unique type names.
type GoType struct {
	Name    string
	Kind    string
	Fields  []*GoField
	Element *GoType
	Import  *config.ImportInfo
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

// NewAutoImportType creates a new GoType representing an opaque type with hidden internals
func NewAutoImportType(importType *config.ImportedTypeConfig) *GoType {
	return &GoType{
		Name:   title(importType.Name),
		Kind:   AutoImportKind,
		Import: &importType.ImportInfo,
	}
}

// AddImportInfo allows to attach teh import information to a type
func AddImportInfo(gt *GoType, alias, packagePath string) *GoType {
	effectiveAlias := alias
	if effectiveAlias == "" {
		effectiveAlias = path.Base(packagePath)
	}
	gt.Import = &config.ImportInfo{Path: packagePath, Alias: effectiveAlias}
	return gt
}

// SetAlias allows to attach an alias to the import information of a type
func SetAlias(gt *GoType, alias string) *GoType {
	if gt.Import != nil {
		gt.Import.Alias = alias
	}
	return gt
}

// IsPrimitive checks if the GoType is a primitive type
func (gt *GoType) IsPrimitive() bool {
	switch gt.Kind {
	case StringKind, IntKind, FloatKind, BoolKind:
		return true
	default:
		return false
	}
}

// signature generates a unique signature for a GoType reflecting its structure.
// This is leveraged by TypeDict to check if a type is already registered with
// the same internal structure, regardless of the name.
func (gt *GoType) Signature() string {
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
			fieldSignatures = append(fieldSignatures, field.Signature())
		}
		return fmt.Sprintf("{%s}", strings.Join(fieldSignatures, ","))
	}
	if gt.Kind == ArrayKind {
		return fmt.Sprintf("[%s]", gt.Element.Signature())
	}
	return gt.Kind
}

// BaseType returns the base type of the GoType.
// If a type is an array, it returns the element type,
// traversing until a non-array type is found.
func (gt *GoType) BaseType() *GoType {
	if gt == nil {
		return nil
	}
	if gt.Kind == ArrayKind || gt.Kind == MapKind {
		return gt.Element.BaseType()
	}
	return gt
}

// CloneStructure copies the structure of another type,
// but leaved the name and import info intact
func (gt *GoType) CloneStructure(ot *GoType) {
	gt.Kind = ot.Kind
	gt.Fields = ot.Fields
	gt.Element = ot.Element
}

// orderFieldsByName sorts the fields of a GoType by name
func orderFieldsByName(fields []*GoField) []*GoField {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

// title capitalizes the first letter of a string and returns it using Go cases library
func title(s string) string {
	if s == "" {
		return ""
	}
	s = strings.TrimLeft(s, "_") // remove leading underscores
	return cases.Upper(language.English).String(s[0:1]) + s[1:]
}
