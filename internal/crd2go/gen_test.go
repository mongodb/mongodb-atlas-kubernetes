package crd2go_test

import (
	"io"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	f, err := samples.Open("samples/group.crd.yaml")
	require.NoError(t, err)
	defer f.Close()

	sf, err := samples.Open("samples/v1/crd.go")
	require.NoError(t, err)
	sample, err := io.ReadAll(sf)
	require.NoError(t, err)
	want := string(sample)
	defer sf.Close()

	crd, err := crd2go.ParseCRD(f)
	require.NoError(t, err)
	stmt, err := crd2go.Generate(crd, crd2go.FirstVersion)
	require.NoError(t, err)
	require.NotNil(t, stmt)
	require.Equal(t, want, stmt.GoString())
}
