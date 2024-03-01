package cmp

import (
	"reflect"
	"slices"
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
