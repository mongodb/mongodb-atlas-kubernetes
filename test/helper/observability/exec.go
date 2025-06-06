package observability

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func ExecCommand(logger io.Writer, cmdArgs ...string) error {
	fmt.Fprintln(logger, cmdArgs)
	//nolint:gosec
	out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Output()
	fmt.Fprintln(logger, string(out))

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		logger.Write(exitErr.Stderr)
	}

	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}
	return nil
}
