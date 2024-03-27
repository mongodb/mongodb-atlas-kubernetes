package cmp

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"time"
)

type Normalizer[T any] interface {
	Normalize() T
}

func SemanticEqual[T Normalizer[T]](this, that T) bool {
	return reflect.DeepEqual(this.Normalize(), that.Normalize())
}

func NormalizeSlice[S ~[]E, E any](slice S, cmp func(a, b E) int) S {
	if len(slice) == 0 {
		return nil
	}
	slices.SortFunc(slice, cmp)
	return slice
}

func Normalize(data any) error {
	var err error
	traverse(data, func(slice reflect.Value) {
		sort.Slice(slice.Interface(), func(i, j int) bool {
			result, e := ByJSON(
				slice.Index(i).Interface(),
				slice.Index(j).Interface(),
			)
			if e != nil {
				err = fmt.Errorf("error converting slice %v to JSON: %w", slice, e)
				return false
			}
			return result < 0
		})
	})
	return err
}

func PermuteOrder(data any) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	traverse(data, func(slice reflect.Value) {
		sliceIface := slice.Interface()
		for i, j := range r.Perm(slice.Len()) {
			reflect.Swapper(sliceIface)(i, j)
		}
	})
}

func traverse(data any, f func(slice reflect.Value)) {
	traverseValue(reflect.ValueOf(data), f)
}

func traverseValue(value reflect.Value, f func(slice reflect.Value)) {
	switch value.Kind() {
	case reflect.Pointer:
		// if it is a pointer, traverse over its dereferenced value
		traverseValue(value.Elem(), f)

	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			// skip unexported fields
			if value.Type().Field(i).PkgPath != "" {
				continue
			}
			// traverse over each field in the struct
			traverseValue(value.Field(i), f)
		}

	case reflect.Slice:
		// omit zero length slices
		if value.Len() == 0 {
			return
		}
		// skip byte slices
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return
		}
		// traverse over each element in the slice
		for j := 0; j < value.Len(); j++ {
			traverseValue(value.Index(j), f)
		}
		// base case: we can apply the given function
		f(value)
	}
}
