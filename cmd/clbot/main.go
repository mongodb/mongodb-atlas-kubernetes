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

// Command clbot is the changelog CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/clbot"
)

func main() {
	if err := clbot.Run(context.Background(), clbot.RuntimeInputs{
		Args:       os.Args[1:],
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		DataSource: clbot.NewGitDataSource("."),
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
