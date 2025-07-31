package crd2go_test

import (
	"bytes"
	"embed"
	"io"
	"path/filepath"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed samples/*
var samples embed.FS

const (
	expectedSources = 19
)

var disabledKinds = []string{} // use ito skip problematic CRD kinds temporarily

var extraReserved = []string{} // use to fix problematic name picks, usually due to skips

func TestGenerateFromCRDs(t *testing.T) {
	buffers := make(map[string]*bytes.Buffer)

	in, err := samples.Open("samples/crds.yaml")
	require.NoError(t, err)
	req := crd2go.Request{
		CodeWriterFn: BufferForCRD(buffers),
		TypeDict: crd2go.NewTypeDict(nil, preloadedTypes()...),
		CoreConfig: crd2go.CoreConfig{
			Version:  crd2go.FirstVersion,
			SkipList: disabledKinds,
		},
	}
	require.NoError(t, crd2go.Generate(&req, in))

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
	req := crd2go.Request{
		CodeWriterFn: BufferForCRD(buffers),
		TypeDict: crd2go.NewTypeDict(nil, preloadedTypes()...),
		CoreConfig: crd2go.CoreConfig{
			Version:  crd2go.FirstVersion,
			SkipList: disabledKinds,
		},
	}
	_, err = crd2go.GenerateStream(&req, in)
	require.NoError(t, err)

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
	return func(filename string, overwrite bool) (io.WriteCloser, error) {
		buffers[filename] = bytes.NewBufferString("")
		return newWriteNopCloser(buffers[filename]), nil
	}
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
