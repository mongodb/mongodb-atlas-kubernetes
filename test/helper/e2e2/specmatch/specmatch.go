// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package specmatch checks that the state Atlas reports matches the spec that was requested.
//
// The spec is the expected result: every field set in the spec must come back from Atlas
// with the same value. Fields Atlas adds (ids, computed state) are not compared.
// Both sides are compared as plain JSON, so the check does not depend on the operator's
// own CR to Atlas translation, where fields can be silently dropped.
package specmatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
)

// Option adjusts how a field is compared, for the cases where Atlas legitimately
// returns something other than what was requested.
type Option func(*options)

type options struct {
	ignored    []string
	unordered  []string
	normalizer map[string]func(any) any
}

// IgnoreFields skips the given JSON paths, e.g. "$.password" for fields Atlas never returns.
func IgnoreFields(paths ...string) Option {
	return func(o *options) { o.ignored = append(o.ignored, paths...) }
}

// UnorderedLists compares the lists at the given JSON paths without regard to order.
func UnorderedLists(paths ...string) Option {
	return func(o *options) { o.unordered = append(o.unordered, paths...) }
}

// Normalize applies fn to both the requested and the returned value at path before comparing.
// The path can point at a scalar, an object or a whole list; use "$.list[]" to reach list elements.
func Normalize(path string, fn func(any) any) Option {
	return func(o *options) { o.normalizer[path] = fn }
}

// AtlasMatchesSpec returns an error listing every field set in spec whose value differs
// from, or is missing in, atlas. Both arguments are marshaled to JSON before comparing.
func AtlasMatchesSpec(spec, atlas any, opts ...Option) error {
	o := &options{normalizer: map[string]func(any) any{}}
	for _, opt := range opts {
		opt(o)
	}

	want, err := toJSONValue(spec)
	if err != nil {
		return fmt.Errorf("failed to convert spec to JSON: %w", err)
	}
	got, err := toJSONValue(atlas)
	if err != nil {
		return fmt.Errorf("failed to convert Atlas response to JSON: %w", err)
	}

	diffs := o.compare("$", want, got)
	if len(diffs) == 0 {
		return nil
	}
	sort.Strings(diffs)
	return errors.New("Atlas does not match the requested spec:\n" + strings.Join(diffs, "\n"))
}

func (o *options) compare(path string, want, got any) []string {
	if slices.Contains(o.ignored, path) {
		return nil
	}
	if fn, ok := o.normalizer[path]; ok {
		want, got = fn(want), fn(got)
	}

	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return []string{fmt.Sprintf("%s: want object, got %v", path, got)}
		}
		var diffs []string
		for key, wv := range w {
			fieldPath := path + "." + key
			gv, ok := g[key]
			if !ok && !slices.Contains(o.ignored, fieldPath) {
				diffs = append(diffs, fmt.Sprintf("%s: want %v, missing in Atlas", fieldPath, wv))
				continue
			}
			diffs = append(diffs, o.compare(fieldPath, wv, gv)...)
		}
		return diffs
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return []string{fmt.Sprintf("%s: want %v, got %v", path, want, got)}
		}
		if slices.Contains(o.unordered, path) {
			return o.compareUnordered(path, w, g)
		}
		var diffs []string
		for i := range w {
			diffs = append(diffs, o.compare(path+"[]", w[i], g[i])...)
		}
		return diffs
	default:
		if !reflect.DeepEqual(want, got) {
			return []string{fmt.Sprintf("%s: want %v, got %v", path, want, got)}
		}
		return nil
	}
}

// compareUnordered matches every requested element to a distinct returned element.
// Returned elements may carry extra fields, so one returned element can match several
// requested ones; augmenting paths find a full pairing whenever one exists.
func (o *options) compareUnordered(path string, want, got []any) []string {
	matches := make([][]bool, len(want))
	for i, wv := range want {
		matches[i] = make([]bool, len(got))
		for j, gv := range got {
			matches[i][j] = len(o.compare(path+"[]", wv, gv)) == 0
		}
	}
	owner := make([]int, len(got))
	for j := range owner {
		owner[j] = -1
	}
	var assign func(i int, seen []bool) bool
	assign = func(i int, seen []bool) bool {
		for j := range got {
			if !matches[i][j] || seen[j] {
				continue
			}
			seen[j] = true
			if owner[j] == -1 || assign(owner[j], seen) {
				owner[j] = i
				return true
			}
		}
		return false
	}
	for i := range want {
		if !assign(i, make([]bool, len(got))) {
			return []string{fmt.Sprintf("%s: want %v, got %v", path, want, got)}
		}
	}
	return nil
}

func toJSONValue(v any) (any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
