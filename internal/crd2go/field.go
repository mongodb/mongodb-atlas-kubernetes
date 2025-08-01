package crd2go

import "fmt"

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
		Name:   name,
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

// RenameType renames the GoType of the field to ensure it is unique within the
// TypeDict. It uses the parent names as needed to create a unique name for the
// type, if the type is not a primitive and its name is already taken.
func (f *GoField) RenameType(td *TypeDict, parentNames []string) error {
	if f.GoType == nil {
		return fmt.Errorf("failed to rename type for field %s: GoType is nil", f.Name)
	}
	if err := td.RenameType(parentNames, f.GoType); err != nil {
		return fmt.Errorf("failed to rename field type: %w", err)
	}
	return nil
}
