package crd2go

import (
	"fmt"
	"reflect"
)

func MustTypeFrom(t reflect.Type) *GoType {
	gt, err := TypeFrom(t)
	if err != nil {
		panic(fmt.Errorf("failed to translate type %v: %w", t.Name(), err))
	}
	return gt
}

func TypeFrom(t reflect.Type) (*GoType, error) {
	builtInType := toBuiltInType(t)
	if builtInType != nil {
		return builtInType, nil
	}
	kind := goKind(t.Kind())
	switch kind {
	case StructKind:
		return structTypeFrom(t)
	case ArrayKind:
		return arrayTypeFrom(t)
	case StringKind, IntKind, Uint64Kind, FloatKind, BoolKind:
		return NewPrimitive(t.Name(), kind), nil
	default:
		return nil, fmt.Errorf("unsupported kind %v", kind)
	}
}

func structTypeFrom(t reflect.Type) (*GoType, error) {
	fields := []*GoField{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		gt, err := TypeFrom(f.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to translate field's %s type %v: %w",
				f.Name, f.Type, err)
		}
		fields = append(fields, NewGoField(f.Name, gt))
	}
	return AddImportInfo(NewStruct(t.Name(), fields), "", t.PkgPath()), nil
}

func arrayTypeFrom(t reflect.Type) (*GoType, error) {
	gt, err := TypeFrom(t.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to translate array element type %v: %w",
			t.Elem(), err)
	}
	return AddImportInfo(NewArray(gt), "", t.Key().PkgPath()), nil
}

func goKind(k reflect.Kind) string {
	switch k {
	case reflect.Array:
		return ArrayKind
	case reflect.Bool:
		return BoolKind
	case reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64:
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return IntKind
	case reflect.String:
		return StringKind
	case reflect.Struct:
		return StructKind
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
	default:
		panic(fmt.Sprintf("%s reflect.Kind: %#v", UnsupportedKind, k))
	}
	return ""
}
