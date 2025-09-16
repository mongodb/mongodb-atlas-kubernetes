package gotype

import (
	"fmt"
	"log"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/pkg/config"
)

// TypeDict is a dictionary of Go types, used to track and ensure unique type names.
// It also keeps track of generated types to avoid re-genrating the same type again.
type TypeDict struct {
	bySignature map[string]*GoType
	byName      map[string]*GoType
	generated   map[string]bool
	renames     map[string]string
}

// Request holds the runtime information to handle a CRD generation request
type Request struct {
	config.CoreConfig
	CodeWriterFn config.CodeWriterFunc
	TypeDict     *TypeDict
}

// NewTypeDict creates a new TypeDict with the given renames and Go types
func NewTypeDict(renames map[string]string, goTypes ...*GoType) *TypeDict {
	td := TypeDict{
		bySignature: make(map[string]*GoType),
		byName:      make(map[string]*GoType),
		generated:   make(map[string]bool),
		renames:     renames,
	}
	for _, gt := range goTypes {
		td.Add(gt)
	}
	return &td
}

// Has checks if the TypeDict contains a GoType with the same signature
func (td TypeDict) Has(gt *GoType) bool {
	signature := gt.Signature()
	_, ok := td.bySignature[signature]
	return ok
}

// Get retrieves a GoType by its name from the TypeDict
func (td TypeDict) Get(name string) (*GoType, bool) {
	gt, ok := td.byName[name]
	return gt, ok
}

// AddAll adds all the given tipes to the dictionary
func (td TypeDict) AddAll(goTypes ...*GoType) {
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
	td.bySignature[gt.Signature()] = gt
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

// RenameField renames the GoType of the field to ensure it is unique within the
// TypeDict. It uses the parent names as needed to create a unique name for the
// type, if the type is not a primitive and its name is already taken.
func (td *TypeDict) RenameField(gf *GoField, parentNames []string) error {
	if gf.GoType == nil {
		return fmt.Errorf("failed to rename type for field %s: GoType is nil", gf.Name)
	}
	if err := td.RenameType(parentNames, gf.GoType); err != nil {
		return fmt.Errorf("failed to rename field type: %w", err)
	}
	return nil
}

// RenameType renames the given GoType to ensure it is unique within the
// TypeDict. It uses the parent names as needed to create a unique name for the
// type, if the type is not a primitive and its name is already taken.
func (td TypeDict) RenameType(parentNames []string, gt *GoType) error {
	goType := gt.BaseType()
	if goType.IsPrimitive() {
		return nil
	}
	goType.Name = td.rename(goType.Name)
	if importInfo := td.matchImport(goType); importInfo != nil {
		goType.Import = importInfo
		return nil
	}
	if td.Has(goType) {
		existingType := td.bySignature[goType.Signature()]
		if existingType == nil {
			return fmt.Errorf("failed to find existing type for %v", gt)
		}
		goType.Name = existingType.Name
		goType.Import = existingType.Import
		return nil
	}

	typeName := goType.Name
	for i := len(parentNames) - 1; i >= 0; i-- {
		_, used := td.Get(typeName)
		if !used {
			break
		}
		typeName = fmt.Sprintf("%s%s", title(parentNames[i]), typeName)
	}

	_, used := td.Get(typeName)
	if used {
		return fmt.Errorf("failed to find a free type name for type %v", gt)
	}
	goType.Name = typeName
	if goType.Name == "string" {
		log.Printf("here")
	}
	td.Add(goType)

	return nil
}

// rename applies custom renames before automated renaming logic
func (td TypeDict) rename(name string) string {
	if len(td.renames) > 0 {
		if newName, ok := td.renames[name]; ok {
			return newName
		}
	}
	return name
}

// matchImport checks if the given type matches a registered auto import type.
// If so, it updated the type structure in the type dictionary entry,
// so that it can be matched by signature going forward and
// returns the matching import info.
func (td TypeDict) matchImport(gt *GoType) *config.ImportInfo {
	entry, ok := td.Get(gt.Name)
	if !ok || entry.Kind != AutoImportKind {
		return nil
	}
	entry.CloneStructure(gt)
	td.byName[entry.Name] = entry
	td.bySignature[entry.Signature()] = entry
	return entry.Import
}
