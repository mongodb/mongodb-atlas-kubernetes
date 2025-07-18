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

package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/onsi/ginkgo/v2/dsl/core"
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

func RunManagerBinary(deletionProtection bool) (*exec.Cmd, error) {
	args := []string{
		"--log-level=-9",
		fmt.Sprintf("--object-deletion-protection=%v", deletionProtection),
		"--log-encoder=console",
		`--atlas-domain=https://cloud-qa.mongodb.com`,
	}

	cmdLine := append(operatorCommand(), args...)
	//nolint:gosec
	cmd := exec.Command(cmdLine[0], cmdLine[1:]...)

	// works around  https://github.com/golang/go/issues/40467
	// to be able to propagate SIGTERM to the child process.
	// See https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = core.GinkgoWriter
	cmd.Stderr = core.GinkgoWriter
	cmd.Env = append(
		os.Environ(),
		`OPERATOR_NAMESPACE=mongodb-atlas-system`,
		`OPERATOR_POD_NAME=mongodb-atlas-operator`,
	)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
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

func envVarOrDefault(name, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if ok {
		return value
	}
	return defaultValue
}

func repositoryDir() string {
	// Caller(0) returns the path to the calling test file rather than the path to this framework file. That
	// precludes assuming how many directories are between the file and the repo root. It's therefore necessary
	// to search in the hierarchy for an indication of a path that looks like the repo root.
	//nolint:dogsled
	_, sourceFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(sourceFile)
	for {
		// go.mod should always exist in the repo root
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		} else if errors.Is(err, os.ErrNotExist) {
			currentDir, err = filepath.Abs(filepath.Join(currentDir, ".."))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return currentDir
}
