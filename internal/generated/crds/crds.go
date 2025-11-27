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

package crds

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

//go:embed crds.yaml
var crdsYAML []byte

// EmbeddedCRD tried to load the given kind from a set of embedded CRDs
func EmbeddedCRD(kind string) (*apiextensionsv1.CustomResourceDefinition, error) {
	for {
		crd, err := ParseCRD(bufio.NewScanner(bytes.NewBuffer(crdsYAML)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CRDs YAML for %q: %w", kind, err)
		}
		if crd.Spec.Names.Kind == kind {
			return crd, nil
		}
	}
}
