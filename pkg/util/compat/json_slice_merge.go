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
func JSONSliceMerge(dst, src interface{}) error {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)

	if dstVal.Kind() != reflect.Ptr {
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

	minLen := dstVal.Len()
	if srcVal.Len() < minLen {
		minLen = srcVal.Len()
	}

	// merge common elements
	for i := 0; i < minLen; i++ {
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
