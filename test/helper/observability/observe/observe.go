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

package observe

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability/jsonwriter"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability/loki_reporter"
)

func Observe(cmdArgs ...string) error {
	var (
		lokiURL string
		jobName string
	)

	flag.StringVar(&jobName, "job-name", "ako", "The \"job\" label value to report to Loki")
	flag.StringVar(&lokiURL, "loki-url", "http://localhost:30002", "The URL of the Loki instance to report to")
	flag.Parse()

	target, err := loki_reporter.New(lokiURL, jobName, os.Stderr)
	if err != nil {
		return fmt.Errorf("error setting up loki: %w", err)
	}
	defer target.Stop()

	return forwardTo(target, cmdArgs...)
}

func forwardTo(target io.Writer, cmdArgs ...string) error {
	//nolint:gosec
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get cmdStdout pipe: %w", err)
	}

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get cmdStdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	for _, dest := range []struct {
		writer io.Writer
		reader io.Reader
	}{
		{
			writer: io.MultiWriter(os.Stdout, jsonwriter.NewJSONWriter(target, "INFO", "stdout")),
			reader: cmdStdout,
		},
		{
			writer: io.MultiWriter(os.Stderr, jsonwriter.NewJSONWriter(target, "ERROR", "stderr")),
			reader: cmdStderr,
		},
	} {
		go func() {
			defer wg.Done()
			scanner := bufio.NewScanner(dest.reader)
			for scanner.Scan() {
				_, err := dest.writer.Write(append(scanner.Bytes(), '\n'))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
			}
		}()
	}
	wg.Wait()

	return nil
}
