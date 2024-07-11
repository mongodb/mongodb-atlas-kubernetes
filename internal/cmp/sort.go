package cmp

import (
	"cmp"
	"encoding/json"
	"fmt"
	"strings"
)

type Sortable interface {
	Key() string
}

func PointerKey[T Sortable](in *T) string {
	if in == nil {
		return "nil"
	}
	return (*in).Key()
}

func SliceKey[K Sortable](in []K) string {
	result := make([]string, 0, len(in))
	for i := range in {
		result = append(result, in[i].Key())
	}
	return "[" + strings.Join(result, ",") + "]"
}

func ByKey[S Sortable](x, y S) int {
	return cmp.Compare(x.Key(), y.Key())
}

func ByJSON[T any](x, y T) (int, error) {
	xJSON, err := marshalJSON(x)
	if err != nil {
		return -1, fmt.Errorf("error converting %v to JSON: %w", x, err)
	}
	yJSON, err := marshalJSON(y)
	if err != nil {
		return -1, fmt.Errorf("error converting %v to JSON: %w", x, err)
	}
	return cmp.Compare(xJSON, yJSON), nil
}

func marshalJSON[T any](obj T) (string, error) {
	jObj, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jObj), nil
}

func JSON[T any](obj T) []byte {
	jObj, err := json.MarshalIndent(obj, "  ", "  ")
	if err != nil {
		return ([]byte)(err.Error())
	}
	return jObj
}

func JSONize[T any](obj T) string {
	return string(JSON(obj))
}
