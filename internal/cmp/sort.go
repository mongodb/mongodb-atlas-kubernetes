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
