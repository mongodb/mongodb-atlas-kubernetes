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
	"fmt"
	"hash/fnv"

	"k8s.io/apimachinery/pkg/util/rand"
)

// hashNames will return a hash corresponding to a name and optional arguments
func hashNames(name string, args ...string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(name))
	for _, arg := range args {
		hasher.Write([]byte(arg))
	}
	rawHash := hasher.Sum64()

	return rand.SafeEncodeString(fmt.Sprint(rawHash))
}

// prefixedName produces {prefix}-{hash} by using HashNames
func prefixedName(prefix string, name string, args ...string) string {
	return fmt.Sprintf("%s-%s", prefix, hashNames(name, args...))
}
