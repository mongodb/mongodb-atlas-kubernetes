package gotype

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GoField is a field in a Go struct
type GoField struct {
	Comment       string
	Required      bool
	Name          string
	Key           string
	GoType        *GoType
	CustomJSONTag string
}

type GoFieldOptionFunc func(*GoField)

// JSONTag option sets a csutom JSON tag
func JSONTag(jsonTag string) GoFieldOptionFunc {
	return func(gt *GoField) {
		gt.CustomJSONTag = jsonTag
	}
}

// Required option sets the required flag
func Required(required bool) GoFieldOptionFunc {
	return func(gt *GoField) {
		gt.Required = required
	}
}

// NewEmbeddedField creates and embedded field, without a name
func NewEmbeddedField(gt *GoType) *GoField {
	return NewGoField("", gt)
}

// NewGoField creates a new GoField with the given name and GoType
func NewGoField(name string, gt *GoType) *GoField {
	return NewGoFieldWithKey(name, untitle(name), gt)
}

// NewGoFieldWithKey creates a new GoField with the given name, key, and GoType
func NewGoFieldWithKey(name, key string, gt *GoType) *GoField {
	return &GoField{
		Name:   title(name),
		Key:    key,
		GoType: gt,
	}
}

func (gt *GoField) WithOptions(opts ...GoFieldOptionFunc) *GoField {
	for _, opt := range opts {
		opt(gt)
	}
	return gt
}

// Signature generates a unique signature for a GoField using the type Signature
func (gt *GoField) Signature() string {
	if gt == nil {
		return "nil"
	}
	return fmt.Sprintf("%s:%s", gt.Name, gt.GoType.Signature())
}

// IsEmbedded returns whether nor not his field is an embedded one
func (gt *GoField) IsEmbedded() bool {
	return gt.Name == ""
}

// untitle de-capitalizes the first letter of a string and returns it using Go cases library
func untitle(s string) string {
	if s == "" {
		return ""
	}
	return cases.Lower(language.English).String(s[0:1]) + s[1:]
}
