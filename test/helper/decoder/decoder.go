package decoder

import (
	"bufio"
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
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		obj, _, err := decoder.Decode(buf, nil, nil)
		require.NoError(t, err)
		result = append(result, obj)
	}

	return result
}
