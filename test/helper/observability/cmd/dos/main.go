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

package main

import (
	"fmt"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability/install"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability/observe"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability/snapshot"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, `available commands: "install", "teardown", "snapshot", "install-snapshot", "observe"`)
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "install":
		err = install.Install(os.Stdout)
	case "teardown":
		err = install.Teardown(os.Stdout)
	case "snapshot":
		err = snapshot.Snapshot(os.Stdout)
	case "install-snapshot":
		err = install.InstallSnapshot(os.Stdout)
	case "observe":
		err = observe.Observe(os.Args[2:]...)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
