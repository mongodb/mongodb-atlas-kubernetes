package crd2go_test

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

//go:embed samples/*
var samples embed.FS

const (
	expectedSources = 13
)

var disabledKinds = []string{"Cluster", "CustomRole", "Group"}

var extraReserved = []string{"ConnectionStrings"}

func TestGenerateFromCRDStream(t *testing.T) {
	buffers := make(map[string]*bytes.Buffer)

	in, err := samples.Open("samples/crds.yaml")
	require.NoError(t, err)
	cfg := crd2go.GenerateConfig{
		Version:        crd2go.FirstVersion,
		Skips:          disabledKinds,
		PreloadedTypes: preloadedTypes(),
	}
	require.NoError(t, crd2go.GenerateStream(BufferForCRD(buffers), in, &cfg))

	assert.NotEmpty(t, buffers)
	assert.Len(t, buffers, expectedSources)
	for key, buf := range buffers {
		want := readTestFile(t, filepath.Join("samples", "v1", key))
		require.Equal(t, want, buf.String())
	}
}

func TestRefs(t *testing.T) {
	buffers := make(map[string]*bytes.Buffer)

	in, err := samples.Open("samples/samplerefs.yaml")
	require.NoError(t, err)
	cfg := crd2go.GenerateConfig{
		Version:        crd2go.FirstVersion,
		Skips:          disabledKinds,
		PreloadedTypes: preloadedTypes(),
	}
	require.NoError(t, crd2go.GenerateStream(BufferForCRD(buffers), in, &cfg))

	assert.NotEmpty(t, buffers)
	assert.Len(t, buffers, 1)
	for key, buf := range buffers {
		want := readTestFile(t, filepath.Join("samples", "refs", "v1", key))
		require.Equal(t, want, buf.String())
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	f, err := samples.Open(path)
	require.NoError(t, err)
	defer f.Close()

	b, err := io.ReadAll(f)
	require.NoError(t, err)

	return string(b)
}

func BufferForCRD(buffers map[string]*bytes.Buffer) crd2go.CodeWriterFunc {
	return func(crd *apiextensions.CustomResourceDefinition) (io.WriteCloser, error) {
		crdName := lowercase(crd.Spec.Names.Kind)
		key := fmt.Sprintf("%s.go", crdName)
		buffers[key] = bytes.NewBufferString("")
		return newWriteNopCloser(buffers[key]), nil
	}
}

// lowercase converts a string to lowercase using Go cases library
func lowercase(s string) string {
	return cases.Lower(language.English).String(s)
}

// WriteNopCloser wraps an io.Writer and adds a no-op Close method.
type writeNopCloser struct {
	io.Writer
}

// Close is a no-op to satisfy the io.WriteCloser interface.
func (w writeNopCloser) Close() error {
	return nil
}

// Helper function to create a WriteNopCloser
func newWriteNopCloser(w io.Writer) io.WriteCloser {
	return writeNopCloser{Writer: w}
}

func preloadedTypes() []*crd2go.GoType {
	return append(crd2go.KnownTypes(), reservedTypeNames(extraReserved)...)
}

func reservedTypeNames(reservedNames []string) []*crd2go.GoType {
	reserved := make([]*crd2go.GoType, 0, len(reservedNames))
	for _, reservedName := range reservedNames {
		reserved = append(reserved, ReserveTypeName(reservedName))
	}
	return reserved
}

func ReserveTypeName(name string) *crd2go.GoType {
	return crd2go.NewOpaqueType(name)
}
