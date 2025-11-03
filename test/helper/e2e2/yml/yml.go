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

package yml

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/scale/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

var (
	// ErrNoCR indicates the parsed YAML is not valid CR
	ErrNoCR = errors.New("YAML definition is not a CR")
)

type autoCloser struct {
	io.ReadCloser
	closed bool
}

func (ac *autoCloser) Read(b []byte) (int, error) {
	if ac.closed {
		return 0, io.EOF
	}
	n, err := ac.ReadCloser.Read(b)
	if err == io.EOF {
		if err := ac.ReadCloser.Close(); err != nil {
			log.Printf("autoCloser failed to close %v: %v", ac.ReadCloser, err)
		}
	}
	return n, err
}

func MustOpen(fsys fs.FS, path string) io.Reader {
	f, err := fsys.Open(path)
	if err != nil {
		panic(fmt.Errorf("Fatal: could not open virtual file system path %q: %w", path, err))
	}
	return &autoCloser{ReadCloser: f}
}

func MustParseObjects(ymls io.Reader) []client.Object {
	objs, err := ParseObjects(ymls)
	if err != nil {
		panic(fmt.Errorf("Fatal: could not parse CRs: %w", err))
	}
	return objs
}

func ParseObjects(ymls io.Reader) ([]client.Object, error) {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(ymls)
	objs := []client.Object{}
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if len(strings.TrimSpace(buffer.String())) > 0 {
				obj, err := DecodeObject(buffer.Bytes())
				if errors.Is(err, ErrNoCR) {
					buffer.Reset()
					continue
				}
				if err != nil {
					return nil, err
				}
				objs = append(objs, obj)
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
		obj, err := DecodeObject(buffer.Bytes())
		if err != nil && !errors.Is(err, ErrNoCR) {
			return nil, err
		}
		objs = append(objs, obj)
	}
	return objs, nil
}

func DecodeObject(content []byte) (client.Object, error) {
	sch := runtime.NewScheme()
	utilruntime.Must(scheme.AddToScheme(sch))
	utilruntime.Must(apiextensions.AddToScheme(sch))
	utilruntime.Must(apiextensionsv1.AddToScheme(sch))
	utilruntime.Must(apiextensionsv1.RegisterConversions(sch))
	utilruntime.Must(apiextensionsv1beta1.AddToScheme(sch))
	utilruntime.Must(apiextensionsv1beta1.RegisterConversions(sch))
	utilruntime.Must(corev1.AddToScheme(sch))

	utilruntime.Must(akov2.AddToScheme(sch))
	utilruntime.Must(akov2next.AddToScheme(sch))

	decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode

	rtObj, _, err := decode(content, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	obj, ok := rtObj.(client.Object)
	if !ok {
		return nil, fmt.Errorf("decoded object is not a client.Object: %T", rtObj)
	}

	return obj, nil
}
