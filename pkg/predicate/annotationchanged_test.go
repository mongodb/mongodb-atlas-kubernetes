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

package predicate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/predicate"
)

func TestAnnotationChanged(t *testing.T) {
	key := "test-key"
	p := predicate.AnnotationChanged(key)

	tests := []struct {
		name string
		old  map[string]string
		new  map[string]string
		want bool
	}{
		{
			name: "No annotations on both objects",
			old:  nil,
			new:  nil,
			want: false,
		},
		{
			name: "Annotation added",
			old:  nil,
			new:  map[string]string{key: "value"},
			want: true,
		},
		{
			name: "Annotation removed",
			old:  map[string]string{key: "value"},
			new:  nil,
			want: true,
		},
		{
			name: "Annotation value changed",
			old:  map[string]string{key: "old-value"},
			new:  map[string]string{key: "new-value"},
			want: true,
		},
		{
			name: "Annotation unchanged",
			old:  map[string]string{key: "value"},
			new:  map[string]string{key: "value"},
			want: false,
		},
		{
			name: "Different annotation key changed",
			old:  map[string]string{"other-key": "value"},
			new:  map[string]string{"other-key": "new-value"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldObj := &unstructured.Unstructured{}
			oldObj.SetAnnotations(tt.old)

			newObj := &unstructured.Unstructured{}
			newObj.SetAnnotations(tt.new)

			e := event.UpdateEvent{
				ObjectOld: oldObj,
				ObjectNew: newObj,
			}

			result := p.Update(e)
			assert.Equal(t, tt.want, result)
		})
	}
}
