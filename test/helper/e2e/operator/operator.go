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

package operator

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func NoGoRunEnvSet() bool {
	envSet, _ := strconv.ParseBool(os.Getenv("NO_GORUN"))
	return envSet
}

func RunDelveEnvSet() bool {
	envSet, _ := strconv.ParseBool(os.Getenv("RUN_DELVE"))
	return envSet
}

func operatorCommand() []string {
	if RunDelveEnvSet() {
		return []string{"dlv", "exec", "--api-version=2", "--headless=true", "--listen=:2345", filepath.Join(repositoryDir(), "bin", "manager"), "--"}
	}

	if NoGoRunEnvSet() {
		return []string{filepath.Join(repositoryDir(), "bin", "manager")}
	}

	return []string{"go", "run", filepath.Join(repositoryDir(), "cmd")}
}

type testingT interface {
	Logf(format string, a ...any)
	Fatalf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, a ...any)
}

type Operator struct {
	cmd     *exec.Cmd
	cmdLine []string
}

func NewOperator(namespace string, stdout, stderr io.Writer, cmdArgs ...string) *Operator {
	cmdLine := append(operatorCommand(), cmdArgs...)
	//nolint:gosec
	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)

	// works around  https://github.com/golang/go/issues/40467
	// to be able to propagate SIGTERM to the child process.
	// See https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf(`WATCH_NAMESPACE=%s`, namespace),
		fmt.Sprintf(`JOB_NAMESPACE=%s`, namespace),
		fmt.Sprintf(`OPERATOR_NAMESPACE=%s`, namespace),
		`OPERATOR_POD_NAME=mongodb-atlas-operator`,
	)

	return &Operator{
		cmd:     cmd,
		cmdLine: cmdLine,
	}
}

func (o *Operator) Start(t testingT) {
	t.Logf("starting operator command: %q", strings.Join(o.cmdLine, " "))
	if err := o.cmd.Start(); err != nil {
		t.Fatalf("failed to start operator: %v", err)
	}
}

func (o *Operator) Wait(t testingT) {
	t.Logf("waiting for operator to stop")
	if err := o.cmd.Wait(); err != nil {
		t.Errorf("error waiting for command: %v", err)
	}
}

func (o *Operator) Stop(t testingT) {
	// Ensure child process is killed on cleanup - send the negative of the pid, which is the process group id.
	// See https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	pid := 0

	if o.cmd != nil && o.cmd.Process != nil {
		pid = -o.cmd.Process.Pid
	}

	if pid != 0 {
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			t.Errorf("error trying to kill command: %v", err)
		}
	}

	if err := o.cmd.Wait(); err != nil {
		t.Errorf("error stopping operator: %v", err)
	}
}
