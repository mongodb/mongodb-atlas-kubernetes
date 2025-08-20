package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/atlas2crd/cmd"
)

func main() {
	if err := cmd.RunCmd(context.Background()).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
