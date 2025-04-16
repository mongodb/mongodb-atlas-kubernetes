// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmp

import (
	"cmp"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sort"
)

type Normalizer[T any] interface {
	Normalize() (T, error)
}

func SemanticEqual[T Normalizer[T]](this, that T) (bool, error) {
	thisResult, thisError := this.Normalize()
	thatResult, thatError := that.Normalize()
	if thisError != nil {
		return false, thisError
	}
	if thatError != nil {
		return false, thatError
	}
	return reflect.DeepEqual(thisResult, thatResult), nil
}

func NormalizeSlice[S ~[]E, E any](slice S, cmp func(a, b E) int) S {
	if len(slice) == 0 {
		return nil
	}
	slices.SortStableFunc(slice, cmp)
	return slice
}

func Normalize(data any) error {
	var err error
	traverse(data, func(slice reflect.Value) {
		sort.Slice(slice.Interface(), func(i, j int) bool {
			iIface, jIface := slice.Index(i).Interface(), slice.Index(j).Interface()

			if ok, result := compareSortable(iIface, jIface); ok {
				return result < 0
			}

			result, e := ByJSON(iIface, jIface)
			if e != nil {
				err = fmt.Errorf("error converting slice %v to JSON: %w", slice, e)
				return false
			}
			return result < 0
		})
	})
	return err
}

func compareSortable(i, j any) (bool, int) {
	iSortable, iSortableOK := i.(Sortable)
	jSortable, jSortableOK := j.(Sortable)

	if iSortableOK && jSortableOK {
		return true, cmp.Compare(iSortable.Key(), jSortable.Key())
	}

	return false, -1
}

func PermuteOrder(data any, r *rand.Rand) {
	traverse(data, func(slice reflect.Value) {
		r.Shuffle(slice.Len(), func(i, j int) {
			reflect.Swapper(slice.Interface())(i, j)
		})
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
		// skip []byte slices
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
