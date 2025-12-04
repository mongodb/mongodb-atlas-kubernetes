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
	"errors"
	"fmt"
	"strings"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	// ErrNoCRD failure when the YAML is not a CRD
	ErrNoCRD = errors.New("not a CRD")
)

// ParseCRDs scans a YAML stream and returns the next CRD found.
// If more than one CRD is present in the stream, calling again
// on the same stream will return the next CRD found.
func ParseCRDs(scanner *bufio.Scanner) ([]*apiextensionsv1.CustomResourceDefinition, error) {
	var (
		buffer bytes.Buffer
		result []*apiextensionsv1.CustomResourceDefinition
	)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if len(strings.TrimSpace(buffer.String())) > 0 {
				crd, err := DecodeCRD(buffer.Bytes())
				if errors.Is(err, ErrNoCRD) {
					buffer.Reset()
					continue
				}
				if err != nil {
					return nil, err
				}
				result = append(result, crd)
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		buffer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	if buffer.Len() > 0 {
		crd, err := DecodeCRD(buffer.Bytes())
		if err != nil && !errors.Is(err, ErrNoCRD) {
			return nil, err
		}
		result = append(result, crd)
	}

	return result, nil
}

func DecodeCRD(content []byte) (*apiextensionsv1.CustomResourceDefinition, error) {
	sch := runtime.NewScheme()
	_ = scheme.AddToScheme(sch)
	_ = apiextensions.AddToScheme(sch)
	_ = apiextensionsv1.AddToScheme(sch)
	_ = apiextensionsv1.RegisterConversions(sch)
	_ = apiextensionsv1beta1.AddToScheme(sch)
	_ = apiextensionsv1beta1.RegisterConversions(sch)

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	obj, _, err := decode(content, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	kind := obj.GetObjectKind().GroupVersionKind().Kind
	if kind != "CustomResourceDefinition" {
		return nil, fmt.Errorf("unexpected kind %q: %w", kind, err)
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	err = sch.Convert(obj, crd, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to convert CRD object: %w", err)
	}

	return crd, nil
}
