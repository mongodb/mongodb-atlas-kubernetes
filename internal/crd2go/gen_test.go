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

func TestGenerateCode(t *testing.T) {
	buffers := make(map[string]*bytes.Buffer)
	for _, tc := range []struct {
		name string
		crd  string
		src  string
	}{
		{
			name: "group",
			crd:  "samples/group.crd.yaml",
			src:  "samples/v1/group.go",
		},
		{
			name: "group",
			crd:  "samples/networkpermission.crd.yaml",
			src:  "samples/v1/networkpermissionentries.go",
		},
		{
			name: "groupalertconfigs",
			crd:  "samples/groupalertconfigs.crd.yaml",
			src:  "samples/v1/groupalertsconfig.go",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in, err := samples.Open(tc.crd)
			require.NoError(t, err)
			defer in.Close()

			want := readTestFile(t, tc.src)

			require.NoError(t, crd2go.GenerateStream(BufferForCRD(buffers), in, crd2go.FirstVersion))
			expectedKey := filepath.Base(tc.src)
			require.NotEmpty(t, buffers[expectedKey], "missing buffer for %s", expectedKey)
			require.Equal(t, want, buffers[expectedKey].String())
		})
	}
}

func TestGenerateFromCRDStream(t *testing.T) {
	buffers := make(map[string]*bytes.Buffer)

	in, err := samples.Open("samples/crds.yml")
	require.NoError(t, err)
	require.NoError(t, crd2go.GenerateStream(BufferForCRD(buffers), in, crd2go.FirstVersion))

	assert.NotEmpty(t, buffers)
	assert.Len(t, buffers, 8)
	for key, buf := range buffers {
		want := readTestFile(t, filepath.Join("samples", "v1", key))
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
