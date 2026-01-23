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
//

package objmap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/objmap"
)

func TestAsPath(t *testing.T) {
	for _, tc := range []struct {
		title string
		input string
		want  []string
	}{
		{
			title: "nothing",
			input: "",
			want:  []string{""},
		},
		{
			title: "single",
			input: "dir",
			want:  []string{"dir"},
		},
		{
			title: "double",
			input: "dir0.dir1",
			want:  []string{"dir0", "dir1"},
		},
		{
			title: "rooted double",
			input: ".dir0.dir1",
			want:  []string{"dir0", "dir1"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, objmap.AsPath(tc.input))
		})
	}
}

func TestBase(t *testing.T) {
	for _, tc := range []struct {
		title string
		input []string
		want  string
	}{
		{
			title: "none",
			input: []string{},
			want:  "",
		},
		{
			title: "single",
			input: []string{"dir"},
			want:  "dir",
		},
		{
			title: "double",
			input: []string{"dir0", "dir1"},
			want:  "dir1",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.want, objmap.Base(tc.input))
		})
	}
}
