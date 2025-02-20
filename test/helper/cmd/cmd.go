package cmd

import (
	"bytes"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// RunCommand executes the given command with the given arguments
// and returns the resulting stdout and stderr as an io.Reader.
//
// If the command fails to run, the given test is being failed immediately.
func RunCommand(t *testing.T, name string, args ...string) io.Reader {
	var result bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &result
	cmd.Stderr = &result
	err := cmd.Run()
	if err != nil {
		t.Log(result.String())
	}
	require.NoError(t, err)
	return &result
}
