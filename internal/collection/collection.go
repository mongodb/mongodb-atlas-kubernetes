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

package collection

func CopyWithSkip[T comparable](list []T, skip T) []T {
	newList := make([]T, 0, len(list))

	for _, item := range list {
		if item != skip {
			newList = append(newList, item)
		}
	}

	return newList
}

func Keys[K comparable, V any](m map[K]V) []K {
	s := make([]K, 0, len(m))

	for k := range m {
		s = append(s, k)
	}

	return s
}

func MapDiff[K comparable, V any](a, b map[K]V) map[K]V {
	d := make(map[K]V, len(a))
	for i, val := range a {
		if _, ok := b[i]; !ok {
			d[i] = val
		}
	}

	return d
}
