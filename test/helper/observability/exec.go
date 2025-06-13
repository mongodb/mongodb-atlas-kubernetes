// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
