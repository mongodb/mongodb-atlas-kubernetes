package gotype

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
		Name:   title(name),
		Key:    key,
		GoType: gt,
	}
}

// Signature generates a unique signature for a GoField using the type Signature
func (g *GoField) Signature() string {
	if g == nil {
		return "nil"
	}
	return fmt.Sprintf("%s:%s", g.Name, g.GoType.Signature())
}

// untitle de-capitalizes the first letter of a string and returns it using Go cases library
func untitle(s string) string {
	if s == "" {
		return ""
	}
	return cases.Lower(language.English).String(s[0:1]) + s[1:]
}
