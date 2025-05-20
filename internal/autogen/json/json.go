package json

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func MustUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func Convert[T any](v any) *T {
	var result T
	MustUnmarshal(MustMarshal(v), &result)
	return &result
}

func ConvertNestedField[T any](obj map[string]interface{}, fields ...string) *T {
	u, _, _ := unstructured.NestedFieldCopy(obj, fields...)
	return Convert[T](u)
}
