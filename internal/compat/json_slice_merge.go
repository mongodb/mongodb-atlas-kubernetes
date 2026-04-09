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

package compat

import (
	"errors"
	"fmt"
	"reflect"
)

// JSONSliceMerge will merge two slices using JSONCopy according to these rules:
//
// 1. If `dst` and `src` are the same length, all elements are merged
//
// 2. If `dst` is longer, only the first `len(src)` elements are merged
//
// 3. If `src` is longer, first `len(dst)` elements are merged, then remaining elements are appended to `dst`
func JSONSliceMerge(dst, src any) error {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)

	if dstVal.Kind() != reflect.Pointer {
		return errors.New("dst must be a pointer to slice")
	}

	dstVal = reflect.Indirect(dstVal)
	srcVal = reflect.Indirect(srcVal)

	if dstVal.Kind() != reflect.Slice {
		return errors.New("dst must be pointing to a slice")
	}

	if srcVal.Kind() != reflect.Slice {
		return errors.New("src must be a slice or a pointer to slice")
	}

	minLen := min(srcVal.Len(), dstVal.Len())

	// merge common elements
	for i := range minLen {
		dstX := dstVal.Index(i).Addr().Interface()
		if err := JSONCopy(dstX, srcVal.Index(i).Interface()); err != nil {
			return fmt.Errorf("cannot copy value at index %d: %w", i, err)
		}
	}

	// append extra elements (if any)
	dstType := reflect.TypeOf(dst).Elem().Elem()
	for i := minLen; i < srcVal.Len(); i++ {
		newVal := reflect.New(dstType).Interface()
		if err := JSONCopy(&newVal, srcVal.Index(i).Interface()); err != nil {
			return fmt.Errorf("cannot copy value at index %d: %w", i, err)
		}
		dstVal.Set(reflect.Append(dstVal, reflect.ValueOf(newVal).Elem()))
	}

	return nil
}
