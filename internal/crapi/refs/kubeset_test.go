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

package refs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	akoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

var (
	testNamespace = "my-namespace"
	mainObj       = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "main-cm", Namespace: testNamespace},
	}
	dep1 = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "dep1-cm", Namespace: testNamespace},
	}
	dep2 = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "dep2-cm", Namespace: testNamespace},
	}
)

func TestNewKubeset(t *testing.T) {
	ctx := newKubeset(testScheme(t), mainObj, []client.Object{dep1, dep2})

	require.NotNil(t, ctx, "context should not be nil")
	assert.Equal(t, mainObj, ctx.main, "main object should be set correctly")
	assert.Len(t, ctx.m, 2, "map should contain the two initial dependencies")
	assert.Len(t, ctx.added, 0, "added slice should be empty on initialization")

	assert.Contains(t, ctx.m, client.ObjectKeyFromObject(dep1))
	assert.Equal(t, dep1, ctx.m[client.ObjectKeyFromObject(dep1)])

	assert.Contains(t, ctx.m, client.ObjectKeyFromObject(dep2))
	assert.Equal(t, dep2, ctx.m[client.ObjectKeyFromObject(dep2)])
}

func TestKubesetFindAndHas(t *testing.T) {
	ctx := newKubeset(testScheme(t), mainObj, []client.Object{dep1})

	testCases := []struct {
		name              string
		searchName        string
		expectedFound     client.Object
		expectedHasResult bool
	}{
		{
			name:              "should find existing object",
			searchName:        "dep1-cm",
			expectedFound:     dep1,
			expectedHasResult: true,
		},
		{
			name:              "should not find non-existent object",
			searchName:        "non-existent-cm",
			expectedFound:     nil,
			expectedHasResult: false,
		},
		{
			name:              "should not find with empty name",
			searchName:        "",
			expectedFound:     nil,
			expectedHasResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			foundObj := ctx.find(tc.searchName)
			hasResult := ctx.has(tc.searchName)

			assert.Equal(t, tc.expectedFound, foundObj, "find() returned unexpected object")
			assert.Equal(t, tc.expectedHasResult, hasResult, "has() returned unexpected result")
		})
	}
}

func TestKubesetAdd(t *testing.T) {
	ctx := newKubeset(testScheme(t), mainObj, []client.Object{dep1})
	require.Len(t, ctx.m, 1, "pre-condition failed: map should have 1 item")
	require.Len(t, ctx.added, 0, "pre-condition failed: added slice should be empty")

	newDep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "new-dep-cm", Namespace: testNamespace},
	}

	require.False(t, ctx.has(newDep.GetName()), "pre-condition failed: new object should not exist yet")

	ctx.add(newDep)

	assert.Len(t, ctx.m, 2, "map length should be 2 after adding")
	assert.Len(t, ctx.added, 1, "added slice should have 1 item after adding")
	assert.True(t, ctx.has(newDep.GetName()), "has() should find the newly added object")
	assert.Equal(t, newDep, ctx.find(newDep.GetName()), "find() should return the newly added object")
	assert.Equal(t, newDep, ctx.added[0], "the newly added object should be in the 'added' slice")

	anotherDep := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "another-dep-cm", Namespace: testNamespace},
	}
	ctx.add(anotherDep)

	assert.Len(t, ctx.m, 3, "map length should be 3 after second add")
	assert.Len(t, ctx.added, 2, "added slice should have 2 items after second add")
	assert.True(t, ctx.has(anotherDep.GetName()), "has() should find the second added object")
	assert.Contains(t, ctx.added, newDep, "added slice should still contain the first added object")
	assert.Contains(t, ctx.added, anotherDep, "added slice should contain the second added object")
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(scheme))
	require.NoError(t, akoscheme.AddToScheme(scheme))
	return scheme
}
