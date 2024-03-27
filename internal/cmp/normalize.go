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
				reflect.ValueOf(slice.Interface()).Index(i).Interface(),
				reflect.ValueOf(slice.Interface()).Index(j).Interface(),
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

			// skip byte slices
			if fieldValue.Type().Elem().Kind() == reflect.Uint8 {
				continue
			}

			// must be a slice
			f(fieldValue)
		}
	}
}
