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
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/run"
)

const (
	DefaultDelveListenPort = ":2345"
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
	operatorBinary := envVarOrDefault("AKO_BINARY", filepath.Join("bin", "manager"))
	if RunDelveEnvSet() {
		return []string{
			"dlv", "exec",
			"--api-version=2",
			"--headless=true",
			fmt.Sprintf("--listen=%s", envVarOrDefault("DELVE_LISTEN", DefaultDelveListenPort)),
			filepath.Join(repositoryDir(), operatorBinary),
			"--",
		}
	}

	if NoGoRunEnvSet() {
		return []string{filepath.Join(repositoryDir(), operatorBinary)}
	}

	if os.Getenv("EXPERIMENTAL") == "true" {
		return []string{
			"go",
			"run",
			"-ldflags=-X github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version.Experimental=true",
			filepath.Join(repositoryDir(), "cmd"),
		}
	}
	return []string{"go", "run", filepath.Join(repositoryDir(), "cmd")}
}

type testingT interface {
	Logf(format string, a ...any)
	Fatalf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, a ...any)
}

type Operator interface {
	Start(ctx context.Context, t testingT)
	Running() bool
	Wait(t testingT)
	Stop(t testingT)
}

type OperatorProcess struct {
	mutex          sync.Mutex
	env            []string
	stdout, stderr io.Writer
	command        []string
	cmd            *exec.Cmd
	cancel         context.CancelFunc
}

func DefaultOperatorEnv(namespace string) []string {
	return append(
		os.Environ(),
		fmt.Sprintf(`WATCH_NAMESPACE=%s`, namespace),
		fmt.Sprintf(`JOB_NAMESPACE=%s`, namespace),
		fmt.Sprintf(`OPERATOR_NAMESPACE=%s`, namespace),
		`OPERATOR_POD_NAME=mongodb-atlas-operator`,
	)
}

func AllNamespacesOperatorEnv(operatorNamespace string) []string {
	return append(
		os.Environ(),
		fmt.Sprintf(`OPERATOR_NAMESPACE=%s`, operatorNamespace),
		`OPERATOR_POD_NAME=mongodb-atlas-operator`,
	)
}

func NewOperator(env []string, stdout, stderr io.Writer, cmdArgs ...string) Operator {
	if RunEmbeddedSet() {
		return NewEmbeddedOperator(run.Run, cmdArgs)
	}
	return &OperatorProcess{
		env:     env,
		stdout:  stdout,
		stderr:  stderr,
		command: append(operatorCommand(), cmdArgs...),
	}
}

func (o *OperatorProcess) Start(ctx context.Context, t testingT) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.cmd != nil {
		return
	}
	localCtx, cancel := context.WithCancel(ctx)
	// nolint:gosec // G204: cmdArgs are controlled by test authors, not external user input.
	cmd := exec.CommandContext(localCtx, o.command[0], o.command[1:]...)

	// works around  https://github.com/golang/go/issues/40467
	// to be able to propagate SIGTERM to the child process.
	// See https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = o.stdout
	cmd.Stderr = o.stderr
	cmd.Env = o.env

	o.cmd = cmd
	o.cancel = cancel

	t.Logf("starting operator command: %q", strings.Join(o.command, " "))
	if err := o.cmd.Start(); err != nil {
		t.Fatalf("failed to start operator: %v", err)
	}
}

func (o *OperatorProcess) Running() bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.cmd == nil {
		return false
	}
	return o.cmd.ProcessState == nil
}

func (o *OperatorProcess) Wait(t testingT) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.cmd == nil {
		return
	}
	t.Logf("waiting for operator to stop")
	if err := o.cmd.Wait(); err != nil {
		t.Errorf("error waiting for command: %v", err)
	}
}

func (o *OperatorProcess) Stop(t testingT) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.cmd == nil || o.cancel == nil {
		return
	}
	// Check if process is already terminated
	if o.cmd.ProcessState != nil && o.cmd.ProcessState.Exited() {
		return
	}

	o.cancel()

	// Send SIGTERM to the process group to ensure cleanup
	if o.cmd != nil && o.cmd.Process != nil {
		pid := -o.cmd.Process.Pid
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			// If process doesn't exist or we don't have permission, that's fine
			if err != syscall.ESRCH && err != syscall.EPERM {
				t.Errorf("error trying to kill command: %v", err)
			}
		}
	}

	// Wait for process to exit - ignore errors if we successfully terminated it
	_ = o.cmd.Wait()
}

func envVarOrDefault(name, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if ok {
		return value
	}
	return defaultValue
}
