package crds

import (
	"bufio"
	"bytes"
	"fmt"

	_ "embed"

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
