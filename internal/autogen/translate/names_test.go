// Copyright 2025 Google LLC
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

package translate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/translate"
)

func TestPrefixedName(t *testing.T) {
	for _, tc := range []struct {
		title  string
		prefix string
		name   string
		args   []string
		want   string
	}{
		{
			title:  "just A",
			prefix: "just",
			name:   "A",
			want:   "just-56b7d6667d8f6cc88c8d",
		},
		{
			title:  "a very long name",
			prefix: "a",
			name:   "very logn name with several parts",
			want:   "a-54fd46599bbf74b8db44",
		},
		{
			title:  "names",
			prefix: "several",
			name:   "names",
			args:   []string{"name0", "name1"},
			want:   "several-75db99b58df57d54bd6",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := translate.PrefixedName(tc.prefix, tc.name, tc.args...)
			assert.Equal(t, tc.want, got)
		})
	}
}
