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

package crapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/testdata/samples/v1"
)

func TestIsNil(t *testing.T) {
	for _, tc := range []struct {
		title string
		value any
		want  bool
	}{
		{title: "nil Group", value: (*v1.Group)(nil), want: true},
		{title: "nil interface", value: nil, want: true},
		{title: "nil pointer", value: (*int)(nil), want: true},
		{title: "nil slice", value: []int(nil), want: true},
		{title: "nil map", value: map[string]int(nil), want: true},
		{title: "non-nil interface", value: 1, want: false},
		{title: "non-nil pointer", value: pointer.MakePtr(1), want: false},
		{title: "non-nil slice", value: []int{1}, want: false},
		{title: "non-nil map", value: map[string]int{"1": 1}, want: false},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := isNil(tc.value)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsNilObject(t *testing.T) {
	for _, tc := range []struct {
		title string
		value client.Object
		want  bool
	}{
		{title: "nil struct", value: (*v1.Group)(nil), want: true},
		{title: "non-nil map", value: &v1.Group{}, want: false},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := isNil(tc.value)
			require.False(t, tc.value == nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
