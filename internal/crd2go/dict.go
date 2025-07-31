package crd2go

import "fmt"

// TypeDict is a dictionary of Go types, used to track and ensure unique type names.
// It also keeps track of generated types to avoid re-genrating the same type again.
type TypeDict struct {
	bySignature map[string]*GoType
	byName      map[string]*GoType
	generated   map[string]bool
	renames     map[string]string
}

// NewTypeDict creates a new TypeDict with the given renames and Go types
func NewTypeDict(renames map[string]string, goTypes ...*GoType) TypeDict {
	td := TypeDict{
		bySignature: make(map[string]*GoType),
		byName:      make(map[string]*GoType),
		generated:   make(map[string]bool),
		renames:    renames,
	}
	for _, gt := range goTypes {
		td.Add(gt)
	}
	return td
}

// Has checks if the TypeDict contains a GoType with the same signature
func (td TypeDict) Has(gt *GoType) bool {
	signature := gt.signature()
	_, ok := td.bySignature[signature]
	return ok
}

// Get retrieves a GoType by its name from the TypeDict
func (td TypeDict) Get(name string) (*GoType, bool) {
	gt, ok := td.byName[name]
	return gt, ok
}

// AddAll adds all the given tipes to the dictionary
func (td TypeDict) AddAll(goTypes ... *GoType) {
	for _, gt := range goTypes {
		td.Add(gt)
	}
}

// Add adds a GoType to the TypeDict, ensuring that the type name is unique
func (td TypeDict) Add(gt *GoType) {
	titledName := title(gt.Name)
	if gt.Name != titledName {
		panic(fmt.Sprintf("type name %s is not titled", gt.Name))
	}
	td.bySignature[gt.signature()] = gt
	td.byName[gt.Name] = gt
}

// MarkGenerated marks a GoType as generated
func (td TypeDict) MarkGenerated(gt *GoType) {
	if !td.Has(gt) {
		td.Add(gt)
	}
	td.generated[gt.Name] = true
}

// WasGenerated checks if a GoType was marked as generated
func (td TypeDict) WasGenerated(gt *GoType) bool {
	if td.Has(gt) {
		return td.generated[gt.Name]
	}
	return false
}

func (td TypeDict) Rename(name string) string {
	if len(td.renames) > 0 {
		if newName, ok := td.renames[name]; ok {
			return newName
		}
	}
	return name
}
