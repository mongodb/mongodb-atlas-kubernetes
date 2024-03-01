package cmp

import (
	"cmp"
	"encoding/json"
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

func ByJSON[T any](x, y T) int {
	return cmp.Compare(mustJSONMarshal(x), mustJSONMarshal(y))
}

func mustJSONMarshal[T any](obj T) string {
	jObj, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(jObj)
}
