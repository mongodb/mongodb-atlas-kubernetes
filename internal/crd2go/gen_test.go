package crd2go_test

import (
	"embed"
	"io"
	"testing"

	"github.com/josvazg/crd2go/internal/crd2go"
	"github.com/stretchr/testify/require"
)

//go:embed samples/*
var samples embed.FS

func TestGenerateCode(t *testing.T) {
	for _, tc := range []struct {
		name string
		crd string
		src string
	}{
		{
			name: "group",
			crd:  "samples/group.crd.yaml",
			src: "samples/v1/group.go",
		},
		{
			name: "group",
			crd:  "samples/networkpermission.crd.yaml",
			src: "samples/v1/networkpermission.go",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, err := samples.Open(tc.crd)
			require.NoError(t, err)
			defer f.Close()
		
			sf, err := samples.Open(tc.src)
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
		})
	}
}
