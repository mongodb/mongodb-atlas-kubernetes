// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"fmt"
	"hash/fnv"
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ComputeStateTracker(obj metav1.Object, dependencies ...client.Object) string {
	stateMap := make(map[string][]byte)
	stateMap["generation"] = []byte(fmt.Sprint(obj.GetGeneration()))

	for _, dep := range dependencies {
		stateData := []byte(string(dep.GetUID()) + "." + dep.GetResourceVersion())

		if secret, ok := dep.(*corev1.Secret); ok {
			stateMap[fmt.Sprintf("secret/%s/%s", secret.GetNamespace(), secret.GetName())] = stateData
		}

		if cm, ok := dep.(*corev1.ConfigMap); ok {
			stateMap[fmt.Sprintf("configmap/%s/%s", cm.GetNamespace(), cm.GetName())] = stateData
		}
	}

	keys := make([]string, 0, len(stateMap))
	for k := range stateMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	hasher := fnv.New64a()
	for _, k := range keys {
		hasher.Write([]byte(k))
		hasher.Write(stateMap[k])
	}

	rawHash := hasher.Sum64()
	currentStateHash := rand.SafeEncodeString(fmt.Sprint(rawHash))

	return currentStateHash
}
