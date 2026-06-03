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

package decoder

import (
	"bufio"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

// DecodeAll decodes runtime objects using the core Kubernetes scheme
// and the AKO scheme.
//
// If scheme registration or decoding fails, the given test fails immediately.
func DecodeAll(t *testing.T, from io.Reader) []runtime.Object {
	s := runtime.NewScheme()
	require.NoError(t, scheme.AddToScheme(s))
	require.NoError(t, akov2.AddToScheme(s))

	decoder := serializer.NewCodecFactory(s).UniversalDeserializer()
	reader := yaml.NewYAMLReader(bufio.NewReader(from))

	var result []runtime.Object
	for {
		buf, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		obj, _, err := decoder.Decode(buf, nil, nil)
		require.NoError(t, err)
		result = append(result, obj)
	}

	return result
}
