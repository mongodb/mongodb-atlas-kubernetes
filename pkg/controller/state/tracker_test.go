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

package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestComputeStateTracker(t *testing.T) {
	tests := []struct {
		name         string
		obj          metav1.Object
		dependencies []client.Object
	}{
		{
			name: "generation only",
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test-pod",
					Namespace:  "default",
					Generation: 1,
				},
			},
			dependencies: nil,
		},
		{
			name: "with secret dependency",
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test-pod",
					Namespace:  "default",
					Generation: 1,
				},
			},
			dependencies: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "my-secret",
						Namespace:       "default",
						UID:             types.UID("secret-uid-123"),
						ResourceVersion: "1",
					},
				},
			},
		},
		{
			name: "with configmap dependency",
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test-pod",
					Namespace:  "default",
					Generation: 1,
				},
			},
			dependencies: []client.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "my-configmap",
						Namespace:       "default",
						UID:             types.UID("cm-uid-456"),
						ResourceVersion: "2",
					},
				},
			},
		},
		{
			name: "with multiple dependencies",
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test-pod",
					Namespace:  "default",
					Generation: 5,
				},
			},
			dependencies: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "secret-1",
						Namespace:       "default",
						UID:             types.UID("uid-1"),
						ResourceVersion: "10",
					},
				},
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "secret-2",
						Namespace:       "other-ns",
						UID:             types.UID("uid-2"),
						ResourceVersion: "20",
					},
				},
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "config-1",
						Namespace:       "default",
						UID:             types.UID("uid-3"),
						ResourceVersion: "30",
					},
				},
			},
		},
		{
			name: "different generation produces different hash",
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "test-pod",
					Namespace:  "default",
					Generation: 2, // different from first test case
				},
			},
			dependencies: nil,
		},
	}

	// Store results to check uniqueness
	results := make(map[string]string)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hash := ComputeStateTracker(tc.obj, tc.dependencies...)

			// Hash should be non-empty
			assert.NotEmpty(t, hash, "hash should not be empty")

			// Hash should be deterministic (calling twice should produce same result)
			hash2 := ComputeStateTracker(tc.obj, tc.dependencies...)
			assert.Equal(t, hash, hash2, "hash should be deterministic")

			// Store for uniqueness check
			results[tc.name] = hash
		})
	}

	// Verify that different inputs produce different hashes
	t.Run("different inputs produce different hashes", func(t *testing.T) {
		// "generation only" (gen=1) vs "different generation produces different hash" (gen=2)
		assert.NotEqual(t,
			results["generation only"],
			results["different generation produces different hash"],
			"different generations should produce different hashes")

		// With dependencies vs without
		assert.NotEqual(t,
			results["generation only"],
			results["with secret dependency"],
			"adding a dependency should change the hash")
	})
}
