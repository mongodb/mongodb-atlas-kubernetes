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
