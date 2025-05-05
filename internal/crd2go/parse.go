package crd2go

import (
	"errors"
	"fmt"
	"io"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/scale/scheme"
)

var (
	// ErrNoCRD failure when the YAML is not a CRD
	ErrNoCRD = errors.New("not a CRD")
)

func ParseCRD(r io.Reader) (*apiextensionsv1.CustomResourceDefinition, error) {
	yml, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read all YAML input: %w", err)
	}
	if len(yml) == 0 {
		return nil, io.EOF
	}
	return DecodeCRD(yml)
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
	return crd, err
}
