// Copyright 2026 MongoDB Inc
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

package clbot

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

// RuntimeInputs contains the arguments and runtime dependencies supplied to Run.
type RuntimeInputs struct {
	Args       []string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	DataSource DataSource
}

// Run executes the changelog CLI with signature Run(ctx context.Context, inputs RuntimeInputs).
// Inputs supplies Args, Stdin, Stdout, Stderr, and DataSource for the CLI invocation.
func Run(ctx context.Context, inputs RuntimeInputs) error {
	flags := flag.NewFlagSet("clbot", flag.ContinueOnError)
	flags.SetOutput(inputs.Stderr)
	// Keep parse diagnostics on stderr; help is handled below and written to stdout.
	flags.Usage = func() {}

	if err := flags.Parse(inputs.Args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return writeUsage(inputs.Stdout)
		}
		return err
	}

	parsed := flags.Args()
	if len(parsed) == 0 {
		if _, err := fmt.Fprintln(inputs.Stdout, Version); err != nil {
			return err
		}
		return writeUsage(inputs.Stdout)
	}

	command := parsed[0]
	commandArgs := parsed[1:]
	switch command {
	case "version":
		if err := requireArgs(command, commandArgs, 0); err != nil {
			return err
		}
		version, err := NextVersion(ctx, inputs.DataSource)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(inputs.Stdout, version)
		return err
	case "released", "latest":
		if err := requireArgs(command, commandArgs, 0); err != nil {
			return err
		}
		release, err := inputs.DataSource.LatestRelease(ctx)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(inputs.Stdout, strings.TrimPrefix(release, "v"))
		return err
	case "list":
		if err := requireArgs(command, commandArgs, 0); err != nil {
			return err
		}
		changes, err := List(ctx, inputs.DataSource)
		var unsetErr *UnsetChangesError
		if err != nil && !errors.As(err, &unsetErr) {
			return err
		}
		for _, change := range changes {
			if _, err := fmt.Fprintf(inputs.Stdout, "%s %s %s\n", change.ID, change.Type, change.Title); err != nil {
				return err
			}
			if change.Description != "" {
				if _, err := fmt.Fprintln(inputs.Stdout, change.Description); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(inputs.Stdout); err != nil {
				return err
			}
		}
		return err
	case "validate", "verify":
		if len(commandArgs) != 1 {
			return fmt.Errorf("%s requires exactly one input (got %d)", command, len(commandArgs))
		}
		if commandArgs[0] == "-" {
			if inputs.Stdin == nil {
				return fmt.Errorf("validate stdin: input reader is nil")
			}
			input, err := io.ReadAll(inputs.Stdin)
			if err != nil {
				return fmt.Errorf("validate stdin: read input: %w", err)
			}
			title, description, _ := strings.Cut(string(input), "\n")
			title = strings.TrimSpace(title)
			description = strings.TrimSpace(description)
			return validateChange(Change{
				Type:        changeType(title, description),
				Title:       title,
				Description: description,
			})
		}
		return Validate(ctx, inputs.DataSource, commandArgs[0])
	default:
		return fmt.Errorf("unknown command %q; run clbot with no arguments to see usage", command)
	}
}

func requireArgs(command string, args []string, want int) error {
	if len(args) != want {
		return fmt.Errorf("%s requires no arguments (got %d)", command, len(args))
	}
	return nil
}

func writeUsage(stdout io.Writer) error {
	_, err := fmt.Fprintln(stdout, `Usage: clbot <command> [arguments]

Commands:
  version           Print the calculated next release version
  released, latest   Print the latest release
  list               List changes since the latest release
  validate <input|-> Validate a change input (- reads stdin)
  verify <input|->   Alias for validate`)
	return err
}
