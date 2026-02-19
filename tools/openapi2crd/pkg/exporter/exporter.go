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
//

package exporter

import (
	"os"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
)

const (
	createOrUpdate = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	defaultDirectoryPermission os.FileMode = 0o755
	defaultFilePermission      os.FileMode = 0o644
)

type Exporter interface {
	Open() error
	Export(crd *apiextensions.CustomResourceDefinition) error
	Close() error
}

func marshalCrd(crd *apiextensions.CustomResourceDefinition) ([]byte, error) {
	obj, err := convert(crd)
	if err != nil {
		return nil, err
	}

	obj.Kind = "CustomResourceDefinition"
	obj.APIVersion = "apiextensions.k8s.io/v1"

	yamlBytes, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return yamlBytes, nil
}

func convert(crd *apiextensions.CustomResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
	sch := runtime.NewScheme()
	_ = scheme.AddToScheme(sch)
	_ = apiextensionsv1.AddToScheme(sch)
	_ = apiextensionsv1.AddToScheme(sch)
	_ = apiextensionsv1.RegisterConversions(sch)

	out := &apiextensionsv1.CustomResourceDefinition{}
	err := sch.Convert(crd, out, nil)
	if err != nil {
		return nil, err
	}

	return out, err
}
