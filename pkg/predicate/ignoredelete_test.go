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

func TestIgnoreDeletedPredicate(t *testing.T) {
	p := predicate.IgnoreDeletedPredicate[*unstructured.Unstructured]()
	obj := &unstructured.Unstructured{}
	e := event.TypedDeleteEvent[*unstructured.Unstructured]{
		Object: obj,
	}

	result := p.Delete(e)
	assert.False(t, result, "DeleteFunc should always return false")
}
