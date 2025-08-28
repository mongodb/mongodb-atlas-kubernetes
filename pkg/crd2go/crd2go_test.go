package crd2go_test

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/josvazg/crd2go/internal/gen"
	"github.com/josvazg/crd2go/internal/gotype"
	"github.com/josvazg/crd2go/pkg/config"
	"github.com/josvazg/crd2go/pkg/crd2go"
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
	req := gotype.Request{
		CodeWriterFn: BufferForCRD(buffers),
		TypeDict:     gotype.NewTypeDict(nil, preloadedTypes()...),
		CoreConfig: config.CoreConfig{
			Version:  gen.FirstVersion,
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
	req := gotype.Request{
		CodeWriterFn: BufferForCRD(buffers),
		TypeDict:     gotype.NewTypeDict(nil, preloadedTypes()...),
		CoreConfig: config.CoreConfig{
			Version:  gen.FirstVersion,
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

func TestLoadConfig(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   string
		want    *config.Config
		wantErr string
	}{
		{
			name:  "empty config",
			input: "{}",
			want:  &config.Config{},
		},
		{
			name: "defaults is empty lists and maps",
			input: `skipList: []
reserved: []
renames: {}
imports: []`,
			want: &config.Config{
				CoreConfig: config.CoreConfig{
					Reserved: []string{},
					SkipList: []string{},
					Renames:  map[string]string{},
					Imports:  []config.ImportedTypeConfig{},
				},
			},
		},
		{
			name: "just input and output",
			input: `input: ./pkg/crd2go/samples/crds.yaml
output: ./pkg/crd2go/samples/v1
skipList: []
reserved: []
renames: {}
imports: []`,
			want: &config.Config{
				Input:  "./pkg/crd2go/samples/crds.yaml",
				Output: "./pkg/crd2go/samples/v1",
				CoreConfig: config.CoreConfig{
					Reserved: []string{},
					SkipList: []string{},
					Renames:  map[string]string{},
					Imports:  []config.ImportedTypeConfig{},
				},
			},
		},
		{
			name:    "bad yaml",
			input:   "this is not a good YAML config",
			wantErr: "cannot unmarshal",
		},
		{
			name:    "no input fails",
			input:   "",
			wantErr: "fake error",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			inputReader := newFakeFailureReader()
			if tc.input != "" {
				inputReader = bytes.NewBufferString(tc.input)
			}
			cfg, err := crd2go.LoadConfig(inputReader)
			if tc.wantErr == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.want, cfg)
			} else {
				require.Nil(t, cfg)
				assert.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestCodeFileForCRDAtPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp(".", "test-code-file-for-crd-path")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	cwFn := crd2go.CodeWriterAtPath(tmpDir)
	require.NotNil(t, cwFn)

	w1, err := cwFn("testfile.go", false)
	assert.NoError(t, err)
	defer w1.Close()

	w2, err := cwFn("testfile.go", true)
	assert.NoError(t, err)
	defer w2.Close()

	_, err = cwFn("..", true)
	assert.ErrorContains(t, err, "failed to create file")
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

func BufferForCRD(buffers map[string]*bytes.Buffer) config.CodeWriterFunc {
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

func preloadedTypes() []*gotype.GoType {
	return append(gotype.KnownTypes(), reservedTypeNames(extraReserved)...)
}

func reservedTypeNames(reservedNames []string) []*gotype.GoType {
	reserved := make([]*gotype.GoType, 0, len(reservedNames))
	for _, reservedName := range reservedNames {
		reserved = append(reserved, ReserveTypeName(reservedName))
	}
	return reserved
}

func ReserveTypeName(name string) *gotype.GoType {
	return gotype.NewOpaqueType(name)
}

type fakeFailureReader struct{}

func (ffr *fakeFailureReader) Read(buf []byte) (int, error) {
	return 0, errors.New("fake error")
}

func newFakeFailureReader() io.Reader {
	return &fakeFailureReader{}
}
