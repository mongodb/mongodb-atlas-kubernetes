package crd2go_test

import (
	"bytes"
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
			in, err := samples.Open(tc.crd)
			require.NoError(t, err)
			defer in.Close()
		
			sf, err := samples.Open(tc.src)
			require.NoError(t, err)
			sample, err := io.ReadAll(sf)
			require.NoError(t, err)
			want := string(sample)
			defer sf.Close()
		
			out := bytes.NewBufferString("")
			require.NoError(t, crd2go.GenerateStream(out, in, crd2go.FirstVersion))
			require.Equal(t, want, out.String())
		})
	}
}
