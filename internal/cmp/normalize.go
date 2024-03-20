package cmp

import (
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

func NormalizeSliceUsingJSON[S ~[]E, E any](slice S) S {
	if len(slice) == 0 {
		return nil
	}
	slices.SortFunc(slice, ByJSON[E])
	return slice
}

func Normalize(data any) {
	traverse(data, func(slice reflect.Value) {
		sort.Slice(slice.Interface(), func(i, j int) bool {
			return ByJSON(
				reflect.ValueOf(slice.Interface()).Index(i).Interface(),
				reflect.ValueOf(slice.Interface()).Index(j).Interface(),
			) < 0
		})
	})
}

func PermuteOrder(data any) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	traverse(data, func(slice reflect.Value) {
		for i, j := range r.Perm(slice.Len()) {
			reflect.Swapper(slice.Interface())(i, j)
		}
	})
}

func traverse(data any, f func(slice reflect.Value)) {
	value := reflect.ValueOf(data)

	// if the value is a pointer, dereference it
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// if it's not a struct, return
	if value.Kind() != reflect.Struct {
		return
	}

	// this must be a struct
	for i := 0; i < value.NumField(); i++ {
		if value.Type().Field(i).PkgPath != "" {
			// skip unexported fields
			continue
		}
		fieldValue := value.Field(i)

		if fieldValue.Kind() == reflect.Struct {
			traverse(fieldValue.Addr().Interface(), f)
		}

		if fieldValue.Kind() == reflect.Slice {
			for j := 0; j < fieldValue.Len(); j++ {
				traverse(fieldValue.Index(j).Addr().Interface(), f)
			}

			if fieldValue.Len() == 0 {
				continue
			}

			// must be a slice
			f(fieldValue)
		}
	}
}
